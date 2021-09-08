package main

import (
	"strconv"
	"time"
	"testing"
)

func TgenerateCase() ([]NodeStatus, []Task) {
	numTask := 6

	var nodes []NodeStatus
	var tasks []Task

	for i := 0; i < numTask*3; i++ {
		node := NodeStatus{ClientID: strconv.Itoa(i), Rack: "Rack-" + strconv.Itoa(i%40), Domain: "Domain-" + strconv.Itoa(i%4)}
		node.NumCPU = 24
		node.UtilCPU = 2.0
		node.MemTotal = 188
		node.MemAvailable = 20
		node.TotalBW = 100
		cnt := 4
		//cnt := rand.Intn(3) + 1
		for i := 0; i < cnt; i++ {
			node.Status = append(node.Status, GPUStatus{MemoryTotal: 11439, MemoryAllocated: 0, UUID: node.ClientID + "-" + strconv.Itoa(i)})
		}
		nodes = append(nodes, node)
	}
	for i := 0; i < numTask; i++ {
		isPS := false
		if i < numTask/3 {
			isPS = true
		}
		task := Task{Name: "task-" + strconv.Itoa(i), IsPS: isPS}
		task.Memory = 4
		task.NumberCPU = 2
		task.NumberGPU = 1
		task.MemoryGPU = 4096
		tasks = append(tasks, task)
	}
	return nodes, tasks
}

func TestBestFit(t *testing.T) {
	return
	nodes, tasks := TgenerateCase()
	for _, node := range nodes {
		log.Info(node)
	}
	s := time.Now()
	allocation := InstanceOfAllocator().fastBestFit(nodes, tasks)
	log.Info(time.Since(s))
	log.Info(allocation)
}

func TestGA(t *testing.T) {

	nodes, tasks := TgenerateCase()

	allocation := InstanceOfAllocator().GA(nodes, tasks, true)

	log.Info(allocation.TasksOnNode)
	log.Info(allocation.Nodes)

	allocation = InstanceOfAllocator().fastBestFit(nodes, tasks)

	InstanceOfResourcePool().Start()
	allocatedNodes := InstanceOfResourcePool().acquireResource(Job{Tasks: tasks})
	log.Info(allocatedNodes)
}
