package orchestrator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/redis/v9"
)

func runExecution(gs libOrch.BaseGlobalState, options *libWorker.Options, scope map[string]string, childJobs map[string]jobDistribution, jobId string) (string, error) ***REMOVED***
	// TODO: implement credit check

	// Check if has credits
	/*hasCredits, err := checkIfHasCredits(gs.Ctx(), scope, job)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	if !hasCredits ***REMOVED***
		libOrch.UpdateStatus(gs.Ctx(), gs.Client(), gs.JobId(), gs.OrchestratorId(), "NO_CREDITS")
		return "", errors.New("not enough credits to execute that job")
	***REMOVED****/

	workerSubscriptions := make(map[string]*redis.PubSub)
	for location, jobDistribution := range childJobs ***REMOVED***
		if jobDistribution.jobs != nil && len(*jobDistribution.jobs) > 0 ***REMOVED***
			workerSubscriptions[location] = jobDistribution.workerClient.Subscribe(gs.Ctx(), fmt.Sprintf("worker:executionUpdates:%s", jobId))
		***REMOVED***
	***REMOVED***

	workerChannels := make(map[string]<-chan *redis.Message)
	for location, subscription := range workerSubscriptions ***REMOVED***
		workerChannels[location] = subscription.Channel()
	***REMOVED***

	// Update the status
	libOrch.UpdateStatus(gs.Ctx(), gs.Client(), jobId, gs.OrchestratorId(), "RUNNING")

	// Check if workerSubscriptions is empty
	if len(workerSubscriptions) == 0 ***REMOVED***
		libOrch.DispatchMessage(gs.Ctx(), gs.Client(), gs.JobId(), gs.OrchestratorId(), "No child jobs were created", "INFO")
		return "SUCCESS", nil
	***REMOVED***

	for _, jobDistribution := range childJobs ***REMOVED***
		for _, job := range *jobDistribution.jobs ***REMOVED***
			err := dispatchJob(gs, jobDistribution.workerClient, job, options)
			if err != nil ***REMOVED***
				return "", err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	unifiedChannel := make(chan *redis.Message)

	for _, channel := range workerChannels ***REMOVED***
		go func(channel <-chan *redis.Message) ***REMOVED***
			for msg := range channel ***REMOVED***
				unifiedChannel <- msg
			***REMOVED***
		***REMOVED***(channel)
	***REMOVED***

	for msg := range unifiedChannel ***REMOVED***
		workerMessage := libOrch.WorkerMessage***REMOVED******REMOVED***
		err := json.Unmarshal([]byte(msg.Payload), &workerMessage)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***

		if workerMessage.MessageType == "STATUS" ***REMOVED***
			if workerMessage.Message == "FAILED" ***REMOVED***
				return "FAILURE", nil
			***REMOVED*** else if workerMessage.Message == "SUCCESS" ***REMOVED***
				return "SUCCESS", nil
			***REMOVED*** else ***REMOVED***
				libOrch.DispatchWorkerMessage(gs.Ctx(), gs.Client(), gs.JobId(), workerMessage.WorkerId, workerMessage.Message, "STATUS")
			***REMOVED***
		***REMOVED*** else if workerMessage.MessageType == "METRICS" ***REMOVED***
			(*gs.MetricsStore()).AddMessage(workerMessage, "portsmouth")
		***REMOVED*** else if workerMessage.MessageType == "DEBUG" ***REMOVED***
			// TODO: make this configurable
			continue
		***REMOVED*** else ***REMOVED***
			libOrch.DispatchWorkerMessage(gs.Ctx(), gs.Client(), gs.JobId(), workerMessage.WorkerId, workerMessage.Message, workerMessage.MessageType)
		***REMOVED***

		// Could handle these differently, but for now just dispatch them

		/*else if workerMessage.MessageType == "MARK" ***REMOVED***
			libOrch.DispatchWorkerMessage(gs.Ctx, orchestratorClient, job.Id, workerMessage.WorkerId, workerMessage.Message, "MARK")
		***REMOVED*** else if workerMessage.MessageType == "CONSOLE" ***REMOVED***
			libOrch.DispatchWorkerMessage(gs.Ctx, orchestratorClient, job.Id, workerMessage.WorkerId, workerMessage.Message, "CONSOLE")
		***REMOVED*** else if workerMessage.MessageType == "METRICS" ***REMOVED***
			libOrch.DispatchWorkerMessage(gs.Ctx, orchestratorClient, job.Id, workerMessage.WorkerId, workerMessage.Message, "METRICS")
		***REMOVED*** else if workerMessage.MessageType == "SUMMARY_METRICS" ***REMOVED***
			libOrch.DispatchWorkerMessage(gs.Ctx, orchestratorClient, job.Id, workerMessage.WorkerId, workerMessage.Message, "SUMMARY_METRICS")
			workerSubscription.Close()
		***REMOVED*** else if workerMessage.MessageType == "ERROR" ***REMOVED***
			libOrch.DispatchWorkerMessage(gs.Ctx, orchestratorClient, job.Id, workerMessage.WorkerId, workerMessage.Message, "ERROR")
			libOrch.HandleStringError(gs.Ctx, orchestratorClient, job.Id, gs.orchestratorId, workerMessage.Message)
			return
		***REMOVED*** else if workerMessage.MessageType == "DEBUG" ***REMOVED***
			libOrch.DispatchWorkerMessage(gs.Ctx, orchestratorClient, job.Id, workerMessage.WorkerId, workerMessage.Message, "DEBUG")
		***REMOVED****/
	***REMOVED***

	// Shouldn't get here
	return "", errors.New("an unexpected error occurred")
***REMOVED***

func dispatchJob(gs libOrch.BaseGlobalState, workerClient *redis.Client, job libOrch.ChildJob, options *libWorker.Options) error ***REMOVED***
	// Convert options to json
	marshalledChildJob, err := json.Marshal(job)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	workerClient.HSet(gs.Ctx(), job.ChildJobId, "job", marshalledChildJob)

	workerClient.Publish(gs.Ctx(), "worker:execution", job.ChildJobId)
	workerClient.SAdd(gs.Ctx(), "worker:executionHistory", job.ChildJobId)

	return nil
***REMOVED***

/*
Check i the scope has the required credits to execute the job
*/
func checkIfHasCredits(ctx context.Context, scope map[string]string, job libOrch.Job) (bool, error) ***REMOVED***
	// TODO: implement fully
	return true, nil
	/*
	   // Check max requests has not been reached
	   maxRequests := scope["maxRequests"]

	   	if maxRequests != "" ***REMOVED***
	   		return false, fmt.Errorf("maxRequests not found")
	   	***REMOVED***

	   currentRequests := scope["currentRequests"]

	   	if currentRequests != "" ***REMOVED***
	   		return false, fmt.Errorf("currentRequests not found")
	   	***REMOVED***

	   // TODO: More checks

	   return currentRequests < maxRequests, nil
	*/
***REMOVED***
