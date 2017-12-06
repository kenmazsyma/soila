'use_strict';

var Client = require('fabric-client');
var path = require('path');
var util = require('util');
var fs = require('fs');
var os = require('os');
var cmn = require('../common/util');
var log = require('../common/logger')('blockchain.fabricclientbase');


function convertPayload(data) {
	let enc = cmn.str.b642s(data);
	return JSON.parse(enc);
}

FabricCliBase = class {

	constructor(chid, conf) {
		this.client = null;
		this.channelID = chid;
		this.channel = null;
		this.genesis_block = null;
		this.targets = null;
		this.eventhub = [];
		this.isConnect = false;
		this.conf = conf;
		this.orderer = null;
	}

	async init() {
		this.client = new Client();
		this.channel = this.client.newChannel(this.channelID);
		this.orderer = this.client.newOrderer(this.conf.orderer.url);
		this.channel.addOrderer(this.orderer);
		let store = await Client.newDefaultKeyValueStore({
			path: path.join(os.tmpdir(), 'fabricclibase/orderer1')
		});
		this.client.setStateStore(store);
		let keyPath = path.join(__dirname, this.conf.orderer.admin.keystore);
		let keyPEM = Buffer.from(readAllFiles(keyPath)[0]).toString();
		let crtPath = path.join(__dirname, this.conf.orderer.admin.signcerts);
		let crtPEM = readAllFiles(crtPath)[0];
		let admin = await this.client.createUser({
			username : this.conf.orderer.admin.username,
			mspid : this.conf.orderer.mspid,
			cryptoContent: {
				privateKeyPEM: keyPEM.toString(),
				signedCertPEM: crtPEM.toString()
			}
		});
		this.genesis_block = await this.channel.getGenesisBlock({
			txId: this.client.newTransactionID(),
			orderer: {
				url : this.conf.orderer.url
			}
		});
		this.client._userContext = null;
		let cres = [];
		this.conf.org.forEach(async (d) => {
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
			cres.push(this.client.createUser({
				username: d.admin.username,
				mspid: d.mspid,
				cryptoContent: {
					privateKeyPEM: keyPEM.toString(),
					signedCertPEM: crtPEM.toString()
				}
			}));
		});
		await Promise.all(cres);
		this.targets = [];
		this.eventhub = [];
		this.conf.org.forEach(d => {
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
				block : /*this.genesis_block*/ null,
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
	}

	async joinChannel() {
		let request = {
			targets : this.targets,
			block : this.genesis_block,
			txId : this.client.newTransactionID()
		};
		let txPromise = [];
		this.eventhub.forEach(hub => {
			txPromise.push(new Promise((reso, reje) => {
				let handle = setTimeout(reje, 30000);
				hub.registerBlockEvent((block) => {
					clearTimeout(handle);
					if(block.data.data.length === 1) {
						let ch_header = block.data.data[0].payload.header.channel_header;
						if (ch_header.channel_id === 'soila') {
							log.info('The new channel has been successfully joined on peer '+ hub.getPeerAddr());
							reso();
						}
						else {
							log.error('The new channel has not been succesfully joined');
						}
					} else {
						log.error('Response of registerBlockEvent is not correct.');
					}
					reje();
				});
			}));
		});
		let sendPromise = this.channel.joinChannel(request);
		txPromise.push(sendPromise);
		await Promise.all(txPromise);
		this.isConnect = true;
	}
	
	async prepareChannel() {
		try {
			await this.channel.initialize();
			log.info('channel is successfully initiated:');
			this.isConnect = true;
		} catch (e) {
			log.error('FabricClientBase.prepareChannel:');
			throw e;
		}
	}

	async invoke(ccid, funcname, args) {
		log.info('data:' + funcname);
		let request = {
			chaincodeId : ccid,
			targets : this.targets,
			fcn: funcname,
			args: args,
			chainId : 'soila',
			txId: this.client.newTransactionID()
		};
		try {
			let results = await this.channel.sendTransactionProposal(request);
			await this.sendTransaction(request.txId, results);
			let rslt0 = results[0][0].response;
			return {
				status : rslt0.status,
				message : rslt0.message,
				data : convertPayload(rslt0.payload)
			};
		} catch (e) {
			log.error(e);
			return {
				status : 500,
				message : 'failed to convert payload data received from legder.'
			};
		}
	}

	async install(ccid, path, ver) {
		let request = {
			targets: this.targets,
			chaincodePath: path,
			chaincodeId: ccid,
			chaincodeVersion: ver
		};
		return this.client.installChaincode(request);
	}

	async instantiate(ccid, ver, args) {
		let request = {
			txId : this.client.newTransactionID(),
			chaincodeId : ccid,
			chaincodeVersion: ver,
			targets: this.targets,
			args: args
		};
		let results = await this.channel.sendInstantiateProposal(request);
		return this.sendTransaction(request.txId, results);
	}

	async upgrade(ccid, ver, args) {
		let request = {
			txId : this.client.newTransactionID(),
			chaincodeId : ccid,
			chaincodeVersion: ver,
			targets: this.targets,
			args: args
		};
		let results = await this.channel.sendUpgradeProposal(request);
		return this.sendTransaction(request.txId, results);
	}

	async term() {
		return new Promise(resolve => {
			setTimeout(()=>{
				resolve();
			}, 2000);
			this.eventhub.forEach(hub => {
				hub.disconnect();
			});
			this.client = null;
			this.channelID = null;
			this.channel = null;
			this.genesis_block = null;
			this.targets = null;
			this.eventhub = [];
			this.isConnect = false;
		});
	}

	async sendTransaction(txid, results) {
		var proposalResponses = results[0];
		var proposal = results[1];
		var all_good = true;
		for (var i in proposalResponses) {
			let one_good = false;
			if (proposalResponses && proposalResponses[i].response &&
				proposalResponses[i].response.status === 200) {
				one_good = true;
				log.info('proposal was good');
			} else {
				log.error('proposal was bad');
			}
			all_good = all_good & one_good;
		}
		if (all_good) {
			let txPromise = [];
			this.eventhub.forEach((hub) => {
				txPromise.push(new Promise((resolve, reject) => {
					let handle = setTimeout(() => {
						this.term();
						reject();
					}, 30000);
					let deployId = txid;
					hub.registerTxEvent(deployId, (tx, code) => {
						clearTimeout(handle);
						hub.unregisterTxEvent(deployId);
						hub.disconnect();
						if (code !== 'VALID') {
							log.error('transaction(' + tx + ') was invalid, code = ' + code);
							reject();
						} else {
							log.info('transactioni(' + tx + ') was valid.');
							resolve();
						}
					});
				}));
			});
			let req = {
				proposalResponses: results[0],
				proposal: results[1]
			};
			let sendPromise = this.channel.sendTransaction(req);
			return Promise.all([sendPromise].concat(txPromise));
		}
		throw 'failed to send proposal or receive valid response. response null or status is not 200. exiting...'
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
