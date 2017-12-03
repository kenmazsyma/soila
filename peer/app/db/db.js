'use_strict'

const { Pool } = require('pg');
let allconf = require('config');

let pool = null;

module.exports = {
	init : () => {
		pool = new Pool ({
			user: allconf.DBEnv.user,
			host: allconf.DBEnv.host,
			database: allconf.DBEnv.dbname,
			password: allconf.DBEnv.pass,
			port: allconf.DBEnv.port
		});
	},
	term : async () => {
		if (pool) {
			let p = pool;
			pool = null;
			return p.end();
		} else {
			return 'ok';
		}
	},
	query : (sql, param) => {
		if (param===undefined) {
			return pool.query(sql);
		}
		return pool.query(sql, param);
	}
};

