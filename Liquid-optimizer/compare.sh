#!/usr/bin/env bash

FILE=./dataset.csv

if ! [[ -f "$FILE" ]]; then
    echo "$FILE not exist."
    exit 1
fi

awk 'FNR>1' ${FILE} | shuf > ./data/train.csv
head -n 1   ${FILE} > ./tmp.csv
head -n 1   ${FILE} > ./data/test.csv
head -n -50 ./data/train.csv >> ./tmp.csv
tail -n 50  ./data/train.csv >> ./data/test.csv

cnt=`wc ${FILE} | awk '{print $1}'`
step=50
maxn=$((cnt / step * step + 1))

step=51
while [[ ${step} -le ${maxn} ]]; do
	echo "step=${step}"
	tail -n ${step} tmp.csv > ./data/train.csv
	echo 'lr:'
	display_diff=0 algorithm=lr   python3 rf.py

	echo 'rf:'
	display_diff=0 algorithm=rf   python3 rf.py

	echo 'dt:'
	display_diff=0 algorithm=dt   python3 rf.py

	echo 'ada:'
	display_diff=0 algorithm=ada  python3 rf.py

	echo 'gbdt:'
	display_diff=0 algorithm=gbdt python3 rf.py
	echo -e "\n"
	step=$(($step + 50))
done

rm ./data/train.csv
rm ./data/test.csv
rm tmp.csv