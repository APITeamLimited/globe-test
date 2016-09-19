export function request(method, url, body, params = ***REMOVED******REMOVED***) ***REMOVED***
	if (typeof body === "object") ***REMOVED***
		let newbody = "";
		for (let entry of body) ***REMOVED***
			if (newbody !== "") ***REMOVED***
				newbody += "&";
			***REMOVED***
			newbody += entry[0] + "=" + encodeURIComponent(entry[1]);
		***REMOVED***
		body = newbody;
	***REMOVED***
	return __jsapi__.HTTPRequest(method, url, body, params);
***REMOVED***;

export function get(url, body, params) ***REMOVED***
	return request("GET", url, body, params);
***REMOVED***;

export function post(url, body, params) ***REMOVED***
	return request("POST", url, body, params);
***REMOVED***;

export function put(url, body, params) ***REMOVED***
	return request("PUT", url, body, params);
***REMOVED***;

export function delete_(url, body, params) ***REMOVED***
	return request("DELETE", url, body, params);
***REMOVED***;

export function patch(url, body, params) ***REMOVED***
	return request("PATCH", url, body, params);
***REMOVED***;

export default ***REMOVED***
	request: request,
	get: get,
	post: post,
	put: put,
	delete_: delete_,
	patch: patch,
***REMOVED***;
