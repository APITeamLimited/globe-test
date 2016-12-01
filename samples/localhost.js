import ***REMOVED*** group, check ***REMOVED*** from "k6";
import http from "k6/http";

export let options = ***REMOVED***
	thresholds: ***REMOVED***
		http_req_duration: ["avg<=100"],
	***REMOVED***
***REMOVED***;

export default function() ***REMOVED***
	group("front page", function() ***REMOVED***
		check(http.get("http://localhost:8080/"), ***REMOVED***
			"status is 200": (res) => res.status === 200,
		***REMOVED***);
	***REMOVED***);
	group("stylesheet", function() ***REMOVED***
		check(http.get("http://localhost:8080/style.css"), ***REMOVED***
			"status is 200": (res) => res.status === 200,
		***REMOVED***);
	***REMOVED***);
	group("image", function() ***REMOVED***
		check(http.get("http://localhost:8080/teddy.jpg"), ***REMOVED***
			"status is 200": (res) => res.status === 200,
		***REMOVED***);
	***REMOVED***);
***REMOVED***
