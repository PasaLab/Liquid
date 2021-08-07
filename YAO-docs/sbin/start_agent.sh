#!/bin/bash

docker service create \
	--name yao-agent \
	--hostname {{.Node.Hostname}} \
	--network name=yao-net,alias={{.Node.Hostname}} \
	--mode global \
	--detach=true \
	--env ClientID={{.Node.Hostname}} \
	--env ClientHost={{.Node.Hostname}} \
	--env KafkaBrokers=kafka-node1:9092,kafka-node2:9092,kafka-node3:9092 \
	--mount type=bind,source=/etc/localtime,target=/etc/localtime,readonly \
	--mount type=bind,src=/var/run/docker.sock,dst=/var/run/docker.sock \
	quickdeploy/yao-agent
