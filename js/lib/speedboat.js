export function group(name, fn, cond) ***REMOVED***
	if (cond !== undefined && !cond) ***REMOVED***
		return
	***REMOVED***

	return __vu_impl__.DoGroup(name, fn);
***REMOVED***

export function test(name, ...sets) ***REMOVED***
	return __vu_impl__.DoTest(name, ...sets);
***REMOVED***

export function sleep(secs) ***REMOVED***
	__vu_impl__.Sleep(secs);
***REMOVED***
