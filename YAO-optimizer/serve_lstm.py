#!/usr/bin/python
from threading import Thread
from threading import Lock
from http.server import BaseHTTPRequestHandler, HTTPServer
import cgi
import json
from urllib import parse
import pandas as pd
import csv
from pandas import DataFrame
from pandas import Series
from pandas import concat
from pandas import read_csv
from sklearn.metrics import mean_squared_error
from sklearn.preprocessing import MinMaxScaler
from keras.models import Sequential
from keras.layers import Dense
from keras.layers import LSTM
from math import sqrt
import numpy
import random
import traceback
from keras.models import load_model
from sklearn.externals import joblib

PORT_NUMBER = 8080
lock = Lock()
models = {}


# frame a sequence as a supervised learning problem
def timeseries_to_supervised(data, lag=1):
	df = DataFrame(data)
	columns = [df.shift(i) for i in range(1, lag + 1)]
	columns.append(df)
	df = concat(columns, axis=1)
	df = df.drop(0)
	return df


# create a differenced series
def difference(dataset, interval=1):
	diff = list()
	for i in range(interval, len(dataset)):
		value = dataset[i] - dataset[i - interval]
		diff.append(value)
	return Series(diff)


# invert differenced value
def inverse_difference(history, yhat, interval=1):
	return yhat + history[-interval]


# inverse scaling for a forecasted value
def invert_scale(scaler, X, yhat):
	new_row = [x for x in X] + [yhat]
	array = numpy.array(new_row)
	array = array.reshape(1, len(array))
	inverted = scaler.inverse_transform(array)
	return inverted[0, -1]


# fit an LSTM network to training data
def fit_lstm(train, batch_size2, nb_epoch, neurons):
	X, y = train[:, 0:-1], train[:, -1]
	X = X.reshape(X.shape[0], 1, X.shape[1])
	model = Sequential()
	model.add(LSTM(neurons, batch_input_shape=(batch_size2, X.shape[1], X.shape[2]), stateful=True))
	model.add(Dense(1))
	model.compile(loss='mean_squared_error', optimizer='adam')
	for i in range(nb_epoch):
		model.fit(X, y, epochs=1, batch_size=batch_size2, verbose=0, shuffle=False)
		# loss = model.evaluate(X, y)
		# print("Epoch {}/{}, loss = {}".format(i, nb_epoch, loss))
		print("Epoch {}/{}".format(i, nb_epoch))
		model.reset_states()
	return model


def train_models(job):
	lock.acquire()
	if job not in models:
		models[job] = {
			'lock': Lock()
		}
	lock.release()

	models[job]['lock'].acquire()

	# load dataset
	series = read_csv('./data/' + job + '.csv', header=0, index_col=0, squeeze=True)

	# transform data to be stationary
	raw_values = series.values
	diff_values = difference(raw_values, 1)
	# transform data to be supervised learning
	lag = 4
	supervised = timeseries_to_supervised(diff_values, lag)
	supervised_values = supervised.values

	batch_size = 32
	if supervised_values.shape[0] < 100:
		batch_size = 16
	if supervised_values.shape[0] < 60:
		batch_size = 8

	# split data into train and test-sets
	train = supervised_values
	# transform the scale of the data

	# scale data to [-1, 1]
	# fit scaler
	scaler = MinMaxScaler(feature_range=(-1, 1))
	scaler = scaler.fit(train)
	# transform train
	train = train.reshape(train.shape[0], train.shape[1])
	train_scaled = scaler.transform(train)

	# fit the model
	t1 = train.shape[0] % batch_size

	train_trimmed = train_scaled[t1:, :]
	model = fit_lstm(train_trimmed, batch_size, 30, 4)

	model.save('./data/checkpoint-' + job)
	scaler_filename = './data/checkpoint-' + job + "-scaler.save"
	joblib.dump(scaler, scaler_filename)

	models[job]['batch_size'] = batch_size

	models[job]['lock'].release()


def predict(job, seq):
	if job not in models or 'batch_size' not in models[job]:
		return -1, False

	batch_size = int(models[job]['batch_size'])

	data = {
		'seq': seq,
		'value': 0,
	}
	model = load_model('./data/checkpoint-' + job)
	scaler_filename = './data/checkpoint-' + job + "-scaler.save"
	scaler = joblib.load(scaler_filename)

	file = './data/' + job + '.' + str(random.randint(1000, 9999)) + '.csv'
	df = pd.read_csv('./data/' + job + '.csv', usecols=['seq', 'value'])
	df = df.tail(batch_size * 2 - 1)
	df = df.append(data, ignore_index=True)
	df.to_csv(file, index=False)

	# load dataset
	df = read_csv(file, header=0, index_col=0, squeeze=True)

	# transform data to be stationary
	raw_values = df.values
	diff_values = difference(raw_values, 1)

	# transform data to be supervised learning
	lag = 4
	supervised = timeseries_to_supervised(diff_values, lag)
	supervised_values = supervised[-batch_size:]
	test = supervised_values.values

	test = test.reshape(test.shape[0], test.shape[1])
	test_scaled = scaler.transform(test)

	# forecast the entire training dataset to build up state for forecasting
	test_reshaped = test_scaled[:, 0:-1]
	test_reshaped = test_reshaped.reshape(len(test_reshaped), 1, lag)
	output = model.predict(test_reshaped, batch_size=batch_size)
	predictions = list()
	for i in range(len(output)):
		yhat = output[i, 0]
		X = test_scaled[i, 0:-1]
		# invert scaling
		yhat = invert_scale(scaler, X, yhat)
		# invert differencing
		yhat = inverse_difference(raw_values, yhat, len(test_scaled) + 1 - i)
		# store forecast
		predictions.append(yhat)
	# report performance

	rmse = sqrt(mean_squared_error(raw_values[-batch_size:], predictions))
	print(predictions, raw_values[-batch_size:])
	return predictions[-1], True


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
				seq = query.get('seq')[0]
				msg = {'code': 0, 'error': ""}

				pred, success = predict(job, int(seq))

				if not success:
					msg = {'code': 2, 'error': "Job " + job + " not exist"}
				else:
					msg = {'code': 0, 'error': "", "total": int(pred)}
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
				seq = query.get('seq')[0]
				value = query.get('value')[0]

				if int(seq) == 1:
					with open('./data/' + job + '.csv', 'w', newline='') as csvfile:
						spamwriter = csv.writer(
							csvfile, delimiter=',',
							quotechar='|', quoting=csv.QUOTE_MINIMAL
						)
						spamwriter.writerow(["seq", "value"])

				with open('./data/' + job + '.csv', 'a+', newline='') as csvfile:
					spamwriter = csv.writer(
						csvfile, delimiter=',',
						quotechar='|', quoting=csv.QUOTE_MINIMAL
					)
					spamwriter.writerow([seq, value])
				msg = {'code': 0, 'error': ""}
			except Exception as e:
				msg = {'code': 1, 'error': str(e)}
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
