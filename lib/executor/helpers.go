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
	"math/big"
	"time"

	"github.com/loadimpact/k6/ui/pb"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/stats"
	"github.com/sirupsen/logrus"
)

func sumStagesDuration(stages []Stage) (result time.Duration) ***REMOVED***
	for _, s := range stages ***REMOVED***
		result += time.Duration(s.Duration.Duration)
	***REMOVED***
	return
***REMOVED***

func getStagesUnscaledMaxTarget(unscaledStartValue int64, stages []Stage) int64 ***REMOVED***
	max := unscaledStartValue
	for _, s := range stages ***REMOVED***
		if s.Target.Int64 > max ***REMOVED***
			max = s.Target.Int64
		***REMOVED***
	***REMOVED***
	return max
***REMOVED***

// A helper function to avoid code duplication
func validateStages(stages []Stage) []error ***REMOVED***
	var errors []error
	if len(stages) == 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("at least one stage has to be specified"))
	***REMOVED*** else ***REMOVED***
		for i, s := range stages ***REMOVED***
			stageNum := i + 1
			if !s.Duration.Valid ***REMOVED***
				errors = append(errors, fmt.Errorf("stage %d doesn't have a duration", stageNum))
			***REMOVED*** else if s.Duration.Duration < 0 ***REMOVED***
				errors = append(errors, fmt.Errorf("the duration for stage %d shouldn't be negative", stageNum))
			***REMOVED***
			if !s.Target.Valid ***REMOVED***
				errors = append(errors, fmt.Errorf("stage %d doesn't have a target", stageNum))
			***REMOVED*** else if s.Target.Int64 < 0 ***REMOVED***
				errors = append(errors, fmt.Errorf("the target for stage %d shouldn't be negative", stageNum))
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return errors
***REMOVED***

// getIterationRunner is a helper function that returns an iteration executor
// closure. It takes care of updating metrics, execution state statistics, and
// warning messages.
func getIterationRunner(
	executionState *lib.ExecutionState, logger *logrus.Entry, out chan<- stats.SampleContainer,
) func(context.Context, lib.VU) ***REMOVED***
	return func(ctx context.Context, vu lib.VU) ***REMOVED***
		err := vu.RunOnce(ctx)

		//TODO: track (non-ramp-down) errors from script iterations as a metric,
		// and have a default threshold that will abort the script when the error
		// rate exceeds a certain percentage

		select ***REMOVED***
		case <-ctx.Done():
			// Don't log errors or emit iterations metrics from cancelled iterations
			executionState.AddInterruptedIterations(1)
		default:
			if err != nil ***REMOVED***
				if s, ok := err.(fmt.Stringer); ok ***REMOVED***
					logger.Error(s.String())
				***REMOVED*** else ***REMOVED***
					logger.Error(err.Error())
				***REMOVED***
				//TODO: investigate context cancelled errors
			***REMOVED***

			out <- stats.Sample***REMOVED***
				Time:   time.Now(),
				Metric: metrics.Iterations,
				Value:  1,
				Tags:   executionState.Options.RunTags,
			***REMOVED***
			executionState.AddFullIterations(1)
		***REMOVED***
	***REMOVED***
***REMOVED***

