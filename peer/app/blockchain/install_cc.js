'use_strict';

let $$ = require('./prepare.js');
let log = require('../common/logger')('blockchain.install_cc');

$$.cli.init().then(() => {
	return $$.cli.prepareChannel();
}).then(() => {
	return $$.cli.install('soila_chain', 'github.com/kenmazsyma/soila/chaincode', '1');
}).then((ret) => {
	log.info('successfully installed chaincode');
	$$.cli.term();
}).catch((err) => {
	log.error('failed to install chaincode:' + err);
	$$.cli.term();
});


