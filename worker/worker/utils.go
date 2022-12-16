package worker

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/lib/agent"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/redis/v9"
)

func fetchChildJob(ctx context.Context, client *redis.Client, childJobId string) (*libOrch.ChildJob, error) ***REMOVED***
	childJobRaw, err := client.HGet(ctx, childJobId, "job").Result()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Check child job not empty
	if childJobRaw == "" ***REMOVED***
		return nil, fmt.Errorf("child job %s is empty", childJobId)
	***REMOVED***

	childJob := libOrch.ChildJob***REMOVED******REMOVED***

	// Parse job as libOrch.ChildJob
	err = json.Unmarshal([]byte(childJobRaw), &childJob)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("error unmarshalling child job %s", childJobId)
	***REMOVED***

	return &childJob, nil
***REMOVED***

func startScheduling(ctx context.Context, client *redis.Client, workerId string,
	executionList *ExecutionList, creditsClient *redis.Client, standalone bool) ***REMOVED***
	jobsCheckScheduler := time.NewTicker(1 * time.Second)

	go func() ***REMOVED***
		for range jobsCheckScheduler.C ***REMOVED***
			client.HSet(ctx, fmt.Sprintf("worker:%s:info", workerId), "lastHeartbeat", time.Now().UnixMilli())
			client.HSet(ctx, fmt.Sprintf("worker:%s:info", workerId), "currentJobsCount", executionList.currentJobsCount)
			client.HSet(ctx, fmt.Sprintf("worker:%s:info", workerId), "currentVUsCount", executionList.currentVUsCount)

			client.HSet(ctx, fmt.Sprintf("worker:%s:info", workerId), "launchTime", time.Now().UnixMilli())
			client.HSet(ctx, fmt.Sprintf("worker:%s:info", workerId), "maxJobs", executionList.maxJobs)
			client.HSet(ctx, fmt.Sprintf("worker:%s:info", workerId), "maxVUs", executionList.maxVUs)

			client.SAdd(ctx, "workers", workerId)

			// Capacity may have freed up, check for queued jobs
			checkForQueuedJobs(ctx, client, workerId, executionList, creditsClient, standalone)
		***REMOVED***
	***REMOVED***()
***REMOVED***

func getWorkerClient(standalone bool) *redis.Client ***REMOVED***
	if !standalone ***REMOVED***
		return redis.NewClient(&redis.Options***REMOVED***
			Addr:     fmt.Sprintf("%s:%s", agent.WorkerRedisHost, agent.WorkerRedisPort),
			Username: "default",
			Password: "",
		***REMOVED***)
	***REMOVED***

	clientHost := lib.GetEnvVariable("CLIENT_HOST", "localhost")
	clientPort := lib.GetEnvVariable("CLIENT_PORT", "6978")

	options := &redis.Options***REMOVED***
		Addr:     fmt.Sprintf("%s:%s", clientHost, clientPort),
		Username: "default",
		Password: lib.GetEnvVariable("CLIENT_PASSWORD", ""),
	***REMOVED***

	isSecure := lib.GetEnvVariable("CLIENT_IS_SECURE", "false") == "true"

	if isSecure ***REMOVED***
		clientCert := lib.GetEnvVariable("CLIENT_CERT", "")
		clientKey := lib.GetEnvVariable("CLIENT_KEY", "")

		cert, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
		if err != nil ***REMOVED***
			panic(fmt.Errorf("error loading client cert: %s", err))
		***REMOVED***

		// Load CA cert
		caCertPool := x509.NewCertPool()
		caCert := lib.GetEnvVariable("CLIENT_CA_CERT", "")
		ok := caCertPool.AppendCertsFromPEM([]byte(caCert))
		if !ok ***REMOVED***
			panic("failed to parse root certificate")
		***REMOVED***

		options.TLSConfig = &tls.Config***REMOVED***
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: lib.GetEnvVariable("CLIENT_INSECURE_SKIP_VERIFY", "false") == "true",
			Certificates:       []tls.Certificate***REMOVED***cert***REMOVED***,
			RootCAs:            caCertPool,
		***REMOVED***
	***REMOVED***

	return redis.NewClient(options)
***REMOVED***

func getMaxJobs(standalone bool) int ***REMOVED***
	if !standalone ***REMOVED***
		// Orchestrator may spit jobs up, so set this high(ish)
		return 100
	***REMOVED***

	maxJobs, err := strconv.Atoi(lib.GetEnvVariable("WORKER_MAX_JOBS", "1000"))
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	return maxJobs
***REMOVED***

func getMaxVUs(standalone bool) int64 ***REMOVED***
	if !standalone ***REMOVED***
		return 5000
	***REMOVED***

	maxVUs, err := strconv.ParseInt(lib.GetEnvVariable("WORKER_MAX_VUS", "5000"), 10, 64)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	return maxVUs
***REMOVED***
