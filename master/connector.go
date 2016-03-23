package master

import (
	log "github.com/Sirupsen/logrus"
	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/pub"
	"github.com/go-mangos/mangos/protocol/sub"
	"github.com/go-mangos/mangos/transport/inproc"
	"github.com/go-mangos/mangos/transport/tcp"
	"github.com/loadimpact/speedboat/message"
)

// A bidirectional pub/sub connector, used to connect to a master.
type Connector struct ***REMOVED***
	InSocket  mangos.Socket
	OutSocket mangos.Socket
***REMOVED***

// Creates a bare, unconnected connector.
func NewBareConnector() (conn Connector, err error) ***REMOVED***
	if conn.OutSocket, err = pub.NewSocket(); err != nil ***REMOVED***
		return conn, err
	***REMOVED***

	if conn.InSocket, err = sub.NewSocket(); err != nil ***REMOVED***
		return conn, err
	***REMOVED***

	return conn, nil
***REMOVED***

func NewClientConnector(topic string, inAddr string, outAddr string) (conn Connector, err error) ***REMOVED***
	if conn, err = NewBareConnector(); err != nil ***REMOVED***
		return conn, err
	***REMOVED***
	if err = setupAndDial(conn.InSocket, inAddr); err != nil ***REMOVED***
		return conn, err
	***REMOVED***
	if err = setupAndDial(conn.OutSocket, outAddr); err != nil ***REMOVED***
		return conn, err
	***REMOVED***

	err = conn.InSocket.SetOption(mangos.OptionSubscribe, []byte(topic))
	if err != nil ***REMOVED***
		return conn, err
	***REMOVED***

	return conn, nil
***REMOVED***

func NewServerConnector(outAddr string, inAddr string) (conn Connector, err error) ***REMOVED***
	if conn, err = NewBareConnector(); err != nil ***REMOVED***
		return conn, err
	***REMOVED***
	if err = setupAndListen(conn.OutSocket, outAddr); err != nil ***REMOVED***
		return conn, err
	***REMOVED***
	if err = setupAndListen(conn.InSocket, inAddr); err != nil ***REMOVED***
		return conn, err
	***REMOVED***

	err = conn.InSocket.SetOption(mangos.OptionSubscribe, []byte(""))
	if err != nil ***REMOVED***
		return conn, err
	***REMOVED***

	return conn, nil
***REMOVED***

func setupSocket(sock mangos.Socket) ***REMOVED***
	sock.AddTransport(inproc.NewTransport())
	sock.AddTransport(tcp.NewTransport())
***REMOVED***

func setupAndListen(sock mangos.Socket, addr string) error ***REMOVED***
	setupSocket(sock)
	if err := sock.Listen(addr); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func setupAndDial(sock mangos.Socket, addr string) error ***REMOVED***
	setupSocket(sock)
	if err := sock.Dial(addr); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// Provides a channel-based interface around the underlying socket API.
func (c *Connector) Run() (<-chan message.Message, chan message.Message, <-chan error) ***REMOVED***
	errors := make(chan error)
	in := make(chan message.Message)
	out := make(chan message.Message)

	// Read incoming messages
	go func() ***REMOVED***
		for ***REMOVED***
			msg, err := c.Read()
			if err != nil ***REMOVED***
				errors <- err
				continue
			***REMOVED***
			in <- msg
		***REMOVED***
	***REMOVED***()

	// Write outgoing messages
	go func() ***REMOVED***
		for ***REMOVED***
			msg := <-out
			err := c.Write(msg)
			if err != nil ***REMOVED***
				errors <- err
				continue
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return in, out, errors
***REMOVED***

func (c *Connector) Read() (msg message.Message, err error) ***REMOVED***
	data, err := c.InSocket.Recv()
	if err != nil ***REMOVED***
		return msg, err
	***REMOVED***
	log.WithField("data", string(data)).Debug("Read data")
	err = message.Decode(data, &msg)
	if err != nil ***REMOVED***
		return msg, err
	***REMOVED***
	log.WithFields(log.Fields***REMOVED***
		"type": msg.Type,
		"body": msg.Body,
	***REMOVED***).Debug("Decoded message")
	return msg, nil
***REMOVED***

func (c *Connector) Write(msg message.Message) (err error) ***REMOVED***
	data, err := msg.Encode()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	log.WithField("data", string(data)).Debug("Writing data")
	err = c.OutSocket.Send(data)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***
