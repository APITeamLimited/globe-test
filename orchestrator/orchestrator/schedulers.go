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
	workerClients libOrch.WorkerClients, storeMongoDB *mongo.Database, creditsClient *redis.Client, standalone bool) ***REMOVED***

	scheduler := time.NewTicker(1 * time.Second)

	go func() ***REMOVED***
		for range scheduler.C ***REMOVED***
			orchestratorClient.HSet(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId), "lastHeartbeat", time.Now().UnixMilli())
			orchestratorClient.HSet(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId), "currentJobsCount", executionList.currentJobsCount)
			orchestratorClient.HSet(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId), "currentManagedVUsCount", executionList.currentManagedVUsCount)

			orchestratorClient.HSet(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId), "launchTime", time.Now().UnixMilli())
			orchestratorClient.HSet(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId), "maxJobs", executionList.maxJobs)
			orchestratorClient.HSet(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId), "maxManagedVUs", executionList.maxManagedVUs)

			orchestratorClient.SAdd(ctx, "orchestrators", orchestratorId)

			// Capacity may have freed up, check for queued jobs
			checkForQueuedJobs(ctx, orchestratorClient, workerClients, orchestratorId, executionList, storeMongoDB, creditsClient, standalone)
		***REMOVED***
	***REMOVED***()

***REMOVED***

// Provides management of all orchestrators, including aggregating total capacities and deleting offline nodes
func createMasterScheduler(ctx context.Context, orchestratorClient *redis.Client, workerClients libOrch.WorkerClients) ***REMOVED***
	createOrchestratorScheduler(ctx, orchestratorClient)
	createWorkerScheduler(ctx, orchestratorClient, workerClients)
***REMOVED***

// Manages orchestrator capacity statistics
func createOrchestratorScheduler(ctx context.Context, orchestratorClient *redis.Client) ***REMOVED***
	// Delete existing orchestrators
	orchestratorClient.Del(ctx, "orchestrators")

	// Find existing existingKeys using keys
	existingKeys := orchestratorClient.Keys(ctx, "orchestrator:*:info").Val()

	// Delete existing orchestratorInfo
	for _, key := range existingKeys ***REMOVED***
		orchestratorClient.Del(ctx, key)
	***REMOVED***

	scheduler := time.NewTicker(1 * time.Second)

	go func() ***REMOVED***
		for range scheduler.C ***REMOVED***
			orchestratorIds := orchestratorClient.SMembers(ctx, "orchestrators").Val()

			totalMaxJobs := 0
			totalMaxManagedVUs := 0

			totalCurrentJobsCount := 0
			totalCurrentManagedVUsCount := 0

			for _, orchestratorId := range orchestratorIds ***REMOVED***
				lastHeartbeat, err := orchestratorClient.HGet(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId), "lastHeartbeat").Int64()
				if err != nil ***REMOVED***
					fmt.Println("Error getting last heartbeat", err)
					continue
				***REMOVED***

				// If the last heartbeat was more than 10 seconds ago, delete the orchestrator
				if time.Now().UnixMilli()-lastHeartbeat > 10000 ***REMOVED***
					orchestratorClient.SRem(ctx, "orchestrators", orchestratorId)
					orchestratorClient.Del(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId))

					continue
				***REMOVED***

				// Get the capacity of the orchestrator
				maxJobs, err := orchestratorClient.HGet(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId), "maxJobs").Int()
				if err != nil ***REMOVED***
					fmt.Println("Error getting max jobs", err)
					continue
				***REMOVED***

				maxManagedVUs, err := orchestratorClient.HGet(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId), "maxManagedVUs").Int()
				if err != nil ***REMOVED***
					fmt.Println("Error getting max managed vus", err)
					continue
				***REMOVED***

				currentJobsCount, err := orchestratorClient.HGet(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId), "currentJobsCount").Int()
				if err != nil ***REMOVED***
					fmt.Println("Error getting current jobs count", err)
					continue
				***REMOVED***

				currentManagedVUsCount, err := orchestratorClient.HGet(ctx, fmt.Sprintf("orchestrator:%s:info", orchestratorId), "currentManagedVUsCount").Int()
				if err != nil ***REMOVED***
					fmt.Println("Error getting current managed vus count", err)
					continue
				***REMOVED***

				totalMaxJobs += maxJobs
				totalMaxManagedVUs += maxManagedVUs

				totalCurrentJobsCount += currentJobsCount
				totalCurrentManagedVUsCount += currentManagedVUsCount
			***REMOVED***

			// Set the total capacity
			orchestratorClient.HSet(ctx, "orchestrator:master:info", "maxJobs", totalMaxJobs)
			orchestratorClient.HSet(ctx, "orchestrator:master:info", "maxManagedVUs", totalMaxManagedVUs)

			orchestratorClient.HSet(ctx, "orchestrator:master:info", "currentJobsCount", totalCurrentJobsCount)
			orchestratorClient.HSet(ctx, "orchestrator:master:info", "currentManagedVUsCount", totalCurrentManagedVUsCount)
		***REMOVED***
	***REMOVED***()
