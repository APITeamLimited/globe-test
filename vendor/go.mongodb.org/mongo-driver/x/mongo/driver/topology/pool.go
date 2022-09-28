// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package topology

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/address"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
)

// Connection pool state constants.
const (
	poolPaused int = iota
	poolReady
	poolClosed
)

// ErrPoolNotPaused is returned when attempting to mark a connection pool "ready" that is not
// currently "paused".
var ErrPoolNotPaused = PoolError("only a paused pool can be marked ready")

// ErrPoolClosed is returned when attempting to check out a connection from a closed pool.
var ErrPoolClosed = PoolError("attempted to check out a connection from closed connection pool")

// ErrConnectionClosed is returned from an attempt to use an already closed connection.
var ErrConnectionClosed = ConnectionError***REMOVED***ConnectionID: "<closed>", message: "connection is closed"***REMOVED***

// ErrWrongPool is return when a connection is returned to a pool it doesn't belong to.
var ErrWrongPool = PoolError("connection does not belong to this pool")

// PoolError is an error returned from a Pool method.
type PoolError string

func (pe PoolError) Error() string ***REMOVED*** return string(pe) ***REMOVED***

// poolClearedError is an error returned when the connection pool is cleared or currently paused. It
// is a retryable error.
type poolClearedError struct ***REMOVED***
	err     error
	address address.Address
***REMOVED***

func (pce poolClearedError) Error() string ***REMOVED***
	return fmt.Sprintf(
		"connection pool for %v was cleared because another operation failed with: %v",
		pce.address,
		pce.err)
***REMOVED***

// Retryable returns true. All poolClearedErrors are retryable.
func (poolClearedError) Retryable() bool ***REMOVED*** return true ***REMOVED***

// Assert that poolClearedError is a driver.RetryablePoolError.
var _ driver.RetryablePoolError = poolClearedError***REMOVED******REMOVED***

// poolConfig contains all aspects of the pool that can be configured
type poolConfig struct ***REMOVED***
	Address          address.Address
	MinPoolSize      uint64
	MaxPoolSize      uint64
	MaxConnecting    uint64
	MaxIdleTime      time.Duration
	MaintainInterval time.Duration
	PoolMonitor      *event.PoolMonitor
	handshakeErrFn   func(error, uint64, *primitive.ObjectID)
***REMOVED***

type pool struct ***REMOVED***
	// The following integer fields must be accessed using the atomic package
	// and should be at the beginning of the struct.
	// - atomic bug: https://pkg.go.dev/sync/atomic#pkg-note-BUG
	// - suggested layout: https://go101.org/article/memory-layout.html

	nextID                       uint64 // nextID is the next pool ID for a new connection.
	pinnedCursorConnections      uint64
	pinnedTransactionConnections uint64

	address       address.Address
	minSize       uint64
	maxSize       uint64
	maxConnecting uint64
	monitor       *event.PoolMonitor

	// handshakeErrFn is used to handle any errors that happen during connection establishment and
	// handshaking.
	handshakeErrFn func(error, uint64, *primitive.ObjectID)

	connOpts   []ConnectionOption
	generation *poolGenerationMap

	maintainInterval time.Duration   // maintainInterval is the maintain() loop interval.
	maintainReady    chan struct***REMOVED******REMOVED***   // maintainReady is a signal channel that starts the maintain() loop when ready() is called.
	backgroundDone   *sync.WaitGroup // backgroundDone waits for all background goroutines to return.

	stateMu      sync.RWMutex // stateMu guards state, lastClearErr
	state        int          // state is the current state of the connection pool.
	lastClearErr error        // lastClearErr is the last error that caused the pool to be cleared.

	// createConnectionsCond is the condition variable that controls when the createConnections()
	// loop runs or waits. Its lock guards cancelBackgroundCtx, conns, and newConnWait. Any changes
	// to the state of the guarded values must be made while holding the lock to prevent undefined
	// behavior in the createConnections() waiting logic.
	createConnectionsCond *sync.Cond
	cancelBackgroundCtx   context.CancelFunc     // cancelBackgroundCtx is called to signal background goroutines to stop.
	conns                 map[uint64]*connection // conns holds all currently open connections.
	newConnWait           wantConnQueue          // newConnWait holds all wantConn requests for new connections.

	idleMu       sync.Mutex    // idleMu guards idleConns, idleConnWait
	idleConns    []*connection // idleConns holds all idle connections.
	idleConnWait wantConnQueue // idleConnWait holds all wantConn requests for idle connections.
***REMOVED***

// getState returns the current state of the pool. Callers must not hold the stateMu lock.
func (p *pool) getState() int ***REMOVED***
	p.stateMu.RLock()
	defer p.stateMu.RUnlock()

	return p.state
***REMOVED***

// connectionPerished checks if a given connection is perished and should be removed from the pool.
func connectionPerished(conn *connection) (string, bool) ***REMOVED***
	switch ***REMOVED***
	case conn.closed():
		// A connection would only be closed if it encountered a network error during an operation and closed itself.
		return event.ReasonError, true
	case conn.idleTimeoutExpired():
		return event.ReasonIdle, true
	case conn.pool.stale(conn):
		return event.ReasonStale, true
	***REMOVED***
	return "", false
***REMOVED***

