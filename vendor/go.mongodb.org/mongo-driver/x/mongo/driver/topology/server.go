// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package topology

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/address"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
)

const minHeartbeatInterval = 500 * time.Millisecond

// Server state constants.
const (
	serverDisconnected int64 = iota
	serverDisconnecting
	serverConnected
)

func serverStateString(state int64) string ***REMOVED***
	switch state ***REMOVED***
	case serverDisconnected:
		return "Disconnected"
	case serverDisconnecting:
		return "Disconnecting"
	case serverConnected:
		return "Connected"
	***REMOVED***

	return ""
***REMOVED***

var (
	// ErrServerClosed occurs when an attempt to Get a connection is made after
	// the server has been closed.
	ErrServerClosed = errors.New("server is closed")
	// ErrServerConnected occurs when at attempt to Connect is made after a server
	// has already been connected.
	ErrServerConnected = errors.New("server is connected")

	errCheckCancelled = errors.New("server check cancelled")
	emptyDescription  = description.NewDefaultServer("")
)

// SelectedServer represents a specific server that was selected during server selection.
// It contains the kind of the topology it was selected from.
type SelectedServer struct ***REMOVED***
	*Server

	Kind description.TopologyKind
***REMOVED***

// Description returns a description of the server as of the last heartbeat.
func (ss *SelectedServer) Description() description.SelectedServer ***REMOVED***
	sdesc := ss.Server.Description()
	return description.SelectedServer***REMOVED***
		Server: sdesc,
		Kind:   ss.Kind,
	***REMOVED***
***REMOVED***

// Server is a single server within a topology.
type Server struct ***REMOVED***
	// The following integer fields must be accessed using the atomic package and should be at the
	// beginning of the struct.
	// - atomic bug: https://pkg.go.dev/sync/atomic#pkg-note-BUG
	// - suggested layout: https://go101.org/article/memory-layout.html

	state          int64
	operationCount int64

	cfg     *serverConfig
	address address.Address

	// connection related fields
	pool *pool

	// goroutine management fields
	done          chan struct***REMOVED******REMOVED***
	checkNow      chan struct***REMOVED******REMOVED***
	disconnecting chan struct***REMOVED******REMOVED***
	closewg       sync.WaitGroup

	// description related fields
	desc                   atomic.Value // holds a description.Server
	updateTopologyCallback atomic.Value
	topologyID             primitive.ObjectID

	// subscriber related fields
	subLock             sync.Mutex
	subscribers         map[uint64]chan description.Server
	currentSubscriberID uint64
	subscriptionsClosed bool

	// heartbeat and cancellation related fields
	// globalCtx should be created in NewServer and cancelled in Disconnect to signal that the server is shutting down.
	// heartbeatCtx should be used for individual heartbeats and should be a child of globalCtx so that it will be
	// cancelled automatically during shutdown.
	heartbeatLock      sync.Mutex
	conn               *connection
	globalCtx          context.Context
	globalCtxCancel    context.CancelFunc
	heartbeatCtx       context.Context
	heartbeatCtxCancel context.CancelFunc

	processErrorLock sync.Mutex
	rttMonitor       *rttMonitor
***REMOVED***

// updateTopologyCallback is a callback used to create a server that should be called when the parent Topology instance
// should be updated based on a new server description. The callback must return the server description that should be
// stored by the server.
type updateTopologyCallback func(description.Server) description.Server

// ConnectServer creates a new Server and then initializes it using the
// Connect method.
func ConnectServer(addr address.Address, updateCallback updateTopologyCallback, topologyID primitive.ObjectID, opts ...ServerOption) (*Server, error) ***REMOVED***
	srvr := NewServer(addr, topologyID, opts...)
	err := srvr.Connect(updateCallback)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return srvr, nil
***REMOVED***

// NewServer creates a new server. The mongodb server at the address will be monitored
// on an internal monitoring goroutine.
func NewServer(addr address.Address, topologyID primitive.ObjectID, opts ...ServerOption) *Server ***REMOVED***
	cfg := newServerConfig(opts...)
	globalCtx, globalCtxCancel := context.WithCancel(context.Background())
	s := &Server***REMOVED***
		state: serverDisconnected,

		cfg:     cfg,
		address: addr,

		done:          make(chan struct***REMOVED******REMOVED***),
		checkNow:      make(chan struct***REMOVED******REMOVED***, 1),
		disconnecting: make(chan struct***REMOVED******REMOVED***),

		topologyID: topologyID,

		subscribers:     make(map[uint64]chan description.Server),
		globalCtx:       globalCtx,
		globalCtxCancel: globalCtxCancel,
	***REMOVED***
	s.desc.Store(description.NewDefaultServer(addr))
	rttCfg := &rttConfig***REMOVED***
		interval:           cfg.heartbeatInterval,
		minRTTWindow:       5 * time.Minute,
		createConnectionFn: s.createConnection,
		createOperationFn:  s.createBaseOperation,
	***REMOVED***
	s.rttMonitor = newRTTMonitor(rttCfg)

	pc := poolConfig***REMOVED***
		Address:          addr,
		MinPoolSize:      cfg.minConns,
		MaxPoolSize:      cfg.maxConns,
		MaxConnecting:    cfg.maxConnecting,
		MaxIdleTime:      cfg.poolMaxIdleTime,
		MaintainInterval: cfg.poolMaintainInterval,
		PoolMonitor:      cfg.poolMonitor,
		handshakeErrFn:   s.ProcessHandshakeError,
	***REMOVED***

	connectionOpts := copyConnectionOpts(cfg.connectionOpts)
	s.pool = newPool(pc, connectionOpts...)
	s.publishServerOpeningEvent(s.address)

	return s
