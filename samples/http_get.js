import http from 'k6/http';
import ***REMOVED*** Counter***REMOVED*** from "k6/metrics";

export function setup() ***REMOVED***
  open("asdf");
***REMOVED***
// export let CounterErrors = new Counter("Errors");

// export let options = ***REMOVED***
// 	thresholds: ***REMOVED***
//     "Errors": [***REMOVED*** threshold: "count<5", abortOnFail: true ***REMOVED***]
//   ***REMOVED***
// ***REMOVED***;

export default function() ***REMOVED***
  const response = http.get("http://test.loadimpact.com");
  // CounterErrors.add(true);
***REMOVED***;
