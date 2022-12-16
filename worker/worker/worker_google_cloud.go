package worker

import (
	"context"
	"net/http"

	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/google/uuid"
)

func init() {
	functions.HTTP("worker", RunGoogleCloud)
}

func RunGoogleCloud(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	workerId := uuid.NewString()
	client := getWorkerClient(true)
	maxJobs := 1
	maxVUs := int64(-1) // unlimited

	creditsClient := lib.GetCreditsClient(true)

	executionList := &ExecutionList{
		currentJobs: make(map[string]libOrch.ChildJob),
		maxJobs:     maxJobs,
		maxVUs:      maxVUs,
	}

	// Get the childJobId from the request

	childJobId := r.URL.Query().Get("childJobId")

	if childJobId == "" {
		panic("childJobId is required")
	}

	err := checkIfCanExecute(ctx, client, childJobId, workerId, executionList, creditsClient, true)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
