<?php

require_once('predis/autoload.php');
require_once('util4p/MysqlPDO.class.php');
require_once('util4p/RedisDAO.class.php');
require_once('util4p/CRLogger.class.php');
require_once('util4p/ReSession.class.php');
require_once('util4p/CRObject.class.php');
require_once('util4p/AccessController.class.php');

require_once('config.inc.php');

init_mysql();
init_redis();
init_logger();
init_Session();
init_accessMap();

function init_mysql()
{
	$config = new CRObject();
	$config->set('host', DB_HOST);
	$config->set('port', DB_PORT);
	$config->set('db', DB_NAME);
	$config->set('user', DB_USER);
	$config->set('password', DB_PASSWORD);
	$config->set('show_error', DB_SHOW_ERROR);
	MysqlPDO::configure($config);
}

function init_redis()
{
	$config = new CRObject();
	$config->set('scheme', REDIS_SCHEME);
	$config->set('host', REDIS_HOST);
	$config->set('port', REDIS_PORT);
	$config->set('show_error', REDIS_SHOW_ERROR);
	RedisDAO::configure($config);
}

function init_logger()
{
	$config = new CRObject();
	$config->set('db_table', 'yao_log');
	CRLogger::configure($config);
}

function init_Session()
{
	$config = new CRObject();
	$config->set('time_out', SESSION_TIME_OUT);
	$config->set('bind_ip', BIND_SESSION_WITH_IP);
	$config->set('PK', 'username');
	Session::configure($config);
}

function init_accessMap()
{
	// $operation => array of roles
	$map = array(
		/* user */
		'user.get' => array('root', 'admin', 'developer', 'normal'),
		'user.get_others' => array('root', 'admin'),

		/* logs */
		'logs.get' => array('root', 'admin', 'developer', 'normal'),
		'logs.get_others' => array('root', 'admin'),

		/* job */
		'job.list' => array('root', 'admin', 'developer', 'normal'),
		'job.submit' => array('root', 'admin', 'developer', 'normal'),
		'job.stop' => array('root', 'admin', 'developer', 'normal'),
		'job.stop_others' => array('root', 'admin', 'developer', 'normal'),

		/* system */
		'system.summary' => array('root', 'admin', 'developer', 'normal'),

		/* agent */
		'agent.list' => array('root', 'admin', 'normal'),
		'agent.add' => array('root', 'admin'),
		'agent.remove' => array('root', 'admin'),

		/* workspace */
		'workspace.list' => array('root', 'admin', 'normal'),
		'workspace.add' => array('root', 'admin', 'normal'),
		'workspace.update' => array('root', 'admin', 'normal'),
		'workspace.remove' => array('root', 'admin', 'normal'),

		/* cluster */
		'cluster.list' => array('root', 'admin', 'normal'),
		'cluster.add' => array('root', 'admin'),
		'cluster.update' => array('root', 'admin'),
		'cluster.remove' => array('root', 'admin'),

		/* ucenter entry show control */
		'ucenter.home' => array('root', 'admin', 'developer', 'normal'),
		'ucenter.jobs' => array('root', 'admin', 'developer', 'normal'),
		'ucenter.workspaces' => array('root', 'admin', 'developer', 'normal'),
		'ucenter.jobs_all' => array('root', 'admin'),
		'ucenter.workspaces_all' => array('root', 'admin'),
		'ucenter.logs' => array('root', 'admin', 'developer', 'normal'),
		'ucenter.logs_all' => array('root', 'admin'),
		'ucenter.agents' => array('root', 'admin'),
		'ucenter.clusters' => array('root', 'admin'),
		'ucenter.admin' => array('root', 'admin'),
		'ucenter.users' => array('root', 'admin'),
		'ucenter.resources' => array('root', 'admin'),
		'ucenter.summary' => array('root', 'admin'),
		'ucenter.visitors' => array('root', 'admin')
	);
	AccessController::setMap($map);
}