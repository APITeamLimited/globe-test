/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2019 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package executor

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	null "gopkg.in/guregu/null.v3"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/stats"
	"github.com/loadimpact/k6/ui/pb"
)

const variableLoopingVUsType = "variable-looping-vus"

func init() ***REMOVED***
	lib.RegisterExecutorConfigType(
		variableLoopingVUsType,
		func(name string, rawJSON []byte) (lib.ExecutorConfig, error) ***REMOVED***
			config := NewVariableLoopingVUsConfig(name)
			err := lib.StrictJSONUnmarshal(rawJSON, &config)
			return config, err
		***REMOVED***,
	)
***REMOVED***

// Stage contains
type Stage struct ***REMOVED***
	Duration types.NullDuration `json:"duration"`
	Target   null.Int           `json:"target"` // TODO: maybe rename this to endVUs? something else?
	//TODO: add a progression function?
***REMOVED***

// VariableLoopingVUsConfig stores the configuration for the stages executor
type VariableLoopingVUsConfig struct ***REMOVED***
	BaseConfig
	StartVUs         null.Int           `json:"startVUs"`
	Stages           []Stage            `json:"stages"`
	GracefulRampDown types.NullDuration `json:"gracefulRampDown"`
***REMOVED***

// NewVariableLoopingVUsConfig returns a VariableLoopingVUsConfig with its default values
func NewVariableLoopingVUsConfig(name string) VariableLoopingVUsConfig ***REMOVED***
	return VariableLoopingVUsConfig***REMOVED***
		BaseConfig:       NewBaseConfig(name, variableLoopingVUsType),
		StartVUs:         null.NewInt(1, false),
		GracefulRampDown: types.NewNullDuration(30*time.Second, false),
	***REMOVED***
***REMOVED***

// Make sure we implement the lib.ExecutorConfig interface
var _ lib.ExecutorConfig = &VariableLoopingVUsConfig***REMOVED******REMOVED***

// GetStartVUs is just a helper method that returns the scaled starting VUs.
func (vlvc VariableLoopingVUsConfig) GetStartVUs(et *lib.ExecutionTuple) int64 ***REMOVED***
	return et.ES.Scale(vlvc.StartVUs.Int64)
***REMOVED***

// GetGracefulRampDown is just a helper method that returns the graceful
// ramp-down period as a standard Go time.Duration value...
func (vlvc VariableLoopingVUsConfig) GetGracefulRampDown() time.Duration ***REMOVED***
	return time.Duration(vlvc.GracefulRampDown.Duration)
***REMOVED***

// GetDescription returns a human-readable description of the executor options
func (vlvc VariableLoopingVUsConfig) GetDescription(et *lib.ExecutionTuple) string ***REMOVED***
	maxVUs := et.ES.Scale(getStagesUnscaledMaxTarget(vlvc.StartVUs.Int64, vlvc.Stages))
	return fmt.Sprintf("Up to %d looping VUs for %s over %d stages%s",
		maxVUs, sumStagesDuration(vlvc.Stages), len(vlvc.Stages),
		vlvc.getBaseInfo(fmt.Sprintf("gracefulRampDown: %s", vlvc.GetGracefulRampDown())))
***REMOVED***

