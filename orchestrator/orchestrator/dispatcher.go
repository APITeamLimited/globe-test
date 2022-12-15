package orchestrator

import (
	"encoding/json"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/redis/v9"
)

func dispatchChildJobs(gs libOrch.BaseGlobalState, options *libWorker.Options, childJobs map[string]jobDistribution) error ***REMOVED***
	for _, jobDistribution := range childJobs ***REMOVED***
		for _, job := range jobDistribution.Jobs ***REMOVED***
			err := dispatchChildJob(gs, jobDistribution.workerClient, job, options)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func dispatchChildJob(gs libOrch.BaseGlobalState, workerClient *redis.Client, job libOrch.ChildJob, options *libWorker.Options) error ***REMOVED***
	// Convert options to json
	marshalledChildJob, err := json.Marshal(job)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	workerClient.HSet(gs.Ctx(), job.ChildJobId, "job", marshalledChildJob)

	workerClient.SAdd(gs.Ctx(), "worker:executionHistory", job.ChildJobId)
	workerClient.Publish(gs.Ctx(), "worker:execution", job.ChildJobId)

	if gs.FuncModeInfo() != nil ***REMOVED***
		// If we're in function mode, need to create a google cloud function call
		panic("Not implemented yet")
	***REMOVED***

	return nil
***REMOVED***
