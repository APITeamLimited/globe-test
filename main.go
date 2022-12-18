package main

import (
	"fmt"
	"log"
	"os"

	"flag"

	"net/http"
	_ "net/http/pprof"

	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/orchestrator/orchestrator"
	"github.com/APITeamLimited/globe-test/worker/worker"
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)

func main() {
	mode := flag.String("mode", "dev_worker_function", "worker, orchestrator, orchestrator_func_mode, or dev_worker_function")
	pProfPort := flag.Int("pprof-port", 0, "Enable pprof on the given port")

	flag.Parse()

	// If pprof is enabled, start the profiling server
	if *pProfPort != 0 {
		fmt.Printf("Starting pprof server on port %d\n", *pProfPort)
		go func() {
			log.Println(http.ListenAndServe(fmt.Sprintf("localhost:%d", *pProfPort), nil))
		}()
	}

	switch *mode {
	case "worker":
		worker.Run(true)
	case "orchestrator":
		orchestrator.Run(true, false)
	case "orchestrator_func_mode":
		orchestrator.Run(true, true)
	case "dev_worker_function":
		devWorkerServerPort := lib.GetEnvVariableRaw("DEV_WORKER_FUNCTION_PORT", "8090", true)
		fmt.Printf("Starting dev worker function on port %s\n", devWorkerServerPort)
		os.Setenv("FUNCTION_TARGET", "worker")

		funcframework.RegisterHTTPFunction("worker", worker.RunGoogleCloud)

		if err := funcframework.Start(devWorkerServerPort); err != nil {
			log.Fatalf("funcframework.Start: %v\n", err)
		}
	default:
		fmt.Println("Invalid GlobeTest mode, please specify one of: worker, orchestrator, orchestrator_func_mode, or dev_worker_function, default is dev_worker_function")
		os.Exit(1)
	}
}