// newPool creates a new pool. It will use the provided options when creating connections.
func newPool(config poolConfig, connOpts ...ConnectionOption) *pool ***REMOVED***
	if config.MaxIdleTime != time.Duration(0) ***REMOVED***
		connOpts = append(connOpts, WithIdleTimeout(func(_ time.Duration) time.Duration ***REMOVED*** return config.MaxIdleTime ***REMOVED***))
	***REMOVED***

	var maxConnecting uint64 = 2
	if config.MaxConnecting > 0 ***REMOVED***
		maxConnecting = config.MaxConnecting
	***REMOVED***

	maintainInterval := 10 * time.Second
	if config.MaintainInterval != 0 ***REMOVED***
		maintainInterval = config.MaintainInterval
	***REMOVED***

	pool := &pool***REMOVED***
		address:               config.Address,
		minSize:               config.MinPoolSize,
		maxSize:               config.MaxPoolSize,
		maxConnecting:         maxConnecting,
		monitor:               config.PoolMonitor,
		handshakeErrFn:        config.handshakeErrFn,
		connOpts:              connOpts,
		generation:            newPoolGenerationMap(),
		state:                 poolPaused,
		maintainInterval:      maintainInterval,
		maintainReady:         make(chan struct***REMOVED******REMOVED***, 1),
		backgroundDone:        &sync.WaitGroup***REMOVED******REMOVED***,
		createConnectionsCond: sync.NewCond(&sync.Mutex***REMOVED******REMOVED***),
		conns:                 make(map[uint64]*connection, config.MaxPoolSize),
		idleConns:             make([]*connection, 0, config.MaxPoolSize),
	***REMOVED***
	// minSize must not exceed maxSize if maxSize is not 0
	if pool.maxSize != 0 && pool.minSize > pool.maxSize ***REMOVED***
		pool.minSize = pool.maxSize
	***REMOVED***
	pool.connOpts = append(pool.connOpts, withGenerationNumberFn(func(_ generationNumberFn) generationNumberFn ***REMOVED*** return pool.getGenerationForNewConnection ***REMOVED***))

	pool.generation.connect()

	// Create a Context with cancellation that's used to signal the createConnections() and
	// maintain() background goroutines to stop. Also create a "backgroundDone" WaitGroup that is
	// used to wait for the background goroutines to return.
	var ctx context.Context
	ctx, pool.cancelBackgroundCtx = context.WithCancel(context.Background())

	for i := 0; i < int(pool.maxConnecting); i++ ***REMOVED***
		pool.backgroundDone.Add(1)
		go pool.createConnections(ctx, pool.backgroundDone)
	***REMOVED***

	// If maintainInterval is not positive, don't start the maintain() goroutine. Expect that
	// negative values are only used in testing; this config value is not user-configurable.
	if maintainInterval > 0 ***REMOVED***
		pool.backgroundDone.Add(1)
		go pool.maintain(ctx, pool.backgroundDone)
	***REMOVED***

	if pool.monitor != nil ***REMOVED***
		pool.monitor.Event(&event.PoolEvent***REMOVED***
			Type: event.PoolCreated,
			PoolOptions: &event.MonitorPoolOptions***REMOVED***
				MaxPoolSize: config.MaxPoolSize,
				MinPoolSize: config.MinPoolSize,
			***REMOVED***,
			Address: pool.address.String(),
		***REMOVED***)
	***REMOVED***

	return pool
***REMOVED***

// stale checks if a given connection's generation is below the generation of the pool
func (p *pool) stale(conn *connection) bool ***REMOVED***
	return conn == nil || p.generation.stale(conn.desc.ServiceID, conn.generation)
***REMOVED***

// ready puts the pool into the "ready" state and starts the background connection creation and
// monitoring goroutines. ready must be called before connections can be checked out. An unused,
// connected pool must be closed or it will leak goroutines and will not be garbage collected.
func (p *pool) ready() error ***REMOVED***
	// While holding the stateMu lock, set the pool to "ready" if it is currently "paused".
	p.stateMu.Lock()
	if p.state == poolReady ***REMOVED***
		p.stateMu.Unlock()
		return nil
	***REMOVED***
	if p.state != poolPaused ***REMOVED***
		p.stateMu.Unlock()
		return ErrPoolNotPaused
	***REMOVED***
	p.lastClearErr = nil
	p.state = poolReady
	p.stateMu.Unlock()

	// Signal maintain() to wake up immediately when marking the pool "ready".
	select ***REMOVED***
	case p.maintainReady <- struct***REMOVED******REMOVED******REMOVED******REMOVED***:
	default:
	***REMOVED***

	if p.monitor != nil ***REMOVED***
		p.monitor.Event(&event.PoolEvent***REMOVED***
			Type:    event.PoolReady,
			Address: p.address.String(),
		***REMOVED***)
	***REMOVED***

	return nil
***REMOVED***

