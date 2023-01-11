package orchestrator

import (
	"encoding/json"
	"fmt"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/redis/v9"
)

func dispatchChildJobs(gs libOrch.BaseGlobalState, childJobs map[string]jobDistribution) (*([](chan libOrch.FunctionResult)), error) {
	responseChannels := []chan libOrch.FunctionResult{}

	for location, jobDistribution := range childJobs {
		for _, job := range jobDistribution.Jobs {
			channels, err := dispatchChildJob(gs, jobDistribution.workerClient, job, location)
			if err != nil {
				return nil, err
			}

			responseChannels = append(responseChannels, *channels...)
		}
	}

	return &responseChannels, nil
}

func dispatchChildJob(gs libOrch.BaseGlobalState, workerClient *redis.Client, job libOrch.ChildJob, location string) (*([](chan libOrch.FunctionResult)), error) {
	// Convert options to json
	marshalledChildJob, err := json.Marshal(job)
	if err != nil {
		return nil, err
	}

	responseChannels := []chan libOrch.FunctionResult{}

	if !gs.FuncMode() {
		workerClient.HSet(gs.Ctx(), job.ChildJobId, "job", marshalledChildJob)

		workerClient.SAdd(gs.Ctx(), "worker:executionHistory", job.ChildJobId)
		workerClient.Publish(gs.Ctx(), "worker:execution", job.ChildJobId)
	} else {
		// If we're in function mode, need to create a google cloud function call
		responseCh, err := gs.FuncAuthClient().ExecuteFunction(location, job)
		if err != nil {
			return nil, err
		}

		responseChannels = append(responseChannels, *responseCh)
	}

	fmt.Printf("Dispatched child job %s to worker %s\n", job.ChildJobId, location)

	return &responseChannels, nil
}
