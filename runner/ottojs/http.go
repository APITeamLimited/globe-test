package ottojs

import (
	"github.com/loadimpact/speedboat/runner"
	"github.com/robertkrimen/otto"
	"github.com/valyala/fasthttp"
	"time"
)

func (vu *VUContext) HTTPGet(call otto.FunctionCall) otto.Value ***REMOVED***
	result := make(chan runner.Result, 1)
	go func() ***REMOVED***
		req := fasthttp.AcquireRequest()
		defer fasthttp.ReleaseRequest(req)

		res := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseResponse(res)

		url := call.Argument(0).String()
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

	return otto.UndefinedValue()
***REMOVED***
