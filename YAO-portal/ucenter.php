<?php
require_once('predis/autoload.php');

require_once('util4p/util.php');
require_once('util4p/ReSession.class.php');
require_once('util4p/AccessController.class.php');

require_once('global.inc.php');

require_once('config.inc.php');
require_once('secure.inc.php');
require_once('init.inc.php');


if (Session::get('uid') === null) {
	header('location:/?notloged');
	exit;
}

$page_type = 'summary';
$uid = Session::get('uid');
$nickname = Session::get('nickname');

if (isset($_GET['logs'])) {
	$page_type = 'logs';

} elseif (isset($_GET['logs_all'])) {
	$page_type = 'logs_all';

} elseif (isset($_GET['summary'])) {
	$page_type = 'summary';

} elseif (isset($_GET['jobs'])) {
	$page_type = 'jobs';

} elseif (isset($_GET['job_status'])) {
	$page_type = 'job_status';

} elseif (isset($_GET['jobs_all'])) {
	$page_type = 'jobs_all';

} elseif (isset($_GET['resources'])) {
	$page_type = 'resources';

} elseif (isset($_GET['agents'])) {
	$page_type = 'agents';

} elseif (isset($_GET['clusters'])) {
	$page_type = 'clusters';

} elseif (isset($_GET['workspaces'])) {
	$page_type = 'workspaces';

}

$entries = array(
	array('summary', 'Summary'),
	array('jobs', 'Jobs'),
	array('workspaces', 'Projects'),
	array('logs', 'Activities'),
	array('resources', 'Admin:Resources'),
	//array('agents', 'Admin:Agents'),
	array('clusters', 'Admin:Queues'),
	array('logs_all', 'Admin:Audit')
);
$visible_entries = array();
foreach ($entries as $entry) {
	if (AccessController::hasAccess(Session::get('role', 'visitor'), 'ucenter.' . $entry[0])) {
		$visible_entries[] = array($entry[0], $entry[1]);
	}
}
?>
<!DOCTYPE html>
<html lang="en-US">
<head>
	<?php require('head.php'); ?>
	<title>Management | YAO</title>
	<link href="https://cdn.jsdelivr.net/npm/bootstrap-table@1.12.1/dist/bootstrap-table.min.css" rel="stylesheet">
	<script type="text/javascript">
		var page_type = "<?=$page_type?>";
	</script>
</head>

