package runner

import (
	"time"
)

// A single metric for a test execution.
type Metric struct ***REMOVED***
	Time     time.Time
	Error    error
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
	Run(filename, src string) <-chan Result
***REMOVED***
