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

func fetchChildJob(ctx context.Context, client *redis.Client, childJobId string) (*libOrch.ChildJob, error) {
	childJobRaw, err := client.HGet(ctx, childJobId, "job").Result()
	if err != nil {
		return nil, err
	}

	// Check child job not empty
	if childJobRaw == "" {
		return nil, fmt.Errorf("child job %s is empty", childJobId)
	}

	childJob := libOrch.ChildJob{}

	// Parse job as libOrch.ChildJob
	err = json.Unmarshal([]byte(childJobRaw), &childJob)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling child job %s", childJobId)
	}

	return &childJob, nil
}

func startScheduling(ctx context.Context, client *redis.Client, workerId string,
	executionList *ExecutionList, creditsClient *redis.Client, standalone bool) {
	jobsCheckScheduler := time.NewTicker(1 * time.Second)

	go func() {
		for range jobsCheckScheduler.C {
			client.HSet(ctx, fmt.Sprintf("worker:%s:info", workerId), "lastHeartbeat", time.Now().UnixMilli())
			client.HSet(ctx, fmt.Sprintf("worker:%s:info", workerId), "currentJobsCount", executionList.currentJobsCount)
			client.HSet(ctx, fmt.Sprintf("worker:%s:info", workerId), "currentVUsCount", executionList.currentVUsCount)

			client.HSet(ctx, fmt.Sprintf("worker:%s:info", workerId), "launchTime", time.Now().UnixMilli())
			client.HSet(ctx, fmt.Sprintf("worker:%s:info", workerId), "maxJobs", executionList.maxJobs)
			client.HSet(ctx, fmt.Sprintf("worker:%s:info", workerId), "maxVUs", executionList.maxVUs)

			client.SAdd(ctx, "workers", workerId)

			// Capacity may have freed up, check for queued jobs
			checkForQueuedJobs(ctx, client, workerId, executionList, creditsClient, standalone)
		}
	}()
}

func getWorkerClient(standalone bool) *redis.Client {
	if !standalone {
		return redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", agent.WorkerRedisHost, agent.WorkerRedisPort),
			Username: "default",
			Password: "",
		})
	}

	clientHost := lib.GetEnvVariable("CLIENT_HOST", "localhost")
	clientPort := lib.GetEnvVariable("CLIENT_PORT", "6978")

	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", clientHost, clientPort),
		Username: "default",
		Password: lib.GetEnvVariable("CLIENT_PASSWORD", ""),
	}

	isSecure := lib.GetEnvVariableBool("CLIENT_IS_SECURE", false)

	if isSecure {
		clientCert := lib.GetHexEnvVariable("CLIENT_CERT_HEX", "")
		clientKey := lib.GetHexEnvVariable("CLIENT_KEY_HEX", "")

		cert, err := tls.X509KeyPair(clientCert, clientKey)
		if err != nil {
			panic(fmt.Errorf("error loading client cert: %s", err))
		}

		// Load CA cert
		caCertPool := x509.NewCertPool()
		caCert := lib.GetHexEnvVariable("CLIENT_CA_CERT_HEX", "")
		ok := caCertPool.AppendCertsFromPEM(caCert)
		if !ok {
			panic("failed to parse root certificate")
		}

		options.TLSConfig = &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: lib.GetEnvVariableBool("CLIENT_INSECURE_SKIP_VERIFY", false),
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
		}
	}

	return redis.NewClient(options)
}

func getMaxJobs(standalone bool) int {
	if !standalone {
		// Orchestrator may spit jobs up, so set this high(ish)
		return 100
	}

	maxJobs, err := strconv.Atoi(lib.GetEnvVariable("WORKER_MAX_JOBS", "1000"))
	if err != nil {
		panic(err)
	}

	return maxJobs
}

func getMaxVUs(standalone bool) int64 {
	if !standalone {
		return 5000
	}

	maxVUs, err := strconv.ParseInt(lib.GetEnvVariable("WORKER_MAX_VUS", "5000"), 10, 64)
	if err != nil {
		panic(err)
	}

	return maxVUs
}
