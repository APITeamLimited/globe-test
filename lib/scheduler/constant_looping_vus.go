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
	"time"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/stats"
	"github.com/loadimpact/k6/ui/pb"
	"github.com/sirupsen/logrus"
	null "gopkg.in/guregu/null.v3"
)

const constantLoopingVUsType = "constant-looping-vus"

func init() ***REMOVED***
	lib.RegisterSchedulerConfigType(
		constantLoopingVUsType,
		func(name string, rawJSON []byte) (lib.SchedulerConfig, error) ***REMOVED***
			config := NewConstantLoopingVUsConfig(name)
			err := lib.StrictJSONUnmarshal(rawJSON, &config)
			return config, err
		***REMOVED***,
	)
***REMOVED***

// The minimum duration we'll allow users to schedule. This doesn't affect the stages
// configuration, where 0-duration virtual stages are allowed for instantaneous VU jumps
const minDuration = 1 * time.Second

// ConstantLoopingVUsConfig stores VUs and duration
type ConstantLoopingVUsConfig struct ***REMOVED***
	BaseConfig
	VUs      null.Int           `json:"vus"`
	Duration types.NullDuration `json:"duration"`
***REMOVED***

// NewConstantLoopingVUsConfig returns a ConstantLoopingVUsConfig with default values
func NewConstantLoopingVUsConfig(name string) ConstantLoopingVUsConfig ***REMOVED***
	return ConstantLoopingVUsConfig***REMOVED***
		BaseConfig: NewBaseConfig(name, constantLoopingVUsType),
		VUs:        null.NewInt(1, false),
	***REMOVED***
***REMOVED***

// Make sure we implement the lib.SchedulerConfig interface
var _ lib.SchedulerConfig = &ConstantLoopingVUsConfig***REMOVED******REMOVED***

// GetVUs returns the scaled VUs for the scheduler.
func (clvc ConstantLoopingVUsConfig) GetVUs(es *lib.ExecutionSegment) int64 ***REMOVED***
	return es.Scale(clvc.VUs.Int64)
***REMOVED***

// GetDescription returns a human-readable description of the scheduler options
func (clvc ConstantLoopingVUsConfig) GetDescription(es *lib.ExecutionSegment) string ***REMOVED***
	return fmt.Sprintf("%d looping VUs for %s%s",
		clvc.GetVUs(es), clvc.Duration.Duration, clvc.getBaseInfo())
***REMOVED***

// Validate makes sure all options are configured and valid
func (clvc ConstantLoopingVUsConfig) Validate() []error ***REMOVED***
	errors := clvc.BaseConfig.Validate()
	if clvc.VUs.Int64 <= 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("the number of VUs should be more than 0"))
	***REMOVED***

	if !clvc.Duration.Valid ***REMOVED***
		errors = append(errors, fmt.Errorf("the duration is unspecified"))
	***REMOVED*** else if time.Duration(clvc.Duration.Duration) < minDuration ***REMOVED***
		errors = append(errors, fmt.Errorf(
			"the duration should be at least %s, but is %s", minDuration, clvc.Duration,
		))
	***REMOVED***

	return errors
***REMOVED***

// GetExecutionRequirements just reserves the number of specified VUs for the
// whole duration of the scheduler, including the maximum waiting time for
// iterations to gracefully stop.
func (clvc ConstantLoopingVUsConfig) GetExecutionRequirements(es *lib.ExecutionSegment) []lib.ExecutionStep ***REMOVED***
	return []lib.ExecutionStep***REMOVED***
		***REMOVED***
			TimeOffset: 0,
			PlannedVUs: uint64(clvc.GetVUs(es)),
		***REMOVED***,
		***REMOVED***
			TimeOffset: time.Duration(clvc.Duration.Duration + clvc.GracefulStop.Duration),
			PlannedVUs: 0,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// NewScheduler creates a new ConstantLoopingVUs scheduler
func (clvc ConstantLoopingVUsConfig) NewScheduler(es *lib.ExecutorState, logger *logrus.Entry) (lib.Scheduler, error) ***REMOVED***
	return ConstantLoopingVUs***REMOVED***
		BaseScheduler: NewBaseScheduler(clvc, es, logger),
		config:        clvc,
	***REMOVED***, nil
***REMOVED***

// ConstantLoopingVUs maintains a constant number of VUs running for the
// specified duration.
type ConstantLoopingVUs struct ***REMOVED***
	*BaseScheduler
	config ConstantLoopingVUsConfig
***REMOVED***

// Make sure we implement the lib.Scheduler interface.
var _ lib.Scheduler = &ConstantLoopingVUs***REMOVED******REMOVED***

// Run constantly loops through as many iterations as possible on a fixed number
// of VUs for the specified duration.
func (clv ConstantLoopingVUs) Run(ctx context.Context, out chan<- stats.SampleContainer) (err error) ***REMOVED***
	segment := clv.executorState.Options.ExecutionSegment
	numVUs := clv.config.GetVUs(segment)
	duration := time.Duration(clv.config.Duration.Duration)
	gracefulStop := clv.config.GetGracefulStop()

	startTime, maxDurationCtx, regDurationCtx, cancel := getDurationContexts(ctx, duration, gracefulStop)
	defer cancel()

	// Make sure the log and the progress bar have accurate information
	clv.logger.WithFields(
		logrus.Fields***REMOVED***"vus": numVUs, "duration": duration, "type": clv.config.GetType()***REMOVED***,
	).Debug("Starting scheduler run...")

	progresFn := func() (float64, string) ***REMOVED***
		spent := time.Since(startTime)
		if spent > duration ***REMOVED***
			return 1, fmt.Sprintf("constant looping %d VUs for %s", numVUs, duration)
		***REMOVED***
		return float64(spent) / float64(duration), fmt.Sprintf(
			"constant looping %d VUs, %s/%s", numVUs, pb.GetFixedLengthDuration(spent, duration), duration,
		)
	***REMOVED***
	clv.progress.Modify(pb.WithProgress(progresFn))
	go trackProgress(ctx, maxDurationCtx, regDurationCtx, clv, progresFn)

	// Actually schedule the VUs and iterations...
	activeVUs := &sync.WaitGroup***REMOVED******REMOVED***
	defer activeVUs.Wait()

	regDurationDone := regDurationCtx.Done()
	runIteration := getIterationRunner(clv.executorState, clv.logger, out)

	handleVU := func(vu lib.VU) ***REMOVED***
		defer clv.executorState.ReturnVU(vu)
		defer activeVUs.Done()

		for ***REMOVED***
			select ***REMOVED***
			case <-regDurationDone:
				return // don't make more iterations
			default:
				// continue looping
			***REMOVED***
			runIteration(maxDurationCtx, vu)
		***REMOVED***
	***REMOVED***

	for i := int64(0); i < numVUs; i++ ***REMOVED***
		vu, err := clv.executorState.GetPlannedVU(clv.logger)
		if err != nil ***REMOVED***
			cancel()
			return err
		***REMOVED***
		activeVUs.Add(1)
		go handleVU(vu)
	***REMOVED***

	return nil
***REMOVED***
