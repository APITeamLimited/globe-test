package statsd

import (
	"errors"
	"net"
	"time"
)

// udpWriter is an internal class wrapping around management of UDP connection
type udpWriter struct ***REMOVED***
	conn net.Conn
***REMOVED***

// New returns a pointer to a new udpWriter given an addr in the format "hostname:port".
func newUDPWriter(addr string) (*udpWriter, error) ***REMOVED***
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	writer := &udpWriter***REMOVED***conn: conn***REMOVED***
	return writer, nil
***REMOVED***

// SetWriteTimeout is not needed for UDP, returns error
func (w *udpWriter) SetWriteTimeout(d time.Duration) error ***REMOVED***
	return errors.New("SetWriteTimeout: not supported for UDP connections")
***REMOVED***

// Write data to the UDP connection with no error handling
func (w *udpWriter) Write(data []byte) (int, error) ***REMOVED***
	return w.conn.Write(data)
***REMOVED***

func (w *udpWriter) Close() error ***REMOVED***
	return w.conn.Close()
***REMOVED***
