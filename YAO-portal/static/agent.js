function register_events_agent() {
	$('#btn-agent-add').click(function (e) {
		$('#form-agent-submit-type').val('add');
		$('#modal-agent').modal('show');
	});

	$("#form-agent-submit").click(function (e) {
		var ip = $('#form-agent-ip').val();
		var alias = $('#form-agent-alias').val();
		var cluster = $('#form-agent-cluster').val();

		/* TODO validate form */

		$('#modal-agent').modal('hide');
		if ($('#form-agent-submit-type').val() !== 'add')
			return;

		var ajax = $.ajax({
			url: "service?action=agent_add",
			type: 'POST',
			data: {
				ip: ip,
				alias: alias,
				cluster: cluster
			}
		});
		ajax.done(function (res) {
			if (res["errno"] !== 0) {
				$("#modal-msg-content").html(res["msg"]);
				$("#modal-msg").modal('show');
			}
			$('#table-agent').bootstrapTable("refresh");
		});
		ajax.fail(function (jqXHR, textStatus) {
			$("#modal-msg-content").html("Request failed : " + jqXHR.statusText);
			$("#modal-msg").modal('show');
			$('#table-agent').bootstrapTable("refresh");
		});
	});

}

function load_agents(cluster) {
	$("#table-agent").bootstrapTable({
		url: 'service?action=agent_list&who=' + cluster,
		responseHandler: agentResponseHandler,
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
			field: 'alias',
			title: 'Alias',
			align: 'center',
			valign: 'middle',
			escape: true
		}, {
			field: 'ip',
			title: 'IP',
			align: 'center',
			valign: 'middle',
			formatter: long2ip
		}, {
			field: 'cluster',
			title: 'Cluster',
			align: 'center',
			valign: 'middle'
		}, {
			field: 'operate',
			title: 'Operate',
			align: 'center',
			events: agentOperateEvents,
			formatter: agentOperateFormatter
		}]
	});
}

function agentResponseHandler(res) {
	if (res['errno'] === 0) {
		var tmp = {};
		tmp["total"] = res["count"];
		tmp["rows"] = res["agents"];
		return tmp;
	}
	$("#modal-msg-content").html(res["msg"]);
	$("#modal-msg").modal('show');
	return [];
}

function agentOperateFormatter(value, row, index) {
	var div = '<div class="btn-group" role="group" aria-label="...">';
	div += '<button class="btn btn-default view"><i class="glyphicon glyphicon-eye-open"></i>&nbsp;</button>';
	div += '<button class="btn btn-default remove"><i class="glyphicon glyphicon-remove"></i>&nbsp;</button>';
	div += '</div>';
	return div;
}

window.agentOperateEvents = {
	'click .view': function (e, value, row, index) {
		$('#form-agent-id').val(row.id);
		$('#form-agent-ip').val(long2ip(row.ip));
		$('#form-agent-alias').val(row.alias);
		$('#form-agent-token').val(row.token);
		$('#form-agent-submit-type').val('view');
		$('#modal-agent').modal('show');
	},
	'click .remove': function (e, value, row, index) {
		if (!confirm('Are you sure to remove this agent?')) {
			return;
		}
		var ajax = $.ajax({
			url: "service?action=agent_remove",
			type: 'POST',
			data: {id: row.id}
		});
		ajax.done(function (res) {
			if (res["errno"] !== 0) {
				$("#modal-msg-content").html(res["msg"]);
				$("#modal-msg").modal('show');
			}
			$('#table-agent').bootstrapTable("refresh");
		});
		ajax.fail(function (jqXHR, textStatus) {
			$("#modal-msg-content").html("Request failed : " + jqXHR.statusText);
			$("#modal-msg").modal('show');
			$('#table-agent').bootstrapTable("refresh");
		});
	}
};