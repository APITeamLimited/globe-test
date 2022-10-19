package worker

import (
	"context"
	"fmt"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/redis/v9"
	"github.com/google/uuid"
)

func Run() ***REMOVED***
	ctx := context.Background()
	workerId := uuid.NewString()
	client := getWorkerClient()
	maxJobs := getMaxJobs()
	maxVUs := getMaxVUs()

	executionList := &ExecutionList***REMOVED***
		currentJobs: make(map[string]libOrch.ChildJob),
		maxJobs:     maxJobs,
		maxVUs:      maxVUs,
	***REMOVED***

	// Create a scheduler for regular updates and checks
	startJobScheduling(ctx, client, workerId, executionList)

	fmt.Print("\n\033[1;35mAPITEAM Worker\033[0m\n\n")
	fmt.Printf("Starting worker %s\n", workerId)
	fmt.Printf("Listening for new jobs on %s...\n\n", client.Options().Addr)

	// Subscribe to the execution channel
	channel := client.Subscribe(ctx, "worker:execution").Channel()

	for msg := range channel ***REMOVED***
		childJobId, err := uuid.Parse(msg.Payload)
		if err != nil ***REMOVED***
			fmt.Println("Error, got did not parse job id")
			return
		***REMOVED***
		go checkIfCanExecute(ctx, client, childJobId.String(), workerId, executionList)
	***REMOVED***
***REMOVED***

// Check for queued jobs that were deferered as they couldn't be executed when they
// were queued as no workers were available.
func checkForQueuedJobs(ctx context.Context, client *redis.Client, workerId string, executionList *ExecutionList) ***REMOVED***
	// Check for job keys in the "worker:executionHistory" set
	historyIds, err := client.SMembers(ctx, "worker:executionHistory").Result()
	if err != nil ***REMOVED***
		fmt.Println("Error getting history ids", err)
	***REMOVED***

	for _, childJobId := range historyIds ***REMOVED***
		go checkIfCanExecute(ctx, client, childJobId, workerId, executionList)
	***REMOVED***
***REMOVED***

func checkIfCanExecute(ctx context.Context, client *redis.Client, childJobId string, workerId string, executionList *ExecutionList) ***REMOVED***
	// Try to HGetAll the worker id
	job, err := fetchChildJob(ctx, client, childJobId)
	if err != nil || job == nil ***REMOVED***
		if err != nil ***REMOVED***
			fmt.Println("Error getting job")
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
		fmt.Println("Error setting worker")
		executionList.mutex.Unlock()
		return
	***REMOVED***

	// We got the job
	executionList.addJob(*job)
	executionList.mutex.Unlock()

	libWorker.UpdateStatus(ctx, client, childJobId, workerId, "ASSIGNED")
	handleExecution(ctx, client, *job, workerId)
	executionList.removeJob(childJobId)
***REMOVED***
