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
func handleTestAbortSignals(gs *globalState, gracefulStopHandler, onHardStop func(os.Signal)) (stop func()) ***REMOVED***
	sigC := make(chan os.Signal, 2)
	done := make(chan struct***REMOVED******REMOVED***)
	gs.signalNotify(sigC, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() ***REMOVED***
		select ***REMOVED***
		case sig := <-sigC:
			gracefulStopHandler(sig)
		case <-done:
			return
		***REMOVED***

		select ***REMOVED***
		case sig := <-sigC:
			if onHardStop != nil ***REMOVED***
				onHardStop(sig)
			***REMOVED***
			// If we get a second signal, we immediately exit, so something like
			// https://github.com/k6io/k6/issues/971 never happens again
			gs.osExit(int(exitcodes.ExternalAbort))
		case <-done:
			return
		***REMOVED***
	***REMOVED***()

	return func() ***REMOVED***
		close(done)
		gs.signalStop(sigC)
	***REMOVED***
***REMOVED***

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

func loadWorkerInfo(ctx context.Context,
	client *redis.Client, job libOrch.ChildJob, workerId string) *libWorker.WorkerInfo ***REMOVED***
	workerInfo := &libWorker.WorkerInfo***REMOVED***
		Client:         client,
		JobId:          job.Id,
		ChildJobId:     job.ChildJobId,
		ScopeId:        job.ScopeId,
		OrchestratorId: job.AssignedOrchestrator,
		WorkerId:       workerId,
		Ctx:            ctx,
		WorkerOptions:  job.Options,
	***REMOVED***

	if job.CollectionContext != nil && job.CollectionContext.Name != "" ***REMOVED***
		collectionVariables := make(map[string]string)

		for _, variable := range job.CollectionContext.Variables ***REMOVED***
			collectionVariables[variable.Key] = variable.Value
		***REMOVED***

		workerInfo.Collection = &libWorker.Collection***REMOVED***
			Variables: collectionVariables,
			Name:      job.CollectionContext.Name,
		***REMOVED***
	***REMOVED***

	if job.EnvironmentContext != nil && job.EnvironmentContext.Name != "" ***REMOVED***
		environmentVariables := make(map[string]string)

		for _, variable := range job.EnvironmentContext.Variables ***REMOVED***
			environmentVariables[variable.Key] = variable.Value
		***REMOVED***

		workerInfo.Environment = &libWorker.Environment***REMOVED***
			Variables: environmentVariables,
			Name:      job.EnvironmentContext.Name,
		***REMOVED***
	***REMOVED***

	workerInfo.FinalRequest = job.FinalRequest
	workerInfo.UnderlyingRequest = job.UnderlyingRequest

	return workerInfo
***REMOVED***

func startJobScheduling(ctx context.Context, client *redis.Client, workerId string, executionList *ExecutionList) ***REMOVED***
	client.SAdd(ctx, "workers", workerId)

	jobsCheckScheduler := time.NewTicker(1 * time.Second)

	go func() ***REMOVED***
		for range jobsCheckScheduler.C ***REMOVED***
			client.HSet(ctx, fmt.Sprintf("worker:%s:info", workerId), "lastHeartbeat", time.Now().UnixMilli())
			client.HSet(ctx, fmt.Sprintf("worker:%s:info", workerId), "currentJobsCount", executionList.currentJobsCount)
			client.HSet(ctx, fmt.Sprintf("worker:%s:info", workerId), "currentVUsCount", executionList.currentVUsCount)

			// Capacity may have freed up, check for queued jobs
			checkForQueuedJobs(ctx, client, workerId, executionList)
		***REMOVED***
	***REMOVED***()

	client.HSet(ctx, fmt.Sprintf("worker:%s:info", workerId), "launchTime", time.Now().UnixMilli())
	client.HSet(ctx, fmt.Sprintf("worker:%s:info", workerId), "maxJobs", executionList.maxJobs)
	client.HSet(ctx, fmt.Sprintf("worker:%s:info", workerId), "maxVUs", executionList.maxVUs)
***REMOVED***

func getWorkerClient() *redis.Client ***REMOVED***
	return redis.NewClient(&redis.Options***REMOVED***
		Addr:     fmt.Sprintf("%s:%s", libWorker.GetEnvVariable("CLIENT_HOST", "localhost"), libWorker.GetEnvVariable("CLIENT_PORT", "6978")),
		Username: "default",
		Password: libWorker.GetEnvVariable("CLIENT_PASSWORD", ""),
	***REMOVED***)
***REMOVED***

func getMaxJobs() int ***REMOVED***
	maxJobs, err := strconv.Atoi(libOrch.GetEnvVariable("WORKER_MAX_JOBS", "1000"))
	if err != nil ***REMOVED***
		maxJobs = 1000
	***REMOVED***

	return maxJobs
***REMOVED***

func getMaxVUs() int64 ***REMOVED***
	maxVUs, err := strconv.ParseInt(libOrch.GetEnvVariable("WORKER_MAX_VUS", "5000"), 10, 64)
	if err != nil ***REMOVED***
		maxVUs = 10000
	***REMOVED***

	return maxVUs
***REMOVED***
