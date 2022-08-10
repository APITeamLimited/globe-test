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

func Run() {
	ctx := context.Background()

	workerId := uuid.New()

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

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

	for msg := range channel {
		jobId, err := uuid.Parse(msg.Payload)
		if err != nil {
			fmt.Println("Error, got did not parse job id")
			return
		}
		checkIfCanExecute(ctx, client, jobId.String(), workerId.String())
	}
}

/*
Check for jobs existing on startup
*/
func startupJobCheck(ctx context.Context, client *redis.Client, workerId string) {
	// Check for job keys in the "k6:executionHistory" set
	historyIds, err := client.SMembers(ctx, "k6:executionHistory").Result()
	if err != nil {
		fmt.Println("Error getting history ids", err)
	}

	for _, jobId := range historyIds {
		checkIfCanExecute(ctx, client, jobId, workerId)
	}
}

func checkIfCanExecute(ctx context.Context, client *redis.Client, jobId string, workerId string) {
	// Try to HGetAll the worker id
	job, err := client.HGetAll(ctx, jobId).Result()

	if err != nil {
		fmt.Println("Error getting worker")
		return
	}

	// TODO: check if has capacity to execute here

	// Check worker['assignedWorker'] is nil
	if job["assignedWorker"] != "" {
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
	go updateStatus(ctx, client, jobId, workerId, "ASSIGNED")
	go handleExecution(ctx, client, job, workerId)
}
