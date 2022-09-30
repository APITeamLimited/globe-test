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
	currentJobs map[string]map[string]string
	mutex       sync.Mutex
	maxJobs     int
***REMOVED***

func Run() ***REMOVED***
	ctx := context.Background()

	orchestratorId := uuid.New()

	// Orchestrator orchestratorClient deals with macro jos and connection to the rest of
	// APITEAM services
	orchestratorClient := redis.NewClient(&redis.Options***REMOVED***
		Addr:     fmt.Sprintf("%s:%s", libOrch.GetEnvVariable("ORCHESTRATOR_REDIS_HOST", "localhost"), libOrch.GetEnvVariable("ORCHESTRATOR_REDIS_PORT", "10000")),
		Username: "default",
		Password: libOrch.GetEnvVariable("ORCHESTRATOR_REDIS_PASSWORD", ""),
	***REMOVED***)

	scopesClient := redis.NewClient(&redis.Options***REMOVED***
		Addr:     fmt.Sprintf("%s:%s", libOrch.GetEnvVariable("CORE_CACHE_REDIS_HOST", "localhost"), libOrch.GetEnvVariable("CORE_CACHE_REDIS_PORT", "10001")),
		Username: "default",
		Password: libOrch.GetEnvVariable("CORE_CACHE_REDIS_PASSWORD", ""),
	***REMOVED***)

	storeMongoDB := getStoreMongoDB(ctx)

	workerClients := connectWorkerClients(ctx)

	currentTime := time.Now().UnixMilli()

	executionList := &ExecutionList***REMOVED***
		currentJobs: make(map[string]map[string]string),
		maxJobs:     -1,
	***REMOVED***

	//Set the orchestrator id and current time
	orchestratorClient.HSet(ctx, "orchestrators", orchestratorId.String(), currentTime)

	fmt.Print("\n\033[1;35mAPITEAM Orchestrator\033[0m\n\n")
	fmt.Printf("Starting orchestrator %s\n", orchestratorId.String())
	fmt.Printf("Listening for new jobs on %s...\n", orchestratorClient.Options().Addr)

	go checkForQueuedJobs(ctx, orchestratorClient, scopesClient, workerClients, orchestratorId.String(), executionList, storeMongoDB)

	// Subscribe to the execution channel
	psc := orchestratorClient.Subscribe(ctx, "orchestrator:execution")

	channel := psc.Channel()

	for msg := range channel ***REMOVED***
		jobId, err := uuid.Parse(msg.Payload)
		if err != nil ***REMOVED***
			fmt.Println("Error, got did not parse job id")
			return
		***REMOVED***
		go checkIfCanExecute(ctx, orchestratorClient, scopesClient, workerClients, jobId.String(), orchestratorId.String(), executionList, storeMongoDB)
	***REMOVED***
***REMOVED***

/*
Check for queued jobs that were deferered as they couldn't be executed when they
were queued as no workers were available.
*/
func checkForQueuedJobs(ctx context.Context, orchestratorClient, scopesClient *redis.Client, workerClients map[string]*redis.Client, orchestratorId string, executionList *ExecutionList, storeMongoDB *mongo.Database) ***REMOVED***
	// Check for job keys in the "orchestrator:executionHistory" set
	historyIds, err := orchestratorClient.SMembers(ctx, "orchestrator:executionHistory").Result()
	if err != nil ***REMOVED***
		fmt.Println("Error getting history ids", err)
	***REMOVED***

	for _, jobId := range historyIds ***REMOVED***
		go checkIfCanExecute(ctx, orchestratorClient, scopesClient, workerClients, jobId, orchestratorId, executionList, storeMongoDB)
	***REMOVED***
***REMOVED***

func checkIfCanExecute(ctx context.Context, orchestratorClient, scopesClient *redis.Client, workerClients map[string]*redis.Client, jobId string, orchestratorId string, executionList *ExecutionList, storeMongoDB *mongo.Database) ***REMOVED***
	// Try to HGetAll the orchestrator id
	job, err := orchestratorClient.HGetAll(ctx, jobId).Result()

	if err != nil ***REMOVED***
		fmt.Println("Error getting orchestrator")
		return
	***REMOVED***

	// TODO: check if has capacity to execute here

	// Check orchestrator['assignedOrchestrator'] is nil
	if _, ok := job["assignedOrchestrator"]; ok ***REMOVED***
		return
	***REMOVED***

	// Check if currently full execution list
	if !checkExecutionCapacity(executionList) ***REMOVED***
		return
	***REMOVED***

	// HSetNX assignedOrchestrator to the orchestratorId
	assignmentResult, err := orchestratorClient.HSetNX(ctx, jobId, "assignedOrchestrator", orchestratorId).Result()

	if err != nil ***REMOVED***
		fmt.Println("Error setting orchestrator")
		return
	***REMOVED***

	// If result is 0, orchestrator is already assigned
	if !assignmentResult ***REMOVED***
		return
	***REMOVED***

	// We got the job
	executionList.addJob(job)
	manageExecution(ctx, orchestratorClient, scopesClient, workerClients, job, orchestratorId, executionList, storeMongoDB)