***REMOVED***

// Connect initializes the Server by starting background monitoring goroutines.
// This method must be called before a Server can be used.
func (s *Server) Connect(updateCallback updateTopologyCallback) error ***REMOVED***
	if !atomic.CompareAndSwapInt64(&s.state, serverDisconnected, serverConnected) ***REMOVED***
		return ErrServerConnected
	***REMOVED***

	desc := description.NewDefaultServer(s.address)
	if s.cfg.loadBalanced ***REMOVED***
		// LBs automatically start off with kind LoadBalancer because there is no monitoring routine for state changes.
		desc.Kind = description.LoadBalancer
	***REMOVED***
	s.desc.Store(desc)
	s.updateTopologyCallback.Store(updateCallback)

	if !s.cfg.monitoringDisabled && !s.cfg.loadBalanced ***REMOVED***
		s.rttMonitor.connect()
		s.closewg.Add(1)
		go s.update()
	***REMOVED***

	// The CMAP spec describes that pools should only be marked "ready" when the server description
	// is updated to something other than "Unknown". However, we maintain the previous Server
	// behavior here and immediately mark the pool as ready during Connect() to simplify and speed
	// up the Client startup behavior. The risk of marking a pool as ready proactively during
	// Connect() is that we could attempt to create connections to a server that was configured
	// erroneously until the first server check or checkOut() failure occurs, when the SDAM error
	// handler would transition the Server back to "Unknown" and set the pool to "paused".
	return s.pool.ready()
***REMOVED***

// Disconnect closes sockets to the server referenced by this Server.
// Subscriptions to this Server will be closed. Disconnect will shutdown
// any monitoring goroutines, closeConnection the idle connection pool, and will
// wait until all the in use connections have been returned to the connection
// pool and are closed before returning. If the context expires via
// cancellation, deadline, or timeout before the in use connections have been
// returned, the in use connections will be closed, resulting in the failure of
// any in flight read or write operations. If this method returns with no
// errors, all connections associated with this Server have been closed.
func (s *Server) Disconnect(ctx context.Context) error ***REMOVED***
	if !atomic.CompareAndSwapInt64(&s.state, serverConnected, serverDisconnecting) ***REMOVED***
		return ErrServerClosed
	***REMOVED***

	s.updateTopologyCallback.Store((updateTopologyCallback)(nil))

	// Cancel the global context so any new contexts created from it will be automatically cancelled. Close the done
	// channel so the update() routine will know that it can stop. Cancel any in-progress monitoring checks at the end.
	// The done channel is closed before cancelling the check so the update routine() will immediately detect that it
	// can stop rather than trying to create new connections until the read from done succeeds.
	s.globalCtxCancel()
	close(s.done)
	s.cancelCheck()

	s.rttMonitor.disconnect()
	s.pool.close(ctx)

	s.closewg.Wait()
	atomic.StoreInt64(&s.state, serverDisconnected)

	return nil
***REMOVED***

// Connection gets a connection to the server.
func (s *Server) Connection(ctx context.Context) (driver.Connection, error) ***REMOVED***
	if atomic.LoadInt64(&s.state) != serverConnected ***REMOVED***
		return nil, ErrServerClosed
	***REMOVED***

	// Increment the operation count before calling checkOut to make sure that all connection
	// requests are included in the operation count, including those in the wait queue. If we got an
	// error instead of a connection, immediately decrement the operation count.
	atomic.AddInt64(&s.operationCount, 1)
	conn, err := s.pool.checkOut(ctx)
	if err != nil ***REMOVED***
		atomic.AddInt64(&s.operationCount, -1)
		return nil, err
	***REMOVED***

	return &Connection***REMOVED***
		connection: conn,
		cleanupServerFn: func() ***REMOVED***
			// Decrement the operation count whenever the caller is done with the connection. Note
			// that cleanupServerFn() is not called while the connection is pinned to a cursor or
			// transaction, so the operation count is not decremented until the cursor is closed or
			// the transaction is committed or aborted. Use an int64 instead of a uint64 to mitigate
			// the impact of any possible bugs that could cause the uint64 to underflow, which would
			// make the server much less selectable.
			atomic.AddInt64(&s.operationCount, -1)
		***REMOVED***,
	***REMOVED***, nil
