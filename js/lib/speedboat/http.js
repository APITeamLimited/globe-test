export class Response ***REMOVED***
	constructor(data) ***REMOVED***
		Object.assign(this, data);
	***REMOVED***

	json() ***REMOVED***
		if (!this._json) ***REMOVED***
			this._json = JSON.parse(this.body);
		***REMOVED***
		return this._json;
	***REMOVED***
***REMOVED***

export function request(method, url, body, params = ***REMOVED******REMOVED***) ***REMOVED***
	method = method.toUpperCase();
	if (body) ***REMOVED***
		if (typeof body === "object") ***REMOVED***
			let formstring = "";
			for (let entry of body) ***REMOVED***
				if (formstring !== "") ***REMOVED***
					formstring += "&";
				***REMOVED***
				formstring += entry[0] + "=" + encodeURIComponent(entry[1]);
			***REMOVED***
			body = formstring;
		***REMOVED***
		if (method === "GET" || method === "HEAD") ***REMOVED***
			url += (url.includes("?") ? "&" : "?") + body;
			body = "";
		***REMOVED***
	***REMOVED***
	return new Response(__jsapi__.HTTPRequest(method, url, body, params));
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

export function del(url, body, params) ***REMOVED***
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
	del: del,
	patch: patch,
***REMOVED***;
