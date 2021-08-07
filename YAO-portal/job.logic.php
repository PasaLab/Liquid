<?php

require_once('predis/autoload.php');

require_once('util4p/util.php');
require_once('util4p/CRObject.class.php');
require_once('util4p/AccessController.class.php');
require_once('util4p/CRLogger.class.php');

require_once('Code.class.php');
require_once('JobManager.class.php');
require_once('Spider.class.php');

require_once('config.inc.php');
require_once('init.inc.php');

function job_submit(CRObject $job)
{
	if (!AccessController::hasAccess(Session::get('role', 'visitor'), 'job.submit')) {
		$res['errno'] = Code::NO_PRIVILEGE;
		return $res;
	}
	$job->set('created_by', Session::get('uid'));
	$job->set('created_at', time());

	/* notify YAO-scheduler */
	$spider = new Spider();
	$tasks = json_decode($job->get('tasks'), true);
	foreach ($tasks as $i => $task) {
		$task['cpu_number'] = intval($task['cpu_number']);
		$task['memory'] = intval($task['memory']);
		$task['gpu_number'] = intval($task['gpu_number']);
		$task['gpu_memory'] = intval($task['gpu_memory']);
		$task['is_ps'] = $task['is_ps'] == 1;
		$tasks[$i] = $task;
	}
	$job->set('tasks', $tasks);
	$job->set('workspace', $job->get('workspace'));
	$job->set('model_dir', $job->get('model_dir'));
	$job->set('output_dir', $job->get('output_dir'));
	$job->set('group', $job->get('virtual_cluster'));
	$job->set('priority', $job->getInt('priority'));
	$job->set('locality', $job->getInt('locality'));
	$job->set('run_before', $job->getInt('run_before'));
	$job->set('created_by', $job->getInt('created_by'));
	$data['job'] = json_encode($job);
	$spider->doPost(YAO_SCHEDULER_ADDR . '?action=job_submit', $data);
	$msg = json_decode($spider->getBody(), true);

	$res['errno'] = Code::SUCCESS;
	if ($msg['code'] !== 0) {
		$res['errno'] = Code::FAIL;
		$res['msg'] = $msg['error'];
	} else {
		$res['job_name'] = $msg['jobName'];
	}

	$log = new CRObject();
	$log->set('scope', Session::get('uid'));
	$log->set('tag', 'job.submit');
	$content = array('job' => $job, 'response' => $res['errno']);
	$log->set('content', json_encode($content));
	CRLogger::log($log);
	return $res;
}

function job_stop(CRObject $job)
{
	if (!AccessController::hasAccess(Session::get('role', 'visitor'), 'job.stop')) {
		$res['errno'] = Code::NO_PRIVILEGE;
		return $res;
	}
	/* TODO: permission check */
	$spider = new Spider();
	$data['id'] = $job->get('id', '');
	$spider->doPost(YAO_SCHEDULER_ADDR . '?action=job_stop', $data);
	$msg = json_decode($spider->getBody(), true);

	$res['errno'] = Code::SUCCESS;
	if ($msg['code'] !== 0) {
		$res['errno'] = Code::FAIL;
		$res['msg'] = $msg['error'];
		return $res;
	}
	$log = new CRObject();
	$log->set('scope', Session::get('uid'));
	$log->set('tag', 'job.stop');
	$content = array('id' => $job->get('id'), 'response' => $res['errno']);
	$log->set('content', json_encode($content));
	CRLogger::log($log);
	return $res;
}

function job_list(CRObject $rule)
{
	if (!AccessController::hasAccess(Session::get('role', 'visitor'), 'job.list')) {
		$res['errno'] = Code::NO_PRIVILEGE;
		return $res;
	}
	if ($rule->get('who') !== 'all') {
		$rule->set('who', 'self');
		$rule->set('created_by', Session::get('uid'));
	}
	if ($rule->get('who') === 'all' && !AccessController::hasAccess(Session::get('role', 'visitor'), 'job.list_others')) {
		$res['errno'] = Code::NO_PRIVILEGE;
		return $res;
	}


	$spider = new Spider();
	$spider->doGet(YAO_SCHEDULER_ADDR . '?action=jobs');
	$msg = json_decode($spider->getBody(), true);


	if ($msg['code'] !== 0) {
		$res['errno'] = $msg['code'] !== null ? $msg['code'] : Code::UNKNOWN_ERROR;
		$res['msg'] = $msg['error'];
		return $res;
	}

	if ($msg['jobs'] !== null) {
		$res['jobs'] = array_reverse($msg['jobs']);
	} else {
		$res['jobs'] = [];
	}
	for ($i = 0; $i < sizeof($res['jobs']); $i++) {
		$res['jobs'][$i]['tasks'] = json_encode($res['jobs'][$i]['tasks']);
		if ($res['jobs'][$i]['run_before'] === 0) {
			$res['jobs'][$i]['run_before'] = null;
		}
	}
	$res['errno'] = Code::SUCCESS;
	return $res;
}

