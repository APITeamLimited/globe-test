package simple

import (
	"github.com/loadimpact/speedboat/loadtest"
	"github.com/loadimpact/speedboat/runner"
	"github.com/valyala/fasthttp"
	"golang.org/x/net/context"
	"time"
)

type SimpleRunner struct ***REMOVED***
	Client *fasthttp.Client
***REMOVED***

func New() *SimpleRunner ***REMOVED***
	return &SimpleRunner***REMOVED***
		Client: &fasthttp.Client***REMOVED***
			MaxIdleConnDuration: time.Duration(0),
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (r *SimpleRunner) Run(ctx context.Context, t loadtest.LoadTest, id int64) <-chan runner.Result ***REMOVED***
	ch := make(chan runner.Result)

	go func() ***REMOVED***
		defer close(ch)

		result := make(chan runner.Result, 1)
		for ***REMOVED***
			go func() ***REMOVED***
				req := fasthttp.AcquireRequest()
				defer fasthttp.ReleaseRequest(req)

				res := fasthttp.AcquireResponse()
				defer fasthttp.ReleaseResponse(res)

				req.SetRequestURI(t.URL)

				startTime := time.Now()
				err := r.Client.Do(req, res)
				duration := time.Since(startTime)

				result <- runner.Result***REMOVED***Error: err, Time: duration***REMOVED***
			***REMOVED***()

			select ***REMOVED***
			case <-ctx.Done():
				return
			case res := <-result:
				ch <- res
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return ch
***REMOVED***
