# Seagull
This is the code repository for the deep learning job scheduling paper titled 'Efficient Distributed Deep Learning Job Scheduling on GPU Clusters Based on Resource Requirement Prediction'.

The project is based on Docker.


# Prerequisites
- OS Centos Linux release7.6.1810
- Nvidia Driver 410.129
- CUDA 10.0
- Docker 19.03
- Nvidia-docker 2.2.2
 


# Steps to bring up the YAO components


## Init a docker swarm cluster
```bash
# on master node
docker swarm init

# Add other nodes to the cluster
docker swarm join --token A-LONG-TOKEN-STRING-HERE 192.168.0.1:2377
docker swarm leave
docker swarm leave --force
```


## Create an overlay network named `yao`
```bash
docker network create --driver overlay --attachable yao-net

# docker network create --driver overlay --attachable --opt encrypted yao-net
```

*Note: try remove encrypted when the containers cannot communicate cross nodes*


## Start HDFS cluster (Optional)
```bash
YAO-docs/sbin/run_hdfs.sh
```

## Start GlusterFS cluster (Optional)
```bash
YAO-docs/sbin/run_glusterfs.sh
```


## Start the agents in each YAO-Worker
```bash
YAO-docs/sbin/run_agent_helper.sh

YAO-docs/sbin/run_agent.sh
```

## Start the agent-master on YAO-Master
```bash
YAO-docs/sbin/start_agent_master.sh
```


## Start mysql
```bash
YAO-docs/sbin/start_mysql.sh
```

## Start yao-optimizer on Master Node
```bash
YAO-docs/sbin/run_optimizer.sh
```

## Start yao-scheduler
```bash
YAO-docs/sbin/start_scheduler.sh
```

## Start Redis
```bash
YAO-docs/sbin/start_redis.sh
```

## Start the web portal
```bash
YAO-docs/sbin/start_portal.sh
```
## Start gitea
```bash
YAO-docs/sbin/start_gitea.sh
```


Visit `http://YOUR_IP/install.php`

