package comm

import (
	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/pub"
	"github.com/go-mangos/mangos/protocol/sub"
	"github.com/go-mangos/mangos/transport/inproc"
	"github.com/go-mangos/mangos/transport/tcp"
)

// A bidirectional pub/sub connector, used for master-based communication.
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

// Creates a connection to the specified master address. It will subscribe to a certain topic and
// filter out anything unrelated to it, and handles reconnections automatically under the hood.
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

// Creates a listening connector for a master server. It will accept any incoming messages.
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

// Common setup for a Mangos socket.
func setupSocket(sock mangos.Socket) ***REMOVED***
	sock.AddTransport(inproc.NewTransport())
	sock.AddTransport(tcp.NewTransport())
***REMOVED***

// Performs standard setup and listens on the specified address.
func setupAndListen(sock mangos.Socket, addr string) error ***REMOVED***
	setupSocket(sock)
	if err := sock.Listen(addr); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// Performs standard setup and dials the specified address.
func setupAndDial(sock mangos.Socket, addr string) error ***REMOVED***
	setupSocket(sock)
	if err := sock.Dial(addr); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// Provides a channel-based interface around the underlying socket API.
func (c *Connector) Run() (<-chan Message, chan Message) ***REMOVED***
	in := make(chan Message)
	out := make(chan Message)

	// Read incoming messages
	go func() ***REMOVED***
		for ***REMOVED***
			msg, err := c.Read()
			if err != nil ***REMOVED***
				in <- ToClient("error").WithError(err)
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
				in <- ToClient("error").WithError(err)
				continue
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return in, out
***REMOVED***

// Reads a single message from a connector; CANNOT be used together with Run().
// Run() will call this by itself under the hood, and if you call it outside as well, you'll create
// a race condition where only one of Run() and Read() will receive a comm.
func (c *Connector) Read() (msg Message, err error) ***REMOVED***
	data, err := c.InSocket.Recv()
	if err != nil ***REMOVED***
		return msg, err
	***REMOVED***
	err = Decode(data, &msg)
	if err != nil ***REMOVED***
		return msg, err
	***REMOVED***
	return msg, nil
***REMOVED***

// Writes a single message to a connector.
func (c *Connector) Write(msg Message) (err error) ***REMOVED***
	data, err := msg.Encode()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	err = c.OutSocket.Send(data)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***
