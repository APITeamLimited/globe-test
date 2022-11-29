package agent

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/APITeamLimited/globe-test/agent/libAgent"
	"github.com/APITeamLimited/globe-test/orchestrator/orchestrator"
	"github.com/APITeamLimited/globe-test/worker/worker"
)

func setupChildProcesses() {
	spawnChildServers()
	runOrchestrator()
	runWorker()
}

// Spawns child redis processes, these are terminated automatically when the agent exits
func spawnChildServers() {
	orchestratorRedis := exec.Command(getServerCommandBase(), "--port", libAgent.OrchestratorRedisPort, "--save", "", "--appendonly", "no")
	err := orchestratorRedis.Start()
	if err != nil {
		panic(err)
	}

	clientRedis := exec.Command(getServerCommandBase(), "--port", libAgent.WorkerRedisPort, "--save", "", "--appendonly", "no")
	err = clientRedis.Start()
	if err != nil {
		panic(err)
	}
}

func getServerCommandBase() string {
	system := runtime.GOOS

	switch system {
	case "darwin":
		return "./redis-server"
	case "linux":
		return "./redis-server"
	case "windows":
		return "./redis-server.exe"
	default:
		panic(fmt.Sprintf("unsupported system: %s", system))

	}
}

func runOrchestrator() {
	// Set some environment variables
	_ = os.Setenv("ORCHESTRATOR_REDIS_HOST", libAgent.OrchestratorRedisHost)
	_ = os.Setenv("ORCHESTRATOR_REDIS_PORT", libAgent.OrchestratorRedisPort)
	_ = os.Setenv("ORCHESTRATOR_REDIS_PASSWORD", "")

	go orchestrator.Run(false)
}

func runWorker() {
	// Set some environment variables
	_ = os.Setenv("CLIENT_HOST", libAgent.WorkerRedisHost)
	_ = os.Setenv("CLIENT_PORT", libAgent.WorkerRedisPort)
	_ = os.Setenv("CLIENT_PASSWORD", "")

	go worker.Run(false)
}
