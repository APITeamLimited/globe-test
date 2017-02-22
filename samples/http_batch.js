import ***REMOVED*** check ***REMOVED*** from 'k6';
import http from 'k6/http';

export default function() ***REMOVED***
  const responses = http.batch(***REMOVED***
    "main": "http://test.loadimpact.com"
  ***REMOVED***);

  check(responses.main, ***REMOVED***
    "main page 200": res => res.status === 200
  ***REMOVED***)

  // check(responses[0], ***REMOVED***
  //   "main page 200": res => res.status === 200,
  // ***REMOVED***);

  // check(responses[1], ***REMOVED***
  //   "pi page 200": res => res.status === 200,
  //   "pi page has right content": res => res.body === "3.14",
  // ***REMOVED***);
***REMOVED***;
