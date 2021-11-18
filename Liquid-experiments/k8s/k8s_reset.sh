#!/bin/bash

# https://github.com/kubernetes/kubernetes/issues/39557
# https://stackoverflow.com/questions/46276796/kubernetes-cannot-cleanup-flannel
sudo kubeadm reset
sudo rm -rf /var/lib/cni/
sudo rm -rf /var/lib/kubelet/*
sudo rm -rf /etc/cni/

sudo ip link delete cni0
sudo ip link delete flannel.1