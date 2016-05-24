package js

import (
	"errors"
	"github.com/loadimpact/speedboat/runner"
	"github.com/valyala/fasthttp"
	"gopkg.in/olebedev/go-duktape.v2"
	"time"
)

type apiFunc func(r *Runner, c *duktape.Context, ch chan<- runner.Result) int

func apiHTTPGet(r *Runner, c *duktape.Context, ch chan<- runner.Result) int ***REMOVED***
	url := argString(c, 0)
	if url == "" ***REMOVED***
		ch <- runner.Result***REMOVED***Error: errors.New("Missing URL in http.get()")***REMOVED***
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

	ch <- runner.Result***REMOVED***Error: err, Time: duration***REMOVED***

	return 0
***REMOVED***
