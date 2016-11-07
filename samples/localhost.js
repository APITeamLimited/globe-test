import ***REMOVED*** check ***REMOVED*** from "speedboat";
import http from "speedboat/http";

export default function() ***REMOVED***
	check(http.get("http://localhost:8080/"), ***REMOVED***
		"status is 200": (v) => v.status === 200,
	***REMOVED***);
***REMOVED***
