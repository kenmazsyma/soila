'use_strict';

let $$ = require('./prepare.js');

$$.cli.init().then(() => {
	return $$.cli.prepareChannel();
}).then(() => {
	console.log('###instantiate###');
	return $$.cli.instantiate('soila_chain', '0', ['init', 'a', '1']);
}).then((ret) => {
	console.log('successfully instantiated chaincode' + JSON.stringify(ret));
	$$.cli.term();
}).catch((err) => {
	console.log('failed to instantiate chaincode:' + err);
	$$.cli.term();
});


