FROM tensorflow/tensorflow:1.14.0-gpu

MAINTAINER Newnius <newnius.cn@gmail.com>

RUN apt update && \
	apt install -y python3 python3-pip && \
	rm -rf /var/lib/apt/lists/*

RUN apt update && \
	apt install -y git vim httpie && \
	rm -rf /var/lib/apt/lists/*

RUN pip3 install pandas sklearn tensorflow-gpu==1.14 keras

ADD bootstrap.sh /etc/bootstrap.sh

RUN git clone https://code.newnius.com/newnius/yao-optimizer.git /root/yao-optimizer/

WORKDIR /root/yao-optimizer

CMD ["/etc/bootstrap.sh"]