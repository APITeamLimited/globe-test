package worker

import (
	"context"
	"encoding/json"

	"github.com/APITeamLimited/redis/v9"
	"go.k6.io/k6/lib"
)

func loadWorkerInfo(ctx context.Context,
	client *redis.Client, job map[string]string, workerId string) (*lib.WorkerInfo, error) {

	workerInfo := &lib.WorkerInfo{
		Client:         client,
		JobId:          job["id"],
		ScopeId:        job["scopeId"],
		OrchestratorId: job["orchestratorId"],
		WorkerId:       workerId,
		Ctx:            ctx,
	}

	err := parseJobEnvironment(workerInfo, job)
	if err != nil {
		return nil, err
	}

	err = parseJobCollection(workerInfo, job)
	if err != nil {
		return nil, err
	}

	return workerInfo, nil
}

func parseJobEnvironment(workerInfo *lib.WorkerInfo, job map[string]string) error {
	// Check environmentContext actually exists in the job
	if job["environmentContext"] != "" {
		// Parse the environmentContext json
		enviromentContext := []lib.KeyValueItem{}
		err := json.Unmarshal([]byte(job["environmentContext"]), &enviromentContext)

		if err != nil {
			return err
		}

		// Init map, need to assign it first to get an address
		workerInfoEnvironment := make(map[string]lib.KeyValueItem, len(enviromentContext))
		workerInfo.Environment = &workerInfoEnvironment
	}

	return nil
}

type parseCollectionContext struct {
	Variables map[string]lib.KeyValueItem `json:"variables"`
}

func parseJobCollection(workerInfo *lib.WorkerInfo, job map[string]string) error {
	// Check collectionContext actually exists in the job
	if job["collectionContext"] != "" {
		// Parse the collectionContext json
		collectionContext := parseCollectionContext{}
		err := json.Unmarshal([]byte(job["collectionContext"]), &collectionContext)

		if err != nil {
			return err
		}

		// Init map, need to assign it first to get an address
		workerInfoVariables := make(map[string]lib.KeyValueItem, len(collectionContext.Variables))

		collection := lib.Collection{
			Variables: &workerInfoVariables,
		}

		workerInfo.Collection = &collection
	}

	return nil
}
