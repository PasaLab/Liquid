#!/usr/bin/env bash

while true
do
	utils=`nvidia-smi -q -x | grep gpu_util | awk '{print $1}' | awk -F'>' '{print $2}' | xargs echo`

	echo "$(date):${utils}"
	sleep 1
done