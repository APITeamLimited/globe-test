import http from "k6/http";
import ***REMOVED*** check ***REMOVED*** from "k6";

export let options = ***REMOVED***
    // When this option is enabled (set to true), all of the verifications
    // that would otherwise be done to establish trust in a server provided
    // TLS certificate will be ignored.
    insecureSkipTLSVerify: true
***REMOVED***;

export default function() ***REMOVED***
    let res = http.get("https://httpbin.org/");
    check(res, ***REMOVED*** "status is 200": (r) => r.status === 200 ***REMOVED***);
***REMOVED***
