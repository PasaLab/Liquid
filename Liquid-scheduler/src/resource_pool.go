package main

import (
	"sync"
	"time"
	"net/url"
	"math/rand"
	"strconv"
	"sort"
	"hash/fnv"
)

var resourcePoolInstance *ResourcePool
var resourcePoolInstanceLock sync.Mutex

func InstanceOfResourcePool() *ResourcePool {
	defer resourcePoolInstanceLock.Unlock()
	resourcePoolInstanceLock.Lock()

	if resourcePoolInstance == nil {
		resourcePoolInstance = &ResourcePool{}
	}
	return resourcePoolInstance
}

type ResourcePool struct {
	poolsCount int
	pools      []PoolSeg
	poolsMu    sync.Mutex

	history   []PoolStatus
	historyMu sync.Mutex

	heartBeat    map[string]time.Time
	heartBeatMu  sync.Mutex
	versions     map[string]float64
	versionsMu   sync.Mutex
	counter      int
	counterTotal int

	subscriptions   map[string]map[string]int
	subscriptionsMu sync.Mutex

	networks     map[string]bool
	networksFree map[string]bool
	networkMu    sync.Mutex

	bindings      map[string]map[string]Job
	bindingsMu    sync.Mutex
	exclusiveJobs map[string]bool

	TotalGPU    int
	TotalCPU    int
	TotalMemory int
	TotalMu     sync.Mutex
	UsingGPU    int
	UsingMu     sync.Mutex

	enableBatch      bool
	batchJobs        map[string]Job
	batchMu          sync.Mutex
	batchAllocations map[string][]NodeStatus
	batchInterval    int
}

func (pool *ResourcePool) Start() {
	log.Info("RM started ")

	pool.networks = map[string]bool{}
	pool.networksFree = map[string]bool{}

	pool.bindings = map[string]map[string]Job{}
	pool.exclusiveJobs = map[string]bool{}

	pool.TotalGPU = 0
	pool.UsingGPU = 0

	pool.TotalCPU = 0
	pool.TotalMemory = 0

	pool.enableBatch = false
	pool.batchAllocations = map[string][]NodeStatus{}
	pool.batchJobs = map[string]Job{}
	pool.batchInterval = 15

	/* init pools */
	pool.poolsCount = 300
	for i := 0; i < pool.poolsCount; i++ {
		pool.pools = append(pool.pools, PoolSeg{Lock: sync.Mutex{}, ID: i})
	}
	/* generate working segs */
	for i := 0; i < 10; i++ {
		pool.pools[rand.Intn(pool.poolsCount)].Nodes = map[string]*NodeStatus{}
	}
	/* init Next pointer */
	var pre *PoolSeg
	for i := pool.poolsCount*2 - 1; ; i-- {
		if pool.pools[i%pool.poolsCount].Next != nil {
			break
		}
		pool.pools[i%pool.poolsCount].Next = pre
		if pool.pools[i%pool.poolsCount].Nodes != nil {
			pre = &pool.pools[i%pool.poolsCount]
		}
	}

	pool.versions = map[string]float64{}
	pool.subscriptions = map[string]map[string]int{}
	pool.heartBeat = map[string]time.Time{}
	go func() {
		pool.checkDeadNodes()
	}()

	pool.history = []PoolStatus{}
	go func() {
		pool.saveStatusHistory()
	}()

	go func() {
		/* batch allocation */
		for {
			time.Sleep(time.Second * time.Duration(pool.batchInterval))
			if !pool.enableBatch {
				continue
			}
			pool.batchMu.Lock()

			var nodes []NodeStatus
			var left []Job
			for {
				var tasks []Task
				for _, job := range pool.batchJobs {
					for _, task := range job.Tasks {
						tasks = append(tasks, task)
					}
				}
				//log.Info(tasks)
				job := Job{Tasks: tasks}
				if len(tasks) == 0 {
					break
				}
				nodes = pool.doAcquireResource(job)
				if len(nodes) == 0 {
					for jobName := range pool.batchJobs {
						left = append(left, pool.batchJobs[jobName])
						delete(pool.batchJobs, jobName)
						log.Info("cannot find a valid allocation, remove a job randomly: ", jobName)
						break
					}
					continue
				}
				for i, task := range job.Tasks {
					if _, ok := pool.batchAllocations[task.Job]; !ok {
						pool.batchAllocations[task.Job] = []NodeStatus{}
					}
					pool.batchAllocations[task.Job] = append(pool.batchAllocations[task.Job], nodes[i])
				}
				break
			}
			pool.batchJobs = map[string]Job{}
			for _, job := range left {
				pool.batchJobs[job.Name] = job
			}
			pool.batchMu.Unlock()
		}
	}()

	/* create overlay networks ahead */
	go func() {
		for {
			pool.networkMu.Lock()
			if len(pool.networksFree) < 5 {
				var network string
				for {
					network = "yao-net-" + strconv.Itoa(rand.Intn(999999))
					if _, ok := pool.networks[network]; !ok {
						break
					}
				}
				v := url.Values{}
				v.Set("name", network)
				spider := Spider{}
				spider.Method = "POST"
				spider.URL = "http://yao-agent-master:8000/create"
				spider.Data = v
				spider.ContentType = "application/x-www-form-urlencoded"
				err := spider.do()
				if err != nil {
					log.Warn(err.Error())
					continue
				}
				resp := spider.getResponse()
				resp.Body.Close()
				pool.networksFree[network] = true
				pool.networks[network] = true
			}
			pool.networkMu.Unlock()
			time.Sleep(time.Second * 1)
		}
	}()
}

