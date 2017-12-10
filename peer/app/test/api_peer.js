let peer = require('../api/peer');
let assert = require("assert");
let bc = require('../blockchain/prepare');

async function term() {
	await bc.cli.term();
}

describe('peer', () => {
	before(function(done) {
		this.timeout(30000);
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
	
//	describe('not found case', () => {
//		it('update', function(done) {
//			this.timeout(30000);
//			peer.update().then(rslt => {
//				done('data found in spite of not registering yet');
//			}).catch(e => {
//				done();
//			});
//		});
//	});

	let key;
	describe('register', () => {
		it('success', function(done) {
			this.timeout(30000);
			peer.register().then(rslt => {
				key = rslt;
				console.log('KEY:' + key);
				done();
			}).catch(e => {
				done(e);
			});
		}, e => {
			console.log(e);
		});
//		it('duplicate', function(done) {
//			this.timeout(30000);
//			peer.register().then(rslt => {
//				done('succeeded in spite of duplicating');
//			}).catch(e => {
//				done();
//			});
//		}, e => {
//			console.log(e);
//		});
	});

//	describe('update2', () => {
//		it('success', function(done) {
//			this.timeout(30000);
//			peer.update().then(rslt => {
//				done();
//			}).catch(e => {
//				done(e);
//			});
//		});
//	});

	describe('get', () => {
		it('success', function(done) {
			this.timeout(30000);
			peer.get(key).then(rslt => {
				if (rslt&&rslt.address) {
					done();
				} else {
					done('failed to get data:' + JSON.stringify(rslt));
				}
			}).catch(e => {
				done(e);
			});
		});
//		it('not found', function(done) {
//			this.timeout(30000);
//			peer.get('test').then(rslt => {
//				done('found in spite of not existing');
//			}).catch(e => {
//				done();
//			});
//		});
	});


//	describe('deregister', () => {
//		it('success', function(done) {
//			this.timeout(30000);
//			peer.deregister().then(rslt => {
//				done();
//			}).catch(e => {
//				done(e);
//			});
//		});
//		it('not found', function(done) {
//			this.timeout(30000);
//			peer.deregister().then(rslt => {
//				done('succeeded in spite of already deregistered');
//			}).catch(e => {
//				done();
//			});
//		});
//	});

});
