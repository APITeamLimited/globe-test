/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package core

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/errext"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/metrics"
	"go.k6.io/k6/output"
	"go.k6.io/k6/stats"
)

const (
	metricsRate    = 1 * time.Second
	collectRate    = 50 * time.Millisecond
	thresholdsRate = 2 * time.Second
)

// The Engine is the beating heart of k6.
type Engine struct ***REMOVED***
	// TODO: Make most of the stuff here private! And think how to refactor the
	// engine to be less stateful... it's currently one big mess of moving
	// pieces, and you implicitly first have to call Init() and then Run() -
	// maybe we should refactor it so we have a `Session` dauther-object that
	// Init() returns? The only problem with doing this is the REST API - it
	// expects to be able to get information from the Engine and is initialized
	// before the Init() call...

	ExecutionScheduler lib.ExecutionScheduler
	executionState     *lib.ExecutionState

	options        lib.Options
	runtimeOptions lib.RuntimeOptions
	outputs        []output.Output

	logger   *logrus.Entry
	stopOnce sync.Once
	stopChan chan struct***REMOVED******REMOVED***

	Metrics     map[string]*stats.Metric // TODO: refactor, this doesn't need to be a map
	MetricsLock sync.Mutex

	registry       *metrics.Registry
	builtinMetrics *metrics.BuiltinMetrics
	Samples        chan stats.SampleContainer

	// These can be both top-level metrics or sub-metrics
	metricsWithThresholds []*stats.Metric

	// Are thresholds tainted?
	thresholdsTainted bool
***REMOVED***

// NewEngine instantiates a new Engine, without doing any heavy initialization.
func NewEngine(
	ex lib.ExecutionScheduler, opts lib.Options, rtOpts lib.RuntimeOptions, outputs []output.Output, logger *logrus.Logger,
	registry *metrics.Registry, builtinMetrics *metrics.BuiltinMetrics,
) (*Engine, error) ***REMOVED***
	if ex == nil ***REMOVED***
		return nil, errors.New("missing ExecutionScheduler instance")
	***REMOVED***

	e := &Engine***REMOVED***
		ExecutionScheduler: ex,
		executionState:     ex.GetState(),

		options:        opts,
		runtimeOptions: rtOpts,
		outputs:        outputs,
		Metrics:        make(map[string]*stats.Metric),
		Samples:        make(chan stats.SampleContainer, opts.MetricSamplesBufferSize.Int64),
		stopChan:       make(chan struct***REMOVED******REMOVED***),
		logger:         logger.WithField("component", "engine"),
		registry:       registry,
		builtinMetrics: builtinMetrics,
	***REMOVED***

	if !(e.runtimeOptions.NoSummary.Bool && e.runtimeOptions.NoThresholds.Bool) ***REMOVED***
		err := e.initSubMetricsAndThresholds()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	return e, nil
***REMOVED***

func (e *Engine) getOrInitPotentialSubmetric(name string) (*stats.Metric, error) ***REMOVED***
	// TODO: replace with strings.Cut after Go 1.18
	nameParts := strings.SplitN(name, "***REMOVED***", 2)

	metric := e.registry.Get(nameParts[0])
	if metric == nil ***REMOVED***
		return nil, fmt.Errorf("metric '%s' does not exist in the script", nameParts[0])
	***REMOVED***
	if len(nameParts) == 1 ***REMOVED*** // no sub-metric
		return metric, nil
	***REMOVED***

	if nameParts[1][len(nameParts[1])-1] != '***REMOVED***' ***REMOVED***
		return nil, fmt.Errorf("missing ending bracket, sub-metric format needs to be 'metric***REMOVED***key:value***REMOVED***'")
	***REMOVED***
	sm, err := metric.AddSubmetric(nameParts[1][:len(nameParts[1])-1])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return sm.Metric, nil
***REMOVED***

