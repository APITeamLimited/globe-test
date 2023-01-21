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
	"github.com/APITeamLimited/redis/v9"
)

type locatedMesaage struct {
	location string
	msg      *redis.Message
}

const (
	JOB_USER_UPDATES_CHANNEL  = "jobUserUpdates"
	NO_CREDITS_ABORT_CHANNEL  = "noCreditsAbort"
	FUNC_ERROR_ABORT_CHANNEL  = "funcErrorAbort"
	OUT_OF_TIME_ABORT_CHANNEL = "outOfTimeAbort"

	unifiedRedis = "unified"
)

type childJobIdStruct struct {
	ChildJobId string `json:"childJobId"`
}

func handleExecution(gs libOrch.BaseGlobalState, job libOrch.Job, childJobs map[string]libOrch.ChildJobDistribution) (string, error) {
	libOrch.UpdateStatus(gs, "LOADING")

	// Create a handler for aborts
	jobUserUpdatesSubscription := gs.OrchestratorClient().Subscribe(gs.Ctx(), fmt.Sprintf("jobUserUpdates:%s:%s:%s", job.Scope.Variant, job.Scope.VariantTargetId, job.Id))
	defer jobUserUpdatesSubscription.Close()

	workerSubscriptions := make(map[string]*redis.PubSub)

	if gs.IndependentWorkerRedisHosts() {
		for location, jobDistribution := range childJobs {
			if jobDistribution.Jobs != nil && len(jobDistribution.Jobs) > 0 && workerSubscriptions[location] == nil {
				workerSubscriptions[location] = jobDistribution.WorkerClient.Subscribe(gs.Ctx(), fmt.Sprintf("worker:executionUpdates:%s", job.Id))
				defer workerSubscriptions[location].Close()
			}
		}
	} else {
		for _, jobDistribution := range childJobs {
			// Only get first client as unified for whole job
			workerSubscriptions[unifiedRedis] = jobDistribution.WorkerClient.Subscribe(gs.Ctx(), fmt.Sprintf("worker:executionUpdates:%s", job.Id))
			defer workerSubscriptions[unifiedRedis].Close()
			break
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

	unifiedChannel := make(chan locatedMesaage)

	functionChannels, err := dispatchChildJobs(gs, childJobs)
	if err != nil {
		return abortAndFailAll(gs, childJobs, err)
	}

	// race main thread with the context cancellation from job.maxTestDurationMinutes
	go func() {
		// Sleep for the max test duration
		time.Sleep(time.Duration(job.MaxTestDurationMinutes) * time.Minute)

		if gs.GetStatus() != "LOADING" && gs.GetStatus() != "RUNNING" {
			return
		}

		unifiedChannel <- locatedMesaage{
			location: OUT_OF_TIME_ABORT_CHANNEL,
			msg:      nil,
		}
	}()

	// Handle aborts on functionChannels
	for _, channel := range *functionChannels {
		go func(channel <-chan libOrch.FunctionResult) {
			for msg := range channel {
				if msg.Error != nil {
					unifiedChannel <- locatedMesaage{
						location: FUNC_ERROR_ABORT_CHANNEL,
						msg:      nil,
					}

					fmt.Println(msg.Error)
				} else if msg.Response.StatusCode != 200 {
					unifiedChannel <- locatedMesaage{
						location: FUNC_ERROR_ABORT_CHANNEL,
						msg:      nil,
					}
				}
			}
		}(channel)
	}

	for location, channel := range workerChannels {
		// Variable declaration here for locationLoop seems to be required to avoid
		// capturing the loop variable error

		go func(channel <-chan *redis.Message, location string) {
			childJobIdToLocation := make(map[string]string)

			if location == unifiedRedis {
				for location, jobDistribution := range childJobs {
					for _, job := range jobDistribution.Jobs {
						childJobIdToLocation[job.ChildJobId] = location
					}
				}
			}

			actualLocation := location

			for msg := range channel {
				if location == unifiedRedis {
					var childJobId childJobIdStruct
					err := json.Unmarshal([]byte(msg.Payload), &childJobId)

					if err != nil {
						fmt.Println("Error unmarshalling childJobId", err)
						continue
					}

					actualLocation = childJobIdToLocation[childJobId.ChildJobId]

					if actualLocation == "" {
						fmt.Println("Could not find location for childJobId", childJobId.ChildJobId)
						continue
					}
				}

				newMessage := locatedMesaage{
					location: actualLocation,
					msg:      msg,
				}

				unifiedChannel <- newMessage
			}
		}(channel, location)
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

	summaryBank := orchMetrics.NewSummaryBank(gs, job.Options)
	defer summaryBank.Cleanup()

	for locatedMessage := range unifiedChannel {
		if locatedMessage.location == OUT_OF_TIME_ABORT_CHANNEL {
			return abortAndFailAll(gs, childJobs, errors.New("max test duration exceeded"))
		} else if locatedMessage.location == FUNC_ERROR_ABORT_CHANNEL {
			return abortAndFailAll(gs, childJobs, errors.New("aborting job due to worker error"))
		} else if locatedMessage.location == NO_CREDITS_ABORT_CHANNEL {
			return abortAndFailAll(gs, childJobs, errors.New("aborting job due to no credits"))
		}

		// Handle user updates separately
		if locatedMessage.location == JOB_USER_UPDATES_CHANNEL {
			jobUserUpdate := lib.JobUserUpdate{}

			err := json.Unmarshal([]byte(locatedMessage.msg.Payload), &jobUserUpdate)
			if err != nil {
				return abortAndFailAll(gs, childJobs, err)
			}

			if jobUserUpdate.UpdateType == "CANCEL" {
				fmt.Println("Aborting job due to user request")

				// Cancel all child jobs
				return abortAndFailAll(gs, childJobs, errors.New("job cancelled by user"))
			}

			// jobUserUpdates are different to other messages, so we can skip the rest of the loop
			continue
		}

		workerMessage := libOrch.WorkerMessage{}
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
						alreadyInitialised = true
						break
					}
				}

				if !alreadyInitialised {
					jobsInitialised = append(jobsInitialised, workerMessage.ChildJobId)

					if len(jobsInitialised) == childJobCount {
						// Broadcast the start message to all child jobs

						var startTime time.Time

						if childJobCount == 1 {
							startTime = time.Now()
						} else {
							// Set one second in the future to allow all workers to start
							startTime = time.Now().Add(time.Second)
						}

						for _, jobDistribution := range childJobs {
							for _, job := range jobDistribution.Jobs {
								jobDistribution.WorkerClient.Publish(gs.Ctx(), fmt.Sprintf("%s:go", job.ChildJobId), startTime.Format(time.RFC3339))
							}
						}

						libOrch.UpdateStatus(gs, "RUNNING")
					}
				}

				jobsMutex.Unlock()
			} else if workerMessage.Message == "SUCCESS" || workerMessage.Message == "FAILURE" {
				fmt.Println("Job", workerMessage.ChildJobId, "finished with status", workerMessage.Message)

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
				return abortAndFailAll(gs, childJobs, fmt.Errorf("could not find child job with id %s to add metrics to", workerMessage.ChildJobId))
			}

			go (*gs.MetricsStore()).AddMessage(workerMessage, locatedMessage.location, childJob.SubFraction)
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

func findChildJob(childJobs map[string]libOrch.ChildJobDistribution, location string, childJobId string) *libOrch.ChildJob {
	for _, childJob := range (childJobs)[location].Jobs {
		if childJob.ChildJobId == childJobId {
			return &childJob
		}
	}

	return nil
}
