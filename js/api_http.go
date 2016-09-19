package js

import (
	// "github.com/robertkrimen/otto"
	"net/http"
	"strings"
)

type HTTPResponse struct ***REMOVED***
	Status int
***REMOVED***

func (a JSAPI) HTTPRequest(method, urlStr, body string, params map[string]interface***REMOVED******REMOVED***) map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	req, err := http.NewRequest(method, urlStr, strings.NewReader(body))
	if err != nil ***REMOVED***
		throw(a.vu.vm, err)
	***REMOVED***

	res, err := a.vu.HTTPClient.Do(req)
	if err != nil ***REMOVED***
		throw(a.vu.vm, err)
	***REMOVED***

	return map[string]interface***REMOVED******REMOVED******REMOVED***
		"status": res.StatusCode,
	***REMOVED***
***REMOVED***
