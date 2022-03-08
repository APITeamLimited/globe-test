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
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/types"
	"go.k6.io/k6/metrics"
	"go.k6.io/k6/stats"
	"go.k6.io/k6/ui/pb"
)

const perVUIterationsType = "per-vu-iterations"

func init() ***REMOVED***
	lib.RegisterExecutorConfigType(perVUIterationsType, func(name string, rawJSON []byte) (lib.ExecutorConfig, error) ***REMOVED***
		config := NewPerVUIterationsConfig(name)
		err := lib.StrictJSONUnmarshal(rawJSON, &config)
		return config, err
	***REMOVED***)
***REMOVED***

// PerVUIterationsConfig stores the number of VUs iterations, as well as maxDuration settings
type PerVUIterationsConfig struct ***REMOVED***
	BaseConfig
	VUs         null.Int           `json:"vus"`
	Iterations  null.Int           `json:"iterations"`
	MaxDuration types.NullDuration `json:"maxDuration"`
***REMOVED***

// NewPerVUIterationsConfig returns a PerVUIterationsConfig with default values
func NewPerVUIterationsConfig(name string) PerVUIterationsConfig ***REMOVED***
	return PerVUIterationsConfig***REMOVED***
		BaseConfig:  NewBaseConfig(name, perVUIterationsType),
		VUs:         null.NewInt(1, false),
		Iterations:  null.NewInt(1, false),
		MaxDuration: types.NewNullDuration(10*time.Minute, false), // TODO: shorten?
	***REMOVED***
***REMOVED***

// Make sure we implement the lib.ExecutorConfig interface
var _ lib.ExecutorConfig = &PerVUIterationsConfig***REMOVED******REMOVED***

// GetVUs returns the scaled VUs for the executor.
func (pvic PerVUIterationsConfig) GetVUs(et *lib.ExecutionTuple) int64 ***REMOVED***
	return et.ScaleInt64(pvic.VUs.Int64)
***REMOVED***

// GetIterations returns the UNSCALED iteration count for the executor. It's
// important to note that scaling per-VU iteration executor affects only the
// number of VUs. If we also scaled the iterations, scaling would have quadratic
// effects instead of just linear.
func (pvic PerVUIterationsConfig) GetIterations() int64 ***REMOVED***
	return pvic.Iterations.Int64
***REMOVED***

// GetDescription returns a human-readable description of the executor options
func (pvic PerVUIterationsConfig) GetDescription(et *lib.ExecutionTuple) string ***REMOVED***
	return fmt.Sprintf("%d iterations for each of %d VUs%s",
		pvic.GetIterations(), pvic.GetVUs(et),
		pvic.getBaseInfo(fmt.Sprintf("maxDuration: %s", pvic.MaxDuration.Duration)))
***REMOVED***

// Validate makes sure all options are configured and valid
func (pvic PerVUIterationsConfig) Validate() []error ***REMOVED***
	errors := pvic.BaseConfig.Validate()
	if pvic.VUs.Int64 <= 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("the number of VUs should be more than 0"))
	***REMOVED***

	if pvic.Iterations.Int64 <= 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("the number of iterations should be more than 0"))
	***REMOVED***

	if pvic.MaxDuration.TimeDuration() < minDuration ***REMOVED***
		errors = append(errors, fmt.Errorf(
			"the maxDuration should be at least %s, but is %s", minDuration, pvic.MaxDuration,
		))
	***REMOVED***

	return errors
***REMOVED***

// GetExecutionRequirements returns the number of required VUs to run the
// executor for its whole duration (disregarding any startTime), including the
// maximum waiting time for any iterations to gracefully stop. This is used by
// the execution scheduler in its VU reservation calculations, so it knows how
// many VUs to pre-initialize.
func (pvic PerVUIterationsConfig) GetExecutionRequirements(et *lib.ExecutionTuple) []lib.ExecutionStep ***REMOVED***
	return []lib.ExecutionStep***REMOVED***
		***REMOVED***
			TimeOffset: 0,
			PlannedVUs: uint64(pvic.GetVUs(et)),
		***REMOVED***,
		***REMOVED***
			TimeOffset: pvic.MaxDuration.TimeDuration() + pvic.GracefulStop.TimeDuration(),
			PlannedVUs: 0,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// NewExecutor creates a new PerVUIterations executor
func (pvic PerVUIterationsConfig) NewExecutor(
	es *lib.ExecutionState, logger *logrus.Entry,
) (lib.Executor, error) ***REMOVED***
	return PerVUIterations***REMOVED***
		BaseExecutor: NewBaseExecutor(pvic, es, logger),
		config:       pvic,
	***REMOVED***, nil
***REMOVED***

// HasWork reports whether there is any work to be done for the given execution segment.
func (pvic PerVUIterationsConfig) HasWork(et *lib.ExecutionTuple) bool ***REMOVED***
	return pvic.GetVUs(et) > 0 && pvic.GetIterations() > 0
