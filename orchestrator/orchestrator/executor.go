package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func run(gs *libOrch.GlobalState, orchestratorId string, orchestratorClient, scopesClient *redis.Client, workerClients map[string]*redis.Client, job map[string]string, storeMongoDB *mongo.Database) ***REMOVED***
	// Get the scope
	scope, err := orchestratorClient.HGetAll(gs.Ctx, job["scopeId"]).Result()
	if err != nil ***REMOVED***
		libOrch.HandleStringError(gs.Ctx, orchestratorClient, job["id"], orchestratorId, fmt.Sprintf("Error getting scope: %s", err.Error()))
		return
	***REMOVED***

	// Check if has credits
	hasCredits, err := checkIfHasCredits(gs.Ctx, scope, job)
	if err != nil ***REMOVED***
		libOrch.HandleStringError(gs.Ctx, orchestratorClient, job["id"], orchestratorId, fmt.Sprintf("Error checking if has credits: %s", err.Error()))
		return
	***REMOVED***
	if !hasCredits ***REMOVED***
		libOrch.DispatchMessage(gs.Ctx, orchestratorClient, job["id"], orchestratorId, "Not enough credits to execute that job", "MESSAGE")
		libOrch.UpdateStatus(gs.Ctx, orchestratorClient, job["id"], orchestratorId, "NO_CREDITS")
		return
	***REMOVED***

	options, err := determineRuntimeOptions(job, gs)
	if err != nil ***REMOVED***
		fmt.Println("Error determining runtime options", err)
		libOrch.HandleError(gs.Ctx, orchestratorClient, job["id"], orchestratorId, err)
		return
	***REMOVED***

	workerClient := workerClients["portsmouth"]
	workerSubscription := workerClient.Subscribe(gs.Ctx, fmt.Sprintf("worker:executionUpdates:%s", job["id"]))

	//marshalledOptions, err := json.Marshal(options)
	if err != nil ***REMOVED***
		libOrch.HandleStringError(gs.Ctx, orchestratorClient, job["id"], orchestratorId, fmt.Sprintf("Error marshalling options: %s", err.Error()))
		return
	***REMOVED***

	// Update the status
	libOrch.UpdateStatus(gs.Ctx, orchestratorClient, job["id"], orchestratorId, "RUNNING")

	err = dispatchJob(gs.Ctx, workerClient, job, "PENDING", orchestratorId, options)
	if err != nil ***REMOVED***
		libOrch.HandleStringError(gs.Ctx, orchestratorClient, job["id"], orchestratorId, fmt.Sprintf("Error dispatching job: %s", err.Error()))
		return
	***REMOVED***

	for msg := range workerSubscription.Channel() ***REMOVED***
		workerMessage := libOrch.WorkerMessage***REMOVED******REMOVED***
		err := json.Unmarshal([]byte(msg.Payload), &workerMessage)
		if err != nil ***REMOVED***
			libOrch.HandleStringError(gs.Ctx, orchestratorClient, job["id"], orchestratorId, fmt.Sprintf("Error unmarshalling worker message: %s", err.Error()))
			return
		***REMOVED***

		if workerMessage.MessageType == "STATUS" ***REMOVED***
			if workerMessage.Message == "FAILED" ***REMOVED***
				libOrch.UpdateStatus(gs.Ctx, orchestratorClient, job["id"], orchestratorId, "FAILED")
				return
			***REMOVED*** else if workerMessage.Message == "SUCCESS" ***REMOVED***
				libOrch.UpdateStatus(gs.Ctx, orchestratorClient, job["id"], orchestratorId, "SUCCESS")
				return
			***REMOVED*** else ***REMOVED***
				libOrch.DispatchWorkerMessage(gs.Ctx, orchestratorClient, job["id"], workerMessage.WorkerId, workerMessage.Message, "STATUS")
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			libOrch.DispatchWorkerMessage(gs.Ctx, orchestratorClient, job["id"], workerMessage.WorkerId, workerMessage.Message, workerMessage.MessageType)
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
			libOrch.HandleStringError(gs.Ctx, orchestratorClient, job["id"], orchestratorId, workerMessage.Message)
			return
		***REMOVED*** else if workerMessage.MessageType == "DEBUG" ***REMOVED***
			libOrch.DispatchWorkerMessage(gs.Ctx, orchestratorClient, job["id"], workerMessage.WorkerId, workerMessage.Message, "DEBUG")
		***REMOVED****/
	***REMOVED***
***REMOVED***

func dispatchJob(ctx context.Context, workerClient *redis.Client, job map[string]string, status string, orchestratorId string, options *libWorker.Options) error ***REMOVED***
	// Convert options to json
	marshalledOptions, err := json.Marshal(options)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	workerClient.HSet(ctx, job["id"], "id", job["id"])
	workerClient.HSet(ctx, job["id"], "source", job["source"])
	workerClient.HSet(ctx, job["id"], "sourceName", job["sourceName"])
	workerClient.HSet(ctx, job["id"], "status", status)
	workerClient.HSet(ctx, job["id"], "scopeId", job["scopeId"])
	workerClient.HSet(ctx, job["id"], "orchestratorId", orchestratorId)
	workerClient.HSet(ctx, job["id"], "options", string(marshalledOptions))

	workerClient.Publish(ctx, "k6:execution", job["id"])
	workerClient.SAdd(ctx, "k6:executionHistory", job["id"])

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