/* check dead nodes periodically */
func (pool *ResourcePool) checkDeadNodes() {
	for {
		pool.heartBeatMu.Lock()
		var nodesToDel []string
		for k, v := range pool.heartBeat {
			if v.Add(time.Second * 30).Before(time.Now()) {
				segID := pool.getNodePool(k)
				seg := &pool.pools[segID]
				if seg.Nodes == nil {
					seg = seg.Next
				}

				seg.Lock.Lock()
				pool.TotalMu.Lock()
				if _, ok := seg.Nodes[k]; ok {
					pool.TotalGPU -= len(seg.Nodes[k].Status)
					pool.TotalCPU -= seg.Nodes[k].NumCPU
					pool.TotalMemory -= seg.Nodes[k].MemTotal
				}
				pool.TotalMu.Unlock()
				delete(seg.Nodes, k)
				seg.Lock.Unlock()
				pool.versionsMu.Lock()
				delete(pool.versions, k)
				pool.versionsMu.Unlock()
				nodesToDel = append(nodesToDel, k)
				log.Warn("node ", k, " is offline")
			}
		}
		for _, v := range nodesToDel {
			segID := pool.getNodePool(v)
			seg := &pool.pools[segID]
			if seg.Nodes == nil {
				seg = seg.Next
			}
			seg.Lock.Lock()
			delete(seg.Nodes, v)
			seg.Lock.Unlock()
			delete(pool.heartBeat, v)
		}
		pool.heartBeatMu.Unlock()
		time.Sleep(time.Second * 10)
	}
}

func (pool *ResourcePool) GPUModelToPower(model string) int {
	mapper := map[string]int{
		"K40": 2, "Tesla K40": 2,
		"K80": 3, "Tesla K80": 3,
		"P100": 4, "Tesla P100": 4,
	}
	if power, err := mapper[model]; !err {
		return power
	}
	return 1
}

func (pool *ResourcePool) getNodePool(name string) int {
	h := fnv.New32a()
	h.Write([]byte(name))
	return int(h.Sum32()) % pool.poolsCount
}

