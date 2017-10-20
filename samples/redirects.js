import http from "k6/http";
import ***REMOVED***check***REMOVED*** from "k6";

export let options = ***REMOVED***
    // Max redirects to follow (default is 10)
    maxRedirects: 5
***REMOVED***;

export default function() ***REMOVED***
    // If redirecting more than options.maxRedirects times, the last response will be returned
    let res = http.get("https://httpbin.org/redirect/6");
    check(res, ***REMOVED***
        "is status 302": (r) => r.status === 302
    ***REMOVED***);

    // The number of redirects to follow can be controlled on a per-request level as well
    res = http.get("https://httpbin.org/redirect/1", ***REMOVED***redirects: 1***REMOVED***);
    console.log(res.status);
    check(res, ***REMOVED***
        "is status 200": (r) => r.status === 200,
        "url is correct": (r) => r.url === "https://httpbin.org/get"
    ***REMOVED***);
***REMOVED***
