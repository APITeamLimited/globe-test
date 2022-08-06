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

func Run() ***REMOVED***
	ctx := context.Background()

	nodeId := uuid.New()

	client := redis.NewClient(&redis.Options***REMOVED***
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	***REMOVED***)

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

	for msg := range channel ***REMOVED***
		checkIfCanExecute(ctx, client, msg, nodeId.String())
	***REMOVED***
***REMOVED***

func checkIfCanExecute(ctx context.Context, client *redis.Client, msg *redis.Message, nodeId string) ***REMOVED***
	// Check if redis message is a uuid
	jobId, err := uuid.Parse(msg.Payload)
	if err != nil ***REMOVED***
		fmt.Println("Not a uuid")
		return
	***REMOVED***

	// Try to HGetAll the node id

	job, err := client.HGetAll(ctx, jobId.String()).Result()

	if err != nil ***REMOVED***
		fmt.Println("Error getting node")
		return
	***REMOVED***

	fmt.Println("job", job, job["assignedNode"])

	// TODO: check if has capacity to execute here

	// Check node['assignedNode'] is nil
	if job["assignedNode"] != "" ***REMOVED***
		fmt.Println("Node already assigned")
		return
	***REMOVED***

	// HSetNX assignedNode to the nodeId
	assignmentResult, err := client.HSetNX(ctx, jobId.String(), "assignedNode", nodeId).Result()

	if err != nil ***REMOVED***
		fmt.Println("Error setting node")
		return
	***REMOVED***

	// If result is 0, node is already assigned
	if !assignmentResult ***REMOVED***
		fmt.Println("Node already assigned")
		return
	***REMOVED***

	// We got the job
	handleExecution(ctx, client, job, nodeId)
***REMOVED***

func handleExecution(ctx context.Context,
	client *redis.Client, job map[string]string, nodeId string) ***REMOVED***
	// Check if redis message is a uuid

	fmt.Println("\033[1;32mGot job", job["id"], "\033[0m")

	client.Publish(ctx, fmt.Sprintf("k6:executionUpdates:%s", job["id"]), fmt.Sprintf("Node %s: Execution Started", nodeId))

	jobScript := job["script"]

	if jobScript == "" ***REMOVED***
		fmt.Println("No script")
		client.Publish(ctx, fmt.Sprintf("k6:executionUpdates:%s", job["id"]), fmt.Sprintf("Node %s: Execution Error, job not found", nodeId))
		return
	***REMOVED***

	// Run the job

***REMOVED***
