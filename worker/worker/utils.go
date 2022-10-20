package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/errext/exitcodes"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
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

func loadWorkerInfo(ctx context.Context,
	client *redis.Client, job libOrch.ChildJob, workerId string, gs libWorker.BaseGlobalState) *libWorker.WorkerInfo {
	workerInfo := &libWorker.WorkerInfo{
		Client:         client,
		JobId:          job.Id,
		ChildJobId:     job.ChildJobId,
		ScopeId:        job.ScopeId,
		OrchestratorId: job.AssignedOrchestrator,
		WorkerId:       workerId,
		Ctx:            ctx,
		WorkerOptions:  job.Options,
		Gs:             &gs,
	}

	if job.CollectionContext != nil && job.CollectionContext.Name != "" {
		collectionVariables := make(map[string]string)

		for _, variable := range job.CollectionContext.Variables {
			collectionVariables[variable.Key] = variable.Value
		}

		workerInfo.Collection = &libWorker.Collection{
			Variables: collectionVariables,
			Name:      job.CollectionContext.Name,
		}
	}

	if job.EnvironmentContext != nil && job.EnvironmentContext.Name != "" {
		environmentVariables := make(map[string]string)

		for _, variable := range job.EnvironmentContext.Variables {
			environmentVariables[variable.Key] = variable.Value
		}

		workerInfo.Environment = &libWorker.Environment{
			Variables: environmentVariables,
			Name:      job.EnvironmentContext.Name,
		}
	}

	workerInfo.FinalRequest = job.FinalRequest
	workerInfo.UnderlyingRequest = job.UnderlyingRequest

	return workerInfo
}

func startJobScheduling(ctx context.Context, client *redis.Client, workerId string, executionList *ExecutionList) {
	client.SAdd(ctx, "workers", workerId)

	jobsCheckScheduler := time.NewTicker(1 * time.Second)

	go func() {
		for range jobsCheckScheduler.C {
			client.HSet(ctx, fmt.Sprintf("worker:%s:info", workerId), "lastHeartbeat", time.Now().UnixMilli())
			client.HSet(ctx, fmt.Sprintf("worker:%s:info", workerId), "currentJobsCount", executionList.currentJobsCount)
			client.HSet(ctx, fmt.Sprintf("worker:%s:info", workerId), "currentVUsCount", executionList.currentVUsCount)

			// Capacity may have freed up, check for queued jobs
			checkForQueuedJobs(ctx, client, workerId, executionList)
		}
	}()

	client.HSet(ctx, fmt.Sprintf("worker:%s:info", workerId), "launchTime", time.Now().UnixMilli())
	client.HSet(ctx, fmt.Sprintf("worker:%s:info", workerId), "maxJobs", executionList.maxJobs)
	client.HSet(ctx, fmt.Sprintf("worker:%s:info", workerId), "maxVUs", executionList.maxVUs)
}

func getWorkerClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", libWorker.GetEnvVariable("CLIENT_HOST", "localhost"), libWorker.GetEnvVariable("CLIENT_PORT", "6978")),
		Username: "default",
		Password: libWorker.GetEnvVariable("CLIENT_PASSWORD", ""),
	})
}

func getMaxJobs() int {
	maxJobs, err := strconv.Atoi(libOrch.GetEnvVariable("WORKER_MAX_JOBS", "1000"))
	if err != nil {
		maxJobs = 1000
	}

	return maxJobs
}

func getMaxVUs() int64 {
	maxVUs, err := strconv.ParseInt(libOrch.GetEnvVariable("WORKER_MAX_VUS", "5000"), 10, 64)
	if err != nil {
		maxVUs = 10000
	}

	return maxVUs
}
