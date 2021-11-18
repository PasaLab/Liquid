# 提前调度优化实验

#### 训练
python3 main.py --lab=pre_schedule --mode=train --case=default
多节点训练

#### 测试
要求：把其他计算节点都下线，只保留一个，确保所有作业是在同一个GPU上执行的（删除yao-agent等待自动下线）
提前进入yao-agent容器，执行监控脚本记录GPU使用信息
python3 monitor.py

python3 main.py --lab=pre_schedule --mode=test --case=cnn-pre
python3 main.py --lab=pre_schedule --mode=test --case=cnn-default
python3 main.py --lab=pre_schedule --mode=test --case=lstm-pre
python3 main.py --lab=pre_schedule --mode=test --case=lstm-default
python3 main.py --lab=pre_schedule --mode=test --case=resnet50-pre
python3 main.py --lab=pre_schedule --mode=test --case=resnet50-default
python3 main.py --lab=pre_schedule --mode=test --case=mixed-pre
python3 main.py --lab=pre_schedule --mode=test --case=mixed-default
python3 main.py --lab=pre_schedule --mode=test --case=neumf-pre
python3 main.py --lab=pre_schedule --mode=test --case=neumf-default
在实验结束后，把残留的几个sleep作业手动停止。