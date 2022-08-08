// Copyright 2012 Gary Burd
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package redis

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"errors"
	"io"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/garyburd/redigo/internal"
)

var (
	_ ConnWithTimeout = (*activeConn)(nil)
	_ ConnWithTimeout = (*errorConn)(nil)
)

var nowFunc = time.Now // for testing

// ErrPoolExhausted is returned from a pool connection method (Do, Send,
// Receive, Flush, Err) when the maximum number of database connections in the
// pool has been reached.
var ErrPoolExhausted = errors.New("redigo: connection pool exhausted")

var (
	errPoolClosed = errors.New("redigo: connection pool closed")
	errConnClosed = errors.New("redigo: connection closed")
)

// Pool maintains a pool of connections. The application calls the Get method
// to get a connection from the pool and the connection's Close method to
// return the connection's resources to the pool.
//
// The following example shows how to use a pool in a web application. The
// application creates a pool at application startup and makes it available to
// request handlers using a package level variable. The pool configuration used
// here is an example, not a recommendation.
//
//  func newPool(addr string) *redis.Pool ***REMOVED***
//    return &redis.Pool***REMOVED***
//      MaxIdle: 3,
//      IdleTimeout: 240 * time.Second,
//      Dial: func () (redis.Conn, error) ***REMOVED*** return redis.Dial("tcp", addr) ***REMOVED***,
//    ***REMOVED***
//  ***REMOVED***
//
//  var (
//    pool *redis.Pool
//    redisServer = flag.String("redisServer", ":6379", "")
//  )
//
//  func main() ***REMOVED***
//    flag.Parse()
//    pool = newPool(*redisServer)
//    ...
//  ***REMOVED***
//
// A request handler gets a connection from the pool and closes the connection
// when the handler is done:
//
//  func serveHome(w http.ResponseWriter, r *http.Request) ***REMOVED***
//      conn := pool.Get()
//      defer conn.Close()
//      ...
//  ***REMOVED***
//
// Use the Dial function to authenticate connections with the AUTH command or
// select a database with the SELECT command:
//
//  pool := &redis.Pool***REMOVED***
//    // Other pool configuration not shown in this example.
//    Dial: func () (redis.Conn, error) ***REMOVED***
//      c, err := redis.Dial("tcp", server)
//      if err != nil ***REMOVED***
//        return nil, err
//      ***REMOVED***
//      if _, err := c.Do("AUTH", password); err != nil ***REMOVED***
//        c.Close()
//        return nil, err
//      ***REMOVED***
//      if _, err := c.Do("SELECT", db); err != nil ***REMOVED***
//        c.Close()
//        return nil, err
//      ***REMOVED***
//      return c, nil
//    ***REMOVED***,
//  ***REMOVED***
//
// Use the TestOnBorrow function to check the health of an idle connection
// before the connection is returned to the application. This example PINGs
// connections that have been idle more than a minute:
//
//  pool := &redis.Pool***REMOVED***
//    // Other pool configuration not shown in this example.
//    TestOnBorrow: func(c redis.Conn, t time.Time) error ***REMOVED***
//      if time.Since(t) < time.Minute ***REMOVED***
//        return nil
//      ***REMOVED***
//      _, err := c.Do("PING")
//      return err
//    ***REMOVED***,
//  ***REMOVED***
//
type Pool struct ***REMOVED***
	// Dial is an application supplied function for creating and configuring a
	// connection.
	//
	// The connection returned from Dial must not be in a special state
	// (subscribed to pubsub channel, transaction started, ...).
	Dial func() (Conn, error)

	// TestOnBorrow is an optional application supplied function for checking
	// the health of an idle connection before the connection is used again by
	// the application. Argument t is the time that the connection was returned
	// to the pool. If the function returns an error, then the connection is
	// closed.
	TestOnBorrow func(c Conn, t time.Time) error

	// Maximum number of idle connections in the pool.
	MaxIdle int

	// Maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	MaxActive int

	// Close connections after remaining idle for this duration. If the value
	// is zero, then idle connections are not closed. Applications should set
	// the timeout to a value less than the server's timeout.
	IdleTimeout time.Duration

	// If Wait is true and the pool is at the MaxActive limit, then Get() waits
	// for a connection to be returned to the pool before returning.
	Wait bool

	// Close connections older than this duration. If the value is zero, then
	// the pool does not close connections based on age.
	MaxConnLifetime time.Duration

	chInitialized uint32 // set to 1 when field ch is initialized

	mu     sync.Mutex    // mu protects the following fields
	closed bool          // set to true when the pool is closed.
	active int           // the number of open connections in the pool
	ch     chan struct***REMOVED******REMOVED*** // limits open connections when p.Wait is true
	idle   idleList      // idle connections
***REMOVED***

