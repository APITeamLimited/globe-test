package simple

import (
	"github.com/loadimpact/speedboat/runner"
	"golang.org/x/net/context"
	"net/http"
	"time"
)

type SimpleRunner struct ***REMOVED***
	URL    string
	Client *http.Client
***REMOVED***

func New() *SimpleRunner ***REMOVED***
	return &SimpleRunner***REMOVED***
		Client: &http.Client***REMOVED***
			Transport: &http.Transport***REMOVED***
				DisableKeepAlives: true,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (r *SimpleRunner) Run(ctx context.Context) <-chan runner.Result ***REMOVED***
	ch := make(chan runner.Result)

	go func() ***REMOVED***
		defer close(ch)

		for ***REMOVED***
			// Note that we abort if we cannot create a request. This means we're either out of
			// memory, or we have invalid user input, neither of which are recoverable.
			req, err := http.NewRequest("GET", r.URL, nil)
			if err != nil ***REMOVED***
				ch <- runner.Result***REMOVED***Error: err***REMOVED***
				return
			***REMOVED***
			req.Close = true
			req.Cancel = ctx.Done()

			startTime := time.Now()
			res, err := r.Client.Do(req)
			duration := time.Since(startTime)

			select ***REMOVED***
			case <-ctx.Done():
				return
			default:
				if err != nil ***REMOVED***
					ch <- runner.Result***REMOVED***Error: err, Time: duration***REMOVED***
					continue
				***REMOVED***
				res.Body.Close()
				ch <- runner.Result***REMOVED***Time: duration***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return ch
***REMOVED***
