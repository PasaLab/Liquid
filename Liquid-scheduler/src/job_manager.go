package main

import (
	"time"
	"net/url"
	"strings"
	"io/ioutil"
	"encoding/json"
	"sync"
	"strconv"
	"math/rand"
)

type JobManager struct {
	/* meta */
	scheduler Scheduler
	job       Job

	/* resource */
	network     string
	resources   map[string]NodeStatus
	resourcesMu sync.Mutex

	/* status */
	jobStatus     JobStatus
	isRunning     bool
	lastHeartBeat map[string]int64
	statusMu      sync.Mutex

	/* history info */
	stats [][]TaskStatus
}

func (jm *JobManager) start() {
	log.Info("start job ", jm.job.Name, " at ", time.Now())
	jm.isRunning = true
	jm.lastHeartBeat = map[string]int64{}
	jm.jobStatus = JobStatus{Name: jm.job.Name, tasks: map[string]TaskStatus{}}
	jm.resources = map[string]NodeStatus{}

	/* register in JHL */
	InstanceJobHistoryLogger().submitJob(jm.job)

	/* request for resources */
	jm.resourcesMu.Lock()
	jm.network = InstanceOfResourcePool().acquireNetwork()
	for {
		if !jm.isRunning {
			break
		}
		resources := jm.scheduler.AcquireResource(jm.job)
		if len(resources) > 0 {
			for i, node := range resources {
				jm.resources[jm.job.Tasks[i].Name] = node
			}
			log.Info(jm.job.Name, " receive resource ", jm.resources)
			break
		}
		/* sleep random Millisecond to avoid deadlock */
		time.Sleep(time.Millisecond * time.Duration(500+rand.Intn(500)))
	}
	jm.resourcesMu.Unlock()

	if InstanceOfConfiguration().mock {
		if jm.isRunning {
			jm.scheduler.UpdateProgress(jm.job, Running)
			duration := InstanceOfMocker().GetDuration(jm.job, jm.resources)
			log.Info("mock ", jm.job.Name, ", wait ", duration)
			time.Sleep(time.Second * time.Duration(duration))
			jm.isRunning = false
			jm.scheduler.UpdateProgress(jm.job, Finished)
		}
		jm.returnResource()
		log.Info(jm.job.Name, "JobMaster exited ")
		return
	}

	isShare := false
	isScheduleAhead := false
	if jm.isRunning {
		/* switch to Running state */
		jm.scheduler.UpdateProgress(jm.job, Running)

		/* bring up containers */
		wg := sync.WaitGroup{}
		success := true
		for i, task := range jm.job.Tasks {
			wg.Add(1)

			go func(task Task, node NodeStatus) {
				defer wg.Done()
				var UUIDs []string
				shouldWait := "0"
				for _, GPU := range node.Status {
					UUIDs = append(UUIDs, GPU.UUID)
					if GPU.MemoryUsed == GPU.MemoryTotal {
						shouldWait = "1"
						isScheduleAhead = true
					} else if GPU.MemoryAllocated > 0 {
						isShare = true
					}
					/* attach to GPUs */
					InstanceOfResourcePool().attach(GPU.UUID, jm.job)
				}
				GPUs := strings.Join(UUIDs, ",")

				v := url.Values{}
				v.Set("image", task.Image)
				v.Set("cmd", task.Cmd)
				v.Set("name", task.Name)
				v.Set("workspace", jm.job.Workspace)
				v.Set("gpus", GPUs)
				v.Set("mem_limit", strconv.Itoa(task.Memory)+"m")
				v.Set("cpu_limit", strconv.Itoa(task.NumberCPU))
				v.Set("network", jm.network)
				v.Set("should_wait", shouldWait)
				v.Set("output_dir", "/tmp/")
				v.Set("hdfs_address", InstanceOfConfiguration().HDFSAddress)
				v.Set("hdfs_dir", InstanceOfConfiguration().HDFSBaseDir+jm.job.Name)
				v.Set("gpu_mem", strconv.Itoa(task.MemoryGPU))
				if InstanceOfConfiguration().DFSBaseDir != "" {
					v.Set("dfs_src", InstanceOfConfiguration().DFSBaseDir+jm.job.Name+"/task-"+task.Name)
				} else {
					v.Set("dfs_src", "")
				}
				v.Set("dfs_dst", "/tmp")

				spider := Spider{}
				spider.Method = "POST"
				spider.URL = "http://" + node.ClientHost + ":8000/create"
				spider.Data = v
				spider.ContentType = "application/x-www-form-urlencoded"
				err := spider.do()
				if err != nil {
					log.Warn(err.Error())
					success = false
					return
				}
				resp := spider.getResponse()

				body, err := ioutil.ReadAll(resp.Body)
				resp.Body.Close()
				if err != nil {
					log.Warn(err)
					success = false
					return
				}

				var res MsgCreate
				err = json.Unmarshal([]byte(string(body)), &res)
				if err != nil || res.Code != 0 {
					log.Warn(res)
					success = false
					return
				}
				taskStatus := TaskStatus{Id: res.Id, Node: node.ClientHost, HostName: jm.job.Tasks[i].Name}
				jm.statusMu.Lock()
				jm.jobStatus.tasks[task.Name] = taskStatus
				jm.lastHeartBeat[task.Name] = time.Now().Unix()
				jm.statusMu.Unlock()

			}(task, jm.resources[task.Name])
		}
		wg.Wait()
		/* start failed */
		if !success {
			jm.isRunning = false
			jm.scheduler.UpdateProgress(jm.job, Failed)
			jm.stop()
		} else {
			log.Info(jm.job.Name, " all tasks launched success")
		}
	}

	/* monitor job execution */
	for {
		if !jm.isRunning {
			break
		}
		now := time.Now().Unix()
		jm.statusMu.Lock()
		for task, pre := range jm.lastHeartBeat {
			if now-pre > 30 {
				log.Warn(jm.job.Name, "-", task, " heartbeat longer than 30s")
			}
		}
		jm.statusMu.Unlock()
		time.Sleep(time.Second * 1)
	}

	/* release again to make sure resources are released */
	jm.stop()
	jm.returnResource()

	/* feed data to optimizer */
	isExclusive := InstanceOfResourcePool().isExclusive(jm.job.Name)

	var stats [][]TaskStatus
	for _, vals := range jm.stats {
		var stat []TaskStatus
		for i, task := range jm.job.Tasks {
			if task.IsPS {
				stat = append(stat, vals[i])
			}
		}
		if len(stat) > 0 {
			stats = append(stats, stat)
		}
	}
	if isExclusive {
		InstanceOfOptimizer().FeedStats(jm.job, "PS", stats)
	}
	stats = [][]TaskStatus{}
	for _, vals := range jm.stats {
		var stat []TaskStatus
		for i, task := range jm.job.Tasks {
			if !task.IsPS {
				stat = append(stat, vals[i])
			}
		}
		if len(stat) > 0 {
			stats = append(stats, stat)
		}
	}
	if isExclusive {
		InstanceOfOptimizer().FeedStats(jm.job, "Worker", stats)
	}

	if len(jm.job.Tasks) == 1 && !isShare && !isScheduleAhead && jm.job.Status == Finished && isExclusive {
		InstanceOfOptimizer().FeedTime(jm.job, stats)
	}

	/* clear, to reduce memory usage */
	jm.stats = [][]TaskStatus{}

	/* remove exited containers */
	//for _, task := range jm.jobStatus.tasks {
	//	go func(container TaskStatus) {
	//		v := url.Values{}
	//		v.Set("id", container.Id)
	//
	//		spider := Spider{}
	//		spider.Method = "POST"
	//		spider.URL = "http://" + container.Node + ":8000/remove"
	//		spider.Data = v
	//		spider.ContentType = "application/x-www-form-urlencoded"
	//		err := spider.do()
	//		if err != nil {
	//			log.Warn(err.Error())
	//		}
	//	}(task)
	//}

	log.Info(jm.job.Name, " JobMaster exited ")
}

