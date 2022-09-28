// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package topology

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go.mongodb.org/mongo-driver/internal"
	"go.mongodb.org/mongo-driver/mongo/address"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/ocsp"
	"go.mongodb.org/mongo-driver/x/mongo/driver/wiremessage"
)

// Connection state constants.
const (
	connDisconnected int64 = iota
	connConnected
	connInitialized
)

var globalConnectionID uint64 = 1

var (
	defaultMaxMessageSize        uint32 = 48000000
	errResponseTooLarge                 = errors.New("length of read message too large")
	errLoadBalancedStateMismatch        = errors.New("driver attempted to initialize in load balancing mode, but the server does not support this mode")
)

func nextConnectionID() uint64 ***REMOVED*** return atomic.AddUint64(&globalConnectionID, 1) ***REMOVED***

type connection struct ***REMOVED***
	// state must be accessed using the atomic package and should be at the beginning of the struct.
	// - atomic bug: https://pkg.go.dev/sync/atomic#pkg-note-BUG
	// - suggested layout: https://go101.org/article/memory-layout.html
	state int64

	id                   string
	nc                   net.Conn // When nil, the connection is closed.
	addr                 address.Address
	idleTimeout          time.Duration
	idleDeadline         atomic.Value // Stores a time.Time
	readTimeout          time.Duration
	writeTimeout         time.Duration
	desc                 description.Server
	helloRTT             time.Duration
	compressor           wiremessage.CompressorID
	zliblevel            int
	zstdLevel            int
	connectDone          chan struct***REMOVED******REMOVED***
	config               *connectionConfig
	cancelConnectContext context.CancelFunc
	connectContextMade   chan struct***REMOVED******REMOVED***
	canStream            bool
	currentlyStreaming   bool
	connectContextMutex  sync.Mutex
	cancellationListener cancellationListener
	serverConnectionID   *int32 // the server's ID for this client's connection

	// pool related fields
	pool       *pool
	poolID     uint64
	generation uint64
***REMOVED***

// newConnection handles the creation of a connection. It does not connect the connection.
func newConnection(addr address.Address, opts ...ConnectionOption) *connection ***REMOVED***
	cfg := newConnectionConfig(opts...)

	id := fmt.Sprintf("%s[-%d]", addr, nextConnectionID())

	c := &connection***REMOVED***
		id:                   id,
		addr:                 addr,
		idleTimeout:          cfg.idleTimeout,
		readTimeout:          cfg.readTimeout,
		writeTimeout:         cfg.writeTimeout,
		connectDone:          make(chan struct***REMOVED******REMOVED***),
		config:               cfg,
		connectContextMade:   make(chan struct***REMOVED******REMOVED***),
		cancellationListener: internal.NewCancellationListener(),
	***REMOVED***
	// Connections to non-load balanced deployments should eagerly set the generation numbers so errors encountered
	// at any point during connection establishment can be processed without the connection being considered stale.
	if !c.config.loadBalanced ***REMOVED***
		c.setGenerationNumber()
	***REMOVED***
	atomic.StoreInt64(&c.state, connInitialized)

	return c
***REMOVED***

// setGenerationNumber sets the connection's generation number if a callback has been provided to do so in connection
// configuration.
func (c *connection) setGenerationNumber() ***REMOVED***
	if c.config.getGenerationFn != nil ***REMOVED***
		c.generation = c.config.getGenerationFn(c.desc.ServiceID)
	***REMOVED***
***REMOVED***

// hasGenerationNumber returns true if the connection has set its generation number. If so, this indicates that the
// generationNumberFn provided via the connection options has been called exactly once.
func (c *connection) hasGenerationNumber() bool ***REMOVED***
	if !c.config.loadBalanced ***REMOVED***
		// The generation is known for all non-LB clusters once the connection object has been created.
		return true
	***REMOVED***

	// For LB clusters, we set the generation after the initial handshake, so we know it's set if the connection
	// description has been updated to reflect that it's behind an LB.
	return c.desc.LoadBalanced()
***REMOVED***

