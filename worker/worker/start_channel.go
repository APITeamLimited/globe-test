package worker

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
)

const (
	WAITING = iota
	STARTED
	FAILED
)

// Starts test on command from orchestrator or cancels test if not received start
// command within timeout of 1 minute
func testStartChannel(gs libWorker.BaseGlobalState, eventChannels eventChannels) chan *time.Time {
	startChan := make(chan *time.Time)

	status := WAITING
	statusMutex := &sync.Mutex{}

	// Listen for start command from orchestrator
	go func() {
		message := <-eventChannels.goMessageChannel

		statusMutex.Lock()
		defer statusMutex.Unlock()

		if status != WAITING {
			return
		}

		if message == "" {
			startChan <- nil
			status = FAILED
			return
		}

		// Parse start time from message

		startTime, err := time.Parse(time.RFC3339, message)
		if err != nil {
			startChan <- nil
			fmt.Println("Error parsing start time from message", err)
		}

		// Send start command to test runner
		if status == WAITING {
			startChan <- &startTime
			status = STARTED
		}
	}()

	// Race against timeout of 1 minute
	go func() {
		time.Sleep(1 * time.Minute)

		statusMutex.Lock()
		defer statusMutex.Unlock()

		if status == WAITING {
			startChan <- nil
			libWorker.HandleStringError(gs, "failed to start test, failed to receive start signal from orchestrator after 1 minute")
			status = FAILED
		}
	}()

	// Listen for abort command from orchestrator
	go func() {
		for msg := range eventChannels.childUpdatesChannel {
			var updateMessage = JobUserUpdate{}
			if err := json.Unmarshal([]byte(msg), &updateMessage); err != nil {
				libWorker.HandleStringError(gs, fmt.Sprintf("Error unmarshalling abort message: %s", err.Error()))
				continue
			}

			if updateMessage.UpdateType == "CANCEL" {
				fmt.Println("Aborting child job due to a request from the orchestrator")

				statusMutex.Lock()

				if status == WAITING {
					startChan <- nil

					statusMutex.Unlock()
					status = FAILED
				} else {
					statusMutex.Unlock()
					return
				}
			}
		}
	}()

	return startChan
}
