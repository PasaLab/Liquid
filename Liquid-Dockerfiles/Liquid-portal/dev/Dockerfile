FROM newnius/php:7.3

MAINTAINER Newnius <newnius.cn@gmail.com>

RUN apt update && \
	apt install -y git vim httpie && \
	rm -rf /var/lib/apt/lists/*

RUN git clone https://github.com/newnius/YAO-portal.git /var/www/html/

ADD bootstrap.sh /etc/

CMD /etc/bootstrap.sh
