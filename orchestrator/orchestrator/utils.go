package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/redis/v9"
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

func getMaxJobs(standalone bool) int {
	if !standalone {
		return 5
	}

	maxJobs, err := strconv.Atoi(lib.GetEnvVariable("ORCHESTRATOR_MAX_JOBS", "1000"))
	if err != nil {
		maxJobs = 1000
	}

	return maxJobs
}

func getMaxManagedVUs(standalone bool) int64 {
	if !standalone {
		return 5000
	}

	maxManagedVUs, err := strconv.ParseInt(lib.GetEnvVariable("ORCHESTRATOR_MAX_MANAGED_VUS", "10000"), 10, 64)
	if err != nil {
		maxManagedVUs = 10000
	}

	return maxManagedVUs
}

func getLoadZones() []string {
	workerNames := []string{}

	for i := 0; i < 100; i++ {
		workerName := lib.GetEnvVariableRaw(fmt.Sprintf("WORKER_%d_DISPLAY_NAME", i), "NONE", true)
		if workerName == "NONE" {
			break
		}

		workerNames = append(workerNames, workerName)
	}

	return workerNames
}
