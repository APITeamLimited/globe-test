package executor

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/APITeamLimited/globe-test/worker/workerMetrics"
	"github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"

	"github.com/APITeamLimited/globe-test/lib/types"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
)

const perVUIterationsType = "per-vu-iterations"

func init() {
	libWorker.RegisterExecutorConfigType(perVUIterationsType, func(name string, rawJSON []byte) (libWorker.ExecutorConfig, error) {
		config := NewPerVUIterationsConfig(name)
		err := libWorker.StrictJSONUnmarshal(rawJSON, &config)
		return config, err
	})
}

// PerVUIterationsConfig stores the number of VUs iterations, as well as maxDuration settings
type PerVUIterationsConfig struct {
	BaseConfig
	VUs         null.Int           `json:"vus"`
	Iterations  null.Int           `json:"iterations"`
	MaxDuration types.NullDuration `json:"maxDuration"`
}

// NewPerVUIterationsConfig returns a PerVUIterationsConfig with default values
func NewPerVUIterationsConfig(name string) PerVUIterationsConfig {
	return PerVUIterationsConfig{
		BaseConfig:  NewBaseConfig(name, perVUIterationsType),
		VUs:         null.NewInt(1, false),
		Iterations:  null.NewInt(1, false),
		MaxDuration: types.NewNullDuration(10*time.Minute, false), // TODO: shorten?
	}
}

// Make sure we implement the libWorker.ExecutorConfig interface
var _ libWorker.ExecutorConfig = &PerVUIterationsConfig{}

// GetVUs returns the scaled VUs for the executor.
func (pvic PerVUIterationsConfig) GetVUs(et *libWorker.ExecutionTuple) int64 {
	return et.ScaleInt64(pvic.VUs.Int64)
}

// GetIterations returns the UNSCALED iteration count for the executor. It's
// important to note that scaling per-VU iteration executor affects only the
// number of VUs. If we also scaled the iterations, scaling would have quadratic
// effects instead of just linear.
func (pvic PerVUIterationsConfig) GetIterations() int64 {
	return pvic.Iterations.Int64
}

// GetDescription returns a human-readable description of the executor options
func (pvic PerVUIterationsConfig) GetDescription(et *libWorker.ExecutionTuple) string {
	return fmt.Sprintf("%d iterations for each of %d VUs%s",
		pvic.GetIterations(), pvic.GetVUs(et),
		pvic.getBaseInfo(fmt.Sprintf("maxDuration: %s", pvic.MaxDuration.Duration)))
}

// Validate makes sure all options are configured and valid
func (pvic PerVUIterationsConfig) Validate() []error {
	errors := pvic.BaseConfig.Validate()
	if pvic.VUs.Int64 <= 0 {
		errors = append(errors, fmt.Errorf("the number of VUs must be more than 0"))
	}

	if pvic.Iterations.Int64 <= 0 {
		errors = append(errors, fmt.Errorf("the number of iterations must be more than 0"))
	}

	if pvic.MaxDuration.TimeDuration() < minDuration {
		errors = append(errors, fmt.Errorf(
			"the maxDuration must be at least %s, but is %s", minDuration, pvic.MaxDuration,
		))
	}

	return errors
}

// GetExecutionRequirements returns the number of required VUs to run the
// executor for its whole duration (disregarding any startTime), including the
// maximum waiting time for any iterations to gracefully stop. This is used by
// the execution scheduler in its VU reservation calculations, so it knows how
// many VUs to pre-initialize.
func (pvic PerVUIterationsConfig) GetExecutionRequirements(et *libWorker.ExecutionTuple) []libWorker.ExecutionStep {
	return []libWorker.ExecutionStep{
		{
			TimeOffset: 0,
			PlannedVUs: uint64(pvic.GetVUs(et)),
		},
		{
			TimeOffset: pvic.MaxDuration.TimeDuration() + pvic.GracefulStop.TimeDuration(),
			PlannedVUs: 0,
		},
	}
}

// NewExecutor creates a new PerVUIterations executor
func (pvic PerVUIterationsConfig) NewExecutor(
	es *libWorker.ExecutionState, logger *logrus.Entry,
) (libWorker.Executor, error) {
	return PerVUIterations{
		BaseExecutor: NewBaseExecutor(pvic, es, logger),
		config:       pvic,
	}, nil
}

