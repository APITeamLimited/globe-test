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
			req, err := http.NewRequest(http.MethodGet, r.URL, nil)
			if err != nil ***REMOVED***
				ch <- runner.Result***REMOVED***Error: err***REMOVED***
				continue
			***REMOVED***
			req.Close = true

			cancel := make(chan struct***REMOVED******REMOVED***)
			req.Cancel = cancel

			go func() ***REMOVED***
				startTime := time.Now()
				res, err := r.Client.Do(req)
				duration := time.Since(startTime)

				if err != nil ***REMOVED***
					ch <- runner.Result***REMOVED***Error: err, Time: duration***REMOVED***
					return
				***REMOVED***
				res.Body.Close()

				ch <- runner.Result***REMOVED***Time: duration***REMOVED***
			***REMOVED***()

			_, keepGoing := <-ctx.Done()
			if !keepGoing ***REMOVED***
				close(cancel)
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return ch
***REMOVED***
