import http from "k6/http";
import ***REMOVED*** check ***REMOVED*** from "k6";

export default function() ***REMOVED***
    // Passing username and password as part of URL will authenticate using HTTP Basic Auth
    let r = http.get("http://user:passwd@httpbin.org/basic-auth/user/passwd");

    // Verify response
    let j = r.json();
    check(r, ***REMOVED***
        "status is 200": (r) => r.status === 200,
        "is authenticated": (r) => j["authenticated"] === true,
        "is correct user": (r) => j["user"] === "user"
    ***REMOVED***);
***REMOVED***
