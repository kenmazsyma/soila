'use_strict';

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
	//return cli.prepareChannel();
	return cli.joinChannel();
}).then(() => {
	//return cli.install('soila_chain', 'github.com/kenmazsyma/soila/chaincode/chaincode', '0');
	return cli.install('soila_chain', 'github.com/kenmazsyma/soila/chaincode', '0');
}).then((ret) => {
	console.log('successfully installed chaincode' + JSON.stringify(ret));
	return sleep(5000);
}).catch((err) => {
	console.log(err);
});


