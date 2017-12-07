'use_strict'

let db = require('../db/db');
let cmn = require('./cmn');
let log = require('../common/logger')('api.peer');
let bc = require('../blockchain/prepare');
let util = require('../common/util');

module.exports = {
	register : async prm => {
		let ip = util.net.local();
		if (!ip.v4||ip.v4.length===0) {
			throw 'failed to get local ip address';
		}
		let rslt = await bc.cli.invoke('soila_chain', 'peer.register', 
						util.str.s2b64([ip.v4[0]]));
		// update PERSON data based on the result of registering on ledger
		if (rslt.status&&rslt.status===200) {
			log.info('success to register peer');
		} else {
			throw rslt.message;
		}
	}
}
