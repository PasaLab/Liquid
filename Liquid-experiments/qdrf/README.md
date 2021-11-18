## 基于改进DRF的公平调度
（说明：实验用了4个节点，使用更多节点需要增加作业数量 16 = 4 * 4）
设置5个队列，其中4个队列的每个作业需求2张GPU卡，另一个队列的每个作业需求16张
为了保证整体的GPU需求数量一致，每个队列的所有作业GPU需求总和为32张

为了避免在实验时由于资源分配任务放置方案的不同对结果造成影响，本实验中采用sleep来替代，
其中大作业的执行时间为10分钟，小作业为5分钟

#### 训练
（创建5个队列，只需要执行一次）
python3 main.py --lab=qdrf --mode=train --case=

#### 测试
python3 main.py --lab=qdrf --mode=test --case=qdrf-balanced
python3 main.py --lab=qdrf --mode=test --case=qdrf-unbalanced


为了避免由于系统其他模块的影响，本文实现了YARN中的DRF策略，并集成到本平台中
需要设置代码中以下两个变量（yao-scheduler:/scheduler_fair.go）
scheduler.drfyarn = true
scheduler.enableBorrow = false
然后重新启动调度器组件

调度器组件启动命令换成
```bash
docker service create \
	--name yao-scheduler \
	--hostname yao-scheduler \
	--constraint node.hostname==gj-slave103 \
	--network yao-net \
	--replicas 1 \
	--detach=true \
	--env SchedulerPolicy=fair \
	--env ListenAddr='0.0.0.0:8080' \
	--env HDFSAddress='' \
	--env HDFSBaseDir='/user/yao/output/' \
	--env DFSBaseDir='/dfs/yao-jobs/' \
	--env EnableShareRatio=1.75 \
	--env ShareMaxUtilization=1.30 \
	--env EnablePreScheduleRatio=1.75 \
	--env PreScheduleExtraTime=15 \
	--env PreScheduleTimeout=300 \
	--mount type=bind,source=/etc/localtime,target=/etc/localtime,readonly \
	quickdeploy/yao-scheduler:dev sleep infinity
```

然后进入到yao-scheduler容器，利用以下命令启动

cd src && go run .

然后做drf的实验，做完之后按照原来的方式恢复。

python3 main.py --lab=qdrf --mode=test --case=drf-balanced
python3 main.py --lab=qdrf --mode=test --case=drf-unbalanced

