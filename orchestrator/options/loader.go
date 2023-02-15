package options

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/APITeamLimited/globe-test/js"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/loader"
	"github.com/spf13/afero"
)

func getCompiledOptions(job libOrch.Job, gs libOrch.BaseGlobalState) (*libWorker.Options, error) {
	sourceData, err := loader.LoadTestData(job.TestData)
	if err != nil {
		return nil, fmt.Errorf("failed to load test data: %w", err)
	}

	filesystems := make(map[string]afero.Fs, 1)
	filesystems["file"] = afero.NewMemMapFs()

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
		Registry:       nil,
		BuiltinMetrics: nil,
	}

	bundle, err := js.NewBundle(preInitState, sourceData, filesystems, orchestratorInfo, true, job.TestData)
	if err != nil {
		output := fmt.Sprintf("failed to parse options: %s", err.Error())
		unescaped, err := url.PathUnescape(output)
		if err != nil {
			return nil, errors.New(output)
		}

		return nil, errors.New(unescaped)
	}

	// Get the options export frrom the exports
	return &bundle.Options, nil
}
