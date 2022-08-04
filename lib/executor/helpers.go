package executor

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/sirupsen/logrus"

	"go.k6.io/k6/errext"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/types"
	"go.k6.io/k6/ui/pb"
)

func sumStagesDuration(stages []Stage) (result time.Duration) ***REMOVED***
	for _, s := range stages ***REMOVED***
		result += s.Duration.TimeDuration()
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
		return errors
	***REMOVED***

	for i, s := range stages ***REMOVED***
		stageNum := i + 1
		if !s.Duration.Valid ***REMOVED***
			errors = append(errors, fmt.Errorf("stage %d doesn't have a duration", stageNum))
		***REMOVED*** else if s.Duration.Duration < 0 ***REMOVED***
			errors = append(errors, fmt.Errorf("the duration for stage %d can't be negative", stageNum))
		***REMOVED***
		if !s.Target.Valid ***REMOVED***
			errors = append(errors, fmt.Errorf("stage %d doesn't have a target", stageNum))
		***REMOVED*** else if s.Target.Int64 < 0 ***REMOVED***
			errors = append(errors, fmt.Errorf("the target for stage %d can't be negative", stageNum))
		***REMOVED***
	***REMOVED***
	return errors
***REMOVED***

// cancelKey is the key used to store the cancel function for the context of an
// executor. This is a work around to avoid excessive changes for the ability of
// nested functions to cancel the passed context.
type cancelKey struct***REMOVED******REMOVED***

type cancelExec struct ***REMOVED***
	cancel context.CancelFunc
	reason error
***REMOVED***

// Context returns context.Context that can be cancelled by calling
// CancelExecutorContext. Use this to initialize context that will be passed to
// executors.
//
// This allows executors to globally halt any executions that uses this context.
// Example use case is when a script calls test.abort().
func Context(ctx context.Context) context.Context ***REMOVED***
	ctx, cancel := context.WithCancel(ctx)
	return context.WithValue(ctx, cancelKey***REMOVED******REMOVED***, &cancelExec***REMOVED***cancel: cancel***REMOVED***)
***REMOVED***

// cancelExecutorContext cancels executor context found in ctx, ctx can be a
// child of a context that was created with Context function.
func cancelExecutorContext(ctx context.Context, err error) ***REMOVED***
	if x := ctx.Value(cancelKey***REMOVED******REMOVED***); x != nil ***REMOVED***
		if v, ok := x.(*cancelExec); ok ***REMOVED***
			v.reason = err
			v.cancel()
		***REMOVED***
	***REMOVED***
***REMOVED***

// CancelReason returns a reason the executor context was cancelled. This will
// return nil if ctx is not an executor context(ctx or any of its parents was
// never created by Context function).
func CancelReason(ctx context.Context) error ***REMOVED***
	if x := ctx.Value(cancelKey***REMOVED******REMOVED***); x != nil ***REMOVED***
		if v, ok := x.(*cancelExec); ok ***REMOVED***
			return v.reason
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// handleInterrupt returns true if err is InterruptError and if so it
// cancels the executor context passed with ctx.
func handleInterrupt(ctx context.Context, err error) bool ***REMOVED***
	if err != nil ***REMOVED***
		if errext.IsInterruptError(err) ***REMOVED***
			cancelExecutorContext(ctx, err)
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// getIterationRunner is a helper function that returns an iteration executor
// closure. It takes care of updating the execution state statistics and
// warning messages. And returns whether a full iteration was finished or not
//
// TODO: emit the end-of-test iteration metrics here (https://github.com/k6io/k6/issues/1250)
func getIterationRunner(
	executionState *lib.ExecutionState, logger *logrus.Entry,
) func(context.Context, lib.ActiveVU) bool ***REMOVED***
	return func(ctx context.Context, vu lib.ActiveVU) bool ***REMOVED***
		err := vu.RunOnce()

		// TODO: track (non-ramp-down) errors from script iterations as a metric,
		// and have a default threshold that will abort the script when the error
		// rate exceeds a certain percentage

		select ***REMOVED***
		case <-ctx.Done():
			// Don't log errors or emit iterations metrics from cancelled iterations
			executionState.AddInterruptedIterations(1)
			return false
		default:
			if err != nil ***REMOVED***
				if handleInterrupt(ctx, err) ***REMOVED***
					executionState.AddInterruptedIterations(1)
					return false
				***REMOVED***

				var exception errext.Exception
				if errors.As(err, &exception) ***REMOVED***
					// TODO don't count this as a full iteration?
					logger.WithField("source", "stacktrace").Error(exception.StackTrace())
				***REMOVED*** else ***REMOVED***
					logger.Error(err.Error())
				***REMOVED***
				// TODO: investigate context cancelled errors
			***REMOVED***

			// TODO: move emission of end-of-iteration metrics here?
			executionState.AddFullIterations(1)
			return true
		***REMOVED***
	***REMOVED***
