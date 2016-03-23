package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/loadimpact/speedboat/client"
	"github.com/loadimpact/speedboat/master"
	"github.com/loadimpact/speedboat/worker"
	"time"
)

func init() ***REMOVED***
	registerCommand(cli.Command***REMOVED***
		Name:   "ping",
		Usage:  "Test command, will be removed",
		Action: actionPing,
		Flags: []cli.Flag***REMOVED***
			cli.BoolFlag***REMOVED***
				Name:  "worker",
				Usage: "Pings a worker instead of the master",
			***REMOVED***,
		***REMOVED***,
	***REMOVED***)
	registerHandler(handlePing)
	registerProcessor(processPing)
***REMOVED***

func processPing(w *worker.Worker, msg master.Message, out chan master.Message) bool ***REMOVED***
	switch msg.Type ***REMOVED***
	case "ping.worker.ping":
		out <- master.Message***REMOVED***
			Type: "ping.worker.pong",
			Body: msg.Body,
		***REMOVED***
		return true
	default:
		return false
	***REMOVED***
***REMOVED***

func handlePing(m *master.Master, msg master.Message, out chan master.Message) bool ***REMOVED***
	switch msg.Type ***REMOVED***
	case "ping.ping":
		out <- master.Message***REMOVED***
			Type: "ping.pong",
			Body: msg.Body,
		***REMOVED***
		return true
	default:
		return false
	***REMOVED***
***REMOVED***

func actionPing(c *cli.Context) ***REMOVED***
	client, err := client.New("tcp://127.0.0.1:9595", "tcp://127.0.0.1:9596")
	if err != nil ***REMOVED***
		log.WithError(err).Fatal("Couldn't create a client")
	***REMOVED***
	in, out, errors := client.Connector.Run()

	// Send a bunch of noise to filter through
	out <- master.Message***REMOVED***Type: "ping.noise"***REMOVED***
	out <- master.Message***REMOVED***Type: "ping.noise"***REMOVED***

	// Send a ping message, target should reply with a pong
	msgType := "ping.ping"
	if c.Bool("worker") ***REMOVED***
		msgType = "ping.worker.ping"
	***REMOVED***
	out <- master.Message***REMOVED***
		Type: msgType,
		Body: time.Now().Format("15:04:05 2006-01-02 MST"),
	***REMOVED***

readLoop:
	for ***REMOVED***
		select ***REMOVED***
		case msg := <-in:
			switch msg.Type ***REMOVED***
			case "ping.pong":
				log.WithFields(log.Fields***REMOVED***
					"body": msg.Body,
				***REMOVED***).Info("Pong!")
				break readLoop
			case "ping.worker.pong":
				log.WithFields(log.Fields***REMOVED***
					"body": msg.Body,
				***REMOVED***).Info("Worker Pong!")
				break readLoop
			***REMOVED***
		case err := <-errors:
			log.WithError(err).Error("Ping failed")
		***REMOVED***
	***REMOVED***
***REMOVED***
