package main

import (
	"sync"
	"math"
	"time"
	"github.com/MaxHalford/eaopt"
	"math/rand"
)

var allocatorInstance *Allocator
var allocatorInstanceLock sync.Mutex

func InstanceOfAllocator() *Allocator {
	defer allocatorInstanceLock.Unlock()
	allocatorInstanceLock.Lock()

	if allocatorInstance == nil {
		allocatorInstance = &Allocator{}
	}
	return allocatorInstance
}

type Allocator struct {
	allocationStrategy string
}

func (allocator *Allocator) init(conf Configuration) {

}

func (allocator *Allocator) updateStrategy(strategy string) bool {
	validStrategies := map[string]bool{
		"bestfit": true,
		"ga":      true,
		"mixed":   true,
		"random":  true,
		"spread":  true,
		"pack":    true,
	}
	if _, ok := validStrategies[strategy]; ok {
		allocator.allocationStrategy = strategy
		log.Info("Allocator strategy switched to ", strategy)
		return true
	}
	log.Warn("Invalid Allocator strategy ", strategy)
	return false
}

func (allocator *Allocator) allocate(nodes []NodeStatus, tasks []Task) Allocation {
	//log.Info(nodes)
	//log.Info(tasks)
	var allocation Allocation
	switch allocator.allocationStrategy {
	case "random":
		allocation = allocator.randomFitStrategy(nodes, tasks)
		break
	case "spread":
		allocation = allocator.spreadStrategy(nodes, tasks, false)
		break
	case "pack":
		allocation = allocator.spreadStrategy(nodes, tasks, true)
		break
	case "bestfit":
		allocation = allocator.fastBestFit(nodes, tasks)
		break
	case "ga":
		if len(tasks) >= 3 {
			allocation = allocator.GA(nodes, tasks, false)
		} else {
			allocation = allocator.fastBestFit(nodes, tasks)
		}
		break
	case "mixed":
		if len(tasks) > 3 {
			allocation = allocator.GA(nodes, tasks, true)
		} else {
			allocation = allocator.fastBestFit(nodes, tasks)
		}
		break
	default:
		allocation = allocator.fastBestFit(nodes, tasks)
	}
	return allocation
}

func (allocator *Allocator) randomFitStrategy(nodes []NodeStatus, tasks []Task) Allocation {
	allocation := Allocation{TasksOnNode: map[string][]Task{}, Nodes: map[string]NodeStatus{}, Flags: map[string]bool{"valid": true}}
	var nodesT []NodeStatus
	for _, node := range nodes {
		/* copy in order not to modify original data */
		nodesT = append(nodesT, node.Copy())
	}
	for _, node := range nodesT {
		allocation.Nodes[node.ClientID] = node
	}
	for _, task := range tasks {
		allocation.Tasks = append(allocation.Tasks, task)
	}
	for _, node := range nodesT {
		allocation.Nodes[node.ClientID] = node
		allocation.NodeIDs = append(allocation.NodeIDs, node.ClientID)
	}
	/* random-fit */
	for _, task := range tasks {
		if nodeID, ok := randomFit(allocation, task); ok {
			allocation.TasksOnNode[nodeID] = append(allocation.TasksOnNode[nodeID], task)
			cnt := task.NumberGPU
			for i := range allocation.Nodes[nodeID].Status {
				if allocation.Nodes[nodeID].Status[i].MemoryAllocated == 0 {
					allocation.Nodes[nodeID].Status[i].MemoryAllocated += task.MemoryGPU
					cnt--
				}
				if cnt == 0 {
					break
				}
			}
		} else {
			allocation.Flags["valid"] = false
			break
		}
	}
	return allocation
}

func (allocator *Allocator) spreadStrategy(nodes []NodeStatus, tasks []Task, isPack bool) Allocation {
	allocation := Allocation{TasksOnNode: map[string][]Task{}, Nodes: map[string]NodeStatus{}, Flags: map[string]bool{"valid": true}}
	var nodesT []NodeStatus
	for _, node := range nodes {
		/* copy in order not to modify original data */
		nodesT = append(nodesT, node.Copy())
	}
	for _, node := range nodesT {
		allocation.Nodes[node.ClientID] = node
	}
	for _, task := range tasks {
		allocation.Tasks = append(allocation.Tasks, task)
	}
	for _, node := range nodesT {
		allocation.Nodes[node.ClientID] = node
		allocation.NodeIDs = append(allocation.NodeIDs, node.ClientID)
	}
	for _, task := range tasks {
		if nodeID, ok := allocator.pickNode(allocation, task, isPack); ok {
			allocation.TasksOnNode[nodeID] = append(allocation.TasksOnNode[nodeID], task)
			cnt := task.NumberGPU
			for i := range allocation.Nodes[nodeID].Status {
				if allocation.Nodes[nodeID].Status[i].MemoryAllocated == 0 {
					allocation.Nodes[nodeID].Status[i].MemoryAllocated += task.MemoryGPU
					cnt--
				}
				if cnt == 0 {
					break
				}
			}
		} else {
			allocation.Flags["valid"] = false
			break
		}
	}
	return allocation
}

