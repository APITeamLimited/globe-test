/**
 * @module k6/http
 */
import ***REMOVED*** parseHTML ***REMOVED*** from "k6/html";

export class Response ***REMOVED***
	/**
	 * Represents an HTTP response.
	 * @memberOf module:k6/http
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
 * @param  ***REMOVED***string|Object***REMOVED*** body Request body; objects will be query encoded.
 * @param  ***REMOVED***Object***REMOVED*** params      Additional parameters.
 * @return ***REMOVED***module:k6/http.Response***REMOVED***
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
	***REMOVED*** else ***REMOVED***
		body = ''
	***REMOVED***
	return new Response(__jsapi__.HTTPRequest(method, url, body, JSON.stringify(params)));
***REMOVED***;

/**
 * Makes a GET request.
 * @see    module:k6/http.request
 * @param  ***REMOVED***string***REMOVED*** url         Request URL (eg. "http://example.com/")
 * @param  ***REMOVED***Object***REMOVED*** params      Additional parameters.
 * @return ***REMOVED***module:k6/http.Response***REMOVED***
 */
export function get(url, params) ***REMOVED***
	return request("GET", url, null, params);
***REMOVED***;

/**
 * Makes a POST request.
 * @see    module:k6/http.request
 * @param  ***REMOVED***string***REMOVED*** url         Request URL (eg. "http://example.com/")
 * @param  ***REMOVED***string|Object***REMOVED*** body Request body; objects will be query encoded.
 * @param  ***REMOVED***Object***REMOVED*** params      Additional parameters.
 * @return ***REMOVED***module:k6/http.Response***REMOVED***
 */
export function post(url, body, params) ***REMOVED***
	return request("POST", url, body, params);
***REMOVED***;

/**
 * Makes a PUT request.
 * @see    module:k6/http.request
 * @param  ***REMOVED***string***REMOVED*** url         Request URL (eg. "http://example.com/")
 * @param  ***REMOVED***string|Object***REMOVED*** body Request body; objects will be query encoded.
 * @param  ***REMOVED***Object***REMOVED*** params      Additional parameters.
 * @return ***REMOVED***module:k6/http.Response***REMOVED***
 */
export function put(url, body, params) ***REMOVED***
	return request("PUT", url, body, params);
***REMOVED***;

/**
 * Makes a DELETE request.
 * @see    module:k6/http.request
 * @param  ***REMOVED***string***REMOVED*** url         Request URL (eg. "http://example.com/")
 * @param  ***REMOVED***string|Object***REMOVED*** body Request body; objects will be query encoded.
 * @param  ***REMOVED***Object***REMOVED*** params      Additional parameters.
 * @return ***REMOVED***module:k6/http.Response***REMOVED***
 */
export function del(url, body, params) ***REMOVED***
	return request("DELETE", url, body, params);
***REMOVED***;

/**
 * Makes a PATCH request.
 * @see    module:k6/http.request
 * @param  ***REMOVED***string***REMOVED*** url         Request URL (eg. "http://example.com/")
 * @param  ***REMOVED***string|Object***REMOVED*** body Request body; objects will be query encoded.
 * @param  ***REMOVED***Object***REMOVED*** params      Additional parameters.
 * @return ***REMOVED***module:k6/http.Response***REMOVED***
 */
export function patch(url, body, params) ***REMOVED***
	return request("PATCH", url, body, params);
***REMOVED***;

export default ***REMOVED***
	Response: Response,
	request: request,
	get: get,
	post: post,
	put: put,
	del: del,
	patch: patch,
***REMOVED***;
