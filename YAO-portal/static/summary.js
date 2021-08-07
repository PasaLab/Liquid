function register_events_summary() {

}

function summary_render() {
	var ctx_cpu = document.getElementById('summary-chart-cpu').getContext('2d');
	var ctx_mem = document.getElementById('summary-chart-mem').getContext('2d');
	var ctx_jobs = document.getElementById('summary-chart-jobs').getContext('2d');
	var ctx_gpu = document.getElementById('summary-chart-gpu').getContext('2d');
	var ctx_gpu_util = document.getElementById('summary-chart-gpu-util').getContext('2d');
	var ctx_gpu_mem = document.getElementById('summary-chart-gpu-mem').getContext('2d');

	var ajax = $.ajax({
		url: "service?action=summary_get",
		type: 'GET',
		data: {}
	});
	ajax.done(function (res) {
		if (res["errno"] !== 0) {
			$("#modal-msg-content").html(res["msg"]);
			$("#modal-msg").modal('show');
		}


		/* Jobs */
		var data = {
			datasets: [{
				data: Object.values(res['jobs']),
				backgroundColor: ["rgb(54, 162, 235)", "rgb(255, 99, 132)", "rgb(255, 205, 86)"]
			}],
			labels: Object.keys(res['jobs'])
		};
		var myPieChart = new Chart(ctx_jobs, {
			type: 'pie',
			data: data,
			options: {
				title: {
					display: true,
					text: 'Jobs'
				},
				legend: {
					display: false
				}
			}
		});

		/* GPUs */
		var data2 = {
			datasets: [{
				data: Object.values(res['gpu']),
				backgroundColor: ["rgb(54, 162, 235)", "rgb(255, 99, 132)"]
			}],
			labels: Object.keys(res['gpu'])
		};
		var myPieChart2 = new Chart(ctx_gpu, {
			type: 'pie',
			data: data2,
			options: {
				title: {
					display: true,
					text: 'GPUs'
				},
				legend: {
					display: false
				}
			}
		});

	});
	ajax.fail(function (jqXHR, textStatus) {
		$("#modal-msg-content").html("Request failed : " + jqXHR.statusText);
		$("#modal-msg").modal('show');
	});


	var ajax_pool = $.ajax({
		url: "service?action=summary_get_pool_history",
		type: 'GET',
		data: {}
	});
	ajax_pool.done(function (res) {
		if (res["errno"] !== 0) {
			$("#modal-msg-content").html(res["msg"]);
			$("#modal-msg").modal('show');
		}

		var cpu_util = [];
		var cpu_total = [];
		var mem_available = [];
		var mem_total = [];
		var mem_using = [];
		var gpu_util = [];
		var gpu_total = [];
		var gpu_mem_available = [];
		var gpu_mem_total = [];
		var gpu_mem_using = [];
		var timestamps = [];
		$.each(res["data"], function (i, item) {
			cpu_util.push(item['cpu_util'].toFixed(2));
			cpu_total.push(item['cpu_total']);
			mem_available.push(item['mem_available']);
			mem_total.push(item['mem_total']);
			mem_using.push(item['mem_total'] - item['mem_available']);
			gpu_util.push(item['gpu_util']);
			gpu_total.push(item['gpu_total']);
			gpu_mem_available.push(item['gpu_mem_available']);
			gpu_mem_total.push(item['gpu_mem_total']);
			gpu_mem_using.push(item['gpu_mem_total'] - item['gpu_mem_available']);
			timestamps.push(moment(item['ts']).format('HH:mm:ss'));
		});

		/* CPU Load */
		ctx_cpu.canvas.height = 200;
		new Chart(ctx_cpu, {
			"type": "line",
			"data": {
				"labels": timestamps,
				"datasets": [{
					"label": "CPU Load",
					"data": cpu_util,
					"fill": true,
					"borderColor": "rgb(75, 192, 192)",
					"lineTension": 0.1
				}]
			},
			"options": {
				title: {
					display: true,
					text: 'CPU Load'
				},
				legend: {
					display: false
				},
				maintainAspectRatio: false
			}
		});


		/* Mem Using */
		ctx_mem.canvas.height = 200;
		new Chart(ctx_mem, {
			"type": "line",
			"data": {
				"labels": timestamps,
				"datasets": [{
					"label": "Using",
					"data": mem_using,
					"fill": true,
					"borderColor": "rgb(75, 192, 192)",
					"lineTension": 0.1
				}, {
					"label": "Total",
					"data": mem_total,
					"fill": true,
					"borderColor": "rgb(75, 192, 192)",
					"lineTension": 0.1
				}]
			},
			"options": {
				title: {
					display: true,
					text: 'MEM Using'
				},
				legend: {
					display: false
				},
				maintainAspectRatio: false
			}
		});

		/* GPU Util */
		ctx_gpu_util.canvas.height = 200;
		new Chart(ctx_gpu_util, {
			"type": "line",
			"data": {
				"labels": timestamps,
				"datasets": [{
					"label": "GPU Util",
					"data": gpu_util,
					"fill": true,
					"borderColor": "rgb(75, 192, 192)",
					"lineTension": 0.1
				}]
			},
			"options": {
				title: {
					display: true,
					text: 'GPU Utilization'
				},
				legend: {
					display: false
				},
				maintainAspectRatio: false
			}
		});


		/* GPU Mem Using */
		ctx_gpu_mem.canvas.height = 200;
		new Chart(ctx_gpu_mem, {
			"type": "line",
			"data": {
				"labels": timestamps,
				"datasets": [{
					"label": "Using",
					"data": gpu_mem_using,
					"fill": true,
					"borderColor": "rgb(75, 192, 192)",
					"lineTension": 0.1
				}, {
					"label": "Total",
					"data": gpu_mem_total,
					"fill": true,
					"borderColor": "rgb(75, 192, 192)",
					"lineTension": 0.1
				}]
			},
			"options": {
				title: {
					display: true,
					text: 'GPU MEM Using'
				},
				legend: {
					display: false
				},
				maintainAspectRatio: false
			}
		});
	});
	ajax_pool.fail(function (jqXHR, textStatus) {
		$("#modal-msg-content").html("Request failed : " + jqXHR.statusText);
		$("#modal-msg").modal('show');
	});
}