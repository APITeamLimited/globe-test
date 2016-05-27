package js

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	neturl "net/url"
	"time"
)

type httpArgs struct ***REMOVED***
	Quiet   bool              `json:"quiet"`
	Headers map[string]string `json:"headers"`
***REMOVED***

type httpResponse struct ***REMOVED***
	Status  int               `json:"status"`
	Body    string            `json:"body"`
	Headers map[string]string `json:"headers"`
***REMOVED***

func httpDo(c *fasthttp.Client, method, url, body string, args httpArgs) (httpResponse, time.Duration, error) ***REMOVED***
	if method == "GET" ***REMOVED***
		u, err := neturl.Parse(url)
		if err != nil ***REMOVED***
			return httpResponse***REMOVED******REMOVED***, 0, err
		***REMOVED***

		var params map[string]interface***REMOVED******REMOVED***
		if err = json.Unmarshal([]byte(body), &params); err != nil ***REMOVED***
			return httpResponse***REMOVED******REMOVED***, 0, err
		***REMOVED***

		for key, val := range params ***REMOVED***
			u.Query().Set(key, fmt.Sprint(val))
		***REMOVED***
		url = u.String()
	***REMOVED***

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod(method)
	req.SetRequestURI(url)

	if method != "GET" ***REMOVED***
		req.SetBodyString(body)
	***REMOVED***

	for key, value := range args.Headers ***REMOVED***
		req.Header.Set(key, value)
	***REMOVED***

	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	startTime := time.Now()
	err := c.Do(req, res)
	duration := time.Since(startTime)

	if err != nil ***REMOVED***
		return httpResponse***REMOVED******REMOVED***, duration, err
	***REMOVED***

	resHeaders := make(map[string]string)
	res.Header.VisitAll(func(key, value []byte) ***REMOVED***
		resHeaders[string(key)] = string(value)
	***REMOVED***)

	return httpResponse***REMOVED***
		Status:  res.StatusCode(),
		Body:    string(res.Body()),
		Headers: resHeaders,
	***REMOVED***, duration, nil
***REMOVED***
