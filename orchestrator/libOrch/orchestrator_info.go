package libOrch

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/APITeamLimited/redis/v9"
)

func DispatchMessage(ctx context.Context, orchestratorClient *redis.Client, jobId string, orchestratorId string, message string, messageType string) ***REMOVED***
	var messageStruct = OrchestratorMessage***REMOVED***
		JobId:          jobId,
		Time:           time.Now(),
		OrchestratorId: orchestratorId,
		Message:        message,
		MessageType:    messageType,
	***REMOVED***

	messageJson, err := json.Marshal(messageStruct)
	if err != nil ***REMOVED***
		fmt.Println("Error marshalling message")
		return
	***REMOVED***

	// Update main job
	updatesKey := fmt.Sprintf("%s:updates", jobId)
	orchestratorClient.SAdd(ctx, updatesKey, messageJson)

	// Dispatch to channel
	orchestratorClient.Publish(ctx, fmt.Sprintf("orchestrator:executionUpdates:%s", jobId), string(messageJson))
***REMOVED***

func DispatchMessageNoSet(ctx context.Context, orchestratorClient *redis.Client, jobId string, orchestratorId string, message string, messageType string) ***REMOVED***
	var messageStruct = OrchestratorMessage***REMOVED***
		JobId:          jobId,
		Time:           time.Now(),
		OrchestratorId: orchestratorId,
		Message:        message,
		MessageType:    messageType,
	***REMOVED***

	messageJson, err := json.Marshal(messageStruct)
	if err != nil ***REMOVED***
		fmt.Println("Error marshalling message")
		return
	***REMOVED***

	// Dispatch to channel
	orchestratorClient.Publish(ctx, fmt.Sprintf("orchestrator:executionUpdates:%s", jobId), string(messageJson))
***REMOVED***

func DispatchWorkerMessage(ctx context.Context, orchestratorClient *redis.Client, jobId string, workerId string, message string, messageType string) ***REMOVED***
	var messageStruct = WorkerMessage***REMOVED***
		JobId:       jobId,
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
	updatesKey := fmt.Sprintf("%s:updates", jobId)
	orchestratorClient.SAdd(ctx, updatesKey, messageJson)

	// Dispatch to channel
	orchestratorClient.Publish(ctx, fmt.Sprintf("orchestrator:executionUpdates:%s", jobId), string(messageJson))
***REMOVED***

func UpdateStatus(ctx context.Context, orchestratorClient *redis.Client, jobId string, orchestratorId string, status string) ***REMOVED***
	orchestratorClient.HSet(ctx, jobId, "status", status)
	DispatchMessage(ctx, orchestratorClient, jobId, orchestratorId, status, "STATUS")
***REMOVED***

func UpdateStatusNoSet(ctx context.Context, orchestratorClient *redis.Client, jobId string, orchestratorId string, status string) ***REMOVED***
	DispatchMessageNoSet(ctx, orchestratorClient, jobId, orchestratorId, status, "STATUS")
***REMOVED***

func HandleStringError(ctx context.Context, orchestratorClient *redis.Client, jobId string, orchestratorId string, errString string) ***REMOVED***
	DispatchMessage(ctx, orchestratorClient, jobId, orchestratorId, errString, "ERROR")
	UpdateStatus(ctx, orchestratorClient, jobId, orchestratorId, "FAILURE")
***REMOVED***

func HandleError(ctx context.Context, orchestratorClient *redis.Client, jobId string, orchestratorId string, err error) ***REMOVED***
	DispatchMessage(ctx, orchestratorClient, jobId, orchestratorId, err.Error(), "ERROR")
	UpdateStatus(ctx, orchestratorClient, jobId, orchestratorId, "FAILURE")
***REMOVED***

func HandleErrorNoSet(ctx context.Context, orchestratorClient *redis.Client, jobId string, orchestratorId string, err error) ***REMOVED***
	DispatchMessageNoSet(ctx, orchestratorClient, jobId, orchestratorId, err.Error(), "ERROR")
	UpdateStatusNoSet(ctx, orchestratorClient, jobId, orchestratorId, "FAILURE")
***REMOVED***
