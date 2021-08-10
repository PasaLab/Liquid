#!/bin/bash

docker run \
	--name yao-optimizer \
	--hostname yao-optimizer \
	--constraint node.hostname==gj-slave103 \
	--network yao-net \
	--network-alias yao-optimizer \
	-d \
	--restart always \
	--detach=true \
	--env PYTHONUNBUFFERED=1 \
	--mount type=bind,source=/etc/localtime,target=/etc/localtime,readonly \
	onceicy/yao-optimizer:test3
