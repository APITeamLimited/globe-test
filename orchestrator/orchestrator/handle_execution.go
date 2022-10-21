package orchestrator

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/orchestrator/orchMetrics"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/redis/v9"
)

type locatedMesaage struct {
	location string
	msg      *redis.Message
}

func handleExecution(gs libOrch.BaseGlobalState, options *libWorker.Options, scope libOrch.Scope, childJobs map[string]jobDistribution, jobId string) (string, error) {
	libOrch.UpdateStatus(gs, "LOADING")

	workerSubscriptions := make(map[string]*redis.PubSub)
	for location, jobDistribution := range childJobs {
		if jobDistribution.Jobs != nil && len(jobDistribution.Jobs) > 0 {
			workerSubscriptions[location] = jobDistribution.workerClient.Subscribe(gs.Ctx(), fmt.Sprintf("worker:executionUpdates:%s", jobId))
		}
	}

	workerChannels := make(map[string]<-chan *redis.Message)
	for location, subscription := range workerSubscriptions {
		workerChannels[location] = subscription.Channel()
	}

	// Update the status
	libOrch.UpdateStatus(gs, "RUNNING")

	// Check if workerSubscriptions is empty
	if len(workerSubscriptions) == 0 {
		libOrch.DispatchMessage(gs, "No child jobs were created", "INFO")
		return "SUCCESS", nil
	}

	for _, jobDistribution := range childJobs {
		for _, job := range jobDistribution.Jobs {
			err := dispatchJob(gs, jobDistribution.workerClient, job, options)
			if err != nil {
				return "FAILURE", err
			}
		}
	}

	unifiedChannel := make(chan locatedMesaage)

	for location, channel := range workerChannels {
		// Seems to be required to avoid capturing the loop variable
		locationLoop := location
		go func(channel <-chan *redis.Message) {
			for msg := range channel {
				newMessage := locatedMesaage{
					location: locationLoop,
					msg:      msg,
				}

				unifiedChannel <- newMessage
			}
		}(channel)
	}

	childJobCount := 0
	for _, jobDistribution := range childJobs {
		childJobCount += len(jobDistribution.Jobs)
	}

	jobsInitialised := []string{}
	jobsMutex := &sync.Mutex{}

	successCount := 0
	failureCount := 0
	resolutionMutex := sync.Mutex{}

	summaryBank := orchMetrics.NewSummaryBank(gs, options)

	for locatedMessage := range unifiedChannel {
		var workerMessage = libOrch.WorkerMessage{}
		err := json.Unmarshal([]byte(locatedMessage.msg.Payload), &workerMessage)
		if err != nil {
			return "FAILURE", err
		}

		if workerMessage.MessageType == "STATUS" {
			gs.SetChildJobState(workerMessage.WorkerId, workerMessage.ChildJobId, workerMessage.Message)

			if workerMessage.Message == "READY" {
				jobsMutex.Lock()

				alreadyInitialised := false

				// Check if the job has already been initialised
				for _, initialisedJob := range jobsInitialised {
					if initialisedJob == workerMessage.ChildJobId {
						jobsMutex.Unlock()
						alreadyInitialised = true
						break
					}
				}

				if !alreadyInitialised {
					jobsInitialised = append(jobsInitialised, workerMessage.ChildJobId)
					jobsMutex.Unlock()

					if len(jobsInitialised) == childJobCount {
						// Broadcast the start message to all child jobs
						for _, jobDistribution := range childJobs {
							for _, job := range jobDistribution.Jobs {
								jobDistribution.workerClient.Publish(gs.Ctx(), fmt.Sprintf("%s:go", job.ChildJobId), "GO TIME")
							}
						}

						libOrch.UpdateStatus(gs, "RUNNING")
					}
				}
			} else if workerMessage.Message == "SUCCESS" || workerMessage.Message == "FAILURE" {
				resolutionMutex.Lock()
				if workerMessage.Message == "SUCCESS" {
					successCount++
				} else {
					failureCount++
				}
				resolutionMutex.Unlock()

				if successCount+failureCount == childJobCount {
					// All jobs have finished
					if failureCount > 0 {
						libOrch.UpdateStatus(gs, "FAILURE")
						return "FAILURE", nil
					} else {
						libOrch.UpdateStatus(gs, "SUCCESS")
						return "SUCCESS", nil
					}
				}
			}
			// Ignore other kinds of status messages
			//libOrch.DispatchWorkerMessage(gs.Ctx(), gs.Client(), gs.JobId(), workerMessage.WorkerId, workerMessage.Message, "STATUS")

			// Sometimes errors don't stop the execution automatically so stop them here
		} else if workerMessage.MessageType == "ERROR" {
			libOrch.DispatchWorkerMessage(gs, workerMessage.WorkerId, workerMessage.ChildJobId, workerMessage.Message, workerMessage.MessageType)

			resolutionMutex.Lock()
			failureCount++
			resolutionMutex.Unlock()

			if successCount+failureCount == childJobCount {
				// All jobs have finished
				libOrch.UpdateStatus(gs, "FAILURE")
				return "FAILURE", nil
			}
		} else if workerMessage.MessageType == "METRICS" {
			childJob := findChildJob(&childJobs, locatedMessage.location, workerMessage.ChildJobId)
			if childJob == nil {
				return "FAILURE", fmt.Errorf("could not find child job with id %s to add summary metrics to", workerMessage.ChildJobId)
			}

			(*gs.MetricsStore()).AddMessage(workerMessage, locatedMessage.location, childJob.SubFraction)
		} else if workerMessage.MessageType == "SUMMARY_METRICS" {
			childJob := findChildJob(&childJobs, locatedMessage.location, workerMessage.ChildJobId)
			if childJob == nil {
				return "FAILURE", fmt.Errorf("could not find child job with id %s to add summary metrics to", workerMessage.ChildJobId)
			}

			summaryBank.AddMessage(workerMessage, locatedMessage.location, childJob.SubFraction)

			if summaryBank.Size() == childJobCount {
				err := summaryBank.CalculateAndDispatchSummaryMetrics()
				if err != nil {
					return "FAILURE", err
				}
			}
		} else if workerMessage.MessageType == "DEBUG" {
			// TODO: make this configurable
			continue
		} else {
			libOrch.DispatchWorkerMessage(gs, workerMessage.WorkerId, workerMessage.ChildJobId, workerMessage.Message, workerMessage.MessageType)
		}
	}

	// Should never get here
	return "FAILURE", errors.New("an unexpected error occurred")
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

func findChildJob(childJobs *map[string]jobDistribution, location string, childJobId string) *libOrch.ChildJob {
	for _, childJob := range (*childJobs)[location].Jobs {
		if childJob.ChildJobId == childJobId {
			return &childJob
		}
	}

	return nil
}