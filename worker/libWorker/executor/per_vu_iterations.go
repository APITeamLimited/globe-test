package executor

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
	"github.com/APITeamLimited/globe-test/worker/metrics"
	"github.com/APITeamLimited/globe-test/worker/pb"
)

const perVUIterationsType = "per-vu-iterations"

func init() ***REMOVED***
	libWorker.RegisterExecutorConfigType(perVUIterationsType, func(name string, rawJSON []byte) (libWorker.ExecutorConfig, error) ***REMOVED***
		config := NewPerVUIterationsConfig(name)
		err := libWorker.StrictJSONUnmarshal(rawJSON, &config)
		return config, err
	***REMOVED***)
***REMOVED***

// PerVUIterationsConfig stores the number of VUs iterations, as well as maxDuration settings
type PerVUIterationsConfig struct ***REMOVED***
	BaseConfig
	VUs         null.Int           `json:"vus"`
	Iterations  null.Int           `json:"iterations"`
	MaxDuration types.NullDuration `json:"maxDuration"`
***REMOVED***

// NewPerVUIterationsConfig returns a PerVUIterationsConfig with default values
func NewPerVUIterationsConfig(name string) PerVUIterationsConfig ***REMOVED***
	return PerVUIterationsConfig***REMOVED***
		BaseConfig:  NewBaseConfig(name, perVUIterationsType),
		VUs:         null.NewInt(1, false),
		Iterations:  null.NewInt(1, false),
		MaxDuration: types.NewNullDuration(10*time.Minute, false), // TODO: shorten?
	***REMOVED***
***REMOVED***

// Make sure we implement the libWorker.ExecutorConfig interface
var _ libWorker.ExecutorConfig = &PerVUIterationsConfig***REMOVED******REMOVED***

// GetVUs returns the scaled VUs for the executor.
func (pvic PerVUIterationsConfig) GetVUs(et *libWorker.ExecutionTuple) int64 ***REMOVED***
	return et.ScaleInt64(pvic.VUs.Int64)
***REMOVED***

// GetIterations returns the UNSCALED iteration count for the executor. It's
// important to note that scaling per-VU iteration executor affects only the
// number of VUs. If we also scaled the iterations, scaling would have quadratic
// effects instead of just linear.
func (pvic PerVUIterationsConfig) GetIterations() int64 ***REMOVED***
	return pvic.Iterations.Int64
***REMOVED***

// GetDescription returns a human-readable description of the executor options
func (pvic PerVUIterationsConfig) GetDescription(et *libWorker.ExecutionTuple) string ***REMOVED***
	return fmt.Sprintf("%d iterations for each of %d VUs%s",
		pvic.GetIterations(), pvic.GetVUs(et),
		pvic.getBaseInfo(fmt.Sprintf("maxDuration: %s", pvic.MaxDuration.Duration)))
***REMOVED***

// Validate makes sure all options are configured and valid
func (pvic PerVUIterationsConfig) Validate() []error ***REMOVED***
	errors := pvic.BaseConfig.Validate()
	if pvic.VUs.Int64 <= 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("the number of VUs must be more than 0"))
	***REMOVED***

	if pvic.Iterations.Int64 <= 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("the number of iterations must be more than 0"))
	***REMOVED***

	if pvic.MaxDuration.TimeDuration() < minDuration ***REMOVED***
		errors = append(errors, fmt.Errorf(
			"the maxDuration must be at least %s, but is %s", minDuration, pvic.MaxDuration,
		))
	***REMOVED***

	return errors
***REMOVED***

// GetExecutionRequirements returns the number of required VUs to run the
// executor for its whole duration (disregarding any startTime), including the
// maximum waiting time for any iterations to gracefully stop. This is used by
// the execution scheduler in its VU reservation calculations, so it knows how
// many VUs to pre-initialize.
func (pvic PerVUIterationsConfig) GetExecutionRequirements(et *libWorker.ExecutionTuple) []libWorker.ExecutionStep ***REMOVED***
	return []libWorker.ExecutionStep***REMOVED***
		***REMOVED***
			TimeOffset: 0,
			PlannedVUs: uint64(pvic.GetVUs(et)),
		***REMOVED***,
		***REMOVED***
			TimeOffset: pvic.MaxDuration.TimeDuration() + pvic.GracefulStop.TimeDuration(),
			PlannedVUs: 0,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// NewExecutor creates a new PerVUIterations executor
func (pvic PerVUIterationsConfig) NewExecutor(
	es *libWorker.ExecutionState, logger *logrus.Entry,
) (libWorker.Executor, error) ***REMOVED***
	return PerVUIterations***REMOVED***
		BaseExecutor: NewBaseExecutor(pvic, es, logger),
		config:       pvic,
	***REMOVED***, nil
***REMOVED***

// HasWork reports whether there is any work to be done for the given execution segment.
func (pvic PerVUIterationsConfig) HasWork(et *libWorker.ExecutionTuple) bool ***REMOVED***
	return pvic.GetVUs(et) > 0 && pvic.GetIterations() > 0
***REMOVED***

