speedboat = ***REMOVED***
	_modules: ***REMOVED******REMOVED***,
	_require: ***REMOVED******REMOVED***,
***REMOVED***;

speedboat._require.float64 = function(v) ***REMOVED***
	out = parseFloat(v);
	if (isNaN(out)) ***REMOVED***
		throw new Error("not a float: " + v);
	***REMOVED***
	return out
***REMOVED***
speedboat._require.float32 = speedboat._require.float64

speedboat._require.int = function(v) ***REMOVED***
	out = parseInt(v);
	if (isNaN(out)) ***REMOVED***
		throw new Error("not an int: " + v);
	***REMOVED***
	return out
***REMOVED***

speedboat._require.string = function(v) ***REMOVED***
	return (v || "").toString();
***REMOVED***

$recvSync(function(raw) ***REMOVED***
	if (raw == 'run') ***REMOVED***
		__run__();
	***REMOVED***
***REMOVED***);
