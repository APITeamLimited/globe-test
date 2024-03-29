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
func manageExecution(gs libOrch.BaseGlobalState, orchestratorClient *redis.Client, job libOrch.Job,
	orchestratorId string, executionList *ExecutionList, storeMongoDB *mongo.Database, optionsErr error) bool {
	// Setup the job

	healthy := optionsErr == nil

	childJobs, err := determineChildJobs(healthy, job, job.Options)
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

		gs.MetricsStore().InitAggregator(childJobs, job.Options.Thresholds)
		gs.MetricsStore().StartConsoleLogging()

		defer gs.MetricsStore().StopAndCleanup()
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
	testInfoStoreReceipt, err := testInfoStoreReceipt(gs)
	if err != nil {
		libOrch.HandleError(gs, err)
		return false
	}

	// Clean up the job and store result in Mongo
	err = cleanup(gs, job, childJobs, storeMongoDB, job.Scope, testInfoStoreReceipt)
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

func testInfoStoreReceipt(gs libOrch.BaseGlobalState) (*primitive.ObjectID, error) {
	testInfoStoreReceipt := primitive.NewObjectID()
	testInfoStoreReceiptMessage := &libOrch.MarkMessage{
		Mark:    "TestInfoStoreReceipt",
		Message: testInfoStoreReceipt.Hex(),
	}

	marshalledGlobeTestReceipt, err := json.Marshal(testInfoStoreReceiptMessage)
	if err != nil {
		fmt.Println("Error marshalling testInfoStoreReceipt", err)
		libOrch.HandleError(gs, err)
		return nil, err
	}

	if gs.Standalone() {
		libOrch.DispatchMessage(gs, string(marshalledGlobeTestReceipt), "MARK")
	}

	return &testInfoStoreReceipt, nil
}
