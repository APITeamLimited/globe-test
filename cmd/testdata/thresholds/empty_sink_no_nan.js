// In scenarios where a threshold would apply to a rate metric
// that would not receive any samples (settting abortToFail emphasis the issue),
// division by zero could occur and lead to NaN values being returned.
//
// Protects from #2520 regressions.
import ***REMOVED*** Rate ***REMOVED*** from "k6/metrics";

const rate = new Rate("rate");

export const options = ***REMOVED***
	thresholds: ***REMOVED***
		"rate***REMOVED***type:read***REMOVED***": [***REMOVED*** threshold: "rate>0.9", abortOnFail: true ***REMOVED***],
	***REMOVED***,
***REMOVED***;

export default function () ***REMOVED***
	console.log("not interacting with rate metric");
***REMOVED***
