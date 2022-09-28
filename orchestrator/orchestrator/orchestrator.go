package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/APITeamLimited/redis/v9"
	"github.com/google/uuid"
	"gitlab.com/apiteamcloud/orchestrator/lib"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ExecutionList struct {
	currentJobs map[string]map[string]string
	mutex       sync.Mutex
	maxJobs     int
}

func Run() {
	ctx := context.Background()

	orchestratorId := uuid.New()

	// Orchestrator orchestratorClient deals with macro jos and connection to the rest of
	// APITEAM services
	orchestratorClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", lib.GetEnvVariable("ORCHESTRATOR_REDIS_HOST", "localhost"), lib.GetEnvVariable("ORCHESTRATOR_REDIS_PORT", "10000")),
		Username: "default",
		Password: lib.GetEnvVariable("ORCHESTRATOR_REDIS_PASSWORD", ""),
	})

	scopesClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", lib.GetEnvVariable("CORE_CACHE_REDIS_HOST", "localhost"), lib.GetEnvVariable("CORE_CACHE_REDIS_PORT", "10001")),
		Username: "default",
		Password: lib.GetEnvVariable("CORE_CACHE_REDIS_PASSWORD", ""),
	})

	storeMongoDB := getStoreMongoDB(ctx)

	workerClients := connectWorkerClients(ctx)

	currentTime := time.Now().UnixMilli()

	executionList := &ExecutionList{
		currentJobs: make(map[string]map[string]string),
		maxJobs:     -1,
	}

	//Set the orchestrator id and current time
	orchestratorClient.HSet(ctx, "orchestrators", orchestratorId.String(), currentTime)

	fmt.Print("\n\033[1;35mAPITEAM Orchestrator\033[0m\n\n")
	fmt.Printf("Starting orchestrator %s\n", orchestratorId.String())
	fmt.Printf("Listening for new jobs on %s...\n", orchestratorClient.Options().Addr)

	go checkForQueuedJobs(ctx, orchestratorClient, scopesClient, workerClients, orchestratorId.String(), executionList, storeMongoDB)

	// Subscribe to the execution channel
	psc := orchestratorClient.Subscribe(ctx, "orchestrator:execution")

	channel := psc.Channel()

	for msg := range channel {
		jobId, err := uuid.Parse(msg.Payload)
		if err != nil {
			fmt.Println("Error, got did not parse job id")
			return
		}
		go checkIfCanExecute(ctx, orchestratorClient, scopesClient, workerClients, jobId.String(), orchestratorId.String(), executionList, storeMongoDB)
	}
}

/*
Check for queued jobs that were deferered as they couldn't be executed when they
were queued as no workers were available.
*/
func checkForQueuedJobs(ctx context.Context, orchestratorClient, scopesClient *redis.Client, workerClients map[string]*redis.Client, orchestratorId string, executionList *ExecutionList, storeMongoDB *mongo.Database) {
	// Check for job keys in the "orchestrator:executionHistory" set
	historyIds, err := orchestratorClient.SMembers(ctx, "orchestrator:executionHistory").Result()
	if err != nil {
		fmt.Println("Error getting history ids", err)
	}

	for _, jobId := range historyIds {
		go checkIfCanExecute(ctx, orchestratorClient, scopesClient, workerClients, jobId, orchestratorId, executionList, storeMongoDB)
	}
}

