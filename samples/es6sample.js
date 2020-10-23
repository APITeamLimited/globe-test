import ***REMOVED*** group, check, sleep ***REMOVED*** from "k6";
import ***REMOVED*** Counter, Rate ***REMOVED*** from "k6/metrics";
import http from "k6/http";

export let options = ***REMOVED***
	vus: 5,
	thresholds: ***REMOVED***
		my_rate: ["rate>=0.4"], // Require my_rate's success rate to be >=40%
		http_req_duration: ["avg<1000"], // Require http_req_duration's average to be <1000ms
	***REMOVED***
***REMOVED***;

let mCounter = new Counter("my_counter");
let mRate = new Rate("my_rate");

export default function() ***REMOVED***
	check(Math.random(), ***REMOVED***
		"top-level test": (v) => v < 1/3
	***REMOVED***);
	group("my group", function() ***REMOVED***
		mCounter.add(1, ***REMOVED*** tag: "test" ***REMOVED***);

		check(Math.random(), ***REMOVED***
			"random value is < 0.5": (v) => mRate.add(v < 0.5),
		***REMOVED***);

		group("json", function() ***REMOVED***
			let res = http.get("https://httpbin.org/get", ***REMOVED***
				headers: ***REMOVED*** "X-Test": "abc123" ***REMOVED***,
			***REMOVED***);

			check(res, ***REMOVED***
				"status is 200": (res) => res.status === 200,
				"X-Test header is correct": (res) => res.json().headers['X-Test'] === "abc123",
			***REMOVED***);
		***REMOVED***);

		group("html", function() ***REMOVED***
			check(http.get("http://test.k6.io/"), ***REMOVED***
				"status is 200": (res) => res.status === 200,
				"content type is html": (res) => res.headers['Content-Type'].startsWith("text/html"),
				"welcome message is correct": (res) => res.html("p.description").text() === "Collection of simple web-pages suitable for load testing.",
			***REMOVED***);
		***REMOVED***);

		group("nested", function() ***REMOVED***
			check(null, ***REMOVED***
				"always passes": true,
				"always fails": false,
			***REMOVED***);
		***REMOVED***);
	***REMOVED***);
	sleep(10 * Math.random());
***REMOVED***;
