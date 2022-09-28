package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/APITeamLimited/redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
Cleans up the worker and orchestrator clients, storing all results in storeMongo
*/
func cleanup(ctx context.Context, job map[string]string, childJobs map[string][]map[string]string, orchestratorClient *redis.Client, workerClients map[string]*redis.Client, orchestratorId string, storeMongoDB *mongo.Database, scope map[string]string, globeTestLogsId primitive.ObjectID) {
	// Clean up worker
	// Set job in orchestrator redis

	marshalledJobInfo, err := json.Marshal(job)
	if err != nil {
		HandleStringError(ctx, orchestratorClient, job["jobId"], orchestratorId, "Error marshalling job info")
	}

	DispatchMessage(ctx, orchestratorClient, job["id"], orchestratorId, string(marshalledJobInfo), "JOB_INFO")

	go func() {
		for _, value := range workerClients {
			for _, zone := range childJobs {
				for _, childJob := range zone {
					updatesKey := fmt.Sprintf("%s:updates", childJob["id"])
					value.Del(ctx, updatesKey)
					value.Del(ctx, childJob["id"])

					// Remove childJob["id"] from k6:executionHistory set
					value.SRem(ctx, "k6:executionHistory", childJob["id"])
				}
			}
		}
	}()

	// Store results in MongoDB
	bucketName := fmt.Sprintf("%s:%s", scope["variant"], scope["variantTargetId"])
	jobBucket, err := gridfs.NewBucket(storeMongoDB, options.GridFSBucket().SetName(bucketName))
	if err != nil {
		HandleStringError(ctx, orchestratorClient, job["id"], orchestratorId, fmt.Sprintf("Error creating job output bucket: %s", err.Error()))
		return
	}

	updatesKey := fmt.Sprintf("%s:updates", job["id"])

	unparsedMessages, err := orchestratorClient.SMembers(ctx, updatesKey).Result()
	if err != nil {
		HandleStringError(ctx, orchestratorClient, job["id"], orchestratorId, fmt.Sprintf("Error getting unparsed messages: %s", err.Error()))
		return
	}

	var logs []OrchestratorOrWorkerMessage

	for _, value := range unparsedMessages {
		var message OrchestratorOrWorkerMessage
		err := json.Unmarshal([]byte(value), &message)
		if err != nil {
			HandleStringError(ctx, orchestratorClient, job["id"], orchestratorId, fmt.Sprintf("Error unmarshalling message: %s", err.Error()))
			return
		}

		// TODO: make a seperate store datatype for large data
		/*
			if message.MessageType == "STORE" {
				parsedStoreMessage := storeMessage{}

				err := json.Unmarshal([]byte(message.Message), &parsedStoreMessage)
				if err != nil {
					HandleStringError(ctx, orchestratorClient, job["id"], orchestratorId, fmt.Sprintf("Error unmarshalling message: %s", err.Error()))
					return
				}

				storeTag := fmt.Sprintf("%s:store:%s", job["id"], parsedStoreMessage.Filename)

				setInBucket(jobBucket, storeTag, []byte(parsedStoreMessage.Message))

				// Remove stored value from logs
				message.MessageType = "STORE_RECEIPT"
				message.Message = storeTag
			}
		*/
		logs = append(logs, message)
	}

	// Convert logs to JSON and set in bucket
	logsJSON, err := json.Marshal(logs)
	if err != nil {
		HandleStringError(ctx, orchestratorClient, job["id"], orchestratorId, fmt.Sprintf("Error marshalling logs: %s", err.Error()))
		return
	}

	err = setInBucket(jobBucket, fmt.Sprintf("GlobeTest:%s:messages.json", job["id"]), logsJSON, "application/json", globeTestLogsId)
	if err != nil {
		// Can't alert client here, as the client has already been cleaned up
		fmt.Printf("Error setting logs in bucket: %s\n", err.Error())
		return
	}

	// Clean up orchestrator
	orchestratorClient.Del(ctx, updatesKey)
	orchestratorClient.Del(ctx, job["id"])
	orchestratorClient.SRem(ctx, "orchestrator:executionHistory", job["id"])
}
