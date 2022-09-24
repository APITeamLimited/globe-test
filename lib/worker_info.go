package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/APITeamLimited/redis/v9"
)

type WorkerInfo struct ***REMOVED***
	Client         *redis.Client
	JobId          string
	ScopeId        string
	OrchestratorId string
	WorkerId       string
	Ctx            context.Context
	environment    *map[string]string
	body           string
	headers        *map[string]string
	parameters     *map[string]string
***REMOVED***

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
	JobId       string `json:"jobId"`
	Time        int64  `json:"time"`
	WorkerId    string `json:"workerId"`
	Message     string `json:"message"`
	MessageType string `json:"messageType"`
***REMOVED***

func DispatchMessage(ctx context.Context, client *redis.Client, jobId string, workerId string, message string, messageType string) ***REMOVED***
	timestamp := time.Now().UnixMilli()
	stampedTag := fmt.Sprintf("%d:%s", timestamp, workerId)

	var messageStruct = Message***REMOVED***
		JobId:       jobId,
		Time:        timestamp,
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
	updatesKey := fmt.Sprintf("%s:updates", jobId)
	client.HSet(ctx, updatesKey, stampedTag, messageJson)

	// Dispatch to channel
	client.Publish(ctx, fmt.Sprintf("worker:executionUpdates:%s", jobId), messageJson)
***REMOVED***

func UpdateStatus(ctx context.Context, client *redis.Client, jobId string, workerId string, status string) ***REMOVED***
	client.HSet(ctx, jobId, "status", status)
	DispatchMessage(ctx, client, jobId, workerId, status, "STATUS")
***REMOVED***

func HandleStringError(ctx context.Context, client *redis.Client, jobId string, workerId string, errString string) ***REMOVED***
	DispatchMessage(ctx, client, jobId, workerId, errString, "ERROR")
	UpdateStatus(ctx, client, jobId, workerId, "FAILED")
***REMOVED***

func HandleError(ctx context.Context, client *redis.Client, jobId string, workerId string, err error) ***REMOVED***
	DispatchMessage(ctx, client, jobId, workerId, err.Error(), "ERROR")
	UpdateStatus(ctx, client, jobId, workerId, "FAILED")
***REMOVED***
