import random
import time
import json


def train(jm=None, api=None, seed=2020, case=None):
	queues = ['small1', 'small2', 'small3', 'small4', 'large']
	for name in queues:
		queue = {
			'name': name,
			'weight': 10,
			'reserved': 'false',
			'quota_gpu': 0,
			'quota_gpu_mem': 10240,
			'quota_cpu': 0,
			'quota_mem': 1024
		}
		msg = api.create_queue(queue)
		print(msg)
	pass


def test(jm=None, api=None, seed=2020, case=1):
	cases = ["drf-balanced", "drf-unbalanced", "qdrf-balanced", "qdrf-unbalanced"]
	if case not in cases:
		print('[WARN] case not in ' + str(cases))
		return

	queues = ['small1', 'small2', 'small3', 'small4', 'large']

	api.conf_reset()
	if case == 'drf-balanced' or case == 'qdrf-balanced':
		for i in range(64):
			job = jm.get_job('small', seed=seed)
			job['cluster'] = queues[i % 4]
			msg = api.submit_job(job)
			print(i, msg)
	elif case == 'drf-unbalanced' or case == 'qdrf-unbalanced':
		for i in range(4):
			job = jm.get_job('small', seed=seed)
			job['cluster'] = queues[i % 4]
			msg = api.submit_job(job)
			print(i, msg)
		for i in range(2):
			job = jm.get_job('large', seed=seed)
			job['cluster'] = 'large'
			msg = api.submit_job(job)
			print(i, msg)
		for i in range(60):
			job = jm.get_job('small', seed=seed)
			job['cluster'] = queues[i % 4]
			msg = api.submit_job(job)
			print(i, msg)
	else:
		print('[ERROR] case not in ' + str(cases))
		return


if __name__ == '__main__':
	print("qdrf.launcher")
	pass
