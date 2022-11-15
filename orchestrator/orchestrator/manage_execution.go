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
func manageExecution(gs *globalState, orchestratorClient *redis.Client, workerClients libOrch.WorkerClients, job libOrch.Job,
	orchestratorId string, executionList *ExecutionList, storeMongoDB *mongo.Database, optionsErr error) bool ***REMOVED***
	// Setup the job

	healthy := optionsErr == nil

	if healthy ***REMOVED***
		marshalledOptions, err := json.Marshal(job.Options)
		if err != nil ***REMOVED***
			libOrch.HandleStringError(gs, fmt.Sprintf("Error marshalling runtime options: %s", err.Error()))
			healthy = false
		***REMOVED***

		libOrch.DispatchMessage(gs, string(marshalledOptions), "OPTIONS")

		(*gs.MetricsStore()).InitMetricsStore(job.Options)
	***REMOVED***

	scope := job.Scope

	childJobs, err := determineChildJobs(healthy, job, job.Options, workerClients)
	if err != nil ***REMOVED***
		libOrch.HandleError(gs, err)
		healthy = false
	***REMOVED***

	// Run the job

	result := "FAILURE"

	if healthy ***REMOVED***
		result, err = handleExecution(gs, job.Options, scope, childJobs, job.Id)
		if err != nil ***REMOVED***
			libOrch.HandleError(gs, err)
		***REMOVED***
	***REMOVED***

	libOrch.UpdateStatus(gs, result)

	// Storing and cleaning up

	(*gs.MetricsStore()).Stop()

	// Create GlobeTest logs store receipt, note this must be sent after cleanup
	globeTestLogsReceipt := primitive.NewObjectID()
	globeTestLogsReceiptMessage := &libOrch.MarkMessage***REMOVED***
		Mark:    "GlobeTestLogsStoreReceipt",
		Message: globeTestLogsReceipt.Hex(),
	***REMOVED***

	marshalledGlobeTestReceipt, err := json.Marshal(globeTestLogsReceiptMessage)
	if err != nil ***REMOVED***
		fmt.Println("Error marshalling GlobeTestLogsStoreReceipt", err)
		libOrch.HandleError(gs, err)
		return false
	***REMOVED***
	libOrch.DispatchMessage(gs, string(marshalledGlobeTestReceipt), "MARK")

	//Create Metrics Store receipt, note this must be sent after cleanup
	metricsStoreReceipt := primitive.NewObjectID()
	metricsStoreReceiptMessage := &libOrch.MarkMessage***REMOVED***
		Mark:    "MetricsStoreReceipt",
		Message: metricsStoreReceipt.Hex(),
	***REMOVED***

	marshalledMetricsStoreReceipt, err := json.Marshal(metricsStoreReceiptMessage)
	if err != nil ***REMOVED***
		fmt.Println("Error marshalling metrics store receipt", err)
		libOrch.HandleError(gs, err)
		return false
	***REMOVED***
	libOrch.DispatchMessage(gs, string(marshalledMetricsStoreReceipt), "MARK")

	// Clean up the job and store result in Mongo
	err = cleanup(gs, job, childJobs, storeMongoDB, scope, globeTestLogsReceipt, metricsStoreReceipt)
	if err != nil ***REMOVED***
		fmt.Println("Error cleaning up", err)
		libOrch.HandleErrorNoSet(gs, err)
		libOrch.UpdateStatusNoSet(gs, result)
	***REMOVED*** else ***REMOVED***
		libOrch.UpdateStatusNoSet(gs, fmt.Sprintf("COMPLETED_%s", result))
	***REMOVED***

	return healthy
***REMOVED***
