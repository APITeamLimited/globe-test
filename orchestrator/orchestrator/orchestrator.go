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

type ExecutionList struct {
	currentJobs map[string]libOrch.Job
	mutex       sync.Mutex
	maxJobs     int
}

func Run() {
	ctx := context.Background()

	orchestratorId := uuid.NewString()

	// Orchestrator orchestratorClient deals with macro jos and connection to the rest of
	// APITEAM services
	orchestratorClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", libOrch.GetEnvVariable("ORCHESTRATOR_REDIS_HOST", "localhost"), libOrch.GetEnvVariable("ORCHESTRATOR_REDIS_PORT", "10000")),
		Username: "default",
		Password: libOrch.GetEnvVariable("ORCHESTRATOR_REDIS_PASSWORD", ""),
	})

	storeMongoDB := getStoreMongoDB(ctx)

	workerClients := connectWorkerClients(ctx)

	// Set the orchestrator id and the start time in the index
	orchestratorClient.SAdd(ctx, "orchestrators", orchestratorId)

	// Every second set a heartbeat update
	heartbeatTicker := time.NewTicker(1 * time.Second)
	go func() {
		for range heartbeatTicker.C {
			orchestratorClient.Set(ctx, fmt.Sprintf("orchestrator:%s:lastHeartbeat", orchestratorId), time.Now().UnixMilli(), time.Second*10)
		}
	}()

	fmt.Print("\n\033[1;35mAPITEAM Orchestrator\033[0m\n\n")
	fmt.Printf("Starting orchestrator %s\n", orchestratorId)
	fmt.Printf("Listening for new jobs on %s...\n\n", orchestratorClient.Options().Addr)

	executionList := &ExecutionList{
		currentJobs: make(map[string]libOrch.Job),
		maxJobs:     -1,
	}

	go checkForQueuedJobs(ctx, orchestratorClient, workerClients, orchestratorId, executionList, storeMongoDB)

	// Subscribe to the execution channel
	pubSub := orchestratorClient.Subscribe(ctx, "orchestrator:execution")

	channel := pubSub.Channel()

	for msg := range channel {
		jobId, err := uuid.Parse(msg.Payload)
		if err != nil {
			fmt.Println("Error, got did not parse job id")
			return
		}
		go checkIfCanExecute(ctx, orchestratorClient, workerClients, jobId.String(), orchestratorId, executionList, storeMongoDB)
	}
}

/*
Check for queued jobs that were deferered as they couldn't be executed when they
were queued as no workers were available.
*/
func checkForQueuedJobs(ctx context.Context, orchestratorClient *redis.Client, workerClients map[string]*redis.Client, orchestratorId string, executionList *ExecutionList, storeMongoDB *mongo.Database) {
	// Check for job keys in the "orchestrator:executionHistory" set
	historyIds, err := orchestratorClient.SMembers(ctx, "orchestrator:executionHistory").Result()
	if err != nil {
		fmt.Println("Error getting history ids", err)
	}

	for _, jobId := range historyIds {
		go checkIfCanExecute(ctx, orchestratorClient, workerClients, jobId, orchestratorId, executionList, storeMongoDB)
	}
}

func checkIfCanExecute(ctx context.Context, orchestratorClient *redis.Client, workerClients map[string]*redis.Client, jobId string, orchestratorId string, executionList *ExecutionList, storeMongoDB *mongo.Database) {
	// Try to HGetAll the orchestrator id
	job, err := fetchJob(ctx, orchestratorClient, jobId)
	if err != nil || job == nil {
		fmt.Println("Error getting job")
		return
	}

	// Check if orchestrtator has already been assigned
	if job.AssignedOrchestrator != "" {
		return
	}

	// TODO: check if has capacity to execute here

	// Check if currently full execution list
	if !checkExecutionCapacity(executionList) {
		return
	}

	// HSetNX assignedOrchestrator to the orchestratorId
	assignmentResult, err := orchestratorClient.HSetNX(ctx, job.Id, "assignedOrchestrator", orchestratorId).Result()

	if err != nil {
		fmt.Println("Error setting orchestrator")
		return
	}

	// If result is 0, orchestrator is already assigned
	if !assignmentResult {
		return
	}

	// We got the job
	job.AssignedOrchestrator = orchestratorId
	executionList.addJob(*job)
	manageExecution(ctx, orchestratorClient, workerClients, *job, orchestratorId, executionList, storeMongoDB)

	// Capacity was freed, so check for queued jobs
	checkForQueuedJobs(ctx, orchestratorClient, workerClients, orchestratorId, executionList, storeMongoDB)
}

