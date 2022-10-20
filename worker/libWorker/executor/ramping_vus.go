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

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
	"github.com/APITeamLimited/globe-test/worker/pb"
)

const rampingVUsType = "ramping-vus"

func init() ***REMOVED***
	libWorker.RegisterExecutorConfigType(
		rampingVUsType,
		func(name string, rawJSON []byte) (libWorker.ExecutorConfig, error) ***REMOVED***
			config := NewRampingVUsConfig(name)
			err := libWorker.StrictJSONUnmarshal(rawJSON, &config)
			return config, err
		***REMOVED***,
	)
***REMOVED***

// Stage contains
type Stage struct ***REMOVED***
	Duration types.NullDuration `json:"duration"`
	Target   null.Int           `json:"target"` // TODO: maybe rename this to endVUs? something else?
	// TODO: add a progression function?
***REMOVED***

// RampingVUsConfig stores the configuration for the stages executor
type RampingVUsConfig struct ***REMOVED***
	BaseConfig
	StartVUs         null.Int           `json:"startVUs"`
	Stages           []Stage            `json:"stages"`
	GracefulRampDown types.NullDuration `json:"gracefulRampDown"`
***REMOVED***

// NewRampingVUsConfig returns a RampingVUsConfig with its default values
func NewRampingVUsConfig(name string) RampingVUsConfig ***REMOVED***
	return RampingVUsConfig***REMOVED***
		BaseConfig:       NewBaseConfig(name, rampingVUsType),
		StartVUs:         null.NewInt(1, false),
		GracefulRampDown: types.NewNullDuration(30*time.Second, false),
	***REMOVED***
***REMOVED***

// Make sure we implement the libWorker.ExecutorConfig interface
var _ libWorker.ExecutorConfig = &RampingVUsConfig***REMOVED******REMOVED***

// GetStartVUs is just a helper method that returns the scaled starting VUs.
func (vlvc RampingVUsConfig) GetStartVUs(et *libWorker.ExecutionTuple) int64 ***REMOVED***
	return et.ScaleInt64(vlvc.StartVUs.Int64)
***REMOVED***

// GetGracefulRampDown is just a helper method that returns the graceful
// ramp-down period as a standard Go time.Duration value...
func (vlvc RampingVUsConfig) GetGracefulRampDown() time.Duration ***REMOVED***
	return vlvc.GracefulRampDown.TimeDuration()
***REMOVED***

// GetDescription returns a human-readable description of the executor options
func (vlvc RampingVUsConfig) GetDescription(et *libWorker.ExecutionTuple) string ***REMOVED***
	maxVUs := et.ScaleInt64(getStagesUnscaledMaxTarget(vlvc.StartVUs.Int64, vlvc.Stages))
	return fmt.Sprintf("Up to %d looping VUs for %s over %d stages%s",
		maxVUs, sumStagesDuration(vlvc.Stages), len(vlvc.Stages),
		vlvc.getBaseInfo(fmt.Sprintf("gracefulRampDown: %s", vlvc.GetGracefulRampDown())))
***REMOVED***

// Validate makes sure all options are configured and valid
func (vlvc RampingVUsConfig) Validate() []error ***REMOVED***
	errors := vlvc.BaseConfig.Validate()
	if vlvc.StartVUs.Int64 < 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("the number of start VUs can't be negative"))
	***REMOVED***

	if getStagesUnscaledMaxTarget(vlvc.StartVUs.Int64, vlvc.Stages) <= 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("either startVUs or one of the stages' target values must be greater than 0"))
	***REMOVED***

	return append(errors, validateStages(vlvc.Stages)...)
***REMOVED***

