package options

import (
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
)

// Scales load options for child jobs proportionately based on their subFraction
func DetermineChildDerivedOptions(loadZone types.LoadZone, workerClient *libOrch.NamedClient, parentOptions libWorker.Options, subFraction float32) libWorker.Options {
	// Copy parent options
	options := parentOptions

	/*
		Don't think this is necessary as only scenarios are used for the execution

		if options.VUs.Valid {
			options.VUs = null.IntFrom(int64(subFraction * float32(options.VUs.ValueOrZero())))
		}

		if options.Iterations.Valid {
			options.Iterations = null.IntFrom(int64(subFraction * float32(options.Iterations.ValueOrZero())))
		}

		for stage := range options.Stages {
			if options.Stages[stage].Target.Valid {
				options.Stages[stage].Target = null.IntFrom(int64(subFraction * float32(options.Stages[stage].Target.ValueOrZero())))
			}
		}

		if options.RPS.Valid {
			options.RPS = null.IntFrom(int64(subFraction * float32(options.RPS.ValueOrZero())))
		}

		if options.Batch.Valid {
			options.Batch = null.IntFrom(int64(subFraction * float32(options.Batch.ValueOrZero())))
		}

		if options.BatchPerHost.Valid {
			options.BatchPerHost = null.IntFrom(int64(subFraction * float32(options.BatchPerHost.ValueOrZero())))
		}
	*/

	for _, scenario := range options.Scenarios {
		scenario.ScaleOptions(subFraction)
	}

	return options
}