***REMOVED***

// ProcessHandshakeError implements SDAM error handling for errors that occur before a connection
// finishes handshaking.
func (s *Server) ProcessHandshakeError(err error, startingGenerationNumber uint64, serviceID *primitive.ObjectID) ***REMOVED***
	// Ignore the error if the server is behind a load balancer but the service ID is unknown. This indicates that the
	// error happened when dialing the connection or during the MongoDB handshake, so we don't know the service ID to
	// use for clearing the pool.
	if err == nil || s.cfg.loadBalanced && serviceID == nil ***REMOVED***
		return
	***REMOVED***
	// Ignore the error if the connection is stale.
	if startingGenerationNumber < s.pool.generation.getGeneration(serviceID) ***REMOVED***
		return
	***REMOVED***

	wrappedConnErr := unwrapConnectionError(err)
	if wrappedConnErr == nil ***REMOVED***
		return
	***REMOVED***

	// Must hold the processErrorLock while updating the server description and clearing the pool.
	// Not holding the lock leads to possible out-of-order processing of pool.clear() and
	// pool.ready() calls from concurrent server description updates.
	s.processErrorLock.Lock()
	defer s.processErrorLock.Unlock()

	// Since the only kind of ConnectionError we receive from pool.Get will be an initialization error, we should set
	// the description.Server appropriately. The description should not have a TopologyVersion because the staleness
	// checking logic above has already determined that this description is not stale.
	s.updateDescription(description.NewServerFromError(s.address, wrappedConnErr, nil))
	s.pool.clear(err, serviceID)
	s.cancelCheck()
***REMOVED***

// Description returns a description of the server as of the last heartbeat.
func (s *Server) Description() description.Server ***REMOVED***
	return s.desc.Load().(description.Server)
***REMOVED***

// SelectedDescription returns a description.SelectedServer with a Kind of
// Single. This can be used when performing tasks like monitoring a batch
// of servers and you want to run one off commands against those servers.
func (s *Server) SelectedDescription() description.SelectedServer ***REMOVED***
	sdesc := s.Description()
	return description.SelectedServer***REMOVED***
		Server: sdesc,
		Kind:   description.Single,
	***REMOVED***
***REMOVED***

// Subscribe returns a ServerSubscription which has a channel on which all
// updated server descriptions will be sent. The channel will have a buffer
// size of one, and will be pre-populated with the current description.
func (s *Server) Subscribe() (*ServerSubscription, error) ***REMOVED***
	if atomic.LoadInt64(&s.state) != serverConnected ***REMOVED***
		return nil, ErrSubscribeAfterClosed
	***REMOVED***
	ch := make(chan description.Server, 1)
	ch <- s.desc.Load().(description.Server)

	s.subLock.Lock()
	defer s.subLock.Unlock()
	if s.subscriptionsClosed ***REMOVED***
		return nil, ErrSubscribeAfterClosed
	***REMOVED***
	id := s.currentSubscriberID
	s.subscribers[id] = ch
	s.currentSubscriberID++

	ss := &ServerSubscription***REMOVED***
		C:  ch,
		s:  s,
		id: id,
	***REMOVED***

	return ss, nil
***REMOVED***

// RequestImmediateCheck will cause the server to send a heartbeat immediately
// instead of waiting for the heartbeat timeout.
func (s *Server) RequestImmediateCheck() ***REMOVED***
	select ***REMOVED***
	case s.checkNow <- struct***REMOVED******REMOVED******REMOVED******REMOVED***:
	default:
	***REMOVED***
***REMOVED***

// getWriteConcernErrorForProcessing extracts a driver.WriteConcernError from the provided error. This function returns
// (error, true) if the error is a WriteConcernError and the falls under the requirements for SDAM error
// handling and (nil, false) otherwise.
func getWriteConcernErrorForProcessing(err error) (*driver.WriteConcernError, bool) ***REMOVED***
	writeCmdErr, ok := err.(driver.WriteCommandError)
	if !ok ***REMOVED***
		return nil, false
	***REMOVED***

	wcerr := writeCmdErr.WriteConcernError
	if wcerr != nil && (wcerr.NodeIsRecovering() || wcerr.NotPrimary()) ***REMOVED***
		return wcerr, true
	***REMOVED***
	return nil, false
***REMOVED***

