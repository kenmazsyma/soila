'use_strict';

//var sv = require('./cli/sv');
//sv.init();
var FabricCliBase = require('./fabricclibase');

var cli = new FabricCliBase('soila');
cli.init(() => {
	cli.prepareChannel(() => {
		console.log('finish');
	});
},(err) => {
});
