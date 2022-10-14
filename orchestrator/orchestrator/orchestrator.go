package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/orchestrator/options"
	"github.com/APITeamLimited/redis/v9"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ExecutionList struct ***REMOVED***
	currentJobs map[string]libOrch.Job
	mutex       sync.Mutex
	maxJobs     int
***REMOVED***

func Run() ***REMOVED***
	ctx := context.Background()

	orchestratorId := uuid.NewString()

	// Orchestrator orchestratorClient deals with macro jos and connection to the rest of
	// APITEAM services
	orchestratorClient := redis.NewClient(&redis.Options***REMOVED***
		Addr:     fmt.Sprintf("%s:%s", libOrch.GetEnvVariable("ORCHESTRATOR_REDIS_HOST", "localhost"), libOrch.GetEnvVariable("ORCHESTRATOR_REDIS_PORT", "10000")),
		Username: "default",
		Password: libOrch.GetEnvVariable("ORCHESTRATOR_REDIS_PASSWORD", ""),
	***REMOVED***)

	storeMongoDB := getStoreMongoDB(ctx)

	workerClients := connectWorkerClients(ctx)

	// Set the orchestrator id and the start time in the index
	orchestratorClient.SAdd(ctx, "orchestrators", orchestratorId)

	// Every second set a heartbeat update
	heartbeatTicker := time.NewTicker(1 * time.Second)
	go func() ***REMOVED***
		for range heartbeatTicker.C ***REMOVED***
			orchestratorClient.Set(ctx, fmt.Sprintf("orchestrator:%s:lastHeartbeat", orchestratorId), time.Now().UnixMilli(), time.Second*10)
		***REMOVED***
	***REMOVED***()

	fmt.Print("\n\033[1;35mAPITEAM Orchestrator\033[0m\n\n")
	fmt.Printf("Starting orchestrator %s\n", orchestratorId)
	fmt.Printf("Listening for new jobs on %s...\n\n", orchestratorClient.Options().Addr)

	executionList := &ExecutionList***REMOVED***
		currentJobs: make(map[string]libOrch.Job),
		maxJobs:     -1,
	***REMOVED***

	go checkForQueuedJobs(ctx, orchestratorClient, workerClients, orchestratorId, executionList, storeMongoDB)

	// Subscribe to the execution channel
	pubSub := orchestratorClient.Subscribe(ctx, "orchestrator:execution")

	channel := pubSub.Channel()

	for msg := range channel ***REMOVED***
		jobId, err := uuid.Parse(msg.Payload)
		if err != nil ***REMOVED***
			fmt.Println("Error, got did not parse job id")
			return
		***REMOVED***
		go checkIfCanExecute(ctx, orchestratorClient, workerClients, jobId.String(), orchestratorId, executionList, storeMongoDB)
	***REMOVED***
***REMOVED***

/*
Check for queued jobs that were deferered as they couldn't be executed when they
were queued as no workers were available.
*/
func checkForQueuedJobs(ctx context.Context, orchestratorClient *redis.Client, workerClients map[string]*redis.Client, orchestratorId string, executionList *ExecutionList, storeMongoDB *mongo.Database) ***REMOVED***
	// Check for job keys in the "orchestrator:executionHistory" set
	historyIds, err := orchestratorClient.SMembers(ctx, "orchestrator:executionHistory").Result()
	if err != nil ***REMOVED***
		fmt.Println("Error getting history ids", err)
	***REMOVED***

	for _, jobId := range historyIds ***REMOVED***
		go checkIfCanExecute(ctx, orchestratorClient, workerClients, jobId, orchestratorId, executionList, storeMongoDB)
	***REMOVED***
***REMOVED***

func checkIfCanExecute(ctx context.Context, orchestratorClient *redis.Client, workerClients map[string]*redis.Client, jobId string, orchestratorId string, executionList *ExecutionList, storeMongoDB *mongo.Database) ***REMOVED***
	// Try to HGetAll the orchestrator id
	job, err := fetchJob(ctx, orchestratorClient, jobId)
	if err != nil || job == nil ***REMOVED***
		fmt.Println("Error getting job")
		return
	***REMOVED***

	// Check if orchestrtator has already been assigned
	if job.AssignedOrchestrator != "" ***REMOVED***
		return
	***REMOVED***

	// TODO: check if has capacity to execute here

	// Check if currently full execution list
	if !checkExecutionCapacity(executionList) ***REMOVED***
		return
	***REMOVED***

	// HSetNX assignedOrchestrator to the orchestratorId
	assignmentResult, err := orchestratorClient.HSetNX(ctx, job.Id, "assignedOrchestrator", orchestratorId).Result()

	if err != nil ***REMOVED***
		fmt.Println("Error setting orchestrator")
		return
	***REMOVED***

	// If result is 0, orchestrator is already assigned
	if !assignmentResult ***REMOVED***
		return
	***REMOVED***

	// We got the job
	job.AssignedOrchestrator = orchestratorId
	executionList.addJob(*job)
	manageExecution(ctx, orchestratorClient, workerClients, *job, orchestratorId, executionList, storeMongoDB)

	// Capacity was freed, so check for queued jobs
	checkForQueuedJobs(ctx, orchestratorClient, workerClients, orchestratorId, executionList, storeMongoDB)
