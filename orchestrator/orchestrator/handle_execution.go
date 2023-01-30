package orchestrator

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/orchestrator/orchMetrics"
	"github.com/gorilla/websocket"
)

type locatedMesaage struct {
	location string
	msg      string
}

const (
	JOB_USER_UPDATES_CHANNEL  = "jobUserUpdates"
	NO_CREDITS_ABORT_CHANNEL  = "noCreditsAbort"
	FUNC_ERROR_ABORT_CHANNEL  = "funcErrorAbort"
	OUT_OF_TIME_ABORT_CHANNEL = "outOfTimeAbort"

	maxConsoleLogs = 100
)

var otherMessageTypes = []string{"MESSAGE", "MARK", "OPTIONS", "JOB_INFO", "COLLECTION_VARIABLES", "ENVIRONMENT_VARIABLES", "LOCALHOST_FILE"}

func handleExecution(gs libOrch.BaseGlobalState, job libOrch.Job, childJobs map[string]libOrch.ChildJobDistribution) (string, error) {
	libOrch.UpdateStatus(gs, "LOADING")

	// Create a handler for aborts
	jobUserUpdatesSubscription := gs.OrchestratorClient().Subscribe(gs.Ctx(), fmt.Sprintf("jobUserUpdates:%s:%s:%s", job.Scope.Variant, job.Scope.VariantTargetId, job.Id))
	defer jobUserUpdatesSubscription.Close()

	// workerSubscriptions := make(map[string]*redis.PubSub)

	// if gs.IndependentWorkerRedisHosts() {
	// 	for location, jobDistribution := range childJobs {
	// 		if jobDistribution.ChildJobs != nil && len(jobDistribution.ChildJobs) > 0 && workerSubscriptions[location] == nil {
	// 			workerSubscriptions[location] = jobDistribution.WorkerConnection.Subscribe(gs.Ctx(), fmt.Sprintf("worker:executionUpdates:%s", job.Id))
	// 			defer workerSubscriptions[location].Close()
	// 		}
	// 	}
	// } else {
	// 	for _, jobDistribution := range childJobs {
	// 		// Only get first client as unified for whole job
	// 		workerSubscriptions[unifiedRedis] = jobDistribution.WorkerConnection.Subscribe(gs.Ctx(), fmt.Sprintf("worker:executionUpdates:%s", job.Id))
	// 		defer workerSubscriptions[unifiedRedis].Close()
	// 		break
	// 	}
	// }

	// workerChannels := make(map[string]<-chan *redis.Message)
	// for location, subscription := range workerSubscriptions {
	// 	workerChannels[location] = subscription.Channel()
	// }

	err := dispatchChildJobs(gs, &childJobs)
	if err != nil {
		fmt.Println("Error dispatching child jobs: ", err)
		return abortAndFailAll(gs, childJobs, fmt.Errorf("internal error occurred dispatching child jobs: %s", err.Error()))
	}

	defer func() {
		for _, jobDistribution := range childJobs {
			for _, childJob := range jobDistribution.ChildJobs {
				if childJob.WorkerConnection == nil {
					continue
				}

				childJob.WorkerConnection.Close()
			}
		}
	}()

	// Check if workerSubscriptions is empty
	if childJobCount(childJobs) == 0 {
		libOrch.DispatchMessage(gs, "No child jobs were created", "MESSAGE")
		return "SUCCESS", nil
	}

	// Create channels for all the functiosn

	unifiedChannel := make(chan locatedMesaage)

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
			msg:      "",
		}
	}()

	// Handle aborts on functionChannels
	// for _, dispatchChildJobs := range dispatchedChildJobs {
	// 	go func(channel <-chan libOrch.FunctionResult) {
	// 		for msg := range channel {
	// 			if msg.Error != nil {
	// 				fmt.Println("Error executing function: ", msg.Error)

	// 				unifiedChannel <- locatedMesaage{
	// 					location: FUNC_ERROR_ABORT_CHANNEL,
	// 					msg:      nil,
	// 				}

	// 				fmt.Println(msg.Error)
	// 			} else if msg.Response.StatusCode != 200 {
	// 				// Read the body
	// 				body, err := ioutil.ReadAll(msg.Response.Body)
	// 				if err != nil {
	// 					fmt.Println("Error reading body: ", err)
	// 				}

	// 				fmt.Println("Error executing function: ", string(body), msg.Response.StatusCode)

	// 				unifiedChannel <- locatedMesaage{
	// 					location: FUNC_ERROR_ABORT_CHANNEL,
	// 					msg:      nil,
	// 				}
	// 			}
	// 		}
	// 	}(channel)
	//}

	for location, childJobDistribution := range childJobs {
		for _, childJob := range childJobDistribution.ChildJobs {
			go func(childJob *libOrch.ChildJob, location string) {
				for {
					childJob.ConnReadMutex.Lock()
					messageKind, p, err := childJob.WorkerConnection.ReadMessage()
					childJob.ConnReadMutex.Unlock()

					if messageKind == websocket.CloseMessage {
						return
					}

					if err != nil {
						fmt.Println("Error reading message: ", err)

						// If websocket is closed, return
						if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseAbnormalClosure) {
							return
						}

						// See if error message contains "use of closed network connection"
						if strings.Contains(err.Error(), "use of closed network connection") {
							return
						}

						continue
					}

					newMessage := locatedMesaage{
						location: location,
						msg:      string(p),
					}

					unifiedChannel <- newMessage
				}
			}(childJob, location)
		}
	}

	// Every second check credits
	go func() {
		ticker := time.NewTicker(1 * time.Second)

		for range ticker.C {
			credits := gs.CreditsManager().GetCredits()

			if credits <= 0 {
				unifiedChannel <- locatedMesaage{
					location: NO_CREDITS_ABORT_CHANNEL,
					msg:      "",
				}
			}
		}
	}()

	// Add jobUserUpdatesSubscription to the unifiedChannel
	go func() {
		for msg := range jobUserUpdatesSubscription.Channel() {
			unifiedChannel <- locatedMesaage{
				location: JOB_USER_UPDATES_CHANNEL,
				msg:      msg.String(),
			}
		}
	}()

	childJobCount := 0
	for _, jobDistribution := range childJobs {
		childJobCount += len(jobDistribution.ChildJobs)
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
			for _, childJob := range jobDistribution.ChildJobs {
				fmt.Println("Publishing start time to", fmt.Sprintf("%s:go", childJob.ChildJobId))

				eventMessage := lib.EventMessage{
					Variant: lib.GO_MESSAGE_TYPE,
					Data:    startTime.Format(time.RFC3339),
				}

				marshalledEventMessage, err := json.Marshal(eventMessage)
				if err != nil {
					fmt.Printf("error marshalling start time to %s: %s\n", childJob.ChildJobId, err)
					continue
				}

				childJob.ConnWriteMutex.Lock()
				err = childJob.WorkerConnection.WriteMessage(websocket.TextMessage, marshalledEventMessage)
				childJob.ConnWriteMutex.Unlock()
				if err != nil {
					fmt.Printf("error publishing start time to %s: %s\n", childJob.ChildJobId, err)
				}
			}
		}

		libOrch.UpdateStatus(gs, "RUNNING")
	}

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

			err := json.Unmarshal([]byte(locatedMessage.msg), &jobUserUpdate)
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
		err := json.Unmarshal([]byte(locatedMessage.msg), &workerMessage)
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
	for _, childJob := range (childJobs)[location].ChildJobs {
		if childJob.ChildJobId == childJobId {
			return childJob
		}
	}

	return nil
}

func childJobCount(childJobs map[string]libOrch.ChildJobDistribution) int {
	count := 0

	for _, childJobDistribution := range childJobs {
		count += len(childJobDistribution.ChildJobs)
	}

	return count
}
