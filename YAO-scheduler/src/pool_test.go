package main

import (
	"testing"
	"strconv"
	"time"
)

func TestPool(t *testing.T) {
	return
	InstanceOfResourcePool().Start()

	for j := 0; j < 100; j++ {
		for i := 0; i < 1000; i++ {
			node := NodeStatus{ClientID: strconv.Itoa(i)}
			InstanceOfResourcePool().update(node)
		}
	}

	count := 0
	for _, seg := range InstanceOfResourcePool().pools {
		log.Info(seg.ID, "<--->", len(seg.Nodes), " ", seg.Nodes == nil, " Next:", seg.Next.ID)
		count += len(seg.Nodes)
	}
	log.Info(count)

	counter := map[int]int{}
	for i := 0; i < 1000; i++ {
		seg := InstanceOfResourcePool().getNodePool(strconv.Itoa(i))
		counter[seg]++
	}
	//for k, v := range counter {
	//fmt.Println(k, "-->",v)
	//}

}

func TestAllocate(t *testing.T) {
	InstanceOfResourcePool().Start()

	job := Job{Name: strconv.Itoa(int(time.Now().Unix() % 1000000000))}
	job.Group = "default"
	var tasks []Task
	task := Task{}
	task.Name = "node1"
	task.NumberGPU = 1
	task.NumberCPU = 2
	task.Memory = 4096
	task.IsPS = false
	task.MemoryGPU = 4096
	tasks = append(tasks, task)
	job.Tasks = tasks

	allocation := InstanceOfResourcePool().acquireResource(job)
	log.Info(allocation)
}
