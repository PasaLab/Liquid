## 批量调度策略效果评估
3个节点
#### 背景
当出现一定数量空闲的资源时，才调度，通过增强亲和性来增加性能

#### 训练
需要预先提交若干作业来训练资源需求向量预测模型，与单作业调度策略共用，所以不需要再训练。

#### 测试
python3 main.py -lab batch -mode test -case multi-random
python3 main.py -lab batch -mode test -case multi-spread
python3 main.py -lab batch -mode test -case multi-pack
python3 main.py -lab batch -mode test -case multi-topology-aware
python3 main.py -lab batch -mode test -case multi-batch
python3 main.py -lab batch -mode test -case multi-k8s-default
python3 main.py -lab batch -mode test -case multi-k8s-affinity
python3 main.py -lab batch -mode test -case multi-k8s-kubeflow
python3 main.py -lab batch -mode test -case multi-k8s-hived