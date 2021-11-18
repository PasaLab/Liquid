import json
import random


def convert_job_status(status_code):
	status_map = [
		'Created',  # 0
		'Starting',  # 1
		'Running',  # 2
		'Stopped',  # 3
		'Finished',  # 4
		'Failed',  # 5
	]
	if 0 <= status_code < len(status_map):
		return status_map[status_code]
	return 'Unknown'


def get_job(job_name, seed=2020):
	random.seed(seed)
	job = {}
	if job_name == 'sleep':
		job = {
			'name': 'sleep',
			'workspace': 'http://code.pasalab.jluapp.com/newnius/yao-job-benchmarks.git',
			'cluster': 'default',
			'priority': '25',
			'run_before': '',
			'locality': '0',
			'tasks': json.dumps([{
				"name": "node1",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": "sleep infinity",
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}]),
		}
	elif job_name == 'small':
		cmd = "sleep 300"
		job = {
			'name': job_name,
			'workspace': 'http://code.pasalab.jluapp.com/newnius/yao-job-CNN.git',
			'cluster': 'default',
			'priority': '25',
			'run_before': '',
			'locality': '0',
			'tasks': json.dumps([{
				"name": "node1",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": cmd,
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}, {
				"name": "node2",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": cmd,
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}]),
		}
	elif job_name == 'large':
		cmd = "sleep 600"
		job = {
			'name': job_name,
			'workspace': 'http://code.pasalab.jluapp.com/newnius/yao-job-CNN.git',
			'cluster': 'default',
			'priority': '25',
			'run_before': '',
			'locality': '0',
			'tasks': json.dumps([{
				"name": "node1",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": cmd,
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}, {
				"name": "node2",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": cmd,
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}, {
				"name": "node3",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": cmd,
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}, {
				"name": "node4",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": cmd,
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}, {
				"name": "node5",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": cmd,
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}, {
				"name": "node6",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": cmd,
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}, {
				"name": "node7",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": cmd,
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}, {
				"name": "node8",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": cmd,
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}, {
				"name": "node9",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": cmd,
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}, {
				"name": "node10",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": cmd,
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}, {
				"name": "node11",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": cmd,
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}, {
				"name": "node12",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": cmd,
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}, {
				"name": "node13",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": cmd,
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}, {
				"name": "node14",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": cmd,
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}, {
				"name": "node15",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": cmd,
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}, {
				"name": "node16",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": cmd,
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}]),
		}
	elif job_name == 'job1':
		job = {
			'name': 'job1',
			'workspace': 'http://code.pasalab.jluapp.com/newnius/yao-job-benchmarks.git',
			'cluster': 'default',
			'priority': '25',
			'run_before': '',
			'locality': '0',
			'tasks': json.dumps([{
				"name": "node1",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": "sleep infinity",
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}]),
		}
	elif job_name == 'job2':
		job = {
			'name': 'job2',
			'workspace': 'http://code.pasalab.jluapp.com/newnius/yao-job-benchmarks.git',
			'cluster': 'default',
			'priority': '25',
			'run_before': '',
			'locality': '0',
			'tasks': json.dumps([{
				"name": "node1",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": "sleep infinity",
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}, {
				"name": "node2",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": "sleep infinity",
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}]),
		}
	elif job_name == 'job3':
		job = {
			'name': 'job3',
			'workspace': 'http://code.pasalab.jluapp.com/newnius/yao-job-benchmarks.git',
			'cluster': 'default',
			'priority': '25',
			'run_before': '',
			'locality': '0',
			'tasks': json.dumps([{
				"name": "node1",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": "sleep infinity",
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "2",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}, {
				"name": "node2",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": "sleep infinity",
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "2",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}, {
				"name": "node3",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": "sleep infinity",
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "2",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}]),
		}
	elif job_name == 'job5':
		job = {
			'name': 'job5',
			'workspace': 'http://code.pasalab.jluapp.com/newnius/yao-job-benchmarks.git',
			'cluster': 'default',
			'priority': '25',
			'run_before': '',
			'locality': '0',
			'tasks': json.dumps([{
				"name": "node1",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": "sleep infinity",
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}, {
				"name": "node2",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": "sleep infinity",
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}, {
				"name": "node3",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": "sleep infinity",
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}, {
				"name": "node4",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": "sleep infinity",
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}, {
				"name": "node5",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": "sleep infinity",
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}]),
		}
	elif job_name == 'job10':
		job = {
			'name': 'job10',
			'workspace': 'http://code.pasalab.jluapp.com/newnius/yao-job-benchmarks.git',
			'cluster': 'default',
			'priority': '25',
			'run_before': '',
			'locality': '0',
			'tasks': json.dumps([{
				"name": "node1",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": "sleep infinity",
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}, {
				"name": "node2",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": "sleep infinity",
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}, {
				"name": "node3",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": "sleep infinity",
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}, {
				"name": "node4",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": "sleep infinity",
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}, {
				"name": "node5",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": "sleep infinity",
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}, {
				"name": "node6",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": "sleep infinity",
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}, {
				"name": "node7",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": "sleep infinity",
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}, {
				"name": "node8",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": "sleep infinity",
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}, {
				"name": "node9",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": "sleep infinity",
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}, {
				"name": "node10",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": "sleep infinity",
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}]),
		}
	elif job_name == 'cnn':
		cmd = "PYTHONPATH=\"$PYTHONPATH:/workspace\""
		cmd += " python /workspace/official/r1/mnist/mnist.py"
		cmd += " --data_dir=/workspace/data/"
		job = {
			'name': job_name,
			'workspace': 'http://code.pasalab.jluapp.com/newnius/yao-job-CNN.git',
			'cluster': 'default',
			'priority': '25',
			'run_before': '',
			'locality': '0',
			'tasks': json.dumps([{
				"name": "node1",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": cmd,
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "4096",
				"is_ps": "0",
				"gpu_model": "t4",
			}]),
		}
	elif job_name == 'lstm':
		cmd = "PYTHONPATH=\"$PYTHONPATH:/workspace\""
		cmd += " python3 /workspace/official/staging/shakespeare/shakespeare_main.py"
		cmd += " --training_data=/workspace/data/shakespeare.txt"
		job = {
			'name': job_name,
			'workspace': 'http://code.pasalab.jluapp.com/newnius/yao-job-LSTM.git',
			'cluster': 'default',
			'priority': '25',
			'run_before': '',
			'locality': '0',
			'tasks': json.dumps([{
				"name": "node1",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:2.1-gpu",
				"cmd": cmd,
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}]),
		}
	elif job_name == 'neumf':
		batch_size = 1000
		# batch_size = random.randint(1, 3) * 1000
		cmd = "PYTHONPATH=\"$PYTHONPATH:/workspace\""
		cmd += " python /workspace/official/recommendation/ncf_keras_main.py"
		cmd += " --batch_size=" + str(batch_size)
		cmd += " --data_dir=/workspace/data/"
		cmd += "  --dataset=ml-20m"
		job = {
			'name': 'neumf',
			'workspace': 'http://code.pasalab.jluapp.com/newnius/yao-job-NeuMF.git',
			'cluster': 'default',
			'priority': '25',
			'run_before': '',
			'locality': '0',
			'tasks': json.dumps([{
				"name": "node1",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:2.1-gpu",
				"cmd": cmd,
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": "1",
				"gpu_memory": "4096",
				"is_ps": "0",
				"gpu_model": "t4",
			}]),
		}
	elif job_name == 'resnet50' or job_name == 'vgg16' or job_name == 'inception3':
		num_gpus = random.randint(1, 2)
		#num_gpus = 3
		batch_size = random.randint(1, 8) * 4
		num_batches = random.randint(4, 20) * 50
		if job_name == 'vgg16':
			batch_size = random.randint(4, 8) * 4
			num_batches = random.randint(4, 10) * 50
		cmd = "python /workspace/scripts/tf_cnn_benchmarks/tf_cnn_benchmarks.py"
		cmd += " --model=" + job_name
		cmd += " --num_gpus=" + str(num_gpus)
		cmd += " --batch_size=" + str(batch_size)
		cmd += " --num_batches=" + str(num_batches)
		cmd += " --train_dir=/tmp"
		cmd += " --variable_update=parameter_server"
		cmd += " --save_model_steps=0"
		job = {
			'name': job_name,
			'workspace': 'http://code.pasalab.jluapp.com/newnius/yao-job-benchmarks.git',
			'cluster': 'default',
			'priority': '25',
			'run_before': '',
			'locality': '0',
			'tasks': json.dumps([{
				"name": "node1",
				"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
				"cmd": cmd,
				"cpu_number": "4",
				"memory": "4096",
				"gpu_number": str(num_gpus),
				"gpu_memory": "8192",
				"is_ps": "0",
				"gpu_model": "t4",
			}]),
		}
	elif job_name == 'resnet50_d' or job_name == 'vgg16_d' or job_name == 'inception3_d':
		model = job_name.split('_')[0]
		batch_size = random.randint(1, 8) * 4
		num_batches = random.randint(4, 20) * 50
		cmd = "python /workspace/scripts/tf_cnn_benchmarks/tf_cnn_benchmarks.py"
		cmd += " --model=" + model
		cmd += " --num_gpus=1"
		cmd += " --batch_size=" + str(batch_size)
		cmd += " --num_batches=" + str(num_batches)
		cmd += " --train_dir=/tmp"
		cmd += " --variable_update=distributed_replicated"
		cmd += " --ps_hosts=ps1:2222"
		cmd += " --worker_hosts=worker1:2222,worker2:2222"
		cmd += " --save_model_steps=0"
		job = {
			'name': job_name,
			'workspace': 'http://code.pasalab.jluapp.com/newnius/yao-job-benchmarks.git',
			'cluster': 'default',
			'priority': '25',
			'run_before': '',
			'locality': '0',
			'tasks': json.dumps([
				{
					"name": "worker1",
					"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
					"cmd": cmd + " --job_name=worker" + " --task_index=0",
					"cpu_number": "4",
					"memory": "4096",
					"gpu_number": "1",
					"gpu_memory": "8192",
					"is_ps": "0",
					"gpu_model": "t4",
				},
				{
					"name": "worker2",
					"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
					"cmd": cmd + " --job_name=worker" + " --task_index=1",
					"cpu_number": "4",
					"memory": "4096",
					"gpu_number": "1",
					"gpu_memory": "8192",
					"is_ps": "0",
					"gpu_model": "t4",
				},
				{
					"name": "ps1",
					"image": "registry.cn-beijing.aliyuncs.com/quickdeploy0/yao-tensorflow:1.14-gpu",
					"cmd": cmd + " --job_name=ps" + " --task_index=0",
					"cpu_number": "4",
					"memory": "4096",
					"gpu_number": "1",
					"gpu_memory": "8192",
					"is_ps": "1",
					"gpu_model": "t4",
				},
			]),
		}
	else:
		print("[WARN] job {} not exist".format(job_name))
	return job
