"use strict";

__internal__.modules.log.debug = function() ***REMOVED***
	return __internal__.modules.log.type.apply(this, _.concat(['debug'], arguments));
***REMOVED***

__internal__.modules.log.info = function() ***REMOVED***
	return __internal__.modules.log.type.apply(this, _.concat(['info'], arguments));
***REMOVED***

__internal__.modules.log.warn = function() ***REMOVED***
	return __internal__.modules.log.type.apply(this, _.concat(['warn'], arguments));
***REMOVED***

__internal__.modules.log.warning = function() ***REMOVED***
	return __internal__.modules.log.warn(arguments);
***REMOVED***

__internal__.modules.log.error = function() ***REMOVED***
	return __internal__.modules.log.type.apply(this, _.concat(['error'], arguments));
***REMOVED***
