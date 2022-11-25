package libOrch

import (
	"encoding/json"
	"fmt"
	"time"
)

func DispatchMessage(gs BaseGlobalState, message string, messageType string) ***REMOVED***
	var messageStruct = OrchestratorMessage***REMOVED***
		JobId:          gs.JobId(),
		Time:           time.Now(),
		OrchestratorId: gs.OrchestratorId(),
		Message:        message,
		MessageType:    messageType,
	***REMOVED***

	messageJson, err := json.Marshal(messageStruct)
	if err != nil ***REMOVED***
		fmt.Println("Error marshalling message")
		return
	***REMOVED***

	// Update main job
	gs.OrchestratorClient().SAdd(gs.Ctx(), fmt.Sprintf("%s:updates", gs.JobId()), messageJson)

	// Dispatch to channel
	gs.OrchestratorClient().Publish(gs.Ctx(), fmt.Sprintf("orchestrator:executionUpdates:%s", gs.JobId()), string(messageJson))
***REMOVED***

func DispatchMessageNoSet(gs BaseGlobalState, message string, messageType string) ***REMOVED***
	var messageStruct = OrchestratorMessage***REMOVED***
		JobId:          gs.JobId(),
		Time:           time.Now(),
		OrchestratorId: gs.OrchestratorId(),
		Message:        message,
		MessageType:    messageType,
	***REMOVED***

	messageJson, err := json.Marshal(messageStruct)
	if err != nil ***REMOVED***
		fmt.Println("Error marshalling message")
		return
	***REMOVED***

	// Dispatch to channel
	gs.OrchestratorClient().Publish(gs.Ctx(), fmt.Sprintf("orchestrator:executionUpdates:%s", gs.JobId()), string(messageJson))
***REMOVED***

func DispatchWorkerMessage(gs BaseGlobalState, workerId string, childJobId string, message string, messageType string) ***REMOVED***
	var messageStruct = WorkerMessage***REMOVED***
		JobId:       gs.JobId(),
		ChildJobId:  childJobId,
		Time:        time.Now(),
		WorkerId:    workerId,
		Message:     message,
		MessageType: messageType,
	***REMOVED***

	messageJson, err := json.Marshal(messageStruct)
	if err != nil ***REMOVED***
		fmt.Println("Error marshalling message")
		return
	***REMOVED***

	// Update main job
	gs.OrchestratorClient().SAdd(gs.Ctx(), fmt.Sprintf("%s:updates", gs.JobId()), messageJson)

	// Dispatch to channel
	gs.OrchestratorClient().Publish(gs.Ctx(), fmt.Sprintf("orchestrator:executionUpdates:%s", gs.JobId()), string(messageJson))
***REMOVED***

func UpdateStatus(gs BaseGlobalState, status string) ***REMOVED***
	if gs.GetStatus() != status ***REMOVED***
		gs.OrchestratorClient().HSet(gs.Ctx(), gs.JobId(), "status", status)
		gs.SetStatus(status)
		DispatchMessage(gs, status, "STATUS")
	***REMOVED***
***REMOVED***

func UpdateStatusNoSet(gs BaseGlobalState, status string) ***REMOVED***
	if gs.GetStatus() != status ***REMOVED***
		DispatchMessageNoSet(gs, status, "STATUS")
	***REMOVED***
***REMOVED***

func HandleStringError(gs BaseGlobalState, errString string) ***REMOVED***
	DispatchMessage(gs, errString, "ERROR")
	UpdateStatus(gs, "FAILURE")
***REMOVED***

func HandleError(gs BaseGlobalState, err error) ***REMOVED***
	DispatchMessage(gs, err.Error(), "ERROR")
	UpdateStatus(gs, "FAILURE")
***REMOVED***

func HandleErrorNoSet(gs BaseGlobalState, err error) ***REMOVED***
	DispatchMessageNoSet(gs, err.Error(), "ERROR")
	UpdateStatusNoSet(gs, "FAILURE")
***REMOVED***
