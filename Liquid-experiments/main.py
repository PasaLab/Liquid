import time
import os
import argparse
import util.jobModel as jm
import util.api as api


def test():
	apiInstance = api.API(base_url='http://yao.pasalab.jluapp.com')
	apiInstance.login()

	apiInstance.get_sys_status()

	job = jm.get_job('cnn', 0)
	msg = apiInstance.submit_job(job)
	if msg['errno'] == 0:
		apiInstance.job_status(msg['job_name'])

	apiInstance.job_list()

	msg = apiInstance.job_list()
	print("\nList of jobs:")
	if len(msg['jobs']) > 0:
		for job in msg['jobs']:
			name = job['name']
			status = jm.convert_job_status(job['status'])
			print("Status of job: {} is {}".format(name, status))
	print("\n")

	job_name = ''
	msg = apiInstance.job_status(job_name)
	print("Status of tasks in {}:".format(job_name))
	if len(msg['tasks']) > 0:
		for task in msg['tasks']:
			print("{} ({}) is/was running on {}".format(task['hostname'], task['status'], task['node']))


if __name__ == '__main__':
	# test()
	os.environ["TZ"] = 'Asia/Shanghai'
	if hasattr(time, 'tzset'):
		time.tzset()

	parser = argparse.ArgumentParser()
	parser.add_argument("-lab", "--lab", help="experiment name")
	parser.add_argument("-mode", "--mode", help="train/test")
	parser.add_argument("-case", "--case", help="case name")
	parser.add_argument("-seed", "--seed", help="random seed", default=2020)
	args = parser.parse_args()

	lab = args.lab
	mode = args.mode
	case = args.case
	seed = int(args.seed)

	if lab is None or mode is None or case is None:
		exit('[WARN] invalid param, use -h to see all params')

	apiInstance = api.API(base_url='http://yao.pasalab.jluapp.com')
	apiInstance.login()

	if lab == 'share':
		import share.launcher as launcher
	elif lab == 'pre_schedule':
		import pre_schedule.launcher as launcher
	elif lab == 'placement':
		import placement.launcher as launcher
	elif lab == 'predict':
		import predict.launcher as launcher
	elif lab == 'topology-aware':
		import topology_aware.launcher as launcher
	elif lab == 'preempt':
		import preempt.launcher as launcher
	elif lab == 'checkpoint':
		import checkpoint.launcher as launcher
	elif lab == 'qdrf':
		import qdrf.launcher as launcher
	elif lab == 'batch':
		import batch.launcher as launcher
	elif lab == 'heartbeat':
		import heartbeat.launcher as launcher
	elif lab == 'concurrent':
		import concurrent.launcher as launcher
	elif lab == 'hash':
		import hash.launcher as launcher
	elif lab == 'conf_list':
		print(apiInstance.conf_list())
		exit()
	else:
		exit('[WARN] Invalid lab param')

	if mode == 'train':
		launcher.train(jm=jm, api=apiInstance, seed=seed, case=case)
	elif mode == 'test':
		launcher.test(jm=jm, api=apiInstance, seed=seed, case=case)
	else:
		exit('[WARN] mode not in ["train", "test"]')
