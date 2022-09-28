package orchestrator

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"

	"github.com/APITeamLimited/k6-worker/js"
	libWorker "github.com/APITeamLimited/k6-worker/lib"
	"github.com/APITeamLimited/k6-worker/loader"
	"github.com/APITeamLimited/k6-worker/metrics"
	"github.com/spf13/afero"
	"gitlab.com/apiteamcloud/orchestrator/lib"
	"gopkg.in/guregu/null.v3"
)

// Import script to determine options on the orchestrator
func determineRuntimeOptions(job map[string]string, gs *lib.GlobalState) (*libWorker.Options, error) {
	source, sourceName, err := validateSource(job, gs)
	if err != nil {
		return nil, err
	}

	options, err := compileAndGetOptions(source, sourceName, gs)
	if err != nil {
		return nil, err
	}

	// Prevent the user from accessing loopback ranges

	localhostIPNets := make([]*net.IPNet, 0, 4)

	localhostIPNets = append(localhostIPNets, &net.IPNet{
		IP:   net.IPv4(10, 0, 0, 0),
		Mask: net.IPv4Mask(255, 0, 0, 0),
	})

	localhostIPNets = append(localhostIPNets, &net.IPNet{
		IP:   net.IPv4(172, 16, 0, 0),
		Mask: net.IPv4Mask(255, 240, 0, 0),
	})

	localhostIPNets = append(localhostIPNets, &net.IPNet{
		IP:   net.IPv4(192, 168, 0, 0),
		Mask: net.IPv4Mask(255, 255, 0, 0),
	})

	localhostIPNets = append(localhostIPNets, &net.IPNet{
		IP:   net.IPv6loopback,
		Mask: net.CIDRMask(128, 128),
	})

	for _, ip := range options.Hosts {
		for _, localhostIPNet := range localhostIPNets {
			if localhostIPNet.Contains(ip) {
				return nil, fmt.Errorf("invalid host %s", ip)
			}
		}
	}

	// Add to blacklist
	for _, ipNet := range localhostIPNets {
		// Blacklist ips takes a struct
		ipStruct := &libWorker.IPNet{
			*ipNet,
		}

		options.BlacklistIPs = append(options.BlacklistIPs, ipStruct)
	}

	marshalledOptions, err := json.Marshal(options)
	if err != nil {
		panic(err)
	}
	fmt.Println("Options", string(marshalledOptions))

	return options, nil
}

func validateSource(job map[string]string, gs *lib.GlobalState) (string, string, error) {
	// Check sourceName is set
	if _, ok := job["sourceName"]; !ok {
		return "", "", errors.New("sourceName not set")
	}

	sourceName, ok := job["sourceName"]
	if !ok {
		return "", "", errors.New("sourceName is not a string")
	}

	if len(sourceName) < 3 {
		return "", "", errors.New("sourceName must be a .js file")
	}

	if sourceName[len(sourceName)-3:] != ".js" {
		return "", "", errors.New("sourceName must be a .js file")
	}

	source, ok := job["source"]

	// Check source in options, if it is return it
	if !ok {
		return "", "", errors.New("source not set")
	}

	return source, sourceName, nil
}

func compileAndGetOptions(source string, sourceName string, gs *lib.GlobalState) (*libWorker.Options, error) {
	runtimeOptions := libWorker.RuntimeOptions{
		TestType:             null.StringFrom("js"),
		IncludeSystemEnvVars: null.BoolFrom(false),
		CompatibilityMode:    null.StringFrom("extended"),
		NoThresholds:         null.BoolFrom(false),
		NoSummary:            null.BoolFrom(false),
		SummaryExport:        null.StringFrom(""),
		Env:                  make(map[string]string),
	}

	registry := metrics.NewRegistry()

	preInitState := &libWorker.TestPreInitState{
		// These gs will need to be changed as on the cloud
		Logger:         gs.Logger,
		RuntimeOptions: runtimeOptions,
		Registry:       registry,
		BuiltinMetrics: metrics.RegisterBuiltinMetrics(registry),
	}

	sourceData := &loader.SourceData{
		Data: []byte(source),
		URL:  &url.URL{Path: sourceName},
	}

	filesytems := make(map[string]afero.Fs, 1)
	filesytems["file"] = afero.NewMemMapFs()

	// Pass orchestratorId as workerId, so that will dispatch as a worker message
	orchestratorInfo := &libWorker.WorkerInfo{
		Client:         gs.Client,
		JobId:          gs.JobId,
		OrchestratorId: gs.OrchestratorId,
		WorkerId:       gs.OrchestratorId,
		Ctx:            gs.Ctx,
		Environment:    nil,
		Collection:     nil,
	}

	bundle, err := js.NewBundle(preInitState, sourceData, filesytems, orchestratorInfo)
	if err != nil {
		return nil, err
	}

	// Get the options export frrom the exports
	return &bundle.Options, nil
}
