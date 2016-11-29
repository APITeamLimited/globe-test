import ***REMOVED*** group, check, sleep ***REMOVED*** from "speedboat";
import ***REMOVED*** Counter, Rate ***REMOVED*** from "speedboat/metrics";
import http from "speedboat/http";

export let options = ***REMOVED***
	vus: 5,
	thresholds: ***REMOVED***
		my_rate: ["rate>=0.4"],
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
			let res = http.get("https://httpbin.org/get", null, ***REMOVED***
				headers: ***REMOVED*** "X-Test": "abc123" ***REMOVED***,
			***REMOVED***);
			check(res, ***REMOVED***
				"status is 200": (res) => res.status === 200,
				"X-Test header is correct": (res) => res.json().headers['X-Test'] === "abc123",
			***REMOVED***);
			// console.log(res.body);
		***REMOVED***);

		group("html", function() ***REMOVED***
			check(http.get("http://test.loadimpact.com/"), ***REMOVED***
				"status is 200": (res) => res.status === 200,
				"content type is html": (res) => res.headers['Content-Type'] === "text/html",
				"welcome message is correct": (res) => res.html("h2").text() === "Welcome to the LoadImpact.com demo site!",
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
