package js

import (
	"errors"
	"github.com/loadimpact/speedboat/runner"
	"github.com/valyala/fasthttp"
	"gopkg.in/olebedev/go-duktape.v2"
	"time"
)

type apiFunc func(r *Runner, c *duktape.Context, ch chan<- runner.Result) int

func apiHTTPDo(r *Runner, c *duktape.Context, ch chan<- runner.Result) int ***REMOVED***
	method := argString(c, 0)
	if method == "" ***REMOVED***
		ch <- runner.Result***REMOVED***Error: errors.New("Missing method in http call")***REMOVED***
		return 0
	***REMOVED***

	url := argString(c, 1)
	if url == "" ***REMOVED***
		ch <- runner.Result***REMOVED***Error: errors.New("Missing URL in http call")***REMOVED***
		return 0
	***REMOVED***

	args := struct ***REMOVED***
		Report  bool              `json:"report"`
		Headers map[string]string `json:"headers"`
	***REMOVED******REMOVED******REMOVED***
	if err := argJSON(c, 2, &args); err != nil ***REMOVED***
		ch <- runner.Result***REMOVED***Error: errors.New("Invalid arguments to http call")***REMOVED***
		return 0
	***REMOVED***

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	req.Header.SetMethod(method)
	req.SetRequestURI(url)

	for key, value := range args.Headers ***REMOVED***
		req.Header.Set(key, value)
	***REMOVED***

	startTime := time.Now()
	err := r.Client.Do(req, res)
	duration := time.Since(startTime)

	if args.Report ***REMOVED***
		ch <- runner.Result***REMOVED***Error: err, Time: duration***REMOVED***
	***REMOVED***

	c.PushObject()
	***REMOVED***
		c.PushNumber(float64(res.StatusCode()))
		c.PutPropString(-2, "status")

		c.PushString(string(res.Body()))
		c.PutPropString(-2, "body")

		c.PushObject()
		res.Header.VisitAll(func(key, value []byte) ***REMOVED***
			c.PushString(string(value))
			c.PutPropString(-2, string(key))
		***REMOVED***)
		c.PutPropString(-2, "headers")
	***REMOVED***

	return 1
***REMOVED***

func apiHTTPSetMaxConnectionsPerHost(r *Runner, c *duktape.Context, ch chan<- runner.Result) int ***REMOVED***
	num := int(argNumber(c, 0))
	if num < 1 ***REMOVED***
		ch <- runner.Result***REMOVED***Error: errors.New("Max connections per host must be at least 1")***REMOVED***
	***REMOVED***
	r.Client.MaxConnsPerHost = num
	return 0
***REMOVED***
