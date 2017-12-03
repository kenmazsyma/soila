let person = require('../api/person');
let db = require('../db/db');
let assert = require("assert");
let bc = require('../blockchain/prepare');

async function term() {
	await db.term();
	await bc.cli.term();
}

describe('register', () => {
	before(done => {
		bc.cli.init().then(() => {
			return bc.cli.prepareChannel();
		}).then(()=> {
			db.init();
			db.query("delete from person where id='test'").then(res => {
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
			console.log(JSON.stringify());
			done(rslt);
		}).catch(e => {
			console.log(e);
			done();
		});
	});
});

