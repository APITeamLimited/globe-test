import http from "k6/http";
import ***REMOVED*** check, group ***REMOVED*** from "k6";

export let options = ***REMOVED***
    maxRedirects: 3
***REMOVED***;

export default function() ***REMOVED***
    // VU cookie jar
    group("Simple cookies send with VU jar", function() ***REMOVED***
        let cookies = ***REMOVED***
            name: "value1",
            name2: "value2"
        ***REMOVED***;
        let res = http.get("http://httpbin.org/cookies", ***REMOVED*** cookies: cookies ***REMOVED***);
        check(res, ***REMOVED***
            "status is 200": (r) => r.status === 200,
            "has cookie 'name'": (r) => r.json().cookies.name.length > 0,
            "has cookie 'name2'": (r) => r.json().cookies.name2.length > 0
        ***REMOVED***);

        // Since the cookies are set as "request cookies" they won't be added to VU cookie jar
        let vuJar = http.cookieJar();
        let cookiesForURL = vuJar.cookiesForURL(res.url);
        check(null, ***REMOVED***
            "vu jar doesn't have cookie 'name'": () => cookiesForURL.name === undefined,
            "vu jar doesn't have cookie 'name2'": () => cookiesForURL.name2 === undefined
        ***REMOVED***);
    ***REMOVED***);

    group("Simple cookies set with VU jar", function() ***REMOVED***
        // Since this request redirects the `res.cookies` property won't contain the cookies
        let res = http.get("http://httpbin.org/cookies/set?name3=value3&name4=value4");
        check(res, ***REMOVED***
            "status is 200": (r) => r.status === 200
        ***REMOVED***);

        // Make sure cookies have been added to VU cookie jar
        let vuJar = http.cookieJar();
        let cookiesForURL = vuJar.cookiesForURL(res.url);
        check(null, ***REMOVED***
            "vu jar has cookie 'name3'": () => cookiesForURL.name3.length > 0,
            "vu jar has cookie 'name4'": () => cookiesForURL.name4.length > 0
        ***REMOVED***);
    ***REMOVED***);

    // Local cookie jar
    group("Simple cookies send with local jar", function() ***REMOVED***
        let jar = new http.CookieJar();
        let cookies = ***REMOVED***
            name5: "value5",
            name6: "value6"
        ***REMOVED***;
        let res = http.get("http://httpbin.org/cookies", ***REMOVED*** cookies: cookies, jar: jar ***REMOVED***);
        check(res, ***REMOVED***
            "status is 200": (r) => r.status === 200,
            "has cookie 'name5'": (r) => r.json().cookies.name5.length > 0,
            "has cookie 'name6'": (r) => r.json().cookies.name6.length > 0
        ***REMOVED***);

        // Since the cookies are set as "request cookies" they won't be added to VU cookie jar
        let cookiesForURL = jar.cookiesForURL(res.url);
        check(null, ***REMOVED***
            "local jar doesn't have cookie 'name5'": () => cookiesForURL.name5 === undefined,
            "local jar doesn't have cookie 'name6'": () => cookiesForURL.name6 === undefined
        ***REMOVED***);

        // Make sure cookies have NOT been added to VU cookie jar
        let vuJar = http.cookieJar();
        cookiesForURL = vuJar.cookiesForURL(res.url);
        check(null, ***REMOVED***
            "vu jar doesn't have cookie 'name5'": () => cookiesForURL.name === undefined,
            "vu jar doesn't have cookie 'name6'": () => cookiesForURL.name2 === undefined
        ***REMOVED***);
    ***REMOVED***);

    group("Advanced send with local jar", function() ***REMOVED***
        let jar = new http.CookieJar();
        jar.set("http://httpbin.org/cookies", "name7", "value7");
        jar.set("http://httpbin.org/cookies", "name8", "value8");
        let res = http.get("http://httpbin.org/cookies", ***REMOVED*** jar: jar ***REMOVED***);
        let cookiesForURL = jar.cookiesForURL(res.url);
        check(res, ***REMOVED***
            "status is 200": (r) => r.status === 200,
            "has cookie 'name7'": (r) => r.json().cookies.name7.length > 0,
            "has cookie 'name8'": (r) => r.json().cookies.name8.length > 0
        ***REMOVED***);

        cookiesForURL = jar.cookiesForURL(res.url);
        check(null, ***REMOVED***
            "local jar has cookie 'name7'": () => cookiesForURL.name7.length > 0,
            "local jar has cookie 'name8'": () => cookiesForURL.name8.length > 0
        ***REMOVED***);

        // Make sure cookies have NOT been added to VU cookie jar
        let vuJar = http.cookieJar();
        cookiesForURL = vuJar.cookiesForURL(res.url);
        check(null, ***REMOVED***
            "vu jar doesn't have cookie 'name7'": () => cookiesForURL.name7 === undefined,
            "vu jar doesn't have cookie 'name8'": () => cookiesForURL.name8 === undefined
        ***REMOVED***);
    ***REMOVED***);

    group("Advanced cookie attributes", function() ***REMOVED***
        let jar = http.cookieJar();
        jar.set("http://httpbin.org/cookies", "name9", "value9", ***REMOVED*** domain: "httpbin.org", path: "/cookies" ***REMOVED***);

        let res = http.get("http://httpbin.org/cookies", ***REMOVED*** jar: jar ***REMOVED***);
        check(res, ***REMOVED***
            "status is 200": (r) => r.status === 200,
            "has cookie 'name9'": (r) => r.json().cookies.name9 === "value9"
        ***REMOVED***);

        jar.set("http://httpbin.org/cookies", "name10", "value10", ***REMOVED*** domain: "example.com", path: "/" ***REMOVED***);
        res = http.get("http://httpbin.org/cookies", ***REMOVED*** jar: jar ***REMOVED***);
        check(res, ***REMOVED***
            "status is 200": (r) => r.status === 200,
            "doesn't have cookie 'name10'": (r) => r.json().cookies.name10 === undefined
        ***REMOVED***);
    ***REMOVED***);
***REMOVED***
