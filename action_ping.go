package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/loadimpact/speedboat/client"
	"os"
)

func init() ***REMOVED***
	registerCommand(cli.Command***REMOVED***
		Name:   "ping",
		Usage:  "Test command, will be removed",
		Action: actionPing,
	***REMOVED***)
***REMOVED***

func actionPing(c *cli.Context) ***REMOVED***
	client, err := client.New("tcp://127.0.0.1:9595", "tcp://127.0.0.1:9596")
	if err != nil ***REMOVED***
		log.WithError(err).Fatal("Failed to ping")
	***REMOVED***

	go func() ***REMOVED***
		client.Connector.Send("ping")
		msg := <-client.Connector.InChannel
		log.WithField("msg", msg).Info("Response")
		os.Exit(0)
	***REMOVED***()

	ch, errors := client.Connector.Run()
	select ***REMOVED***
	case msg := <-ch:
		log.WithField("msg", msg).Info("Response")
	case err := <-errors:
		log.WithError(err).Error("Failed to ping master")
	***REMOVED***
***REMOVED***
