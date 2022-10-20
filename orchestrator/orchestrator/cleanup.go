package orchestrator

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
Cleans up the worker and orchestrator clients, storing all results in storeMongo
*/
func cleanup(gs libOrch.BaseGlobalState, job libOrch.Job, childJobs map[string]jobDistribution, storeMongoDB *mongo.Database,
	scope libOrch.Scope, globeTestLogsReceipt primitive.ObjectID, metricsStoreReceipt primitive.ObjectID) error ***REMOVED***

	// Clean up worker
	// Set job in orchestrator redis

	marshalledJobInfo, err := json.Marshal(job)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	libOrch.DispatchMessage(gs, string(marshalledJobInfo), "JOB_INFO")

	go func() ***REMOVED***
		for _, jobDistribution := range childJobs ***REMOVED***
			client := jobDistribution.workerClient

			for _, childJob := range jobDistribution.jobs ***REMOVED***
				client.Del(gs.Ctx(), childJob.ChildJobId)

				// Remove childJob["id"] from worker:executionHistory set
				client.SRem(gs.Ctx(), "worker:executionHistory", childJob.ChildJobId)

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

	unparsedMessages, err := gs.Client().SMembers(gs.Ctx(), updatesKey).Result()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	var globeTestLogs []libOrch.OrchestratorOrWorkerMessage
	var metrics []libOrch.OrchestratorOrWorkerMessage

	for _, value := range unparsedMessages ***REMOVED***
		// Declare here else fields will be inherited from previous iteration
		var message libOrch.OrchestratorOrWorkerMessage

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
		***REMOVED*** else if message.MessageType == "COLLECTION_VARIABLES" || message.MessageType == "ENVIRONMENT_VARIABLES" ***REMOVED***
			continue
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
	// Set types to expire so lagging users can access environment variables
	gs.Client().Expire(gs.Ctx(), updatesKey, time.Second*10)
	gs.Client().Expire(gs.Ctx(), job.Id, time.Second*10)
	gs.Client().SRem(gs.Ctx(), "orchestrator:executionHistory", job.Id)

	return nil
***REMOVED***
