function register_events_job() {
	$('#btn-job-add').click(function (e) {
		var cb = function (workspaces) {
			$('#form-job-workspace').children().remove();
			$.each(workspaces, function (i, workspace) {
				var newGroupOption = '<option value="' + workspace.git_repo + '">' + workspace.name + '</option>';
				$('#form-job-workspace').append(newGroupOption);
			});
		};
		wordspace_gets(null, cb);

		var cb_cluster = function (clusters) {
			$('#form-job-cluster').children().remove();
			$.each(clusters, function (i, cluster) {
				var newGroupOption = '<option value="' + cluster.name + '">' + cluster.name + '</option>';
				$('#form-job-cluster').append(newGroupOption);
			});
		};
		cluster_gets(cb_cluster);

		//$('#form-job-name').val('');
		$('#form-job-priority').val(25);
		//$('#form-job-cluster').val(1);
		$('#modal-job').modal('show');
	});

	$('#form-job-tasks').on('click', '.task-remove', function (e) {
		if ($('#form-job-tasks').find('.row').length <= 1) {
			return;
		}
		var task = $(this).parent().parent().parent();
		task.remove();
	});

	$("#form-job-predict-req").click(function (e) {
		var name = $('#form-job-name').val();
		var workspace = $('#form-job-workspace').val();
		var cluster = $('#form-job-cluster').val();
		var priority = $('#form-job-priority').val();
		var run_before = $('#form-job-run-before').val();
		var locality = $('#form-job-locality').val();
		if (run_before.length !== 0) {
			run_before = moment(run_before).unix();
		}
		var tasks = [];
		$('#form-job-tasks').find('.row').each(function () {
			var task = {};
			task['name'] = $(this).find('.task-name').eq(0).val();
			task['image'] = $(this).find('.task-image').eq(0).val();
			task['cmd'] = $(this).find('.task-cmd').eq(0).val();
			task['cpu_number'] = $(this).find('.task-cpu').eq(0).val();
			task['memory'] = $(this).find('.task-mem').eq(0).val();
			task['gpu_number'] = $(this).find('.task-gpu-num').eq(0).val();
			task['gpu_memory'] = $(this).find('.task-gpu-mem').eq(0).val();
			task['is_ps'] = $(this).find('.task-is-ps').eq(0).val();
			task['gpu_model'] = $(this).find('.task-gpu-model').eq(0).val();
			tasks.push(task);
		});

		/* TODO validate form */
		if (name.length === 0) {
			return true;
		}
		$.each(tasks, function (i, task) {
			if (task['name'].length === 0) {
				return true;
			}
		});

		var roles = ['PS', 'Worker'];
		$.each(roles, function (i, role) {
			var ajax = $.ajax({
				url: "service?action=job_predict_req&role=" + role,
				type: 'POST',
				data: {
					name: name,
					workspace: workspace,
					cluster: cluster,
					priority: priority,
					run_before: run_before,
					locality: locality,
					tasks: JSON.stringify(tasks)
				}
			});
			ajax.done(function (res) {
				if (res["errno"] !== 0) {
					$("#modal-msg-content").html(res["msg"]);
					$("#modal-msg").modal('show');
				} else {
					$('#form-job-tasks').find('.row').each(function () {
						var taskRole = parseInt($(this).find('.task-is-ps').eq(0).val());
						console.log(taskRole);
						if ((role === 'PS' && taskRole === 1) || (role === 'Worker' && taskRole === 0)) {
							$(this).find('.task-cpu').eq(0).val(res['cpu']);
							$(this).find('.task-mem').eq(0).val(res['mem']);
							$(this).find('.task-gpu-mem').eq(0).val(res['gpu_mem']);
						}
					});
				}
			});
			ajax.fail(function (jqXHR, textStatus) {
				$("#modal-msg-content").html("Request failed : " + jqXHR.statusText);
				$("#modal-msg").modal('show');
				$('#table-job').bootstrapTable("refresh");
			});

		});
	});

	$("#form-job-predict-time").click(function (e) {
		var name = $('#form-job-name').val();
		var workspace = $('#form-job-workspace').val();
		var model_dir = $('#form-job-model_dir').val();
		var output_dir = $('#form-job-output_dir').val();
		var cluster = $('#form-job-cluster').val();
		var priority = $('#form-job-priority').val();
		var run_before = $('#form-job-run-before').val();
		var locality = $('#form-job-locality').val();
		if (run_before.length !== 0) {
			run_before = moment(run_before).unix();
		}
		var tasks = [];
		$('#form-job-tasks').find('.row').each(function () {
			var task = {};
			task['name'] = $(this).find('.task-name').eq(0).val();
			task['image'] = $(this).find('.task-image').eq(0).val();
			task['cmd'] = $(this).find('.task-cmd').eq(0).val();
			task['cpu_number'] = $(this).find('.task-cpu').eq(0).val();
			task['memory'] = $(this).find('.task-mem').eq(0).val();
			task['gpu_number'] = $(this).find('.task-gpu-num').eq(0).val();
			task['gpu_memory'] = $(this).find('.task-gpu-mem').eq(0).val();
			task['is_ps'] = $(this).find('.task-is-ps').eq(0).val();
			task['gpu_model'] = $(this).find('.task-gpu-model').eq(0).val();
			tasks.push(task);
		});

		/* TODO validate form */
		if (name.length === 0) {
			return true;
		}
		$.each(tasks, function (i, task) {
			if (task['name'].length === 0) {
				return true;
			}
		});

		var ajax = $.ajax({
			url: "service?action=job_predict_time",
			type: 'POST',
			data: {
				name: name,
				workspace: workspace,
				model_dir: model_dir,
				output_dir: output_dir,
				cluster: cluster,
				priority: priority,
				run_before: run_before,
				locality: locality,
				tasks: JSON.stringify(tasks)
			}
		});
		ajax.done(function (res) {
			if (res["errno"] !== 0) {
				$("#modal-msg-content").html(res["msg"]);
				$("#modal-msg").modal('show');
			} else {
				console.log(res);
			}
		});
		ajax.fail(function (jqXHR, textStatus) {
			$("#modal-msg-content").html("Request failed : " + jqXHR.statusText);
			$("#modal-msg").modal('show');
		});
	});

	$('#form-job-task-add').click(function (e) {
		var tasks = $('#form-job-tasks');
		var newTask = $('#form-job-tasks').find('.row').eq(0).clone();
		tasks.append(newTask);
	});

	$("#form-job-submit").click(function (e) {
		var name = $('#form-job-name').val();
		var workspace = $('#form-job-workspace').val();
		var cluster = $('#form-job-cluster').val();
		var priority = $('#form-job-priority').val();
		var run_before = $('#form-job-run-before').val();
		var locality = $('#form-job-locality').val();
		if (run_before.length !== 0) {
			run_before = moment(run_before).unix();
		}
		var tasks = [];
		$('#form-job-tasks').find('.row').each(function () {
			var task = {};
			task['name'] = $(this).find('.task-name').eq(0).val();
			task['image'] = $(this).find('.task-image').eq(0).val();
			task['cmd'] = $(this).find('.task-cmd').eq(0).val();
			task['cpu_number'] = $(this).find('.task-cpu').eq(0).val();
			task['memory'] = $(this).find('.task-mem').eq(0).val();
			task['gpu_number'] = $(this).find('.task-gpu-num').eq(0).val();
			task['gpu_memory'] = $(this).find('.task-gpu-mem').eq(0).val();
			task['is_ps'] = $(this).find('.task-is-ps').eq(0).val();
			task['gpu_model'] = $(this).find('.task-gpu-model').eq(0).val();
			tasks.push(task);
		});

		/* TODO validate form */
		if (name.length === 0) {
			return true;
		}
		$.each(tasks, function (i, task) {
			if (task['name'].length === 0) {
				return true;
			}
		});


		$('#modal-job').modal('hide');
		var ajax = $.ajax({
			url: "service?action=job_submit",
			type: 'POST',
			data: {
				name: name,
				workspace: workspace,
				cluster: cluster,
				priority: priority,
				run_before: run_before,
				locality: locality,
				tasks: JSON.stringify(tasks)
			}
		});
		ajax.done(function (res) {
			if (res["errno"] !== 0) {
				$("#modal-msg-content").html(res["msg"]);
				$("#modal-msg").modal('show');
			}
			$('#table-job').bootstrapTable("refresh");

		});
		ajax.fail(function (jqXHR, textStatus) {
			$("#modal-msg-content").html("Request failed : " + jqXHR.statusText);
			$("#modal-msg").modal('show');
			$('#table-job').bootstrapTable("refresh");
		});
	});

}

