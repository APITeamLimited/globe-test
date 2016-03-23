package worker

import (
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/master"
)

// A Worker executes distributed tasks, communicating over a Master.
type Worker struct ***REMOVED***
	Connector  master.Connector
	Processors []func(*Worker, master.Message, chan master.Message) bool
***REMOVED***

// Creates a new Worker, connecting to a master listening on the given in/out addresses.
func New(inAddr string, outAddr string) (w Worker, err error) ***REMOVED***
	w.Connector, err = master.NewClientConnector(inAddr, outAddr)
	if err != nil ***REMOVED***
		return w, err
	***REMOVED***

	return w, nil
***REMOVED***

// Runs the main loop for a worker.
func (w *Worker) Run() ***REMOVED***
	in, out, errors := w.Connector.Run()
	for ***REMOVED***
		select ***REMOVED***
		case msg := <-in:
			log.WithFields(log.Fields***REMOVED***
				"type": msg.Type,
				"body": msg.Body,
			***REMOVED***).Info("Message Received")

			// Call handlers until we find one that responds
			for _, processor := range w.Processors ***REMOVED***
				if processor(w, msg, out) ***REMOVED***
					break
				***REMOVED***
			***REMOVED***

		case err := <-errors:
			log.WithError(err).Error("Error")
		***REMOVED***
	***REMOVED***
***REMOVED***
