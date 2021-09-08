<?php
require_once('config.inc.php');
require_once('secure.inc.php');
?>
<!DOCTYPE html>
<html lang="en-US">
<head>
	<?php require('head.php'); ?>
	<title>Help | YAO</title>
</head>
<body>
<div class="wrapper">
	<?php require('header.php'); ?>
	<?php require('modals.php'); ?>
	<div class="container">
		<div class="row">
			<div class="col-sm-4 col-md-3 hidden-xs">
				<div id="help-nav" class="panel panel-default">
					<div class="panel-heading">List</div>
					<ul class="nav nav-pills nav-stacked panel-body">
						<li role="presentation">
							<a href="#introduction">YAO</a>
						</li>
						<li role="presentation">
							<a href="#about">About</a>
						</li>
						<li role="presentation">
							<a href="#TOS">TOS</a>
						</li>
						<li role="presentation">
							<a href="#privacy">Privacy</a>
						</li>
						<li role="presentation">
							<a href="#feedback">Feedback</a>
						</li>
					</ul>
				</div>
			</div>
			<div class="col-xs-12 col-sm-8 col-md-8 col-md-offset-1 ">
				<div id="introduction" class="panel panel-default">
					<div class="panel-heading">YAO</div>
					<div class="panel-body">
						<p>Yet Another Octopus</p>
					</div>
				</div>
				<div id="about" class="panel panel-default">
					<div class="panel-heading">About</div>
					<div class="panel-body">
						<ul>
							<li>one</li>
							<li>two</li>
						</ul>
					</div>
				</div>
				<div id="TOS" class="panel panel-default">
					<div class="panel-heading">TOS</div>
					<div class="panel-body">
						<p>Term of service</p>
					</div>
				</div>
				<div id="privacy" class="panel panel-default">
					<div class="panel-heading">Privacy</div>
					<div class="panel-body">
						<p>Privacy</p>
					</div>
				</div>
				<div id="feedback" class="panel panel-default">
					<div class="panel-heading">Feedback</div>
					<div class="panel-body">
						<p>This is feedback.</p>
					</div>
				</div>

			</div>
		</div>
	</div> <!-- /container -->
	<!--This div exists to avoid footer from covering main body-->
	<div class="push"></div>
</div>
<?php require('footer.php'); ?>
</body>
</html>
