// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.7

package http2

import (
	"context"
	"net"
	"net/http"
	"net/http/httptrace"
	"time"
)

type contextContext interface ***REMOVED***
	context.Context
***REMOVED***

func serverConnBaseContext(c net.Conn, opts *ServeConnOpts) (ctx contextContext, cancel func()) ***REMOVED***
	ctx, cancel = context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, http.LocalAddrContextKey, c.LocalAddr())
	if hs := opts.baseConfig(); hs != nil ***REMOVED***
		ctx = context.WithValue(ctx, http.ServerContextKey, hs)
	***REMOVED***
	return
***REMOVED***

func contextWithCancel(ctx contextContext) (_ contextContext, cancel func()) ***REMOVED***
	return context.WithCancel(ctx)
***REMOVED***

func requestWithContext(req *http.Request, ctx contextContext) *http.Request ***REMOVED***
	return req.WithContext(ctx)
***REMOVED***

type clientTrace httptrace.ClientTrace

func reqContext(r *http.Request) context.Context ***REMOVED*** return r.Context() ***REMOVED***

func (t *Transport) idleConnTimeout() time.Duration ***REMOVED***
	if t.t1 != nil ***REMOVED***
		return t.t1.IdleConnTimeout
	***REMOVED***
	return 0
***REMOVED***

func setResponseUncompressed(res *http.Response) ***REMOVED*** res.Uncompressed = true ***REMOVED***

func traceGotConn(req *http.Request, cc *ClientConn) ***REMOVED***
	trace := httptrace.ContextClientTrace(req.Context())
	if trace == nil || trace.GotConn == nil ***REMOVED***
		return
	***REMOVED***
	ci := httptrace.GotConnInfo***REMOVED***Conn: cc.tconn***REMOVED***
	cc.mu.Lock()
	ci.Reused = cc.nextStreamID > 1
	ci.WasIdle = len(cc.streams) == 0 && ci.Reused
	if ci.WasIdle && !cc.lastActive.IsZero() ***REMOVED***
		ci.IdleTime = time.Now().Sub(cc.lastActive)
	***REMOVED***
	cc.mu.Unlock()

	trace.GotConn(ci)
***REMOVED***

func traceWroteHeaders(trace *clientTrace) ***REMOVED***
	if trace != nil && trace.WroteHeaders != nil ***REMOVED***
		trace.WroteHeaders()
	***REMOVED***
***REMOVED***

func traceGot100Continue(trace *clientTrace) ***REMOVED***
	if trace != nil && trace.Got100Continue != nil ***REMOVED***
		trace.Got100Continue()
	***REMOVED***
***REMOVED***

func traceWait100Continue(trace *clientTrace) ***REMOVED***
	if trace != nil && trace.Wait100Continue != nil ***REMOVED***
		trace.Wait100Continue()
	***REMOVED***
***REMOVED***

func traceWroteRequest(trace *clientTrace, err error) ***REMOVED***
	if trace != nil && trace.WroteRequest != nil ***REMOVED***
		trace.WroteRequest(httptrace.WroteRequestInfo***REMOVED***Err: err***REMOVED***)
	***REMOVED***
***REMOVED***

func traceFirstResponseByte(trace *clientTrace) ***REMOVED***
	if trace != nil && trace.GotFirstResponseByte != nil ***REMOVED***
		trace.GotFirstResponseByte()
	***REMOVED***
***REMOVED***

func requestTrace(req *http.Request) *clientTrace ***REMOVED***
	trace := httptrace.ContextClientTrace(req.Context())
	return (*clientTrace)(trace)
***REMOVED***

// Ping sends a PING frame to the server and waits for the ack.
func (cc *ClientConn) Ping(ctx context.Context) error ***REMOVED***
	return cc.ping(ctx)
***REMOVED***
