'use strict';

// Backward-compatible default workload. The original sample pointed to a
// non-existent RegisterProduct function, so route it to the real SC01 workload.
module.exports = require('./registerActor');
