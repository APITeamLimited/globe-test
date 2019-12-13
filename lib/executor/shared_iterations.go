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

const sharedIterationsType = "shared-iterations"

func init() ***REMOVED***
	lib.RegisterExecutorConfigType(
		sharedIterationsType,
		func(name string, rawJSON []byte) (lib.ExecutorConfig, error) ***REMOVED***
			config := NewSharedIterationsConfig(name)
			err := lib.StrictJSONUnmarshal(rawJSON, &config)
			return config, err
		***REMOVED***,
	)
***REMOVED***

// SharedIterationsConfig stores the number of VUs iterations, as well as maxDuration settings
type SharedIterationsConfig struct ***REMOVED***
	BaseConfig
	VUs         null.Int           `json:"vus"`
	Iterations  null.Int           `json:"iterations"`
	MaxDuration types.NullDuration `json:"maxDuration"`
***REMOVED***

// NewSharedIterationsConfig returns a SharedIterationsConfig with default values
func NewSharedIterationsConfig(name string) SharedIterationsConfig ***REMOVED***
	return SharedIterationsConfig***REMOVED***
		BaseConfig:  NewBaseConfig(name, sharedIterationsType),
		VUs:         null.NewInt(1, false),
		Iterations:  null.NewInt(1, false),
		MaxDuration: types.NewNullDuration(10*time.Minute, false), //TODO: shorten?
	***REMOVED***
***REMOVED***

// Make sure we implement the lib.ExecutorConfig interface
var _ lib.ExecutorConfig = &SharedIterationsConfig***REMOVED******REMOVED***

// GetVUs returns the scaled VUs for the executor.
func (sic SharedIterationsConfig) GetVUs(es *lib.ExecutionSegment) int64 ***REMOVED***
	return es.Scale(sic.VUs.Int64)
***REMOVED***

// GetIterations returns the scaled iteration count for the executor.
func (sic SharedIterationsConfig) GetIterations(es *lib.ExecutionSegment) int64 ***REMOVED***
	return es.Scale(sic.Iterations.Int64)
***REMOVED***

// GetDescription returns a human-readable description of the executor options
func (sic SharedIterationsConfig) GetDescription(es *lib.ExecutionSegment) string ***REMOVED***
	return fmt.Sprintf("%d iterations shared among %d VUs%s",
		sic.GetIterations(es), sic.GetVUs(es),
		sic.getBaseInfo(fmt.Sprintf("maxDuration: %s", sic.MaxDuration.Duration)))
***REMOVED***

// Validate makes sure all options are configured and valid
func (sic SharedIterationsConfig) Validate() []error ***REMOVED***
	errors := sic.BaseConfig.Validate()
	if sic.VUs.Int64 <= 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("the number of VUs should be more than 0"))
	***REMOVED***

	if sic.Iterations.Int64 < sic.VUs.Int64 ***REMOVED***
		errors = append(errors, fmt.Errorf(
			"the number of iterations (%d) shouldn't be less than the number of VUs (%d)",
			sic.Iterations.Int64, sic.VUs.Int64,
		))
	***REMOVED***

	if time.Duration(sic.MaxDuration.Duration) < minDuration ***REMOVED***
		errors = append(errors, fmt.Errorf(
			"the maxDuration should be at least %s, but is %s", minDuration, sic.MaxDuration,
		))
	***REMOVED***

	return errors
***REMOVED***

// GetExecutionRequirements returns the number of required VUs to run the
// executor for its whole duration (disregarding any startTime), including the
// maximum waiting time for any iterations to gracefully stop. This is used by
// the execution scheduler in its VU reservation calculations, so it knows how
// many VUs to pre-initialize.
func (sic SharedIterationsConfig) GetExecutionRequirements(es *lib.ExecutionSegment) []lib.ExecutionStep ***REMOVED***
	return []lib.ExecutionStep***REMOVED***
		***REMOVED***
			TimeOffset: 0,
			PlannedVUs: uint64(sic.GetVUs(es)),
		***REMOVED***,
		***REMOVED***
			TimeOffset: time.Duration(sic.MaxDuration.Duration + sic.GracefulStop.Duration),
			PlannedVUs: 0,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// NewExecutor creates a new SharedIterations executor
