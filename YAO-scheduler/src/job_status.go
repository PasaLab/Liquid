package main

type JobStatus struct {
	Name  string
	tasks map[string]TaskStatus
}

type TaskStatus struct {
	Id          string                 `json:"id"`
	HostName    string                 `json:"hostname"`
	Node        string                 `json:"node"`
	Image       string                 `json:"image"`
	ImageDigest string                 `json:"image_digest"`
	Command     string                 `json:"command"`
	CreatedAt   string                 `json:"created_at"`
	FinishedAt  string                 `json:"finished_at"`
	Status      string                 `json:"status"`
	State       map[string]interface{} `json:"state"`
	UtilCPU     float64                `json:"cpu"`
	Mem         float64                `json:"mem"`
	BwRX        float64                `json:"bw_rx"`
	BWTx        float64                `json:"bw_tx"`
	UtilGPU     int                    `json:"gpu_util"`
	UtilMemGPU  int                    `json:"gpu_mem_util"`
	MemGPU      int                    `json:"gpu_mem"`
	TimeStamp   int64                  `json:"timestamp"`
}