// getRawExecutionSteps calculates and returns as execution steps the number of
// actively running VUs the executor should have at every moment.
//
// It doesn't take into account graceful ramp-downs. It also doesn't deal with
// the end-of-executor drop to 0 VUs, whether graceful or not. These are
// handled by GetExecutionRequirements(), which internally uses this method and
// reserveVUsForGracefulRampDowns().
//
// The zeroEnd argument tells the method if we should artificially add a step
// with 0 VUs at offset sum(stages.duration), i.e. when the executor is
// supposed to end.
//
// It's also important to note how scaling works. Say, we ramp up from 0 to 10
// VUs over 10 seconds and then back to 0, and we want to split the execution in
// 2 equal segments (i.e. execution segments "0:0.5" and "0.5:1"). The original
// execution steps would look something like this:
//
// VUs  ^
//
//	10|          *
//	 9|         ***
//	 8|        *****
//	 7|       *******
//	 6|      *********
//	 5|     ***********
//	 4|    *************
//	 3|   ***************
//	 2|  *****************
//	 1| *******************
//	 0------------------------> time(s)
//	   01234567890123456789012   (t%10)
//	   00000000001111111111222   (t/10)
//
// The chart for one of the execution segments would look like this:
//
// VUs  ^
//
//	5|         XXX
//	4|       XXXXXXX
//	3|     XXXXXXXXXXX
//	2|   XXXXXXXXXXXXXXX
//	1| XXXXXXXXXXXXXXXXXXX
//	0------------------------> time(s)
//	  01234567890123456789012   (t%10)
//	  00000000001111111111222   (t/10)
//
// And the chart for the other execution segment would look like this:
//
// VUs  ^
//
//	5|          Y
//	4|        YYYYY
//	3|      YYYYYYYYY
//	2|    YYYYYYYYYYYYY
//	1|  YYYYYYYYYYYYYYYYY
//	0------------------------> time(s)
//	  01234567890123456789012   (t%10)
//	  00000000001111111111222   (t/10)
//
// Notice the time offsets and the slower ramping up and down. All of that is
// because the sum of the two execution segments has to produce exactly the
// original shape, as if the test ran on a single machine:
//
// VUs  ^
//
//	10|          Y
//	 9|         XXX
//	 8|        YYYYY
//	 7|       XXXXXXX
//	 6|      YYYYYYYYY
//	 5|     XXXXXXXXXXX
//	 4|    YYYYYYYYYYYYY
//	 3|   XXXXXXXXXXXXXXX
//	 2|  YYYYYYYYYYYYYYYYY
//	 1| XXXXXXXXXXXXXXXXXXX
//	 0------------------------> time(s)
//	   01234567890123456789012   (t%10)
//	   00000000001111111111222   (t/10)
//
// More information: https://github.com/k6io/k6/issues/997#issuecomment-484416866
func (vlvc RampingVUsConfig) getRawExecutionSteps(et *libWorker.ExecutionTuple, zeroEnd bool) []libWorker.ExecutionStep ***REMOVED***
	var (
		timeTillEnd time.Duration
		fromVUs     = vlvc.StartVUs.Int64
		steps       = make([]libWorker.ExecutionStep, 0, vlvc.precalculateTheRequiredSteps(et, zeroEnd))
		index       = libWorker.NewSegmentedIndex(et)
	)

	// Reserve the scaled StartVUs at the beginning
	scaled, unscaled := index.GoTo(fromVUs)
	steps = append(steps, libWorker.ExecutionStep***REMOVED***TimeOffset: 0, PlannedVUs: uint64(scaled)***REMOVED***)
	addStep := func(timeOffset time.Duration, plannedVUs uint64) ***REMOVED***
		if steps[len(steps)-1].PlannedVUs != plannedVUs ***REMOVED***
			steps = append(steps, libWorker.ExecutionStep***REMOVED***TimeOffset: timeOffset, PlannedVUs: plannedVUs***REMOVED***)
		***REMOVED***
	***REMOVED***

	for _, stage := range vlvc.Stages ***REMOVED***
		stageEndVUs := stage.Target.Int64
		stageDuration := stage.Duration.TimeDuration()
		timeTillEnd += stageDuration

		stageVUDiff := stageEndVUs - fromVUs
		if stageVUDiff == 0 ***REMOVED***
			continue
		***REMOVED***
		if stageDuration == 0 ***REMOVED***
			scaled, unscaled = index.GoTo(stageEndVUs)
			addStep(timeTillEnd, uint64(scaled))
			fromVUs = stageEndVUs
			continue
		***REMOVED***

		// VU reservation for gracefully ramping down is handled as a
		// separate method: reserveVUsForGracefulRampDowns()
		if unscaled > stageEndVUs ***REMOVED*** // ramp down
			// here we don't want to emit for the equal to stageEndVUs as it doesn't go below it
			// it will just go to it
			for ; unscaled > stageEndVUs; scaled, unscaled = index.Prev() ***REMOVED***
				addStep(
					// this is the time that we should go up 1 if we are ramping up
					// but we are ramping down so we should go 1 down, but because we want to not
					// stop VUs immediately we stop it on the next unscaled VU's time
					timeTillEnd-time.Duration(int64(stageDuration)*(stageEndVUs-unscaled+1)/stageVUDiff),
					uint64(scaled-1),
				)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			for ; unscaled <= stageEndVUs; scaled, unscaled = index.Next() ***REMOVED***
				addStep(
					timeTillEnd-time.Duration(int64(stageDuration)*(stageEndVUs-unscaled)/stageVUDiff),
					uint64(scaled),
				)
			***REMOVED***
		***REMOVED***
		fromVUs = stageEndVUs
	***REMOVED***

	if zeroEnd && steps[len(steps)-1].PlannedVUs != 0 ***REMOVED***
		// If the last PlannedVUs value wasn't 0, add a last step with 0
		steps = append(steps, libWorker.ExecutionStep***REMOVED***TimeOffset: timeTillEnd, PlannedVUs: 0***REMOVED***)
	***REMOVED***
	return steps
