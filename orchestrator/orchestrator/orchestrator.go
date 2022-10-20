package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/orchestrator/options"
	"github.com/APITeamLimited/redis/v9"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Run() {
	ctx := context.Background()
	orchestratorId := uuid.NewString()
	orchestratorClient := getOrchestratorClient()
	storeMongoDB := getStoreMongoDB(ctx)
	workerClients := connectWorkerClients(ctx)
	maxJobs := getMaxJobs()
	maxManagedVUs := getMaxManagedVUs()

	executionList := &ExecutionList{
		currentJobs:   make(map[string]libOrch.Job),
		maxJobs:       maxJobs,
		maxManagedVUs: maxManagedVUs,
	}

	// Create a scheduler for regular updates and checks
	startJobScheduling(ctx, orchestratorClient, orchestratorId, executionList, workerClients, storeMongoDB)

	// Periodically check for and delete offline orchestrators
	createDeletionScheduler(ctx, orchestratorClient, workerClients)

	// Change process title
	fmt.Printf("\033]0;GlobeTest Orchestrator: %s\007", orchestratorId)

	fmt.Print("\n\033[1;35mGlobeTest Orchestrator\033[0m\n\n")
	fmt.Printf("Starting orchestrator %s\n", orchestratorId)
	fmt.Printf("Listening for new jobs on %s...\n\n", orchestratorClient.Options().Addr)

	// Subscribe to the execution channel and listen for new jobs
	channel := orchestratorClient.Subscribe(ctx, "orchestrator:execution").Channel()

	for msg := range channel {
		jobId, err := uuid.Parse(msg.Payload)
		if err != nil {
			fmt.Println("Error, got did not parse job id")
			return
		}
		go checkIfCanExecute(ctx, orchestratorClient, workerClients, jobId.String(), orchestratorId, executionList, storeMongoDB)
	}
}

// Check for queued jobs that were deferered as they couldn't be executed when they
// were queued as no workers were available.
func checkForQueuedJobs(ctx context.Context, orchestratorClient *redis.Client, workerClients libOrch.WorkerClients, orchestratorId string, executionList *ExecutionList, storeMongoDB *mongo.Database) {
	// Check for job keys in the "orchestrator:executionHistory" set
	historyIds, err := orchestratorClient.SMembers(ctx, "orchestrator:executionHistory").Result()
	if err != nil {
		fmt.Println("Error getting history ids", err)
	}

	for _, jobId := range historyIds {
		go checkIfCanExecute(ctx, orchestratorClient, workerClients, jobId, orchestratorId, executionList, storeMongoDB)
	}
}

// Ensures job has no already been assigned and determines if this node has capacity to execute
func checkIfCanExecute(ctx context.Context, orchestratorClient *redis.Client, workerClients libOrch.WorkerClients, jobId string, orchestratorId string, executionList *ExecutionList, storeMongoDB *mongo.Database) {
	// Try to HGetAll the orchestrator id
	job, err := fetchJob(ctx, orchestratorClient, jobId)
	if err != nil || job == nil {
		if err != nil {
			fmt.Println("Error getting job")
		}
		return
	}

	executionList.mutex.Lock()

	// Don't even bother calculating options if we don't have capacity
	if !executionList.checkExecutionCapacity(nil) {
		return
	}

	if value, _ := orchestratorClient.HGet(ctx, job.Id, "assignedOrchestrator").Result(); value != "" {
		// If the job has been assigned to another orchestrator, return

		executionList.mutex.Unlock()
		return
	}

	executionList.mutex.Unlock()

	gs := NewGlobalState(ctx, orchestratorClient, job.Id, orchestratorId)
	options, optionsErr := options.DetermineRuntimeOptions(*job, gs, workerClients)
	job.Options = options

	// Check execution capacity again, bearing in mind options
	executionList.mutex.Lock()
	if !executionList.checkExecutionCapacity(options) {
		executionList.mutex.Unlock()
		return
	}

	// HSetNX assignedOrchestrator to the orchestratorId
	assignmentResult, err := orchestratorClient.HSetNX(ctx, job.Id, "assignedOrchestrator", orchestratorId).Result()

	// If result is 0, orchestrator is already assigned
	if !assignmentResult {
		executionList.mutex.Unlock()
		return
	}

	if err != nil {
		fmt.Println("Error setting orchestrator")
		executionList.mutex.Unlock()
		return
	}

	// We got the job and have confirmed capacity for it

	job.AssignedOrchestrator = orchestratorId

	executionList.addJob(job)
	executionList.mutex.Unlock()

	manageExecution(gs, orchestratorClient, workerClients, *job, orchestratorId, executionList, storeMongoDB, optionsErr)
}

