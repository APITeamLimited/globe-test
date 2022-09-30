package executor

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/APITeamLimited/globe-test/worker/workerMetrics"

	"github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
	"github.com/APITeamLimited/globe-test/worker/pb"
)

const constantArrivalRateType = "constant-arrival-rate"

func init() ***REMOVED***
	libWorker.RegisterExecutorConfigType(
		constantArrivalRateType,
		func(name string, rawJSON []byte) (libWorker.ExecutorConfig, error) ***REMOVED***
			config := NewConstantArrivalRateConfig(name)
			err := libWorker.StrictJSONUnmarshal(rawJSON, &config)
			return config, err
		***REMOVED***,
	)
***REMOVED***

// ConstantArrivalRateConfig stores config for the constant arrival-rate executor
type ConstantArrivalRateConfig struct ***REMOVED***
	BaseConfig
	Rate     null.Int           `json:"rate"`
	TimeUnit types.NullDuration `json:"timeUnit"`
	Duration types.NullDuration `json:"duration"`

	// Initialize `PreAllocatedVUs` number of VUs, and if more than that are needed,
	// they will be dynamically allocated, until `MaxVUs` is reached, which is an
	// absolutely hard limit on the number of VUs the executor will use
	PreAllocatedVUs null.Int `json:"preAllocatedVUs"`
	MaxVUs          null.Int `json:"maxVUs"`
***REMOVED***

// NewConstantArrivalRateConfig returns a ConstantArrivalRateConfig with default values
func NewConstantArrivalRateConfig(name string) *ConstantArrivalRateConfig ***REMOVED***
	return &ConstantArrivalRateConfig***REMOVED***
		BaseConfig: NewBaseConfig(name, constantArrivalRateType),
		TimeUnit:   types.NewNullDuration(1*time.Second, false),
	***REMOVED***
***REMOVED***

// Make sure we implement the libWorker.ExecutorConfig interface
var _ libWorker.ExecutorConfig = &ConstantArrivalRateConfig***REMOVED******REMOVED***

// GetPreAllocatedVUs is just a helper method that returns the scaled pre-allocated VUs.
func (carc ConstantArrivalRateConfig) GetPreAllocatedVUs(et *libWorker.ExecutionTuple) int64 ***REMOVED***
	return et.ScaleInt64(carc.PreAllocatedVUs.Int64)
***REMOVED***

// GetMaxVUs is just a helper method that returns the scaled max VUs.
func (carc ConstantArrivalRateConfig) GetMaxVUs(et *libWorker.ExecutionTuple) int64 ***REMOVED***
	return et.ScaleInt64(carc.MaxVUs.Int64)
***REMOVED***

// GetDescription returns a human-readable description of the executor options
func (carc ConstantArrivalRateConfig) GetDescription(et *libWorker.ExecutionTuple) string ***REMOVED***
	preAllocatedVUs, maxVUs := carc.GetPreAllocatedVUs(et), carc.GetMaxVUs(et)
	maxVUsRange := fmt.Sprintf("maxVUs: %d", preAllocatedVUs)
	if maxVUs > preAllocatedVUs ***REMOVED***
		maxVUsRange += fmt.Sprintf("-%d", maxVUs)
	***REMOVED***

	timeUnit := carc.TimeUnit.TimeDuration()
	var arrRatePerSec float64
	if maxVUs != 0 ***REMOVED*** // TODO: do something better?
		ratio := big.NewRat(maxVUs, carc.MaxVUs.Int64)
		arrRate := big.NewRat(carc.Rate.Int64, int64(timeUnit))
		arrRate.Mul(arrRate, ratio)
		arrRatePerSec, _ = getArrivalRatePerSec(arrRate).Float64()
	***REMOVED***

	return fmt.Sprintf("%.2f iterations/s for %s%s", arrRatePerSec, carc.Duration.Duration,
		carc.getBaseInfo(maxVUsRange))
***REMOVED***

// Validate makes sure all options are configured and valid
func (carc *ConstantArrivalRateConfig) Validate() []error ***REMOVED***
	errors := carc.BaseConfig.Validate()
	if !carc.Rate.Valid ***REMOVED***
		errors = append(errors, fmt.Errorf("the iteration rate isn't specified"))
	***REMOVED*** else if carc.Rate.Int64 <= 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("the iteration rate must be more than 0"))
	***REMOVED***

	if carc.TimeUnit.TimeDuration() <= 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("the timeUnit must be more than 0"))
	***REMOVED***

	if !carc.Duration.Valid ***REMOVED***
		errors = append(errors, fmt.Errorf("the duration is unspecified"))
	***REMOVED*** else if carc.Duration.TimeDuration() < minDuration ***REMOVED***
		errors = append(errors, fmt.Errorf(
			"the duration must be at least %s, but is %s", minDuration, carc.Duration,
		))
	***REMOVED***

	if !carc.PreAllocatedVUs.Valid ***REMOVED***
		errors = append(errors, fmt.Errorf("the number of preAllocatedVUs isn't specified"))
	***REMOVED*** else if carc.PreAllocatedVUs.Int64 < 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("the number of preAllocatedVUs can't be negative"))
	***REMOVED***

	if !carc.MaxVUs.Valid ***REMOVED***
		// TODO: don't change the config while validating
		carc.MaxVUs.Int64 = carc.PreAllocatedVUs.Int64
	***REMOVED*** else if carc.MaxVUs.Int64 < carc.PreAllocatedVUs.Int64 ***REMOVED***
		errors = append(errors, fmt.Errorf("maxVUs can't be less than preAllocatedVUs"))
	***REMOVED***

	return errors