***REMOVED***

func absInt64(a int64) int64 ***REMOVED***
	if a < 0 ***REMOVED***
		return -a
	***REMOVED***
	return a
***REMOVED***

func (vlvc RampingVUsConfig) precalculateTheRequiredSteps(et *libWorker.ExecutionTuple, zeroEnd bool) int ***REMOVED***
	p := et.ScaleInt64(vlvc.StartVUs.Int64)
	var result int64
	result++ // for the first one

	if zeroEnd ***REMOVED***
		result++ // for the last one - this one can be more then needed
	***REMOVED***
	for _, stage := range vlvc.Stages ***REMOVED***
		stageEndVUs := et.ScaleInt64(stage.Target.Int64)
		if stage.Duration.Duration == 0 ***REMOVED***
			result++
		***REMOVED*** else ***REMOVED***
			result += absInt64(p - stageEndVUs)
		***REMOVED***
		p = stageEndVUs
	***REMOVED***
	return int(result)
***REMOVED***

// If the graceful ramp-downs are enabled, we need to reserve any VUs that may
// potentially have to finish running iterations when we're scaling their number
// down. This would prevent attempts from other executors to use them while the
// iterations are finishing up during their allotted gracefulRampDown periods.
//
// But we also need to be careful to not over-allocate more VUs than we actually
// need. We should never have more PlannedVUs than the max(startVUs,
// stage[n].target), even if we're quickly scaling VUs up and down multiple
// times, one after the other. In those cases, any previously reserved VUs
// finishing up interrupted iterations should be reused by the executor,
// instead of new ones being requested from the execution state.
//
// Here's an example with graceful ramp-down (i.e. "uninterruptible"
// iterations), where stars represent actively scheduled VUs and dots are used
// for VUs that are potentially finishing up iterations:
//
//	^
//	|
//
// VUs 6|  *..............................
//
//	5| ***.......*..............................
//	4|*****.....***.....**..............................
//	3|******...*****...***..............................
//	2|*******.*******.****..............................
//	1|***********************..............................
//	0--------------------------------------------------------> time(s)
//	  012345678901234567890123456789012345678901234567890123   (t%10)
//	  000000000011111111112222222222333333333344444444445555   (t/10)
//
// We start with 4 VUs, scale to 6, scale down to 1, scale up to 5, scale down
// to 1 again, scale up to 4, back to 1, and finally back down to 0. If our
// gracefulStop timeout was 30s (the default), then we'll stay with 6 PlannedVUs
// until t=32 in the test above, and the actual executor could run until t=52.
// See TestRampingVUsConfigExecutionPlanExample() for the above example
// as a unit test.
//
// The algorithm we use below to reserve VUs so that ramping-down VUs can finish
// their last iterations is pretty simple. It just traverses the raw execution
// steps and whenever there's a scaling down of VUs, it prevents the number of
// VUs from decreasing for the configured gracefulRampDown period.
//
// Finishing up the test, i.e. making sure we have a step with 0 VUs at time
// executorEndOffset, is not handled here. Instead GetExecutionRequirements()
// takes care of that. But to make its job easier, this method won't add any
// steps with an offset that's greater or equal to executorEndOffset.
func (vlvc RampingVUsConfig) reserveVUsForGracefulRampDowns( //nolint:funlen
	rawSteps []libWorker.ExecutionStep, executorEndOffset time.Duration,
) []libWorker.ExecutionStep ***REMOVED***
	rawStepsLen := len(rawSteps)
	gracefulRampDownPeriod := vlvc.GetGracefulRampDown()
	newSteps := []libWorker.ExecutionStep***REMOVED******REMOVED***

	lastPlannedVUs := uint64(0)
	lastDownwardSlopeStepNum := 0
	for rawStepNum := 0; rawStepNum < rawStepsLen; rawStepNum++ ***REMOVED***
		rawStep := rawSteps[rawStepNum]
		// Add the first step or any step where the number of planned VUs is
		// greater than the ones in the previous step. We don't need to worry
		// about reserving time for ramping-down VUs when the number of planned
		// VUs is growing. That's because the gracefulRampDown period is a fixed
		// value and any timeouts from early steps with fewer VUs will get
		// overshadowed by timeouts from latter steps with more VUs.
		if rawStepNum == 0 || rawStep.PlannedVUs > lastPlannedVUs ***REMOVED***
			newSteps = append(newSteps, rawStep)
			lastPlannedVUs = rawStep.PlannedVUs
			continue
		***REMOVED***

		// We simply skip steps with the same number of planned VUs
		if rawStep.PlannedVUs == lastPlannedVUs ***REMOVED***
			continue
		***REMOVED***

		// If we're here, we have a downward "slope" - the lastPlannedVUs are
		// more than the current rawStep's planned VUs. We're going to look
		// forward in time (up to gracefulRampDown) and inspect the rawSteps.
		// There are a 3 possibilities:
		//  - We find a new step within the gracefulRampDown period which has
		//    the same number of VUs or greater than lastPlannedVUs. Which
		//    means that we can just advance rawStepNum to that number and we
		//    don't need to worry about any of the raw steps in the middle!
		//    Both their planned VUs and their gracefulRampDown periods will
		//    be lower than what we're going to set from that new rawStep -
		//    we've basically found a new upward slope or equal value again.
		//  - We reach executorEndOffset, in which case we are done - we can't
		//    add any new steps, since those will be after the executor end
		//    offset.
		//  - We reach the end of the rawSteps, or we don't find any higher or
		//    equal steps to prevStep in the next gracefulRampDown period. So
		//    we'll simply add a new entry into newSteps with the values
		//    ***REMOVED***timeOffsetWithTimeout, rawStep.PlannedVUs***REMOVED***, in which
		//    timeOffsetWithTimeout = (prevStep.TimeOffset + gracefulRampDown),
		//    after which we'll continue with traversing the following rawSteps.
		//
		//  In this last case, we can also keep track of the downward slope. We
		//  can save the last index we've seen, to optimize future lookaheads by
		//  starting them from the furhest downward slope step index we've seen so
		//  far. This can be done because the gracefulRampDown is constant, so
		//  when we're on a downward slope, raw steps will always fall in the
		//  gracefulRampDown "shadow" of their preceding steps on the slope.

		skippedToNewRawStep := false
		timeOffsetWithTimeout := rawStep.TimeOffset + gracefulRampDownPeriod

		advStepStart := rawStepNum + 1
		if lastDownwardSlopeStepNum > advStepStart ***REMOVED***
			advStepStart = lastDownwardSlopeStepNum
		***REMOVED***

		wasRampingDown := true
		for advStepNum := advStepStart; advStepNum < rawStepsLen; advStepNum++ ***REMOVED***
			advStep := rawSteps[advStepNum]
			if advStep.TimeOffset > timeOffsetWithTimeout ***REMOVED***
				break
			***REMOVED***
			if advStep.PlannedVUs >= lastPlannedVUs ***REMOVED***
				rawStepNum = advStepNum - 1
				skippedToNewRawStep = true
				break
			***REMOVED***
			if wasRampingDown ***REMOVED***
				if rawSteps[advStepNum-1].PlannedVUs > advStep.PlannedVUs ***REMOVED***
					// Still ramping down
					lastDownwardSlopeStepNum = advStepNum
				***REMOVED*** else ***REMOVED***
					// No longer ramping down
					wasRampingDown = false
				***REMOVED***
			***REMOVED***
		***REMOVED***

		// Nothing more to do here, found a new "slope" with equal or grater
		// PlannedVUs in the gracefulRampDownPeriod window, so we go to it.
		if skippedToNewRawStep ***REMOVED***
			continue
		***REMOVED***

		// We've reached the absolute executor end offset, and we were already
		// on a downward "slope" (i.e. the previous planned VUs are more than
		// the current planned VUs), so nothing more we can do here.
		if timeOffsetWithTimeout >= executorEndOffset ***REMOVED***
			break
		***REMOVED***

		newSteps = append(newSteps, libWorker.ExecutionStep***REMOVED***
			TimeOffset: timeOffsetWithTimeout,
			PlannedVUs: rawStep.PlannedVUs,
		***REMOVED***)
		lastPlannedVUs = rawStep.PlannedVUs
	***REMOVED***

	return newSteps
