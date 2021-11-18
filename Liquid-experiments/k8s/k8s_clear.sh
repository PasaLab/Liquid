#!/bin/bash

kubectl get ns | grep -v 'default' | grep -v 'kube-' | awk 'NR>1' | awk '{print $1}' | xargs kubectl delete ns

pkill k8s_launch.sh