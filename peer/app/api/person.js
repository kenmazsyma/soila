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
					resolve({
						result : 'OK'
					});
				}).catch(e=> {
					log.debug('query for inserting to person table failed.');
					resolve({
						result : 'ERROR',
						error : JSON.stringify(e)
					});
				});
			} catch (e) {
				log.error('query for inserting to person table raise exception.');
				reject(e);
			}
		}).then(() => {
			try {
			return bc.cli.invoke('soila_chain', 'person.put', [prm.id, JSON.stringify(prm)])
				.then(rslt => {
					for ( var i in rslt) {
						log.info(i + ':' + rslt[i]);
					}
					//console.log(util.str.b2s(rslt.payload.data));
					return Promise.resolve(rslt);
				});
			} catch (e) {
				log.error(e);
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
