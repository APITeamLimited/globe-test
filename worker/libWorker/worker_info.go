package libWorker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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
	***REMOVED***

	MarkMessage struct ***REMOVED***
		Mark    string                 `json:"mark"`
		Message map[string]interface***REMOVED******REMOVED*** `json:"message"`
	***REMOVED***
)

func GetTestWorkerInfo() *WorkerInfo ***REMOVED***
	return &WorkerInfo***REMOVED***
		Client: redis.NewClient(&redis.Options***REMOVED***
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		***REMOVED***),
		JobId:          "4d2b8a88-07e6-4e70-9a53-de45c273b3d6",
		ScopeId:        "7faae966-211d-4b41-a9da-d9ae634ad085",
		OrchestratorId: "33f39131-3cec-4e9c-aff9-66d7c7b0e4b8",
		WorkerId:       "46221780-2f61-4733-a181-9d34684734b9",
		Ctx:            context.Background(),
	***REMOVED***
***REMOVED***

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
