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

type ExecutionList struct ***REMOVED***
	currentJobs            map[string]libOrch.ChildJob
	mutex                  sync.Mutex
	maxJobs                int
	maxIterationsPerSecond int
	maxVUs                 int
***REMOVED***

func Run() ***REMOVED***
	ctx := context.Background()

	workerId := uuid.NewString()

	fmt.Println("Client host", libWorker.GetEnvVariable("CLIENT_HOST", "localhost"))
	fmt.Println("Client port", libWorker.GetEnvVariable("CLIENT_PORT", "6379"))
	fmt.Println("Client password", libWorker.GetEnvVariable("CLIENT_PASSWORD", ""))

	client := redis.NewClient(&redis.Options***REMOVED***
		Addr:     fmt.Sprintf("%s:%s", libWorker.GetEnvVariable("CLIENT_HOST", "localhost"), libWorker.GetEnvVariable("CLIENT_PORT", "6978")),
		Username: "default",
		Password: libWorker.GetEnvVariable("CLIENT_PASSWORD", ""),
	***REMOVED***)

	// Set the worker id and current time
	client.SAdd(ctx, "workers", workerId)

	// Every second set a heartbeat update
	heartbeatTicker := time.NewTicker(1 * time.Second)
	go func() ***REMOVED***
		for range heartbeatTicker.C ***REMOVED***
			client.Set(ctx, fmt.Sprintf("worker:%s:lastHeartbeat", workerId), time.Now().UnixMilli(), time.Second*10)
		***REMOVED***
	***REMOVED***()

	fmt.Print("\n\033[1;35mAPITEAM Worker\033[0m\n\n")
	fmt.Printf("Starting worker %s\n", workerId)
	fmt.Printf("Listening for new jobs on %s...\n\n", client.Options().Addr)

	executionList := &ExecutionList***REMOVED***
		currentJobs:            make(map[string]libOrch.ChildJob),
		maxJobs:                -1,
		maxIterationsPerSecond: -1,
		maxVUs:                 -1,
	***REMOVED***

	go checkForQueuedJobs(ctx, client, workerId, executionList)

	// Subscribe to the execution channel
	pubSub := client.Subscribe(ctx, "worker:execution")

	channel := pubSub.Channel()

	for msg := range channel ***REMOVED***
		childJobId, err := uuid.Parse(msg.Payload)
		if err != nil ***REMOVED***
			fmt.Println("Error, got did not parse job id")
			return
		***REMOVED***
		go checkIfCanExecute(ctx, client, childJobId.String(), workerId, executionList)
	***REMOVED***
***REMOVED***

/*
Check for queued jobs that were deferered as they couldn't be executed when they
were queued as no workers were available.
*/
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
		fmt.Println("Error fetching child job", err)
		return
	***REMOVED***

	if job.Id == "" ***REMOVED***
		_, err = client.Del(ctx, childJobId).Result()
		if err != nil ***REMOVED***
			fmt.Println("Error deleting job from redis")
		***REMOVED***
		return
	***REMOVED***

	assignedWorker, _ := client.HGet(ctx, childJobId, "assignedWorker").Result()
	if assignedWorker != "" ***REMOVED***
		return
	***REMOVED***

	// TODO: check if has capacity to execute here

	// Check if currently full execution list
	if !checkExecutionCapacity(executionList) ***REMOVED***
		return
	***REMOVED***

	// HSetNX assignedWorker to the workerId
	assignmentResult, err := client.HSetNX(ctx, childJobId, "assignedWorker", workerId).Result()

	if err != nil ***REMOVED***
		fmt.Println("Error setting worker")
		return
	***REMOVED***

	// If result is 0, worker is already assigned
	if !assignmentResult ***REMOVED***
		return
	***REMOVED***

	// We got the job
	executionList.addJob(*job)

	go libWorker.UpdateStatus(ctx, client, childJobId, workerId, "ASSIGNED")
	handleExecution(ctx, client, *job, workerId)
	executionList.removeJob(childJobId)
	// Capacity was freed, so check for queued jobs
	checkForQueuedJobs(ctx, client, workerId, executionList)
***REMOVED***

func (e *ExecutionList) addJob(job libOrch.ChildJob) ***REMOVED***
	e.mutex.Lock()
	e.currentJobs[job.Id] = job
	e.mutex.Unlock()
***REMOVED***

func (e *ExecutionList) removeJob(childJobId string) ***REMOVED***
	e.mutex.Lock()
	delete(e.currentJobs, childJobId)
	e.mutex.Unlock()
***REMOVED***

func checkExecutionCapacity(executionList *ExecutionList) bool ***REMOVED***
	// TODO: check if has capacity to execute here

	// If more than max jobs, return false
	if executionList.maxJobs >= 0 && len(executionList.currentJobs) >= executionList.maxJobs ***REMOVED***
		return false
	***REMOVED***

	// TODO: implement more capacity checks

	return true
***REMOVED***