/* save pool status periodically */
func (pool *ResourcePool) saveStatusHistory() {
	/* waiting for nodes */
	time.Sleep(time.Second * 30)
	for {
		log.Debug("pool.saveStatusHistory")
		summary := PoolStatus{}

		UtilCPU := 0.0
		TotalCPU := 0
		TotalMem := 0
		AvailableMem := 0

		TotalGPU := 0
		UtilGPU := 0
		TotalMemGPU := 0
		AvailableMemGPU := 0
		nodesCount := 0

		start := pool.pools[0]
		if start.Nodes == nil {
			start = *start.Next
		}
		for cur := start; ; {
			cur.Lock.Lock()
			for _, node := range cur.Nodes {
				UtilCPU += node.UtilCPU
				TotalCPU += node.NumCPU
				TotalMem += node.MemTotal
				AvailableMem += node.MemAvailable

				for _, GPU := range node.Status {
					UtilGPU += GPU.UtilizationGPU
					TotalGPU ++
					TotalMemGPU += GPU.MemoryTotal
					AvailableMemGPU += GPU.MemoryFree
				}
			}
			nodesCount += len(cur.Nodes)
			cur.Lock.Unlock()
			cur = *cur.Next
			if cur.ID == start.ID {
				break
			}
		}
		summary.TimeStamp = time.Now().Format("2006-01-02 15:04:05")
		summary.UtilCPU = UtilCPU / (float64(nodesCount) + 0.001)
		summary.TotalCPU = TotalCPU
		summary.TotalMem = TotalMem
		summary.AvailableMem = AvailableMem
		summary.TotalGPU = TotalGPU
		if TotalGPU == 0 {
			summary.UtilGPU = 0.0
		} else {
			summary.UtilGPU = UtilGPU / TotalGPU
		}
		summary.TotalMemGPU = TotalMemGPU
		summary.AvailableMemGPU = AvailableMemGPU

		pool.historyMu.Lock()
		pool.history = append(pool.history, summary)
		if len(pool.history) > 60 {
			pool.history = pool.history[len(pool.history)-60:]
		}
		pool.historyMu.Unlock()

		pool.TotalMu.Lock()
		pool.TotalGPU = TotalGPU
		pool.TotalCPU = TotalCPU
		pool.TotalMemory = TotalMemGPU
		pool.TotalMu.Unlock()
		time.Sleep(time.Second * 60)
	}
}

/* update node info */
func (pool *ResourcePool) update(node NodeStatus) {
	pool.poolsMu.Lock()
	defer pool.poolsMu.Unlock()
	segID := pool.getNodePool(node.ClientID)
	seg := &pool.pools[segID]
	if seg.Nodes == nil {
		seg = seg.Next
	}
	seg.Lock.Lock()
	defer seg.Lock.Unlock()

	/* init bindings */
	go func(node NodeStatus) {
		pool.subscriptionsMu.Lock()
		defer pool.subscriptionsMu.Unlock()
		pool.bindingsMu.Lock()
		defer pool.bindingsMu.Unlock()
		for _, gpu := range node.Status {
			if _, ok := pool.subscriptions[gpu.UUID]; ok {
				for jobName := range pool.subscriptions[gpu.UUID] {
					go func(name string) {
						/* ask to update job status */
						scheduler.QueryState(name)
					}(jobName)
				}
			}
		}
		pool.heartBeatMu.Lock()
		pool.heartBeat[node.ClientID] = time.Now()
		pool.heartBeatMu.Unlock()
	}(node)

	pool.counterTotal++
	pool.versionsMu.Lock()
	if version, ok := pool.versions[node.ClientID]; ok && version == node.Version {
		//pool.versionsMu.Unlock()
		//return
	}
	pool.versionsMu.Unlock()
	pool.counter++
	log.Debug(node.Version, "!=", pool.versions[node.ClientID])

	status, ok := seg.Nodes[node.ClientID]
	if ok {
		/* keep allocation info */
		for i, GPU := range status.Status {
			if GPU.UUID == node.Status[i].UUID {
				node.Status[i].MemoryAllocated = GPU.MemoryAllocated
			}
		}
	} else {
		/* TODO: double check node do belong to this seg */
		pool.TotalMu.Lock()
		pool.TotalGPU += len(node.Status)
		pool.TotalCPU += node.NumCPU
		pool.TotalMemory += node.MemTotal
		pool.TotalMu.Unlock()
		log.Info("node ", node.ClientID, " is online")
	}
	seg.Nodes[node.ClientID] = &node
	if len(seg.Nodes) > 10 {
		go func() {
			pool.scaleSeg(seg)
		}()
	}
	pool.versions[node.ClientID] = node.Version
}

