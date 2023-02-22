package local

import (
	"context"
	"fmt"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/APITeamLimited/globe-test/worker/metrics"
	"github.com/sirupsen/logrus"

	"github.com/APITeamLimited/globe-test/worker/errext"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/executor"
)

// ExecutionScheduler is the local implementation of libWorker.ExecutionScheduler
type ExecutionScheduler struct {
	executorConfigs []libWorker.ExecutorConfig // sorted by (startTime, ID)
	executors       []libWorker.Executor       // sorted by (startTime, ID), excludes executors with no work
	executionPlan   []libWorker.ExecutionStep
	maxDuration     time.Duration // cached value derived from the execution plan
	maxPossibleVUs  uint64        // cached value derived from the execution plan
	state           *libWorker.ExecutionState

	// TODO: remove these when we don't have separate Init() and Run() methods
	// and can use a context + a WaitGroup (or something like that)
	stopVUsEmission, vusEmissionStopped chan struct{}
}

// Check to see if we implement the libWorker.ExecutionScheduler interface
var _ libWorker.ExecutionScheduler = &ExecutionScheduler{}

// NewExecutionScheduler creates and returns a new local libWorker.ExecutionScheduler
// instance, without initializing it beyond the bare minimum. Specifically, it
// creates the needed executor instances and a lot of state placeholders, but it
// doesn't initialize the executors and it doesn't initialize or run VUs.
func NewExecutionScheduler(trs *libWorker.TestRunState) (*ExecutionScheduler, error) {
	options := trs.Options
	et, err := libWorker.NewExecutionTuple(options.ExecutionSegment, options.ExecutionSegmentSequence)
	if err != nil {
		return nil, err
	}
	executionPlan := options.Scenarios.GetFullExecutionRequirements(et)
	maxPlannedVUs := libWorker.GetMaxPlannedVUs(executionPlan)
	maxPossibleVUs := libWorker.GetMaxPossibleVUs(executionPlan)

	executionState := libWorker.NewExecutionState(trs, et, maxPlannedVUs, maxPossibleVUs)
	maxDuration, _ := libWorker.GetEndOffset(executionPlan) // we don't care if the end offset is final

	executorConfigs := options.Scenarios.GetSortedConfigs()

	executors := make([]libWorker.Executor, 0, len(executorConfigs))
	// Only take executors which have work.
	for _, sc := range executorConfigs {
		if !sc.HasWork(et) {
			trs.Logger.Warnf(
				"Executor '%s' is disabled for segment %s due to lack of work!",
				sc.GetName(), options.ExecutionSegment,
			)
			continue
		}
		s, err := sc.NewExecutor(executionState, trs.Logger.WithFields(logrus.Fields{
			"scenario": sc.GetName(),
			"executor": sc.GetType(),
		}))
		if err != nil {
			return nil, err
		}
		executors = append(executors, s)
	}

	if options.Paused.Bool {
		if err := executionState.Pause(); err != nil {
			return nil, err
		}
	}

	return &ExecutionScheduler{
		executors:       executors,
		executorConfigs: executorConfigs,
		executionPlan:   executionPlan,
		maxDuration:     maxDuration,
		maxPossibleVUs:  maxPossibleVUs,
		state:           executionState,

		stopVUsEmission:    make(chan struct{}),
		vusEmissionStopped: make(chan struct{}),
	}, nil
}

// GetRunner returns the wrapped libWorker.Runner instance.
func (e *ExecutionScheduler) GetRunner() libWorker.Runner { // TODO: remove
	return e.state.Test.Runner
}

// GetState returns a pointer to the execution state struct for the local
// execution scheduler. It's guaranteed to be initialized and present, though
// see the documentation in lib/execution.go for caveats about its usage. The
// most important one is that none of the methods beyond the pause-related ones
// should be used for synchronization.
func (e *ExecutionScheduler) GetState() *libWorker.ExecutionState {
	return e.state
}

