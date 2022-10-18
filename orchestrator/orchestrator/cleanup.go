package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
Cleans up the worker and orchestrator clients, storing all results in storeMongo
*/
func cleanup(ctx context.Context, job libOrch.Job, childJobs map[string]jobDistribution,
	orchestratorClient *redis.Client, orchestratorId string, storeMongoDB *mongo.Database,
	scope libOrch.Scope, globeTestLogsReceipt primitive.ObjectID,
	metricsStoreReceipt primitive.ObjectID) error {
	// Clean up worker
	// Set job in orchestrator redis

	marshalledJobInfo, err := json.Marshal(job)
	if err != nil {
		return err
	}

	libOrch.DispatchMessage(ctx, orchestratorClient, job.Id, orchestratorId, string(marshalledJobInfo), "JOB_INFO")

	go func() {
		for _, jobDistribution := range childJobs {
			client := jobDistribution.workerClient

			for _, childJob := range *jobDistribution.jobs {
				client.Del(ctx, childJob.ChildJobId)

				// Remove childJob["id"] from worker:executionHistory set
				client.SRem(ctx, "worker:executionHistory", childJob.ChildJobId)

			}
		}
	}()

	// Store results in MongoDB
	bucketName := fmt.Sprintf("%s:%s", scope.Variant, scope.VariantTargetId)
	jobBucket, err := gridfs.NewBucket(storeMongoDB, options.GridFSBucket().SetName(bucketName))
	if err != nil {
		return err
	}

	updatesKey := fmt.Sprintf("%s:updates", job.Id)

	unparsedMessages, err := orchestratorClient.SMembers(ctx, updatesKey).Result()
	if err != nil {
		return err
	}

	var globeTestLogs []libOrch.OrchestratorOrWorkerMessage
	var metrics []libOrch.OrchestratorOrWorkerMessage

	for _, value := range unparsedMessages {
		// Declare here else fields will be inherited from previous iteration
		var message libOrch.OrchestratorOrWorkerMessage

		err := json.Unmarshal([]byte(value), &message)
		if err != nil {
			return err
		}

		// TODO: make a seperate store datatype for large data
		/*
			if message.MessageType == "STORE" {
				parsedStoreMessage := storeMessage{}

				err := json.Unmarshal([]byte(message.Message), &parsedStoreMessage)
				if err != nil {
					libOrch.HandleStringError(ctx, orchestratorClient, job.Id, orchestratorId, fmt.Sprintf("Error unmarshalling message: %s", err.Error()))
					return
				}

				storeTag := fmt.Sprintf("%s:store:%s", job.Id, parsedStoreMessage.Filename)

				setInBucket(jobBucket, storeTag, []byte(parsedStoreMessage.Message))

				// Remove stored value from logs
				message.MessageType = "STORE_RECEIPT"
				message.Message = storeTag
			}
		*/

		if message.MessageType == "METRICS" {
			metrics = append(metrics, message)
		} else if message.MessageType == "COLLECTION_VARIABLES" || message.MessageType == "ENVIRONMENT_VARIABLES" {
			continue
		} else {
			globeTestLogs = append(globeTestLogs, message)
		}
	}

	channel := make(chan error)

	go func() {
		// Convert logs to JSON and set in bucket
		globeTestLogsMarshalled, err := json.Marshal(globeTestLogs)
		if err != nil {
			// Can't alert client here, as the client has already been cleaned up
			channel <- fmt.Errorf("error marshalling logs: %s", err.Error())
			return
		}

		err = libOrch.SetInBucket(jobBucket, fmt.Sprintf("GlobeTest:%s:messages.json", job.Id), globeTestLogsMarshalled, "application/json", globeTestLogsReceipt)
		if err != nil {
			// Can't alert client here, as the client has already been cleaned up
			channel <- fmt.Errorf("error setting logs in bucket: %s", err.Error())
			return
		}

		channel <- nil
	}()

	go func() {
		// Convert metrics to JSON and set in bucket
		metricsMarshalled, err := json.Marshal(metrics)
		if err != nil {
			channel <- fmt.Errorf("error marshalling metrics: %s", err.Error())
			return
		}

		err = libOrch.SetInBucket(jobBucket, fmt.Sprintf("GlobeTest:%s:metrics.json", job.Id), metricsMarshalled, "application/json", metricsStoreReceipt)
		if err != nil {
			channel <- fmt.Errorf("error setting metrics in bucket: %s", err.Error())
			return
		}

		channel <- nil
	}()

	for i := 0; i < 2; i++ {
		err := <-channel
		if err != nil {
			return err
		}
	}

	// Clean up orchestrator
	// Set types to expire so lagging users can access environment variables
	orchestratorClient.Expire(ctx, updatesKey, time.Second*10)
	orchestratorClient.Expire(ctx, job.Id, time.Second*10)
	orchestratorClient.SRem(ctx, "orchestrator:executionHistory", job.Id)

	return nil
}
