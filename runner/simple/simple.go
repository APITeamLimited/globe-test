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
			startTime := time.Now()
			res, err := r.Client.Get(r.URL)
			duration := time.Since(startTime)
			if err != nil ***REMOVED***
				ch <- runner.Result***REMOVED***Error: err***REMOVED***
			***REMOVED***
			res.Body.Close()

			select ***REMOVED***
			case <-ctx.Done():
				return
			default:
				ch <- runner.Result***REMOVED***Time: duration***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return ch
***REMOVED***
