package executor

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"

	"github.com/APITeamLimited/k6-worker/lib"
	"github.com/APITeamLimited/k6-worker/lib/types"
	"github.com/APITeamLimited/k6-worker/metrics"
	"github.com/APITeamLimited/k6-worker/pb"
)

const sharedIterationsType = "shared-iterations"

func init() ***REMOVED***
	lib.RegisterExecutorConfigType(
		sharedIterationsType,
		func(name string, rawJSON []byte) (lib.ExecutorConfig, error) ***REMOVED***
			config := NewSharedIterationsConfig(name)
			err := lib.StrictJSONUnmarshal(rawJSON, &config)
			return config, err
		***REMOVED***,
	)
***REMOVED***

// SharedIterationsConfig stores the number of VUs iterations, as well as maxDuration settings
type SharedIterationsConfig struct ***REMOVED***
	BaseConfig
	VUs         null.Int           `json:"vus"`
	Iterations  null.Int           `json:"iterations"`
	MaxDuration types.NullDuration `json:"maxDuration"`
***REMOVED***

// NewSharedIterationsConfig returns a SharedIterationsConfig with default values
func NewSharedIterationsConfig(name string) SharedIterationsConfig ***REMOVED***
	return SharedIterationsConfig***REMOVED***
		BaseConfig:  NewBaseConfig(name, sharedIterationsType),
		VUs:         null.NewInt(1, false),
		Iterations:  null.NewInt(1, false),
		MaxDuration: types.NewNullDuration(10*time.Minute, false), // TODO: shorten?
	***REMOVED***
***REMOVED***

// Make sure we implement the lib.ExecutorConfig interface
var _ lib.ExecutorConfig = &SharedIterationsConfig***REMOVED******REMOVED***

// GetVUs returns the scaled VUs for the executor.
func (sic SharedIterationsConfig) GetVUs(et *lib.ExecutionTuple) int64 ***REMOVED***
	return et.ScaleInt64(sic.VUs.Int64)
***REMOVED***

// GetIterations returns the scaled iteration count for the executor.
func (sic SharedIterationsConfig) GetIterations(et *lib.ExecutionTuple) int64 ***REMOVED***
	// TODO: Optimize this by probably changing the whole Config API
	newTuple, err := et.GetNewExecutionTupleFromValue(sic.VUs.Int64)
	if err != nil ***REMOVED***
		return 0
	***REMOVED***
	return newTuple.ScaleInt64(sic.Iterations.Int64)
***REMOVED***

// GetDescription returns a human-readable description of the executor options
func (sic SharedIterationsConfig) GetDescription(et *lib.ExecutionTuple) string ***REMOVED***
	return fmt.Sprintf("%d iterations shared among %d VUs%s",
		sic.GetIterations(et), sic.GetVUs(et),
		sic.getBaseInfo(fmt.Sprintf("maxDuration: %s", sic.MaxDuration.Duration)))
***REMOVED***

// Validate makes sure all options are configured and valid
func (sic SharedIterationsConfig) Validate() []error ***REMOVED***
	errors := sic.BaseConfig.Validate()
	if sic.VUs.Int64 <= 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("the number of VUs must be more than 0"))
	***REMOVED***

	if sic.Iterations.Int64 < sic.VUs.Int64 ***REMOVED***
		errors = append(errors, fmt.Errorf(
			"the number of iterations (%d) can't be less than the number of VUs (%d)",
			sic.Iterations.Int64, sic.VUs.Int64,
		))
	***REMOVED***

	if sic.MaxDuration.TimeDuration() < minDuration ***REMOVED***
		errors = append(errors, fmt.Errorf(
			"the maxDuration must be at least %s, but is %s", minDuration, sic.MaxDuration,
		))
	***REMOVED***

	return errors
***REMOVED***

// GetExecutionRequirements returns the number of required VUs to run the
// executor for its whole duration (disregarding any startTime), including the
// maximum waiting time for any iterations to gracefully stop. This is used by
// the execution scheduler in its VU reservation calculations, so it knows how
// many VUs to pre-initialize.
func (sic SharedIterationsConfig) GetExecutionRequirements(et *lib.ExecutionTuple) []lib.ExecutionStep ***REMOVED***
	vus := sic.GetVUs(et)
	if vus == 0 ***REMOVED***
		return []lib.ExecutionStep***REMOVED***
			***REMOVED***
				TimeOffset: 0,
				PlannedVUs: 0,
			***REMOVED***,
		***REMOVED***
	***REMOVED***

	return []lib.ExecutionStep***REMOVED***
		***REMOVED***
			TimeOffset: 0,
			PlannedVUs: uint64(vus),
		***REMOVED***,
		***REMOVED***
			TimeOffset: sic.MaxDuration.TimeDuration() + sic.GracefulStop.TimeDuration(),
			PlannedVUs: 0,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// NewExecutor creates a new SharedIterations executor
func (sic SharedIterationsConfig) NewExecutor(
	es *lib.ExecutionState, logger *logrus.Entry,
) (lib.Executor, error) ***REMOVED***
	return &SharedIterations***REMOVED***
		BaseExecutor: NewBaseExecutor(sic, es, logger),
		config:       sic,
	***REMOVED***, nil
***REMOVED***

// SharedIterations executes a specific total number of iterations, which are
// all shared by the configured VUs.
type SharedIterations struct ***REMOVED***
	*BaseExecutor
	config SharedIterationsConfig
	et     *lib.ExecutionTuple
***REMOVED***