***REMOVED***

type jobDistribution struct ***REMOVED***
	jobs         *[]libOrch.ChildJob
	workerClient *redis.Client
***REMOVED***

// Over-arching function that manages the execution of a job and handles its state and lifecycle
// This is the highest level function with global state
func manageExecution(ctx context.Context, orchestratorClient *redis.Client, workerClients map[string]*redis.Client, job libOrch.Job, orchestratorId string, executionList *ExecutionList, storeMongoDB *mongo.Database) ***REMOVED***
	// Get the job id and check if it is a string
	fmt.Println("Assigned job", job.Id)
	libOrch.UpdateStatus(ctx, orchestratorClient, job.Id, orchestratorId, "ASSIGNED")

	// Setup the job

	healthy := true
	gs := NewGlobalState(ctx, orchestratorClient, job.Id, orchestratorId)

	options, err := options.DetermineRuntimeOptions(job, gs)
	if err != nil ***REMOVED***
		libOrch.HandleStringError(ctx, orchestratorClient, job.Id, orchestratorId, fmt.Sprintf("Error determining runtime options: %s", err.Error()))
		healthy = false
	***REMOVED***

	if healthy ***REMOVED***
		marshalledOptions, err := json.Marshal(options)
		if err != nil ***REMOVED***
			libOrch.HandleStringError(ctx, orchestratorClient, job.Id, orchestratorId, fmt.Sprintf("Error marshalling runtime options: %s", err.Error()))
			healthy = false
		***REMOVED***

		libOrch.DispatchMessage(ctx, orchestratorClient, job.Id, orchestratorId, string(marshalledOptions), "OPTIONS")
	***REMOVED***

	scope := job.Scope

	childJobs := make(map[string]jobDistribution)

	if healthy ***REMOVED***
		childJob := libOrch.ChildJob***REMOVED***
			Job:               job,
			ChildJobId:        uuid.NewString(),
			Options:           *options,
			UnderlyingRequest: job.UnderlyingRequest,
			FinalRequest:      job.FinalRequest,
		***REMOVED***

		childJobs["portsmouth"] = jobDistribution***REMOVED***
			jobs:         &[]libOrch.ChildJob***REMOVED***childJob***REMOVED***,
			workerClient: workerClients["portsmouth"],
		***REMOVED***
	***REMOVED***

	// Run the job

	result := "FAILURE"

	if healthy ***REMOVED***
		result, err = runExecution(gs, options, scope, childJobs, job.Id)
		if err != nil ***REMOVED***
			fmt.Println("Error running execution", err)
			libOrch.HandleError(ctx, orchestratorClient, job.Id, orchestratorId, err)
		***REMOVED***
	***REMOVED***

	libOrch.UpdateStatus(ctx, orchestratorClient, job.Id, orchestratorId, result)

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
		libOrch.HandleError(ctx, orchestratorClient, job.Id, orchestratorId, err)
		return
	***REMOVED***
	libOrch.DispatchMessage(ctx, orchestratorClient, job.Id, orchestratorId, string(marshalledGlobeTestReceipt), "MARK")

	//Create Metrics Store receipt, note this must be sent after cleanup
	metricsStoreReceipt := primitive.NewObjectID()
	metricsStoreReceiptMessage := &libOrch.MarkMessage***REMOVED***
		Mark:    "MetricsStoreReceipt",
		Message: metricsStoreReceipt.Hex(),
	***REMOVED***

	marshalledMetricsStoreReceipt, err := json.Marshal(metricsStoreReceiptMessage)
	if err != nil ***REMOVED***
		fmt.Println("Error marshalling metrics store receipt", err)
		libOrch.HandleError(ctx, orchestratorClient, job.Id, orchestratorId, err)
		return
	***REMOVED***
	libOrch.DispatchMessage(ctx, orchestratorClient, job.Id, orchestratorId, string(marshalledMetricsStoreReceipt), "MARK")

	// Clean up the job and store result in Mongo
	err = cleanup(ctx, job, childJobs, orchestratorClient, orchestratorId, storeMongoDB, scope, globeTestLogsReceipt, metricsStoreReceipt)
	if err != nil ***REMOVED***
		fmt.Println("Error cleaning up", err)
		libOrch.HandleErrorNoSet(ctx, orchestratorClient, job.Id, orchestratorId, err)
		libOrch.UpdateStatusNoSet(ctx, orchestratorClient, job.Id, orchestratorId, result)
	***REMOVED*** else ***REMOVED***
		libOrch.UpdateStatusNoSet(ctx, orchestratorClient, job.Id, orchestratorId, fmt.Sprintf("COMPLETED_%s", result))
	***REMOVED***

	executionList.removeJob(job.Id)
***REMOVED***

/*
Check if the exectutor has the physical capacity to execute this job, this does
not concern whether the user has the required credits to execute the job.
*/
func checkExecutionCapacity(executionList *ExecutionList) bool ***REMOVED***
	// TODO: check if has capacity to execute here

	// If more than max jobs, return false
	if executionList.maxJobs >= 0 && len(executionList.currentJobs) >= executionList.maxJobs ***REMOVED***
		return false
	***REMOVED***

	// TODO: implement more capacity checks

	return true
***REMOVED***
