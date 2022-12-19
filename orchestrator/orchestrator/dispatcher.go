package orchestrator

import (
	"encoding/json"
	"fmt"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/redis/v9"
)

func dispatchChildJobs(gs libOrch.BaseGlobalState, childJobs map[string]jobDistribution) error {
	fmt.Println("dispatchChildJobs", len(childJobs))

	for location, jobDistribution := range childJobs {
		for _, job := range jobDistribution.Jobs {
			err := dispatchChildJob(gs, jobDistribution.workerClient, job, location)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func dispatchChildJob(gs libOrch.BaseGlobalState, workerClient *redis.Client, job libOrch.ChildJob, location string) error {
	// Convert options to json
	marshalledChildJob, err := json.Marshal(job)
	if err != nil {
		return err
	}

	fmt.Printf("Dispatched job %s to worker %s\n", job.ChildJobId, location)

	if gs.FuncAuthClient() == nil {
		workerClient.HSet(gs.Ctx(), job.ChildJobId, "job", marshalledChildJob)

		workerClient.SAdd(gs.Ctx(), "worker:executionHistory", job.ChildJobId)
		workerClient.Publish(gs.Ctx(), "worker:execution", job.ChildJobId)

	} else {
		// If we're in function mode, need to create a google cloud function call
		_, err := gs.FuncAuthClient().ExecuteFunction(location, job)
		if err != nil {
			return err
		}
	}

	return nil
}