type jobDistribution struct {
	jobs         []libOrch.ChildJob
	workerClient *redis.Client
}

// Over-arching function that manages the execution of a job and handles its state and lifecycle
// This is the highest level function with global state
func manageExecution(gs *globalState, orchestratorClient *redis.Client, workerClients libOrch.WorkerClients, job libOrch.Job,
	orchestratorId string, executionList *ExecutionList, storeMongoDB *mongo.Database, optionsErr error) {
	// Get the job id and check if it is a string
	fmt.Println("Assigned job", job.Id)
	libOrch.UpdateStatus(gs, "ASSIGNED")

	// Setup the job

	healthy := optionsErr == nil

	options, err := options.DetermineRuntimeOptions(job, gs, workerClients)
	if err != nil {
		libOrch.HandleStringError(gs, fmt.Sprintf("Error determining runtime options: %s", err.Error()))
		healthy = false
	}

	if healthy {
		marshalledOptions, err := json.Marshal(options)
		if err != nil {
			libOrch.HandleStringError(gs, fmt.Sprintf("Error marshalling runtime options: %s", err.Error()))
			healthy = false
		}

		libOrch.DispatchMessage(gs, string(marshalledOptions), "OPTIONS")
	}

	scope := job.Scope

	childJobs, err := determineChildJobs(healthy, job, options, workerClients)
	if err != nil {
		libOrch.HandleError(gs, err)
		healthy = false
	}

	// Run the job

	result := "FAILURE"

	if healthy {
		result, err = handleExecution(gs, options, scope, childJobs, job.Id)
		if err != nil {
			fmt.Println("Error running execution", err)
			libOrch.HandleError(gs, err)
		}
	}

	libOrch.UpdateStatus(gs, result)

	// Storing and cleaning up

	(*gs.MetricsStore()).Stop()

	// Create GlobeTest logs store receipt, note this must be sent after cleanup
	globeTestLogsReceipt := primitive.NewObjectID()
	globeTestLogsReceiptMessage := &libOrch.MarkMessage{
		Mark:    "GlobeTestLogsStoreReceipt",
		Message: globeTestLogsReceipt.Hex(),
	}

	marshalledGlobeTestReceipt, err := json.Marshal(globeTestLogsReceiptMessage)
	if err != nil {
		fmt.Println("Error marshalling GlobeTestLogsStoreReceipt", err)
		libOrch.HandleError(gs, err)
		return
	}
	libOrch.DispatchMessage(gs, string(marshalledGlobeTestReceipt), "MARK")

	//Create Metrics Store receipt, note this must be sent after cleanup
	metricsStoreReceipt := primitive.NewObjectID()
	metricsStoreReceiptMessage := &libOrch.MarkMessage{
		Mark:    "MetricsStoreReceipt",
		Message: metricsStoreReceipt.Hex(),
	}

	marshalledMetricsStoreReceipt, err := json.Marshal(metricsStoreReceiptMessage)
	if err != nil {
		fmt.Println("Error marshalling metrics store receipt", err)
		libOrch.HandleError(gs, err)
		return
	}
	libOrch.DispatchMessage(gs, string(marshalledMetricsStoreReceipt), "MARK")

	// Clean up the job and store result in Mongo
	err = cleanup(gs, job, childJobs, storeMongoDB, scope, globeTestLogsReceipt, metricsStoreReceipt)
	if err != nil {
		fmt.Println("Error cleaning up", err)
		libOrch.HandleErrorNoSet(gs, err)
		libOrch.UpdateStatusNoSet(gs, result)
	} else {
		libOrch.UpdateStatusNoSet(gs, fmt.Sprintf("COMPLETED_%s", result))
	}

	if healthy {
		executionList.removeJob(job.Id)
	}
}
