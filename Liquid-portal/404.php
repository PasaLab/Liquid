<?php
require_once('Code.class.php');
require_once('secure.inc.php');

$error = '404 Not Found';
?>
<!DOCTYPE html>
<html lang="en-US">
<head>
	<?php require_once('head.php'); ?>
	<title>404 | YAO</title>
</head>

<body>
<div class="wrapper">
	<?php require_once('header.php'); ?>
	<div class="container">
		<div class="container">
			<h2 style="text-align: center"><?= $error ?></h2>
		</div>
	</div> <!-- /container -->
	<!--This div exists to avoid footer from covering main body-->
	<div class="push"></div>
</div>
<?php require_once('footer.php'); ?>
</body>
</html>
