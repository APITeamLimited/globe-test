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
	scope libOrch.Scope, globeTestLogsReceipt primitive.ObjectID, metricsStoreReceipt primitive.ObjectID) error {
	// Clean up worker
	// Set job in orchestrator redis

	marshalledJobInfo, err := json.Marshal(job)
	if err != nil {
		return err
	}

	libOrch.DispatchMessage(gs, string(marshalledJobInfo), "JOB_INFO")

	go func() {
		for _, jobDistribution := range childJobs {
			client := jobDistribution.workerClient

			for _, childJob := range jobDistribution.Jobs {
				client.Del(gs.Ctx(), childJob.ChildJobId)

				// Remove childJob["id"] from worker:executionHistory set
				client.SRem(gs.Ctx(), "worker:executionHistory", childJob.ChildJobId)

			}
		}
	}()

	// Store results in MongoDB
	bucketName := fmt.Sprintf("%s:%s", scope.Variant, scope.VariantTargetId)

	var jobBucket *gridfs.Bucket
	if gs.Standalone() {
		jobBucket, err = gridfs.NewBucket(storeMongoDB, options.GridFSBucket().SetName(bucketName))
		if err != nil {
			return err
		}
	}

	updatesKey := fmt.Sprintf("%s:updates", job.Id)

	unparsedMessages, err := gs.OrchestratorClient().SMembers(gs.Ctx(), updatesKey).Result()
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

		globeTestLogsFilename := fmt.Sprintf("GlobeTest:%s:messages.json", job.Id)

		if gs.Standalone() {
			err = libOrch.SetInBucket(jobBucket, globeTestLogsFilename, globeTestLogsMarshalled, "application/json", globeTestLogsReceipt)
			if err != nil {
				// Can't alert client here, as the client has already been cleaned up
				channel <- fmt.Errorf("error setting logs in bucket: %s", err.Error())
				return
			}
		} else {
			// TODO
			localhostFile := libOrch.LocalhostFile{
				FileName: globeTestLogsFilename,
				Contents: string(globeTestLogsMarshalled),
			}

			marshalledLocalhostFile, err := json.Marshal(localhostFile)
			if err != nil {
				// Can't alert client here, as the client has already been cleaned up
				channel <- fmt.Errorf("error setting logs in bucket: %s", err.Error())
				return
			}

			libOrch.DispatchMessage(gs, string(marshalledLocalhostFile), "LOCALHOST_FILE")
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

		metricsFilename := fmt.Sprintf("GlobeTest:%s:metrics.json", job.Id)

		if gs.Standalone() {
			err = libOrch.SetInBucket(jobBucket, metricsFilename, metricsMarshalled, "application/json", metricsStoreReceipt)
			if err != nil {
				channel <- fmt.Errorf("error setting metrics in bucket: %s", err.Error())
				return
			}
		} else {
			// TODO
			localhostFile := libOrch.LocalhostFile{
				FileName: metricsFilename,
				Contents: string(metricsMarshalled),
			}

			marshalledLocalhostFile, err := json.Marshal(localhostFile)
			if err != nil {
				// Can't alert client here, as the client has already been cleaned up
				channel <- fmt.Errorf("error setting logs in bucket: %s", err.Error())
				return
			}

			libOrch.DispatchMessage(gs, string(marshalledLocalhostFile), "LOCALHOST_FILE")
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
	gs.OrchestratorClient().Expire(gs.Ctx(), updatesKey, time.Second*10)
	gs.OrchestratorClient().Expire(gs.Ctx(), job.Id, time.Second*10)
	gs.OrchestratorClient().SRem(gs.Ctx(), "orchestrator:executionHistory", job.Id)

	return nil
}
