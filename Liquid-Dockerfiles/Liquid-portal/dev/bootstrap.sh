#!/bin/bash

if [[ ! -f /var/www/html/config.inc.php ]]; then
	cp /var/www/html/config-sample.inc.php /var/www/html/config.inc.php
fi

if [[ -d /config/ ]]; then
	if [[ ! -f /config/config.inc.php ]]; then
		cp /var/www/html/config.inc.php /config/config.inc.php
	fi
fi

if [[ -f /config/config.inc.php ]]; then
	rm /var/www/html/config.inc.php
	ln -s /config/config.inc.php /var/www/html/config.inc.php
fi

apache2-foreground
