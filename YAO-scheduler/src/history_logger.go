package main

import (
	"sync"
)

type JobHistoryLogger struct {
	scheduler  Scheduler
	jobs       map[string]Job
	tasks      map[string][]TaskStatus
	jobStatus  JobStatus
	resources  []NodeStatus
	killedFlag bool
	mu         sync.Mutex
}

var jobHistoryLoggerInstance *JobHistoryLogger
var jobHistoryLoggerInstanceLock sync.Mutex

func InstanceJobHistoryLogger() *JobHistoryLogger {
	defer jobHistoryLoggerInstanceLock.Unlock()
	jobHistoryLoggerInstanceLock.Lock()

	if jobHistoryLoggerInstance == nil {
		jobHistoryLoggerInstance = &JobHistoryLogger{}
	}
	return jobHistoryLoggerInstance
}

func (jhl *JobHistoryLogger) Start() {
	log.Info("jhl init")
	jhl.jobs = map[string]Job{}
	jhl.tasks = map[string][]TaskStatus{}
	/* retrieve list */
}

func (jhl *JobHistoryLogger) submitJob(job Job) {
	jhl.mu.Lock()
	defer jhl.mu.Unlock()
	log.Debug("submit job", job.Name)
	jhl.jobs[job.Name] = job
	jhl.tasks[job.Name] = []TaskStatus{}
}

func (jhl *JobHistoryLogger) updateJobStatus(jobName string, state State) {
	jhl.mu.Lock()
	defer jhl.mu.Unlock()
	log.Debug("update job status", jobName)
	if job, ok := jhl.jobs[jobName]; ok {
		job.Status = state
	}
}

func (jhl *JobHistoryLogger) submitTaskStatus(jobName string, task TaskStatus) {
	jhl.mu.Lock()
	defer jhl.mu.Unlock()
	log.Debug("submit job task status", jobName)
	if tasks, ok := jhl.tasks[jobName]; ok {
		jhl.tasks[jobName] = append(tasks, task)
	}
}

func (jhl *JobHistoryLogger) getTaskStatus(jobName string) []TaskStatus {
	jhl.mu.Lock()
	defer jhl.mu.Unlock()
	if _, ok := jhl.tasks[jobName]; ok {
		return jhl.tasks[jobName]
	}
	return []TaskStatus{}
}
