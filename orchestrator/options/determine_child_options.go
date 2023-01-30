package options

import (
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
	"gopkg.in/guregu/null.v3"
)

// Scales load options for child jobs proportionately based on their subFraction
func DetermineChildDerivedOptions(loadZone types.LoadZone, options libWorker.Options, subFraction float64) libWorker.Options {
	// Options been pass by value, so we can modify them

	options.MaxPossibleVUs.Int64 = int64(float64(options.MaxPossibleVUs.Int64) * subFraction)

	if options.VUs.Valid {
		options.VUs = null.IntFrom(int64(subFraction * float64(options.VUs.ValueOrZero())))

		if options.VUs.Int64 < 1 {
			options.VUs = null.IntFrom(1)
		}
	}

	if options.Iterations.Valid {
		options.Iterations = null.IntFrom(int64(subFraction * float64(options.Iterations.ValueOrZero())))

		if options.Iterations.Int64 < 1 {
			options.Iterations = null.IntFrom(1)
		}
	}

	for stage := range options.Stages {
		if options.Stages[stage].Target.Valid {
			options.Stages[stage].Target = null.IntFrom(int64(subFraction * float64(options.Stages[stage].Target.ValueOrZero())))

			if options.Stages[stage].Target.Int64 < 1 {
				options.Stages[stage].Target = null.IntFrom(1)
			}
		}
	}

	if options.RPS.Valid {
		options.RPS = null.IntFrom(int64(subFraction * float64(options.RPS.ValueOrZero())))

		if options.RPS.Int64 < 1 {
			options.RPS = null.IntFrom(1)
		}
	}

	if options.Batch.Valid {
		options.Batch = null.IntFrom(int64(subFraction * float64(options.Batch.ValueOrZero())))

		if options.Batch.Int64 < 1 {
			options.Batch = null.IntFrom(1)
		}
	}

	if options.BatchPerHost.Valid {
		options.BatchPerHost = null.IntFrom(int64(subFraction * float64(options.BatchPerHost.ValueOrZero())))

		if options.BatchPerHost.Int64 < 1 {
			options.BatchPerHost = null.IntFrom(1)
		}
	}

	for scenarioName, scenario := range options.Scenarios {
		options.Scenarios[scenarioName] = scenario.ScaleOptions(subFraction)
	}

	return options
}
