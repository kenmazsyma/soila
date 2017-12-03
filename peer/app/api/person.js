'use_strict'

let db = require('../db/db');
let cmn = require('./cmn');
let log = require('../common/logger')('api.person');
let bc = require('../blockchain/prepare');
let util = require('../common/util');

module.exports = {
	create : function(prm) {
		return new Promise((resolve, reject) => {
			let flds = {'id':true, 'name':true, 'pass':true, 'profile':false };
			let chk = cmn.chkParams(prm, flds);
			if (chk!==null) {
				reject(chk + ' is mandaroty.');
				return;
			}
			prm = cmn.fillParams(prm, flds);
			try {
				db.query("insert into person(id, name, pass, profile) values($1,$2,$3,$4)",
						[prm.id, prm.name, prm.pass, prm.profile]
				).then(res => {
					log.debug('query for inserting to person table succeeded.');
					resolve();
				}).catch(e=> {
					log.debug('query for inserting to person table failed.');
					reject(e);
				});
			} catch (e) {
				log.error('query for inserting to person table raise exception.');
				reject(e);
			}
		}).then(() => {
			try {
				return bc.cli.invoke('soila_chain', 'person.register', 
								[prm.id, JSON.stringify({
									id : prm.id,
									name : prm.name,
									profile : prm.profile
								})]);
			} catch (e) {
				log.error(e);
				return Promise.resolve({
					status:500,
					message:e
				});
			}
		}).then(rslt => {
			if (rslt.status===200) {
				try {
					return db.query("update person set ledgerkey=$1 where id=$2", [rslt.data[0], prm.id]).then(() => {
						return Promise.resolve('OK');
					}).catch(e=> {
						// TODO:necessary to remove data from ledger
						log.error(e);
						return Promise.reject(e);
					});
				} catch (e) {
					log.debug('12347');
					log.error(e);
					return Promise.reject(e);
				}
			} else {
				try {
					return db.query("delete from person where id=$1", [prm.id]).then(() => {
						return Promise.reject(rslt.message);
					}).catch(e=> {
						log.error(e);
						// TODO:necessary to resolve mismatch between db and ledger
						return Promise.reject(e);
					});
				} catch (e) {
					log.error(e);
					// TODO:necessary to resolve mismatch between db and ledger
					return Promise.reject(e);
				}
			}
			
		});
	},
	delete : function(prm) {
	},
	get : function(prm) {
	},
	update : function(prm) {
	}
};
