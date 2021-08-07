<?php

require_once('predis/autoload.php');

require_once('util4p/util.php');
require_once('util4p/CRObject.class.php');
require_once('util4p/AccessController.class.php');
require_once('util4p/CRLogger.class.php');

require_once('Code.class.php');

require_once('config.inc.php');
require_once('init.inc.php');

function cluster_add(CRObject $cluster)
{
	if (!AccessController::hasAccess(Session::get('role', 'visitor'), 'cluster.add')) {
		$res['errno'] = Code::NO_PRIVILEGE;
		return $res;
	}

	$spider = new Spider();
	$cluster->set('weight', $cluster->getInt('weight', 0));
	$cluster->set('reserved', $cluster->getBool('reserved', false));
	$cluster->set('quota_cpu', $cluster->getInt('quota_cpu', 0));
	$cluster->set('quota_mem', $cluster->getInt('quota_mem', 0));
	$cluster->set('quota_gpu', $cluster->getInt('quota_gpu', 0));
	$cluster->set('quota_gpu_mem', $cluster->getInt('quota_gpu_mem', 0));
	$data['group'] = json_encode($cluster);
	$spider->doPost(YAO_SCHEDULER_ADDR . '?action=group_add', $data);
	$msg = json_decode($spider->getBody(), true);

	$res['errno'] = Code::SUCCESS;
	if ($msg['code'] !== 0) {
		$res['errno'] = Code::FAIL;
		$res['msg'] = $msg['error'];
	}

	$log = new CRObject();
	$log->set('scope', Session::get('uid'));
	$log->set('tag', 'cluster.add');
	$content = array('cluster' => $cluster, 'response' => $res['errno']);
	$log->set('content', json_encode($content));
	CRLogger::log($log);
	return $res;
}

function cluster_remove(CRObject $cluster)
{
	if (!AccessController::hasAccess(Session::get('role', 'visitor'), 'cluster.remove')) {
		$res['errno'] = Code::NO_PRIVILEGE;
		return $res;
	}
	// TODO: check owner

	$spider = new Spider();
	$data['group'] = json_encode($cluster);
	$spider->doPost(YAO_SCHEDULER_ADDR . '?action=group_remove', $data);
	$msg = json_decode($spider->getBody(), true);

	$res['errno'] = Code::SUCCESS;
	if ($msg['code'] !== 0) {
		$res['errno'] = Code::FAIL;
		$res['msg'] = $msg['error'];
	}

	$log = new CRObject();
	$log->set('scope', Session::get('uid'));
	$log->set('tag', 'cluster.remove');
	$content = array('cluster' => $cluster, 'response' => $res['errno']);
	$log->set('content', json_encode($content));
	CRLogger::log($log);
	return $res;
}

function cluster_update(CRObject $cluster)
{
	if (!AccessController::hasAccess(Session::get('role', 'visitor'), 'cluster.update')) {
		$res['errno'] = Code::NO_PRIVILEGE;
		return $res;
	}
	//TODO: check owner

	$spider = new Spider();
	$cluster->set('weight', $cluster->getInt('weight', 0));
	$cluster->set('reserved', $cluster->getBool('reserved', false));
	$cluster->set('quota_cpu', $cluster->getInt('quota_cpu', 0));
	$cluster->set('quota_mem', $cluster->getInt('quota_mem', 0));
	$cluster->set('quota_gpu', $cluster->getInt('quota_gpu', 0));
	$cluster->set('quota_gpu_mem', $cluster->getInt('quota_gpu_mem', 0));
	$data['group'] = json_encode($cluster);
	$spider->doPost(YAO_SCHEDULER_ADDR . '?action=group_update', $data);
	$msg = json_decode($spider->getBody(), true);

	$res['errno'] = Code::SUCCESS;
	if ($msg['code'] !== 0) {
		$res['errno'] = Code::FAIL;
		$res['msg'] = $msg['error'];
	}

	$log = new CRObject();
	$log->set('scope', Session::get('uid'));
	$log->set('tag', 'cluster.update');
	$content = array('cluster' => $cluster, 'response' => $res['errno']);
	$log->set('content', json_encode($content));
	CRLogger::log($log);
	return $res;
}

function cluster_list(CRObject $rule)
{
	if (!AccessController::hasAccess(Session::get('role', 'visitor'), 'cluster.list')) {
		$res['errno'] = Code::NO_PRIVILEGE;
		return $res;
	}
	$spider = new Spider();
	$spider->doGet(YAO_SCHEDULER_ADDR . '?action=group_list');
	$msg = json_decode($spider->getBody(), true);


	if ($msg['code'] !== 0) {
		$res['errno'] = $msg['code'] !== null ? $msg['code'] : Code::UNKNOWN_ERROR;
		$res['msg'] = $msg['error'];
		return $res;
	}

	$res['clusters'] = $msg['groups'];
	$res['errno'] = Code::SUCCESS;
	return $res;
}
