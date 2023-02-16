package orchestrator

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/APITeamLimited/globe-test/errext"
	"github.com/APITeamLimited/globe-test/errext/exitcodes"
	"github.com/APITeamLimited/globe-test/metrics"
	"github.com/APITeamLimited/globe-test/metrics/engine"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Over-arching function that manages the execution of a job and handles its state and lifecycle
func manageExecution(gs libOrch.BaseGlobalState, orchestratorClient *redis.Client, job libOrch.Job,
	orchestratorId string, storeMongoDB *mongo.Database, optionsErr error) bool {
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
	}

	registry := metrics.NewRegistry()
	metrics.RegisterBuiltinMetrics(registry)

	for metricName, thresholdsDefinition := range job.Options.Thresholds {
		err := thresholdsDefinition.Parse()
		if err != nil {
			libOrch.HandleError(gs, errext.WithExitCodeIfNone(err, exitcodes.InvalidConfig))
			healthy = false
		}

		// TODO: move this to the orcj
		err = thresholdsDefinition.Validate(metricName, registry)
		if err != nil {
			libOrch.HandleError(gs, errext.WithExitCodeIfNone(err, exitcodes.InvalidConfig))
			healthy = false
		}
	}

	metricsEngine, err := engine.NewMetricsEngine(job.Options, gs.Logger(), registry, func() time.Duration {
		return gs.GetCurrentTestRunDuration()
	})
	if err != nil {
		libOrch.HandleError(gs, err)
		healthy = false
	}

	// Run the job

	result := "FAILURE"

	metricsEngine.Start()

	if healthy {
		result, err = handleExecution(gs, job, childJobs, registry)
		if err != nil {
			libOrch.HandleError(gs, err)
		}
	}

	metricsEngine.Stop()

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

func metricsStoreReceipt(gs libOrch.BaseGlobalState) (*primitive.ObjectID, error) {
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

func globeTestLogsStoreReceipt(gs libOrch.BaseGlobalState) (*primitive.ObjectID, error) {
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
