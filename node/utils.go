package node

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
)

func dispatchMessage(ctx context.Context, client *redis.Client, jobId string, nodeId string, message string) ***REMOVED***
	timestamp := time.Now().UnixMilli()
	stampedTag := fmt.Sprintf("%d:%s", timestamp, nodeId)

	// Update main job
	updatesKey := fmt.Sprintf("%s:updates", jobId)
	client.HSet(ctx, updatesKey, stampedTag, message)

	// Dispatch to channel
	publishTaggedMessage := fmt.Sprintf("%s: %s", stampedTag, message)
	client.Publish(ctx, fmt.Sprintf("k6:executionUpdates:%s", jobId), publishTaggedMessage)
***REMOVED***

func updateStatus(ctx context.Context, client *redis.Client, jobId string, nodeId string, status string) ***REMOVED***
	client.HSet(ctx, jobId, "status", status)
	dispatchMessage(ctx, client, jobId, nodeId, status)
***REMOVED***

func handleError(ctx context.Context, client *redis.Client, jobId string, nodeId string, err error) ***REMOVED***
	dispatchMessage(ctx, client, jobId, nodeId, err.Error())
	updateStatus(ctx, client, jobId, nodeId, "FAILED")
***REMOVED***
