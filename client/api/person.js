'use_strict'

let db = require('../db/db');
let cmn = require('./cmn');
let log = require('../common/logger')('api.person');
let bc = require('../blockchain/prepare');

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
						error : (e.code!==undefined) ? e.code + ':' + e.detail : e
					});
				});
			} catch (e) {
				log.error('query for inserting to person table raise exception.');
				reject(e);
			}
		}).then(() => {
			return bc.invoke('soila_chain', 'person.put',[prm.id, JSON.stringify(prm)])
		});
	},
	delete : function(prm) {
	},
	get : function(prm) {
	},
	update : function(prm) {
	}
};
