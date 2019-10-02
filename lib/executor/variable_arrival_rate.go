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
	"math"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/stats"
	"github.com/loadimpact/k6/ui/pb"
	"github.com/sirupsen/logrus"
	null "gopkg.in/guregu/null.v3"
)

const variableArrivalRateType = "variable-arrival-rate"

// How often we can make arrival rate adjustments when processing stages
// TODO: make configurable, in some bounds?
const minIntervalBetweenRateAdjustments = 250 * time.Millisecond

func init() ***REMOVED***
	lib.RegisterExecutorConfigType(
		variableArrivalRateType,
		func(name string, rawJSON []byte) (lib.ExecutorConfig, error) ***REMOVED***
			config := NewVariableArrivalRateConfig(name)
			err := lib.StrictJSONUnmarshal(rawJSON, &config)
			return config, err
		***REMOVED***,
	)
***REMOVED***

// VariableArrivalRateConfig stores config for the variable arrival-rate executor
type VariableArrivalRateConfig struct ***REMOVED***
	BaseConfig
	StartRate null.Int           `json:"startRate"`
	TimeUnit  types.NullDuration `json:"timeUnit"`
	Stages    []Stage            `json:"stages"`

	// Initialize `PreAllocatedVUs` number of VUs, and if more than that are needed,
	// they will be dynamically allocated, until `MaxVUs` is reached, which is an
	// absolutely hard limit on the number of VUs the executor will use
	PreAllocatedVUs null.Int `json:"preAllocatedVUs"`
	MaxVUs          null.Int `json:"maxVUs"`
***REMOVED***

// NewVariableArrivalRateConfig returns a VariableArrivalRateConfig with default values
func NewVariableArrivalRateConfig(name string) VariableArrivalRateConfig ***REMOVED***
	return VariableArrivalRateConfig***REMOVED***
		BaseConfig: NewBaseConfig(name, variableArrivalRateType),
		TimeUnit:   types.NewNullDuration(1*time.Second, false),
	***REMOVED***
***REMOVED***

// Make sure we implement the lib.ExecutorConfig interface
var _ lib.ExecutorConfig = &VariableArrivalRateConfig***REMOVED******REMOVED***

// GetPreAllocatedVUs is just a helper method that returns the scaled pre-allocated VUs.
func (varc VariableArrivalRateConfig) GetPreAllocatedVUs(es *lib.ExecutionSegment) int64 ***REMOVED***
	return es.Scale(varc.PreAllocatedVUs.Int64)
***REMOVED***

// GetMaxVUs is just a helper method that returns the scaled max VUs.
func (varc VariableArrivalRateConfig) GetMaxVUs(es *lib.ExecutionSegment) int64 ***REMOVED***
	return es.Scale(varc.MaxVUs.Int64)
***REMOVED***

// GetDescription returns a human-readable description of the executor options
func (varc VariableArrivalRateConfig) GetDescription(es *lib.ExecutionSegment) string ***REMOVED***
	//TODO: something better? always show iterations per second?
	maxVUsRange := fmt.Sprintf("maxVUs: %d", es.Scale(varc.PreAllocatedVUs.Int64))
	if varc.MaxVUs.Int64 > varc.PreAllocatedVUs.Int64 ***REMOVED***
		maxVUsRange += fmt.Sprintf("-%d", es.Scale(varc.MaxVUs.Int64))
	***REMOVED***
	maxUnscaledRate := getStagesUnscaledMaxTarget(varc.StartRate.Int64, varc.Stages)
	maxArrRatePerSec, _ := getArrivalRatePerSec(
		getScaledArrivalRate(es, maxUnscaledRate, time.Duration(varc.TimeUnit.Duration)),
	).Float64()

	return fmt.Sprintf("Up to %.2f iterations/s for %s over %d stages%s",
		maxArrRatePerSec, sumStagesDuration(varc.Stages),
		len(varc.Stages), varc.getBaseInfo(maxVUsRange))
***REMOVED***

