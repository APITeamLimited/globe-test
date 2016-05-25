package js

import (
	"errors"
	"github.com/loadimpact/speedboat/runner"
	"github.com/valyala/fasthttp"
	"gopkg.in/olebedev/go-duktape.v2"
	"time"
)

type apiFunc func(r *Runner, c *duktape.Context, ch chan<- runner.Result) int

type apiHTTPArgs struct ***REMOVED***
	Report bool `json:"report"`
***REMOVED***

func apiHTTPGet(r *Runner, c *duktape.Context, ch chan<- runner.Result) int ***REMOVED***
	url := argString(c, 0)
	if url == "" ***REMOVED***
		ch <- runner.Result***REMOVED***Error: errors.New("Missing URL in http.get()")***REMOVED***
		return 0
	***REMOVED***
	args := apiHTTPArgs***REMOVED******REMOVED***
	if err := argJSON(c, 1, &args); err != nil ***REMOVED***
		ch <- runner.Result***REMOVED***Error: errors.New("Invalid arguments to http.get()")***REMOVED***
		return 0
	***REMOVED***

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	req.SetRequestURI(url)

	startTime := time.Now()
	err := r.Client.Do(req, res)
	duration := time.Since(startTime)

	if args.Report ***REMOVED***
		ch <- runner.Result***REMOVED***Error: err, Time: duration***REMOVED***
	***REMOVED***

	return 0
***REMOVED***
