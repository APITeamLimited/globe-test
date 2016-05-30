package main

import (
	"fmt"
	"github.com/loadimpact/speedboat/sampler"
	stdlog "log"
	"time"
)

func printMetrics(l *stdlog.Logger) ***REMOVED***
	for name, m := range sampler.DefaultSampler.Metrics ***REMOVED***
		l.Printf("%s\n", name)
		switch m.Type ***REMOVED***
		case sampler.CounterType:
			last := m.Entries[len(m.Entries)-1]
			l.Printf("  value=%s\n", applyIntent(m, last.Value))
		case sampler.StatsType:
			l.Printf("  min=%-15s max=%s\n", applyIntent(m, m.Min()), applyIntent(m, m.Max()))
			l.Printf("  avg=%-15s med=%s\n", applyIntent(m, m.Avg()), applyIntent(m, m.Med()))
		***REMOVED***
	***REMOVED***
***REMOVED***

func applyIntent(m *sampler.Metric, v int64) interface***REMOVED******REMOVED*** ***REMOVED***
	if m.Intent == sampler.TimeIntent ***REMOVED***
		return time.Duration(v)
	***REMOVED***
	return fmt.Sprint(v)
***REMOVED***
