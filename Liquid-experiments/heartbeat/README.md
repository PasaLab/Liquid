
1. 背景：
心跳机制广泛应用于现有的集群管理系统中，尤其是主从架构。但是固定间隔的心跳机制不利于资源的快速回收利用，所以采用事件驱动与心跳机制相结合的消息通信机制。

对于调度系统来说，一个资源被重新利用需要经过以下几个步骤：
1）从节点定期检测作业执行状态，以心跳的形式发送给调度器
2）调度器检测到作业退出，释放资源
3）下一个作业申请资源并得到该资源
4）通知从节点执行下一个作业
5）从节点启动作业

可以看到，在前一个作业执行完毕到下一个作业开始执行，经过了很多步骤，带来了很大的延迟，其中一个原因是调度器无法及时感知到作业退出。
通过降低心跳间隔可以一定程度上缓解，但是会造成调度器负载增加，
为了解决这个问题，在现有的心跳机制上，结合事件驱动的形式，来尽可能让被释放的资源尽快得到重新利用。

2. 实现：
原有的心跳机制是定期检查节点和作业执行状态，然后推送给调度器，
现在在本地增加检测的频率，但心跳包的间隔不变，而当检测到作业退出时，不再等待，立即发送新的心跳包

在具体实现上，利用docker的event API，检测容器运行状态的变化，当出现容器退出事件时，
当监听到的事件为die时，立即触发一次心跳操作

3. 实验
实验中只保留一个节点，然后进入最后一个yao-agent容器，记录GPU使用情况
python3 monitor.py

python3 main.py --lab=heartbeat --mode=test --case=after

将yao-agent的
--env EnableEventTrigger='true'
换成
--env EnableEventTrigger='false'
然后重新启动，进行实验

python3 main.py --lab=heartbeat --mode=test --case=before
在实验完毕后，需要恢复yao-agent中的参数