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


# scale train and test data to [-1, 1]
def scale(train, test):
	# fit scaler
	scaler = MinMaxScaler(feature_range=(-1, 1))
	scaler = scaler.fit(train)
	# transform train
	train = train.reshape(train.shape[0], train.shape[1])
	train_scaled = scaler.transform(train)
	# transform test
	test = test.reshape(test.shape[0], test.shape[1])
	test_scaled = scaler.transform(test)
	return scaler, train_scaled, test_scaled


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
		model.reset_states()
	return model


# run a experiment
def experiment(series):
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
	test_data_num = batch_size

	# split data into train and test-sets
	train, test = supervised_values[0:-test_data_num], supervised_values[-test_data_num:]
	# transform the scale of the data

	scaler, train_scaled, test_scaled = scale(train, test)

	# run experiment
	error_scores = list()
	# fit the model
	t1 = train.shape[0] % batch_size
	t2 = test.shape[0] % batch_size

	train_trimmed = train_scaled[t1:, :]
	lstm_model = fit_lstm(train_trimmed, batch_size, 30, 4)

	# forecast the entire training dataset to build up state for forecasting
	test_reshaped = test_scaled[:, 0:-1]
	test_reshaped = test_reshaped.reshape(len(test_reshaped), 1, lag)
	output = lstm_model.predict(test_reshaped, batch_size=batch_size)
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
	rmse = sqrt(mean_squared_error(raw_values[-test_data_num:], predictions))
	print(predictions, raw_values[-test_data_num:])
	error_scores.append(rmse)
	return error_scores


# load dataset
series = read_csv('data.csv', header=0, index_col=0, squeeze=True)

with_seed = experiment(series)
