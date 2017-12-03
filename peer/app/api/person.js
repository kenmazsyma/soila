'use_strict'

let db = require('../db/db');
let cmn = require('./cmn');
let log = require('../common/logger')('api.person');
let bc = require('../blockchain/prepare');
let util = require('../common/util');

module.exports = {
	register : async function(prm) {
		// register PERSON data to db
		let flds = {'id':true, 'name':true, 'pass':true, 'profile':false };
		let chk = cmn.chkParams(prm, flds);
		if (chk!==null) {
			throw chk + ' is mandaroty.';
		}
		prm = cmn.fillParams(prm, flds);
		try {
			await db.query(
				"insert into person(id, name, pass, profile) values($1,$2,$3,$4)",
				[prm.id, prm.name, prm.pass, prm.profile]
			);
			log.debug('query for inserting to person table succeeded.');
		} catch (e) {
			log.error('query for inserting to person table errored.');
			throw e;
		}
		// register PERSON data to ledger
		let rslt = await bc.cli.invoke('soila_chain', 'person.register', 
						[prm.id, JSON.stringify({
							id : prm.id,
							name : prm.name,
							profile : prm.profile
						})]);
		// update PERSON data on db according to the result of registering on ledger
		if (rslt.status===200) {
			try {
				return await db.query(
						"update person set ledgerkey=$1 where id=$2", 
						[rslt.data[0], prm.id]
				);
			} catch (e) {
				log.error('failed to update ledgerkey on db');
				throw e;
			}
		} else {
			try {
				return await db.query(
						"delete from person where id=$1", 
						[prm.id]
				);
			} catch (e) {
				log.error('failed to delete PERSON record on db');
				// TODO:necessary to resolve mismatch between db and ledger
				throw e;
			}
		}
	},
	delete : async function(prm) {
	},
	get : async function(prm) {
	},
	update : async function(prm) {
	}
};