/* release all resource */
func (jm *JobManager) returnResource() {
	jm.resourcesMu.Lock()
	defer jm.resourcesMu.Unlock()
	/* return resource */
	for i := range jm.resources {
		jm.scheduler.ReleaseResource(jm.job, jm.resources[i])
		for _, t := range jm.resources[i].Status {
			InstanceOfResourcePool().detach(t.UUID, jm.job)
		}
	}
	jm.resources = map[string]NodeStatus{}
	if jm.network != "" {
		InstanceOfResourcePool().releaseNetwork(jm.network)
		jm.network = ""
	}
}

/* monitor all tasks, update job status */
func (jm *JobManager) checkStatus(status []TaskStatus) {
	if !jm.isRunning {
		return
	}
	flagRunning := false
	onlyPS := true
	for i := range status {
		if status[i].Status == "ready" || status[i].Status == "running" || status[i].Status == "launching" {
			flagRunning = true
			if !jm.job.Tasks[i].IsPS {
				onlyPS = false
			}
			InstanceJobHistoryLogger().submitTaskStatus(jm.job.Name, status[i])
		} else if status[i].Status == "unknown" {
			flagRunning = true
			if !jm.job.Tasks[i].IsPS {
				onlyPS = false
			}
		} else {
			log.Info(jm.job.Name, "-", i, " ", status[i].Status)
			if exitCode, ok := status[i].State["ExitCode"].(float64); ok && exitCode != 0 && jm.isRunning {
				log.Warn(jm.job.Name+"-"+jm.job.Tasks[i].Name+" exited unexpected, exitCode=", exitCode)
				jm.isRunning = false
				jm.scheduler.UpdateProgress(jm.job, Failed)
				jm.stop()
			} else if jm.isRunning {
				log.Info(jm.job.Name, " Some instance exited, close others")
				jm.isRunning = false
				jm.scheduler.UpdateProgress(jm.job, Finished)
				jm.stop()
			}

			jm.resourcesMu.Lock()
			nodeID := jm.job.Tasks[i].Name
			if _, ok := jm.resources[nodeID]; ok {
				jm.scheduler.ReleaseResource(jm.job, jm.resources[nodeID])
				log.Info(jm.job.Name, " return resource ", jm.resources[nodeID].ClientID)

				for _, t := range jm.resources[nodeID].Status {
					InstanceOfResourcePool().detach(t.UUID, jm.job)
				}
				InstanceJobHistoryLogger().submitTaskStatus(jm.job.Name, status[i])
				delete(jm.resources, nodeID)
			}
			jm.resourcesMu.Unlock()
		}
		jm.statusMu.Lock()
		jm.lastHeartBeat[jm.job.Tasks[i].Name] = time.Now().Unix()
		jm.statusMu.Unlock()
	}
	if flagRunning && onlyPS && jm.isRunning {
		log.Info(jm.job.Name, " Only PS is running, stop ")
		jm.isRunning = false
		jm.scheduler.UpdateProgress(jm.job, Finished)
		jm.stop()
	}

	if !flagRunning && jm.isRunning {
		log.Info(jm.job.Name, " finish job ")
		jm.isRunning = false
		jm.scheduler.UpdateProgress(jm.job, Finished)
	}
}

