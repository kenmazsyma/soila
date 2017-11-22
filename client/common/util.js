'use_strict'

module.exports = {
	str : {
		b2s : function(src) {
			return new Buffer(src).toString('utf-8');
		}
	}
}