***REMOVED***

// GetExecutionRequirements very dynamically reserves exactly the number of
// required VUs for this executor at every moment of the test.
//
// If gracefulRampDown is specified, it will also be taken into account, and the
// number of needed VUs to handle that will also be reserved. See the
// documentation of reserveVUsForGracefulRampDowns() for more details.
//
// On the other hand, gracefulStop is handled here. To facilitate it, we'll
// ensure that the last execution step will have 0 VUs and will be at time
// offset (sum(stages.Duration)+gracefulStop). Any steps that would've been
// added after it will be ignored. Thus:
//   - gracefulStop can be less than gracefulRampDown and can cut the graceful
//     ramp-down periods of the last VUs short.
//   - gracefulRampDown can be more than gracefulStop:
//   - If the user manually ramped down VUs at the end of the test (i.e. the
//     last stage's target is 0), then this will have no effect.
//   - If the last stage's target is more than 0, the VUs at the end of the
//     executor's life will have more time to finish their last iterations.
func (vlvc RampingVUsConfig) GetExecutionRequirements(et *libWorker.ExecutionTuple) []libWorker.ExecutionStep ***REMOVED***
	steps := vlvc.getRawExecutionSteps(et, false)

	executorEndOffset := sumStagesDuration(vlvc.Stages) + vlvc.GracefulStop.TimeDuration()
	// Handle graceful ramp-downs, if we have them
	if vlvc.GracefulRampDown.Duration > 0 ***REMOVED***
		steps = vlvc.reserveVUsForGracefulRampDowns(steps, executorEndOffset)
	***REMOVED***

	// If the last PlannedVUs value wasn't 0, add a last step with 0
	if steps[len(steps)-1].PlannedVUs != 0 ***REMOVED***
		steps = append(steps, libWorker.ExecutionStep***REMOVED***TimeOffset: executorEndOffset, PlannedVUs: 0***REMOVED***)
	***REMOVED***

	return steps
