'use_strict'

module.exports = {
	chkParams : (prm, req) => {
		for (let i in req) {
			if (req[i]===true&&prm[i]===undefined) return i;
		}
		return null;
	},

	fillParams : (prm, req) => {
		for (let i in req) {
			if (req[i]===false&&prm[i]===undefined) {
				prm[i] = '';
			}
		}
		return prm;
	},

	chkEither : (prm, req) => {
		for (let i in req) {
			if (req[i]===false&&prm[i]!==undefined) {
				return true;
			}
		}
		return false;
	},

	buildUpdate : (tb, prm, flds) => {
		let ix = 1;
		let nms = [], ret = [], keys = [];
		for (let i in prm) {
			if (flds[i]===false) {
				nms.push(i + '=$' + (ix++));
				ret.push(prm[i]);
			} else if (flds[i]===true) {
				keys.push(i + '=$' + (ix++));
				ret.push(prm[i]);
			}
		}
		return {
			sql : 'update ' + tb + ' set ' + nms.join(', ') + ' where ' + keys.join(' and '),
			prm : ret
		};
	}
};
