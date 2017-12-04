let person = require('../api/person');
let db = require('../db/db');
let assert = require("assert");
let bc = require('../blockchain/prepare');
//bc.cli = require('./fabricstub');

async function term() {
	await db.term();
	await bc.cli.term();
}

describe('person', () => {
	before(done => {
		//bc.cli.chaincode = {
		//	kv : {},
		//	'person.register' : (args) => {
		//	},
		//	'person.get' : (args) => {
		//	},
		//	'peer.get' : (args) => {
		//	}
		//};
		bc.cli.init().then(() => {
			return bc.cli.prepareChannel();
		}).then(()=> {
			db.init();
			db.none("delete from person where id='test'").then(res => {
				done();
			}).catch(e => {
				done(e);
			});
		}).catch(e => {
			term();
			done();
		});
	});
	after(done => {
		term();
		done();
	});
	describe('register', () => {
		it('success', done => {
			person.register({
				"id":"test",
				"pass":"testtest",
				"name":"tester",
				"profile":"This is a test data."
			}).then(rslt => {
				console.log(JSON.stringify());
				done();
			}).catch(e => {
				console.log(e);
				done(e);
			});
		});
		it('duplicate', done => {
			person.register({
				"id":"test",
				"pass":"testtest",
				"name":"tester",
				"profile":"This is a test data."
			}).then(rslt => {
				console.log(JSON.stringify(rslt));
				done(rslt);
			}).catch(e => {
				console.log(e);
				done();
			});
		});
	});
	let key = '';
	describe('getbyid', () => {
		it('success', done => {
			person.getbykey({id:'test'}).then(rslt => {
				if (rslt&&rslt.id) {
					done();
					key = rslt.key;
				} else {
					done('failed to get data');
				}
			}).catch(e => {
				done(e);
			});
		});

		it('not found', done => {
			person.getbykey({id:'test1'}).then(rslt => {
				if (rslt.length>0) {
					done('got unexpected data.');
				} else {
					done();
				}
			}).catch(e => {
				done(e);
			});
		});
	});
	describe('get', () => {
		it('found', done => {
			person.get({key:key}).then(rslt => {
				console.log(rslt);
				if (rslt&&rslt.found) {
					done();
				} else {
					done('failed to get data');
				}
			}).catch(e => {
				done(e);
			});
		});
		it('not found', done => {
			person.get({key:'1'}).then(rslt => {
				if (rslt&&rslt.found) {
					console.log(rslt);
					done('got unexpected data.');
				} else {
					console.log(rslt);
					done('got unexpected data.');
				}
			}).catch (e => {
				console.log(e);
				done();
			});
		});
	});
	describe('update', () => {
		it('success', done => {
			person.update({id:'test',pass:'testtest2', name:'あいうえお'}).then(rslt => {
				done();
			}).catch (e => {
				done(e);
			});
		});
	});

});