// Validate makes sure all options are configured and valid
func (vlvc VariableLoopingVUsConfig) Validate() []error ***REMOVED***
	errors := vlvc.BaseConfig.Validate()
	if vlvc.StartVUs.Int64 < 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("the number of start VUs shouldn't be negative"))
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
//    10|          *
//     9|         ***
//     8|        *****
//     7|       *******
//     6|      *********
//     5|     ***********
//     4|    *************
//     3|   ***************
//     2|  *****************
//     1| *******************
//     0------------------------> time(s)
//       01234567890123456789012   (t%10)
//       00000000001111111111222   (t/10)
//
// The chart for one of the execution segments would look like this:
//
// VUs  ^
//     5|         XXX
//     4|       XXXXXXX
//     3|     XXXXXXXXXXX
//     2|   XXXXXXXXXXXXXXX
//     1| XXXXXXXXXXXXXXXXXXX
//     0------------------------> time(s)
//       01234567890123456789012   (t%10)
//       00000000001111111111222   (t/10)
//
// And the chart for the other execution segment would look like this:
//
// VUs  ^
//     5|          Y
//     4|        YYYYY
//     3|      YYYYYYYYY
//     2|    YYYYYYYYYYYYY
//     1|  YYYYYYYYYYYYYYYYY
//     0------------------------> time(s)
//       01234567890123456789012   (t%10)
//       00000000001111111111222   (t/10)
//
// Notice the time offsets and the slower ramping up and down. All of that is
// because the sum of the two execution segments has to produce exactly the
// original shape, as if the test ran on a single machine:
//
// VUs  ^
//    10|          Y
//     9|         YYY
//     8|        YYYYY
//     7|       YYYYYYY
//     6|      YYYYYYYYY
//     5|     YYYYXXXYYYY
//     4|    YYYXXXXXXXYYY
//     3|   YYXXXXXXXXXXXYY
//     2|  YXXXXXXXXXXXXXXXY
//     1| XXXXXXXXXXXXXXXXXXX
//     0------------------------> time(s)
//       01234567890123456789012   (t%10)
//       00000000001111111111222   (t/10)
//
// More information: https://github.com/loadimpact/k6/issues/997#issuecomment-484416866
func (vlvc VariableLoopingVUsConfig) getRawExecutionSteps(et *lib.ExecutionTuple, zeroEnd bool) []lib.ExecutionStep ***REMOVED***
	var (
		timeTillEnd         time.Duration
		fromVUs             = vlvc.StartVUs.Int64
		start, offsets, lcd = et.GetStripedOffsets(et.ES)
		steps               = make([]lib.ExecutionStep, 0, vlvc.precalculateTheRequiredSteps(et, zeroEnd))
		index               = segmentedIndex***REMOVED***start: start, lcd: lcd, offsets: offsets***REMOVED***
	)

	// Reserve the scaled StartVUs at the beginning
	steps = append(steps, lib.ExecutionStep***REMOVED***TimeOffset: 0, PlannedVUs: uint64(index.goTo(fromVUs))***REMOVED***)
	addStep := func(timeOffset time.Duration, plannedVUs uint64) ***REMOVED***
		if steps[len(steps)-1].PlannedVUs != plannedVUs ***REMOVED***
			steps = append(steps, lib.ExecutionStep***REMOVED***TimeOffset: timeOffset, PlannedVUs: plannedVUs***REMOVED***)
		***REMOVED***
	***REMOVED***

	for _, stage := range vlvc.Stages ***REMOVED***
		stageEndVUs := stage.Target.Int64
		stageDuration := time.Duration(stage.Duration.Duration)
		timeTillEnd += stageDuration

		stageVUDiff := stageEndVUs - fromVUs
		if stageVUDiff == 0 ***REMOVED***
			continue
		***REMOVED***
		if stageDuration == 0 ***REMOVED***
			addStep(timeTillEnd, uint64(index.goTo(stageEndVUs)))
			fromVUs = stageEndVUs
			continue
		***REMOVED***

		// VU reservation for gracefully ramping down is handled as a
		// separate method: reserveVUsForGracefulRampDowns()
		if index.unscaled > stageEndVUs ***REMOVED*** // ramp down
			// here we don't want to emit for the equal to stageEndVUs as it doesn't go below it
			// it will just go to it
			for ; index.unscaled > stageEndVUs; index.prev() ***REMOVED***
				addStep(
					// this is the time that we should go up 1 if we are ramping up
					// but we are ramping down so we should go 1 down, but because we want to not
					// stop VUs immediately we stop it on the next unscaled VU's time
					timeTillEnd-time.Duration(int64(stageDuration)*(stageEndVUs-index.unscaled+1)/stageVUDiff),
					uint64(index.scaled-1),
				)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			for ; index.unscaled <= stageEndVUs; index.next() ***REMOVED***
				addStep(
					timeTillEnd-time.Duration(int64(stageDuration)*(stageEndVUs-index.unscaled)/stageVUDiff),
					uint64(index.scaled),
				)
			***REMOVED***
		***REMOVED***
		fromVUs = stageEndVUs
	***REMOVED***

	if zeroEnd && steps[len(steps)-1].PlannedVUs != 0 ***REMOVED***
		// If the last PlannedVUs value wasn't 0, add a last step with 0
		steps = append(steps, lib.ExecutionStep***REMOVED***TimeOffset: timeTillEnd, PlannedVUs: 0***REMOVED***)
	***REMOVED***
	return steps