// Validate makes sure all options are configured and valid
func (varc VariableArrivalRateConfig) Validate() []error ***REMOVED***
	errors := varc.BaseConfig.Validate()

	if varc.StartRate.Int64 < 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("the startRate value shouldn't be negative"))
	***REMOVED***

	if time.Duration(varc.TimeUnit.Duration) < 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("the timeUnit should be more than 0"))
	***REMOVED***

	errors = append(errors, validateStages(varc.Stages)...)

	if !varc.PreAllocatedVUs.Valid ***REMOVED***
		errors = append(errors, fmt.Errorf("the number of preAllocatedVUs isn't specified"))
	***REMOVED*** else if varc.PreAllocatedVUs.Int64 < 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("the number of preAllocatedVUs shouldn't be negative"))
	***REMOVED***

	if !varc.MaxVUs.Valid ***REMOVED***
		errors = append(errors, fmt.Errorf("the number of maxVUs isn't specified"))
	***REMOVED*** else if varc.MaxVUs.Int64 < varc.PreAllocatedVUs.Int64 ***REMOVED***
		errors = append(errors, fmt.Errorf("maxVUs shouldn't be less than preAllocatedVUs"))
	***REMOVED***

	return errors
***REMOVED***

// GetExecutionRequirements returns the number of required VUs to run the
// executor for its whole duration (disregarding any startTime), including the
// maximum waiting time for any iterations to gracefully stop. This is used by
// the execution scheduler in its VU reservation calculations, so it knows how
// many VUs to pre-initialize.
func (varc VariableArrivalRateConfig) GetExecutionRequirements(es *lib.ExecutionSegment) []lib.ExecutionStep ***REMOVED***
	return []lib.ExecutionStep***REMOVED***
		***REMOVED***
			TimeOffset:      0,
			PlannedVUs:      uint64(es.Scale(varc.PreAllocatedVUs.Int64)),
			MaxUnplannedVUs: uint64(es.Scale(varc.MaxVUs.Int64 - varc.PreAllocatedVUs.Int64)),
		***REMOVED***,
		***REMOVED***
			TimeOffset:      sumStagesDuration(varc.Stages) + time.Duration(varc.GracefulStop.Duration),
			PlannedVUs:      0,
			MaxUnplannedVUs: 0,
		***REMOVED***,
	***REMOVED***
***REMOVED***

type rateChange struct ***REMOVED***
	// At what time should the rate below be applied.
	timeOffset time.Duration
	// Equals 1/rate: if rate was "1/5s", then this value, which is intended to
	// be passed to time.NewTicker(), will be 5s. There's a special case when
	// the rate is 0, for which we'll set Valid=false. That's because 0 isn't a
	// valid ticker period and shouldn't be passed to time.NewTicker(). Instead,
	// an empty or stopped ticker should be used.
	tickerPeriod types.NullDuration
***REMOVED***

// A helper method to generate the plan how the rate changes would happen.
func (varc VariableArrivalRateConfig) getPlannedRateChanges(segment *lib.ExecutionSegment) []rateChange ***REMOVED***
	timeUnit := time.Duration(varc.TimeUnit.Duration)
	// Important note for accuracy: we must work with and scale only the
	// rational numbers, never the raw target values directly. It matters most
	// for the accuracy of the intermediate rate change values, but it's
	// important even here.
	//
	// Say we have a desired rate growth from 1/sec to 2/sec over 1 minute, and
	// we split the test into two segments of 20% and 80%. If we used the whole
	// numbers for scaling, then the instance executing the first segment won't
	// ever do even a single request, since scale(20%, 1) would be 0, whereas
	// the rational value for scale(20%, 1/sec) is 0.2/sec, or rather 1/5sec...
	currentRate := getScaledArrivalRate(segment, varc.StartRate.Int64, timeUnit)

	rateChanges := []rateChange***REMOVED******REMOVED***
	timeFromStart := time.Duration(0)

	for _, stage := range varc.Stages ***REMOVED***
		stageTargetRate := getScaledArrivalRate(segment, stage.Target.Int64, timeUnit)
		stageDuration := time.Duration(stage.Duration.Duration)

		if currentRate.Cmp(stageTargetRate) == 0 ***REMOVED***
			// We don't have to do anything but update the time offset
			// if the rate wasn't changed in this stage
			timeFromStart += stageDuration
			continue
		***REMOVED***

		// Handle 0-duration stages, i.e. instant rate jumps
		if stageDuration == 0 ***REMOVED***
			// check if the last set change is for the same time and overwrite it
			if len(rateChanges) > 0 && rateChanges[len(rateChanges)-1].timeOffset == timeFromStart ***REMOVED***
				rateChanges[len(rateChanges)-1].tickerPeriod = getTickerPeriod(stageTargetRate)
			***REMOVED*** else ***REMOVED***
				rateChanges = append(rateChanges, rateChange***REMOVED***
					timeOffset:   timeFromStart,
					tickerPeriod: getTickerPeriod(stageTargetRate),
				***REMOVED***)
			***REMOVED***
			currentRate = stageTargetRate
			continue
		***REMOVED***
		// Basically, find out how many regular intervals with size of at least
		// minIntervalBetweenRateAdjustments are in the stage's duration, and
		// then use that number to calculate the actual step. All durations have
		// nanosecond precision, so there isn't any actual loss of precision...
		stepNumber := (stageDuration / minIntervalBetweenRateAdjustments)
		if stepNumber > 1 ***REMOVED***
			stepInterval := stageDuration / stepNumber
			for t := stepInterval; ; t += stepInterval ***REMOVED***
				if stageDuration-t < minIntervalBetweenRateAdjustments ***REMOVED***
					break
				***REMOVED***

				rateDiff := new(big.Rat).Sub(stageTargetRate, currentRate)
				tArrivalRate := new(big.Rat).Add(
					currentRate,
					rateDiff.Mul(rateDiff, big.NewRat(int64(t), int64(stageDuration))),
				)

				rateChanges = append(rateChanges, rateChange***REMOVED***
					timeOffset:   timeFromStart + t,
					tickerPeriod: getTickerPeriod(tArrivalRate),
				***REMOVED***)
			***REMOVED***
		***REMOVED***
		timeFromStart += stageDuration
		rateChanges = append(rateChanges, rateChange***REMOVED***
			timeOffset:   timeFromStart,
			tickerPeriod: getTickerPeriod(stageTargetRate),
		***REMOVED***)
		currentRate = stageTargetRate
	***REMOVED***

	return rateChanges
