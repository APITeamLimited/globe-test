//go:build go1.17
// +build go1.17

package websocket

import (
	"context"
	"crypto/tls"
)

func doHandshake(ctx context.Context, tlsConn *tls.Conn, cfg *tls.Config) error ***REMOVED***
	if err := tlsConn.HandshakeContext(ctx); err != nil ***REMOVED***
		return err
	***REMOVED***
	if !cfg.InsecureSkipVerify ***REMOVED***
		if err := tlsConn.VerifyHostname(cfg.ServerName); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