// GetExecutors returns the slice of configured executor instances which
// have work, sorted by their (startTime, name) in an ascending order.
func (e *ExecutionScheduler) GetExecutors() []libWorker.Executor {
	return e.executors
}

// GetExecutorConfigs returns the slice of all executor configs, sorted by
// their (startTime, name) in an ascending order.
func (e *ExecutionScheduler) GetExecutorConfigs() []libWorker.ExecutorConfig {
	return e.executorConfigs
}

// GetExecutionPlan is a helper method so users of the local execution scheduler
// don't have to calculate the execution plan again.
func (e *ExecutionScheduler) GetExecutionPlan() []libWorker.ExecutionStep {
	return e.executionPlan
}

// initVU is a helper method that's used to both initialize the planned VUs
// in the Init() method, and also passed to executors so they can initialize
// any unplanned VUs themselves.
func (e *ExecutionScheduler) initVU(
	samplesOut chan<- metrics.SampleContainer, logger logrus.FieldLogger, workerInfo *libWorker.WorkerInfo,
) (libWorker.InitializedVU, error) {
	// Get the VU IDs here, so that the VUs are (mostly) ordered by their
	// number in the channel buffer
	vuIDLocal, vuIDGlobal := e.state.GetUniqueVUIdentifiers()
	vu, err := e.state.Test.Runner.NewVU(vuIDLocal, vuIDGlobal, samplesOut, workerInfo)
	if err != nil {
		return nil, errext.WithHint(err, fmt.Sprintf("error while initializing VU #%d", vuIDGlobal))
	}

	logger.Debugf("Initialized VU #%d", vuIDGlobal)
	return vu, nil
}

func (e *ExecutionScheduler) initVUsConcurrently(
	ctx context.Context, samplesOut chan<- metrics.SampleContainer, count uint64,
	concurrency int, logger logrus.FieldLogger,
	workerInfo *libWorker.WorkerInfo,
) chan error {
	doneInits := make(chan error, count) // poor man's early-return waitgroup
	limiter := make(chan struct{})

	for i := 0; i < concurrency; i++ {
		go func() {
			for range limiter {
				newVU, err := e.initVU(samplesOut, logger, workerInfo)
				if err == nil {
					e.state.AddInitializedVU(newVU)
				}
				doneInits <- err
			}
		}()
	}

	go func() {
		defer close(limiter)
		for vuNum := uint64(0); vuNum < count; vuNum++ {
			select {
			case limiter <- struct{}{}:
			case <-ctx.Done():
				return
			}
		}
	}()

	return doneInits
}

func (e *ExecutionScheduler) emitVUsAndVUsMax(ctx context.Context, out chan<- metrics.SampleContainer) {
	e.state.Test.Logger.Debug("Starting emission of VUs and VUsMax metrics...")
	runTags := metrics.NewSampleTags(e.state.Test.Options.RunTags)

	emitMetrics := func() {
		t := time.Now()
		samples := metrics.ConnectedSamples{
			Samples: []metrics.Sample{
				{
					Time:   t,
					Metric: e.state.Test.BuiltinMetrics.VUs,
					Value:  float64(e.state.GetCurrentlyActiveVUsCount()),
					Tags:   runTags,
				}, {
					Time:   t,
					Metric: e.state.Test.BuiltinMetrics.VUsMax,
					Value:  float64(e.state.GetInitializedVUsCount()),
					Tags:   runTags,
				},
			},
			Tags: runTags,
			Time: t,
		}
		metrics.PushIfNotDone(ctx, out, samples)
	}

	ticker := time.NewTicker(1 * time.Second)
	go func() {
		defer func() {
			ticker.Stop()
			e.state.Test.Logger.Debug("Metrics emission of VUs and VUsMax metrics stopped")
			close(e.vusEmissionStopped)
		}()

		for {
			select {
			case <-ticker.C:
				emitMetrics()
			case <-ctx.Done():
				return
			case <-e.stopVUsEmission:
				return
			}
		}
	}()
}

