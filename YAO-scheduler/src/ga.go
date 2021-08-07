package main

import (
	"math/rand"
	"github.com/MaxHalford/eaopt"
	"math"
)

// A resource allocation
type Allocation struct {
	TasksOnNode map[string][]Task // tasks on nodes[id]
	Nodes       map[string]NodeStatus
	NodeIDs     []string
	Flags       map[string]bool
	Evaluator   Evaluator
	Tasks       []Task
}

/* Evaluate the allocation */
func (X Allocation) Evaluate() (float64, error) {
	if !X.Flags["valid"] {
		//fmt.Println("Invalid allocation")
		return math.MaxFloat64, nil
	}

	var nodes []NodeStatus
	for _, node := range X.Nodes {
		nodes = append(nodes, node)
	}

	eva := Evaluator{}
	eva.init(nodes, X.Tasks)
	for node, tasks := range X.TasksOnNode {
		for _, task := range tasks {
			eva.add(X.Nodes[node], task)
		}
	}

	cost := eva.calculate()
	//log.Info(cost)
	//return float64(cost) + float64(len(X.Nodes)), nil
	return float64(cost) + float64(len(X.Nodes))/float64(len(X.Tasks)), nil
	//return float64(cost), nil
}

// Mutate a Vector by resampling each element from a normal distribution with probability 0.8.
func (X Allocation) Mutate(rng *rand.Rand) {
	/* remove a node randomly */
	// make sure n > 0 && round >0
	round := rng.Intn(1+len(X.Nodes)/100)%50 + 1
	for i := 0; i < round; i++ {
		if !X.Flags["valid"] {
			return
		}
		//fmt.Println("Mutate")
		//fmt.Println("Before", X)

		var nodeIDs []string
		for nodeID := range X.Nodes {
			nodeIDs = append(nodeIDs, nodeID)
		}
		randIndex := rng.Intn(len(X.Nodes))
		nodeID := nodeIDs[randIndex]

		/* reschedule tasks on tgt node */
		var tasks []Task
		if _, ok := X.TasksOnNode[nodeID]; ok {
			for _, task := range X.TasksOnNode[nodeID] {
				tasks = append(tasks, task)
			}
			delete(X.TasksOnNode, nodeID)
		}
		//log.Info("Delete node ", nodeID)
		//log.Info("Before ", X.Nodes)
		delete(X.Nodes, nodeID)
		//log.Info("After ", X.Nodes)

		//fmt.Println(tasks)

		/* random-fit */
		for _, task := range tasks {
			if nodeID, ok := randomFit(X, task); ok {
				if len(X.TasksOnNode[nodeID]) == 0 {
					X.TasksOnNode[nodeID] = []Task{}
				}
				X.TasksOnNode[nodeID] = append(X.TasksOnNode[nodeID], task)
				cnt := task.NumberGPU
				//log.Info("Add task ", task.Name, " in ", nodeID)
				//log.Info("Before ", X.Nodes[nodeID].Status)
				for i := range X.Nodes[nodeID].Status {
					if X.Nodes[nodeID].Status[i].MemoryAllocated == 0 {
						X.Nodes[nodeID].Status[i].MemoryAllocated += task.MemoryGPU
						cnt--
						if cnt == 0 {
							break
						}
					}
				}
				if cnt != 0 {
					log.Warn("task ", task.Name, " still need ", cnt)
				}

				//log.Info("After ", X.Nodes[nodeID].Status)
			} else {
				X.Flags["valid"] = false
				break
			}
		}
	}

	return
	/* move tasks */
	if !X.Flags["valid"] {
		//fmt.Println("Invalid allocation")
		return
	}
	var nodeIDs []string
	for nodeID := range X.Nodes {
		nodeIDs = append(nodeIDs, nodeID)
	}
	randIndex1 := rng.Intn(len(nodeIDs))
	nodeID1 := nodeIDs[randIndex1]
	if tasks, ok := X.TasksOnNode[nodeID1]; ok && len(tasks) > 0 {
		idx := rng.Intn(len(tasks))
		task := tasks[idx]
		copy(X.TasksOnNode[nodeID1][idx:], X.TasksOnNode[nodeID1][idx+1:])
		X.TasksOnNode[nodeID1] = X.TasksOnNode[nodeID1][:len(X.TasksOnNode[nodeID1])-1]

		if nodeID, ok := firstFit(X, task); ok {
			X.TasksOnNode[nodeID] = append(X.TasksOnNode[nodeID], task)
		} else {
			X.Flags["valid"] = false
		}
	}

}

