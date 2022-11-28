package main

import (
	"fmt"
	"log"
	"os"

	"flag"

	"github.com/APITeamLimited/globe-test/agent"
	"github.com/APITeamLimited/globe-test/orchestrator/orchestrator"
	"github.com/APITeamLimited/globe-test/worker/worker"

	"net/http"
	_ "net/http/pprof"
)

func main() ***REMOVED***
	mode := flag.String("mode", "worker", "worker or orchestrator")
	pProfPort := flag.Int("pprof-port", 0, "Enable pprof on the given port")

	flag.Parse()

	// If pprof is enabled, start the profiling server
	if *pProfPort != 0 ***REMOVED***
		fmt.Printf("Starting pprof server on port %d\n", *pProfPort)
		go func() ***REMOVED***
			log.Println(http.ListenAndServe(fmt.Sprintf("localhost:%d", *pProfPort), nil))
		***REMOVED***()
	***REMOVED***

	switch *mode ***REMOVED***
	case "worker":
		worker.Run(true)
	case "orchestrator":
		orchestrator.Run(true)
	case "agent":
		agent.Run()
	default:
		fmt.Println("Invalid GlobeTest mode, please specify either worker or orchestrator (default is worker)")
		os.Exit(1)
	***REMOVED***
***REMOVED***
