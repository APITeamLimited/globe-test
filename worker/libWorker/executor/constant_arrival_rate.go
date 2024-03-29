package executor

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/APITeamLimited/globe-test/worker/metrics"

	"github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
)

const constantArrivalRateType = "constant-arrival-rate"

func init() {
	libWorker.RegisterExecutorConfigType(
		constantArrivalRateType,
		func(name string, rawJSON []byte) (libWorker.ExecutorConfig, error) {
			config := NewConstantArrivalRateConfig(name)
			err := libWorker.StrictJSONUnmarshal(rawJSON, &config)
			return config, err
		},
	)
}

// ConstantArrivalRateConfig stores config for the constant arrival-rate executor
type ConstantArrivalRateConfig struct {
	BaseConfig
	Rate     null.Int           `json:"rate"`
	TimeUnit types.NullDuration `json:"timeUnit"`
	Duration types.NullDuration `json:"duration"`

	// Initialize `PreAllocatedVUs` number of VUs, and if more than that are needed,
	// they will be dynamically allocated, until `MaxVUs` is reached, which is an
	// absolutely hard limit on the number of VUs the executor will use
	PreAllocatedVUs null.Int `json:"preAllocatedVUs"`
	MaxVUs          null.Int `json:"maxVUs"`
}

// NewConstantArrivalRateConfig returns a ConstantArrivalRateConfig with default values
func NewConstantArrivalRateConfig(name string) *ConstantArrivalRateConfig {
	return &ConstantArrivalRateConfig{
		BaseConfig: NewBaseConfig(name, constantArrivalRateType),
		TimeUnit:   types.NewNullDuration(1*time.Second, false),
	}
}

// Make sure we implement the libWorker.ExecutorConfig interface
var _ libWorker.ExecutorConfig = &ConstantArrivalRateConfig{}

// GetPreAllocatedVUs is just a helper method that returns the scaled pre-allocated VUs.
func (carc ConstantArrivalRateConfig) GetPreAllocatedVUs(et *libWorker.ExecutionTuple) int64 {
	return et.ScaleInt64(carc.PreAllocatedVUs.Int64)
}

// GetMaxVUs is just a helper method that returns the scaled max VUs.
func (carc ConstantArrivalRateConfig) GetMaxVUs(et *libWorker.ExecutionTuple) int64 {
	return et.ScaleInt64(carc.MaxVUs.Int64)
}

// GetDescription returns a human-readable description of the executor options
func (carc ConstantArrivalRateConfig) GetDescription(et *libWorker.ExecutionTuple) string {
	preAllocatedVUs, maxVUs := carc.GetPreAllocatedVUs(et), carc.GetMaxVUs(et)
	maxVUsRange := fmt.Sprintf("maxVUs: %d", preAllocatedVUs)
	if maxVUs > preAllocatedVUs {
		maxVUsRange += fmt.Sprintf("-%d", maxVUs)
	}

	timeUnit := carc.TimeUnit.TimeDuration()
	var arrRatePerSec float64
	if maxVUs != 0 { // TODO: do something better?
		ratio := big.NewRat(maxVUs, carc.MaxVUs.Int64)
		arrRate := big.NewRat(carc.Rate.Int64, int64(timeUnit))
		arrRate.Mul(arrRate, ratio)
		arrRatePerSec, _ = getArrivalRatePerSec(arrRate).Float64()
	}

	return fmt.Sprintf("%.2f iterations/s for %s%s", arrRatePerSec, carc.Duration.Duration,
		carc.getBaseInfo(maxVUsRange))
}

// Validate makes sure all options are configured and valid
func (carc ConstantArrivalRateConfig) Validate() []error {
	errors := carc.BaseConfig.Validate()
	if !carc.Rate.Valid {
		errors = append(errors, fmt.Errorf("the iteration rate isn't specified"))
	} else if carc.Rate.Int64 <= 0 {
		errors = append(errors, fmt.Errorf("the iteration rate must be more than 0"))
	}

	if carc.TimeUnit.TimeDuration() <= 0 {
		errors = append(errors, fmt.Errorf("the timeUnit must be more than 0"))
	}

	if !carc.Duration.Valid {
		errors = append(errors, fmt.Errorf("the duration is unspecified"))
	} else if carc.Duration.TimeDuration() < minDuration {
		errors = append(errors, fmt.Errorf(
			"the duration must be at least %s, but is %s", minDuration, carc.Duration,
		))
	}

	if !carc.PreAllocatedVUs.Valid {
		errors = append(errors, fmt.Errorf("the number of preAllocatedVUs isn't specified"))
	} else if carc.PreAllocatedVUs.Int64 < 0 {
		errors = append(errors, fmt.Errorf("the number of preAllocatedVUs can't be negative"))
	}

	if !carc.MaxVUs.Valid {
		// TODO: don't change the config while validating
		carc.MaxVUs.Int64 = carc.PreAllocatedVUs.Int64
	} else if carc.MaxVUs.Int64 < carc.PreAllocatedVUs.Int64 {
		errors = append(errors, fmt.Errorf("maxVUs can't be less than preAllocatedVUs"))
	}

	return errors
}

