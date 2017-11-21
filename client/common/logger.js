let log4js = require('log4js');
let allconf = require('config');

log4js.configure(allconf.Logger);

module.exports = (name => {
    return log4js.getLogger(name);
    //error: log4js.getLogger('error')
});
