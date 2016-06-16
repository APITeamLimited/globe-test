package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/sampler"
	"time"
)

func printMetrics() ***REMOVED***
	for name, m := range sampler.DefaultSampler.Metrics ***REMOVED***
		switch m.Type ***REMOVED***
		case sampler.GaugeType:
			last := m.Last()
			if last == 0 ***REMOVED***
				continue
			***REMOVED***
			log.WithField("val", applyIntent(m, last)).Infof("Metric: %s", name)
		case sampler.CounterType:
			sum := m.Sum()
			if sum == 0 ***REMOVED***
				continue
			***REMOVED***
			log.WithField("num", applyIntent(m, sum)).Infof("Metric: %s", name)
		case sampler.StatsType:
			log.WithFields(log.Fields***REMOVED***
				"min":   applyIntent(m, m.Min()),
				"max":   applyIntent(m, m.Max()),
				"avg":   applyIntent(m, m.Avg()),
				"med":   applyIntent(m, m.Med()),
				"count": m.Count(),
			***REMOVED***).Infof("Metric: %s", name)
		***REMOVED***
	***REMOVED***
***REMOVED***

func commitMetrics() ***REMOVED***
	if err := sampler.DefaultSampler.Commit(); err != nil ***REMOVED***
		log.WithError(err).Error("Couldn't write samples!")
	***REMOVED***
***REMOVED***

func closeMetrics() ***REMOVED***
	if err := sampler.DefaultSampler.Commit(); err != nil ***REMOVED***
		log.WithError(err).Error("Couldn't close sampler!")
	***REMOVED***
***REMOVED***

func applyIntent(m *sampler.Metric, v int64) interface***REMOVED******REMOVED*** ***REMOVED***
	if m.Intent == sampler.TimeIntent ***REMOVED***
		return time.Duration(v)
	***REMOVED***
	return v
***REMOVED***
