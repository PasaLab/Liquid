FROM tensorflow/tensorflow:1.14.0-gpu

MAINTAINER Newnius <newnius.cn@gmail.com>

RUN apt update && \
	apt install -y python3 python3-pip && \
	rm -rf /var/lib/apt/lists/*

RUN apt update && \
	apt install -y git vim httpie && \
	rm -rf /var/lib/apt/lists/*

RUN pip3 install docker kafka psutil

RUN git clone https://github.com/newnius/YAO-agent.git /yao-agent

WORKDIR /yao-agent

ADD bootstrap.sh /etc/

CMD ["/etc/bootstrap.sh"]
