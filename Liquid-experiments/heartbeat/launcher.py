import random
import json
import time


def train(jm=None, api=None, seed=2020, case=None):
	pass


def test(jm=None, api=None, seed=2020, case=1):
	cases = ["before", "after"]
	if case not in cases:
		print('[WARN] case not in ' + str(cases))
		return

	api.conf_reset()

	conf = {}
	if case == 'before':
		conf = {
			'candidates': ['neumf'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 8,
		}
	elif case == 'after':
		conf = {
			'candidates': ['neumf'],
			'jm': jm,
			'api': api,
			'seed': seed,
			'job_num': 8,
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
		while True:
			next_seed = random.randint(0, 999999)
			job_name = random.choice(candidates)
			job = jm.get_job(job_name, seed=next_seed)
			tasks = json.loads(job['tasks'])
			if len(tasks) == 1 and tasks[0]['gpu_number'] == '1':
				break
		msg = api.submit_job(job)
		print(i, msg)
		time.sleep(3)


if __name__ == '__main__':
	print("pre_schedule.launcher")
	pass