***REMOVED***

// NewExecutor creates a new RampingVUs executor
func (vlvc RampingVUsConfig) NewExecutor(es *libWorker.ExecutionState, logger *logrus.Entry) (libWorker.Executor, error) ***REMOVED***
	return &RampingVUs***REMOVED***
		BaseExecutor: NewBaseExecutor(vlvc, es, logger),
		config:       vlvc,
	***REMOVED***, nil
***REMOVED***

// HasWork reports whether there is any work to be done for the given execution segment.
func (vlvc RampingVUsConfig) HasWork(et *libWorker.ExecutionTuple) bool ***REMOVED***
	return libWorker.GetMaxPlannedVUs(vlvc.GetExecutionRequirements(et)) > 0
***REMOVED***

// RampingVUs handles the old "stages" execution configuration - it loops
// iterations with a variable number of VUs for the sum of all of the specified
// stages' duration.
type RampingVUs struct ***REMOVED***
	*BaseExecutor
	config RampingVUsConfig

	rawSteps, gracefulSteps []libWorker.ExecutionStep
***REMOVED***

// Make sure we implement the libWorker.Executor interface.
var _ libWorker.Executor = &RampingVUs***REMOVED******REMOVED***

// Init initializes the rampingVUs executor by precalculating the raw
// and graceful steps.
func (vlv *RampingVUs) Init(_ context.Context) error ***REMOVED***
	vlv.rawSteps = vlv.config.getRawExecutionSteps(
		vlv.executionState.ExecutionTuple, true,
	)
	vlv.gracefulSteps = vlv.config.GetExecutionRequirements(
		vlv.executionState.ExecutionTuple,
	)
	return nil
