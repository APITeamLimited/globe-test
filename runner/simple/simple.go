package simple

import (
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

func (r *SimpleRunner) Run(ctx context.Context) <-chan time.Duration ***REMOVED***
	ch := make(chan time.Duration)

	go func() ***REMOVED***
		defer close(ch)
		for ***REMOVED***
			startTime := time.Now()
			res, err := r.Client.Get(r.URL)
			duration := time.Since(startTime)
			if err != nil ***REMOVED***
				panic(err)
			***REMOVED***
			res.Body.Close()

			select ***REMOVED***
			case <-ctx.Done():
				return
			default:
				ch <- duration
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return ch
***REMOVED***
