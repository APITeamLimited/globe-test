package node

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

	nodeId := uuid.New()

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	currentTime := time.Now().UnixMilli()

	//Set the node id and current time
	client.HSet(ctx, "k6:nodes", nodeId.String(), currentTime)

	c := color.New(color.FgCyan)
	fmt.Print("\n", c.Sprint(consts.Banner()), "\n\n")
	fmt.Print("\033[1;35mStarting node", nodeId, "\033[0m\n\n")

	fmt.Printf("Listening for new jobs on %s...\n", client.Options().Addr)

	// Subscribe to the execution channel
	psc := client.Subscribe(ctx, "k6:execution")

	channel := psc.Channel()

	for msg := range channel {
		checkIfCanExecute(ctx, client, msg, nodeId.String())
	}
}

func checkIfCanExecute(ctx context.Context, client *redis.Client, msg *redis.Message, nodeId string) {
	// Check if redis message is a uuid
	jobId, err := uuid.Parse(msg.Payload)
	if err != nil {
		fmt.Println("Not a uuid")
		return
	}

	// Try to HGetAll the node id

	job, err := client.HGetAll(ctx, jobId.String()).Result()

	if err != nil {
		fmt.Println("Error getting node")
		return
	}

	// TODO: check if has capacity to execute here

	// Check node['assignedNode'] is nil
	if job["assignedNode"] != "" {
		fmt.Println("Node already assigned")
		return
	}

	// HSetNX assignedNode to the nodeId
	assignmentResult, err := client.HSetNX(ctx, jobId.String(), "assignedNode", nodeId).Result()

	if err != nil {
		fmt.Println("Error setting node")
		return
	}

	// If result is 0, node is already assigned
	if !assignmentResult {
		fmt.Println("Node already assigned")
		return
	}

	// We got the job
	updateStatus(ctx, client, jobId.String(), nodeId, "ASSIGNED")
	go handleExecution(ctx, client, job, nodeId)
}
