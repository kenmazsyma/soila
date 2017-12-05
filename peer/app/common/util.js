'use_strict'

let to = src => {
	return (typeof(src)==='string') ? new Buffer(src).toString('base64') : src;
}

let fr = src => {
	return ((typeof(src)==='string') || (src.constructor === Buffer)) ? new Buffer(src, 'base64').toString() : src;
}

module.exports = {
	str : {
		b2s : src => {
			return new Buffer(src).toString('utf-8');
		},
		s2b64 : src => {
			let ret;
			if (src.constructor === Array) {
				ret = [];
				src.forEach(elm => {
					ret.push(to(elm));
				});
			} else {
				ret = to(src);
			}
			return ret;
		},
		b642s : src => {
			let ret;
			if (src.constructor === Array) {
				ret = [];
				src.forEach(elm => {
					ret.push(fr(elm));
				});
			} else {
				ret = fr(src);
			}
			return ret;
		}
	},
}
