import http from "k6/http";
import ***REMOVED*** check ***REMOVED*** from "k6";

export default function() ***REMOVED***
    let res = http.get("https://stackoverflow.com");
    check(res, ***REMOVED***
        "is OCSP response good": (r) => r.ocsp.status === http.OCSP_STATUS_GOOD
    ***REMOVED***);
***REMOVED***;
