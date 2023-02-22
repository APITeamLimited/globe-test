package worker

import (
	"fmt"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/output"
	"github.com/APITeamLimited/globe-test/worker/output/globetest"
)

func createOutputs(gs libWorker.BaseGlobalState, location string) ([]output.Output, error) {
	result := make([]output.Output, 0, 1)

	// Using globetest output only
	globetestOutput, err := globetest.New(gs, location)
	if err != nil {
		return nil, fmt.Errorf("could not create the 'globetest' output: %w", err)
	}

	result = append(result, globetestOutput)

	return result, nil
}
