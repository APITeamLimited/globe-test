// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Transport code's client connection pooling.

package http2

import (
	"crypto/tls"
	"net/http"
	"sync"
)

// ClientConnPool manages a pool of HTTP/2 client connections.
type ClientConnPool interface ***REMOVED***
	GetClientConn(req *http.Request, addr string) (*ClientConn, error)
	MarkDead(*ClientConn)
***REMOVED***

// clientConnPoolIdleCloser is the interface implemented by ClientConnPool
// implementations which can close their idle connections.
type clientConnPoolIdleCloser interface ***REMOVED***
	ClientConnPool
	closeIdleConnections()
***REMOVED***

var (
	_ clientConnPoolIdleCloser = (*clientConnPool)(nil)
	_ clientConnPoolIdleCloser = noDialClientConnPool***REMOVED******REMOVED***
)

// TODO: use singleflight for dialing and addConnCalls?
type clientConnPool struct ***REMOVED***
	t *Transport

	mu sync.Mutex // TODO: maybe switch to RWMutex
	// TODO: add support for sharing conns based on cert names
	// (e.g. share conn for googleapis.com and appspot.com)
	conns        map[string][]*ClientConn // key is host:port
	dialing      map[string]*dialCall     // currently in-flight dials
	keys         map[*ClientConn][]string
	addConnCalls map[string]*addConnCall // in-flight addConnIfNeede calls
***REMOVED***

func (p *clientConnPool) GetClientConn(req *http.Request, addr string) (*ClientConn, error) ***REMOVED***
	return p.getClientConn(req, addr, dialOnMiss)
***REMOVED***

const (
	dialOnMiss   = true
	noDialOnMiss = false
)

func (p *clientConnPool) getClientConn(req *http.Request, addr string, dialOnMiss bool) (*ClientConn, error) ***REMOVED***
	if isConnectionCloseRequest(req) && dialOnMiss ***REMOVED***
		// It gets its own connection.
		const singleUse = true
		cc, err := p.t.dialClientConn(addr, singleUse)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return cc, nil
	***REMOVED***
	p.mu.Lock()
	for _, cc := range p.conns[addr] ***REMOVED***
		if cc.CanTakeNewRequest() ***REMOVED***
			p.mu.Unlock()
			return cc, nil
		***REMOVED***
	***REMOVED***
	if !dialOnMiss ***REMOVED***
		p.mu.Unlock()
		return nil, ErrNoCachedConn
	***REMOVED***
	call := p.getStartDialLocked(addr)
	p.mu.Unlock()
	<-call.done
	return call.res, call.err
***REMOVED***

// dialCall is an in-flight Transport dial call to a host.
type dialCall struct ***REMOVED***
	p    *clientConnPool
	done chan struct***REMOVED******REMOVED*** // closed when done
	res  *ClientConn   // valid after done is closed
	err  error         // valid after done is closed
***REMOVED***

// requires p.mu is held.
func (p *clientConnPool) getStartDialLocked(addr string) *dialCall ***REMOVED***
	if call, ok := p.dialing[addr]; ok ***REMOVED***
		// A dial is already in-flight. Don't start another.
		return call
	***REMOVED***
	call := &dialCall***REMOVED***p: p, done: make(chan struct***REMOVED******REMOVED***)***REMOVED***
	if p.dialing == nil ***REMOVED***
		p.dialing = make(map[string]*dialCall)
	***REMOVED***
	p.dialing[addr] = call
	go call.dial(addr)
	return call
***REMOVED***

// run in its own goroutine.
func (c *dialCall) dial(addr string) ***REMOVED***
	const singleUse = false // shared conn
	c.res, c.err = c.p.t.dialClientConn(addr, singleUse)
	close(c.done)

	c.p.mu.Lock()
	delete(c.p.dialing, addr)
	if c.err == nil ***REMOVED***
		c.p.addConnLocked(addr, c.res)
	***REMOVED***
	c.p.mu.Unlock()
***REMOVED***

