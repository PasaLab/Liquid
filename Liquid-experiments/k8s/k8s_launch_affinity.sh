#!/bin/bash

f(){
	echo "Format: JOB_NAME=<JOB_NAME> MODEL=<MODEL> BATCH_SIZE=<BATCH_SIZE> NUM_BATCHES=<NUM_BATCHES> $0"
	echo -e "\n"
}

if [[ -z "${JOB_NAME}" ]]; then
	echo "[ERROR] JOB_NAME not specified."
	f
	JOB_NAME='resnet50-1'
	exit 1
fi
if [[ -z "${MODEL}" ]]; then
	echo "[ERROR] MODEL not specified."
	f
	MODEL='resnet50'
	exit 1
fi
if [[ -z "${BATCH_SIZE}" ]]; then
	echo "[ERROR] BATCH_SIZE not specified."
	f
	BATCH_SIZE=32
	exit 1
fi
if [[ -z "${NUM_BATCHES}" ]]; then
	echo "[ERROR] NUM_BATCHES not specified."
	f
	NUM_BATCHES=200
	exit 1
fi


if [[ -d "/dfs/yao-jobs/${JOB_NAME}" ]]
then
    echo "Remove directory /dfs/yao-jobs/${JOB_NAME} first." 
    exit;
fi

mkdir -p /dfs/yao-jobs/${JOB_NAME}/ps1/
mkdir -p /dfs/yao-jobs/${JOB_NAME}/worker1/
mkdir -p /dfs/yao-jobs/${JOB_NAME}/worker2/

echo "Job ${JOB_NAME} submitted at $(date)"


cp /dfs/k8s/job-template_affinity_namespace.yaml            tmp-${JOB_NAME}_namespace.yaml
sed -i "s/JOB_NAMESPACE/${JOB_NAME}/g"                      tmp-${JOB_NAME}_namespace.yaml
kubectl apply -f tmp-${JOB_NAME}_namespace.yaml
rm tmp-${JOB_NAME}_namespace.yaml


cp /dfs/k8s/job-template_affinity_ps1.yaml                  tmp-${JOB_NAME}_ps1.yaml
sed -i "s/JOB_NAMESPACE/${JOB_NAME}/g"                      tmp-${JOB_NAME}_ps1.yaml
sed -i "s/JOB_NAME/${JOB_NAME}/g"                           tmp-${JOB_NAME}_ps1.yaml
sed -i "s/--model=resnet50/--model=${MODEL}/g"              tmp-${JOB_NAME}_ps1.yaml
sed -i "s/--batch_size=32/--batch_size=${BATCH_SIZE}/g"     tmp-${JOB_NAME}_ps1.yaml
sed -i "s/--num_batches=200/--num_batches=${NUM_BATCHES}/g" tmp-${JOB_NAME}_ps1.yaml
kubectl apply -f tmp-${JOB_NAME}_ps1.yaml
rm tmp-${JOB_NAME}_ps1.yaml
sleep 1


cp /dfs/k8s/job-template_affinity_worker1.yaml              tmp-${JOB_NAME}_worker1.yaml
sed -i "s/JOB_NAMESPACE/${JOB_NAME}/g"                      tmp-${JOB_NAME}_worker1.yaml
sed -i "s/JOB_NAME/${JOB_NAME}/g"                           tmp-${JOB_NAME}_worker1.yaml
sed -i "s/--model=resnet50/--model=${MODEL}/g"              tmp-${JOB_NAME}_worker1.yaml
sed -i "s/--batch_size=32/--batch_size=${BATCH_SIZE}/g"     tmp-${JOB_NAME}_worker1.yaml
sed -i "s/--num_batches=200/--num_batches=${NUM_BATCHES}/g" tmp-${JOB_NAME}_worker1.yaml
kubectl apply -f tmp-${JOB_NAME}_worker1.yaml
rm tmp-${JOB_NAME}_worker1.yaml
sleep 1

cp /dfs/k8s/job-template_affinity_worker2.yaml              tmp-${JOB_NAME}_worker2.yaml
sed -i "s/JOB_NAMESPACE/${JOB_NAME}/g"                      tmp-${JOB_NAME}_worker2.yaml
sed -i "s/JOB_NAME/${JOB_NAME}/g"                           tmp-${JOB_NAME}_worker2.yaml
sed -i "s/--model=resnet50/--model=${MODEL}/g"              tmp-${JOB_NAME}_worker2.yaml
sed -i "s/--batch_size=32/--batch_size=${BATCH_SIZE}/g"     tmp-${JOB_NAME}_worker2.yaml
sed -i "s/--num_batches=200/--num_batches=${NUM_BATCHES}/g" tmp-${JOB_NAME}_worker2.yaml
kubectl apply -f tmp-${JOB_NAME}_worker2.yaml
rm tmp-${JOB_NAME}_worker2.yaml

while true; do
  if [[ $(kubectl get pods -n ${JOB_NAME} | grep Running) ]]; then
      echo ${JOB_NAME} 'launched, wait for completion' $(date)
      break
  fi
  sleep 3
done

while true; do
  if [[ $(kubectl get pods -n ${JOB_NAME} | grep Completed) ]]; then
      echo "Some pods completed, stop the jobs"
      break
  fi
  sleep 3
done

kubectl delete --namespace ${JOB_NAME} job ps1 worker1 worker2

kubectl delete --namespace ${JOB_NAME} service ps1 worker1 worker2

kubectl delete namespace ${JOB_NAME}


echo "Job ${JOB_NAME} completed at $(date)"
