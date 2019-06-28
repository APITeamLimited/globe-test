import http from "k6/http";
import ***REMOVED*** check, group, sleep ***REMOVED*** from "k6";
import ***REMOVED*** Rate ***REMOVED*** from "k6/metrics";

// A custom metric to track failure rates
var failureRate = new Rate("check_failure_rate");

// Options
export let options = ***REMOVED***
    stages: [
        // Linearly ramp up from 1 to 50 VUs during first minute
        ***REMOVED*** target: 50, duration: "1m" ***REMOVED***,
        // Hold at 50 VUs for the next 3 minutes and 30 seconds
        ***REMOVED*** target: 50, duration: "3m30s" ***REMOVED***,
        // Linearly ramp down from 50 to 0 50 VUs over the last 30 seconds
        ***REMOVED*** target: 0, duration: "30s" ***REMOVED***
        // Total execution time will be ~5 minutes
    ],
    thresholds: ***REMOVED***
        // We want the 95th percentile of all HTTP request durations to be less than 500ms
        "http_req_duration": ["p(95)<500"],
        // Requests with the staticAsset tag should finish even faster
        "http_req_duration***REMOVED***staticAsset:yes***REMOVED***": ["p(99)<250"],
        // Thresholds based on the custom metric we defined and use to track application failures
        "check_failure_rate": [
            // Global failure rate should be less than 1%
            "rate<0.01",
            // Abort the test early if it climbs over 5%
            ***REMOVED*** threshold: "rate<=0.05", abortOnFail: true ***REMOVED***,
        ],
    ***REMOVED***,
***REMOVED***;

// Main function
export default function () ***REMOVED***
    let response = http.get("https://test.loadimpact.com/");

    // check() returns false if any of the specified conditions fail
    let checkRes = check(response, ***REMOVED***
        "http2 is used": (r) => r.proto === "HTTP/2.0",
        "status is 200": (r) => r.status === 200,
        "content is present": (r) => r.body.indexOf("Welcome to the LoadImpact.com demo site!") !== -1,
    ***REMOVED***);

    // We reverse the check() result since we want to count the failures
    failureRate.add(!checkRes);

    // Load static assets, all requests
    group("Static Assets", function () ***REMOVED***
        // Execute multiple requests in parallel like a browser, to fetch some static resources
        let resps = http.batch([
            ["GET", "https://test.loadimpact.com/style.css", null, ***REMOVED*** tags: ***REMOVED*** staticAsset: "yes" ***REMOVED*** ***REMOVED***],
            ["GET", "https://test.loadimpact.com/images/logo.png", null, ***REMOVED*** tags: ***REMOVED*** staticAsset: "yes" ***REMOVED*** ***REMOVED***]
        ]);
        // Combine check() call with failure tracking
        failureRate.add(!check(resps, ***REMOVED***
            "status is 200": (r) => r[0].status === 200 && r[1].status === 200,
            "reused connection": (r) => r[0].timings.connecting == 0,
        ***REMOVED***));
    ***REMOVED***);

    sleep(Math.random() * 3 + 2); // Random sleep between 2s and 5s
***REMOVED***