/* spilt seg */
func (pool *ResourcePool) scaleSeg(seg *PoolSeg) {
	log.Info("Scaling seg ", seg.ID)

	pool.poolsMu.Lock()
	defer pool.poolsMu.Unlock()

	var segIDs []int
	segIDs = append(segIDs, seg.ID)

	/* find previous seg */
	var pre *PoolSeg
	for i := seg.ID + pool.poolsCount - 1; i >= 0; i-- {
		segIDs = append(segIDs, i%pool.poolsCount)
		if pool.pools[i%pool.poolsCount].Next.ID != seg.ID {
			break
		}
		pre = &pool.pools[i%pool.poolsCount]
	}

	distance := seg.ID - pre.ID
	if distance < 0 {
		distance += pool.poolsCount
	}
	if distance <= 1 {
		log.Warn("Unable to scale, ", seg.ID, ", already full")
		return
	}

	candidate := pre
	/* walk to the nearest middle */
	if pre.ID < seg.ID {
		candidate = &pool.pools[(pre.ID+seg.ID)/2]
	} else {
		candidate = &pool.pools[(pre.ID+seg.ID+pool.poolsCount)/2%pool.poolsCount]
	}
	candidate.Next = seg
	candidate.Nodes = map[string]*NodeStatus{}

	/* lock in asc sequence to avoid deadlock */
	sort.Ints(segIDs)
	for _, id := range segIDs {
		pool.pools[id].Lock.Lock()
	}
	//log.Println(segIDs)

	/* update Next */
	for i := 0; ; i++ {
		id := (pre.ID + i) % pool.poolsCount
		if id == candidate.ID {
			break
		}
		pool.pools[id].Next = candidate
	}

	/* move nodes */
	nodesToMove := map[string]*NodeStatus{}
	for _, node := range seg.Nodes {
		seg2ID := pool.getNodePool(node.ClientID)
		seg2 := &pool.pools[seg2ID]
		if seg2.Nodes == nil {
			seg2 = seg2.Next
		}
		if seg2.ID != seg.ID {
			nodesToMove[node.ClientID] = node
		}
	}
	for _, node := range nodesToMove {
		delete(seg.Nodes, node.ClientID)
	}
	candidate.Nodes = nodesToMove
	//log.Info("pre=", pre.ID, " active=", candidate.ID, " seg=", seg.ID)
	for _, id := range segIDs {
		pool.pools[id].Lock.Unlock()
	}
}

/* get node by ClientID */
func (pool *ResourcePool) getByID(id string) NodeStatus {
	poolID := pool.getNodePool(id)
	seg := &pool.pools[poolID]
	if seg.Nodes == nil {
		seg = seg.Next
	}
	seg.Lock.Lock()
	defer seg.Lock.Unlock()

	status, ok := seg.Nodes[id]
	if ok {
		return *status
	}
	return NodeStatus{}
}

/* get all nodes */
func (pool *ResourcePool) list() MsgResource {
	nodes := map[string]NodeStatus{}

	start := pool.pools[0]
	if start.Nodes == nil {
		start = *start.Next
	}
	for cur := start; ; {
		cur.Lock.Lock()
		for k, node := range cur.Nodes {
			nodes[k] = *node
		}
		cur.Lock.Unlock()
		cur = *cur.Next
		if cur.ID == start.ID {
			break
		}
	}
	return MsgResource{Code: 0, Resource: nodes}
}

func (pool *ResourcePool) statusHistory() MsgPoolStatusHistory {
	pool.historyMu.Lock()
	defer pool.historyMu.Unlock()
	history := pool.history
	return MsgPoolStatusHistory{Code: 0, Data: history}
}

func (pool *ResourcePool) getCounter() map[string]int {
	return map[string]int{"counter": pool.counter, "counterTotal": pool.counterTotal}
}

