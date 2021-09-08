import requests
import time
import os
import json

'''Configuration'''
BASE_URL = 'http://yao.example.com'
''''''

sess = requests.Session()
sess.headers.update({'Referer': BASE_URL})

status_map = [
	'Created',  # 0
	'Starting',  # 1
	'Running',  # 2
	'Stopped',  # 3
	'Finished',  # 4
	'Failed',  # 5
]


def login(user='', pwd=''):
	# Get CSRF Token
	r = sess.get(BASE_URL)
	# print(r.content)

	# Login
	url = BASE_URL + '/service?action=user_login'
	r = sess.post(url, data={})
	# print(r.content)
	return


def get_sys_status():
	# Retrieve Status
	r = sess.get(BASE_URL + '/service?action=summary_get')
	print(r.content)
	# b'{"jobs":{"finished":1,"running":0,"pending":0},"gpu":{"free":20,"using":0},"errno":0,"msg":"Success !"}'

	# Get pool Util history
	r = sess.get(BASE_URL + '/service?action=summary_get_pool_history')
	# print(r.content)
	return


def submit_job(job):
	r = sess.post(BASE_URL + '/service?action=job_submit', data=job)
	data = str(r.content, 'utf-8')
	msg = json.loads(data)
	print(msg)
	return msg


def job_list():
	print("\nList of jobs:")
	r = sess.get(BASE_URL + '/service?action=job_list&who=self&sort=nobody&order=desc&offset=0&limit=10')
	data = str(r.content, 'utf-8')
	msg = json.loads(data)

	if len(msg['jobs']) > 0:
		for job in msg['jobs']:
			name = job['name']
			status = status_map[job['status']]
			print("Status of job: {} is {}".format(name, status))
	print("\n")


def job_status(job_name):
	r = sess.get(BASE_URL + '/service?action=job_status&name=' + job_name)
	data = str(r.content, 'utf-8')
	msg = json.loads(data)
	print("Status of tasks in {}:".format(job_name))
	if len(msg['tasks']) > 0:
		for task in msg['tasks']:
			print("{} ({}) is/was running on {}".format(task['hostname'], task['status'], task['node']))


if __name__ == '__main__':
	os.environ["TZ"] = 'Asia/Shanghai'
	if hasattr(time, 'tzset'):
		time.tzset()

	login()
	get_sys_status()

	tasks = [{
		"name": "node1",
		"image": "quickdeploy/yao-tensorflow:1.14-gpu",
		"cmd": "python /workspace/scripts/tf_cnn_benchmarks/tf_cnn_benchmarks.py \\ --num_gpus=1 \\ --batch_size=32 \\ --model=resnet50 \\ --num_batches=200 \\ --train_dir=/tmp \\ --variable_update=parameter_server",
		"cpu_number": "4",
		"memory": "4096",
		"gpu_number": "1",
		"gpu_memory": "8192",
		"is_ps": "0",
		"gpu_model": "k80",
	}]
	# print(json.dumps(tasks))
	job = {
		'name': 'test',
		'workspace': 'https://github.com/tensorflow/benchmarks.git',
		'cluster': 'default',
		'priority': '25',
		'run_before': '',
		'locality': '0',
		'tasks': json.dumps(tasks),
	}

	msg = submit_job(job)
	if msg['errno'] == 0:
		job_status(msg['job_name'])

	job_list()

