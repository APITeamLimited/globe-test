package orchestrator

import (
	"encoding/json"
	"fmt"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/redis/v9"
)

func dispatchChildJobs(gs libOrch.BaseGlobalState, childJobs map[string]libOrch.ChildJobDistribution) (*([](chan libOrch.FunctionResult)), error) {
	responseChannels := []chan libOrch.FunctionResult{}

	fmt.Printf("Dispatching %d child jobs\n", len(childJobs))

	for location, jobDistribution := range childJobs {
		for _, job := range jobDistribution.Jobs {
			channels, err := dispatchChildJob(gs, jobDistribution.WorkerClient, job, location)
			if err != nil {
				fmt.Printf("Error dispatching child job %s to worker %s: %s\n", job.ChildJobId, location, err)
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

	if gs.FuncMode() {
		// If we're in function mode, need to create a google cloud function call
		responseCh, err := gs.FuncAuthClient().ExecuteFunction(location, job)
		if err != nil {
			return nil, err
		}

		responseChannels = append(responseChannels, *responseCh)
	} else {
		workerClient.HSet(gs.Ctx(), job.ChildJobId, "job", marshalledChildJob)

		workerClient.SAdd(gs.Ctx(), "worker:executionHistory", job.ChildJobId)
		workerClient.Publish(gs.Ctx(), "worker:execution", job.ChildJobId)
	}

	fmt.Printf("Dispatched child job %s to worker %s\n", job.ChildJobId, location)

	return &responseChannels, nil
}
