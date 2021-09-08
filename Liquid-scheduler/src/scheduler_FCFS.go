package main

import (
	"sync"
	"time"
)

type SchedulerFCFS struct {
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

func (scheduler *SchedulerFCFS) Start() {
	scheduler.jobMasters = map[string]*JobManager{}
	scheduler.history = []*Job{}
	scheduler.enabled = true
	scheduler.parallelism = 1
	scheduler.schedulingJobs = map[string]bool{}

	go func() {
		for {
			log.Debug("Scheduling")

			scheduler.schedulingMu.Lock()
			if len(scheduler.schedulingJobs) >= scheduler.parallelism {
				scheduler.schedulingMu.Unlock()
				time.Sleep(time.Millisecond * 100)
				continue
			}
			scheduler.schedulingMu.Unlock()

			scheduler.queueMu.Lock()
			if len(scheduler.queue) > 0 {

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
			}
			scheduler.queueMu.Unlock()
		}
	}()
}

func (scheduler *SchedulerFCFS) UpdateProgress(job Job, state State) {
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

func (scheduler *SchedulerFCFS) Schedule(job Job) {
	scheduler.queueMu.Lock()
	defer scheduler.queueMu.Unlock()

	scheduler.queue = append(scheduler.queue, job)
	job.Status = Created
}

func (scheduler *SchedulerFCFS) AcquireResource(job Job) []NodeStatus {
	res := InstanceOfResourcePool().acquireResource(job)
	return res
}

func (scheduler *SchedulerFCFS) ReleaseResource(job Job, agent NodeStatus) {
	InstanceOfResourcePool().releaseResource(job, agent)
}

func (scheduler *SchedulerFCFS) QueryState(jobName string) MsgJobStatus {
	jm, ok := scheduler.jobMasters[jobName]
	if !ok {
		return MsgJobStatus{Code: 1, Error: "Job not exist!"}
	}
	return jm.status()
}

func (scheduler *SchedulerFCFS) Stop(jobName string) MsgStop {
	jm, ok := scheduler.jobMasters[jobName]
	if !ok {
		return MsgStop{Code: 1, Error: "Job not exist!"}
	}
	return jm.stop()
}

func (scheduler *SchedulerFCFS) QueryLogs(jobName string, taskName string) MsgLog {
	jm, ok := scheduler.jobMasters[jobName]
	if !ok {
		return MsgLog{Code: 1, Error: "Job not exist!"}
	}
	return jm.logs(taskName)
}

func (scheduler *SchedulerFCFS) ListJobs() MsgJobList {
	var tmp []Job
	for _, job := range scheduler.history {
		tmp = append(tmp, *job)
	}
	tmp = append(tmp, scheduler.queue...)
	return MsgJobList{Code: 0, Jobs: tmp}
}

func (scheduler *SchedulerFCFS) Summary() MsgSummary {
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

func (scheduler *SchedulerFCFS) SetEnabled(enabled bool) bool {
	scheduler.enabled = enabled
	log.Info("scheduler is set to ", enabled)
	return true
}

func (scheduler *SchedulerFCFS) UpdateParallelism(parallelism int) bool {
	scheduler.parallelism = parallelism
	log.Info("parallelism is updated to", parallelism)
	return true
}

func (scheduler *SchedulerFCFS) updateGroup(group Group) bool {
	return true
}

func (scheduler *SchedulerFCFS) DebugDump() map[string]interface{} {
	res := map[string]interface{}{}
	return res
}
