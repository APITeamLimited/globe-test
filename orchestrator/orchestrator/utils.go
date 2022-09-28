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

func DispatchMessage(ctx context.Context, orchestratorClient *redis.Client, jobId string, orchestratorId string, message string, messageType string) ***REMOVED***
	timestamp := time.Now().UnixMilli()

	var messageStruct = OrchestratorMessage***REMOVED***
		JobId:          jobId,
		Time:           timestamp,
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

func DispatchWorkerMessage(ctx context.Context, orchestratorClient *redis.Client, jobId string, workerId string, message string, messageType string) ***REMOVED***
	timestamp := time.Now().UnixMilli()

	var messageStruct = WorkerMessage***REMOVED***
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
	orchestratorClient.SAdd(ctx, updatesKey, messageJson)

	// Dispatch to channel
	orchestratorClient.Publish(ctx, fmt.Sprintf("orchestrator:executionUpdates:%s", jobId), string(messageJson))
***REMOVED***

func UpdateStatus(ctx context.Context, orchestratorClient *redis.Client, jobId string, orchestratorId string, status string) ***REMOVED***
	orchestratorClient.HSet(ctx, jobId, "status", status)
	DispatchMessage(ctx, orchestratorClient, jobId, orchestratorId, status, "STATUS")
***REMOVED***

func HandleStringError(ctx context.Context, orchestratorClient *redis.Client, jobId string, orchestratorId string, errString string) ***REMOVED***
	DispatchMessage(ctx, orchestratorClient, jobId, orchestratorId, errString, "ERROR")
	UpdateStatus(ctx, orchestratorClient, jobId, orchestratorId, "FAILED")
***REMOVED***

func HandleError(ctx context.Context, orchestratorClient *redis.Client, jobId string, orchestratorId string, err error) ***REMOVED***
	DispatchMessage(ctx, orchestratorClient, jobId, orchestratorId, err.Error(), "ERROR")
	UpdateStatus(ctx, orchestratorClient, jobId, orchestratorId, "FAILED")
***REMOVED***

func GetEnvVariable(key, defaultValue string) string ***REMOVED***
	SYSTEM_ENV := os.Getenv("SYSTEM_ENV")

	if SYSTEM_ENV == "" ***REMOVED***
		// Assume development environment
		SYSTEM_ENV = "development"
	***REMOVED***

	if value := os.Getenv(key); value != "" ***REMOVED***
		return value
	***REMOVED***

	if SYSTEM_ENV == "development" ***REMOVED***
		return defaultValue
	***REMOVED***

	panic(fmt.Sprintf("%s is not set, and environment was %s, not development", key, SYSTEM_ENV))
***REMOVED***

func setInBucket(bucket *gridfs.Bucket, filename string, data []byte, contentType string, globeTestLogsId primitive.ObjectID) error ***REMOVED***
	bucketOptions := options.GridFSUpload()
	bucketOptions.Metadata = map[string]interface***REMOVED******REMOVED******REMOVED***
		"filename":    filename,
		"contentType": contentType,
	***REMOVED***

	uploadStream, err := bucket.OpenUploadStreamWithID(globeTestLogsId, filename, bucketOptions)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer uploadStream.Close()
	_, err = uploadStream.Write(data)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (e *ExecutionList) addJob(job map[string]string) ***REMOVED***
	e.mutex.Lock()
	e.currentJobs[job["id"]] = job
	e.mutex.Unlock()
***REMOVED***

func (e *ExecutionList) removeJob(jobId string) ***REMOVED***
	e.mutex.Lock()
	delete(e.currentJobs, jobId)
	e.mutex.Unlock()
***REMOVED***
