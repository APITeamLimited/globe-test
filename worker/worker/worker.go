package worker

import (
	"context"
	"fmt"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/redis/v9"
	"github.com/google/uuid"
)

func Run() {
	ctx := context.Background()
	workerId := uuid.NewString()
	client := getWorkerClient()
	maxJobs := getMaxJobs()
	maxVUs := getMaxVUs()

	executionList := &ExecutionList{
		currentJobs: make(map[string]libOrch.ChildJob),
		maxJobs:     maxJobs,
		maxVUs:      maxVUs,
	}

	// Create a scheduler for regular updates and checks
	startJobScheduling(ctx, client, workerId, executionList)

	// Change process title
	fmt.Printf("\033]0;GlobeTest Worker %s\007", workerId)

	fmt.Print("\n\033[1;35mGlobeTest Worker\033[0m\n\n")
	fmt.Printf("Starting worker %s\n", workerId)
	fmt.Printf("Listening for new jobs on %s...\n\n", client.Options().Addr)

	// Subscribe to the execution channel
	channel := client.Subscribe(ctx, "worker:execution").Channel()

	for msg := range channel {
		childJobId, err := uuid.Parse(msg.Payload)
		if err != nil {
			fmt.Println("Error, got did not parse job id")
			return
		}
		go checkIfCanExecute(ctx, client, childJobId.String(), workerId, executionList)
	}
}

// Check for queued jobs that were deferered as they couldn't be executed when they
// were queued as no workers were available.
func checkForQueuedJobs(ctx context.Context, client *redis.Client, workerId string, executionList *ExecutionList) {
	// Check for job keys in the "worker:executionHistory" set
	historyIds, err := client.SMembers(ctx, "worker:executionHistory").Result()
	if err != nil {
		fmt.Println("Error getting history ids", err)
	}

	for _, childJobId := range historyIds {
		go checkIfCanExecute(ctx, client, childJobId, workerId, executionList)
	}
}

func checkIfCanExecute(ctx context.Context, client *redis.Client, childJobId string, workerId string, executionList *ExecutionList) {
	// Try to HGetAll the worker id
	job, err := fetchChildJob(ctx, client, childJobId)
	if err != nil || job == nil {
		if err != nil {
			fmt.Println("Error getting job")
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
		fmt.Println("Error setting worker")
		executionList.mutex.Unlock()
		return
	}

	// We got the job
	executionList.addJob(*job)
	executionList.mutex.Unlock()

	handleExecution(ctx, client, *job, workerId)

	executionList.removeJob(childJobId)
}
