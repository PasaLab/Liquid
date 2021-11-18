

## kfctl

```bash
# https://github.com/kubeflow/kfctl/releases
wget -O kfctl_linux.tar.gz https://github.com/kubeflow/kfctl/releases/download/v1.1.0/kfctl_v1.1.0-0-g9a3621e_linux.tar.gz

tar -xvf kfctl_linux.tar.gz

sudo mv ./kfctl /usr/local/bin/kfctl

rm kfctl_linux.tar.gz
```

```bash
# https://www.kubeflow.org/docs/started/getting-started/
# kfctl_k8s_istio.v1.0.2.yaml
wget -O kfctl_k8s_istio.yaml https://raw.githubusercontent.com/kubeflow/manifests/v1.0-branch/kfdef/kfctl_k8s_istio.v1.0.2.yaml
```

```bash
wget -O v1.0.2.tar.gz https://github.com/kubeflow/manifests/archive/v1.0.2.tar.gz
```


> If no Internet connection,
> change `https://github.com/kubeflow/manifests/archive/v1.0.2.tar.gz` 
> to `file:/home/newnius/v1.0.2.tar.gz` in file `kfctl_k8s_istio.yaml`
> [kfctl - How to access GitHub behind corporate proxy?](https://github.com/kubeflow/kubeflow/issues/4753)

```bash
kfctl build -V -f kfctl_k8s_istio.yaml
```

```bash
kfctl apply -V -f kfctl_k8s_istio.v1.0.2.yaml
```

mysql等需要存储，创建本地存储

```bash
cat << EOF > local_pv.yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: pipeline-mysql-pv
  namespace: kubeflow
  labels:
    type: local
    app: pipeline-mysql-pv
    key: kubeflow-pv
spec:
  capacity:
    storage: 20Gi
  accessModes:
    - ReadWriteOnce
  hostPath:
    path: /data/pipeline-mysql
    type: DirectoryOrCreate
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: pipeline-minio-pv
  namespace: kubeflow
  labels:
    type: local
    app: pipeline-minio-pv
    key: kubeflow-pv
spec:
  capacity:
    storage: 20Gi
  accessModes:
    - ReadWriteOnce
  hostPath:
    path: /data/pipeline-minio
    type: DirectoryOrCreate
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: katib-mysql
  namespace: kubeflow
  labels:
    type: local
    app: katib-mysql
spec:
  capacity:
    storage: 20Gi
  accessModes:
    - ReadWriteOnce
  hostPath:
    path: /data/katib-mysql
    type: DirectoryOrCreate
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: metadata-mysql-pv
  namespace: kubeflow
  labels:
    type: local
    app: metadata-mysql-pv
    key: kubeflow-pv
spec:
  capacity:
    storage: 20Gi
  accessModes:
    - ReadWriteOnce
  hostPath:
    path: /data/metadata-mysql
    type: DirectoryOrCreate
EOF
kubectl create -f local_pv.yaml
```



## View

```bash
kubectl get tfjob -n kubeflow

kubectl get pod -n kubeflow
```

以环境变量形式传入 `ps_hosts`, `job_name`等变量

```bash
$ kubectl get pod/mnist-train-dist-worker-0 -n kubeflow -oyaml

...
- name: TF_CONFIG
      value: '{"cluster":{
	"chief":["mnist-train-dist-chief-0.kubeflow.svc:2222"],
	"ps":["mnist-train-dist-ps-0.kubeflow.svc:2222"],
	"worker":["mnist-train-dist-worker-0.kubeflow.svc:2222","mnist-train-dist-worker-1.kubeflow.svc:2222"]},
	"task":{"type":"worker","index":0},
	"environment":"cloud"}'
...
```


## Ref

[kubeflow1.0.2部署](https://my.oschina.net/u/3825598/blog/4276681)

[Kubeflow 1.0 上线： 体验生产级的机器学习平台](https://developer.aliyun.com/article/758776)

[完整 Kubeflow 使用教學 — 開發 ML 模型、進行分散式訓練與部署服務](https://medium.com/infuseai/%E5%AE%8C%E6%95%B4-kubeflow-%E4%BD%BF%E7%94%A8%E6%95%99%E5%AD%B8-%E9%96%8B%E7%99%BC-ml-%E6%A8%A1%E5%9E%8B-%E9%80%B2%E8%A1%8C%E5%88%86%E6%95%A3%E5%BC%8F%E8%A8%93%E7%B7%B4%E8%88%87%E9%83%A8%E7%BD%B2%E6%9C%8D%E5%8B%99-ca1348b1cb8b)