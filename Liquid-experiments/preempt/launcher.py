import random


def train(jm=None, api=None, seed=2020, case=None):
	pass


def test(jm=None, api=None, seed=2020, case=1):
	cases = ["FCFS", "priority", "preempt", "batch_preempt", "reset_strategy"]
	if case not in cases:
		print('[WARN] case not in ' + str(cases))
		return
	api.conf_reset()

	conf = {}
	if case == 'FCFS':
		api.conf_update('scheduler.strategy', 'FCFS')
		conf = {
			'candidates': ['inception3_d', 'resnet50_d', 'vgg16_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 8,
		}
	elif case == 'priority':
		api.conf_update('scheduler.strategy', 'priority')
		conf = {
			'candidates': ['inception3_d', 'resnet50_d', 'vgg16_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 8,
		}
	elif case == 'preempt':
		api.conf_update('scheduler.strategy', 'priority')
		api.conf_update('scheduler.preempt_enabled', 'true')
		conf = {
			'candidates': ['inception3_d', 'resnet50_d', 'vgg16_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 8,
		}
	elif case == 'batch_preempt':
		api.conf_update('scheduler.preempt_enabled', 'true')
		conf = {
			'candidates': ['inception3_d', 'resnet50_d', 'vgg16_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 8,
		}
	elif case == 'reset_strategy':
		api.conf_update('scheduler.strategy', 'fair')
		return
	else:
		print('[ERROR] case not in ' + str(cases))
		return
	launch(conf)


def launch(conf=None):
	if conf is None:
		conf = {}

	candidates = conf['candidates']
	jm = conf['jm']
	job_num = conf['job_num']
	api = conf['api']

	import re

	random.seed(conf['seed'])
	for i in range(job_num):
		next_seed = random.randint(0, 999999)
		job_name = random.choice(candidates)
		job = jm.get_job(job_name, seed=next_seed)
		if i == 5 or i == 7:
			job['priority'] = 50
		job['tasks'] = job['tasks'].replace('--save_model_steps=0', '--save_model_steps=50')

		obj = re.search(r'--num_batches=(\w+)', job['tasks'], re.M | re.I)
		if obj:
			num_batches = int(obj.group(1))
			job['tasks'] = job['tasks'].replace(
				'--num_batches={}'.format(num_batches),
				'--num_batches={}'.format(num_batches * 2)
			)
		msg = api.submit_job(job)
		print(i, msg)
		import time
		time.sleep(5)


if __name__ == '__main__':
	print("preempt.launcher")
	pass