// close closes the pool, closes all connections associated with the pool, and stops all background
// goroutines. All subsequent checkOut requests will return an error. An unused, ready pool must be
// closed or it will leak goroutines and will not be garbage collected.
func (p *pool) close(ctx context.Context) ***REMOVED***
	p.stateMu.Lock()
	if p.state == poolClosed ***REMOVED***
		p.stateMu.Unlock()
		return
	***REMOVED***
	p.state = poolClosed
	p.stateMu.Unlock()

	// Call cancelBackgroundCtx() to exit the maintain() and createConnections() background
	// goroutines. Broadcast to the createConnectionsCond to wake up all createConnections()
	// goroutines. We must hold the createConnectionsCond lock here because we're changing the
	// condition by cancelling the "background goroutine" Context, even tho cancelling the Context
	// is also synchronized by a lock. Otherwise, we run into an intermittent bug that prevents the
	// createConnections() goroutines from exiting.
	p.createConnectionsCond.L.Lock()
	p.cancelBackgroundCtx()
	p.createConnectionsCond.Broadcast()
	p.createConnectionsCond.L.Unlock()

	// Wait for all background goroutines to exit.
	p.backgroundDone.Wait()

	p.generation.disconnect()

	if ctx == nil ***REMOVED***
		ctx = context.Background()
	***REMOVED***

	// If we have a deadline then we interpret it as a request to gracefully shutdown. We wait until
	// either all the connections have been checked back into the pool (i.e. total open connections
	// equals idle connections) or until the Context deadline is reached.
	if _, ok := ctx.Deadline(); ok ***REMOVED***
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

	graceful:
		for ***REMOVED***
			if p.totalConnectionCount() == p.availableConnectionCount() ***REMOVED***
				break graceful
			***REMOVED***

			select ***REMOVED***
			case <-ticker.C:
			case <-ctx.Done():
				break graceful
			default:
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Empty the idle connections stack and try to deliver ErrPoolClosed to any waiting wantConns
	// from idleConnWait while holding the idleMu lock.
	p.idleMu.Lock()
	p.idleConns = p.idleConns[:0]
	for ***REMOVED***
		w := p.idleConnWait.popFront()
		if w == nil ***REMOVED***
			break
		***REMOVED***
		w.tryDeliver(nil, ErrPoolClosed)
	***REMOVED***
	p.idleMu.Unlock()

	// Collect all conns from the pool and try to deliver ErrPoolClosed to any waiting wantConns
	// from newConnWait while holding the createConnectionsCond lock. We can't call removeConnection
	// on the connections while holding any locks, so do that after we release the lock.
	p.createConnectionsCond.L.Lock()
	conns := make([]*connection, 0, len(p.conns))
	for _, conn := range p.conns ***REMOVED***
		conns = append(conns, conn)
	***REMOVED***
	for ***REMOVED***
		w := p.newConnWait.popFront()
		if w == nil ***REMOVED***
			break
		***REMOVED***
		w.tryDeliver(nil, ErrPoolClosed)
	***REMOVED***
	p.createConnectionsCond.L.Unlock()

	// Now that we're not holding any locks, remove all of the connections we collected from the
	// pool.
	for _, conn := range conns ***REMOVED***
		_ = p.removeConnection(conn, event.ReasonPoolClosed)
		_ = p.closeConnection(conn) // We don't care about errors while closing the connection.
	***REMOVED***

	if p.monitor != nil ***REMOVED***
		p.monitor.Event(&event.PoolEvent***REMOVED***
			Type:    event.PoolClosedEvent,
			Address: p.address.String(),
		***REMOVED***)
	***REMOVED***
***REMOVED***

func (p *pool) pinConnectionToCursor() ***REMOVED***
	atomic.AddUint64(&p.pinnedCursorConnections, 1)
***REMOVED***

func (p *pool) unpinConnectionFromCursor() ***REMOVED***
	// See https://golang.org/pkg/sync/atomic/#AddUint64 for an explanation of the ^uint64(0) syntax.
	atomic.AddUint64(&p.pinnedCursorConnections, ^uint64(0))
***REMOVED***

func (p *pool) pinConnectionToTransaction() ***REMOVED***
	atomic.AddUint64(&p.pinnedTransactionConnections, 1)
***REMOVED***

func (p *pool) unpinConnectionFromTransaction() ***REMOVED***
	// See https://golang.org/pkg/sync/atomic/#AddUint64 for an explanation of the ^uint64(0) syntax.
	atomic.AddUint64(&p.pinnedTransactionConnections, ^uint64(0))
***REMOVED***

