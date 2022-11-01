package orchestrator

import (
	"encoding/json"
	"fmt"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
)

func abortChildJobs(gs libOrch.BaseGlobalState, childJobs map[string]jobDistribution) error {
	var err error

	cancelMessage := jobUserUpdate{
		UpdateType: "CANCEL",
	}

	marshalledCancelMessage, err := json.Marshal(cancelMessage)
	if err != nil {
		return err
	}
	stringMarshalledCancelMessage := string(marshalledCancelMessage)

	for _, jobDistribution := range childJobs {
		for _, job := range jobDistribution.Jobs {
			err = jobDistribution.workerClient.Publish(gs.Ctx(), fmt.Sprintf("childJobUserUpdates:%s", job.ChildJobId), stringMarshalledCancelMessage).Err()
		}
	}

	return err
}

func abortAndFailAll(gs libOrch.BaseGlobalState, childJobs map[string]jobDistribution, err error) (string, error) {
	abortErr := abortChildJobs(gs, childJobs)
	if abortErr != nil {
		return "FAILURE", abortErr
	}

	libOrch.UpdateStatus(gs, "FAILURE")

	return "FAILURE", err
}
