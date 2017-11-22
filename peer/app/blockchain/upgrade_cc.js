'use_strict';

let $$ = require('./prepare.js');
let log = require('../common/logger')('blockchain.instantiate_cc');

$$.cli.init().then(() => {
	return $$.cli.prepareChannel();
}).then(() => {
	log.info('successfully prepared channel.');
	return $$.cli.upgrade('soila_chain', '1', ['init', 'a', '1']);
}).then((ret) => {
	log.info('successfully upgraded chaincode');
	$$.cli.term();
}).catch((err) => {
	log.error('failed to upgrade chaincode:' + err);
	$$.cli.term();
});