***REMOVED***

type loadZoneSummary struct ***REMOVED***
	name                  string
	totalMaxJobs          int
	totalMaxVUs           int
	totalCurrentJobsCount int
	totalCurrentVUsCount  int
***REMOVED***

// Manages worker capacity statistics
func createWorkerScheduler(ctx context.Context, orchestratorClient *redis.Client, workerClients libOrch.WorkerClients) ***REMOVED***
	// Delete existing worker info
	orchestratorClient.Del(ctx, "workers")
	orchestratorClient.Del(ctx, "loadZones")

	for _, client := range workerClients.Clients ***REMOVED***
		client.Client.Del(ctx, "workers")

		// Find existing workerInfo using keys
		existingKeys := client.Client.Keys(ctx, "worker:*:info").Val()

		// Delete existing workerInfo
		for _, key := range existingKeys ***REMOVED***
			client.Client.Del(ctx, key)
		***REMOVED***
	***REMOVED***

	scheduler := time.NewTicker(1 * time.Second)

	go func() ***REMOVED***
		for range scheduler.C ***REMOVED***
			loadZones := make(map[string]loadZoneSummary)

			for loadZone, client := range workerClients.Clients ***REMOVED***
				clientInstance := client.Client
				zoneSummary := loadZoneSummary***REMOVED***
					name: loadZone,
				***REMOVED***

				// Aggregate the total capacity of all workers in the load zone

				workerIds := clientInstance.SMembers(ctx, "workers").Val()

				for _, workerId := range workerIds ***REMOVED***
					// Check if the worker is still alive
					lastHeartbeat, err := clientInstance.HGet(ctx, fmt.Sprintf("worker:%s:info", workerId), "lastHeartbeat").Int64()
					if err != nil ***REMOVED***
						fmt.Println("Error getting last heartbeat", err, "for worker", workerId)
						continue
					***REMOVED***

					// If the last heartbeat was more than 10 seconds ago, delete the worker
					if time.Now().UnixMilli()-lastHeartbeat > 10000 ***REMOVED***
						clientInstance.SRem(ctx, "workers", workerId)
						clientInstance.Del(ctx, fmt.Sprintf("worker:%s:info", workerId))

						continue
					***REMOVED***

					maxJobs, err := clientInstance.HGet(ctx, fmt.Sprintf("worker:%s:info", workerId), "maxJobs").Int()
					if err != nil ***REMOVED***
						fmt.Println("Error getting max jobs", err)
						continue
					***REMOVED***

					maxVUs, err := clientInstance.HGet(ctx, fmt.Sprintf("worker:%s:info", workerId), "maxVUs").Int()
					if err != nil ***REMOVED***
						fmt.Println("Error getting max vus", err)
						continue
					***REMOVED***

					currentJobsCount, err := clientInstance.HGet(ctx, fmt.Sprintf("worker:%s:info", workerId), "currentJobsCount").Int()
					if err != nil ***REMOVED***
						fmt.Println("Error getting current jobs count", err)
						continue
					***REMOVED***

					currentVUsCount, err := clientInstance.HGet(ctx, fmt.Sprintf("worker:%s:info", workerId), "currentVUsCount").Int()
					if err != nil ***REMOVED***
						fmt.Println("Error getting current vus count", err)
						continue
					***REMOVED***

					zoneSummary.totalMaxJobs += maxJobs
					zoneSummary.totalMaxVUs += maxVUs

					zoneSummary.totalCurrentJobsCount += currentJobsCount
					zoneSummary.totalCurrentVUsCount += currentVUsCount

					// Realay worker stats to orchestrator
					workerStats := clientInstance.HGetAll(ctx, fmt.Sprintf("worker:%s:info", workerId)).Val()

					for key, value := range workerStats ***REMOVED***
						orchestratorClient.HSet(ctx, fmt.Sprintf("worker:%s:info", workerId), key, value)
					***REMOVED***
				***REMOVED***

				// Add the load zone summary to the map
				loadZones[zoneSummary.name] = zoneSummary
			***REMOVED***

			// Calculate the global capacity as the sum of all load zones
			globalZone := loadZoneSummary***REMOVED***
				name: libOrch.GlobalName,
			***REMOVED***

			for _, zone := range loadZones ***REMOVED***
				globalZone.totalMaxJobs += zone.totalMaxJobs
				globalZone.totalMaxVUs += zone.totalMaxVUs

				globalZone.totalCurrentJobsCount += zone.totalCurrentJobsCount
				globalZone.totalCurrentVUsCount += zone.totalCurrentVUsCount
			***REMOVED***

			loadZones[globalZone.name] = globalZone

			loadZoneNames := make([]string, 0, len(loadZones))
			for loadZone := range loadZones ***REMOVED***
				loadZoneNames = append(loadZoneNames, loadZone)
			***REMOVED***

			existingLoadZones := orchestratorClient.SMembers(ctx, "loadZones").Val()

			// Find the load zones that have been removed in loadZoneNames
			removedLoadZones := make([]string, 0)

			for _, loadZone := range existingLoadZones ***REMOVED***
				found := false

				for _, newLoadZone := range loadZoneNames ***REMOVED***
					if loadZone == newLoadZone ***REMOVED***
						found = true
						break
					***REMOVED***
				***REMOVED***

				if !found ***REMOVED***
					removedLoadZones = append(removedLoadZones, loadZone)
				***REMOVED***
			***REMOVED***

			// Remove the load zones that have been removed
			if len(removedLoadZones) > 0 ***REMOVED***
				orchestratorClient.SRem(ctx, "loadZones", removedLoadZones)
			***REMOVED***

			orchestratorClient.SAdd(ctx, "loadZones", loadZoneNames)

			// Dispatch updates for each load zone
			for _, zone := range loadZones ***REMOVED***
				orchestratorClient.HSet(ctx, fmt.Sprintf("loadZone:%s:info", zone.name), "maxJobs", zone.totalMaxJobs)
				orchestratorClient.HSet(ctx, fmt.Sprintf("loadZone:%s:info", zone.name), "maxVUs", zone.totalMaxVUs)

				orchestratorClient.HSet(ctx, fmt.Sprintf("loadZone:%s:info", zone.name), "currentJobsCount", zone.totalCurrentJobsCount)
				orchestratorClient.HSet(ctx, fmt.Sprintf("loadZone:%s:info", zone.name), "currentVUsCount", zone.totalCurrentVUsCount)
			***REMOVED***
		***REMOVED***
	***REMOVED***()
***REMOVED***
