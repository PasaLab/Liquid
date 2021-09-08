# Liquid-agent

https://docker-py.readthedocs.io/en/stable/containers.html

```bash
bin/kafka-topics.sh \
	--create \
	--zookeeper zookeeper-node1:2181,zookeeper-node2:2181,zookeeper-node3:2181 \
	--replication-factor 3 \
	--partitions 6 \
	--topic yao
```

```bash
bin/kafka-topics.sh \
	--describe \
	--zookeeper zookeeper-node1:2181,zookeeper-node2:2181,zookeeper-node3:2181 \
	--topic yao
```

```bash
bin/kafka-console-consumer.sh \
	--bootstrap-server kafka-node1:9092,kafka-node2:9092,kafka-node3:9092 \
	--topic yao \
	--from-beginning
```

```bash
bin/kafka-console-producer.sh \
	--broker-list kafka-node1:9092,kafka-node2:9092,kafka-node3:9092 \
	--topic yao
```

```bash
bin/kafka-topics.sh \
	--delete \
	--zookeeper zookeeper-node1:2181,zookeeper-node2:2181,zookeeper-node3:2181 \
	--topic yao
```