// checkOut checks out a connection from the pool. If an idle connection is not available, the
// checkOut enters a queue waiting for either the next idle or new connection. If the pool is not
// ready, checkOut returns an error.
// Based partially on https://cs.opensource.google/go/go/+/refs/tags/go1.16.6:src/net/http/transport.go;l=1324
func (p *pool) checkOut(ctx context.Context) (conn *connection, err error) ***REMOVED***
	// TODO(CSOT): If a Timeout was specified at any level, respect the Timeout is server selection, connection
	// TODO checkout.
	if p.monitor != nil ***REMOVED***
		p.monitor.Event(&event.PoolEvent***REMOVED***
			Type:    event.GetStarted,
			Address: p.address.String(),
		***REMOVED***)
	***REMOVED***

	// Check the pool state while holding a stateMu read lock. If the pool state is not "ready",
	// return an error. Do all of this while holding the stateMu read lock to prevent a state change between
	// checking the state and entering the wait queue. Not holding the stateMu read lock here may
	// allow a checkOut() to enter the wait queue after clear() pauses the pool and clears the wait
	// queue, resulting in createConnections() doing work while the pool is "paused".
	p.stateMu.RLock()
	switch p.state ***REMOVED***
	case poolClosed:
		p.stateMu.RUnlock()
		if p.monitor != nil ***REMOVED***
			p.monitor.Event(&event.PoolEvent***REMOVED***
				Type:    event.GetFailed,
				Address: p.address.String(),
				Reason:  event.ReasonPoolClosed,
			***REMOVED***)
		***REMOVED***
		return nil, ErrPoolClosed
	case poolPaused:
		err := poolClearedError***REMOVED***err: p.lastClearErr, address: p.address***REMOVED***
		p.stateMu.RUnlock()
		if p.monitor != nil ***REMOVED***
			p.monitor.Event(&event.PoolEvent***REMOVED***
				Type:    event.GetFailed,
				Address: p.address.String(),
				Reason:  event.ReasonConnectionErrored,
			***REMOVED***)
		***REMOVED***
		return nil, err
	***REMOVED***

	if ctx == nil ***REMOVED***
		ctx = context.Background()
	***REMOVED***

	// Create a wantConn, which we will use to request an existing idle or new connection. Always
	// cancel the wantConn if checkOut() returned an error to make sure any delivered connections
	// are returned to the pool (e.g. if a connection was delivered immediately after the Context
	// timed out).
	w := newWantConn()
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			w.cancel(p, err)
		***REMOVED***
	***REMOVED***()

	// Get in the queue for an idle connection. If getOrQueueForIdleConn returns true, it was able to
	// immediately deliver an idle connection to the wantConn, so we can return the connection or
	// error from the wantConn without waiting for "ready".
	if delivered := p.getOrQueueForIdleConn(w); delivered ***REMOVED***
		// If delivered = true, we didn't enter the wait queue and will return either a connection
		// or an error, so unlock the stateMu lock here.
		p.stateMu.RUnlock()

		if w.err != nil ***REMOVED***
			if p.monitor != nil ***REMOVED***
				p.monitor.Event(&event.PoolEvent***REMOVED***
					Type:    event.GetFailed,
					Address: p.address.String(),
					Reason:  event.ReasonConnectionErrored,
				***REMOVED***)
			***REMOVED***
			return nil, w.err
		***REMOVED***

		if p.monitor != nil ***REMOVED***
			p.monitor.Event(&event.PoolEvent***REMOVED***
				Type:         event.GetSucceeded,
				Address:      p.address.String(),
				ConnectionID: w.conn.poolID,
			***REMOVED***)
		***REMOVED***
		return w.conn, nil
	***REMOVED***

	// If we didn't get an immediately available idle connection, also get in the queue for a new
	// connection while we're waiting for an idle connection.
	p.queueForNewConn(w)
	p.stateMu.RUnlock()

	// Wait for either the wantConn to be ready or for the Context to time out.
	select ***REMOVED***
	case <-w.ready:
		if w.err != nil ***REMOVED***
			if p.monitor != nil ***REMOVED***
				p.monitor.Event(&event.PoolEvent***REMOVED***
					Type:    event.GetFailed,
					Address: p.address.String(),
					Reason:  event.ReasonConnectionErrored,
				***REMOVED***)
			***REMOVED***
			return nil, w.err
		***REMOVED***

		if p.monitor != nil ***REMOVED***
			p.monitor.Event(&event.PoolEvent***REMOVED***
				Type:         event.GetSucceeded,
				Address:      p.address.String(),
				ConnectionID: w.conn.poolID,
			***REMOVED***)
		***REMOVED***
		return w.conn, nil
	case <-ctx.Done():
		if p.monitor != nil ***REMOVED***
			p.monitor.Event(&event.PoolEvent***REMOVED***
				Type:    event.GetFailed,
				Address: p.address.String(),
				Reason:  event.ReasonTimedOut,
			***REMOVED***)
		***REMOVED***
		return nil, WaitQueueTimeoutError***REMOVED***
			Wrapped:                      ctx.Err(),
			PinnedCursorConnections:      atomic.LoadUint64(&p.pinnedCursorConnections),
			PinnedTransactionConnections: atomic.LoadUint64(&p.pinnedTransactionConnections),
			maxPoolSize:                  p.maxSize,
			totalConnectionCount:         p.totalConnectionCount(),
		***REMOVED***
	***REMOVED***
***REMOVED***

// closeConnection closes a connection.
func (p *pool) closeConnection(conn *connection) error ***REMOVED***
	if conn.pool != p ***REMOVED***
		return ErrWrongPool
	***REMOVED***

	if atomic.LoadInt64(&conn.state) == connConnected ***REMOVED***
		conn.closeConnectContext()
		conn.wait() // Make sure that the connection has finished connecting.
	***REMOVED***

	err := conn.close()
	if err != nil ***REMOVED***
		return ConnectionError***REMOVED***ConnectionID: conn.id, Wrapped: err, message: "failed to close net.Conn"***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (p *pool) getGenerationForNewConnection(serviceID *primitive.ObjectID) uint64 ***REMOVED***
	return p.generation.addConnection(serviceID)
***REMOVED***

