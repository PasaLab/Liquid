function register_events_resource() {

}

function load_resources() {
	$("#table-resource").bootstrapTable({
		url: 'service?action=resource_list',
		responseHandler: resourceResponseHandler,
		sidePagination: 'server',
		cache: true,
		striped: true,
		pagination: false,
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
			field: 'host',
			title: 'Node',
			align: 'center',
			valign: 'middle',
			escape: true
		}, {
			field: 'CPU',
			title: 'CPU',
			align: 'center',
			valign: 'middle'
		}, {
			field: 'MEM',
			title: 'MEM',
			align: 'center',
			valign: 'middle',
			escape: true,
			visible: true
		}, {
			field: 'GPU',
			title: 'GPU',
			align: 'center',
			valign: 'middle',
			escape: true,
			visible: true
		}, {
			field: 'operate',
			title: 'Operate',
			align: 'center',
			events: resourceOperateEvents,
			formatter: resourceOperateFormatter
		}]
	});
}

function resourceResponseHandler(res) {
	if (res['errno'] === 0) {
		var tmp = {};
		tmp["rows"] = [];
		$.each(res["resources"], function (i, node) {
			var item = {
				'id': node.id,
				'host': node.host,
				'CPU': node.cpu_num,
				'MEM': node.mem_available + ' / ' + node.mem_total + ' (GB)',
				'GPU': node.status.length,
				'status': node.status
			};
			tmp["rows"].push(item);
		});

		return tmp['rows'];
	}
	$("#modal-msg-content").html(res["msg"]);
	$("#modal-msg").modal('show');
	return [];
}

function resourceOperateFormatter(value, row, index) {
	var div = '<div class="btn-group" role="group" aria-label="...">';
	div += '<button class="btn btn-default view"><i class="glyphicon glyphicon-eye-open"></i>&nbsp;</button>';
	div += '</div>';
	return div;
}


window.resourceOperateEvents = {
	'click .view': function (e, value, row, index) {
		var tmp = jQuery.extend(true, {}, row.status);
		var formattedData = JSON.stringify(tmp, null, '\t');
		$('#modal-resource-gpu-detail-content').text(formattedData);
		$('#modal-resource-gpu-detail').modal('show');
	}
};