// NewPool creates a new pool.
//
// Deprecated: Initialize the Pool directory as shown in the example.
func NewPool(newFn func() (Conn, error), maxIdle int) *Pool ***REMOVED***
	return &Pool***REMOVED***Dial: newFn, MaxIdle: maxIdle***REMOVED***
***REMOVED***

// Get gets a connection. The application must close the returned connection.
// This method always returns a valid connection so that applications can defer
// error handling to the first use of the connection. If there is an error
// getting an underlying connection, then the connection Err, Do, Send, Flush
// and Receive methods return that error.
func (p *Pool) Get() Conn ***REMOVED***
	pc, err := p.get(nil)
	if err != nil ***REMOVED***
		return errorConn***REMOVED***err***REMOVED***
	***REMOVED***
	return &activeConn***REMOVED***p: p, pc: pc***REMOVED***
***REMOVED***

// PoolStats contains pool statistics.
type PoolStats struct ***REMOVED***
	// ActiveCount is the number of connections in the pool. The count includes
	// idle connections and connections in use.
	ActiveCount int
	// IdleCount is the number of idle connections in the pool.
	IdleCount int
***REMOVED***

// Stats returns pool's statistics.
func (p *Pool) Stats() PoolStats ***REMOVED***
	p.mu.Lock()
	stats := PoolStats***REMOVED***
		ActiveCount: p.active,
		IdleCount:   p.idle.count,
	***REMOVED***
	p.mu.Unlock()

	return stats
***REMOVED***

// ActiveCount returns the number of connections in the pool. The count
// includes idle connections and connections in use.
func (p *Pool) ActiveCount() int ***REMOVED***
	p.mu.Lock()
	active := p.active
	p.mu.Unlock()
	return active
***REMOVED***

// IdleCount returns the number of idle connections in the pool.
func (p *Pool) IdleCount() int ***REMOVED***
	p.mu.Lock()
	idle := p.idle.count
	p.mu.Unlock()
	return idle
***REMOVED***

// Close releases the resources used by the pool.
func (p *Pool) Close() error ***REMOVED***
	p.mu.Lock()
	if p.closed ***REMOVED***
		p.mu.Unlock()
		return nil
	***REMOVED***
	p.closed = true
	p.active -= p.idle.count
	pc := p.idle.front
	p.idle.count = 0
	p.idle.front, p.idle.back = nil, nil
	if p.ch != nil ***REMOVED***
		close(p.ch)
	***REMOVED***
	p.mu.Unlock()
	for ; pc != nil; pc = pc.next ***REMOVED***
		pc.c.Close()
	***REMOVED***
	return nil
***REMOVED***

func (p *Pool) lazyInit() ***REMOVED***
	// Fast path.
	if atomic.LoadUint32(&p.chInitialized) == 1 ***REMOVED***
		return
	***REMOVED***
	// Slow path.
	p.mu.Lock()
	if p.chInitialized == 0 ***REMOVED***
		p.ch = make(chan struct***REMOVED******REMOVED***, p.MaxActive)
		if p.closed ***REMOVED***
			close(p.ch)
		***REMOVED*** else ***REMOVED***
			for i := 0; i < p.MaxActive; i++ ***REMOVED***
				p.ch <- struct***REMOVED******REMOVED******REMOVED******REMOVED***
			***REMOVED***
		***REMOVED***
		atomic.StoreUint32(&p.chInitialized, 1)
	***REMOVED***
	p.mu.Unlock()
***REMOVED***

