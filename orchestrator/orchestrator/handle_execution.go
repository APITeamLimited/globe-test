package orchestrator

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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

	maxConsoleLogs = 100
)

var otherMessageTypes = []string{"MESSAGE", "MARK", "OPTIONS", "JOB_INFO", "COLLECTION_VARIABLES", "ENVIRONMENT_VARIABLES", "LOCALHOST_FILE"}

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
		fmt.Println("Error dispatching child jobs: ", err)
		return abortAndFailAll(gs, childJobs, err)
	}

	// race main thread with the context cancellation from job.maxTestDurationMinutes
	go func() {
		if job.MaxTestDurationMinutes == 0 {
			return
		}

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
					fmt.Println("Error executing function: ", msg.Error)

					unifiedChannel <- locatedMesaage{
						location: FUNC_ERROR_ABORT_CHANNEL,
						msg:      nil,
					}

					fmt.Println(msg.Error)
				} else if msg.Response.StatusCode != 200 {
					// Read the body
					body, err := ioutil.ReadAll(msg.Response.Body)
					if err != nil {
						fmt.Println("Error reading body: ", err)
					}

					fmt.Println("Error executing function: ", string(body), msg.Response.StatusCode)

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

	consoleLogCount := 0
	sentMaxLogsReached := false
	consoleLogCountMutex := sync.Mutex{}

	summaryBank := orchMetrics.NewSummaryBank(gs, job.Options)
	defer summaryBank.Cleanup()

	checkIfCanStart := func() {
		jobsMutex.Lock()
		defer jobsMutex.Unlock()

		if len(jobsInitialised) != childJobCount {
			return
		}

		// Broadcast the start message to all child jobs

		var startTime time.Time

		if childJobCount == 1 {
			startTime = time.Now()
		} else {
			// Set one second in the future to allow all workers to start
			startTime = time.Now().Add(time.Second)
		}

		fmt.Printf("Broadcasting start time %s to all child jobs", startTime.Format(time.RFC3339))

		if gs.GetStatus() == "FAILURE" {
			// If the job has been cancelled, don't start it
			return
		}

		for _, jobDistribution := range childJobs {
			for _, job := range jobDistribution.Jobs {
				fmt.Println("Publishing start time to", fmt.Sprintf("%s:go", job.ChildJobId))
				go jobDistribution.WorkerClient.Set(gs.Ctx(), fmt.Sprintf("%s:go:set", job.ChildJobId), startTime.Format(time.RFC3339), time.Minute)
				go jobDistribution.WorkerClient.Publish(gs.Ctx(), fmt.Sprintf("%s:go", job.ChildJobId), startTime.Format(time.RFC3339))
			}
		}

		libOrch.UpdateStatus(gs, "RUNNING")
	}

	// Somtimes start message appears to be missed, so periodically poll for it
	go func() {
		for {
			time.Sleep(200 * time.Millisecond)

			if gs.GetStatus() != "LOADING" {
				return
			}

			// Check statuses of all child jobs
			for _, jobDistribution := range childJobs {
				for _, childJob := range jobDistribution.Jobs {
					status, err := jobDistribution.WorkerClient.HGet(gs.Ctx(), childJob.ChildJobId, "status").Result()
					if err != nil {
						// Check for nil
						if err != redis.Nil {
							fmt.Println("Error getting status", err)
						}

						continue
					}

					if status == "READY" {
						jobsMutex.Lock()

						alreadyInitialised := false

						// Check if we've already initialised this job
						for _, initialisedJobId := range jobsInitialised {
							if initialisedJobId == childJob.ChildJobId {
								alreadyInitialised = true
								break
							}
						}

						if alreadyInitialised {
							jobsMutex.Unlock()
							continue
						}

						jobsInitialised = append(jobsInitialised, childJob.ChildJobId)
						jobsMutex.Unlock()

						checkIfCanStart()
					}
				}
			}
		}
	}()

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
				fmt.Println("Error unmarshalling jobUserUpdate", err)
				continue
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
			fmt.Println("Error unmarshalling workerMessage", err)
			continue
		}

		if workerMessage.MessageType == "STATUS" {
			fmt.Println("Received status message from", workerMessage.WorkerId, workerMessage.ChildJobId, workerMessage.Message)

			gs.SetChildJobState(workerMessage.WorkerId, workerMessage.ChildJobId, workerMessage.Message)

			if workerMessage.Message == "READY" {
				// Check if the job has already been initialised

				alreadyInitialised := false

				jobsMutex.Lock()
				for _, initialisedJob := range jobsInitialised {
					if initialisedJob == workerMessage.ChildJobId {
						alreadyInitialised = true
						break
					}
				}

				fmt.Println("Job", workerMessage.ChildJobId, "is ready", "alreadyInitialised:", alreadyInitialised, "new jobs count:", len(jobsInitialised)+1, "childJobCount:", childJobCount)

				jobsMutex.Unlock()

				if !alreadyInitialised {
					jobsMutex.Lock()
					jobsInitialised = append(jobsInitialised, workerMessage.ChildJobId)
					jobsMutex.Unlock()
					checkIfCanStart()
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
		} else if workerMessage.MessageType == "CONSOLE" {
			consoleLogCountMutex.Lock()

			if consoleLogCount >= maxConsoleLogs {
				if !sentMaxLogsReached {
					sentMaxLogsReached = true
					consoleLogCountMutex.Unlock()

					libOrch.DispatchWorkerMessage(gs, workerMessage.WorkerId, workerMessage.ChildJobId, "MAX_CONSOLE_LOGS_REACHED", "MESSAGE")
				} else {
					consoleLogCountMutex.Unlock()
				}

				continue
			}

			consoleLogCount++
			consoleLogCountMutex.Unlock()

			libOrch.DispatchWorkerMessage(gs, workerMessage.WorkerId, workerMessage.ChildJobId, workerMessage.Message, workerMessage.MessageType)
		} else {
			// Check if the message is a custom message
			for _, messageType := range otherMessageTypes {
				if messageType == workerMessage.MessageType {

					libOrch.DispatchWorkerMessage(gs, workerMessage.WorkerId, workerMessage.ChildJobId, workerMessage.Message, workerMessage.MessageType)
					break
				}
			}

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
