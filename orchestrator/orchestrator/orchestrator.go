package orchestrator

import (
	"context"
	"fmt"
	"strings"

	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/orchestrator/options"
	"github.com/APITeamLimited/redis/v9"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

func Run() {
	ctx := context.Background()
	orchestratorId := uuid.NewString()

	fmt.Print("\n\033[1;35mGlobeTest Orchestrator\033[0m\n\n")
	fmt.Printf("Starting orchestrator %s\n", orchestratorId)

	orchestratorClient := getOrchestratorOrchestratorClient()
	storeMongoDB := getStoreMongoDB(ctx)
	workerClients := connectWorkerClients(ctx)
	maxJobs := getMaxJobs()
	maxManagedVUs := getMaxManagedVUs()
	creditsClient := lib.GetCreditsClient()

	executionList := &ExecutionList{
		currentJobs:   make(map[string]libOrch.Job),
		maxJobs:       maxJobs,
		maxManagedVUs: maxManagedVUs,
	}

	// Create a scheduler for regular updates and checks
	startJobScheduling(ctx, orchestratorClient, orchestratorId, executionList, workerClients, storeMongoDB, creditsClient)

	// Periodically check for and delete offline orchestrators
	if strings.ToLower(lib.GetEnvVariable("IS_MASTER", "false")) == "true" {
		createMasterScheduler(ctx, orchestratorClient, workerClients)
	}

	fmt.Printf("Listening for new jobs on %s...\n\n", orchestratorClient.Options().Addr)

	// Subscribe to the execution channel and listen for new jobs
	channel := orchestratorClient.Subscribe(ctx, "orchestrator:execution").Channel()

	for msg := range channel {
		jobId, err := uuid.Parse(msg.Payload)
		if err != nil {
			fmt.Println("Error, got did not parse job id")
			return
		}
		go checkIfCanExecute(ctx, orchestratorClient, workerClients, jobId.String(),
			orchestratorId, executionList, storeMongoDB, creditsClient)
	}
}

// Ensures job has no already been assigned and determines if this node has capacity to execute
func checkIfCanExecute(ctx context.Context, orchestratorClient *redis.Client, workerClients libOrch.WorkerClients,
	jobId string, orchestratorId string, executionList *ExecutionList, storeMongoDB *mongo.Database,
	creditsClient *redis.Client) {
	// Try to HGetAll the orchestrator id
	job, err := fetchJob(ctx, orchestratorClient, jobId)
	if err != nil || job == nil {
		if err != nil {
			fmt.Println("Error getting job from orchestrator:executionHistory set, it will be deleted:", err)

			// Remove the job from the history set
			orchestratorClient.SRem(ctx, "orchestrator:executionHistory", jobId).Result()
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

	gs := NewGlobalState(ctx, orchestratorClient, job, orchestratorId, creditsClient)
	options, optionsErr := options.DetermineRuntimeOptions(*job, gs, workerClients)
	job.Options = options

	if optionsErr != nil {
		libOrch.HandleError(gs, optionsErr)
		// Continue as need to delete job
	}

	// Check execution capacity again, bearing in mind options
	executionList.mutex.Lock()
	if !executionList.checkExecutionCapacity(options) {
		executionList.mutex.Unlock()
		gs.CreditsManager().StopCreditsCapturing()
		return
	}

	// HSetNX assignedOrchestrator to the orchestratorId
	assignmentResult, err := orchestratorClient.HSetNX(ctx, job.Id, "assignedOrchestrator", orchestratorId).Result()

	// If result is 0, orchestrator is already assigned
	if !assignmentResult {
		executionList.mutex.Unlock()
		gs.CreditsManager().StopCreditsCapturing()
		return
	}

	if err != nil {
		fmt.Println("Error setting orchestrator")
		executionList.mutex.Unlock()
		gs.CreditsManager().StopCreditsCapturing()
		return
	}

	// We got the job and have confirmed capacity for it

	job.AssignedOrchestrator = orchestratorId

	executionList.addJob(job)
	executionList.mutex.Unlock()
	defer executionList.removeJob(job.Id)

	fmt.Println("Assigned job:", job.Id)

	successfullExecution := manageExecution(gs, orchestratorClient, workerClients,
		*job, orchestratorId, executionList, storeMongoDB, optionsErr)

	if successfullExecution {
		fmt.Printf("Completed job successfully: %s\n", job.Id)
	} else {
		fmt.Printf("Job failed: %s\n", job.Id)
	}
}

// Check for queued jobs that were deferered as they couldn't be executed when they
// were queued as no workers were available.
func checkForQueuedJobs(ctx context.Context, orchestratorClient *redis.Client,
	workerClients libOrch.WorkerClients, orchestratorId string, executionList *ExecutionList,
	storeMongoDB *mongo.Database, creditsClient *redis.Client) {
	// Check for job keys in the "orchestrator:executionHistory" set
	historyIds, err := orchestratorClient.SMembers(ctx, "orchestrator:executionHistory").Result()
	if err != nil {
		fmt.Println("Error getting history ids", err)
	}

	for _, jobId := range historyIds {
		go checkIfCanExecute(ctx, orchestratorClient, workerClients, jobId,
			orchestratorId, executionList, storeMongoDB, creditsClient)
	}
}