func (e *Engine) initSubMetricsAndThresholds() error ***REMOVED***
	for metricName, thresholds := range e.options.Thresholds ***REMOVED***
		metric, err := e.getOrInitPotentialSubmetric(metricName)

		if e.runtimeOptions.NoThresholds.Bool ***REMOVED***
			if err != nil ***REMOVED***
				e.logger.WithError(err).Warnf("Invalid metric '%s' in threshold definitions", metricName)
			***REMOVED***
			continue
		***REMOVED***

		if err != nil ***REMOVED***
			return fmt.Errorf("invalid metric '%s' in threshold definitions: %w", metricName, err)
		***REMOVED***

		metric.Thresholds = thresholds
		e.metricsWithThresholds = append(e.metricsWithThresholds, metric)
	***REMOVED***

	// TODO: refactor out of here when https://github.com/grafana/k6/issues/1321
	// lands and there is a better way to enable a metric with tag
	if e.options.SystemTags.Has(stats.TagExpectedResponse) ***REMOVED***
		_, err := e.getOrInitPotentialSubmetric("http_req_duration***REMOVED***expected_response:true***REMOVED***")
		if err != nil ***REMOVED***
			return err // shouldn't happen, but ¯\_(ツ)_/¯
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// Init is used to initialize the execution scheduler and all metrics processing
// in the engine. The first is a costly operation, since it initializes all of
// the planned VUs and could potentially take a long time.
//
// This method either returns an error immediately, or it returns test run() and
// wait() functions.
//
// Things to note:
//  - The first lambda, Run(), synchronously executes the actual load test.
//  - It can be prematurely aborted by cancelling the runCtx - this won't stop
//    the metrics collection by the Engine.
//  - Stopping the metrics collection can be done at any time after Run() has
//    returned by cancelling the globalCtx
//  - The second returned lambda can be used to wait for that process to finish.
func (e *Engine) Init(globalCtx, runCtx context.Context) (run func() error, wait func(), err error) ***REMOVED***
	e.logger.Debug("Initialization starting...")
	// TODO: if we ever need metrics processing in the init context, we can move
	// this below the other components... or even start them concurrently?
	if err := e.ExecutionScheduler.Init(runCtx, e.Samples); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	// TODO: move all of this in a separate struct? see main TODO above
	runSubCtx, runSubCancel := context.WithCancel(runCtx)

	resultCh := make(chan error)
	processMetricsAfterRun := make(chan struct***REMOVED******REMOVED***)
	runFn := func() error ***REMOVED***
		e.logger.Debug("Execution scheduler starting...")
		err := e.ExecutionScheduler.Run(globalCtx, runSubCtx, e.Samples, e.builtinMetrics)
		e.logger.WithError(err).Debug("Execution scheduler terminated")

		select ***REMOVED***
		case <-runSubCtx.Done():
			// do nothing, the test run was aborted somehow
		default:
			resultCh <- err // we finished normally, so send the result
		***REMOVED***

		// Make the background jobs process the currently buffered metrics and
		// run the thresholds, then wait for that to be done.
		processMetricsAfterRun <- struct***REMOVED******REMOVED******REMOVED******REMOVED***
		<-processMetricsAfterRun

		return err
	***REMOVED***

	waitFn := e.startBackgroundProcesses(globalCtx, runCtx, resultCh, runSubCancel, processMetricsAfterRun)
	return runFn, waitFn, nil
***REMOVED***

// This starts a bunch of goroutines to process metrics, thresholds, and set the
// test run status when it ends. It returns a function that can be used after
// the provided context is called, to wait for the complete winding down of all
// started goroutines.
func (e *Engine) startBackgroundProcesses(
	globalCtx, runCtx context.Context, runResult <-chan error, runSubCancel func(), processMetricsAfterRun chan struct***REMOVED******REMOVED***,
) (wait func()) ***REMOVED***
	processes := new(sync.WaitGroup)

	// Siphon and handle all produced metric samples
	processes.Add(1)
	go func() ***REMOVED***
		defer processes.Done()
		e.processMetrics(globalCtx, processMetricsAfterRun)
	***REMOVED***()

	// Run VU metrics emission, only while the test is running.
	// TODO: move? this seems like something the ExecutionScheduler should emit...
	processes.Add(1)
	go func() ***REMOVED***
		defer processes.Done()
		e.logger.Debug("Starting emission of VU metrics...")
		e.runMetricsEmission(runCtx)
		e.logger.Debug("Metrics emission terminated")
	***REMOVED***()

	// Update the test run status when the test finishes
	processes.Add(1)
	thresholdAbortChan := make(chan struct***REMOVED******REMOVED***)
	go func() ***REMOVED***
		defer processes.Done()
		select ***REMOVED***
		case err := <-runResult:
			if err != nil ***REMOVED***
				e.logger.WithError(err).Debug("run: execution scheduler returned an error")
				var serr errext.Exception
				switch ***REMOVED***
				case errors.As(err, &serr):
					e.setRunStatus(lib.RunStatusAbortedScriptError)
				case common.IsInterruptError(err):
					e.setRunStatus(lib.RunStatusAbortedUser)
				default:
					e.setRunStatus(lib.RunStatusAbortedSystem)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				e.logger.Debug("run: execution scheduler terminated")
				e.setRunStatus(lib.RunStatusFinished)
			***REMOVED***
		case <-runCtx.Done():
			e.logger.Debug("run: context expired; exiting...")
			e.setRunStatus(lib.RunStatusAbortedUser)
		case <-e.stopChan:
			runSubCancel()
			e.logger.Debug("run: stopped by user; exiting...")
			e.setRunStatus(lib.RunStatusAbortedUser)
		case <-thresholdAbortChan:
			e.logger.Debug("run: stopped by thresholds; exiting...")
			runSubCancel()
			e.setRunStatus(lib.RunStatusAbortedThreshold)
		***REMOVED***
	***REMOVED***()

	// Run thresholds, if not disabled.
	if !e.runtimeOptions.NoThresholds.Bool ***REMOVED***
		processes.Add(1)
		go func() ***REMOVED***
			defer processes.Done()
			defer e.logger.Debug("Engine: Thresholds terminated")
			ticker := time.NewTicker(thresholdsRate)
			defer ticker.Stop()

			for ***REMOVED***
				select ***REMOVED***
				case <-ticker.C:
					if e.processThresholds() ***REMOVED***
						close(thresholdAbortChan)
						return
					***REMOVED***
				case <-runCtx.Done():
					return
				***REMOVED***
			***REMOVED***
		***REMOVED***()
	***REMOVED***

	return processes.Wait
