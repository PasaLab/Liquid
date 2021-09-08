package main

import (
	"sync"
	"time"
	"sort"
	"math"
)

type SchedulerPriority struct {
	history   []*Job
	historyMu sync.Mutex

	queue   []Job
	queueMu sync.Mutex

	schedulingJobs map[string]bool
	schedulingMu   sync.Mutex

	jobMasters  map[string]*JobManager
	enabled     bool
	parallelism int
}

func (scheduler *SchedulerPriority) Start() {
	scheduler.jobMasters = map[string]*JobManager{}
	scheduler.history = []*Job{}
	scheduler.enabled = true
	scheduler.parallelism = 1
	scheduler.schedulingJobs = map[string]bool{}

	go func() {
		flag := true
		for {
			log.Debug("Scheduling")
			if !flag { /* no more job */
				time.Sleep(time.Millisecond * 100)
			}
			flag = false
			scheduler.schedulingMu.Lock()
			if len(scheduler.schedulingJobs) >= scheduler.parallelism {
				scheduler.schedulingMu.Unlock()
				time.Sleep(time.Millisecond * 100)
				continue
			}
			scheduler.schedulingMu.Unlock()

			scheduler.queueMu.Lock()
			if len(scheduler.queue) > 0 {

				numberGPU := 0
				for _, task := range scheduler.queue[0].Tasks {
					numberGPU += task.NumberGPU
				}
				if numberGPU <= (InstanceOfResourcePool().TotalGPU - InstanceOfResourcePool().UsingGPU) {

					jm := JobManager{}
					jm.job = scheduler.queue[0]
					scheduler.queue = scheduler.queue[1:]
					jm.scheduler = scheduler
					scheduler.jobMasters[jm.job.Name] = &jm

					jm.job.Status = Starting
					scheduler.schedulingMu.Lock()
					scheduler.schedulingJobs[jm.job.Name] = true
					scheduler.schedulingMu.Unlock()
					scheduler.historyMu.Lock()
					scheduler.history = append(scheduler.history, &jm.job)
					scheduler.historyMu.Unlock()

					go func() {
						jm.start()
					}()
				} else if InstanceOfConfiguration().PreemptEnabled {
					/* start preempt */
					var jobs []Job
					preemptee := scheduler.queue[0]
					lowest := preemptee.Priority - 1
					scheduler.historyMu.Lock()
					for _, job := range scheduler.history {
						if job.Status != Running {
							continue
						}
						if job.Priority < lowest {
							jobs = []Job{*job}
							lowest = job.Priority
						} else if job.Priority == lowest {
							jobs = append(jobs, *job)
						}
					}
					sort.Sort(JobSorter(jobs))
					if len(jobs) > 0 {
						preempted := jobs[0]
						minScore := math.MaxFloat64
						for _, jobT := range jobs {
							score := 0.0
							numberGPUt := 0
							for _, task := range jobT.Tasks {
								numberGPUt += task.NumberGPU
							}

							needGPU := numberGPU - (InstanceOfResourcePool().TotalGPU - InstanceOfResourcePool().UsingGPU)
							score = float64(time.Now().UnixNano()/100000-int64(jobT.CreatedAt)) * math.Abs(0.01+float64(numberGPU-needGPU)) / float64(numberGPUt)

							if score < minScore {
								minScore = score
								preempted = jobT
							}
						}

						/* Remove from history */
						idx := -1
						for i, job := range scheduler.history {
							if job.Name == preempted.Name {
								idx = i
							}
						}
						if idx != -1 {
							copy(scheduler.history[idx:], scheduler.history[idx+1:])
							scheduler.history = scheduler.history[:len(scheduler.history)-1]
						}
						scheduler.historyMu.Unlock()

						before := InstanceOfResourcePool().UsingGPU
						log.Info("Start preempt ", preempted.Name)
						scheduler.Stop(preempted.Name)

						/* add back */
						idx = len(scheduler.queue)
						for i, job := range scheduler.queue {
							if job.Priority > preempted.Priority {
								continue
							}
							idx = i
							break
						}
						scheduler.queue = append(scheduler.queue, Job{})

						preempted.Status = Created
						copy(scheduler.queue[idx+1:], scheduler.queue[idx:])
						scheduler.queue[idx] = preempted

						delete(scheduler.jobMasters, preempted.Name)

						for {
							after := InstanceOfResourcePool().UsingGPU
							if after != before {
								break
							}
							log.Info("before=", before, " after=", after)
							time.Sleep(time.Millisecond * 100)
						}
					} else {
						scheduler.historyMu.Unlock()
					}
				}
			}
			scheduler.queueMu.Unlock()
		}
	}()
}

