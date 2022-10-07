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

func getCompiledOptions(job libOrch.Job, gs libOrch.BaseGlobalState) (*libWorker.Options, error) ***REMOVED***
	source, sourceName, err := validateSource(job, gs)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return compileAndGetOptions(source, sourceName, gs)
***REMOVED***

func validateSource(job libOrch.Job, gs libOrch.BaseGlobalState) (string, string, error) ***REMOVED***
	// Check job.SourceName is set
	if job.SourceName == "" ***REMOVED***
		return "", "", errors.New("job.SourceName not set")
	***REMOVED***

	if len(job.SourceName) < 3 ***REMOVED***
		return "", "", errors.New("job.SourceName must be a .js file")
	***REMOVED***

	if job.SourceName[len(job.SourceName)-3:] != ".js" ***REMOVED***
		return "", "", errors.New("job.SourceName must be a .js file")
	***REMOVED***

	// Check source in options, if it is return it
	if job.Source == "" ***REMOVED***
		return "", "", errors.New("source not set")
	***REMOVED***

	return job.Source, job.SourceName, nil
***REMOVED***

func compileAndGetOptions(source string, sourceName string, gs libOrch.BaseGlobalState) (*libWorker.Options, error) ***REMOVED***
	runtimeOptions := libWorker.RuntimeOptions***REMOVED***
		TestType:             null.StringFrom("js"),
		IncludeSystemEnvVars: null.BoolFrom(false),
		CompatibilityMode:    null.StringFrom("extended"),
		NoThresholds:         null.BoolFrom(false),
		SummaryExport:        null.StringFrom(""),
		Env:                  make(map[string]string),
	***REMOVED***

	registry := workerMetrics.NewRegistry()

	preInitState := &libWorker.TestPreInitState***REMOVED***
		// These gs will need to be changed as on the cloud
		Logger:         gs.Logger(),
		RuntimeOptions: runtimeOptions,
		Registry:       registry,
		BuiltinMetrics: workerMetrics.RegisterBuiltinMetrics(registry),
	***REMOVED***

	sourceData := &loader.SourceData***REMOVED***
		Data: []byte(source),
		URL:  &url.URL***REMOVED***Path: sourceName***REMOVED***,
	***REMOVED***

	filesytems := make(map[string]afero.Fs, 1)
	filesytems["file"] = afero.NewMemMapFs()

	// Pass orchestratorId as workerId, so that will dispatch as a worker message
	orchestratorInfo := &libWorker.WorkerInfo***REMOVED***
		Client:         gs.Client(),
		JobId:          gs.JobId(),
		OrchestratorId: gs.OrchestratorId(),
		WorkerId:       gs.OrchestratorId(),
		Ctx:            gs.Ctx(),
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
