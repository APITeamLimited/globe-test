// +build !go1.8

package websocket

import (
	"crypto/tls"
	"net/http/httptrace"
)

func doHandshakeWithTrace(trace *httptrace.ClientTrace, tlsConn *tls.Conn, cfg *tls.Config) error ***REMOVED***
	return doHandshake(tlsConn, cfg)
***REMOVED***