***REMOVED***

// Over-arching function that manages the execution of a job and handles its state and lifecycle
// This is the highest level function with global state
func manageExecution(ctx context.Context, orchestratorClient, scopesClient *redis.Client, workerClients map[string]*redis.Client, job map[string]string, orchestratorId string, executionList *ExecutionList, storeMongoDB *mongo.Database) ***REMOVED***
	fmt.Println("Assigned job", job["id"])
	libOrch.UpdateStatus(ctx, orchestratorClient, job["id"], orchestratorId, "ASSIGNED")

	// Setup the job

	healthy := true

	gs := NewGlobalState(ctx, orchestratorClient, job["id"], orchestratorId)

	options, err := options.DetermineRuntimeOptions(job, gs)
	if err != nil ***REMOVED***
		libOrch.HandleStringError(ctx, orchestratorClient, job["id"], orchestratorId, fmt.Sprintf("Error determining runtime options: %s", err.Error()))
		healthy = false
	***REMOVED***

	marshalledOptions, _ := json.Marshal(options)
	if err != nil ***REMOVED***
		libOrch.HandleStringError(ctx, orchestratorClient, job["id"], orchestratorId, fmt.Sprintf("Error marshalling runtime options: %s", err.Error()))
		healthy = false
	***REMOVED***

	libOrch.DispatchMessage(ctx, orchestratorClient, job["id"], orchestratorId, string(marshalledOptions), "OPTIONS")

	scope, err := fetchScope(ctx, scopesClient, job["scopeId"])
	if err != nil ***REMOVED***
		libOrch.HandleError(ctx, orchestratorClient, job["id"], orchestratorId, err)
		healthy = false
	***REMOVED***

	// Running the job

	result := "FAILED"

	if healthy ***REMOVED***
		result, err = runExecution(gs, options, scope, workerClients, job)
		if err != nil ***REMOVED***
			libOrch.HandleError(ctx, orchestratorClient, job["id"], orchestratorId, err)
		***REMOVED***
	***REMOVED***

	libOrch.UpdateStatus(ctx, orchestratorClient, job["id"], orchestratorId, result)

	// Storing and cleaning up

	(*gs.MetricsStore()).Stop()

	// Create globe test log id, note this must be sent after cleanup
	globeTestLogsId := primitive.NewObjectID()
	globeTestLogsIdMessage := &libOrch.MarkMessage***REMOVED***
		Mark:    "GlobeTestLogsStoreReceipt",
		Message: globeTestLogsId.Hex(),
	***REMOVED***

	marshalledLogs, err := json.Marshal(globeTestLogsIdMessage)
	if err != nil ***REMOVED***
		libOrch.HandleError(ctx, orchestratorClient, job["id"], orchestratorId, err)
		return
	***REMOVED***
	libOrch.DispatchMessage(ctx, orchestratorClient, job["id"], orchestratorId, string(marshalledLogs), "MARK")

	// Temporary object storing map[string][]map[string]string, the job in production
	// should be separate to allow for parallel child jobs
	childJobs := make(map[string][]map[string]string)
	childJobs["portsmouth"] = append(childJobs["portsmouth"], job)

	// Clean up the job and store result in Mongo
	cleanup(ctx, job, childJobs, orchestratorClient, workerClients, orchestratorId, storeMongoDB, scope, globeTestLogsId)

	libOrch.UpdateStatus(ctx, orchestratorClient, job["id"], orchestratorId, fmt.Sprintf("COMPLETED_%s", result))

	executionList.removeJob(job["id"])

	// Capacity was freed, so check for queued jobs
	checkForQueuedJobs(ctx, orchestratorClient, scopesClient, workerClients, orchestratorId, executionList, storeMongoDB)
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
