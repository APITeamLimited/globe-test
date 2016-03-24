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

// An envelope for a result.
type Result struct ***REMOVED***
	Type     string
	Error    error
	LogEntry LogEntry
	Metric   Metric
***REMOVED***

type Runner interface ***REMOVED***
	Load(filename, src string) error
	RunVU() <-chan Result
***REMOVED***

func NewError(err error) ***REMOVED***
	return Result***REMOVED***Type: "error", Error: err***REMOVED***
***REMOVED***

func NewLogEntry(entry LogEntry) ***REMOVED***
	return Result***REMOVED***Type: "log", LogEntry: entry***REMOVED***
***REMOVED***

func NewMetric(metric Metric) ***REMOVED***
	return Result***REMOVED***Type: "metric", Metric: metric***REMOVED***
***REMOVED***

func Run(r Runner, vus int) <-chan Result ***REMOVED***
	ch := make(chan Result)

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