// addConnIfNeeded makes a NewClientConn out of c if a connection for key doesn't
// already exist. It coalesces concurrent calls with the same key.
// This is used by the http1 Transport code when it creates a new connection. Because
// the http1 Transport doesn't de-dup TCP dials to outbound hosts (because it doesn't know
// the protocol), it can get into a situation where it has multiple TLS connections.
// This code decides which ones live or die.
// The return value used is whether c was used.
// c is never closed.
func (p *clientConnPool) addConnIfNeeded(key string, t *Transport, c *tls.Conn) (used bool, err error) ***REMOVED***
	p.mu.Lock()
	for _, cc := range p.conns[key] ***REMOVED***
		if cc.CanTakeNewRequest() ***REMOVED***
			p.mu.Unlock()
			return false, nil
		***REMOVED***
	***REMOVED***
	call, dup := p.addConnCalls[key]
	if !dup ***REMOVED***
		if p.addConnCalls == nil ***REMOVED***
			p.addConnCalls = make(map[string]*addConnCall)
		***REMOVED***
		call = &addConnCall***REMOVED***
			p:    p,
			done: make(chan struct***REMOVED******REMOVED***),
		***REMOVED***
		p.addConnCalls[key] = call
		go call.run(t, key, c)
	***REMOVED***
	p.mu.Unlock()

	<-call.done
	if call.err != nil ***REMOVED***
		return false, call.err
	***REMOVED***
	return !dup, nil
***REMOVED***

type addConnCall struct ***REMOVED***
	p    *clientConnPool
	done chan struct***REMOVED******REMOVED*** // closed when done
	err  error
***REMOVED***

func (c *addConnCall) run(t *Transport, key string, tc *tls.Conn) ***REMOVED***
	cc, err := t.NewClientConn(tc)

	p := c.p
	p.mu.Lock()
	if err != nil ***REMOVED***
		c.err = err
	***REMOVED*** else ***REMOVED***
		p.addConnLocked(key, cc)
	***REMOVED***
	delete(p.addConnCalls, key)
	p.mu.Unlock()
	close(c.done)
***REMOVED***

func (p *clientConnPool) addConn(key string, cc *ClientConn) ***REMOVED***
	p.mu.Lock()
	p.addConnLocked(key, cc)
	p.mu.Unlock()
***REMOVED***

// p.mu must be held
func (p *clientConnPool) addConnLocked(key string, cc *ClientConn) ***REMOVED***
	for _, v := range p.conns[key] ***REMOVED***
		if v == cc ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	if p.conns == nil ***REMOVED***
		p.conns = make(map[string][]*ClientConn)
	***REMOVED***
	if p.keys == nil ***REMOVED***
		p.keys = make(map[*ClientConn][]string)
	***REMOVED***
	p.conns[key] = append(p.conns[key], cc)
	p.keys[cc] = append(p.keys[cc], key)
***REMOVED***

func (p *clientConnPool) MarkDead(cc *ClientConn) ***REMOVED***
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, key := range p.keys[cc] ***REMOVED***
		vv, ok := p.conns[key]
		if !ok ***REMOVED***
			continue
		***REMOVED***
		newList := filterOutClientConn(vv, cc)
		if len(newList) > 0 ***REMOVED***
			p.conns[key] = newList
		***REMOVED*** else ***REMOVED***
			delete(p.conns, key)
		***REMOVED***
	***REMOVED***
	delete(p.keys, cc)
***REMOVED***

func (p *clientConnPool) closeIdleConnections() ***REMOVED***
	p.mu.Lock()
	defer p.mu.Unlock()
	// TODO: don't close a cc if it was just added to the pool
	// milliseconds ago and has never been used. There's currently
	// a small race window with the HTTP/1 Transport's integration
	// where it can add an idle conn just before using it, and
	// somebody else can concurrently call CloseIdleConns and
	// break some caller's RoundTrip.
	for _, vv := range p.conns ***REMOVED***
		for _, cc := range vv ***REMOVED***
			cc.closeIfIdle()
		***REMOVED***
	***REMOVED***
***REMOVED***

func filterOutClientConn(in []*ClientConn, exclude *ClientConn) []*ClientConn ***REMOVED***
	out := in[:0]
	for _, v := range in ***REMOVED***
		if v != exclude ***REMOVED***
			out = append(out, v)
		***REMOVED***
	***REMOVED***
	// If we filtered it out, zero out the last item to prevent
	// the GC from seeing it.
	if len(in) != len(out) ***REMOVED***
		in[len(in)-1] = nil
	***REMOVED***
	return out
***REMOVED***

// noDialClientConnPool is an implementation of http2.ClientConnPool
// which never dials. We let the HTTP/1.1 client dial and use its TLS
// connection instead.
type noDialClientConnPool struct***REMOVED*** *clientConnPool ***REMOVED***

func (p noDialClientConnPool) GetClientConn(req *http.Request, addr string) (*ClientConn, error) ***REMOVED***
	return p.getClientConn(req, addr, noDialOnMiss)
***REMOVED***
