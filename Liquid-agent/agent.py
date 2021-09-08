import os
from threading import Thread
from threading import Lock
import time
import subprocess
import json
from xml.dom.minidom import parse
import xml.dom.minidom
import multiprocessing
import psutil
import math
from http.server import BaseHTTPRequestHandler, HTTPServer
from socketserver import ThreadingMixIn
import cgi
import docker
from urllib import parse
import random
import string
from pathlib import Path
import requests
import traceback

ClientID = os.getenv('ClientID', 1)
ClientHost = os.getenv('ClientHost', "localhost")
ClientExtHost = os.getenv('ClientExtHost', "localhost")
ReportAddress = os.getenv('ReportAddress', "http://yao-scheduler:8080/?action=agent_report")

RackID = os.getenv('RackID', "default")
DomainID = os.getenv('DomainID', "default")

PORT = int(os.getenv('Port', 8000))
HeartbeatInterval = int(os.getenv('HeartbeatInterval', 5))

EnableEventTrigger = os.getenv('EnableEventTrigger', 'true')

lock = Lock()
pending_tasks = {}
id2token = {}

counter = {}

event_counter = 0

client = docker.from_env()

taskStats = {}
taskStatsLock = Lock()

active_stats = {0: {
	'util': 0,
	'mem_util': 0,
	'mem': 0
}}


def generate_token(string_length=8):
	letters = string.ascii_lowercase
	return ''.join(random.choice(letters) for i in range(string_length))


def monitor_task(container_id):
	print(container_id)
	try:
		container = client.containers.get(container_id)
	except Exception as e:
		print("[WARN]", str(e))
		return

	pid = 0
	maxCPU = 0
	maxMem = 0
	last_bw_rx = 0
	last_bw_tx = 0
	last_time = time.time() - 1
	for statR in container.stats():
		if pid == 0:
			try:
				res = container.top()['Processes']
				for x in res:
					if "/workspace" in x[7] and int(x[1]) in active_stats:
						pid = int(x[1])
						break
			except Exception as e:
				print("[WARN]", str(e))

		stat = json.loads(statR)
		# print(stat)
		if stat['read'] == '0001-01-01T00:00:00Z':
			print('container ', container_id, ' exited')
			break

		taskStatsLock.acquire()
		# CPU load, x%
		cur = stat['cpu_stats']['cpu_usage']['total_usage']
		last = stat['precpu_stats']['cpu_usage']['total_usage']
		utilCPU = (cur - last) / 10000000

		# Memory, MB
		mem = stat['memory_stats']['stats']['active_anon']
		mem = mem / 1024
		mem = mem / 1024

		# Bandwidth, KB/s
		cur = stat['networks']['eth0']['rx_bytes'] / 1024
		bw_rx = cur - last_bw_rx
		last_bw_rx = cur
		cur = stat['networks']['eth0']['tx_bytes'] / 1024
		bw_tx = cur - last_bw_tx
		last_bw_tx = cur
		now = time.time()
		dur = now - last_time
		last_time = now
		bw_rx /= dur
		bw_tx /= dur

		taskStats[container_id] = {
			'cpu': utilCPU,
			'mem': mem,
			'bw_rx': bw_rx,
			'bw_tx': bw_tx,
			'gpu_util': active_stats[pid]['util'],
			'gpu_mem_util': active_stats[pid]['mem_util'],
			'gpu_mem': active_stats[pid]['mem'],
		}
		# print(taskStats[container_id])
		# print(utilCPU, mem, maxCPU, maxMem, bw_rx, bw_tx)
		taskStatsLock.release()
		if stat['preread'] == '0001-01-01T00:00:00Z':
			continue
		if utilCPU > maxCPU:
			maxCPU = utilCPU
		if mem > maxMem:
			maxMem = mem
		# When container exited, break & clear taskStats after 30s
		if pid != 0 and container.status != 'running':
			time.sleep(30)
			taskStatsLock.acquire()
			taskStats.pop(container_id, None)
			taskStatsLock.release()
			break


