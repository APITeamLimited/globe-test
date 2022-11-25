package libWorker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/redis/v9"
)

type (
	Collection struct ***REMOVED***
		Variables map[string]string
		Name      string
	***REMOVED***

	Environment struct ***REMOVED***
		Variables map[string]string
		Name      string
	***REMOVED***

	KeyValueItem struct ***REMOVED***
		Key   string `json:"key"`
		Value string `json:"value"`
	***REMOVED***

	WorkerInfo struct ***REMOVED***
		Client            *redis.Client
		JobId             string
		ChildJobId        string
		ScopeId           string
		OrchestratorId    string
		WorkerId          string
		Ctx               context.Context
		Environment       *Environment
		Collection        *Collection
		WorkerOptions     Options
		FinalRequest      map[string]interface***REMOVED******REMOVED***
		UnderlyingRequest map[string]interface***REMOVED******REMOVED***
		Gs                *BaseGlobalState
		VerifiedDomains   []string
		SubFraction       float64
		CreditsManager    *lib.CreditsManager
	***REMOVED***

	MarkMessage struct ***REMOVED***
		Mark    string                 `json:"mark"`
		Message map[string]interface***REMOVED******REMOVED*** `json:"message"`
	***REMOVED***
)

type Message struct ***REMOVED***
	JobId       string    `json:"jobId"`
	ChildJobId  string    `json:"childJobId"`
	Time        time.Time `json:"time"`
	WorkerId    string    `json:"workerId"`
	Message     string    `json:"message"`
	MessageType string    `json:"messageType"`
***REMOVED***

func DispatchMessage(gs BaseGlobalState, message string, messageType string) ***REMOVED***
	var messageStruct = Message***REMOVED***
		JobId:       gs.JobId(),
		ChildJobId:  gs.ChildJobId(),
		Time:        time.Now(),
		WorkerId:    gs.WorkerId(),
		Message:     message,
		MessageType: messageType,
	***REMOVED***

	messageJson, err := json.Marshal(messageStruct)
	if err != nil ***REMOVED***
		fmt.Println("Error marshalling message")
		return
	***REMOVED***

	// Worker doesn't need to set the message, it's just for the orchestrator and will be
	// instantly received by the orchestrator

	// Dispatch to channel
	gs.Client().Publish(gs.Ctx(), fmt.Sprintf("worker:executionUpdates:%s", gs.JobId()), messageJson)
***REMOVED***

func UpdateStatus(gs BaseGlobalState, status string) ***REMOVED***
	if gs.GetWorkerStatus() != status ***REMOVED***
		gs.Client().HSet(gs.Ctx(), gs.JobId(), "status", status)
		DispatchMessage(gs, status, "STATUS")
		gs.SetWorkerStatus(status)
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
