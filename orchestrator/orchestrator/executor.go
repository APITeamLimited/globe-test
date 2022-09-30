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

func runExecution(gs libOrch.BaseGlobalState, options *libWorker.Options, scope map[string]string, workerClients map[string]*redis.Client, job map[string]string) (string, error) ***REMOVED***
	// Check if has credits
	hasCredits, err := checkIfHasCredits(gs.Ctx(), scope, job)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	if !hasCredits ***REMOVED***
		libOrch.UpdateStatus(gs.Ctx(), gs.Client(), gs.JobId(), gs.OrchestratorId(), "NO_CREDITS")
		return "", errors.New("not enough credits to execute that job")
	***REMOVED***

	workerClient := workerClients["portsmouth"]
	workerSubscription := workerClient.Subscribe(gs.Ctx(), fmt.Sprintf("worker:executionUpdates:%s", job["id"]))

	//marshalledOptions, err := json.Marshal(options)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	// Update the status
	libOrch.UpdateStatus(gs.Ctx(), gs.Client(), job["id"], gs.OrchestratorId(), "RUNNING")

	err = dispatchJob(gs, workerClient, job, "PENDING", options)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	for msg := range workerSubscription.Channel() ***REMOVED***
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
			libOrch.DispatchWorkerMessage(gs.Ctx, orchestratorClient, job["id"], workerMessage.WorkerId, workerMessage.Message, "MARK")
		***REMOVED*** else if workerMessage.MessageType == "CONSOLE" ***REMOVED***
			libOrch.DispatchWorkerMessage(gs.Ctx, orchestratorClient, job["id"], workerMessage.WorkerId, workerMessage.Message, "CONSOLE")
		***REMOVED*** else if workerMessage.MessageType == "METRICS" ***REMOVED***
			libOrch.DispatchWorkerMessage(gs.Ctx, orchestratorClient, job["id"], workerMessage.WorkerId, workerMessage.Message, "METRICS")
		***REMOVED*** else if workerMessage.MessageType == "SUMMARY_METRICS" ***REMOVED***
			libOrch.DispatchWorkerMessage(gs.Ctx, orchestratorClient, job["id"], workerMessage.WorkerId, workerMessage.Message, "SUMMARY_METRICS")
			workerSubscription.Close()
		***REMOVED*** else if workerMessage.MessageType == "ERROR" ***REMOVED***
			libOrch.DispatchWorkerMessage(gs.Ctx, orchestratorClient, job["id"], workerMessage.WorkerId, workerMessage.Message, "ERROR")
			libOrch.HandleStringError(gs.Ctx, orchestratorClient, job["id"], gs.orchestratorId, workerMessage.Message)
			return
		***REMOVED*** else if workerMessage.MessageType == "DEBUG" ***REMOVED***
			libOrch.DispatchWorkerMessage(gs.Ctx, orchestratorClient, job["id"], workerMessage.WorkerId, workerMessage.Message, "DEBUG")
		***REMOVED****/
	***REMOVED***

	// Shouldn't get here
	return "", errors.New("an unexpected error occurred")
***REMOVED***

func dispatchJob(gs libOrch.BaseGlobalState, workerClient *redis.Client, job map[string]string, status string, options *libWorker.Options) error ***REMOVED***
	// Convert options to json
	marshalledOptions, err := json.Marshal(options)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	workerClient.HSet(gs.Ctx(), job["id"], "id", job["id"])
	workerClient.HSet(gs.Ctx(), job["id"], "source", job["source"])
	workerClient.HSet(gs.Ctx(), job["id"], "sourceName", job["sourceName"])
	workerClient.HSet(gs.Ctx(), job["id"], "status", status)
	workerClient.HSet(gs.Ctx(), job["id"], "scopeId", job["scopeId"])
	workerClient.HSet(gs.Ctx(), job["id"], "gs.orchestratorId", gs.OrchestratorId)
	workerClient.HSet(gs.Ctx(), job["id"], "options", string(marshalledOptions))

	workerClient.Publish(gs.Ctx(), "k6:execution", job["id"])
	workerClient.SAdd(gs.Ctx(), "k6:executionHistory", job["id"])

	return nil
***REMOVED***

/*
Check i the scope has the required credits to execute the job
*/
func checkIfHasCredits(ctx context.Context, scope map[string]string, job map[string]string) (bool, error) ***REMOVED***
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