func (pool *ResourcePool) acquireNetwork() string {
	pool.networkMu.Lock()
	defer pool.networkMu.Unlock()
	var network string
	log.Debug(pool.networksFree)
	if len(pool.networksFree) == 0 {
		for {
			for {
				network = "yao-net-" + strconv.Itoa(rand.Intn(999999))
				if _, ok := pool.networks[network]; !ok {
					break
				}
			}
			v := url.Values{}
			v.Set("name", network)
			spider := Spider{}
			spider.Method = "POST"
			spider.URL = "http://yao-agent-master:8000/create"
			spider.Data = v
			spider.ContentType = "application/x-www-form-urlencoded"
			err := spider.do()
			if err != nil {
				log.Warn(err.Error())
				continue
			}
			resp := spider.getResponse()
			resp.Body.Close()
			pool.networksFree[network] = true
			pool.networks[network] = true
			break
		}
	}

	for k := range pool.networksFree {
		network = k
		delete(pool.networksFree, k)
		break
	}
	return network
}

func (pool *ResourcePool) releaseNetwork(network string) {
	pool.networkMu.Lock()
	//pool.networksFree[network] = true
	pool.networkMu.Unlock()
}

func (pool *ResourcePool) attach(GPU string, job Job) {
	pool.subscriptionsMu.Lock()
	defer pool.subscriptionsMu.Unlock()
	pool.bindingsMu.Lock()
	defer pool.bindingsMu.Unlock()

	if _, ok := pool.subscriptions[GPU]; !ok {
		pool.subscriptions[GPU] = map[string]int{}
	}
	pool.subscriptions[GPU][job.Name] = int(time.Now().Unix())

	if _, ok := pool.bindings[GPU]; !ok {
		pool.bindings[GPU] = map[string]Job{}
	}
	pool.bindings[GPU][job.Name] = job
}

func (pool *ResourcePool) detach(GPU string, job Job) {
	pool.subscriptionsMu.Lock()
	defer pool.subscriptionsMu.Unlock()
	pool.bindingsMu.Lock()
	defer pool.bindingsMu.Unlock()

	if _, ok := pool.subscriptions[GPU]; ok {
		delete(pool.subscriptions[GPU], job.Name)
	}

	if list, ok := pool.bindings[GPU]; ok {
		delete(list, job.Name)
	}
}

/* return free & using GPUs */
func (pool *ResourcePool) countGPU() (int, int) {
	return pool.TotalGPU - pool.UsingGPU, pool.UsingGPU
}

func (pool *ResourcePool) acquireResource(job Job) []NodeStatus {
	for i := range job.Tasks {
		job.Tasks[i].Job = job.Name
	}
	if !pool.enableBatch {
		return pool.doAcquireResource(job)
	}
	pool.batchMu.Lock()
	pool.batchJobs[job.Name] = job
	pool.batchMu.Unlock()
	for {
		/* wait until request is satisfied */
		pool.batchMu.Lock()
		if _, ok := pool.batchAllocations[job.Name]; ok {
			pool.batchMu.Unlock()
			break
		} else {
			pool.batchMu.Unlock()
			time.Sleep(time.Millisecond * 100)
		}
	}
	pool.batchMu.Lock()
	nodes := pool.batchAllocations[job.Name]
	delete(pool.batchAllocations, job.Name)
	pool.batchMu.Unlock()
	return nodes
}

