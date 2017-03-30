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

function parseBody(body) ***REMOVED***
	if (body) ***REMOVED***
		if (typeof body === "object") ***REMOVED***
			let formstring = "";
			for (let key in body) ***REMOVED***
				if (formstring !== "") ***REMOVED***
					formstring += "&";
				***REMOVED***
				if (Array.isArray(body[key])) ***REMOVED***
					let l = body[key].length;
					for (let i = 0; i < l; i++) ***REMOVED***
						formstring += key + "=" + encodeURIComponent(body[key][i]);
						if (formstring !== "") ***REMOVED***
							formstring += "&";
						***REMOVED***
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					formstring += key + "=" + encodeURIComponent(body[key]);
				***REMOVED***
			***REMOVED***
			return formstring;
		***REMOVED***
		return body;
	***REMOVED*** else ***REMOVED***
		return '';
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
	if (typeof body === "object") ***REMOVED***
		if (typeof params["headers"] !== "object") ***REMOVED***
			params["headers"] = ***REMOVED******REMOVED***;
		***REMOVED***
		params["headers"]["Content-Type"] = "application/x-www-form-urlencoded";
	***REMOVED***
	body = parseBody(body);
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

/**
 * Makes a CONNECT request.
 * @see    module:k6/http.request
 * @param  ***REMOVED***string***REMOVED*** url         Request URL (eg. "http://example.com/")
 * @param  ***REMOVED***string|Object***REMOVED*** body Request body; objects will be query encoded.
 * @param  ***REMOVED***Object***REMOVED*** params      Additional parameters.
 * @return ***REMOVED***module:k6/http.Response***REMOVED***
 */
export function connect(url, body, params) ***REMOVED***
	return request("CONNECT", url, body, params);
***REMOVED***;

/**
 * Makes a OPTIONS request.
 * @see    module:k6/http.request
 * @param  ***REMOVED***string***REMOVED*** url         Request URL (eg. "http://example.com/")
 * @param  ***REMOVED***string|Object***REMOVED*** body Request body; objects will be query encoded.
 * @param  ***REMOVED***Object***REMOVED*** params      Additional parameters.
 * @return ***REMOVED***module:k6/http.Response***REMOVED***
 */
export function options(url, body, params) ***REMOVED***
	return request("OPTIONS", url, body, params);
***REMOVED***;

/**
 * Makes a TRACE request.
 * @see    module:k6/http.request
 * @param  ***REMOVED***string***REMOVED*** url         Request URL (eg. "http://example.com/")
 * @param  ***REMOVED***string|Object***REMOVED*** body Request body; objects will be query encoded.
 * @param  ***REMOVED***Object***REMOVED*** params      Additional parameters.
 * @return ***REMOVED***module:k6/http.Response***REMOVED***
 */
export function trace(url, body, params) ***REMOVED***
	return request("TRACE", url, body, params);
***REMOVED***;

/**
 * Batches multiple requests together.
 * @see    module:k6/http.request
 * @param  ***REMOVED***Array|Object***REMOVED*** requests	An array or object of requests, in string or object form.
 * @return ***REMOVED***Array.<module:k6/http.Response>|Object***REMOVED***
 */
export function batch(requests) ***REMOVED***
	function stringToObject(str) ***REMOVED***
		return ***REMOVED***
			"method": "GET",
			"url": str,
			"body": null,
			"params": JSON.stringify(***REMOVED******REMOVED***)
		***REMOVED***
	***REMOVED***

	function formatObject(obj) ***REMOVED***
		obj.params = !obj.params ? ***REMOVED******REMOVED*** :obj.params
		obj.body = parseBody(obj.body)
		obj.params = JSON.stringify(obj.params)
		return obj
	***REMOVED***

	let result
	if (requests.length > 0) ***REMOVED***
		result = requests.map(e => ***REMOVED***
			if (typeof e === 'string') ***REMOVED***
				return stringToObject(e)
			***REMOVED*** else ***REMOVED***
				return formatObject(e)
			***REMOVED***
		***REMOVED***)
	***REMOVED*** else ***REMOVED***
		result = ***REMOVED******REMOVED***
		Object.keys(requests).map(e => ***REMOVED***
			let val = requests[e]
			if (typeof val === 'string') ***REMOVED***
				result[e] = stringToObject(val)
			***REMOVED*** else ***REMOVED***
				result[e] = formatObject(val)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
	
	let response = __jsapi__.BatchHTTPRequest(result);
	return response
***REMOVED***;

export default ***REMOVED***
	Response,
	request,
	get,
	post,
	put,
	del,
	patch,
	batch,
***REMOVED***;
