package runner

import (
	"sync"
	"time"
)

// A single metric for a test execution.
type Metric struct ***REMOVED***
	Start    time.Time
	Duration time.Duration
***REMOVED***

// A user-printed log message.
type LogEntry struct ***REMOVED***
	Text string
***REMOVED***

type Runner interface ***REMOVED***
	Load(filename, src string) error
	RunVU(stop <-chan interface***REMOVED******REMOVED***) <-chan interface***REMOVED******REMOVED***
***REMOVED***

func NewError(err error) interface***REMOVED******REMOVED*** ***REMOVED***
	return err
***REMOVED***

func NewLogEntry(text string) interface***REMOVED******REMOVED*** ***REMOVED***
	return LogEntry***REMOVED***Text: text***REMOVED***
***REMOVED***

func NewMetric(start time.Time, duration time.Duration) interface***REMOVED******REMOVED*** ***REMOVED***
	return Metric***REMOVED***Start: start, Duration: duration***REMOVED***
***REMOVED***

func Run(r Runner, vus int, stop <-chan interface***REMOVED******REMOVED***) <-chan interface***REMOVED******REMOVED*** ***REMOVED***
	ch := make(chan interface***REMOVED******REMOVED***)

	go func() ***REMOVED***
		wg := sync.WaitGroup***REMOVED******REMOVED***
		for i := 0; i < vus; i++ ***REMOVED***
			wg.Add(1)
			go func() ***REMOVED***
				defer wg.Done()
				for res := range r.RunVU(stop) ***REMOVED***
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