// removeConnection removes a connection from the pool and emits a "ConnectionClosed" event.
func (p *pool) removeConnection(conn *connection, reason string) error ***REMOVED***
	if conn == nil ***REMOVED***
		return nil
	***REMOVED***

	if conn.pool != p ***REMOVED***
		return ErrWrongPool
	***REMOVED***

	p.createConnectionsCond.L.Lock()
	_, ok := p.conns[conn.poolID]
	if !ok ***REMOVED***
		// If the connection has been removed from the pool already, exit without doing any
		// additional state changes.
		p.createConnectionsCond.L.Unlock()
		return nil
	***REMOVED***
	delete(p.conns, conn.poolID)
	// Signal the createConnectionsCond so any goroutines waiting for a new connection slot in the
	// pool will proceed.
	p.createConnectionsCond.Signal()
	p.createConnectionsCond.L.Unlock()

	// Only update the generation numbers map if the connection has retrieved its generation number.
	// Otherwise, we'd decrement the count for the generation even though it had never been
	// incremented.
	if conn.hasGenerationNumber() ***REMOVED***
		p.generation.removeConnection(conn.desc.ServiceID)
	***REMOVED***

	if p.monitor != nil ***REMOVED***
		p.monitor.Event(&event.PoolEvent***REMOVED***
			Type:         event.ConnectionClosed,
			Address:      p.address.String(),
			ConnectionID: conn.poolID,
			Reason:       reason,
		***REMOVED***)
	***REMOVED***

	return nil
***REMOVED***

// checkIn returns an idle connection to the pool. If the connection is perished or the pool is
// closed, it is removed from the connection pool and closed.
func (p *pool) checkIn(conn *connection) error ***REMOVED***
	if conn == nil ***REMOVED***
		return nil
	***REMOVED***
	if conn.pool != p ***REMOVED***
		return ErrWrongPool
	***REMOVED***

	if p.monitor != nil ***REMOVED***
		p.monitor.Event(&event.PoolEvent***REMOVED***
			Type:         event.ConnectionReturned,
			ConnectionID: conn.poolID,
			Address:      conn.addr.String(),
		***REMOVED***)
	***REMOVED***

	return p.checkInNoEvent(conn)
***REMOVED***

// checkInNoEvent returns a connection to the pool. It behaves identically to checkIn except it does
// not publish events. It is only intended for use by pool-internal functions.
func (p *pool) checkInNoEvent(conn *connection) error ***REMOVED***
	if conn == nil ***REMOVED***
		return nil
	***REMOVED***
	if conn.pool != p ***REMOVED***
		return ErrWrongPool
	***REMOVED***

	// Bump the connection idle deadline here because we're about to make the connection "available".
	// The idle deadline is used to determine when a connection has reached its max idle time and
	// should be closed. A connection reaches its max idle time when it has been "available" in the
	// idle connections stack for more than the configured duration (maxIdleTimeMS). Set it before
	// we call connectionPerished(), which checks the idle deadline, because a newly "available"
	// connection should never be perished due to max idle time.
	conn.bumpIdleDeadline()

	if reason, perished := connectionPerished(conn); perished ***REMOVED***
		_ = p.removeConnection(conn, reason)
		go func() ***REMOVED***
			_ = p.closeConnection(conn)
		***REMOVED***()
		return nil
	***REMOVED***

	if conn.pool.getState() == poolClosed ***REMOVED***
		_ = p.removeConnection(conn, event.ReasonPoolClosed)
		go func() ***REMOVED***
			_ = p.closeConnection(conn)
		***REMOVED***()
		return nil
	***REMOVED***

	p.idleMu.Lock()
	defer p.idleMu.Unlock()

	for ***REMOVED***
		w := p.idleConnWait.popFront()
		if w == nil ***REMOVED***
			break
		***REMOVED***
		if w.tryDeliver(conn, nil) ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***

	for _, idle := range p.idleConns ***REMOVED***
		if idle == conn ***REMOVED***
			return fmt.Errorf("duplicate idle conn %p in idle connections stack", conn)
		***REMOVED***
	***REMOVED***

	p.idleConns = append(p.idleConns, conn)
	return nil
***REMOVED***

// clear marks all connections as stale by incrementing the generation number, stops all background
// goroutines, removes all requests from idleConnWait and newConnWait, and sets the pool state to
// "paused". If serviceID is nil, clear marks all connections as stale. If serviceID is not nil,
// clear marks only connections associated with the given serviceID stale (for use in load balancer
// mode).
func (p *pool) clear(err error, serviceID *primitive.ObjectID) ***REMOVED***
	if p.getState() == poolClosed ***REMOVED***
		return
	***REMOVED***

	p.generation.clear(serviceID)

	// If serviceID is nil (i.e. not in load balancer mode), transition the pool to a paused state
	// by stopping all background goroutines, clearing the wait queues, and setting the pool state
	// to "paused".
	sendEvent := true
	if serviceID == nil ***REMOVED***
		// While holding the stateMu lock, set the pool state to "paused" if it's currently "ready",
		// and set lastClearErr to the error that caused the pool to be cleared. If the pool is
		// already paused, don't send another "ConnectionPoolCleared" event.
		p.stateMu.Lock()
		if p.state == poolPaused ***REMOVED***
			sendEvent = false
		***REMOVED***
		if p.state == poolReady ***REMOVED***
			p.state = poolPaused
		***REMOVED***
		p.lastClearErr = err
		p.stateMu.Unlock()

		pcErr := poolClearedError***REMOVED***err: err, address: p.address***REMOVED***

		// Clear the idle connections wait queue.
		p.idleMu.Lock()
		for ***REMOVED***
			w := p.idleConnWait.popFront()
			if w == nil ***REMOVED***
				break
			***REMOVED***
			w.tryDeliver(nil, pcErr)
		***REMOVED***
		p.idleMu.Unlock()

		// Clear the new connections wait queue. This effectively pauses the createConnections()
		// background goroutine because newConnWait is empty and checkOut() won't insert any more
		// wantConns into newConnWait until the pool is marked "ready" again.
		p.createConnectionsCond.L.Lock()
		for ***REMOVED***
			w := p.newConnWait.popFront()
			if w == nil ***REMOVED***
				break
			***REMOVED***
			w.tryDeliver(nil, pcErr)
		***REMOVED***
		p.createConnectionsCond.L.Unlock()
	***REMOVED***

	if sendEvent && p.monitor != nil ***REMOVED***
		p.monitor.Event(&event.PoolEvent***REMOVED***
			Type:      event.PoolCleared,
			Address:   p.address.String(),
			ServiceID: serviceID,
		***REMOVED***)
	***REMOVED***
