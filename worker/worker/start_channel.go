package worker

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/redis/v9"
)

const (
	WAITING = iota
	STARTED
	FAILED
)

// Starts test on command from orchestrator or cancels test if not received start
// command within timeout of 1 minute
func testStartChannel(workerInfo *libWorker.WorkerInfo, startSubChannel <-chan *redis.Message, gs libWorker.BaseGlobalState) chan *time.Time {
	startChan := make(chan *time.Time)

	status := WAITING
	statusMutex := &sync.Mutex{}

	// Listen for start command from orchestrator
	go func() {
		message := <-startSubChannel

		statusMutex.Lock()
		defer statusMutex.Unlock()

		if status != WAITING {
			return
		}

		if message == nil || message.Payload == "" {
			startChan <- nil
			status = FAILED
			return
		}

		// Parse start time from message

		startTime, err := time.Parse(time.RFC3339, message.Payload)
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

	// Sometimes the event is missed, so poll the set value
	go func() {
		for {
			time.Sleep(100 * time.Millisecond)

			statusMutex.Lock()
			statusValue := status
			statusMutex.Unlock()
			if statusValue != WAITING {
				return
			}

			setKey := fmt.Sprintf("%s:go:set", workerInfo.ChildJobId)

			startTime, err := workerInfo.Client.Get(workerInfo.Ctx, setKey).Result()
			if err != nil {
				if err != redis.Nil {
					fmt.Println("Error getting start time from set", err)
					return
				}

				// Set not set yet, continue polling
				continue
			}

			startTimeParsed, err := time.Parse(time.RFC3339, startTime)
			if err != nil {
				fmt.Println("Error parsing start time from set", err)
				return
			}

			statusMutex.Lock()
			statusValue = status
			statusMutex.Unlock()

			if statusValue == WAITING {
				startChan <- &startTimeParsed
				status = STARTED
			}

			return
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
		cancelKey := fmt.Sprintf("childjobUserUpdates:%s", workerInfo.ChildJobId)
		cancelSubscription := workerInfo.Client.Subscribe(workerInfo.Ctx, cancelKey)
		cancelChannel := cancelSubscription.Channel()
		defer cancelSubscription.Close()

		for msg := range cancelChannel {
			var updateMessage = JobUserUpdate{}
			if err := json.Unmarshal([]byte(msg.Payload), &updateMessage); err != nil {
				libWorker.HandleStringError(*workerInfo.Gs, fmt.Sprintf("Error unmarshalling abort message: %s", err.Error()))
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
