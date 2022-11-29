package worker

import (
	"context"
	"fmt"

	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/redis/v9"
	"github.com/google/uuid"
)

func Run(standalone bool) {
	ctx := context.Background()
	workerId := uuid.NewString()

	if standalone {
		fmt.Print("\n\033[1;35mGlobeTest Worker\033[0m\n\n")
	}
	fmt.Printf("Starting worker %s\n", workerId)

	client := getWorkerClient(standalone)
	maxJobs := getMaxJobs(standalone)
	maxVUs := getMaxVUs(standalone)
	creditsClient := lib.GetCreditsClient(standalone)

	executionList := &ExecutionList{
		currentJobs: make(map[string]libOrch.ChildJob),
		maxJobs:     maxJobs,
		maxVUs:      maxVUs,
	}

	// Create a scheduler for regular updates and checks
	startScheduling(ctx, client, workerId, executionList, creditsClient, standalone)

	fmt.Printf("Worker listening for new jobs on %s...\n", client.Options().Addr)

	// Subscribe to the execution channel
	channel := client.Subscribe(ctx, "worker:execution").Channel()

	for msg := range channel {
		childJobId, err := uuid.Parse(msg.Payload)
		if err != nil {
			fmt.Println("Error, got did not parse job id")
			return
		}
		go checkIfCanExecute(ctx, client, childJobId.String(), workerId, executionList, creditsClient, standalone)
	}
}

// Check for queued jobs that were deferered as they couldn't be executed when they
// were queued as no workers were available.
func checkForQueuedJobs(ctx context.Context, client *redis.Client, workerId string,
	executionList *ExecutionList, creditsClient *redis.Client, standalone bool) {
	// Check for job keys in the "worker:executionHistory" set
	historyIds, err := client.SMembers(ctx, "worker:executionHistory").Result()
	if err != nil {
		fmt.Println("Error getting history ids", err)
	}

	for _, childJobId := range historyIds {
		go checkIfCanExecute(ctx, client, childJobId, workerId, executionList, creditsClient, standalone)
	}
}

func checkIfCanExecute(ctx context.Context, client *redis.Client, childJobId string,
	workerId string, executionList *ExecutionList, creditsClient *redis.Client, standalone bool) {
	// Try to HGetAll the worker id
	job, err := fetchChildJob(ctx, client, childJobId)
	if err != nil || job == nil {
		if err != nil {
			fmt.Println("Error getting child job from worker:executionHistory set, it will be deleted:", err)

			// Remove the job from the history set
			client.SRem(ctx, "worker:executionHistory", childJobId).Result()
		}
		return
	}

	assignedWorker, _ := client.HGet(ctx, childJobId, "assignedWorker").Result()
	if assignedWorker != "" {
		return
	}

	executionList.mutex.Lock()

	// Check if this node has execution capacity for this job
	if !executionList.checkExecutionCapacity(job.Options) {
		executionList.mutex.Unlock()
		return
	}

	// HSetNX assignedWorker to the workerId
	assignmentResult, err := client.HSetNX(ctx, childJobId, "assignedWorker", workerId).Result()

	// If result is 0, worker is already assigned
	if !assignmentResult {
		executionList.mutex.Unlock()
		return
	}

	if err != nil {
		executionList.mutex.Unlock()
		return
	}

	// We got the job
	executionList.addJob(*job)
	executionList.mutex.Unlock()
	defer executionList.removeJob(childJobId)

	fmt.Printf("Got child job: %s\n", job.ChildJobId)

	successfullExecution := handleExecution(ctx, client, *job, workerId, creditsClient, standalone)

	if successfullExecution {
		fmt.Printf("Completed child job successfully: %s\n", job.ChildJobId)
	} else {
		fmt.Printf("Error executing child job: %s\n", job.ChildJobId)
	}

}
