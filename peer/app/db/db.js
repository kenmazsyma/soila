'use_strict'

const pgp = require('pg-promise')();
const allconf = require('config');
const db = pgp({
	host: allconf.DBEnv.host,
	port: allconf.DBEnv.port,
	user: allconf.DBEnv.user,
	password: allconf.DBEnv.pass,
	database: allconf.DBEnv.dbname
});

module.exports = {
	init : async () => {
	},
	term : async () => {
	},
	any : async (sql, param) => {
		if (param===undefined) {
			return db.any(sql);
		}
		return db.any(sql, param);
	},
	none : async (sql, param) => {
		if (param===undefined) {
			return db.none(sql);
		}
		return db.none(sql, param);
	},
};

