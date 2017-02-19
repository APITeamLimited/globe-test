import ***REMOVED*** check ***REMOVED*** from 'k6'
import http from 'k6/http'

export default function() ***REMOVED***
  const responses = http.batch([
    "http://test.loadimpact.com",
    "http://test.loadimpact.com/pi.php",
    ***REMOVED***
      "method": "GET",
      "url": "http://test.loadimpact.com/pi.php?decimals=50"
    ***REMOVED***
  ]);

  check(responses[0], ***REMOVED***
    "main page 200": res => res.status === 200,
  ***REMOVED***)

  check(responses[1], ***REMOVED***
    "pi page 200": res => res.status === 200,
  ***REMOVED***)

  check(responses[2], ***REMOVED***
    "pi page has 50 digits": res => ***REMOVED***
      return res.body === "3.14159265358979323846264338327950288"
    ***REMOVED***
  ***REMOVED***)
***REMOVED***;
