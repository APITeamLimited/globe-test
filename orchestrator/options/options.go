package options

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/orchestrator/options/validators"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/executor"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
	"github.com/APITeamLimited/globe-test/worker/metrics"
)

// Import script to determine options on the orchestrator
func DetermineRuntimeOptions(job map[string]string, gs *libOrch.GlobalState) (*libWorker.Options, error) ***REMOVED***
	options, err := getCompiledOptions(job, gs)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Prevent the user from accessing internal ip ranges
	localhostIPNets := generateBannedIPNets()

	// validate the options

	validators.BlacklistIPs(options, localhostIPNets)

	err = validators.Batch(options)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	err = validators.BatchPerHost(options)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	err = validators.Duration(options)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	err = validators.Hosts(options, localhostIPNets)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	err = validators.MinIterationDuration(options)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	validators.InsecureSkipTLSVerify(options)

	err = validators.TeardownTimeout(options)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Check the generated and user supplied options are valid
	checkedOptions, err := executor.DeriveScenariosFromShortcuts(applyDefault(options), gs.Logger)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	marshalledOptions, err := json.Marshal(checkedOptions)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	fmt.Println(string(marshalledOptions))

	return &checkedOptions, nil
***REMOVED***

func applyDefault(options *libWorker.Options) libWorker.Options ***REMOVED***
	if options.SystemTags == nil ***REMOVED***
		options.SystemTags = &metrics.DefaultSystemTagSet
	***REMOVED***
	if options.SummaryTrendStats == nil ***REMOVED***
		options.SummaryTrendStats = libWorker.DefaultSummaryTrendStats
	***REMOVED***
	defDNS := types.DefaultDNSConfig()
	if !options.DNS.TTL.Valid ***REMOVED***
		options.DNS.TTL = defDNS.TTL
	***REMOVED***
	if !options.DNS.Select.Valid ***REMOVED***
		options.DNS.Select = defDNS.Select
	***REMOVED***
	if !options.DNS.Policy.Valid ***REMOVED***
		options.DNS.Policy = defDNS.Policy
	***REMOVED***
	if !options.SetupTimeout.Valid ***REMOVED***
		options.SetupTimeout.Duration = types.Duration(60 * time.Second)
	***REMOVED***
	if !options.TeardownTimeout.Valid ***REMOVED***
		options.TeardownTimeout.Duration = types.Duration(60 * time.Second)
	***REMOVED***

	return *options
***REMOVED***
