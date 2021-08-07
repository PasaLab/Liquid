package main

import (
	"sync"
	"strings"
	"io/ioutil"
	"strconv"
	"encoding/json"
	"math"
	"hash/fnv"
	"time"
)

type Optimizer struct {
	versions map[string]int
	reqCache map[string]MsgJobReq
	cache    map[string]OptimizerJobExecutionTime
	cacheMu  sync.Mutex
}

var optimizerInstance *Optimizer
var OptimizerInstanceLock sync.Mutex

func InstanceOfOptimizer() *Optimizer {
	defer OptimizerInstanceLock.Unlock()
	OptimizerInstanceLock.Lock()

	if optimizerInstance == nil {
		optimizerInstance = &Optimizer{}
		optimizerInstance.versions = map[string]int{}
		optimizerInstance.cache = map[string]OptimizerJobExecutionTime{}
		optimizerInstance.reqCache = map[string]MsgJobReq{}

		go func() {
			/* remove expired cache */
			for {
				time.Sleep(time.Second * 30)
				optimizerInstance.cacheMu.Lock()
				var expired []string
				for k, v := range optimizerInstance.cache {
					if time.Now().Unix()-v.Version > 300 {
						expired = append(expired, k)
					}
				}
				for _, k := range expired {
					delete(optimizerInstance.cache, k)
				}
				expired = []string{}
				for k, v := range optimizerInstance.reqCache {
					if time.Now().Unix()-v.Version > 300 {
						expired = append(expired, k)
					}
				}
				for _, k := range expired {
					delete(optimizerInstance.reqCache, k)
				}
				optimizerInstance.cacheMu.Unlock()
			}
		}()
	}
	return optimizerInstance
}

func (optimizer *Optimizer) Start() {
	log.Info("optimizer started")
}

func (optimizer *Optimizer) FeedTime(job Job, stats [][]TaskStatus) {
	//log.Info("optimizer feedTime", job)
	if len(stats) == 0 || len(job.Tasks) != 1 {
		return
	}

	go func() {
		str := strings.Split(job.Name, "-")
		if len(str) == 2 {
			jobName := str[0]

			var UtilGPUs []UtilGPUTimeSeries
			for _, stat := range stats {
				for _, task := range stat {
					UtilGPUs = append(UtilGPUs, UtilGPUTimeSeries{Time: task.TimeStamp, Util: task.UtilGPU})
				}
			}
			var preTime int64
			for i := 0; i < len(UtilGPUs); i++ {
				if UtilGPUs[i].Util > 15 {
					preTime = UtilGPUs[i].Time - UtilGPUs[0].Time
					break
				}
			}

			var postTime int64
			for i := len(UtilGPUs) - 1; i >= 0; i-- {
				if UtilGPUs[i].Util > 15 {
					postTime = UtilGPUs[len(UtilGPUs)-1].Time - UtilGPUs[i].Time
					break
				}
			}
			totalTime := UtilGPUs[len(UtilGPUs)-1].Time - UtilGPUs[0].Time
			if preTime+postTime >= totalTime { /* in case GPU is not used */
				preTime /= 2
				postTime /= 2
			}

			tmp := map[string]float64{
				"pre":   float64(preTime),
				"post":  float64(postTime),
				"total": float64(totalTime),
			}
			labels, _ := json.Marshal(tmp)

			cmd := job.Tasks[0].Cmd
			params := map[string]int{}

			exceptions := map[string]bool{}
			exceptions["train_dir"] = true
			exceptions["data__dir"] = true
			exceptions["tmp__dir"] = true
			exceptions["variable_update"] = true
			exceptions["ps_hosts"] = true
			exceptions["worker_hosts"] = true
			exceptions["task_index"] = true
			exceptions["job_name"] = true

			pairs := strings.Split(cmd, " ")
			for _, pair := range pairs {
				v := strings.Split(pair, "=")
				if len(v) == 2 && v[0][:2] == "--" {
					var param string
					var value int
					param = v[0][2:]

					if val, err := strconv.Atoi(v[1]); err == nil {
						value = val
					} else {
						h := fnv.New32a()
						h.Write([]byte(v[1]))
						value = int(h.Sum32())
					}
					if _, ok := exceptions[param]; !ok {
						params[param] = value
					}
				}
			}
			params["cpu"] = job.Tasks[0].NumberCPU
			params["mem"] = job.Tasks[0].Memory
			params["gpu"] = job.Tasks[0].NumberGPU
			params["gpu_mem"] = job.Tasks[0].MemoryGPU
			//log.Info(job.Name, params)
			features, _ := json.Marshal(params)

			spider := Spider{}
			spider.Method = "GET"
			spider.URL = "http://yao-optimizer:8080/feed?job=" + jobName + "-time" + "&features=" + string(features) + "&labels=" + string(labels)

			err := spider.do()
			if err != nil {
				log.Warn(err)
				return
			}

			resp := spider.getResponse()
			if _, err := ioutil.ReadAll(resp.Body); err != nil {
				log.Warn(err)
			}
			resp.Body.Close()
			if err != nil {
				log.Warn(err)
				return
			}

			if optimizer.versions[jobName]%3 == 0 {
				optimizer.trainTime(jobName)
			}
		}
	}()
}

