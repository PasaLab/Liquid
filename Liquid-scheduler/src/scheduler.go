package main

type Scheduler interface {
	Start()

	Schedule(Job)

	UpdateProgress(job Job, state State)

	AcquireResource(Job) []NodeStatus

	ReleaseResource(Job, NodeStatus)

	QueryState(jobName string) MsgJobStatus

	QueryLogs(jobName string, taskName string) MsgLog

	Stop(jobName string) MsgStop

	ListJobs() MsgJobList

	Summary() MsgSummary

	SetEnabled(enabled bool) bool

	UpdateParallelism(parallelism int) bool

	/* TODO: rearrange jobs to other queues */
	updateGroup(group Group) bool

	DebugDump() map[string]interface{}
}
