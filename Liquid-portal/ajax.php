<?php

require_once('util4p/util.php');
require_once('util4p/CRObject.class.php');
require_once('util4p/Random.class.php');

require_once('Code.class.php');
require_once('Securer.class.php');

require_once('user.logic.php');
require_once('job.logic.php');
require_once('agent.logic.php');
require_once('workspace.logic.php');
require_once('cluster.logic.php');
require_once('debug.logic.php');

require_once('config.inc.php');
require_once('init.inc.php');


function csrf_check($action)
{
	/* check referer, just in case I forget to add the method to $post_methods */
	$referer = cr_get_SERVER('HTTP_REFERER', '');
	$url = parse_url($referer);
	$host = isset($url['host']) ? $url['host'] : '';
	$host .= isset($url['port']) && $url['port'] !== 80 ? ':' . $url['port'] : '';
	if ($host !== cr_get_SERVER('HTTP_HOST')) {
		return false;
	}
	$post_methods = array(
		'signout',
		'oauth_get_url'
	);
	if (in_array($action, $post_methods)) {
		return Securer::validate_csrf_token();
	}
	return true;
}

function print_response($res)
{
	if (!isset($res['msg']))
		$res['msg'] = Code::getErrorMsg($res['errno']);
	$json = json_encode($res);
	header('Content-type: application/json');
	echo $json;
}


$res = array('errno' => Code::UNKNOWN_REQUEST);

$action = cr_get_GET('action');

if (!csrf_check($action)) {
	$res['errno'] = 99;
	$res['msg'] = 'invalid csrf_token';
	print_response($res);
	exit(0);
}