***REMOVED***

// GetExecutionRequirements returns the number of required VUs to run the
// executor for its whole duration (disregarding any startTime), including the
// maximum waiting time for any iterations to gracefully stop. This is used by
// the execution scheduler in its VU reservation calculations, so it knows how
// many VUs to pre-initialize.
func (carc ConstantArrivalRateConfig) GetExecutionRequirements(et *libWorker.ExecutionTuple) []libWorker.ExecutionStep ***REMOVED***
	return []libWorker.ExecutionStep***REMOVED***
		***REMOVED***
			TimeOffset:      0,
			PlannedVUs:      uint64(et.ScaleInt64(carc.PreAllocatedVUs.Int64)),
			MaxUnplannedVUs: uint64(et.ScaleInt64(carc.MaxVUs.Int64) - et.ScaleInt64(carc.PreAllocatedVUs.Int64)),
		***REMOVED***, ***REMOVED***
			TimeOffset:      carc.Duration.TimeDuration() + carc.GracefulStop.TimeDuration(),
			PlannedVUs:      0,
			MaxUnplannedVUs: 0,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// NewExecutor creates a new ConstantArrivalRate executor
func (carc ConstantArrivalRateConfig) NewExecutor(
	es *libWorker.ExecutionState, logger *logrus.Entry,
) (libWorker.Executor, error) ***REMOVED***
	return &ConstantArrivalRate***REMOVED***
		BaseExecutor: NewBaseExecutor(&carc, es, logger),
		config:       carc,
	***REMOVED***, nil
***REMOVED***

// HasWork reports whether there is any work to be done for the given execution segment.
func (carc ConstantArrivalRateConfig) HasWork(et *libWorker.ExecutionTuple) bool ***REMOVED***
	return carc.GetMaxVUs(et) > 0
***REMOVED***

// ConstantArrivalRate tries to execute a specific number of iterations for a
// specific period.
type ConstantArrivalRate struct ***REMOVED***
	*BaseExecutor
	config ConstantArrivalRateConfig
	et     *libWorker.ExecutionTuple
***REMOVED***

// Make sure we implement the libWorker.Executor interface.
var _ libWorker.Executor = &ConstantArrivalRate***REMOVED******REMOVED***

// Init values needed for the execution
func (car *ConstantArrivalRate) Init(ctx context.Context) error ***REMOVED***
	// err should always be nil, because Init() won't be called for executors
	// with no work, as determined by their config's HasWork() method.
	et, err := car.BaseExecutor.executionState.ExecutionTuple.GetNewExecutionTupleFromValue(car.config.MaxVUs.Int64)
	car.et = et
	car.iterSegIndex = libWorker.NewSegmentedIndex(et)

	return err
