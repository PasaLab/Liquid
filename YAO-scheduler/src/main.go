package main

import (
	"net/http"
	"encoding/json"
	"time"
	"strconv"
	"math/rand"
	"os"
	"fmt"
)

var log Logger

var scheduler Scheduler

func serverAPI(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Query().Get("action") {
	/* resource pool */
	case "agent_report":
		log.Debug("agent_report")
		msgAgentReport := MsgAgentReport{Code: 0}
		var nodeStatus NodeStatus
		err := json.Unmarshal([]byte(string(r.PostFormValue("data"))), &nodeStatus)
		if err != nil {
			msgAgentReport.Code = 1
			msgAgentReport.Error = err.Error()
			log.Warn(err)
		} else {
			go func() {
				InstanceOfResourcePool().update(nodeStatus)
			}()
		}
		js, err := json.Marshal(msgAgentReport)
		if err != nil {
			log.Warn(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		break

	case "resource_list":
		js, _ := json.Marshal(InstanceOfResourcePool().list())
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		break

	case "resource_get_by_node":
		id := r.URL.Query().Get("id")
		js, _ := json.Marshal(InstanceOfResourcePool().getByID(id))
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		break

	case "pool_status_history":
		log.Debug("pool_status_history")
		js, _ := json.Marshal(InstanceOfResourcePool().statusHistory())
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		break

	case "get_counter":
		log.Debug("get_counters")
		js, _ := json.Marshal(InstanceOfResourcePool().getCounter())
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		break

	case "debug_pool_dump":
		log.Debug("debug_pool_dump")
		js, _ := json.Marshal(InstanceOfResourcePool().Dump())
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		break

		/* scheduler */
	case "job_submit":
		var job Job
		log.Debug("job_submit")
		msgSubmit := MsgSubmit{Code: 0}
		err := json.Unmarshal([]byte(string(r.PostFormValue("job"))), &job)

		if err != nil {
			msgSubmit.Code = 1
			msgSubmit.Error = err.Error()
		} else if len(job.Tasks) == 0 {
			msgSubmit.Code = 2
			msgSubmit.Error = "task not found in the job"
		} else if InstanceOfGroupManager().get(job.Group) == nil {
			msgSubmit.Code = 2
			msgSubmit.Error = "Queue not found"
		} else if len(job.Workspace) == 0 {
			msgSubmit.Code = 2
			msgSubmit.Error = "Docker image cannot be empty"
		} else {
			job.Name = job.Name + "-"
			job.Name += strconv.FormatInt(time.Now().UnixNano(), 10)
			job.Name += strconv.Itoa(10000 + rand.Intn(89999))
			bwWorker := InstanceOfOptimizer().PredictReq(job, "Worker").BW
			for i := range job.Tasks {
				job.Tasks[i].ID = job.Name + ":" + job.Tasks[i].Name
				job.Tasks[i].Job = job.Name
				job.Tasks[i].BW = bwWorker
			}
			job.CreatedAt = int(time.Now().Unix())
			msgSubmit.JobName = job.Name
			scheduler.Schedule(job)
		}
		js, err := json.Marshal(msgSubmit)
		if err != nil {
			log.Warn(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		break

	case "job_status":
		log.Debug("job_status")
		js, _ := json.Marshal(scheduler.QueryState(r.URL.Query().Get("id")))
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		break

	case "job_stop":
		log.Debug("job_stop")
		js, _ := json.Marshal(scheduler.Stop(string(r.PostFormValue("id"))))
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		break

	case "task_logs":
		log.Debug("task_logs")
		js, _ := json.Marshal(scheduler.QueryLogs(r.URL.Query().Get("job"), r.URL.Query().Get("task")))
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		break

	case "jobs":
		log.Debug("job_list")
		js, _ := json.Marshal(scheduler.ListJobs())
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		break

	case "debug_scheduler_dump":
		log.Debug("debug_scheduler_dump")
		js, _ := json.Marshal(scheduler.DebugDump())
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		break

	case "summary":
		log.Debug("summary")
		js, _ := json.Marshal(scheduler.Summary())
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		break

		/* optimizer */
	case "job_predict_req":
		log.Debug("job_predict_req")
		var job Job
		role := r.URL.Query().Get("role")
		err := json.Unmarshal([]byte(string(r.PostFormValue("job"))), &job)
		msgJobReq := MsgJobReq{Code: 0}
		if err != nil {
			msgJobReq.Code = 1
			msgJobReq.Error = err.Error()
		} else {
			job.Name = job.Name + "-"
			job.Name += strconv.FormatInt(time.Now().UnixNano(), 10)
			job.Name += strconv.Itoa(10000 + rand.Intn(89999))
			msgJobReq = InstanceOfOptimizer().PredictReq(job, role)
		}
		js, err := json.Marshal(msgJobReq)
		if err != nil {
			log.Warn(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		break

	case "job_predict_time":
		log.Debug("job_predict_time")
		var job Job
		err := json.Unmarshal([]byte(string(r.PostFormValue("job"))), &job)
		msgJobReq := MsgOptimizerPredict{Code: 0}
		if err != nil {
			msgJobReq.Code = 1
			msgJobReq.Error = err.Error()
		} else {
			job.Name = job.Name + "-"
			job.Name += strconv.FormatInt(time.Now().UnixNano(), 10)
			job.Name += strconv.Itoa(10000 + rand.Intn(89999))
			msg := InstanceOfOptimizer().PredictTime(job)
			msgJobReq.Pre = msg.Pre
			msgJobReq.Post = msg.Post
			msgJobReq.Total = msg.Total
		}
		js, err := json.Marshal(msgJobReq)
		if err != nil {
			log.Warn(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		break

		/* job history logger */
	case "jhl_job_status":
		log.Debug("jhl_job_status")
		js, _ := json.Marshal(InstanceJobHistoryLogger().getTaskStatus(r.URL.Query().Get("job")))
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		break

		/* group */
	case "group_list":
		log.Debug("group_list")
		js, _ := json.Marshal(InstanceOfGroupManager().List())
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		break

	case "group_add":
		log.Debug("group_add")
		var group Group
		msg := MsgGroupCreate{Code: 0}
		err := json.Unmarshal([]byte(string(r.PostFormValue("group"))), &group)
		if err != nil {
			msg.Code = 1
			msg.Error = err.Error()
		} else {
			msg = InstanceOfGroupManager().Add(group)
			scheduler.updateGroup(group)
		}
		js, _ := json.Marshal(msg)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		break

	case "group_update":
		log.Debug("group_update")
		var group Group
		msg := MsgGroupCreate{Code: 0}
		err := json.Unmarshal([]byte(string(r.PostFormValue("group"))), &group)
		if err != nil {
			msg.Code = 1
			msg.Error = err.Error()
		} else {
			msg = InstanceOfGroupManager().Update(group)
			scheduler.updateGroup(group)
		}
		js, _ := json.Marshal(msg)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		break

	case "group_remove":
		log.Debug("group_remove")
		var group Group
		msg := MsgGroupCreate{Code: 0}
		err := json.Unmarshal([]byte(string(r.PostFormValue("group"))), &group)
		if err != nil {
			msg.Code = 1
			msg.Error = err.Error()
		} else {
			msg = InstanceOfGroupManager().Remove(group)
			scheduler.updateGroup(group)
		}
		js, _ := json.Marshal(msg)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		break

		/* configuration */
	case "conf_list":
		log.Debug("conf_list")
		var msg MsgConfList
		msg.Code = 0
		msg.Options = InstanceOfConfiguration().Dump()
		js, _ := json.Marshal(msg)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		break

	case "conf_update":
		log.Debug("conf_update")
		option := r.URL.Query().Get("option")
		value := r.URL.Query().Get("value")
		ok := false
		switch option {
		/* pool.share */
		case "pool.share.enable_threshold":
			if threshold, err := strconv.ParseFloat(value, 32); err == nil {
				ok = InstanceOfConfiguration().SetShareRatio(threshold)
			}
			break

		case "pool.share.max_utilization":
			util, err := strconv.ParseFloat(value, 32)
			if err == nil {
				ok = InstanceOfConfiguration().SetShareMaxUtilization(util)
			}
			break

			/* pool.pre_schedule */
		case "pool.pre_schedule.enable_threshold":
			if threshold, err := strconv.ParseFloat(value, 32); err == nil {
				ok = InstanceOfConfiguration().SetPreScheduleRatio(threshold)
			}
			break

			/* pool.batch */
		case "pool.batch.enabled":
			ok = InstanceOfResourcePool().SetBatchEnabled(value == "true")
			break

		case "pool.batch.interval":
			if interval, err := strconv.Atoi(value); err == nil {
				ok = InstanceOfResourcePool().SetBatchInterval(interval)
			}
			break

			/* scheduler.strategy */
			/* TODO: move jobs */
		case "scheduler.strategy":
			ok = InstanceOfConfiguration().SetSchedulePolicy(value)
			scheduler = InstanceOfConfiguration().GetScheduler()
			scheduler.Start()
			break

			/* scheduler.mock */
		case "scheduler.mock.enabled":
			ok = InstanceOfConfiguration().SetMockEnabled(value == "true")
			break

			/* scheduler.enabled */
		case "scheduler.enabled":
			ok = scheduler.SetEnabled(value == "true")
			break

			/* scheduler.parallelism */
		case "scheduler.parallelism":
			if parallelism, err := strconv.Atoi(value); err == nil {
				ok = scheduler.UpdateParallelism(parallelism)
			}
			break

			/* scheduler.preempt_enabled */
		case "scheduler.preempt_enabled":
			ok = InstanceOfConfiguration().SetPreemptEnabled(value == "true")
			break

			/* allocator.strategy */
		case "allocator.strategy":
			ok = InstanceOfAllocator().updateStrategy(value)
			break

			/* logger */
		case "logger.level":
			ok = log.SetLoggerLevel(value)
			break

		case "logger.enable_module":
			ok = log.LoggerEnableModule(value)
			break

		case "logger.disable_module":
			ok = log.LoggerDisableModule(value)
			break

		case "scheduler.job_max_retries":
			if maxRetries, err := strconv.Atoi(value); err == nil {
				ok = InstanceOfConfiguration().SetJobMaxRetries(maxRetries)
			}
			break

		}
		var msg MsgConfUpdate
		msg.Code = 0
		if !ok {
			msg.Code = 1
			msg.Error = fmt.Sprintf("Option (%s) not exist or invalid value (%s)", option, value)
		}
		js, _ := json.Marshal(msg)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		break

	default:
		http.Error(w, "Not Found", http.StatusNotFound)
		break
	}
}

func main() {
	log = Logger{}
	log.Init()
	loggerDir := os.Getenv("LoggerOutputDir")
	if len(loggerDir) != 0 {
		if _, err := os.Stat(loggerDir); os.IsNotExist(err) {
			os.Mkdir(loggerDir, os.ModePerm)
		}
		t := time.Now()
		file := t.Format("20060102.15_04_05") + ".log"
		f, err := os.OpenFile(loggerDir+file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		defer f.Close()
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		log.SetOutput(f)
	}

	config := InstanceOfConfiguration()
	config.InitFromEnv()

	/* init components */
	InstanceOfResourcePool().Start()
	InstanceJobHistoryLogger().Start()
	InstanceOfOptimizer().Start()
	InstanceOfGroupManager().Start()

	scheduler = config.GetScheduler()
	scheduler.Start()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serverAPI(w, r)
	})

	err := http.ListenAndServe(config.ListenAddr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
