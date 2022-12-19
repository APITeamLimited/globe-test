package orchestrator

import (
	"encoding/json"
	"fmt"

	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
)

func abortChildJobs(gs libOrch.BaseGlobalState, childJobs map[string]jobDistribution) error {
	cancelMessage := lib.JobUserUpdate{
		UpdateType: "CANCEL",
	}

	marshalledCancelMessage, err := json.Marshal(cancelMessage)
	if err != nil {
		return err
	}
	stringMarshalledCancelMessage := string(marshalledCancelMessage)

	for _, jobDistribution := range childJobs {
		for _, job := range jobDistribution.Jobs {
			sendError := jobDistribution.workerClient.Publish(gs.Ctx(), fmt.Sprintf("childjobUserUpdates:%s", job.ChildJobId), stringMarshalledCancelMessage).Err()

			if sendError != nil {
				libOrch.HandleError(gs, sendError)
			}
		}
	}

	return err
}

func abortAndFailAll(gs libOrch.BaseGlobalState, childJobs map[string]jobDistribution, err error) (string, error) {
	abortChildJobs(gs, childJobs)

	libOrch.UpdateStatus(gs, "FAILURE")

	return "FAILURE", err
}
