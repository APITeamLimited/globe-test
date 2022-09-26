package worker

import (
	"context"
	"encoding/json"

	"github.com/APITeamLimited/redis/v9"
	"go.k6.io/k6/lib"
)

func loadWorkerInfo(ctx context.Context,
	client *redis.Client, job map[string]string, workerId string) (*lib.WorkerInfo, error) ***REMOVED***

	workerInfo := &lib.WorkerInfo***REMOVED***
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

func parseJobEnvironment(workerInfo *lib.WorkerInfo, job map[string]string) error ***REMOVED***
	// Check environmentContext actually exists in the job
	if job["environmentContext"] != "" ***REMOVED***
		// Parse the environmentContext json
		enviromentContext := []lib.KeyValueItem***REMOVED******REMOVED***
		err := json.Unmarshal([]byte(job["environmentContext"]), &enviromentContext)

		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// Init map, need to assign it first to get an address
		workerInfoEnvironment := make(map[string]lib.KeyValueItem, len(enviromentContext))
		workerInfo.Environment = &workerInfoEnvironment
	***REMOVED***

	return nil
***REMOVED***

type parseCollectionContext struct ***REMOVED***
	Variables map[string]lib.KeyValueItem `json:"variables"`
***REMOVED***

func parseJobCollection(workerInfo *lib.WorkerInfo, job map[string]string) error ***REMOVED***
	// Check collectionContext actually exists in the job
	if job["collectionContext"] != "" ***REMOVED***
		// Parse the collectionContext json
		collectionContext := parseCollectionContext***REMOVED******REMOVED***
		err := json.Unmarshal([]byte(job["collectionContext"]), &collectionContext)

		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// Init map, need to assign it first to get an address
		workerInfoVariables := make(map[string]lib.KeyValueItem, len(collectionContext.Variables))

		collection := lib.Collection***REMOVED***
			Variables: &workerInfoVariables,
		***REMOVED***

		workerInfo.Collection = &collection
	***REMOVED***

	return nil
***REMOVED***
