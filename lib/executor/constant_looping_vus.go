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
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/stats"
	"github.com/loadimpact/k6/ui/pb"
)

const constantLoopingVUsType = "constant-looping-vus"

func init() ***REMOVED***
	lib.RegisterExecutorConfigType(
		constantLoopingVUsType,
		func(name string, rawJSON []byte) (lib.ExecutorConfig, error) ***REMOVED***
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

// Make sure we implement the lib.ExecutorConfig interface
var _ lib.ExecutorConfig = &ConstantLoopingVUsConfig***REMOVED******REMOVED***

// GetVUs returns the scaled VUs for the executor.
func (clvc ConstantLoopingVUsConfig) GetVUs(et *lib.ExecutionTuple) int64 ***REMOVED***
	return et.Segment.Scale(clvc.VUs.Int64)
***REMOVED***

// GetDescription returns a human-readable description of the executor options
func (clvc ConstantLoopingVUsConfig) GetDescription(et *lib.ExecutionTuple) string ***REMOVED***
	return fmt.Sprintf("%d looping VUs for %s%s",
		clvc.GetVUs(et), clvc.Duration.Duration, clvc.getBaseInfo())
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

// GetExecutionRequirements returns the number of required VUs to run the
// executor for its whole duration (disregarding any startTime), including the
// maximum waiting time for any iterations to gracefully stop. This is used by
// the execution scheduler in its VU reservation calculations, so it knows how
// many VUs to pre-initialize.
func (clvc ConstantLoopingVUsConfig) GetExecutionRequirements(et *lib.ExecutionTuple) []lib.ExecutionStep ***REMOVED***
	return []lib.ExecutionStep***REMOVED***
		***REMOVED***
			TimeOffset: 0,
			PlannedVUs: uint64(clvc.GetVUs(et)),
		***REMOVED***,
		***REMOVED***
			TimeOffset: time.Duration(clvc.Duration.Duration + clvc.GracefulStop.Duration),
			PlannedVUs: 0,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// HasWork reports whether there is any work to be done for the given execution segment.
func (clvc ConstantLoopingVUsConfig) HasWork(et *lib.ExecutionTuple) bool ***REMOVED***
	return clvc.GetVUs(et) > 0
***REMOVED***

// NewExecutor creates a new ConstantLoopingVUs executor
func (clvc ConstantLoopingVUsConfig) NewExecutor(es *lib.ExecutionState, logger *logrus.Entry) (lib.Executor, error) ***REMOVED***
	return ConstantLoopingVUs***REMOVED***
		BaseExecutor: NewBaseExecutor(clvc, es, logger),
		config:       clvc,
	***REMOVED***, nil
***REMOVED***

// ConstantLoopingVUs maintains a constant number of VUs running for the
// specified duration.
type ConstantLoopingVUs struct ***REMOVED***
	*BaseExecutor
	config ConstantLoopingVUsConfig
***REMOVED***

// Make sure we implement the lib.Executor interface.
var _ lib.Executor = &ConstantLoopingVUs***REMOVED******REMOVED***

// Run constantly loops through as many iterations as possible on a fixed number
// of VUs for the specified duration.
func (clv ConstantLoopingVUs) Run(ctx context.Context, out chan<- stats.SampleContainer) (err error) ***REMOVED***
	numVUs := clv.config.GetVUs(clv.executionState.ExecutionTuple)
	duration := time.Duration(clv.config.Duration.Duration)
	gracefulStop := clv.config.GetGracefulStop()

	startTime, maxDurationCtx, regDurationCtx, cancel := getDurationContexts(ctx, duration, gracefulStop)
	defer cancel()

	// Make sure the log and the progress bar have accurate information
	clv.logger.WithFields(
		logrus.Fields***REMOVED***"vus": numVUs, "duration": duration, "type": clv.config.GetType()***REMOVED***,
	).Debug("Starting executor run...")

	progresFn := func() (float64, []string) ***REMOVED***
		spent := time.Since(startTime)
		right := []string***REMOVED***fmt.Sprintf("%d VUs", numVUs)***REMOVED***
		if spent > duration ***REMOVED***
			right = append(right, duration.String())
			return 1, right
		***REMOVED***
		right = append(right, fmt.Sprintf("%s/%s",
			pb.GetFixedLengthDuration(spent, duration), duration))
		return float64(spent) / float64(duration), right
	***REMOVED***
	clv.progress.Modify(pb.WithProgress(progresFn))
	go trackProgress(ctx, maxDurationCtx, regDurationCtx, clv, progresFn)

	// Actually schedule the VUs and iterations...
	activeVUs := &sync.WaitGroup***REMOVED******REMOVED***
	defer activeVUs.Wait()

	regDurationDone := regDurationCtx.Done()
	runIteration := getIterationRunner(clv.executionState, clv.logger)

	activationParams := getVUActivationParams(maxDurationCtx, clv.config.BaseConfig,
		func(u lib.InitializedVU) ***REMOVED***
			clv.executionState.ReturnVU(u, true)
			activeVUs.Done()
		***REMOVED***)
	handleVU := func(initVU lib.InitializedVU) ***REMOVED***
		ctx, cancel := context.WithCancel(maxDurationCtx)
		defer cancel()

		newParams := *activationParams
		newParams.RunContext = ctx

		activeVU := initVU.Activate(&newParams)

		for ***REMOVED***
			select ***REMOVED***
			case <-regDurationDone:
				return // don't make more iterations
			default:
				// continue looping
			***REMOVED***
			runIteration(maxDurationCtx, activeVU)
		***REMOVED***
	***REMOVED***

	for i := int64(0); i < numVUs; i++ ***REMOVED***
		initVU, err := clv.executionState.GetPlannedVU(clv.logger, true)
		if err != nil ***REMOVED***
			cancel()
			return err
		***REMOVED***
		activeVUs.Add(1)
		go handleVU(initVU)
	***REMOVED***

	return nil
***REMOVED***