***REMOVED***

// getOrQueueForIdleConn attempts to deliver an idle connection to the given wantConn. If there is
// an idle connection in the idle connections stack, it pops an idle connection, delivers it to the
// wantConn, and returns true. If there are no idle connections in the idle connections stack, it
// adds the wantConn to the idleConnWait queue and returns false.
func (p *pool) getOrQueueForIdleConn(w *wantConn) bool ***REMOVED***
	p.idleMu.Lock()
	defer p.idleMu.Unlock()

	// Try to deliver an idle connection from the idleConns stack first.
	for len(p.idleConns) > 0 ***REMOVED***
		conn := p.idleConns[len(p.idleConns)-1]
		p.idleConns = p.idleConns[:len(p.idleConns)-1]

		if conn == nil ***REMOVED***
			continue
		***REMOVED***

		if reason, perished := connectionPerished(conn); perished ***REMOVED***
			_ = conn.pool.removeConnection(conn, reason)
			go func() ***REMOVED***
				_ = conn.pool.closeConnection(conn)
			***REMOVED***()
			continue
		***REMOVED***

		if !w.tryDeliver(conn, nil) ***REMOVED***
			// If we couldn't deliver the conn to w, put it back in the idleConns stack.
			p.idleConns = append(p.idleConns, conn)
		***REMOVED***

		// If we got here, we tried to deliver an idle conn to w. No matter if tryDeliver() returned
		// true or false, w is no longer waiting and doesn't need to be added to any wait queues, so
		// return delivered = true.
		return true
	***REMOVED***

	p.idleConnWait.cleanFront()
	p.idleConnWait.pushBack(w)
	return false
***REMOVED***

func (p *pool) queueForNewConn(w *wantConn) ***REMOVED***
	p.createConnectionsCond.L.Lock()
	defer p.createConnectionsCond.L.Unlock()

	p.newConnWait.cleanFront()
	p.newConnWait.pushBack(w)
	p.createConnectionsCond.Signal()
***REMOVED***

func (p *pool) totalConnectionCount() int ***REMOVED***
	p.createConnectionsCond.L.Lock()
	defer p.createConnectionsCond.L.Unlock()

	return len(p.conns)
***REMOVED***

func (p *pool) availableConnectionCount() int ***REMOVED***
	p.idleMu.Lock()
	defer p.idleMu.Unlock()

	return len(p.idleConns)
***REMOVED***

// createConnections creates connections for wantConn requests on the newConnWait queue.
func (p *pool) createConnections(ctx context.Context, wg *sync.WaitGroup) ***REMOVED***
	defer wg.Done()

	// condition returns true if the createConnections() loop should continue and false if it should
	// wait. Note that the condition also listens for Context cancellation, which also causes the
	// loop to continue, allowing for a subsequent check to return from createConnections().
	condition := func() bool ***REMOVED***
		checkOutWaiting := p.newConnWait.len() > 0
		poolHasSpace := p.maxSize == 0 || uint64(len(p.conns)) < p.maxSize
		cancelled := ctx.Err() != nil
		return (checkOutWaiting && poolHasSpace) || cancelled
	***REMOVED***

	// wait waits for there to be an available wantConn and for the pool to have space for a new
	// connection. When the condition becomes true, it creates a new connection and returns the
	// waiting wantConn and new connection. If the Context is cancelled or there are any
	// errors, wait returns with "ok = false".
	wait := func() (*wantConn, *connection, bool) ***REMOVED***
		p.createConnectionsCond.L.Lock()
		defer p.createConnectionsCond.L.Unlock()

		for !condition() ***REMOVED***
			p.createConnectionsCond.Wait()
		***REMOVED***

		if ctx.Err() != nil ***REMOVED***
			return nil, nil, false
		***REMOVED***

		p.newConnWait.cleanFront()
		w := p.newConnWait.popFront()
		if w == nil ***REMOVED***
			return nil, nil, false
		***REMOVED***

		conn := newConnection(p.address, p.connOpts...)
		conn.pool = p
		conn.poolID = atomic.AddUint64(&p.nextID, 1)
		p.conns[conn.poolID] = conn

		return w, conn, true
	***REMOVED***

	for ctx.Err() == nil ***REMOVED***
		w, conn, ok := wait()
		if !ok ***REMOVED***
			continue
		***REMOVED***

		if p.monitor != nil ***REMOVED***
			p.monitor.Event(&event.PoolEvent***REMOVED***
				Type:         event.ConnectionCreated,
				Address:      p.address.String(),
				ConnectionID: conn.poolID,
			***REMOVED***)
		***REMOVED***

		// Pass the createConnections context to connect to allow pool close to cancel connection
		// establishment so shutdown doesn't block indefinitely if connectTimeout=0.
		err := conn.connect(ctx)
		if err != nil ***REMOVED***
			w.tryDeliver(nil, err)

			// If there's an error connecting the new connection, call the handshake error handler
			// that implements the SDAM handshake error handling logic. This must be called after
			// delivering the connection error to the waiting wantConn. If it's called before, the
			// handshake error handler may clear the connection pool, leading to a different error
			// message being delivered to the same waiting wantConn in idleConnWait when the wait
			// queues are cleared.
			if p.handshakeErrFn != nil ***REMOVED***
				p.handshakeErrFn(err, conn.generation, conn.desc.ServiceID)
			***REMOVED***

			_ = p.removeConnection(conn, event.ReasonError)
			_ = p.closeConnection(conn)
			continue
		***REMOVED***

		if p.monitor != nil ***REMOVED***
			p.monitor.Event(&event.PoolEvent***REMOVED***
				Type:         event.ConnectionReady,
				Address:      p.address.String(),
				ConnectionID: conn.poolID,
			***REMOVED***)
		***REMOVED***

		if w.tryDeliver(conn, nil) ***REMOVED***
			continue
		***REMOVED***

		_ = p.checkInNoEvent(conn)
	***REMOVED***
