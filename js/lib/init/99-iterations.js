// The first iteration is iteration 1.
__data__.iteration = 1;

// Wrap the script in a function that increments the iteration counter.
__script__ = function(script) ***REMOVED***
	return function() ***REMOVED***
		script();
		__data__.iteration++;
	***REMOVED***
***REMOVED***(__script__);
