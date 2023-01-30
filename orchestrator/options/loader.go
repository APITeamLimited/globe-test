package options

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/js"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/loader"
	"github.com/spf13/afero"
	"gopkg.in/guregu/null.v3"
)

func getCompiledOptions(job libOrch.Job, gs libOrch.BaseGlobalState) (*libWorker.Options, error) {
	source, sourceName, err := validateSource(job, gs)
	if err != nil {
		return nil, err
	}

	return compileAndGetOptions(source, sourceName, gs)
}

func validateSource(job libOrch.Job, gs libOrch.BaseGlobalState) (string, string, error) {
	// Check job.SourceName is set
	if job.SourceName == "" {
		return "", "", errors.New("job.SourceName not set")
	}

	if len(job.SourceName) < 3 {
		return "", "", errors.New("job.SourceName must be a .js file")
	}

	if job.SourceName[len(job.SourceName)-3:] != ".js" {
		return "", "", errors.New("job.SourceName must be a .js file")
	}

	// Check source in options, if it is return it
	if job.Source == "" {
		return "", "", errors.New("source not set")
	}

	return job.Source, job.SourceName, nil
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

	sourceData := &loader.SourceData{
		Data: []byte(source),
		URL:  &url.URL{Path: sourceName},
	}

	filesytems := make(map[string]afero.Fs, 1)
	filesytems["file"] = afero.NewMemMapFs()

	// Pass orchestratorId as workerId, so that will dispatch as a worker message
	orchestratorInfo := &libWorker.WorkerInfo{
		Conn:           nil,
		JobId:          gs.JobId(),
		OrchestratorId: gs.OrchestratorId(),
		WorkerId:       gs.OrchestratorId(),
		Ctx:            gs.Ctx(),
		Environment:    nil,
		Collection:     nil,
	}

	preInitState := &libWorker.TestPreInitState{
		// These gs will need to be changed as on the cloud
		Logger:         gs.Logger(),
		RuntimeOptions: runtimeOptions,
		Registry:       nil,
		BuiltinMetrics: nil,
	}

	bundle, err := js.NewBundleUnsafe(preInitState, sourceData, filesytems, orchestratorInfo, true)
	if err != nil {
		return nil, fmt.Errorf("failed to parse options: %w", err)
	}

	// Get the options export frrom the exports
	return &bundle.Options, nil
}
