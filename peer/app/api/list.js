'use_strict';

let log = require('../common/logger')('api.list');

module.exports = {
	cls : {
		person : require('./person')
	},
	call : async (cn, mn, json) => {
		log.debug('api.call:' + cn + ':' + mn);
		let cls = this.cls[cn];
		if (cls===undefined||cls[mn]===undefined) {
			throw 'api not found';
		}
		try {
			let rslt = await cls[mn](JSON.parse(json));
			return JSON.stringify(rslt);
		} catch (e) {
			throw e;
		}
	}
};
