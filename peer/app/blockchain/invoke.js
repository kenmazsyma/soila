'use_strict';

let $$ = require('./prepare');
let log = require('../common/logger')('bc_main');

$$.cli.init().then(() => {
	return $$.cli.prepareChannel();
}).then((ret) => {
	log.info('channel is successfully prepared.');
	return $$.cli.invoke('soila_chain', 'person.put', ['123', 'testdata']);
}).then((ret) => {
	log.info('invoking person.put is successfully called.' + JSON.stringify(ret));
	$$.cli.term();
}).catch((err) => {
	log.error('failed to invoke person.put:' + err);
	$$.cli.term();
});


