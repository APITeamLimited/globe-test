__internal__._register = function(mod, obj) ***REMOVED***
	if (!(mod in __internal__._modules)) ***REMOVED***
		__internal__._modules[mod] = ***REMOVED******REMOVED***;
	***REMOVED***
	for (k in Object.keys(obj)) ***REMOVED***
		__internal__._modules[mod][k] = obj[k];
	***REMOVED***
***REMOVED***

function require(mod) ***REMOVED***
	if (!(mod in __internal__._modules)) ***REMOVED***
		throw new Error("module not found: " + mod);
	***REMOVED***
	return __internal__._modules[mod];
***REMOVED***
