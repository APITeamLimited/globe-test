package orchestrator

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/redis/v9"
)

type locatedMesaage struct {
	location string
	msg      *redis.Message
}

func runExecution(gs libOrch.BaseGlobalState, options *libWorker.Options, scope libOrch.Scope, childJobs map[string]jobDistribution, jobId string) (string, error) {
	libOrch.UpdateStatus(gs.Ctx(), gs.Client(), jobId, gs.OrchestratorId(), "LOADING")

	workerSubscriptions := make(map[string]*redis.PubSub)
	for location, jobDistribution := range childJobs {
		if jobDistribution.jobs != nil && len(jobDistribution.jobs) > 0 {
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
		for _, job := range jobDistribution.jobs {
			err := dispatchJob(gs, jobDistribution.workerClient, job, options)
			if err != nil {
				return "", err
			}
		}
	}

	unifiedChannel := make(chan locatedMesaage)

	for location, channel := range workerChannels {
		go func(channel <-chan *redis.Message) {
			for msg := range channel {
				unifiedChannel <- locatedMesaage{
					location,
					msg,
				}
			}
		}(channel)
	}

	chilJobCount := len(childJobs)
	jobsInitialised := 0

	for locatedMessage := range unifiedChannel {
		workerMessage := libOrch.WorkerMessage{}
		err := json.Unmarshal([]byte(locatedMessage.msg.Payload), &workerMessage)
		if err != nil {
			return "", err
		}

		if workerMessage.MessageType == "STATUS" {
			if workerMessage.Message == "READY" {
				jobsInitialised++
				if jobsInitialised == chilJobCount {
					// Broadcast the start message to all child jobs

					for _, jobDistribution := range childJobs {
						for _, job := range jobDistribution.jobs {
							jobDistribution.workerClient.Publish(gs.Ctx(), fmt.Sprintf("%s:go", job.ChildJobId), "GO TIME")
						}
					}

				}

				libOrch.UpdateStatus(gs.Ctx(), gs.Client(), jobId, gs.OrchestratorId(), "RUNNING")
			} else if workerMessage.Message == "FAILURE" {
				return "FAILURE", nil
			} else if workerMessage.Message == "SUCCESS" {
				return "SUCCESS", nil
			}
			// Ignore other kinds of status messages
			//libOrch.DispatchWorkerMessage(gs.Ctx(), gs.Client(), gs.JobId(), workerMessage.WorkerId, workerMessage.Message, "STATUS")

			// Sometimes errors don't stop the execution automatically so stop them here
		} else if workerMessage.MessageType == "ERROR" {
			libOrch.DispatchWorkerMessage(gs.Ctx(), gs.Client(), gs.JobId(), workerMessage.WorkerId, workerMessage.Message, workerMessage.MessageType)
			return "FAILURE", nil
		} else if workerMessage.MessageType == "METRICS" {
			(*gs.MetricsStore()).AddMessage(workerMessage, locatedMessage.location)
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

	workerClient.SAdd(gs.Ctx(), "worker:executionHistory", job.ChildJobId)
	workerClient.Publish(gs.Ctx(), "worker:execution", job.ChildJobId)

	return nil
}
