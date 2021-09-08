from hdfs import *
import os
import time


if __name__ == '__main__':
  os.environ["TZ"] = 'Asia/Shanghai'
  if hasattr(time, 'tzset'):
    time.tzset()
  try:
    hdfs_address = os.environ['hdfs_address']
    hdfs_dir = os.environ['hdfs_dir']
    output_dir = os.environ['output_dir']
    
    client = Client(hdfs_address)
    client.upload(hdfs_dir, output_dir)

    print('Save ' + output_dir + ' to' + hdfs_address + ' ' + hdfs_dir)
  except Exception as e:
    print('Unable to persist data to HDFS,', str(e))
