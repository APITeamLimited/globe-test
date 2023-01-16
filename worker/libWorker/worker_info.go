package libWorker

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/APITeamLimited/globe-test/lib"
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
		CreditsManager    *lib.CreditsManager
		Standalone        bool
	}

	MarkMessage struct {
		Mark    string                 `json:"mark"`
		Message map[string]interface{} `json:"message"`
	}
)

type Message struct {
	JobId       string    `json:"jobId"`
	ChildJobId  string    `json:"childJobId"`
	Time        time.Time `json:"time"`
	WorkerId    string    `json:"workerId"`
	Message     string    `json:"message"`
	MessageType string    `json:"messageType"`
}

type MessageQueue struct {
	Mutex sync.Mutex

	// The count of currently actively being sent messages
	QueueCount    int
	NewQueueCount chan int
}

func DispatchMessage(gs BaseGlobalState, message string, messageType string) {
	go func() {
		serializedMessage, err := serializeMessage(gs, message, messageType)
		if err != nil {
			fmt.Println("DispatchMessage: Error marshalling message", err)
			return
		}

		messageQueue := gs.MessageQueue()

		isTerminal := messageType == "STATUS" && (message == "FAILURE" || message == "SUCCESS")
		if !isTerminal {
			messageQueue.Mutex.Lock()
			messageQueue.QueueCount++
			messageQueue.Mutex.Unlock()

			err := gs.Client().Publish(gs.Ctx(), fmt.Sprintf("worker:executionUpdates:%s", gs.JobId()), serializedMessage).Err()

			messageQueue.Mutex.Lock()
			messageQueue.QueueCount--
			messageQueue.Mutex.Unlock()

			if err != nil {
				fmt.Println("DispatchMessage: Error publishing message", err)
			}

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

		err = gs.Client().Publish(gs.Ctx(), fmt.Sprintf("worker:executionUpdates:%s", gs.JobId()), serializedMessage).Err()
		if err != nil {
			fmt.Println("DispatchMessage: Error publishing message", err)
		}
	}()
}

func serializeMessage(gs BaseGlobalState, message string, messageType string) ([]byte, error) {
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
		return nil, err
	}

	return messageJson, nil
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