// ProcessError handles SDAM error handling and implements driver.ErrorProcessor.
func (s *Server) ProcessError(err error, conn driver.Connection) driver.ProcessErrorResult ***REMOVED***
	// ignore nil error
	if err == nil ***REMOVED***
		return driver.NoChange
	***REMOVED***

	// Must hold the processErrorLock while updating the server description and clearing the pool.
	// Not holding the lock leads to possible out-of-order processing of pool.clear() and
	// pool.ready() calls from concurrent server description updates.
	s.processErrorLock.Lock()
	defer s.processErrorLock.Unlock()

	// ignore stale error
	if conn.Stale() ***REMOVED***
		return driver.NoChange
	***REMOVED***
	// Invalidate server description if not primary or node recovering error occurs.
	// These errors can be reported as a command error or a write concern error.
	desc := conn.Description()
	if cerr, ok := err.(driver.Error); ok && (cerr.NodeIsRecovering() || cerr.NotPrimary()) ***REMOVED***
		// ignore stale error
		if desc.TopologyVersion.CompareToIncoming(cerr.TopologyVersion) >= 0 ***REMOVED***
			return driver.NoChange
		***REMOVED***

		// updates description to unknown
		s.updateDescription(description.NewServerFromError(s.address, err, cerr.TopologyVersion))
		s.RequestImmediateCheck()

		res := driver.ServerMarkedUnknown
		// If the node is shutting down or is older than 4.2, we synchronously clear the pool
		if cerr.NodeIsShuttingDown() || desc.WireVersion == nil || desc.WireVersion.Max < 8 ***REMOVED***
			res = driver.ConnectionPoolCleared
			s.pool.clear(err, desc.ServiceID)
		***REMOVED***

		return res
	***REMOVED***
	if wcerr, ok := getWriteConcernErrorForProcessing(err); ok ***REMOVED***
		// ignore stale error
		if desc.TopologyVersion.CompareToIncoming(wcerr.TopologyVersion) >= 0 ***REMOVED***
			return driver.NoChange
		***REMOVED***

		// updates description to unknown
		s.updateDescription(description.NewServerFromError(s.address, err, wcerr.TopologyVersion))
		s.RequestImmediateCheck()

		res := driver.ServerMarkedUnknown
		// If the node is shutting down or is older than 4.2, we synchronously clear the pool
		if wcerr.NodeIsShuttingDown() || desc.WireVersion == nil || desc.WireVersion.Max < 8 ***REMOVED***
			res = driver.ConnectionPoolCleared
			s.pool.clear(err, desc.ServiceID)
		***REMOVED***
		return res
	***REMOVED***

	wrappedConnErr := unwrapConnectionError(err)
	if wrappedConnErr == nil ***REMOVED***
		return driver.NoChange
	***REMOVED***

	// Ignore transient timeout errors.
	if netErr, ok := wrappedConnErr.(net.Error); ok && netErr.Timeout() ***REMOVED***
		return driver.NoChange
	***REMOVED***
	if wrappedConnErr == context.Canceled || wrappedConnErr == context.DeadlineExceeded ***REMOVED***
		return driver.NoChange
	***REMOVED***

	// For a non-timeout network error, we clear the pool, set the description to Unknown, and cancel the in-progress
	// monitoring check. The check is cancelled last to avoid a post-cancellation reconnect racing with
	// updateDescription.
	s.updateDescription(description.NewServerFromError(s.address, err, nil))
	s.pool.clear(err, desc.ServiceID)
	s.cancelCheck()
	return driver.ConnectionPoolCleared
***REMOVED***