// connect handles the I/O for a connection. It will dial, configure TLS, and perform initialization
// handshakes. All errors returned by connect are considered "before the handshake completes" and
// must be handled by calling the appropriate SDAM handshake error handler.
func (c *connection) connect(ctx context.Context) (err error) ***REMOVED***
	if !atomic.CompareAndSwapInt64(&c.state, connInitialized, connConnected) ***REMOVED***
		return nil
	***REMOVED***

	defer close(c.connectDone)

	// If connect returns an error, set the connection status as disconnected and close the
	// underlying net.Conn if it was created.
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			atomic.StoreInt64(&c.state, connDisconnected)

			if c.nc != nil ***REMOVED***
				_ = c.nc.Close()
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	// Create separate contexts for dialing a connection and doing the MongoDB/auth handshakes.
	//
	// handshakeCtx is simply a cancellable version of ctx because there's no default timeout that needs to be applied
	// to the full handshake. The cancellation allows consumers to bail out early when dialing a connection if it's no
	// longer required. This is done in lock because it accesses the shared cancelConnectContext field.
	//
	// dialCtx is equal to handshakeCtx if connectTimeoutMS=0. Otherwise, it is derived from handshakeCtx so the
	// cancellation still applies but with an added timeout to ensure the connectTimeoutMS option is applied to socket
	// establishment and the TLS handshake as a whole. This is created outside of the connectContextMutex lock to avoid
	// holding the lock longer than necessary.
	c.connectContextMutex.Lock()
	var handshakeCtx context.Context
	handshakeCtx, c.cancelConnectContext = context.WithCancel(ctx)
	c.connectContextMutex.Unlock()

	dialCtx := handshakeCtx
	var dialCancel context.CancelFunc
	if c.config.connectTimeout != 0 ***REMOVED***
		dialCtx, dialCancel = context.WithTimeout(handshakeCtx, c.config.connectTimeout)
		defer dialCancel()
	***REMOVED***

	defer func() ***REMOVED***
		var cancelFn context.CancelFunc

		c.connectContextMutex.Lock()
		cancelFn = c.cancelConnectContext
		c.cancelConnectContext = nil
		c.connectContextMutex.Unlock()

		if cancelFn != nil ***REMOVED***
			cancelFn()
		***REMOVED***
	***REMOVED***()

	close(c.connectContextMade)

	// Assign the result of DialContext to a temporary net.Conn to ensure that c.nc is not set in an error case.
	tempNc, err := c.config.dialer.DialContext(dialCtx, c.addr.Network(), c.addr.String())
	if err != nil ***REMOVED***
		return ConnectionError***REMOVED***Wrapped: err, init: true***REMOVED***
	***REMOVED***
	c.nc = tempNc

	if c.config.tlsConfig != nil ***REMOVED***
		tlsConfig := c.config.tlsConfig.Clone()

		// store the result of configureTLS in a separate variable than c.nc to avoid overwriting c.nc with nil in
		// error cases.
		ocspOpts := &ocsp.VerifyOptions***REMOVED***
			Cache:                   c.config.ocspCache,
			DisableEndpointChecking: c.config.disableOCSPEndpointCheck,
		***REMOVED***
		tlsNc, err := configureTLS(dialCtx, c.config.tlsConnectionSource, c.nc, c.addr, tlsConfig, ocspOpts)
		if err != nil ***REMOVED***
			return ConnectionError***REMOVED***Wrapped: err, init: true***REMOVED***
		***REMOVED***
		c.nc = tlsNc
	***REMOVED***

	// running hello and authentication is handled by a handshaker on the configuration instance.
	handshaker := c.config.handshaker
	if handshaker == nil ***REMOVED***
		return nil
	***REMOVED***

	var handshakeInfo driver.HandshakeInformation
	handshakeStartTime := time.Now()
	handshakeConn := initConnection***REMOVED***c***REMOVED***
	handshakeInfo, err = handshaker.GetHandshakeInformation(handshakeCtx, c.addr, handshakeConn)
	if err == nil ***REMOVED***
		// We only need to retain the Description field as the connection's description. The authentication-related
		// fields in handshakeInfo are tracked by the handshaker if necessary.
		c.desc = handshakeInfo.Description
		c.serverConnectionID = handshakeInfo.ServerConnectionID
		c.helloRTT = time.Since(handshakeStartTime)

		// If the application has indicated that the cluster is load balanced, ensure the server has included serviceId
		// in its handshake response to signal that it knows it's behind an LB as well.
		if c.config.loadBalanced && c.desc.ServiceID == nil ***REMOVED***
			err = errLoadBalancedStateMismatch
		***REMOVED***
	***REMOVED***
	if err == nil ***REMOVED***
		// For load-balanced connections, the generation number depends on the service ID, which isn't known until the
		// initial MongoDB handshake is done. To account for this, we don't attempt to set the connection's generation
		// number unless GetHandshakeInformation succeeds.
		if c.config.loadBalanced ***REMOVED***
			c.setGenerationNumber()
		***REMOVED***

		// If we successfully finished the first part of the handshake and verified LB state, continue with the rest of
		// the handshake.
		err = handshaker.FinishHandshake(handshakeCtx, handshakeConn)
	***REMOVED***

	// We have a failed handshake here
	if err != nil ***REMOVED***
		return ConnectionError***REMOVED***Wrapped: err, init: true***REMOVED***
	***REMOVED***

	if len(c.desc.Compression) > 0 ***REMOVED***
	clientMethodLoop:
		for _, method := range c.config.compressors ***REMOVED***
			for _, serverMethod := range c.desc.Compression ***REMOVED***
				if method != serverMethod ***REMOVED***
					continue
				***REMOVED***

				switch strings.ToLower(method) ***REMOVED***
				case "snappy":
					c.compressor = wiremessage.CompressorSnappy
				case "zlib":
					c.compressor = wiremessage.CompressorZLib
					c.zliblevel = wiremessage.DefaultZlibLevel
					if c.config.zlibLevel != nil ***REMOVED***
						c.zliblevel = *c.config.zlibLevel
					***REMOVED***
				case "zstd":
					c.compressor = wiremessage.CompressorZstd
					c.zstdLevel = wiremessage.DefaultZstdLevel
					if c.config.zstdLevel != nil ***REMOVED***
						c.zstdLevel = *c.config.zstdLevel
					***REMOVED***
				***REMOVED***
				break clientMethodLoop
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (c *connection) wait() ***REMOVED***
	if c.connectDone != nil ***REMOVED***
		<-c.connectDone
	***REMOVED***
