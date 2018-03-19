import http from "k6/http";
import ***REMOVED*** check ***REMOVED*** from "k6";

export default function() ***REMOVED***
    // Passing username and password as part of URL plus the auth option will authenticate using HTTP Digest authentication
    let res = http.get("http://user:passwd@httpbin.org/digest-auth/auth/user/passwd", ***REMOVED***auth: "digest"***REMOVED***);

    // Verify response
    check(res, ***REMOVED***
        "status is 200": (r) => r.status === 200,
        "is authenticated": (r) => r.json().authenticated === true,
        "is correct user": (r) => r.json().user === "user"
    ***REMOVED***);
***REMOVED***
