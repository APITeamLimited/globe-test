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
			test(http.get("http://localhost:8080"), ***REMOVED***
				"status is 200": (res) => res.status === 200,
			***REMOVED***);
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