***REMOVED***

// NewExecutor creates a new VariableArrivalRate executor
func (varc VariableArrivalRateConfig) NewExecutor(
	es *lib.ExecutionState, logger *logrus.Entry,
) (lib.Executor, error) ***REMOVED***
	return VariableArrivalRate***REMOVED***
		BaseExecutor:       NewBaseExecutor(varc, es, logger),
		config:             varc,
		plannedRateChanges: varc.getPlannedRateChanges(es.Options.ExecutionSegment),
	***REMOVED***, nil
***REMOVED***

// VariableArrivalRate tries to execute a specific number of iterations for a
// specific period.
//TODO: combine with the ConstantArrivalRate?
type VariableArrivalRate struct ***REMOVED***
	*BaseExecutor
	config             VariableArrivalRateConfig
	plannedRateChanges []rateChange
***REMOVED***

// Make sure we implement the lib.Executor interface.
var _ lib.Executor = &VariableArrivalRate***REMOVED******REMOVED***

// streamRateChanges is a helper method that emits rate change events at their
// proper time.
func (varr VariableArrivalRate) streamRateChanges(ctx context.Context, startTime time.Time) <-chan rateChange ***REMOVED***
	ch := make(chan rateChange)
	go func() ***REMOVED***
		for _, step := range varr.plannedRateChanges ***REMOVED***
			offsetDiff := step.timeOffset - time.Since(startTime)
			if offsetDiff > 0 ***REMOVED*** // wait until time of event arrives
				select ***REMOVED***
				case <-ctx.Done():
					return // exit if context is cancelled
				case <-time.After(offsetDiff): //TODO: reuse a timer?
					// do nothing
				***REMOVED***
			***REMOVED***
			select ***REMOVED***
			case <-ctx.Done():
				return // exit if context is cancelled
			case ch <- step: // send the step
			***REMOVED***
		***REMOVED***
	***REMOVED***()
	return ch
***REMOVED***

