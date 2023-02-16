package core

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/APITeamLimited/globe-test/metrics"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/sirupsen/logrus"
)

const (
	collectRate = 1000 * time.Millisecond
)

type Engine struct {
	ExecutionScheduler libWorker.ExecutionScheduler
	logger             *logrus.Entry
	stopOnce           sync.Once
	stopChan           chan struct{}

	Samples chan metrics.SampleContainer

	// Are thresholds tainted?
	thresholdsTaintedLock sync.Mutex
	thresholdsTainted     bool
}

// NewEngine instantiates a new Engine, without doing any heavy initialization.
func NewEngine(testState *libWorker.TestRunState, ex libWorker.ExecutionScheduler, samplesChan chan metrics.SampleContainer) (*Engine, error) {
	if ex == nil {
		return nil, errors.New("missing ExecutionScheduler instance")
	}

	e := &Engine{
		ExecutionScheduler: ex,
		Samples:            samplesChan,
		stopChan:           make(chan struct{}),
		logger:             testState.Logger.WithField("component", "engine"),
	}
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
	if err := e.ExecutionScheduler.Init(runCtx, e.Samples, workerInfo); err != nil {
		return nil, nil, err
	}

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

	//waitFn := e.startBackgroundProcesses(globalCtx, runCtx, resultCh, runSubCancel, processMetricsAfterRun)
	return runFn, func() {}, nil
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
/*func (e *Engine) startBackgroundProcesses(
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
}*/

// // Stop closes a signal channel, forcing a running Engine to return
// func (e *Engine) Stop() {
// 	e.stopOnce.Do(func() {
// 		close(e.stopChan)
// 	})
// }

// // IsStopped returns a bool indicating whether the Engine has been stopped
// func (e *Engine) IsStopped() bool {
// 	select {
// 	case <-e.stopChan:
// 		return true
// 	default:
// 		return false
// 	}
// }
