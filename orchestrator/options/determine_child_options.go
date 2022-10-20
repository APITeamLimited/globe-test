package options

import (
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
)

// Scales load options for child jobs proportionately based on their subFraction
func DetermineChildDerivedOptions(loadZone types.LoadZone, workerClient *libOrch.NamedClient, parentOptions libWorker.Options, subFraction float32) libWorker.Options ***REMOVED***
	// Copy parent options
	options := parentOptions

	/*
		Don't think this is necessary as only scenarios are used for the execution

		if options.VUs.Valid ***REMOVED***
			options.VUs = null.IntFrom(int64(subFraction * float32(options.VUs.ValueOrZero())))
		***REMOVED***

		if options.Iterations.Valid ***REMOVED***
			options.Iterations = null.IntFrom(int64(subFraction * float32(options.Iterations.ValueOrZero())))
		***REMOVED***

		for stage := range options.Stages ***REMOVED***
			if options.Stages[stage].Target.Valid ***REMOVED***
				options.Stages[stage].Target = null.IntFrom(int64(subFraction * float32(options.Stages[stage].Target.ValueOrZero())))
			***REMOVED***
		***REMOVED***

		if options.RPS.Valid ***REMOVED***
			options.RPS = null.IntFrom(int64(subFraction * float32(options.RPS.ValueOrZero())))
		***REMOVED***

		if options.Batch.Valid ***REMOVED***
			options.Batch = null.IntFrom(int64(subFraction * float32(options.Batch.ValueOrZero())))
		***REMOVED***

		if options.BatchPerHost.Valid ***REMOVED***
			options.BatchPerHost = null.IntFrom(int64(subFraction * float32(options.BatchPerHost.ValueOrZero())))
		***REMOVED***
	*/

	for _, scenario := range options.Scenarios ***REMOVED***
		scenario.ScaleOptions(subFraction)
	***REMOVED***

	return options
***REMOVED***
