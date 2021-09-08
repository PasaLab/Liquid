# TensorFlow benchmarks
This repository contains various TensorFlow benchmarks. Currently, it consists of two projects:


1. [PerfZero](https://github.com/tensorflow/benchmarks/tree/master/perfzero): A benchmark framework for TensorFlow.

2. [scripts/tf_cnn_benchmarks](https://github.com/tensorflow/benchmarks/tree/master/scripts/tf_cnn_benchmarks) (no longer maintained): The TensorFlow CNN benchmarks contain TensorFlow 1 benchmarks for several convolutional neural networks.

If you want to run TensorFlow models and measure their performance, also consider the [TensorFlow Official Models](https://github.com/tensorflow/models/tree/master/official)

```python
python /workspace/scripts/tf_cnn_benchmarks/tf_cnn_benchmarks.py \
--num_gpus=1 \
--batch_size=32 \
--model=resnet50 \
--num_batches=200 \
--train_dir=/tmp \
--variable_update=parameter_server
```

```python
python /workspace/scripts/tf_cnn_benchmarks/tf_cnn_benchmarks.py \
--num_gpus=1 \
--batch_size=32 \
--model=resnet50 \
--num_batches=200 \
--train_dir=/tmp \
--variable_update=distributed_replicated
--ps_hosts=ps1:2222 \
--worker_hosts=worker1:2222,worker2:2222 \
--job_name=ps \
--task_index=0
```

```
--save_model_secs: How often to save trained models. Pass 0 to disable saving checkpoints every N seconds. A checkpoint is saved after training completes regardless of this
option.
(default: '0')
(an integer)
--save_model_steps: How often to save trained models. If specified, save_model_secs must not be specified.
(an integer)
```