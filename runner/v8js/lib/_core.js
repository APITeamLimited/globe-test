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

speedboat._require.interface = function(v) ***REMOVED***
	return v;
***REMOVED***

speedboat._require.struct = function(spec, v) ***REMOVED***
	for (key in Object.keys(v)) ***REMOVED***
		v[key] = speedboat._require[spec[v]];
	***REMOVED***
	return v;
***REMOVED***

speedboat._invoke = function(mod, fn, args, async) ***REMOVED***
	var send = async ? $send : $sendSync
	var res = send(JSON.stringify(***REMOVED*** m: mod, f: fn, a: args ***REMOVED***));
	if (res) ***REMOVED***
		var obj = JSON.parse(res);
		if (obj._error) ***REMOVED***
			throw new Error(obj._error);
		***REMOVED***
		return obj;
	***REMOVED***
***REMOVED***

$recvSync(function(raw) ***REMOVED***
	if (raw == 'run') ***REMOVED***
		__run__();
	***REMOVED***
***REMOVED***);