function job_status(CRObject $job)
{
	if (!AccessController::hasAccess(Session::get('role', 'visitor'), 'job.list')) {
		$res['errno'] = Code::NO_PRIVILEGE;
		return $res;
	}

	$spider = new Spider();
	$spider->doGet(YAO_SCHEDULER_ADDR . '?action=job_status&id=' . $job->get('name'));
	$msg = json_decode($spider->getBody(), true);

	if ($msg['code'] !== 0) {
		$res['errno'] = $msg['code'];
		$res['msg'] = $msg['error'];
		return $res;
	}

	$res['tasks'] = $msg['status'];
	$res['errno'] = Code::SUCCESS;
	return $res;
}

function job_predict_req(CRObject $job, $role)
{
	if (!AccessController::hasAccess(Session::get('role', 'visitor'), 'job.list')) {
		$res['errno'] = Code::NO_PRIVILEGE;
		return $res;
	}

	$spider = new Spider();
	$tasks = json_decode($job->get('tasks'), true);
	foreach ($tasks as $i => $task) {
		$task['cpu_number'] = intval($task['cpu_number']);
		$task['memory'] = intval($task['memory']);
		$task['gpu_number'] = intval($task['gpu_number']);
		$task['gpu_memory'] = intval($task['gpu_memory']);
		$task['is_ps'] = $task['is_ps'] == 1;
		$tasks[$i] = $task;
	}
	$job->set('tasks', $tasks);
	$job->set('workspace', $job->get('workspace'));
	$job->set('group', $job->get('virtual_cluster'));
	$job->set('priority', $job->getInt('priority'));
	$job->set('locality', $job->getInt('locality'));
	$job->set('run_before', $job->getInt('run_before'));
	$job->set('created_by', $job->getInt('created_by'));
	$data['job'] = json_encode($job);

	$spider->doPost(YAO_SCHEDULER_ADDR . '?action=job_predict_req&role=' . $role, $data);
	$msg = json_decode($spider->getBody(), true);

	if ($msg === NULL) {
		$res['errno'] = Code::UNKNOWN_ERROR;
		$res['msg'] = 'response is null';
		return $res;
	}

	if ($msg['code'] !== 0) {
		$res['errno'] = $msg['code'];
		$res['msg'] = $msg['error'];
		return $res;
	}

	$res['cpu'] = $msg['cpu'];
	$res['mem'] = $msg['mem'];
	$res['gpu_util'] = $msg['gpu_util'];
	$res['gpu_mem'] = $msg['gpu_mem'];
	$res['bw'] = $msg['bw'];
	$res['errno'] = Code::SUCCESS;
	return $res;
}

function job_predict_time(CRObject $job, $role)
{
	if (!AccessController::hasAccess(Session::get('role', 'visitor'), 'job.list')) {
		$res['errno'] = Code::NO_PRIVILEGE;
		return $res;
	}

	$spider = new Spider();
	$tasks = json_decode($job->get('tasks'), true);
	foreach ($tasks as $i => $task) {
		$task['cpu_number'] = intval($task['cpu_number']);
		$task['memory'] = intval($task['memory']);
		$task['gpu_number'] = intval($task['gpu_number']);
		$task['gpu_memory'] = intval($task['gpu_memory']);
		$task['is_ps'] = $task['is_ps'] == 1;
		$tasks[$i] = $task;
	}
	$job->set('tasks', $tasks);
	$job->set('workspace', $job->get('workspace'));
	$job->set('group', $job->get('virtual_cluster'));
	$job->set('priority', $job->getInt('priority'));
	$job->set('locality', $job->getInt('locality'));
	$job->set('run_before', $job->getInt('run_before'));
	$job->set('created_by', $job->getInt('created_by'));
	$data['job'] = json_encode($job);

	$spider->doPost(YAO_SCHEDULER_ADDR . '?action=job_predict_time', $data);
	$msg = json_decode($spider->getBody(), true);

	if ($msg === NULL) {
		$res['errno'] = Code::UNKNOWN_ERROR;
		$res['msg'] = 'response is null';
		return $res;
	}

	if ($msg['code'] !== 0) {
		$res['errno'] = $msg['code'];
		$res['msg'] = $msg['error'];
		return $res;
	}

	$res['total'] = $msg['total'];
	$res['pre'] = $msg['pre'];
	$res['post'] = $msg['post'];
	$res['errno'] = Code::SUCCESS;
	return $res;
}

