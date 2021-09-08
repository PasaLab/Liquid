import subprocess
import docker


def run():
	client = docker.from_env()
	try:
		print(client.containers.run(image="alpine", command="pwd", environment={"KEY": "value"}))
	except Exception as e:
		print(e.__class__.__name__, e)


def run_in_background():
	client = docker.from_env()
	container = client.containers.run("alpine", ["echo", "hello", "world"], detach=True)
	print(container.id)


def list_containers():
	client = docker.from_env()
	for container in client.containers.list():
		print(container.id)


def get_logs(id):
	try:
		client = docker.from_env()
		container = client.containers.get(id)
		print(container.logs().decode())
	except Exception as e:
		print(e)


def get_status(id):
	client = docker.from_env()
	container = client.containers.list(all=True, filters={'id': id})
	status = {}
	if len(container) > 0:
		container = container[0]
		status['id'] = container.short_id
		status['image'] = container.attrs['Config']['Image']
		status['image_digest'] = container.attrs['Image']
		status['command'] = container.attrs['Config']['Cmd']
		status['createdAt'] = container.attrs['Created']
		status['finishedAt'] = container.attrs['State']['FinishedAt']
		status['status'] = container.status
		if status['command'] is not None:
			status['command'] = ' '.join(container.attrs['Config']['Cmd'])
	print(status)
	print(container.attrs)


def create_network():
	client = docker.from_env()
	client.networks.create(name='yao-net-1024', driver='overlay', attachable=True)


def list_networks():
	client = docker.from_env()
	networks = client.networks.list(filters={'name': 'yao-net-'})
	result = []
	for network in networks:
		result.append(network.name)
	print(result)


def remove_network():
	client = docker.from_env()
	client.networks.prune(filters={'name': 'yao-net-1024'})


def create_container():
	client = docker.APIClient(base_url='unix://var/run/docker.sock')

	host_config = client.create_host_config(
		mem_limit='512m',
		cpu_shares=1 * 1024
	)
	networking_config = client.create_networking_config(
		endpoints_config={
			# 'yao-net-1201': client.create_endpoint_config(
			# 	aliases=['node1'],
			# )
		}
	)

	container = client.create_container(
		image='alpine',
		command='pwd',
		hostname='node1',
		detach=True,
		host_config=host_config,
		environment={"repo": '', "NVIDIA_VISIBLE_DEVICES": '0'},
		networking_config=networking_config,
		runtime='nvidia'
	)
	client.start(container)
	print(container)


def exec_run():
	client = docker.from_env()
	container = client.containers.get('yao-agent-helper')
	exit_code, output = container.exec_run(
		cmd="sh -c 'docker run --gpus all --detach=True tensorflow/tensorflow:1.14.0-gpu nvidia-smi'")
	if exit_code == 0:
		print(output.decode('utf-8').rstrip('\n'))


def report():
	try:
		status, msg_gpu = execute(['nvidia-smi', 'pmon', '-c', '1', '-s', 'um'])
		if not status:
			print("execute failed, ", msg_gpu, status)
		lists = msg_gpu.split('\n')
		for p in lists:
			if "#" not in p and "-" not in p:
				tmp = p.split()
				data = {
					'idx': int(tmp[0]),
					'pid': int(tmp[1]),
					'util': int(tmp[3]),
					'mem_util': int(tmp[4]),
					'mem': int(tmp[7])
				}
				print(data)
	except Exception as e:
		print(e)


def execute(cmd):
	try:
		result = subprocess.run(cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
		if result.returncode == 0:
			return True, result.stdout.decode('utf-8').rstrip('\n')
		return False, result.stderr.decode('utf-8').rstrip('\n')
	except Exception as e:
		return False, e


def getPID(container_id):
	client = docker.from_env()
	container = client.containers.get(container_id)
	res = container.top()['Processes']
	for x in res:
		if "/workspace" in x[7]:
			print(res[1])
			break


# create_network()
# list_networks()

# remove_network()
# get_status('af121babda9b')
# exec_run()
# run()
# create_container()
# report()
getPID('a6543cef3c85')
