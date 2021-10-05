// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !go1.15
// +build !go1.15

package http2

import (
	"context"
	"crypto/tls"
)

// dialTLSWithContext opens a TLS connection.
func (t *Transport) dialTLSWithContext(ctx context.Context, network, addr string, cfg *tls.Config) (*tls.Conn, error) ***REMOVED***
	cn, err := tls.Dial(network, addr, cfg)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := cn.Handshake(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if cfg.InsecureSkipVerify ***REMOVED***
		return cn, nil
	***REMOVED***
	if err := cn.VerifyHostname(cfg.ServerName); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return cn, nil
***REMOVED***
