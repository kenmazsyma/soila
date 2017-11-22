'use_strict'

var assert = require('assert');
var db = require('./db');


describe('db', () => {
	before(() => {
		db.init();
	});

	after(async function()  {
		await db.term();
	});

	it('query', (done) => {
		db.query("select 'a' as a").then((rslt)=> {
			assert.equal('a', rslt.rows[0].a);
			done();
		}).catch((err)=>{
			done(err);
		});
	});
});

