# Liquid
This is the code repository for the deep learning job scheduling paper titled 'Liquid: Intelligent Resource Requirement Estimation and Scheduling for Deep Learning Jobs on Distributed GPU Clusters'.

The project is based on Docker.


# Prerequisites
- OS Centos Linux release7.6.1810
- Nvidia Driver 410.129
- CUDA 10.0
- Docker 19.03
- Nvidia-docker 2.2.2
 


# Steps to bring up the Liquid components


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
Liquid-docs/sbin/run_hdfs.sh
```

## Start GlusterFS cluster (Optional)
```bash
Liquid-docs/sbin/run_glusterfs.sh
```


## Start the agents in each Liquid-Worker
```bash
Liquid-docs/sbin/run_agent_helper.sh

Liquid-docs/sbin/run_agent.sh
```

## Start the agent-master on Liquid-Master
```bash
Liquid-docs/sbin/start_agent_master.sh
```


## Start mysql
```bash
Liquid-docs/sbin/start_mysql.sh
```

## Start Liquid-optimizer on Master Node
```bash
Liquid-docs/sbin/run_optimizer.sh
```

## Start Liquid-scheduler
```bash
Liquid-docs/sbin/start_scheduler.sh
```

## Start Redis
```bash
Liquid-docs/sbin/start_redis.sh
```

## Start the web portal
```bash
Liquid-docs/sbin/start_portal.sh
```
## Start gitea
```bash
Liquid-docs/sbin/start_gitea.sh
```


Visit `http://YOUR_IP/install.php`