***REMOVED***

// Run constantly loops through as many iterations as possible on a variable
// number of VUs for the specified stages.
func (vlv *RampingVUs) Run(ctx context.Context, _ chan<- workerMetrics.SampleContainer, workerInfo *libWorker.WorkerInfo) error ***REMOVED***
	regularDuration, isFinal := libWorker.GetEndOffset(vlv.rawSteps)
	if !isFinal ***REMOVED***
		return fmt.Errorf("%s expected raw end offset at %s to be final", vlv.config.GetName(), regularDuration)
	***REMOVED***
	maxDuration, isFinal := libWorker.GetEndOffset(vlv.gracefulSteps)
	if !isFinal ***REMOVED***
		return fmt.Errorf("%s expected graceful end offset at %s to be final", vlv.config.GetName(), maxDuration)
	***REMOVED***
	waitOnProgressChannel := make(chan struct***REMOVED******REMOVED***)
	startTime, maxDurationCtx, regularDurationCtx, cancel := getDurationContexts(
		ctx, regularDuration, maxDuration-regularDuration,
	)
	defer func() ***REMOVED***
		cancel()
		<-waitOnProgressChannel
	***REMOVED***()

	maxVUs := libWorker.GetMaxPlannedVUs(vlv.gracefulSteps)

	vlv.logger.WithFields(logrus.Fields***REMOVED***
		"type":      vlv.config.GetType(),
		"startVUs":  vlv.config.GetStartVUs(vlv.executionState.ExecutionTuple),
		"maxVUs":    maxVUs,
		"duration":  regularDuration,
		"numStages": len(vlv.config.Stages),
	***REMOVED***).Debug("Starting executor run...")

	runState := &rampingVUsRunState***REMOVED***
		executor:       vlv,
		vuHandles:      make([]*vuHandle, maxVUs),
		maxVUs:         maxVUs,
		activeVUsCount: new(int64),
		started:        startTime,
		runIteration:   getIterationRunner(vlv.executionState, vlv.logger),
	***REMOVED***

	progressFn := runState.makeProgressFn(regularDuration)
	maxDurationCtx = libWorker.WithScenarioState(maxDurationCtx, &libWorker.ScenarioState***REMOVED***
		Name:       vlv.config.Name,
		Executor:   vlv.config.Type,
		StartTime:  runState.started,
		ProgressFn: progressFn,
	***REMOVED***)
	vlv.progress.Modify(pb.WithProgress(progressFn))
	go func() ***REMOVED***
		trackProgress(ctx, maxDurationCtx, regularDurationCtx, vlv, progressFn)
		close(waitOnProgressChannel)
	***REMOVED***()
	defer runState.wg.Wait()
	// this will populate stopped VUs and run runLoopsIfPossible on each VU
	// handle in a new goroutine
	runState.runLoopsIfPossible(maxDurationCtx, cancel)

	var (
		handleNewMaxAllowedVUs = runState.maxAllowedVUsHandlerStrategy()
		handleNewScheduledVUs  = runState.scheduledVUsHandlerStrategy()
	)
	handledGracefulSteps := runState.iterateSteps(
		ctx,
		handleNewMaxAllowedVUs,
		handleNewScheduledVUs,
	)
	go runState.runRemainingGracefulSteps(
		ctx,
		handleNewMaxAllowedVUs,
		handledGracefulSteps,
	)
	return nil
