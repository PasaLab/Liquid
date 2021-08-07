<?php

require_once('util4p/CRObject.class.php');
require_once('util4p/MysqlPDO.class.php');
require_once('util4p/SQLBuilder.class.php');

class JobManager
{
	private static $table = 'yao_job';

	/*
	 * do submit job
	 */
	public static function add(CRObject $job)
	{
		$name = $job->get('name');
		$virtual_cluster = $job->getInt('virtual_cluster');
		$run_before = $job->getInt('run_before');
		$tasks = $job->get('tasks');
		$image = $job->get('image');
		$workspace = $job->getInt('workspace');
		$priority = $job->getInt('priority');
		$created_at = time();
		$created_by = $job->getInt('created_by');

		$key_values = array(
			'name' => '?', 'image' => '?', 'workspace' => '?', 'virtual_cluster' => '?', 'priority' => '?',
			'run_before' => '?', 'created_at' => '?', 'created_by' => '?', 'tasks' => '?'
		);
		$builder = new SQLBuilder();
		$builder->insert(self::$table, $key_values);
		$sql = $builder->build();
		$params = array($name, $image, $workspace, $virtual_cluster, $priority, $run_before, $created_at, $created_by, $tasks);
		return (new MysqlPDO())->execute($sql, $params);
	}

	/* */
	public static function gets(CRObject $rule)
	{
		$virtual_cluster = $rule->getInt('virtual_cluster', null);
		$status = $rule->getInt('status', null);
		$offset = $rule->getInt('offset', 0);
		$limit = $rule->getInt('limit', -1);
		$selected_rows = array();
		$where = array();
		$params = array();
		if ($virtual_cluster !== null) {
			$where['virtual_cluster'] = '?';
			$params[] = $virtual_cluster;
		}
		if ($status !== null) {
			$where['status'] = '?';
			$params[] = $status;
		}
		$order_by = array('created_at' => 'DESC');
		$builder = new SQLBuilder();
		$builder->select(self::$table, $selected_rows);
		$builder->where($where);
		$builder->order($order_by);
		$builder->limit($offset, $limit);
		$sql = $builder->build();
		$jobs = (new MysqlPDO())->executeQuery($sql, $params);
		return $jobs;
	}

	/* */
	public static function count(CRObject $rule)
	{
		$virtual_cluster = $rule->getInt('virtual_cluster', null);
		$status = $rule->getInt('status', null);
		$selected_rows = array('COUNT(1) as cnt');
		$where = array();
		$params = array();
		if ($virtual_cluster !== null) {
			$where['virtual_cluster'] = '?';
			$params[] = $virtual_cluster;
		}
		if ($status !== null) {
			$where['status'] = '?';
			$params[] = $status;
		}
		$builder = new SQLBuilder();
		$builder->select(self::$table, $selected_rows);
		$builder->where($where);
		$sql = $builder->build();
		$res = (new MysqlPDO())->executeQuery($sql, $params);
		return $res === null ? 0 : intval($res[0]['cnt']);
	}

	/* get job by id */
	public static function get(CRObject $rule)
	{
		$id = $rule->getInt('id');
		$selected_rows = array();
		$where = array('id' => '?');
		$params = array($id);
		$builder = new SQLBuilder();
		$builder->select(self::$table, $selected_rows);
		$builder->where($where);
		$sql = $builder->build();
		$jobs = (new MysqlPDO())->executeQuery($sql, $params);
		return $jobs !== null && count($jobs) === 1 ? $jobs[0] : null;
	}

	/* */
	public static function update(CRObject $link)
	{
		$id = $link->getId('id');
		$url = $link->get('url', '');
		$remark = $link->get('remark');
		$valid_from = $link->getInt('valid_from');
		$valid_to = $link->getInt('valid_to');
		$status = $link->getInt('status', 0);

		$key_values = array(
			'url' => '?', 'remark' => '?', 'valid_from' => '?', 'valid_to' => '?', 'status' => '?'
		);
		$where = array('token' => '?');
		$builder = new SQLBuilder();
		$builder->update(self::$table, $key_values);
		$builder->where($where);
		$sql = $builder->build();
		$params = array($url, $remark, $valid_from, $valid_to, $status, $token);
		return (new MysqlPDO())->execute($sql, $params);
	}

}
