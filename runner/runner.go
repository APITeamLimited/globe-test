package runner

import (
	"sync"
	"time"
)

// A single metric for a test execution.
type Metric struct ***REMOVED***
	Time     time.Time
	Duration time.Duration
***REMOVED***

// A user-printed log message.
type LogEntry struct ***REMOVED***
	Time time.Time
	Text string
***REMOVED***

type Runner interface ***REMOVED***
	Load(filename, src string) error
	RunVU() <-chan interface***REMOVED******REMOVED***
***REMOVED***

func NewError(err error) interface***REMOVED******REMOVED*** ***REMOVED***
	return err
***REMOVED***

func NewLogEntry(entry LogEntry) interface***REMOVED******REMOVED*** ***REMOVED***
	return entry
***REMOVED***

func NewMetric(metric Metric) interface***REMOVED******REMOVED*** ***REMOVED***
	return metric
***REMOVED***

func Run(r Runner, vus int) <-chan interface***REMOVED******REMOVED*** ***REMOVED***
	ch := make(chan interface***REMOVED******REMOVED***)

	go func() ***REMOVED***
		wg := sync.WaitGroup***REMOVED******REMOVED***
		for i := 0; i < vus; i++ ***REMOVED***
			wg.Add(1)
			go func() ***REMOVED***
				defer wg.Done()
				for res := range r.RunVU() ***REMOVED***
					ch <- res
				***REMOVED***
			***REMOVED***()
		***REMOVED***

		go func() ***REMOVED***
			wg.Wait()
			close(ch)
		***REMOVED***()
	***REMOVED***()

	return ch
***REMOVED***
