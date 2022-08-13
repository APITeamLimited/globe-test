package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
)

type Message struct ***REMOVED***
	JobId       string `json:"jobId"`
	Time        int64  `json:"time"`
	WorkerId    string `json:"workerId"`
	Message     string `json:"message"`
	MessageType string `json:"messageType"`
***REMOVED***

func dispatchMessage(ctx context.Context, client *redis.Client, jobId string, workerId string, message string, messageType string) ***REMOVED***
	timestamp := time.Now().UnixMilli()
	stampedTag := fmt.Sprintf("%d:%s", timestamp, workerId)

	fmt.Println("Dispatching message: ", message)

	var messageStruct = Message***REMOVED***
		JobId:       jobId,
		Time:        timestamp,
		WorkerId:    workerId,
		Message:     message,
		MessageType: messageType,
	***REMOVED***

	messageJson, err := json.Marshal(messageStruct)
	if err != nil ***REMOVED***
		fmt.Println("Error marshalling message")
		return
	***REMOVED***

	// Update main job
	updatesKey := fmt.Sprintf("%s:updates", jobId)
	client.HSet(ctx, updatesKey, stampedTag, messageJson)

	// Dispatch to channel
	client.Publish(ctx, fmt.Sprintf("k6:executionUpdates:%s", jobId), messageJson)
***REMOVED***

func updateStatus(ctx context.Context, client *redis.Client, jobId string, workerId string, status string) ***REMOVED***
	client.HSet(ctx, jobId, "status", status)
	dispatchMessage(ctx, client, jobId, workerId, status, "STATUS")
***REMOVED***

func handleStringError(ctx context.Context, client *redis.Client, jobId string, workerId string, errString string) ***REMOVED***
	dispatchMessage(ctx, client, jobId, workerId, errString, "ERROR")
	updateStatus(ctx, client, jobId, workerId, "FAILED")
***REMOVED***

func handleError(ctx context.Context, client *redis.Client, jobId string, workerId string, err error) ***REMOVED***
	dispatchMessage(ctx, client, jobId, workerId, err.Error(), "ERROR")
	updateStatus(ctx, client, jobId, workerId, "FAILED")
***REMOVED***

// Splits an outputted console.log message into a message type and message
//func splitLogMessage(message string) (string[] string[]) ***REMOVED***
//	var level = ""
//
//	if (string.Contains("level=info") ***REMOVED***
//		level = "INFO"
//	***REMOVED*** e
//***REMOVED***
