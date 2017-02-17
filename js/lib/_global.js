// https://developer.mozilla.org/en/docs/Web/API/console
console = ***REMOVED***
	log(msg, ...args) ***REMOVED*** console.info(msg, ...args); ***REMOVED***,

	debug(msg, ...args) ***REMOVED*** __console__.Log(0, msg, args); ***REMOVED***,
	info(msg, ...args) ***REMOVED*** __console__.Log(1, msg, args); ***REMOVED***,
	warn(msg, ...args) ***REMOVED*** __console__.Log(2, msg, args); ***REMOVED***,
	error(msg, ...args) ***REMOVED*** __console__.Log(3, msg, args); ***REMOVED***,
***REMOVED***;

// https://developer.mozilla.org/en/docs/Web/API/performance
performance = ***REMOVED***
	now() ***REMOVED*** return __jsapi__.ElapsedMs(); ***REMOVED***
***REMOVED***;

require = function(name) ***REMOVED***
	if (!__initapi__) ***REMOVED***
		throw new Error("imports are only permitted during initialization")
	***REMOVED***
	return __initapi__.Require(name);
***REMOVED***

open = function(name) ***REMOVED***
	if (!__initapi__) ***REMOVED***
		throw new Error("open() is only permitted during initialization")
	***REMOVED***
	return __initapi__.Open(name);
***REMOVED***