// getDurationContexts is used to create sub-contexts that can restrict a
// executor to only run for its allotted time.
//
// If the executor doesn't have a graceful stop period for iterations, then
// both returned sub-contexts will be the same one, with a timeout equal to
// supplied regular executor duration.
//
// But if a graceful stop is enabled, then the first returned context (and the
// cancel func) will be for the "outer" sub-context. Its timeout will include
// both the regular duration and the specified graceful stop period. The second
// context will be a sub-context of the first one and its timeout will include
// only the regular duration.
//
// In either case, the usage of these contexts should be like this:
//  - As long as the regDurationCtx isn't done, new iterations can be started.
//  - After regDurationCtx is done, no new iterations should be started; every
//    VU that finishes an iteration from now on can be returned to the buffer
//    pool in the ExecutionState struct.
//  - After maxDurationCtx is done, any VUs with iterations will be
//    interrupted by the context's closing and will be returned to the buffer.
//  - If you want to interrupt the execution of all VUs prematurely (e.g. there
//    was an error or something like that), trigger maxDurationCancel().
//  - If the whole test is aborted, the parent context will be cancelled, so
//    that will also cancel these contexts, thus the "general abort" case is
//    handled transparently.
func getDurationContexts(parentCtx context.Context, regularDuration, gracefulStop time.Duration) (
	startTime time.Time, maxDurationCtx, regDurationCtx context.Context, maxDurationCancel func(),
) ***REMOVED***
	startTime = time.Now()
	maxEndTime := startTime.Add(regularDuration + gracefulStop)

	maxDurationCtx, maxDurationCancel = context.WithDeadline(parentCtx, maxEndTime)
	if gracefulStop == 0 ***REMOVED***
		return startTime, maxDurationCtx, maxDurationCtx, maxDurationCancel
	***REMOVED***
	regDurationCtx, _ = context.WithDeadline(maxDurationCtx, startTime.Add(regularDuration)) //nolint:govet
	return startTime, maxDurationCtx, regDurationCtx, maxDurationCancel
***REMOVED***

// trackProgress is a helper function that monitors certain end-events in a
// executor and updates its progressbar accordingly.
func trackProgress(
	parentCtx, maxDurationCtx, regDurationCtx context.Context,
	sched lib.Executor, snapshot func() (float64, string),
) ***REMOVED***
	progressBar := sched.GetProgress()
	logger := sched.GetLogger()

	<-regDurationCtx.Done() // Wait for the regular context to be over
	gracefulStop := sched.GetConfig().GetGracefulStop()
	if parentCtx.Err() == nil && gracefulStop > 0 ***REMOVED***
		p, right := snapshot()
		logger.WithField("gracefulStop", gracefulStop).Debug(
			"Regular duration is done, waiting for iterations to gracefully finish",
		)
		progressBar.Modify(pb.WithConstProgress(p, right+", gracefully stopping..."))
	***REMOVED***

	<-maxDurationCtx.Done()
	p, right := snapshot()
	select ***REMOVED***
	case <-parentCtx.Done():
		progressBar.Modify(pb.WithConstProgress(p, right+" interrupted!"))
	default:
		progressBar.Modify(pb.WithConstProgress(p, right+" done!"))
	***REMOVED***
***REMOVED***

// getScaledArrivalRate returns a rational number containing the scaled value of
// the given rate over the given period. This should generally be the first
// function that's called, before we do any calculations with the users-supplied
// rates in the arrival-rate executors.
func getScaledArrivalRate(es *lib.ExecutionSegment, rate int64, period time.Duration) *big.Rat ***REMOVED***
	return es.InPlaceScaleRat(big.NewRat(rate, int64(period)))
***REMOVED***

// just a cached value to avoid allocating it every getTickerPeriod() call
var zero = big.NewInt(0) //nolint:gochecknoglobals

// getTickerPeriod is just a helper function that returns the ticker interval*
// we need for given arrival-rate parameters.
//
// It's possible for this function to return a zero duration (i.e. valid=false)
// and 0 isn't a valid ticker period. This happens so we don't divide by 0 when
// the arrival-rate period is 0. This case has to be handled separately.
func getTickerPeriod(scaledArrivalRate *big.Rat) types.NullDuration ***REMOVED***
	if scaledArrivalRate.Num().Cmp(zero) == 0 ***REMOVED***
		return types.NewNullDuration(0, false)
	***REMOVED***
	// Basically, the ticker rate is time.Duration(1/arrivalRate). Considering
	// that time.Duration is represented as int64 nanoseconds, no meaningful
	// precision is likely to be lost here...
	result, _ := new(big.Rat).SetFrac(scaledArrivalRate.Denom(), scaledArrivalRate.Num()).Float64()
	return types.NewNullDuration(time.Duration(result), true)
***REMOVED***

// getArrivalRatePerSec returns the iterations per second rate.
func getArrivalRatePerSec(scaledArrivalRate *big.Rat) *big.Rat ***REMOVED***
	perSecRate := big.NewRat(int64(time.Second), 1)
	return perSecRate.Mul(perSecRate, scaledArrivalRate)
***REMOVED***
