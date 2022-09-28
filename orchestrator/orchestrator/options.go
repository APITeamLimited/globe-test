package orchestrator

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/js"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/loader"
	"github.com/APITeamLimited/globe-test/worker/metrics"
	"github.com/spf13/afero"
	"gopkg.in/guregu/null.v3"
)

// Import script to determine options on the orchestrator
func determineRuntimeOptions(job map[string]string, gs *libOrch.GlobalState) (*libWorker.Options, error) ***REMOVED***
	source, sourceName, err := validateSource(job, gs)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	options, err := compileAndGetOptions(source, sourceName, gs)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Prevent the user from accessing loopback ranges

	localhostIPNets := make([]*net.IPNet, 0, 4)

	localhostIPNets = append(localhostIPNets, &net.IPNet***REMOVED***
		IP:   net.IPv4(10, 0, 0, 0),
		Mask: net.IPv4Mask(255, 0, 0, 0),
	***REMOVED***)

	localhostIPNets = append(localhostIPNets, &net.IPNet***REMOVED***
		IP:   net.IPv4(172, 16, 0, 0),
		Mask: net.IPv4Mask(255, 240, 0, 0),
	***REMOVED***)

	localhostIPNets = append(localhostIPNets, &net.IPNet***REMOVED***
		IP:   net.IPv4(192, 168, 0, 0),
		Mask: net.IPv4Mask(255, 255, 0, 0),
	***REMOVED***)

	localhostIPNets = append(localhostIPNets, &net.IPNet***REMOVED***
		IP:   net.IPv6loopback,
		Mask: net.CIDRMask(128, 128),
	***REMOVED***)

	for _, ip := range options.Hosts ***REMOVED***
		for _, localhostIPNet := range localhostIPNets ***REMOVED***
			netIp := net.ParseIP(string(ip.IP))
			if netIp != nil && localhostIPNet.Contains(netIp) ***REMOVED***
				return nil, fmt.Errorf("invalid host %s", ip)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Add to blacklist
	for _, ipNet := range localhostIPNets ***REMOVED***
		// Blacklist ips takes a struct
		ipStruct := &libWorker.IPNet***REMOVED***
			*ipNet,
		***REMOVED***

		options.BlacklistIPs = append(options.BlacklistIPs, ipStruct)
	***REMOVED***

	marshalledOptions, err := json.Marshal(options)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	fmt.Println("Options", string(marshalledOptions))

	return options, nil
***REMOVED***

func validateSource(job map[string]string, gs *libOrch.GlobalState) (string, string, error) ***REMOVED***
	// Check sourceName is set
	if _, ok := job["sourceName"]; !ok ***REMOVED***
		return "", "", errors.New("sourceName not set")
	***REMOVED***

	sourceName, ok := job["sourceName"]
	if !ok ***REMOVED***
		return "", "", errors.New("sourceName is not a string")
	***REMOVED***

	if len(sourceName) < 3 ***REMOVED***
		return "", "", errors.New("sourceName must be a .js file")
	***REMOVED***

	if sourceName[len(sourceName)-3:] != ".js" ***REMOVED***
		return "", "", errors.New("sourceName must be a .js file")
	***REMOVED***

	source, ok := job["source"]

	// Check source in options, if it is return it
	if !ok ***REMOVED***
		return "", "", errors.New("source not set")
	***REMOVED***

	return source, sourceName, nil
***REMOVED***

func compileAndGetOptions(source string, sourceName string, gs *libOrch.GlobalState) (*libWorker.Options, error) ***REMOVED***
	runtimeOptions := libWorker.RuntimeOptions***REMOVED***
		TestType:             null.StringFrom("js"),
		IncludeSystemEnvVars: null.BoolFrom(false),
		CompatibilityMode:    null.StringFrom("extended"),
		NoThresholds:         null.BoolFrom(false),
		NoSummary:            null.BoolFrom(false),
		SummaryExport:        null.StringFrom(""),
		Env:                  make(map[string]string),
	***REMOVED***

	registry := metrics.NewRegistry()

	preInitState := &libWorker.TestPreInitState***REMOVED***
		// These gs will need to be changed as on the cloud
		Logger:         gs.Logger,
		RuntimeOptions: runtimeOptions,
		Registry:       registry,
		BuiltinMetrics: metrics.RegisterBuiltinMetrics(registry),
	***REMOVED***

	sourceData := &loader.SourceData***REMOVED***
		Data: []byte(source),
		URL:  &url.URL***REMOVED***Path: sourceName***REMOVED***,
	***REMOVED***

	filesytems := make(map[string]afero.Fs, 1)
	filesytems["file"] = afero.NewMemMapFs()

	// Pass orchestratorId as workerId, so that will dispatch as a worker message
	orchestratorInfo := &libWorker.WorkerInfo***REMOVED***
		Client:         gs.Client,
		JobId:          gs.JobId,
		OrchestratorId: gs.OrchestratorId,
		WorkerId:       gs.OrchestratorId,
		Ctx:            gs.Ctx,
		Environment:    nil,
		Collection:     nil,
	***REMOVED***

	bundle, err := js.NewBundle(preInitState, sourceData, filesytems, orchestratorInfo)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Get the options export frrom the exports
	return &bundle.Options, nil
***REMOVED***