***REMOVED***

func (e *Engine) processMetrics(globalCtx context.Context, processMetricsAfterRun chan struct***REMOVED******REMOVED***) ***REMOVED***
	sampleContainers := []stats.SampleContainer***REMOVED******REMOVED***

	defer func() ***REMOVED***
		// Process any remaining metrics in the pipeline, by this point Run()
		// has already finished and nothing else should be producing metrics.
		e.logger.Debug("Metrics processing winding down...")

		close(e.Samples)
		for sc := range e.Samples ***REMOVED***
			sampleContainers = append(sampleContainers, sc)
		***REMOVED***
		e.processSamples(sampleContainers)

		if !e.runtimeOptions.NoThresholds.Bool ***REMOVED***
			e.processThresholds() // Process the thresholds one final time
		***REMOVED***
	***REMOVED***()

	ticker := time.NewTicker(collectRate)
	defer ticker.Stop()

	e.logger.Debug("Metrics processing started...")
	processSamples := func() ***REMOVED***
		if len(sampleContainers) > 0 ***REMOVED***
			e.processSamples(sampleContainers)
			// Make the new container with the same size as the previous
			// one, assuming that we produce roughly the same amount of
			// metrics data between ticks...
			sampleContainers = make([]stats.SampleContainer, 0, cap(sampleContainers))
		***REMOVED***
	***REMOVED***
	for ***REMOVED***
		select ***REMOVED***
		case <-ticker.C:
			processSamples()
		case <-processMetricsAfterRun:
		getCachedMetrics:
			for ***REMOVED***
				select ***REMOVED***
				case sc := <-e.Samples:
					sampleContainers = append(sampleContainers, sc)
				default:
					break getCachedMetrics
				***REMOVED***
			***REMOVED***
			e.logger.Debug("Processing metrics and thresholds after the test run has ended...")
			processSamples()
			if !e.runtimeOptions.NoThresholds.Bool ***REMOVED***
				e.processThresholds()
			***REMOVED***
			processMetricsAfterRun <- struct***REMOVED******REMOVED******REMOVED******REMOVED***

		case sc := <-e.Samples:
			sampleContainers = append(sampleContainers, sc)
		case <-globalCtx.Done():
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *Engine) setRunStatus(status lib.RunStatus) ***REMOVED***
	for _, out := range e.outputs ***REMOVED***
		if statUpdOut, ok := out.(output.WithRunStatusUpdates); ok ***REMOVED***
			statUpdOut.SetRunStatus(status)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *Engine) IsTainted() bool ***REMOVED***
	return e.thresholdsTainted
***REMOVED***

// Stop closes a signal channel, forcing a running Engine to return
func (e *Engine) Stop() ***REMOVED***
	e.stopOnce.Do(func() ***REMOVED***
		close(e.stopChan)
	***REMOVED***)
***REMOVED***

