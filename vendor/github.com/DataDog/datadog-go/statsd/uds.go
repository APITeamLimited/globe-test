package statsd

import (
	"net"
	"time"
)

/*
UDSTimeout holds the default timeout for UDS socket writes, as they can get
blocking when the receiving buffer is full.
*/
const defaultUDSTimeout = 1 * time.Millisecond

// udsWriter is an internal class wrapping around management of UDS connection
type udsWriter struct ***REMOVED***
	// Address to send metrics to, needed to allow reconnection on error
	addr net.Addr
	// Established connection object, or nil if not connected yet
	conn net.Conn
	// write timeout
	writeTimeout time.Duration
***REMOVED***

// New returns a pointer to a new udsWriter given a socket file path as addr.
func newUdsWriter(addr string) (*udsWriter, error) ***REMOVED***
	udsAddr, err := net.ResolveUnixAddr("unixgram", addr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	// Defer connection to first Write
	writer := &udsWriter***REMOVED***addr: udsAddr, conn: nil, writeTimeout: defaultUDSTimeout***REMOVED***
	return writer, nil
***REMOVED***

// SetWriteTimeout allows the user to set a custom write timeout
func (w *udsWriter) SetWriteTimeout(d time.Duration) error ***REMOVED***
	w.writeTimeout = d
	return nil
***REMOVED***

// Write data to the UDS connection with write timeout and minimal error handling:
// create the connection if nil, and destroy it if the statsd server has disconnected
func (w *udsWriter) Write(data []byte) (int, error) ***REMOVED***
	// Try connecting (first packet or connection lost)
	if w.conn == nil ***REMOVED***
		conn, err := net.Dial(w.addr.Network(), w.addr.String())
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		w.conn = conn
	***REMOVED***
	w.conn.SetWriteDeadline(time.Now().Add(w.writeTimeout))
	n, e := w.conn.Write(data)
	if e != nil ***REMOVED***
		// Statsd server disconnected, retry connecting at next packet
		w.conn = nil
		return 0, e
	***REMOVED***
	return n, e
***REMOVED***

func (w *udsWriter) Close() error ***REMOVED***
	if w.conn != nil ***REMOVED***
		return w.conn.Close()
	***REMOVED***
	return nil
***REMOVED***
