package executor

import (
	"context"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
	"github.com/APITeamLimited/globe-test/worker/pb"
	"github.com/APITeamLimited/globe-test/worker/workerMetrics"
	"github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"
)

const sharedIterationsType = "shared-iterations"

func init() {
	libWorker.RegisterExecutorConfigType(
		sharedIterationsType,
		func(name string, rawJSON []byte) (libWorker.ExecutorConfig, error) {
			config := NewSharedIterationsConfig(name)
			err := libWorker.StrictJSONUnmarshal(rawJSON, &config)
			return config, err
		},
	)
}

// SharedIterationsConfig stores the number of VUs iterations, as well as maxDuration settings
type SharedIterationsConfig struct {
	BaseConfig
	VUs         null.Int           `json:"vus"`
	Iterations  null.Int           `json:"iterations"`
	MaxDuration types.NullDuration `json:"maxDuration"`
}

// NewSharedIterationsConfig returns a SharedIterationsConfig with default values
func NewSharedIterationsConfig(name string) SharedIterationsConfig {
	return SharedIterationsConfig{
		BaseConfig:  NewBaseConfig(name, sharedIterationsType),
		VUs:         null.NewInt(1, false),
		Iterations:  null.NewInt(1, false),
		MaxDuration: types.NewNullDuration(10*time.Minute, false), // TODO: shorten?
	}
}

// Make sure we implement the libWorker.ExecutorConfig interface
var _ libWorker.ExecutorConfig = &SharedIterationsConfig{}

// GetVUs returns the scaled VUs for the executor.
func (sic SharedIterationsConfig) GetVUs(et *libWorker.ExecutionTuple) int64 {
	return et.ScaleInt64(sic.VUs.Int64)
}

// GetIterations returns the scaled iteration count for the executor.
func (sic SharedIterationsConfig) GetIterations(et *libWorker.ExecutionTuple) int64 {
	// TODO: Optimize this by probably changing the whole Config API
	newTuple, err := et.GetNewExecutionTupleFromValue(sic.VUs.Int64)
	if err != nil {
		return 0
	}
	return newTuple.ScaleInt64(sic.Iterations.Int64)
}

// GetDescription returns a human-readable description of the executor options
func (sic SharedIterationsConfig) GetDescription(et *libWorker.ExecutionTuple) string {
	return fmt.Sprintf("%d iterations shared among %d VUs%s",
		sic.GetIterations(et), sic.GetVUs(et),
		sic.getBaseInfo(fmt.Sprintf("maxDuration: %s", sic.MaxDuration.Duration)))
}

// Validate makes sure all options are configured and valid
func (sic SharedIterationsConfig) Validate() []error {
	errors := sic.BaseConfig.Validate()
	if sic.VUs.Int64 <= 0 {
		errors = append(errors, fmt.Errorf("the number of VUs must be more than 0"))
	}

	if sic.Iterations.Int64 < sic.VUs.Int64 {
		errors = append(errors, fmt.Errorf(
			"the number of iterations (%d) can't be less than the number of VUs (%d)",
			sic.Iterations.Int64, sic.VUs.Int64,
		))
	}

	if sic.MaxDuration.TimeDuration() < minDuration {
		errors = append(errors, fmt.Errorf(
			"the maxDuration must be at least %s, but is %s", minDuration, sic.MaxDuration,
		))
	}

	return errors
}

// GetExecutionRequirements returns the number of required VUs to run the
// executor for its whole duration (disregarding any startTime), including the
// maximum waiting time for any iterations to gracefully stop. This is used by
// the execution scheduler in its VU reservation calculations, so it knows how
// many VUs to pre-initialize.
func (sic SharedIterationsConfig) GetExecutionRequirements(et *libWorker.ExecutionTuple) []libWorker.ExecutionStep {
	vus := sic.GetVUs(et)
	if vus == 0 {
		return []libWorker.ExecutionStep{
			{
				TimeOffset: 0,
				PlannedVUs: 0,
			},
		}
	}

	return []libWorker.ExecutionStep{
		{
			TimeOffset: 0,
			PlannedVUs: uint64(vus),
		},
		{
			TimeOffset: sic.MaxDuration.TimeDuration() + sic.GracefulStop.TimeDuration(),
			PlannedVUs: 0,
		},
	}
}

// NewExecutor creates a new SharedIterations executor
func (sic SharedIterationsConfig) NewExecutor(
	es *libWorker.ExecutionState, logger *logrus.Entry,
) (libWorker.Executor, error) {
	return &SharedIterations{
		BaseExecutor: NewBaseExecutor(sic, es, logger),
		config:       sic,
	}, nil
}

// SharedIterations executes a specific total number of iterations, which are
// all shared by the configured VUs.
type SharedIterations struct {
	*BaseExecutor
	config SharedIterationsConfig
	et     *libWorker.ExecutionTuple
}

// Make sure we implement the libWorker.Executor interface.
var _ libWorker.Executor = &SharedIterations{}

// HasWork reports whether there is any work to be done for the given execution segment.
func (sic SharedIterationsConfig) HasWork(et *libWorker.ExecutionTuple) bool {
	return sic.GetVUs(et) > 0 && sic.GetIterations(et) > 0
}

