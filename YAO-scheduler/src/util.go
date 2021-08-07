package main

import (
	"strconv"
)

type Job struct {
	ID           int         `json:"id"`
	Name         string      `json:"name"`
	Tasks        []Task      `json:"tasks"`
	Workspace    string      `json:"workspace"`
	Group        string      `json:"group"`
	BasePriority float64     `json:"base_priority"`
	Priority     JobPriority `json:"priority"`
	RunBefore    int         `json:"run_before"`
	CreatedAt    int         `json:"created_at"`
	StartedAt    int64       `json:"started_at"`
	UpdatedAt    int         `json:"updated_at"`
	CreatedBy    int         `json:"created_by"`
	Locality     int         `json:"locality"`
	Status       State       `json:"status"`
	NumberGPU    int         `json:"number_GPU"`
	Retries      int         `json:"retries"`
}

type Task struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Job       string `json:"job_name"`
	Image     string `json:"image"`
	Cmd       string `json:"cmd"`
	NumberCPU int    `json:"cpu_number"`
	Memory    int    `json:"memory"`
	NumberGPU int    `json:"gpu_number"`
	MemoryGPU int    `json:"gpu_memory"`
	BW        int    `json:"bw"`
	IsPS      bool   `json:"is_ps"`
	ModelGPU  string `json:"gpu_model"`
}

type UtilGPUTimeSeries struct {
	Time int64 `json:"time"`
	Util int   `json:"util"`
}

type OptimizerJobExecutionTime struct {
	Pre     int   `json:"pre"`
	Post    int   `json:"post"`
	Total   int   `json:"total"`
	Main    int   `json:"main"`
	Version int64 `json:"version"`
}

type OptimizerUtilGPU struct {
	Util    int `json:"util"`
	Version int `json:"version"`
}

type ResourceCount struct {
	NumberGPU int
	MemoryGPU int
	CPU       int
	Memory    int
}

func str2int(s string, defaultValue int) int {
	i, err := strconv.Atoi(s)
	if err == nil {
		return i
	}
	return defaultValue
}