***REMOVED***

func (p *pool) maintain(ctx context.Context, wg *sync.WaitGroup) ***REMOVED***
	defer wg.Done()

	ticker := time.NewTicker(p.maintainInterval)
	defer ticker.Stop()

	// remove removes the *wantConn at index i from the slice and returns the new slice. The order
	// of the slice is not maintained.
	remove := func(arr []*wantConn, i int) []*wantConn ***REMOVED***
		end := len(arr) - 1
		arr[i], arr[end] = arr[end], arr[i]
		return arr[:end]
	***REMOVED***

	// removeNotWaiting removes any wantConns that are no longer waiting from given slice of
	// wantConns. That allows maintain() to use the size of its wantConns slice as an indication of
	// how many new connection requests are outstanding and subtract that from the number of
	// connections to ask for when maintaining minPoolSize.
	removeNotWaiting := func(arr []*wantConn) []*wantConn ***REMOVED***
		for i := len(arr) - 1; i >= 0; i-- ***REMOVED***
			w := arr[i]
			if !w.waiting() ***REMOVED***
				arr = remove(arr, i)
			***REMOVED***
		***REMOVED***

		return arr
	***REMOVED***

	wantConns := make([]*wantConn, 0, p.minSize)
	defer func() ***REMOVED***
		for _, w := range wantConns ***REMOVED***
			w.tryDeliver(nil, ErrPoolClosed)
		***REMOVED***
	***REMOVED***()

	for ***REMOVED***
		select ***REMOVED***
		case <-ticker.C:
		case <-p.maintainReady:
		case <-ctx.Done():
			return
		***REMOVED***

		// Only maintain the pool while it's in the "ready" state. If the pool state is not "ready",
		// wait for the next tick or "ready" signal. Do all of this while holding the stateMu read
		// lock to prevent a state change between checking the state and entering the wait queue.
		// Not holding the stateMu read lock here may allow maintain() to request wantConns after
		// clear() pauses the pool and clears the wait queue, resulting in createConnections()
		// doing work while the pool is "paused".
		p.stateMu.RLock()
		if p.state != poolReady ***REMOVED***
			p.stateMu.RUnlock()
			continue
		***REMOVED***

		p.removePerishedConns()

		// Remove any wantConns that are no longer waiting.
		wantConns = removeNotWaiting(wantConns)

		// Figure out how many more wantConns we need to satisfy minPoolSize. Assume that the
		// outstanding wantConns (i.e. the ones that weren't removed from the slice) will all return
		// connections when they're ready, so only add wantConns to make up the difference. Limit
		// the number of connections requested to max 10 at a time to prevent overshooting
		// minPoolSize in case other checkOut() calls are requesting new connections, too.
		total := p.totalConnectionCount()
		n := int(p.minSize) - total - len(wantConns)
		if n > 10 ***REMOVED***
			n = 10
		***REMOVED***

		for i := 0; i < n; i++ ***REMOVED***
			w := newWantConn()
			p.queueForNewConn(w)
			wantConns = append(wantConns, w)

			// Start a goroutine for each new wantConn, waiting for it to be ready.
			go func() ***REMOVED***
				<-w.ready
				if w.conn != nil ***REMOVED***
					_ = p.checkInNoEvent(w.conn)
				***REMOVED***
			***REMOVED***()
		***REMOVED***
		p.stateMu.RUnlock()
	***REMOVED***
***REMOVED***

func (p *pool) removePerishedConns() ***REMOVED***
	p.idleMu.Lock()
	defer p.idleMu.Unlock()

	for i := range p.idleConns ***REMOVED***
		conn := p.idleConns[i]
		if conn == nil ***REMOVED***
			continue
		***REMOVED***

		if reason, perished := connectionPerished(conn); perished ***REMOVED***
			p.idleConns[i] = nil

			_ = p.removeConnection(conn, reason)
			go func() ***REMOVED***
				_ = p.closeConnection(conn)
			***REMOVED***()
		***REMOVED***
	***REMOVED***

	p.idleConns = compact(p.idleConns)
