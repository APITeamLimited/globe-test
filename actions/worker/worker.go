package actions

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/loadimpact/speedboat/client"
	"github.com/loadimpact/speedboat/util"
	"github.com/loadimpact/speedboat/worker"
)

func init() ***REMOVED***
	desc := "A worker executes distributed tasks, and reports back to its master."

	client.RegisterCommand(cli.Command***REMOVED***
		Name:        "worker",
		Usage:       "Runs a worker server for distributed tests",
		Description: desc,
		Flags: []cli.Flag***REMOVED***
			util.MasterHostFlag,
			util.MasterPortFlag,
		***REMOVED***,
		Action: actionWorker,
	***REMOVED***)
***REMOVED***

// Runs a standalone worker.
func actionWorker(c *cli.Context) ***REMOVED***
	inAddr, outAddr, local := util.ParseMasterParams(c)

	// Running a standalone worker without a master doesn't make any sense
	if local ***REMOVED***
		log.Fatal("No master specified")
	***REMOVED***

	w, err := worker.New(inAddr, outAddr)
	if err != nil ***REMOVED***
		log.WithError(err).Fatal("Failed to start worker")
	***REMOVED***

	log.WithFields(log.Fields***REMOVED***
		"master": inAddr,
	***REMOVED***).Info("Worker running")

	w.Processors = worker.GlobalProcessors
	w.Run()
***REMOVED***
