FROM ubuntu:16.04

RUN apt update && \
	apt install -y wget vim git gcc httpie dnsutils

RUN wget https://dl.google.com/go/go1.12.4.linux-amd64.tar.gz && \
	tar -C /usr/local -xzf go1.12.4.linux-amd64.tar.gz && \
	rm go1.12.4.linux-amd64.tar.gz

ENV PATH $PATH:/usr/local/go/bin

RUN go get github.com/Shopify/sarama
RUN go get github.com/sirupsen/logrus
RUN go get github.com/MaxHalford/eaopt

RUN git clone https://github.com/newnius/YAO-scheduler.git /root/yao-scheduler/

RUN mkdir -p /etc/yao/

ADD conf/config.json /etc/yao/

WORKDIR /root/yao-scheduler/

ADD bootstrap.sh /etc/

CMD ["/etc/bootstrap.sh"]
