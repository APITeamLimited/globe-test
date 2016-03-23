package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/loadimpact/speedboat/client"
	"github.com/loadimpact/speedboat/master"
	"time"
)

func init() ***REMOVED***
	registerCommand(cli.Command***REMOVED***
		Name:   "ping",
		Usage:  "Test command, will be removed",
		Action: actionPing,
	***REMOVED***)
	registerHandler(handlePing)
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
	out <- master.Message***REMOVED***
		Type: "ping.ping",
		Body: time.Now().Format("15:04:05 2006-01-02 MST"),
	***REMOVED***

	select ***REMOVED***
	case reply := <-in:
		log.WithFields(log.Fields***REMOVED***
			"type": reply.Type,
			"body": reply.Body,
		***REMOVED***).Info("Reply")
	case err := <-errors:
		log.WithError(err).Error("Ping failed")
	***REMOVED***
***REMOVED***
