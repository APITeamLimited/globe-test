"use strict";

$log = $log || ***REMOVED******REMOVED***;

$log.debug = function() ***REMOVED***
	return $log.log.apply(this, _.concat(['debug'], arguments));
***REMOVED***

$log.info = function() ***REMOVED***
	return $log.log.apply(this, _.concat(['info'], arguments));
***REMOVED***

$log.warn = function() ***REMOVED***
	return $log.log.apply(this, _.concat(['warn'], arguments));
***REMOVED***

$log.warning = function() ***REMOVED***
	return $log.warn.apply(this, arguments);
***REMOVED***

$log.error = function() ***REMOVED***
	return $log.log.apply(this, _.concat(['error'], arguments));
***REMOVED***