***REMOVED***

type segmentedIndex struct ***REMOVED*** // TODO: rename ... although this is probably the best name so far :D
	start, lcd       int64
	offsets          []int64
	scaled, unscaled int64 // for both the first element(vu) is 1 not 0
***REMOVED***

// goes to the next scaled index and move the unscaled one accordingly
func (s *segmentedIndex) next() ***REMOVED***
	if s.scaled == 0 ***REMOVED*** // the 1 element(VU) is at the start
		s.unscaled += s.start + 1 // the first element of the start 0, but the here we need it to be 1 so we add 1
	***REMOVED*** else ***REMOVED*** // if we are not at the first element we need to go through the offsets, looping over them
		s.unscaled += s.offsets[int(s.scaled-1)%len(s.offsets)] // slice's index start at 0 ours start at 1
	***REMOVED***
	s.scaled++
***REMOVED***

// prev goest to the previous scaled value and sets the unscaled one accordingly
// calling prev when s.scaled == 0 is undefined
func (s *segmentedIndex) prev() ***REMOVED***
	if s.scaled == 1 ***REMOVED*** // we are the first need to go to the 0th element which means we need to remove the start
		s.unscaled -= s.start + 1 // this could've been just settign to 0
	***REMOVED*** else ***REMOVED*** // not at the first element - need to get the previously added offset so
		s.unscaled -= s.offsets[int(s.scaled-2)%len(s.offsets)] // slice's index start 0 our start at 1
	***REMOVED***
	s.scaled--
***REMOVED***

// goTo sets the scaled index to it's biggest value for which the corresponding unscaled index is
// is smaller or equal to value
func (s *segmentedIndex) goTo(value int64) int64 ***REMOVED*** // TODO optimize
	var gi int64
	// Because of the cyclical nature of the striping algorithm (with a cycle
	// length of LCD, the least common denominator), when scaling large values
	// (i.e. many multiples of the LCD), we can quickly calculate how many times
	// the cycle repeats.
	wholeCycles := (value / s.lcd)
	// So we can set some approximate initial values quickly, since we also know
	// precisely how many scaled values there are per cycle length.
	s.scaled = wholeCycles * int64(len(s.offsets))
	s.unscaled = wholeCycles*s.lcd + s.start + 1 // our indexes are from 1 the start is from 0
	// Approach the final value using the slow algorithm with the step by step loop
	// TODO: this can be optimized by another array with size offsets that instead of the offsets
	// from the previous is the offset from either 0 or start
	i := s.start
	for ; i < value%s.lcd; gi, i = gi+1, i+s.offsets[gi] ***REMOVED***
		s.scaled++
		s.unscaled += s.offsets[gi]
	***REMOVED***

	if gi > 0 ***REMOVED*** // there were more values after the wholecycles
		// the last offset actually shouldn't have been added
		s.unscaled -= s.offsets[gi-1]
	***REMOVED*** else if s.scaled > 0 ***REMOVED*** // we didn't actually have more values after the wholecycles but we still had some
		// in this case the unscaled value needs to move back by the last offset as it would've been
		// the one to get it from the value it needs to be to it's current one
		s.unscaled -= s.offsets[len(s.offsets)-1]
	***REMOVED***

	if s.scaled == 0 ***REMOVED***
		s.unscaled = 0 // we would've added the start and 1
	***REMOVED***

	return s.scaled
***REMOVED***

func absInt64(a int64) int64 ***REMOVED***
	if a < 0 ***REMOVED***
		return -a
	***REMOVED***
	return a
***REMOVED***

