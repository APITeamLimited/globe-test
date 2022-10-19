package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func fetchJob(ctx context.Context, orchestratorClient *redis.Client, jobId string) (*libOrch.Job, error) ***REMOVED***
	jobRaw, err := orchestratorClient.HGet(ctx, jobId, "job").Result()

	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Check job not empty
	if jobRaw == "" ***REMOVED***
		return nil, fmt.Errorf("job %s is empty", jobId)
	***REMOVED***

	job := libOrch.Job***REMOVED******REMOVED***
	// Parse job as libOrch.Job
	err = json.Unmarshal([]byte(jobRaw), &job)
	if err != nil ***REMOVED***
		fmt.Println("error unmarshalling job", err)
		return nil, fmt.Errorf("error unmarshalling job %s", jobId)
	***REMOVED***

	// Sensitive field, ensure it is nil
	job.Options = nil

	return &job, nil
***REMOVED***

// Deletes old orchestrator info periodically about offline nodes
func createDeletionScheduler(ctx context.Context, orchestratorClient *redis.Client, workerClients libOrch.WorkerClients) ***REMOVED***
	deletionScheduler := time.NewTicker(10 * time.Second)

	go func() ***REMOVED***
		for range deletionScheduler.C ***REMOVED***
			orchestratorIds := orchestratorClient.SMembers(ctx, "orchestrators").Val()

			for _, orchestratorId := range orchestratorIds ***REMOVED***
				// Check if the orchestrator is still alive
				lastHeartbeat, err := orchestratorClient.HGet(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId), "lastHeartbeat").Int64()
				if err != nil ***REMOVED***
					fmt.Println("Error getting last heartbeat", err)
					continue
				***REMOVED***

				// If the last heartbeat was more than 10 seconds ago, delete the orchestrator
				if time.Now().UnixMilli()-lastHeartbeat > 10000 ***REMOVED***
					orchestratorClient.SRem(ctx, "orchestrators", orchestratorId)
					orchestratorClient.Del(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId))
				***REMOVED***
			***REMOVED***

		***REMOVED***
	***REMOVED***()
***REMOVED***

func startJobScheduling(ctx context.Context, orchestratorClient *redis.Client, orchestratorId string, executionList *ExecutionList,
	workerClients libOrch.WorkerClients, storeMongoDB *mongo.Database) ***REMOVED***
	orchestratorClient.SAdd(ctx, "orchestrators", orchestratorId)

	jobsCheckScheduler := time.NewTicker(1 * time.Second)

	go func() ***REMOVED***
		for range jobsCheckScheduler.C ***REMOVED***
			orchestratorClient.HSet(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId), "lastHeartbeat", time.Now().UnixMilli())
			orchestratorClient.HSet(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId), "currentJobsCount", executionList.currentJobsCount)
			orchestratorClient.HSet(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId), "currentManagedVUs", executionList.currentManagedVUsCount)

			// Capacity may have freed up, check for queued jobs
			checkForQueuedJobs(ctx, orchestratorClient, workerClients, orchestratorId, executionList, storeMongoDB)
		***REMOVED***
	***REMOVED***()

	orchestratorClient.HSet(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId), "launchTime", time.Now().UnixMilli())
	orchestratorClient.HSet(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId), "maxJobs", executionList.maxJobs)
***REMOVED***

func getOrchestratorClient() *redis.Client ***REMOVED***
	return redis.NewClient(&redis.Options***REMOVED***
		Addr:     fmt.Sprintf("%s:%s", libOrch.GetEnvVariable("ORCHESTRATOR_REDIS_HOST", "localhost"), libOrch.GetEnvVariable("ORCHESTRATOR_REDIS_PORT", "10000")),
		Username: "default",
		Password: libOrch.GetEnvVariable("ORCHESTRATOR_REDIS_PASSWORD", ""),
	***REMOVED***)
***REMOVED***

func getMaxJobs() int ***REMOVED***
	maxJobs, err := strconv.Atoi(libOrch.GetEnvVariable("ORCHESTRATOR_MAX_JOBS", "1000"))
	if err != nil ***REMOVED***
		maxJobs = 1000
	***REMOVED***

	return maxJobs
***REMOVED***

func getMaxManagedVUs() int64 ***REMOVED***
	maxManagedVUs, err := strconv.ParseInt(libOrch.GetEnvVariable("ORCHESTRATOR_MAX_MANAGED_VUS", "10000"), 10, 64)
	if err != nil ***REMOVED***
		maxManagedVUs = 10000
	***REMOVED***

	return maxManagedVUs
***REMOVED***
