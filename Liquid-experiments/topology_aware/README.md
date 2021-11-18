## 基于拓扑感知的调度策略实验

（说明：实验用了3个节点，使用更多节点需要增加作业数量）

#### 训练
无

#### 测试
python3 main.py --lab=topology-aware --mode=test --case=single-topology-aware-resnet50
python3 main.py --lab=topology-aware --mode=test --case=single-topology-aware-vgg16
python3 main.py --lab=topology-aware --mode=test --case=single-topology-aware-inception3

python3 main.py --lab=topology-aware --mode=test --case=single-random-resnet50
python3 main.py --lab=topology-aware --mode=test --case=single-random-vgg16
python3 main.py --lab=topology-aware --mode=test --case=single-random-inception3

python3 main.py --lab=topology-aware --mode=test --case=single-spread-resnet50
python3 main.py --lab=topology-aware --mode=test --case=single-spread-vgg16
python3 main.py --lab=topology-aware --mode=test --case=single-spread-inception3

python3 main.py --lab=topology-aware --mode=test --case=single-pack-resnet50
python3 main.py --lab=topology-aware --mode=test --case=single-pack-vgg16
python3 main.py --lab=topology-aware --mode=test --case=single-pack-inception3

python3 main.py --lab=topology-aware --mode=test --case=single-k8s-default-resnet50
python3 main.py --lab=topology-aware --mode=test --case=single-k8s-default-vgg16
python3 main.py --lab=topology-aware --mode=test --case=single-k8s-default-inception3

python3 main.py --lab=topology-aware --mode=test --case=single-k8s-affinity-resnet50
python3 main.py --lab=topology-aware --mode=test --case=single-k8s-affinity-vgg16
python3 main.py --lab=topology-aware --mode=test --case=single-k8s-affinity-inception3

python3 main.py --lab=topology-aware --mode=test --case=single-k8s-kubeflow-resnet50
python3 main.py --lab=topology-aware --mode=test --case=single-k8s-kubeflow-vgg16
python3 main.py --lab=topology-aware --mode=test --case=single-k8s-kubeflow-inception3

python3 main.py --lab=topology-aware --mode=test --case=single-k8s-hived-resnet50
python3 main.py --lab=topology-aware --mode=test --case=single-k8s-hived-vgg16
python3 main.py --lab=topology-aware --mode=test --case=single-k8s-hived-inception3

python3 main.py --lab=topology-aware --mode=test --case=multi-topology-aware
python3 main.py --lab=topology-aware --mode=test --case=multi-random
python3 main.py --lab=topology-aware --mode=test --case=multi-spread
python3 main.py --lab=topology-aware --mode=test --case=multi-pack
python3 main.py --lab=topology-aware --mode=test --case=multi-k8s-default
python3 main.py --lab=topology-aware --mode=test --case=multi-k8s-affinity
python3 main.py --lab=topology-aware --mode=test --case=multi-k8s-kubeflow