***REMOVED***

func (c *connection) closeConnectContext() ***REMOVED***
	<-c.connectContextMade
	var cancelFn context.CancelFunc

	c.connectContextMutex.Lock()
	cancelFn = c.cancelConnectContext
	c.cancelConnectContext = nil
	c.connectContextMutex.Unlock()

	if cancelFn != nil ***REMOVED***
		cancelFn()
	***REMOVED***
***REMOVED***

func transformNetworkError(ctx context.Context, originalError error, contextDeadlineUsed bool) error ***REMOVED***
	if originalError == nil ***REMOVED***
		return nil
	***REMOVED***

	// If there was an error and the context was cancelled, we assume it happened due to the cancellation.
	if ctx.Err() == context.Canceled ***REMOVED***
		return context.Canceled
	***REMOVED***

	// If there was a timeout error and the context deadline was used, we convert the error into
	// context.DeadlineExceeded.
	if !contextDeadlineUsed ***REMOVED***
		return originalError
	***REMOVED***
	if netErr, ok := originalError.(net.Error); ok && netErr.Timeout() ***REMOVED***
		return context.DeadlineExceeded
	***REMOVED***

	return originalError
***REMOVED***

func (c *connection) cancellationListenerCallback() ***REMOVED***
	_ = c.close()
***REMOVED***

func (c *connection) writeWireMessage(ctx context.Context, wm []byte) error ***REMOVED***
	var err error
	if atomic.LoadInt64(&c.state) != connConnected ***REMOVED***
		return ConnectionError***REMOVED***ConnectionID: c.id, message: "connection is closed"***REMOVED***
	***REMOVED***

	var deadline time.Time
	if c.writeTimeout != 0 ***REMOVED***
		deadline = time.Now().Add(c.writeTimeout)
	***REMOVED***

	var contextDeadlineUsed bool
	if dl, ok := ctx.Deadline(); ok && (deadline.IsZero() || dl.Before(deadline)) ***REMOVED***
		contextDeadlineUsed = true
		deadline = dl
	***REMOVED***

	if err := c.nc.SetWriteDeadline(deadline); err != nil ***REMOVED***
		return ConnectionError***REMOVED***ConnectionID: c.id, Wrapped: err, message: "failed to set write deadline"***REMOVED***
	***REMOVED***

	err = c.write(ctx, wm)
	if err != nil ***REMOVED***
		c.close()
		return ConnectionError***REMOVED***
			ConnectionID: c.id,
			Wrapped:      transformNetworkError(ctx, err, contextDeadlineUsed),
			message:      "unable to write wire message to network",
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (c *connection) write(ctx context.Context, wm []byte) (err error) ***REMOVED***
	go c.cancellationListener.Listen(ctx, c.cancellationListenerCallback)
	defer func() ***REMOVED***
		// There is a race condition between Write and StopListening. If the context is cancelled after c.nc.Write
		// succeeds, the cancellation listener could fire and close the connection. In this case, the connection has
		// been invalidated but the error is nil. To account for this, overwrite the error to context.Cancelled if
		// the abortedForCancellation flag was set.

		if aborted := c.cancellationListener.StopListening(); aborted && err == nil ***REMOVED***
			err = context.Canceled
		***REMOVED***
	***REMOVED***()

	_, err = c.nc.Write(wm)
	return err