***REMOVED***

// rampingVUsRunState is created and initialized by the Run() method
// of the ramping VUs executor. It is used to track and modify various
// details of the execution.
type rampingVUsRunState struct ***REMOVED***
	executor       *RampingVUs
	vuHandles      []*vuHandle // handles for manipulating and tracking all of the VUs
	maxVUs         uint64      // the scaled number of initially configured MaxVUs
	activeVUsCount *int64      // the current number of active VUs, used only for the progress display
	started        time.Time
	wg             sync.WaitGroup

	runIteration func(context.Context, libWorker.ActiveVU) bool // a helper closure function that runs a single iteration
***REMOVED***

func (rs *rampingVUsRunState) makeProgressFn(regular time.Duration) (progressFn func() (float64, []string)) ***REMOVED***
	vusFmt := pb.GetFixedLengthIntFormat(int64(rs.maxVUs))
	regularDuration := pb.GetFixedLengthDuration(regular, regular)

	return func() (float64, []string) ***REMOVED***
		spent := time.Since(rs.started)
		cur := atomic.LoadInt64(rs.activeVUsCount)
		progVUs := fmt.Sprintf(vusFmt+"/"+vusFmt+" VUs", cur, rs.maxVUs)
		if spent > regular ***REMOVED***
			return 1, []string***REMOVED***progVUs, regular.String()***REMOVED***
		***REMOVED***
		status := pb.GetFixedLengthDuration(spent, regular) + "/" + regularDuration
		return float64(spent) / float64(regular), []string***REMOVED***progVUs, status***REMOVED***
	***REMOVED***
***REMOVED***

func (rs *rampingVUsRunState) runLoopsIfPossible(ctx context.Context, cancel func()) ***REMOVED***
	getVU := func() (libWorker.InitializedVU, error) ***REMOVED***
		pvu, err := rs.executor.executionState.GetPlannedVU(rs.executor.logger, false)
		if err != nil ***REMOVED***
			rs.executor.logger.WithError(err).Error("Cannot get a VU from the buffer")
			cancel()
			return pvu, err
		***REMOVED***
		rs.wg.Add(1)
		atomic.AddInt64(rs.activeVUsCount, 1)
		rs.executor.executionState.ModCurrentlyActiveVUsCount(+1)
		return pvu, err
	***REMOVED***
	returnVU := func(initVU libWorker.InitializedVU) ***REMOVED***
		rs.executor.executionState.ReturnVU(initVU, false)
		atomic.AddInt64(rs.activeVUsCount, -1)
		rs.wg.Done()
		rs.executor.executionState.ModCurrentlyActiveVUsCount(-1)
	***REMOVED***
	for i := uint64(0); i < rs.maxVUs; i++ ***REMOVED***
		rs.vuHandles[i] = newStoppedVUHandle(
			ctx, getVU, returnVU, rs.executor.nextIterationCounters,
			&rs.executor.config.BaseConfig, rs.executor.logger.WithField("vuNum", i))
		go rs.vuHandles[i].runLoopsIfPossible(rs.runIteration)
	***REMOVED***
***REMOVED***

// iterateSteps iterates over rawSteps and gracefulSteps in order according to
// their TimeOffsets, prioritizing rawSteps. It stops iterating once rawSteps
// are over. And it returns the number of handled gracefulSteps.
func (rs *rampingVUsRunState) iterateSteps(
	ctx context.Context,
	handleNewMaxAllowedVUs, handleNewScheduledVUs func(libWorker.ExecutionStep),
) (handledGracefulSteps int) ***REMOVED***
	wait := waiter(ctx, rs.started)
	i, j := 0, 0
	for i != len(rs.executor.rawSteps) ***REMOVED***
		r, g := rs.executor.rawSteps[i], rs.executor.gracefulSteps[j]
		if g.TimeOffset < r.TimeOffset ***REMOVED***
			if wait(g.TimeOffset) ***REMOVED***
				break
			***REMOVED***
			handleNewMaxAllowedVUs(g)
			j++
		***REMOVED*** else ***REMOVED***
			if wait(r.TimeOffset) ***REMOVED***
				break
			***REMOVED***
			handleNewScheduledVUs(r)
			i++
		***REMOVED***
	***REMOVED***
	return j
