var http = require('http');
var res = http.get('http://httpbin.org/get', ***REMOVED***'a': 1, 'b': 2***REMOVED***);
print("URL: " + res.json().url);