// HasWork reports whether there is any work to be done for the given execution segment.
func (pvic PerVUIterationsConfig) HasWork(et *libWorker.ExecutionTuple) bool {
	return pvic.GetVUs(et) > 0 && pvic.GetIterations() > 0
}

// PerVUIterations executes a specific number of iterations with each VU.
type PerVUIterations struct {
	*BaseExecutor
	config PerVUIterationsConfig
}

// Make sure we implement the libWorker.Executor interface.
var _ libWorker.Executor = &PerVUIterations{}

// Run executes a specific number of iterations with each configured VU.
//
//nolint:funlen
func (pvi PerVUIterations) Run(parentCtx context.Context, out chan<- workerMetrics.SampleContainer, workerInfo *libWorker.WorkerInfo) (err error) {
	numVUs := pvi.config.GetVUs(pvi.executionState.ExecutionTuple)
	iterations := pvi.config.GetIterations()
	duration := pvi.config.MaxDuration.TimeDuration()
	gracefulStop := pvi.config.GetGracefulStop()

	waitOnProgressChannel := make(chan struct{})
	startTime, maxDurationCtx, regDurationCtx, cancel := getDurationContexts(parentCtx, duration, gracefulStop)
	defer func() {
		cancel()
		<-waitOnProgressChannel
	}()

	// Make sure the log and the progress bar have accurate information
	pvi.logger.WithFields(logrus.Fields{
		"vus": numVUs, "iterations": iterations, "maxDuration": duration, "type": pvi.config.GetType(),
	}).Debug("Starting executor run...")

	doneIters := new(uint64)

	maxDurationCtx = libWorker.WithScenarioState(maxDurationCtx, &libWorker.ScenarioState{
		Name:      pvi.config.Name,
		Executor:  pvi.config.Type,
		StartTime: startTime,
	})

	handleVUsWG := &sync.WaitGroup{}
	defer handleVUsWG.Wait()
	// Actually schedule the VUs and iterations...
	activeVUs := &sync.WaitGroup{}
	defer activeVUs.Wait()

	regDurationDone := regDurationCtx.Done()
	runIteration := getIterationRunner(pvi.executionState, pvi.logger)

	returnVU := func(u libWorker.InitializedVU) {
		pvi.executionState.ReturnVU(u, true)
		activeVUs.Done()
	}

	droppedIterationMetric := pvi.executionState.Test.BuiltinMetrics.DroppedIterations
	handleVU := func(initVU libWorker.InitializedVU) {
		defer handleVUsWG.Done()
		ctx, cancel := context.WithCancel(maxDurationCtx)
		defer cancel()

		vuID := initVU.GetID()
		activeVU := initVU.Activate(
			getVUActivationParams(ctx, pvi.config.BaseConfig, returnVU,
				pvi.nextIterationCounters))

		for i := int64(0); i < iterations; i++ {
			select {
			case <-regDurationDone:
				workerMetrics.PushIfNotDone(parentCtx, out, workerMetrics.Sample{
					Value: float64(iterations - i), Metric: droppedIterationMetric,
					Tags: pvi.getMetricTags(&vuID), Time: time.Now(),
				})
				return // don't make more iterations
			default:
				// continue looping
			}
			runIteration(maxDurationCtx, activeVU)
			atomic.AddUint64(doneIters, 1)
		}
	}

	for i := int64(0); i < numVUs; i++ {
		initializedVU, err := pvi.executionState.GetPlannedVU(pvi.logger, true)
		if err != nil {
			cancel()
			return err
		}
		activeVUs.Add(1)
		handleVUsWG.Add(1)
		go handleVU(initializedVU)
	}

	return nil
}

func (pvic PerVUIterationsConfig) GetMaxExecutorVUs() int64 {
	return pvic.VUs.Int64
}

func (pvic PerVUIterationsConfig) ScaleOptions(subFraction float64) libWorker.ExecutorConfig {
	newConfig := pvic

	if newConfig.VUs.Valid {
		newConfig.VUs.Int64 = int64(float64(newConfig.VUs.Int64) * subFraction)

		if newConfig.VUs.Int64 < 1 {
			newConfig.VUs.Int64 = 1
		}
	}

	if newConfig.Iterations.Valid {
		newConfig.Iterations.Int64 = int64(float64(newConfig.Iterations.Int64) * subFraction)

		if newConfig.Iterations.Int64 < 1 {
			newConfig.Iterations.Int64 = 1
		}
	}

	return newConfig
}
