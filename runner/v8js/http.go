package v8js

import (
	"github.com/loadimpact/speedboat/runner"
	"github.com/valyala/fasthttp"
	"time"
)

func (vu *VUContext) HTTPGet(url string) ***REMOVED***
	result := make(chan runner.Result, 1)
	go func() ***REMOVED***
		req := fasthttp.AcquireRequest()
		defer fasthttp.ReleaseRequest(req)

		res := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseResponse(res)

		req.SetRequestURI(url)

		startTime := time.Now()
		err := vu.r.Client.Do(req, res)
		duration := time.Since(startTime)

		result <- runner.Result***REMOVED***Error: err, Time: duration***REMOVED***
	***REMOVED***()

	select ***REMOVED***
	case <-vu.ctx.Done():
	case res := <-result:
		vu.ch <- res
	***REMOVED***
***REMOVED***
