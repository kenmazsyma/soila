Rpc = {
	SV_URL : '',
	getId : function() {
		return Math.floor( Math.random() * 0x8000000000000000 );
	},
	call : function(method, param, cb, cberr, failcb, bNoModal) {
		var name = this.call.caller.name;
		var getMsg = function(val) {
			var sep = '&nbsp;:&nbsp;';
			if (isEmpty(val)) return val;
			return val + '&nbsp;(' + name + sep + method + ')';
		}
		var data = {
			'jsonrpc':'2.0',
			'id':Rpc.getId(),
			'method':method,
			'params':param
		};
		if (failcb===failcb||type(failcb)!=='function') {
			failcb = function(jqXHR, textStatus, errorThrown) {};
		}
		var retry = 0;
		var call = function() {
			$.ajax({
				type: "POST",
				url: 'http://' + location.host + '/api',
				data: JSON.stringify(data),
				dataType: 'json',
				contentType: "application/json",
				success: function(ret) {
					if (ret.error) {
						ret.statusText = ret.error;
					} else {
						cb(ret);
					}
				},
				error : function(ret) {
					var msg = getMsg(cberr(ret));
					retry++;
					if (ret.statusText=='Not Found') {
						if (retry==2) {
							msg = 'エラーが発生しました。';
						} else {
							window.setTimeout(call, 500);
							return;
						}
					} else if (ret.statusText=='error') {
						if (retry==3) {
							msg = 'エラーが発生しました。';
						} else {
							window.setTimeout(call, 500);
							return;
						}
					}
					alert(msg);
				},
			}).fail(function(ret) {Modal.close(); Modal.msg(getMsg(failcb(ret)))});
		}
		call();
	}
};

var Base64 = {
	encode: function(str) {
		return btoa(unescape(encodeURIComponent(str)));
	},
	decode: function(str) {
		return decodeURIComponent(escape(atob(str)));
	}
};