func (optimizer *Optimizer) trainTime(jobName string) {
	spider := Spider{}
	spider.Method = "GET"
	params := "job=" + jobName + "-time"
	spider.URL = "http://yao-optimizer:8080/train?" + params

	err := spider.do()
	if err != nil {
		return
	}

	resp := spider.getResponse()
	if _, err := ioutil.ReadAll(resp.Body); err != nil {
		log.Warn(err)
	}
	resp.Body.Close()
	if err != nil {
		return
	}
}

func (optimizer *Optimizer) PredictTime(job Job) OptimizerJobExecutionTime {
	optimizer.cacheMu.Lock()
	if val, ok := optimizer.cache[job.Name]; ok {
		optimizer.cacheMu.Unlock()
		return val
	}
	optimizer.cacheMu.Unlock()

	res := OptimizerJobExecutionTime{Pre: 0, Post: 0, Total: math.MaxInt64}
	var jobName string
	str := strings.Split(job.Name, "-")
	if len(str) == 2 {
		jobName = str[0]
	} else if len(str) == 1 {
		jobName = job.Name
	} else {
		return res
	}
	cmd := job.Tasks[0].Cmd
	params := map[string]int{}

	exceptions := map[string]bool{}
	exceptions["train_dir"] = true
	exceptions["variable_update"] = true
	exceptions["ps_hosts"] = true
	exceptions["worker_hosts"] = true
	exceptions["task_index"] = true
	exceptions["job_name"] = true

	pairs := strings.Split(cmd, " ")
	for _, pair := range pairs {
		v := strings.Split(pair, "=")
		if len(v) == 2 && v[0][:2] == "--" {
			var param string
			var value int
			param = v[0][2:]

			if val, err := strconv.Atoi(v[1]); err == nil {
				value = val
			} else {
				h := fnv.New32a()
				h.Write([]byte(v[1]))
				value = int(h.Sum32())
			}
			if _, ok := exceptions[param]; !ok {
				params[param] = value
			}
		}
	}
	params["cpu"] = job.Tasks[0].NumberCPU
	params["mem"] = job.Tasks[0].Memory
	params["gpu"] = job.Tasks[0].NumberGPU
	params["gpu_mem"] = job.Tasks[0].MemoryGPU
	//log.Info(job.Name, params)

	features, _ := json.Marshal(params)

	spider := Spider{}
	spider.Method = "GET"
	spider.URL = "http://yao-optimizer:8080/predict?job=" + jobName + "-time" + "&features=" + string(features)

	err := spider.do()
	if err != nil {
		return res
	}

	resp := spider.getResponse()
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Warn(err)
		return res
	}

	var msg MsgJobReqPredict
	err = json.Unmarshal([]byte(string(body)), &msg)
	if err == nil && msg.Code == 0 {
		tmp := msg.Labels
		if v, ok := tmp["pre"]; ok {
			res.Pre = int(math.Ceil(v))
		}
		if v, ok := tmp["post"]; ok {
			res.Post = int(math.Ceil(v))
		}
		if v, ok := tmp["total"]; ok {
			res.Total = int(math.Ceil(v))
		}
	}
	res.Version = time.Now().Unix()
	optimizer.cacheMu.Lock()
	optimizer.cache[job.Name] = res
	optimizer.cacheMu.Unlock()
	return res
}

