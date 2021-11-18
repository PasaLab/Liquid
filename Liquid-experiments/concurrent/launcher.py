import random
import json
import time


def train(jm=None, api=None, seed=2020, case=None):
	pass


def test(jm=None, api=None, seed=2020, case=1):
	cases = ["parallelism-1", "parallelism-2", "parallelism-5", "parallelism-10"]
	if case not in cases:
		print('[WARN] case not in ' + str(cases))
		return

	api.conf_reset()

	conf = {}
	if case == 'parallelism-1':
		api.conf_update('scheduler.parallelism', 1)
		conf = {
			'candidates': ['job1', 'job2', 'job3', 'job5', 'job10'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 100,
		}
	elif case == 'parallelism-2':
		api.conf_update('scheduler.parallelism', 2)
		conf = {
			'candidates': ['job1', 'job2', 'job3', 'job5', 'job10'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 100,
		}
	elif case == 'parallelism-5':
		api.conf_update('scheduler.parallelism', 5)
		conf = {
			'candidates': ['job1', 'job2', 'job3', 'job5', 'job10'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 100,
		}
	elif case == 'parallelism-10':
		api.conf_update('scheduler.parallelism', 10)
		conf = {
			'candidates': ['job1', 'job2', 'job3', 'job5', 'job10'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 100,
		}
	else:
		print('[ERROR] case not in ' + str(cases))
		return

	api.conf_update('scheduler.enabled', 'false')
	time.sleep(5)
	launch(conf)
	time.sleep(5)
	api.conf_update('scheduler.enabled', 'true')


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
	print("concurrent.launcher")
	pass
