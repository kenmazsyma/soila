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
	return cli.install('soila_chain', 'github.com/kenmazsyma/soila/chaincode/chaincode', '0');
}).then((ret) => {
	console.log('successfully installed chaincode' + JSON.stringify(ret));
	return sleep(5000);
}).then(() => {
	return cli.instantiate('soila_chain', '0');
}).then((ret) => {
	console.log('successfully instantiated chaincode' + JSON.stringify(ret));
	return sleep(5000);
}).then(() => {
	return cli.invoke('soila_chain', 'person.put', ['123', 'testdata']);
}).then((ret) => {
	console.log('successfully invoked person.put' + JSON.stringify(ret));
	return sleep(5000);
}).then(() => {
	return cli.invoke('soila_chain', 'person.get', ['123']);
}).then((ret) => {
	console.log('successfully invoked person.get' + JSON.stringify(ret));
}).catch((err) => {
	console.log(err);
});