***REMOVED***

// readWireMessage reads a wiremessage from the connection. The dst parameter will be overwritten.
func (c *connection) readWireMessage(ctx context.Context, dst []byte) ([]byte, error) ***REMOVED***
	if atomic.LoadInt64(&c.state) != connConnected ***REMOVED***
		return dst, ConnectionError***REMOVED***ConnectionID: c.id, message: "connection is closed"***REMOVED***
	***REMOVED***

	var deadline time.Time
	if c.readTimeout != 0 ***REMOVED***
		deadline = time.Now().Add(c.readTimeout)
	***REMOVED***

	var contextDeadlineUsed bool
	if dl, ok := ctx.Deadline(); ok && (deadline.IsZero() || dl.Before(deadline)) ***REMOVED***
		contextDeadlineUsed = true
		deadline = dl
	***REMOVED***

	if err := c.nc.SetReadDeadline(deadline); err != nil ***REMOVED***
		return nil, ConnectionError***REMOVED***ConnectionID: c.id, Wrapped: err, message: "failed to set read deadline"***REMOVED***
	***REMOVED***

	dst, errMsg, err := c.read(ctx, dst)
	if err != nil ***REMOVED***
		// We closeConnection the connection because we don't know if there are other bytes left to read.
		c.close()
		message := errMsg
		if err == io.EOF ***REMOVED***
			message = "socket was unexpectedly closed"
		***REMOVED***
		return nil, ConnectionError***REMOVED***
			ConnectionID: c.id,
			Wrapped:      transformNetworkError(ctx, err, contextDeadlineUsed),
			message:      message,
		***REMOVED***
	***REMOVED***

	return dst, nil
***REMOVED***

func (c *connection) read(ctx context.Context, dst []byte) (bytesRead []byte, errMsg string, err error) ***REMOVED***
	go c.cancellationListener.Listen(ctx, c.cancellationListenerCallback)
	defer func() ***REMOVED***
		// If the context is cancelled after we finish reading the server response, the cancellation listener could fire
		// even though the socket reads succeed. To account for this, we overwrite err to be context.Canceled if the
		// abortedForCancellation flag is set.

		if aborted := c.cancellationListener.StopListening(); aborted && err == nil ***REMOVED***
			errMsg = "unable to read server response"
			err = context.Canceled
		***REMOVED***
	***REMOVED***()

	// We use an array here because it only costs 4 bytes on the stack and means we'll only need to
	// reslice dst once instead of twice.
	var sizeBuf [4]byte

	// We do a ReadFull into an array here instead of doing an opportunistic ReadAtLeast into dst
	// because there might be more than one wire message waiting to be read, for example when
	// reading messages from an exhaust cursor.
	_, err = io.ReadFull(c.nc, sizeBuf[:])
	if err != nil ***REMOVED***
		return nil, "incomplete read of message header", err
	***REMOVED***

	// read the length as an int32
	size := (int32(sizeBuf[0])) | (int32(sizeBuf[1]) << 8) | (int32(sizeBuf[2]) << 16) | (int32(sizeBuf[3]) << 24)

	// In the case of a hello response where MaxMessageSize has not yet been set, use the hard-coded
	// defaultMaxMessageSize instead.
	maxMessageSize := c.desc.MaxMessageSize
	if maxMessageSize == 0 ***REMOVED***
		maxMessageSize = defaultMaxMessageSize
	***REMOVED***
	if uint32(size) > maxMessageSize ***REMOVED***
		return nil, errResponseTooLarge.Error(), errResponseTooLarge
	***REMOVED***

	if int(size) > cap(dst) ***REMOVED***
		// Since we can't grow this slice without allocating, just allocate an entirely new slice.
		dst = make([]byte, 0, size)
	***REMOVED***
	// We need to ensure we don't accidentally read into a subsequent wire message, so we set the
	// size to read exactly this wire message.
	dst = dst[:size]
	copy(dst, sizeBuf[:])

	_, err = io.ReadFull(c.nc, dst[4:])
	if err != nil ***REMOVED***
		return nil, "incomplete read of full message", err
	***REMOVED***

	return dst, "", nil
