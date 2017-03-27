import http from "k6/http";
import ***REMOVED*** check ***REMOVED*** from "k6";

export default function() ***REMOVED***
    // Send a JSON encoded POST request
    let body = JSON.stringify(***REMOVED*** key: "value" ***REMOVED***);
    let r = http.post("http://httpbin.org/post", body, ***REMOVED*** headers: ***REMOVED*** "Content-Type": "application/json" ***REMOVED******REMOVED***);

    // Use JSON.parse to deserialize the JSON (instead of using the r.json() method)
    let j = JSON.parse(r.body);

    // Verify response
    check(r, ***REMOVED***
        "status is 200": (r) => r.status === 200,
        "is key correct": (r) => j["json"]["key"] === "value",
    ***REMOVED***);
***REMOVED***
