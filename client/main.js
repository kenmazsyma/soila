'use_strict';

let $$ = require('./prepare.js');
//var sv = require('./cli/sv');
//sv.init();

$$.cli.init().then(() => {
	return $$.cli.prepareChannel();
}).then(() => {
	console.log('###invoke1###');
	return $$.cli.invoke('soila_chain', 'person.put', ['123', 'testdata']);
}).then((ret) => {
	console.log('successfully invoked person.put' + JSON.stringify(ret));
//	return sleep(5000);
//}).then(() => {
//	console.log('###invoke2###');
//	return cli.invoke('soila_chain', 'person.get', ['123']);
//}).then((ret) => {
//	console.log('successfully invoked person.get' + JSON.stringify(ret));
	$$.cli.term();
}).catch((err) => {
	console.log(err);
	$$.cli.term();
});