// get prunes stale connections and returns a connection from the idle list or
// creates a new connection.
func (p *Pool) get(ctx interface ***REMOVED***
	Done() <-chan struct***REMOVED******REMOVED***
	Err() error
***REMOVED***) (*poolConn, error) ***REMOVED***

	// Handle limit for p.Wait == true.
	if p.Wait && p.MaxActive > 0 ***REMOVED***
		p.lazyInit()
		if ctx == nil ***REMOVED***
			<-p.ch
		***REMOVED*** else ***REMOVED***
			select ***REMOVED***
			case <-p.ch:
			case <-ctx.Done():
				return nil, ctx.Err()
			***REMOVED***
		***REMOVED***
	***REMOVED***

	p.mu.Lock()

	// Prune stale connections at the back of the idle list.
	if p.IdleTimeout > 0 ***REMOVED***
		n := p.idle.count
		for i := 0; i < n && p.idle.back != nil && p.idle.back.t.Add(p.IdleTimeout).Before(nowFunc()); i++ ***REMOVED***
			pc := p.idle.back
			p.idle.popBack()
			p.mu.Unlock()
			pc.c.Close()
			p.mu.Lock()
			p.active--
		***REMOVED***
	***REMOVED***

	// Get idle connection from the front of idle list.
	for p.idle.front != nil ***REMOVED***
		pc := p.idle.front
		p.idle.popFront()
		p.mu.Unlock()
		if (p.TestOnBorrow == nil || p.TestOnBorrow(pc.c, pc.t) == nil) &&
			(p.MaxConnLifetime == 0 || nowFunc().Sub(pc.created) < p.MaxConnLifetime) ***REMOVED***
			return pc, nil
		***REMOVED***
		pc.c.Close()
		p.mu.Lock()
		p.active--
	***REMOVED***

	// Check for pool closed before dialing a new connection.
	if p.closed ***REMOVED***
		p.mu.Unlock()
		return nil, errors.New("redigo: get on closed pool")
	***REMOVED***

	// Handle limit for p.Wait == false.
	if !p.Wait && p.MaxActive > 0 && p.active >= p.MaxActive ***REMOVED***
		p.mu.Unlock()
		return nil, ErrPoolExhausted
	***REMOVED***

	p.active++
	p.mu.Unlock()
	c, err := p.Dial()
	if err != nil ***REMOVED***
		c = nil
		p.mu.Lock()
		p.active--
		if p.ch != nil && !p.closed ***REMOVED***
			p.ch <- struct***REMOVED******REMOVED******REMOVED******REMOVED***
		***REMOVED***
		p.mu.Unlock()
	***REMOVED***
	return &poolConn***REMOVED***c: c, created: nowFunc()***REMOVED***, err
***REMOVED***

func (p *Pool) put(pc *poolConn, forceClose bool) error ***REMOVED***
	p.mu.Lock()
	if !p.closed && !forceClose ***REMOVED***
		pc.t = nowFunc()
		p.idle.pushFront(pc)
		if p.idle.count > p.MaxIdle ***REMOVED***
			pc = p.idle.back
			p.idle.popBack()
		***REMOVED*** else ***REMOVED***
			pc = nil
		***REMOVED***
	***REMOVED***

	if pc != nil ***REMOVED***
		p.mu.Unlock()
		pc.c.Close()
		p.mu.Lock()
		p.active--
	***REMOVED***

	if p.ch != nil && !p.closed ***REMOVED***
		p.ch <- struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***
	p.mu.Unlock()
	return nil
***REMOVED***

type activeConn struct ***REMOVED***
	p     *Pool
	pc    *poolConn
	state int
***REMOVED***

var (
	sentinel     []byte
	sentinelOnce sync.Once
)

func initSentinel() ***REMOVED***
	p := make([]byte, 64)
	if _, err := rand.Read(p); err == nil ***REMOVED***
		sentinel = p
	***REMOVED*** else ***REMOVED***
		h := sha1.New()
		io.WriteString(h, "Oops, rand failed. Use time instead.")
		io.WriteString(h, strconv.FormatInt(time.Now().UnixNano(), 10))
		sentinel = h.Sum(nil)
	***REMOVED***
***REMOVED***

func (ac *activeConn) Close() error ***REMOVED***
	pc := ac.pc
	if pc == nil ***REMOVED***
		return nil
	***REMOVED***
	ac.pc = nil

	if ac.state&internal.MultiState != 0 ***REMOVED***
		pc.c.Send("DISCARD")
		ac.state &^= (internal.MultiState | internal.WatchState)
	***REMOVED*** else if ac.state&internal.WatchState != 0 ***REMOVED***
		pc.c.Send("UNWATCH")
		ac.state &^= internal.WatchState
	***REMOVED***
	if ac.state&internal.SubscribeState != 0 ***REMOVED***
		pc.c.Send("UNSUBSCRIBE")
		pc.c.Send("PUNSUBSCRIBE")
		// To detect the end of the message stream, ask the server to echo
		// a sentinel value and read until we see that value.
		sentinelOnce.Do(initSentinel)
		pc.c.Send("ECHO", sentinel)
		pc.c.Flush()
		for ***REMOVED***
			p, err := pc.c.Receive()
			if err != nil ***REMOVED***
				break
			***REMOVED***
			if p, ok := p.([]byte); ok && bytes.Equal(p, sentinel) ***REMOVED***
				ac.state &^= internal.SubscribeState
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
	pc.c.Do("")
	ac.p.put(pc, ac.state != 0 || pc.c.Err() != nil)
	return nil
***REMOVED***

func (ac *activeConn) Err() error ***REMOVED***
	pc := ac.pc
	if pc == nil ***REMOVED***
		return errConnClosed
	***REMOVED***
	return pc.c.Err()
***REMOVED***

func (ac *activeConn) Do(commandName string, args ...interface***REMOVED******REMOVED***) (reply interface***REMOVED******REMOVED***, err error) ***REMOVED***
	pc := ac.pc
	if pc == nil ***REMOVED***
		return nil, errConnClosed
	***REMOVED***
	ci := internal.LookupCommandInfo(commandName)
	ac.state = (ac.state | ci.Set) &^ ci.Clear
	return pc.c.Do(commandName, args...)
