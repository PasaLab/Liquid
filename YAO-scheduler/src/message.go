package main

type MsgAgentReport struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

type MsgSubmit struct {
	Code    int    `json:"code"`
	Error   string `json:"error"`
	JobName string `json:"jobName"`
}

type MsgPoolStatusHistory struct {
	Code  int          `json:"code"`
	Error string       `json:"error"`
	Data  []PoolStatus `json:"data"`
}

type MsgStop struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

type MsgSummary struct {
	Code         int    `json:"code"`
	Error        string `json:"error"`
	JobsFinished int    `json:"jobs_finished"`
	JobsRunning  int    `json:"jobs_running"`
	JobsPending  int    `json:"jobs_pending"`
	FreeGPU      int    `json:"gpu_free"`
	UsingGPU     int    `json:"gpu_using"`
}

type MsgJobList struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
	Jobs  []Job  `json:"jobs"`
}

type MsgLog struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
	Logs  string `json:"logs"`
}

type MsgTaskStatus struct {
	Code   int        `json:"code"`
	Error  string     `json:"error"`
	Status TaskStatus `json:"status"`
}

type MsgJobStatus struct {
	Code   int          `json:"code"`
	Error  string       `json:"error"`
	Status []TaskStatus `json:"status"`
}

type MsgCreate struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
	Id    string `json:"id"`
}

type MsgResource struct {
	Code     int                   `json:"code"`
	Error    string                `json:"error"`
	Resource map[string]NodeStatus `json:"resources"`
}

type MsgGroupCreate struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

type MsgGroupList struct {
	Code   int     `json:"code"`
	Error  string  `json:"error"`
	Groups []Group `json:"groups"`
}

type MsgOptimizerPredict struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
	Total int    `json:"total"`
	Pre   int    `json:"pre"`
	Main  int    `json:"main"`
	Post  int    `json:"post"`
}

type MsgJobReq struct {
	Code    int    `json:"code"`
	Error   string `json:"error"`
	CPU     int    `json:"cpu"`
	Mem     int    `json:"mem"`
	UtilGPU int    `json:"gpu_util"`
	MemGPU  int    `json:"gpu_mem"`
	BW      int    `json:"bw"`
	Version int64  `json:"version"`
}

type MsgJobReqPredict struct {
	Code   int                `json:"code"`
	Error  string             `json:"error"`
	Labels map[string]float64 `json:"labels"`
}

type MsgConfUpdate struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

type MsgConfList struct {
	Code    int                    `json:"code"`
	Error   string                 `json:"error"`
	Options map[string]interface{} `json:"options"`
}
