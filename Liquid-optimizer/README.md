

## Feed

/feed?job=lstm&seq=1&value=3

## train
/train?job=lstm

## predict
/predict?job=lstm&seq=2


## Make Data

```bash
cat data/stock_data.csv | awk -F',' '{print NR "," $5 }' > data/stock.csv
```



# Ref

[Python Keras + LSTM 进行单变量时间序列预测](https://edmondfrank.github.io/blog/2018/02/22/python-keras-plus-lstm-jin-xing-dan-bian-liang-shi-jian-xu-lie-yu-ce/)

[How to Seed State for LSTMs for Time Series Forecasting in Python](https://machinelearningmastery.com/seed-state-lstms-time-series-forecasting-python/)

[单维时间序列预测](https://github.com/xchadesi/Space-time-Sequence/blob/master/4.%E5%8D%95%E7%BB%B4%E6%97%B6%E9%97%B4%E5%BA%8F%E5%88%97%E9%A2%84%E6%B5%8B.md)

[Python中利用LSTM模型进行时间序列预测分析](https://www.cnblogs.com/arkenstone/p/5794063.html)

[stock_predict_with_LSTM](https://github.com/hichenway/stock_predict_with_LSTM)

