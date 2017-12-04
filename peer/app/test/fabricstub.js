'use_strict'

module.exports = {
	chaincode : {
		kv : {},
	},
	init : async () => {
	},
	joinChannel : async () => {
	},
	prepareChannel : async () => {
	},
	invoke : async (ccid, funcname, args) => {
		let rslt = this.chaincode[funname]&&this.chaincode[funcname](args);
		if (!rslt) {
			return {
				status:500,
				message:'failed to run',
				payload:[{}]
			}
		}
		return rslt;
	},
	install : async (ccid, path, ver) => {
	},
	instantiate : async (ccid, ver, args) => {
	},
	upgrade : async (ccid, ver, args) => {
	},
	term : async () => {
	}
};
