# 基于一致性哈希的全局资源池管理优化性能对比实验

#### 实验准备
停止所有的yao-agent组件，然后启动若干个模拟从节点，在虚拟节点上执行。
分为两类作业，一个优化后（使用现有的代码就可以），二是优化前（按照下面的说明修改代码）

分别进行以下4组实验，然后计算调度阶段时长（第一个作业的启动时间到最后一个作业的启动时间）
python3 main.py --lab=hash --mode=test --case=parallelism-1
python3 main.py --lab=hash --mode=test --case=parallelism-2
python3 main.py --lab=hash --mode=test --case=parallelism-5
python3 main.py --lab=hash --mode=test --case=parallelism-10

每组实验结束，获得需要的数据后，重启yao-scheduler组件清理待调度队列

#### 优化前代码修改
修改以下两处代码，并且重新启动yao-scheduler（不要重启容器，代码会恢复）

1. yao-scheduler:/resource_pool.go#87

	pool.poolsCount = 5 # 固定5
	for i := 0; i < pool.poolsCount; i++ {
		pool.pools = append(pool.pools, PoolSeg{Lock: sync.Mutex{}, ID: i})
	}
	/* generate working segs */
	for i := 0; i < 1; i++ { # 固定1
		pool.pools[rand.Intn(pool.poolsCount)].Nodes = map[string]*NodeStatus{}
	}


2. yao-scheduler:/resource_pool.go#412
在scaleSeg函数的第一行添加
return
使得函数不执行

然后按照优化后的步骤再次执行