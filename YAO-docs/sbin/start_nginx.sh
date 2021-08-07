#! /bin/bash

docker service create \
	--name nginx \
	--hostname nginx \
	--constraint node.hostname==gj-slave101 \
	--network yao-net \
	--replicas 1 \
	--detach=true \
	--publish mode=host,published=80,target=80 \
	--publish mode=host,published=443,target=443 \
	--mount type=bind,src=/etc/localtime,dst=/etc/localtime,readonly \
	--mount type=bind,src=/data/nginx/conf.d/,dst=/etc/nginx/conf.d/,readonly \
	nginx
