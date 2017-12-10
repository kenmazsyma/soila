'use_strict'

let db = require('../db/db');
let cmn = require('./cmn');
let log = require('../common/logger')('api.peer');
let bc = require('../blockchain/prepare');
let util = require('../common/util');

module.exports = {
	register : async () => {
		let ip = util.net.local();
		if (ip.v4.length===0) {
			throw 'failed to get local ip address';
		}
		let rslt = await bc.cli.invoke('soila_chain', 'peer.register', 
						util.str.s2b64([ip.v4[0]]));
		if (rslt.status&&rslt.status===200) {
			log.info('success to register peer:' + rslt.data[0]);
			return rslt.data[0];
		} else {
			throw rslt.message;
		}
	},

	update : async () => {
		let ip = util.net.local();
		if (ip.v4.length===0) {
			throw 'failed to get local ip address';
		}
		let rslt = await bc.cli.invoke('soila_chain', 'peer.update', 
						util.str.s2b64([ip.v4[0]]));
		if (rslt.status&&rslt.status===200) {
			log.info('success to update peer:' + rslt.data[0]);
			return rslt.data[0]
		} else {
			throw rslt.message;
		}
	},

	get : async arg => {
		if (arg.constructor!==String) {
			return 'parameter is not as expect.'
		}
		let rslt = await bc.cli.invoke('soila_chain', 'peer.get', [arg]);
		if (rslt.status&&rslt.status===200) {
			log.info('success to get peer');
			return rslt.data[1];
		} else {
			throw rslt.message;
		}
	},

	deregister : async () => {
		let rslt = await bc.cli.invoke('soila_chain', 'peer.deregister', []);
		if (rslt.status&&rslt.status===200) {
			log.info('success to deregister peer');
			return;
		} else {
			throw rslt.message;
		}
	}

}
