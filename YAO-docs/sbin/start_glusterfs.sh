#!/bin/bash

# Setup
docker run \
	--name gluster_server \
	--restart always \
	-d \
	--privileged=true \
	--net=host \
	-v /data/gluster/configuration:/etc/glusterfs:z \
	-v /data/gluster/metadata:/var/lib/glusterd:z \
	-v /data/gluster/logs:/var/log/glusterfs:z \
	-v /data/gluster/data:/data \
	gluster/gluster-centos

exit 0;


# Init
gluster peer probe 192.168.100.101
gluster peer probe 192.168.100.102
gluster peer probe 192.168.100.103
gluster peer probe 192.168.100.104
gluster peer probe 192.168.100.105
gluster peer probe 192.168.100.106
gluster peer status

gluster peer detach 192.168.100.104 force
# Create & Start Volume
gluster volume create yao replica 3 192.168.100.101:/data/yao 192.168.100.102:/data/yao 192.168.100.103:/data/yao 192.168.100.104:/data/yao 192.168.100.105:/data/yao 192.168.100.106:/data/yao

gluster volume start yao

gluster volume status

systemctl status glusterfsd

gluster volume stop yao

gluster volume delete yao
gluster volume remove-brick yao replica 3 192.168.100.101:/data/yao force
gluster volume remove-brick yao 192.168.100.104:/data/yao 192.168.100.105:/data/yao 192.168.100.106:/data/yao force
rm -rf /data/yao/yao-jobs/*
setfattr -x trusted.glusterfs.volume-id /data/yao
setfattr -x trusted.gfid /data/yao
# Client Mount
sudo yum install glusterfs glusterfs-fuse attr -y

sudo rmdir /dfs

sudo mkdir -p /dfs && sudo chmod -R 777 /dfs

sudo mount -t glusterfs 192.168.100.102:/yao /dfs

mkdir -p /dfs/yao-jobs

sudo umount -v /dfs

sudo chmod 777 /dfs

cd /data/gluster/data/yao/.glusterfs

rm -rf /data/gluster/data/yao/.glusterfs/??

sudo rm -rf /data/gluster/configuration

sudo rm -rf /data/gluster/metadata

sudo rm -rf /data/gluster/logs

sudo rm -rf /data/gluster/data

sudo mkdir -p /data/gluster/configuration

sudo mkdir -p /data/gluster/metadata

sudo mkdir -p /data/gluster/logs

sudo mkdir -p /data/gluster/data