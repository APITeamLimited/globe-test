// Package engine contains the internal metrics engine responsible for
// aggregating metrics during the test and evaluating thresholds against them.
package engine

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/APITeamLimited/globe-test/metrics"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"
)

const (
	thresholdsCheckInterval = 2 * time.Second
)

// MetricsEngine is the internal metrics engine that k6 uses to keep track of
// aggregated metric sample values. They are used to generate the end-of-test
// summary and to evaluate the test thresholds.
type MetricsEngine struct {
	getCurrentTestRunDuration func() time.Duration

	options  *libWorker.Options
	logger   logrus.FieldLogger
	registry *metrics.Registry

	// These can be both top-level metrics or sub-metrics
	metricsWithThresholds []*metrics.Metric

	// TODO: completely refactor:
	//   - make these private,
	//   - do not use an unnecessary map for the observed metrics
	//   - have one lock per metric instead of a a global one, when
	//     the metrics are decoupled from their types
	MetricsLock     sync.Mutex
	ObservedMetrics map[string]*metrics.Metric

	threasholdTicker   *time.Ticker
	thresholdAbortChan chan struct{}
}

// NewMetricsEngine creates a new metrics Engine with the given parameters.
func NewMetricsEngine(options *libWorker.Options, logger *logrus.Logger,
	registry *metrics.Registry, getCurrentTestRunDuration func() time.Duration) (*MetricsEngine, error) {
	me := &MetricsEngine{
		options:  options,
		logger:   logger.WithField("component", "metrics-engine"),
		registry: registry,

		ObservedMetrics: make(map[string]*metrics.Metric),
	}

	err := me.initSubMetricsAndThresholds()
	if err != nil {
		return nil, err
	}

	return me, nil
}

func (me *MetricsEngine) getThresholdMetricOrSubmetric(name string) (*metrics.Metric, error) {
	nameParts := strings.SplitN(name, "{", 2)

	metric := me.registry.Get(nameParts[0])
	if metric == nil {
		return nil, fmt.Errorf("metric '%s' does not exist in the script", nameParts[0])
	}
	if len(nameParts) == 1 { // no sub-metric
		return metric, nil
	}

	submetricDefinition := nameParts[1]
	if submetricDefinition[len(submetricDefinition)-1] != '}' {
		return nil, fmt.Errorf("missing ending bracket, sub-metric format needs to be 'metric{key:value}'")
	}
	sm, err := metric.AddSubmetric(submetricDefinition[:len(submetricDefinition)-1])
	if err != nil {
		return nil, err
	}

	if sm.Metric.Observed {
		// Do not repeat warnings for the same sub-metrics
		return sm.Metric, nil
	}

	return sm.Metric, nil
}

func (me *MetricsEngine) markObserved(metric *metrics.Metric) {
	if !metric.Observed {
		metric.Observed = true
		me.ObservedMetrics[metric.Name] = metric
	}
}

func (me *MetricsEngine) initSubMetricsAndThresholds() error {
	for metricName, thresholds := range me.options.Thresholds {
		metric, err := me.getThresholdMetricOrSubmetric(metricName)

		if err != nil {
			return fmt.Errorf("invalid metric '%s' in threshold definitions: %w", metricName, err)
		}

		metric.Thresholds = thresholds
		me.metricsWithThresholds = append(me.metricsWithThresholds, metric)

		// Mark the metric (and the parent metric, if we're dealing with a
		// submetric) as observed, so they are shown in the end-of-test summary,
		// even if they don't have any metric samples during the test run
		me.markObserved(metric)
		if metric.Sub != nil {
			me.markObserved(metric.Sub.Parent)
		}
	}

	// TODO: refactor out of here when https://github.com/grafana/k6/issues/1321
	// lands and there is a better way to enable a metric with tag
	if me.options.SystemTags.Has(metrics.TagExpectedResponse) {
		_, err := me.getThresholdMetricOrSubmetric("http_req_duration{expected_response:true}")
		if err != nil {
			return err // shouldn't happen, but ¯\_(ツ)_/¯
		}
	}

	return nil
}

