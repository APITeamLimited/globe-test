package worker

import (
	"fmt"

	"go.k6.io/k6/lib"
	"go.k6.io/k6/output"
	"go.k6.io/k6/output/json"
)

func createOutputs(
	gs *globalState, test *workerLoadedAndConfiguredTest, executionPlan []lib.ExecutionStep,
) ([]output.Output, error) {
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

	jsonOutput, err := json.New(baseParams)
	if err != nil {
		return nil, fmt.Errorf("could not create the 'json' output: %w", err)
	}

	result = append(result, jsonOutput)

	return result, nil
}