***REMOVED***

// PerVUIterations executes a specific number of iterations with each VU.
type PerVUIterations struct ***REMOVED***
	*BaseExecutor
	config PerVUIterationsConfig
***REMOVED***

// Make sure we implement the lib.Executor interface.
var _ lib.Executor = &PerVUIterations***REMOVED******REMOVED***

// Run executes a specific number of iterations with each configured VU.
// nolint:funlen
func (pvi PerVUIterations) Run(
	parentCtx context.Context, out chan<- stats.SampleContainer, builtinMetrics *metrics.BuiltinMetrics,
) (err error) ***REMOVED***
	numVUs := pvi.config.GetVUs(pvi.executionState.ExecutionTuple)
	iterations := pvi.config.GetIterations()
	duration := pvi.config.MaxDuration.TimeDuration()
	gracefulStop := pvi.config.GetGracefulStop()

	startTime, maxDurationCtx, regDurationCtx, cancel := getDurationContexts(parentCtx, duration, gracefulStop)
	defer cancel()

	// Make sure the log and the progress bar have accurate information
	pvi.logger.WithFields(logrus.Fields***REMOVED***
		"vus": numVUs, "iterations": iterations, "maxDuration": duration, "type": pvi.config.GetType(),
	***REMOVED***).Debug("Starting executor run...")

	totalIters := uint64(numVUs * iterations)
	doneIters := new(uint64)

	vusFmt := pb.GetFixedLengthIntFormat(numVUs)
	itersFmt := pb.GetFixedLengthIntFormat(int64(totalIters))
	progressFn := func() (float64, []string) ***REMOVED***
		spent := time.Since(startTime)
		progVUs := fmt.Sprintf(vusFmt+" VUs", numVUs)
		currentDoneIters := atomic.LoadUint64(doneIters)
		progIters := fmt.Sprintf(itersFmt+"/"+itersFmt+" iters, %d per VU",
			currentDoneIters, totalIters, iterations)
		right := []string***REMOVED***progVUs, duration.String(), progIters***REMOVED***
		if spent > duration ***REMOVED***
			return 1, right
		***REMOVED***

		spentDuration := pb.GetFixedLengthDuration(spent, duration)
		progDur := fmt.Sprintf("%s/%s", spentDuration, duration)
		right[1] = progDur

		return float64(currentDoneIters) / float64(totalIters), right
	***REMOVED***
	pvi.progress.Modify(pb.WithProgress(progressFn))
	go trackProgress(parentCtx, maxDurationCtx, regDurationCtx, pvi, progressFn)

	handleVUsWG := &sync.WaitGroup***REMOVED******REMOVED***
	defer handleVUsWG.Wait()
	// Actually schedule the VUs and iterations...
	activeVUs := &sync.WaitGroup***REMOVED******REMOVED***
	defer activeVUs.Wait()

	regDurationDone := regDurationCtx.Done()
	runIteration := getIterationRunner(pvi.executionState, pvi.logger)

	maxDurationCtx = lib.WithScenarioState(maxDurationCtx, &lib.ScenarioState***REMOVED***
		Name:       pvi.config.Name,
		Executor:   pvi.config.Type,
		StartTime:  startTime,
		ProgressFn: progressFn,
	***REMOVED***)

	returnVU := func(u lib.InitializedVU) ***REMOVED***
		pvi.executionState.ReturnVU(u, true)
		activeVUs.Done()
	***REMOVED***

	droppedIterationMetric := builtinMetrics.DroppedIterations
	handleVU := func(initVU lib.InitializedVU) ***REMOVED***
		defer handleVUsWG.Done()
		ctx, cancel := context.WithCancel(maxDurationCtx)
		defer cancel()

		vuID := initVU.GetID()
		activeVU := initVU.Activate(
			getVUActivationParams(ctx, pvi.config.BaseConfig, returnVU,
				pvi.nextIterationCounters))

		for i := int64(0); i < iterations; i++ ***REMOVED***
			select ***REMOVED***
			case <-regDurationDone:
				stats.PushIfNotDone(parentCtx, out, stats.Sample***REMOVED***
					Value: float64(iterations - i), Metric: droppedIterationMetric,
					Tags: pvi.getMetricTags(&vuID), Time: time.Now(),
				***REMOVED***)
				return // don't make more iterations
			default:
				// continue looping
			***REMOVED***
			runIteration(maxDurationCtx, activeVU)
			atomic.AddUint64(doneIters, 1)
		***REMOVED***
	***REMOVED***

	for i := int64(0); i < numVUs; i++ ***REMOVED***
		initializedVU, err := pvi.executionState.GetPlannedVU(pvi.logger, true)
		if err != nil ***REMOVED***
			cancel()
			return err
		***REMOVED***
		activeVUs.Add(1)
		handleVUsWG.Add(1)
		go handleVU(initializedVU)
	***REMOVED***

	return nil
***REMOVED***
