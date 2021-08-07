# _*_coding:utf-8_*_
import numpy as np
import pandas as pd
import os


def load_data(trainfile, testfile):
	traindata = pd.read_csv(trainfile)
	testdata = pd.read_csv(testfile)
	feature_data = traindata.iloc[:, 1:-1]
	label_data = traindata.iloc[:, -1]
	test_feature = testdata.iloc[:, 1:-1]
	test_label = testdata.iloc[:, -1]
	return feature_data, label_data, test_feature, test_label


def random_forest_train(feature_data, label_data, test_feature):
	from sklearn.ensemble import RandomForestRegressor
	from sklearn.model_selection import train_test_split
	from sklearn.metrics import mean_squared_error

	X_train, X_test, y_train, y_test = train_test_split(feature_data, label_data, test_size=0.01)
	params = {
		'n_estimators': 70,
		'max_depth': 13,
		'min_samples_split': 10,
		'min_samples_leaf': 5,  # 10
		'max_features': len(X_train.columns)
	}
	# print(X_test)
	model = RandomForestRegressor(**params)
	model.fit(X_train, y_train)
	# 对测试集进行预测
	y_pred = model.predict(X_test)
	# 计算准确率
	MSE = mean_squared_error(y_test, y_pred)
	RMSE = np.sqrt(MSE)
	# print(abs(y_test - y_pred) / y_test)
	# print(RMSE)
	'''
	submit = pd.read_csv(submitfile)
	print(submit)
	submit['CPU'] = model.predict(test_feature)
	submit.to_csv('my_random_forest_prediction1.csv', index=False)
	print(submit)
	print(model.predict(test_feature))
	'''
	return model.predict(test_feature)


def linear_regression_train(feature_data, label_data, test_feature):
	from sklearn.linear_model import LinearRegression
	from sklearn.model_selection import train_test_split
	from sklearn.metrics import mean_squared_error

	X_train, X_test, y_train, y_test = train_test_split(feature_data, label_data, test_size=0.01)
	params = {}
	# print(X_test)
	model = LinearRegression(**params)
	model.fit(X_train, y_train)
	# 对测试集进行预测
	y_pred = model.predict(X_test)
	# 计算准确率
	MSE = mean_squared_error(y_test, y_pred)
	RMSE = np.sqrt(MSE)
	# print(abs(y_test - y_pred) / y_test)
	# print(RMSE)
	return model.predict(test_feature)


def adaboost_train(feature_data, label_data, test_feature):
	from sklearn.ensemble import AdaBoostRegressor
	from sklearn.model_selection import train_test_split
	from sklearn.metrics import mean_squared_error

	X_train, X_test, y_train, y_test = train_test_split(feature_data, label_data, test_size=0.01)
	params = {}
	# print(X_test)
	model = AdaBoostRegressor(**params)
	model.fit(X_train, y_train)
	# 对测试集进行预测
	y_pred = model.predict(X_test)
	# 计算准确率
	MSE = mean_squared_error(y_test, y_pred)
	RMSE = np.sqrt(MSE)
	# print(abs(y_test - y_pred) / y_test)
	# print(RMSE)
	return model.predict(test_feature)


def gbdt_train(feature_data, label_data, test_feature):
	from sklearn.ensemble import GradientBoostingRegressor
	from sklearn.model_selection import train_test_split
	from sklearn.metrics import mean_squared_error

	X_train, X_test, y_train, y_test = train_test_split(feature_data, label_data, test_size=0.01)
	params = {
		'loss': 'ls',
		'n_estimators': 70,
		'max_depth': 13,
		'min_samples_split': 10,
		'min_samples_leaf': 5,  # 10
		'max_features': len(X_train.columns)
	}
	# print(X_test)
	model = GradientBoostingRegressor(**params)
	model.fit(X_train, y_train)
	# 对测试集进行预测
	y_pred = model.predict(X_test)
	# 计算准确率
	MSE = mean_squared_error(y_test, y_pred)
	RMSE = np.sqrt(MSE)
	# print(abs(y_test - y_pred) / y_test)
	# print(RMSE)
	return model.predict(test_feature)


def decision_tree_train(feature_data, label_data, test_feature):
	from sklearn.tree import DecisionTreeRegressor
	from sklearn.model_selection import train_test_split
	from sklearn.metrics import mean_squared_error

	X_train, X_test, y_train, y_test = train_test_split(feature_data, label_data, test_size=0.01)
	params = {
		'max_depth': 13,
	}
	# print(X_test)
	model = DecisionTreeRegressor(**params)
	model.fit(X_train, y_train)
	# 对测试集进行预测
	y_pred = model.predict(X_test)
	# 计算准确率
	MSE = mean_squared_error(y_test, y_pred)
	RMSE = np.sqrt(MSE)
	# print(abs(y_test - y_pred) / y_test)
	# print(RMSE)
	return model.predict(test_feature)


