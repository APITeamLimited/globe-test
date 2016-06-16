"use strict";

function require(name) ***REMOVED***
	var mod = __modules__[name];
	if (!mod) ***REMOVED***
		throw new Error("Unknown module: " + name);
	***REMOVED***
	return mod;
***REMOVED***
