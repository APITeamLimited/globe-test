package main

import (
	"fmt"
	"os"

	"flag"

	"github.com/APITeamLimited/globe-test/orchestrator/orchestrator"
	"github.com/APITeamLimited/globe-test/worker/worker"
)

func main() {
	mode := flag.String("mode", "worker", "worker or orchestrator")

	flag.Parse()

	switch *mode {
	case "worker":
		worker.Run()
	case "orchestrator":
		orchestrator.Run()
	default:
		fmt.Println("Invalid GlobeTest mode, please specify either worker or orchestrator")
		os.Exit(1)
	}
}
