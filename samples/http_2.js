import http from "k6/http";
import ***REMOVED*** check ***REMOVED*** from "k6";

export default function () ***REMOVED***
  check(http.get("https://test-api.k6.io/"), ***REMOVED***
    "status is 200": (r) => r.status == 200,
    "protocol is HTTP/2": (r) => r.proto == "HTTP/2.0",
  ***REMOVED***);
***REMOVED***
