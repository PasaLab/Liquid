package main

import (
	"sync"
	"github.com/Shopify/sarama"
	"encoding/json"
	"time"
)

var collectorInstance *Collector
var collectorInstanceLock sync.Mutex

func InstanceOfCollector() *Collector {
	defer collectorInstanceLock.Unlock()
	collectorInstanceLock.Lock()

	if collectorInstance == nil {
		collectorInstance = &Collector{}
	}
	return collectorInstance
}

type Collector struct {
	wg sync.WaitGroup
}

func (collector *Collector) init(conf Configuration) {
	go func() {
		consumer, err := sarama.NewConsumer(conf.KafkaBrokers, nil)
		for {
			if err == nil {
				break
			}
			log.Warn(err)
			time.Sleep(time.Second * 5)
			consumer, err = sarama.NewConsumer(conf.KafkaBrokers, nil)
		}

		partitionList, err := consumer.Partitions(conf.KafkaTopic)
		if err != nil {
			panic(err)
		}

		for partition := range partitionList {
			pc, err := consumer.ConsumePartition(conf.KafkaTopic, int32(partition), sarama.OffsetNewest)
			if err != nil {
				panic(err)
			}
			defer pc.AsyncClose()

			collector.wg.Add(1)

			go func(sarama.PartitionConsumer) {
				defer collector.wg.Done()
				for msg := range pc.Messages() {
					go func(msg *sarama.ConsumerMessage) {
						var nodeStatus NodeStatus
						err = json.Unmarshal([]byte(string(msg.Value)), &nodeStatus)
						if err != nil {
							log.Warn(err)
							return
						}
						InstanceOfResourcePool().update(nodeStatus)
					}(msg)
				}

			}(pc)
		}
		collector.wg.Wait()
		consumer.Close()
	}()
}
