package orchestrator

import (
	"strings"
	"time"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/redis/v9"
	"github.com/gorilla/websocket"
)

// Listens for status updates and can automatically abort the job
func listenForOrchestratorErrors(gs libOrch.BaseGlobalState, unifiedChannel chan locatedMesaage) {
	for message := range gs.StatusUpdatesChannel() {
		if message == "FAILED" || message == "SUCCESS" {
			unifiedChannel <- locatedMesaage{
				location: OTHER_FAILURE_ABORT_CHANNEL,
				msg:      "",
			}

			return
		}
	}
}

// Automatically aborts the job if it takes too long
func abortIfMaxDurationExceeded(gs libOrch.BaseGlobalState, job libOrch.Job, unifiedChannel chan locatedMesaage) {
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
}

// Checks for credits every second
func checkCreditsPeriodically(gs libOrch.BaseGlobalState, unifiedChannel chan locatedMesaage) {
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
}

// Listens for messages from child jobs
func listenForChildJobMessages(gs libOrch.BaseGlobalState, childJobs map[string]libOrch.ChildJobDistribution, unifiedChannel chan locatedMesaage) {
	for location, childJobDistribution := range childJobs {
		for _, childJob := range childJobDistribution.ChildJobs {
			go func(childJob *libOrch.ChildJob, location string) {
				for {
					childJob.ConnReadMutex.Lock()
					messageKind, p, err := childJob.WorkerConnection.ReadMessage()

					if messageKind == websocket.CloseMessage {
						childJob.ConnReadMutex.Unlock()
						return
					}

					if err != nil {
						if strings.Contains(err.Error(), "use of closed network connection") || websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseAbnormalClosure) {
							childJob.ConnReadMutex.Unlock()

							newMessage := locatedMesaage{
								location: FUNC_ERROR_ABORT_CHANNEL,
								msg:      "",
							}

							unifiedChannel <- newMessage

							return
						}

						childJob.ConnReadMutex.Unlock()
						continue
					}

					childJob.ConnReadMutex.Unlock()

					newMessage := locatedMesaage{
						location: location,
						msg:      string(p),
					}

					unifiedChannel <- newMessage
				}
			}(childJob, location)
		}
	}
}

// Listens for job updates from the end user
func listenForJobUserUpdates(gs libOrch.BaseGlobalState, jobUserUpdatesSubscription *redis.PubSub, unifiedChannel chan locatedMesaage) {
	for msg := range jobUserUpdatesSubscription.Channel() {
		unifiedChannel <- locatedMesaage{
			location: JOB_USER_UPDATES_CHANNEL,
			msg:      msg.Payload,
		}
	}
}
