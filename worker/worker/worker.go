package worker

import (
	"context"
	"fmt"

	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/redis/v9"
	"github.com/google/uuid"
)

func Run(standalone bool) ***REMOVED***
	ctx := context.Background()
	workerId := uuid.NewString()

	fmt.Print("\n\033[1;35mGlobeTest Worker\033[0m\n\n")
	fmt.Printf("Starting worker %s\n", workerId)

	client := getWorkerClient(standalone)
	maxJobs := getMaxJobs(standalone)
	maxVUs := getMaxVUs(standalone)
	creditsClient := lib.GetCreditsClient(standalone)

	executionList := &ExecutionList***REMOVED***
		currentJobs: make(map[string]libOrch.ChildJob),
		maxJobs:     maxJobs,
		maxVUs:      maxVUs,
	***REMOVED***

	// Create a scheduler for regular updates and checks
	startScheduling(ctx, client, workerId, executionList, creditsClient)

	fmt.Printf("Listening for new jobs on %s...\n\n", client.Options().Addr)

	// Subscribe to the execution channel
	channel := client.Subscribe(ctx, "worker:execution").Channel()

	for msg := range channel ***REMOVED***
		childJobId, err := uuid.Parse(msg.Payload)
		if err != nil ***REMOVED***
			fmt.Println("Error, got did not parse job id")
			return
		***REMOVED***
		go checkIfCanExecute(ctx, client, childJobId.String(), workerId, executionList, creditsClient)
	***REMOVED***
***REMOVED***

// Check for queued jobs that were deferered as they couldn't be executed when they
// were queued as no workers were available.
func checkForQueuedJobs(ctx context.Context, client *redis.Client, workerId string,
	executionList *ExecutionList, creditsClient *redis.Client) ***REMOVED***
	// Check for job keys in the "worker:executionHistory" set
	historyIds, err := client.SMembers(ctx, "worker:executionHistory").Result()
	if err != nil ***REMOVED***
		fmt.Println("Error getting history ids", err)
	***REMOVED***

	for _, childJobId := range historyIds ***REMOVED***
		go checkIfCanExecute(ctx, client, childJobId, workerId, executionList, creditsClient)
	***REMOVED***
***REMOVED***

func checkIfCanExecute(ctx context.Context, client *redis.Client, childJobId string,
	workerId string, executionList *ExecutionList, creditsClient *redis.Client) ***REMOVED***
	// Try to HGetAll the worker id
	job, err := fetchChildJob(ctx, client, childJobId)
	if err != nil || job == nil ***REMOVED***
		if err != nil ***REMOVED***
			fmt.Println("Error getting child job from worker:executionHistory set, it will be deleted:", err)

			// Remove the job from the history set
			client.SRem(ctx, "worker:executionHistory", childJobId).Result()
		***REMOVED***
		return
	***REMOVED***

	assignedWorker, _ := client.HGet(ctx, childJobId, "assignedWorker").Result()
	if assignedWorker != "" ***REMOVED***
		return
	***REMOVED***

	executionList.mutex.Lock()

	// Check if this node has execution capacity for this job
	if !executionList.checkExecutionCapacity(job.Options) ***REMOVED***
		executionList.mutex.Unlock()
		return
	***REMOVED***

	// HSetNX assignedWorker to the workerId
	assignmentResult, err := client.HSetNX(ctx, childJobId, "assignedWorker", workerId).Result()

	// If result is 0, worker is already assigned
	if !assignmentResult ***REMOVED***
		executionList.mutex.Unlock()
		return
	***REMOVED***

	if err != nil ***REMOVED***
		executionList.mutex.Unlock()
		return
	***REMOVED***

	// We got the job
	executionList.addJob(*job)
	executionList.mutex.Unlock()
	defer executionList.removeJob(childJobId)

	fmt.Printf("Got child job: %s\n", job.ChildJobId)

	successfullExecution := handleExecution(ctx, client, *job, workerId, creditsClient)

	if successfullExecution ***REMOVED***
		fmt.Printf("Completed child job successfully: %s\n", job.ChildJobId)
	***REMOVED*** else ***REMOVED***
		fmt.Printf("Error executing child job: %s\n", job.ChildJobId)
	***REMOVED***

***REMOVED***