// Run executes a variable number of iterations per second.
func (varr VariableArrivalRate) Run(ctx context.Context, out chan<- stats.SampleContainer) (err error) ***REMOVED*** //nolint:funlen
	segment := varr.executionState.Options.ExecutionSegment
	gracefulStop := varr.config.GetGracefulStop()
	duration := sumStagesDuration(varr.config.Stages)
	preAllocatedVUs := varr.config.GetPreAllocatedVUs(segment)
	maxVUs := varr.config.GetMaxVUs(segment)

	timeUnit := time.Duration(varr.config.TimeUnit.Duration)
	startArrivalRate := getScaledArrivalRate(segment, varr.config.StartRate.Int64, timeUnit)

	maxUnscaledRate := getStagesUnscaledMaxTarget(varr.config.StartRate.Int64, varr.config.Stages)
	maxArrivalRatePerSec, _ := getArrivalRatePerSec(getScaledArrivalRate(segment, maxUnscaledRate, timeUnit)).Float64()
	startTickerPeriod := getTickerPeriod(startArrivalRate)

	startTime, maxDurationCtx, regDurationCtx, cancel := getDurationContexts(ctx, duration, gracefulStop)
	defer cancel()
	ticker := &time.Ticker***REMOVED******REMOVED***
	if startTickerPeriod.Valid ***REMOVED***
		ticker = time.NewTicker(time.Duration(startTickerPeriod.Duration))
	***REMOVED***

	// Make sure the log and the progress bar have accurate information
	varr.logger.WithFields(logrus.Fields***REMOVED***
		"maxVUs": maxVUs, "preAllocatedVUs": preAllocatedVUs, "duration": duration, "numStages": len(varr.config.Stages),
		"startTickerPeriod": startTickerPeriod.Duration, "type": varr.config.GetType(),
	***REMOVED***).Debug("Starting executor run...")

	// Pre-allocate the VUs local shared buffer
	vus := make(chan lib.VU, maxVUs)

	initialisedVUs := uint64(0)
	// Make sure we put back planned and unplanned VUs back in the global
	// buffer, and as an extra incentive, this replaces a waitgroup.
	defer func() ***REMOVED***
		// no need for atomics, since initialisedVUs is mutated only in the select***REMOVED******REMOVED***
		for i := uint64(0); i < initialisedVUs; i++ ***REMOVED***
			varr.executionState.ReturnVU(<-vus, true)
		***REMOVED***
	***REMOVED***()

	// Get the pre-allocated VUs in the local buffer
	for i := int64(0); i < preAllocatedVUs; i++ ***REMOVED***
		vu, err := varr.executionState.GetPlannedVU(varr.logger, true)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		initialisedVUs++
		vus <- vu
	***REMOVED***

	tickerPeriod := new(int64)
	*tickerPeriod = int64(startTickerPeriod.Duration)

	fmtStr := pb.GetFixedLengthFloatFormat(maxArrivalRatePerSec, 2) + " iters/s, " +
		pb.GetFixedLengthIntFormat(maxVUs) + " out of " + pb.GetFixedLengthIntFormat(maxVUs) + " VUs active"
	progresFn := func() (float64, string) ***REMOVED***
		currentInitialisedVUs := atomic.LoadUint64(&initialisedVUs)
		currentTickerPeriod := atomic.LoadInt64(tickerPeriod)
		vusInBuffer := uint64(len(vus))

		itersPerSec := 0.0
		if currentTickerPeriod > 0 ***REMOVED***
			itersPerSec = float64(time.Second) / float64(currentTickerPeriod)
		***REMOVED***
		return math.Min(1, float64(time.Since(startTime))/float64(duration)), fmt.Sprintf(fmtStr,
			itersPerSec, currentInitialisedVUs-vusInBuffer, currentInitialisedVUs,
		)
	***REMOVED***
	varr.progress.Modify(pb.WithProgress(progresFn))
	go trackProgress(ctx, maxDurationCtx, regDurationCtx, varr, progresFn)

	regDurationDone := regDurationCtx.Done()
	runIterationBasic := getIterationRunner(varr.executionState, varr.logger, out)
	runIteration := func(vu lib.VU) ***REMOVED***
		runIterationBasic(maxDurationCtx, vu)
		vus <- vu
	***REMOVED***

	remainingUnplannedVUs := maxVUs - preAllocatedVUs
	rateChangesStream := varr.streamRateChanges(maxDurationCtx, startTime)

	for ***REMOVED***
		select ***REMOVED***
		case rateChange := <-rateChangesStream:
			newPeriod := rateChange.tickerPeriod
			ticker.Stop()
			if newPeriod.Valid ***REMOVED***
				ticker = time.NewTicker(time.Duration(newPeriod.Duration))
			***REMOVED***
			atomic.StoreInt64(tickerPeriod, int64(newPeriod.Duration))
		case <-ticker.C:
			select ***REMOVED***
			case vu := <-vus:
				// ideally, we get the VU from the buffer without any issues
				go runIteration(vu)
			default:
				if remainingUnplannedVUs == 0 ***REMOVED***
					//TODO: emit an error metric?
					varr.logger.Warningf("Insufficient VUs, reached %d active VUs and cannot allocate more", maxVUs)
					break
				***REMOVED***
				vu, err := varr.executionState.GetUnplannedVU(maxDurationCtx, varr.logger)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				remainingUnplannedVUs--
				atomic.AddUint64(&initialisedVUs, 1)
				go runIteration(vu)
			***REMOVED***
		case <-regDurationDone:
			return nil
		***REMOVED***
	***REMOVED***
***REMOVED***
