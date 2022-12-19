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

	// Update main job
	gs.OrchestratorClient().SAdd(gs.Ctx(), fmt.Sprintf("%s:updates", gs.JobId()), messageJson)

	// Dispatch to channel
	gs.OrchestratorClient().Publish(gs.Ctx(), fmt.Sprintf("orchestrator:executionUpdates:%s", gs.JobId()), string(messageJson))
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

	// Dispatch to channel
	gs.OrchestratorClient().Publish(gs.Ctx(), fmt.Sprintf("orchestrator:executionUpdates:%s", gs.JobId()), messageJson)
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

	// Update main job
	gs.OrchestratorClient().SAdd(gs.Ctx(), fmt.Sprintf("%s:updates", gs.JobId()), messageJson)

	// Dispatch to channel
	gs.OrchestratorClient().Publish(gs.Ctx(), fmt.Sprintf("orchestrator:executionUpdates:%s", gs.JobId()), messageJson)
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
