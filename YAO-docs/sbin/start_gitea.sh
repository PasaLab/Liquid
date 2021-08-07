#!/bin/bash

docker service create \
	--name gitea \
	--hostname gitea \
	--constraint node.hostname==gj-slave105 \
	--network yao-net \
	--replicas 1 \
	--detach=true \
	--mount type=bind,source=/etc/localtime,target=/etc/localtime,readonly \
	--mount type=bind,src=/data/gitea/,target=/data \
	gitea/gitea
