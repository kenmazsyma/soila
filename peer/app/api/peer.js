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
			log.info('success to register peer');
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
			log.info('success to register peer');
		} else {
			throw rslt.message;
		}
	}


}