<body>
<div class="wrapper">
	<?php require('header.php'); ?>
	<?php require('modals.php'); ?>
	<div class="container">
		<div class="row">

			<div class="hidden-xs hidden-sm col-md-2 col-lg-2">
				<div class="panel panel-default">
					<div class="panel-heading">Menu Bar</div>
					<ul class="nav nav-pills nav-stacked panel-body">
						<?php foreach ($visible_entries as $entry) { ?>
							<li role="presentation" <?php if ($page_type == $entry[0]) echo 'class="disabled"'; ?> >
								<a href="?<?= $entry[0] ?>"><?= $entry[1] ?></a>
							</li>
						<?php } ?>
					</ul>
				</div>
			</div>

			<div class="col-xs-12 col-sm-12 col-md-10 col-lg-10">
				<div class="visible-xs visible-sm">
					<div class="panel panel-default">
						<div class="panel-heading">Menu Bar</div>
						<ul class="nav nav-pills panel-body">
							<?php foreach ($visible_entries as $entry) { ?>
								<li role="presentation" <?php if ($page_type == $entry[0]) echo 'class="disabled"'; ?> >
									<a href="?<?= $entry[0] ?>"><?= $entry[1] ?></a>
								</li>
							<?php } ?>
						</ul>
					</div>
				</div>

				<?php if ($page_type === 'summary') { ?>
					<div id="jobs">
						<div class="panel panel-default">
							<div class="panel-heading">Summary</div>
							<div class="panel-body">
								<div class="row">
									<div class="col-md-4">
										<canvas id="summary-chart-cpu"></canvas>
									</div>
									<div class="col-md-4">
										<canvas id="summary-chart-jobs"></canvas>
									</div>
									<div class="col-md-4">
										<canvas id="summary-chart-mem"></canvas>
									</div>
								</div>
								<div class="row">
									<div class="col-md-4">
										<canvas id="summary-chart-gpu-util"></canvas>
									</div>
									<div class="col-md-4">
										<canvas id="summary-chart-gpu"></canvas>
									</div>
									<div class="col-md-4">
										<canvas id="summary-chart-gpu-mem"></canvas>
									</div>
								</div>
							</div>
						</div>
					</div>

				<?php } elseif ($page_type === 'jobs' || $page_type === 'jobs_all') { ?>
					<div id="jobs">
						<div class="panel panel-default">
							<div class="panel-heading">Job</div>
							<div class="panel-body">
								<div class="table-responsive">
									<div id="toolbar">
										<button id="btn-job-add" class="btn btn-primary">
											<i class="glyphicon glyphicon-plus"></i> Submit
										</button>
									</div>
									<table id="table-job" data-toolbar="#toolbar" class="table table-striped">
									</table>
								</div>
							</div>
						</div>
					</div>

				<?php } elseif ($page_type === 'job_status') { ?>
					<div id="jobs">
						<div class="panel panel-default">
							<div class="panel-heading">Job Status</div>
							<div class="panel-body">
								<div class="table-responsive">
									<div id="toolbar"></div>
									<table id="table-task" data-toolbar="#toolbar" class="table table-striped">
									</table>
								</div>
							</div>
						</div>
					</div>

				<?php } elseif ($page_type === 'resources') { ?>
					<div id="resources">
						<div class="panel panel-default">
							<div class="panel-heading">Resources</div>
							<div class="panel-body">
								<div class="table-responsive">
									<div id="toolbar"></div>
									<table id="table-resource" data-toolbar="#toolbar" class="table table-striped">
									</table>
								</div>
							</div>
						</div>
					</div>

				<?php } elseif ($page_type === 'logs' || $page_type === 'logs_all') { ?>
					<div id="logs">
						<div class="panel panel-default">
							<div class="panel-heading">Activities</div>
							<div class="panel-body">
								<div class="table-responsive">
									<div id="toolbar"></div>
									<table id="table-log" data-toolbar="#toolbar" class="table table-striped">
									</table>
								</div>
							</div>
						</div>
					</div>

				<?php } elseif ($page_type === 'clusters') { ?>
					<div id="clusters">
						<div class="panel panel-default">
							<div class="panel-heading">Virtual Clusters</div>
							<div class="panel-body">
								<div class="table-responsive">
									<div id="toolbar">
										<button id="btn-cluster-add" class="btn btn-primary">
											<i class="glyphicon glyphicon-plus"></i> Create
										</button>
									</div>
									<table id="table-cluster" data-toolbar="#toolbar" class="table table-striped">
									</table>
								</div>
							</div>
						</div>
					</div>

				<?php } elseif ($page_type === 'agents') { ?>
					<div id="agents">
						<div class="panel panel-default">
							<div class="panel-heading">Agents</div>
							<div class="panel-body">
								<div class="table-responsive">
									<div id="toolbar">
										<button id="btn-agent-add" class="btn btn-primary">
											<i class="glyphicon glyphicon-plus"></i> Add
										</button>
									</div>
									<table id="table-agent" data-toolbar="#toolbar" class="table table-striped">
									</table>
								</div>
							</div>
						</div>
					</div>

				<?php } elseif ($page_type === 'workspaces') { ?>
					<div id="workspaces">
						<div class="panel panel-default">
							<div class="panel-heading">Workspaces</div>
							<div class="panel-body">
								<div class="table-responsive">
									<div id="toolbar">
										<button id="btn-workspace-add" class="btn btn-primary">
											<i class="glyphicon glyphicon-plus"></i> New
										</button>
									</div>
									<table id="table-workspace" data-toolbar="#toolbar" class="table table-striped">
									</table>
								</div>
							</div>
						</div>
					</div>

				<?php } ?>

			</div>
		</div>
	</div> <!-- /container -->

	<!--This div exists to avoid footer from covering main body-->
	<div class="push"></div>
</div>
<?php require('footer.php'); ?>
<script src="https://cdn.jsdelivr.net/npm/chart.js@2.7.3/dist/Chart.min.js"></script>
<script src="https://cdn.jsdelivr.net/npm/bootstrap-table@1.12.1/dist/bootstrap-table.min.js"></script>
<script src="https://cdn.jsdelivr.net/npm/bootstrap-table@1.12.1/dist/extensions/mobile/bootstrap-table-mobile.min.js"></script>
<script src="https://cdn.jsdelivr.net/npm/bootstrap-table@1.12.1/dist/extensions/export/bootstrap-table-export.min.js"></script>
<script src="https://cdn.jsdelivr.net/npm/tableexport.jquery.plugin@1.10.1/tableExport.min.js"></script>
<script src="https://cdn.jsdelivr.net/npm/chart.js@2.8.0"></script>
<script src="https://cdn.jsdelivr.net/npm/downloadjs@1.4.7/download.min.js"></script>

<script src="static/workspace.js"></script>
<script src="static/job.js"></script>
<script src="static/cluster.js"></script>
<script src="static/agent.js"></script>
<script src="static/resource.js"></script>
<script src="static/summary.js"></script>
<script src="static/ucenter.js"></script>
</body>
</html>