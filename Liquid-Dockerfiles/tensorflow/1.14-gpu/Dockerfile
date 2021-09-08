FROM tensorflow/tensorflow:1.14.0-gpu

MAINTAINER Newnius <newnius.cn@gmail.com>

USER root

RUN apt update && \
	apt install -y git && \
	rm -rf /var/lib/apt/lists/*

RUN pip install requests psutil py-cpuinfo pandas hdfs

ADD bootstrap.sh /etc/bootstrap.sh

ADD save.py /etc/save.py

ENTRYPOINT ["/etc/bootstrap.sh"]
