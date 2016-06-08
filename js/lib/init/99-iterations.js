// Wrap the script in a function that increments the iteration counter.
__data__.iteration = 1;
__script__ = function(script) ***REMOVED***
	return function() ***REMOVED***
		__data__.iteration++;
		script();
	***REMOVED***
***REMOVED***(__script__);