// update handles performing heartbeats and updating any subscribers of the
// newest description.Server retrieved.
func (s *Server) update() ***REMOVED***
	defer s.closewg.Done()
	heartbeatTicker := time.NewTicker(s.cfg.heartbeatInterval)
	rateLimiter := time.NewTicker(minHeartbeatInterval)
	defer heartbeatTicker.Stop()
	defer rateLimiter.Stop()
	checkNow := s.checkNow
	done := s.done

	var doneOnce bool
	defer func() ***REMOVED***
		if r := recover(); r != nil ***REMOVED***
			if doneOnce ***REMOVED***
				return
			***REMOVED***
			// We keep this goroutine alive attempting to read from the done channel.
			<-done
		***REMOVED***
	***REMOVED***()

	closeServer := func() ***REMOVED***
		doneOnce = true
		s.subLock.Lock()
		for id, c := range s.subscribers ***REMOVED***
			close(c)
			delete(s.subscribers, id)
		***REMOVED***
		s.subscriptionsClosed = true
		s.subLock.Unlock()

		// We don't need to take s.heartbeatLock here because closeServer is called synchronously when the select checks
		// below detect that the server is being closed, so we can be sure that the connection isn't being used.
		if s.conn != nil ***REMOVED***
			_ = s.conn.close()
		***REMOVED***
	***REMOVED***

	waitUntilNextCheck := func() ***REMOVED***
		// Wait until heartbeatFrequency elapses, an application operation requests an immediate check, or the server
		// is disconnecting.
		select ***REMOVED***
		case <-heartbeatTicker.C:
		case <-checkNow:
		case <-done:
			// Return because the next update iteration will check the done channel again and clean up.
			return
		***REMOVED***

		// Ensure we only return if minHeartbeatFrequency has elapsed or the server is disconnecting.
		select ***REMOVED***
		case <-rateLimiter.C:
		case <-done:
			return
		***REMOVED***
	***REMOVED***

	for ***REMOVED***
		// Check if the server is disconnecting. Even if waitForNextCheck has already read from the done channel, we
		// can safely read from it again because Disconnect closes the channel.
		select ***REMOVED***
		case <-done:
			closeServer()
			return
		default:
		***REMOVED***

		previousDescription := s.Description()

		// Perform the next check.
		desc, err := s.check()
		if err == errCheckCancelled ***REMOVED***
			if atomic.LoadInt64(&s.state) != serverConnected ***REMOVED***
				continue
			***REMOVED***

			// If the server is not disconnecting, the check was cancelled by an application operation after an error.
			// Wait before running the next check.
			waitUntilNextCheck()
			continue
		***REMOVED***

		// Must hold the processErrorLock while updating the server description and clearing the
		// pool. Not holding the lock leads to possible out-of-order processing of pool.clear() and
		// pool.ready() calls from concurrent server description updates.
		s.processErrorLock.Lock()
		s.updateDescription(desc)
		if err := desc.LastError; err != nil ***REMOVED***
			// Clear the pool once the description has been updated to Unknown. Pass in a nil service ID to clear
			// because the monitoring routine only runs for non-load balanced deployments in which servers don't return
			// IDs.
			s.pool.clear(err, nil)
		***REMOVED***
		s.processErrorLock.Unlock()

		// If the server supports streaming or we're already streaming, we want to move to streaming the next response
		// without waiting. If the server has transitioned to Unknown from a network error, we want to do another
		// check without waiting in case it was a transient error and the server isn't actually down.
		serverSupportsStreaming := desc.Kind != description.Unknown && desc.TopologyVersion != nil
		connectionIsStreaming := s.conn != nil && s.conn.getCurrentlyStreaming()
		transitionedFromNetworkError := desc.LastError != nil && unwrapConnectionError(desc.LastError) != nil &&
			previousDescription.Kind != description.Unknown

		if serverSupportsStreaming || connectionIsStreaming || transitionedFromNetworkError ***REMOVED***
			continue
		***REMOVED***

		// The server either does not support the streamable protocol or is not in a healthy state, so we wait until
		// the next check.
		waitUntilNextCheck()
	***REMOVED***
***REMOVED***

// updateDescription handles updating the description on the Server, notifying
// subscribers, and potentially draining the connection pool. The initial
// parameter is used to determine if this is the first description from the
// server.
func (s *Server) updateDescription(desc description.Server) ***REMOVED***
	if s.cfg.loadBalanced ***REMOVED***
		// In load balanced mode, there are no updates from the monitoring routine. For errors encountered in pooled
		// connections, the server should not be marked Unknown to ensure that the LB remains selectable.
		return
	***REMOVED***

	defer func() ***REMOVED***
		//  ¯\_(ツ)_/¯
		_ = recover()
	***REMOVED***()

	// Anytime we update the server description to something other than "unknown", set the pool to
	// "ready". Do this before updating the description so that connections can be checked out as
	// soon as the server is selectable. If the pool is already ready, this operation is a no-op.
	// Note that this behavior is roughly consistent with the current Go driver behavior (connects
	// to all servers, even non-data-bearing nodes) but deviates slightly from CMAP spec, which
	// specifies a more restricted set of server descriptions and topologies that should mark the
	// pool ready. We don't have access to the topology here, so prefer the current Go driver
	// behavior for simplicity.
	if desc.Kind != description.Unknown ***REMOVED***
		_ = s.pool.ready()
	***REMOVED***

	// Use the updateTopologyCallback to update the parent Topology and get the description that should be stored.
	callback, ok := s.updateTopologyCallback.Load().(updateTopologyCallback)
	if ok && callback != nil ***REMOVED***
		desc = callback(desc)
	***REMOVED***
	s.desc.Store(desc)

	s.subLock.Lock()
	for _, c := range s.subscribers ***REMOVED***
		select ***REMOVED***
		// drain the channel if it isn't empty
		case <-c:
		default:
		***REMOVED***
		c <- desc
	***REMOVED***
	s.subLock.Unlock()
***REMOVED***