***REMOVED***

func (ac *activeConn) DoWithTimeout(timeout time.Duration, commandName string, args ...interface***REMOVED******REMOVED***) (reply interface***REMOVED******REMOVED***, err error) ***REMOVED***
	pc := ac.pc
	if pc == nil ***REMOVED***
		return nil, errConnClosed
	***REMOVED***
	cwt, ok := pc.c.(ConnWithTimeout)
	if !ok ***REMOVED***
		return nil, errTimeoutNotSupported
	***REMOVED***
	ci := internal.LookupCommandInfo(commandName)
	ac.state = (ac.state | ci.Set) &^ ci.Clear
	return cwt.DoWithTimeout(timeout, commandName, args...)
***REMOVED***

func (ac *activeConn) Send(commandName string, args ...interface***REMOVED******REMOVED***) error ***REMOVED***
	pc := ac.pc
	if pc == nil ***REMOVED***
		return errConnClosed
	***REMOVED***
	ci := internal.LookupCommandInfo(commandName)
	ac.state = (ac.state | ci.Set) &^ ci.Clear
	return pc.c.Send(commandName, args...)
***REMOVED***

func (ac *activeConn) Flush() error ***REMOVED***
	pc := ac.pc
	if pc == nil ***REMOVED***
		return errConnClosed
	***REMOVED***
	return pc.c.Flush()
***REMOVED***

func (ac *activeConn) Receive() (reply interface***REMOVED******REMOVED***, err error) ***REMOVED***
	pc := ac.pc
	if pc == nil ***REMOVED***
		return nil, errConnClosed
	***REMOVED***
	return pc.c.Receive()
***REMOVED***

func (ac *activeConn) ReceiveWithTimeout(timeout time.Duration) (reply interface***REMOVED******REMOVED***, err error) ***REMOVED***
	pc := ac.pc
	if pc == nil ***REMOVED***
		return nil, errConnClosed
	***REMOVED***
	cwt, ok := pc.c.(ConnWithTimeout)
	if !ok ***REMOVED***
		return nil, errTimeoutNotSupported
	***REMOVED***
	return cwt.ReceiveWithTimeout(timeout)
***REMOVED***

type errorConn struct***REMOVED*** err error ***REMOVED***

func (ec errorConn) Do(string, ...interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED*** return nil, ec.err ***REMOVED***
func (ec errorConn) DoWithTimeout(time.Duration, string, ...interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	return nil, ec.err
***REMOVED***
func (ec errorConn) Send(string, ...interface***REMOVED******REMOVED***) error                     ***REMOVED*** return ec.err ***REMOVED***
func (ec errorConn) Err() error                                            ***REMOVED*** return ec.err ***REMOVED***
func (ec errorConn) Close() error                                          ***REMOVED*** return nil ***REMOVED***
func (ec errorConn) Flush() error                                          ***REMOVED*** return ec.err ***REMOVED***
func (ec errorConn) Receive() (interface***REMOVED******REMOVED***, error)                         ***REMOVED*** return nil, ec.err ***REMOVED***
func (ec errorConn) ReceiveWithTimeout(time.Duration) (interface***REMOVED******REMOVED***, error) ***REMOVED*** return nil, ec.err ***REMOVED***

type idleList struct ***REMOVED***
	count       int
	front, back *poolConn
***REMOVED***

type poolConn struct ***REMOVED***
	c          Conn
	t          time.Time
	created    time.Time
	next, prev *poolConn
***REMOVED***

func (l *idleList) pushFront(pc *poolConn) ***REMOVED***
	pc.next = l.front
	pc.prev = nil
	if l.count == 0 ***REMOVED***
		l.back = pc
	***REMOVED*** else ***REMOVED***
		l.front.prev = pc
	***REMOVED***
	l.front = pc
	l.count++
	return
***REMOVED***

func (l *idleList) popFront() ***REMOVED***
	pc := l.front
	l.count--
	if l.count == 0 ***REMOVED***
		l.front, l.back = nil, nil
	***REMOVED*** else ***REMOVED***
		pc.next.prev = nil
		l.front = pc.next
	***REMOVED***
	pc.next, pc.prev = nil, nil
***REMOVED***

func (l *idleList) popBack() ***REMOVED***
	pc := l.back
	l.count--
	if l.count == 0 ***REMOVED***
		l.front, l.back = nil, nil
	***REMOVED*** else ***REMOVED***
		pc.prev.next = nil
		l.back = pc.prev
	***REMOVED***
	pc.next, pc.prev = nil, nil
***REMOVED***
