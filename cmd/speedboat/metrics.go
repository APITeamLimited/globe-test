package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/sampler"
	"time"
)

func printMetrics() ***REMOVED***
	for name, metric := range sampler.DefaultSampler.Metrics ***REMOVED***
		text := fmt.Sprintf("Metric: %s", name)
		switch metric.Type ***REMOVED***
		case sampler.CounterType:
			last := metric.Entries[len(metric.Entries)-1]
			log.WithField("val", applyIntent(metric, last.Value)).Info(text)
		case sampler.StatsType:
			log.WithFields(log.Fields***REMOVED***
				"min": applyIntent(metric, metric.Min()),
				"max": applyIntent(metric, metric.Max()),
				"avg": applyIntent(metric, metric.Avg()),
				"med": applyIntent(metric, metric.Med()),
			***REMOVED***).Info(text)
		***REMOVED***
	***REMOVED***
***REMOVED***

func applyIntent(m *sampler.Metric, v int64) interface***REMOVED******REMOVED*** ***REMOVED***
	if m.Intent == sampler.TimeIntent ***REMOVED***
		return time.Duration(v)
	***REMOVED***
	return v
***REMOVED***