def random_forest_parameter_tuning1(feature_data, label_data):
	from sklearn.ensemble import RandomForestRegressor
	from sklearn.model_selection import train_test_split
	from sklearn.metrics import mean_squared_error
	from sklearn.model_selection import GridSearchCV

	X_train, X_test, y_train, y_test = train_test_split(feature_data, label_data, test_size=0.01)
	param_test1 = {
		'n_estimators': range(10, 71, 10)
	}
	model = GridSearchCV(estimator=RandomForestRegressor(
		min_samples_split=100, min_samples_leaf=20, max_depth=8, max_features='sqrt',
		random_state=10), param_grid=param_test1, cv=5
	)
	model.fit(X_train, y_train)
	# 对测试集进行预测
	y_pred = model.predict(X_test)
	# 计算准确率
	MSE = mean_squared_error(y_test, y_pred)
	RMSE = np.sqrt(MSE)
	print(RMSE)
	return model.best_score_, model.best_params_


def random_forest_parameter_tuning2(feature_data, label_data):
	from sklearn.ensemble import RandomForestRegressor
	from sklearn.model_selection import train_test_split
	from sklearn.metrics import mean_squared_error
	from sklearn.model_selection import GridSearchCV

	X_train, X_test, y_train, y_test = train_test_split(feature_data, label_data, test_size=0.01)
	param_test2 = {
		'max_depth': range(3, 14, 2),
		'min_samples_split': range(50, 201, 20)
	}
	model = GridSearchCV(estimator=RandomForestRegressor(
		n_estimators=70, min_samples_leaf=20, max_features='sqrt', oob_score=True,
		random_state=10), param_grid=param_test2, cv=5
	)
	model.fit(X_train, y_train)
	# 对测试集进行预测
	y_pred = model.predict(X_test)
	# 计算准确率
	MSE = mean_squared_error(y_test, y_pred)
	RMSE = np.sqrt(MSE)
	print(RMSE)
	return model.best_score_, model.best_params_


def random_forest_parameter_tuning3(feature_data, label_data, test_feature):
	from sklearn.ensemble import RandomForestRegressor
	from sklearn.model_selection import train_test_split
	from sklearn.metrics import mean_squared_error
	from sklearn.model_selection import GridSearchCV

	X_train, X_test, y_train, y_test = train_test_split(feature_data, label_data, test_size=0.01)
	param_test3 = {
		'min_samples_split': range(10, 90, 20),
		'min_samples_leaf': range(10, 60, 10),
	}
	model = GridSearchCV(estimator=RandomForestRegressor(
		n_estimators=70, max_depth=13, max_features='sqrt', oob_score=True,
		random_state=10), param_grid=param_test3, cv=5
	)
	model.fit(X_train, y_train)
	# 对测试集进行预测
	y_pred = model.predict(X_test)
	# 计算准确率
	MSE = mean_squared_error(y_test, y_pred)
	RMSE = np.sqrt(MSE)
	print(RMSE)
	return model.best_score_, model.best_params_


def random_forest_parameter_tuning4(feature_data, label_data, test_feature):
	from sklearn.ensemble import RandomForestRegressor
	from sklearn.model_selection import train_test_split
	from sklearn.metrics import mean_squared_error
	from sklearn.model_selection import GridSearchCV

	X_train, X_test, y_train, y_test = train_test_split(feature_data, label_data, test_size=0.01)
	param_test4 = {
		'max_features': range(3, 9, 2)
	}
	model = GridSearchCV(estimator=RandomForestRegressor(
		n_estimators=70, max_depth=13, min_samples_split=10, min_samples_leaf=10, oob_score=True,
		random_state=10), param_grid=param_test4, cv=5
	)
	model.fit(X_train, y_train)
	# 对测试集进行预测
	y_pred = model.predict(X_test)
	# 计算准确率
	MSE = mean_squared_error(y_test, y_pred)
	RMSE = np.sqrt(MSE)
	print(RMSE)
	return model.best_score_, model.best_params_


if __name__ == '__main__':
	algorithm = os.getenv('algorithm', 'rf')
	trainfile = 'data/train.csv'
	testfile = 'data/test.csv'
	feature_data, label_data, test_feature, test_label = load_data(trainfile, testfile)
	if algorithm == 'lr':
		y_pred = linear_regression_train(feature_data, label_data, test_feature)
	elif algorithm == 'ada':
		y_pred = adaboost_train(feature_data, label_data, test_feature)
	elif algorithm == 'gbdt':
		y_pred = gbdt_train(feature_data, label_data, test_feature)
	elif algorithm == 'dt':
		y_pred = decision_tree_train(feature_data, label_data, test_feature)
	else:
		y_pred = random_forest_train(feature_data, label_data, test_feature)

	from sklearn.metrics import mean_squared_error

	MSE = mean_squared_error(test_label, y_pred)
	RMSE = np.sqrt(MSE)
	var = np.var(test_label)
	r2 = 1 - MSE / var
	# print(abs(test_label - y_pred) / test_label)
	print(RMSE, r2)
	display_diff = os.getenv('display_diff', '0')
	if display_diff == '1':
		for i in range(20):
			print("{},{},{}".format(test_label[i], y_pred[i], (y_pred[i] - test_label[i]) / test_label[i]))
