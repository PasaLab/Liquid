import random


def train(jm=None, api=None, seed=2020, case=None):
	pass


def test(jm=None, api=None, seed=2020, case=''):
	cases = [
		"single-random-resnet50", "single-random-vgg16", "single-random-inception3", "multi-random",
		"single-spread-resnet50", "single-spread-vgg16", "single-spread-inception3", "multi-spread",
		"single-pack-resnet50", "single-pack-vgg16", "single-pack-inception3", "multi-pack",
		"single-topology-aware-resnet50", "single-topology-aware-vgg16", "single-topology-aware-inception3",
		"multi-topology-aware",
		"single-k8s-default-resnet50", "single-k8s-default-vgg16", "single-k8s-default-inception3", "multi-k8s-default",
		"single-k8s-affinity-resnet50", "single-k8s-affinity-vgg16", "single-k8s-affinity-inception3",
		"multi-k8s-affinity",
		"single-k8s-hived-resnet50", "single-k8s-hived-vgg16", "single-k8s-hived-inception3",
		"single-k8s-kubeflow-resnet50", "single-k8s-kubeflow-vgg16", "single-k8s-kubeflow-inception3",
		"multi-k8s-kubeflow"
	]
	if case not in cases:
		print('[WARN] case not in ' + str(cases))
		return

	api.conf_reset()
	if case == 'single-topology-aware-resnet50':
		conf = {
			'candidates': ['resnet50_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 1,
		}

	elif case == 'single-topology-aware-vgg16':
		conf = {
			'candidates': ['vgg16_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 1,
		}

	elif case == 'single-topology-aware-inception3':
		conf = {
			'candidates': ['inception3_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 1,
		}

	elif case == 'single-random-resnet50':
		api.conf_update('allocator.strategy', 'random')
		conf = {
			'candidates': ['resnet50_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 1,
		}

	elif case == 'single-random-vgg16':
		api.conf_update('allocator.strategy', 'random')
		conf = {
			'candidates': ['vgg16_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 1,
		}

	elif case == 'single-random-inception3':
		api.conf_update('allocator.strategy', 'random')
		conf = {
			'candidates': ['inception3_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 1,
		}

	elif case == 'single-spread-resnet50':
		api.conf_update('allocator.strategy', 'spread')
		conf = {
			'candidates': ['resnet50_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 1,
		}

	elif case == 'single-spread-vgg16':
		api.conf_update('allocator.strategy', 'spread')
		conf = {
			'candidates': ['vgg16_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 1,
		}

	elif case == 'single-spread-inception3':
		api.conf_update('allocator.strategy', 'spread')
		conf = {
			'candidates': ['inception3_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 1,
		}

	elif case == 'single-pack-resnet50':
		api.conf_update('allocator.strategy', 'pack')
		conf = {
			'candidates': ['resnet50_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 1,
		}

	elif case == 'single-pack-vgg16':
		api.conf_update('allocator.strategy', 'pack')
		conf = {
			'candidates': ['vgg16_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 1,
		}

	elif case == 'single-pack-inception3':
		api.conf_update('allocator.strategy', 'pack')
		conf = {
			'candidates': ['inception3_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 1,
		}

	elif case == 'multi-random':
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

	elif case == 'single-k8s-default-resnet50':
		conf = {
			'candidates': ['resnet50_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 1,
			'runner': 'k8s',
		}

	elif case == 'single-k8s-default-vgg16':
		conf = {
			'candidates': ['vgg16_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 1,
			'runner': 'k8s',
		}

	elif case == 'single-k8s-default-inception3':
		conf = {
			'candidates': ['inception3_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 1,
			'runner': 'k8s',
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

	elif case == 'single-k8s-affinity-resnet50':
		conf = {
			'candidates': ['resnet50_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 1,
			'runner': 'k8s_affinity',
		}

	elif case == 'single-k8s-affinity-vgg16':
		conf = {
			'candidates': ['vgg16_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 1,
			'runner': 'k8s_affinity',
		}

	elif case == 'single-k8s-affinity-inception3':
		conf = {
			'candidates': ['inception3_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 1,
			'runner': 'k8s_affinity',
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

	elif case == 'single-k8s-kubeflow-resnet50':
		conf = {
			'candidates': ['resnet50_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 1,
			'runner': 'kubeflow',
		}

	elif case == 'single-k8s-kubeflow-vgg16':
		conf = {
			'candidates': ['vgg16_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 1,
			'runner': 'kubeflow',
		}

	elif case == 'single-k8s-kubeflow-inception3':
		conf = {
			'candidates': ['inception3_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 1,
			'runner': 'kubeflow',
		}
	elif case == 'single-k8s-hived-resnet50':
		conf = {
			'candidates': ['resnet50_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 1,
			'runner': 'hived',
		}

	elif case == 'single-k8s-hived-vgg16':
		conf = {
			'candidates': ['vgg16_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 1,
			'runner': 'hived',
		}

	elif case == 'single-k8s-hived-inception3':
		conf = {
			'candidates': ['inception3_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 1,
			'runner': 'hived',
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

	candidates = conf['candidates']
	jm = conf['jm']
	job_num = conf['job_num']
	api = conf['api']

	random.seed(conf['seed'])
	for i in range(job_num):
		next_seed = random.randint(0, 999999)
		job_name = random.choice(candidates)
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

			job_name = model + '-' + str(time.time()).replace('.', '-')

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
