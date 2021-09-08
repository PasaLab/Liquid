#! /bin/bash

docker service create \
	--name mysql \
	--hostname mysql \
	--constraint node.hostname==gj-slave101 \
	--network yao-net \
	--replicas 1 \
	--detach=true \
	--endpoint-mode dnsrr \
	-e MYSQL_ROOT_PASSWORD=123456 \
	-e MYSQL_DATABASE=yao \
	--mount type=bind,source=/etc/localtime,target=/etc/localtime,readonly \
	mysql:5.7

#--mount type=bind,source=/data/mysql,target=/var/lib/mysql \
