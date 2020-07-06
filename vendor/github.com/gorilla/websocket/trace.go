// +build go1.8

package websocket

import (
	"crypto/tls"
	"net/http/httptrace"
)

func doHandshakeWithTrace(trace *httptrace.ClientTrace, tlsConn *tls.Conn, cfg *tls.Config) error ***REMOVED***
	if trace.TLSHandshakeStart != nil ***REMOVED***
		trace.TLSHandshakeStart()
	***REMOVED***
	err := doHandshake(tlsConn, cfg)
	if trace.TLSHandshakeDone != nil ***REMOVED***
		trace.TLSHandshakeDone(tlsConn.ConnectionState(), err)
	***REMOVED***
	return err
***REMOVED***
