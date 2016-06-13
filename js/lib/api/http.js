"use strict";

__modules__.http = __modules__.http || ***REMOVED******REMOVED***;

__modules__.http.get = function() ***REMOVED***
	return __modules__.http.request.apply(this, _.concat(['GET'], arguments));
***REMOVED***

__modules__.http.post = function() ***REMOVED***
	return __modules__.http.request.apply(this, _.concat(['POST'], arguments));
***REMOVED***

__modules__.http.put = function() ***REMOVED***
	return __modules__.http.request.apply(this, _.concat(['PUT'], arguments));
***REMOVED***

__modules__.http.delete = function() ***REMOVED***
	return __modules__.http.request.apply(this, _.concat(['DELETE'], arguments));
***REMOVED***

__modules__.http.patch = function() ***REMOVED***
	return __modules__.http.request.apply(this, _.concat(['PATCH'], arguments));
***REMOVED***

__modules__.http.options = function() ***REMOVED***
	return __modules__.http.request.apply(this, _.concat(['OPTIONS'], arguments));
***REMOVED***
