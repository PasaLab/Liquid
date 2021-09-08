<?php

$settings = [
	'site' => [
		'version' => '0.2.2',
		'scheduler_addr' => 'http://127.0.0.1',
		'base_url' => 'http://127.0.0.1', /* make absolute url for SEO and avoid hijack, no '/' at the end */
		'timezone' => 'Asia/Shanghai',
	],
	'mysql' => [
		'host' => 'localhost',
		'port' => 3306,
		'database' => 'yao',
		'user' => 'root', /* It is not recommended to use `root` in production environment */
		'password' => '',
		'show_error' => false, /* set to true to see detailed Mysql errors __only__ for debug purpose */
	],
	'redis' => [ /* Make sure that your Redis only listens to Intranet */
		'scheme' => 'tcp',
		'host' => 'localhost',
		'port' => 6379,
		'show_error' => false, /* set to true to see detailed Redis errors __only__ for debug purpose */
	],
	'oauth' => [
		'site' => 'https://quickauth.newnius.com',
		'client_id' => 'XgaII6NxeE08LtKB',
		'client_secret' => 'L9hdi4dQToM0GsDLtcYYQ3k4ZDEjuGVOtPS3nOVKlo6cxLcVjH9TqvmTBiHAgLp2',
	]
];


foreach ($settings as $category => $values) {
	foreach ($values as $option => $value) {
		$env = getenv(strtoupper($category . '_' . $option));
		if ($env !== false) {
			$settings[$category][$option] = $env;
		};
	}
}

define('YAO_VERSION', strval($settings['site']['version']));

define('YAO_SCHEDULER_ADDR', strval($settings['site']['scheduler_addr']));

/* Mysql */
define('DB_HOST', strval($settings['mysql']['host']));
define('DB_PORT', intval($settings['mysql']['port']));
define('DB_NAME', strval($settings['mysql']['database']));
define('DB_USER', strval($settings['mysql']['user']));
define('DB_PASSWORD', strval($settings['mysql']['password']));
define('DB_SHOW_ERROR', boolval($settings['mysql']['show_error']));

/* Redis */
define('REDIS_SCHEME', strval($settings['redis']['scheme']));
define('REDIS_HOST', strval($settings['redis']['host']));
define('REDIS_PORT', intval($settings['redis']['port']));
define('REDIS_SHOW_ERROR', boolval($settings['redis']['show_error']));

/* Site */
define('BASE_URL', strval($settings['site']['base_url']));
define('WEB_ROOT', __DIR__);
define('FEEDBACK_EMAIL', 'mail@example.com');

/* Auth */
define('AUTH_CODE_TIMEOUT', 300); // 5 min
define('AUTH_TOKEN_TIMEOUT', 604800); // 7 day

/* Session */
define('ENABLE_MULTIPLE_LOGIN', true);
define('BIND_SESSION_WITH_IP', false);  // current session will be logged when ip changes
define('SESSION_TIME_OUT', 1800);// 30 minutes 30*60=1800
define('ENABLE_COOKIE', true);

/* Rate Limit */
define('ENABLE_RATE_LIMIT', false);
define('RATE_LIMIT_PREFIX', 'rl');

/* OAuth */
/* The default conf is only usable when this runs on localhost */
define('OAUTH_SITE', strval($settings['oauth']['site']));
define('OAUTH_CLIENT_ID', strval($settings['oauth']['client_id']));
define('OAUTH_CLIENT_SECRET', strval($settings['oauth']['client_secret']));

header("content-type:text/html; charset=utf-8");

date_default_timezone_set(strval($settings['site']['timezone']));