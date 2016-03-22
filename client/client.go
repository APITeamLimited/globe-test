package client

import (
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/master"
)

// A Client controls load test execution.
type Client struct ***REMOVED***
	Connector master.Connector
***REMOVED***

// Creates a new Client, connecting to a Master listening on the given in/out addresses.
// The in/out addresses may be tcp:// or inproc:// addresses; see the documentation for
// mangos/nanomsg for more information.
func New(inAddr string, outAddr string) (c Client, err error) ***REMOVED***
	c.Connector, err = master.NewClientConnector(inAddr, outAddr)
	if err != nil ***REMOVED***
		return c, err
	***REMOVED***

	return c, err
***REMOVED***

// Runs the main loop for a client. This is probably going to go away.
func (c *Client) Run() ***REMOVED***
	ch, errors := c.Connector.Run()
	for ***REMOVED***
		select ***REMOVED***
		case msg := <-ch:
			log.WithFields(log.Fields***REMOVED***
				"msg": msg,
			***REMOVED***).Info("Client: Message received")
		case err := <-errors:
			log.WithFields(log.Fields***REMOVED***
				"error": err,
			***REMOVED***).Error("Client: Error receiving")
		***REMOVED***
	***REMOVED***
***REMOVED***