def launch_tasks(stats):
	utils = {}
	mem_frees = {}
	for stat in stats:
		utils[stat['uuid']] = stat['utilization_gpu']
		if int(stat['utilization_gpu']) < 10:
			if stat['uuid'] not in counter:
				counter[stat['uuid']] = 0
			counter[stat['uuid']] += 1
		else:
			counter[stat['uuid']] = 0
		mem_frees[stat['uuid']] = stat['memory_free']

	entries_to_remove = []
	lock.acquire()
	for token, task in pending_tasks.items():
		if int(utils[task['gpus'][0]]) < 10 and counter[task['gpus'][0]] >= 2 \
				and (mem_frees[task['gpus'][0]] > task['gpu_mem'] or mem_frees[task['gpus'][0]] < 100):
			entries_to_remove.append(token)

	for k in entries_to_remove:
		pending_tasks.pop(k, None)
	lock.release()


class MyHandler(BaseHTTPRequestHandler):
	# Handler for the GET requests
	def do_GET(self):
		req = parse.urlparse(self.path)
		query = parse.parse_qs(req.query)

		if req.path == "/ping":
			self.send_response(200)
			self.send_header('Content-type', 'application/json')
			self.end_headers()
			self.wfile.write(bytes("pong", "utf-8"))

		elif req.path == "/debug":
			msg = {
				'pending_tasks': pending_tasks,
				'id2token': id2token,
				'event_counter': event_counter,
				'taskStats': taskStats
			}
			self.send_response(200)
			self.send_header('Content-type', 'application/json')
			self.end_headers()
			self.wfile.write(bytes(json.dumps(msg), "utf-8"))

		elif req.path == "/can_run":
			res = "1"
			try:
				token = query.get('token')[0]
				for i in range(0, 50):
					if token in pending_tasks:
						res = "0"
					else:
						break
					time.sleep(0.1)
			except Exception as e:
				print(e)
			self.send_response(200)
			self.send_header('Content-type', 'application/json')
			self.end_headers()
			self.wfile.write(bytes(res, "utf-8"))

		elif req.path == "/logs":
			try:
				container_id = query.get('id')[0]
				container = client.containers.get(container_id)
				msg = {'code': 0, 'logs': str(container.logs().decode())}
			except Exception as e:
				msg = {'code': 1, 'error': str(e)}
			self.send_response(200)
			self.send_header('Content-type', 'application/json')
			self.end_headers()
			self.wfile.write(bytes(json.dumps(msg), "utf-8"))

		elif req.path == "/status":
			try:
				container_id = query.get('id')[0]
				container = client.containers.list(all=True, filters={'id': container_id})
				if len(container) > 0:
					container = container[0]
					status = {
						'id': container.short_id,
						'image': container.attrs['Config']['Image'],
						'image_digest': container.attrs['Image'],
						'command': container.attrs['Config']['Cmd'],
						'created_at': container.attrs['Created'],
						'finished_at': container.attrs['State']['FinishedAt'],
						'status': container.status,
						'hostname': container.attrs['Config']['Hostname'],
						'state': container.attrs['State'],
						'cpu': 0,
						'mem': 0,
						'bw_rx': 0,
						'bw_tx': 0
					}
					if container_id in taskStats:
						taskStatsLock.acquire()
						status['cpu'] = taskStats[container_id]['cpu']
						status['mem'] = taskStats[container_id]['mem']
						status['bw_rx'] = taskStats[container_id]['bw_rx']
						status['bw_tx'] = taskStats[container_id]['bw_tx']
						status['bw_tx'] = taskStats[container_id]['bw_tx']
						status['gpu_util'] = taskStats[container_id]['gpu_util']
						status['gpu_mem_util'] = taskStats[container_id]['gpu_mem_util']
						status['gpu_mem'] = taskStats[container_id]['gpu_mem']
						taskStatsLock.release()
					if container_id in id2token:
						token = id2token[container_id]
						if token in pending_tasks:
							status['status'] = 'ready'
						else:
							id2token.pop(container_id, None)
					if status['command'] is not None:
						status['command'] = ' '.join(container.attrs['Config']['Cmd'])
					msg = {'code': 0, 'status': status}
				else:
					msg = {'code': 1, 'error': "container not exist"}
			except Exception as e:
				msg = {'code': 2, 'error': str(e)}
			self.send_response(200)
			self.send_header('Content-type', 'application/json')
			self.end_headers()
			self.wfile.write(bytes(json.dumps(msg), "utf-8"))

		else:
			self.send_error(404, 'File Not Found: %s' % self.path)

	# Handler for the POST requests
	def do_POST(self):
		if self.path == "/create":
			form = cgi.FieldStorage(
				fp=self.rfile,
				headers=self.headers,
				environ={
					'REQUEST_METHOD': 'POST',
					'CONTENT_TYPE': self.headers['Content-Type'],
				})
			docker_image = form.getvalue('image')
			docker_name = form.getvalue('name')
			docker_cmd = form.getvalue('cmd', '')
			docker_workspace = form.getvalue('workspace')
			docker_gpus = form.getvalue('gpus', '')
			docker_mem_limit = form.getvalue('mem_limit', 4096)
			docker_cpu_limit = form.getvalue('cpu_limit', 4)
			docker_network = form.getvalue('network')
			docker_wait = form.getvalue('should_wait', '0')
			docker_output = form.getvalue('output_dir', '')
			docker_hdfs_address = form.getvalue('hdfs_address', '')
			docker_hdfs_dir = form.getvalue('hdfs_dir', '')
			docker_gpu_mem = form.getvalue('gpu_mem', 8192)
			dfs_src = form.getvalue('dfs_src', '')
			dfs_dst = form.getvalue('dfs_dst', '')
			token = generate_token(16)

			if len(dfs_src) > 0:
				failed_cnt = 0
				while True:
					failed_cnt += 1
					if failed_cnt > 3:
						break
					try:
						# Docker wouldn't create dir by default on bind mode,
						# see https://github.com/moby/moby/issues/13121
						path = Path(dfs_src)
						path.mkdir(parents=True, exist_ok=True)
						break
					except OSError as e:
						print("Creation of the directory %s failed" % dfs_src)
						print(e)
						print(traceback.format_exc())

			try:
				# set PYTHONUNBUFFERED=1 to output immediately
				# see https://tarunlalwani.com/post/why-delayed-output-python-docker/
				script = " ".join([
					"docker run",
					"--gpus '\"device=" + docker_gpus + "\"'",
					"--detach=True",
					"--hostname " + docker_name,
					"--network " + docker_network,
					"--network-alias " + docker_name,
					"--memory-reservation " + docker_mem_limit,
					"--cpus " + docker_cpu_limit,
					"--env repo=" + docker_workspace,
					"--env should_wait=" + docker_wait,
					"--env should_cb=" + 'http://' + ClientExtHost + ':' + str(PORT) + '/can_run?token=' + token,
					"--env output_dir=" + docker_output,
					"--env hdfs_address=" + docker_hdfs_address if len(docker_hdfs_address) > 0 else '',
					"--env hdfs_dir=" + docker_hdfs_dir if len(docker_hdfs_address) > 0 else '',
					"--env gpu_mem=" + docker_gpu_mem,
					"--env PYTHONUNBUFFERED=1",
					"--mount type=bind,src=" + dfs_src + ",dst=" + dfs_dst if len(dfs_src) > 0 else '',
					docker_image,
					docker_cmd
				])

				container = client.containers.get('yao-agent-helper')
				exit_code, output = container.exec_run(['sh', '-c', script])
				msg = {"code": 0, "id": output.decode('utf-8').rstrip('\n')}

				Thread(target=monitor_task, name='monitor_task', args=(msg['id'],)).start()

				if docker_wait == "1":
					lock.acquire()
					pending_tasks[token] = {'gpus': str(docker_gpus).split(','), 'gpu_mem': int(docker_gpu_mem)}
					id2token[msg['id']] = token
					lock.release()
				if exit_code != 0:
					msg["code"] = 1
					msg["error"] = output.decode('utf-8').rstrip('\n')
					print(msg["error"])
			except Exception as e:
				msg = {"code": 1, "error": str(e)}
				print(str(e))

			self.send_response(200)
			self.send_header('Content-type', 'application/json')
			self.end_headers()
			self.wfile.write(bytes(json.dumps(msg), "utf-8"))

		elif self.path == "/stop":
			form = cgi.FieldStorage(
				fp=self.rfile,
				headers=self.headers,
				environ={
					'REQUEST_METHOD': 'POST',
					'CONTENT_TYPE': self.headers['Content-Type'],
				})
			container_id = form.getvalue('id')

			try:
				container = client.containers.get(container_id)
				container.stop(timeout=1)
				msg = {"code": 0, "error": "Success"}
			except Exception as e:
				msg = {"code": 1, "error": str(e)}

			self.send_response(200)
			self.send_header('Content-type', 'application/json')
			self.end_headers()
			self.wfile.write(bytes(json.dumps(msg), "utf-8"))

		elif self.path == "/remove":
			form = cgi.FieldStorage(
				fp=self.rfile,
				headers=self.headers,
				environ={
					'REQUEST_METHOD': 'POST',
					'CONTENT_TYPE': self.headers['Content-Type'],
				})
			container_id = form.getvalue('id')

			try:
				container = client.containers.get(container_id)
				container.remove(force=True)
				msg = {"code": 0, "error": "Success"}
			except Exception as e:
				msg = {"code": 1, "error": str(e)}

			self.send_response(200)
			self.send_header('Content-type', 'application/json')
			self.end_headers()
			self.wfile.write(bytes(json.dumps(msg), "utf-8"))
		else:
			self.send_error(404, 'File Not Found: %s' % self.path)