func (scheduler *SchedulerPriority) UpdateProgress(job Job, state State) {
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

func (scheduler *SchedulerPriority) Schedule(job Job) {
	scheduler.queueMu.Lock()
	defer scheduler.queueMu.Unlock()

	index := 0

	left := 0
	right := len(scheduler.queue) - 1
	for ; left <= right; {
		mid := (left + right) / 2
		if scheduler.queue[left].Priority < job.Priority {
			index = left
			break
		}
		if scheduler.queue[right].Priority >= job.Priority {
			index = right + 1
			break
		}
		if scheduler.queue[mid].Priority >= job.Priority {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	scheduler.queue = append(scheduler.queue, Job{})

	copy(scheduler.queue[index+1:], scheduler.queue[index:])
	scheduler.queue[index] = job

	job.Status = Created
}

func (scheduler *SchedulerPriority) AcquireResource(job Job) []NodeStatus {
	res := InstanceOfResourcePool().acquireResource(job)
	return res
}

func (scheduler *SchedulerPriority) ReleaseResource(job Job, agent NodeStatus) {
	InstanceOfResourcePool().releaseResource(job, agent)
}

func (scheduler *SchedulerPriority) QueryState(jobName string) MsgJobStatus {
	jm, ok := scheduler.jobMasters[jobName]
	if !ok {
		return MsgJobStatus{Code: 1, Error: "Job not exist!"}
	}
	return jm.status()
}

func (scheduler *SchedulerPriority) Stop(jobName string) MsgStop {
	jm, ok := scheduler.jobMasters[jobName]
	if !ok {
		return MsgStop{Code: 1, Error: "Job not exist!"}
	}
	return jm.stop()
}

func (scheduler *SchedulerPriority) QueryLogs(jobName string, taskName string) MsgLog {
	jm, ok := scheduler.jobMasters[jobName]
	if !ok {
		return MsgLog{Code: 1, Error: "Job not exist!"}
	}
	return jm.logs(taskName)
}

func (scheduler *SchedulerPriority) ListJobs() MsgJobList {
	var tmp []Job
	for _, job := range scheduler.history {
		tmp = append(tmp, *job)
	}
	tmp = append(tmp, scheduler.queue...)
	return MsgJobList{Code: 0, Jobs: tmp}
}

func (scheduler *SchedulerPriority) Summary() MsgSummary {
	summary := MsgSummary{}
	summary.Code = 0

	finishedJobsCounter := 0
	runningJobsCounter := 0
	pendingJobsCounter := 0

	var tmp []Job
	for _, job := range scheduler.history {
		tmp = append(tmp, *job)
	}
	tmp = append(tmp, scheduler.queue...)

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

func (scheduler *SchedulerPriority) SetEnabled(enabled bool) bool {
	scheduler.enabled = enabled
	log.Info("scheduler is set to ", enabled)
	return true
}

func (scheduler *SchedulerPriority) UpdateParallelism(parallelism int) bool {
	scheduler.parallelism = parallelism
	log.Info("parallelism is updated to", parallelism)
	return true
}

func (scheduler *SchedulerPriority) updateGroup(group Group) bool {
	return true
}

func (scheduler *SchedulerPriority) DebugDump() map[string]interface{} {
	res := map[string]interface{}{
		"queue":          scheduler.queue,
		"schedulingJobs": scheduler.schedulingJobs,
		"history":        scheduler.history}
	return res
}
