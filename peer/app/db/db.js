'use_strict'

const { Pool } = require('pg');
let allconf = require('config');

let pool = null;

module.exports = {
	init : function() {
		pool = new Pool ({
			user: allconf.DBEnv.user,
			host: allconf.DBEnv.host,
			database: allconf.DBEnv.dbname,
			password: allconf.DBEnv.pass,
			port: allconf.DBEnv.port
		});
	},
	term : async function() {
		if (pool) {
			let p = pool;
			pool = null;
			return p.end();
		} else {
			return 'ok';
		}
	},
	query : function(sql, param) {
		if (param===undefined) {
			return pool.query(sql);
		}
		return pool.query(sql, param);
	}
};