// GetExecutionRequirements returns the number of required VUs to run the
// executor for its whole duration (disregarding any startTime), including the
// maximum waiting time for any iterations to gracefully stop. This is used by
// the execution scheduler in its VU reservation calculations, so it knows how
// many VUs to pre-initialize.
func (carc ConstantArrivalRateConfig) GetExecutionRequirements(et *libWorker.ExecutionTuple) []libWorker.ExecutionStep {
	return []libWorker.ExecutionStep{
		{
			TimeOffset:      0,
			PlannedVUs:      uint64(et.ScaleInt64(carc.PreAllocatedVUs.Int64)),
			MaxUnplannedVUs: uint64(et.ScaleInt64(carc.MaxVUs.Int64) - et.ScaleInt64(carc.PreAllocatedVUs.Int64)),
		}, {
			TimeOffset:      carc.Duration.TimeDuration() + carc.GracefulStop.TimeDuration(),
			PlannedVUs:      0,
			MaxUnplannedVUs: 0,
		},
	}
}

// NewExecutor creates a new ConstantArrivalRate executor
func (carc ConstantArrivalRateConfig) NewExecutor(
	es *libWorker.ExecutionState, logger *logrus.Entry,
) (libWorker.Executor, error) {
	return &ConstantArrivalRate{
		BaseExecutor: NewBaseExecutor(&carc, es, logger),
		config:       carc,
	}, nil
}

// HasWork reports whether there is any work to be done for the given execution segment.
func (carc ConstantArrivalRateConfig) HasWork(et *libWorker.ExecutionTuple) bool {
	return carc.GetMaxVUs(et) > 0
}

// ConstantArrivalRate tries to execute a specific number of iterations for a
// specific period.
type ConstantArrivalRate struct {
	*BaseExecutor
	config ConstantArrivalRateConfig
	et     *libWorker.ExecutionTuple
}

// Make sure we implement the libWorker.Executor interface.
var _ libWorker.Executor = &ConstantArrivalRate{}

// Init values needed for the execution
func (car *ConstantArrivalRate) Init(ctx context.Context) error {
	// err should always be nil, because Init() won't be called for executors
	// with no work, as determined by their config's HasWork() method.
	et, err := car.BaseExecutor.executionState.ExecutionTuple.GetNewExecutionTupleFromValue(car.config.MaxVUs.Int64)
	car.et = et
	car.iterSegIndex = libWorker.NewSegmentedIndex(et)

	return err
}