switch ($action) {
	case 'job_list':
		$rule = new CRObject();
		$rule->set('who', cr_get_GET('who', 'self'));
		$rule->set('offset', cr_get_GET('offset'));
		$rule->set('limit', cr_get_GET('limit'));
		$rule->set('order', 'latest');
		$res = job_list($rule);
		break;

	case 'job_submit':
		$job = new CRObject();
		$job->set('name', cr_get_POST('name', 'jobName'));
		$job->set('virtual_cluster', cr_get_POST('cluster'));
		$job->set('workspace', cr_get_POST('workspace'));
		$job->set('model_dir', cr_get_POST('model_dir'));
		$job->set('output_dir', cr_get_POST('output_dir'));
		$job->set('priority', cr_get_POST('priority'));
		$job->set('run_before', cr_get_POST('run_before'));
		$job->set('locality', cr_get_POST('locality'));
		$job->set('tasks', cr_get_POST('tasks'));
		$res = job_submit($job);
		break;

	case 'job_stop':
		$job = new CRObject();
		$job->set('id', cr_get_POST('id'));
		$res = job_stop($job);
		break;

	case 'job_describe':
		$job = new CRObject();
		$job->set('id', cr_get_POST('id'));
		$res = job_describe($job);
		break;

	case 'job_status':
		$job = new CRObject();
		$job->set('name', cr_get_GET('name'));
		$res = job_status($job);
		break;

	case 'job_predict_req':
		$job = new CRObject();
		$job->set('name', cr_get_POST('name', 'jobName'));
		$job->set('virtual_cluster', cr_get_POST('cluster'));
		$job->set('workspace', cr_get_POST('workspace'));
		$job->set('priority', cr_get_POST('priority'));
		$job->set('run_before', cr_get_POST('run_before'));
		$job->set('locality', cr_get_POST('locality'));
		$job->set('tasks', cr_get_POST('tasks'));
		$role = cr_get_GET('role');
		$res = job_predict_req($job, $role);
		break;

	case 'job_predict_time':
		$job = new CRObject();
		$job->set('name', cr_get_POST('name', 'jobName'));
		$job->set('virtual_cluster', cr_get_POST('cluster'));
		$job->set('workspace', cr_get_POST('workspace'));
		$job->set('priority', cr_get_POST('priority'));
		$job->set('run_before', cr_get_POST('run_before'));
		$job->set('locality', cr_get_POST('locality'));
		$job->set('tasks', cr_get_POST('tasks'));
		$res = job_predict_time($job, $role);
		break;

	case 'summary_get':
		$res = summary_get();
		break;

	case 'summary_get_pool_history':
		$res = summary_get_pool_history();
		break;

	case 'task_logs':
		$task = new CRObject();
		$task->set('job', cr_get_GET('job'));
		$task->set('task', cr_get_GET('task'));
		$res = task_logs($task);
		break;

	case 'agent_list':
		$rule = new CRObject();
		$rule->set('offset', cr_get_GET('offset'));
		$rule->set('limit', cr_get_GET('limit'));
		$res = agent_list($rule);
		break;

	case 'agent_add':
		$agent = new CRObject();
		$agent->set('ip', cr_get_POST('ip'));
		$agent->set('alias', cr_get_POST('alias'));
		$agent->set('cluster', cr_get_POST('cluster'));
		$res = agent_add($agent);
		break;

	case 'agent_remove':
		$job = new CRObject();
		$job->set('id', cr_get_POST('id'));
		$res = agent_remove($job);
		break;

	case 'resource_list':
		$res = resource_list();
		break;

	case 'workspace_list':
		$rule = new CRObject();
		$rule->set('offset', cr_get_GET('offset'));
		$rule->set('limit', cr_get_GET('limit'));
		$res = workspace_list($rule);
		break;

	case 'workspace_add':
		$workspace = new CRObject();
		$workspace->set('name', cr_get_POST('name'));
		$workspace->set('type', cr_get_POST('type'));
		$workspace->set('git_repo', cr_get_POST('git_repo'));
		$res = workspace_add($workspace);
		break;

	case 'workspace_update':
		$workspace = new CRObject();
		$workspace->set('id', cr_get_POST('id'));
		$workspace->set('name', cr_get_POST('name'));
		$workspace->set('type', cr_get_POST('type'));
		$workspace->set('git_repo', cr_get_POST('git_repo'));
		$res = workspace_update($workspace);
		break;

	case 'workspace_remove':
		$workspace = new CRObject();
		$workspace->set('id', cr_get_POST('id'));
		$res = workspace_remove($workspace);
		break;

	case 'cluster_list':
		$rule = new CRObject();
		$rule->set('offset', cr_get_GET('offset'));
		$rule->set('limit', cr_get_GET('limit'));
		$rule->set('order', 'latest');
		$res = cluster_list($rule);
		break;

	case 'cluster_add':
		$cluster = new CRObject();
		$cluster->set('name', cr_get_POST('name'));
		$cluster->set('weight', cr_get_POST('weight'));
		$cluster->set('reserved', cr_get_POST('reserved'));
		$cluster->set('quota_cpu', cr_get_POST('quota_cpu'));
		$cluster->set('quota_mem', cr_get_POST('quota_mem'));
		$cluster->set('quota_gpu', cr_get_POST('quota_gpu'));
		$cluster->set('quota_gpu_mem', cr_get_POST('quota_gpu_mem'));
		$res = cluster_add($cluster);
		break;

	case 'cluster_update':
		$cluster = new CRObject();
		$cluster->set('name', cr_get_POST('name'));
		$cluster->set('weight', cr_get_POST('weight'));
		$cluster->set('reserved', cr_get_POST('reserved'));
		$cluster->set('quota_cpu', cr_get_POST('quota_cpu'));
		$cluster->set('quota_mem', cr_get_POST('quota_mem'));
		$cluster->set('quota_gpu', cr_get_POST('quota_gpu'));
		$cluster->set('quota_gpu_mem', cr_get_POST('quota_gpu_mem'));
		$res = cluster_update($cluster);
		break;

	case 'cluster_remove':
		$cluster = new CRObject();
		$cluster->set('name', cr_get_POST('name'));
		$res = cluster_remove($cluster);
		break;

	case 'get_counter':
		$res = debug_get_counter();
		break;

	case 'get_bindings':
		$res = debug_get_bindings();
		break;

	case 'conf_update':
		$option = cr_get_POST('option', '');
		$value = cr_get_POST('value', '');
		$res = conf_update($option, $value);
		break;

	case 'conf_list':
		$res = conf_list();
		break;

	case 'user_signout':
		$res = user_signout();
		break;

	case 'log_gets':
		$rule = new CRObject();
		$rule->set('who', cr_get_GET('who', 'self'));
		$rule->set('offset', cr_get_GET('offset'));
		$rule->set('limit', cr_get_GET('limit'));
		$rule->set('order', 'latest');
		$res = log_gets($rule);
		break;

	case 'oauth_get_url':
		$res = oauth_get_url();
		break;

	case 'user_login':
		$user = new CRObject();
		$res = user_login($user);
		break;

	default:
		break;
}

print_response($res);