class ThreadingSimpleServer(ThreadingMixIn, HTTPServer):
	pass


def event_trigger():
	global event_counter
	for event in client.events(decode=True, filters={'event': 'die'}):
		Thread(target=report).start()
		event_counter += 1
		print(event)


def report():
	try:
		status, msg_gpu = execute(['nvidia-smi', '-q', '-x', '-f', 'status.xml'])
		if not status:
			print("execute failed, ", msg_gpu)
		stats = get_gpu_status()
		report_msg(stats)
		Thread(target=launch_tasks, name='launch_tasks', args=(stats,)).start()
	except Exception as e:
		print("[WARN]", str(e))


def reporter():
	while True:
		report()
		time.sleep(HeartbeatInterval)


def pmon():
	while True:
		try:
			status, msg_gpu = execute(['nvidia-smi', 'pmon', '-c', '1', '-s', 'um'])
			if not status:
				print("[WARN] execute failed, ", msg_gpu, status)
			lists = msg_gpu.split('\n')
			for p in lists:
				if "#" not in p and "-" not in p:
					tmp = p.split()
					data = {
						'idx': int(tmp[0]),
						'pid': int(tmp[1]),
						'util': int(tmp[3]),
						'mem_util': int(tmp[4]),
						'mem': int(tmp[7])
					}
					active_stats[int(tmp[1])] = data
		except Exception as e:
			print("[WARN]", str(e))
		time.sleep(HeartbeatInterval)


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
		'id': ClientID,
		'rack': RackID,
		'domain': DomainID,
		'host': ClientHost,
		'status': stats,
		'cpu_num': multiprocessing.cpu_count(),
		'cpu_load': os.getloadavg()[0],
		'mem_total': math.floor(mem.total / (1024. ** 3)),
		'mem_available': math.floor(mem.available / (1024. ** 3)),
		'version': time.time()
	}

	data = json.dumps(post_fields)

	try:
		url = ReportAddress
		params = {'data': data}
		result = requests.post(url, data=params)
	except Exception as e:
		pass
	'''
	producer = KafkaProducer(bootstrap_servers=KafkaBrokers)
	future = producer.send('yao', value=data.encode(), partition=0)
	result = future.get(timeout=5)
	'''


def listener():
	global server
	try:
		# Create a web server and define the handler to manage the
		# incoming request
		server = ThreadingSimpleServer(('', PORT), MyHandler)
		print('[INFO] Started http server on port ', PORT)

		# Wait forever for incoming http requests
		server.serve_forever()
	except KeyboardInterrupt:
		print('^C received, shutting down the web server')
		server.socket.close()


if __name__ == '__main__':
	os.environ["TZ"] = 'Asia/Shanghai'
	if hasattr(time, 'tzset'):
		time.tzset()

	Thread(target=reporter).start()
	Thread(target=listener).start()
	Thread(target=pmon).start()
	if EnableEventTrigger == 'true':
		print('[INFO] start event trigger')
		Thread(target=event_trigger).start()

	while True:
		time.sleep(5)