***REMOVED***

func (c *connection) close() error ***REMOVED***
	// Overwrite the connection state as the first step so only the first close call will execute.
	if !atomic.CompareAndSwapInt64(&c.state, connConnected, connDisconnected) ***REMOVED***
		return nil
	***REMOVED***

	var err error
	if c.nc != nil ***REMOVED***
		err = c.nc.Close()
	***REMOVED***

	return err
***REMOVED***

func (c *connection) closed() bool ***REMOVED***
	return atomic.LoadInt64(&c.state) == connDisconnected
***REMOVED***

func (c *connection) idleTimeoutExpired() bool ***REMOVED***
	now := time.Now()
	if c.idleTimeout > 0 ***REMOVED***
		idleDeadline, ok := c.idleDeadline.Load().(time.Time)
		if ok && now.After(idleDeadline) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***

	return false
***REMOVED***

func (c *connection) bumpIdleDeadline() ***REMOVED***
	if c.idleTimeout > 0 ***REMOVED***
		c.idleDeadline.Store(time.Now().Add(c.idleTimeout))
	***REMOVED***
***REMOVED***

func (c *connection) setCanStream(canStream bool) ***REMOVED***
	c.canStream = canStream
***REMOVED***

func (c initConnection) supportsStreaming() bool ***REMOVED***
	return c.canStream
***REMOVED***

func (c *connection) setStreaming(streaming bool) ***REMOVED***
	c.currentlyStreaming = streaming
***REMOVED***

func (c *connection) getCurrentlyStreaming() bool ***REMOVED***
	return c.currentlyStreaming
***REMOVED***

func (c *connection) setSocketTimeout(timeout time.Duration) ***REMOVED***
	c.readTimeout = timeout
	c.writeTimeout = timeout
***REMOVED***

func (c *connection) ID() string ***REMOVED***
	return c.id
***REMOVED***

func (c *connection) ServerConnectionID() *int32 ***REMOVED***
	return c.serverConnectionID
***REMOVED***

// initConnection is an adapter used during connection initialization. It has the minimum
// functionality necessary to implement the driver.Connection interface, which is required to pass a
// *connection to a Handshaker.
type initConnection struct***REMOVED*** *connection ***REMOVED***

var _ driver.Connection = initConnection***REMOVED******REMOVED***
var _ driver.StreamerConnection = initConnection***REMOVED******REMOVED***

func (c initConnection) Description() description.Server ***REMOVED***
	if c.connection == nil ***REMOVED***
		return description.Server***REMOVED******REMOVED***
	***REMOVED***
	return c.connection.desc
***REMOVED***
func (c initConnection) Close() error             ***REMOVED*** return nil ***REMOVED***
func (c initConnection) ID() string               ***REMOVED*** return c.id ***REMOVED***
func (c initConnection) Address() address.Address ***REMOVED*** return c.addr ***REMOVED***
func (c initConnection) Stale() bool              ***REMOVED*** return false ***REMOVED***
func (c initConnection) LocalAddress() address.Address ***REMOVED***
	if c.connection == nil || c.nc == nil ***REMOVED***
		return address.Address("0.0.0.0")
	***REMOVED***
	return address.Address(c.nc.LocalAddr().String())
***REMOVED***
func (c initConnection) WriteWireMessage(ctx context.Context, wm []byte) error ***REMOVED***
	return c.writeWireMessage(ctx, wm)
***REMOVED***
func (c initConnection) ReadWireMessage(ctx context.Context, dst []byte) ([]byte, error) ***REMOVED***
	return c.readWireMessage(ctx, dst)
***REMOVED***
func (c initConnection) SetStreaming(streaming bool) ***REMOVED***
	c.setStreaming(streaming)
***REMOVED***
func (c initConnection) CurrentlyStreaming() bool ***REMOVED***
	return c.getCurrentlyStreaming()