// EvaluateThresholds processes all of the thresholds.
func (me *MetricsEngine) evaluateThresholds(ignoreEmptySinks bool) (thresholdsTainted, shouldAbort bool) {
	me.MetricsLock.Lock()
	defer me.MetricsLock.Unlock()

	t := me.getCurrentTestRunDuration()

	for _, m := range me.metricsWithThresholds {
		// If either the metric has no thresholds defined, or its sinks
		// are empty, let's ignore its thresholds execution at this point.
		if len(m.Thresholds.Thresholds) == 0 || (ignoreEmptySinks && m.Sink.IsEmpty()) {
			continue
		}
		m.Tainted = null.BoolFrom(false)

		me.logger.WithField("metric_name", m.Name).Debug("running thresholds")
		succ, err := m.Thresholds.Run(m.Sink, t)
		if err != nil {
			me.logger.WithField("metric_name", m.Name).WithError(err).Error("Threshold error")
			continue
		}
		if succ {
			continue // threshold passed
		}
		me.logger.WithField("metric_name", m.Name).Debug("Thresholds failed")
		m.Tainted = null.BoolFrom(true)
		thresholdsTainted = true
		if m.Thresholds.Abort {
			shouldAbort = true
		}
	}

	return thresholdsTainted, shouldAbort
}

func (me *MetricsEngine) Start() {
	me.thresholdAbortChan = make(chan struct{})
	me.threasholdTicker = time.NewTicker(thresholdsCheckInterval)

	go func() {
		// MAKE sure to listen to ticker close
		for {
			select {
			case <-me.thresholdAbortChan:
				threasholdsTainted, shouldAbort := me.evaluateThresholds(true)

				if shouldAbort && threasholdsTainted {
					me.thresholdAbortChan <- struct{}{}
				}
			}
		}
	}()
}

func (me *MetricsEngine) Stop() {
	me.threasholdTicker.Stop()
}

// processMetrics process the execution's metrics samples as they are collected.
// The processing of samples happens at a fixed rate defined by the `collectRate`
// constant.
//
// The `processMetricsAfterRun` channel argument is used by the caller to signal
// that the test run is finished, no more metric samples will be produced, and that
// the metrics samples remaining in the pipeline should be should be processed.
func (e *MetricsEngine) processMetrics(globalCtx context.Context, processMetricsAfterRun chan struct{}) {
	sampleContainers := []metrics.SampleContainer{}

	defer func() {
		// Process any remaining metrics in the pipeline, by this point Run()
		// has already finished and nothing else should be producing metrics.
		e.logger.Debug("Metrics processing winding down...")

		close(e.Samples)
		for sc := range e.Samples {
			sampleContainers = append(sampleContainers, sc)
		}
		e.OutputManager.AddMetricSamples(sampleContainers)

		// Process the thresholds one final time
		thresholdsTainted, _ := e.MetricsEngine.EvaluateThresholds(false)
		e.thresholdsTaintedLock.Lock()
		e.thresholdsTainted = thresholdsTainted
		e.thresholdsTaintedLock.Unlock()
	}()

	ticker := time.NewTicker(collectRate)
	defer ticker.Stop()

	processSamples := func() {
		if len(sampleContainers) > 0 {
			e.OutputManager.AddMetricSamples(sampleContainers)
			// Make the new container with the same size as the previous
			// one, assuming that we produce roughly the same amount of
			// metrics data between ticks...
			sampleContainers = make([]metrics.SampleContainer, 0, cap(sampleContainers))
		}
	}
	for {
		select {
		case <-ticker.C:
			processSamples()
		case <-processMetricsAfterRun:
		getCachedMetrics:
			for {
				select {
				case sc := <-e.Samples:
					sampleContainers = append(sampleContainers, sc)
				default:
					break getCachedMetrics
				}
			}
			processSamples()
			// Ensure the ingester flushes any buffered metrics
			_ = e.ingester.Stop()
			thresholdsTainted, _ := e.MetricsEngine.EvaluateThresholds(false)
			e.thresholdsTaintedLock.Lock()
			e.thresholdsTainted = thresholdsTainted
			e.thresholdsTaintedLock.Unlock()

			processMetricsAfterRun <- struct{}{}

		case sc := <-e.Samples:
			sampleContainers = append(sampleContainers, sc)
		case <-globalCtx.Done():
			return
		}
	}
}