func (allocator *Allocator) pickNode(allocation Allocation, task Task, isPack bool) (string, bool) {
	flag := false
	bestNodeID := ""
	bestUtil := math.MaxFloat64
	for _, nodeID := range allocation.NodeIDs {
		if _, ok := allocation.Nodes[nodeID]; !ok {
			continue
		}
		available := 0
		total := 0
		for _, gpu := range allocation.Nodes[nodeID].Status {
			if gpu.MemoryAllocated == 0 {
				available += 1
			}
			total += 1
		}
		util := float64(available) / float64(total)
		if task.NumberGPU <= available {
			if bestNodeID == "" || (!isPack && util > bestUtil) || (isPack && util < bestUtil) {
				flag = true
				bestNodeID = nodeID
				bestUtil = util
			}
		}
	}
	return bestNodeID, flag
}

func (allocator *Allocator) fastBestFit(nodes []NodeStatus, tasks []Task) Allocation {
	eva := Evaluator{}
	eva.init(nodes, tasks)

	allocation := Allocation{Flags: map[string]bool{"valid": true}}
	allocation.TasksOnNode = map[string][]Task{}
	for _, task := range tasks {
		minCost := math.MaxFloat64
		var best *NodeStatus
		if task.IsPS {
			eva.factorSpread = 1.0
		} else {
			eva.factorSpread = -1.0
		}
		for i, node := range nodes {
			if _, ok := allocation.TasksOnNode[node.ClientID]; !ok {
				allocation.TasksOnNode[node.ClientID] = []Task{}
			}
			available := 0
			for _, gpu := range node.Status {
				if gpu.MemoryAllocated == 0 {
					available += 1
				}
			}
			if task.NumberGPU > available {
				continue
			}
			eva.add(node, task)
			cost := eva.calculate()
			eva.remove(node, task)
			//log.Info(node, cost)
			if cost < minCost || best == nil {
				minCost = cost
				best = &nodes[i]
			}
		}
		if best == nil {
			allocation.Flags["valid"] = false
			break
		} else {
			//log.Info(task, " choose ", best.ClientID)
			//fmt.Println(task, nodeID, allocation.TasksOnNode, minCost)
			allocation.TasksOnNode[best.ClientID] = append(allocation.TasksOnNode[best.ClientID], task)
			eva.add(*best, task)
			cnt := 0
			for i := range best.Status {
				if best.Status[i].MemoryAllocated == 0 {
					best.Status[i].MemoryAllocated += task.MemoryGPU
					cnt++
				}
				if cnt >= task.NumberGPU {
					break
				}
			}
		}
	}
	log.Info("BestFit Cost:", eva.calculate())
	return allocation
}

