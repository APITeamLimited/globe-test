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

func setupChildProcesses() ***REMOVED***
	spawnChildServers()
	runOrchestrator()
	runWorker()
***REMOVED***

// Spawns child redis processes, these are terminated automatically when the agent exits
func spawnChildServers() ***REMOVED***
	orchestratorRedis := exec.Command(getServerCommandBase(), "--port", libAgent.OrchestratorPort)
	err := orchestratorRedis.Start()
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	clientRedis := exec.Command(getServerCommandBase(), "--port", libAgent.WorkerPort)
	err = clientRedis.Start()
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

func getServerCommandBase() string ***REMOVED***
	system := runtime.GOOS

	switch system ***REMOVED***
	case "darwin":
		return "./redis-server"
	case "linux":
		return "./redis-server"
	case "windows":
		return "./redis-server.exe"
	default:
		panic(fmt.Sprintf("unsupported system: %s", system))

	***REMOVED***
***REMOVED***

func runOrchestrator() ***REMOVED***
	// Set some environment variables
	_ = os.Setenv("ORCHESTRATOR_REDIS_HOST", libAgent.OrchestratorHost)
	_ = os.Setenv("ORCHESTRATOR_REDIS_PORT", libAgent.OrchestratorPort)
	_ = os.Setenv("ORCHESTRATOR_REDIS_PASSWORD", "")

	go orchestrator.Run(false)
***REMOVED***

func runWorker() ***REMOVED***
	// Set some environment variables
	_ = os.Setenv("CLIENT_HOST", libAgent.WorkerHost)
	_ = os.Setenv("CLIENT_PORT", libAgent.WorkerPort)
	_ = os.Setenv("CLIENT_PASSWORD", "")

	go worker.Run(false)
***REMOVED***
