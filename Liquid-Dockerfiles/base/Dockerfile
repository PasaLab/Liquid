FROM nvidia/cuda:9.0-cudnn7-runtime

MAINTAINER Newnius <newnius.cn@gmail.com>

USER root

RUN apt update && \
	apt install -y git && \
	rm -rf /var/lib/apt/lists/*

ADD bootstrap.sh /etc/bootstrap.sh

ENTRYPOINT ["/etc/bootstrap.sh"]

