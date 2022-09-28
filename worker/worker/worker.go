package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/redis/v9"
	"github.com/google/uuid"
)

type ExecutionList struct ***REMOVED***
	currentJobs            map[string]map[string]string
	mutex                  sync.Mutex
	maxJobs                int
	maxIterationsPerSecond int
	maxVUs                 int
***REMOVED***

func Run() ***REMOVED***
	ctx := context.Background()

	workerId := uuid.New()

	client := redis.NewClient(&redis.Options***REMOVED***
		Addr:     fmt.Sprintf("%s:%s", libWorker.GetEnvVariable("CLIENT_HOST", "localhost"), libWorker.GetEnvVariable("CLIENT_PORT", "6978")),
		Password: libWorker.GetEnvVariable("CLIENT_PASSWORD", ""),
		DB:       0, // use default DB
	***REMOVED***)

	currentTime := time.Now().UnixMilli()

	executionList := &ExecutionList***REMOVED***
		currentJobs:            make(map[string]map[string]string),
		maxJobs:                -1,
		maxIterationsPerSecond: -1,
		maxVUs:                 -1,
	***REMOVED***

	//Set the worker id and current time
	client.HSet(ctx, "k6:workers", workerId.String(), currentTime)

	fmt.Print("\n\033[1;35mAPITEAM Worker\033[0m\n\n")
	fmt.Printf("Starting worker %s\n", workerId.String())
	fmt.Printf("Listening for new jobs on %s...\n", client.Options().Addr)

	go checkForQueuedJobs(ctx, client, workerId.String(), executionList)

	// Subscribe to the execution channel
	psc := client.Subscribe(ctx, "k6:execution")

	channel := psc.Channel()

	for msg := range channel ***REMOVED***
		jobId, err := uuid.Parse(msg.Payload)
		if err != nil ***REMOVED***
			fmt.Println("Error, got did not parse job id")
			return
		***REMOVED***
		go checkIfCanExecute(ctx, client, jobId.String(), workerId.String(), executionList)
	***REMOVED***
***REMOVED***

/*
Check for queued jobs that were deferered as they couldn't be executed when they
were queued as no workers were available.
*/
func checkForQueuedJobs(ctx context.Context, client *redis.Client, workerId string, executionList *ExecutionList) ***REMOVED***
	// Check for job keys in the "k6:executionHistory" set
	historyIds, err := client.SMembers(ctx, "k6:executionHistory").Result()
	if err != nil ***REMOVED***
		fmt.Println("Error getting history ids", err)
	***REMOVED***

	for _, jobId := range historyIds ***REMOVED***
		go checkIfCanExecute(ctx, client, jobId, workerId, executionList)
	***REMOVED***
***REMOVED***

func checkIfCanExecute(ctx context.Context, client *redis.Client, jobId string, workerId string, executionList *ExecutionList) ***REMOVED***
	// Try to HGetAll the worker id
	job, err := client.HGetAll(ctx, jobId).Result()

	if err != nil ***REMOVED***
		fmt.Println("Error getting job from redis")
		return
	***REMOVED***

	// TODO: check if has capacity to execute here

	// Check worker['assignedWorker'] is nil
	if job["assignedWorker"] != "" ***REMOVED***
		return
	***REMOVED***

	if job["id"] == "" ***REMOVED***
		_, err = client.Del(ctx, jobId).Result()
		if err != nil ***REMOVED***
			fmt.Println("Error deleting job from redis")
		***REMOVED***
		return
	***REMOVED***

	// Check if currently full execution list
	if !checkExecutionCapacity(executionList) ***REMOVED***
		return
	***REMOVED***

	// HSetNX assignedWorker to the workerId
	assignmentResult, err := client.HSetNX(ctx, jobId, "assignedWorker", workerId).Result()

	if err != nil ***REMOVED***
		fmt.Println("Error setting worker")
		return
	***REMOVED***

	// If result is 0, worker is already assigned
	if !assignmentResult ***REMOVED***
		return
	***REMOVED***

	// We got the job
	executionList.addJob(job)

	go libWorker.UpdateStatus(ctx, client, jobId, workerId, "ASSIGNED")
	handleExecution(ctx, client, job, workerId)
	executionList.removeJob(jobId)
	// Capacity was freed, so check for queued jobs
	checkForQueuedJobs(ctx, client, workerId, executionList)
***REMOVED***

func (e *ExecutionList) addJob(job map[string]string) ***REMOVED***
	e.mutex.Lock()
	e.currentJobs[job["id"]] = job
	e.mutex.Unlock()
***REMOVED***

func (e *ExecutionList) removeJob(jobId string) ***REMOVED***
	e.mutex.Lock()
	delete(e.currentJobs, jobId)
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
