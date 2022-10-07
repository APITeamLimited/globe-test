/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"syscall"

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

func fetchChildJob(ctx context.Context, orchestratorClient *redis.Client, childJobId string) (*libOrch.ChildJob, error) {
	childJobRaw, err := orchestratorClient.HGet(ctx, childJobId, "job").Result()
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
	client *redis.Client, job libOrch.ChildJob, workerId string) *libWorker.WorkerInfo {
	workerInfo := &libWorker.WorkerInfo{
		Client:         client,
		JobId:          job.Id,
		ChildJobId:     job.ChildJobId,
		ScopeId:        job.ScopeId,
		OrchestratorId: job.AssignedOrchestrator,
		WorkerId:       workerId,
		Ctx:            ctx,
	}

	if job.CollectionContext != nil {
		colletionVariables := make(map[string]string)

		for _, variable := range job.CollectionContext.Variables {
			colletionVariables[variable.Key] = variable.Value
		}

		workerInfo.Collection = &libWorker.Collection{
			Variables: colletionVariables,
		}
	}

	if job.EnvironmentContext != nil {
		environmentVariables := make(map[string]string)

		for _, variable := range job.EnvironmentContext.Variables {
			environmentVariables[variable.Key] = variable.Value
		}

		workerInfo.Environment = environmentVariables
	}

	return workerInfo
}
