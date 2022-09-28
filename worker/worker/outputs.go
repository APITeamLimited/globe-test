package worker

import (
	"fmt"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/output"
	"github.com/APITeamLimited/globe-test/worker/output/globetest"
)

func createOutputs(
	gs *globalState, test *workerLoadedAndConfiguredTest, executionPlan []libWorker.ExecutionStep,
	workerInfo *libWorker.WorkerInfo) ([]output.Output, error) {
	baseParams := output.Params{
		ScriptPath:     test.source.URL,
		Logger:         gs.logger,
		Environment:    gs.envVars,
		StdOut:         gs.stdOut,
		StdErr:         gs.stdErr,
		FS:             gs.fs,
		ScriptOptions:  test.derivedConfig.Options,
		RuntimeOptions: test.preInitState.RuntimeOptions,
		ExecutionPlan:  executionPlan,
	}
	result := make([]output.Output, 0, 1)

	// Using globetest output only
	globetestOutput, err := globetest.New(baseParams, workerInfo)

	if err != nil {
		return nil, fmt.Errorf("could not create the 'globetest' output: %w", err)
	}

	result = append(result, globetestOutput)

	return result, nil
}
