# 基于一致性哈希的全局资源池管理优化性能对比实验

#### 实验准备
停止所有的yao-agent组件，然后启动若干个模拟从节点，在虚拟节点上执行

有两个参数需要调节，一是虚拟节点的数量，通过以下命令来启动
（每个物理机上最多执行一次，尽可能将虚拟节点部署到6台物理机上）

NUM=100 bash sbin/run_mock.sh
（这条命令启动了100个虚拟节点）

（在所有的虚拟节点都上线之后！！！在前端面板上可以看）
提交若干作业（保证可以一次性调度完）
每条命令对应一组实验，在上一组实验结束后再进行下一个

python3 main.py --lab=concurrent --mode=test --case=parallelism-1
python3 main.py --lab=concurrent --mode=test --case=parallelism-2
python3 main.py --lab=concurrent --mode=test --case=parallelism-5
python3 main.py --lab=concurrent --mode=test --case=parallelism-10

每组实验结束，获得需要的数据后，重启yao-scheduler组件清理待调度队列

共需要进行下表所示的实验数量

| nodes | parallelism-1 | parallelism-2 | parallelism-5 | parallelism-10 |
| ----- | ------------- | ------------- | ------------- | -------------- |
|  200  |       ?       |       ?       |       ?       |       ?        |
|  300  |       ?       |       ?       |       ?       |       ?        |
|  400  |       ?       |       ?       |       ?       |       ?        |
|  500  |       ?       |       ?       |       ?       |       ?        |
|  600  |       ?       |       ?       |       ?       |       ?        |
|  700  |       ?       |       ?       |       ?       |       ?        |
|  800  |       ?       |       ?       |       ?       |       ?        |
|  900  |       ?       |       ?       |       ?       |       ?        |
|  1000 |       ?       |       ?       |       ?       |       ?        |
|  1500 |       ?       |       ?       |       ?       |       ?        |
|  2000 |       ?       |       ?       |       ?       |       ?        |
|  2500 |       ?       |       ?       |       ?       |       ?        |
|  3000 |       ?       |       ?       |       ?       |       ?        |

在每组实验结束后，计算第一个作业启动时间到最后一个作业启动时间，即调度阶段时长。