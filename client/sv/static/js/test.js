$(function() {
	$('#run').click(function() {
		Rpc.call('Person.Create',
				 [{
					 id : $('#user').val(),
					 name : $('#name').val(),
					 pass : $('#pass').val(),
					 profile : $('#desc').val()
				 }],
			function(ret) {
				alert(JSON.stringify(ret.result));
			},
			function(ret) {
				return ret.statusText;
			},
			function(ret) {
				return ret.statusText;
			}
		);
	});
});
