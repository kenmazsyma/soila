'use_strict';

let http = require('http');
let fs = require('fs');
let allconf = require('config');
let api = require('../api/list');

module.exports = {
	run : function() {
		http.createServer()
		.on('request', request)
		.listen(
			allconf.WebEnv.port, 
			allconf.WebEnv.url
		);
	}
}

const cont_type = {
	'html': 'text/html',
	'htm': 'text/html',
	'css':  'text/css',
	'js':   'application/x-javascript',
	'json': 'application/json',
	'jpg':  'image/jpeg',
	'jpeg': 'image/jpeg',
	'png':  'image/png',
	'gif':  'image/gif',
	'svg':  'image/svg+xml'
}

function getType(url) {
	let typ = url.split('.');
	typ = typ[typ.length-1];
	return typ||'';
}

function conttype(url) {
	return cont_type[getType(url)]||'';
}

function request(req, res) {
	if (getType(req.url)==='json') {
		return callApi(req, res);
	}
	fs.readFile(__dirname + '/static' + req.url, 'utf-8', function(err, data){
		if (err) {
			return return404(res);
		}
		res.writeHead(200, {'Content-Type' : conttype(req.url)});　
		res.end(data);
	});
}

function return404(res, e) {
	res.writeHead(404, {'Content-Type': 'text/plain'});
	if (e===undefined) {
		e = 'Not Found';
	}
	res.write(JSON.stringify(e));
	return res.end();　
}

function callApi(req, res) {
	new Promise((resolve, reject) => {
		setTimeout(() => {
			reject();
		}, 2000);
		req.on('readable', function() {
			var data = req.read();
			if (data !== null) {
				resolve(new Buffer(data).toString('utf-8'));
			}
		});
	}).then((data)=> {
		/([^\./]*)\.([^\./]*)\.json$/.exec(req.url);
		api.call(RegExp.$1, RegExp.$2, data).then(rslt =>  {
			res.writeHead(200, {'Content-Type': 'application/json'});
			res.write(rslt);
			res.end();
		}).catch(e => {
			return404(res,e);
		});
	}).catch(e => {
		return404(res,e);
	});
}