func (optimizer *Optimizer) FeedStats(job Job, role string, stats [][]TaskStatus) {
	if len(stats) == 0 {
		return
	}
	str := strings.Split(job.Name, "-")
	if len(str) == 1 {
		return
	}
	jobName := str[0]
	go func() {

		var UtilsCPU []float64
		var Mems []float64
		var BwRxs []float64
		var BwTxs []float64
		var UtilGPUs []float64
		var MemGPUs []float64
		for _, stat := range stats {
			for _, task := range stat {
				UtilsCPU = append(UtilsCPU, task.UtilCPU)
				Mems = append(Mems, task.Mem)
				BwRxs = append(BwRxs, task.BwRX)
				BwTxs = append(BwTxs, task.BWTx)
				UtilGPUs = append(UtilGPUs, float64(task.UtilGPU))
				MemGPUs = append(MemGPUs, float64(task.MemGPU))
			}
		}
		tmp := map[string]float64{
			"cpu":          optimizer.max(UtilsCPU),
			"cpu_std":      optimizer.std(UtilsCPU),
			"cpu_mean":     optimizer.mean(UtilsCPU),
			"mem":          optimizer.max(Mems),
			"bw_rx":        optimizer.mean(BwRxs),
			"bw_tx":        optimizer.mean(BwTxs),
			"gpu_util":     optimizer.mean(UtilGPUs),
			"gpu_util_std": optimizer.std(UtilGPUs),
			"gpu_mem":      optimizer.max(MemGPUs),
		}
		for k, v := range tmp {
			tmp[k] = float64(int(v))
		}
		labels, _ := json.Marshal(tmp)

		cmd := job.Tasks[0].Cmd
		params := map[string]int{}

		psNumber := 0
		workerNumber := 0
		for _, task := range job.Tasks {
			if (role == "PS" && task.IsPS) || (role == "Worker" && !task.IsPS) {
				params["num_gpus"] = task.NumberGPU
				cmd = task.Cmd
			}
			if task.IsPS {
				psNumber++
			} else {
				workerNumber++
			}
		}
		params["ps_number"] = psNumber
		params["worker_number"] = workerNumber
		if role == "PS" {
			params["role"] = 1
		} else {
			params["role"] = 0
		}

		exceptions := map[string]bool{}
		exceptions["train_dir"] = true
		exceptions["variable_update"] = true
		exceptions["ps_hosts"] = true
		exceptions["worker_hosts"] = true
		exceptions["task_index"] = true
		exceptions["job_name"] = true

		pairs := strings.Split(cmd, " ")
		for _, pair := range pairs {
			v := strings.Split(pair, "=")
			if len(v) == 2 && v[0][:2] == "--" {
				var param string
				var value int
				param = v[0][2:]

				if val, err := strconv.Atoi(v[1]); err == nil {
					value = val
				} else {
					h := fnv.New32a()
					h.Write([]byte(v[1]))
					value = int(h.Sum32())
				}
				if _, ok := exceptions[param]; !ok {
					params[param] = value
				}
			}
		}
		//log.Info(job.Name, params)

		features, _ := json.Marshal(params)

		spider := Spider{}
		spider.Method = "GET"
		spider.URL = "http://yao-optimizer:8080/feed?job=" + jobName + "&features=" + string(features) + "&labels=" + string(labels)

		err := spider.do()
		if err != nil {
			log.Warn(err)
			return
		}

		resp := spider.getResponse()
		if _, err := ioutil.ReadAll(resp.Body); err != nil {
			log.Warn(err)
		}
		resp.Body.Close()
		if err != nil {
			log.Warn(err)
			return
		}

		optimizer.versions[jobName]++
		if optimizer.versions[jobName]%3 == 0 {
			optimizer.trainReq(jobName)
		}
	}()
}

func (optimizer *Optimizer) trainReq(jobName string) {
	spider := Spider{}
	spider.Method = "GET"
	params := "job=" + jobName
	spider.URL = "http://yao-optimizer:8080/train?" + params

	err := spider.do()
	if err != nil {
		return
	}

	resp := spider.getResponse()
	if _, err := ioutil.ReadAll(resp.Body); err != nil {
		log.Warn(err)
	}
	resp.Body.Close()
	if err != nil {
		return
	}
}

