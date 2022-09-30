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

func runExecution(gs libOrch.BaseGlobalState, options *libWorker.Options, scope map[string]string, workerClients map[string]*redis.Client, job map[string]string) (string, error) {
	// Check if has credits
	hasCredits, err := checkIfHasCredits(gs.Ctx(), scope, job)
	if err != nil {
		return "", err
	}
	if !hasCredits {
		libOrch.UpdateStatus(gs.Ctx(), gs.Client(), gs.JobId(), gs.OrchestratorId(), "NO_CREDITS")
		return "", errors.New("not enough credits to execute that job")
	}

	workerClient := workerClients["portsmouth"]
	workerSubscription := workerClient.Subscribe(gs.Ctx(), fmt.Sprintf("worker:executionUpdates:%s", job["id"]))

	//marshalledOptions, err := json.Marshal(options)
	if err != nil {
		return "", err
	}

	// Update the status
	libOrch.UpdateStatus(gs.Ctx(), gs.Client(), job["id"], gs.OrchestratorId(), "RUNNING")

	err = dispatchJob(gs, workerClient, job, "PENDING", options)
	if err != nil {
		return "", err
	}

	if err != nil {
		return "", err
	}

	for msg := range workerSubscription.Channel() {
		workerMessage := libOrch.WorkerMessage{}
		err := json.Unmarshal([]byte(msg.Payload), &workerMessage)
		if err != nil {
			return "", err
		}

		if workerMessage.MessageType == "STATUS" {
			if workerMessage.Message == "FAILED" {
				return "FAILURE", nil
			} else if workerMessage.Message == "SUCCESS" {
				return "SUCCESS", nil
			} else {
				libOrch.DispatchWorkerMessage(gs.Ctx(), gs.Client(), gs.JobId(), workerMessage.WorkerId, workerMessage.Message, "STATUS")
			}
		} else if workerMessage.MessageType == "METRICS" {
			(*gs.MetricsStore()).AddMessage(workerMessage, "portsmouth")
		} else if workerMessage.MessageType == "DEBUG" {
			// TODO: make this configurable
			continue
		} else {
			libOrch.DispatchWorkerMessage(gs.Ctx(), gs.Client(), gs.JobId(), workerMessage.WorkerId, workerMessage.Message, workerMessage.MessageType)
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
			libOrch.HandleStringError(gs.Ctx, orchestratorClient, job["id"], gs.orchestratorId, workerMessage.Message)
			return
		} else if workerMessage.MessageType == "DEBUG" {
			libOrch.DispatchWorkerMessage(gs.Ctx, orchestratorClient, job["id"], workerMessage.WorkerId, workerMessage.Message, "DEBUG")
		}*/
	}

	// Shouldn't get here
	return "", errors.New("an unexpected error occurred")
}

func dispatchJob(gs libOrch.BaseGlobalState, workerClient *redis.Client, job map[string]string, status string, options *libWorker.Options) error {
	// Convert options to json
	marshalledOptions, err := json.Marshal(options)
	if err != nil {
		return err
	}

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
