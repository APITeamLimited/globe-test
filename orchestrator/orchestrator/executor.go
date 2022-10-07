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

func runExecution(gs libOrch.BaseGlobalState, options *libWorker.Options, scope map[string]string, childJobs map[string]jobDistribution, jobId string) (string, error) {
	// TODO: implement credit check

	// Check if has credits
	/*hasCredits, err := checkIfHasCredits(gs.Ctx(), scope, job)
	if err != nil {
		return "", err
	}
	if !hasCredits {
		libOrch.UpdateStatus(gs.Ctx(), gs.Client(), gs.JobId(), gs.OrchestratorId(), "NO_CREDITS")
		return "", errors.New("not enough credits to execute that job")
	}*/

	workerSubscriptions := make(map[string]*redis.PubSub)
	for location, jobDistribution := range childJobs {
		if jobDistribution.jobs != nil && len(*jobDistribution.jobs) > 0 {
			workerSubscriptions[location] = jobDistribution.workerClient.Subscribe(gs.Ctx(), fmt.Sprintf("worker:executionUpdates:%s", jobId))
		}
	}

	workerChannels := make(map[string]<-chan *redis.Message)
	for location, subscription := range workerSubscriptions {
		workerChannels[location] = subscription.Channel()
	}

	// Update the status
	libOrch.UpdateStatus(gs.Ctx(), gs.Client(), jobId, gs.OrchestratorId(), "RUNNING")

	// Check if workerSubscriptions is empty
	if len(workerSubscriptions) == 0 {
		libOrch.DispatchMessage(gs.Ctx(), gs.Client(), gs.JobId(), gs.OrchestratorId(), "No child jobs were created", "INFO")
		return "SUCCESS", nil
	}

	for _, jobDistribution := range childJobs {
		for _, job := range *jobDistribution.jobs {
			err := dispatchJob(gs, jobDistribution.workerClient, job, options)
			if err != nil {
				return "", err
			}
		}
	}

	unifiedChannel := make(chan *redis.Message)

	for _, channel := range workerChannels {
		go func(channel <-chan *redis.Message) {
			for msg := range channel {
				unifiedChannel <- msg
			}
		}(channel)
	}

	for msg := range unifiedChannel {
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
			libOrch.DispatchWorkerMessage(gs.Ctx, orchestratorClient, job.Id, workerMessage.WorkerId, workerMessage.Message, "MARK")
		} else if workerMessage.MessageType == "CONSOLE" {
			libOrch.DispatchWorkerMessage(gs.Ctx, orchestratorClient, job.Id, workerMessage.WorkerId, workerMessage.Message, "CONSOLE")
		} else if workerMessage.MessageType == "METRICS" {
			libOrch.DispatchWorkerMessage(gs.Ctx, orchestratorClient, job.Id, workerMessage.WorkerId, workerMessage.Message, "METRICS")
		} else if workerMessage.MessageType == "SUMMARY_METRICS" {
			libOrch.DispatchWorkerMessage(gs.Ctx, orchestratorClient, job.Id, workerMessage.WorkerId, workerMessage.Message, "SUMMARY_METRICS")
			workerSubscription.Close()
		} else if workerMessage.MessageType == "ERROR" {
			libOrch.DispatchWorkerMessage(gs.Ctx, orchestratorClient, job.Id, workerMessage.WorkerId, workerMessage.Message, "ERROR")
			libOrch.HandleStringError(gs.Ctx, orchestratorClient, job.Id, gs.orchestratorId, workerMessage.Message)
			return
		} else if workerMessage.MessageType == "DEBUG" {
			libOrch.DispatchWorkerMessage(gs.Ctx, orchestratorClient, job.Id, workerMessage.WorkerId, workerMessage.Message, "DEBUG")
		}*/
	}

	// Shouldn't get here
	return "", errors.New("an unexpected error occurred")
}

func dispatchJob(gs libOrch.BaseGlobalState, workerClient *redis.Client, job libOrch.ChildJob, options *libWorker.Options) error {
	// Convert options to json
	marshalledChildJob, err := json.Marshal(job)
	if err != nil {
		return err
	}

	workerClient.HSet(gs.Ctx(), job.ChildJobId, "job", marshalledChildJob)

	workerClient.Publish(gs.Ctx(), "worker:execution", job.ChildJobId)
	workerClient.SAdd(gs.Ctx(), "worker:executionHistory", job.ChildJobId)

	return nil
}

/*
Check i the scope has the required credits to execute the job
*/
func checkIfHasCredits(ctx context.Context, scope map[string]string, job libOrch.Job) (bool, error) {
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
