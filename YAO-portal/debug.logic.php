<?php

require_once('predis/autoload.php');

require_once('util4p/util.php');
require_once('util4p/CRObject.class.php');
require_once('util4p/Random.class.php');
require_once('util4p/AccessController.class.php');
require_once('util4p/CRLogger.class.php');

require_once('Code.class.php');

require_once('config.inc.php');
require_once('init.inc.php');


function debug_get_counter()
{
	if (!AccessController::hasAccess(Session::get('role', 'visitor'), 'system.summary')) {
		$res['errno'] = Code::NO_PRIVILEGE;
		return $res;
	}

	$spider = new Spider();
	$spider->doGet(YAO_SCHEDULER_ADDR . '?action=get_counter');
	$msg = json_decode($spider->getBody(), true);

	$res['counter'] = $msg['counter'];
	$res['counterTotal'] = $msg['counterTotal'];
	$res['errno'] = Code::SUCCESS;
	return $res;
}

function debug_get_bindings()
{
	if (!AccessController::hasAccess(Session::get('role', 'visitor'), 'system.summary')) {
		$res['errno'] = Code::NO_PRIVILEGE;
		return $res;
	}

	$spider = new Spider();
	$spider->doGet(YAO_SCHEDULER_ADDR . '?action=get_bindings');
	$msg = json_decode($spider->getBody(), true);

	$res['data'] = $msg;
	$res['errno'] = Code::SUCCESS;
	return $res;
}