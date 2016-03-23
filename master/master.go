package master

import (
	log "github.com/Sirupsen/logrus"
)

// A Master serves as a semi-intelligent message bus between clients and workers.
type Master struct ***REMOVED***
	Connector Connector
	Handlers  []func(*Master, Message, chan Message) bool
***REMOVED***

// Creates a new Master, listening on the given in/out addresses.
// The in/out addresses may be tcp:// or inproc:// addresses.
// Note that positions of the in/out parameters are swapped compared to client.New(), to make
// `client.New(a, b)` connect to a master created with `master.New(a, b)`.
func New(outAddr string, inAddr string) (m Master, err error) ***REMOVED***
	m.Connector, err = NewServerConnector(outAddr, inAddr)
	if err != nil ***REMOVED***
		return m, err
	***REMOVED***

	return m, nil
***REMOVED***

// Runs the main loop for a master.
func (m *Master) Run() ***REMOVED***
	in, out, errors := m.Connector.Run()
	for ***REMOVED***
		select ***REMOVED***
		case msg := <-in:
			log.WithFields(log.Fields***REMOVED***
				"type": msg.Type,
				"body": msg.Body,
			***REMOVED***).Info("Message Received")

			// Call handlers until we find one that responds
			for _, handler := range m.Handlers ***REMOVED***
				if handler(m, msg, out) ***REMOVED***
					break
				***REMOVED***
			***REMOVED***

		case err := <-errors:
			log.WithError(err).Error("Error")
		***REMOVED***
	***REMOVED***
***REMOVED***
