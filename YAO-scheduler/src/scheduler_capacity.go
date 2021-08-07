package main

import (
	"sync"
	"time"
	"sort"
	"math"
)

type SchedulerCapacity struct {
	history   []*Job
	historyMu sync.Mutex

	nextQueue string
	jobMasters      map[string]*JobManager
	queues    map[string][]Job
	queuesMu  sync.Mutex

	schedulingJobs map[string]bool
	schedulingMu   sync.Mutex

	resourceAllocations   map[string]*ResourceCount
	resourceAllocationsMu sync.Mutex

	enabled     bool
	parallelism int

	allocatingGPU   int
	allocatingGPUMu sync.Mutex
}

func (scheduler *SchedulerCapacity) Start() {
	log.Info("JS (capacity) started")

	scheduler.jobMasters = map[string]*JobManager{}
	scheduler.history = []*Job{}
	scheduler.nextQueue = "default"
	scheduler.queues = map[string][]Job{}
	scheduler.queues["default"] = []Job{}
	scheduler.resourceAllocations = map[string]*ResourceCount{}
	scheduler.enabled = true
	scheduler.schedulingJobs = map[string]bool{}
	scheduler.allocatingGPU = 0

	scheduler.parallelism = 1

	go func() {
		flag := true
		for {
			log.Debug("Scheduling")
			if !flag {
				time.Sleep(time.Millisecond * 100)
			}
			flag = false
			if !scheduler.enabled {
				time.Sleep(time.Millisecond * 100)
				continue
			}
			scheduler.schedulingMu.Lock()
			if len(scheduler.schedulingJobs) >= scheduler.parallelism {
				scheduler.schedulingMu.Unlock()
				time.Sleep(time.Millisecond * 100)
				continue
			}
			scheduler.schedulingMu.Unlock()

			scheduler.queuesMu.Lock()
			queue := scheduler.nextQueue
			go func() {
				scheduler.UpdateNextQueue()
			}()
			if len(scheduler.queues[queue]) > 0 {
				jm := JobManager{}
				jm.job = scheduler.queues[queue][0]

				cnt := 0
				for _, task := range jm.job.Tasks {
					cnt += task.NumberGPU
				}

				pool := InstanceOfResourcePool()
				log.Info(cnt, pool.TotalGPU, pool.UsingGPU, scheduler.allocatingGPU)
				if len(scheduler.schedulingJobs) > 1 && (cnt*10+(scheduler.allocatingGPU)*13 > (pool.TotalGPU-pool.UsingGPU)*10) {
					scheduler.queuesMu.Unlock()
					continue
				}

				flag = true
				scheduler.allocatingGPUMu.Lock()
				scheduler.allocatingGPU += cnt
				scheduler.allocatingGPUMu.Unlock()
				log.Info("allocatingGPU is ", scheduler.allocatingGPU)
				log.Info("schedulingJobs are ", scheduler.schedulingJobs)

				scheduler.queues[queue] = scheduler.queues[queue][1:]
				jm.scheduler = scheduler
				scheduler.jobMasters[jm.job.Name] = &jm

				jm.job.Status = Starting
				scheduler.historyMu.Lock()
				scheduler.history = append(scheduler.history, &jm.job)
				scheduler.historyMu.Unlock()

				scheduler.schedulingMu.Lock()
				scheduler.schedulingJobs[jm.job.Name] = true
				scheduler.schedulingMu.Unlock()
				go func() {
					jm.start()
				}()
			} else {
				log.Debug("No more jobs to scheduling ", time.Now())
			}
			scheduler.queuesMu.Unlock()
		}
	}()
}

func (scheduler *SchedulerCapacity) UpdateProgress(job Job, state State) {
	scheduler.historyMu.Lock()
	defer scheduler.historyMu.Unlock()

	scheduler.schedulingMu.Lock()
	delete(scheduler.schedulingJobs, job.Name)
	scheduler.schedulingMu.Unlock()

	switch state {
	case Running:
		for i := range scheduler.history {
			if scheduler.history[i].Name == job.Name {
				scheduler.history[i].Status = Running
				scheduler.history[i].UpdatedAt = int(time.Now().Unix())
			}
		}
		break
	case Finished:
		for i := range scheduler.history {
			if scheduler.history[i].Name == job.Name {
				scheduler.history[i].Status = Finished
				scheduler.history[i].UpdatedAt = int(time.Now().Unix())
			}
		}
		break
	case Stopped:
		for i := range scheduler.history {
			if scheduler.history[i].Name == job.Name {
				scheduler.history[i].Status = Stopped
				scheduler.history[i].UpdatedAt = int(time.Now().Unix())
			}
		}
		break
	case Failed:
		for i := range scheduler.history {
			if scheduler.history[i].Name == job.Name {
				scheduler.history[i].Status = Failed
				scheduler.history[i].UpdatedAt = int(time.Now().Unix())
			}
		}
		break
	}
}

