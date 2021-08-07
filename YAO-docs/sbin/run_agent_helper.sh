#!/bin/bash

docker run \
	--name yao-agent-helper \
	-d \
	--restart always \
	--detach=true \
	--mount type=bind,source=/etc/localtime,target=/etc/localtime,readonly \
	--mount type=bind,src=/var/run/docker.sock,dst=/var/run/docker.sock \
	--mount type=bind,src=/dfs/yao-jobs/,dst=/dfs/yao-jobs/ \
	docker:latest sleep 86400000
