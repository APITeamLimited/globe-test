package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/sampler"
	stdlog "log"
	"time"
)

func printMetrics(l *stdlog.Logger) ***REMOVED***
	for name, m := range sampler.DefaultSampler.Metrics ***REMOVED***
		l.Printf("%s\n", name)
		switch m.Type ***REMOVED***
		case sampler.GaugeType:
			l.Printf("  value=%s\n", applyIntent(m, m.Last()))
		case sampler.CounterType:
			l.Printf("  num=%s\n", applyIntent(m, m.Sum()))
		case sampler.StatsType:
			l.Printf("  min=%s\n", applyIntent(m, m.Min()))
			l.Printf("  max=%s\n", applyIntent(m, m.Max()))
			l.Printf("  avg=%s\n", applyIntent(m, m.Avg()))
			l.Printf("  med=%s\n", applyIntent(m, m.Med()))
		***REMOVED***
	***REMOVED***
***REMOVED***

func commitMetrics() ***REMOVED***
	if err := sampler.DefaultSampler.Commit(); err != nil ***REMOVED***
		log.WithError(err).Error("Couldn't write samples!")
	***REMOVED***
***REMOVED***

func applyIntent(m *sampler.Metric, v int64) interface***REMOVED******REMOVED*** ***REMOVED***
	if m.Intent == sampler.TimeIntent ***REMOVED***
		return time.Duration(v)
	***REMOVED***
	return fmt.Sprint(v)
***REMOVED***