function summary_get()
{
	if (!AccessController::hasAccess(Session::get('role', 'visitor'), 'system.summary')) {
		$res['errno'] = Code::NO_PRIVILEGE;
		return $res;
	}

	$spider = new Spider();
	$spider->doGet(YAO_SCHEDULER_ADDR . '?action=summary');
	$msg = json_decode($spider->getBody(), true);

	if ($msg['code'] !== 0) {
		$res['errno'] = $msg['code'] !== null ? $msg['code'] : Code::UNKNOWN_ERROR;
		$res['msg'] = $msg['error'];
		return $res;
	}

	$res['jobs']['finished'] = $msg['jobs_finished'];
	$res['jobs']['running'] = $msg['jobs_running'];
	$res['jobs']['pending'] = $msg['jobs_pending'];
	$res['gpu']['free'] = $msg['gpu_free'];
	$res['gpu']['using'] = $msg['gpu_using'];
	$res['errno'] = Code::SUCCESS;
	return $res;
}

function summary_get_pool_history()
{
	if (!AccessController::hasAccess(Session::get('role', 'visitor'), 'system.summary')) {
		$res['errno'] = Code::NO_PRIVILEGE;
		return $res;
	}

	$spider = new Spider();
	$spider->doGet(YAO_SCHEDULER_ADDR . '?action=pool_status_history');
	$msg = json_decode($spider->getBody(), true);

	if ($msg['code'] !== 0) {
		$res['errno'] = $msg['code'] !== null ? $msg['code'] : Code::UNKNOWN_ERROR;
		$res['msg'] = $msg['error'];
		return $res;
	}

	$res['data'] = $msg['data'];
	$res['errno'] = Code::SUCCESS;
	return $res;
}

function task_logs(CRObject $job)
{
	if (!AccessController::hasAccess(Session::get('role', 'visitor'), 'job.list')) {
		$res['errno'] = Code::NO_PRIVILEGE;
		return $res;
	}

	$spider = new Spider();
	$spider->doGet(YAO_SCHEDULER_ADDR . '?action=task_logs&job=' . $job->get('job') . '&task=' . $job->get('task'));
	$msg = json_decode($spider->getBody(), true);

	if ($msg['code'] !== 0) {
		$res['errno'] = $msg['code'];
		$res['msg'] = $msg['error'];
		return $res;
	}

	$res['logs'] = $msg['logs'];
	$res['errno'] = Code::SUCCESS;
	return $res;
}

function job_describe(CRObject $rule)
{
	if (!AccessController::hasAccess(Session::get('role', 'visitor'), 'job.describe')) {
		$res['errno'] = Code::NO_PRIVILEGE;
		return $res;
	}
	$res['errno'] = Code::FAIL;
	$origin = JobManager::get($rule);
	if ($origin === null) {
		$res['errno'] = Code::RECORD_NOT_EXIST;
	} else if ($origin['created_by'] !== Session::get('uid') && !AccessController::hasAccess(Session::get('role', 'visitor'), 'job.describe_others')) {
		$res['errno'] = Code::NO_PRIVILEGE;
	}
	return $res;
}

function conf_update($option, $value)
{
	if (!AccessController::hasAccess(Session::get('role', 'visitor'), 'job.list')) {
		$res['errno'] = Code::NO_PRIVILEGE;
		return $res;
	}

	$spider = new Spider();
	$spider->doGet(YAO_SCHEDULER_ADDR . '?action=conf_update&option=' . $option . '&value=' . $value);
	$msg = json_decode($spider->getBody(), true);

	if ($msg['code'] !== 0) {
		$res['errno'] = $msg['code'];
		$res['msg'] = $msg['error'];
		return $res;
	}

	$res['errno'] = Code::SUCCESS;
	return $res;
}

function conf_list()
{
	if (!AccessController::hasAccess(Session::get('role', 'visitor'), 'job.list')) {
		$res['errno'] = Code::NO_PRIVILEGE;
		return $res;
	}

	$spider = new Spider();
	$spider->doGet(YAO_SCHEDULER_ADDR . '?action=conf_list');
	$msg = json_decode($spider->getBody(), true);

	if ($msg['code'] !== 0) {
		$res['errno'] = $msg['code'];
		$res['msg'] = $msg['error'];
		return $res;
	}

	$res['errno'] = Code::SUCCESS;
	$res['options'] = $msg['options'];
	return $res;
}