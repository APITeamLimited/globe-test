package options

import (
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
	"gopkg.in/guregu/null.v3"
)

// Scales load options for child jobs proportionately based on their subFraction
func DetermineChildDerivedOptions(loadZone types.LoadZone, workerClient *libOrch.NamedClient, options libWorker.Options, subFraction float64) libWorker.Options ***REMOVED***
	// Options been pass by value, so we can modify them

	options.MaxPossibleVUs.Int64 = int64(float64(options.MaxPossibleVUs.Int64) * subFraction)

	if options.VUs.Valid ***REMOVED***
		options.VUs = null.IntFrom(int64(subFraction * float64(options.VUs.ValueOrZero())))

		if options.VUs.Int64 < 1 ***REMOVED***
			options.VUs = null.IntFrom(1)
		***REMOVED***
	***REMOVED***

	if options.Iterations.Valid ***REMOVED***
		options.Iterations = null.IntFrom(int64(subFraction * float64(options.Iterations.ValueOrZero())))

		if options.Iterations.Int64 < 1 ***REMOVED***
			options.Iterations = null.IntFrom(1)
		***REMOVED***
	***REMOVED***

	for stage := range options.Stages ***REMOVED***
		if options.Stages[stage].Target.Valid ***REMOVED***
			options.Stages[stage].Target = null.IntFrom(int64(subFraction * float64(options.Stages[stage].Target.ValueOrZero())))

			if options.Stages[stage].Target.Int64 < 1 ***REMOVED***
				options.Stages[stage].Target = null.IntFrom(1)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if options.RPS.Valid ***REMOVED***
		options.RPS = null.IntFrom(int64(subFraction * float64(options.RPS.ValueOrZero())))

		if options.RPS.Int64 < 1 ***REMOVED***
			options.RPS = null.IntFrom(1)
		***REMOVED***
	***REMOVED***

	if options.Batch.Valid ***REMOVED***
		options.Batch = null.IntFrom(int64(subFraction * float64(options.Batch.ValueOrZero())))

		if options.Batch.Int64 < 1 ***REMOVED***
			options.Batch = null.IntFrom(1)
		***REMOVED***
	***REMOVED***

	if options.BatchPerHost.Valid ***REMOVED***
		options.BatchPerHost = null.IntFrom(int64(subFraction * float64(options.BatchPerHost.ValueOrZero())))

		if options.BatchPerHost.Int64 < 1 ***REMOVED***
			options.BatchPerHost = null.IntFrom(1)
		***REMOVED***
	***REMOVED***

	for scenarioName, scenario := range options.Scenarios ***REMOVED***
		options.Scenarios[scenarioName] = scenario.ScaleOptions(subFraction)
	***REMOVED***

	return options
***REMOVED***