// Crossover a Vector with another Vector by applying uniform crossover.
func (X Allocation) Crossover(Y eaopt.Genome, rng *rand.Rand) {
	// make sure n > 0 && round > 0
	cnt := 0
	for _, tasks := range X.TasksOnNode {
		for range tasks {
			cnt++
		}
	}
	if cnt != len(X.Tasks) && X.Flags["valid"] {
		log.Warn("1:", cnt, len(X.Tasks))
	}
	round := rng.Intn(1+len(X.Nodes)/100)%10 + 1
	for i := 0; i < round; i++ {
		if !Y.(Allocation).Flags["valid"] || !X.Flags["valid"] {
			return
		}
		taskToNode := map[string]string{}
		for nodeID, tasks := range X.TasksOnNode {
			for _, task := range tasks {
				taskToNode[task.ID] = nodeID
			}
		}

		var nodeIDs []string
		for nodeID := range Y.(Allocation).Nodes {
			nodeIDs = append(nodeIDs, nodeID)
		}

		//fmt.Println(nodeIDs, Y.(Allocation))
		randIndex := rng.Intn(len(nodeIDs))
		nodeID := nodeIDs[randIndex]

		/* remove duplicated tasks */
		for _, task := range Y.(Allocation).TasksOnNode[nodeID] {
			//fmt.Println(Y.(Allocation).TasksOnNode[nodeID])
			idx := -1
			nodeID2, ok := taskToNode[task.ID]
			if !ok {
				log.Warn("Error", taskToNode, X.TasksOnNode, task.ID)
			}
			for i, task2 := range X.TasksOnNode[nodeID2] {
				if task2.ID == task.ID {
					idx = i
				}
			}
			if idx == -1 {
				log.Warn("Error 2", taskToNode, X.TasksOnNode, task.ID)
			}
			//fmt.Println(X.TasksOnNode)
			copy(X.TasksOnNode[nodeID2][idx:], X.TasksOnNode[nodeID2][idx+1:])
			X.TasksOnNode[nodeID2] = X.TasksOnNode[nodeID2][:len(X.TasksOnNode[nodeID2])-1]
			cnt := task.NumberGPU
			//log.Info("Remove task ", task.Name, " in ", nodeID2)
			//log.Info("Before ", X.Nodes[nodeID2].Status)
			for i := range X.Nodes[nodeID2].Status {
				/* TODO: determine correct GPU */
				if X.Nodes[nodeID2].Status[i].MemoryAllocated == task.MemoryGPU {
					X.Nodes[nodeID2].Status[i].MemoryAllocated -= task.MemoryGPU
					cnt--
					if cnt == 0 {
						break
					}
				}
			}
			if cnt != 0 {
				log.Warn("cross add, no enough GPU left on ", nodeID, ", still need", cnt)
			}
			//log.Info("After ", X.Nodes[nodeID].Status)
			//fmt.Println(X.TasksOnNode)
		}

		/* reschedule tasks on tgt node */
		var tasks []Task
		if _, ok := X.TasksOnNode[nodeID]; ok {
			for _, task := range X.TasksOnNode[nodeID] {
				tasks = append(tasks, task)
			}
			delete(X.TasksOnNode, nodeID)
		}

		if _, ok := X.Nodes[nodeID]; ok {
			delete(X.Nodes, nodeID)
		}
		X.Nodes[nodeID] = Y.(Allocation).Nodes[nodeID].Copy()

		var newTasksOnNode []Task
		for _, task := range Y.(Allocation).TasksOnNode[nodeID] {
			newTasksOnNode = append(newTasksOnNode, task)
		}
		X.TasksOnNode[nodeID] = newTasksOnNode

		/* random-fit */
		for _, task := range tasks {
			if nodeID, ok := randomFit(X, task); ok {
				if len(X.TasksOnNode[nodeID]) == 0 {
					X.TasksOnNode[nodeID] = []Task{}
				}
				X.TasksOnNode[nodeID] = append(X.TasksOnNode[nodeID], task)
				cnt := task.NumberGPU
				//log.Info("Remove task ", task.Name, " in ", nodeID)
				//log.Info("Before ", X.Nodes[nodeID].Status)
				for i := range X.Nodes[nodeID].Status {
					if X.Nodes[nodeID].Status[i].MemoryAllocated == 0 {
						X.Nodes[nodeID].Status[i].MemoryAllocated += task.MemoryGPU
						cnt--
						if cnt == 0 {
							break
						}
					}
				}
				//log.Info("After ", X.Nodes[nodeID].Status)
				if cnt != 0 {
					log.Warn("cross add, no enough GPU left on ", nodeID, ", still need", cnt)
				}
			} else {
				X.Flags["valid"] = false
				break
			}
		}
	}
	cnt = 0
	for _, tasks := range X.TasksOnNode {
		for range tasks {
			cnt++
		}
	}
	if cnt != len(X.Tasks) && X.Flags["valid"] {
		log.Warn("2:", cnt, len(X.Tasks))
	}
	//fmt.Println()
	//fmt.Println("crossover", X.TasksOnNode)
}

// Clone a Vector to produce a new one that points to a different slice.
func (X Allocation) Clone() eaopt.Genome {
	if !X.Flags["valid"] {
		//fmt.Println(X.Valid)
	}
	Y := Allocation{TasksOnNode: map[string][]Task{}, Nodes: map[string]NodeStatus{}, Flags: map[string]bool{"valid": X.Flags["valid"]}}
	for id, node := range X.Nodes {
		Y.Nodes[id] = node.Copy()
		Y.NodeIDs = append(Y.NodeIDs, node.ClientID)
	}
	for id, tasks := range X.TasksOnNode {
		var t []Task
		for _, task := range tasks {
			t = append(t, task)
		}
		Y.TasksOnNode[id] = t
	}
	Y.Tasks = X.Tasks
	return Y
}
