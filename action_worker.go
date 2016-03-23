package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/loadimpact/speedboat/common"
	"github.com/loadimpact/speedboat/worker"
)

func init() ***REMOVED***
	desc := "A worker executes distributed tasks, and reports back to its master."

	registerCommand(cli.Command***REMOVED***
		Name:        "worker",
		Usage:       "Runs a worker server for distributed tests",
		Description: desc,
		Flags: []cli.Flag***REMOVED***
			common.MasterHostFlag,
			common.MasterPortFlag,
		***REMOVED***,
		Action: actionWorker,
	***REMOVED***)
***REMOVED***

func actionWorker(c *cli.Context) ***REMOVED***
	inAddr, outAddr, local := common.ParseMasterParams(c)

	// Running a standalone worker without a master doesn't make any sense
	if local ***REMOVED***
		log.Fatal("No master specified")
	***REMOVED***

	worker, err := worker.New(inAddr, outAddr)
	if err != nil ***REMOVED***
		log.WithError(err).Fatal("Failed to start worker")
	***REMOVED***

	log.WithFields(log.Fields***REMOVED***
		"master": inAddr,
	***REMOVED***).Info("Worker running")

	worker.Processors = globalProcessors
	worker.Run()
***REMOVED***
