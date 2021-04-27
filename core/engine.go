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
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/output"
	"github.com/loadimpact/k6/stats"
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

	Options        lib.Options
	runtimeOptions lib.RuntimeOptions
	outputs        []output.Output

	logger   *logrus.Entry
	stopOnce sync.Once
	stopChan chan struct***REMOVED******REMOVED***

	Metrics     map[string]*stats.Metric
	MetricsLock sync.Mutex

	Samples chan stats.SampleContainer

	// Assigned to metrics upon first received sample.
	thresholds map[string]stats.Thresholds
	submetrics map[string][]*stats.Submetric

	// Are thresholds tainted?
	thresholdsTainted bool
***REMOVED***

// NewEngine instantiates a new Engine, without doing any heavy initialization.
func NewEngine(
	ex lib.ExecutionScheduler, opts lib.Options, rtOpts lib.RuntimeOptions, outputs []output.Output, logger *logrus.Logger,
) (*Engine, error) ***REMOVED***
	if ex == nil ***REMOVED***
		return nil, errors.New("missing ExecutionScheduler instance")
	***REMOVED***

	e := &Engine***REMOVED***
		ExecutionScheduler: ex,
		executionState:     ex.GetState(),

		Options:        opts,
		runtimeOptions: rtOpts,
		outputs:        outputs,
		Metrics:        make(map[string]*stats.Metric),
		Samples:        make(chan stats.SampleContainer, opts.MetricSamplesBufferSize.Int64),
		stopChan:       make(chan struct***REMOVED******REMOVED***),
		logger:         logger.WithField("component", "engine"),
	***REMOVED***

	e.thresholds = opts.Thresholds
	e.submetrics = make(map[string][]*stats.Submetric)
	for name := range e.thresholds ***REMOVED***
		if !strings.Contains(name, "***REMOVED***") ***REMOVED***
			continue
		***REMOVED***

		parent, sm := stats.NewSubmetric(name)
		e.submetrics[parent] = append(e.submetrics[parent], sm)
	***REMOVED***

	// TODO: refactor this out of here when https://github.com/loadimpact/k6/issues/1832 lands and
	// there is a better way to enable a metric with tag
	if opts.SystemTags.Has(stats.TagExpectedResponse) ***REMOVED***
		for _, name := range []string***REMOVED***
			"http_req_duration***REMOVED***expected_response:true***REMOVED***",
		***REMOVED*** ***REMOVED***
			if _, ok := e.thresholds[name]; ok ***REMOVED***
				continue
			***REMOVED***
			parent, sm := stats.NewSubmetric(name)
			e.submetrics[parent] = append(e.submetrics[parent], sm)
		***REMOVED***
	***REMOVED***

	return e, nil
***REMOVED***

// StartOutputs spins up all configured outputs, giving the thresholds to any
// that can accept them. And if some output fails, stop the already started
// ones. This may take some time, since some outputs make initial network
// requests to set up whatever remote services are going to listen to them.
//
// TODO: this doesn't really need to be in the Engine, so take it out?
func (e *Engine) StartOutputs() error ***REMOVED***
	e.logger.Debugf("Starting %d outputs...", len(e.outputs))
	for i, out := range e.outputs ***REMOVED***
		if thresholdOut, ok := out.(output.WithThresholds); ok ***REMOVED***
			thresholdOut.SetThresholds(e.thresholds)
		***REMOVED***

		if stopOut, ok := out.(output.WithTestRunStop); ok ***REMOVED***
			stopOut.SetTestRunStopCallback(
				func(err error) ***REMOVED***
					e.logger.WithError(err).Error("Received error to stop from output")
					e.Stop()
				***REMOVED***)
		***REMOVED***

		if err := out.Start(); err != nil ***REMOVED***
			e.stopOutputs(i)
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// StopOutputs stops all configured outputs.
func (e *Engine) StopOutputs() ***REMOVED***
	e.stopOutputs(len(e.outputs))
***REMOVED***

func (e *Engine) stopOutputs(upToID int) ***REMOVED***
	e.logger.Debugf("Stopping %d outputs...", upToID)
	for i := 0; i < upToID; i++ ***REMOVED***
		if err := e.outputs[i].Stop(); err != nil ***REMOVED***
			e.logger.WithError(err).Errorf("Stopping output %d failed", i)
		***REMOVED***
	***REMOVED***
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
		err := e.ExecutionScheduler.Run(globalCtx, runSubCtx, e.Samples)
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
				var serr types.ScriptException
				if errors.As(err, &serr) ***REMOVED***
					e.setRunStatus(lib.RunStatusAbortedScriptError)
				***REMOVED*** else ***REMOVED***
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
				Metric: metrics.VUs,
				Value:  float64(executionState.GetCurrentlyActiveVUsCount()),
				Tags:   e.Options.RunTags,
			***REMOVED***, ***REMOVED***
				Time:   t,
				Metric: metrics.VUsMax,
				Value:  float64(executionState.GetInitializedVUsCount()),
				Tags:   e.Options.RunTags,
			***REMOVED***,
		***REMOVED***,
		Tags: e.Options.RunTags,
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

func (e *Engine) processSamplesForMetrics(sampleContainers []stats.SampleContainer) ***REMOVED***
	for _, sampleContainer := range sampleContainers ***REMOVED***
		samples := sampleContainer.GetSamples()

		if len(samples) == 0 ***REMOVED***
			continue
		***REMOVED***

		for _, sample := range samples ***REMOVED***
			m, ok := e.Metrics[sample.Metric.Name]
			if !ok ***REMOVED***
				m = stats.New(sample.Metric.Name, sample.Metric.Type, sample.Metric.Contains)
				m.Thresholds = e.thresholds[m.Name]
				m.Submetrics = e.submetrics[m.Name]
				e.Metrics[m.Name] = m
			***REMOVED***
			m.Sink.Add(sample)

			for _, sm := range m.Submetrics ***REMOVED***
				if !sample.Tags.Contains(sm.Tags) ***REMOVED***
					continue
				***REMOVED***

				if sm.Metric == nil ***REMOVED***
					sm.Metric = stats.New(sm.Name, sample.Metric.Type, sample.Metric.Contains)
					sm.Metric.Sub = *sm
					sm.Metric.Thresholds = e.thresholds[sm.Name]
					e.Metrics[sm.Name] = sm.Metric
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
		e.processSamplesForMetrics(sampleContainers)
	***REMOVED***

	for _, out := range e.outputs ***REMOVED***
		out.AddMetricSamples(sampleContainers)
	***REMOVED***
***REMOVED***
