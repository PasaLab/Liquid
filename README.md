# YAO
This is the code repository for the deep learning job scheduling paper titled 'Efficient Distributed Deep Learning Job Scheduling on GPU Clusters Based onResource Requirement Prediction'.

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
sbin/run_hdfs.sh
```

## Start GlusterFS cluster (Optional)
```bash
sbin/start_glusterfs.sh
```


## Start the agents in each YAO-Worker
```bash
sbin/run_agent_helper.sh

sbin/run_agent.sh
```

## Start the agent-master on YAO-Master
```bash
sbin/start_agent_master.sh
```


## Start mysql
```bash
sbin/start_mysql.sh
```

## Start yao-optimizer on Master Node
```bash
sbin/run_optimizer.sh
```

## Start yao-scheduler
```bash
sbin/start_scheduler.sh
```

## Start Redis
```bash
sbin/start_redis.sh
```

## Start the web portal
```bash
sbin/start_portal.sh
```
## Start gitea
```bash
sbin/start_gitea.sh
```


Visit `http://YOUR_IP/install.php`
> After that, you can visit the webpage, like `slave001:8888`, and then import packages for tests.
> The default port is `8888`. You can change it according to your requirements.