***REMOVED***

// runRemainingGracefulSteps runs the remaining gracefulSteps concurrently
// before the gracefulStop timeout period stops VUs.
//
// This way we will have run the gracefulSteps at the same time while
// waiting for the VUs to finish.
//
// (gracefulStop = maxDuration-regularDuration)
func (rs *rampingVUsRunState) runRemainingGracefulSteps(
	ctx context.Context,
	handleNewMaxAllowedVUs func(libWorker.ExecutionStep),
	handledGracefulSteps int,
) ***REMOVED***
	wait := waiter(ctx, rs.started)
	for _, s := range rs.executor.gracefulSteps[handledGracefulSteps:] ***REMOVED***
		if wait(s.TimeOffset) ***REMOVED***
			return
		***REMOVED***
		handleNewMaxAllowedVUs(s)
	***REMOVED***
***REMOVED***

func (rs *rampingVUsRunState) maxAllowedVUsHandlerStrategy() func(libWorker.ExecutionStep) ***REMOVED***
	var cur uint64 // current number of planned graceful VUs
	return func(graceful libWorker.ExecutionStep) ***REMOVED***
		pv := graceful.PlannedVUs
		for ; pv < cur; cur-- ***REMOVED***
			rs.vuHandles[cur-1].hardStop()
		***REMOVED***
		cur = pv
	***REMOVED***
***REMOVED***

func (rs *rampingVUsRunState) scheduledVUsHandlerStrategy() func(libWorker.ExecutionStep) ***REMOVED***
	var cur uint64 // current number of planned raw VUs
	return func(raw libWorker.ExecutionStep) ***REMOVED***
		pv := raw.PlannedVUs
		for ; cur < pv; cur++ ***REMOVED***
			_ = rs.vuHandles[cur].start() // TODO: handle the error
		***REMOVED***
		for ; pv < cur; cur-- ***REMOVED***
			rs.vuHandles[cur-1].gracefulStop()
		***REMOVED***
	***REMOVED***
***REMOVED***

// waiter returns a function that will sleep/wait for the required time since the startTime and then
// return. If the context was done before that it will return true otherwise it will return false
// TODO use elsewhere
// TODO set start here?
// TODO move it to a struct type or something and benchmark if that makes a difference
func waiter(ctx context.Context, start time.Time) func(offset time.Duration) bool ***REMOVED***
	timer := time.NewTimer(time.Hour * 24)
	return func(offset time.Duration) bool ***REMOVED***
		diff := offset - time.Since(start)
		if diff > 0 ***REMOVED*** // wait until time of event arrives // TODO have a mininum
			timer.Reset(diff)
			select ***REMOVED***
			case <-ctx.Done():
				return true // exit if context is cancelled
			case <-timer.C:
				// now we do a step
			***REMOVED***
		***REMOVED***
		return false
	***REMOVED***
***REMOVED***

func (vlvc RampingVUsConfig) GetMaxExecutorVUs() int64 ***REMOVED***
	// Lop through stages and find the max number of VUs
	maxVUs := int64(vlvc.StartVUs.ValueOrZero())

	for _, stage := range vlvc.Stages ***REMOVED***
		if stage.Target.ValueOrZero() > maxVUs ***REMOVED***
			maxVUs = stage.Target.ValueOrZero()
		***REMOVED***
	***REMOVED***

	return maxVUs
***REMOVED***

func (vlvc RampingVUsConfig) ScaleOptions(subFraction float32) ***REMOVED***
	if vlvc.StartVUs.Valid ***REMOVED***
		vlvc.StartVUs.Int64 = int64(float32(vlvc.StartVUs.Int64) * subFraction)
	***REMOVED***

	for i := range vlvc.Stages ***REMOVED***
		if vlvc.Stages[i].Target.Valid ***REMOVED***
			vlvc.Stages[i].Target.Int64 = int64(float32(vlvc.Stages[i].Target.Int64) * subFraction)
		***REMOVED***
	***REMOVED***
***REMOVED***
