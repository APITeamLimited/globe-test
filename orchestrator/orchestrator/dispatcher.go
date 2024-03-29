package orchestrator

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/gorilla/websocket"
)

type childJobDispatchResult struct {
	childJob *libOrch.ChildJob
	err      error
}

func dispatchChildJobs(gs libOrch.BaseGlobalState, childJobs *map[string]libOrch.ChildJobDistribution) error {
	unifiedDispatchResultCh := make(chan childJobDispatchResult)

	childJobCount := 0

	for location, jobDistribution := range *childJobs {
		for _, childJob := range jobDistribution.ChildJobs {
			dispatchResultCh := dispatchChildJob(gs, childJob, location)
			childJobCount++

			go func(dispatchCh chan childJobDispatchResult) {
				for v := range dispatchCh {
					unifiedDispatchResultCh <- v
					break
				}
			}(dispatchResultCh)
		}
	}

	dispatchedChildJobs := []*libOrch.ChildJob{}
	successFullDispatches := 0

	startTime := time.Now()
	for dispatchResult := range unifiedDispatchResultCh {
		if dispatchResult.err != nil {
			return dispatchResult.err
		}

		if dispatchResult.childJob != nil {
			dispatchedChildJobs = append(dispatchedChildJobs, dispatchResult.childJob)
		}

		successFullDispatches++

		if successFullDispatches == childJobCount {
			break
		}
	}

	fmt.Printf("Dispatched all child jobs in %s\n", time.Since(startTime))

	// Dispatch all child job info instantaneously after got all connections
	for _, childJob := range dispatchedChildJobs {
		go func(childJob *libOrch.ChildJob) error {
			serialializedChildJob, err := json.Marshal(childJob)
			if err != nil {
				fmt.Printf("Error marshalling child job %s: %s", childJob.ChildJobId, err)
			}

			childJobEvent := lib.EventMessage{
				Variant: lib.CHILD_JOB_INFO,
				Data:    string(serialializedChildJob),
			}

			marshalledEvent, err := json.Marshal(childJobEvent)
			if err != nil {
				fmt.Printf("Error marshalling child job event %s: %s", childJob.ChildJobId, err)
				return err
			}

			childJob.ConnWriteMutex.Lock()
			err = childJob.WorkerConnection.WriteMessage(websocket.TextMessage, marshalledEvent)
			childJob.ConnWriteMutex.Unlock()
			if err != nil {
				fmt.Printf("Error sending child job info to worker %s: %s", childJob.ChildJobId, err)
			}

			return nil
		}(childJob)
	}

	// Loop through childJobs and set WorkerConnection
	for location, jobDistribution := range *childJobs {
		for index, childJob := range jobDistribution.ChildJobs {
			addedConnection := false

			// Find the dispatched child job
			for _, dispatchedChildJob := range dispatchedChildJobs {
				if dispatchedChildJob.ChildJobId == childJob.ChildJobId {
					addedConnection = true
					(*childJobs)[location].ChildJobs[index].WorkerConnection = dispatchedChildJob.WorkerConnection
					(*childJobs)[location].ChildJobs[index].ConnWriteMutex = dispatchedChildJob.ConnWriteMutex
					(*childJobs)[location].ChildJobs[index].ConnReadMutex = dispatchedChildJob.ConnReadMutex
					break
				}
			}

			if !addedConnection {
				return errors.New("could not find dispatched child job")
			}
		}
	}

	return nil
}

func dispatchChildJob(gs libOrch.BaseGlobalState, childJob *libOrch.ChildJob, location string) chan childJobDispatchResult {
	dispatchResultCh := make(chan childJobDispatchResult)

	// Convert options to json
	go func() {
		// If we're in function mode, need to create a google cloud function call
		startTime := time.Now()
		conn, err := gs.FuncAuthClient().ExecuteService(gs, location)
		if err != nil {
			dispatchResultCh <- childJobDispatchResult{
				childJob: nil,
				err:      err,
			}

			return
		}

		newChildJob := *childJob
		newChildJob.WorkerConnection = conn
		newChildJob.ConnWriteMutex = &sync.Mutex{}
		newChildJob.ConnReadMutex = &sync.Mutex{}

		dispatchResultCh <- childJobDispatchResult{
			childJob: &newChildJob,
			err:      nil,
		}

		fmt.Println(("time taken to dispatch child job: " + time.Since(startTime).String()))

		fmt.Printf("Dispatched child job %s to function %s\n", childJob.ChildJobId, location)
	}()

	return dispatchResultCh
}
