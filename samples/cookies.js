import http from "k6/http";
import ***REMOVED*** check, group ***REMOVED*** from "k6";

export default function() ***REMOVED***
    group("Simple cookies", function() ***REMOVED***
        let cookies = ***REMOVED***
            name: "value1",
            name2: "value2"
        ***REMOVED***;
        let r = http.get("http://httpbin.org/cookies", ***REMOVED*** cookies: cookies ***REMOVED***);
        check(r, ***REMOVED***
            "status is 200": (r) => r.status === 200,
            "has cookie": (r) => r.cookies["name"].length > 0
        ***REMOVED***);
    ***REMOVED***);

    group("Advanced cookies", function() ***REMOVED***
        let cookie = ***REMOVED*** name: "name3", value: "value3", domain: "httpbin.org", path: "/cookies" ***REMOVED***;
        let r = http.get("http://httpbin.org/cookies", ***REMOVED*** cookies: [cookie] ***REMOVED***);
        check(r, ***REMOVED***
            "status is 200": (r) => r.status === 200,
            "has cookie": (r) => r.cookies["name3"].length > 0
        ***REMOVED***);
    ***REMOVED***);
***REMOVED***
