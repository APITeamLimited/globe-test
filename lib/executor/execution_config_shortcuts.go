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
	"github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/types"
)

// ExecutionConflictError is a custom error type used for all of the errors in
// the DeriveScenariosFromShortcuts() function.
type ExecutionConflictError string

func (e ExecutionConflictError) Error() string ***REMOVED***
	return string(e)
***REMOVED***

var _ error = ExecutionConflictError("")

func getConstantLoopingVUsScenario(duration types.NullDuration, vus null.Int) lib.ExecutorConfigMap ***REMOVED***
	ds := NewConstantLoopingVUsConfig(lib.DefaultExecutorName)
	ds.VUs = vus
	ds.Duration = duration
	return lib.ExecutorConfigMap***REMOVED***lib.DefaultExecutorName: ds***REMOVED***
***REMOVED***

func getVariableLoopingVUsScenario(stages []lib.Stage, startVUs null.Int) lib.ExecutorConfigMap ***REMOVED***
	ds := NewVariableLoopingVUsConfig(lib.DefaultExecutorName)
	ds.StartVUs = startVUs
	for _, s := range stages ***REMOVED***
		if s.Duration.Valid ***REMOVED***
			ds.Stages = append(ds.Stages, Stage***REMOVED***Duration: s.Duration, Target: s.Target***REMOVED***)
		***REMOVED***
	***REMOVED***
	return lib.ExecutorConfigMap***REMOVED***lib.DefaultExecutorName: ds***REMOVED***
***REMOVED***

func getSharedIterationsScenario(iters null.Int, duration types.NullDuration, vus null.Int) lib.ExecutorConfigMap ***REMOVED***
	ds := NewSharedIterationsConfig(lib.DefaultExecutorName)
	ds.VUs = vus
	ds.Iterations = iters
	if duration.Valid ***REMOVED***
		ds.MaxDuration = duration
	***REMOVED***
	return lib.ExecutorConfigMap***REMOVED***lib.DefaultExecutorName: ds***REMOVED***
***REMOVED***

// DeriveScenariosFromShortcuts checks for conflicting options and turns any
// shortcut options (i.e. duration, iterations, stages) into the proper
// long-form scenario/executor configuration in the scenarios property.
func DeriveScenariosFromShortcuts(opts lib.Options) (lib.Options, error) ***REMOVED***
	result := opts

	switch ***REMOVED***
	case opts.Iterations.Valid:
		if len(opts.Stages) > 0 ***REMOVED*** // stages isn't nil (not set) and isn't explicitly set to empty
			return result, ExecutionConflictError(
				"using multiple execution config shortcuts (`iterations` and `stages`) simultaneously is not allowed",
			)
		***REMOVED***
		if opts.Scenarios != nil ***REMOVED***
			return opts, ExecutionConflictError(
				"using an execution configuration shortcut (`iterations`) and `scenarios` simultaneously is not allowed",
			)
		***REMOVED***
		result.Scenarios = getSharedIterationsScenario(opts.Iterations, opts.Duration, opts.VUs)

	case opts.Duration.Valid:
		if len(opts.Stages) > 0 ***REMOVED*** // stages isn't nil (not set) and isn't explicitly set to empty
			return result, ExecutionConflictError(
				"using multiple execution config shortcuts (`duration` and `stages`) simultaneously is not allowed",
			)
		***REMOVED***
		if opts.Scenarios != nil ***REMOVED***
			return result, ExecutionConflictError(
				"using an execution configuration shortcut (`duration`) and `scenarios` simultaneously is not allowed",
			)
		***REMOVED***
		if opts.Duration.Duration <= 0 ***REMOVED***
			//TODO: move this validation to Validate()?
			return result, ExecutionConflictError(
				"`duration` should be more than 0, for infinite duration use the externally-controlled executor",
			)
		***REMOVED***
		result.Scenarios = getConstantLoopingVUsScenario(opts.Duration, opts.VUs)

	case len(opts.Stages) > 0: // stages isn't nil (not set) and isn't explicitly set to empty
		if opts.Scenarios != nil ***REMOVED***
			return opts, ExecutionConflictError(
				"using an execution configuration shortcut (`stages`) and `scenarios` simultaneously is not allowed",
			)
		***REMOVED***
		result.Scenarios = getVariableLoopingVUsScenario(opts.Stages, opts.VUs)

	case len(opts.Scenarios) > 0:
		// Do nothing, scenarios was explicitly specified

	default:
		// Check if we should emit some warnings
		if opts.VUs.Valid && opts.VUs.Int64 != 1 ***REMOVED***
			logrus.Warnf(
				"the `vus=%d` option will be ignored, it only works in conjunction with `iterations`, `duration`, or `stages`",
				opts.VUs.Int64,
			)
		***REMOVED***
		if opts.Stages != nil && len(opts.Stages) == 0 ***REMOVED***
			// No someone explicitly set stages to empty
			logrus.Warnf("`stages` was explicitly set to an empty value, running the script with 1 iteration in 1 VU")
		***REMOVED***
		if opts.Scenarios != nil && len(opts.Scenarios) == 0 ***REMOVED***
			// No shortcut, and someone explicitly set execution to empty
			logrus.Warnf("`scenarios` was explicitly set to an empty value, running the script with 1 iteration in 1 VU")
		***REMOVED***
		// No execution parameters whatsoever were specified, so we'll create a per-VU iterations config
		// with 1 VU and 1 iteration.
		result.Scenarios = lib.ExecutorConfigMap***REMOVED***
			lib.DefaultExecutorName: NewPerVUIterationsConfig(lib.DefaultExecutorName),
		***REMOVED***
	***REMOVED***

	//TODO: validate the config; questions:
	// - separately validate the duration, iterations and stages for better error messages?
	// - or reuse the execution validation somehow, at the end? or something mixed?
	// - here or in getConsolidatedConfig() or somewhere else?

	return result, nil
***REMOVED***
