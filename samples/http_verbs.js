import http from "k6/http";
import ***REMOVED*** check, group ***REMOVED*** from "k6";

/*
 * k6 supports all standard HTTP verbs/methods:
 * CONNECT, DELETE, GET, HEAD, OPTIONS, PATCH, POST, PUT and TRACE.
 * 
 * Below are examples showing how to use the most common of these.
 */

export default function() ***REMOVED***
    // GET request
    group("GET", function() ***REMOVED***
        let r = http.get("http://httpbin.org/get?verb=get");
        let j = r.json();
        check(r, ***REMOVED***
            "status is 200": (r) => r.status === 200,
            "is verb correct": (r) => j["args"]["verb"] === "get",
        ***REMOVED***);
    ***REMOVED***);

    // POST request
    group("POST", function() ***REMOVED***
        let r = http.post("http://httpbin.org/post", ***REMOVED*** verb: "post" ***REMOVED***);
        let j = r.json();
        check(r, ***REMOVED***
            "status is 200": (r) => r.status === 200,
            "is verb correct": (r) => j["form"]["verb"] === "post",
        ***REMOVED***);
    ***REMOVED***);

    // PUT request
    group("PUT", function() ***REMOVED***
        let r = http.put("http://httpbin.org/put", JSON.stringify(***REMOVED*** verb: "put" ***REMOVED***), ***REMOVED*** headers: ***REMOVED*** "Content-Type": "application/json" ***REMOVED******REMOVED***);
        let j = r.json();
        check(r, ***REMOVED***
            "status is 200": (r) => r.status === 200,
            "is verb correct": (r) => j["json"]["verb"] === "put",
        ***REMOVED***);
    ***REMOVED***);

    // PATCH request
    group("PATCH", function() ***REMOVED***
        let r = http.patch("http://httpbin.org/patch", JSON.stringify(***REMOVED*** verb: "patch" ***REMOVED***), ***REMOVED*** headers: ***REMOVED*** "Content-Type": "application/json" ***REMOVED******REMOVED***);
        let j = r.json();
        check(r, ***REMOVED***
            "status is 200": (r) => r.status === 200,
            "is verb correct": (r) => j["json"]["verb"] === "patch",
        ***REMOVED***);
    ***REMOVED***);

    // DELETE request
    group("DELETE", function() ***REMOVED***
        let r = http.del("http://httpbin.org/delete?verb=delete");
        let j = r.json();
        check(r, ***REMOVED***
            "status is 200": (r) => r.status === 200,
            "is verb correct": (r) => j["args"]["verb"] === "delete",
        ***REMOVED***);
    ***REMOVED***);
***REMOVED***
