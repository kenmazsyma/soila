'use_strict';

var Client = require('fabric-client');
var path = require('path');
var util = require('util');
var fs = require('fs');
var os = require('os');


FabricCliBase = class {

	constructor(chid, conf) {
		this.client = null;
		this.channelID = chid;
		this.channel = null;
		this.genesis_block = null;
		this.targets = null;
		this.eventhub = null;
		this.isConnect = false;
		this.conf = conf;
	}

	init() {
		return new Promise((resolve, reject) => {
			Promise.resolve().then(()=> {
				this.client = new Client();
				this.channel = this.client.newChannel(this.channelID);
				this.channel.addOrderer(
					this.client.newOrderer(this.conf.orderer.url)
				);
				return Client.newDefaultKeyValueStore({
					path: path.join(os.tmpdir(), 'fabricclibase/orderer1')
				});
			}).then((store) => {
				this.client.setStateStore(store);
				let keyPath = path.join(__dirname, this.conf.orderer.admin.keystore);
				let keyPEM = Buffer.from(readAllFiles(keyPath)[0]).toString();
				let crtPath = path.join(__dirname, this.conf.orderer.admin.signcerts);
				let crtPEM = readAllFiles(crtPath)[0];
				return this.client.createUser({
					username : this.conf.orderer.admin.username,
					mspid : this.conf.orderer.mspid,
					cryptoContent: {
						privateKeyPEM: keyPEM.toString(),
						signedCertPEM: crtPEM.toString()
					}
				});
			}).then((admin) => {
				return this.channel.getGenesisBlock({
					txId: this.client.newTransactionID(),
					orderer: {
						url : this.conf.orderer.url
					}
				});
			}).then((block)=> {
				this.genesis_block = block;
				this.client._userContext = null;
				var cres = [];
			//	[this.conf.org[0]].forEach((d) => {
				this.conf.org.forEach((d) => {
					let keyPath = path.join(__dirname, d.admin.keystore);
					let keyPEM = Buffer.from(readAllFiles(keyPath)[0]).toString();
					let crtPath = path.join(__dirname, d.admin.signcerts);
					let crtPEM = readAllFiles(crtPath)[0];
					var cryptoSuite = Client.newCryptoSuite();
					cryptoSuite.setCryptoKeyStore(Client.newCryptoKeyStore(
						{
							path: path.join(os.tmpdir(), 'hfc/'+d.mspid)
						}
					));
					this.client.setCryptoSuite(cryptoSuite);
					cres.push(
						this.client.createUser({
							username: d.admin.username,
							mspid: d.mspid,
							cryptoContent: {
								privateKeyPEM: keyPEM.toString(),
								signedCertPEM: crtPEM.toString()
							}
						})
					);
				});
				return Promise.all(cres);
			}).then((admin) => {
				this.targets = [];
				this.eventhub = [];
				//[this.conf.org[0]].forEach((d) => {
				this.conf.org.forEach((d) => {
					let data = fs.readFileSync(path.join(__dirname, d.cert));
					this.targets.push(
						this.client.newPeer(
							d.rpc,
							{
								pem: Buffer.from(data).toString(),
								'ssl-target-name-override': d.host
							}
						)
					);
					let request = {
						targets : this.targets,
						block : this.genesis_block,
						txId : this.client.newTransactionID()
					};
					var hub = this.client.newEventHub();
					this.eventhub.push(hub);
					hub.setPeerAddr(
						d.evtrpc,
						{
							pem: Buffer.from(data).toString(),
							'ssl-target-name-override': d.host
						}
					);
					hub.connect();
				});
				resolve();
			}).catch((err) => {
				console.log(err);
				reject(err);
			});
		});
	}

	joinChannel() {
		return new Promise((resolve, reject) => {
			let request = {
				targets : this.targets,
				block : this.genesis_block,
				txId : this.client.newTransactionID()
			};
			let txPromise = [];
			this.eventhub.forEach((hub) => {
				txPromise.push(new Promise((reso, reje) => {
					let handle = setTimeout(reje, 30000);
					hub.registerBlockEvent((block) => {
						clearTimeout(handle);
						if(block.data.data.length === 1) {
							let ch_header = block.data.data[0].payload.header.channel_header;
							if (ch_header.channel_id === 'soila') {
								console.log('The new channel has been successfully joined on peer '+ hub.getPeerAddr());
								reso();
							}
							else {
								console.log('The new channel has not been succesfully joined');
							}
						}
						reje();
					});
				}));
			});
			let sendPromise = this.channel.joinChannel(request);
			txPromise.push(sendPromise);
			Promise.all(txPromise).then((result) => {
				this.isConnect = true;
				resolve();
			}, (rslt)=>{ 
				reject(rslt) 
			});	
		});
	}
	
	prepareChannel() {
		return new Promise((resolve, reject) => {
			this.channel.initialize().then((result) => {
				this.isConnect = true;
				resolve();
			}, (err) => {
				console.log(err);
				reject(err);
			});
		});
	}

	invoke(ccid, funcname, args) {
		let request = {
			chaincodeId : ccid,
			targets : this.targets,
			fcn: funcname,
			args: args,
			txId: this.client.newTransactionID()
		};
		return this.channel.sendTransactionProposal(request);
	}

	install(ccid, path, ver) {
		let request = {
			targets: this.targets,
			chaincodePath: path,
			chaincodeId: ccid,
			chaincodeVersion: ver
		};
		return this.client.installChaincode(request);
//		this.client.installChaincode(request).then((results) => {
//			let proposalResponse = results[0];
//			//let proposal = results[1];
//			proposalResponse&&proposalResponse.forEach((elm) => {
//				let status = elm.response&&elm.response.status;
//				if (status === 200) {
//					console.log('successfully installed proposal');
//				} else {
//					console.log('failed to install proposal:' + status);
//				}
//			});
//			cb&&cb(results);
//		}).catch((err) => {
//			console.log(err);
//			cberr&&cberr(err);
//		});
	}

	instantiate(ccid, ver, args) {
		let request = {
			txId : this.client.newTransactionID(),
			chaincodeId : ccid,
			chaincodeVersion: ver,
			targets: this.targets,
			args: args
		};
		return this.channel.sendInstantiateProposal(request);
	}

	term() {
		this.eventhub.forEach((hub) => {
			hub.disconnect();
		});
		this.client = null;
		this.channelID = null;
		this.channel = null;
		this.genesis_block = null;
		this.targets = null;
		this.eventhub = null;
		this.isConnect = false;
	}
};

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

module.exports = FabricCliBase;
