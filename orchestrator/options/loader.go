package options

import (
	"errors"
	"net/url"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/js"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/loader"
	"github.com/APITeamLimited/globe-test/worker/workerMetrics"
	"github.com/spf13/afero"
	"gopkg.in/guregu/null.v3"
)

func getCompiledOptions(job map[string]string, gs libOrch.BaseGlobalState) (*libWorker.Options, error) {
	source, sourceName, err := validateSource(job, gs)
	if err != nil {
		return nil, err
	}

	return compileAndGetOptions(source, sourceName, gs)
}

func validateSource(job map[string]string, gs libOrch.BaseGlobalState) (string, string, error) {
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

func compileAndGetOptions(source string, sourceName string, gs libOrch.BaseGlobalState) (*libWorker.Options, error) {
	runtimeOptions := libWorker.RuntimeOptions{
		TestType:             null.StringFrom("js"),
		IncludeSystemEnvVars: null.BoolFrom(false),
		CompatibilityMode:    null.StringFrom("extended"),
		NoThresholds:         null.BoolFrom(false),
		SummaryExport:        null.StringFrom(""),
		Env:                  make(map[string]string),
	}

	registry := workerMetrics.NewRegistry()

	preInitState := &libWorker.TestPreInitState{
		// These gs will need to be changed as on the cloud
		Logger:         gs.Logger(),
		RuntimeOptions: runtimeOptions,
		Registry:       registry,
		BuiltinMetrics: workerMetrics.RegisterBuiltinMetrics(registry),
	}

	sourceData := &loader.SourceData{
		Data: []byte(source),
		URL:  &url.URL{Path: sourceName},
	}

	filesytems := make(map[string]afero.Fs, 1)
	filesytems["file"] = afero.NewMemMapFs()

	// Pass orchestratorId as workerId, so that will dispatch as a worker message
	orchestratorInfo := &libWorker.WorkerInfo{
		Client:         gs.Client(),
		JobId:          gs.JobId(),
		OrchestratorId: gs.OrchestratorId(),
		WorkerId:       gs.OrchestratorId(),
		Ctx:            gs.Ctx(),
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
