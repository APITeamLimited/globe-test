import ***REMOVED*** group, test, sleep ***REMOVED*** from "speedboat";
import http from "speedboat/http";

export let options = ***REMOVED***
	vus: 5,
***REMOVED***;

export default function() ***REMOVED***
	test(Math.random(), ***REMOVED***
		"top-level test": (v) => v < 1/3
	***REMOVED***);
	group("my group", function() ***REMOVED***
		test(Math.random(), ***REMOVED***
			"random value is < 0.5": (v) => v < 0.5
		***REMOVED***);

		group("http", function() ***REMOVED***
			let res = http.get("https://httpbin.org/get", null, ***REMOVED***
				headers: ***REMOVED*** "X-Test": "abc123" ***REMOVED***,
			***REMOVED***);
			test(res, ***REMOVED***
				"status is 200": (res) => res.status === 200,
				"X-Test header is correct": (res) => res.json().headers['X-Test'] === "abc123",
			***REMOVED***);
			// console.log(res.body);
		***REMOVED***);

		group("nested", function() ***REMOVED***
			test(null, ***REMOVED***
				"always passes": true,
				"always fails": false,
			***REMOVED***);
		***REMOVED***);
	***REMOVED***);
	sleep(10 * Math.random());
***REMOVED***;
