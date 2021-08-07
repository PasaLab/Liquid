$(function () {
	$("#btn-signout").click(function (e) {
		e.preventDefault();
		var ajax = $.ajax({
			url: "service?action=user_signout",
			type: 'POST',
			data: {}
		});
		ajax.done(function (res) {
			window.location.pathname = "/";
		});
	});

	$("#btn-oauth-login").click(function (e) {
		e.preventDefault();
		var ajax = $.ajax({
			url: "service?action=user_login",
			type: 'POST',
			data: {}
		});
		ajax.done(function (res) {
			if (res["errno"] !== 0) {
				$("#modal-msg-content").html(res["msg"]);
				$("#modal-msg").modal('show');
			} else {
				window.location.pathname = "/ucenter";
			}
		});
	});

	$('.date-picker').datetimepicker();
});