package main

type Group struct {
	Name      string `json:"name"`
	Weight    int    `json:"weight"`
	Reserved  bool   `json:"reserved"`
	NumGPU    int    `json:"quota_gpu"`
	MemoryGPU int    `json:"quota_gpu_mem"`
	CPU       int    `json:"quota_cpu"`
	Memory    int    `json:"quota_mem"`
}
