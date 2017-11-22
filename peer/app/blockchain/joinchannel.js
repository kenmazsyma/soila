'use_strict';

let $$ = require('./prepare.js');

$$.cli.init().then(() => {
	return $$.cli.joinChannel();
}).then(() => {
	console.log('successfully joined.');
	$$.cli.term();
}).catch((err) => {
	console.log(err);
	$$.cli.term();
});


