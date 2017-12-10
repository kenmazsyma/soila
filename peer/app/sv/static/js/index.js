_ = {
	regPeer : function() {
		Rpc.call('peer.register', [], function(res) {
			alert(res);
		}, function(err) {
		}, function(fail) {
		});
	}
};

$(function(){
	$('#regpeer').click(_.regPeer());
});

