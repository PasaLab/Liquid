Setup k8s cluster (CentOS)

# On Master

## Init
```bash
# `--pod-network-cidr` is required
# see https://github.com/coreos/flannel/issues/728
sudo kubeadm init --pod-network-cidr "10.244.0.0/16"
```


## Copy conf files
```bash
sudo rm -rf $HOME/.kube

mkdir -p $HOME/.kube

sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config

sudo chown $(id -u):$(id -g) $HOME/.kube/config
```


## Create a Pod network

```bash
curl -o kube-flannel.yml https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml
kubectl apply -f kube-flannel.yml
```

## Setup nvidia gpu plugin
[ref](https://kubernetes.io/docs/tasks/manage-gpus/scheduling-gpus/#official-nvidia-gpu-device-plugin)

```bash
kubectl create -f https://raw.githubusercontent.com/NVIDIA/k8s-device-plugin/1.0.0-beta4/nvidia-device-plugin.yml
```


```bas
cat <<EOF > /etc/docker/daemon.json
{
    "default-runtime": "nvidia",
    "runtimes": {
        "nvidia": {
            "path": "/usr/bin/nvidia-container-runtime",
            "runtimeArgs": []
        }
    }
}
EOF

systemctl restart docker.service

docker run nvidia/cuda:10.2-runtime-ubuntu18.04 nvidia-smi
```

# On other nodes

## Setup nvidia-docker

[manual install nvidia-container-runtime](https://github.com/NVIDIA/k8s-device-plugin/issues/166)

```bash
yum install nvidia-container-runtime
```

## Add nodes to cluster

```bash
kubeadm join 192.168.100.104:6443 --token 9bok6x.80iyvjrgzlivhwri \
    --discovery-token-ca-cert-hash sha256:002776e203046b08e0d76d9be84ea6257768db60fe0f070dbab4744c6412a3a0
```

## Auto Completion

```bash
source <(kubectl completion bash)
echo "source <(kubectl completion bash)" >> ~/.bashrc
```

# Revert

sudo kubeadm reset

# Ref
[炼丹师的工程修养之四： TensorFlow的分布式训练和K8S](https://zhuanlan.zhihu.com/p/56699786)

[基于KubeAdm快速搭建K8s集群](https://www.jianshu.com/p/06d487aea2c1)

[kubernetes系列03—kubeadm安装部署K8S集群](https://www.cnblogs.com/along21/p/10303495.html)