***REMOVED***
func (c initConnection) SupportsStreaming() bool ***REMOVED***
	return c.supportsStreaming()
***REMOVED***

// Connection implements the driver.Connection interface to allow reading and writing wire
// messages and the driver.Expirable interface to allow expiring.
type Connection struct ***REMOVED***
	*connection
	refCount      int
	cleanupPoolFn func()

	// cleanupServerFn resets the server state when a connection is returned to the connection pool
	// via Close() or expired via Expire().
	cleanupServerFn func()

	mu sync.RWMutex
***REMOVED***

var _ driver.Connection = (*Connection)(nil)
var _ driver.Expirable = (*Connection)(nil)
var _ driver.PinnedConnection = (*Connection)(nil)

// WriteWireMessage handles writing a wire message to the underlying connection.
func (c *Connection) WriteWireMessage(ctx context.Context, wm []byte) error ***REMOVED***
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.connection == nil ***REMOVED***
		return ErrConnectionClosed
	***REMOVED***
	return c.writeWireMessage(ctx, wm)
***REMOVED***

// ReadWireMessage handles reading a wire message from the underlying connection. The dst parameter
// will be overwritten with the new wire message.
func (c *Connection) ReadWireMessage(ctx context.Context, dst []byte) ([]byte, error) ***REMOVED***
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.connection == nil ***REMOVED***
		return dst, ErrConnectionClosed
	***REMOVED***
	return c.readWireMessage(ctx, dst)
***REMOVED***

// CompressWireMessage handles compressing the provided wire message using the underlying
// connection's compressor. The dst parameter will be overwritten with the new wire message. If
// there is no compressor set on the underlying connection, then no compression will be performed.
func (c *Connection) CompressWireMessage(src, dst []byte) ([]byte, error) ***REMOVED***
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.connection == nil ***REMOVED***
		return dst, ErrConnectionClosed
	***REMOVED***
	if c.connection.compressor == wiremessage.CompressorNoOp ***REMOVED***
		if len(dst) == 0 ***REMOVED***
			return src, nil
		***REMOVED***
		return append(dst, src...), nil
	***REMOVED***
	_, reqid, respto, origcode, rem, ok := wiremessage.ReadHeader(src)
	if !ok ***REMOVED***
		return dst, errors.New("wiremessage is too short to compress, less than 16 bytes")
	***REMOVED***
	idx, dst := wiremessage.AppendHeaderStart(dst, reqid, respto, wiremessage.OpCompressed)
	dst = wiremessage.AppendCompressedOriginalOpCode(dst, origcode)
	dst = wiremessage.AppendCompressedUncompressedSize(dst, int32(len(rem)))
	dst = wiremessage.AppendCompressedCompressorID(dst, c.connection.compressor)
	opts := driver.CompressionOpts***REMOVED***
		Compressor: c.connection.compressor,
		ZlibLevel:  c.connection.zliblevel,
		ZstdLevel:  c.connection.zstdLevel,
	***REMOVED***
	compressed, err := driver.CompressPayload(rem, opts)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	dst = wiremessage.AppendCompressedCompressedMessage(dst, compressed)
	return bsoncore.UpdateLength(dst, idx, int32(len(dst[idx:]))), nil
***REMOVED***

// Description returns the server description of the server this connection is connected to.
func (c *Connection) Description() description.Server ***REMOVED***
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.connection == nil ***REMOVED***
		return description.Server***REMOVED******REMOVED***
	***REMOVED***
	return c.desc
***REMOVED***

// Close returns this connection to the connection pool. This method may not closeConnection the underlying
// socket.
func (c *Connection) Close() error ***REMOVED***
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.connection == nil || c.refCount > 0 ***REMOVED***
		return nil
	***REMOVED***

	return c.cleanupReferences()
***REMOVED***

// Expire closes this connection and will closeConnection the underlying socket.
func (c *Connection) Expire() error ***REMOVED***
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.connection == nil ***REMOVED***
		return nil
	***REMOVED***

	_ = c.close()
	return c.cleanupReferences()
***REMOVED***

func (c *Connection) cleanupReferences() error ***REMOVED***
	err := c.pool.checkIn(c.connection)
	if c.cleanupPoolFn != nil ***REMOVED***
		c.cleanupPoolFn()
		c.cleanupPoolFn = nil
	***REMOVED***
	if c.cleanupServerFn != nil ***REMOVED***
		c.cleanupServerFn()
		c.cleanupServerFn = nil
	***REMOVED***
	c.connection = nil
	return err
