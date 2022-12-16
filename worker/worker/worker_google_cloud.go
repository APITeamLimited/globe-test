package worker

import (
	"context"
	"net/http"

	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/google/uuid"
)

func init() ***REMOVED***
	functions.HTTP("worker", RunGoogleCloud)
***REMOVED***

func RunGoogleCloud(w http.ResponseWriter, r *http.Request) ***REMOVED***
	ctx := context.Background()
	workerId := uuid.NewString()
	client := getWorkerClient(true)
	maxJobs := 1
	maxVUs := int64(-1) // unlimited

	creditsClient := lib.GetCreditsClient(true)

	executionList := &ExecutionList***REMOVED***
		currentJobs: make(map[string]libOrch.ChildJob),
		maxJobs:     maxJobs,
		maxVUs:      maxVUs,
	***REMOVED***

	// Get the childJobId from the request

	childJobId := r.URL.Query().Get("childJobId")

	if childJobId == "" ***REMOVED***
		panic("childJobId is required")
	***REMOVED***

	err := checkIfCanExecute(ctx, client, childJobId, workerId, executionList, creditsClient, true)
	if err != nil ***REMOVED***
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	***REMOVED***

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
***REMOVED***
