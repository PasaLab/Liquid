package main

import (
	"sync"
	"time"
	"math"
	"sort"
)

type SchedulerFair struct {
	history   []*Job
	historyMu sync.Mutex

	jobMasters map[string]*JobManager
	queues     map[string]JobList
	queuesMu   sync.Mutex

	drfyarn       bool
	enableBorrow  bool
	IOUs          map[string]map[string]*ResourceCount
	queuesQuota   map[string]*ResourceCount
	queuesQuotaMu sync.Mutex

	schedulingJobs map[string]bool
	schedulingMu   sync.Mutex

	resourceAllocations   map[string]*ResourceCount
	resourceAllocationsMu sync.Mutex

	enabled     bool
	parallelism int

	allocatingGPU   int
	allocatingGPUMu sync.Mutex
}

type JobList []Job

func (jobs JobList) Len() int {
	return len(jobs)
}

func (jobs JobList) Less(i, j int) bool {
	if jobs[i].Priority != jobs[j].Priority {
		return jobs[i].Priority < jobs[j].Priority
	}
	/* lower jobs, which unable to be scheduled */
	if InstanceOfResourcePool().TotalGPU < jobs[i].NumberGPU {
		return true
	}
	return jobs[i].BasePriority/float64(jobs[i].NumberGPU) < jobs[j].BasePriority/float64(jobs[j].NumberGPU)
}

func (jobs JobList) Swap(i, j int) {
	jobs[i], jobs[j] = jobs[j], jobs[i]
}

