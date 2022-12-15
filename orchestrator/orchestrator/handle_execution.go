package orchestrator

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/orchestrator/orchMetrics"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/redis/v9"
)

type locatedMesaage struct {
	location string
	msg      *redis.Message
}

const (
	JOB_USER_UPDATES_CHANNEL = "jobUserUpdates"
	NO_CREDITS_ABORT_CHANNEL = "noCreditsAbort"
)

func handleExecution(gs libOrch.BaseGlobalState, options *libWorker.Options, scope libOrch.Scope, childJobs map[string]jobDistribution, jobId string) (string, error) {
	libOrch.UpdateStatus(gs, "LOADING")

	// Create a handler for aborts
	jobUserUpdatesSubscription := gs.OrchestratorClient().Subscribe(gs.Ctx(), fmt.Sprintf("jobUserUpdates:%s:%s:%s", scope.Variant, scope.VariantTargetId, jobId))
	defer jobUserUpdatesSubscription.Close()

	workerSubscriptions := make(map[string]*redis.PubSub)
	for location, jobDistribution := range childJobs {
		if jobDistribution.Jobs != nil && len(jobDistribution.Jobs) > 0 {
			workerSubscriptions[location] = jobDistribution.workerClient.Subscribe(gs.Ctx(), fmt.Sprintf("worker:executionUpdates:%s", jobId))
			defer workerSubscriptions[location].Close()
		}
	}

	workerChannels := make(map[string]<-chan *redis.Message)
	for location, subscription := range workerSubscriptions {
		workerChannels[location] = subscription.Channel()
	}

	// Check if workerSubscriptions is empty
	if len(workerSubscriptions) == 0 {
		libOrch.DispatchMessage(gs, "No child jobs were created", "MESSAGE")
		return "SUCCESS", nil
	}

	if err := dispatchChildJobs(gs, options, childJobs); err != nil {
		return abortAndFailAll(gs, childJobs, err)
	}

	unifiedChannel := make(chan locatedMesaage)

	for location, channel := range workerChannels {
		// Variable declaration here for locationLoop seems to be required to avoid
		// capturing the loop variable error
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

	// Every second check credits
	go func() {
		ticker := time.NewTicker(1 * time.Second)

		for range ticker.C {
			credits := gs.CreditsManager().GetCredits()

			if credits <= 0 {
				unifiedChannel <- locatedMesaage{
					location: NO_CREDITS_ABORT_CHANNEL,
					msg:      nil,
				}
			}
		}
	}()

	// Add jobUserUpdatesSubscription to the unifiedChannel
	go func() {
		for msg := range jobUserUpdatesSubscription.Channel() {
			unifiedChannel <- locatedMesaage{
				location: JOB_USER_UPDATES_CHANNEL,
				msg:      msg,
			}
		}
	}()

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
		if locatedMessage.location == NO_CREDITS_ABORT_CHANNEL {
			fmt.Println("Aborting job due to no credits")

			err := abortChildJobs(gs, childJobs)
			return abortAndFailAll(gs, childJobs, err)
		}

		// Handle user updates separately
		if locatedMessage.location == JOB_USER_UPDATES_CHANNEL {
			var JobUserUpdate = lib.JobUserUpdate{}

			err := json.Unmarshal([]byte(locatedMessage.msg.Payload), &JobUserUpdate)
			if err != nil {
				return abortAndFailAll(gs, childJobs, err)
			}

			if JobUserUpdate.UpdateType == "CANCEL" {
				fmt.Println("Aborting job due to user request")

				// Cancel all child jobs
				err := abortChildJobs(gs, childJobs)
				if err != nil {
					libOrch.HandleError(gs, err)
				}

				return abortAndFailAll(gs, childJobs, errors.New("job cancelled by user"))
			}

			// jobUserUpdates are different to other messages, so we can skip the rest of the loop
			continue
		}

		var workerMessage = libOrch.WorkerMessage{}
		err := json.Unmarshal([]byte(locatedMessage.msg.Payload), &workerMessage)
		if err != nil {
			return abortAndFailAll(gs, childJobs, err)
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
						// If one job fails, cancel all other jobs
						return abortAndFailAll(gs, childJobs, nil)
					} else {
						libOrch.UpdateStatus(gs, "SUCCESS")
						return "SUCCESS", nil
					}
				}
			}

			// Sometimes errors don't stop the execution automatically so stop them here
		} else if workerMessage.MessageType == "ERROR" {
			libOrch.DispatchWorkerMessage(gs, workerMessage.WorkerId, workerMessage.ChildJobId, workerMessage.Message, workerMessage.MessageType)

			resolutionMutex.Lock()
			failureCount++
			resolutionMutex.Unlock()

			if successCount+failureCount == childJobCount {
				// All jobs have finished
				return abortAndFailAll(gs, childJobs, nil)
			}
		} else if workerMessage.MessageType == "METRICS" {
			childJob := findChildJob(childJobs, locatedMessage.location, workerMessage.ChildJobId)
			if childJob == nil {
				return abortAndFailAll(gs, childJobs, fmt.Errorf("could not find child job with id %s to add summary metrics to", workerMessage.ChildJobId))
			}

			(*gs.MetricsStore()).AddMessage(workerMessage, locatedMessage.location, childJob.SubFraction)
		} else if workerMessage.MessageType == "SUMMARY_METRICS" {
			childJob := findChildJob(childJobs, locatedMessage.location, workerMessage.ChildJobId)
			if childJob == nil {
				return abortAndFailAll(gs, childJobs, fmt.Errorf("could not find child job with id %s to add summary metrics to", workerMessage.ChildJobId))
			}

			summaryBank.AddMessage(workerMessage, locatedMessage.location, childJob.SubFraction)

			if summaryBank.Size() == childJobCount {
				err := summaryBank.CalculateAndDispatchSummaryMetrics()
				if err != nil {
					return abortAndFailAll(gs, childJobs, err)
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
	return abortAndFailAll(gs, childJobs, errors.New("an unexpected error occurred"))
}

func findChildJob(childJobs map[string]jobDistribution, location string, childJobId string) *libOrch.ChildJob {
	for _, childJob := range (childJobs)[location].Jobs {
		if childJob.ChildJobId == childJobId {
			return &childJob
		}
	}

	return nil
}
