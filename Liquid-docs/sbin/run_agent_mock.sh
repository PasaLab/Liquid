#!/bin/bash

# docker ps -a | grep yao-agent:mock | awk '{print $1}' | xargs docker rm -f

#for ((i=1; i<=${NUM}; i++));
#do
#	docker run \
#		--name yao-agent-mock-$i \
#		--network yao-net \
#		--network-alias node-mock-$i \
#		--hostname node-mock-$i \
#		-d \
#		--restart always \
#		--detach=true \
#		--env NUMS=100 \
#		--env ClientHost=$(hostname) \
#		--env HeartbeatInterval=5 \
#		--env ReportAddress='http://yao-scheduler:8080/?action=agent_report' \
#		--env PYTHONUNBUFFERED=1 \
#		--mount type=bind,source=/etc/localtime,target=/etc/localtime,readonly \
#		--mount type=bind,src=/var/run/docker.sock,dst=/var/run/docker.sock \
#		quickdeploy/yao-agent:mock

#		# reduce instant load
#		sleep 1
#done


#!/bin/bash

# Example: NUM=100 bash sbin/run_agent_mock.sh

# docker ps -a | grep yao-agent:mock | awk '{print $1}' | xargs docker rm -f

docker run \
	--name yao-agent-mock-$(hostname) \
	--network yao-net \
	--network-alias node-mock-$(hostname) \
	--hostname node-mock-$(hostname) \
	-d \
	--restart always \
	--detach=true \
	--env NUMS=${NUM} \
	--env ClientHost=node-mock-$(hostname) \
	--env ReportAddress='http://yao-scheduler:8080/?action=agent_report' \
	--env PYTHONUNBUFFERED=1 \
	--mount type=bind,source=/etc/localtime,target=/etc/localtime,readonly \
	--mount type=bind,src=/var/run/docker.sock,dst=/var/run/docker.sock \
	quickdeploy/yao-agent:mock