***REMOVED***

// getDurationContexts is used to create sub-contexts that can restrict an
// executor to only run for its allotted time.
//
// If the executor doesn't have a graceful stop period for iterations, then
// both returned sub-contexts will be the same one, with a timeout equal to
// the supplied regular executor duration.
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

// trackProgress is a helper function that monitors certain end-events in an
// executor and updates its progressbar accordingly.
func trackProgress(
	parentCtx, maxDurationCtx, regDurationCtx context.Context,
	exec lib.Executor, snapshot func() (float64, []string),
) ***REMOVED***
	progressBar := exec.GetProgress()
	logger := exec.GetLogger()

	<-regDurationCtx.Done() // Wait for the regular context to be over
	gracefulStop := exec.GetConfig().GetGracefulStop()
	if parentCtx.Err() == nil && gracefulStop > 0 ***REMOVED***
		p, right := snapshot()
		logger.WithField("gracefulStop", gracefulStop).Debug(
			"Regular duration is done, waiting for iterations to gracefully finish",
		)
		progressBar.Modify(
			pb.WithStatus(pb.Stopping),
			pb.WithConstProgress(p, right...),
		)
	***REMOVED***

	<-maxDurationCtx.Done()
	p, right := snapshot()
	constProg := pb.WithConstProgress(p, right...)
	select ***REMOVED***
	case <-parentCtx.Done():
		progressBar.Modify(pb.WithStatus(pb.Interrupted), constProg)
	default:
		status := pb.WithStatus(pb.Done)
		if p < 1 ***REMOVED***
			status = pb.WithStatus(pb.Interrupted)
		***REMOVED***
		progressBar.Modify(status, constProg)
	***REMOVED***
***REMOVED***

// getScaledArrivalRate returns a rational number containing the scaled value of
// the given rate over the given period. This should generally be the first
// function that's called, before we do any calculations with the users-supplied
// rates in the arrival-rate executors.
func getScaledArrivalRate(es *lib.ExecutionSegment, rate int64, period time.Duration) *big.Rat ***REMOVED***
	return es.InPlaceScaleRat(big.NewRat(rate, int64(period)))
***REMOVED***

// getTickerPeriod is just a helper function that returns the ticker interval
// we need for given arrival-rate parameters.
//
// It's possible for this function to return a zero duration (i.e. valid=false)
// and 0 isn't a valid ticker period. This happens so we don't divide by 0 when
// the arrival-rate period is 0. This case has to be handled separately.
func getTickerPeriod(scaledArrivalRate *big.Rat) types.NullDuration ***REMOVED***
	if scaledArrivalRate.Sign() == 0 ***REMOVED***
		return types.NewNullDuration(0, false)
	***REMOVED***
	// Basically, the ticker rate is time.Duration(1/arrivalRate). Considering
	// that time.Duration is represented as int64 nanoseconds, no meaningful
	// precision is likely to be lost here...
	result, _ := new(big.Rat).Inv(scaledArrivalRate).Float64()
	return types.NewNullDuration(time.Duration(result), true)
***REMOVED***

// getArrivalRatePerSec returns the iterations per second rate.
func getArrivalRatePerSec(scaledArrivalRate *big.Rat) *big.Rat ***REMOVED***
	perSecRate := big.NewRat(int64(time.Second), 1)
	return perSecRate.Mul(perSecRate, scaledArrivalRate)
***REMOVED***

// TODO: Refactor this, maybe move all scenario things to an embedded struct?
func getVUActivationParams(
	ctx context.Context, conf BaseConfig, deactivateCallback func(lib.InitializedVU),
	nextIterationCounters func() (uint64, uint64),
) *lib.VUActivationParams ***REMOVED***
	return &lib.VUActivationParams***REMOVED***
		RunContext:               ctx,
		Scenario:                 conf.Name,
		Exec:                     conf.GetExec(),
		Env:                      conf.GetEnv(),
		Tags:                     conf.GetTags(),
		DeactivateCallback:       deactivateCallback,
		GetNextIterationCounters: nextIterationCounters,
	***REMOVED***
***REMOVED***
