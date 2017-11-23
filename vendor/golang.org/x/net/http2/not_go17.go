// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !go1.7

package http2

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

type contextContext interface ***REMOVED***
	Done() <-chan struct***REMOVED******REMOVED***
	Err() error
***REMOVED***

type fakeContext struct***REMOVED******REMOVED***

func (fakeContext) Done() <-chan struct***REMOVED******REMOVED*** ***REMOVED*** return nil ***REMOVED***
func (fakeContext) Err() error            ***REMOVED*** panic("should not be called") ***REMOVED***

func reqContext(r *http.Request) fakeContext ***REMOVED***
	return fakeContext***REMOVED******REMOVED***
***REMOVED***

func setResponseUncompressed(res *http.Response) ***REMOVED***
	// Nothing.
***REMOVED***

type clientTrace struct***REMOVED******REMOVED***

func requestTrace(*http.Request) *clientTrace ***REMOVED*** return nil ***REMOVED***
func traceGotConn(*http.Request, *ClientConn) ***REMOVED******REMOVED***
func traceFirstResponseByte(*clientTrace)     ***REMOVED******REMOVED***
func traceWroteHeaders(*clientTrace)          ***REMOVED******REMOVED***
func traceWroteRequest(*clientTrace, error)   ***REMOVED******REMOVED***
func traceGot100Continue(trace *clientTrace)  ***REMOVED******REMOVED***
func traceWait100Continue(trace *clientTrace) ***REMOVED******REMOVED***

func nop() ***REMOVED******REMOVED***

func serverConnBaseContext(c net.Conn, opts *ServeConnOpts) (ctx contextContext, cancel func()) ***REMOVED***
	return nil, nop
***REMOVED***

func contextWithCancel(ctx contextContext) (_ contextContext, cancel func()) ***REMOVED***
	return ctx, nop
***REMOVED***

func requestWithContext(req *http.Request, ctx contextContext) *http.Request ***REMOVED***
	return req
***REMOVED***

// temporary copy of Go 1.6's private tls.Config.clone:
func cloneTLSConfig(c *tls.Config) *tls.Config ***REMOVED***
	return &tls.Config***REMOVED***
		Rand:                     c.Rand,
		Time:                     c.Time,
		Certificates:             c.Certificates,
		NameToCertificate:        c.NameToCertificate,
		GetCertificate:           c.GetCertificate,
		RootCAs:                  c.RootCAs,
		NextProtos:               c.NextProtos,
		ServerName:               c.ServerName,
		ClientAuth:               c.ClientAuth,
		ClientCAs:                c.ClientCAs,
		InsecureSkipVerify:       c.InsecureSkipVerify,
		CipherSuites:             c.CipherSuites,
		PreferServerCipherSuites: c.PreferServerCipherSuites,
		SessionTicketsDisabled:   c.SessionTicketsDisabled,
		SessionTicketKey:         c.SessionTicketKey,
		ClientSessionCache:       c.ClientSessionCache,
		MinVersion:               c.MinVersion,
		MaxVersion:               c.MaxVersion,
		CurvePreferences:         c.CurvePreferences,
	***REMOVED***
***REMOVED***

func (cc *ClientConn) Ping(ctx contextContext) error ***REMOVED***
	return cc.ping(ctx)
***REMOVED***

func (t *Transport) idleConnTimeout() time.Duration ***REMOVED*** return 0 ***REMOVED***