// createConnection creates a new connection instance but does not call connect on it. The caller must call connect
// before the connection can be used for network operations.
func (s *Server) createConnection() *connection ***REMOVED***
	opts := copyConnectionOpts(s.cfg.connectionOpts)
	opts = append(opts,
		WithConnectTimeout(func(time.Duration) time.Duration ***REMOVED*** return s.cfg.heartbeatTimeout ***REMOVED***),
		WithReadTimeout(func(time.Duration) time.Duration ***REMOVED*** return s.cfg.heartbeatTimeout ***REMOVED***),
		WithWriteTimeout(func(time.Duration) time.Duration ***REMOVED*** return s.cfg.heartbeatTimeout ***REMOVED***),
		// We override whatever handshaker is currently attached to the options with a basic
		// one because need to make sure we don't do auth.
		WithHandshaker(func(h Handshaker) Handshaker ***REMOVED***
			return operation.NewHello().AppName(s.cfg.appname).Compressors(s.cfg.compressionOpts).
				ServerAPI(s.cfg.serverAPI)
		***REMOVED***),
		// Override any monitors specified in options with nil to avoid monitoring heartbeats.
		WithMonitor(func(*event.CommandMonitor) *event.CommandMonitor ***REMOVED*** return nil ***REMOVED***),
	)

	return newConnection(s.address, opts...)
***REMOVED***

func copyConnectionOpts(opts []ConnectionOption) []ConnectionOption ***REMOVED***
	optsCopy := make([]ConnectionOption, len(opts))
	copy(optsCopy, opts)
	return optsCopy
***REMOVED***

func (s *Server) setupHeartbeatConnection() error ***REMOVED***
	conn := s.createConnection()

	// Take the lock when assigning the context and connection because they're accessed by cancelCheck.
	s.heartbeatLock.Lock()
	if s.heartbeatCtxCancel != nil ***REMOVED***
		// Ensure the previous context is cancelled to avoid a leak.
		s.heartbeatCtxCancel()
	***REMOVED***
	s.heartbeatCtx, s.heartbeatCtxCancel = context.WithCancel(s.globalCtx)
	s.conn = conn
	s.heartbeatLock.Unlock()

	return s.conn.connect(s.heartbeatCtx)
***REMOVED***

// cancelCheck cancels in-progress connection dials and reads. It does not set any fields on the server.
func (s *Server) cancelCheck() ***REMOVED***
	var conn *connection

	// Take heartbeatLock for mutual exclusion with the checks in the update function.
	s.heartbeatLock.Lock()
	if s.heartbeatCtx != nil ***REMOVED***
		s.heartbeatCtxCancel()
	***REMOVED***
	conn = s.conn
	s.heartbeatLock.Unlock()

	if conn == nil ***REMOVED***
		return
	***REMOVED***

	// If the connection exists, we need to wait for it to be connected because conn.connect() and
	// conn.close() cannot be called concurrently. If the connection wasn't successfully opened, its
	// state was set back to disconnected, so calling conn.close() will be a no-op.
	conn.closeConnectContext()
	conn.wait()
	_ = conn.close()
***REMOVED***

func (s *Server) checkWasCancelled() bool ***REMOVED***
	return s.heartbeatCtx.Err() != nil
***REMOVED***

func (s *Server) createBaseOperation(conn driver.Connection) *operation.Hello ***REMOVED***
	return operation.
		NewHello().
		ClusterClock(s.cfg.clock).
		Deployment(driver.SingleConnectionDeployment***REMOVED***conn***REMOVED***).
		ServerAPI(s.cfg.serverAPI)
***REMOVED***