***REMOVED***

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
func (car ConstantArrivalRate) Run(parentCtx context.Context, out chan<- workerMetrics.SampleContainer, workerInfo *libWorker.WorkerInfo) (err error) ***REMOVED***
	gracefulStop := car.config.GetGracefulStop()
	duration := car.config.Duration.TimeDuration()
	preAllocatedVUs := car.config.GetPreAllocatedVUs(car.executionState.ExecutionTuple)
	maxVUs := car.config.GetMaxVUs(car.executionState.ExecutionTuple)
	// TODO: refactor and simplify
	arrivalRate := getScaledArrivalRate(car.et.Segment, car.config.Rate.Int64, car.config.TimeUnit.TimeDuration())
	tickerPeriod := getTickerPeriod(arrivalRate).TimeDuration()
	arrivalRatePerSec, _ := getArrivalRatePerSec(arrivalRate).Float64()

	// Make sure the log and the progress bar have accurate information
	car.logger.WithFields(logrus.Fields***REMOVED***
		"maxVUs": maxVUs, "preAllocatedVUs": preAllocatedVUs, "duration": duration,
		"tickerPeriod": tickerPeriod, "type": car.config.GetType(),
	***REMOVED***).Debug("Starting executor run...")

	activeVUsWg := &sync.WaitGroup***REMOVED******REMOVED***

	returnedVUs := make(chan struct***REMOVED******REMOVED***)
	waitOnProgressChannel := make(chan struct***REMOVED******REMOVED***)
	startTime, maxDurationCtx, regDurationCtx, cancel := getDurationContexts(parentCtx, duration, gracefulStop)
	defer func() ***REMOVED***
		cancel()
		<-waitOnProgressChannel
	***REMOVED***()

	vusPool := newActiveVUPool()
	defer func() ***REMOVED***
		// Make sure all VUs aren't executing iterations anymore, for the cancel()
		// below to deactivate them.
		<-returnedVUs
		// first close the vusPool so we wait for the gracefulShutdown
		vusPool.Close()
		cancel()
		activeVUsWg.Wait()
	***REMOVED***()
	activeVUsCount := uint64(0)

	vusFmt := pb.GetFixedLengthIntFormat(maxVUs)
	progIters := fmt.Sprintf(
		pb.GetFixedLengthFloatFormat(arrivalRatePerSec, 2)+" iters/s", arrivalRatePerSec)
	progressFn := func() (float64, []string) ***REMOVED***
		spent := time.Since(startTime)
		currActiveVUs := atomic.LoadUint64(&activeVUsCount)
		progVUs := fmt.Sprintf(vusFmt+"/"+vusFmt+" VUs",
			vusPool.Running(), currActiveVUs)

		right := []string***REMOVED***progVUs, duration.String(), progIters***REMOVED***

		if spent > duration ***REMOVED***
			return 1, right
		***REMOVED***

		spentDuration := pb.GetFixedLengthDuration(spent, duration)
		progDur := fmt.Sprintf("%s/%s", spentDuration, duration)
		right[1] = progDur

		return math.Min(1, float64(spent)/float64(duration)), right
	***REMOVED***
	car.progress.Modify(pb.WithProgress(progressFn))
	maxDurationCtx = libWorker.WithScenarioState(maxDurationCtx, &libWorker.ScenarioState***REMOVED***
		Name:       car.config.Name,
		Executor:   car.config.Type,
		StartTime:  startTime,
		ProgressFn: progressFn,
	***REMOVED***)

	go func() ***REMOVED***
		trackProgress(parentCtx, maxDurationCtx, regDurationCtx, &car, progressFn)
		close(waitOnProgressChannel)
	***REMOVED***()

	returnVU := func(u libWorker.InitializedVU) ***REMOVED***
		car.executionState.ReturnVU(u, true)
		activeVUsWg.Done()
	***REMOVED***

	runIterationBasic := getIterationRunner(car.executionState, car.logger)
	activateVU := func(initVU libWorker.InitializedVU) libWorker.ActiveVU ***REMOVED***
		activeVUsWg.Add(1)
		activeVU := initVU.Activate(getVUActivationParams(
			maxDurationCtx, car.config.BaseConfig, returnVU,
			car.nextIterationCounters,
		))
		car.executionState.ModCurrentlyActiveVUsCount(+1)
		atomic.AddUint64(&activeVUsCount, 1)
		vusPool.AddVU(maxDurationCtx, activeVU, runIterationBasic)
		return activeVU
	***REMOVED***

	remainingUnplannedVUs := maxVUs - preAllocatedVUs
	makeUnplannedVUCh := make(chan struct***REMOVED******REMOVED***)
	defer close(makeUnplannedVUCh)
	go func() ***REMOVED***
		defer close(returnedVUs)
		for range makeUnplannedVUCh ***REMOVED***
			car.logger.Debug("Starting initialization of an unplanned VU...")
			initVU, err := car.executionState.GetUnplannedVU(maxDurationCtx, car.logger, workerInfo)
			if err != nil ***REMOVED***
				// TODO figure out how to return it to the Run goroutine
				car.logger.WithError(err).Error("Error while allocating unplanned VU")
			***REMOVED*** else ***REMOVED***
				car.logger.Debug("The unplanned VU finished initializing successfully!")
				activateVU(initVU)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	// Get the pre-allocated VUs in the local buffer
	for i := int64(0); i < preAllocatedVUs; i++ ***REMOVED***
		initVU, err := car.executionState.GetPlannedVU(car.logger, false)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		activateVU(initVU)
	***REMOVED***

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
	for li, gi := 0, start; ; li, gi = li+1, gi+offsets[li%len(offsets)] ***REMOVED***
		t := notScaledTickerPeriod*time.Duration(gi) - time.Since(startTime)
		timer.Reset(t)
		select ***REMOVED***
		case <-timer.C:
			if vusPool.TryRunIteration() ***REMOVED***
				continue
			***REMOVED***

			// Since there aren't any free VUs available, consider this iteration
			// dropped - we aren't going to try to recover it, but

			workerMetrics.PushIfNotDone(parentCtx, out, workerMetrics.Sample***REMOVED***
				Value: 1, Metric: droppedIterationMetric,
				Tags: metricTags, Time: time.Now(),
			***REMOVED***)

			// We'll try to start allocating another VU in the background,
			// non-blockingly, if we have remainingUnplannedVUs...
			if remainingUnplannedVUs == 0 ***REMOVED***
				if !shownWarning ***REMOVED***
					car.logger.Warningf("Insufficient VUs, reached %d active VUs and cannot initialize more", maxVUs)
					shownWarning = true
				***REMOVED***
				continue
			***REMOVED***

			select ***REMOVED***
			case makeUnplannedVUCh <- struct***REMOVED******REMOVED******REMOVED******REMOVED***: // great!
				remainingUnplannedVUs--
			default: // we're already allocating a new VU
			***REMOVED***

		case <-regDurationCtx.Done():
			return nil
		***REMOVED***
	***REMOVED***
***REMOVED***
