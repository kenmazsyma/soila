'use_strict';

let sv = require('./sv/sv');
let db = require('./db/db');
let bc = require('./blockchain/prepare');
let log = require('./common/logger')('main');


bc.cli.init().then(() => {
	return bc.cli.prepareChannel();
}).then(()=> {
	log.info('blockchain is successfully started.');
	sv.run();
	db.init();
}).catch(e => {
	log.error('starting blockchain failed.');
	log.error(e);
	term();
});

async function term() {
	await db.term();
	await bc.cli.term();
	log.info('all connections are successfully terminated.');
	process.exit(0);
}

process.on('beforeExit', term).on('SIGINT', term);
