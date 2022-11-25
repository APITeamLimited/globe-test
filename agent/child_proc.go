package agent

import (
	"context"
	"os"
	"os/exec"

	"github.com/APITeamLimited/globe-test/orchestrator/orchestrator"
	"github.com/APITeamLimited/globe-test/worker/worker"
	"github.com/APITeamLimited/redis/v9"
)

const orchestratorHost = "localhost"
const orchestratorPort = "59126"

const workerHost = "localhost"
const workerPort = "59127"

// Spawns child redis process
func spawnChildServers() {
	// Spawn a child redis server
	orchestratorRedis := exec.Command("redis-server", "--port", orchestratorPort)
	err := orchestratorRedis.Start()
	if err != nil {
		panic(err)
	}

	clientRedis := exec.Command("redis-server", "--port", workerPort)
	err = clientRedis.Start()
	if err != nil {
		panic(err)
	}
}

func closeChildServers(ctx context.Context, orchestratorClient *redis.Client, clientClient *redis.Client) {
	// Terminate the child redis server shutting it down gracefully
	orchestratorClient.ShutdownNoSave(ctx)
	clientClient.ShutdownNoSave(ctx)
}

func runWorker() {
	// Set some environment variables
	_ = os.Setenv("CLIENT_HOST", workerHost)
	_ = os.Setenv("CLIENT_PORT", workerPort)
	_ = os.Setenv("CLIENT_PASSWORD", "")

	go worker.Run()
}

func runOrchestrator() {
	// Set some environment variables
	_ = os.Setenv("ORCHESTRATOR_REDIS_HOST", orchestratorHost)
	_ = os.Setenv("ORCHESTRATOR_REDIS_PORT", orchestratorPort)
	_ = os.Setenv("ORCHESTRATOR_REDIS_PASSWORD", "")

	go orchestrator.Run()
}
