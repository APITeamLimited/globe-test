package js

import (
	// "github.com/robertkrimen/otto"
	"io"
	"net/http"
	"strings"
)

type HTTPResponse struct ***REMOVED***
	Status int
***REMOVED***

func (a JSAPI) HTTPRequest(method, url, body string, params map[string]interface***REMOVED******REMOVED***) map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	bodyReader := io.Reader(nil)
	if body != "" ***REMOVED***
		bodyReader = strings.NewReader(body)
	***REMOVED***
	req, err := http.NewRequest(method, url, bodyReader)
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
