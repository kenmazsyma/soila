'use_strict';

let $$ = require('./prepare.js');

$$.cli.init().then(() => {
	return $$.cli.prepareChannel();
}).then(() => {
	//return cli.install('soila_chain', 'github.com/kenmazsyma/soila/chaincode/chaincode', '0');
	return $$.cli.install('soila_chain', 'github.com/kenmazsyma/soila/chaincode', '0');
}).then((ret) => {
	console.log('successfully installed chaincode' + JSON.stringify(ret));
	$$.cli.term();
}).catch((err) => {
	console.log(err);
	$$.cli.term();
});


