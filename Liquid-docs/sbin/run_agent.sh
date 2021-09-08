#!/bin/bash

ip=`hostname --ip-address`

docker run \
	--gpus all \
	--name yao-agent \
	--pid=host \
	--network yao-net \
	--network-alias $(hostname) \
	--hostname $(hostname) \
	-d \
	--restart always \
	--detach=true \
	--publish 8000:8000 \
	--env ClientID=$(hostname) \
	--env ClientHost=$(hostname) \
	--env ClientExtHost=${ip} \
	--env PORT=8000 \
	--env HeartbeatInterval=5 \
	--env ReportAddress='http://yao-scheduler:8080/?action=agent_report' \
	--env EnableEventTrigger='true' \
	--env PYTHONUNBUFFERED=1 \
	--mount type=bind,source=/etc/localtime,target=/etc/localtime,readonly \
	--mount type=bind,src=/var/run/docker.sock,dst=/var/run/docker.sock \
	--mount type=bind,src=/dfs/yao-jobs/,dst=/dfs/yao-jobs/ \
	quickdeploy/yao-agent
