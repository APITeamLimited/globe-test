package orchestrator

import (
	"encoding/json"
	"fmt"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Over-arching function that manages the execution of a job and handles its state and lifecycle
// This is the highest level function with global state
// Avoids use of credits as this will cause undesired side effects
func manageExecution(gs *globalState, orchestratorClient *redis.Client, workerClients libOrch.WorkerClients, job libOrch.Job,
	orchestratorId string, executionList *ExecutionList, storeMongoDB *mongo.Database, optionsErr error) bool {
	// Setup the job

	healthy := optionsErr == nil

	childJobs, err := determineChildJobs(healthy, job, job.Options, workerClients)
	if err != nil {
		libOrch.HandleError(gs, err)
		healthy = false
	}

	if healthy {
		marshalledOptions, err := json.Marshal(job.Options)
		if err != nil {
			libOrch.HandleStringError(gs, fmt.Sprintf("Error marshalling runtime options: %s", err.Error()))
			healthy = false
		}

		libOrch.DispatchMessage(gs, string(marshalledOptions), "OPTIONS")

		(*gs.MetricsStore()).InitMetricsStore(childJobs)
		defer (*gs.MetricsStore()).Cleanup()
	}

	// Run the job

	result := "FAILURE"

	if healthy {
		result, err = handleExecution(gs, job, childJobs)
		if err != nil {
			libOrch.HandleError(gs, err)
		}
	}

	libOrch.UpdateStatus(gs, result)

	// Storing and cleaning up

	// Create GlobeTest logs store receipt, note this must be sent after cleanup
	globeTestLogsStoreReceipt, err := globeTestLogsStoreReceipt(gs)
	if err != nil {
		libOrch.HandleError(gs, err)
		return false
	}

	metricsStoreReceipt, err := metricsStoreReceipt(gs)
	if err != nil {
		libOrch.HandleError(gs, err)
		return false
	}

	// Clean up the job and store result in Mongo
	err = cleanup(gs, job, childJobs, storeMongoDB, job.Scope, globeTestLogsStoreReceipt, metricsStoreReceipt)
	if err != nil {
		fmt.Println("Error cleaning up", err)
		libOrch.HandleErrorNoSet(gs, err)
		// This is needed
		libOrch.UpdateStatusNoSet(gs, result)
	} else {
		libOrch.UpdateStatusNoSet(gs, fmt.Sprintf("COMPLETED_%s", result))
	}

	return healthy
}

func metricsStoreReceipt(gs *globalState) (*primitive.ObjectID, error) {
	if !gs.Standalone() {
		return nil, nil
	}

	// Create Metrics Store receipt, note this must be sent after cleanup
	metricsStoreReceipt := primitive.NewObjectID()
	metricsStoreReceiptMessage := &libOrch.MarkMessage{
		Mark:    "MetricsStoreReceipt",
		Message: metricsStoreReceipt.Hex(),
	}

	marshalledMetricsStoreReceipt, err := json.Marshal(metricsStoreReceiptMessage)
	if err != nil {
		return nil, err
	}

	libOrch.DispatchMessage(gs, string(marshalledMetricsStoreReceipt), "MARK")

	return &metricsStoreReceipt, nil
}

func globeTestLogsStoreReceipt(gs *globalState) (*primitive.ObjectID, error) {
	globeTestLogsReceipt := primitive.NewObjectID()
	globeTestLogsReceiptMessage := &libOrch.MarkMessage{
		Mark:    "GlobeTestLogsStoreReceipt",
		Message: globeTestLogsReceipt.Hex(),
	}

	marshalledGlobeTestReceipt, err := json.Marshal(globeTestLogsReceiptMessage)
	if err != nil {
		fmt.Println("Error marshalling GlobeTestLogsStoreReceipt", err)
		libOrch.HandleError(gs, err)
		return nil, err
	}

	if gs.Standalone() {
		libOrch.DispatchMessage(gs, string(marshalledGlobeTestReceipt), "MARK")
	}

	return &globeTestLogsReceipt, nil
}
