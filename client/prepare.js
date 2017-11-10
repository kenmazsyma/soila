'use_strict';

let allconf = require('config');
let FabricCliBase = require('./fabricclibase');
let conf = {
	orderer : allconf.ChainEnv.orderer,
	org : [ allconf.ChainEnv[process.argv[2]||'org1'] ]
};

module.exports = {
	cli : new FabricCliBase('soila', conf),
	sleep : function(time) {
		return new Promise((resolve, reject) => {
			setTimeout(() => {
				resolve();
			}, time);
		});
	}
};
