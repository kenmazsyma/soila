'use strict';

var Client = require('fabric-client');
var path = require('path');
var util = require('util');
var fs = require('fs');
var os = require('os');

var client;
var channel;
var genesis_block;
var targets;

Promise.resolve().then(()=> {
	client = new Client();
	Client.setLogger({
		debug : function(txt, txt2, txt3) {
			/*var log = 'DEBUG';
			for ( var i=0; i<arguments.length; i++) {
				if (typeof arguments[i]==='string') {
					log += ' : ' + arguments[i];
				} else {
					try {
						log += JSON.stringify(arguments[i]);
						continue;
					} catch (e) {
					}
					if (arguments[i].toString) {
						log += ' : ' + arguments[i].toString();
					}
				}
			}
			console.log(log);*/
		},
		info : function(txt) {
			console.log('INFO : ' + txt)
		},
		warn: function(txt) {
			console.log('WARN : ' + txt)
		},
		error : function(txt) {
			console.log('ERROR : ' + txt)
		},
	});
	channel = client.newChannel('soila');
	//let data = fs.readFileSync(path.join(__dirname, 'key/crypto-config/ordererOrganizations/soila.com/orderers/orderer.soila.com/msp/tlscacerts/tlsca.soila.com-cert.pem'));
	//let data = fs.readFileSync(path.join(__dirname, 'key/crypto-config/ordererOrganizations/soila.com/msp/tlscacerts/tlsca.soila.com-cert.pem'));
	//let data = fs.readFileSync(path.join(__dirname, 'key/crypto-config/ordererOrganizations/soila.com/tlsca/tlsca.soila.com-cert.pem'));
	//let data = fs.readFileSync(path.join(__dirname, 'key/crypto-config/ordererOrganizations/soila.com/ca/ca.soila.com-cert.pem'));
	//let caroots = Buffer.from(data).toString();

	channel.addOrderer(
		client.newOrderer(
			'grpc://localhost:7050'/*,
			{
				'pem': caroots,
				'ssl-target-name-override': 'orderer.soila.com'
			}*/
		)
	);
	console.log('point1');
	return Client.newDefaultKeyValueStore({
		path: path.join(os.tmpdir(), 'hfc/hfc_org')
	});
}).then((store) => {
	console.log('point2');
	client.setStateStore(store);
	var keyPath = path.join(__dirname, 'key/crypto-config/ordererOrganizations/soila.com/users/Admin@soila.com/msp/keystore');
	var keyPEM = Buffer.from(readAllFiles(keyPath)[0]).toString();
	var certPath = path.join(__dirname, 'key/crypto-config/ordererOrganizations/soila.com/users/Admin@soila.com/msp/signcerts');
	var certPEM = readAllFiles(certPath)[0];

	return client.createUser({
		username : 'ordererAdmin',
		mspid : 'OrdererMSP',
		cryptoContent: {
			privateKeyPEM: keyPEM.toString(),
			signedCertPEM: certPEM.toString()
		}
	});
}).then((admin) => {
	return channel.getGenesisBlock({
		txId: client.newTransactionID(),
		orderer: {
			url : 'grpc://localhost:7050'
		}
	});
}).then((block)=> {
	console.log('successfully got the genesis block');
	genesis_block = block;
	
	// get the peer org's admin required to send join channel requests
	client._userContext = null;

	var keyPath = path.join(__dirname, 'key/crypto-config/peerOrganizations/org1.soila.com/users/Admin@org1.soila.com/msp/keystore');
	var keyPEM = Buffer.from(readAllFiles(keyPath)[0]).toString();
	var certPath = path.join(__dirname, 'key/crypto-config/peerOrganizations/org1.soila.com/users/Admin@org1.soila.com/msp/signcerts');
	var certPEM = readAllFiles(certPath)[0];
	var cryptoSuite = Client.newCryptoSuite();
	cryptoSuite.setCryptoKeyStore(Client.newCryptoKeyStore({path: path.join(os.tmpdir(), 'hfc/hfc_org1')}));
	client.setCryptoSuite(cryptoSuite);
	return client.createUser({
		username: 'peerOrg1Admin',
		mspid: 'Org1MSP',
		cryptoContent: {
			privateKeyPEM: keyPEM.toString(),
			signedCertPEM: certPEM.toString()
		}
	});
}).then((admin) => {
	console.log('successfully prepared admin user:' );
	let peers = {
		org1 : {
			rpc : 'grpc://localhost:7051',
			evtrpc : 'grpc://localhost:7053',
			host : 'peer0.org1.soila.com',
			cert : 'key/crypto-config/peerOrganizations/org1.soila.com/tlsca/tlsca.org1.soila.com-cert.pem'
		},
		org2 : {
			rpc : 'grpc://locahost:8051',
			evtrpc : 'grpc://locahost:8053',
			host : 'peer0.org2.soila.com',
			cert : 'key/crypto-config/peerOrganizations/org2.soila.com/tlsca/tlsca.org2.soila.com-cert.pem'
		},
		org3 : {
			rpc : 'grpc://locahost:8551',
			evtrpc : 'grpc://locahost:8553',
			host : 'peer0.org3.soila.com',
			cert : 'key/crypto-config/peerOrganizations/org3.soila.com/tlsca/tlsca.org3.soila.com-cert.pem'
		}
	};
	let p = peers['org1'];
	let data = fs.readFileSync(path.join(__dirname, p['cert']));
	targets = [
		client.newPeer(
			p['rpc'],
			{
				pem: Buffer.from(data).toString(),
				'ssl-target-name-override': p['host']
			}
		)
	];
	let request = {
		targets : targets,
		block : genesis_block,
		txId : client.newTransactionID()
	};
	let eh = client.newEventHub();
	eh.setPeerAddr(
		p['evtrpc'],
		{
			pem: Buffer.from(data).toString(),
			'ssl-target-name-override': p['host']
		}
	);
	eh.connect();
	/*let txPromise = new Promise((resolve, reject) => {
		let handle = setTimeout(reject, 30000);
		eh.registerBlockEvent((block) => {
			console.log('callback for registerBlockEvent has been called');
			clearTimeout(handle);
			// in real-world situations, a peer may have more than one channel so
			// we must check that this block came from the channel we asked the peer to join
			if(block.data.data.length === 1) {
				// Config block must only contain one transaction
				var channel_header = block.data.data[0].payload.header.channel_header;
				if (channel_header.channel_id === 'soila') {
					console.log('The new channel has been successfully joined on peer '+ eh.getPeerAddr());
					resolve();
				}
				else {
					console.log('The new channel has not been succesfully joined');
					reject();
				}
			}
		});
	});
	let sendPromise = channel.joinChannel(request);
	return Promise.all([txPromise, sendPromise]);*/
	return channel.initialize();
}).then((results)=> {
	console.log('successfully prepared.');
	var tx_id = client.newTransactionID();
	var request = {
		chaincodeId : 'soila',
		targets : targets,
		fcn: 'move',
		args: ['a', 'b','100'],
		txId: tx_id,
	};
	return channel.sendTransactionProposal(request);
}).then((results)=> {
	console.log('response : ' + JSON.stringify(results));

}).catch((err)=> {
	console.log(err);
});

function readAllFiles(dir) {
	var files = fs.readdirSync(dir);
	var certs = [];
	files.forEach((file_name) => {
		let file_path = path.join(dir,file_name);
		// console.log(' looking at file ::'+file_path);
		let data = fs.readFileSync(file_path);
		certs.push(data);
	});
	return certs;
}