***REMOVED***

// Alive returns if the connection is still alive.
func (c *Connection) Alive() bool ***REMOVED***
	return c.connection != nil
***REMOVED***

// ID returns the ID of this connection.
func (c *Connection) ID() string ***REMOVED***
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.connection == nil ***REMOVED***
		return "<closed>"
	***REMOVED***
	return c.id
***REMOVED***

// Stale returns if the connection is stale.
func (c *Connection) Stale() bool ***REMOVED***
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.pool.stale(c.connection)
***REMOVED***

// Address returns the address of this connection.
func (c *Connection) Address() address.Address ***REMOVED***
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.connection == nil ***REMOVED***
		return address.Address("0.0.0.0")
	***REMOVED***
	return c.addr
***REMOVED***

// LocalAddress returns the local address of the connection
func (c *Connection) LocalAddress() address.Address ***REMOVED***
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.connection == nil || c.nc == nil ***REMOVED***
		return address.Address("0.0.0.0")
	***REMOVED***
	return address.Address(c.nc.LocalAddr().String())
***REMOVED***

// PinToCursor updates this connection to reflect that it is pinned to a cursor.
func (c *Connection) PinToCursor() error ***REMOVED***
	return c.pin("cursor", c.pool.pinConnectionToCursor, c.pool.unpinConnectionFromCursor)
***REMOVED***

// PinToTransaction updates this connection to reflect that it is pinned to a transaction.
func (c *Connection) PinToTransaction() error ***REMOVED***
	return c.pin("transaction", c.pool.pinConnectionToTransaction, c.pool.unpinConnectionFromTransaction)
***REMOVED***

func (c *Connection) pin(reason string, updatePoolFn, cleanupPoolFn func()) error ***REMOVED***
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.connection == nil ***REMOVED***
		return fmt.Errorf("attempted to pin a connection for a %s, but the connection has already been returned to the pool", reason)
	***REMOVED***

	// Only use the provided callbacks for the first reference to avoid double-counting pinned connection statistics
	// in the pool.
	if c.refCount == 0 ***REMOVED***
		updatePoolFn()
		c.cleanupPoolFn = cleanupPoolFn
	***REMOVED***
	c.refCount++
	return nil
***REMOVED***

// UnpinFromCursor updates this connection to reflect that it is no longer pinned to a cursor.
func (c *Connection) UnpinFromCursor() error ***REMOVED***
	return c.unpin("cursor")
***REMOVED***

// UnpinFromTransaction updates this connection to reflect that it is no longer pinned to a transaction.
func (c *Connection) UnpinFromTransaction() error ***REMOVED***
	return c.unpin("transaction")
***REMOVED***

func (c *Connection) unpin(reason string) error ***REMOVED***
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.connection == nil ***REMOVED***
		// We don't error here because the resource could have been forcefully closed via Expire.
		return nil
	***REMOVED***
	if c.refCount == 0 ***REMOVED***
		return fmt.Errorf("attempted to unpin a connection from a %s, but the connection is not pinned by any resources", reason)
	***REMOVED***

	c.refCount--
	return nil
***REMOVED***

func configureTLS(ctx context.Context,
	tlsConnSource tlsConnectionSource,
	nc net.Conn,
	addr address.Address,
	config *tls.Config,
	ocspOpts *ocsp.VerifyOptions,
) (net.Conn, error) ***REMOVED***
	// Ensure config.ServerName is always set for SNI.
	if config.ServerName == "" ***REMOVED***
		hostname := addr.String()
		colonPos := strings.LastIndex(hostname, ":")
		if colonPos == -1 ***REMOVED***
			colonPos = len(hostname)
		***REMOVED***

		hostname = hostname[:colonPos]
		config.ServerName = hostname
	***REMOVED***

	client := tlsConnSource.Client(nc, config)
	if err := clientHandshake(ctx, client); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Only do OCSP verification if TLS verification is requested.
	if !config.InsecureSkipVerify ***REMOVED***
		if ocspErr := ocsp.Verify(ctx, client.ConnectionState(), ocspOpts); ocspErr != nil ***REMOVED***
			return nil, ocspErr
		***REMOVED***
	***REMOVED***
	return client, nil
***REMOVED***
