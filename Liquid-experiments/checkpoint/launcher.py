import random


def train(jm=None, api=None, seed=2020, case=None):
	pass


def test(jm=None, api=None, seed=2020, case=1):
	cases = ["local", "distributed"]
	if case not in cases:
		print('[WARN] case not in ' + str(cases))
		return
	api.conf_reset()

	conf = {}
	if case == 'checkpoint_high':
		conf = {
			'candidates': ['cnn', 'lstm', 'resnet50', 'vgg16', 'inception3'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 1,
		}
	elif case == 'checkpoint_low':
		conf = {
			'candidates': ['resnet50_d', 'vgg16_d', 'inception3_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 1,
		}
	elif case == 'checkpoint_auto':
		conf = {
			'candidates': ['resnet50_d', 'vgg16_d', 'inception3_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 1,
		}
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

	random.seed(conf['seed'])
	for i in range(job_num):
		next_seed = random.randint(0, 999999)
		job_name = random.choice(candidates)
		job = jm.get_job(job_name, seed=next_seed)
		job['tasks'] = job['tasks'].replace('--save_model_steps=0', '--save_model_steps=50')
		msg = api.submit_job(job)
		print(i, msg)


if __name__ == '__main__':
	print("checkpoint.launcher")
	pass
