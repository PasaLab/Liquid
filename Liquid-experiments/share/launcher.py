import random
import time
import json


def train(jm=None, api=None, seed=2020, case=None):
	api.conf_reset()
	api.conf_update('pool.share.enable_threshold', 1.75)

	conf = {
		'candidates': ['cnn'],
		'jm': jm,
		'api': api,
		'seed': seed,
		'job_num': 200,
	}
	launch(conf)


def test(jm=None, api=None, seed=2020, case=1):
	cases = ["cnn-share", "cnn-exclusive", "neumf-share", "neumf-exclusive", "mixed-share", "mixed-exclusive"]
	if case not in cases:
		print('[WARN] case not in ' + str(cases))
		return

	api.conf_reset()
	if case == 'cnn-share' or case == 'neumf-share' or case == 'mixed-share':
		api.conf_update('pool.share.enable_threshold', 0.75)
	else:
		api.conf_update('pool.share.enable_threshold', 1.75)
	conf = {}
	if case == 'cnn-share' or case == 'cnn-exclusive':
		conf = {
			'candidates': ['cnn'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 8,
		}
	elif case == 'neumf-share' or case == 'neumf-exclusive':
		conf = {
			'candidates': ['neumf'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 8,
		}
	elif case == 'mixed-share' or case == 'mixed-exclusive':
		conf = {
			'candidates': ['cnn', 'neumf'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 8,
		}
	else:
		print('[ERROR] case not in ' + str(cases))
		return

	job = jm.get_job('sleep', seed=0)
	for i in range(1):
		msg = api.submit_job(job)
		print(i, msg)
	time.sleep(3)
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
		while True:
			next_seed = random.randint(0, 999999)
			job_name = random.choice(candidates)
			job = jm.get_job(job_name, seed=next_seed)
			tasks = json.loads(job['tasks'])
			if len(tasks) == 1 and tasks[0]['gpu_number'] == '1':
				break
		msg = api.submit_job(job)
		print(i, msg)
		time.sleep(5)


if __name__ == '__main__':
	print("share.launcher")
	pass
