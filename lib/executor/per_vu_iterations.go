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
		MaxDuration: types.NewNullDuration(10*time.Minute, false), //TODO: shorten?
	***REMOVED***
***REMOVED***

// Make sure we implement the lib.ExecutorConfig interface
var _ lib.ExecutorConfig = &PerVUIterationsConfig***REMOVED******REMOVED***

// GetVUs returns the scaled VUs for the executor.
func (pvic PerVUIterationsConfig) GetVUs(es *lib.ExecutionSegment) int64 ***REMOVED***
	return es.Scale(pvic.VUs.Int64)
***REMOVED***

// GetIterations returns the UNSCALED iteration count for the executor. It's
// important to note that scaling per-VU iteration executor affects only the
// number of VUs. If we also scaled the iterations, scaling would have quadratic
// effects instead of just linear.
func (pvic PerVUIterationsConfig) GetIterations() int64 ***REMOVED***
	return pvic.Iterations.Int64
***REMOVED***

// GetDescription returns a human-readable description of the executor options
func (pvic PerVUIterationsConfig) GetDescription(es *lib.ExecutionSegment) string ***REMOVED***
	return fmt.Sprintf("%d iterations for each of %d VUs%s",
		pvic.GetIterations(), pvic.GetVUs(es),
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

	if time.Duration(pvic.MaxDuration.Duration) < minDuration ***REMOVED***
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
func (pvic PerVUIterationsConfig) GetExecutionRequirements(es *lib.ExecutionSegment) []lib.ExecutionStep ***REMOVED***
	return []lib.ExecutionStep***REMOVED***
		***REMOVED***
			TimeOffset: 0,
			PlannedVUs: uint64(pvic.GetVUs(es)),
		***REMOVED***,
		***REMOVED***
			TimeOffset: time.Duration(pvic.MaxDuration.Duration + pvic.GracefulStop.Duration),
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
func (pvic PerVUIterationsConfig) HasWork(es *lib.ExecutionSegment) bool ***REMOVED***
	return pvic.GetVUs(es) > 0 && pvic.GetIterations() > 0
***REMOVED***

// PerVUIterations executes a specific number of iterations with each VU.
type PerVUIterations struct ***REMOVED***
	*BaseExecutor
	config PerVUIterationsConfig
***REMOVED***

// Make sure we implement the lib.Executor interface.
var _ lib.Executor = &PerVUIterations***REMOVED******REMOVED***

// Run executes a specific number of iterations with each configured VU.
func (pvi PerVUIterations) Run(ctx context.Context, out chan<- stats.SampleContainer) (err error) ***REMOVED***
	segment := pvi.executionState.Options.ExecutionSegment
	numVUs := pvi.config.GetVUs(segment)
	iterations := pvi.config.GetIterations()
	duration := time.Duration(pvi.config.MaxDuration.Duration)
	gracefulStop := pvi.config.GetGracefulStop()

	_, maxDurationCtx, regDurationCtx, cancel := getDurationContexts(ctx, duration, gracefulStop)
	defer cancel()

	// Make sure the log and the progress bar have accurate information
	pvi.logger.WithFields(logrus.Fields***REMOVED***
		"vus": numVUs, "iterations": iterations, "maxDuration": duration, "type": pvi.config.GetType(),
	***REMOVED***).Debug("Starting executor run...")

	totalIters := uint64(numVUs * iterations)
	doneIters := new(uint64)

	vusFmt := pb.GetFixedLengthIntFormat(numVUs)
	itersFmt := pb.GetFixedLengthIntFormat(int64(totalIters))
	progresFn := func() (float64, []string) ***REMOVED***
		currentDoneIters := atomic.LoadUint64(doneIters)
		return float64(currentDoneIters) / float64(totalIters), []string***REMOVED***
			fmt.Sprintf(vusFmt+" VUs", numVUs),
			fmt.Sprintf(itersFmt+"/"+itersFmt+" iters, %d per VU",
				currentDoneIters, totalIters, iterations),
		***REMOVED***
	***REMOVED***
	pvi.progress.Modify(pb.WithProgress(progresFn))
	go trackProgress(ctx, maxDurationCtx, regDurationCtx, pvi, progresFn)

	// Actually schedule the VUs and iterations...
	activeVUs := &sync.WaitGroup***REMOVED******REMOVED***
	defer activeVUs.Wait()

	regDurationDone := regDurationCtx.Done()
	runIteration := getIterationRunner(pvi.executionState, pvi.logger, out)

	handleVU := func(vu lib.VU) ***REMOVED***
		defer activeVUs.Done()
		defer pvi.executionState.ReturnVU(vu, true)

		for i := int64(0); i < iterations; i++ ***REMOVED***
			select ***REMOVED***
			case <-regDurationDone:
				return // don't make more iterations
			default:
				// continue looping
			***REMOVED***
			runIteration(maxDurationCtx, vu)
			atomic.AddUint64(doneIters, 1)
		***REMOVED***
	***REMOVED***

	for i := int64(0); i < numVUs; i++ ***REMOVED***
		vu, err := pvi.executionState.GetPlannedVU(pvi.logger, true)
		if err != nil ***REMOVED***
			cancel()
			return err
		***REMOVED***
		activeVUs.Add(1)
		go handleVU(vu)
	***REMOVED***

	return nil
***REMOVED***
