package worker

import (
	"context"
	"encoding/json"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/redis/v9"
)

func loadWorkerInfo(ctx context.Context,
	client *redis.Client, job map[string]string, workerId string) (*libWorker.WorkerInfo, error) ***REMOVED***

	workerInfo := &libWorker.WorkerInfo***REMOVED***
		Client:         client,
		JobId:          job["id"],
		ScopeId:        job["scopeId"],
		OrchestratorId: job["orchestratorId"],
		WorkerId:       workerId,
		Ctx:            ctx,
	***REMOVED***

	err := parseJobEnvironment(workerInfo, job)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	err = parseJobCollection(workerInfo, job)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return workerInfo, nil
***REMOVED***

func parseJobEnvironment(workerInfo *libWorker.WorkerInfo, job map[string]string) error ***REMOVED***
	// Check environmentContext actually exists in the job
	if job["environmentContext"] != "" ***REMOVED***
		// Parse the environmentContext json
		enviromentContext := []libWorker.KeyValueItem***REMOVED******REMOVED***
		err := json.Unmarshal([]byte(job["environmentContext"]), &enviromentContext)

		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// Init map, need to assign it first to get an address
		workerInfoEnvironment := make(map[string]libWorker.KeyValueItem, len(enviromentContext))
		workerInfo.Environment = &workerInfoEnvironment
	***REMOVED***

	return nil
***REMOVED***

type parseCollectionContext struct ***REMOVED***
	Variables map[string]libWorker.KeyValueItem `json:"variables"`
***REMOVED***

func parseJobCollection(workerInfo *libWorker.WorkerInfo, job map[string]string) error ***REMOVED***
	// Check collectionContext actually exists in the job
	if job["collectionContext"] != "" ***REMOVED***
		// Parse the collectionContext json
		collectionContext := parseCollectionContext***REMOVED******REMOVED***
		err := json.Unmarshal([]byte(job["collectionContext"]), &collectionContext)

		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// Init map, need to assign it first to get an address
		workerInfoVariables := make(map[string]libWorker.KeyValueItem, len(collectionContext.Variables))

		collection := libWorker.Collection***REMOVED***
			Variables: &workerInfoVariables,
		***REMOVED***

		workerInfo.Collection = &collection
	***REMOVED***

	return nil
***REMOVED***
