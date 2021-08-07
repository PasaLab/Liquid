#!/bin/bash

# Run on all HDFS nodes, and update the conf files

exit;

docker run \
        --name hdfs \
        --hostname $(hostname) \
        -d \
        --restart always \
        --net host \
	--add-host=$(node1):$(ip1) \
	--add-host=$(node2):$(ip2) \
	--add-host=$(node3):$(ip3) \
        --mount type=bind,src=/data/hadoop/config,dst=/config/hadoop \
        --mount type=bind,src=/data/hadoop/hdfs,dst=/tmp/hadoop-root \
        --mount type=bind,src=/data/hadoop/logs,dst=/usr/local/hadoop/logs \
        newnius/hadoop:2.7.4
