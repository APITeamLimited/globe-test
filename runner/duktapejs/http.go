package duktapejs

import (
	"github.com/loadimpact/speedboat/runner"
	"github.com/valyala/fasthttp"
	"gopkg.in/olebedev/go-duktape.v2"
	"time"
)

func (vu *VUContext) HTTPGet(c *duktape.Context) int ***REMOVED***
	result := make(chan runner.Result, 1)
	go func() ***REMOVED***
		req := fasthttp.AcquireRequest()
		defer fasthttp.ReleaseRequest(req)

		res := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseResponse(res)

		url := c.SafeToString(-1)
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

	return 0
***REMOVED***
