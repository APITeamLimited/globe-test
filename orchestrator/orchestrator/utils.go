package orchestrator

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/lib/agent"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/redis/v9"
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

func getOrchestratorClient(standalone bool) *redis.Client ***REMOVED***
	if !standalone ***REMOVED***
		return redis.NewClient(&redis.Options***REMOVED***
			Addr:     fmt.Sprintf("%s:%s", agent.OrchestratorRedisHost, agent.OrchestratorRedisPort),
			Username: "default",
			Password: "",
		***REMOVED***,
		)
	***REMOVED***

	orchestratorHost := lib.GetEnvVariable("ORCHESTRATOR_REDIS_HOST", "localhost")

	options := &redis.Options***REMOVED***
		Addr:     fmt.Sprintf("%s:%s", orchestratorHost, lib.GetEnvVariable("ORCHESTRATOR_REDIS_PORT", "10000")),
		Username: "default",
		Password: lib.GetEnvVariable("ORCHESTRATOR_REDIS_PASSWORD", ""),
	***REMOVED***

	isSecure := lib.GetEnvVariableRaw("ORCHESTRATOR_REDIS_IS_SECURE", "false", true) == "true"

	if isSecure ***REMOVED***
		clientCert := lib.GetEnvVariable("ORCHESTRATOR_REDIS_CERT", "")
		clientKey := lib.GetEnvVariable("ORCHESTRATOR_REDIS_KEY", "")

		cert, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
		if err != nil ***REMOVED***
			panic(fmt.Errorf("error loading orchestrator cert: %s", err))
		***REMOVED***

		options.TLSConfig = &tls.Config***REMOVED***
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: lib.GetEnvVariable("ORCHESTRATOR_REDIS_INSECURE_SKIP_VERIFY", "false") == "true",
			Certificates:       []tls.Certificate***REMOVED***cert***REMOVED***,
		***REMOVED***
	***REMOVED***

	return redis.NewClient(options)
***REMOVED***

func getMaxJobs(standalone bool) int ***REMOVED***
	if !standalone ***REMOVED***
		return 5
	***REMOVED***

	maxJobs, err := strconv.Atoi(lib.GetEnvVariable("ORCHESTRATOR_MAX_JOBS", "1000"))
	if err != nil ***REMOVED***
		maxJobs = 1000
	***REMOVED***

	return maxJobs
***REMOVED***

func getMaxManagedVUs(standalone bool) int64 ***REMOVED***
	if !standalone ***REMOVED***
		return 5000
	***REMOVED***

	maxManagedVUs, err := strconv.ParseInt(lib.GetEnvVariable("ORCHESTRATOR_MAX_MANAGED_VUS", "10000"), 10, 64)
	if err != nil ***REMOVED***
		maxManagedVUs = 10000
	***REMOVED***

	return maxManagedVUs
***REMOVED***
