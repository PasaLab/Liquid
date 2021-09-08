<?php
require_once('config.inc.php');
require_once('secure.inc.php');
?>
<!DOCTYPE html>
<html lang="en-US">
<head>
	<?php require('head.php'); ?>
	<title>Yet Another Octopus | YAO</title>
</head>

<body>
<div class="wrapper">
	<?php require('header.php'); ?>
	<?php require('modals.php'); ?>
	<div class="container">

		<div id="main" class="form ui-widget load-overlay container">
			<h2 style="text-align: center">YAO---Yet Another Octopus</h2>
		</div>
	</div> <!-- /container -->
	<!--This div exists to avoid footer from covering main body-->
	<div class="push"></div>
</div>
<?php require('footer.php'); ?>
<script src="static/main.js"></script>
</body>
</html>