// Run executes a constant number of iterations per second.
//
// TODO: Split this up and make an independent component that can be reused
// between the constant and ramping arrival rate executors - that way we can
// keep the complexity in one well-architected part (with short methods and few
// lambdas :D), while having both config frontends still be present for maximum
// UX benefits. Basically, keep the progress bars and scheduling (i.e. at what
// time should iteration X begin) different, but keep everything else the same.
// This will allow us to implement https://github.com/k6io/k6/issues/1386
// and things like all of the TODOs below in one place only.
//
//nolint:funlen,cyclop
func (car ConstantArrivalRate) Run(parentCtx context.Context, out chan<- metrics.SampleContainer, workerInfo *libWorker.WorkerInfo) (err error) {
	gracefulStop := car.config.GetGracefulStop()
	duration := car.config.Duration.TimeDuration()
	preAllocatedVUs := car.config.GetPreAllocatedVUs(car.executionState.ExecutionTuple)
	maxVUs := car.config.GetMaxVUs(car.executionState.ExecutionTuple)
	// TODO: refactor and simplify
	arrivalRate := getScaledArrivalRate(car.et.Segment, car.config.Rate.Int64, car.config.TimeUnit.TimeDuration())
	tickerPeriod := getTickerPeriod(arrivalRate).TimeDuration()

	// Make sure the log and the progress bar have accurate information
	car.logger.WithFields(logrus.Fields{
		"maxVUs": maxVUs, "preAllocatedVUs": preAllocatedVUs, "duration": duration,
		"tickerPeriod": tickerPeriod, "type": car.config.GetType(),
	}).Debug("Starting executor run...")

	activeVUsWg := &sync.WaitGroup{}

	returnedVUs := make(chan struct{})
	startTime, maxDurationCtx, regDurationCtx, cancel := getDurationContexts(parentCtx, duration, gracefulStop)
	defer cancel()

	vusPool := newActiveVUPool()
	defer func() {
		// Make sure all VUs aren't executing iterations anymore, for the cancel()
		// below to deactivate them.
		<-returnedVUs
		// first close the vusPool so we wait for the gracefulShutdown
		vusPool.Close()
		cancel()
		activeVUsWg.Wait()
	}()
	activeVUsCount := uint64(0)

	maxDurationCtx = libWorker.WithScenarioState(maxDurationCtx, &libWorker.ScenarioState{
		Name:      car.config.Name,
		Executor:  car.config.Type,
		StartTime: startTime,
	})

	returnVU := func(u libWorker.InitializedVU) {
		car.executionState.ReturnVU(u, true)
		activeVUsWg.Done()
	}

	runIterationBasic := getIterationRunner(car.executionState, car.logger)
	activateVU := func(initVU libWorker.InitializedVU) libWorker.ActiveVU {
		activeVUsWg.Add(1)
		activeVU := initVU.Activate(getVUActivationParams(
			maxDurationCtx, car.config.BaseConfig, returnVU,
			car.nextIterationCounters,
		))
		car.executionState.ModCurrentlyActiveVUsCount(+1)
		atomic.AddUint64(&activeVUsCount, 1)
		vusPool.AddVU(maxDurationCtx, activeVU, runIterationBasic)
		return activeVU
	}

	remainingUnplannedVUs := maxVUs - preAllocatedVUs
	makeUnplannedVUCh := make(chan struct{})
	defer close(makeUnplannedVUCh)
	go func() {
		defer close(returnedVUs)
		for range makeUnplannedVUCh {
			car.logger.Debug("Starting initialization of an unplanned VU...")
			initVU, err := car.executionState.GetUnplannedVU(maxDurationCtx, car.logger, workerInfo)
			if err != nil {
				// TODO figure out how to return it to the Run goroutine
				car.logger.WithError(err).Error("Error while allocating unplanned VU")
			} else {
				car.logger.Debug("The unplanned VU finished initializing successfully!")
				activateVU(initVU)
			}
		}
	}()

	// Get the pre-allocated VUs in the local buffer
	for i := int64(0); i < preAllocatedVUs; i++ {
		initVU, err := car.executionState.GetPlannedVU(car.logger, false)
		if err != nil {
			return err
		}
		activateVU(initVU)
	}

	start, offsets, _ := car.et.GetStripedOffsets()
	timer := time.NewTimer(time.Hour * 24)
	// here the we need the not scaled one
	notScaledTickerPeriod := getTickerPeriod(
		big.NewRat(
			car.config.Rate.Int64,
			int64(car.config.TimeUnit.TimeDuration()),
		)).TimeDuration()

	droppedIterationMetric := car.executionState.Test.BuiltinMetrics.DroppedIterations
	shownWarning := false
	metricTags := car.getMetricTags(nil)
	for li, gi := 0, start; ; li, gi = li+1, gi+offsets[li%len(offsets)] {
		t := notScaledTickerPeriod*time.Duration(gi) - time.Since(startTime)
		timer.Reset(t)
		select {
		case <-timer.C:
			if vusPool.TryRunIteration() {
				continue
			}

			// Since there aren't any free VUs available, consider this iteration
			// dropped - we aren't going to try to recover it, but

			metrics.PushIfNotDone(parentCtx, out, metrics.Sample{
				Value: 1, Metric: droppedIterationMetric,
				Tags: metricTags, Time: time.Now(),
			})

			// We'll try to start allocating another VU in the background,
			// non-blockingly, if we have remainingUnplannedVUs...
			if remainingUnplannedVUs == 0 {
				if !shownWarning {
					car.logger.Warningf("Insufficient VUs, reached %d active VUs and cannot initialize more", maxVUs)
					shownWarning = true
				}
				continue
			}

			select {
			case makeUnplannedVUCh <- struct{}{}: // great!
				remainingUnplannedVUs--
			default: // we're already allocating a new VU
			}

		case <-regDurationCtx.Done():
			return nil
		}
	}
}

func (config ConstantArrivalRateConfig) GetMaxExecutorVUs() int64 {
	return config.MaxVUs.Int64
}

func (config ConstantArrivalRateConfig) ScaleOptions(subFraction float64) libWorker.ExecutorConfig {
	newConfig := config

	if newConfig.Rate.Valid {
		newConfig.Rate.Int64 = int64(float64(newConfig.Rate.Int64) * float64(subFraction))

		if newConfig.Rate.Int64 < 1 {
			newConfig.Rate.Int64 = 1
		}
	}

	if newConfig.PreAllocatedVUs.Valid {
		newConfig.PreAllocatedVUs.Int64 = int64(float64(newConfig.PreAllocatedVUs.Int64) * float64(subFraction))

		if newConfig.PreAllocatedVUs.Int64 < 1 {
			newConfig.PreAllocatedVUs.Int64 = 1
		}
	}

	if newConfig.MaxVUs.Valid {
		newConfig.MaxVUs.Int64 = int64(float64(newConfig.MaxVUs.Int64) * float64(subFraction))

		if newConfig.MaxVUs.Int64 < 1 {
			newConfig.MaxVUs.Int64 = 1
		}
	}

	return newConfig
}
