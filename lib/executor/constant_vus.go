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

	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/types"
	"go.k6.io/k6/stats"
	"go.k6.io/k6/ui/pb"
)

const constantVUsType = "constant-vus"

func init() ***REMOVED***
	lib.RegisterExecutorConfigType(
		constantVUsType,
		func(name string, rawJSON []byte) (lib.ExecutorConfig, error) ***REMOVED***
			config := NewConstantVUsConfig(name)
			err := lib.StrictJSONUnmarshal(rawJSON, &config)
			return config, err
		***REMOVED***,
	)
***REMOVED***

// The minimum duration we'll allow users to schedule. This doesn't affect the stages
// configuration, where 0-duration virtual stages are allowed for instantaneous VU jumps
const minDuration = 1 * time.Second

// ConstantVUsConfig stores VUs and duration
type ConstantVUsConfig struct ***REMOVED***
	BaseConfig
	VUs      null.Int           `json:"vus"`
	Duration types.NullDuration `json:"duration"`
***REMOVED***

// NewConstantVUsConfig returns a ConstantVUsConfig with default values
func NewConstantVUsConfig(name string) ConstantVUsConfig ***REMOVED***
	return ConstantVUsConfig***REMOVED***
		BaseConfig: NewBaseConfig(name, constantVUsType),
		VUs:        null.NewInt(1, false),
	***REMOVED***
***REMOVED***

// Make sure we implement the lib.ExecutorConfig interface
var _ lib.ExecutorConfig = &ConstantVUsConfig***REMOVED******REMOVED***

// GetVUs returns the scaled VUs for the executor.
func (clvc ConstantVUsConfig) GetVUs(et *lib.ExecutionTuple) int64 ***REMOVED***
	return et.Segment.Scale(clvc.VUs.Int64)
***REMOVED***

// GetDescription returns a human-readable description of the executor options
func (clvc ConstantVUsConfig) GetDescription(et *lib.ExecutionTuple) string ***REMOVED***
	return fmt.Sprintf("%d looping VUs for %s%s",
		clvc.GetVUs(et), clvc.Duration.Duration, clvc.getBaseInfo())
***REMOVED***

// Validate makes sure all options are configured and valid
func (clvc ConstantVUsConfig) Validate() []error ***REMOVED***
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
func (clvc ConstantVUsConfig) GetExecutionRequirements(et *lib.ExecutionTuple) []lib.ExecutionStep ***REMOVED***
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
func (clvc ConstantVUsConfig) HasWork(et *lib.ExecutionTuple) bool ***REMOVED***
	return clvc.GetVUs(et) > 0
***REMOVED***

// NewExecutor creates a new ConstantVUs executor
func (clvc ConstantVUsConfig) NewExecutor(es *lib.ExecutionState, logger *logrus.Entry) (lib.Executor, error) ***REMOVED***
	return ConstantVUs***REMOVED***
		BaseExecutor: NewBaseExecutor(clvc, es, logger),
		config:       clvc,
	***REMOVED***, nil
***REMOVED***

// ConstantVUs maintains a constant number of VUs running for the
// specified duration.
type ConstantVUs struct ***REMOVED***
	*BaseExecutor
	config ConstantVUsConfig
***REMOVED***

// Make sure we implement the lib.Executor interface.
var _ lib.Executor = &ConstantVUs***REMOVED******REMOVED***

// Run constantly loops through as many iterations as possible on a fixed number
// of VUs for the specified duration.
func (clv ConstantVUs) Run(parentCtx context.Context, out chan<- stats.SampleContainer) (err error) ***REMOVED***
	numVUs := clv.config.GetVUs(clv.executionState.ExecutionTuple)
	duration := time.Duration(clv.config.Duration.Duration)
	gracefulStop := clv.config.GetGracefulStop()

	startTime, maxDurationCtx, regDurationCtx, cancel := getDurationContexts(parentCtx, duration, gracefulStop)
	defer cancel()

	// Make sure the log and the progress bar have accurate information
	clv.logger.WithFields(
		logrus.Fields***REMOVED***"vus": numVUs, "duration": duration, "type": clv.config.GetType()***REMOVED***,
	).Debug("Starting executor run...")

	progressFn := func() (float64, []string) ***REMOVED***
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
	clv.progress.Modify(pb.WithProgress(progressFn))
	go trackProgress(parentCtx, maxDurationCtx, regDurationCtx, clv, progressFn)

	// Actually schedule the VUs and iterations...
	activeVUs := &sync.WaitGroup***REMOVED******REMOVED***
	defer activeVUs.Wait()

	regDurationDone := regDurationCtx.Done()
	runIteration := getIterationRunner(clv.executionState, clv.logger)

	maxDurationCtx = lib.WithScenarioState(maxDurationCtx, &lib.ScenarioState***REMOVED***
		Name:       clv.config.Name,
		Executor:   clv.config.Type,
		StartTime:  startTime,
		ProgressFn: progressFn,
	***REMOVED***)

	returnVU := func(u lib.InitializedVU) ***REMOVED***
		clv.executionState.ReturnVU(u, true)
		activeVUs.Done()
	***REMOVED***

	handleVU := func(initVU lib.InitializedVU) ***REMOVED***
		ctx, cancel := context.WithCancel(maxDurationCtx)
		defer cancel()

		activeVU := initVU.Activate(
			getVUActivationParams(ctx, clv.config.BaseConfig, returnVU, clv.GetNextLocalVUID))

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
