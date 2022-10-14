package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"

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
	metricsStoreReceipt primitive.ObjectID) error ***REMOVED***
	// Clean up worker
	// Set job in orchestrator redis

	marshalledJobInfo, err := json.Marshal(job)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	libOrch.DispatchMessage(ctx, orchestratorClient, job.Id, orchestratorId, string(marshalledJobInfo), "JOB_INFO")

	go func() ***REMOVED***
		for _, jobDistribution := range childJobs ***REMOVED***
			client := jobDistribution.workerClient

			for _, childJob := range *jobDistribution.jobs ***REMOVED***
				client.Del(ctx, childJob.ChildJobId)

				// Remove childJob["id"] from worker:executionHistory set
				client.SRem(ctx, "worker:executionHistory", childJob.ChildJobId)

			***REMOVED***
		***REMOVED***
	***REMOVED***()

	// Store results in MongoDB
	bucketName := fmt.Sprintf("%s:%s", scope.Variant, scope.VariantTargetId)
	jobBucket, err := gridfs.NewBucket(storeMongoDB, options.GridFSBucket().SetName(bucketName))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	updatesKey := fmt.Sprintf("%s:updates", job.Id)

	unparsedMessages, err := orchestratorClient.SMembers(ctx, updatesKey).Result()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	var globeTestLogs []libOrch.OrchestratorOrWorkerMessage
	var metrics []libOrch.OrchestratorOrWorkerMessage

	var message libOrch.OrchestratorOrWorkerMessage

	for _, value := range unparsedMessages ***REMOVED***
		err := json.Unmarshal([]byte(value), &message)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// TODO: make a seperate store datatype for large data
		/*
			if message.MessageType == "STORE" ***REMOVED***
				parsedStoreMessage := storeMessage***REMOVED******REMOVED***

				err := json.Unmarshal([]byte(message.Message), &parsedStoreMessage)
				if err != nil ***REMOVED***
					libOrch.HandleStringError(ctx, orchestratorClient, job.Id, orchestratorId, fmt.Sprintf("Error unmarshalling message: %s", err.Error()))
					return
				***REMOVED***

				storeTag := fmt.Sprintf("%s:store:%s", job.Id, parsedStoreMessage.Filename)

				setInBucket(jobBucket, storeTag, []byte(parsedStoreMessage.Message))

				// Remove stored value from logs
				message.MessageType = "STORE_RECEIPT"
				message.Message = storeTag
			***REMOVED***
		*/

		if message.MessageType == "METRICS" ***REMOVED***
			metrics = append(metrics, message)
		***REMOVED*** else ***REMOVED***
			globeTestLogs = append(globeTestLogs, message)
		***REMOVED***
	***REMOVED***

	channel := make(chan error)

	go func() ***REMOVED***
		// Convert logs to JSON and set in bucket
		globeTestLogsMarshalled, err := json.Marshal(globeTestLogs)
		if err != nil ***REMOVED***
			// Can't alert client here, as the client has already been cleaned up
			channel <- fmt.Errorf("error marshalling logs: %s", err.Error())
			return
		***REMOVED***

		err = libOrch.SetInBucket(jobBucket, fmt.Sprintf("GlobeTest:%s:messages.json", job.Id), globeTestLogsMarshalled, "application/json", globeTestLogsReceipt)
		if err != nil ***REMOVED***
			// Can't alert client here, as the client has already been cleaned up
			channel <- fmt.Errorf("error setting logs in bucket: %s", err.Error())
			return
		***REMOVED***

		channel <- nil
	***REMOVED***()

	go func() ***REMOVED***
		// Convert metrics to JSON and set in bucket
		metricsMarshalled, err := json.Marshal(metrics)
		if err != nil ***REMOVED***
			channel <- fmt.Errorf("error marshalling metrics: %s", err.Error())
			return
		***REMOVED***

		err = libOrch.SetInBucket(jobBucket, fmt.Sprintf("GlobeTest:%s:metrics.json", job.Id), metricsMarshalled, "application/json", metricsStoreReceipt)
		if err != nil ***REMOVED***
			channel <- fmt.Errorf("error setting metrics in bucket: %s", err.Error())
			return
		***REMOVED***

		channel <- nil
	***REMOVED***()

	for i := 0; i < 2; i++ ***REMOVED***
		err := <-channel
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// Clean up orchestrator
	orchestratorClient.Del(ctx, updatesKey)
	orchestratorClient.Del(ctx, job.Id)
	orchestratorClient.SRem(ctx, "orchestrator:executionHistory", job.Id)

	return nil
***REMOVED***
