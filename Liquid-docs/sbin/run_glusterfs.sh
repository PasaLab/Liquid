#!/bin/bash

docker run \
	--name gluster_server \
	-d \
	--restart always \
	--net=host \
	--privileged=true \
	-v /data/gluster/configuration:/etc/glusterfs:z \
	-v /data/gluster/metadata:/var/lib/glusterd:z \
	-v /data/gluster/logs:/var/log/glusterfs:z \
	-v /data/gluster/data:/data \
	gluster/gluster-centos
