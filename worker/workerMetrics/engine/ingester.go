package engine

import (
	"time"

	"github.com/APITeamLimited/globe-test/worker/output"
	"github.com/sirupsen/logrus"
)

const collectRate = 50 * time.Millisecond

var _ output.Output = &outputIngester***REMOVED******REMOVED***

// outputIngester implements the output.Output interface and can be used to
// "feed" the MetricsEngine data from a `k6 run` test run.
type outputIngester struct ***REMOVED***
	output.SampleBuffer
	logger logrus.FieldLogger

	metricsEngine   *MetricsEngine
	periodicFlusher *output.PeriodicFlusher
***REMOVED***

// Description returns a human-readable description of the output.
func (oi *outputIngester) Description() string ***REMOVED***
	return "engine"
***REMOVED***

// Start the engine by initializing a new output.PeriodicFlusher
func (oi *outputIngester) Start() error ***REMOVED***
	oi.logger.Debug("Starting...")

	pf, err := output.NewPeriodicFlusher(collectRate, oi.flushMetrics)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	oi.logger.Debug("Started!")
	oi.periodicFlusher = pf

	return nil
***REMOVED***

// Stop flushes any remaining metrics and stops the goroutine.
func (oi *outputIngester) Stop() error ***REMOVED***
	oi.logger.Debug("Stopping...")
	defer oi.logger.Debug("Stopped!")
	oi.periodicFlusher.Stop()
	return nil
***REMOVED***

// flushMetrics Writes samples to the MetricsEngine
func (oi *outputIngester) flushMetrics() ***REMOVED***
	sampleContainers := oi.GetBufferedSamples()
	if len(sampleContainers) == 0 ***REMOVED***
		return
	***REMOVED***

	oi.metricsEngine.MetricsLock.Lock()
	defer oi.metricsEngine.MetricsLock.Unlock()

	// TODO: split metric samples in buckets with a *workerMetrics.Metric key; this will
	// allow us to have a per-bucket lock, instead of one global one, and it
	// will allow us to split apart the metric Name and Type from its Sink and
	// Observed fields...
	//
	// And, to further optimize things, if every metric (and sub-metric) had a
	// sequential integer ID, we would be able to use a slice for these buckets
	// and eliminate the map loopkups altogether!

	for _, sampleContainer := range sampleContainers ***REMOVED***
		samples := sampleContainer.GetSamples()

		if len(samples) == 0 ***REMOVED***
			continue
		***REMOVED***

		for _, sample := range samples ***REMOVED***
			m := sample.Metric               // this should have come from the Registry, no need to look it up
			oi.metricsEngine.markObserved(m) // mark it as observed so it shows in the end-of-test summary
			m.Sink.Add(sample)               // finally, add its value to its own sink

			// and also to the same for any submetrics that match the metric sample
			for _, sm := range m.Submetrics ***REMOVED***
				if !sample.Tags.Contains(sm.Tags) ***REMOVED***
					continue
				***REMOVED***
				oi.metricsEngine.markObserved(sm.Metric)
				sm.Metric.Sink.Add(sample)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
