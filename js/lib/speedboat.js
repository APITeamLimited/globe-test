export function group(name, fn, cond) ***REMOVED***
	if (cond !== undefined && !cond) ***REMOVED***
		return
	***REMOVED***

	return __jsapi__.DoGroup(name, fn);
***REMOVED***

export function test(name, ...sets) ***REMOVED***
	return __jsapi__.DoTest(name, ...sets);
***REMOVED***

export function sleep(secs) ***REMOVED***
	__jsapi__.Sleep(secs);
***REMOVED***

export default ***REMOVED***
	group: group,
	test: test,
	sleep: sleep,
***REMOVED***;
