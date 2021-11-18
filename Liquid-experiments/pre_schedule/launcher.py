import random
import json
import time


def train(jm=None, api=None, seed=2020, case=None):
	api.conf_reset()
	api.conf_update('pool.pre_schedule.enable_threshold', 1.75)

	conf = {
		'candidates': ['neumf'],
		'jm': jm,
		'api': api,
		'seed': seed,
		'job_num': 100,
	}
	launch(conf)


def test(jm=None, api=None, seed=2020, case=1):
	cases = ["cnn-pre", "cnn-default", "lstm-pre", "lstm-default", "neumf-default", "neumf-pre"]
	cases.extend(["resnet50-pre", "resnet50-default", "mixed-pre", "mixed-default"])
	if case not in cases:
		print('[WARN] case not in ' + str(cases))
		return

	msg = api.conf_reset()
	print(msg)
	if case == 'cnn-pre' or case == 'lstm-pre' or case == 'resnet50-pre' or case == 'mixed-pre' or case == 'neumf-pre':
		msg = api.conf_update('pool.pre_schedule.enable_threshold', 0.5)
		print(msg)
	else:
		api.conf_update('pool.pre_schedule.enable_threshold', 1.75)

	conf = {}
	if case == 'cnn-default' or case == 'cnn-pre':
		conf = {
			'candidates': ['cnn'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 4,
		}
	elif case == 'lstm-default' or case == 'lstm-pre':
		conf = {
			'candidates': ['lstm'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 4,
		}
	elif case == 'resnet50-default' or case == 'resnet50-pre':
		conf = {
			'candidates': ['resnet50'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 4,
		}
	elif case == 'mixed-default' or case == 'mixed-pre':
		conf = {
			'candidates': ['cnn', 'lstm', 'resnet50'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 4,
		}
	elif case == 'neumf-default' or case == 'neumf-pre':
		conf = {
			'candidates': ['neumf'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 4,
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


if __name__ == '__main__':
	print("pre_schedule.launcher")
	pass
