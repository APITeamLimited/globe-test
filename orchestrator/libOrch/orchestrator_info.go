package libOrch

import (
	"encoding/json"
	"fmt"
	"time"
)

func DispatchMessage(gs BaseGlobalState, message string, messageType string) {
	var messageStruct = OrchestratorMessage{
		JobId:          gs.JobId(),
		Time:           time.Now(),
		OrchestratorId: gs.OrchestratorId(),
		Message:        message,
		MessageType:    messageType,
	}

	messageJson, err := json.Marshal(messageStruct)
	if err != nil {
		fmt.Println("Error marshalling message")
		return
	}

	isTerminal := messageType == "STATUS" && (message == "COMPLETED_SUCCESS" || message == "COMPLETED_FAILURE")

	handleDispatchMessage(gs, func() {
		// Update main job
		gs.OrchestratorClient().SAdd(gs.Ctx(), fmt.Sprintf("%s:updates", gs.JobId()), messageJson)

		// Dispatch to channel
		gs.OrchestratorClient().Publish(gs.Ctx(), fmt.Sprintf("orchestrator:executionUpdates:%s", gs.JobId()), string(messageJson))
	}, isTerminal)
}

func DispatchMessageNoSet(gs BaseGlobalState, message string, messageType string) {
	var messageStruct = OrchestratorMessage{
		JobId:          gs.JobId(),
		Time:           time.Now(),
		OrchestratorId: gs.OrchestratorId(),
		Message:        message,
		MessageType:    messageType,
	}

	messageJson, err := json.Marshal(messageStruct)
	if err != nil {
		fmt.Println("Error marshalling message")
		return
	}

	isTerminal := messageType == "STATUS" && (message == "COMPLETED_SUCCESS" || message == "COMPLETED_FAILURE")

	handleDispatchMessage(gs, func() {
		// Dispatch to channel
		gs.OrchestratorClient().Publish(gs.Ctx(), fmt.Sprintf("orchestrator:executionUpdates:%s", gs.JobId()), messageJson)
	}, isTerminal)
}

func DispatchWorkerMessage(gs BaseGlobalState, workerId string, childJobId string, message string, messageType string) {
	var messageStruct = WorkerMessage{
		JobId:       gs.JobId(),
		ChildJobId:  childJobId,
		Time:        time.Now(),
		WorkerId:    workerId,
		Message:     message,
		MessageType: messageType,
	}

	messageJson, err := json.Marshal(messageStruct)
	if err != nil {
		fmt.Println("Error marshalling message")
		return
	}

	handleDispatchMessage(gs, func() {
		// Update main job
		gs.OrchestratorClient().SAdd(gs.Ctx(), fmt.Sprintf("%s:updates", gs.JobId()), messageJson)

		// Dispatch to channel
		gs.OrchestratorClient().Publish(gs.Ctx(), fmt.Sprintf("orchestrator:executionUpdates:%s", gs.JobId()), messageJson)
	}, false)
}

// Prevents blocking of the main thread unless isTerminal is true
func handleDispatchMessage(gs BaseGlobalState, setFunc func(), isTerminal bool) {
	go func() {
		messageQueue := gs.MessageQueue()

		if !isTerminal {
			messageQueue.Mutex.Lock()
			messageQueue.QueueCount++
			messageQueue.Mutex.Unlock()

			setFunc()

			messageQueue.Mutex.Lock()
			messageQueue.QueueCount--

			// Must unlock the mutex before sending the new count to the channel
			messageQueue.Mutex.Unlock()
			messageQueue.NewQueueCount <- messageQueue.QueueCount

			return
		}

		messageQueue.Mutex.Lock()
		queueCount := messageQueue.QueueCount
		messageQueue.Mutex.Unlock()

		// If the message is terminal, we want to make sure that all messages are sent before we return
		if queueCount > 0 {
			for newCount := range messageQueue.NewQueueCount {
				if newCount == 0 {
					break
				}
			}
		}

		setFunc()
	}()
}

func UpdateStatus(gs BaseGlobalState, status string) {
	if gs.GetStatus() != status {
		gs.OrchestratorClient().HSet(gs.Ctx(), gs.JobId(), "status", status)
		gs.SetStatus(status)
		DispatchMessage(gs, status, "STATUS")
	}
}

func UpdateStatusNoSet(gs BaseGlobalState, status string) {
	if gs.GetStatus() != status {
		DispatchMessageNoSet(gs, status, "STATUS")
	}
}

func HandleStringError(gs BaseGlobalState, errString string) {
	DispatchMessage(gs, errString, "ERROR")
	UpdateStatus(gs, "FAILURE")
}

func HandleError(gs BaseGlobalState, err error) {
	DispatchMessage(gs, err.Error(), "ERROR")
	UpdateStatus(gs, "FAILURE")
}

func HandleErrorNoSet(gs BaseGlobalState, err error) {
	DispatchMessageNoSet(gs, err.Error(), "ERROR")
	UpdateStatusNoSet(gs, "FAILURE")
}
