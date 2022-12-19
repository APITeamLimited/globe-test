package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/google/uuid"
)

func RunWorkerFunction(w http.ResponseWriter, r *http.Request) {
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

	fmt.Printf("Worker %s executing child job %s\n", workerId, childJob.ChildJobId)

	successfullExecution := handleExecution(ctx, client, childJob, workerId, creditsClient, true)
	if !successfullExecution {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error executing child job"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func RunDevWorkerServer() {
	devWorkerServerPort := lib.GetEnvVariableRaw("DEV_WORKER_FUNCTION_PORT", "8090", true)
	fmt.Printf("Starting dev worker function on port %s\n", devWorkerServerPort)
	os.Setenv("FUNCTION_TARGET", "WorkerCloud")

	// Load CLIENT_CERT from tls/client.crt
	// Load CLIENT_KEY from tls/client.key
	// Load CA_CERT from tls/ca.crt

	clientCertFile, err := os.ReadFile("tls/client.crt")
	if err != nil {
		panic(err)
	}

	os.Setenv("CLIENT_CERT", string(clientCertFile))

	clientKeyFile, err := os.ReadFile("tls/client.key")
	if err != nil {
		panic(err)
	}

	os.Setenv("CLIENT_KEY", string(clientKeyFile))

	caCertFile, err := os.ReadFile("tls/ca.crt")
	if err != nil {
		panic(err)
	}

	os.Setenv("CLIENT_CA_CERT", string(caCertFile))

	functions.HTTP("WorkerCloud", RunWorkerFunction)

	if err := funcframework.Start(devWorkerServerPort); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
