package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"go.k6.io/k6/lib/consts"
)

func Run() ***REMOVED***
	ctx := context.Background()

	workerId := uuid.New()

	client := redis.NewClient(&redis.Options***REMOVED***
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	***REMOVED***)

	currentTime := time.Now().UnixMilli()

	//Set the worker id and current time
	client.HSet(ctx, "k6:workers", workerId.String(), currentTime)

	c := color.New(color.FgCyan)
	fmt.Print("\n", c.Sprint(consts.Banner()), "\n\n")
	fmt.Printf("\033[1;35mStarting worker %s \033[0m\n\n", workerId.String())

	fmt.Printf("Listening for new jobs on %s...\n", client.Options().Addr)

	go startupJobCheck(ctx, client, workerId.String())

	// Subscribe to the execution channel
	psc := client.Subscribe(ctx, "k6:execution")

	channel := psc.Channel()

	for msg := range channel ***REMOVED***
		jobId, err := uuid.Parse(msg.Payload)
		if err != nil ***REMOVED***
			fmt.Println("Error, got did not parse job id")
			return
		***REMOVED***
		checkIfCanExecute(ctx, client, jobId.String(), workerId.String())
	***REMOVED***
***REMOVED***

/*
Check for jobs existing on startup
*/
func startupJobCheck(ctx context.Context, client *redis.Client, workerId string) ***REMOVED***
	// Check for job keys in the "k6:executionHistory" set
	historyIds, err := client.SMembers(ctx, "k6:executionHistory").Result()
	if err != nil ***REMOVED***
		fmt.Println("Error getting history ids", err)
	***REMOVED***

	for _, jobId := range historyIds ***REMOVED***
		checkIfCanExecute(ctx, client, jobId, workerId)
	***REMOVED***
***REMOVED***

func checkIfCanExecute(ctx context.Context, client *redis.Client, jobId string, workerId string) ***REMOVED***
	// Try to HGetAll the worker id
	job, err := client.HGetAll(ctx, jobId).Result()

	if err != nil ***REMOVED***
		fmt.Println("Error getting worker")
		return
	***REMOVED***

	// TODO: check if has capacity to execute here

	// Check worker['assignedWorker'] is nil
	if job["assignedWorker"] != "" ***REMOVED***
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
	go updateStatus(ctx, client, jobId, workerId, "ASSIGNED")
	go handleExecution(ctx, client, job, workerId)
***REMOVED***
