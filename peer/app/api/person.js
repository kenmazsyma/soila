'use_strict'

let db = require('../db/db');
let cmn = require('./cmn');
let log = require('../common/logger')('api.person');
let bc = require('../blockchain/prepare');
let util = require('../common/util');

module.exports = {
	register : async prm => {
		// verify parameters
		let flds = {'id':true, 'name':true, 'pass':true, 'profile':false };
		let chk = cmn.chkParams(prm, flds);
		if (chk!==null) {
			throw chk + ' is mandaroty.';
		}
		prm = cmn.fillParams(prm, flds);
		// check if data specified by key is already registered
		try {
			let rslt = await db.any("select id from person where id=$1", [prm.id]);
			if (rslt.length>0) {
				console.log('!!!!!' + JSON.stringify(rslt));
				throw prm.id + ' is already registerd.';
			}
		} catch (e) {
			throw e;
		}
		// register PERSON data to db
		try {
			await db.none(
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
						util.str.s2b64([prm.id, JSON.stringify({
							id : prm.id,
							name : prm.name,
							profile : prm.profile
						})]));
		// update PERSON data based on the result of registering on ledger
		if (rslt.status&&rslt.status===200) {
			try {
				return await db.none(
						"update person set key=$1 where id=$2", 
						[rslt.data[0], prm.id]
				);
			} catch (e) {
				log.error('failed to update key on db');
				throw e;
			}
		} else {
			try {
				await db.none(
						"delete from person where id=$1", 
						[prm.id]
				);
				throw rslt;
			} catch (e) {
				log.error('failed to delete PERSON record on db');
				// TODO:necessary to resolve mismatch between db and ledger
				throw e;
			}
		}
	},
	delete : async prm => {
	},
	get : async prm => {
		// get from db if specifying PERSON data exists in this peer
		let rslt;
		try {
			rslt = await db.any(
					"select id, name, profile from person where key=$1",
					[prm.key]
			);
		} catch (e) {
			log.error('query for gettting PERSON data failed.');
			throw e;
		}
		if (rslt.length>1) {
			throw 'key duplicates on db';
		} else if (rslt.length===1) {
			return { found:true, data:rslt[0]};
		}
		// get PERSON data from ledger
		rslt = await bc.cli.invoke('soila_chain', 'person.get', [prm.key]);
		if (rslt.status!=200) {
			throw 'specifying PERSON is not exists in ledger.';
		}
		if (rslt.data.length<2||!rslt.data[1].peerkey) {
			throw 'getting PERSON data from ledger failed.';
		}
		// get PEER data from ledger
		rslt = await bc.cli.invoke('soila_chain', 'peer.get', [rslt.data[1].peerkey]);
		if (rslt.status!=200) {
			throw 'PEER data according to specifying PERSON is not exists in ledger.:' 
				 + rslt.data[1].peerkey;
		}
		if (rslt.data.length<2||!rslt.data[1].address) {
			throw 'getting PEER data according to specifying PERSON from ledger failed.';
		}
		return {found:false, peer:rslt.data[1].address};
	},
	getbyid : async prm => {
		try {
			let rslt = await db.any(
					"select id, name, profile, key from person where id=$1",
					[prm.id]
			);
			return (rslt.length===0) ? {} : rslt[0];
		} catch (e) {
			log.error('failed to get person data specified by id:' + prm.id);
			throw e
		}
	},
	update : async prm => {
		// verify parameters
		let flds = {'id':true, 'name':false, 'pass':false, 'profile':false };
		let chk = cmn.chkParams(prm, flds);
		if (chk!==null) {
			throw chk + ' is mandaroty.';
		}
		if (!cmn.chkEither(prm, flds)) {
			throw 'at least one parameter without id is needed.'
		}
		// update PERSON data on db
		try {
			let data = cmn.buildUpdate('person', prm, flds);
			await db.none(
				data.sql,
				data.prm
			);
		} catch (e) {
			console.log('failed to update person data specified by id');
			throw e;
		}
		// get updated data fron db
		let rslt;
		try {
			rslt = await db.any(
					"select id, name, profile from person where id=$1",
					[prm.id]
			);
			if (rslt.length!==1) {
				throw 'failed to get updated data from db';
			}
		} catch (e) {
			log.error('query for gettting PERSON data failed.');
			throw e;
		}
		// update PERSON data on ledger
		rslt = await bc.cli.invoke('soila_chain', 'person.update', 
						util.str.s2b64([rslt[0].id, JSON.stringify(rslt[0])]));
		if (rslt.status&&rslt.status===200) {
			return rslt.data[0];
		} else {
			throw 'failed to update ledger';
			// TODO:necessary to rollback
		}
	}
};
