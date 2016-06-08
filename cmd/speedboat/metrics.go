package main

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/sampler"
	"io"
	stdlog "log"
	"time"
)

type LogMetricsOutput struct ***REMOVED***
	Writer io.Writer

	encoder *json.Encoder
***REMOVED***

func (o *LogMetricsOutput) Write(m *sampler.Metric, e *sampler.Entry) error ***REMOVED***
	if o.encoder == nil ***REMOVED***
		o.encoder = json.NewEncoder(o.Writer)
	***REMOVED***
	return o.encoder.Encode(e)
***REMOVED***

func (o *LogMetricsOutput) Commit() error ***REMOVED***
	return nil
***REMOVED***

func printMetrics(l *stdlog.Logger) ***REMOVED***
	for name, m := range sampler.DefaultSampler.Metrics ***REMOVED***
		switch m.Type ***REMOVED***
		case sampler.GaugeType:
			l.Printf("%s val=%v\n", name, applyIntent(m, m.Last()))
		case sampler.CounterType:
			l.Printf("%s num=%v\n", name, applyIntent(m, m.Sum()))
		case sampler.StatsType:
			l.Printf("%s min=%v max=%v avg=%v med=%v\n", name,
				applyIntent(m, m.Min()),
				applyIntent(m, m.Max()),
				applyIntent(m, m.Avg()),
				applyIntent(m, m.Med()),
			)
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