func (pool *ResourcePool) doAcquireResource(job Job) []NodeStatus {
	if len(job.Tasks) == 0 {
		return []NodeStatus{}
	}
	segID := rand.Intn(pool.poolsCount)
	if pool.TotalGPU < 100 {
		segID = 0
	}
	start := &pool.pools[segID]
	if start.Nodes == nil {
		start = start.Next
	}

	config := InstanceOfConfiguration()

	locks := map[int]*sync.Mutex{}

	/* 1-Share, 2-Vacant, 3-PreSchedule */
	allocationType := 0

	var candidates []NodeStatus

	if pool.TotalGPU == 0 {
		return []NodeStatus{}
	}
	var ress []NodeStatus

	loadRatio := float64(pool.UsingGPU) / float64(pool.TotalGPU)
	/* first, choose sharable GPUs */
	task := job.Tasks[0]
	if len(job.Tasks) == 1 && task.NumberGPU == 1 && loadRatio >= config.EnableShareRatio {
		// check sharable
		allocationType = 1
		pred := InstanceOfOptimizer().PredictReq(job, "Worker")
		availables := map[string][]GPUStatus{}
		for cur := start; ; {
			if _, ok := locks[cur.ID]; !ok {
				cur.Lock.Lock()
				locks[cur.ID] = &cur.Lock
			}

			for _, node := range cur.Nodes {
				var available []GPUStatus
				for _, status := range node.Status {
					if status.MemoryAllocated > 0 && status.MemoryTotal > task.MemoryGPU+status.MemoryAllocated {

						pool.bindingsMu.Lock()
						if jobs, ok := pool.bindings[status.UUID]; ok {
							totalUtil := pred.UtilGPU
							for _, job := range jobs {
								utilT := InstanceOfOptimizer().PredictReq(job, "Worker").UtilGPU
								totalUtil += utilT
							}
							if totalUtil < int(InstanceOfConfiguration().ShareMaxUtilization*100) {
								available = append(available, status)
							}
						}
						pool.bindingsMu.Unlock()
					}
				}
				if len(available) >= task.NumberGPU {
					candidates = append(candidates, *node)
					availables[node.ClientHost] = available
					if len(candidates) >= len(job.Tasks)*3+5 {
						break
					}
				}
			}
			if len(candidates) >= len(job.Tasks)*3+5 {
				break
			}
			if cur.ID > cur.Next.ID {
				break
			}
			cur = cur.Next
		}

		if len(candidates) > 0 {
			node := candidates[0]
			res := NodeStatus{}
			res.ClientID = node.ClientID
			res.ClientHost = node.ClientHost
			res.NumCPU = task.NumberCPU
			res.MemTotal = task.Memory
			res.Status = availables[node.ClientHost][0:task.NumberGPU]

			for i := range res.Status {
				pool.bindingsMu.Lock()
				if jobsT, okT := pool.bindings[res.Status[i].UUID]; okT {
					for jobT := range jobsT {
						delete(pool.exclusiveJobs, jobT)
					}
				}
				pool.bindingsMu.Unlock()

				for j := range node.Status {
					if res.Status[i].UUID == node.Status[j].UUID {
						if node.Status[j].MemoryAllocated == 0 {
							pool.UsingMu.Lock()
							pool.UsingGPU ++
							pool.UsingMu.Unlock()
						}
						node.Status[j].MemoryAllocated += task.MemoryGPU
						res.Status[i].MemoryTotal = task.MemoryGPU
					}
				}
			}
			for _, t := range res.Status {
				pool.attach(t.UUID, job)
			}
			ress = append(ress, res)
		}
	}
	//log.Info(candidates)

	/* second round, find vacant gpu */
	if len(candidates) == 0 {
		allocationType = 2
		for cur := start; ; {
			if _, ok := locks[cur.ID]; !ok {
				cur.Lock.Lock()
				locks[cur.ID] = &cur.Lock
			}
			for _, node := range cur.Nodes {
				var available []GPUStatus
				for _, status := range node.Status {
					/* make sure GPU is not used by in-system and outer-system */
					if status.MemoryAllocated == 0 { //} && status.MemoryUsed < 100 {
						available = append(available, status)
					}
				}
				if len(available) >= task.NumberGPU {
					candidates = append(candidates, *node)
					if len(candidates) >= len(job.Tasks)*3+5 {
						break
					}
				}
			}
			if len(candidates) >= len(job.Tasks)*3+5 {
				break
			}
			if cur.ID > cur.Next.ID {
				break
			}
			cur = cur.Next
		}
		//log.Info(candidates)
	}

	/* third round, find gpu to be released */
	if len(candidates) == 0 && len(job.Tasks) == 1 && task.NumberGPU == 1 {
		estimate := InstanceOfOptimizer().PredictTime(job)
		log.Debug(estimate)

		if loadRatio >= config.EnablePreScheduleRatio {
			allocationType = 3
			availables := map[string][]GPUStatus{}
			for cur := start; ; {
				if _, ok := locks[cur.ID]; !ok {
					cur.Lock.Lock()
					locks[cur.ID] = &cur.Lock
				}
				for _, node := range cur.Nodes {
					var available []GPUStatus
					for _, status := range node.Status {
						if jobs, ok := pool.bindings[status.UUID]; ok {
							if len(jobs) > 1 || status.MemoryAllocated == 0 {
								continue
							}
							for _, jobT := range jobs {
								est := InstanceOfOptimizer().PredictTime(jobT)
								now := time.Now().Unix()
								if int(now-jobT.StartedAt) > est.Total-est.Post-estimate.Pre-InstanceOfConfiguration().PreScheduleExtraTime {
									available = append(available, status)
								}
							}
						}
					}
					if len(available) >= task.NumberGPU {
						candidates = append(candidates, *node)
						availables[node.ClientHost] = available
						if len(candidates) >= len(job.Tasks)*3+5 {
							break
						}
					}
				}
				if len(candidates) >= len(job.Tasks)*3+5 {
					break
				}
				if cur.ID > cur.Next.ID {
					break
				}
				cur = cur.Next
			}
			//log.Info(candidates)
			if len(candidates) > 0 {
				node := candidates[0]
				res := NodeStatus{}
				res.ClientID = node.ClientID
				res.ClientHost = node.ClientHost
				res.NumCPU = task.NumberCPU
				res.MemTotal = task.Memory
				res.Status = availables[node.ClientHost][0:task.NumberGPU]

				for i := range res.Status {
					for j := range node.Status {
						if res.Status[i].UUID == node.Status[j].UUID {
							if node.Status[j].MemoryAllocated == 0 {
								pool.UsingMu.Lock()
								pool.UsingGPU ++
								pool.UsingMu.Unlock()
							}
							node.Status[j].MemoryAllocated += task.MemoryGPU
							res.Status[i].MemoryTotal = task.MemoryGPU
							/* being fully used, means ahead */
							res.Status[i].MemoryUsed = res.Status[i].MemoryTotal
						}
					}
				}
				for _, t := range res.Status {
					pool.attach(t.UUID, job)
				}
				ress = append(ress, res)
			}
		}
	}

	if len(candidates) > 0 {
		log.Info("allocationType is ", allocationType)
		//log.Info(candidates)
	}

	/* assign */
	if len(candidates) > 0 && len(ress) == 0 {
		var nodesT []NodeStatus
		for _, node := range candidates {
			nodesT = append(nodesT, node.Copy())
		}

		tasks := make([]Task, len(job.Tasks))
		var tasksPS []Task
		var tasksWorker []Task
		for _, taskT := range job.Tasks {
			if taskT.IsPS {
				tasksPS = append(tasksPS, taskT)
			} else {
				tasksWorker = append(tasksWorker, taskT)
			}
		}
		idxPS := 0
		idxWorker := 0
		factor := float64(len(tasksWorker)) / (float64(len(tasksPS)) + 0.001)
		for i := range tasks {
			if float64(idxPS)*factor <= float64(idxWorker) && idxPS < len(tasksPS) {
				tasks[i] = tasksPS[idxPS]
				idxPS++
			} else if idxWorker < len(tasksWorker) {
				tasks[i] = tasksWorker[idxWorker]
				idxWorker++
			} else {
				tasks[i] = tasksPS[idxPS]
				idxPS++
			}
		}

		//log.Info(tasks, factor)
		allocation := InstanceOfAllocator().allocate(nodesT, tasks)
		//log.Info(allocation)
		if allocation.Flags["valid"] {
			for range job.Tasks { //append would cause uncertain order
				ress = append(ress, NodeStatus{ClientID: "null"})
			}

			cnt := 0
			for nodeID, tasks := range allocation.TasksOnNode {
				var node *NodeStatus
				for i := range candidates {
					if candidates[i].ClientID == nodeID {
						node = &candidates[i]
					}
				}

				var available []GPUStatus
				for _, gpu := range node.Status {
					if gpu.MemoryAllocated == 0 {
						available = append(available, gpu)
					}
				}
				for _, task := range tasks {
					cnt++
					res := NodeStatus{}
					res.ClientID = node.ClientID
					res.ClientHost = node.ClientHost
					res.NumCPU = task.NumberCPU
					res.MemTotal = task.Memory
					res.Status = available[0:task.NumberGPU]
					available = available[task.NumberGPU:]

					for i := range res.Status {
						for j := range node.Status {
							if res.Status[i].UUID == node.Status[j].UUID {
								if node.Status[j].MemoryAllocated == 0 {
									pool.UsingMu.Lock()
									pool.UsingGPU ++
									pool.UsingMu.Unlock()
								}
								node.Status[j].MemoryAllocated += task.MemoryGPU
								res.Status[i].MemoryTotal = task.MemoryGPU
							}
						}
					}
					for _, t := range res.Status {
						pool.attach(t.UUID, job)
					}

					flag := false
					for i := range job.Tasks {
						if job.Tasks[i].ID == task.ID {
							ress[i] = res
							flag = true
							break
						}
					}
					if !flag {
						log.Warn("Unable to find task, ", res)
					}

				}
			}

			if cnt != len(job.Tasks) {
				log.Warn("Allocation is invalid")
				log.Warn(cnt, len(job.Tasks))
				log.Warn(job.Tasks)
				log.Warn(allocation.TasksOnNode)
			}

		}
	}

	pool.bindingsMu.Lock()
	if allocationType == 2 {
		pool.exclusiveJobs[job.Name] = true
	}
	pool.bindingsMu.Unlock()

	for segID, lock := range locks {
		log.Debug("Unlock ", segID)
		lock.Unlock()
	}
	return ress
}

