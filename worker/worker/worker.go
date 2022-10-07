package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/redis/v9"
	"github.com/google/uuid"
)

type ExecutionList struct {
	currentJobs            map[string]libOrch.ChildJob
	mutex                  sync.Mutex
	maxJobs                int
	maxIterationsPerSecond int
	maxVUs                 int
}

func Run() {
	ctx := context.Background()

	workerId := uuid.NewString()

	fmt.Println("Client host", libWorker.GetEnvVariable("CLIENT_HOST", "localhost"))
	fmt.Println("Client port", libWorker.GetEnvVariable("CLIENT_PORT", "6379"))
	fmt.Println("Client password", libWorker.GetEnvVariable("CLIENT_PASSWORD", ""))

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", libWorker.GetEnvVariable("CLIENT_HOST", "localhost"), libWorker.GetEnvVariable("CLIENT_PORT", "6978")),
		Username: "default",
		Password: libWorker.GetEnvVariable("CLIENT_PASSWORD", ""),
	})

	// Set the worker id and current time
	client.SAdd(ctx, "workers", workerId)

	// Every second set a heartbeat update
	heartbeatTicker := time.NewTicker(1 * time.Second)
	go func() {
		for range heartbeatTicker.C {
			client.Set(ctx, fmt.Sprintf("worker:%s:lastHeartbeat", workerId), time.Now().UnixMilli(), time.Second*10)
		}
	}()

	fmt.Print("\n\033[1;35mAPITEAM Worker\033[0m\n\n")
	fmt.Printf("Starting worker %s\n", workerId)
	fmt.Printf("Listening for new jobs on %s...\n\n", client.Options().Addr)

	executionList := &ExecutionList{
		currentJobs:            make(map[string]libOrch.ChildJob),
		maxJobs:                -1,
		maxIterationsPerSecond: -1,
		maxVUs:                 -1,
	}

	go checkForQueuedJobs(ctx, client, workerId, executionList)

	// Subscribe to the execution channel
	pubSub := client.Subscribe(ctx, "worker:execution")

	channel := pubSub.Channel()

	for msg := range channel {
		childJobId, err := uuid.Parse(msg.Payload)
		if err != nil {
			fmt.Println("Error, got did not parse job id")
			return
		}
		go checkIfCanExecute(ctx, client, childJobId.String(), workerId, executionList)
	}
}

/*
Check for queued jobs that were deferered as they couldn't be executed when they
were queued as no workers were available.
*/
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
		fmt.Println("Error fetching child job", err)
		return
	}

	if job.Id == "" {
		_, err = client.Del(ctx, childJobId).Result()
		if err != nil {
			fmt.Println("Error deleting job from redis")
		}
		return
	}

	assignedWorker, _ := client.HGet(ctx, childJobId, "assignedWorker").Result()
	if assignedWorker != "" {
		return
	}

	// TODO: check if has capacity to execute here

	// Check if currently full execution list
	if !checkExecutionCapacity(executionList) {
		return
	}

	// HSetNX assignedWorker to the workerId
	assignmentResult, err := client.HSetNX(ctx, childJobId, "assignedWorker", workerId).Result()

	if err != nil {
		fmt.Println("Error setting worker")
		return
	}

	// If result is 0, worker is already assigned
	if !assignmentResult {
		return
	}

	// We got the job
	executionList.addJob(*job)

	go libWorker.UpdateStatus(ctx, client, childJobId, workerId, "ASSIGNED")
	handleExecution(ctx, client, *job, workerId)
	executionList.removeJob(childJobId)
	// Capacity was freed, so check for queued jobs
	checkForQueuedJobs(ctx, client, workerId, executionList)
}

func (e *ExecutionList) addJob(job libOrch.ChildJob) {
	e.mutex.Lock()
	e.currentJobs[job.Id] = job
	e.mutex.Unlock()
}

func (e *ExecutionList) removeJob(childJobId string) {
	e.mutex.Lock()
	delete(e.currentJobs, childJobId)
	e.mutex.Unlock()
}

func checkExecutionCapacity(executionList *ExecutionList) bool {
	// TODO: check if has capacity to execute here

	// If more than max jobs, return false
	if executionList.maxJobs >= 0 && len(executionList.currentJobs) >= executionList.maxJobs {
		return false
	}

	// TODO: implement more capacity checks

	return true
}