func (sic SharedIterationsConfig) NewExecutor(
	es *lib.ExecutionState, logger *logrus.Entry,
) (lib.Executor, error) ***REMOVED***
	return SharedIterations***REMOVED***
		BaseExecutor: NewBaseExecutor(sic, es, logger),
		config:       sic,
	***REMOVED***, nil
***REMOVED***

// SharedIterations executes a specific total number of iterations, which are
// all shared by the configured VUs.
type SharedIterations struct ***REMOVED***
	*BaseExecutor
	config SharedIterationsConfig
***REMOVED***

// Make sure we implement the lib.Executor interface.
var _ lib.Executor = &SharedIterations***REMOVED******REMOVED***

// Run executes a specific total number of iterations, which are all shared by
// the configured VUs.
func (si SharedIterations) Run(ctx context.Context, out chan<- stats.SampleContainer) (err error) ***REMOVED***
	segment := si.executionState.Options.ExecutionSegment
	numVUs := si.config.GetVUs(segment)
	iterations := si.config.GetIterations(segment)
	duration := time.Duration(si.config.MaxDuration.Duration)
	gracefulStop := si.config.GetGracefulStop()

	_, maxDurationCtx, regDurationCtx, cancel := getDurationContexts(ctx, duration, gracefulStop)
	defer cancel()

	// Make sure the log and the progress bar have accurate information
	si.logger.WithFields(logrus.Fields***REMOVED***
		"vus": numVUs, "iterations": iterations, "maxDuration": duration, "type": si.config.GetType(),
	***REMOVED***).Debug("Starting executor run...")

	totalIters := uint64(iterations)
	doneIters := new(uint64)
	fmtStr := pb.GetFixedLengthIntFormat(int64(totalIters)) + "/%d shared iters among %d VUs"
	progresFn := func() (float64, string) ***REMOVED***
		currentDoneIters := atomic.LoadUint64(doneIters)
		return float64(currentDoneIters) / float64(totalIters), fmt.Sprintf(
			fmtStr, currentDoneIters, totalIters, numVUs,
		)
	***REMOVED***
	si.progress.Modify(pb.WithProgress(progresFn))
	go trackProgress(ctx, maxDurationCtx, regDurationCtx, si, progresFn)

	// Actually schedule the VUs and iterations...
	activeVUs := &sync.WaitGroup***REMOVED******REMOVED***
	defer activeVUs.Wait()

	regDurationDone := regDurationCtx.Done()
	runIteration := getIterationRunner(si.executionState, si.logger, out)

	attemptedIters := new(uint64)
	handleVU := func(vu lib.VU) ***REMOVED***
		defer activeVUs.Done()
		defer si.executionState.ReturnVU(vu, true)

		for ***REMOVED***
			select ***REMOVED***
			case <-regDurationDone:
				return // don't make more iterations
			default:
				// continue looping
			***REMOVED***

			attemptedIterNumber := atomic.AddUint64(attemptedIters, 1)
			if attemptedIterNumber > totalIters ***REMOVED***
				return
			***REMOVED***

			runIteration(maxDurationCtx, vu)
			atomic.AddUint64(doneIters, 1)
		***REMOVED***
	***REMOVED***

	for i := int64(0); i < numVUs; i++ ***REMOVED***
		vu, err := si.executionState.GetPlannedVU(si.logger, true)
		if err != nil ***REMOVED***
			cancel()
			return err
		***REMOVED***
		activeVUs.Add(1)
		go handleVU(vu)
	***REMOVED***

	return nil
***REMOVED***
