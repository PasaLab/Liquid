#!/bin/bash

token=`sudo kubeadm token generate`

sudo kubeadm token create ${token} --print-join-command