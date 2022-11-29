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

func Run(standalone bool) ***REMOVED***
	ctx := context.Background()
	orchestratorId := uuid.NewString()

	if standalone ***REMOVED***
		fmt.Print("\n\033[1;35mGlobeTest Orchestrator\033[0m\n\n")
	***REMOVED***
	fmt.Printf("Starting orchestrator %s\n", orchestratorId)

	orchestratorClient := getOrchestratorOrchestratorClient(standalone)
	storeMongoDB := getStoreMongoDB(ctx, standalone)
	workerClients := connectWorkerClients(ctx, standalone)
	maxJobs := getMaxJobs(standalone)
	maxManagedVUs := getMaxManagedVUs(standalone)

	var creditsClient *redis.Client

	if standalone ***REMOVED***
		lib.GetCreditsClient(standalone)
	***REMOVED***

	executionList := &ExecutionList***REMOVED***
		currentJobs:   make(map[string]libOrch.Job),
		maxJobs:       maxJobs,
		maxManagedVUs: maxManagedVUs,
	***REMOVED***

	// Create a scheduler for regular updates and checks
	startJobScheduling(ctx, orchestratorClient, orchestratorId, executionList, workerClients, storeMongoDB, creditsClient, standalone)

	// Periodically check for and delete offline orchestrators
	if strings.ToLower(lib.GetEnvVariable("IS_MASTER", "false")) == "true" ***REMOVED***
		createMasterScheduler(ctx, orchestratorClient, workerClients)
	***REMOVED***

	fmt.Printf("Orchestrator listening for new jobs on %s...\n", orchestratorClient.Options().Addr)

	// Subscribe to the execution channel and listen for new jobs
	channel := orchestratorClient.Subscribe(ctx, "orchestrator:execution").Channel()

	for msg := range channel ***REMOVED***
		jobId, err := uuid.Parse(msg.Payload)
		if err != nil ***REMOVED***
			fmt.Println("Error, got did not parse job id")
			return
		***REMOVED***
		go checkIfCanExecute(ctx, orchestratorClient, workerClients, jobId.String(),
			orchestratorId, executionList, storeMongoDB, creditsClient, standalone)
	***REMOVED***
***REMOVED***

// Ensures job has no already been assigned and determines if this node has capacity to execute
func checkIfCanExecute(ctx context.Context, orchestratorClient *redis.Client, workerClients libOrch.WorkerClients,
	jobId string, orchestratorId string, executionList *ExecutionList, storeMongoDB *mongo.Database,
	creditsClient *redis.Client, standalone bool) ***REMOVED***
	// Try to HGetAll the orchestrator id
	job, err := fetchJob(ctx, orchestratorClient, jobId)
	if err != nil || job == nil ***REMOVED***
		if err != nil ***REMOVED***
			fmt.Println("Error getting job from orchestrator:executionHistory set, it will be deleted:", err)

			// Remove the job from the history set
			orchestratorClient.SRem(ctx, "orchestrator:executionHistory", jobId).Result()
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

	gs := NewGlobalState(ctx, orchestratorClient, job, orchestratorId, creditsClient, standalone)
	options, optionsErr := options.DetermineRuntimeOptions(*job, gs, workerClients)
	job.Options = options

	if optionsErr != nil ***REMOVED***
		libOrch.HandleError(gs, optionsErr)
		// Continue as need to delete job
	***REMOVED***

	// Check execution capacity again, bearing in mind options
	executionList.mutex.Lock()
	if !executionList.checkExecutionCapacity(options) ***REMOVED***
		executionList.mutex.Unlock()
		gs.CreditsManager().StopCreditsCapturing()
		return
	***REMOVED***

	// HSetNX assignedOrchestrator to the orchestratorId
	assignmentResult, err := orchestratorClient.HSetNX(ctx, job.Id, "assignedOrchestrator", orchestratorId).Result()

	// If result is 0, orchestrator is already assigned
	if !assignmentResult ***REMOVED***
		executionList.mutex.Unlock()
		gs.CreditsManager().StopCreditsCapturing()
		return
	***REMOVED***

	if err != nil ***REMOVED***
		fmt.Println("Error setting orchestrator")
		executionList.mutex.Unlock()
		gs.CreditsManager().StopCreditsCapturing()
		return
	***REMOVED***

	// We got the job and have confirmed capacity for it

	job.AssignedOrchestrator = orchestratorId

	executionList.addJob(job)
	executionList.mutex.Unlock()
	defer executionList.removeJob(job.Id)

	fmt.Println("Assigned job:", job.Id)

	successfullExecution := manageExecution(gs, orchestratorClient, workerClients,
		*job, orchestratorId, executionList, storeMongoDB, optionsErr)

	if successfullExecution ***REMOVED***
		fmt.Printf("Completed job successfully: %s\n", job.Id)
	***REMOVED*** else ***REMOVED***
		fmt.Printf("Job failed: %s\n", job.Id)
	***REMOVED***
***REMOVED***

// Check for queued jobs that were deferered as they couldn't be executed when they
// were queued as no workers were available.
func checkForQueuedJobs(ctx context.Context, orchestratorClient *redis.Client,
	workerClients libOrch.WorkerClients, orchestratorId string, executionList *ExecutionList,
	storeMongoDB *mongo.Database, creditsClient *redis.Client, standalone bool) ***REMOVED***
	// Check for job keys in the "orchestrator:executionHistory" set
	historyIds, err := orchestratorClient.SMembers(ctx, "orchestrator:executionHistory").Result()
	if err != nil ***REMOVED***
		fmt.Println("Error getting history ids", err)
	***REMOVED***

	for _, jobId := range historyIds ***REMOVED***
		go checkIfCanExecute(ctx, orchestratorClient, workerClients, jobId,
			orchestratorId, executionList, storeMongoDB, creditsClient, standalone)
	***REMOVED***
***REMOVED***