func (optimizer *Optimizer) PredictReq(job Job, role string) MsgJobReq {
	optimizer.cacheMu.Lock()
	if val, ok := optimizer.reqCache[job.Name]; ok {
		optimizer.cacheMu.Unlock()
		return val
	}
	optimizer.cacheMu.Unlock()
	res := MsgJobReq{CPU: 4, Mem: 4096, UtilGPU: 100, MemGPU: 8192, BW: 1}

	var jobName string
	str := strings.Split(job.Name, "-")
	if len(str) == 2 {
		jobName = str[0]
	} else if len(str) == 1 {
		jobName = job.Name
	} else {
		return res
	}
	cmd := ""
	params := map[string]int{}

	psNumber := 0
	workerNumber := 0
	flag := false
	for _, task := range job.Tasks {
		if (role == "PS" && task.IsPS) || (role == "Worker" && !task.IsPS) {
			params["num_gpus"] = task.NumberGPU
			cmd = task.Cmd
			flag = true
		}
		if task.IsPS {
			psNumber++
		} else {
			workerNumber++
		}
	}
	params["ps_number"] = psNumber
	params["worker_number"] = workerNumber
	if role == "PS" {
		params["role"] = 1
	} else {
		params["role"] = 0
	}
	if !flag {
		return res
	}

	exceptions := map[string]bool{}
	exceptions["train_dir"] = true
	exceptions["variable_update"] = true
	exceptions["ps_hosts"] = true
	exceptions["worker_hosts"] = true
	exceptions["task_index"] = true
	exceptions["job_name"] = true

	pairs := strings.Split(cmd, " ")
	for _, pair := range pairs {
		v := strings.Split(pair, "=")
		if len(v) == 2 && v[0][:2] == "--" {
			var param string
			var value int
			param = v[0][2:]

			if val, err := strconv.Atoi(v[1]); err == nil {
				value = val
			} else {
				h := fnv.New32a()
				h.Write([]byte(v[1]))
				value = int(h.Sum32())
			}
			if _, ok := exceptions[param]; !ok {
				params[param] = value
			}
		}
	}
	//log.Info(job.Name, params)

	features, _ := json.Marshal(params)

	spider := Spider{}
	spider.Method = "GET"
	spider.URL = "http://yao-optimizer:8080/predict?job=" + jobName + "&features=" + string(features)

	err := spider.do()
	if err != nil {
		return MsgJobReq{Code: 2, Error: err.Error()}
	}

	resp := spider.getResponse()
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Warn(err)
		return MsgJobReq{Code: 3, Error: err.Error()}
	}

	var msg MsgJobReqPredict
	err = json.Unmarshal([]byte(string(body)), &msg)
	if err == nil && msg.Code == 0 {
		tmp := msg.Labels
		if v, ok := tmp["cpu"]; ok {
			res.CPU = int(math.Ceil(v / 100))
		}
		if v, ok := tmp["mem"]; ok {
			res.Mem = int(math.Ceil(v/1024)) * 1024
		}
		if v, ok := tmp["gpu_util"]; ok {
			res.UtilGPU = int(math.Ceil(v)/10) * 10
		}
		if v, ok := tmp["gpu_mem"]; ok {
			res.MemGPU = int(math.Ceil(v/1024)) * 1024
		}
		if v, ok := tmp["bw_rx"]; ok {
			res.BW = int(math.Ceil(v / 100000))
		}
	}
	res.Version = time.Now().Unix()
	optimizer.cacheMu.Lock()
	optimizer.reqCache[job.Name] = res
	optimizer.cacheMu.Unlock()
	return res
}

func (optimizer *Optimizer) max(values []float64) float64 {
	value := 0.0
	for _, v := range values {
		if v > value {
			value = v
		}
	}
	return value
}

func (optimizer *Optimizer) mean(values []float64) float64 {
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func (optimizer *Optimizer) std(values []float64) float64 {
	mean := optimizer.mean(values)
	std := 0.0
	for j := 0; j < len(values); j++ {
		// The use of Pow math function func Pow(x, y float64) float64
		std += math.Pow(values[j]-mean, 2)
	}
	// The use of Sqrt math function func Sqrt(x float64) float64
	std = math.Sqrt(std / float64(len(values)))
	return std
}
