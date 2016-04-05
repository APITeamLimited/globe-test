package client

import (
	// log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/master"
	"github.com/loadimpact/speedboat/message"
)

// A Client controls load test execution.
type Client struct ***REMOVED***
	Connector master.Connector
***REMOVED***

// Creates a new Client, connecting to a Master listening on the given in/out addresses.
// The in/out addresses may be tcp:// or inproc:// addresses; see the documentation for
// mangos/nanomsg for more information.
func New(inAddr string, outAddr string) (c Client, err error) ***REMOVED***
	c.Connector, err = master.NewClientConnector(message.ClientTopic, inAddr, outAddr)
	if err != nil ***REMOVED***
		return c, err
	***REMOVED***

	return c, err
***REMOVED***

func (c *Client) Run() (<-chan message.Message, chan message.Message) ***REMOVED***
	return c.Connector.Run()
***REMOVED***
