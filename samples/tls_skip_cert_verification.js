import http from "k6/http";
import ***REMOVED*** check ***REMOVED*** from "k6";

export let options = ***REMOVED***
    // When this option is enabled (set to true), mismatches in hostname
    // between target system and TLS certificate will be ignored
    insecureSkipTLSVerify: true
***REMOVED***;

export default function() ***REMOVED***
    let r = http.get("https://httpbin.org/");
    check(r, ***REMOVED*** "status is 200": (r) => r.status === 200 ***REMOVED***);
***REMOVED***
