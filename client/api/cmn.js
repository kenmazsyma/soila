'use_strict'

module.exports = {
	chkParams : (prm, req) => {
		for (var i in req) {
			if (req[i]===true&&prm[i]===undefined) return i;
		}
		return null;
	},

	fillParams : (prm, req) => {
		for (var i in req) {
			if (req[i]===false&&prm[i]===undefined) {
				prm[i] = '';
			}
		}
		return prm;
	}
};
