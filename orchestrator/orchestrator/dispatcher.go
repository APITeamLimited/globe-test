package orchestrator

import (
	"encoding/json"
	"fmt"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/redis/v9"
)

type childJobDispatchResult struct {
	responseChannel *chan libOrch.FunctionResult
	err             error
}

func dispatchChildJobs(gs libOrch.BaseGlobalState, childJobs map[string]libOrch.ChildJobDistribution) (*([](chan libOrch.FunctionResult)), error) {
	unifiedDispatchResultCh := make(chan childJobDispatchResult)

	for location, jobDistribution := range childJobs {
		for _, job := range jobDistribution.Jobs {
			dispatchResultCh := dispatchChildJob(gs, jobDistribution.WorkerClient, job, location)

			go func(dispatchCh chan childJobDispatchResult) {
				for v := range dispatchCh {
					unifiedDispatchResultCh <- v
				}
			}(dispatchResultCh)
		}
	}

	responseChannels := []chan libOrch.FunctionResult{}

	successFullDispatches := 0

	for dispatchResult := range unifiedDispatchResultCh {
		if dispatchResult.err != nil {
			return nil, dispatchResult.err
		}

		if dispatchResult.responseChannel != nil {
			responseChannels = append(responseChannels, *dispatchResult.responseChannel)
		}

		successFullDispatches++
	}

	return &responseChannels, nil
}

func dispatchChildJob(gs libOrch.BaseGlobalState, workerClient *redis.Client, job libOrch.ChildJob, location string) chan childJobDispatchResult {
	dispatchResultCh := make(chan childJobDispatchResult)

	// Convert options to json
	go func() {
		marshalledChildJob, err := json.Marshal(job)
		if err != nil {
			dispatchResultCh <- childJobDispatchResult{
				err: err,
			}

			return
		}

		if gs.FuncMode() {
			// If we're in function mode, need to create a google cloud function call
			responseCh, err := gs.FuncAuthClient().ExecuteFunction(location, job)
			if err != nil {
				dispatchResultCh <- childJobDispatchResult{
					err: err,
				}

				return
			}

			dispatchResultCh <- childJobDispatchResult{
				responseChannel: responseCh,
			}

			return
		}
		workerClient.HSet(gs.Ctx(), job.ChildJobId, "job", marshalledChildJob)

		workerClient.SAdd(gs.Ctx(), "worker:executionHistory", job.ChildJobId)
		workerClient.Publish(gs.Ctx(), "worker:execution", job.ChildJobId)

		fmt.Printf("Dispatched child job %s to worker %s\n", job.ChildJobId, location)

		dispatchResultCh <- childJobDispatchResult{}
	}()

	return dispatchResultCh
}
