import http from "k6/http";
import ***REMOVED*** check ***REMOVED*** from "k6";


/*
 * Stages (aka ramping) is how you, in code, specify the ramping of VUs.
 * That is, how many VUs should be active and generating traffic against
 * the target system at any specific point in time for the duration of
 * the test.
 * 
 * The following stages configuration will result in up-flat-down ramping
 * profile over a 20s total test duration.
 */ 

export let options = ***REMOVED***
    stages: [
        // Ramp-up from 1 to 5 VUs in 10s
        ***REMOVED*** duration: "10s", target: 5 ***REMOVED***,

        // Stay at rest on 5 VUs for 5s
        ***REMOVED*** duration: "5s", target: 5 ***REMOVED***,

        // Ramp-down from 5 to 0 VUs for 5s
        ***REMOVED*** duration: "5s", target: 0 ***REMOVED***
    ]
***REMOVED***;

export default function() ***REMOVED***
    let res = http.get("http://httpbin.org/");
    check(res, ***REMOVED*** "status is 200": (r) => r.status === 200 ***REMOVED***);
***REMOVED***
