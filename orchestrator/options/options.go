package options

import (
	"database/sql"
	"errors"
	"time"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/orchestrator/options/validators"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
	"github.com/APITeamLimited/globe-test/worker/workerMetrics"
	"gopkg.in/guregu/null.v3"
)

// Import script to determine options on the orchestrator
func DetermineRuntimeOptions(job libOrch.Job, gs libOrch.BaseGlobalState, workerClients libOrch.WorkerClients) (*libWorker.Options, error) {
	options, err := getCompiledOptions(job, gs)
	if err != nil {
		return nil, err
	}

	// TODO: Check if MetricSamplesBufferSize config option is needed
	options.MetricSamplesBufferSize = null.Int{
		NullInt64: sql.NullInt64{
			Int64: 0,
			Valid: false,
		},
	}

	// Prevent the user from accessing internal ip ranges
	localhostIPNets := generateBannedIPNets()

	// validate the options

	validators.BlacklistIPs(options, localhostIPNets)

	err = validators.Batch(options)
	if err != nil {
		return nil, err
	}

	err = validators.BatchPerHost(options)
	if err != nil {
		return nil, err
	}

	if options.Duration.Valid {
		err = validators.Duration(options)
		if err != nil {
			return nil, err
		}
	}

	err = validators.Hosts(options, localhostIPNets)
	if err != nil {
		return nil, err
	}

	err = validators.MinIterationDuration(options)
	if err != nil {
		return nil, err
	}

	validators.InsecureSkipTLSVerify(options)

	err = validators.TeardownTimeout(options)
	if err != nil {
		return nil, err
	}

	err = validators.ExecutionMode(options)
	if err != nil {
		return nil, err
	}

	err = validators.LoadDistribution(options, workerClients)
	if err != nil {
		return nil, err
	}

	err = validators.OutputConfig(options)
	if err != nil {
		return nil, err
	}

	// Check the generated and user supplied options are valid
	checkedOptions, err := deriveScenariosFromShortcuts(applyDefault(options), gs.Logger())
	if err != nil {
		return nil, err
	}

	// Add max possible VU count
	maxVUsCount := int64(0)
	for _, scenario := range checkedOptions.Scenarios {
		maxVUsCount += scenario.GetMaxExecutorVUs()
	}

	// Restrict VU count to 50 for now
	if maxVUsCount > 50 {
		return nil, errors.New("max VU count is currently limited to 50")
	}

	checkedOptions.MaxPossibleVUs = null.IntFrom(maxVUsCount)

	return &checkedOptions, nil
}

func applyDefault(options *libWorker.Options) libWorker.Options {
	if options.SystemTags == nil {
		options.SystemTags = &workerMetrics.DefaultSystemTagSet
	}
	if options.SummaryTrendStats == nil {
		options.SummaryTrendStats = libWorker.DefaultSummaryTrendStats
	}
	defDNS := types.DefaultDNSConfig()
	if !options.DNS.TTL.Valid {
		options.DNS.TTL = defDNS.TTL
	}
	if !options.DNS.Select.Valid {
		options.DNS.Select = defDNS.Select
	}
	if !options.DNS.Policy.Valid {
		options.DNS.Policy = defDNS.Policy
	}
	if !options.SetupTimeout.Valid {
		options.SetupTimeout.Duration = types.Duration(60 * time.Second)
	}
	if !options.TeardownTimeout.Valid {
		options.TeardownTimeout.Duration = types.Duration(60 * time.Second)
	}

	return *options
}