func (scheduler *SchedulerFair) Start() {
	log.Info("JS (fairness) started")

	scheduler.jobMasters = map[string]*JobManager{}
	scheduler.history = []*Job{}
	scheduler.queues = map[string]JobList{}
	scheduler.queues["default"] = []Job{}
	scheduler.drfyarn = false
	scheduler.enableBorrow = true
	scheduler.IOUs = map[string]map[string]*ResourceCount{}
	scheduler.queuesQuota = map[string]*ResourceCount{}
	scheduler.resourceAllocations = map[string]*ResourceCount{}
	scheduler.enabled = true
	scheduler.schedulingJobs = map[string]bool{}
	scheduler.allocatingGPU = 0

	scheduler.parallelism = 1

	go func() {
		flag := true
		for {
			log.Debug("Scheduling")
			if !flag { /* no more job */
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
			scheduler.queuesQuotaMu.Lock()
			/* choose queue which has the largest job */
			bestQueue := ""
			maxNumberGPU := math.MaxInt64
			maxNumberCPU := math.MaxInt64

			/* drf of yarn/kube-batch */
			if scheduler.drfyarn {
				least := math.MaxInt32
				for queue, jobs := range scheduler.queues {
					if len(jobs) == 0 {
						continue
					}
					if allocate, ok := scheduler.resourceAllocations[queue]; ok {
						if bestQueue == "" || allocate.NumberGPU < least {
							bestQueue = queue
							least = allocate.NumberGPU
						}
					} else {
						bestQueue = queue
						least = 0
						break
					}
				}
			}

			/* phase 1: execute jobs using self quota */
			if bestQueue == "" {
				for queue, jobs := range scheduler.queues {
					/* find largest job */
					if len(jobs) > 0 {
						/* calculate resource request of head job */
						numberGPUtmp := 0
						numberCPUtmp := 0
						for _, task := range jobs[0].Tasks {
							numberGPUtmp += task.NumberGPU
							numberCPUtmp += task.NumberCPU
						}
						/* if queue quota cannot satisfy, skip */
						if quota, ok := scheduler.queuesQuota[queue]; !ok || quota.NumberGPU/1000 < numberGPUtmp {
							continue
						}
						/* the more, the better */
						if bestQueue == "" || numberGPUtmp > maxNumberGPU || (numberGPUtmp == maxNumberGPU && numberCPUtmp > maxNumberCPU) {
							bestQueue = queue
							maxNumberGPU = numberGPUtmp
							maxNumberCPU = numberCPUtmp
						}
					}
				}
			}

			/* phase 2: borrow */
			if bestQueue == "" && scheduler.enableBorrow {
				/* calculate real quotas */
				quotas := map[string]*ResourceCount{}
				for queue, quota := range scheduler.queuesQuota {
					quotas[queue] = &ResourceCount{NumberGPU: quota.NumberGPU}
				}
				for q, IOUs := range scheduler.IOUs {
					for queue, IOU := range IOUs {
						quota := quotas[queue]
						quota.NumberGPU += IOU.NumberGPU
						quota = quotas[q]
						quota.NumberGPU -= IOU.NumberGPU
					}
				}
				/* firstly, check if quota sum can run a job */
				totalGPU := 0
				for _, quota := range scheduler.queuesQuota {
					totalGPU += quota.NumberGPU
				}
				/* find job which is short of least resource */
				minRequestGPU := math.MaxInt32
				minNeedBorrow := math.MaxInt32
				for queue, jobs := range scheduler.queues {
					if len(jobs) == 0 {
						continue
					}
					numberGPUtmp := 0
					for _, task := range jobs[0].Tasks {
						numberGPUtmp += task.NumberGPU
					}
					if _, ok := scheduler.queuesQuota[queue]; !ok {
						scheduler.queuesQuota[queue] = &ResourceCount{}
					}
					if _, ok := quotas[queue]; !ok {
						quotas[queue] = &ResourceCount{}
					}
					needGPU := numberGPUtmp*1000 - quotas[queue].NumberGPU
					/* the less, the better */
					if bestQueue == "" || needGPU < minRequestGPU {
						bestQueue = queue
						minRequestGPU = needGPU
						minNeedBorrow = numberGPUtmp*1000 - scheduler.queuesQuota[queue].NumberGPU
					}
				}
				if quota, ok := scheduler.queuesQuota[bestQueue]; ok {
					totalGPU -= quota.NumberGPU
				}
				/* if totalGPU can satisfy that job, start borrowing */
				if bestQueue != "" && totalGPU >= minNeedBorrow {
					log.Info("start borrow phase")
					log.Info(bestQueue, ": ", "total=", totalGPU, " still need ", minNeedBorrow)
					for {
						/* if all satisfied, break */
						if minNeedBorrow <= 0 {
							break
						}
						least := math.MaxInt32
						for queue, quota := range scheduler.queuesQuota {
							if queue == bestQueue {
								continue
							}
							if quota.NumberGPU > 0 && quota.NumberGPU < least {
								least = quota.NumberGPU
							}
						}
						if minNeedBorrow < least*(len(scheduler.queuesQuota)-1) {
							least = minNeedBorrow / (len(scheduler.queuesQuota) - 1)
						}
						/* start borrow */
						for queue, quota := range scheduler.queuesQuota {
							/* do not self borrow */
							if queue == bestQueue || quota.NumberGPU < least || least == 0 {
								continue
							}
							quota.NumberGPU -= least
							if _, ok := scheduler.IOUs[bestQueue]; !ok {
								scheduler.IOUs[bestQueue] = map[string]*ResourceCount{}
							}
							IOU, ok := scheduler.IOUs[bestQueue][queue]
							if !ok {
								scheduler.IOUs[bestQueue][queue] = &ResourceCount{}
								IOU = scheduler.IOUs[bestQueue][queue]
							}
							IOU.NumberGPU += least
							minNeedBorrow -= least
							scheduler.queuesQuota[bestQueue].NumberGPU += least

							log.Info(bestQueue, " borrow ", least, " from ", queue)
						}
						if least == 0 {
							for queue, quota := range scheduler.queuesQuota {
								/* do not self borrow */
								if queue == bestQueue || quota.NumberGPU == 0 {
									continue
								}
								if quota.NumberGPU < minNeedBorrow {
									least = quota.NumberGPU
								} else {
									least = minNeedBorrow
								}
								quota.NumberGPU -= least
								if _, ok := scheduler.IOUs[bestQueue]; !ok {
									scheduler.IOUs[bestQueue] = map[string]*ResourceCount{}
								}
								IOU, ok := scheduler.IOUs[bestQueue][queue]
								if !ok {
									scheduler.IOUs[bestQueue][queue] = &ResourceCount{}
									IOU = scheduler.IOUs[bestQueue][queue]
								}
								IOU.NumberGPU += least
								scheduler.queuesQuota[bestQueue].NumberGPU += least
								log.Info(bestQueue, " borrow ", minNeedBorrow, " from ", queue, " now ", scheduler.queuesQuota[bestQueue].NumberGPU)
								minNeedBorrow -= least
								break
							}
						}
					}

				} else {
					bestQueue = ""
				}
			}

			/* support schedule ahead & share */
			if bestQueue == "" && len(scheduler.schedulingJobs) == 0 {
				maxQuota := 0
				for queue, jobs := range scheduler.queues {
					if len(jobs) > 0 && len(jobs[0].Tasks) == 1 && jobs[0].Tasks[0].NumberGPU == 1 {
						if quota, ok := scheduler.queuesQuota[queue]; ok && (bestQueue == "" || quota.NumberGPU > maxQuota) {
							maxQuota = quota.NumberGPU
							bestQueue = queue
						}
					}
				}
			}

			/* launch that job */
			if bestQueue != "" {
				numberGPUtmp := 0
				numberCPUtmp := 0
				Memorytmp := 0
				for _, task := range scheduler.queues[bestQueue][0].Tasks {
					numberGPUtmp += task.NumberGPU
					numberCPUtmp += task.NumberCPU
					Memorytmp += task.Memory
				}

				log.Debug("schedulingJobs are ", scheduler.schedulingJobs)
				log.Debug("Before ")
				for queue, quota := range scheduler.queuesQuota {
					log.Debug("Queue<->", queue)
					log.Debug("GPU:", quota.NumberGPU)
					log.Debug("CPU:", quota.CPU)
					log.Debug("Memory:", quota.Memory)
				}
				pool := InstanceOfResourcePool()
				/* Make sure resource it enough */
				if len(scheduler.schedulingJobs) == 0 || (numberGPUtmp*10+(scheduler.allocatingGPU)*10 <= (pool.TotalGPU-pool.UsingGPU)*10) {
					flag = true

					log.Info("Before, ", scheduler.queuesQuota[bestQueue])
					if quota, ok := scheduler.queuesQuota[bestQueue]; ok {
						quota.NumberGPU -= numberGPUtmp * 1000
						quota.CPU -= numberCPUtmp * 1000
						quota.Memory -= Memorytmp
					}
					log.Info("After, ", scheduler.queuesQuota[bestQueue])

					scheduler.resourceAllocationsMu.Lock()
					if _, ok := scheduler.resourceAllocations[bestQueue]; !ok {
						scheduler.resourceAllocations[bestQueue] = &ResourceCount{}
					}
					cnt, _ := scheduler.resourceAllocations[bestQueue]
					cnt.NumberGPU += numberGPUtmp
					cnt.CPU += numberCPUtmp
					cnt.Memory += Memorytmp
					scheduler.resourceAllocationsMu.Unlock()

					scheduler.allocatingGPUMu.Lock()
					scheduler.allocatingGPU += numberGPUtmp
					scheduler.allocatingGPUMu.Unlock()
					log.Info("allocatingGPU is ", scheduler.allocatingGPU)

					jm := JobManager{}
					jm.job = scheduler.queues[bestQueue][0]
					jm.scheduler = scheduler
					jm.job.Status = Starting

					scheduler.jobMasters[jm.job.Name] = &jm
					scheduler.queues[bestQueue] = scheduler.queues[bestQueue][1:]

					scheduler.historyMu.Lock()
					scheduler.history = append(scheduler.history, &jm.job)
					scheduler.historyMu.Unlock()

					scheduler.schedulingMu.Lock()
					scheduler.schedulingJobs[jm.job.Name] = true
					scheduler.schedulingMu.Unlock()
					go func() {
						jm.start()
					}()
				}
			} else {
				log.Debug("No more jobs to scheduling ", time.Now())
				go func() {
					scheduler.UpdateQuota()
				}()
			}
			scheduler.queuesQuotaMu.Unlock()
			scheduler.queuesMu.Unlock()
		}
	}()
}

func (scheduler *SchedulerFair) UpdateProgress(job Job, state State) {
	scheduler.historyMu.Lock()
	defer scheduler.historyMu.Unlock()

	scheduler.schedulingMu.Lock()
	if _, ok := scheduler.schedulingJobs[job.Name]; ok {
		delete(scheduler.schedulingJobs, job.Name)
		scheduler.allocatingGPU -= job.NumberGPU
		log.Info("allocatingGPU is ", scheduler.allocatingGPU)
	}
	scheduler.schedulingMu.Unlock()

	switch state {
	case Running:
		for i := range scheduler.history {
			if scheduler.history[i].Name == job.Name {
				scheduler.history[i].Status = Running
				scheduler.history[i].UpdatedAt = int(time.Now().Unix())
				scheduler.history[i].StartedAt = time.Now().Unix()
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

func (scheduler *SchedulerFair) Schedule(job Job) {
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

	numberGPU := 0
	for _, task := range job.Tasks {
		numberGPU += task.NumberGPU
	}
	job.NumberGPU = numberGPU

	job.Status = Created

	job.BasePriority = -float64(time.Now().UnixNano() / 100000) / 100000000000000
	scheduler.queues[queue][index] = job
}

func (scheduler *SchedulerFair) AcquireResource(job Job) []NodeStatus {
	res := InstanceOfResourcePool().acquireResource(job)
	go func() {
		scheduler.UpdateQuota()
	}()
	return res
}

func (scheduler *SchedulerFair) ReleaseResource(job Job, agent NodeStatus) {
	InstanceOfResourcePool().releaseResource(job, agent)

	scheduler.resourceAllocationsMu.Lock()
	if _, ok := scheduler.resourceAllocations[job.Group]; !ok {
		scheduler.resourceAllocations[job.Group] = &ResourceCount{}
	}
	cnt, _ := scheduler.resourceAllocations[job.Group]
	cnt.CPU -= agent.NumCPU
	cnt.Memory -= agent.MemTotal
	cnt.NumberGPU -= len(agent.Status)
	scheduler.resourceAllocationsMu.Unlock()

	go func() {
		scheduler.UpdateQuota()
	}()
}

/* allocate quota to queues */
func (scheduler *SchedulerFair) UpdateQuota() {
	scheduler.queuesMu.Lock()
	defer scheduler.queuesMu.Unlock()
	scheduler.queuesQuotaMu.Lock()
	defer scheduler.queuesQuotaMu.Unlock()
	//log.Info("Updating queues quota~")

	/* phase 1: DRF */
	usingGPU := 0
	usingCPU := 0
	usingMemory := 0
	allocatedGPU := 0
	allocatedCPU := 0
	allocatedMemory := 0
	scheduler.resourceAllocationsMu.Lock()
	for _, quota := range scheduler.resourceAllocations {
		usingGPU += quota.NumberGPU
		usingCPU += quota.CPU
		usingMemory += quota.Memory
	}
	scheduler.resourceAllocationsMu.Unlock()

	for _, quota := range scheduler.queuesQuota {
		allocatedGPU += quota.NumberGPU
		allocatedCPU += quota.CPU
		allocatedMemory += quota.Memory
	}

	pool := InstanceOfResourcePool()

	availableGPU := pool.TotalGPU*1000 - usingGPU*1000 - allocatedGPU
	availableCPU := pool.TotalCPU*1000 - usingCPU*1000 - allocatedCPU
	//availableMemory := pool.TotalMemory - usingMemory - allocatedMemory
	/* <0 means some nodes exited */
	//log.Info(availableGPU)
	if availableGPU <= 0 {
		return
	}

	var candidates []string
	requests := map[string]ResourceCount{}
	weights := 0

	for queue, jobs := range scheduler.queues {
		if len(jobs) == 0 {
			continue
		}
		weights += InstanceOfGroupManager().groups[queue].Weight
		request := ResourceCount{}
		for i, job := range jobs {
			GPU := 0
			CPU := 0
			Memory := 0
			for _, task := range job.Tasks {
				GPU += task.NumberGPU
				CPU += task.NumberCPU
				Memory += task.Memory
			}
			request.NumberGPU += GPU
			request.CPU += CPU
			request.Memory += Memory
			/* increase priority at most 10 jobs, to avoid small jobs always goes first in a batch */
			if job.Priority == jobs[0].Priority && i < 10 {
				scheduler.queues[queue][i].BasePriority += 1.0
			}
		}
		sort.Sort(sort.Reverse(scheduler.queues[queue]))

		/* minimum is 1, to avoid divide by zero error following */
		if request.NumberGPU == 0 {
			request.NumberGPU = 1
		}
		if quota, ok := scheduler.queuesQuota[queue]; ok && quota.NumberGPU >= request.NumberGPU*1000 {
			continue
		}
		requests[queue] = request
		candidates = append(candidates, queue)
	}

	if len(candidates) == 0 {
		return
	}
	log.Debug("Can allocate ", availableGPU)
	log.Debug("Before ")
	for queue, quota := range scheduler.queuesQuota {
		log.Debug("Queue<->", queue)
		log.Debug("GPU:", quota.NumberGPU)
		log.Debug("CPU:", quota.CPU)
		log.Debug("Memory:", quota.Memory)
	}

	per := availableGPU / weights
	for _, queue := range candidates {
		if _, ok := scheduler.queuesQuota[queue]; !ok {
			scheduler.queuesQuota[queue] = &ResourceCount{}
		}
		weight := InstanceOfGroupManager().groups[queue].Weight
		quota := scheduler.queuesQuota[queue]

		/* if allocate is more than request, reduce weight */
		log.Info(quota.NumberGPU, per, weight, requests[queue].NumberGPU)
		base := per * weight
		if quota.NumberGPU+base > requests[queue].NumberGPU*1000 {
			base = requests[queue].NumberGPU*1000 - quota.NumberGPU
		}

		quota.NumberGPU += base
		availableGPU -= base

		quota.CPU += (requests[queue].CPU * base) / requests[queue].NumberGPU
		availableCPU -= (requests[queue].CPU * base) / requests[queue].NumberGPU
		quota.Memory += ((requests[queue].Memory * base) / requests[queue].NumberGPU) / 1000
	}
	/* avoid resource leak, and reserve full */
	availableGPU = availableGPU % 1000
	if availableGPU > 0 {
		for _, queue := range candidates {
			quota := scheduler.queuesQuota[queue]
			quota.NumberGPU += availableGPU
			quota.CPU += (requests[queue].CPU * availableGPU) / requests[queue].NumberGPU
			quota.Memory += ((requests[queue].Memory * availableGPU) / requests[queue].NumberGPU) / 1000
			break
		}
	}
	log.Debug("After ")
	for queue, quota := range scheduler.queuesQuota {
		log.Debug("Queue<->", queue)
		log.Debug("GPU:", quota.NumberGPU)
		log.Debug("CPU:", quota.CPU)
		log.Debug("Memory:", quota.Memory)
	}

	/* Phase 2: clear IOUs */
	for queue, IOUs := range scheduler.IOUs {
		/* no IOU, skip */
		if t, ok := scheduler.IOUs[queue]; !ok || len(t) == 0 {
			continue
		}
		/* nothing to pay */
		if tmp, ok := scheduler.queuesQuota[queue]; !ok || tmp.NumberGPU == 0 {
			continue
		}
		minIOU := 0
		totalIOU := 0
		for _, IOU := range IOUs {
			if IOU.NumberGPU > minIOU && IOU.NumberGPU != 0 {
				minIOU = IOU.NumberGPU
				totalIOU += IOU.NumberGPU
			}
		}
		quota := scheduler.queuesQuota[queue]
		if quota.NumberGPU >= totalIOU {
			/* can clear all */
			minIOU = totalIOU
		}
		if quota.NumberGPU < minIOU*len(IOUs) {
			minIOU = quota.NumberGPU / len(IOUs)
		}

		for q, IOU := range IOUs {
			if IOU.NumberGPU <= minIOU {
				quota.NumberGPU -= IOU.NumberGPU
				scheduler.queuesQuota[q].NumberGPU += IOU.NumberGPU
				IOU.NumberGPU = 0
			} else {
				quota.NumberGPU -= minIOU
				scheduler.queuesQuota[q].NumberGPU += minIOU
				IOU.NumberGPU -= minIOU
			}
			log.Info(queue, " pay IOU to ", q, " now ", IOU.NumberGPU)
			/* clear */
			if IOU.NumberGPU == 0 {
				delete(scheduler.IOUs[queue], q)
			}
		}

		if minIOU == 0 {
			for q, IOU := range IOUs {
				quota.NumberGPU -= 1
				scheduler.queuesQuota[q].NumberGPU += 1
				IOU.NumberGPU -= 1
				log.Info(queue, " pay IOU to ", q, " now ", IOU.NumberGPU)
				/* clear */
				if IOU.NumberGPU == 0 {
					delete(scheduler.IOUs[queue], q)
				}
				if quota.NumberGPU == 0 {
					break
				}
			}
		}
	}
}

func (scheduler *SchedulerFair) QueryState(jobName string) MsgJobStatus {
	scheduler.queuesMu.Lock()
	jm, ok := scheduler.jobMasters[jobName]
	scheduler.queuesMu.Unlock()
	if !ok {
		return MsgJobStatus{Code: 1, Error: "Job not exist!"}
	}
	return jm.status()
}

func (scheduler *SchedulerFair) Stop(jobName string) MsgStop {
	log.Info("Stop job ", jobName)
	scheduler.queuesMu.Lock()
	jm, ok := scheduler.jobMasters[jobName]
	scheduler.queuesMu.Unlock()
	if ok {
		return jm.stop()
	} else {
		found := false
		for queue := range scheduler.queues {
			index := -1
			for i, job := range scheduler.queues[queue] {
				if job.Name == jobName {
					index = i
				}
			}
			log.Info(index)
			if index != -1 {
				(&scheduler.queues[queue][index]).Status = Stopped
				scheduler.historyMu.Lock()
				job := scheduler.queues[queue][index]
				scheduler.history = append(scheduler.history, &job)
				scheduler.historyMu.Unlock()
				copy(scheduler.queues[queue][index:], scheduler.queues[queue][index+1:])
				scheduler.queues[queue] = scheduler.queues[queue][:len(scheduler.queues[queue])-1]
				found = true
				break
			}
		}
		if found {
			return MsgStop{Code: 0}
		}
	}
	return MsgStop{Code: 1, Error: "Job not exist!"}
}

func (scheduler *SchedulerFair) QueryLogs(jobName string, taskName string) MsgLog {
	scheduler.queuesMu.Lock()
	jm, ok := scheduler.jobMasters[jobName]
	scheduler.queuesMu.Unlock()
	if !ok {
		return MsgLog{Code: 1, Error: "Job not exist!"}
	}
	return jm.logs(taskName)
}

func (scheduler *SchedulerFair) ListJobs() MsgJobList {
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

func (scheduler *SchedulerFair) Summary() MsgSummary {
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

func (scheduler *SchedulerFair) SetEnabled(enabled bool) bool {
	scheduler.enabled = enabled
	log.Info("scheduler is set to ", enabled)
	return true
}

func (scheduler *SchedulerFair) UpdateParallelism(parallelism int) bool {
	if parallelism < 1 {
		parallelism = 1
	}
	scheduler.parallelism = parallelism
	log.Info("parallelism is updated to ", parallelism)
	return true
}

func (scheduler *SchedulerFair) updateGroup(group Group) bool {
	return true
}

func (scheduler *SchedulerFair) DebugDump() map[string]interface{} {
	res := map[string]interface{}{}
	res["queuesQuota"] = scheduler.queuesQuota
	res["schedulingJobs"] = scheduler.schedulingJobs
	res["resourceAllocations"] = scheduler.resourceAllocations
	res["allocatingGPU"] = scheduler.allocatingGPU
	res["IOUs"] = scheduler.IOUs
	res["queues"] = scheduler.queues
	res["history"] = scheduler.history
	return res
}