func (scheduler *SchedulerCapacity) Schedule(job Job) {
	scheduler.queuesMu.Lock()
	defer scheduler.queuesMu.Unlock()

	queue := job.Group
	_, ok := scheduler.queues[queue]
	if !ok {
		if InstanceOfGroupManager().get(queue) != nil {
			scheduler.queues[queue] = []Job{}
		} else {
			queue = "default"
		}
	}

	index := 0
	left := 0
	right := len(scheduler.queues[queue]) - 1
	for ; left <= right; {
		mid := (left + right) / 2
		if scheduler.queues[queue][left].Priority < job.Priority {
			index = left
			break
		}
		if scheduler.queues[queue][right].Priority >= job.Priority {
			index = right + 1
			break
		}
		if scheduler.queues[queue][mid].Priority >= job.Priority {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	scheduler.queues[queue] = append(scheduler.queues[queue], Job{})

	copy(scheduler.queues[queue][index+1:], scheduler.queues[queue][index:])
	scheduler.queues[queue][index] = job

	job.Status = Created
}

func (scheduler *SchedulerCapacity) AcquireResource(job Job) []NodeStatus {
	res := InstanceOfResourcePool().acquireResource(job)

	if len(res) != 0 {
		for _, task := range job.Tasks {

			scheduler.allocatingGPUMu.Lock()
			scheduler.allocatingGPU -= task.NumberGPU
			scheduler.allocatingGPUMu.Unlock()
		}
		log.Info("allocatingGPU is ", scheduler.allocatingGPU)

		go func(nodes []NodeStatus) {
			for _, node := range nodes {
				scheduler.resourceAllocationsMu.Lock()
				if _, ok := scheduler.resourceAllocations[job.Group]; !ok {
					scheduler.resourceAllocations[job.Group] = &ResourceCount{}
				}
				cnt, _ := scheduler.resourceAllocations[job.Group]
				cnt.CPU += node.MemTotal
				cnt.Memory += node.NumCPU
				for _, v := range node.Status {
					cnt.NumberGPU ++
					cnt.MemoryGPU += v.MemoryTotal
				}
				scheduler.resourceAllocationsMu.Unlock()
				scheduler.UpdateNextQueue()
			}

		}(res)
	}

	return res
}

func (scheduler *SchedulerCapacity) ReleaseResource(job Job, agent NodeStatus) {
	InstanceOfResourcePool().releaseResource(job, agent)

	scheduler.resourceAllocationsMu.Lock()
	if _, ok := scheduler.resourceAllocations[job.Group]; !ok {
		scheduler.resourceAllocations[job.Group] = &ResourceCount{}
	}
	cnt, _ := scheduler.resourceAllocations[job.Group]
	cnt.CPU -= agent.MemTotal
	cnt.Memory -= agent.NumCPU
	for _, v := range agent.Status {
		cnt.NumberGPU --
		cnt.MemoryGPU -= v.MemoryTotal
	}
	scheduler.resourceAllocationsMu.Unlock()
	go func(res NodeStatus) {
		scheduler.UpdateNextQueue()
	}(agent)
}

func (scheduler *SchedulerCapacity) QueryState(jobName string) MsgJobStatus {
	scheduler.queuesMu.Lock()
	jm, ok := scheduler.jobMasters[jobName]
	scheduler.queuesMu.Unlock()
	if !ok {
		return MsgJobStatus{Code: 1, Error: "Job not exist!"}
	}
	return jm.status()
}

func (scheduler *SchedulerCapacity) Stop(jobName string) MsgStop {
	scheduler.queuesMu.Lock()
	jm, ok := scheduler.jobMasters[jobName]
	scheduler.queuesMu.Unlock()
	if !ok {
		return MsgStop{Code: 1, Error: "Job not exist!"}
	}
	return jm.stop()
}

func (scheduler *SchedulerCapacity) QueryLogs(jobName string, taskName string) MsgLog {
	scheduler.queuesMu.Lock()
	jm, ok := scheduler.jobMasters[jobName]
	scheduler.queuesMu.Unlock()
	if !ok {
		return MsgLog{Code: 1, Error: "Job not exist!"}
	}
	return jm.logs(taskName)
}

func (scheduler *SchedulerCapacity) ListJobs() MsgJobList {
	var jobs []Job
	scheduler.historyMu.Lock()
	for _, job := range scheduler.history {
		jobs = append(jobs, *job)
	}
	scheduler.historyMu.Unlock()
	var tmp []Job
	for _, v := range scheduler.queues {
		tmp = append(tmp, v...)
	}
	sort.Sort(JobSorter(tmp))
	jobs = append(jobs, tmp...)
	return MsgJobList{Code: 0, Jobs: jobs}
}

func (scheduler *SchedulerCapacity) Summary() MsgSummary {
	summary := MsgSummary{}
	summary.Code = 0

	finishedJobsCounter := 0
	runningJobsCounter := 0
	pendingJobsCounter := 0

	var tmp []Job
	scheduler.historyMu.Lock()
	for _, job := range scheduler.history {
		tmp = append(tmp, *job)
	}
	scheduler.historyMu.Unlock()

	scheduler.queuesMu.Lock()
	for _, v := range scheduler.queues {
		tmp = append(tmp, v...)
	}
	scheduler.queuesMu.Unlock()

	for _, job := range tmp {
		switch job.Status {
		case Created:
			pendingJobsCounter++
		case Starting:
			pendingJobsCounter++
			break
		case Running:
			runningJobsCounter++
			break
		case Finished:
			finishedJobsCounter++
		case Stopped:
			finishedJobsCounter++
		}
	}
	summary.JobsFinished = finishedJobsCounter
	summary.JobsPending = pendingJobsCounter
	summary.JobsRunning = runningJobsCounter

	summary.FreeGPU, summary.UsingGPU = InstanceOfResourcePool().countGPU()
	return summary
}

func (scheduler *SchedulerCapacity) UpdateNextQueue() {
	next := "default"
	quota := math.MaxFloat64

	NumberGPU := float64(InstanceOfResourcePool().TotalGPU) + 0.00001

	scheduler.queuesMu.Lock()
	for k, t := range scheduler.queues {
		if len(t) == 0 {
			continue
		}
		scheduler.resourceAllocationsMu.Lock()
		if _, ok := scheduler.resourceAllocations[k]; !ok {
			scheduler.resourceAllocations[k] = &ResourceCount{}
		}
		v := scheduler.resourceAllocations[k]

		tmp := float64(v.NumberGPU) / NumberGPU
		scheduler.resourceAllocationsMu.Unlock()
		weight := 10
		if g, ok2 := InstanceOfGroupManager().groups[k]; !ok2 {
			weight = g.Weight
		}
		tmp /= float64(weight)
		if tmp < quota {
			quota = tmp
			next = k
		}
	}
	scheduler.nextQueue = next
	scheduler.queuesMu.Unlock()
	log.Debug("updateNextQueue ->", next)
}

func (scheduler *SchedulerCapacity) SetEnabled(enabled bool) bool {
	scheduler.enabled = enabled
	log.Info("scheduler is set to ", enabled)
	return true
}

func (scheduler *SchedulerCapacity) UpdateParallelism(parallelism int) bool {
	scheduler.parallelism = parallelism
	log.Info("parallelism is updated to ", parallelism)
	return true
}

func (scheduler *SchedulerCapacity) updateGroup(group Group) bool {
	return true
}

func (scheduler *SchedulerCapacity) DebugDump() map[string]interface{} {
	res := map[string]interface{}{}
	res["nextQueue"] = scheduler.nextQueue
	res["schedulingJobs"] = scheduler.schedulingJobs
	res["resourceAllocations"] = scheduler.resourceAllocations
	res["allocatingGPU"] = scheduler.allocatingGPU
	return res
}
