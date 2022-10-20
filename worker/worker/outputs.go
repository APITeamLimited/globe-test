package worker

import (
	"fmt"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/output"
	"github.com/APITeamLimited/globe-test/worker/output/globetest"
)

func createOutputs(gs libWorker.BaseGlobalState) ([]output.Output, error) ***REMOVED***
	result := make([]output.Output, 0, 1)

	// Using globetest output only
	globetestOutput, err := globetest.New(gs)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("could not create the 'globetest' output: %w", err)
	***REMOVED***

	result = append(result, globetestOutput)

	return result, nil
***REMOVED***
