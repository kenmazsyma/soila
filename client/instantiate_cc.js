'use_strict';

//var sv = require('./cli/sv');
//sv.init();
var FabricCliBase = require('./fabricclibase');
var cli = new FabricCliBase('soila');

function sleep(time) {
	return new Promise((resolve, reject) => {
		setTimeout(() => {
			resolve();
		}, time);
	});
}


cli.init().then(() => {
	return cli.prepareChannel();
}).then(() => {
	console.log('###instantiate###');
	return cli.instantiate('soila_chain', '0');
}).then((ret) => {
	console.log('successfully instantiated chaincode' + JSON.stringify(ret));
	return sleep(5000);
}).catch((err) => {
	console.log(err);
});


