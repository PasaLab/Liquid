FROM quickdeploy/yao-python3

MAINTAINER Newnius <newnius.cn@gmail.com>

RUN pip3 install docker psutil

ADD bootstrap.sh /etc/bootstrap.sh

ADD executor.py /root/executor.py

WORKDIR /root

CMD ["/etc/bootstrap.sh"]