package main

import (
	"sync"
	"os"
	"strings"
	"strconv"
)

type Configuration struct {
	KafkaBrokers           []string `json:"KafkaBrokers"`
	KafkaTopic             string   `json:"KafkaTopic"`
	SchedulerPolicy        string   `json:"scheduler.strategy"`
	ListenAddr             string   `json:"ListenAddr"`
	HDFSAddress            string   `json:"HDFSAddress"`
	HDFSBaseDir            string   `json:"HDFSBaseDir"`
	DFSBaseDir             string   `json:"DFSBaseDir"`
	EnableShareRatio       float64  `json:"EnableShareRatio"`
	ShareMaxUtilization    float64  `json:"ShareMaxUtilization"`
	EnablePreScheduleRatio float64  `json:"EnablePreScheduleRatio"`
	PreScheduleExtraTime   int      `json:"PreScheduleExtraTime"` /* seconds of schedule ahead except pre+post */
	PreScheduleTimeout     int      `json:"PreScheduleTimeout"`
	JobMaxRetries          int      `json:"scheduler.job_max_retries"`
	PreemptEnabled         bool     `json:"scheduler.preempt_enabled"`

	mock bool
	mu   sync.Mutex
}

var configurationInstance *Configuration
var ConfigurationInstanceLock sync.Mutex

func InstanceOfConfiguration() *Configuration {
	ConfigurationInstanceLock.Lock()
	defer ConfigurationInstanceLock.Unlock()

	if configurationInstance == nil {
		/* set default values */
		configurationInstance = &Configuration{
			mock: false,
			KafkaBrokers: []string{
				"kafka-node1:9092",
				"kafka-node2:9092",
				"kafka-node3:9092",
			},
			KafkaTopic:             "yao",
			SchedulerPolicy:        "fair",
			ListenAddr:             "0.0.0.0:8080",
			HDFSAddress:            "",
			HDFSBaseDir:            "/user/root/",
			DFSBaseDir:             "",
			EnableShareRatio:       1.5,
			ShareMaxUtilization:    1.3, // more than 1.0 to expect more improvement
			EnablePreScheduleRatio: 1.5,
			PreScheduleExtraTime:   15,
			JobMaxRetries:          0,
			PreemptEnabled:         false,
		}
	}
	return configurationInstance
}

/* read conf value from env */
func (config *Configuration) InitFromEnv() {
	value := os.Getenv("KafkaBrokers")
	if len(value) != 0 {
		configurationInstance.KafkaBrokers = strings.Split(value, ",")
	}
	value = os.Getenv("KafkaTopic")
	if len(value) != 0 {
		configurationInstance.KafkaTopic = value
	}
	value = os.Getenv("SchedulerPolicy")
	if len(value) != 0 {
		configurationInstance.SchedulerPolicy = value
	}
	value = os.Getenv("ListenAddr")
	if len(value) != 0 {
		configurationInstance.ListenAddr = value
	}
	value = os.Getenv("HDFSAddress")
	if len(value) != 0 {
		configurationInstance.HDFSAddress = value
	}
	value = os.Getenv("HDFSBaseDir")
	if len(value) != 0 {
		configurationInstance.HDFSBaseDir = value
	}
	value = os.Getenv("DFSBaseDir")
	if len(value) != 0 {
		configurationInstance.DFSBaseDir = value
	}
	value = os.Getenv("EnableShareRatio")
	if len(value) != 0 {
		if val, err := strconv.ParseFloat(value, 32); err == nil {
			configurationInstance.EnableShareRatio = val
		}
	}
	value = os.Getenv("ShareMaxUtilization")
	if len(value) != 0 {
		if val, err := strconv.ParseFloat(value, 32); err == nil {
			configurationInstance.ShareMaxUtilization = val
		}
	}
	value = os.Getenv("EnablePreScheduleRatio")
	if len(value) != 0 {
		if val, err := strconv.ParseFloat(value, 32); err == nil {
			configurationInstance.EnablePreScheduleRatio = val
		}
	}
	value = os.Getenv("PreScheduleExtraTime")
	if len(value) != 0 {
		if val, err := strconv.Atoi(value); err == nil {
			configurationInstance.PreScheduleExtraTime = val
		}
	}
	value = os.Getenv("PreScheduleTimeout")
	if len(value) != 0 {
		if val, err := strconv.Atoi(value); err == nil {
			configurationInstance.PreScheduleTimeout = val
		}
	}
	value = os.Getenv("scheduler.job_max_retries")
	if len(value) != 0 {
		if val, err := strconv.Atoi(value); err == nil && val >= 0 {
			configurationInstance.JobMaxRetries = val
		}
	}
	value = os.Getenv("scheduler.preempt_enabled")
	configurationInstance.PreemptEnabled = value == "true"
}

