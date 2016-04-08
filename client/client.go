package client

import (
	// log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/comm"
)

// A Client controls load test execution.
type Client struct ***REMOVED***
	Connector comm.Connector
***REMOVED***

// Creates a new Client, connecting to a Master listening on the given in/out addresses.
// The in/out addresses may be tcp:// or inproc:// addresses; see the documentation for
// mangos/nanomsg for more information.
func New(inAddr string, outAddr string) (c Client, err error) ***REMOVED***
	c.Connector, err = comm.NewClientConnector(comm.ClientTopic, inAddr, outAddr)
	if err != nil ***REMOVED***
		return c, err
	***REMOVED***

	return c, err
***REMOVED***

func (c *Client) Run() (<-chan comm.Message, chan comm.Message) ***REMOVED***
	return c.Connector.Run()
***REMOVED***
