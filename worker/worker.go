package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/APITeamLimited/redis/v9"
	"github.com/google/uuid"
	"go.k6.io/k6/lib"
)

type ExecutionList struct {
	currentJobs            map[string]map[string]string
	mutex                  sync.Mutex
	maxJobs                int
	maxIterationsPerSecond int
	maxVUs                 int
}

func Run() {
	ctx := context.Background()

	workerId := uuid.New()

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", lib.GetEnvVariable("CLIENT_HOST", "localhost"), lib.GetEnvVariable("CLIENT_PORT", "6978")),
		Password: lib.GetEnvVariable("CLIENT_PASSWORD", ""),
		DB:       0, // use default DB
	})

	currentTime := time.Now().UnixMilli()

	executionList := &ExecutionList{
		currentJobs:            make(map[string]map[string]string),
		maxJobs:                -1,
		maxIterationsPerSecond: -1,
		maxVUs:                 -1,
	}

	//Set the worker id and current time
	client.HSet(ctx, "k6:workers", workerId.String(), currentTime)

	fmt.Print("\n\033[1;35mAPITEAM Worker\033[0m\n\n")
	fmt.Printf("Starting worker %s\n", workerId.String())
	fmt.Printf("Listening for new jobs on %s...\n", client.Options().Addr)

	go checkForQueuedJobs(ctx, client, workerId.String(), executionList)

	// Subscribe to the execution channel
	psc := client.Subscribe(ctx, "k6:execution")

	channel := psc.Channel()

	for msg := range channel {
		jobId, err := uuid.Parse(msg.Payload)
		if err != nil {
			fmt.Println("Error, got did not parse job id")
			return
		}
		go checkIfCanExecute(ctx, client, jobId.String(), workerId.String(), executionList)
	}
}

/*
Check for queued jobs that were deferered as they couldn't be executed when they
were queued as no workers were available.
*/
func checkForQueuedJobs(ctx context.Context, client *redis.Client, workerId string, executionList *ExecutionList) {
	// Check for job keys in the "k6:executionHistory" set
	historyIds, err := client.SMembers(ctx, "k6:executionHistory").Result()
	if err != nil {
		fmt.Println("Error getting history ids", err)
	}

	for _, jobId := range historyIds {
		go checkIfCanExecute(ctx, client, jobId, workerId, executionList)
	}
}

func checkIfCanExecute(ctx context.Context, client *redis.Client, jobId string, workerId string, executionList *ExecutionList) {
	// Try to HGetAll the worker id
	job, err := client.HGetAll(ctx, jobId).Result()

	if err != nil {
		fmt.Println("Error getting job from redis")
		return
	}

	// TODO: check if has capacity to execute here

	// Check worker['assignedWorker'] is nil
	if job["assignedWorker"] != "" {
		return
	}

	if job["id"] == "" {
		_, err = client.Del(ctx, jobId).Result()
		if err != nil {
			fmt.Println("Error deleting job from redis")
		}
		return
	}

	// Check if currently full execution list
	if !checkExecutionCapacity(executionList) {
		return
	}

	// HSetNX assignedWorker to the workerId
	assignmentResult, err := client.HSetNX(ctx, jobId, "assignedWorker", workerId).Result()

	if err != nil {
		fmt.Println("Error setting worker")
		return
	}

	// If result is 0, worker is already assigned
	if !assignmentResult {
		return
	}

	// We got the job
	executionList.addJob(job)

	go lib.UpdateStatus(ctx, client, jobId, workerId, "ASSIGNED")
	handleExecution(ctx, client, job, workerId)
	executionList.removeJob(jobId)
	// Capacity was freed, so check for queued jobs
	checkForQueuedJobs(ctx, client, workerId, executionList)
}

func (e *ExecutionList) addJob(job map[string]string) {
	e.mutex.Lock()
	e.currentJobs[job["id"]] = job
	e.mutex.Unlock()
}

func (e *ExecutionList) removeJob(jobId string) {
	e.mutex.Lock()
	delete(e.currentJobs, jobId)
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
