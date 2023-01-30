package main

import (
	"fmt"
	"log"
	"os"

	"flag"

	"net/http"
	_ "net/http/pprof"

	"github.com/APITeamLimited/globe-test/orchestrator/orchestrator"
	"github.com/APITeamLimited/globe-test/worker/worker"
)

func main() {
	mode := flag.String("mode", "worker_server", "worker_server or orchestrator")
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
	case "orchestrator":
		orchestrator.Run(true)
	case "worker_server":
		worker.RunWorkerServer()
	default:
		fmt.Println("Invalid GlobeTest mode, please specify one of: worker_server (default) or orchestrator")
		os.Exit(1)
	}
}
