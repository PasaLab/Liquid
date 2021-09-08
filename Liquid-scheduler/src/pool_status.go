package main

type PoolStatus struct {
	TimeStamp       string  `json:"ts"`
	UtilCPU         float64 `json:"cpu_util"`
	TotalCPU        int     `json:"cpu_total"`
	TotalMem        int     `json:"mem_total"`
	AvailableMem    int     `json:"mem_available"`
	TotalGPU        int     `json:"TotalGPU"`
	UtilGPU         int     `json:"gpu_util"`
	TotalMemGPU     int     `json:"gpu_mem_total"`
	AvailableMemGPU int     `json:"gpu_mem_available"`
}

type GPUStatus struct {
	UUID             string `json:"uuid"`
	ProductName      string `json:"product_name"`
	PerformanceState string `json:"performance_state"`
	MemoryTotal      int    `json:"memory_total"`
	MemoryFree       int    `json:"memory_free"`
	MemoryAllocated  int    `json:"memory_allocated"`
	MemoryUsed       int    `json:"memory_used"`
	UtilizationGPU   int    `json:"utilization_gpu"`
	UtilizationMem   int    `json:"utilization_mem"`
	TemperatureGPU   int    `json:"temperature_gpu"`
	PowerDraw        int    `json:"power_draw"`
}

type NodeStatus struct {
	ClientID     string      `json:"id"`
	ClientHost   string      `json:"host"`
	Domain       string      `json:"domain"`
	Rack         string      `json:"rack"`
	Version      float64     `json:"version"`
	NumCPU       int         `json:"cpu_num"`
	UtilCPU      float64     `json:"cpu_load"`
	MemTotal     int         `json:"mem_total"`
	MemAvailable int         `json:"mem_available"`
	UsingBW      float64     `json:"bw_using"`
	TotalBW      float64     `json:"bw_total"`
	Status       []GPUStatus `json:"status"`
}

func (X NodeStatus) Copy() NodeStatus {
	res := X
	res.Status = make([]GPUStatus, len(X.Status))
	copy(res.Status, X.Status)
	return res
}
