#! /bin/bash

docker service create \
	--name redis \
	--hostname redis \
	--constraint node.hostname==gj-slave101 \
	--network yao-net \
	--replicas 1 \
	--detach=true \
	--endpoint-mode dnsrr \
	--mount type=bind,source=/etc/localtime,target=/etc/localtime,readonly \
	redis redis-server --appendonly yes
