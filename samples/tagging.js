import http from "k6/http";
import ***REMOVED*** Trend ***REMOVED*** from "k6/metrics";
import ***REMOVED*** check ***REMOVED*** from "k6";

/*
 * Checks, custom metrics and requests can be tagged with any number of tags.
 *
 * Tags can be used for:
 * - Creating metric thresholds by filtering the metric data stream based on tags
 * - Aid result analysis by allowing for more precise filtering of metrics
 */

let myTrend = new Trend("my_trend");

export default function() ***REMOVED***
    // Add tag to request metric data
    let res = http.get("http://httpbin.org/", ***REMOVED*** tags: ***REMOVED*** my_tag: "I'm a tag" ***REMOVED*** ***REMOVED***);

    // Add tag to check
    check(res, ***REMOVED*** "status is 200": (r) => r.status === 200 ***REMOVED***, ***REMOVED*** my_tag: "I'm a tag" ***REMOVED***);

    // Add tag to custom metric
    myTrend.add(res.timings.connecting, ***REMOVED*** my_tag: "I'm a tag" ***REMOVED***);
***REMOVED***