// Init values needed for the execution
func (si *SharedIterations) Init(ctx context.Context) error {
	// err should always be nil, because Init() won't be called for executors
	// with no work, as determined by their config's HasWork() method.
	et, err := si.BaseExecutor.executionState.ExecutionTuple.GetNewExecutionTupleFromValue(si.config.VUs.Int64)
	si.et = et
	si.iterSegIndex = libWorker.NewSegmentedIndex(et)

	return err
}

// Run executes a specific total number of iterations, which are all shared by
// the configured VUs.
//
//nolint:funlen
func (si SharedIterations) Run(parentCtx context.Context, out chan<- workerMetrics.SampleContainer, workerInfo *libWorker.WorkerInfo) (err error) {
	numVUs := si.config.GetVUs(si.executionState.ExecutionTuple)
	iterations := si.et.ScaleInt64(si.config.Iterations.Int64)
	duration := si.config.MaxDuration.TimeDuration()
	gracefulStop := si.config.GetGracefulStop()

	waitOnProgressChannel := make(chan struct{})
	startTime, maxDurationCtx, regDurationCtx, cancel := getDurationContexts(parentCtx, duration, gracefulStop)
	defer func() {
		cancel()
		<-waitOnProgressChannel
	}()

	// Make sure the log and the progress bar have accurate information
	si.logger.WithFields(logrus.Fields{
		"vus": numVUs, "iterations": iterations, "maxDuration": duration, "type": si.config.GetType(),
	}).Debug("Starting executor run...")

	totalIters := uint64(iterations)
	doneIters := new(uint64)
	vusFmt := pb.GetFixedLengthIntFormat(numVUs)
	itersFmt := pb.GetFixedLengthIntFormat(int64(totalIters))
	progressFn := func() (float64, []string) {
		spent := time.Since(startTime)
		progVUs := fmt.Sprintf(vusFmt+" VUs", numVUs)
		currentDoneIters := atomic.LoadUint64(doneIters)
		progIters := fmt.Sprintf(itersFmt+"/"+itersFmt+" shared iters",
			currentDoneIters, totalIters)
		spentDuration := pb.GetFixedLengthDuration(spent, duration)
		progDur := fmt.Sprintf("%s/%s", spentDuration, duration)
		right := []string{progVUs, progDur, progIters}

		return float64(currentDoneIters) / float64(totalIters), right
	}
	si.progress.Modify(pb.WithProgress(progressFn))
	maxDurationCtx = libWorker.WithScenarioState(maxDurationCtx, &libWorker.ScenarioState{
		Name:       si.config.Name,
		Executor:   si.config.Type,
		StartTime:  startTime,
		ProgressFn: progressFn,
	})
	go func() {
		trackProgress(parentCtx, maxDurationCtx, regDurationCtx, &si, progressFn)
		close(waitOnProgressChannel)
	}()

	var attemptedIters uint64

	// Actually schedule the VUs and iterations...
	activeVUs := &sync.WaitGroup{}
	defer func() {
		activeVUs.Wait()
		if attemptedIters < totalIters {
			workerMetrics.PushIfNotDone(parentCtx, out, workerMetrics.Sample{
				Value:  float64(totalIters - attemptedIters),
				Metric: si.executionState.Test.BuiltinMetrics.DroppedIterations,
				Tags:   si.getMetricTags(nil), Time: time.Now(),
			})
		}
	}()

	regDurationDone := regDurationCtx.Done()
	runIteration := getIterationRunner(si.executionState, si.logger)

	returnVU := func(u libWorker.InitializedVU) {
		si.executionState.ReturnVU(u, true)
		activeVUs.Done()
	}

	handleVU := func(initVU libWorker.InitializedVU) {
		ctx, cancel := context.WithCancel(maxDurationCtx)
		defer cancel()

		activeVU := initVU.Activate(getVUActivationParams(
			ctx, si.config.BaseConfig, returnVU, si.nextIterationCounters))

		for {
			select {
			case <-regDurationDone:
				return // don't make more iterations
			default:
				// continue looping
			}

			attemptedIterNumber := atomic.AddUint64(&attemptedIters, 1)
			if attemptedIterNumber > totalIters {
				return
			}

			runIteration(maxDurationCtx, activeVU)
			atomic.AddUint64(doneIters, 1)
		}
	}

	for i := int64(0); i < numVUs; i++ {
		initVU, err := si.executionState.GetPlannedVU(si.logger, true)
		if err != nil {
			cancel()
			return err
		}
		activeVUs.Add(1)
		go handleVU(initVU)
	}

	return nil
}

func (sic SharedIterationsConfig) GetMaxExecutorVUs() int64 {
	return sic.VUs.Int64
}

func (sic SharedIterationsConfig) ScaleOptions(subFraction float64) libWorker.ExecutorConfig {
	newConfig := sic

	if newConfig.VUs.Valid {
		newConfig.VUs.Int64 = int64(math.Round(float64(newConfig.VUs.Int64) * float64(subFraction)))

		if newConfig.VUs.Int64 < 1 {
			newConfig.VUs.Int64 = 1
		}
	}

	if newConfig.Iterations.Valid {
		newConfig.Iterations.Int64 = int64(math.Round(float64(newConfig.Iterations.Int64) * float64(subFraction)))

		if newConfig.Iterations.Int64 < 1 {
			newConfig.Iterations.Int64 = 1
		}
	}

	return newConfig
}
