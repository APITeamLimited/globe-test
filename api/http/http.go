package http

import (
	"github.com/loadimpact/speedboat/runner"
	"github.com/valyala/fasthttp"
	"math"
	"time"
)

type context struct ***REMOVED***
	client *fasthttp.Client
***REMOVED***

func New() map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	ctx := &context***REMOVED***
		client: &fasthttp.Client***REMOVED***
			Dial:                fasthttp.Dial,
			MaxIdleConnDuration: time.Duration(0),
			MaxConnsPerHost:     math.MaxInt64,
		***REMOVED***,
	***REMOVED***
	return map[string]interface***REMOVED******REMOVED******REMOVED***
		"get": ctx.Get,
	***REMOVED***
***REMOVED***

func (ctx *context) Get(url string) <-chan runner.Result ***REMOVED***
	ch := make(chan runner.Result, 1)
	go func() ***REMOVED***
		defer close(ch)

		req := fasthttp.AcquireRequest()
		defer fasthttp.ReleaseRequest(req)

		res := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseResponse(res)

		req.SetRequestURI(url)

		startTime := time.Now()
		err := ctx.client.Do(req, res)
		duration := time.Since(startTime)

		ch <- runner.Result***REMOVED***Error: err, Time: duration***REMOVED***
	***REMOVED***()
	return ch
***REMOVED***