// IsStopped returns a bool indicating whether the Engine has been stopped
func (e *Engine) IsStopped() bool ***REMOVED***
	select ***REMOVED***
	case <-e.stopChan:
		return true
	default:
		return false
	***REMOVED***
***REMOVED***

func (e *Engine) runMetricsEmission(ctx context.Context) ***REMOVED***
	ticker := time.NewTicker(metricsRate)
	for ***REMOVED***
		select ***REMOVED***
		case <-ticker.C:
			e.emitMetrics()
		case <-ctx.Done():
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *Engine) emitMetrics() ***REMOVED***
	t := time.Now()

	executionState := e.ExecutionScheduler.GetState()
	// TODO: optimize and move this, it shouldn't call processSamples() directly
	e.processSamples([]stats.SampleContainer***REMOVED***stats.ConnectedSamples***REMOVED***
		Samples: []stats.Sample***REMOVED***
			***REMOVED***
				Time:   t,
				Metric: e.builtinMetrics.VUs,
				Value:  float64(executionState.GetCurrentlyActiveVUsCount()),
				Tags:   e.options.RunTags,
			***REMOVED***, ***REMOVED***
				Time:   t,
				Metric: e.builtinMetrics.VUsMax,
				Value:  float64(executionState.GetInitializedVUsCount()),
				Tags:   e.options.RunTags,
			***REMOVED***,
		***REMOVED***,
		Tags: e.options.RunTags,
		Time: t,
	***REMOVED******REMOVED***)
***REMOVED***

func (e *Engine) processThresholds() (shouldAbort bool) ***REMOVED***
	e.MetricsLock.Lock()
	defer e.MetricsLock.Unlock()

	t := e.executionState.GetCurrentTestRunDuration()

	e.thresholdsTainted = false
	for _, m := range e.Metrics ***REMOVED***
		if len(m.Thresholds.Thresholds) == 0 ***REMOVED***
			continue
		***REMOVED***
		m.Tainted = null.BoolFrom(false)

		e.logger.WithField("m", m.Name).Debug("running thresholds")
		succ, err := m.Thresholds.Run(m.Sink, t)
		if err != nil ***REMOVED***
			e.logger.WithField("m", m.Name).WithError(err).Error("Threshold error")
			continue
		***REMOVED***
		if !succ ***REMOVED***
			e.logger.WithField("m", m.Name).Debug("Thresholds failed")
			m.Tainted = null.BoolFrom(true)
			e.thresholdsTainted = true
			if m.Thresholds.Abort ***REMOVED***
				shouldAbort = true
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return shouldAbort
***REMOVED***

func (e *Engine) processMetricsInSamples(sampleContainers []stats.SampleContainer) ***REMOVED***
	for _, sampleContainer := range sampleContainers ***REMOVED***
		samples := sampleContainer.GetSamples()

		if len(samples) == 0 ***REMOVED***
			continue
		***REMOVED***

		for _, sample := range samples ***REMOVED***
			m := sample.Metric // this should have come from the Registry, no need to look it up
			if !m.Observed ***REMOVED***
				// But we need to add it here, so we can show data in the
				// end-of-test summary for this metric
				e.Metrics[m.Name] = m
				m.Observed = true
			***REMOVED***
			m.Sink.Add(sample) // add its value to its own sink

			// and also add it to any submetrics that match
			for _, sm := range m.Submetrics ***REMOVED***
				if !sample.Tags.Contains(sm.Tags) ***REMOVED***
					continue
				***REMOVED***
				if !sm.Metric.Observed ***REMOVED***
					// But we need to add it here, so we can show data in the
					// end-of-test summary for this metric
					e.Metrics[sm.Metric.Name] = sm.Metric
					sm.Metric.Observed = true
				***REMOVED***
				sm.Metric.Sink.Add(sample)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *Engine) processSamples(sampleContainers []stats.SampleContainer) ***REMOVED***
	if len(sampleContainers) == 0 ***REMOVED***
		return
	***REMOVED***

	// TODO: optimize this...
	e.MetricsLock.Lock()
	defer e.MetricsLock.Unlock()

	// TODO: run this and the below code in goroutines?
	if !(e.runtimeOptions.NoSummary.Bool && e.runtimeOptions.NoThresholds.Bool) ***REMOVED***
		e.processMetricsInSamples(sampleContainers)
	***REMOVED***

	for _, out := range e.outputs ***REMOVED***
		out.AddMetricSamples(sampleContainers)
	***REMOVED***
***REMOVED***
