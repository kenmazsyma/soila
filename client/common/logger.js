let log4js = require('log4js');
module.exports = {
    debug: log4js.getLogger('system'),
    info: log4js.getLogger('access'),
    error: log4js.getLogger('error')
}
