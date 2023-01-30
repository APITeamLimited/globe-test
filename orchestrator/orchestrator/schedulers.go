package orchestrator

import (
	"context"
	"fmt"
	"time"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func startJobScheduling(ctx context.Context, orchestratorClient *redis.Client, orchestratorId string, executionList *ExecutionList, storeMongoDB *mongo.Database, creditsClient *redis.Client, standalone bool, funcAuthClient libOrch.RunAuthClient, loadZones []string) {

	scheduler := time.NewTicker(1 * time.Second)

	go func() {
		for range scheduler.C {
			orchestratorClient.HSet(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId), "lastHeartbeat", time.Now().UnixMilli())
			orchestratorClient.HSet(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId), "currentJobsCount", executionList.currentJobsCount)
			orchestratorClient.HSet(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId), "currentManagedVUsCount", executionList.currentManagedVUsCount)

			orchestratorClient.HSet(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId), "launchTime", time.Now().UnixMilli())
			orchestratorClient.HSet(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId), "maxJobs", executionList.maxJobs)
			orchestratorClient.HSet(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId), "maxManagedVUs", executionList.maxManagedVUs)

			orchestratorClient.SAdd(ctx, "orchestrators", orchestratorId)

			// Capacity may have freed up, check for queued jobs
			checkForQueuedJobs(ctx, orchestratorClient, orchestratorId, executionList, storeMongoDB, creditsClient, standalone, funcAuthClient, loadZones)
		}
	}()

}

// Provides management of all orchestrators, including aggregating total capacities and deleting offline nodes
func createMasterScheduler(ctx context.Context, orchestratorClient *redis.Client) {
	// Delete existing orchestrators
	orchestratorClient.Del(ctx, "orchestrators")

	// Find existing existingKeys using keys
	existingKeys := orchestratorClient.Keys(ctx, "orchestrator:*:info").Val()

	// Delete existing orchestratorInfo
	for _, key := range existingKeys {
		orchestratorClient.Del(ctx, key)
	}

	scheduler := time.NewTicker(1 * time.Second)

	go func() {
		for range scheduler.C {
			orchestratorIds := orchestratorClient.SMembers(ctx, "orchestrators").Val()

			totalMaxJobs := 0
			totalMaxManagedVUs := 0

			totalCurrentJobsCount := 0
			totalCurrentManagedVUsCount := 0

			for _, orchestratorId := range orchestratorIds {
				lastHeartbeat, err := orchestratorClient.HGet(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId), "lastHeartbeat").Int64()
				if err != nil {
					fmt.Println("Error getting last heartbeat", err)
					continue
				}

				// If the last heartbeat was more than 10 seconds ago, delete the orchestrator
				if time.Now().UnixMilli()-lastHeartbeat > 10000 {
					orchestratorClient.SRem(ctx, "orchestrators", orchestratorId)
					orchestratorClient.Del(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId))

					continue
				}

				// Get the capacity of the orchestrator
				maxJobs, err := orchestratorClient.HGet(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId), "maxJobs").Int()
				if err != nil {
					fmt.Println("Error getting max jobs", err)
					continue
				}

				maxManagedVUs, err := orchestratorClient.HGet(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId), "maxManagedVUs").Int()
				if err != nil {
					fmt.Println("Error getting max managed vus", err)
					continue
				}

				currentJobsCount, err := orchestratorClient.HGet(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId), "currentJobsCount").Int()
				if err != nil {
					fmt.Println("Error getting current jobs count", err)
					continue
				}

				currentManagedVUsCount, err := orchestratorClient.HGet(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId), "currentManagedVUsCount").Int()
				if err != nil {
					fmt.Println("Error getting current managed vus count", err)
					continue
				}

				totalMaxJobs += maxJobs
				totalMaxManagedVUs += maxManagedVUs

				totalCurrentJobsCount += currentJobsCount
				totalCurrentManagedVUsCount += currentManagedVUsCount
			}

			// Set the total capacity
			orchestratorClient.HSet(ctx, "orchestrator:master:info", "maxJobs", totalMaxJobs)
			orchestratorClient.HSet(ctx, "orchestrator:master:info", "maxManagedVUs", totalMaxManagedVUs)

			orchestratorClient.HSet(ctx, "orchestrator:master:info", "currentJobsCount", totalCurrentJobsCount)
			orchestratorClient.HSet(ctx, "orchestrator:master:info", "currentManagedVUsCount", totalCurrentManagedVUsCount)
		}
	}()
}