/* fetch logs of task */
func (jm *JobManager) logs(taskName string) MsgLog {
	spider := Spider{}
	spider.Method = "GET"
	jm.statusMu.Lock()
	spider.URL = "http://" + jm.jobStatus.tasks[taskName].Node + ":8000/logs?id=" + jm.jobStatus.tasks[taskName].Id
	_, ok := jm.jobStatus.tasks[taskName]
	jm.statusMu.Unlock()
	if !ok {
		return MsgLog{Code: -1, Error: "Task not exist"}
	}

	err := spider.do()
	if err != nil {
		return MsgLog{Code: 1, Error: err.Error()}
	}

	resp := spider.getResponse()
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return MsgLog{Code: 2, Error: err.Error()}
	}

	var res MsgLog
	err = json.Unmarshal([]byte(string(body)), &res)
	if err != nil {
		log.Warn(err)
		return MsgLog{Code: 3, Error: "Unknown"}
	}
	return res
}

/* fetch job tasks status */
func (jm *JobManager) status() MsgJobStatus {
	var tasksStatus []TaskStatus
	/* create slice ahead, since append would cause uncertain order */
	for range jm.job.Tasks {
		tasksStatus = append(tasksStatus, TaskStatus{})
	}

	for i, task := range jm.job.Tasks {
		jm.statusMu.Lock()
		taskStatus := jm.jobStatus.tasks[task.Name]
		jm.statusMu.Unlock()

		/* still in launching phase */
		if len(taskStatus.Node) == 0 {
			tasksStatus[i] = TaskStatus{Status: "launching", State: map[string]interface{}{"ExitCode": float64(0)}}
			continue
		}

		spider := Spider{}
		spider.Method = "GET"
		spider.URL = "http://" + taskStatus.Node + ":8000/status?id=" + taskStatus.Id

		err := spider.do()
		if err != nil {
			log.Warn(err)
			tasksStatus[i] = TaskStatus{Status: "unknown", State: map[string]interface{}{"ExitCode": float64(-1)}}
			continue
		}

		resp := spider.getResponse()
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			tasksStatus[i] = TaskStatus{Status: "unknown", State: map[string]interface{}{"ExitCode": float64(-1)}}
			continue
		}

		var res MsgTaskStatus
		err = json.Unmarshal([]byte(string(body)), &res)
		if err != nil {
			log.Warn(err)
			tasksStatus[i] = TaskStatus{Status: "unknown", State: map[string]interface{}{"ExitCode": float64(-1)}}
			continue
		}
		if res.Code == 2 {
			tasksStatus[i] = TaskStatus{Status: "unknown", State: map[string]interface{}{"ExitCode": float64(-2)}}
			log.Warn(res.Error)
			continue
		}
		if res.Code != 0 {
			tasksStatus[i] = TaskStatus{Status: "notexist", State: map[string]interface{}{"ExitCode": float64(res.Code)}}
			continue
		}
		res.Status.Node = taskStatus.Node
		tasksStatus[i] = res.Status
	}
	for i := range jm.job.Tasks {
		tasksStatus[i].TimeStamp = time.Now().Unix()
	}

	if jm.isRunning {
		go func() {
			jm.checkStatus(tasksStatus)
		}()
		jm.stats = append(jm.stats, tasksStatus)

	}
	return MsgJobStatus{Status: tasksStatus}
}

func (jm *JobManager) stop() MsgStop {
	if jm.isRunning {
		jm.isRunning = false
		jm.scheduler.UpdateProgress(jm.job, Stopped)
		log.Info("kill job, ", jm.job.Name)
	}

	jm.statusMu.Lock()
	for _, taskStatus := range jm.jobStatus.tasks {
		/* stop at background */
		go func(task TaskStatus) {
			log.Info("kill ", jm.job.Name, "-", task.Id, " :", task.HostName)
			v := url.Values{}
			v.Set("id", task.Id)

			spider := Spider{}
			spider.Method = "POST"
			spider.URL = "http://" + task.Node + ":8000/stop"
			spider.Data = v
			spider.ContentType = "application/x-www-form-urlencoded"

			err := spider.do()
			if err != nil {
				log.Warn(err.Error())
				return
			}
			resp := spider.getResponse()
			body, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				log.Warn(err)
				return
			}
			var res MsgStop
			err = json.Unmarshal([]byte(string(body)), &res)
			if err != nil || res.Code != 0 {
				log.Warn(res)
				return
			}
			if res.Code != 0 {
				log.Warn(res.Error)
			}
		}(taskStatus)
	}
	jm.statusMu.Unlock()
	return MsgStop{Code: 0}
}
