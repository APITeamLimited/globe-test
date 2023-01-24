package worker

import (
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
func testStartChannel(workerInfo *libWorker.WorkerInfo, startSubChannel <-chan *redis.Message) chan *time.Time {
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
			status = FAILED
		}
	}()

	return startChan
}
