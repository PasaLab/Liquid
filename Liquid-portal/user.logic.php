<?php

require_once('predis/autoload.php');

require_once('util4p/CRObject.class.php');
require_once('util4p/ReSession.class.php');
require_once('util4p/CRLogger.class.php');

require_once('Code.class.php');
require_once('UserManager.class.php');
require_once('config.inc.php');
require_once('init.inc.php');

/* Get or Create User from Open_id */
function user_get(CRObject $info)
{
	$res['user'] = UserManager::getByOpenID($info->get('open_id'));
	if ($res['user'] === null) {
		/* if user(UID=1) not exist, set this user as admin
		* WARN: do not delete user(UID=1) */
		if (UserManager::getByUID(1) === null) {
			$info->set('role', 'admin');
		}
		if (!UserManager::add($info)) {
			$res['errno'] = Code::FAIL;
			return $res;
		}
		$res['user'] = UserManager::getByOpenID($info->get('open_id'));
	}
	$res['errno'] = $res['user'] !== null ? Code::SUCCESS : Code::UNKNOWN_ERROR;
	return $res;
}

function user_signout()
{
	Session::expire();
	$res['errno'] = Code::SUCCESS;
	return $res;
}

function log_gets(CRObject $rule)
{
	if (Session::get('uid') === null) {
		$res['errno'] = Code::NOT_LOGED;
		return $res;
	}
	if ($rule->get('who') !== 'all') {
		$rule->set('who', 'self');
		$rule->set('scope', Session::get('uid'));
		$rule->set('tag', 'user.login');
	}
	if ($rule->get('who') === 'all' && !AccessController::hasAccess(Session::get('role', 'visitor'), 'logs.get_others')) {
		$res['errno'] = Code::NO_PRIVILEGE;
		return $res;
	}
	$res['errno'] = Code::SUCCESS;
	$res['count'] = CRLogger::getCount($rule);
	$res['logs'] = CRLogger::search($rule);
	return $res;
}

function oauth_get_url()
{
	$state = Random::randomString(16);
	Session::put('oauth:state', $state);
	$url = OAUTH_SITE . '/login?response_type=code&client_id=' . OAUTH_CLIENT_ID . '&redirect_uri=' . BASE_URL . '/auth&state=' . $state . '&scope=role';
	$res['errno'] = Code::SUCCESS;
	$res['url'] = $url;
	return $res;
}

function user_login($user)
{
	/* mock */
	$info = array('open_id' => 'admin', 'role' => 'normal', 'nickname' => 'Admin');
	$open_id = ($info !== null && isset($info['open_id'])) ? $info['open_id'] : null;
	$email = ($info !== null && isset($info['email'])) ? $info['email'] : null;
	$role = ($info !== null && isset($info['role'])) ? $info['role'] : 'normal';
	$nickname = ($info !== null && isset($info['nickname'])) ? $info['nickname'] : 'u2913';

	$user = new CRObject();
	$user->set('open_id', $open_id);
	$user->set('email', $email);
	$user->set('role', $role);
	$res = user_get($user);

	if ($res['errno'] === 0) {
		$user = $res['user'];
		Session::put('uid', $user['uid']);
		Session::put('role', $user['role']);
		Session::put('nickname', $nickname);

		$log = new CRObject();
		$log->set('scope', $user['uid']);
		$log->set('tag', 'user.login');
		$content = array('uid' => $user['uid'], 'response' => $res['errno']);
		$log->set('content', json_encode($content));
		CRLogger::log($log);
	}
	return $res;
}