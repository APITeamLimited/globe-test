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

package scheduler

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/stats"
	"github.com/loadimpact/k6/ui/pb"
	"github.com/sirupsen/logrus"
	null "gopkg.in/guregu/null.v3"
)

const perVUIterationsType = "per-vu-iterations"

func init() ***REMOVED***
	lib.RegisterSchedulerConfigType(perVUIterationsType, func(name string, rawJSON []byte) (lib.SchedulerConfig, error) ***REMOVED***
		config := NewPerVUIterationsConfig(name)
		err := lib.StrictJSONUnmarshal(rawJSON, &config)
		return config, err
	***REMOVED***)
***REMOVED***

// PerVUIteationsConfig stores the number of VUs iterations, as well as maxDuration settings
type PerVUIteationsConfig struct ***REMOVED***
	BaseConfig
	VUs         null.Int           `json:"vus"`
	Iterations  null.Int           `json:"iterations"`
	MaxDuration types.NullDuration `json:"maxDuration"`
***REMOVED***

// NewPerVUIterationsConfig returns a PerVUIteationsConfig with default values
func NewPerVUIterationsConfig(name string) PerVUIteationsConfig ***REMOVED***
	return PerVUIteationsConfig***REMOVED***
		BaseConfig:  NewBaseConfig(name, perVUIterationsType),
		VUs:         null.NewInt(1, false),
		Iterations:  null.NewInt(1, false),
		MaxDuration: types.NewNullDuration(10*time.Minute, false), //TODO: shorten?
	***REMOVED***
***REMOVED***

// Make sure we implement the lib.SchedulerConfig interface
var _ lib.SchedulerConfig = &PerVUIteationsConfig***REMOVED******REMOVED***

// GetVUs returns the scaled VUs for the scheduler.
func (pvic PerVUIteationsConfig) GetVUs(es *lib.ExecutionSegment) int64 ***REMOVED***
	return es.Scale(pvic.VUs.Int64)
***REMOVED***

// GetIterations returns the UNSCALED iteration count for the scheduler. It's
// important to note that scaling per-VU iteration scheduler affects only the
// number of VUs. If we also scaled the iterations, scaling would have quadratic
// effects instead of just linear.
func (pvic PerVUIteationsConfig) GetIterations() int64 ***REMOVED***
	return pvic.Iterations.Int64
***REMOVED***

// GetDescription returns a human-readable description of the scheduler options
func (pvic PerVUIteationsConfig) GetDescription(es *lib.ExecutionSegment) string ***REMOVED***
	return fmt.Sprintf("%d iterations for each of %d VUs%s",
		pvic.GetIterations(), pvic.GetVUs(es),
		pvic.getBaseInfo(fmt.Sprintf("maxDuration: %s", pvic.MaxDuration.Duration)))
***REMOVED***

// Validate makes sure all options are configured and valid
func (pvic PerVUIteationsConfig) Validate() []error ***REMOVED***
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

// GetExecutionRequirements just reserves the number of specified VUs for the
// whole duration of the scheduler, including the maximum waiting time for
// iterations to gracefully stop.
func (pvic PerVUIteationsConfig) GetExecutionRequirements(es *lib.ExecutionSegment) []lib.ExecutionStep ***REMOVED***
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

// NewScheduler creates a new PerVUIteations scheduler
func (pvic PerVUIteationsConfig) NewScheduler(
	es *lib.ExecutorState, logger *logrus.Entry) (lib.Scheduler, error) ***REMOVED***

	return PerVUIteations***REMOVED***
		BaseScheduler: NewBaseScheduler(pvic, es, logger),
		config:        pvic,
	***REMOVED***, nil
***REMOVED***

// PerVUIteations executes a specific number of iterations with each VU.
type PerVUIteations struct ***REMOVED***
	*BaseScheduler
	config PerVUIteationsConfig
***REMOVED***

// Make sure we implement the lib.Scheduler interface.
var _ lib.Scheduler = &PerVUIteations***REMOVED******REMOVED***

// Run executes a specific number of iterations with each confugured VU.
func (pvi PerVUIteations) Run(ctx context.Context, out chan<- stats.SampleContainer) (err error) ***REMOVED***
	segment := pvi.executorState.Options.ExecutionSegment
	numVUs := pvi.config.GetVUs(segment)
	iterations := pvi.config.GetIterations()
	duration := time.Duration(pvi.config.MaxDuration.Duration)
	gracefulStop := pvi.config.GetGracefulStop()

	_, maxDurationCtx, regDurationCtx, cancel := getDurationContexts(ctx, duration, gracefulStop)
	defer cancel()

	// Make sure the log and the progress bar have accurate information
	pvi.logger.WithFields(logrus.Fields***REMOVED***
		"vus": numVUs, "iterations": iterations, "maxDuration": duration, "type": pvi.config.GetType(),
	***REMOVED***).Debug("Starting scheduler run...")

	totalIters := uint64(numVUs * iterations)
	doneIters := new(uint64)
	fmtStr := pb.GetFixedLengthIntFormat(int64(totalIters)) + "/%d iters, %d from each of %d VUs"
	progresFn := func() (float64, string) ***REMOVED***
		currentDoneIters := atomic.LoadUint64(doneIters)
		return float64(currentDoneIters) / float64(totalIters), fmt.Sprintf(
			fmtStr, currentDoneIters, totalIters, iterations, numVUs,
		)
	***REMOVED***
	pvi.progress.Modify(pb.WithProgress(progresFn))
	go trackProgress(ctx, maxDurationCtx, regDurationCtx, pvi, progresFn)

	// Actually schedule the VUs and iterations...
	wg := sync.WaitGroup***REMOVED******REMOVED***
	regDurationDone := regDurationCtx.Done()
	runIteration := getIterationRunner(pvi.executorState, pvi.logger, out)

	handleVU := func(vu lib.VU) ***REMOVED***
		defer pvi.executorState.ReturnVU(vu)
		defer wg.Done()

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
		wg.Add(1)
		vu, err := pvi.executorState.GetPlannedVU(ctx, pvi.logger)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		go handleVU(vu)
	***REMOVED***

	wg.Wait()
	return nil
***REMOVED***
