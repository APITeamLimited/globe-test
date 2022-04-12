export const options = ***REMOVED***
	thresholds: ***REMOVED***
		// non existing is neither registered, nor a builtin metric.
		// k6 should catch that.
		"non existing": ["rate>0"],
	***REMOVED***,
***REMOVED***;

export default function () ***REMOVED***
	console.log(
		"asserting that a threshold over a non-existing metric fails with exit code 104 (Invalid config)"
	);
***REMOVED***
