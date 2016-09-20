package js

import (
	// "github.com/robertkrimen/otto"
	"io"
	"io/ioutil"
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

	if h, ok := params["headers"]; ok ***REMOVED***
		headers, ok := h.(map[string]interface***REMOVED******REMOVED***)
		if !ok ***REMOVED***
			panic(a.vu.vm.MakeTypeError("headers must be an object"))
		***REMOVED***
		for key, v := range headers ***REMOVED***
			value, ok := v.(string)
			if !ok ***REMOVED***
				panic(a.vu.vm.MakeTypeError("header values must be strings"))
			***REMOVED***
			req.Header.Set(key, value)
		***REMOVED***
	***REMOVED***

	res, err := a.vu.HTTPClient.Do(req)
	if err != nil ***REMOVED***
		throw(a.vu.vm, err)
	***REMOVED***

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil ***REMOVED***
		throw(a.vu.vm, err)
	***REMOVED***
	res.Body.Close()

	return map[string]interface***REMOVED******REMOVED******REMOVED***
		"status": res.StatusCode,
		"body":   string(resBody),
	***REMOVED***
***REMOVED***
