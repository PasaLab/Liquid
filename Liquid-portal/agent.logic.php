<?php

require_once('predis/autoload.php');

require_once('util4p/util.php');
require_once('util4p/CRObject.class.php');
require_once('util4p/Random.class.php');
require_once('util4p/AccessController.class.php');
require_once('util4p/CRLogger.class.php');

require_once('Code.class.php');
require_once('AgentManager.class.php');

require_once('config.inc.php');
require_once('init.inc.php');

function agent_add(CRObject $agent)
{
	if (!AccessController::hasAccess(Session::get('role', 'visitor'), 'agent.add')) {
		$res['errno'] = Code::NO_PRIVILEGE;
		return $res;
	}
	if (AgentManager::getByIP($agent->get('ip')) !== null) {
		$res['errno'] = Code::RECORD_ALREADY_EXIST;
	} else {
		$token = Random::randomString(32);
		$agent->set('token', $token);
		$res['errno'] = AgentManager::add($agent) ? Code::SUCCESS : Code::UNKNOWN_ERROR;
	}
	$log = new CRObject();
	$log->set('scope', Session::get('uid'));
	$log->set('tag', 'agent.add');
	$content = array('agent' => $agent, 'response' => $res['errno']);
	$log->set('content', json_encode($content));
	CRLogger::log($log);
	return $res;
}

function agent_remove(CRObject $agent)
{
	if (!AccessController::hasAccess(Session::get('role', 'visitor'), 'agent.remove')) {
		$res['errno'] = Code::NO_PRIVILEGE;
		return $res;
	}
	$res['errno'] = AgentManager::remove($agent) ? Code::SUCCESS : Code::UNKNOWN_ERROR;
	$log = new CRObject();
	$log->set('scope', Session::get('uid'));
	$log->set('tag', 'agent.remove');
	$content = array('agent' => $agent, 'response' => $res['errno']);
	$log->set('content', json_encode($content));
	CRLogger::log($log);
	return $res;
}

function agent_list(CRObject $rule)
{
	if (!AccessController::hasAccess(Session::get('role', 'visitor'), 'agent.list')) {
		$res['errno'] = Code::NO_PRIVILEGE;
		return $res;
	}
	$res['agents'] = AgentManager::gets($rule);
	$res['count'] = AgentManager::count($rule);
	$res['errno'] = $res['agents'] === null ? Code::FAIL : Code::SUCCESS;
	return $res;
}

function resource_list()
{
	if (!AccessController::hasAccess(Session::get('role', 'visitor'), 'system.summary')) {
		$res['errno'] = Code::NO_PRIVILEGE;
		return $res;
	}

	$spider = new Spider();
	$spider->doGet(YAO_SCHEDULER_ADDR . '?action=resource_list');
	$msg = json_decode($spider->getBody(), true);

	if ($msg['code'] !== 0) {
		$res['errno'] = $msg['code'] !== null ? $msg['code'] : Code::UNKNOWN_ERROR;
		$res['msg'] = $msg['error'];
		return $res;
	}
	$res['resources'] = $msg['resources'];
	$res['errno'] = Code::SUCCESS;
	return $res;
}