// Make sure we implement the lib.Executor interface.
var _ lib.Executor = &SharedIterations***REMOVED******REMOVED***

// HasWork reports whether there is any work to be done for the given execution segment.
func (sic SharedIterationsConfig) HasWork(et *lib.ExecutionTuple) bool ***REMOVED***
	return sic.GetVUs(et) > 0 && sic.GetIterations(et) > 0
***REMOVED***

// Init values needed for the execution
func (si *SharedIterations) Init(ctx context.Context) error ***REMOVED***
	// err should always be nil, because Init() won't be called for executors
	// with no work, as determined by their config's HasWork() method.
	et, err := si.BaseExecutor.executionState.ExecutionTuple.GetNewExecutionTupleFromValue(si.config.VUs.Int64)
	si.et = et
	si.iterSegIndex = lib.NewSegmentedIndex(et)

	return err
***REMOVED***

// Run executes a specific total number of iterations, which are all shared by
// the configured VUs.
//
//nolint:funlen
func (si SharedIterations) Run(parentCtx context.Context, out chan<- metrics.SampleContainer, workerInfo *lib.WorkerInfo) (err error) ***REMOVED***
	numVUs := si.config.GetVUs(si.executionState.ExecutionTuple)
	iterations := si.et.ScaleInt64(si.config.Iterations.Int64)
	duration := si.config.MaxDuration.TimeDuration()
	gracefulStop := si.config.GetGracefulStop()

	waitOnProgressChannel := make(chan struct***REMOVED******REMOVED***)
	startTime, maxDurationCtx, regDurationCtx, cancel := getDurationContexts(parentCtx, duration, gracefulStop)
	defer func() ***REMOVED***
		cancel()
		<-waitOnProgressChannel
	***REMOVED***()

	// Make sure the log and the progress bar have accurate information
	si.logger.WithFields(logrus.Fields***REMOVED***
		"vus": numVUs, "iterations": iterations, "maxDuration": duration, "type": si.config.GetType(),
	***REMOVED***).Debug("Starting executor run...")

	totalIters := uint64(iterations)
	doneIters := new(uint64)
	vusFmt := pb.GetFixedLengthIntFormat(numVUs)
	itersFmt := pb.GetFixedLengthIntFormat(int64(totalIters))
	progressFn := func() (float64, []string) ***REMOVED***
		spent := time.Since(startTime)
		progVUs := fmt.Sprintf(vusFmt+" VUs", numVUs)
		currentDoneIters := atomic.LoadUint64(doneIters)
		progIters := fmt.Sprintf(itersFmt+"/"+itersFmt+" shared iters",
			currentDoneIters, totalIters)
		spentDuration := pb.GetFixedLengthDuration(spent, duration)
		progDur := fmt.Sprintf("%s/%s", spentDuration, duration)
		right := []string***REMOVED***progVUs, progDur, progIters***REMOVED***

		return float64(currentDoneIters) / float64(totalIters), right
	***REMOVED***
	si.progress.Modify(pb.WithProgress(progressFn))
	maxDurationCtx = lib.WithScenarioState(maxDurationCtx, &lib.ScenarioState***REMOVED***
		Name:       si.config.Name,
		Executor:   si.config.Type,
		StartTime:  startTime,
		ProgressFn: progressFn,
	***REMOVED***)
	go func() ***REMOVED***
		trackProgress(parentCtx, maxDurationCtx, regDurationCtx, &si, progressFn)
		close(waitOnProgressChannel)
	***REMOVED***()

	var attemptedIters uint64

	// Actually schedule the VUs and iterations...
	activeVUs := &sync.WaitGroup***REMOVED******REMOVED***
	defer func() ***REMOVED***
		activeVUs.Wait()
		if attemptedIters < totalIters ***REMOVED***
			metrics.PushIfNotDone(parentCtx, out, metrics.Sample***REMOVED***
				Value:  float64(totalIters - attemptedIters),
				Metric: si.executionState.Test.BuiltinMetrics.DroppedIterations,
				Tags:   si.getMetricTags(nil), Time: time.Now(),
			***REMOVED***)
		***REMOVED***
	***REMOVED***()

	regDurationDone := regDurationCtx.Done()
	runIteration := getIterationRunner(si.executionState, si.logger)

	returnVU := func(u lib.InitializedVU) ***REMOVED***
		si.executionState.ReturnVU(u, true)
		activeVUs.Done()
	***REMOVED***

	handleVU := func(initVU lib.InitializedVU) ***REMOVED***
		ctx, cancel := context.WithCancel(maxDurationCtx)
		defer cancel()

		activeVU := initVU.Activate(getVUActivationParams(
			ctx, si.config.BaseConfig, returnVU, si.nextIterationCounters))

		for ***REMOVED***
			select ***REMOVED***
			case <-regDurationDone:
				return // don't make more iterations
			default:
				// continue looping
			***REMOVED***

			attemptedIterNumber := atomic.AddUint64(&attemptedIters, 1)
			if attemptedIterNumber > totalIters ***REMOVED***
				return
			***REMOVED***

			runIteration(maxDurationCtx, activeVU)
			atomic.AddUint64(doneIters, 1)
		***REMOVED***
	***REMOVED***

	for i := int64(0); i < numVUs; i++ ***REMOVED***
		initVU, err := si.executionState.GetPlannedVU(si.logger, true)
		if err != nil ***REMOVED***
			cancel()
			return err
		***REMOVED***
		activeVUs.Add(1)
		go handleVU(initVU)
	***REMOVED***

	return nil
***REMOVED***
