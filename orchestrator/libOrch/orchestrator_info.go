package libOrch

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/APITeamLimited/redis/v9"
)

func DispatchMessage(ctx context.Context, orchestratorClient *redis.Client, jobId string, orchestratorId string, message string, messageType string) {
	timestamp := time.Now().UnixMilli()

	var messageStruct = OrchestratorMessage{
		JobId:          jobId,
		Time:           timestamp,
		OrchestratorId: orchestratorId,
		Message:        message,
		MessageType:    messageType,
	}

	messageJson, err := json.Marshal(messageStruct)
	if err != nil {
		fmt.Println("Error marshalling message")
		return
	}

	// Update main job
	updatesKey := fmt.Sprintf("%s:updates", jobId)
	orchestratorClient.SAdd(ctx, updatesKey, messageJson)

	// Dispatch to channel
	orchestratorClient.Publish(ctx, fmt.Sprintf("orchestrator:executionUpdates:%s", jobId), string(messageJson))
}

func DispatchWorkerMessage(ctx context.Context, orchestratorClient *redis.Client, jobId string, workerId string, message string, messageType string) {
	timestamp := time.Now().UnixMilli()

	var messageStruct = WorkerMessage{
		JobId:       jobId,
		Time:        timestamp,
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
	updatesKey := fmt.Sprintf("%s:updates", jobId)
	orchestratorClient.SAdd(ctx, updatesKey, messageJson)

	// Dispatch to channel
	orchestratorClient.Publish(ctx, fmt.Sprintf("orchestrator:executionUpdates:%s", jobId), string(messageJson))
}

func UpdateStatus(ctx context.Context, orchestratorClient *redis.Client, jobId string, orchestratorId string, status string) {
	orchestratorClient.HSet(ctx, jobId, "status", status)
	DispatchMessage(ctx, orchestratorClient, jobId, orchestratorId, status, "STATUS")
}

func HandleStringError(ctx context.Context, orchestratorClient *redis.Client, jobId string, orchestratorId string, errString string) {
	DispatchMessage(ctx, orchestratorClient, jobId, orchestratorId, errString, "ERROR")
	UpdateStatus(ctx, orchestratorClient, jobId, orchestratorId, "FAILED")
}

func HandleError(ctx context.Context, orchestratorClient *redis.Client, jobId string, orchestratorId string, err error) {
	DispatchMessage(ctx, orchestratorClient, jobId, orchestratorId, err.Error(), "ERROR")
	UpdateStatus(ctx, orchestratorClient, jobId, orchestratorId, "FAILED")
}
