import os
import time
import subprocess
import json
from xml.dom.minidom import parse
import xml.dom.minidom
import multiprocessing
import psutil
import math
import datetime


def main():
	interval = 1
	while True:
		try:
			status, msg_gpu = execute(['nvidia-smi', '-q', '-x', '-f', 'status.xml'])
			if not status:
				print("execute failed, ", msg_gpu)
			stats = get_gpu_status()
			utils = []
			for stat in stats:
				utils.append(str(stat['utilization_gpu']))
			# report_msg(stats)
			t = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
			print(str(t) + ',' + ','.join(utils))
			time.sleep(interval)
		except Exception as e:
			print(e)
			time.sleep(interval)


def execute(cmd):
	try:
		result = subprocess.run(cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
		if result.returncode == 0:
			return True, result.stdout.decode('utf-8').rstrip('\n')
		return False, result.stderr.decode('utf-8').rstrip('\n')
	except Exception as e:
		return False, e


def get_gpu_status():
	DOMTree = xml.dom.minidom.parse("status.xml")
	collection = DOMTree.documentElement
	gpus = collection.getElementsByTagName("gpu")
	stats = []
	for gpu in gpus:
		stat = {
			'uuid': gpu.getElementsByTagName('uuid')[0].childNodes[0].data,
			'product_name': gpu.getElementsByTagName('product_name')[0].childNodes[0].data,
			'performance_state': gpu.getElementsByTagName('performance_state')[0].childNodes[0].data,
			'memory_total': gpu.getElementsByTagName('fb_memory_usage')[0].getElementsByTagName('total')[0].childNodes[
				0].data,
			'memory_free': gpu.getElementsByTagName('fb_memory_usage')[0].getElementsByTagName('free')[0].childNodes[
				0].data,
			'memory_used': gpu.getElementsByTagName('fb_memory_usage')[0].getElementsByTagName('used')[0].childNodes[
				0].data,
			'utilization_gpu':
				gpu.getElementsByTagName('utilization')[0].getElementsByTagName('gpu_util')[0].childNodes[0].data,
			'utilization_mem':
				gpu.getElementsByTagName('utilization')[0].getElementsByTagName('memory_util')[0].childNodes[0].data,
			'temperature_gpu':
				gpu.getElementsByTagName('temperature')[0].getElementsByTagName('gpu_temp')[0].childNodes[0].data,
			'power_draw':
				gpu.getElementsByTagName('power_readings')[0].getElementsByTagName('power_draw')[0].childNodes[0].data
		}

		stat['memory_total'] = int(float(stat['memory_total'].split(' ')[0]))
		stat['memory_free'] = int(float(stat['memory_free'].split(' ')[0]))
		stat['memory_used'] = int(float(stat['memory_used'].split(' ')[0]))
		stat['utilization_gpu'] = int(float(stat['utilization_gpu'].split(' ')[0]))
		stat['utilization_mem'] = int(float(stat['utilization_mem'].split(' ')[0]))
		stat['temperature_gpu'] = int(float(stat['temperature_gpu'].split(' ')[0]))
		stat['power_draw'] = int(float(stat['power_draw'].split(' ')[0]))

		stats.append(stat)
	return stats


def report_msg(stats):
	mem = psutil.virtual_memory()

	post_fields = {
		'status': stats,
		'cpu_num': multiprocessing.cpu_count(),
		'cpu_load': os.getloadavg()[0],
		'mem_total': math.floor(mem.total / (1024. ** 3)),
		'mem_available': math.floor(mem.available / (1024. ** 3))
	}
	data = json.dumps(post_fields)


if __name__ == '__main__':
	os.environ["TZ"] = 'Asia/Shanghai'
	if hasattr(time, 'tzset'):
		time.tzset()
	main()
