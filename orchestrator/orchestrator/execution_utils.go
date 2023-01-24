package orchestrator

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
)

func abortChildJobs(gs libOrch.BaseGlobalState, childJobs map[string]libOrch.ChildJobDistribution) error {
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
			sendError := jobDistribution.WorkerClient.Publish(gs.Ctx(), fmt.Sprintf("childjobUserUpdates:%s", job.ChildJobId), stringMarshalledCancelMessage).Err()

			if sendError != nil {
				libOrch.HandleError(gs, sendError)
			}
		}
	}

	return err
}

func abortAndFailAll(gs libOrch.BaseGlobalState, childJobs map[string]libOrch.ChildJobDistribution, err error) (string, error) {
	libOrch.UpdateStatus(gs, "FAILURE")

	abortChildJobs(gs, childJobs)

	// Send messages again in case they were not received in 10s and 30s
	go func() {
		time.Sleep(10 * time.Second)
		abortChildJobs(gs, childJobs)
		time.Sleep(20 * time.Second)
		abortChildJobs(gs, childJobs)
	}()

	return "FAILURE", err
}
