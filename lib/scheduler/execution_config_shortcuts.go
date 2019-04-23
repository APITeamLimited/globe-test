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
	"github.com/loadimpact/k6/lib"
	"github.com/sirupsen/logrus"
	null "gopkg.in/guregu/null.v3"
)

// ExecutionConflictError is a custom error type used for all of the errors in
// the BuildExecutionConfig() function.
type ExecutionConflictError string

func (e ExecutionConflictError) Error() string ***REMOVED***
	return string(e)
***REMOVED***

var _ error = ExecutionConflictError("")

// BuildExecutionConfig checks for conflicting options and turns any shortcut
// options (i.e. duration, iterations, stages) into the proper long-form
// scheduler configuration in the execution property.
func BuildExecutionConfig(opts lib.Options) (lib.Options, error) ***REMOVED***
	result := opts

	switch ***REMOVED***
	case opts.Duration.Valid:
		if opts.Iterations.Valid ***REMOVED***
			return result, ExecutionConflictError(
				"using multiple execution config shortcuts (`duration` and `iterations`) simultaneously is not allowed",
			)
		***REMOVED***

		if len(opts.Stages) > 0 ***REMOVED*** // stages isn't nil (not set) and isn't explicitly set to empty
			return result, ExecutionConflictError(
				"using multiple execution config shortcuts (`duration` and `stages`) simultaneously is not allowed",
			)
		***REMOVED***

		if opts.Execution != nil ***REMOVED***
			return result, ExecutionConflictError(
				"using an execution configuration shortcut (`duration`) and `execution` simultaneously is not allowed",
			)
		***REMOVED***

		if opts.Duration.Duration <= 0 ***REMOVED***
			//TODO: move this validation to Validate()?
			return result, ExecutionConflictError(
				"`duration` should be more than 0, for infinite duration use the manual-execution scheduler",
			)
		***REMOVED***

		ds := NewConstantLoopingVUsConfig(lib.DefaultSchedulerName)
		ds.VUs = opts.VUs
		ds.Duration = opts.Duration
		result.Execution = lib.SchedulerConfigMap***REMOVED***lib.DefaultSchedulerName: ds***REMOVED***

	case len(opts.Stages) > 0: // stages isn't nil (not set) and isn't explicitly set to empty
		if opts.Iterations.Valid ***REMOVED***
			return result, ExecutionConflictError(
				"using multiple execution config shortcuts (`stages` and `iterations`) simultaneously is not allowed",
			)
		***REMOVED***

		if opts.Execution != nil ***REMOVED***
			return opts, ExecutionConflictError(
				"using an execution configuration shortcut (`stages`) and `execution` simultaneously is not allowed",
			)
		***REMOVED***

		ds := NewVariableLoopingVUsConfig(lib.DefaultSchedulerName)
		ds.StartVUs = opts.VUs
		for _, s := range opts.Stages ***REMOVED***
			if s.Duration.Valid ***REMOVED***
				ds.Stages = append(ds.Stages, Stage***REMOVED***Duration: s.Duration, Target: s.Target***REMOVED***)
			***REMOVED***
		***REMOVED***
		result.Execution = lib.SchedulerConfigMap***REMOVED***lib.DefaultSchedulerName: ds***REMOVED***

	case opts.Iterations.Valid:
		if opts.Execution != nil ***REMOVED***
			return opts, ExecutionConflictError(
				"using an execution configuration shortcut (`iterations`) and `execution` simultaneously is not allowed",
			)
		***REMOVED***
		// TODO: maybe add a new flag that will be used as a shortcut to per-VU iterations?

		ds := NewSharedIterationsConfig(lib.DefaultSchedulerName)
		ds.VUs = opts.VUs
		ds.Iterations = opts.Iterations
		result.Execution = lib.SchedulerConfigMap***REMOVED***lib.DefaultSchedulerName: ds***REMOVED***

	case len(opts.Execution) > 0:
		// Do nothing, execution was explicitly specified
	default:
		// Check if we should emit some warnings
		if opts.Stages != nil && len(opts.Stages) == 0 ***REMOVED***
			// No someone explicitly set stages to empty
			logrus.Warnf("`stages` was explicitly set to an empty value, running the script with 1 iteration in 1 VU")
		***REMOVED***
		if opts.Execution != nil && len(opts.Execution) == 0 ***REMOVED***
			// No shortcut, and someone explicitly set execution to empty
			logrus.Warnf("`execution` was explicitly set to an empty value, running the script with 1 iteration in 1 VU")
		***REMOVED***
		// No execution parameters whatsoever were specified, so we'll create a per-VU iterations config
		// with 1 VU and 1 iteration. We're choosing the per-VU config, since that one could also
		// be executed both locally, and in the cloud.
		result.Execution = lib.SchedulerConfigMap***REMOVED***
			lib.DefaultSchedulerName: NewPerVUIterationsConfig(lib.DefaultSchedulerName),
		***REMOVED***
		result.Iterations = null.NewInt(1, false)
	***REMOVED***

	//TODO: validate the config; questions:
	// - separately validate the duration, iterations and stages for better error messages?
	// - or reuse the execution validation somehow, at the end? or something mixed?
	// - here or in getConsolidatedConfig() or somewhere else?

	return result, nil
***REMOVED***
