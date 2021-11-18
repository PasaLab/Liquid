import random


def train(jm=None, api=None, seed=2020, case=None):
	pass


def test(jm=None, api=None, seed=2020, case=''):
	cases = [
		"multi-random", "multi-spread", "multi-pack", "multi-topology-aware", "multi-batch",
		"multi-k8s-default", "multi-k8s-affinity", "multi-k8s-kubeflow","multi-k8s-hived"
	]
	if case not in cases:
		print('[WARN] case not in ' + str(cases))
		return

	api.conf_reset()
	if case == 'multi-random':
		api.conf_update('allocator.strategy', 'random')
		conf = {
			'candidates': ['resnet50_d', 'vgg16_d', 'inception3_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 8,
		}

	elif case == 'multi-spread':
		api.conf_update('allocator.strategy', 'spread')
		conf = {
			'candidates': ['resnet50_d', 'vgg16_d', 'inception3_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 8,
		}

	elif case == 'multi-pack':
		api.conf_update('allocator.strategy', 'pack')
		conf = {
			'candidates': ['resnet50_d', 'vgg16_d', 'inception3_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 8,
		}

	elif case == 'multi-topology-aware':
		conf = {
			'candidates': ['resnet50_d', 'vgg16_d', 'inception3_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 8,
		}

	elif case == 'multi-batch':
		api.conf_update('pool.batch.enabled', 'true')
		api.conf_update('pool.batch.interval', '5')
		api.conf_update('scheduler.parallelism', '5')
		conf = {
			'candidates': ['resnet50_d', 'vgg16_d', 'inception3_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 8,
		}

	elif case == 'multi-k8s-default':
		conf = {
			'candidates': ['resnet50_d', 'vgg16_d', 'inception3_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 8,
			'runner': 'k8s',
		}

	elif case == 'multi-k8s-affinity':
		conf = {
			'candidates': ['resnet50_d', 'vgg16_d', 'inception3_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 8,
			'runner': 'k8s_affinity',
		}

	elif case == 'multi-k8s-kubeflow':
		conf = {
			'candidates': ['resnet50_d', 'vgg16_d', 'inception3_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 8,
			'runner': 'kubeflow',
		}
	elif case == 'multi-k8s-hived':
		conf = {
			'candidates': ['resnet50_d', 'vgg16_d', 'inception3_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 8,
			'runner': 'hived',
		}
	else:  # multi topology-aware
		conf = {
			'candidates': ['resnet50_d', 'vgg16_d', 'inception3_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 8,
		}

	launch(conf)


def launch(conf=None):
	if conf is None:
		conf = {}

	jobs = [
		'inception3_d', 'resnet50_d', 'resnet50_d', 'vgg16_d',
		'resnet50_d', 'inception3_d', 'resnet50_d', 'vgg16_d'
	]
	jm = conf['jm']
	job_num = conf['job_num']
	api = conf['api']

	random.seed(conf['seed'])
	for i in range(job_num):
		next_seed = random.randint(0, 999999)
		job_name = jobs[i]
		job = jm.get_job(job_name, seed=next_seed)
		if 'runner' not in conf:
			msg = api.submit_job(job)
			print(i, msg)

		else:
			import re
			import time

			model = 'resnet50'
			batch_size = 32
			num_batches = 200

			obj = re.search(r'--model=(\w+)', job['tasks'], re.M | re.I)
			if obj:
				model = obj.group(1)

			obj = re.search(r'--batch_size=(\d+)', job['tasks'], re.M | re.I)
			if obj:
				batch_size = int(obj.group(1))

			obj = re.search(r'--num_batches=(\d+)', job['tasks'], re.M | re.I)
			if obj:
				num_batches = int(obj.group(1))

			timestamp = time.time()
			job_name = model + '-' + str(timestamp).replace('.', '-')

			time.sleep(1)

			script = './k8s_launch.sh'
			if conf['runner'] == 'k8s':
				script = './k8s_launch.sh'
			elif conf['runner'] == 'k8s_affinity':
				script = './k8s_launch_affinity.sh'
			elif conf['runner'] == 'kubeflow':
				script = './k8s_launch_kubeflow.sh'
			elif conf['runner'] == 'hived':
				script = './k8s_launch_hived.sh'
			print("JOB_NAME={} MODEL={} BATCH_SIZE={} NUM_BATCHES={} nohup {} &".format(
				job_name, model, batch_size, num_batches, script))
			print("sleep 3")


if __name__ == '__main__':
	print("topology_aware.launcher")
	pass
