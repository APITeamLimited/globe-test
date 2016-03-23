package common

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/loadimpact/speedboat/client"
	"github.com/loadimpact/speedboat/master"
	"github.com/loadimpact/speedboat/worker"
)

var MasterHostFlag = cli.StringFlag***REMOVED***
	Name:  "master, m",
	Usage: "Host for the master process",
***REMOVED***
var MasterPortFlag = cli.IntFlag***REMOVED***
	Name:  "port, p",
	Usage: "Base port for the master process",
	Value: 9595,
***REMOVED***

func ParseMasterParams(c *cli.Context) (inAddr, outAddr string, local bool) ***REMOVED***
	switch ***REMOVED***
	case c.IsSet("master"):
		host := c.String("master")
		port := c.Int("port")
		inAddr = fmt.Sprintf("tcp://%s:%d", host, port)
		outAddr = fmt.Sprintf("tcp://%s:%d", host, port+1)
		local = false
	default:
		inAddr = "inproc://master.pub"
		outAddr = "inproc://master.sub"
		local = true
	***REMOVED***
	return inAddr, outAddr, local
***REMOVED***

func RunLocalMaster(inAddr, outAddr string) error ***REMOVED***
	m, err := master.New(inAddr, outAddr)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	go m.Run()
	return nil
***REMOVED***

func RunLocalWorker(inAddr, outAddr string) error ***REMOVED***
	w, err := worker.New(inAddr, outAddr)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	go w.Run()
	return nil
***REMOVED***

func MustGetClient(c *cli.Context) client.Client ***REMOVED***
	inAddr, outAddr, local := ParseMasterParams(c)

	// If we're running locally, ensure a local master and worker are running
	if local ***REMOVED***
		if err := RunLocalMaster(inAddr, outAddr); err != nil ***REMOVED***
			log.WithError(err).Fatal("Failed to start local master")
		***REMOVED***
		if err := RunLocalWorker(inAddr, outAddr); err != nil ***REMOVED***
			log.WithError(err).Fatal("Failed to start local worker")
		***REMOVED***
	***REMOVED***

	client, err := client.New(inAddr, outAddr)
	if err != nil ***REMOVED***
		log.WithError(err).Fatal("Failed to start a client")
	***REMOVED***
	return client
***REMOVED***
