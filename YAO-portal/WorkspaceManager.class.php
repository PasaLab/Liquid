<?php

require_once('util4p/CRObject.class.php');
require_once('util4p/MysqlPDO.class.php');
require_once('util4p/SQLBuilder.class.php');

class WorkspaceManager
{

	private static $table = 'yao_workspace';

	/*
	 * do add workspace
	 */
	public static function add(CRObject $workspace)
	{
		$name = $workspace->get('name');
		$type = $workspace->get('type');
		$git_repo = $workspace->get('git_repo');
		$created_at = $workspace->getInt('created_at', time());
		$updated_at = $workspace->getInt('updated_at', time());
		$created_by = $workspace->get('created_by');

		$key_values = array(
			'name' => '?', 'type' => '?', 'git_repo' => '?',
			'created_at' => '?', 'updated_at' => '?', 'created_by' => '?'
		);
		$builder = new SQLBuilder();
		$builder->insert(self::$table, $key_values);
		$sql = $builder->build();
		$params = array($name, $type, $git_repo, $created_at, $updated_at, $created_by);
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
		$order_by = array('id' => 'DESC');
		$builder = new SQLBuilder();
		$builder->select(self::$table, $selected_rows);
		$builder->where($where);
		$builder->order($order_by);
		$builder->limit($offset, $limit);
		$sql = $builder->build();
		$workspaces = (new MysqlPDO())->executeQuery($sql, $params);
		return $workspaces;
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

	/* get workspace by ip */
	public static function getByID($id)
	{
		$selected_rows = array();
		$where = array('id' => '?');
		$params = array($id);
		$builder = new SQLBuilder();
		$builder->select(self::$table, $selected_rows);
		$builder->where($where);
		$sql = $builder->build();
		$workspaces = (new MysqlPDO())->executeQuery($sql, $params);
		return $workspaces !== null && count($workspaces) === 1 ? $workspaces[0] : null;
	}

	/*
	 * do update workspace
	 */
	public static function update(CRObject $workspace)
	{
		$id = $workspace->getInt('id');
		$name = $workspace->get('name');
		$type = $workspace->get('type');
		$git_repo = $workspace->get('git_repo');
		$updated_at = $workspace->getInt('updated_at', time());

		$key_values = array(
			'name' => '?', 'type' => '?', 'git_repo' => '?', 'updated_at' => '?'
		);
		$where = array('id' => '?');
		$builder = new SQLBuilder();
		$builder->update(self::$table, $key_values);
		$builder->where($where);
		$sql = $builder->build();
		$params = array($name, $type, $git_repo, $updated_at, $id);
		return (new MysqlPDO())->execute($sql, $params);
	}

	/* */
	public static function remove(CRObject $workspace)
	{
		$id = $workspace->getInt('id');
		$where = array('id' => '?');
		$builder = new SQLBuilder();
		$builder->delete(self::$table);
		$builder->where($where);
		$sql = $builder->build();
		$params = array($id);
		return (new MysqlPDO())->execute($sql, $params);
	}

}
