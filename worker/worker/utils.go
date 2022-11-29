package worker

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/APITeamLimited/globe-test/agent/libAgent"
	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/errext/exitcodes"
	"github.com/APITeamLimited/redis/v9"
)

// Trap Interrupts, SIGINTs and SIGTERMs and call the given.
func handleTestAbortSignals(gs *globalState, gracefulStopHandler, onHardStop func(os.Signal)) (stop func()) {
	sigC := make(chan os.Signal, 2)
	done := make(chan struct{})
	gs.signalNotify(sigC, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case sig := <-sigC:
			gracefulStopHandler(sig)
		case <-done:
			return
		}

		select {
		case sig := <-sigC:
			if onHardStop != nil {
				onHardStop(sig)
			}
			// If we get a second signal, we immediately exit, so something like
			// https://github.com/k6io/k6/issues/971 never happens again
			gs.osExit(int(exitcodes.ExternalAbort))
		case <-done:
			return
		}
	}()

	return func() {
		close(done)
		gs.signalStop(sigC)
	}
}

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
			Addr:     fmt.Sprintf("%s:%s", libAgent.WorkerRedisHost, libAgent.WorkerRedisPort),
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

	isSecure := lib.GetEnvVariable("CLIENT_IS_SECURE", "false") == "true"

	if isSecure {
		clientCert := lib.GetEnvVariable("CLIENT_CERT", "")
		clientKey := lib.GetEnvVariable("CLIENT_KEY", "")

		cert, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
		if err != nil {
			panic(fmt.Errorf("error loading client cert: %s", err))
		}

		options.TLSConfig = &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: lib.GetEnvVariable("CLIENT_INSECURE_SKIP_VERIFY", "false") == "true",
			Certificates:       []tls.Certificate{cert},
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
