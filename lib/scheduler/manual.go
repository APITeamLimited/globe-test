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
	"errors"
	"fmt"
	"time"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/types"
	"github.com/sirupsen/logrus"
	null "gopkg.in/guregu/null.v3"
)

const manualExecution = "manual-execution"

// ManualExecutionConfig stores VUs and duration
type ManualExecutionConfig struct ***REMOVED***
	StartVUs null.Int
	MaxVUs   null.Int
	Duration types.NullDuration
***REMOVED***

// NewManualExecutionConfig returns a ManualExecutionConfig with default values
func NewManualExecutionConfig(startVUs, maxVUs null.Int, duration types.NullDuration) ManualExecutionConfig ***REMOVED***
	if !maxVUs.Valid ***REMOVED***
		maxVUs = startVUs
	***REMOVED***
	return ManualExecutionConfig***REMOVED***startVUs, maxVUs, duration***REMOVED***
***REMOVED***

// Make sure we implement the lib.SchedulerConfig interface
var _ lib.SchedulerConfig = &ManualExecutionConfig***REMOVED******REMOVED***

// GetDescription returns a human-readable description of the scheduler options
func (mec ManualExecutionConfig) GetDescription(_ *lib.ExecutionSegment) string ***REMOVED***
	duration := ""
	if mec.Duration.Duration != 0 ***REMOVED***
		duration = fmt.Sprintf(" and duration %s", mec.Duration)
	***REMOVED***
	return fmt.Sprintf(
		"Manual execution with %d starting and %d initialized VUs%s",
		mec.StartVUs.Int64, mec.MaxVUs.Int64, duration,
	)
***REMOVED***

// Validate makes sure all options are configured and valid
func (mec ManualExecutionConfig) Validate() []error ***REMOVED***
	var errors []error
	if mec.StartVUs.Int64 <= 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("the number of VUs should be more than 0"))
	***REMOVED***

	if mec.MaxVUs.Int64 < mec.StartVUs.Int64 ***REMOVED***
		errors = append(errors, fmt.Errorf("the number of MaxVUs should more than or equal to the starting number of VUs"))
	***REMOVED***

	if !mec.Duration.Valid ***REMOVED***
		errors = append(errors, fmt.Errorf("the duration should be specified, for infinite duration use 0"))
	***REMOVED*** else if time.Duration(mec.Duration.Duration) < 0 ***REMOVED***
		errors = append(errors, fmt.Errorf(
			"the duration shouldn't be negative, for infinite duration use 0",
		))
	***REMOVED***

	return errors
***REMOVED***

// GetExecutionRequirements just reserves the number of starting VUs for the whole
// duration of the scheduler, so these VUs can be initialized in the beginning of the
// test.
//
// Importantly, if 0 (i.e. infinite) duration is configured, this scheduler doesn't
// emit the last step to relinquish these VUs.
//
// Also, the manual execution scheduler doesn't set MaxUnplannedVUs in the returned steps,
// since their initialization and usage is directly controlled by the user and is effectively
// bounded only by the resources of the machine k6 is running on.
//
// This is not a problem, because the MaxUnplannedVUs are mostly meant to be used for
// calculating the maximum possble number of initialized VUs at any point during a test
// run. That's used for sizing purposes and for user qouta checking in the cloud execution,
// where the manual scheduler isn't supported.
func (mec ManualExecutionConfig) GetExecutionRequirements(es *lib.ExecutionSegment) []lib.ExecutionStep ***REMOVED***
	startVUs := lib.ExecutionStep***REMOVED***
		TimeOffset:      0,
		PlannedVUs:      uint64(es.Scale(mec.StartVUs.Int64)),
		MaxUnplannedVUs: 0, // intentional, see function comment
	***REMOVED***

	maxDuration := time.Duration(mec.Duration.Duration)
	if maxDuration == 0 ***REMOVED***
		// Infinite duration, don't emit 0 VUs at the end since there's no planned end
		return []lib.ExecutionStep***REMOVED***startVUs***REMOVED***
	***REMOVED***
	return []lib.ExecutionStep***REMOVED***startVUs, ***REMOVED***
		TimeOffset:      maxDuration,
		PlannedVUs:      0,
		MaxUnplannedVUs: 0, // intentional, see function comment
	***REMOVED******REMOVED***
***REMOVED***

// GetName always returns manual-execution, since this config can't be
// specified in the exported script options.
func (ManualExecutionConfig) GetName() string ***REMOVED***
	return manualExecution
***REMOVED***

// GetType always returns manual-execution, since that's this special
// config's type...
func (ManualExecutionConfig) GetType() string ***REMOVED***
	return manualExecution
***REMOVED***

// GetStartTime always returns 0, since the manual execution scheduler
// always starts in the beginning and is always the only scheduler.
func (ManualExecutionConfig) GetStartTime() time.Duration ***REMOVED***
	return 0
***REMOVED***

// GetGracefulStop always returns 0, since we still don't support graceful
// stops or ramp downs in the manual execution mode.
//TODO: implement?
func (ManualExecutionConfig) GetGracefulStop() time.Duration ***REMOVED***
	return 0
***REMOVED***

// GetEnv returns an empty map, since the manual executor doesn't support custom
// environment variables.
func (ManualExecutionConfig) GetEnv() map[string]string ***REMOVED***
	return nil
***REMOVED***

// GetExec always returns nil, for now there's no way to execute custom funcions in
// the manual execution mode.
func (ManualExecutionConfig) GetExec() null.String ***REMOVED***
	return null.NewString("", false)
***REMOVED***

// IsDistributable simply returns false because there's no way to reliably
// distribute the manual execution scheduler.
func (ManualExecutionConfig) IsDistributable() bool ***REMOVED***
	return false
***REMOVED***

// NewScheduler creates a new ManualExecution "scheduler"
func (mec ManualExecutionConfig) NewScheduler(
	es *lib.ExecutorState, logger *logrus.Entry) (lib.Scheduler, error) ***REMOVED***
	return nil, errors.New("not implemented 4") //TODO
***REMOVED***
