<?php

require_once('predis/autoload.php');

require_once('util4p/util.php');
require_once('util4p/CRObject.class.php');
require_once('util4p/AccessController.class.php');
require_once('util4p/CRLogger.class.php');

require_once('Code.class.php');
require_once('WorkspaceManager.class.php');

require_once('config.inc.php');
require_once('init.inc.php');

function workspace_add(CRObject $workspace)
{
	if (!AccessController::hasAccess(Session::get('role', 'visitor'), 'agent.add')) {
		$res['errno'] = Code::NO_PRIVILEGE;
		return $res;
	}
	if (WorkspaceManager::getByID($workspace->get('ip')) !== null) {
		$res['errno'] = Code::RECORD_ALREADY_EXIST;
	} else {
		$workspace->set('created_by', Session::get('uid'));
		$res['errno'] = WorkspaceManager::add($workspace) ? Code::SUCCESS : Code::UNKNOWN_ERROR;
	}
	$log = new CRObject();
	$log->set('scope', Session::get('uid'));
	$log->set('tag', 'workspace.add');
	$content = array('workspace' => $workspace, 'response' => $res['errno']);
	$log->set('content', json_encode($content));
	CRLogger::log($log);
	return $res;
}

function workspace_remove(CRObject $workspace)
{
	if (!AccessController::hasAccess(Session::get('role', 'visitor'), 'workspace.remove')) {
		$res['errno'] = Code::NO_PRIVILEGE;
		return $res;
	}
	$res['errno'] = WorkspaceManager::remove($workspace) ? Code::SUCCESS : Code::UNKNOWN_ERROR;
	$log = new CRObject();
	$log->set('scope', Session::get('uid'));
	$log->set('tag', 'workspace.remove');
	$content = array('workspace' => $workspace, 'response' => $res['errno']);
	$log->set('content', json_encode($content));
	CRLogger::log($log);
	return $res;
}

function workspace_update(CRObject $workspace)
{
	if (!AccessController::hasAccess(Session::get('role', 'visitor'), 'workspace.update')) {
		$res['errno'] = Code::NO_PRIVILEGE;
		return $res;
	}
	$res['errno'] = WorkspaceManager::update($workspace) ? Code::SUCCESS : Code::UNKNOWN_ERROR;
	$log = new CRObject();
	$log->set('scope', Session::get('uid'));
	$log->set('tag', 'workspace.update');
	$content = array('workspace' => $workspace, 'response' => $res['errno']);
	$log->set('content', json_encode($content));
	CRLogger::log($log);
	return $res;
}

function workspace_list(CRObject $rule)
{
	if (!AccessController::hasAccess(Session::get('role', 'visitor'), 'workspace.list')) {
		$res['errno'] = Code::NO_PRIVILEGE;
		return $res;
	}
	$res['workspaces'] = WorkspaceManager::gets($rule);
	$res['count'] = WorkspaceManager::count($rule);
	$res['errno'] = $res['workspaces'] === null ? Code::FAIL : Code::SUCCESS;
	return $res;
}
