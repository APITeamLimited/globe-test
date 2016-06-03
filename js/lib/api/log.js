"use strict";

__modules__.log = __modules__.log || ***REMOVED******REMOVED***;

__modules__.log.debug = function() ***REMOVED***
	return __modules__.log.log.apply(this, _.concat(['debug'], arguments));
***REMOVED***

__modules__.log.info = function() ***REMOVED***
	return __modules__.log.log.apply(this, _.concat(['info'], arguments));
***REMOVED***

__modules__.log.warn = function() ***REMOVED***
	return __modules__.log.log.apply(this, _.concat(['warn'], arguments));
***REMOVED***

__modules__.log.warning = function() ***REMOVED***
	return __modules__.log.warn.apply(this, arguments);
***REMOVED***

__modules__.log.error = function() ***REMOVED***
	return __modules__.log.log.apply(this, _.concat(['error'], arguments));
***REMOVED***
