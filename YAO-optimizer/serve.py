#!/usr/bin/python
from threading import Thread
from threading import Lock
from http.server import BaseHTTPRequestHandler, HTTPServer
import cgi
import json
from urllib import parse
import pandas as pd
import numpy as np
import csv
import random
import traceback
import pickle
import os
from sklearn.ensemble import RandomForestRegressor
from sklearn.model_selection import train_test_split
from sklearn.metrics import mean_squared_error

PORT_NUMBER = int(os.getenv('Port', 8080))
lock = Lock()
models = {}


def load_data(trainfile, testfile):
	traindata = pd.read_csv(trainfile)
	testdata = pd.read_csv(testfile)
	feature_data = traindata.iloc[:, 1:-1]
	label_data = traindata.iloc[:, -1]
	test_feature = testdata.iloc[:, 1:]
	return feature_data, label_data, test_feature


def train_models(job):
	if job not in models or 'features' not in models[job]:
		return
	models[job]['lock'].acquire()
	try:
		for label in models[job]['labels']:
			trainfile = './data/' + job + '_' + label + '.csv'
			traindata = pd.read_csv(trainfile)
			feature_data = traindata.iloc[:, 1:-1]
			label_data = traindata.iloc[:, -1]

			x_train, x_test, y_train, y_test = train_test_split(feature_data, label_data, test_size=0.01)
			params = {
				'n_estimators': 70,
				'max_depth': 13,
				'min_samples_split': 10,
				'min_samples_leaf': 5,  # 10
				'max_features': len(models[job]['features']) - 1  # 7
			}
			# print(params)
			model = RandomForestRegressor(**params)
			model.fit(x_train, y_train)

			# save the model to disk
			modelname = './data/' + job + '_' + label + '.sav'
			pickle.dump(model, open(modelname, 'wb'))

			# 对测试集进行预测
			y_pred = model.predict(x_test)
			# 计算准确率
			MSE = mean_squared_error(y_test, y_pred)
			RMSE = np.sqrt(MSE)
			print('RMSE of {}:{} is {}'.format(job, label, str(RMSE)))
	except Exception as e:
		print(traceback.format_exc())
		print(str(e))

	models[job]['lock'].release()


def predict(job, features):
	if job not in models or 'features' not in models[job]:
		return -1, False

	values = [job]
	for feature in models[job]['features']:
		if feature in features:
			values.append(features[feature])
		else:
			values.append(0)

	testfile = './data/' + job + '.' + str(random.randint(1000, 9999)) + '.csv'
	t = ['job']
	t.extend(models[job]['features'])
	with open(testfile, 'w', newline='') as csvfile:
		spamwriter = csv.writer(
			csvfile, delimiter=',',
			quotechar='|', quoting=csv.QUOTE_MINIMAL
		)
		spamwriter.writerow(t)

	with open(testfile, 'a+', newline='') as csvfile:
		spamwriter = csv.writer(
			csvfile, delimiter=',',
			quotechar='|', quoting=csv.QUOTE_MINIMAL
		)
		spamwriter.writerow(values)

	testdata = pd.read_csv(testfile)
	test_feature = testdata.iloc[:, 1:]

	predictions = {}
	for label in models[job]['labels']:
		# load the model from disk
		modelfile = './data/' + job + '_' + label + '.sav'
		if not os.path.exists(modelfile):
			if os.path.exists(testfile):
				os.remove(testfile)
			return -1, False
		model = pickle.load(open(modelfile, 'rb'))
		preds = model.predict(test_feature)
		predictions[label] = preds[0]

	if os.path.exists(testfile):
		os.remove(testfile)
	return predictions, True


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

		elif req.path == "/predict":
			try:
				job = query.get('job')[0]
				features = json.loads(query.get('features')[0])
				pred, success = predict(job, features)

				if not success:
					msg = {'code': 2, 'error': "Job " + job + " not exist"}
				else:
					msg = {'code': 0, 'error': "", "labels": pred}
			except Exception as e:
				track = traceback.format_exc()
				print(track)
				msg = {'code': 1, 'error': str(e)}

			self.send_response(200)
			self.send_header('Content-type', 'application/json')
			self.end_headers()
			self.wfile.write(bytes(json.dumps(msg), "utf-8"))

		elif req.path == "/feed":
			try:
				job = query.get('job')[0]
				features = json.loads(query.get('features')[0])
				labels = json.loads(query.get('labels')[0])

				lock.acquire()
				flag = False
				if job not in models:
					models[job] = {
						'lock': Lock(),
						'features': list(features.keys()),
						'labels': list(labels.keys())
					}
					flag = True
				lock.release()
				models[job]['lock'].acquire()

				for label in models[job]['labels']:
					values = [job]
					for feature in models[job]['features']:
						if feature in features:
							values.append(features[feature])
						else:
							values.append(0)
					if label in labels:
						values.append(labels[label])
					else:
						values.append(0)

					if flag:
						t = ['job']
						t.extend(models[job]['features'])
						t.append(label)
						with open('./data/' + job + '_' + label + '.csv', 'w', newline='') as csvfile:
							spamwriter = csv.writer(
								csvfile, delimiter=',',
								quotechar='|', quoting=csv.QUOTE_MINIMAL
							)
							spamwriter.writerow(t)

					with open('./data/' + job + '_' + label + '.csv', 'a+', newline='') as csvfile:
						spamwriter = csv.writer(
							csvfile, delimiter=',',
							quotechar='|', quoting=csv.QUOTE_MINIMAL
						)
						spamwriter.writerow(values)

				models[job]['lock'].release()
				msg = {'code': 0, 'error': ""}
			except Exception as e:
				msg = {'code': 1, 'error': str(e)}
				track = traceback.format_exc()
				print(track)
			self.send_response(200)
			self.send_header('Content-type', 'application/json')
			self.end_headers()
			self.wfile.write(bytes(json.dumps(msg), "utf-8"))

		elif req.path == "/train":
			try:
				job = query.get('job')[0]
				t = Thread(target=train_models, name='train_models', args=(job,))
				t.start()
				msg = {'code': 0, 'error': ""}
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
			if self.path == "/train2":
				form = cgi.FieldStorage(
					fp=self.rfile,
					headers=self.headers,
					environ={
						'REQUEST_METHOD': 'POST',
						'CONTENT_TYPE': self.headers['Content-Type'],
					})
				try:
					job = form.getvalue('job')[0]
					seq = form.getvalue('seq')[0]
					t = Thread(target=train_models(), name='train_models', args=(job, seq,))
					t.start()
					msg = {"code": 0, "error": ""}
				except Exception as e:
					msg = {"code": 1, "error": str(e)}
				self.send_response(200)
				self.send_header('Content-type', 'application/json')
				self.end_headers()
				self.wfile.write(bytes(json.dumps(msg), "utf-8"))
			else:
				self.send_error(404, 'File Not Found: %s' % self.path)


if __name__ == '__main__':
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