// Init concurrently initializes all of the planned VUs and then sequentially
// initializes all of the configured executors.
func (e *ExecutionScheduler) Init(ctx context.Context, samplesOut chan<- metrics.SampleContainer, workerInfo *libWorker.WorkerInfo) error {
	e.emitVUsAndVUsMax(ctx, samplesOut)

	logger := e.state.Test.Logger.WithField("phase", "local-execution-scheduler-init")
	vusToInitialize := libWorker.GetMaxPlannedVUs(e.executionPlan)
	logger.WithFields(logrus.Fields{
		"neededVUs":      vusToInitialize,
		"executorsCount": len(e.executors),
	}).Debugf("Start of initialization")

	subctx, cancel := context.WithCancel(ctx)
	defer cancel()

	e.state.SetExecutionStatus(libWorker.ExecutionStatusInitVUs)
	doneInits := e.initVUsConcurrently(subctx, samplesOut, vusToInitialize, runtime.GOMAXPROCS(0), logger, workerInfo)

	initializedVUs := new(uint64)

	for vuNum := uint64(0); vuNum < vusToInitialize; vuNum++ {
		select {
		case err := <-doneInits:
			if err != nil {
				logger.WithError(err).Debug("VU initialization returned with an error, aborting...")
				// the context's cancel() is called in a defer above and will
				// abort any in-flight VU initializations
				return err
			}
			atomic.AddUint64(initializedVUs, 1)
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	e.state.SetInitVUFunc(func(ctx context.Context, logger *logrus.Entry, workerInfo *libWorker.WorkerInfo) (libWorker.InitializedVU, error) {
		return e.initVU(samplesOut, logger, workerInfo)
	})

	e.state.SetExecutionStatus(libWorker.ExecutionStatusInitExecutors)
	logger.Debugf("Finished initializing needed VUs, start initializing executors...")
	for _, exec := range e.executors {
		executorConfig := exec.GetConfig()

		if err := exec.Init(ctx); err != nil {
			return fmt.Errorf("error while initializing executor %s: %w", executorConfig.GetName(), err)
		}
		logger.Debugf("Initialized executor %s", executorConfig.GetName())
	}

	e.state.SetExecutionStatus(libWorker.ExecutionStatusInitDone)
	logger.Debugf("Initialization completed")
	return nil
}

// runExecutor gets called by the public Run() method once per configured
// executor, each time in a new goroutine. It is responsible for waiting out the
// configured startTime for the specific executor and then running its Run()
// method.
func (e *ExecutionScheduler) runExecutor(
	runCtx context.Context, runResults chan<- error, engineOut chan<- metrics.SampleContainer, executor libWorker.Executor, workerInfo *libWorker.WorkerInfo,
) {
	executorConfig := executor.GetConfig()
	executorStartTime := executorConfig.GetStartTime()
	executorLogger := e.state.Test.Logger.WithFields(logrus.Fields{
		"executor":  executorConfig.GetName(),
		"type":      executorConfig.GetType(),
		"startTime": executorStartTime,
	})

	// Check if we have to wait before starting the actual executor execution
	if executorStartTime > 0 {
		executorLogger.Debugf("Waiting for executor start time...")
		select {
		case <-runCtx.Done():
			runResults <- nil // no error since executor hasn't started yet
			return
		case <-time.After(executorStartTime):
			// continue
		}
	}

	executorLogger.Debugf("Starting executor")
	err := executor.Run(runCtx, engineOut, workerInfo) // executor should handle context cancel itself
	if err == nil {
		executorLogger.Debugf("Executor finished successfully")
	} else {
		executorLogger.WithField("error", err).Errorf("Executor error")
	}
	runResults <- err
}

// Run the ExecutionScheduler, funneling all generated metric samples through the supplied
// out channel.
//
//nolint:funlen
func (e *ExecutionScheduler) Run(globalCtx, runCtx context.Context, engineOut chan<- metrics.SampleContainer, workerInfo *libWorker.WorkerInfo) error {
	defer func() {
		close(e.stopVUsEmission)
		<-e.vusEmissionStopped
	}()

	executorsCount := len(e.executors)
	logger := e.state.Test.Logger.WithField("phase", "local-execution-scheduler-run")
	var interrupted bool
	defer func() {
		e.state.MarkEnded()
		if interrupted {
			e.state.SetExecutionStatus(libWorker.ExecutionStatusInterrupted)
		}
	}()

	if e.state.IsPaused() {
		logger.Debug("Execution is paused, waiting for resume or interrupt...")
		e.state.SetExecutionStatus(libWorker.ExecutionStatusPausedBeforeRun)
		select {
		case <-e.state.ResumeNotify():
			// continue
		case <-runCtx.Done():
			return nil
		}
	}

	e.state.MarkStarted()

	logger.WithFields(logrus.Fields{"executorsCount": executorsCount}).Debugf("Start of test run")

	runResults := make(chan error, executorsCount) // nil values are successful runs

	runCtx = libWorker.WithExecutionState(runCtx, e.state)
	runSubCtx, cancel := context.WithCancel(runCtx)
	defer cancel() // just in case, and to shut up go vet...

	// Run setup() before any executors, if it's not disabled
	if !e.state.Test.Options.NoSetup.Bool {
		logger.Debug("Running setup()")
		e.state.SetExecutionStatus(libWorker.ExecutionStatusSetup)
		if err := e.state.Test.Runner.Setup(runSubCtx, engineOut, workerInfo); err != nil {
			logger.WithField("error", err).Debug("setup() aborted by error")
			return err
		}
	}

	// Start all executors at their particular startTime in a separate goroutine...
	logger.Debug("Start all executors...")
	e.state.SetExecutionStatus(libWorker.ExecutionStatusRunning)

	// We are using this context to allow libWorker.Executor implementations to cancel
	// this context effectively stopping all executions.
	//
	// This is for addressing test.abort().
	execCtx := executor.Context(runSubCtx)
	for _, exec := range e.executors {
		go e.runExecutor(execCtx, runResults, engineOut, exec, workerInfo)
	}

	// Wait for all executors to finish
	var firstErr error
	for range e.executors {
		err := <-runResults
		if err != nil && firstErr == nil {
			logger.WithError(err).Debug("Executor returned with an error, cancelling test run...")
			firstErr = err
			cancel()
		}
	}

	// Run teardown() after all executors are done, if it's not disabled
	if !e.state.Test.Options.NoTeardown.Bool {
		logger.Debug("Running teardown()")
		e.state.SetExecutionStatus(libWorker.ExecutionStatusTeardown)

		// We run teardown() with the global context, so it isn't interrupted by
		// aborts caused by thresholds or even Ctrl+C (unless used twice).
		if err := e.state.Test.Runner.Teardown(globalCtx, engineOut, workerInfo); err != nil {
			logger.WithField("error", err).Debug("teardown() aborted by error")
			return err
		}
	}
	if err := executor.CancelReason(execCtx); err != nil && errext.IsInterruptError(err) {
		interrupted = true
		return err
	}
	return firstErr
}

// SetPaused pauses a test, if called with true. And if called with false, tries
// to start/resume it. See the libWorker.ExecutionScheduler interface documentation of
// the methods for the various caveats about its usage.
func (e *ExecutionScheduler) SetPaused(pause bool) error {
	if !e.state.HasStarted() && e.state.IsPaused() {
		if pause {
			return fmt.Errorf("execution is already paused")
		}
		e.state.Test.Logger.Debug("Starting execution")
		return e.state.Resume()
	}

	for _, exec := range e.executors {
		pausableExecutor, ok := exec.(libWorker.PausableExecutor)
		if !ok {
			return fmt.Errorf(
				"%s executor '%s' doesn't support pause and resume operations after its start",
				exec.GetConfig().GetType(), exec.GetConfig().GetName(),
			)
		}
		if err := pausableExecutor.SetPaused(pause); err != nil {
			return err
		}
	}
	if pause {
		return e.state.Pause()
	}
	return e.state.Resume()
}
