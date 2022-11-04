package libWorker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/APITeamLimited/redis/v9"
)

type (
	Collection struct {
		Variables map[string]string
		Name      string
	}

	Environment struct {
		Variables map[string]string
		Name      string
	}

	KeyValueItem struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	WorkerInfo struct {
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
		FinalRequest      map[string]interface{}
		UnderlyingRequest map[string]interface{}
		Gs                *BaseGlobalState
		VerifiedDomains   []string
		SubFraction       float64
	}

	MarkMessage struct {
		Mark    string                 `json:"mark"`
		Message map[string]interface{} `json:"message"`
	}
)

func GetTestWorkerInfo() *WorkerInfo {
	return &WorkerInfo{
		Client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
		JobId:          "4d2b8a88-07e6-4e70-9a53-de45c273b3d6",
		ScopeId:        "7faae966-211d-4b41-a9da-d9ae634ad085",
		OrchestratorId: "33f39131-3cec-4e9c-aff9-66d7c7b0e4b8",
		WorkerId:       "46221780-2f61-4733-a181-9d34684734b9",
		Ctx:            context.Background(),
	}
}

type Message struct {
	JobId       string    `json:"jobId"`
	ChildJobId  string    `json:"childJobId"`
	Time        time.Time `json:"time"`
	WorkerId    string    `json:"workerId"`
	Message     string    `json:"message"`
	MessageType string    `json:"messageType"`
}

func DispatchMessage(gs BaseGlobalState, message string, messageType string) {
	var messageStruct = Message{
		JobId:       gs.JobId(),
		ChildJobId:  gs.ChildJobId(),
		Time:        time.Now(),
		WorkerId:    gs.WorkerId(),
		Message:     message,
		MessageType: messageType,
	}

	messageJson, err := json.Marshal(messageStruct)
	if err != nil {
		fmt.Println("Error marshalling message")
		return
	}

	// Worker doesn't need to set the message, it's just for the orchestrator and will be
	// instantly received by the orchestrator

	// Dispatch to channel
	gs.Client().Publish(gs.Ctx(), fmt.Sprintf("worker:executionUpdates:%s", gs.JobId()), messageJson)
}

func UpdateStatus(gs BaseGlobalState, status string) {
	if gs.GetWorkerStatus() != status {
		gs.Client().HSet(gs.Ctx(), gs.JobId(), "status", status)
		DispatchMessage(gs, status, "STATUS")
		gs.SetWorkerStatus(status)
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
