// Package engine contains the internal metrics engine responsible for
// aggregating metrics during the test and evaluating thresholds against them.
package engine

import (
	"fmt"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/metrics"
	"go.k6.io/k6/output"
	"gopkg.in/guregu/null.v3"
)

// MetricsEngine is the internal metrics engine that k6 uses to keep track of
// aggregated metric sample values. They are used to generate the end-of-test
// summary and to evaluate the test thresholds.
type MetricsEngine struct ***REMOVED***
	registry       *metrics.Registry
	executionState *lib.ExecutionState
	options        lib.Options
	runtimeOptions lib.RuntimeOptions
	logger         logrus.FieldLogger

	// These can be both top-level metrics or sub-metrics
	metricsWithThresholds []*metrics.Metric

	// TODO: completely refactor:
	//   - make these private,
	//   - do not use an unnecessary map for the observed metrics
	//   - have one lock per metric instead of a a global one, when
	//     the metrics are decoupled from their types
	MetricsLock     sync.Mutex
	ObservedMetrics map[string]*metrics.Metric
***REMOVED***

// NewMetricsEngine creates a new metrics Engine with the given parameters.
func NewMetricsEngine(
	registry *metrics.Registry, executionState *lib.ExecutionState,
	opts lib.Options, rtOpts lib.RuntimeOptions, logger logrus.FieldLogger,
) (*MetricsEngine, error) ***REMOVED***
	me := &MetricsEngine***REMOVED***
		registry:       registry,
		executionState: executionState,
		options:        opts,
		runtimeOptions: rtOpts,
		logger:         logger.WithField("component", "metrics-engine"),

		ObservedMetrics: make(map[string]*metrics.Metric),
	***REMOVED***

	if !(me.runtimeOptions.NoSummary.Bool && me.runtimeOptions.NoThresholds.Bool) ***REMOVED***
		err := me.initSubMetricsAndThresholds()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	return me, nil
***REMOVED***

// GetIngester returns a pseudo-Output that uses the given metric samples to
// update the engine's inner state.
func (me *MetricsEngine) GetIngester() output.Output ***REMOVED***
	return &outputIngester***REMOVED***
		logger:        me.logger.WithField("component", "metrics-engine-ingester"),
		metricsEngine: me,
	***REMOVED***
***REMOVED***

func (me *MetricsEngine) getThresholdMetricOrSubmetric(name string) (*metrics.Metric, error) ***REMOVED***
	// TODO: replace with strings.Cut after Go 1.18
	nameParts := strings.SplitN(name, "***REMOVED***", 2)

	metric := me.registry.Get(nameParts[0])
	if metric == nil ***REMOVED***
		return nil, fmt.Errorf("metric '%s' does not exist in the script", nameParts[0])
	***REMOVED***
	if len(nameParts) == 1 ***REMOVED*** // no sub-metric
		return metric, nil
	***REMOVED***

	submetricDefinition := nameParts[1]
	if submetricDefinition[len(submetricDefinition)-1] != '***REMOVED***' ***REMOVED***
		return nil, fmt.Errorf("missing ending bracket, sub-metric format needs to be 'metric***REMOVED***key:value***REMOVED***'")
	***REMOVED***
	sm, err := metric.AddSubmetric(submetricDefinition[:len(submetricDefinition)-1])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return sm.Metric, nil
***REMOVED***

func (me *MetricsEngine) markObserved(metric *metrics.Metric) ***REMOVED***
	if !metric.Observed ***REMOVED***
		metric.Observed = true
		me.ObservedMetrics[metric.Name] = metric
	***REMOVED***
***REMOVED***

func (me *MetricsEngine) initSubMetricsAndThresholds() error ***REMOVED***
	for metricName, thresholds := range me.options.Thresholds ***REMOVED***
		metric, err := me.getThresholdMetricOrSubmetric(metricName)

		if me.runtimeOptions.NoThresholds.Bool ***REMOVED***
			if err != nil ***REMOVED***
				me.logger.WithError(err).Warnf("Invalid metric '%s' in threshold definitions", metricName)
			***REMOVED***
			continue
		***REMOVED***

		if err != nil ***REMOVED***
			return fmt.Errorf("invalid metric '%s' in threshold definitions: %w", metricName, err)
		***REMOVED***

		metric.Thresholds = thresholds
		me.metricsWithThresholds = append(me.metricsWithThresholds, metric)

		// Mark the metric (and the parent metric, if we're dealing with a
		// submetric) as observed, so they are shown in the end-of-test summary,
		// even if they don't have any metric samples during the test run
		me.markObserved(metric)
		if metric.Sub != nil ***REMOVED***
			me.markObserved(metric.Sub.Parent)
		***REMOVED***
	***REMOVED***

	// TODO: refactor out of here when https://github.com/grafana/k6/issues/1321
	// lands and there is a better way to enable a metric with tag
	if me.options.SystemTags.Has(metrics.TagExpectedResponse) ***REMOVED***
		_, err := me.getThresholdMetricOrSubmetric("http_req_duration***REMOVED***expected_response:true***REMOVED***")
		if err != nil ***REMOVED***
			return err // shouldn't happen, but ¯\_(ツ)_/¯
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// EvaluateThresholds processes all of the thresholds.
//
// TODO: refactor, make private, optimize
func (me *MetricsEngine) EvaluateThresholds(ignoreEmptySinks bool) (thresholdsTainted, shouldAbort bool) ***REMOVED***
	me.MetricsLock.Lock()
	defer me.MetricsLock.Unlock()

	t := me.executionState.GetCurrentTestRunDuration()

	for _, m := range me.metricsWithThresholds ***REMOVED***
		// If either the metric has no thresholds defined, or its sinks
		// are empty, let's ignore its thresholds execution at this point.
		if len(m.Thresholds.Thresholds) == 0 || (ignoreEmptySinks && m.Sink.IsEmpty()) ***REMOVED***
			continue
		***REMOVED***
		m.Tainted = null.BoolFrom(false)

		me.logger.WithField("metric_name", m.Name).Debug("running thresholds")
		succ, err := m.Thresholds.Run(m.Sink, t)
		if err != nil ***REMOVED***
			me.logger.WithField("metric_name", m.Name).WithError(err).Error("Threshold error")
			continue
		***REMOVED***
		if succ ***REMOVED***
			continue // threshold passed
		***REMOVED***
		me.logger.WithField("metric_name", m.Name).Debug("Thresholds failed")
		m.Tainted = null.BoolFrom(true)
		thresholdsTainted = true
		if m.Thresholds.Abort ***REMOVED***
			shouldAbort = true
		***REMOVED***
	***REMOVED***

	return thresholdsTainted, shouldAbort
***REMOVED***