function load_jobs(scope) {
	$("#table-job").bootstrapTable({
		url: 'service?action=job_list&who=' + scope,
		responseHandler: jobResponseHandler,
		sidePagination: 'client',
		cache: true,
		striped: true,
		pagination: true,
		pageSize: 10,
		pageList: [10, 25, 50, 100, 200],
		search: false,
		showColumns: true,
		showRefresh: true,
		showToggle: false,
		showPaginationSwitch: true,
		minimumCountColumns: 2,
		clickToSelect: false,
		sortName: 'nobody',
		sortOrder: 'desc',
		smartDisplay: true,
		mobileResponsive: true,
		showExport: true,
		columns: [{
			field: 'created_by',
			title: 'Created By',
			align: 'center',
			valign: 'middle',
			formatter: UIDFormatter,
			visible: scope === 'all'
		}, {
			field: 'name',
			title: 'Name',
			align: 'center',
			valign: 'middle',
			escape: true
		}, {
			field: 'workspace',
			title: 'Workspace',
			align: 'center',
			valign: 'middle',
			visible: false,
			formatter: workspaceFormatter
		}, {
			field: 'group',
			title: 'Group',
			align: 'center',
			valign: 'middle',
			formatter: clusterFormatter
		}, {
			field: 'priority',
			title: 'Priority',
			align: 'center',
			valign: 'middle',
			formatter: priorityFormatter
		}, {
			field: 'run_before',
			title: 'Run Before',
			align: 'center',
			valign: 'middle',
			visible: false,
			formatter: timeFormatter
		}, {
			field: 'created_at',
			title: 'Created At',
			align: 'center',
			valign: 'middle',
			formatter: timeFormatter
		}, {
			field: 'started_at',
			title: 'Started At',
			align: 'center',
			valign: 'middle',
			formatter: timeFormatter,
			visible: false
		}, {
			field: 'updated_at',
			title: 'Updated At',
			align: 'center',
			valign: 'middle',
			formatter: timeFormatter
		}, {
			field: 'status',
			title: 'Status',
			align: 'center',
			valign: 'middle',
			formatter: statusFormatter,
			visible: true
		}, {
			field: 'base_priority',
			title: 'BasePriority',
			align: 'center',
			valign: 'middle',
			visible: false
		}, {
			field: 'operate',
			title: 'Operate',
			align: 'center',
			events: jobOperateEvents,
			formatter: jobOperateFormatter
		}]
	});
}

