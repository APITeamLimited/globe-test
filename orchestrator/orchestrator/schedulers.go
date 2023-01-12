package orchestrator

import (
	"context"
	"fmt"
	"time"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func startJobScheduling(ctx context.Context, orchestratorClient *redis.Client, orchestratorId string, executionList *ExecutionList,
	workerClients libOrch.WorkerClients, storeMongoDB *mongo.Database, creditsClient *redis.Client, standalone bool, funcAuthClient libOrch.FunctionAuthClient, funcMode, independentWorkerRedisHosts bool) {

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
			checkForQueuedJobs(ctx, orchestratorClient, workerClients, orchestratorId, executionList, storeMongoDB, creditsClient, standalone, funcAuthClient, funcMode, independentWorkerRedisHosts)
		}
	}()

}

// Provides management of all orchestrators, including aggregating total capacities and deleting offline nodes
func createMasterScheduler(ctx context.Context, orchestratorClient *redis.Client, workerClients libOrch.WorkerClients, funcMode bool) {
	createOrchestratorScheduler(ctx, orchestratorClient)

	if !funcMode {
		createWorkerScheduler(ctx, orchestratorClient, workerClients)
	}
}

// Manages orchestrator capacity statistics
func createOrchestratorScheduler(ctx context.Context, orchestratorClient *redis.Client) {
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

type loadZoneSummary struct {
	name                  string
	totalMaxJobs          int
	totalMaxVUs           int
	totalCurrentJobsCount int
	totalCurrentVUsCount  int
}

// Manages worker capacity statistics
func createWorkerScheduler(ctx context.Context, orchestratorClient *redis.Client, workerClients libOrch.WorkerClients) {
	// Delete existing worker info
	orchestratorClient.Del(ctx, "workers")
	orchestratorClient.Del(ctx, "loadZones")

	for _, client := range workerClients.Clients {
		client.Client.Del(ctx, "workers")

		// Find existing workerInfo using keys
		existingKeys := client.Client.Keys(ctx, "worker:*:info").Val()

		// Delete existing workerInfo
		for _, key := range existingKeys {
			client.Client.Del(ctx, key)
		}
	}

	scheduler := time.NewTicker(1 * time.Second)

	go func() {
		for range scheduler.C {
			loadZones := make(map[string]loadZoneSummary)

			for loadZone, client := range workerClients.Clients {
				clientInstance := client.Client
				zoneSummary := loadZoneSummary{
					name: loadZone,
				}

				// Aggregate the total capacity of all workers in the load zone

				workerIds := clientInstance.SMembers(ctx, "workers").Val()

				for _, workerId := range workerIds {
					// Check if the worker is still alive
					lastHeartbeat, err := clientInstance.HGet(ctx, fmt.Sprintf("worker:%s:info", workerId), "lastHeartbeat").Int64()
					if err != nil {
						fmt.Println("Error getting last heartbeat", err, "for worker", workerId)
						continue
					}

					// If the last heartbeat was more than 10 seconds ago, delete the worker
					if time.Now().UnixMilli()-lastHeartbeat > 10000 {
						clientInstance.SRem(ctx, "workers", workerId)
						clientInstance.Del(ctx, fmt.Sprintf("worker:%s:info", workerId))

						continue
					}

					maxJobs, err := clientInstance.HGet(ctx, fmt.Sprintf("worker:%s:info", workerId), "maxJobs").Int()
					if err != nil {
						fmt.Println("Error getting max jobs", err)
						continue
					}

					maxVUs, err := clientInstance.HGet(ctx, fmt.Sprintf("worker:%s:info", workerId), "maxVUs").Int()
					if err != nil {
						fmt.Println("Error getting max vus", err)
						continue
					}

					currentJobsCount, err := clientInstance.HGet(ctx, fmt.Sprintf("worker:%s:info", workerId), "currentJobsCount").Int()
					if err != nil {
						fmt.Println("Error getting current jobs count", err)
						continue
					}

					currentVUsCount, err := clientInstance.HGet(ctx, fmt.Sprintf("worker:%s:info", workerId), "currentVUsCount").Int()
					if err != nil {
						fmt.Println("Error getting current vus count", err)
						continue
					}

					zoneSummary.totalMaxJobs += maxJobs
					zoneSummary.totalMaxVUs += maxVUs

					zoneSummary.totalCurrentJobsCount += currentJobsCount
					zoneSummary.totalCurrentVUsCount += currentVUsCount

					// Realay worker stats to orchestrator
					workerStats := clientInstance.HGetAll(ctx, fmt.Sprintf("worker:%s:info", workerId)).Val()

					for key, value := range workerStats {
						orchestratorClient.HSet(ctx, fmt.Sprintf("worker:%s:info", workerId), key, value)
					}
				}

				// Add the load zone summary to the map
				loadZones[zoneSummary.name] = zoneSummary
			}

			// Calculate the global capacity as the sum of all load zones
			globalZone := loadZoneSummary{
				name: libOrch.GlobalName,
			}

			for _, zone := range loadZones {
				globalZone.totalMaxJobs += zone.totalMaxJobs
				globalZone.totalMaxVUs += zone.totalMaxVUs

				globalZone.totalCurrentJobsCount += zone.totalCurrentJobsCount
				globalZone.totalCurrentVUsCount += zone.totalCurrentVUsCount
			}

			loadZones[globalZone.name] = globalZone

			loadZoneNames := make([]string, 0, len(loadZones))
			for loadZone := range loadZones {
				loadZoneNames = append(loadZoneNames, loadZone)
			}

			existingLoadZones := orchestratorClient.SMembers(ctx, "loadZones").Val()

			// Find the load zones that have been removed in loadZoneNames
			removedLoadZones := make([]string, 0)

			for _, loadZone := range existingLoadZones {
				found := false

				for _, newLoadZone := range loadZoneNames {
					if loadZone == newLoadZone {
						found = true
						break
					}
				}

				if !found {
					removedLoadZones = append(removedLoadZones, loadZone)
				}
			}

			// Remove the load zones that have been removed
			if len(removedLoadZones) > 0 {
				orchestratorClient.SRem(ctx, "loadZones", removedLoadZones)
			}

			orchestratorClient.SAdd(ctx, "loadZones", loadZoneNames)

			// Dispatch updates for each load zone
			for _, zone := range loadZones {
				orchestratorClient.HSet(ctx, fmt.Sprintf("loadZone:%s:info", zone.name), "maxJobs", zone.totalMaxJobs)
				orchestratorClient.HSet(ctx, fmt.Sprintf("loadZone:%s:info", zone.name), "maxVUs", zone.totalMaxVUs)

				orchestratorClient.HSet(ctx, fmt.Sprintf("loadZone:%s:info", zone.name), "currentJobsCount", zone.totalCurrentJobsCount)
				orchestratorClient.HSet(ctx, fmt.Sprintf("loadZone:%s:info", zone.name), "currentVUsCount", zone.totalCurrentVUsCount)
			}
		}
	}()
}
