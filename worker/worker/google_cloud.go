package worker

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/google/uuid"
)

func RunGoogleCloud(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	workerId := uuid.NewString()
	client := getWorkerClient(true)

	creditsClient := lib.GetCreditsClient(true)

	// Ensure is POST request
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}

	// Get the childJob from the request body
	decoder := json.NewDecoder(r.Body)
	var childJob libOrch.ChildJob

	err := decoder.Decode(&childJob)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	successfullExecution := handleExecution(ctx, client, childJob, workerId, creditsClient, true)
	if !successfullExecution {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error executing child job"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
