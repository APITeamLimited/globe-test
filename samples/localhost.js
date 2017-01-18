import ***REMOVED*** group, check ***REMOVED*** from "k6";
import http from "k6/http";

export let options = ***REMOVED***
	thresholds: ***REMOVED***
		'http_req_duration***REMOVED***kind:html***REMOVED***': ["avg<=10*ms"],
		'http_req_duration***REMOVED***kind:css***REMOVED***': ["avg<=10*ms"],
		'http_req_duration***REMOVED***kind:img***REMOVED***': ["avg<=100*ms"],
	***REMOVED***
***REMOVED***;

export default function() ***REMOVED***
	group("front page", function() ***REMOVED***
		check(http.get("http://localhost:8080/", null, ***REMOVED***
			tags: ***REMOVED***'kind': 'html' ***REMOVED***,
		***REMOVED***), ***REMOVED***
			"status is 200": (res) => res.status === 200,
		***REMOVED***);
	***REMOVED***);
	group("stylesheet", function() ***REMOVED***
		check(http.get("http://localhost:8080/style.css", null, ***REMOVED***
			tags: ***REMOVED***'kind': 'css' ***REMOVED***,
		***REMOVED***), ***REMOVED***
			"status is 200": (res) => res.status === 200,
		***REMOVED***);
	***REMOVED***);
	group("image", function() ***REMOVED***
		check(http.get("http://localhost:8080/teddy.jpg", null, ***REMOVED***
			tags: ***REMOVED***'kind': 'img' ***REMOVED***,
		***REMOVED***), ***REMOVED***
			"status is 200": (res) => res.status === 200,
		***REMOVED***);
	***REMOVED***);
***REMOVED***
