package main

type Evaluator struct {
	domains map[string]map[string]map[string]int
	racks   map[string]map[string]map[string]int
	nodes   map[string]map[string]map[string]int
	//upstreams   map[string]map[string]string
	totalPSs     map[string]int
	totalWorkers map[string]int
	totalPS      int
	totalWorker  int

	costNetwork float64
	costLoad    float64

	factorNode   float64
	factorRack   float64
	factorDomain float64

	factorSpread float64
}

func (eva *Evaluator) init(nodes []NodeStatus, tasks []Task) {
	eva.domains = map[string]map[string]map[string]int{}
	eva.racks = map[string]map[string]map[string]int{}
	eva.nodes = map[string]map[string]map[string]int{}
	//eva.upstreams = map[string]string{}
	eva.totalPSs = map[string]int{}
	eva.totalWorkers = map[string]int{}
	eva.totalPS = 0
	eva.totalWorker = 0
	eva.factorNode = 1.0
	eva.factorRack = 4.0
	eva.factorDomain = 40.0
	eva.costNetwork = 0.0
	eva.costLoad = 0.0
	eva.factorSpread = -1.0
}

func (eva *Evaluator) add(node NodeStatus, task Task) {
	if _, ok := eva.nodes[task.Job]; !ok {
		eva.nodes[task.Job] = map[string]map[string]int{}
		eva.racks[task.Job] = map[string]map[string]int{}
		eva.domains[task.Job] = map[string]map[string]int{}
		eva.totalPSs = map[string]int{}
		eva.totalWorkers = map[string]int{}
	}
	/* update network cost */
	if _, ok := eva.nodes[task.Job][node.ClientID]; !ok {
		eva.nodes[task.Job][node.ClientID] = map[string]int{"PS": 0, "Worker": 0}
	}
	if _, ok := eva.racks[task.Job][node.Rack]; !ok {
		eva.racks[task.Job][node.Rack] = map[string]int{"PS": 0, "Worker": 0}
	}
	if _, ok := eva.domains[task.Job][node.Domain]; !ok {
		eva.domains[task.Job][node.Domain] = map[string]int{"PS": 0, "Worker": 0}
	}
	bwFactor := float64(task.BW)
	if task.IsPS {
		eva.costNetwork += bwFactor * eva.factorNode * float64(eva.racks[task.Job][node.Rack]["Worker"]-eva.nodes[task.Job][node.ClientID]["Worker"])
		eva.costNetwork += bwFactor * eva.factorRack * float64(eva.domains[task.Job][node.Domain]["Worker"]-eva.racks[task.Job][node.Rack]["Worker"])
		eva.costNetwork += bwFactor * eva.factorDomain * float64(eva.totalWorkers[task.Job]-eva.domains[task.Job][node.Domain]["Worker"])

		eva.nodes[task.Job][node.ClientID]["PS"]++
		eva.racks[task.Job][node.Rack]["PS"]++
		eva.domains[task.Job][node.Domain]["PS"]++
		eva.totalPSs[task.Job]++
		eva.totalPS++
	} else {
		eva.costNetwork += bwFactor * eva.factorNode * float64(eva.racks[task.Job][node.Rack]["PS"]-eva.nodes[task.Job][node.ClientID]["PS"])
		eva.costNetwork += bwFactor * eva.factorRack * float64(eva.domains[task.Job][node.Domain]["PS"]-eva.racks[task.Job][node.Rack]["PS"])
		eva.costNetwork += bwFactor * eva.factorDomain * float64(eva.totalPSs[task.Job]-eva.domains[task.Job][node.Domain]["PS"])

		eva.nodes[task.Job][node.ClientID]["Worker"]++
		eva.racks[task.Job][node.Rack]["Worker"]++
		eva.domains[task.Job][node.Domain]["Worker"]++
		eva.totalWorkers[task.Job]++
		eva.totalWorker++
	}

	/* update node load cost */
	numberGPU := 1
	for _, gpu := range node.Status {
		if gpu.MemoryAllocated != 0 {
			numberGPU += 1
		}
	}
	eva.costLoad += float64(numberGPU+task.NumberGPU) / float64(len(node.Status))

}

func (eva *Evaluator) remove(node NodeStatus, task Task) {
	bwFactor := float64(task.BW)
	/* update network cost */
	if task.IsPS {
		eva.costNetwork -= bwFactor * eva.factorNode * float64(eva.racks[task.Job][node.Rack]["Worker"]-eva.nodes[task.Job][node.ClientID]["Worker"])
		eva.costNetwork -= bwFactor * eva.factorRack * float64(eva.domains[task.Job][node.Domain]["Worker"]-eva.racks[task.Job][node.Rack]["Worker"])
		eva.costNetwork -= bwFactor * eva.factorDomain * float64(eva.totalWorkers[task.Job]-eva.domains[task.Job][node.Domain]["Worker"])

		eva.nodes[task.Job][node.ClientID]["PS"]--
		eva.racks[task.Job][node.Rack]["PS"]--
		eva.domains[task.Job][node.Domain]["PS"]--
		eva.totalPSs[task.Job]--
		eva.totalPS--
	} else {
		eva.costNetwork -= bwFactor * eva.factorNode * float64(eva.racks[task.Job][node.Rack]["PS"]-eva.nodes[task.Job][node.ClientID]["PS"])
		eva.costNetwork -= bwFactor * eva.factorRack * float64(eva.domains[task.Job][node.Domain]["PS"]-eva.racks[task.Job][node.Rack]["PS"])
		eva.costNetwork -= bwFactor * eva.factorDomain * float64(eva.totalPSs[task.Job]-eva.domains[task.Job][node.Domain]["PS"])

		eva.nodes[task.Job][node.ClientID]["Worker"]--
		eva.racks[task.Job][node.Rack]["Worker"]--
		eva.domains[task.Job][node.Domain]["Worker"]--
		eva.totalWorkers[task.Job]--
		eva.totalWorker--
	}

	/* update node load cost */
	numberGPU := 1
	for _, gpu := range node.Status {
		if gpu.MemoryAllocated != 0 {
			numberGPU += 1
		}
	}
	eva.costLoad -= float64(numberGPU+task.NumberGPU) / float64(len(node.Status))
}

func (eva *Evaluator) calculate() float64 {
	usingNodes := 0.0
	for _, job := range eva.nodes {
		for _, pair := range job {
			if v, ok := pair["PS"]; ok && v > 0 {
				usingNodes += 1.0
			} else if v, ok := pair["Worker"]; ok && v > 0 {
				usingNodes += 1.0
			}
		}
	}
	if eva.totalPS+eva.totalWorker == 0 {
		usingNodes = 1.0
	} else {
		usingNodes /= float64(eva.totalWorker + eva.totalPS)
	}
	return eva.costNetwork + eva.factorSpread*eva.costLoad/float64(eva.totalPS+eva.totalWorker) + usingNodes
}
