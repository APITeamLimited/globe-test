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

func Run() ***REMOVED***
	ctx := context.Background()
	orchestratorId := uuid.NewString()
	orchestratorClient := getOrchestratorClient()
	storeMongoDB := getStoreMongoDB(ctx)
	workerClients := connectWorkerClients(ctx)
	maxJobs := getMaxJobs()
	maxManagedVUs := getMaxManagedVUs()

	executionList := &ExecutionList***REMOVED***
		currentJobs:   make(map[string]libOrch.Job),
		maxJobs:       maxJobs,
		maxManagedVUs: maxManagedVUs,
	***REMOVED***

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

	for msg := range channel ***REMOVED***
		jobId, err := uuid.Parse(msg.Payload)
		if err != nil ***REMOVED***
			fmt.Println("Error, got did not parse job id")
			return
		***REMOVED***
		go checkIfCanExecute(ctx, orchestratorClient, workerClients, jobId.String(), orchestratorId, executionList, storeMongoDB)
	***REMOVED***
***REMOVED***

// Check for queued jobs that were deferered as they couldn't be executed when they
// were queued as no workers were available.
func checkForQueuedJobs(ctx context.Context, orchestratorClient *redis.Client, workerClients libOrch.WorkerClients, orchestratorId string, executionList *ExecutionList, storeMongoDB *mongo.Database) ***REMOVED***
	// Check for job keys in the "orchestrator:executionHistory" set
	historyIds, err := orchestratorClient.SMembers(ctx, "orchestrator:executionHistory").Result()
	if err != nil ***REMOVED***
		fmt.Println("Error getting history ids", err)
	***REMOVED***

	for _, jobId := range historyIds ***REMOVED***
		go checkIfCanExecute(ctx, orchestratorClient, workerClients, jobId, orchestratorId, executionList, storeMongoDB)
	***REMOVED***
***REMOVED***

// Ensures job has no already been assigned and determines if this node has capacity to execute
func checkIfCanExecute(ctx context.Context, orchestratorClient *redis.Client, workerClients libOrch.WorkerClients, jobId string, orchestratorId string, executionList *ExecutionList, storeMongoDB *mongo.Database) ***REMOVED***
	// Try to HGetAll the orchestrator id
	job, err := fetchJob(ctx, orchestratorClient, jobId)
	if err != nil || job == nil ***REMOVED***
		if err != nil ***REMOVED***
			fmt.Println("Error getting job")
		***REMOVED***
		return
	***REMOVED***

	executionList.mutex.Lock()

	// Don't even bother calculating options if we don't have capacity
	if !executionList.checkExecutionCapacity(nil) ***REMOVED***
		return
	***REMOVED***

	if value, _ := orchestratorClient.HGet(ctx, job.Id, "assignedOrchestrator").Result(); value != "" ***REMOVED***
		// If the job has been assigned to another orchestrator, return

		executionList.mutex.Unlock()
		return
	***REMOVED***

	executionList.mutex.Unlock()

	gs := NewGlobalState(ctx, orchestratorClient, job.Id, orchestratorId)
	options, optionsErr := options.DetermineRuntimeOptions(*job, gs, workerClients)
	job.Options = options

	// Check execution capacity again, bearing in mind options
	executionList.mutex.Lock()
	if !executionList.checkExecutionCapacity(options) ***REMOVED***
		executionList.mutex.Unlock()
		return
	***REMOVED***

	// HSetNX assignedOrchestrator to the orchestratorId
	assignmentResult, err := orchestratorClient.HSetNX(ctx, job.Id, "assignedOrchestrator", orchestratorId).Result()

	// If result is 0, orchestrator is already assigned
	if !assignmentResult ***REMOVED***
		executionList.mutex.Unlock()
		return
	***REMOVED***

	if err != nil ***REMOVED***
		fmt.Println("Error setting orchestrator")
		executionList.mutex.Unlock()
		return
	***REMOVED***

	// We got the job and have confirmed capacity for it

	job.AssignedOrchestrator = orchestratorId

	executionList.addJob(job)
	executionList.mutex.Unlock()

	manageExecution(gs, orchestratorClient, workerClients, *job, orchestratorId, executionList, storeMongoDB, optionsErr)
***REMOVED***

type jobDistribution struct ***REMOVED***
	jobs         []libOrch.ChildJob
	workerClient *redis.Client
***REMOVED***

// Over-arching function that manages the execution of a job and handles its state and lifecycle
// This is the highest level function with global state
func manageExecution(gs *globalState, orchestratorClient *redis.Client, workerClients libOrch.WorkerClients, job libOrch.Job,
	orchestratorId string, executionList *ExecutionList, storeMongoDB *mongo.Database, optionsErr error) ***REMOVED***
	// Get the job id and check if it is a string
	fmt.Println("Assigned job", job.Id)
	libOrch.UpdateStatus(gs, "ASSIGNED")

	// Setup the job

	healthy := optionsErr == nil

	options, err := options.DetermineRuntimeOptions(job, gs, workerClients)
	if err != nil ***REMOVED***
		libOrch.HandleStringError(gs, fmt.Sprintf("Error determining runtime options: %s", err.Error()))
		healthy = false
	***REMOVED***

	if healthy ***REMOVED***
		marshalledOptions, err := json.Marshal(options)
		if err != nil ***REMOVED***
			libOrch.HandleStringError(gs, fmt.Sprintf("Error marshalling runtime options: %s", err.Error()))
			healthy = false
		***REMOVED***

		libOrch.DispatchMessage(gs, string(marshalledOptions), "OPTIONS")
	***REMOVED***

	scope := job.Scope

	childJobs, err := determineChildJobs(healthy, job, options, workerClients)
	if err != nil ***REMOVED***
		libOrch.HandleError(gs, err)
		healthy = false
	***REMOVED***

	// Run the job

	result := "FAILURE"

	if healthy ***REMOVED***
		result, err = handleExecution(gs, options, scope, childJobs, job.Id)
		if err != nil ***REMOVED***
			fmt.Println("Error running execution", err)
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
		return
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
		return
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

	if healthy ***REMOVED***
		executionList.removeJob(job.Id)
	***REMOVED***
***REMOVED***
