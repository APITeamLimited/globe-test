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

func fetchJob(ctx context.Context, orchestratorClient *redis.Client, jobId string) (*libOrch.Job, error) {
	jobRaw, err := orchestratorClient.HGet(ctx, jobId, "job").Result()

	if err != nil {
		return nil, err
	}

	// Check job not empty
	if jobRaw == "" {
		return nil, fmt.Errorf("job %s is empty", jobId)
	}

	job := libOrch.Job{}
	// Parse job as libOrch.Job
	err = json.Unmarshal([]byte(jobRaw), &job)
	if err != nil {
		fmt.Println("error unmarshalling job", err)
		return nil, fmt.Errorf("error unmarshalling job %s", jobId)
	}

	// Sensitive field, ensure it is nil
	job.Options = nil

	return &job, nil
}

// Deletes old orchestrator info periodically about offline nodes
func createDeletionScheduler(ctx context.Context, orchestratorClient *redis.Client, workerClients libOrch.WorkerClients) {
	deletionScheduler := time.NewTicker(10 * time.Second)

	go func() {
		for range deletionScheduler.C {
			orchestratorIds := orchestratorClient.SMembers(ctx, "orchestrators").Val()

			for _, orchestratorId := range orchestratorIds {
				// Check if the orchestrator is still alive
				lastHeartbeat, err := orchestratorClient.HGet(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId), "lastHeartbeat").Int64()
				if err != nil {
					fmt.Println("Error getting last heartbeat", err)
					continue
				}

				// If the last heartbeat was more than 10 seconds ago, delete the orchestrator
				if time.Now().UnixMilli()-lastHeartbeat > 10000 {
					orchestratorClient.SRem(ctx, "orchestrators", orchestratorId)
					orchestratorClient.Del(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId))
				}
			}

		}
	}()
}

func startJobScheduling(ctx context.Context, orchestratorClient *redis.Client, orchestratorId string, executionList *ExecutionList,
	workerClients libOrch.WorkerClients, storeMongoDB *mongo.Database) {
	orchestratorClient.SAdd(ctx, "orchestrators", orchestratorId)

	jobsCheckScheduler := time.NewTicker(1 * time.Second)

	go func() {
		for range jobsCheckScheduler.C {
			orchestratorClient.HSet(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId), "lastHeartbeat", time.Now().UnixMilli())
			orchestratorClient.HSet(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId), "currentJobsCount", executionList.currentJobsCount)
			orchestratorClient.HSet(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId), "currentManagedVUs", executionList.currentManagedVUsCount)

			// Capacity may have freed up, check for queued jobs
			checkForQueuedJobs(ctx, orchestratorClient, workerClients, orchestratorId, executionList, storeMongoDB)
		}
	}()

	orchestratorClient.HSet(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId), "launchTime", time.Now().UnixMilli())
	orchestratorClient.HSet(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId), "maxJobs", executionList.maxJobs)
}

func getOrchestratorClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", libOrch.GetEnvVariable("ORCHESTRATOR_REDIS_HOST", "localhost"), libOrch.GetEnvVariable("ORCHESTRATOR_REDIS_PORT", "10000")),
		Username: "default",
		Password: libOrch.GetEnvVariable("ORCHESTRATOR_REDIS_PASSWORD", ""),
	})
}

func getMaxJobs() int {
	maxJobs, err := strconv.Atoi(libOrch.GetEnvVariable("ORCHESTRATOR_MAX_JOBS", "1000"))
	if err != nil {
		maxJobs = 1000
	}

	return maxJobs
}

func getMaxManagedVUs() int64 {
	maxManagedVUs, err := strconv.ParseInt(libOrch.GetEnvVariable("ORCHESTRATOR_MAX_MANAGED_VUS", "10000"), 10, 64)
	if err != nil {
		maxManagedVUs = 10000
	}

	return maxManagedVUs
}
