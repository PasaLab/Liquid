<?php

class Code
{
	/* common */
	const SUCCESS = 0;
	const FAIL = 1;
	const NO_PRIVILEGE = 2;
	const UNKNOWN_ERROR = 3;
	const IN_DEVELOP = 4;
	const INVALID_REQUEST = 5;
	const UNKNOWN_REQUEST = 6;
	const CAN_NOT_BE_EMPTY = 7;
	const INCOMPLETE_CONTENT = 8;
	const FILE_NOT_UPLOADED = 9;
	const RECORD_NOT_EXIST = 10;
	const RECORD_ALREADY_EXIST = 34;
	const INVALID_PASSWORD = 11;
	const UNABLE_TO_CONNECT_REDIS = 12;
	const UNABLE_TO_CONNECT_MYSQL = 13;

	/* user */
	const USERNAME_OCCUPIED = 14;
	const EMAIL_OCCUPIED = 15;
	const INVALID_USERNAME = 16;
	const INVALID_EMAIL = 17;
	const WRONG_PASSWORD = 18;
	const NOT_LOGED = 19;
	const USER_NOT_EXIST = 20;
	const USER_IS_BLOCKED = 21;
	const USER_IS_REMOVED = 22;
	const EMAIL_IS_NOT_VERIFIED = 23;

	const USERNAME_MISMATCH_EMAIL = 24;

	const CODE_EXPIRED = 25;
	const EMAIL_ALREADY_VERIFIED = 26;
	const INVALID_COOKIE = 27;

	/* site */
	const INVALID_DOMAIN = 28;
	const NEED_VERIFY = 29;
	const INVALID_PATTERN = 30;

	/* auth */
	const TOKEN_EXPIRED = 31;
	const SITE_NOT_EXIST = 32;
	const INVALID_URL = 33;
	const INVALID_PARAM = 34;
	const DOMAIN_MISMATCH = 35;

	const TOKEN_LENGTH_INVALID = 36;
	const URL_LENGTH_INVALID = 37;

	const RECORD_PAUSED = 38;
	const RECORD_REMOVED = 39;
	const RECORD_DISABLED = 40;
	const RECORD_NOT_IN_VALID_TIME = 41;

	/* rate limit */
	const TOO_FAST = 30;

	public static function getErrorMsg($errno)
	{
		switch ($errno) {
			case Code::SUCCESS:
				return 'Success !';

			default:
				return 'Unknown Error Code(' . $errno . ') !';
		}
	}
}
