import random


def train(jm=None, api=None, seed=2020, case=None):
	pass


def test(jm=None, api=None, seed=2020, case=1):
	cases = [
		"resnet50-pww", "resnet50-p:w:w", "resnet50-pw:w", "resnet50-p:ww",
		"vgg16-pww", "vgg16-p:w:w", "vgg16-pw:w", "vgg16-p:ww",
		"inception3-pww", "inception3-p:w:w", "inception3-pw:w", "inception3-p:ww",
	]
	if case not in cases:
		print('[WARN] case not in ' + str(cases))
		return
	api.conf_reset()

	conf = {}
	if case == 'vgg16-pww':
		conf = {
			'candidates': ['vgg16_d'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 1,
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
	print("placement.launcher")
	pass
