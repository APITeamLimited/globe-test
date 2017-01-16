//
// *** v3_account_login.js ***
// This file is an example of how to write test cases for single REST API end points using k6
// It implements a combined functional and load test for the Load Impact REST API end point /account/login
//

import http from "k6/http";
import ***REMOVED*** group, sleep, check ***REMOVED*** from "k6";
import ***REMOVED*** myTrend, options, urlbase, thinktime1, thinktime2 ***REMOVED*** from "./common";

export ***REMOVED*** options ***REMOVED***;

// (Note that these credentials do not work, this script is not intended to actually be executed)
let username = "testuser@loadimpact.com";
let password = "testpassword";

// We export this function as other test cases might want to use it to authenticate
export function v3_account_login(username, password, debug) ***REMOVED***
        // First we login. We are not interested in performance metrics from these login transactions
        var url = urlbase + "/v3/account/login";
        var payload = ***REMOVED*** email: username, password: password ***REMOVED***;
        var res = http.post(url, JSON.stringify(payload), ***REMOVED*** headers: ***REMOVED*** "Content-Type": "application/json" ***REMOVED*** ***REMOVED***);
	if (typeof debug !== 'undefined')
        	console.log("Login: status=" + String(res.status) + "  Body=" + res.body);
        return res;
***REMOVED***;

// Exercise /login endpoint when this test case is executed
export default function() ***REMOVED***
	group("login", function() ***REMOVED***
		var res = v3_account_login(username, password);
		check(res, ***REMOVED***
			"status is 200": (res) => res.status === 200,
			"content-type is application/json": (res) => res.headers['Content-Type'] === "application/json",
			"login successful": (res) => JSON.parse(res.body).hasOwnProperty('token')
		***REMOVED***);
		myTrend.add(res.timings.duration);
		sleep(thinktime1);
	***REMOVED***);
***REMOVED***;
