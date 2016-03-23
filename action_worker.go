package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/loadimpact/speedboat/worker"
)

func init() ***REMOVED***
	desc := "A worker executes distributed tasks, and reports back to its master."

	registerCommand(cli.Command***REMOVED***
		Name:        "worker",
		Usage:       "Runs a worker server for distributed tests",
		Description: desc,
		Flags: []cli.Flag***REMOVED***
			cli.StringFlag***REMOVED***
				Name:  "host, h",
				Usage: "Host for the master process",
				Value: "127.0.0.1",
			***REMOVED***,
			cli.IntFlag***REMOVED***
				Name:  "port, p",
				Usage: "Base port for the master process",
				Value: 9595,
			***REMOVED***,
		***REMOVED***,
		Action: actionWorker,
	***REMOVED***)
***REMOVED***

func actionWorker(c *cli.Context) ***REMOVED***
	host := c.String("host")
	port := c.Int("port")

	inAddr := fmt.Sprintf("tcp://%s:%d", host, port)
	outAddr := fmt.Sprintf("tcp://%s:%d", host, port+1)
	worker, err := worker.New(inAddr, outAddr)
	if err != nil ***REMOVED***
		log.WithError(err).Fatal("Couldn't start worker")
	***REMOVED***

	log.WithFields(log.Fields***REMOVED***
		"host": host,
		"pub":  port,
		"sub":  port + 1,
	***REMOVED***).Info("Worker running")
	worker.Processors = globalProcessors
	worker.Run()
***REMOVED***
