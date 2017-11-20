'use_strict';

let logger = require('../common/logger');

module.exports = {
	cls : {
		user : require('./user')
	},
	call : function(cn, mn, json) {
		return new Promise((resolve, reject) => {
			logger.debug.trace('api.call:' + cn + ':' + mn);
			let cls = this.cls[cn];
			if (cls===undefined||cls[mn]===undefined) {
				reject('api not found');
			}
			try {
				console.log(json);
				let prm = JSON.parse(json);
				cls[mn](prm).then(rslt => {
					resolve(JSON.stringify(rslt));
				}).catch(e => {
					reject(e);
				});
			} catch (e) {
				reject(e);
			}
		});
	}
};
