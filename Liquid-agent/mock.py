import os
from threading import Thread
import time
import json
from http.server import BaseHTTPRequestHandler, HTTPServer
from socketserver import ThreadingMixIn
from urllib import parse
import requests

NUMS = os.getenv('NUMS', 1)

ClientHost = os.getenv('ClientHost', "localhost")
ReportAddress = os.getenv('ReportAddress', "http://yao-scheduler:8080/?action=agent_report")
PORT = os.getenv('Port', 8000)
HeartbeatInterval = os.getenv('HeartbeatInterval', 5)


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

		elif req.path == "/logs":
			try:
				msg = {'code': 0, 'logs': 'Output from mock container'}
			except Exception as e:
				msg = {'code': 1, 'error': str(e)}
			self.send_response(200)
			self.send_header('Content-type', 'application/json')
			self.end_headers()
			self.wfile.write(bytes(json.dumps(msg), "utf-8"))

		elif req.path == "/status":
			status = {
				'id': 'mock-container-id',
				'image': 'mock-container-image',
				'image_digest': 'mock-image-digest',
				'command': 'mock-container-command',
				'created_at': 'mock-container-created-at',
				'finished_at': 'mock-container-finished-at',
				'status': 'running',
				'hostname': 'mock-container-hostname',
				'state': {'exitCode': 1}
			}
			msg = {'code': 0, 'status': status}

			self.send_response(200)
			self.send_header('Content-type', 'application/json')
			self.end_headers()
			self.wfile.write(bytes(json.dumps(msg), "utf-8"))

		else:
			self.send_error(404, 'File Not Found: %s' % self.path)

	# Handler for the POST requests
	def do_POST(self):
		if self.path == "/create":
			msg = {"code": 0, "id": 'mock-container-id'}
			self.send_response(200)
			self.send_header('Content-type', 'application/json')
			self.end_headers()
			self.wfile.write(bytes(json.dumps(msg), "utf-8"))

		elif self.path == "/stop":
			msg = {"code": 0, "error": "Success"}
			self.send_response(200)
			self.send_header('Content-type', 'application/json')
			self.end_headers()
			self.wfile.write(bytes(json.dumps(msg), "utf-8"))

		elif self.path == "/remove":
			msg = {"code": 0, "error": "Success"}
			self.send_response(200)
			self.send_header('Content-type', 'application/json')
			self.end_headers()
			self.wfile.write(bytes(json.dumps(msg), "utf-8"))
		else:
			self.send_error(404, 'File Not Found: %s' % self.path)


class ThreadingSimpleServer(ThreadingMixIn, HTTPServer):
	pass


def report(ClientID):
	while True:
		try:
			stats = []
			for i in range(0, 4):
				stat = {
					'uuid': 'UUID-' + ClientID + '-' + str(i),
					'product_name': 'K80',
					'performance_state': 'P0',
					'memory_total': 11260,
					'memory_free': 11260,
					'memory_used': 0,
					'utilization_gpu': 0,
					'utilization_mem': 0,
					'temperature_gpu': 45,
					'power_draw': 25
				}

				stats.append(stat)

			post_fields = {
				'id': ClientID,
				'rack': ClientHost,
				'domain': ClientHost,
				'host': ClientHost,
				'status': stats,
				'cpu_num': 64,
				'cpu_load': 3,
				'mem_total': 188,
				'mem_available': 180,
				'version': 0
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
			time.sleep(HeartbeatInterval)
		except Exception as e:
			print(e)
			time.sleep(HeartbeatInterval)


def listener():
	global server
	try:
		# Create a web server and define the handler to manage the
		# incoming request
		server = ThreadingSimpleServer(('', PORT), MyHandler)
		print('Started http server on port ', PORT)

		# Wait forever for incoming http requests
		server.serve_forever()
	except KeyboardInterrupt:
		print('^C received, shutting down the web server')

		server.socket.close()


if __name__ == '__main__':
	os.environ["TZ"] = 'Asia/Shanghai'
	if hasattr(time, 'tzset'):
		time.tzset()
	threads = []
	for clientID in range(0, int(NUMS)):
		t = Thread(target=report, name=ClientHost + '_' + str(clientID), args=(ClientHost + '_' + str(clientID),))
		threads.append(t)

	t2 = Thread(target=listener)
	threads.append(t2)

	# Start all threads
	for t in threads:
		t.start()

	# Wait for all of them to finish
	for t in threads:
		t.join()
