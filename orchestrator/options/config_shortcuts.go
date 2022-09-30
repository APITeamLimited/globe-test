package options

import (
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/executor"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
	"github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"
)

// ExecutionConflictError is a custom error type used for all of the errors in
// the DeriveScenariosFromShortcuts() function.
type ExecutionConflictError string

func (e ExecutionConflictError) Error() string ***REMOVED***
	return string(e)
***REMOVED***

var _ error = ExecutionConflictError("")

func getConstantVUsScenario(duration types.NullDuration, vus null.Int) libWorker.ScenarioConfigs ***REMOVED***
	ds := executor.NewConstantVUsConfig(libWorker.DefaultScenarioName)
	ds.VUs = vus
	ds.Duration = duration
	return libWorker.ScenarioConfigs***REMOVED***libWorker.DefaultScenarioName: ds***REMOVED***
***REMOVED***

func getRampingVUsScenario(stages []libWorker.Stage, startVUs null.Int) libWorker.ScenarioConfigs ***REMOVED***
	ds := executor.NewRampingVUsConfig(libWorker.DefaultScenarioName)
	ds.StartVUs = startVUs
	for _, s := range stages ***REMOVED***
		if s.Duration.Valid ***REMOVED***
			ds.Stages = append(ds.Stages, executor.Stage***REMOVED***Duration: s.Duration, Target: s.Target***REMOVED***)
		***REMOVED***
	***REMOVED***
	return libWorker.ScenarioConfigs***REMOVED***libWorker.DefaultScenarioName: ds***REMOVED***
***REMOVED***

func getSharedIterationsScenario(iters null.Int, duration types.NullDuration, vus null.Int) libWorker.ScenarioConfigs ***REMOVED***
	ds := executor.NewSharedIterationsConfig(libWorker.DefaultScenarioName)
	ds.VUs = vus
	ds.Iterations = iters
	if duration.Valid ***REMOVED***
		ds.MaxDuration = duration
	***REMOVED***
	return libWorker.ScenarioConfigs***REMOVED***libWorker.DefaultScenarioName: ds***REMOVED***
***REMOVED***

// deriveScenariosFromShortcuts checks for conflicting options and turns any
// shortcut options (i.e. duration, iterations, stages) into the proper
// long-form scenario/executor configuration in the scenarios property.
func deriveScenariosFromShortcuts(opts libWorker.Options, logger logrus.FieldLogger) (libWorker.Options, error) ***REMOVED***
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
			// TODO: move this validation to Validate()?
			return result, ExecutionConflictError(
				"`duration` should be more than 0, for infinite duration use the externally-controlled executor",
			)
		***REMOVED***
		result.Scenarios = getConstantVUsScenario(opts.Duration, opts.VUs)

	case len(opts.Stages) > 0: // stages isn't nil (not set) and isn't explicitly set to empty
		if opts.Scenarios != nil ***REMOVED***
			return opts, ExecutionConflictError(
				"using an execution configuration shortcut (`stages`) and `scenarios` simultaneously is not allowed",
			)
		***REMOVED***
		result.Scenarios = getRampingVUsScenario(opts.Stages, opts.VUs)

	case len(opts.Scenarios) > 0:
		// Do nothing, scenarios was explicitly specified

	default:
		// Check if we should emit some warnings
		if opts.VUs.Valid && opts.VUs.Int64 != 1 ***REMOVED***
			logger.Warnf(
				"the `vus=%d` option will be ignored, it only works in conjunction with `iterations`, `duration`, or `stages`",
				opts.VUs.Int64,
			)
		***REMOVED***
		if opts.Stages != nil && len(opts.Stages) == 0 ***REMOVED***
			// No someone explicitly set stages to empty
			logger.Warnf("`stages` was explicitly set to an empty value, running the script with 1 iteration in 1 VU")
		***REMOVED***
		if opts.Scenarios != nil && len(opts.Scenarios) == 0 ***REMOVED***
			// No shortcut, and someone explicitly set execution to empty
			logger.Warnf("`scenarios` was explicitly set to an empty value, running the script with 1 iteration in 1 VU")
		***REMOVED***
		// No execution parameters whatsoever were specified, so we'll create a per-VU iterations config
		// with 1 VU and 1 iteration.
		result.Scenarios = libWorker.ScenarioConfigs***REMOVED***
			libWorker.DefaultScenarioName: executor.NewPerVUIterationsConfig(libWorker.DefaultScenarioName),
		***REMOVED***
	***REMOVED***

	// TODO: validate the config; questions:
	// - separately validate the duration, iterations and stages for better error messages?
	// - or reuse the execution validation somehow, at the end? or something mixed?
	// - here or in getConsolidatedConfig() or somewhere else?

	return result, nil
***REMOVED***