func (s *Server) check() (description.Server, error) ***REMOVED***
	var descPtr *description.Server
	var err error
	var durationNanos int64

	// Create a new connection if this is the first check, the connection was closed after an error during the previous
	// check, or the previous check was cancelled.
	if s.conn == nil || s.conn.closed() || s.checkWasCancelled() ***REMOVED***
		// Create a new connection and add it's handshake RTT as a sample.
		err = s.setupHeartbeatConnection()
		if err == nil ***REMOVED***
			// Use the description from the connection handshake as the value for this check.
			s.rttMonitor.addSample(s.conn.helloRTT)
			descPtr = &s.conn.desc
		***REMOVED***
	***REMOVED***

	if descPtr == nil && err == nil ***REMOVED***
		// An existing connection is being used. Use the server description properties to execute the right heartbeat.

		// Wrap conn in a type that implements driver.StreamerConnection.
		heartbeatConn := initConnection***REMOVED***s.conn***REMOVED***
		baseOperation := s.createBaseOperation(heartbeatConn)
		previousDescription := s.Description()
		streamable := previousDescription.TopologyVersion != nil

		s.publishServerHeartbeatStartedEvent(s.conn.ID(), s.conn.getCurrentlyStreaming() || streamable)
		start := time.Now()
		switch ***REMOVED***
		case s.conn.getCurrentlyStreaming():
			// The connection is already in a streaming state, so we stream the next response.
			err = baseOperation.StreamResponse(s.heartbeatCtx, heartbeatConn)
		case streamable:
			// The server supports the streamable protocol. Set the socket timeout to
			// connectTimeoutMS+heartbeatFrequencyMS and execute an awaitable hello request. Set conn.canStream so
			// the wire message will advertise streaming support to the server.

			// Calculation for maxAwaitTimeMS is taken from time.Duration.Milliseconds (added in Go 1.13).
			maxAwaitTimeMS := int64(s.cfg.heartbeatInterval) / 1e6
			// If connectTimeoutMS=0, the socket timeout should be infinite. Otherwise, it is connectTimeoutMS +
			// heartbeatFrequencyMS to account for the fact that the query will block for heartbeatFrequencyMS
			// server-side.
			socketTimeout := s.cfg.heartbeatTimeout
			if socketTimeout != 0 ***REMOVED***
				socketTimeout += s.cfg.heartbeatInterval
			***REMOVED***
			s.conn.setSocketTimeout(socketTimeout)
			baseOperation = baseOperation.TopologyVersion(previousDescription.TopologyVersion).
				MaxAwaitTimeMS(maxAwaitTimeMS)
			s.conn.setCanStream(true)
			err = baseOperation.Execute(s.heartbeatCtx)
		default:
			// The server doesn't support the awaitable protocol. Set the socket timeout to connectTimeoutMS and
			// execute a regular heartbeat without any additional parameters.

			s.conn.setSocketTimeout(s.cfg.heartbeatTimeout)
			err = baseOperation.Execute(s.heartbeatCtx)
		***REMOVED***
		durationNanos = time.Since(start).Nanoseconds()

		if err == nil ***REMOVED***
			tempDesc := baseOperation.Result(s.address)
			descPtr = &tempDesc
			s.publishServerHeartbeatSucceededEvent(s.conn.ID(), durationNanos, tempDesc, s.conn.getCurrentlyStreaming() || streamable)
		***REMOVED*** else ***REMOVED***
			// Close the connection here rather than below so we ensure we're not closing a connection that wasn't
			// successfully created.
			if s.conn != nil ***REMOVED***
				_ = s.conn.close()
			***REMOVED***
			s.publishServerHeartbeatFailedEvent(s.conn.ID(), durationNanos, err, s.conn.getCurrentlyStreaming() || streamable)
		***REMOVED***
	***REMOVED***

	if descPtr != nil ***REMOVED***
		// The check was successful. Set the average RTT and the 90th percentile RTT and return.
		desc := *descPtr
		desc = desc.SetAverageRTT(s.rttMonitor.getRTT())
		desc.HeartbeatInterval = s.cfg.heartbeatInterval
		return desc, nil
	***REMOVED***

	if s.checkWasCancelled() ***REMOVED***
		// If the previous check was cancelled, we don't want to clear the pool. Return a sentinel error so the caller
		// will know that an actual error didn't occur.
		return emptyDescription, errCheckCancelled
	***REMOVED***

	// An error occurred. We reset the RTT monitor for all errors and return an Unknown description. The pool must also
	// be cleared, but only after the description has already been updated, so that is handled by the caller.
	topologyVersion := extractTopologyVersion(err)
	s.rttMonitor.reset()
	return description.NewServerFromError(s.address, err, topologyVersion), nil
***REMOVED***

func extractTopologyVersion(err error) *description.TopologyVersion ***REMOVED***
	if ce, ok := err.(ConnectionError); ok ***REMOVED***
		err = ce.Wrapped
	***REMOVED***

	switch converted := err.(type) ***REMOVED***
	case driver.Error:
		return converted.TopologyVersion
	case driver.WriteCommandError:
		if converted.WriteConcernError != nil ***REMOVED***
			return converted.WriteConcernError.TopologyVersion
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// MinRTT returns the minimum round-trip time to the server observed over the last 5 minutes.
func (s *Server) MinRTT() time.Duration ***REMOVED***
	return s.rttMonitor.getMinRTT()
***REMOVED***

// RTT90 returns the 90th percentile round-trip time to the server observed over the last 5 minutes.
func (s *Server) RTT90() time.Duration ***REMOVED***
	return s.rttMonitor.getRTT90()
***REMOVED***

// OperationCount returns the current number of in-progress operations for this server.
func (s *Server) OperationCount() int64 ***REMOVED***
	return atomic.LoadInt64(&s.operationCount)
***REMOVED***

// String implements the Stringer interface.
func (s *Server) String() string ***REMOVED***
	desc := s.Description()
	state := atomic.LoadInt64(&s.state)
	str := fmt.Sprintf("Addr: %s, Type: %s, State: %s",
		s.address, desc.Kind, serverStateString(state))
	if len(desc.Tags) != 0 ***REMOVED***
		str += fmt.Sprintf(", Tag sets: %s", desc.Tags)
	***REMOVED***
	if state == serverConnected ***REMOVED***
		str += fmt.Sprintf(", Average RTT: %s, Min RTT: %s", desc.AverageRTT, s.MinRTT())
	***REMOVED***
	if desc.LastError != nil ***REMOVED***
		str += fmt.Sprintf(", Last error: %s", desc.LastError)
	***REMOVED***

	return str
