package engineold

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/APITeamLimited/globe-test/metrics"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/output"
	"github.com/sirupsen/logrus"
)

const (
	collectRate    = 1000 * time.Millisecond
	thresholdsRate = 2 * time.Second
)

type Engine struct {
	ExecutionScheduler libWorker.ExecutionScheduler
	OutputManager      *output.Manager

	ingester output.Output

	logger   *logrus.Entry
	stopOnce sync.Once
	stopChan chan struct{}

	Samples chan metrics.SampleContainer

	// Are thresholds tainted?
	thresholdsTaintedLock sync.Mutex
	thresholdsTainted     bool
}

// NewEngine instantiates a new Engine, without doing any heavy initialization.
func NewEngine(testState *libWorker.TestRunState, ex libWorker.ExecutionScheduler, outputs []output.Output) (*Engine, error) {
	if ex == nil {
		return nil, errors.New("missing ExecutionScheduler instance")
	}

	e := &Engine{
		ExecutionScheduler: ex,
		Samples:            make(chan metrics.SampleContainer, testState.Options.MetricSamplesBufferSize.Int64),
		stopChan:           make(chan struct{}),
		logger:             testState.Logger.WithField("component", "engine"),
	}

	e.ingester = me.GetIngester()
	outputs = append(outputs, e.ingester)

	e.OutputManager = output.NewManager(outputs, testState.Logger, func(err error) {
		if err != nil {
			testState.Logger.WithError(err).Error("Received error to stop from output")
		}
		e.Stop()
	})

	return e, nil
}

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
func (e *Engine) Init(globalCtx, runCtx context.Context, workerInfo *libWorker.WorkerInfo) (run func() error, wait func(), err error) {
	e.logger.Debug("Initialization starting...")
	// TODO: if we ever need metrics processing in the init context, we can move
	// this below the other components... or even start them concurrently?
	if err := e.ExecutionScheduler.Init(runCtx, e.Samples, workerInfo); err != nil {
		return nil, nil, err
	}

	// TODO: move all of this in a separate struct? see main TODO above
	runSubCtx, runSubCancel := context.WithCancel(runCtx)

	resultCh := make(chan error)
	processMetricsAfterRun := make(chan struct{})
	runFn := func() error {
		e.logger.Debug("Execution scheduler starting...")
		err := e.ExecutionScheduler.Run(globalCtx, runSubCtx, e.Samples, workerInfo)
		e.logger.WithError(err).Debug("Execution scheduler terminated")

		select {
		case <-runSubCtx.Done():
			// do nothing, the test run was aborted somehow
		default:
			resultCh <- err // we finished normally, so send the result
		}

		// Make the background jobs process the currently buffered metrics and
		// run the thresholds, then wait for that to be done.
		select {
		case processMetricsAfterRun <- struct{}{}:
			<-processMetricsAfterRun
		case <-globalCtx.Done():
		}

		return err
	}

	waitFn := e.startBackgroundProcesses(globalCtx, runCtx, resultCh, runSubCancel, processMetricsAfterRun)
	return runFn, waitFn, nil
}

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
	globalCtx, runCtx context.Context, runResult <-chan error, runSubCancel func(), processMetricsAfterRun chan struct{},
) (wait func()) {
	processes := new(sync.WaitGroup)

	// Siphon and handle all produced metric samples
	processes.Add(1)
	go func() {
		defer processes.Done()
		e.processMetrics(globalCtx, processMetricsAfterRun)
	}()

	// Update the test run status when the test finishes
	processes.Add(1)
	thresholdAbortChan := make(chan struct{})
	go func() {
		defer processes.Done()
		select {
		case err := <-runResult:
			if err != nil {
				e.logger.WithError(err).Debug("run: execution scheduler returned an error")
			} else {
				e.logger.Debug("run: execution scheduler terminated")
			}
		case <-runCtx.Done():
			e.logger.Debug("run: context expired; exiting...")
		case <-e.stopChan:
			runSubCancel()
			e.logger.Debug("run: stopped by user; exiting...")
		case <-thresholdAbortChan:
			e.logger.Debug("run: stopped by thresholds; exiting...")
			runSubCancel()
		}
	}()

	// Run thresholds, if not disabled.
	processes.Add(1)
	go func() {
		defer processes.Done()
		defer e.logger.Debug("Engine: Thresholds terminated")
		ticker := time.NewTicker(thresholdsRate)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				thresholdsTainted, shouldAbort := e.MetricsEngine.EvaluateThresholds(true)
				e.thresholdsTaintedLock.Lock()
				e.thresholdsTainted = thresholdsTainted
				e.thresholdsTaintedLock.Unlock()
				if shouldAbort {
					close(thresholdAbortChan)
					return
				}
			case <-runCtx.Done():
				return
			}
		}
	}()

	return processes.Wait
}

// processMetrics process the execution's metrics samples as they are collected.
// The processing of samples happens at a fixed rate defined by the `collectRate`
// constant.
//
// The `processMetricsAfterRun` channel argument is used by the caller to signal
// that the test run is finished, no more metric samples will be produced, and that
// the metrics samples remaining in the pipeline should be should be processed.
func (e *Engine) processMetrics(globalCtx context.Context, processMetricsAfterRun chan struct{}) {
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

	e.logger.Debug("Metrics processing started...")
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

func (e *Engine) IsTainted() bool {
	e.thresholdsTaintedLock.Lock()
	defer e.thresholdsTaintedLock.Unlock()
	return e.thresholdsTainted
}

// Stop closes a signal channel, forcing a running Engine to return
func (e *Engine) Stop() {
	e.stopOnce.Do(func() {
		close(e.stopChan)
	})
}

// IsStopped returns a bool indicating whether the Engine has been stopped
func (e *Engine) IsStopped() bool {
	select {
	case <-e.stopChan:
		return true
	default:
		return false
	}
}
