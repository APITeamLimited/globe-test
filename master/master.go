package master

import (
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/comm"
)

// A Master serves as a semi-intelligent message bus between clients and workers.
type Master struct ***REMOVED***
	Connector  comm.Connector
	Processors []func(*Master) comm.Processor
***REMOVED***

// Creates a new Master, listening on the given in/out addresses.
// The in/out addresses may be tcp:// or inproc:// addresses.
// Note that positions of the in/out parameters are swapped compared to client.New(), to make
// `client.New(a, b)` connect to a master created with `comm.New(a, b)`.
func New(outAddr string, inAddr string) (m Master, err error) ***REMOVED***
	m.Connector, err = comm.NewServerConnector(outAddr, inAddr)
	if err != nil ***REMOVED***
		return m, err
	***REMOVED***

	return m, nil
***REMOVED***

// Runs the main loop for a master.
func (m *Master) Run() ***REMOVED***
	in, out := m.Connector.Run()
	pInstances := m.createProcessors()
	for msg := range in ***REMOVED***
		log.WithFields(log.Fields***REMOVED***
			"type":    msg.Type,
			"topic":   msg.Topic,
			"payload": string(msg.Payload),
		***REMOVED***).Debug("Master Received")

		// If it's not intended for the master, rebroadcast
		if msg.Topic != comm.MasterTopic ***REMOVED***
			out <- msg
			continue
		***REMOVED***

		// Let master processors have a stab at them instead
		go func() ***REMOVED***
			for m := range comm.Process(pInstances, msg) ***REMOVED***
				out <- m
			***REMOVED***
		***REMOVED***()
	***REMOVED***
***REMOVED***

func (m *Master) createProcessors() []comm.Processor ***REMOVED***
	pInstances := []comm.Processor***REMOVED******REMOVED***
	for _, fn := range m.Processors ***REMOVED***
		pInstances = append(pInstances, fn(m))
	***REMOVED***
	return pInstances
***REMOVED***
