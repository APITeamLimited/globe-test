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

		// We can reuse the same request across multiple iterations; if something goes awry here,
		// we abort the test, since it normally means a user failure (like a malformed URL).
		req, err := http.NewRequest(http.MethodGet, r.URL, nil)
		if err != nil ***REMOVED***
			ch <- runner.Result***REMOVED***Error: err***REMOVED***
			return
		***REMOVED***
		req.Close = true

		// Close this channel to abort the request on the spot. The old, transport-based way of
		// doing this is deprecated, as it doesn't play nice with HTTP/2 requests.
		// cancelRequest := make(chan struct***REMOVED******REMOVED***)
		// req.Cancel = cancelRequest

		results := make(chan runner.Result, 1)
		for ***REMOVED***
			go func() ***REMOVED***
				startTime := time.Now()
				res, err := r.Client.Do(req)
				duration := time.Since(startTime)

				if err != nil ***REMOVED***
					results <- runner.Result***REMOVED***Error: err, Time: duration***REMOVED***
					return
				***REMOVED***
				res.Body.Close()

				results <- runner.Result***REMOVED***Time: duration***REMOVED***
			***REMOVED***()

			select ***REMOVED***
			case res := <-results:
				ch <- res
			case <-ctx.Done():
				// close(cancelRequest)
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return ch
***REMOVED***
