#!/bin/bash

IF=$1
if [[ -z "$IF" ]]; then
	IF=`ls -1 /sys/class/net/ | head -1`
	echo "Usage: $0 <Interface, eg>"
	echo "eg: $0 docker0"
	exit;
fi

echo "Listening $IF..."

RX_PREV=-1
TX_PREV=-1
while true; do
	RX=`cat /sys/class/net/${IF}/statistics/rx_bytes`
	TX=`cat /sys/class/net/${IF}/statistics/tx_bytes`

	if [[ ${RX_PREV} -ne -1 ]]; then
	    let BW_RX=(${RX}-${RX_PREV})/1024/1024/3
	    let BW_TX=(${TX}-${TX_PREV})/1024/1024/3
		#echo "Received: ${BW_RX} MB/s    Sent: ${BW_TX} MB/s"
		echo "$(date),$BW_RX,$BW_TX"
	fi

	RX_PREV=${RX}
	TX_PREV=${TX}
	sleep 3

done