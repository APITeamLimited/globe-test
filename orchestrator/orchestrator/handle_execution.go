package orchestrator

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

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

func handleExecution(gs libOrch.BaseGlobalState, options *libWorker.Options, scope libOrch.Scope, childJobs map[string]jobDistribution, jobId string) (string, error) {
	if gs.CreditsManager().UseCredits(1) {
		libOrch.UpdateStatus(gs, "LOADING")
	} else {
		return abortAndFailAll(gs, childJobs, errors.New(lib.OUT_OF_CREDITS_MESSAGE))
	}

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
		if gs.CreditsManager().UseCredits(1) {
			libOrch.DispatchMessage(gs, "No child jobs were created", "MESSAGE")
			return "SUCCESS", nil
		} else {
			return abortAndFailAll(gs, childJobs, errors.New(lib.OUT_OF_CREDITS_MESSAGE))
		}
	}

	for _, jobDistribution := range childJobs {
		for _, job := range jobDistribution.Jobs {
			err := dispatchJob(gs, jobDistribution.workerClient, job, options)
			if err != nil {
				return abortAndFailAll(gs, childJobs, err)
			}
		}
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

	// Add jobUserUpdatesSubscription to the unifiedChannel
	go func() {
		for msg := range jobUserUpdatesSubscription.Channel() {
			newMessage := locatedMesaage{
				location: "jobUserUpdates",
				msg:      msg,
			}

			unifiedChannel <- newMessage
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
		// Handle user updates separately
		if locatedMessage.location == "jobUserUpdates" {
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

						if gs.CreditsManager().UseCredits(1) {
							libOrch.UpdateStatus(gs, "RUNNING")
						} else {
							return abortAndFailAll(gs, childJobs, errors.New(lib.OUT_OF_CREDITS_MESSAGE))
						}
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
						if gs.CreditsManager().UseCredits(1) {
							libOrch.UpdateStatus(gs, "SUCCESS")
							return "SUCCESS", nil
						} else {
							return abortAndFailAll(gs, childJobs, errors.New(lib.OUT_OF_CREDITS_MESSAGE))
						}
					}
				}
			}

			// Sometimes errors don't stop the execution automatically so stop them here
		} else if workerMessage.MessageType == "ERROR" {
			if gs.CreditsManager().UseCredits(1) {
				libOrch.DispatchWorkerMessage(gs, workerMessage.WorkerId, workerMessage.ChildJobId, workerMessage.Message, workerMessage.MessageType)
			} else {
				return abortAndFailAll(gs, childJobs, errors.New(lib.OUT_OF_CREDITS_MESSAGE))
			}

			resolutionMutex.Lock()
			failureCount++
			resolutionMutex.Unlock()

			if successCount+failureCount == childJobCount {
				// All jobs have finished
				return abortAndFailAll(gs, childJobs, nil)
			}
		} else if workerMessage.MessageType == "METRICS" {
			childJob := findChildJob(&childJobs, locatedMessage.location, workerMessage.ChildJobId)
			if childJob == nil {
				return abortAndFailAll(gs, childJobs, fmt.Errorf("could not find child job with id %s to add summary metrics to", workerMessage.ChildJobId))
			}

			(*gs.MetricsStore()).AddMessage(workerMessage, locatedMessage.location, childJob.SubFraction)
		} else if workerMessage.MessageType == "SUMMARY_METRICS" {
			childJob := findChildJob(&childJobs, locatedMessage.location, workerMessage.ChildJobId)
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
			if gs.CreditsManager().UseCredits(1) {
				libOrch.DispatchWorkerMessage(gs, workerMessage.WorkerId, workerMessage.ChildJobId, workerMessage.Message, workerMessage.MessageType)
			} else {
				return abortAndFailAll(gs, childJobs, errors.New(lib.OUT_OF_CREDITS_MESSAGE))
			}
		}
	}

	// Should never get here
	return abortAndFailAll(gs, childJobs, errors.New("an unexpected error occurred"))
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
