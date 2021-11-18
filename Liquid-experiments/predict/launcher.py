import random


def train(jm=None, api=None, seed=2020, case=None):
	cases = ["local", "distributed"]
	if case not in cases:
		print('[WARN] case not in ' + str(cases))
		return
	api.conf_reset()

	conf = {}
	if case == 'local':
		conf = {
			'candidates': ['cnn', 'lstm', 'resnet50', 'vgg16', 'inception3'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 150,
		}
	elif case == 'distributed':
		conf = {
			'candidates': ['resnet50_d', 'vgg16_d', 'inception3_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 150,
		}
	else:
		print('[ERROR] case not in ' + str(cases))
		return
	launch(conf)


def test(jm=None, api=None, seed=2020, case=1):
	pass


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
		msg = api.submit_job(job)
		print(i, msg)


if __name__ == '__main__':
	print("predict.launcher")
	pass
