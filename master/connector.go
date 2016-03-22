package master

import (
	log "github.com/Sirupsen/logrus"
	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/pub"
	"github.com/go-mangos/mangos/protocol/sub"
	"github.com/go-mangos/mangos/transport/inproc"
	"github.com/go-mangos/mangos/transport/tcp"
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

func NewClientConnector(inAddr string, outAddr string) (conn Connector, err error) ***REMOVED***
	if conn, err = NewBareConnector(); err != nil ***REMOVED***
		return conn, err
	***REMOVED***
	if err = setupAndDial(conn.InSocket, inAddr); err != nil ***REMOVED***
		return conn, err
	***REMOVED***
	if err = setupAndDial(conn.OutSocket, outAddr); err != nil ***REMOVED***
		return conn, err
	***REMOVED***

	err = conn.InSocket.SetOption(mangos.OptionSubscribe, []byte(""))
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

func (c *Connector) Run() (chan string, <-chan error) ***REMOVED***
	ch := make(chan string)
	errors := make(chan error)

	// Start a read loop
	go func() ***REMOVED***
		log.Debug("-> Connector Read Loop")
		msg, err := c.InSocket.Recv()
		if err != nil ***REMOVED***
			errors <- err
		***REMOVED***
		ch <- string(msg)
		log.Debug("<- Connector Read Loop")
	***REMOVED***()

	// // Start a write loop
	go func() ***REMOVED***
		log.Debug("-> Connector Write Loop")
		msg := <-ch
		if err := c.OutSocket.Send([]byte(msg)); err != nil ***REMOVED***
			errors <- err
		***REMOVED***
		log.Debug("<- Connector Write Loop")
	***REMOVED***()

	return ch, errors
***REMOVED***
