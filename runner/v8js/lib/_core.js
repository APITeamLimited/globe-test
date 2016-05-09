__internal__ = ***REMOVED***
	_modules: ***REMOVED******REMOVED***,
	_data: ***REMOVED******REMOVED***,
***REMOVED***;

__internal__._invoke = function(mod, fn, args, async) ***REMOVED***
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