***REMOVED***

// ServerSubscription represents a subscription to the description.Server updates for
// a specific server.
type ServerSubscription struct ***REMOVED***
	C  <-chan description.Server
	s  *Server
	id uint64
***REMOVED***

// Unsubscribe unsubscribes this ServerSubscription from updates and closes the
// subscription channel.
func (ss *ServerSubscription) Unsubscribe() error ***REMOVED***
	ss.s.subLock.Lock()
	defer ss.s.subLock.Unlock()
	if ss.s.subscriptionsClosed ***REMOVED***
		return nil
	***REMOVED***

	ch, ok := ss.s.subscribers[ss.id]
	if !ok ***REMOVED***
		return nil
	***REMOVED***

	close(ch)
	delete(ss.s.subscribers, ss.id)

	return nil
***REMOVED***

// publishes a ServerOpeningEvent to indicate the server is being initialized
func (s *Server) publishServerOpeningEvent(addr address.Address) ***REMOVED***
	if s == nil ***REMOVED***
		return
	***REMOVED***

	serverOpening := &event.ServerOpeningEvent***REMOVED***
		Address:    addr,
		TopologyID: s.topologyID,
	***REMOVED***

	if s.cfg.serverMonitor != nil && s.cfg.serverMonitor.ServerOpening != nil ***REMOVED***
		s.cfg.serverMonitor.ServerOpening(serverOpening)
	***REMOVED***
***REMOVED***

// publishes a ServerHeartbeatStartedEvent to indicate a hello command has started
func (s *Server) publishServerHeartbeatStartedEvent(connectionID string, await bool) ***REMOVED***
	serverHeartbeatStarted := &event.ServerHeartbeatStartedEvent***REMOVED***
		ConnectionID: connectionID,
		Awaited:      await,
	***REMOVED***

	if s != nil && s.cfg.serverMonitor != nil && s.cfg.serverMonitor.ServerHeartbeatStarted != nil ***REMOVED***
		s.cfg.serverMonitor.ServerHeartbeatStarted(serverHeartbeatStarted)
	***REMOVED***
***REMOVED***

// publishes a ServerHeartbeatSucceededEvent to indicate hello has succeeded
func (s *Server) publishServerHeartbeatSucceededEvent(connectionID string,
	durationNanos int64,
	desc description.Server,
	await bool) ***REMOVED***
	serverHeartbeatSucceeded := &event.ServerHeartbeatSucceededEvent***REMOVED***
		DurationNanos: durationNanos,
		Reply:         desc,
		ConnectionID:  connectionID,
		Awaited:       await,
	***REMOVED***

	if s != nil && s.cfg.serverMonitor != nil && s.cfg.serverMonitor.ServerHeartbeatSucceeded != nil ***REMOVED***
		s.cfg.serverMonitor.ServerHeartbeatSucceeded(serverHeartbeatSucceeded)
	***REMOVED***
***REMOVED***

// publishes a ServerHeartbeatFailedEvent to indicate hello has failed
func (s *Server) publishServerHeartbeatFailedEvent(connectionID string,
	durationNanos int64,
	err error,
	await bool) ***REMOVED***
	serverHeartbeatFailed := &event.ServerHeartbeatFailedEvent***REMOVED***
		DurationNanos: durationNanos,
		Failure:       err,
		ConnectionID:  connectionID,
		Awaited:       await,
	***REMOVED***

	if s != nil && s.cfg.serverMonitor != nil && s.cfg.serverMonitor.ServerHeartbeatFailed != nil ***REMOVED***
		s.cfg.serverMonitor.ServerHeartbeatFailed(serverHeartbeatFailed)
	***REMOVED***
***REMOVED***

// unwrapConnectionError returns the connection error wrapped by err, or nil if err does not wrap a connection error.
func unwrapConnectionError(err error) error ***REMOVED***
	// This is essentially an implementation of errors.As to unwrap this error until we get a ConnectionError and then
	// return ConnectionError.Wrapped.

	connErr, ok := err.(ConnectionError)
	if ok ***REMOVED***
		return connErr.Wrapped
	***REMOVED***

	driverErr, ok := err.(driver.Error)
	if !ok || !driverErr.NetworkError() ***REMOVED***
		return nil
	***REMOVED***

	connErr, ok = driverErr.Wrapped.(ConnectionError)
	if ok ***REMOVED***
		return connErr.Wrapped
	***REMOVED***

	return nil
***REMOVED***