var UIDFormatter = function (UID) {
	return UID;
};

var workspaceFormatter = function (workspace) {
	return workspace;
};

var clusterFormatter = function (cluster) {
	return cluster;
};

var priorityFormatter = function (status) {
	status = parseInt(status);
	switch (status) {
		case 1:
			return '<span class="text-normal">Low</span>';
		case 25:
			return '<span class="text-info">Medium</span>';
		case 50:
			return '<span class="text-success">High</span>';
		case 99:
			return '<span class="text-danger">Urgent</span>';
	}
	return 'Unknown (' + status + ')';
};

var statusFormatter = function (status) {
	status = parseInt(status);
	switch (status) {
		case 0:
			return '<span class="text-normal">Submitted</span>';
		case 1:
			return '<span class="text-info">Starting</span>';
		case 2:
			return '<span class="text-primary">Running</span>';
		case 3:
			return '<span class="text-danger">Stopped</span>';
		case 4:
			return '<span class="text-success">Finished</span>';
		case 5:
			return '<span class="text-warning">Failed</span>';
	}
	return 'Unknown(' + status + ')';
};

function jobResponseHandler(res) {
	if (res['errno'] === 0) {
		return res["jobs"];
	}
	$("#modal-msg-content").html(res["msg"]);
	$("#modal-msg").modal('show');
	return [];
}

function jobOperateFormatter(value, row, index) {
	var div = '<div class="btn-group" role="group" aria-label="...">';
	if (page_type === 'jobs')
		div += '<button class="btn btn-default config"><i class="glyphicon glyphicon-cog"></i>&nbsp;</button>';
	if (page_type === 'jobs')
		div += '<button class="btn btn-default stats"><i class="glyphicon glyphicon-eye-open"></i>&nbsp;</button>';
	if (page_type === 'jobs' && (parseInt(row.status) !== 3 && parseInt(row.status) !== 4))
		div += '<button class="btn btn-default stop"><i class="glyphicon glyphicon-remove"></i>&nbsp;</button>';
	div += '</div>';
	return div;
}

window.jobOperateEvents = {
	'click .config': function (e, value, row, index) {
		var tmp = jQuery.extend(true, {}, row);
		tmp.tasks = JSON.parse(tmp.tasks);
		var formattedData = JSON.stringify(tmp, null, '\t');
		$('#modal-job-description-content').text(formattedData);
		$('#modal-job-description').modal('show');
	},
	'click .stats': function (e, value, row, index) {
		window.open("?job_status&name=" + row.name);
	},
	'click .stop': function (e, value, row, index) {
		if (!confirm('Are you sure to stop this job?')) {
			return;
		}
		var ajax = $.ajax({
			url: "service?action=job_stop",
			type: 'POST',
			data: {id: row.name}
		});
		ajax.done(function (res) {
			if (res["errno"] !== 0) {
				$("#modal-msg-content").html(res["msg"]);
				$("#modal-msg").modal('show');
			}
			$('#table-job').bootstrapTable("refresh");
		});
		ajax.fail(function (jqXHR, textStatus) {
			$("#modal-msg-content").html("Request failed : " + jqXHR.statusText);
			$("#modal-msg").modal('show');
			$('#table-job').bootstrapTable("refresh");
		});
	}
};

