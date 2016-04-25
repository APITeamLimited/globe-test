speedboat = ***REMOVED***
	_internal: ***REMOVED***
		recv: ***REMOVED******REMOVED***,
	***REMOVED***,
***REMOVED***;

$recvSync(function(msg) ***REMOVED***
	d = JSON.parse(msg);
	fn = speedboat._internal.recv[d.call];
	if (fn !== undefined) ***REMOVED***
		fn.apply(speedboat, d.args);
	***REMOVED***
***REMOVED***);
