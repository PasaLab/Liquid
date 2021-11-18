import requests
import json


class API:
	def __init__(self, base_url):
		self.BASE_URL = base_url
		self.sess = requests.Session()
		self.sess.headers.update({'Referer': self.BASE_URL})

	def login(self, user='', pwd=''):
		# Get CSRF Token
		r = self.sess.get(self.BASE_URL)
		# print(r.content)
		# print('[api.login.get_csrf]', r.status_code)

		# Login
		url = self.BASE_URL + '/service?action=user_login'
		r = self.sess.post(url, data={})
		data = str(r.content, 'utf-8')
		msg = json.loads(data)
		# print('[api.login]', msg)
		return msg

	def get_sys_status(self):
		# Retrieve Status
		r = self.sess.get(self.BASE_URL + '/service?action=summary_get')
		# print(r.content)
		summary = str(r.content, 'utf-8')
		# b'{"jobs":{"finished":1,"running":0,"pending":0},"gpu":{"free":20,"using":0},"errno":0,"msg":"Success !"}'

		# Get pool Util history
		r = self.sess.get(self.BASE_URL + '/service?action=summary_get_pool_history')
		pool_history = str(r.content, 'utf-8')
		# print(r.content)
		return summary, pool_history

	def create_queue(self, queue):
		r = self.sess.post(self.BASE_URL + '/service?action=cluster_add', data=queue)
		data = str(r.content, 'utf-8')
		msg = json.loads(data)
		# print('[api.create_queue]', msg['errno'])
		return msg

	def submit_job(self, job):
		r = self.sess.post(self.BASE_URL + '/service?action=job_submit', data=job)
		data = str(r.content, 'utf-8')
		msg = json.loads(data)
		# print('[api.submit_job]', msg['errno'])
		return msg

	def job_list(self):
		r = self.sess.get(self.BASE_URL + '/service?action=job_list&who=self&sort=nobody&order=desc&offset=0&limit=10')
		data = str(r.content, 'utf-8')
		msg = json.loads(data)
		# print('[api.job_list]', msg['errno'])
		return msg

	def job_status(self, job_name):
		r = self.sess.get(self.BASE_URL + '/service?action=job_status&name=' + job_name)
		data = str(r.content, 'utf-8')
		msg = json.loads(data)
		# print('[api.job_status]', msg['errno'])
		return msg

	def conf_update(self, option=None, value=None):
		data = {'option': option, 'value': str(value)}
		r = self.sess.post(self.BASE_URL + '/service?action=conf_update', data=data)
		data = str(r.content, 'utf-8')
		msg = json.loads(data)
		# print('[api.update_conf]', msg['errno'])
		return msg

	def conf_list(self):
		r = self.sess.get(self.BASE_URL + '/service?action=conf_list')
		data = str(r.content, 'utf-8')
		msg = json.loads(data)
		# print('[api.update_conf]', msg['errno'])
		return msg

	def conf_reset(self):
		default_conf = {
			'pool.share.enable_threshold': '1.5',
			'pool.share.max_utilization': '100',
			'pool.pre_schedule.enable_threshold': '1.5',
			'pool.batch.enabled': 'false',
			'pool.batch.interval': '30',
			'scheduler.mock.enabled': 'false',
			'scheduler.enabled': 'true',
			'scheduler.parallelism': '1',
			'allocator.strategy': 'bestfit',
			'scheduler.preempt_enabled': 'false',
		}
		for option in default_conf:
			value = default_conf[option]
			msg = self.conf_update(option, value)
			if msg['errno'] != 0:
				print('[api.update_conf] option={} errno={} msg={}'.format(option, msg['errno'], msg['msg']))
		return


if __name__ == '__main__':
	print('util.api')

