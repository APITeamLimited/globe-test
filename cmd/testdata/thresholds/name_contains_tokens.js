//  The threshold name contains the '***REMOVED***' and '***REMOVED***' characters, which
// are used as tokens when parsing the submetric part of a threshold's
// name. This pattern occurs when following the URL grouping pattern, and
// should not error.");
//
// Protects from #2512 regressions.
export const options = ***REMOVED***
	thresholds: ***REMOVED***
		"http_req_duration***REMOVED***name:http://$***REMOVED******REMOVED***.com***REMOVED***": ["max < 1000"],
	***REMOVED***,
***REMOVED***;

export default function () ***REMOVED***
	console.log(
		"asserting a threshold's name containing parsable tokens is valid"
	);
***REMOVED***
