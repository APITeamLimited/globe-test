package master

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/loadimpact/speedboat/client"
	"github.com/loadimpact/speedboat/master"
)

func init() ***REMOVED***
	desc := "A master server acts as a message bus, between a clients and workers.\n" +
		"\n" +
		"The master works by opening TWO ports: a PUB port and a SUB port. Your firewall " +
		"must allow access to both of these, or clients will not be able to communicate " +
		"properly with the master."

	client.RegisterCommand(cli.Command***REMOVED***
		Name:        "master",
		Usage:       "Runs a master server for distributed tests",
		Description: desc,
		Flags: []cli.Flag***REMOVED***
			cli.StringFlag***REMOVED***
				Name:  "host, h",
				Usage: "Listen on the given address",
				Value: "127.0.0.1",
			***REMOVED***,
			cli.IntFlag***REMOVED***
				Name:  "port, p",
				Usage: "Listen on this port (PUB) + the next (SUB)",
				Value: 9595,
			***REMOVED***,
		***REMOVED***,
		Action: actionMaster,
	***REMOVED***)
***REMOVED***

// Runs a master.
func actionMaster(c *cli.Context) ***REMOVED***
	host := c.String("host")
	port := c.Int("port")

	outAddr := fmt.Sprintf("tcp://%s:%d", host, port)
	inAddr := fmt.Sprintf("tcp://%s:%d", host, port+1)
	m, err := master.New(outAddr, inAddr)
	if err != nil ***REMOVED***
		log.WithError(err).Fatal("Couldn't start master")
	***REMOVED***

	log.WithFields(log.Fields***REMOVED***
		"host": host,
		"pub":  port,
		"sub":  port + 1,
	***REMOVED***).Info("Master running")
	m.Processors = master.GlobalProcessors
	m.Run()
***REMOVED***
