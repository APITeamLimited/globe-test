package core

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/sirupsen/logrus"

	"go.k6.io/k6/errext"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/metrics"
	"go.k6.io/k6/metrics/engine"
	"go.k6.io/k6/output"
)

const (
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

	// TODO: completely remove the engine and use all of these separately, in a
	// much more composable and testable manner
	ExecutionScheduler lib.ExecutionScheduler
	MetricsEngine      *engine.MetricsEngine
	OutputManager      *output.Manager

	runtimeOptions lib.RuntimeOptions

	ingester output.Output

	logger   *logrus.Entry
	stopOnce sync.Once
	stopChan chan struct***REMOVED******REMOVED***

	Samples chan metrics.SampleContainer

	// Are thresholds tainted?
	thresholdsTaintedLock sync.Mutex
	thresholdsTainted     bool
***REMOVED***

// NewEngine instantiates a new Engine, without doing any heavy initialization.
func NewEngine(testState *lib.TestRunState, ex lib.ExecutionScheduler, outputs []output.Output) (*Engine, error) ***REMOVED***
	if ex == nil ***REMOVED***
		return nil, errors.New("missing ExecutionScheduler instance")
	***REMOVED***

	e := &Engine***REMOVED***
		ExecutionScheduler: ex,

		runtimeOptions: testState.RuntimeOptions,
		Samples:        make(chan metrics.SampleContainer, testState.Options.MetricSamplesBufferSize.Int64),
		stopChan:       make(chan struct***REMOVED******REMOVED***),
		logger:         testState.Logger.WithField("component", "engine"),
	***REMOVED***

	me, err := engine.NewMetricsEngine(ex.GetState())
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	e.MetricsEngine = me

	if !(testState.RuntimeOptions.NoSummary.Bool && testState.RuntimeOptions.NoThresholds.Bool) ***REMOVED***
		e.ingester = me.GetIngester()
		outputs = append(outputs, e.ingester)
	***REMOVED***

	e.OutputManager = output.NewManager(outputs, testState.Logger, func(err error) ***REMOVED***
		if err != nil ***REMOVED***
			testState.Logger.WithError(err).Error("Received error to stop from output")
		***REMOVED***
		e.Stop()
	***REMOVED***)

	return e, nil
***REMOVED***

// Init is used to initialize the execution scheduler and all metrics processing
// in the engine. The first is a costly operation, since it initializes all of
// the planned VUs and could potentially take a long time.
//
// This method either returns an error immediately, or it returns test run() and
// wait() functions.
//
// Things to note:
//   - The first lambda, Run(), synchronously executes the actual load test.
//   - It can be prematurely aborted by cancelling the runCtx - this won't stop
//     the metrics collection by the Engine.
//   - Stopping the metrics collection can be done at any time after Run() has
//     returned by cancelling the globalCtx
//   - The second returned lambda can be used to wait for that process to finish.
func (e *Engine) Init(globalCtx, runCtx context.Context, client *redis.Client) (run func() error, wait func(), err error) ***REMOVED***
	e.logger.Debug("Initialization starting...")
	// TODO: if we ever need metrics processing in the init context, we can move
	// this below the other components... or even start them concurrently?
	if err := e.ExecutionScheduler.Init(runCtx, e.Samples, client); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	// TODO: move all of this in a separate struct? see main TODO above
	runSubCtx, runSubCancel := context.WithCancel(runCtx)

	resultCh := make(chan error)
	processMetricsAfterRun := make(chan struct***REMOVED******REMOVED***)
	runFn := func() error ***REMOVED***
		e.logger.Debug("Execution scheduler starting...")
		err := e.ExecutionScheduler.Run(globalCtx, runSubCtx, e.Samples, client)
		e.logger.WithError(err).Debug("Execution scheduler terminated")

		select ***REMOVED***
		case <-runSubCtx.Done():
			// do nothing, the test run was aborted somehow
		default:
			resultCh <- err // we finished normally, so send the result
		***REMOVED***

		// Make the background jobs process the currently buffered metrics and
		// run the thresholds, then wait for that to be done.
		select ***REMOVED***
		case processMetricsAfterRun <- struct***REMOVED******REMOVED******REMOVED******REMOVED***:
			<-processMetricsAfterRun
		case <-globalCtx.Done():
		***REMOVED***

		return err
	***REMOVED***

	waitFn := e.startBackgroundProcesses(globalCtx, runCtx, resultCh, runSubCancel, processMetricsAfterRun)
	return runFn, waitFn, nil
***REMOVED***

// This starts a bunch of goroutines to process metrics, thresholds, and set the
// test run status when it ends. It returns a function that can be used after
// the provided context is called, to wait for the complete winding down of all
// started goroutines.
//
// Because the background process is not aware of the execution's state, `processMetricsAfterRun`
// will be used to signal that the test run is finished, no more metric samples will be produced,
// and that the remaining metrics samples in the pipeline should be processed as the background
// process is about to exit.
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
					e.OutputManager.SetRunStatus(lib.RunStatusAbortedScriptError)
				case errext.IsInterruptError(err):
					e.OutputManager.SetRunStatus(lib.RunStatusAbortedUser)
				default:
					e.OutputManager.SetRunStatus(lib.RunStatusAbortedSystem)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				e.logger.Debug("run: execution scheduler terminated")
				e.OutputManager.SetRunStatus(lib.RunStatusFinished)
			***REMOVED***
		case <-runCtx.Done():
			e.logger.Debug("run: context expired; exiting...")
			e.OutputManager.SetRunStatus(lib.RunStatusAbortedUser)
		case <-e.stopChan:
			runSubCancel()
			e.logger.Debug("run: stopped by user; exiting...")
			e.OutputManager.SetRunStatus(lib.RunStatusAbortedUser)
		case <-thresholdAbortChan:
			e.logger.Debug("run: stopped by thresholds; exiting...")
			runSubCancel()
			e.OutputManager.SetRunStatus(lib.RunStatusAbortedThreshold)
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
					thresholdsTainted, shouldAbort := e.MetricsEngine.EvaluateThresholds(true)
					e.thresholdsTaintedLock.Lock()
					e.thresholdsTainted = thresholdsTainted
					e.thresholdsTaintedLock.Unlock()
					if shouldAbort ***REMOVED***
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

// processMetrics process the execution's metrics samples as they are collected.
// The processing of samples happens at a fixed rate defined by the `collectRate`
// constant.
//
// The `processMetricsAfterRun` channel argument is used by the caller to signal
// that the test run is finished, no more metric samples will be produced, and that
// the metrics samples remaining in the pipeline should be should be processed.
func (e *Engine) processMetrics(globalCtx context.Context, processMetricsAfterRun chan struct***REMOVED******REMOVED***) ***REMOVED***
	sampleContainers := []metrics.SampleContainer***REMOVED******REMOVED***

	defer func() ***REMOVED***
		// Process any remaining metrics in the pipeline, by this point Run()
		// has already finished and nothing else should be producing metrics.
		e.logger.Debug("Metrics processing winding down...")

		close(e.Samples)
		for sc := range e.Samples ***REMOVED***
			sampleContainers = append(sampleContainers, sc)
		***REMOVED***
		e.OutputManager.AddMetricSamples(sampleContainers)

		if !e.runtimeOptions.NoThresholds.Bool ***REMOVED***
			// Process the thresholds one final time
			thresholdsTainted, _ := e.MetricsEngine.EvaluateThresholds(false)
			e.thresholdsTaintedLock.Lock()
			e.thresholdsTainted = thresholdsTainted
			e.thresholdsTaintedLock.Unlock()
		***REMOVED***
	***REMOVED***()

	ticker := time.NewTicker(collectRate)
	defer ticker.Stop()

	e.logger.Debug("Metrics processing started...")
	processSamples := func() ***REMOVED***
		if len(sampleContainers) > 0 ***REMOVED***
			e.OutputManager.AddMetricSamples(sampleContainers)
			// Make the new container with the same size as the previous
			// one, assuming that we produce roughly the same amount of
			// metrics data between ticks...
			sampleContainers = make([]metrics.SampleContainer, 0, cap(sampleContainers))
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
				// Ensure the ingester flushes any buffered metrics
				_ = e.ingester.Stop()
				thresholdsTainted, _ := e.MetricsEngine.EvaluateThresholds(false)
				e.thresholdsTaintedLock.Lock()
				e.thresholdsTainted = thresholdsTainted
				e.thresholdsTaintedLock.Unlock()
			***REMOVED***
			processMetricsAfterRun <- struct***REMOVED******REMOVED******REMOVED******REMOVED***

		case sc := <-e.Samples:
			sampleContainers = append(sampleContainers, sc)
		case <-globalCtx.Done():
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *Engine) IsTainted() bool ***REMOVED***
	e.thresholdsTaintedLock.Lock()
	defer e.thresholdsTaintedLock.Unlock()
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
