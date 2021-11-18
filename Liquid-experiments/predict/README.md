## 作业资源需求向量评估

#### 背景

#### 建模

#### 方法

#### 先提交若干作业训练资源预估模型
分为两类，一个是单机模式的作业，包括cnn、neumf、lstm、resnet50、vgg16、inception3；
另一个是分布式作业，包括resnet50、vgg16、inception13。
可变参数包括GPU个数、batch大小、数据规模、ps与worker个数，
预估的资源包括GPU利用率、显存占用、CPU利用率、内存占用、带宽占用；

python3 main.py --lab=predict --mode=train --case=local
python3 main.py --lab=predict --mode=train --case=distributed

#### 效果评估

分为两个，一是模型的效果查看，
在优化器容器内，复制任意一个作业的数据，复制为dataset.csv，移动到compare.sh脚本同目录，
然后执行compare.sh脚本，脚本随机打乱记录，然后划分训练集和测试集，输出在lr、rf、dt、ada、gbdt五种
模型下r2指标和rmse指标随训练样本的增加而变化的数据。

另外一个是实际的效果，
在提交作业的界面上，通过修改启动命令里的参数值，然后predict，即可自动修改资源需求。
