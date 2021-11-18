import util.api as api

apiInstance = api.API(base_url='http://yao.pasalab.jluapp.com')
apiInstance.login()

print(apiInstance.conf_list())
print(apiInstance.get_sys_status())