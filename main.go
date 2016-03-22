package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/client"
	"github.com/loadimpact/speedboat/master"
)

func main() ***REMOVED***
	master, err := master.New("inproc://master.pub", "inproc://master.sub")
	if err != nil ***REMOVED***
		log.WithFields(log.Fields***REMOVED***
			"error": err,
		***REMOVED***).Fatal("Failed to start master")
	***REMOVED***
	go master.Run()

	client, err := client.New("inproc://master.pub", "inproc://master.sub")
	if err != nil ***REMOVED***
		log.WithFields(log.Fields***REMOVED***
			"error": err,
		***REMOVED***).Fatal("Failed to start client")
	***REMOVED***
	client.Run()

	log.WithFields(log.Fields***REMOVED***
		"thing": "aaaa",
	***REMOVED***).Info("Is this working??")
***REMOVED***
