import ***REMOVED*** group, check ***REMOVED*** from "speedboat";
import http from "speedboat/http";

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
