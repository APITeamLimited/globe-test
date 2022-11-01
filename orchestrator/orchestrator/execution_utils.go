package orchestrator

import (
	"encoding/json"
	"fmt"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
)

func abortChildJobs(gs libOrch.BaseGlobalState, childJobs map[string]jobDistribution) error ***REMOVED***
	var err error

	cancelMessage := jobUserUpdate***REMOVED***
		UpdateType: "CANCEL",
	***REMOVED***

	marshalledCancelMessage, err := json.Marshal(cancelMessage)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	stringMarshalledCancelMessage := string(marshalledCancelMessage)

	for _, jobDistribution := range childJobs ***REMOVED***
		for _, job := range jobDistribution.Jobs ***REMOVED***
			err = jobDistribution.workerClient.Publish(gs.Ctx(), fmt.Sprintf("childJobUserUpdates:%s", job.ChildJobId), stringMarshalledCancelMessage).Err()
		***REMOVED***
	***REMOVED***

	return err
***REMOVED***

func abortAndFailAll(gs libOrch.BaseGlobalState, childJobs map[string]jobDistribution, err error) (string, error) ***REMOVED***
	abortErr := abortChildJobs(gs, childJobs)
	if abortErr != nil ***REMOVED***
		return "FAILURE", abortErr
	***REMOVED***

	libOrch.UpdateStatus(gs, "FAILURE")

	return "FAILURE", err
***REMOVED***
