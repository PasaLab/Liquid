package main

import (
	"sync"
	"strings"
	"math/rand"
)

type Mocker struct {
	mu sync.Mutex
}

var MockerInstance *Mocker
var MockerInstanceLock sync.Mutex

func InstanceOfMocker() *Mocker {
	MockerInstanceLock.Lock()
	defer MockerInstanceLock.Unlock()

	if MockerInstance == nil {
		MockerInstance = &Mocker{}
	}
	return MockerInstance
}

func (mocker *Mocker) GetDuration(job Job, nodes map[string]NodeStatus) int {
	str := strings.Split(job.Name, "-")
	duration := 300

	mode := "unknown"
	if len(job.Tasks) == 1 {
		if job.Tasks[0].NumberGPU == 1 {
			mode = "s1"
		} else if job.Tasks[0].NumberGPU == 2 {
			mode = "s2"
		}
	} else if len(job.Tasks) == 3 {
		var psNodes []string
		var workerNodes []string
		for _, task := range job.Tasks {
			if task.IsPS {
				psNodes = append(psNodes, nodes[task.Name].ClientHost)
			} else {
				workerNodes = append(workerNodes, nodes[task.Name].ClientHost)
			}
		}
		if psNodes[0] == workerNodes[0] {
			if psNodes[0] == workerNodes[1] {
				mode = "pww"
			} else {
				mode = "pw:w"
			}
		} else {
			if psNodes[0] == workerNodes[1] {
				mode = "pw:w"
			} else if workerNodes[0] == workerNodes[1] {
				mode = "p:ww"
			} else {
				mode = "p:w:w"
			}
		}

	}

	if len(str) > 1 {
		jobName := str[0]

		durations := map[string]map[string]int{
			"vgg16": {
				"s1":    220,
				"s2":    227,
				"pww":   510,
				"pw:w":  767,
				"p:ww":  1190,
				"p:w:w": 810,
			},
			"resnet50": {
				"s1":    146,
				"s2":    164,
				"pww":   203,
				"pw:w":  204,
				"p:ww":  255,
				"p:w:w": 210,
			},
			"inception3": {
				"s1":    253,
				"s2":    257,
				"pww":   289,
				"pw:w":  295,
				"p:ww":  310,
				"p:w:w": 290,
			},
		}

		if vals, ok := durations[jobName]; ok {
			if val, ok2 := vals[mode]; ok2 {
				return val * (100 + (rand.Intn(5) - 5)) / 100
			}
		}
	}
	return duration
}
