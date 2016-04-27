function require(mod) ***REMOVED***
	if (!(mod in speedboat._modules)) ***REMOVED***
		throw new Error("module not found: " + mod);
	***REMOVED***
	return speedboat._modules[mod];
***REMOVED***
