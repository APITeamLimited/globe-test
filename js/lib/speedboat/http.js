/**
 * @module speedboat/http
 */
import ***REMOVED*** parseHTML ***REMOVED*** from "speedboat/html";

export class Response ***REMOVED***
	/**
	 * Represents an HTTP response.
	 * @memberOf module:speedboat/http
	 */
	constructor(data) ***REMOVED***
		Object.assign(this, data);
	***REMOVED***

	json() ***REMOVED***
		if (!this._json) ***REMOVED***
			this._json = JSON.parse(this.body);
		***REMOVED***
		return this._json;
	***REMOVED***

	html(sel) ***REMOVED***
		if (!this._html) ***REMOVED***
			this._html = parseHTML(this.body);
		***REMOVED***
		if (sel) ***REMOVED***
			return this._html.find(sel);
		***REMOVED***
		return this._html;
	***REMOVED***
***REMOVED***

/**
 * Makes an HTTP request.
 * @param  ***REMOVED***string***REMOVED*** method      HTTP Method (eg. "GET")
 * @param  ***REMOVED***string***REMOVED*** url         Request URL (eg. "http://example.com/")
 * @param  ***REMOVED***string|Object***REMOVED*** body Request body (query for GET/HEAD); objects will be query encoded.
 * @param  ***REMOVED***Object***REMOVED*** params      Additional parameters.
 * @return ***REMOVED***module:speedboat/http.Response***REMOVED***
 */
export function request(method, url, body, params = ***REMOVED******REMOVED***) ***REMOVED***
	method = method.toUpperCase();
	if (body) ***REMOVED***
		if (typeof body === "object") ***REMOVED***
			let formstring = "";
			for (let key in body) ***REMOVED***
				if (formstring !== "") ***REMOVED***
					formstring += "&";
				***REMOVED***
				formstring += key + "=" + encodeURIComponent(body[key]);
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

/**
 * Makes a GET request.
 * @see    module:speedboat/http.request
 * @param  ***REMOVED***string***REMOVED*** url         Request URL (eg. "http://example.com/")
 * @param  ***REMOVED***string|Object***REMOVED*** body Request body (query for GET/HEAD); objects will be query encoded.
 * @param  ***REMOVED***Object***REMOVED*** params      Additional parameters.
 * @return ***REMOVED***module:speedboat/http.Response***REMOVED***
 */
export function get(url, body, params) ***REMOVED***
	return request("GET", url, body, params);
***REMOVED***;

/**
 * Makes a POST request.
 * @see    module:speedboat/http.request
 * @param  ***REMOVED***string***REMOVED*** url         Request URL (eg. "http://example.com/")
 * @param  ***REMOVED***string|Object***REMOVED*** body Request body (query for GET/HEAD); objects will be query encoded.
 * @param  ***REMOVED***Object***REMOVED*** params      Additional parameters.
 * @return ***REMOVED***module:speedboat/http.Response***REMOVED***
 */
export function post(url, body, params) ***REMOVED***
	return request("POST", url, body, params);
***REMOVED***;

/**
 * Makes a PUT request.
 * @see    module:speedboat/http.request
 * @param  ***REMOVED***string***REMOVED*** url         Request URL (eg. "http://example.com/")
 * @param  ***REMOVED***string|Object***REMOVED*** body Request body (query for GET/HEAD); objects will be query encoded.
 * @param  ***REMOVED***Object***REMOVED*** params      Additional parameters.
 * @return ***REMOVED***module:speedboat/http.Response***REMOVED***
 */
export function put(url, body, params) ***REMOVED***
	return request("PUT", url, body, params);
***REMOVED***;

/**
 * Makes a DELETE request.
 * @see    module:speedboat/http.request
 * @param  ***REMOVED***string***REMOVED*** url         Request URL (eg. "http://example.com/")
 * @param  ***REMOVED***string|Object***REMOVED*** body Request body (query for GET/HEAD); objects will be query encoded.
 * @param  ***REMOVED***Object***REMOVED*** params      Additional parameters.
 * @return ***REMOVED***module:speedboat/http.Response***REMOVED***
 */
export function del(url, body, params) ***REMOVED***
	return request("DELETE", url, body, params);
***REMOVED***;

/**
 * Makes a PATCH request.
 * @see    module:speedboat/http.request
 * @param  ***REMOVED***string***REMOVED*** url         Request URL (eg. "http://example.com/")
 * @param  ***REMOVED***string|Object***REMOVED*** body Request body (query for GET/HEAD); objects will be query encoded.
 * @param  ***REMOVED***Object***REMOVED*** params      Additional parameters.
 * @return ***REMOVED***module:speedboat/http.Response***REMOVED***
 */
export function patch(url, body, params) ***REMOVED***
	return request("PATCH", url, body, params);
***REMOVED***;

/**
 * Sets the maximum number of redirects to follow. A request that encounters more than this many
 * redirects will error. Default: 10.
 * @param ***REMOVED***Number***REMOVED*** n Max number of redirects.
 */
export function setMaxRedirects(n) ***REMOVED***
	__jsapi__.HTTPSetMaxRedirects(n);
***REMOVED***

export default ***REMOVED***
	Response: Response,
	request: request,
	get: get,
	post: post,
	put: put,
	del: del,
	patch: patch,
	setMaxRedirects: setMaxRedirects,
***REMOVED***;