***REMOVED***

// compact removes any nil pointers from the slice and keeps the non-nil pointers, retaining the
// order of the non-nil pointers.
func compact(arr []*connection) []*connection ***REMOVED***
	offset := 0
	for i := range arr ***REMOVED***
		if arr[i] == nil ***REMOVED***
			continue
		***REMOVED***
		arr[offset] = arr[i]
		offset++
	***REMOVED***
	return arr[:offset]
***REMOVED***

// A wantConn records state about a wanted connection (that is, an active call to checkOut).
// The conn may be gotten by creating a new connection or by finding an idle connection, or a
// cancellation may make the conn no longer wanted. These three options are racing against each
// other and use wantConn to coordinate and agree about the winning outcome.
// Based on https://cs.opensource.google/go/go/+/refs/tags/go1.16.6:src/net/http/transport.go;l=1174-1240
type wantConn struct ***REMOVED***
	ready chan struct***REMOVED******REMOVED***

	mu   sync.Mutex // Guards conn, err
	conn *connection
	err  error
***REMOVED***

func newWantConn() *wantConn ***REMOVED***
	return &wantConn***REMOVED***
		ready: make(chan struct***REMOVED******REMOVED***, 1),
	***REMOVED***
***REMOVED***

// waiting reports whether w is still waiting for an answer (connection or error).
func (w *wantConn) waiting() bool ***REMOVED***
	select ***REMOVED***
	case <-w.ready:
		return false
	default:
		return true
	***REMOVED***
***REMOVED***

// tryDeliver attempts to deliver conn, err to w and reports whether it succeeded.
func (w *wantConn) tryDeliver(conn *connection, err error) bool ***REMOVED***
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.conn != nil || w.err != nil ***REMOVED***
		return false
	***REMOVED***

	w.conn = conn
	w.err = err
	if w.conn == nil && w.err == nil ***REMOVED***
		panic("x/mongo/driver/topology: internal error: misuse of tryDeliver")
	***REMOVED***
	close(w.ready)
	return true
***REMOVED***

// cancel marks w as no longer wanting a result (for example, due to cancellation). If a connection
// has been delivered already, cancel returns it with p.checkInNoEvent(). Note that the caller must
// not hold any locks on the pool while calling cancel.
func (w *wantConn) cancel(p *pool, err error) ***REMOVED***
	if err == nil ***REMOVED***
		panic("x/mongo/driver/topology: internal error: misuse of cancel")
	***REMOVED***

	w.mu.Lock()
	if w.conn == nil && w.err == nil ***REMOVED***
		close(w.ready) // catch misbehavior in future delivery
	***REMOVED***
	conn := w.conn
	w.conn = nil
	w.err = err
	w.mu.Unlock()

	if conn != nil ***REMOVED***
		_ = p.checkInNoEvent(conn)
	***REMOVED***
***REMOVED***

// A wantConnQueue is a queue of wantConns.
// Based on https://cs.opensource.google/go/go/+/refs/tags/go1.16.6:src/net/http/transport.go;l=1242-1306
type wantConnQueue struct ***REMOVED***
	// This is a queue, not a deque.
	// It is split into two stages - head[headPos:] and tail.
	// popFront is trivial (headPos++) on the first stage, and
	// pushBack is trivial (append) on the second stage.
	// If the first stage is empty, popFront can swap the
	// first and second stages to remedy the situation.
	//
	// This two-stage split is analogous to the use of two lists
	// in Okasaki's purely functional queue but without the
	// overhead of reversing the list when swapping stages.
	head    []*wantConn
	headPos int
	tail    []*wantConn
***REMOVED***

// len returns the number of items in the queue.
func (q *wantConnQueue) len() int ***REMOVED***
	return len(q.head) - q.headPos + len(q.tail)
***REMOVED***

// pushBack adds w to the back of the queue.
func (q *wantConnQueue) pushBack(w *wantConn) ***REMOVED***
	q.tail = append(q.tail, w)
***REMOVED***

// popFront removes and returns the wantConn at the front of the queue.
func (q *wantConnQueue) popFront() *wantConn ***REMOVED***
	if q.headPos >= len(q.head) ***REMOVED***
		if len(q.tail) == 0 ***REMOVED***
			return nil
		***REMOVED***
		// Pick up tail as new head, clear tail.
		q.head, q.headPos, q.tail = q.tail, 0, q.head[:0]
	***REMOVED***
	w := q.head[q.headPos]
	q.head[q.headPos] = nil
	q.headPos++
	return w
***REMOVED***

// peekFront returns the wantConn at the front of the queue without removing it.
func (q *wantConnQueue) peekFront() *wantConn ***REMOVED***
	if q.headPos < len(q.head) ***REMOVED***
		return q.head[q.headPos]
	***REMOVED***
	if len(q.tail) > 0 ***REMOVED***
		return q.tail[0]
	***REMOVED***
	return nil
***REMOVED***

// cleanFront pops any wantConns that are no longer waiting from the head of the queue.
func (q *wantConnQueue) cleanFront() ***REMOVED***
	for ***REMOVED***
		w := q.peekFront()
		if w == nil || w.waiting() ***REMOVED***
			return
		***REMOVED***
		q.popFront()
	***REMOVED***
***REMOVED***
