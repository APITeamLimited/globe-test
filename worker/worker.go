package worker

import (
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/master"
	"github.com/loadimpact/speedboat/message"
)

// A Worker executes distributed tasks, communicating over a Master.
type Worker struct ***REMOVED***
	Connector  master.Connector
	Processors []func(*Worker) master.Processor

	pInstances []master.Processor
***REMOVED***

// Creates a new Worker, connecting to a master listening on the given in/out addresses.
func New(inAddr string, outAddr string) (w Worker, err error) ***REMOVED***
	w.Connector, err = master.NewClientConnector(message.WorkerTopic, inAddr, outAddr)
	if err != nil ***REMOVED***
		return w, err
	***REMOVED***

	return w, nil
***REMOVED***

// Runs the main loop for a worker.
func (w *Worker) Run() ***REMOVED***
	w.createProcessors()
	in, out, errors := w.Connector.Run()
	for ***REMOVED***
		select ***REMOVED***
		case msg := <-in:
			log.WithFields(log.Fields***REMOVED***
				"type":   msg.Type,
				"fields": msg.Fields,
			***REMOVED***).Info("Message Received")

			for m := range master.Process(w.pInstances, msg) ***REMOVED***
				out <- m
			***REMOVED***

		case err := <-errors:
			log.WithError(err).Error("Error")
		***REMOVED***
	***REMOVED***
***REMOVED***

func (w *Worker) createProcessors() ***REMOVED***
	w.pInstances = []master.Processor***REMOVED******REMOVED***
	for _, fn := range w.Processors ***REMOVED***
		w.pInstances = append(w.pInstances, fn(w))
	***REMOVED***
***REMOVED***
