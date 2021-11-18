## 单卡多任务共享优化实验

#### 训练
python3 main.py --lab=share --mode=train --case=default

#### 测试
要求：把其他计算节点都下线，只保留一个，确保所有作业是在同一个GPU上执行的（删除yao-agent等待自动下线）

python3 main.py --lab=share --mode=test --case=cnn-share
python3 main.py --lab=share --mode=test --case=cnn-exclusive
python3 main.py --lab=share --mode=test --case=neumf-share
python3 main.py --lab=share --mode=test --case=neumf-exclusive
python3 main.py --lab=share --mode=test --case=mined-share
python3 main.py --lab=share --mode=test --case=mixed-exclusive

在实验结束后，把残留的几个sleep作业手动停止。