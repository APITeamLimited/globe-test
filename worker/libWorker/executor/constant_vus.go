package executor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/APITeamLimited/globe-test/worker/metrics"
	"github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
)

const constantVUsType = "constant-vus"

func init() {
	libWorker.RegisterExecutorConfigType(
		constantVUsType,
		func(name string, rawJSON []byte) (libWorker.ExecutorConfig, error) {
			config := NewConstantVUsConfig(name)
			err := libWorker.StrictJSONUnmarshal(rawJSON, &config)
			return config, err
		},
	)
}

// The minimum duration we'll allow users to schedule. This doesn't affect the stages
// configuration, where 0-duration virtual stages are allowed for instantaneous VU jumps
const minDuration = 1 * time.Second

// ConstantVUsConfig stores VUs and duration
type ConstantVUsConfig struct {
	BaseConfig
	VUs      null.Int           `json:"vus"`
	Duration types.NullDuration `json:"duration"`
}

// NewConstantVUsConfig returns a ConstantVUsConfig with default values
func NewConstantVUsConfig(name string) ConstantVUsConfig {
	return ConstantVUsConfig{
		BaseConfig: NewBaseConfig(name, constantVUsType),
		VUs:        null.NewInt(1, false),
	}
}

// Make sure we implement the libWorker.ExecutorConfig interface
var _ libWorker.ExecutorConfig = &ConstantVUsConfig{}

// GetVUs returns the scaled VUs for the executor.
func (clvc ConstantVUsConfig) GetVUs(et *libWorker.ExecutionTuple) int64 {
	return et.ScaleInt64(clvc.VUs.Int64)
}

// GetDescription returns a human-readable description of the executor options
func (clvc ConstantVUsConfig) GetDescription(et *libWorker.ExecutionTuple) string {
	return fmt.Sprintf("%d looping VUs for %s%s",
		clvc.GetVUs(et), clvc.Duration.Duration, clvc.getBaseInfo())
}

// Validate makes sure all options are configured and valid
func (clvc ConstantVUsConfig) Validate() []error {
	errors := clvc.BaseConfig.Validate()
	if clvc.VUs.Int64 <= 0 {
		errors = append(errors, fmt.Errorf("the number of VUs must be more than 0"))
	}

	if !clvc.Duration.Valid {
		errors = append(errors, fmt.Errorf("the duration is unspecified"))
	} else if clvc.Duration.TimeDuration() < minDuration {
		errors = append(errors, fmt.Errorf(
			"the duration must be at least %s, but is %s", minDuration, clvc.Duration,
		))
	}

	return errors
}

// GetExecutionRequirements returns the number of required VUs to run the
// executor for its whole duration (disregarding any startTime), including the
// maximum waiting time for any iterations to gracefully stop. This is used by
// the execution scheduler in its VU reservation calculations, so it knows how
// many VUs to pre-initialize.
func (clvc ConstantVUsConfig) GetExecutionRequirements(et *libWorker.ExecutionTuple) []libWorker.ExecutionStep {
	return []libWorker.ExecutionStep{
		{
			TimeOffset: 0,
			PlannedVUs: uint64(clvc.GetVUs(et)),
		},
		{
			TimeOffset: clvc.Duration.TimeDuration() + clvc.GracefulStop.TimeDuration(),
			PlannedVUs: 0,
		},
	}
}

// HasWork reports whether there is any work to be done for the given execution segment.
func (clvc ConstantVUsConfig) HasWork(et *libWorker.ExecutionTuple) bool {
	return clvc.GetVUs(et) > 0
}

// NewExecutor creates a new ConstantVUs executor
func (clvc ConstantVUsConfig) NewExecutor(es *libWorker.ExecutionState, logger *logrus.Entry) (libWorker.Executor, error) {
	return ConstantVUs{
		BaseExecutor: NewBaseExecutor(clvc, es, logger),
		config:       clvc,
	}, nil
}

// ConstantVUs maintains a constant number of VUs running for the
// specified duration.
type ConstantVUs struct {
	*BaseExecutor
	config ConstantVUsConfig
}

// Make sure we implement the libWorker.Executor interface.
var _ libWorker.Executor = &ConstantVUs{}

// Run constantly loops through as many iterations as possible on a fixed number
// of VUs for the specified duration.
func (clv ConstantVUs) Run(parentCtx context.Context, out chan<- metrics.SampleContainer, workerInfo *libWorker.WorkerInfo) (err error) {
	numVUs := clv.config.GetVUs(clv.executionState.ExecutionTuple)
	duration := clv.config.Duration.TimeDuration()
	gracefulStop := clv.config.GetGracefulStop()

	startTime, maxDurationCtx, regDurationCtx, cancel := getDurationContexts(parentCtx, duration, gracefulStop)
	defer cancel()

	// Make sure the log and the progress bar have accurate information
	clv.logger.WithFields(
		logrus.Fields{"vus": numVUs, "duration": duration, "type": clv.config.GetType()},
	).Debug("Starting executor run...")

	maxDurationCtx = libWorker.WithScenarioState(maxDurationCtx, &libWorker.ScenarioState{
		Name:      clv.config.Name,
		Executor:  clv.config.Type,
		StartTime: startTime,
	})

	// Actually schedule the VUs and iterations...
	activeVUs := &sync.WaitGroup{}
	defer activeVUs.Wait()

	regDurationDone := regDurationCtx.Done()
	runIteration := getIterationRunner(clv.executionState, clv.logger)

	returnVU := func(u libWorker.InitializedVU) {
		clv.executionState.ReturnVU(u, true)
		activeVUs.Done()
	}

	handleVU := func(initVU libWorker.InitializedVU) {
		ctx, cancel := context.WithCancel(maxDurationCtx)
		defer cancel()

		activeVU := initVU.Activate(
			getVUActivationParams(ctx, clv.config.BaseConfig, returnVU,
				clv.nextIterationCounters))

		for {
			select {
			case <-regDurationDone:
				return // don't make more iterations
			default:
				// continue looping
			}
			runIteration(maxDurationCtx, activeVU)
		}
	}

	for i := int64(0); i < numVUs; i++ {
		initVU, err := clv.executionState.GetPlannedVU(clv.logger, true)
		if err != nil {
			cancel()
			return err
		}
		activeVUs.Add(1)
		go handleVU(initVU)
	}

	return nil
}

func (clvc ConstantVUsConfig) GetMaxExecutorVUs() int64 {
	return clvc.VUs.Int64
}

func (clvc ConstantVUsConfig) ScaleOptions(subFraction float64) libWorker.ExecutorConfig {
	newConfig := clvc

	if newConfig.VUs.Valid {
		newConfig.VUs.Int64 = int64(float64(newConfig.VUs.Int64) * subFraction)

		if newConfig.VUs.Int64 < 1 {
			newConfig.VUs.Int64 = 1
		}
	}

	return newConfig
}