func (config *Configuration) SetMockEnabled(enabled bool) bool {
	config.mu.Lock()
	defer config.mu.Unlock()
	config.mock = enabled
	log.Info("configuration.mock = ", enabled)
	return true
}

func (config *Configuration) SetShareRatio(ratio float64) bool {
	config.mu.Lock()
	defer config.mu.Unlock()
	config.EnableShareRatio = ratio
	log.Info("enableShareRatio is updated to ", ratio)
	return true
}

func (config *Configuration) SetPreScheduleRatio(ratio float64) bool {
	config.mu.Lock()
	defer config.mu.Unlock()
	config.EnablePreScheduleRatio = ratio
	log.Info("enablePreScheduleRatio is updated to ", ratio)
	return true
}

func (config *Configuration) SetShareMaxUtilization(value float64) bool {
	config.mu.Lock()
	defer config.mu.Unlock()
	config.ShareMaxUtilization = value
	log.Info("ShareMaxUtilization is set to ", value)
	return true
}

func (config *Configuration) SetJobMaxRetries(value int) bool {
	config.mu.Lock()
	defer config.mu.Unlock()
	config.JobMaxRetries = value
	log.Info("scheduler.job_max_retries is set to ", value)
	return true
}

func (config *Configuration) SetPreemptEnabled(enabled bool) bool {
	config.mu.Lock()
	defer config.mu.Unlock()
	config.PreemptEnabled = enabled
	log.Info("scheduler.preempt_enabled is set to ", enabled)
	return true
}

func (config *Configuration) SetSchedulePolicy(strategy string) bool {
	config.mu.Lock()
	defer config.mu.Unlock()
	available := map[string]bool{
		"FCFS":     true,
		"priority": true,
		"capacity": true,
		"fair":     true,
	}
	if _, ok := available[strategy]; ok {
		config.SchedulerPolicy = strategy
	}
	log.Info("scheduler.strategy is set to ", config.SchedulerPolicy)
	return true
}

func (config *Configuration) GetScheduler() Scheduler {
	config.mu.Lock()
	defer config.mu.Unlock()
	var scheduler Scheduler
	switch config.SchedulerPolicy {
	case "FCFS":
		scheduler = &SchedulerFCFS{}
		break
	case "priority":
		scheduler = &SchedulerPriority{}
		break
	case "capacity":
		scheduler = &SchedulerCapacity{}
		break
	case "fair":
		scheduler = &SchedulerFair{}
		break
	default:
		scheduler = &SchedulerFCFS{}
	}
	return scheduler
}

func (config *Configuration) Dump() map[string]interface{} {
	config.mu.Lock()
	defer config.mu.Unlock()
	res := map[string]interface{}{}
	res["KafkaBrokers"] = config.KafkaBrokers
	res["KafkaTopic"] = config.KafkaTopic
	res["SchedulerPolicy"] = config.SchedulerPolicy
	res["ListenAddr"] = config.ListenAddr
	res["Mock"] = config.mock
	res["HDFSAddress"] = config.HDFSAddress
	res["HDFSBaseDir"] = config.HDFSBaseDir
	res["DFSBaseDir"] = config.DFSBaseDir
	res["EnableShareRatio"] = config.EnableShareRatio
	res["ShareMaxUtilization"] = config.ShareMaxUtilization
	res["EnablePreScheduleRatio"] = config.EnablePreScheduleRatio
	res["PreScheduleExtraTime"] = config.PreScheduleExtraTime
	res["PreScheduleTimeout"] = config.PreScheduleTimeout
	res["scheduler.preempt_enabled"] = config.PreemptEnabled
	res["logger.level"] = log.LoggerLevel
	res["logger.modules_disabled"] = log.LoggerModuleDisabled
	return res
}
