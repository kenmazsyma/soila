let peer = require('../api/peer');
let assert = require("assert");
let bc = require('../blockchain/prepare');

async function term() {
	await bc.cli.term();
}

describe('peer', () => {
	before(done => {
		bc.cli.init().then(() => {
			return bc.cli.prepareChannel();
		}).then(()=> {
			done();
		}).catch(e => {
			term();
			done();
		});
	});
	after(done => {
		term();
		done();
	});
	
	describe('update1', done => {
		it('not found', function(done) {
			this.timeout(30000);
			peer.update().then(rslt => {
				done();
			}).catch(e => {
				done(e);
			});
		});
	});


	describe('register', done => {
		it('success', function(done) {
			this.timeout(30000);
			peer.register().then(rslt => {
				done();
			}).catch(e => {
				done(e);
			});
		}, e => {
			console.log(e);
		});
		it('duplicate', function(done) {
			this.timeout(30000);
			peer.register().then(rslt => {
				done();
			}).catch(e => {
				done(e);
			});
		}, e => {
			console.log(e);
		});
	});

	describe('update2', done => {
		it('success', function(done) {
			this.timeout(30000);
			peer.update().then(rslt => {
				done();
			}).catch(e => {
				done(e);
			});
		});
	});

});
