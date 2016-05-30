var http = require('http');

function dump(x, name, indent) ***REMOVED***
  var ret = name + "[" + typeof x + "] ";
  if (indent === undefined) indent = "";
  if (typeof x === 'object') ***REMOVED***
    for (var prop in x) ***REMOVED***
      ret += ("\n" + dump(x[prop], name + "." + prop, indent + "   "));
    ***REMOVED***
    return indent + ret;
  ***REMOVED***
  return indent + ret + "= " + x;
***REMOVED***

var res = http.get('http://httpbin.org/get', ***REMOVED***'a': 1, 'b': 2***REMOVED***);
var jsonob = res.json();
print(dump(jsonob, "jsonob"));

