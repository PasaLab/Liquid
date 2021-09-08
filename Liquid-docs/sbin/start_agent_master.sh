#!/bin/bash

docker service create \
	--name yao-agent-master \
	--hostname yao-agent-master \
	--constraint node.hostname==gj-slave103 \
	--network yao-net \
	--replicas 1 \
	--detach=true \
	--mount type=bind,source=/etc/localtime,target=/etc/localtime,readonly \
	--mount type=bind,src=/var/run/docker.sock,dst=/var/run/docker.sock \
	quickdeploy/yao-agent-master:dev