// PerVUIterations executes a specific number of iterations with each VU.
type PerVUIterations struct ***REMOVED***
	*BaseExecutor
	config PerVUIterationsConfig
***REMOVED***

// Make sure we implement the libWorker.Executor interface.
var _ libWorker.Executor = &PerVUIterations***REMOVED******REMOVED***

// Run executes a specific number of iterations with each configured VU.
//
//nolint:funlen
func (pvi PerVUIterations) Run(parentCtx context.Context, out chan<- metrics.SampleContainer, workerInfo *libWorker.WorkerInfo) (err error) ***REMOVED***
	numVUs := pvi.config.GetVUs(pvi.executionState.ExecutionTuple)
	iterations := pvi.config.GetIterations()
	duration := pvi.config.MaxDuration.TimeDuration()
	gracefulStop := pvi.config.GetGracefulStop()

	waitOnProgressChannel := make(chan struct***REMOVED******REMOVED***)
	startTime, maxDurationCtx, regDurationCtx, cancel := getDurationContexts(parentCtx, duration, gracefulStop)
	defer func() ***REMOVED***
		cancel()
		<-waitOnProgressChannel
	***REMOVED***()

	// Make sure the log and the progress bar have accurate information
	pvi.logger.WithFields(logrus.Fields***REMOVED***
		"vus": numVUs, "iterations": iterations, "maxDuration": duration, "type": pvi.config.GetType(),
	***REMOVED***).Debug("Starting executor run...")

	totalIters := uint64(numVUs * iterations)
	doneIters := new(uint64)

	vusFmt := pb.GetFixedLengthIntFormat(numVUs)
	itersFmt := pb.GetFixedLengthIntFormat(int64(totalIters))
	progressFn := func() (float64, []string) ***REMOVED***
		spent := time.Since(startTime)
		progVUs := fmt.Sprintf(vusFmt+" VUs", numVUs)
		currentDoneIters := atomic.LoadUint64(doneIters)
		progIters := fmt.Sprintf(itersFmt+"/"+itersFmt+" iters, %d per VU",
			currentDoneIters, totalIters, iterations)
		right := []string***REMOVED***progVUs, duration.String(), progIters***REMOVED***
		if spent > duration ***REMOVED***
			return 1, right
		***REMOVED***

		spentDuration := pb.GetFixedLengthDuration(spent, duration)
		progDur := fmt.Sprintf("%s/%s", spentDuration, duration)
		right[1] = progDur

		return float64(currentDoneIters) / float64(totalIters), right
	***REMOVED***
	pvi.progress.Modify(pb.WithProgress(progressFn))

	maxDurationCtx = libWorker.WithScenarioState(maxDurationCtx, &libWorker.ScenarioState***REMOVED***
		Name:       pvi.config.Name,
		Executor:   pvi.config.Type,
		StartTime:  startTime,
		ProgressFn: progressFn,
	***REMOVED***)
	go func() ***REMOVED***
		trackProgress(parentCtx, maxDurationCtx, regDurationCtx, pvi, progressFn)
		close(waitOnProgressChannel)
	***REMOVED***()

	handleVUsWG := &sync.WaitGroup***REMOVED******REMOVED***
	defer handleVUsWG.Wait()
	// Actually schedule the VUs and iterations...
	activeVUs := &sync.WaitGroup***REMOVED******REMOVED***
	defer activeVUs.Wait()

	regDurationDone := regDurationCtx.Done()
	runIteration := getIterationRunner(pvi.executionState, pvi.logger)

	returnVU := func(u libWorker.InitializedVU) ***REMOVED***
		pvi.executionState.ReturnVU(u, true)
		activeVUs.Done()
	***REMOVED***

	droppedIterationMetric := pvi.executionState.Test.BuiltinMetrics.DroppedIterations
	handleVU := func(initVU libWorker.InitializedVU) ***REMOVED***
		defer handleVUsWG.Done()
		ctx, cancel := context.WithCancel(maxDurationCtx)
		defer cancel()

		vuID := initVU.GetID()
		activeVU := initVU.Activate(
			getVUActivationParams(ctx, pvi.config.BaseConfig, returnVU,
				pvi.nextIterationCounters))

		for i := int64(0); i < iterations; i++ ***REMOVED***
			select ***REMOVED***
			case <-regDurationDone:
				metrics.PushIfNotDone(parentCtx, out, metrics.Sample***REMOVED***
					Value: float64(iterations - i), Metric: droppedIterationMetric,
					Tags: pvi.getMetricTags(&vuID), Time: time.Now(),
				***REMOVED***)
				return // don't make more iterations
			default:
				// continue looping
			***REMOVED***
			runIteration(maxDurationCtx, activeVU)
			atomic.AddUint64(doneIters, 1)
		***REMOVED***
	***REMOVED***

	for i := int64(0); i < numVUs; i++ ***REMOVED***
		initializedVU, err := pvi.executionState.GetPlannedVU(pvi.logger, true)
		if err != nil ***REMOVED***
			cancel()
			return err
		***REMOVED***
		activeVUs.Add(1)
		handleVUsWG.Add(1)
		go handleVU(initializedVU)
	***REMOVED***

	return nil
***REMOVED***
