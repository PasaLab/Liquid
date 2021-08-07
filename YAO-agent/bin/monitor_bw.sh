#!/usr/bin/env bash

IF=$1
if [[ -z "$IF" ]]; then
	IF=`ls -1 /sys/class/net/ | head -1`
fi
RX_PREV=-1
TX_PREV=-1

echo "Listening $IF..."

while [[ 1 == 1 ]] ; do
	RX=`cat /sys/class/net/${IF}/statistics/rx_bytes`
	TX=`cat /sys/class/net/${IF}/statistics/tx_bytes`
	if [ ${RX_PREV} -ne -1 ] ; then
		let BWRX=($RX-$RX_PREV)/1024/1024/3
		let BWTX=($TX-$TX_PREV)/1024/1024/3
		#echo "Received: $BWRX MB/s    Sent: $BWTX MB/s"
		date=`date`
		echo "$date:${BWRX},${BWTX}"
	fi
	RX_PREV=${RX}
	TX_PREV=${TX}
	sleep 3
done