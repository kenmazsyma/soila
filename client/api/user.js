'use_strict'

let db = require('../db/db');

module.exports = {
	create : function(prm) {
		return new Promise((resolve, reject) => {
			try {
				db.query("insert into person(id, name, pass, profile) values($1,$2,$3,$4)",
						[prm.id, prm.name, prm.pass, prm.profile]
				).then(res => {
					console.log('res:' + JSON.stringify(res));
					resolve({
						result : 'OK'
					});
				}).catch(e=> {
					console.log('e:' + JSON.stringify(e));
					resolve({
						result : 'ERROR',
						error : e
					});
				});
			} catch (e) {
				reject(e);
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