func (vlvc VariableLoopingVUsConfig) precalculateTheRequiredSteps(et *lib.ExecutionTuple, zeroEnd bool) int ***REMOVED***
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
//
//      ^
//      |
// VUs 6|  *..............................
//     5| ***.......*..............................
//     4|*****.....***.....**..............................
//     3|******...*****...***..............................
//     2|*******.*******.****..............................
//     1|***********************..............................
//     0--------------------------------------------------------> time(s)
//       012345678901234567890123456789012345678901234567890123   (t%10)
//       000000000011111111112222222222333333333344444444445555   (t/10)
//
// We start with 4 VUs, scale to 6, scale down to 1, scale up to 5, scale down
// to 1 again, scale up to 4, back to 1, and finally back down to 0. If our
// gracefulStop timeout was 30s (the default), then we'll stay with 6 PlannedVUs
// until t=32 in the test above, and the actual executor could run until t=52.
// See TestVariableLoopingVUsConfigExecutionPlanExample() for the above example
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
func (vlvc VariableLoopingVUsConfig) reserveVUsForGracefulRampDowns( //nolint:funlen
	rawSteps []lib.ExecutionStep, executorEndOffset time.Duration,
) []lib.ExecutionStep ***REMOVED***
	rawStepsLen := len(rawSteps)
	gracefulRampDownPeriod := vlvc.GetGracefulRampDown()
	newSteps := []lib.ExecutionStep***REMOVED******REMOVED***

	lastPlannedVUs := uint64(0)
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
		//    we'll simply try to add an entry into newSteps with the values
		//    ***REMOVED***prevStep.TimeOffset + gracefulRampDown, rawStep.PlannedVUs***REMOVED*** and
		//    we'll continue with traversing the following rawSteps.

		skippedToNewRawStep := false
		timeOffsetWithTimeout := rawStep.TimeOffset + gracefulRampDownPeriod

		for advStepNum := rawStepNum + 1; advStepNum < rawStepsLen; advStepNum++ ***REMOVED***
			advStep := rawSteps[advStepNum]
			if advStep.TimeOffset > timeOffsetWithTimeout ***REMOVED***
				break
			***REMOVED***
			if advStep.PlannedVUs >= lastPlannedVUs ***REMOVED***
				rawStepNum = advStepNum - 1
				skippedToNewRawStep = true
				break
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

		newSteps = append(newSteps, lib.ExecutionStep***REMOVED***
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
//     - If the user manually ramped down VUs at the end of the test (i.e. the
//       last stage's target is 0), then this will have no effect.
//     - If the last stage's target is more than 0, the VUs at the end of the
//       executor's life will have more time to finish their last iterations.
func (vlvc VariableLoopingVUsConfig) GetExecutionRequirements(et *lib.ExecutionTuple) []lib.ExecutionStep ***REMOVED***
	steps := vlvc.getRawExecutionSteps(et, false)

	executorEndOffset := sumStagesDuration(vlvc.Stages) + time.Duration(vlvc.GracefulStop.Duration)
	// Handle graceful ramp-downs, if we have them
	if vlvc.GracefulRampDown.Duration > 0 ***REMOVED***
		steps = vlvc.reserveVUsForGracefulRampDowns(steps, executorEndOffset)
	***REMOVED***

	// If the last PlannedVUs value wasn't 0, add a last step with 0
	if steps[len(steps)-1].PlannedVUs != 0 ***REMOVED***
		steps = append(steps, lib.ExecutionStep***REMOVED***TimeOffset: executorEndOffset, PlannedVUs: 0***REMOVED***)
	***REMOVED***

	return steps
***REMOVED***

// NewExecutor creates a new VariableLoopingVUs executor
func (vlvc VariableLoopingVUsConfig) NewExecutor(es *lib.ExecutionState, logger *logrus.Entry) (lib.Executor, error) ***REMOVED***
	return VariableLoopingVUs***REMOVED***
		BaseExecutor: NewBaseExecutor(vlvc, es, logger),
		config:       vlvc,
	***REMOVED***, nil
***REMOVED***

// HasWork reports whether there is any work to be done for the given execution segment.
func (vlvc VariableLoopingVUsConfig) HasWork(et *lib.ExecutionTuple) bool ***REMOVED***
	return lib.GetMaxPlannedVUs(vlvc.GetExecutionRequirements(et)) > 0
***REMOVED***

// VariableLoopingVUs handles the old "stages" execution configuration - it
// loops iterations with a variable number of VUs for the sum of all of the
// specified stages' duration.
type VariableLoopingVUs struct ***REMOVED***
	*BaseExecutor
	config VariableLoopingVUsConfig
***REMOVED***

// Make sure we implement the lib.Executor interface.
var _ lib.Executor = &VariableLoopingVUs***REMOVED******REMOVED***

// Run constantly loops through as many iterations as possible on a variable
// number of VUs for the specified stages.
//
// TODO: split up? since this does a ton of things, unfortunately I can't think
// of a less complex way to implement it (besides the old "increment by 100ms
// and see what happens)... :/ so maybe see how it can be split?
// nolint:funlen,gocognit
func (vlv VariableLoopingVUs) Run(ctx context.Context, out chan<- stats.SampleContainer) (err error) ***REMOVED***
	rawExecutionSteps := vlv.config.getRawExecutionSteps(vlv.executionState.ExecutionTuple, true)
	regularDuration, isFinal := lib.GetEndOffset(rawExecutionSteps)
	if !isFinal ***REMOVED***
		return fmt.Errorf("%s expected raw end offset at %s to be final", vlv.config.GetName(), regularDuration)
	***REMOVED***

	gracefulExecutionSteps := vlv.config.GetExecutionRequirements(vlv.executionState.ExecutionTuple)
	maxDuration, isFinal := lib.GetEndOffset(gracefulExecutionSteps)
	if !isFinal ***REMOVED***
		return fmt.Errorf("%s expected graceful end offset at %s to be final", vlv.config.GetName(), maxDuration)
	***REMOVED***
	maxVUs := lib.GetMaxPlannedVUs(gracefulExecutionSteps)
	gracefulStop := maxDuration - regularDuration

	startTime, maxDurationCtx, regDurationCtx, cancel := getDurationContexts(ctx, regularDuration, gracefulStop)
	defer cancel()

	// Make sure the log and the progress bar have accurate information
	vlv.logger.WithFields(logrus.Fields***REMOVED***
		"type": vlv.config.GetType(), "startVUs": vlv.config.GetStartVUs(vlv.executionState.ExecutionTuple), "maxVUs": maxVUs,
		"duration": regularDuration, "numStages": len(vlv.config.Stages)***REMOVED***,
	).Debug("Starting executor run...")

	activeVUsCount := new(int64)
	vusFmt := pb.GetFixedLengthIntFormat(int64(maxVUs))
	progresFn := func() (float64, []string) ***REMOVED***
		spent := time.Since(startTime)
		currentlyActiveVUs := atomic.LoadInt64(activeVUsCount)
		vus := fmt.Sprintf(vusFmt+"/"+vusFmt+" VUs", currentlyActiveVUs, maxVUs)
		if spent > regularDuration ***REMOVED***
			return 1, []string***REMOVED***vus, regularDuration.String()***REMOVED***
		***REMOVED***
		progVUs := fmt.Sprintf(vusFmt+"/"+vusFmt+" VUs", currentlyActiveVUs, maxVUs)
		progDur := fmt.Sprintf("%s/%s", pb.GetFixedLengthDuration(spent, regularDuration), regularDuration)
		return float64(spent) / float64(regularDuration), []string***REMOVED***progVUs, progDur***REMOVED***
	***REMOVED***
	vlv.progress.Modify(pb.WithProgress(progresFn))
	go trackProgress(ctx, maxDurationCtx, regDurationCtx, vlv, progresFn)

	// Actually schedule the VUs and iterations, likely the most complicated
	// executor among all of them...
	activeVUs := &sync.WaitGroup***REMOVED******REMOVED***
	defer activeVUs.Wait()

	runIteration := getIterationRunner(vlv.executionState, vlv.logger)
	getVU := func() (lib.InitializedVU, error) ***REMOVED***
		initVU, err := vlv.executionState.GetPlannedVU(vlv.logger, false)
		if err != nil ***REMOVED***
			cancel()
		***REMOVED*** else ***REMOVED***
			activeVUs.Add(1)
			atomic.AddInt64(activeVUsCount, 1)
		***REMOVED***
		return initVU, err
	***REMOVED***
	returnVU := func(initVU lib.InitializedVU) ***REMOVED***
		vlv.executionState.ReturnVU(initVU, false)
		atomic.AddInt64(activeVUsCount, -1)
		activeVUs.Done()
	***REMOVED***

	conf := vlv.GetConfig()
	execFn := conf.GetExec().ValueOrZero()
	env := conf.GetEnv()
	tags := conf.GetTags()
	vuHandles := make([]*vuHandle, maxVUs)
	for i := uint64(0); i < maxVUs; i++ ***REMOVED***
		vuHandle := newStoppedVUHandle(
			maxDurationCtx, getVU, returnVU, execFn, env,
			tags, vlv.logger.WithField("vuNum", i))
		go vuHandle.runLoopsIfPossible(runIteration)
		vuHandles[i] = vuHandle
	***REMOVED***

	rawStepEvents := lib.StreamExecutionSteps(ctx, startTime, rawExecutionSteps, true)
	gracefulLimitEvents := lib.StreamExecutionSteps(ctx, startTime, gracefulExecutionSteps, false)

	// 0 <= currentScheduledVUs <= currentMaxAllowedVUs <= maxVUs
	var currentScheduledVUs, currentMaxAllowedVUs uint64

	handleNewScheduledVUs := func(newScheduledVUs uint64) ***REMOVED***
		if newScheduledVUs > currentScheduledVUs ***REMOVED***
			for vuNum := currentScheduledVUs; vuNum < newScheduledVUs; vuNum++ ***REMOVED***
				vuHandles[vuNum].start()
				vlv.executionState.ModCurrentlyActiveVUsCount(+1)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			for vuNum := newScheduledVUs; vuNum < currentScheduledVUs; vuNum++ ***REMOVED***
				vuHandles[vuNum].gracefulStop()
				vlv.executionState.ModCurrentlyActiveVUsCount(-1)
			***REMOVED***
		***REMOVED***
		currentScheduledVUs = newScheduledVUs
	***REMOVED***

	handleNewMaxAllowedVUs := func(newMaxAllowedVUs uint64) ***REMOVED***
		if newMaxAllowedVUs < currentMaxAllowedVUs ***REMOVED***
			for vuNum := newMaxAllowedVUs; vuNum < currentMaxAllowedVUs; vuNum++ ***REMOVED***
				vuHandles[vuNum].hardStop()
			***REMOVED***
		***REMOVED***
		currentMaxAllowedVUs = newMaxAllowedVUs
	***REMOVED***

	handleAllRawSteps := func() bool ***REMOVED***
		for ***REMOVED***
			select ***REMOVED***
			case step, ok := <-rawStepEvents:
				if !ok ***REMOVED***
					return true
				***REMOVED***
				handleNewScheduledVUs(step.PlannedVUs)
			case step := <-gracefulLimitEvents:
				if step.PlannedVUs > currentMaxAllowedVUs ***REMOVED***
					// Handle the case where a value is read from the
					// gracefulLimitEvents channel before rawStepEvents
					handleNewScheduledVUs(step.PlannedVUs)
				***REMOVED***
				handleNewMaxAllowedVUs(step.PlannedVUs)
			case <-ctx.Done():
				return false
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if handleAllRawSteps() ***REMOVED***
		// Handle any remaining graceful stops
		go func() ***REMOVED***
			for ***REMOVED***
				select ***REMOVED***
				case step := <-gracefulLimitEvents:
					handleNewMaxAllowedVUs(step.PlannedVUs)
				case <-maxDurationCtx.Done():
					return
				***REMOVED***
			***REMOVED***
		***REMOVED***()
	***REMOVED***

	return nil
***REMOVED***
