'use_strict';

let log = require('../common/logger')('api.list');

module.exports = {
	cls : {
		person : require('./person')
	},
	call : function(cn, mn, json) {
		return new Promise((resolve, reject) => {
			log.debug('api.call:' + cn + ':' + mn);
			let cls = this.cls[cn];
			if (cls===undefined||cls[mn]===undefined) {
				reject('api not found');
				return;
			}
			try {
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