function load_job_status(name) {
	$("#table-task").bootstrapTable({
		url: 'service?action=job_status&name=' + name,
		responseHandler: jobStatusResponseHandler,
		sidePagination: 'server',
		cache: true,
		striped: true,
		pagination: true,
		pageSize: 10,
		pageList: [10, 25, 50, 100, 200],
		search: false,
		showColumns: true,
		showRefresh: true,
		showToggle: false,
		showPaginationSwitch: true,
		minimumCountColumns: 2,
		clickToSelect: false,
		sortName: 'nobody',
		sortOrder: 'desc',
		smartDisplay: true,
		mobileResponsive: true,
		showExport: true,
		columns: [{
			field: 'id',
			title: 'ID',
			align: 'center',
			valign: 'middle',
			visible: false
		}, {
			field: 'image',
			title: 'Image',
			align: 'center',
			valign: 'middle',
			visible: false
		}, {
			field: 'image_digest',
			title: 'Image Version',
			align: 'center',
			valign: 'middle',
			visible: false
		}, {
			field: 'hostname',
			title: 'Hostname',
			align: 'center',
			valign: 'middle'
		}, {
			field: 'command',
			title: 'Command',
			align: 'center',
			valign: 'middle'
		}, {
			field: 'created_at',
			title: 'Created At',
			align: 'center',
			valign: 'middle'
		}, {
			field: 'finished_at',
			title: 'Finished At',
			align: 'center',
			valign: 'middle',
			visible: false
		}, {
			field: 'node',
			title: 'Node',
			align: 'center',
			valign: 'middle'
		}, {
			field: 'status',
			title: 'Status',
			align: 'center',
			valign: 'middle'
		}, {
			field: 'operate',
			title: 'Logs',
			align: 'center',
			events: jobStatusOperateEvents,
			formatter: jobStatusOperateFormatter
		}]
	});
}

function jobStatusResponseHandler(res) {
	if (res['errno'] === 0) {
		var tmp = {};
		tmp["total"] = res["count"];
		tmp["rows"] = res["tasks"];
		return tmp;
	}
	$("#modal-msg-content").html(res["msg"]);
	$("#modal-msg").modal('show');
	return [];
}

function jobStatusOperateFormatter(value, row, index) {
	var div = '<div class="btn-group" role="group" aria-label="...">';
	div += '<button class="btn btn-default logs"><i class="glyphicon glyphicon-eye-open"></i>&nbsp;</button>';
	div += '<button class="btn btn-default download"><i class="glyphicon glyphicon-download-alt"></i>&nbsp;</button>';
	div += '</div>';
	return div;
}

window.jobStatusOperateEvents = {
	'click .logs': function (e, value, row, index) {
		var job = getParameterByName('name');
		var task = row.hostname;

		var ajax = $.ajax({
			url: "service?action=task_logs",
			type: 'GET',
			data: {
				job: job,
				task: task
			}
		});
		ajax.done(function (res) {
			if (res["errno"] !== 0) {
				$("#modal-msg-content").html(res["msg"]);
				$("#modal-msg").modal('show');
			}
			$('#modal-task-logs-content').text(res['logs']);
			$('#modal-task-logs').modal('show');
		});
		ajax.fail(function (jqXHR, textStatus) {
			$("#modal-msg-content").html("Request failed : " + jqXHR.statusText);
			$("#modal-msg").modal('show');
			$('#table-job').bootstrapTable("refresh");
		});
	},
	'click .download': function (e, value, row, index) {
		var job = getParameterByName('name');
		var task = row.hostname;

		var ajax = $.ajax({
			url: "service?action=task_logs",
			type: 'GET',
			data: {
				job: job,
				task: task
			}
		});
		ajax.done(function (res) {
			if (res["errno"] !== 0) {
				$("#modal-msg-content").html(res["msg"]);
				$("#modal-msg").modal('show');
			} else {
				download(res['logs'], job + '_' + task + '.txt', "text/plain");
			}
		});
		ajax.fail(function (jqXHR, textStatus) {
			$("#modal-msg-content").html("Request failed : " + jqXHR.statusText);
			$("#modal-msg").modal('show');
			$('#table-job').bootstrapTable("refresh");
		});
	}
};