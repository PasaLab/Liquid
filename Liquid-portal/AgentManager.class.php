<?php

require_once('util4p/CRObject.class.php');
require_once('util4p/MysqlPDO.class.php');
require_once('util4p/SQLBuilder.class.php');

class AgentManager
{

	private static $table = 'yao_agent';

	/*
	 * do add agent
	 */
	public static function add(CRObject $agent)
	{
		$ip = $agent->get('ip');
		$alias = $agent->get('alias');
		$cluster = $agent->getInt('cluster');
		$token = $agent->get('token');

		$key_values = array(
			'ip' => '?', 'alias' => '?', 'cluster' => '?', 'token' => '?'
		);
		$builder = new SQLBuilder();
		$builder->insert(self::$table, $key_values);
		$sql = $builder->build();
		$params = array(ip2long($ip), $alias, $cluster, $token);
		return (new MysqlPDO())->execute($sql, $params);
	}

	/* */
	public static function gets(CRObject $rule)
	{
		$offset = $rule->getInt('offset', 0);
		$limit = $rule->getInt('limit', -1);
		$selected_rows = array();
		$where = array();
		$params = array();
		$order_by = array('ip' => 'ASC');
		$builder = new SQLBuilder();
		$builder->select(self::$table, $selected_rows);
		$builder->where($where);
		$builder->order($order_by);
		$builder->limit($offset, $limit);
		$sql = $builder->build();
		$agents = (new MysqlPDO())->executeQuery($sql, $params);
		return $agents;
	}

	/* */
	public static function count(CRObject $rule)
	{
		$selected_rows = array('COUNT(1) as cnt');
		$where = array();
		$params = array();
		$builder = new SQLBuilder();
		$builder->select(self::$table, $selected_rows);
		$builder->where($where);
		$sql = $builder->build();
		$res = (new MysqlPDO())->executeQuery($sql, $params);
		return $res === null ? 0 : intval($res[0]['cnt']);
	}

	/* get agent by ip */
	public static function getByIP($ip)
	{
		$selected_rows = array();
		$where = array('ip' => '?');
		$params = array(ip2long($ip));
		$builder = new SQLBuilder();
		$builder->select(self::$table, $selected_rows);
		$builder->where($where);
		$sql = $builder->build();
		$agents = (new MysqlPDO())->executeQuery($sql, $params);
		return $agents !== null && count($agents) === 1 ? $agents[0] : null;
	}

	/* */
	public static function remove(CRObject $agent)
	{
		$id = $agent->getInt('id');
		$where = array('id' => '?');
		$builder = new SQLBuilder();
		$builder->delete(self::$table);
		$builder->where($where);
		$sql = $builder->build();
		$params = array($id);
		return (new MysqlPDO())->execute($sql, $params);
	}

}
