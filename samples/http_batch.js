import ***REMOVED*** check ***REMOVED*** from 'k6';
import http from 'k6/http';

export default function() ***REMOVED***
  const responses = http.batch([
    "http://test.k6.io",
    "http://test.k6.io/pi.php"
  ]);

  check(responses[0], ***REMOVED***
    "main page 200": res => res.status === 200,
  ***REMOVED***);

  check(responses[1], ***REMOVED***
    "pi page 200": res => res.status === 200,
    "pi page has right content": res => res.body === "3.14",
  ***REMOVED***);
***REMOVED***;
