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

func run(gs *libOrch.GlobalState, orchestratorId string, orchestratorClient, scopesClient *redis.Client, workerClients map[string]*redis.Client, job map[string]string, storeMongoDB *mongo.Database) {
	// Get the scope
	scope, err := orchestratorClient.HGetAll(gs.Ctx, job["scopeId"]).Result()
	if err != nil {
		libOrch.HandleStringError(gs.Ctx, orchestratorClient, job["id"], orchestratorId, fmt.Sprintf("Error getting scope: %s", err.Error()))
		return
	}

	// Check if has credits
	hasCredits, err := checkIfHasCredits(gs.Ctx, scope, job)
	if err != nil {
		libOrch.HandleStringError(gs.Ctx, orchestratorClient, job["id"], orchestratorId, fmt.Sprintf("Error checking if has credits: %s", err.Error()))
		return
	}
	if !hasCredits {
		libOrch.DispatchMessage(gs.Ctx, orchestratorClient, job["id"], orchestratorId, "Not enough credits to execute that job", "MESSAGE")
		libOrch.UpdateStatus(gs.Ctx, orchestratorClient, job["id"], orchestratorId, "NO_CREDITS")
		return
	}

	options, err := determineRuntimeOptions(job, gs)
	if err != nil {
		fmt.Println("Error determining runtime options", err)
		libOrch.HandleError(gs.Ctx, orchestratorClient, job["id"], orchestratorId, err)
		return
	}

	workerClient := workerClients["portsmouth"]
	workerSubscription := workerClient.Subscribe(gs.Ctx, fmt.Sprintf("worker:executionUpdates:%s", job["id"]))

	//marshalledOptions, err := json.Marshal(options)
	if err != nil {
		libOrch.HandleStringError(gs.Ctx, orchestratorClient, job["id"], orchestratorId, fmt.Sprintf("Error marshalling options: %s", err.Error()))
		return
	}

	// Update the status
	libOrch.UpdateStatus(gs.Ctx, orchestratorClient, job["id"], orchestratorId, "RUNNING")

	err = dispatchJob(gs.Ctx, workerClient, job, "PENDING", orchestratorId, options)
	if err != nil {
		libOrch.HandleStringError(gs.Ctx, orchestratorClient, job["id"], orchestratorId, fmt.Sprintf("Error dispatching job: %s", err.Error()))
		return
	}

	for msg := range workerSubscription.Channel() {
		workerMessage := libOrch.WorkerMessage{}
		err := json.Unmarshal([]byte(msg.Payload), &workerMessage)
		if err != nil {
			libOrch.HandleStringError(gs.Ctx, orchestratorClient, job["id"], orchestratorId, fmt.Sprintf("Error unmarshalling worker message: %s", err.Error()))
			return
		}

		if workerMessage.MessageType == "STATUS" {
			if workerMessage.Message == "FAILED" {
				libOrch.UpdateStatus(gs.Ctx, orchestratorClient, job["id"], orchestratorId, "FAILED")
				return
			} else if workerMessage.Message == "SUCCESS" {
				libOrch.UpdateStatus(gs.Ctx, orchestratorClient, job["id"], orchestratorId, "SUCCESS")
				return
			} else {
				libOrch.DispatchWorkerMessage(gs.Ctx, orchestratorClient, job["id"], workerMessage.WorkerId, workerMessage.Message, "STATUS")
			}
		} else {
			libOrch.DispatchWorkerMessage(gs.Ctx, orchestratorClient, job["id"], workerMessage.WorkerId, workerMessage.Message, workerMessage.MessageType)
		}

		// Could handle these differently, but for now just dispatch them

		/*else if workerMessage.MessageType == "MARK" {
			libOrch.DispatchWorkerMessage(gs.Ctx, orchestratorClient, job["id"], workerMessage.WorkerId, workerMessage.Message, "MARK")
		} else if workerMessage.MessageType == "CONSOLE" {
			libOrch.DispatchWorkerMessage(gs.Ctx, orchestratorClient, job["id"], workerMessage.WorkerId, workerMessage.Message, "CONSOLE")
		} else if workerMessage.MessageType == "METRICS" {
			libOrch.DispatchWorkerMessage(gs.Ctx, orchestratorClient, job["id"], workerMessage.WorkerId, workerMessage.Message, "METRICS")
		} else if workerMessage.MessageType == "SUMMARY_METRICS" {
			libOrch.DispatchWorkerMessage(gs.Ctx, orchestratorClient, job["id"], workerMessage.WorkerId, workerMessage.Message, "SUMMARY_METRICS")
			workerSubscription.Close()
		} else if workerMessage.MessageType == "ERROR" {
			libOrch.DispatchWorkerMessage(gs.Ctx, orchestratorClient, job["id"], workerMessage.WorkerId, workerMessage.Message, "ERROR")
			libOrch.HandleStringError(gs.Ctx, orchestratorClient, job["id"], orchestratorId, workerMessage.Message)
			return
		} else if workerMessage.MessageType == "DEBUG" {
			libOrch.DispatchWorkerMessage(gs.Ctx, orchestratorClient, job["id"], workerMessage.WorkerId, workerMessage.Message, "DEBUG")
		}*/
	}
}

func dispatchJob(ctx context.Context, workerClient *redis.Client, job map[string]string, status string, orchestratorId string, options *libWorker.Options) error {
	// Convert options to json
	marshalledOptions, err := json.Marshal(options)
	if err != nil {
		return err
	}

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
}

/*
Check i the scope has the required credits to execute the job
*/
func checkIfHasCredits(ctx context.Context, scope map[string]string, job map[string]string) (bool, error) {
	// TODO: implement fully
	return true, nil
	/*
	   // Check max requests has not been reached
	   maxRequests := scope["maxRequests"]

	   	if maxRequests != "" {
	   		return false, fmt.Errorf("maxRequests not found")
	   	}

	   currentRequests := scope["currentRequests"]

	   	if currentRequests != "" {
	   		return false, fmt.Errorf("currentRequests not found")
	   	}

	   // TODO: More checks

	   return currentRequests < maxRequests, nil
	*/
}
