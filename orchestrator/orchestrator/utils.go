package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/APITeamLimited/redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func GetEnvVariable(key, defaultValue string) string {
	SYSTEM_ENV := os.Getenv("SYSTEM_ENV")

	if SYSTEM_ENV == "" {
		// Assume development environment
		SYSTEM_ENV = "development"
	}

	if value := os.Getenv(key); value != "" {
		return value
	}

	if SYSTEM_ENV == "development" {
		return defaultValue
	}

	panic(fmt.Sprintf("%s is not set, and environment was %s, not development", key, SYSTEM_ENV))
}

func setInBucket(bucket *gridfs.Bucket, filename string, data []byte, contentType string, globeTestLogsId primitive.ObjectID) error {
	bucketOptions := options.GridFSUpload()
	bucketOptions.Metadata = map[string]interface{}{
		"filename":    filename,
		"contentType": contentType,
	}

	uploadStream, err := bucket.OpenUploadStreamWithID(globeTestLogsId, filename, bucketOptions)
	if err != nil {
		return err
	}
	defer uploadStream.Close()
	_, err = uploadStream.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (e *ExecutionList) addJob(job map[string]string) {
	e.mutex.Lock()
	e.currentJobs[job["id"]] = job
	e.mutex.Unlock()
}

func (e *ExecutionList) removeJob(jobId string) {
	e.mutex.Lock()
	delete(e.currentJobs, jobId)
	e.mutex.Unlock()
}
