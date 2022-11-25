package orchestrator

import (
	"context"
	"crypto/tls"
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

func getOrchestratorOrchestratorClient() *redis.Client {
	orchestratorHost := lib.GetEnvVariable("ORCHESTRATOR_REDIS_HOST", "localhost")

	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", orchestratorHost, lib.GetEnvVariable("ORCHESTRATOR_REDIS_PORT", "10000")),
		Username: "default",
		Password: lib.GetEnvVariable("ORCHESTRATOR_REDIS_PASSWORD", ""),
	}

	isSecure := lib.GetEnvVariable("ORCHESTRATOR_REDIS_IS_SECURE", "false") == "true"

	if isSecure {
		clientCert := lib.GetEnvVariable("ORCHESTRATOR_REDIS_CERT", "")
		clientKey := lib.GetEnvVariable("ORCHESTRATOR_REDIS_KEY", "")

		cert, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
		if err != nil {
			panic(fmt.Errorf("error loading orchestrator cert: %s", err))
		}

		options.TLSConfig = &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: lib.GetEnvVariable("ORCHESTRATOR_REDIS_INSECURE_SKIP_VERIFY", "false") == "true",
			Certificates:       []tls.Certificate{cert},
		}
	}

	return redis.NewClient(options)
}

func getMaxJobs() int {
	maxJobs, err := strconv.Atoi(lib.GetEnvVariable("ORCHESTRATOR_MAX_JOBS", "1000"))
	if err != nil {
		maxJobs = 1000
	}

	return maxJobs
}

func getMaxManagedVUs() int64 {
	maxManagedVUs, err := strconv.ParseInt(lib.GetEnvVariable("ORCHESTRATOR_MAX_MANAGED_VUS", "10000"), 10, 64)
	if err != nil {
		maxManagedVUs = 10000
	}

	return maxManagedVUs
}