func checkIfCanExecute(ctx context.Context, orchestratorClient, scopesClient *redis.Client, workerClients map[string]*redis.Client, jobId string, orchestratorId string, executionList *ExecutionList, storeMongoDB *mongo.Database) {
	// Try to HGetAll the orchestrator id
	job, err := orchestratorClient.HGetAll(ctx, jobId).Result()

	if err != nil {
		fmt.Println("Error getting orchestrator")
		return
	}

	// TODO: check if has capacity to execute here

	// Check orchestrator['assignedOrchestrator'] is nil
	if _, ok := job["assignedOrchestrator"]; ok {
		return
	}

	// Check if currently full execution list
	if !checkExecutionCapacity(executionList) {
		return
	}

	// HSetNX assignedOrchestrator to the orchestratorId
	assignmentResult, err := orchestratorClient.HSetNX(ctx, jobId, "assignedOrchestrator", orchestratorId).Result()

	if err != nil {
		fmt.Println("Error setting orchestrator")
		return
	}

	// If result is 0, orchestrator is already assigned
	if !assignmentResult {
		return
	}

	// We got the job
	executionList.addJob(job)
	manageExecution(ctx, orchestratorClient, scopesClient, workerClients, job, orchestratorId, executionList, storeMongoDB)
}

func manageExecution(ctx context.Context, orchestratorClient, scopesClient *redis.Client, workerClients map[string]*redis.Client, job map[string]string, orchestratorId string, executionList *ExecutionList, storeMongoDB *mongo.Database) {
	fmt.Println("Assigned job", job["id"])

	gs := lib.NewGlobalState(ctx, orchestratorClient, job["id"], orchestratorId)

	UpdateStatus(ctx, orchestratorClient, job["id"], orchestratorId, "ASSIGNED")
	run(gs, orchestratorId, orchestratorClient, scopesClient, workerClients, job, storeMongoDB)

	// Create a new objectId for globeTestLogs
	globeTestLogsId := primitive.NewObjectID()

	globeTestLogsIdMessage := &markMessage{
		Mark:    "GlobeTestLogsStoreReceipt",
		Message: globeTestLogsId.Hex(),
	}

	marshalledLogs, err := json.Marshal(globeTestLogsIdMessage)
	if err != nil {
		HandleError(ctx, orchestratorClient, job["id"], orchestratorId, err)
		return
	}

	DispatchMessage(ctx, orchestratorClient, job["id"], orchestratorId, string(marshalledLogs), "MARK")

	// Temporary object storing map[string][]map[string]string, the job in production
	// should be separate to allow for parallel child jobs
	childJobs := make(map[string][]map[string]string)
	childJobs["portsmouth"] = append(childJobs["portsmouth"], job)

	scope, err := scopesClient.Get(ctx, fmt.Sprintf("scope__id:%s", job["scopeId"])).Result()
	if err != nil {
		fmt.Println("Error getting scope", err)
		HandleStringError(ctx, orchestratorClient, job["id"], orchestratorId, fmt.Sprintf("Error getting scope %s", job["scopeId"]))
		return
	}

	// Check scope not empty
	if len(scope) == 0 {
		HandleStringError(ctx, orchestratorClient, job["id"], orchestratorId, fmt.Sprintf("Scope %s is empty", job["scopeId"]))
		return
	}

	// Parse scope as map[string]string
	scopeMap := make(map[string]string)
	err = json.Unmarshal([]byte(scope), &scopeMap)
	if err != nil {
		fmt.Println("Error unmarshalling scope", err)
		HandleStringError(ctx, orchestratorClient, job["id"], orchestratorId, fmt.Sprintf("Error unmarshalling scope %s", job["scopeId"]))
		return
	}

	// Clean up the job and store result in Mongo
	cleanup(ctx, job, childJobs, orchestratorClient, workerClients, orchestratorId, storeMongoDB, scopeMap, globeTestLogsId)

	executionList.removeJob(job["id"])
	// Capacity was freed, so check for queued jobs
	checkForQueuedJobs(ctx, orchestratorClient, scopesClient, workerClients, orchestratorId, executionList, storeMongoDB)
}

/*
Check if the exectutor has the physical capacity to execute this job, this does
not concern whether the user has the required credits to execute the job.
*/
func checkExecutionCapacity(executionList *ExecutionList) bool {
	// TODO: check if has capacity to execute here

	// If more than max jobs, return false
	if executionList.maxJobs >= 0 && len(executionList.currentJobs) >= executionList.maxJobs {
		return false
	}

	// TODO: implement more capacity checks

	return true
}
