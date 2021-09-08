#!/usr/bin/python
from http.server import BaseHTTPRequestHandler, HTTPServer
import cgi
import docker
import json
from urllib import parse

PORT_NUMBER = 8000


# This class will handles any incoming request from
# the browser
class MyHandler(BaseHTTPRequestHandler):
	# Handler for the GET requests
	def do_GET(self):
		req = parse.urlparse(self.path)

		if req.path == "/list":
			try:
				client = docker.from_env()
				networks = client.networks.list(filters={'name': 'yao-net-'})
				result = []
				for network in networks:
					result.append(network.name)
				msg = {'code': 0, 'networks': result}
			except Exception as e:
				msg = {'code': 1, 'error': str(e)}
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
			try:
				network_name = form.getvalue('name')
				client = docker.from_env()
				client.networks.create(
					name=network_name,
					driver='overlay',
					attachable=True
				)
				msg = {"code": 0, "error": 'Success'}
			except Exception as e:
				msg = {"code": 1, "error": str(e)}

			self.send_response(200)
			self.send_header('Content-type', 'application/json')
			self.end_headers()
			self.wfile.write(bytes(json.dumps(msg), "utf-8"))

		else:
			self.send_error(404, 'File Not Found: %s' % self.path)


try:
	# Create a web server and define the handler to manage the
	# incoming request
	server = HTTPServer(('', PORT_NUMBER), MyHandler)
	print('Started http server on port ', PORT_NUMBER)

	# Wait forever for incoming http requests
	server.serve_forever()

except KeyboardInterrupt:
	print('^C received, shutting down the web server')

server.socket.close()
