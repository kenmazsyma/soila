'use_strict';

let $$ = require('./prepare.js');
var log = require('../common/logger')('blockchain.joinchannel');

$$.cli.init().then(() => {
	return $$.cli.joinChannel();
}).then(() => {
	log.info('successfully joined with channel.');
	$$.cli.term();
}).catch((err) => {
	log.error('failed to join with channel:' + err);
	$$.cli.term();
});