/*
TODO:
bug-1: node is offline, unable to retrieve allocation info
bug-2: when node offline & back, allocation info is lost
*/
func (pool *ResourcePool) releaseResource(job Job, agent NodeStatus) {
	segID := pool.getNodePool(agent.ClientID)
	seg := pool.pools[segID]
	if seg.Nodes == nil {
		seg = *seg.Next
	}
	seg.Lock.Lock()
	defer seg.Lock.Unlock()

	node, ok := seg.Nodes[agent.ClientID]
	/* in case node is offline */
	if !ok {
		/* bug-1 */
		log.Warn("node ", agent.ClientID, " not present")
		return
	}
	for _, gpu := range agent.Status {
		for j := range node.Status {
			if gpu.UUID == node.Status[j].UUID {
				node.Status[j].MemoryAllocated -= gpu.MemoryTotal
				log.Debug(node.Status[j].MemoryAllocated)
				if node.Status[j].MemoryAllocated < 0 {
					/* bug-2: a node is offline and then online, the allocation info will be lost */
					log.Warn(node.ClientID, " UUID=", gpu.UUID, " More Memory Allocated")
					node.Status[j].MemoryAllocated = 0
				}
				if node.Status[j].MemoryAllocated == 0 {
					pool.UsingMu.Lock()
					pool.UsingGPU--
					pool.UsingMu.Unlock()
					log.Info(node.Status[j].UUID, " is released")
				}
			}
		}
	}
}

func (pool *ResourcePool) SetBatchEnabled(enabled bool) bool {
	pool.enableBatch = enabled
	log.Info("enableBatch is set to ", enabled)
	return true
}

func (pool *ResourcePool) SetBatchInterval(interval int) bool {
	if interval < 1 {
		interval = 1
	}
	pool.batchInterval = interval
	log.Info("batchInterval is updated to ", interval)
	return true
}

func (pool *ResourcePool) isExclusive(jobName string) bool {
	pool.bindingsMu.Lock()
	defer pool.bindingsMu.Unlock()
	_, ok := pool.exclusiveJobs[jobName]
	/* clear after called */
	delete(pool.exclusiveJobs, jobName)
	return ok
}

func (pool *ResourcePool) Dump() map[string]interface{} {
	res := map[string]interface{}{}
	res["batchJobs"] = pool.batchJobs
	res["bindings"] = pool.bindings
	return res
}
