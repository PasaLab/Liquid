$(function () {
	$('#main-tabs a').click(function (e) {
		e.preventDefault();
		$(this).tab('show');
	});

});