type jobDistribution struct {
	jobs         *[]libOrch.ChildJob
	workerClient *redis.Client
}

// Over-arching function that manages the execution of a job and handles its state and lifecycle
// This is the highest level function with global state
func manageExecution(ctx context.Context, orchestratorClient *redis.Client, workerClients map[string]*redis.Client, job libOrch.Job, orchestratorId string, executionList *ExecutionList, storeMongoDB *mongo.Database) {
	// Get the job id and check if it is a string
	fmt.Println("Assigned job", job.Id)
	libOrch.UpdateStatus(ctx, orchestratorClient, job.Id, orchestratorId, "ASSIGNED")

	// Setup the job

	healthy := true

	gs := NewGlobalState(ctx, orchestratorClient, job.Id, orchestratorId)

	options, err := options.DetermineRuntimeOptions(job, gs)
	if err != nil {
		libOrch.HandleStringError(ctx, orchestratorClient, job.Id, orchestratorId, fmt.Sprintf("Error determining runtime options: %s", err.Error()))
		healthy = false
	}

	if healthy {
		marshalledOptions, err := json.Marshal(options)
		if err != nil {
			libOrch.HandleStringError(ctx, orchestratorClient, job.Id, orchestratorId, fmt.Sprintf("Error marshalling runtime options: %s", err.Error()))
			healthy = false
		}

		libOrch.DispatchMessage(ctx, orchestratorClient, job.Id, orchestratorId, string(marshalledOptions), "OPTIONS")
	}

	scope := job.Scope

	childJobs := make(map[string]jobDistribution)

	if healthy {
		childJob := libOrch.ChildJob{
			Job:        job,
			ChildJobId: uuid.NewString(),
			Options:    *options,
		}

		childJobs["portsmouth"] = jobDistribution{
			jobs:         &[]libOrch.ChildJob{childJob},
			workerClient: workerClients["portsmouth"],
		}
	}

	// Run the job

	result := "FAILED"

	if healthy {
		result, err = runExecution(gs, options, scope, childJobs, job.Id)
		if err != nil {
			fmt.Println("Error running execution", err)
			libOrch.HandleError(ctx, orchestratorClient, job.Id, orchestratorId, err)
		}
	}

	libOrch.UpdateStatus(ctx, orchestratorClient, job.Id, orchestratorId, result)

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
		libOrch.HandleError(ctx, orchestratorClient, job.Id, orchestratorId, err)
		return
	}
	libOrch.DispatchMessage(ctx, orchestratorClient, job.Id, orchestratorId, string(marshalledGlobeTestReceipt), "MARK")

	//Create Metrics Store receipt, note this must be sent after cleanup
	metricsStoreReceipt := primitive.NewObjectID()
	metricsStoreReceiptMessage := &libOrch.MarkMessage{
		Mark:    "MetricsStoreReceipt",
		Message: metricsStoreReceipt.Hex(),
	}

	marshalledMetricsStoreReceipt, err := json.Marshal(metricsStoreReceiptMessage)
	if err != nil {
		fmt.Println("Error marshalling metrics store receipt", err)
		libOrch.HandleError(ctx, orchestratorClient, job.Id, orchestratorId, err)
		return
	}
	libOrch.DispatchMessage(ctx, orchestratorClient, job.Id, orchestratorId, string(marshalledMetricsStoreReceipt), "MARK")

	// Clean up the job and store result in Mongo
	err = cleanup(ctx, job, childJobs, orchestratorClient, orchestratorId, storeMongoDB, scope, globeTestLogsReceipt, metricsStoreReceipt)
	if err != nil {
		fmt.Println("Error cleaning up", err)
		libOrch.HandleErrorNoSet(ctx, orchestratorClient, job.Id, orchestratorId, err)
		libOrch.UpdateStatusNoSet(ctx, orchestratorClient, job.Id, orchestratorId, result)
	} else {
		libOrch.UpdateStatusNoSet(ctx, orchestratorClient, job.Id, orchestratorId, fmt.Sprintf("COMPLETED_%s", result))
	}

	executionList.removeJob(job.Id)
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
