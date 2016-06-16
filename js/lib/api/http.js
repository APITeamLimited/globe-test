"use strict";

$http = $http || ***REMOVED******REMOVED***;

$http.get = function() ***REMOVED***
	return $http.request.apply(this, _.concat(['GET'], arguments));
***REMOVED***

$http.post = function() ***REMOVED***
	return $http.request.apply(this, _.concat(['POST'], arguments));
***REMOVED***

$http.put = function() ***REMOVED***
	return $http.request.apply(this, _.concat(['PUT'], arguments));
***REMOVED***

$http.delete = function() ***REMOVED***
	return $http.request.apply(this, _.concat(['DELETE'], arguments));
***REMOVED***

$http.patch = function() ***REMOVED***
	return $http.request.apply(this, _.concat(['PATCH'], arguments));
***REMOVED***

$http.options = function() ***REMOVED***
	return $http.request.apply(this, _.concat(['OPTIONS'], arguments));
***REMOVED***