func (allocator *Allocator) GA(nodes []NodeStatus, tasks []Task, useBestFit bool) Allocation {
	// Instantiate a GA with a GAConfig
	var ga, err = eaopt.NewDefaultGAConfig().NewGA()
	if err != nil {
		log.Warn(err)
		return Allocation{Flags: map[string]bool{"valid": false}}
	}

	// Set the number of generations to run for
	ga.NGenerations = math.MaxInt32
	ga.NPops = 1
	ga.PopSize = 30 + uint(len(tasks)/2)

	// Add a custom print function to track progress
	ga.Callback = func(ga *eaopt.GA) {
		log.Info("Best fitness at generation ", ga.Generations, ": ", ga.HallOfFame[0].Fitness)
	}

	/* remember best */
	best := Allocation{Flags: map[string]bool{"valid": false}}
	bestFitness := math.MaxFloat64
	count := 0

	ts := time.Now()
	ga.EarlyStop = func(ga *eaopt.GA) bool {
		improvement := -(ga.HallOfFame[0].Fitness - bestFitness)
		if improvement <= 0.000001 {
			if count >= 30+len(tasks) || time.Since(ts) > time.Second*30 {
				//log.Info("Early Stop")
				return true
			} else {
				count++
			}
		} else {
			bestFitness = ga.HallOfFame[0].Fitness
			count = 1
			best = ga.HallOfFame[0].Genome.(Allocation)
		}
		return false
	}

	var Factory = func(rng *rand.Rand) eaopt.Genome {
		allocation := Allocation{TasksOnNode: map[string][]Task{}, Nodes: map[string]NodeStatus{}, Flags: map[string]bool{"valid": true}}

		var nodesT []NodeStatus
		for _, node := range nodes {
			/* copy in order not to modify original data */
			nodesT = append(nodesT, node.Copy())
		}
		for _, node := range nodesT {
			allocation.Nodes[node.ClientID] = node
		}
		for _, task := range tasks {
			allocation.Tasks = append(allocation.Tasks, task)
		}

		/* shuffle */
		for n := len(tasks); n > 0; n-- {
			randIndex := rng.Intn(n)
			allocation.Tasks[n-1], allocation.Tasks[randIndex] = allocation.Tasks[randIndex], allocation.Tasks[n-1]
		}

		for _, node := range nodesT {
			allocation.Nodes[node.ClientID] = node
			allocation.NodeIDs = append(allocation.NodeIDs, node.ClientID)
		}

		t := rng.Int() % 10
		if t == 0 && useBestFit {
			/* best-fit */
			//ts := time.Now()
			allo := allocator.fastBestFit(nodesT, tasks)
			allocation.TasksOnNode = allo.TasksOnNode
			allocation.Flags = allo.Flags
			//log.Println(time.Since(ts))
			//fmt.Println("Best Fit")
		} else if t%2 == 0 {
			/* first-fit */
			for _, task := range tasks {
				if nodeID, ok := randomFit(allocation, task); ok {
					allocation.TasksOnNode[nodeID] = append(allocation.TasksOnNode[nodeID], task)
					cnt := task.NumberGPU
					for i := range allocation.Nodes[nodeID].Status {
						if allocation.Nodes[nodeID].Status[i].MemoryAllocated == 0 {
							allocation.Nodes[nodeID].Status[i].MemoryAllocated += task.MemoryGPU
							cnt--
						}
						if cnt == 0 {
							break
						}
					}
				} else {
					allocation.Flags["valid"] = false
					break
				}
			}
		} else {
			/* random-fit */
			for _, task := range tasks {
				if nodeID, ok := randomFit(allocation, task); ok {
					allocation.TasksOnNode[nodeID] = append(allocation.TasksOnNode[nodeID], task)
					cnt := task.NumberGPU
					for i := range allocation.Nodes[nodeID].Status {
						if allocation.Nodes[nodeID].Status[i].MemoryAllocated == 0 {
							allocation.Nodes[nodeID].Status[i].MemoryAllocated += task.MemoryGPU
							cnt--
						}
						if cnt == 0 {
							break
						}
					}
				} else {
					allocation.Flags["valid"] = false
					break
				}
			}
		}
		//fmt.Println(evaluate(allocation))
		//fmt.Println(allocation)

		cnt := 0
		for _, tasks := range allocation.TasksOnNode {
			for range tasks {
				cnt++
			}
		}
		if cnt != len(allocation.Tasks) && allocation.Flags["valid"] {
			log.Warn("factory:", cnt, len(allocation.Tasks), " t=", t)
		}

		return allocation
	}

	// Find the minimum
	err = ga.Minimize(Factory)
	log.Info("GA uses ", time.Since(ts))
	//log.Println(ga.HallOfFame[0].Genome.(Allocation).TasksOnNode)
	//log.Println(ga.HallOfFame[0].Genome.(Allocation).Flags)
	//log.Println(ga.HallOfFame[0].Genome.(Allocation).Nodes)
	if err != nil {
		log.Warn(err)
		return Allocation{Flags: map[string]bool{"valid": false}}
	}
	return best
}

func randomFit(allocation Allocation, task Task) (string, bool) {
	flag := false
	nodeID := ""
	for nodeID = range allocation.Nodes {
		available := 0
		for _, gpu := range allocation.Nodes[nodeID].Status {
			if gpu.MemoryAllocated == 0 {
				available += 1
			}
		}
		if task.NumberGPU <= available {
			flag = true
			break
		}
	}
	return nodeID, flag
}

func firstFit(allocation Allocation, task Task) (string, bool) {
	flag := false
	nodeID := ""
	for _, nodeID = range allocation.NodeIDs {
		if _, ok := allocation.Nodes[nodeID]; !ok {
			continue
		}
		available := 0
		for _, gpu := range allocation.Nodes[nodeID].Status {
			if gpu.MemoryAllocated == 0 {
				available += 1
			}
		}
		if task.NumberGPU <= available {
			flag = true
			break
		}
	}
	return nodeID, flag
}
