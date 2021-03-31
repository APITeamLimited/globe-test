/*
 *
 * Copyright 2014 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package transport

import (
	"context"
	"fmt"
	"io"
	"math"
	"net"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
	"google.golang.org/grpc/internal/grpcutil"
	imetadata "google.golang.org/grpc/internal/metadata"
	"google.golang.org/grpc/internal/transport/networktype"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/internal"
	"google.golang.org/grpc/internal/channelz"
	"google.golang.org/grpc/internal/syscall"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
)

// clientConnectionCounter counts the number of connections a client has
// initiated (equal to the number of http2Clients created). Must be accessed
// atomically.
var clientConnectionCounter uint64

// http2Client implements the ClientTransport interface with HTTP2.
type http2Client struct ***REMOVED***
	lastRead   int64 // Keep this field 64-bit aligned. Accessed atomically.
	ctx        context.Context
	cancel     context.CancelFunc
	ctxDone    <-chan struct***REMOVED******REMOVED*** // Cache the ctx.Done() chan.
	userAgent  string
	md         metadata.MD
	conn       net.Conn // underlying communication channel
	loopy      *loopyWriter
	remoteAddr net.Addr
	localAddr  net.Addr
	authInfo   credentials.AuthInfo // auth info about the connection

	readerDone chan struct***REMOVED******REMOVED*** // sync point to enable testing.
	writerDone chan struct***REMOVED******REMOVED*** // sync point to enable testing.
	// goAway is closed to notify the upper layer (i.e., addrConn.transportMonitor)
	// that the server sent GoAway on this transport.
	goAway chan struct***REMOVED******REMOVED***

	framer *framer
	// controlBuf delivers all the control related tasks (e.g., window
	// updates, reset streams, and various settings) to the controller.
	controlBuf *controlBuffer
	fc         *trInFlow
	// The scheme used: https if TLS is on, http otherwise.
	scheme string

	isSecure bool

	perRPCCreds []credentials.PerRPCCredentials

	kp               keepalive.ClientParameters
	keepaliveEnabled bool

	statsHandler stats.Handler

	initialWindowSize int32

	// configured by peer through SETTINGS_MAX_HEADER_LIST_SIZE
	maxSendHeaderListSize *uint32

	bdpEst *bdpEstimator
	// onPrefaceReceipt is a callback that client transport calls upon
	// receiving server preface to signal that a succefull HTTP2
	// connection was established.
	onPrefaceReceipt func()

	maxConcurrentStreams  uint32
	streamQuota           int64
	streamsQuotaAvailable chan struct***REMOVED******REMOVED***
	waitingStreams        uint32
	nextID                uint32

	mu            sync.Mutex // guard the following variables
	state         transportState
	activeStreams map[uint32]*Stream
	// prevGoAway ID records the Last-Stream-ID in the previous GOAway frame.
	prevGoAwayID uint32
	// goAwayReason records the http2.ErrCode and debug data received with the
	// GoAway frame.
	goAwayReason GoAwayReason
	// A condition variable used to signal when the keepalive goroutine should
	// go dormant. The condition for dormancy is based on the number of active
	// streams and the `PermitWithoutStream` keepalive client parameter. And
	// since the number of active streams is guarded by the above mutex, we use
	// the same for this condition variable as well.
	kpDormancyCond *sync.Cond
	// A boolean to track whether the keepalive goroutine is dormant or not.
	// This is checked before attempting to signal the above condition
	// variable.
	kpDormant bool

	// Fields below are for channelz metric collection.
	channelzID int64 // channelz unique identification number
	czData     *channelzData

	onGoAway func(GoAwayReason)
	onClose  func()

	bufferPool *bufferPool

	connectionID uint64
***REMOVED***

func dial(ctx context.Context, fn func(context.Context, string) (net.Conn, error), addr resolver.Address, useProxy bool, grpcUA string) (net.Conn, error) ***REMOVED***
	address := addr.Addr
	networkType, ok := networktype.Get(addr)
	if fn != nil ***REMOVED***
		if networkType == "unix" && !strings.HasPrefix(address, "\x00") ***REMOVED***
			// For backward compatibility, if the user dialed "unix:///path",
			// the passthrough resolver would be used and the user's custom
			// dialer would see "unix:///path". Since the unix resolver is used
			// and the address is now "/path", prepend "unix://" so the user's
			// custom dialer sees the same address.
			return fn(ctx, "unix://"+address)
		***REMOVED***
		return fn(ctx, address)
	***REMOVED***
	if !ok ***REMOVED***
		networkType, address = parseDialTarget(address)
	***REMOVED***
	if networkType == "tcp" && useProxy ***REMOVED***
		return proxyDial(ctx, address, grpcUA)
	***REMOVED***
	return (&net.Dialer***REMOVED******REMOVED***).DialContext(ctx, networkType, address)
***REMOVED***

func isTemporary(err error) bool ***REMOVED***
	switch err := err.(type) ***REMOVED***
	case interface ***REMOVED***
		Temporary() bool
	***REMOVED***:
		return err.Temporary()
	case interface ***REMOVED***
		Timeout() bool
	***REMOVED***:
		// Timeouts may be resolved upon retry, and are thus treated as
		// temporary.
		return err.Timeout()
	***REMOVED***
	return true
***REMOVED***

// newHTTP2Client constructs a connected ClientTransport to addr based on HTTP2
// and starts to receive messages on it. Non-nil error returns if construction
// fails.
func newHTTP2Client(connectCtx, ctx context.Context, addr resolver.Address, opts ConnectOptions, onPrefaceReceipt func(), onGoAway func(GoAwayReason), onClose func()) (_ *http2Client, err error) ***REMOVED***
	scheme := "http"
	ctx, cancel := context.WithCancel(ctx)
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			cancel()
		***REMOVED***
	***REMOVED***()

	conn, err := dial(connectCtx, opts.Dialer, addr, opts.UseProxy, opts.UserAgent)
	if err != nil ***REMOVED***
		if opts.FailOnNonTempDialError ***REMOVED***
			return nil, connectionErrorf(isTemporary(err), err, "transport: error while dialing: %v", err)
		***REMOVED***
		return nil, connectionErrorf(true, err, "transport: Error while dialing %v", err)
	***REMOVED***
	// Any further errors will close the underlying connection
	defer func(conn net.Conn) ***REMOVED***
		if err != nil ***REMOVED***
			conn.Close()
		***REMOVED***
	***REMOVED***(conn)
	kp := opts.KeepaliveParams
	// Validate keepalive parameters.
	if kp.Time == 0 ***REMOVED***
		kp.Time = defaultClientKeepaliveTime
	***REMOVED***
	if kp.Timeout == 0 ***REMOVED***
		kp.Timeout = defaultClientKeepaliveTimeout
	***REMOVED***
	keepaliveEnabled := false
	if kp.Time != infinity ***REMOVED***
		if err = syscall.SetTCPUserTimeout(conn, kp.Timeout); err != nil ***REMOVED***
			return nil, connectionErrorf(false, err, "transport: failed to set TCP_USER_TIMEOUT: %v", err)
		***REMOVED***
		keepaliveEnabled = true
	***REMOVED***
	var (
		isSecure bool
		authInfo credentials.AuthInfo
	)
	transportCreds := opts.TransportCredentials
	perRPCCreds := opts.PerRPCCredentials

	if b := opts.CredsBundle; b != nil ***REMOVED***
		if t := b.TransportCredentials(); t != nil ***REMOVED***
			transportCreds = t
		***REMOVED***
		if t := b.PerRPCCredentials(); t != nil ***REMOVED***
			perRPCCreds = append(perRPCCreds, t)
		***REMOVED***
	***REMOVED***
	if transportCreds != nil ***REMOVED***
		// gRPC, resolver, balancer etc. can specify arbitrary data in the
		// Attributes field of resolver.Address, which is shoved into connectCtx
		// and passed to the credential handshaker. This makes it possible for
		// address specific arbitrary data to reach the credential handshaker.
		contextWithHandshakeInfo := internal.NewClientHandshakeInfoContext.(func(context.Context, credentials.ClientHandshakeInfo) context.Context)
		connectCtx = contextWithHandshakeInfo(connectCtx, credentials.ClientHandshakeInfo***REMOVED***Attributes: addr.Attributes***REMOVED***)
		conn, authInfo, err = transportCreds.ClientHandshake(connectCtx, addr.ServerName, conn)
		if err != nil ***REMOVED***
			return nil, connectionErrorf(isTemporary(err), err, "transport: authentication handshake failed: %v", err)
		***REMOVED***
		for _, cd := range perRPCCreds ***REMOVED***
			if cd.RequireTransportSecurity() ***REMOVED***
				if ci, ok := authInfo.(interface ***REMOVED***
					GetCommonAuthInfo() credentials.CommonAuthInfo
				***REMOVED***); ok ***REMOVED***
					secLevel := ci.GetCommonAuthInfo().SecurityLevel
					if secLevel != credentials.InvalidSecurityLevel && secLevel < credentials.PrivacyAndIntegrity ***REMOVED***
						return nil, connectionErrorf(true, nil, "transport: cannot send secure credentials on an insecure connection")
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
		isSecure = true
		if transportCreds.Info().SecurityProtocol == "tls" ***REMOVED***
			scheme = "https"
		***REMOVED***
	***REMOVED***
	dynamicWindow := true
	icwz := int32(initialWindowSize)
	if opts.InitialConnWindowSize >= defaultWindowSize ***REMOVED***
		icwz = opts.InitialConnWindowSize
		dynamicWindow = false
	***REMOVED***
	writeBufSize := opts.WriteBufferSize
	readBufSize := opts.ReadBufferSize
	maxHeaderListSize := defaultClientMaxHeaderListSize
	if opts.MaxHeaderListSize != nil ***REMOVED***
		maxHeaderListSize = *opts.MaxHeaderListSize
	***REMOVED***
	t := &http2Client***REMOVED***
		ctx:                   ctx,
		ctxDone:               ctx.Done(), // Cache Done chan.
		cancel:                cancel,
		userAgent:             opts.UserAgent,
		conn:                  conn,
		remoteAddr:            conn.RemoteAddr(),
		localAddr:             conn.LocalAddr(),
		authInfo:              authInfo,
		readerDone:            make(chan struct***REMOVED******REMOVED***),
		writerDone:            make(chan struct***REMOVED******REMOVED***),
		goAway:                make(chan struct***REMOVED******REMOVED***),
		framer:                newFramer(conn, writeBufSize, readBufSize, maxHeaderListSize),
		fc:                    &trInFlow***REMOVED***limit: uint32(icwz)***REMOVED***,
		scheme:                scheme,
		activeStreams:         make(map[uint32]*Stream),
		isSecure:              isSecure,
		perRPCCreds:           perRPCCreds,
		kp:                    kp,
		statsHandler:          opts.StatsHandler,
		initialWindowSize:     initialWindowSize,
		onPrefaceReceipt:      onPrefaceReceipt,
		nextID:                1,
		maxConcurrentStreams:  defaultMaxStreamsClient,
		streamQuota:           defaultMaxStreamsClient,
		streamsQuotaAvailable: make(chan struct***REMOVED******REMOVED***, 1),
		czData:                new(channelzData),
		onGoAway:              onGoAway,
		onClose:               onClose,
		keepaliveEnabled:      keepaliveEnabled,
		bufferPool:            newBufferPool(),
	***REMOVED***

	if md, ok := addr.Metadata.(*metadata.MD); ok ***REMOVED***
		t.md = *md
	***REMOVED*** else if md := imetadata.Get(addr); md != nil ***REMOVED***
		t.md = md
	***REMOVED***
	t.controlBuf = newControlBuffer(t.ctxDone)
	if opts.InitialWindowSize >= defaultWindowSize ***REMOVED***
		t.initialWindowSize = opts.InitialWindowSize
		dynamicWindow = false
	***REMOVED***
	if dynamicWindow ***REMOVED***
		t.bdpEst = &bdpEstimator***REMOVED***
			bdp:               initialWindowSize,
			updateFlowControl: t.updateFlowControl,
		***REMOVED***
	***REMOVED***
	if t.statsHandler != nil ***REMOVED***
		t.ctx = t.statsHandler.TagConn(t.ctx, &stats.ConnTagInfo***REMOVED***
			RemoteAddr: t.remoteAddr,
			LocalAddr:  t.localAddr,
		***REMOVED***)
		connBegin := &stats.ConnBegin***REMOVED***
			Client: true,
		***REMOVED***
		t.statsHandler.HandleConn(t.ctx, connBegin)
	***REMOVED***
	if channelz.IsOn() ***REMOVED***
		t.channelzID = channelz.RegisterNormalSocket(t, opts.ChannelzParentID, fmt.Sprintf("%s -> %s", t.localAddr, t.remoteAddr))
	***REMOVED***
	if t.keepaliveEnabled ***REMOVED***
		t.kpDormancyCond = sync.NewCond(&t.mu)
		go t.keepalive()
	***REMOVED***
	// Start the reader goroutine for incoming message. Each transport has
	// a dedicated goroutine which reads HTTP2 frame from network. Then it
	// dispatches the frame to the corresponding stream entity.
	go t.reader()

	// Send connection preface to server.
	n, err := t.conn.Write(clientPreface)
	if err != nil ***REMOVED***
		t.Close()
		return nil, connectionErrorf(true, err, "transport: failed to write client preface: %v", err)
	***REMOVED***
	if n != len(clientPreface) ***REMOVED***
		t.Close()
		return nil, connectionErrorf(true, err, "transport: preface mismatch, wrote %d bytes; want %d", n, len(clientPreface))
	***REMOVED***
	var ss []http2.Setting

	if t.initialWindowSize != defaultWindowSize ***REMOVED***
		ss = append(ss, http2.Setting***REMOVED***
			ID:  http2.SettingInitialWindowSize,
			Val: uint32(t.initialWindowSize),
		***REMOVED***)
	***REMOVED***
	if opts.MaxHeaderListSize != nil ***REMOVED***
		ss = append(ss, http2.Setting***REMOVED***
			ID:  http2.SettingMaxHeaderListSize,
			Val: *opts.MaxHeaderListSize,
		***REMOVED***)
	***REMOVED***
	err = t.framer.fr.WriteSettings(ss...)
	if err != nil ***REMOVED***
		t.Close()
		return nil, connectionErrorf(true, err, "transport: failed to write initial settings frame: %v", err)
	***REMOVED***
	// Adjust the connection flow control window if needed.
	if delta := uint32(icwz - defaultWindowSize); delta > 0 ***REMOVED***
		if err := t.framer.fr.WriteWindowUpdate(0, delta); err != nil ***REMOVED***
			t.Close()
			return nil, connectionErrorf(true, err, "transport: failed to write window update: %v", err)
		***REMOVED***
	***REMOVED***

	t.connectionID = atomic.AddUint64(&clientConnectionCounter, 1)

	if err := t.framer.writer.Flush(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	go func() ***REMOVED***
		t.loopy = newLoopyWriter(clientSide, t.framer, t.controlBuf, t.bdpEst)
		err := t.loopy.run()
		if err != nil ***REMOVED***
			if logger.V(logLevel) ***REMOVED***
				logger.Errorf("transport: loopyWriter.run returning. Err: %v", err)
			***REMOVED***
		***REMOVED***
		// If it's a connection error, let reader goroutine handle it
		// since there might be data in the buffers.
		if _, ok := err.(net.Error); !ok ***REMOVED***
			t.conn.Close()
		***REMOVED***
		close(t.writerDone)
	***REMOVED***()
	return t, nil
***REMOVED***

func (t *http2Client) newStream(ctx context.Context, callHdr *CallHdr) *Stream ***REMOVED***
	// TODO(zhaoq): Handle uint32 overflow of Stream.id.
	s := &Stream***REMOVED***
		ct:             t,
		done:           make(chan struct***REMOVED******REMOVED***),
		method:         callHdr.Method,
		sendCompress:   callHdr.SendCompress,
		buf:            newRecvBuffer(),
		headerChan:     make(chan struct***REMOVED******REMOVED***),
		contentSubtype: callHdr.ContentSubtype,
	***REMOVED***
	s.wq = newWriteQuota(defaultWriteQuota, s.done)
	s.requestRead = func(n int) ***REMOVED***
		t.adjustWindow(s, uint32(n))
	***REMOVED***
	// The client side stream context should have exactly the same life cycle with the user provided context.
	// That means, s.ctx should be read-only. And s.ctx is done iff ctx is done.
	// So we use the original context here instead of creating a copy.
	s.ctx = ctx
	s.trReader = &transportReader***REMOVED***
		reader: &recvBufferReader***REMOVED***
			ctx:     s.ctx,
			ctxDone: s.ctx.Done(),
			recv:    s.buf,
			closeStream: func(err error) ***REMOVED***
				t.CloseStream(s, err)
			***REMOVED***,
			freeBuffer: t.bufferPool.put,
		***REMOVED***,
		windowHandler: func(n int) ***REMOVED***
			t.updateWindow(s, uint32(n))
		***REMOVED***,
	***REMOVED***
	return s
***REMOVED***

func (t *http2Client) getPeer() *peer.Peer ***REMOVED***
	return &peer.Peer***REMOVED***
		Addr:     t.remoteAddr,
		AuthInfo: t.authInfo,
	***REMOVED***
***REMOVED***

func (t *http2Client) createHeaderFields(ctx context.Context, callHdr *CallHdr) ([]hpack.HeaderField, error) ***REMOVED***
	aud := t.createAudience(callHdr)
	ri := credentials.RequestInfo***REMOVED***
		Method:   callHdr.Method,
		AuthInfo: t.authInfo,
	***REMOVED***
	ctxWithRequestInfo := internal.NewRequestInfoContext.(func(context.Context, credentials.RequestInfo) context.Context)(ctx, ri)
	authData, err := t.getTrAuthData(ctxWithRequestInfo, aud)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	callAuthData, err := t.getCallAuthData(ctxWithRequestInfo, aud, callHdr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	// TODO(mmukhi): Benchmark if the performance gets better if count the metadata and other header fields
	// first and create a slice of that exact size.
	// Make the slice of certain predictable size to reduce allocations made by append.
	hfLen := 7 // :method, :scheme, :path, :authority, content-type, user-agent, te
	hfLen += len(authData) + len(callAuthData)
	headerFields := make([]hpack.HeaderField, 0, hfLen)
	headerFields = append(headerFields, hpack.HeaderField***REMOVED***Name: ":method", Value: "POST"***REMOVED***)
	headerFields = append(headerFields, hpack.HeaderField***REMOVED***Name: ":scheme", Value: t.scheme***REMOVED***)
	headerFields = append(headerFields, hpack.HeaderField***REMOVED***Name: ":path", Value: callHdr.Method***REMOVED***)
	headerFields = append(headerFields, hpack.HeaderField***REMOVED***Name: ":authority", Value: callHdr.Host***REMOVED***)
	headerFields = append(headerFields, hpack.HeaderField***REMOVED***Name: "content-type", Value: grpcutil.ContentType(callHdr.ContentSubtype)***REMOVED***)
	headerFields = append(headerFields, hpack.HeaderField***REMOVED***Name: "user-agent", Value: t.userAgent***REMOVED***)
	headerFields = append(headerFields, hpack.HeaderField***REMOVED***Name: "te", Value: "trailers"***REMOVED***)
	if callHdr.PreviousAttempts > 0 ***REMOVED***
		headerFields = append(headerFields, hpack.HeaderField***REMOVED***Name: "grpc-previous-rpc-attempts", Value: strconv.Itoa(callHdr.PreviousAttempts)***REMOVED***)
	***REMOVED***

	if callHdr.SendCompress != "" ***REMOVED***
		headerFields = append(headerFields, hpack.HeaderField***REMOVED***Name: "grpc-encoding", Value: callHdr.SendCompress***REMOVED***)
		headerFields = append(headerFields, hpack.HeaderField***REMOVED***Name: "grpc-accept-encoding", Value: callHdr.SendCompress***REMOVED***)
	***REMOVED***
	if dl, ok := ctx.Deadline(); ok ***REMOVED***
		// Send out timeout regardless its value. The server can detect timeout context by itself.
		// TODO(mmukhi): Perhaps this field should be updated when actually writing out to the wire.
		timeout := time.Until(dl)
		headerFields = append(headerFields, hpack.HeaderField***REMOVED***Name: "grpc-timeout", Value: grpcutil.EncodeDuration(timeout)***REMOVED***)
	***REMOVED***
	for k, v := range authData ***REMOVED***
		headerFields = append(headerFields, hpack.HeaderField***REMOVED***Name: k, Value: encodeMetadataHeader(k, v)***REMOVED***)
	***REMOVED***
	for k, v := range callAuthData ***REMOVED***
		headerFields = append(headerFields, hpack.HeaderField***REMOVED***Name: k, Value: encodeMetadataHeader(k, v)***REMOVED***)
	***REMOVED***
	if b := stats.OutgoingTags(ctx); b != nil ***REMOVED***
		headerFields = append(headerFields, hpack.HeaderField***REMOVED***Name: "grpc-tags-bin", Value: encodeBinHeader(b)***REMOVED***)
	***REMOVED***
	if b := stats.OutgoingTrace(ctx); b != nil ***REMOVED***
		headerFields = append(headerFields, hpack.HeaderField***REMOVED***Name: "grpc-trace-bin", Value: encodeBinHeader(b)***REMOVED***)
	***REMOVED***

	if md, added, ok := metadata.FromOutgoingContextRaw(ctx); ok ***REMOVED***
		var k string
		for k, vv := range md ***REMOVED***
			// HTTP doesn't allow you to set pseudoheaders after non pseudoheaders were set.
			if isReservedHeader(k) ***REMOVED***
				continue
			***REMOVED***
			for _, v := range vv ***REMOVED***
				headerFields = append(headerFields, hpack.HeaderField***REMOVED***Name: k, Value: encodeMetadataHeader(k, v)***REMOVED***)
			***REMOVED***
		***REMOVED***
		for _, vv := range added ***REMOVED***
			for i, v := range vv ***REMOVED***
				if i%2 == 0 ***REMOVED***
					k = strings.ToLower(v)
					continue
				***REMOVED***
				// HTTP doesn't allow you to set pseudoheaders after non pseudoheaders were set.
				if isReservedHeader(k) ***REMOVED***
					continue
				***REMOVED***
				headerFields = append(headerFields, hpack.HeaderField***REMOVED***Name: k, Value: encodeMetadataHeader(k, v)***REMOVED***)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	for k, vv := range t.md ***REMOVED***
		if isReservedHeader(k) ***REMOVED***
			continue
		***REMOVED***
		for _, v := range vv ***REMOVED***
			headerFields = append(headerFields, hpack.HeaderField***REMOVED***Name: k, Value: encodeMetadataHeader(k, v)***REMOVED***)
		***REMOVED***
	***REMOVED***
	return headerFields, nil
***REMOVED***

func (t *http2Client) createAudience(callHdr *CallHdr) string ***REMOVED***
	// Create an audience string only if needed.
	if len(t.perRPCCreds) == 0 && callHdr.Creds == nil ***REMOVED***
		return ""
	***REMOVED***
	// Construct URI required to get auth request metadata.
	// Omit port if it is the default one.
	host := strings.TrimSuffix(callHdr.Host, ":443")
	pos := strings.LastIndex(callHdr.Method, "/")
	if pos == -1 ***REMOVED***
		pos = len(callHdr.Method)
	***REMOVED***
	return "https://" + host + callHdr.Method[:pos]
***REMOVED***

func (t *http2Client) getTrAuthData(ctx context.Context, audience string) (map[string]string, error) ***REMOVED***
	if len(t.perRPCCreds) == 0 ***REMOVED***
		return nil, nil
	***REMOVED***
	authData := map[string]string***REMOVED******REMOVED***
	for _, c := range t.perRPCCreds ***REMOVED***
		data, err := c.GetRequestMetadata(ctx, audience)
		if err != nil ***REMOVED***
			if _, ok := status.FromError(err); ok ***REMOVED***
				return nil, err
			***REMOVED***

			return nil, status.Errorf(codes.Unauthenticated, "transport: %v", err)
		***REMOVED***
		for k, v := range data ***REMOVED***
			// Capital header names are illegal in HTTP/2.
			k = strings.ToLower(k)
			authData[k] = v
		***REMOVED***
	***REMOVED***
	return authData, nil
***REMOVED***

func (t *http2Client) getCallAuthData(ctx context.Context, audience string, callHdr *CallHdr) (map[string]string, error) ***REMOVED***
	var callAuthData map[string]string
	// Check if credentials.PerRPCCredentials were provided via call options.
	// Note: if these credentials are provided both via dial options and call
	// options, then both sets of credentials will be applied.
	if callCreds := callHdr.Creds; callCreds != nil ***REMOVED***
		if callCreds.RequireTransportSecurity() ***REMOVED***
			ri, _ := credentials.RequestInfoFromContext(ctx)
			if !t.isSecure || credentials.CheckSecurityLevel(ri.AuthInfo, credentials.PrivacyAndIntegrity) != nil ***REMOVED***
				return nil, status.Error(codes.Unauthenticated, "transport: cannot send secure credentials on an insecure connection")
			***REMOVED***
		***REMOVED***
		data, err := callCreds.GetRequestMetadata(ctx, audience)
		if err != nil ***REMOVED***
			return nil, status.Errorf(codes.Internal, "transport: %v", err)
		***REMOVED***
		callAuthData = make(map[string]string, len(data))
		for k, v := range data ***REMOVED***
			// Capital header names are illegal in HTTP/2
			k = strings.ToLower(k)
			callAuthData[k] = v
		***REMOVED***
	***REMOVED***
	return callAuthData, nil
***REMOVED***

// PerformedIOError wraps an error to indicate IO may have been performed
// before the error occurred.
type PerformedIOError struct ***REMOVED***
	Err error
***REMOVED***

// Error implements error.
func (p PerformedIOError) Error() string ***REMOVED***
	return p.Err.Error()
***REMOVED***

// NewStream creates a stream and registers it into the transport as "active"
// streams.
func (t *http2Client) NewStream(ctx context.Context, callHdr *CallHdr) (_ *Stream, err error) ***REMOVED***
	ctx = peer.NewContext(ctx, t.getPeer())
	headerFields, err := t.createHeaderFields(ctx, callHdr)
	if err != nil ***REMOVED***
		// We may have performed I/O in the per-RPC creds callback, so do not
		// allow transparent retry.
		return nil, PerformedIOError***REMOVED***err***REMOVED***
	***REMOVED***
	s := t.newStream(ctx, callHdr)
	cleanup := func(err error) ***REMOVED***
		if s.swapState(streamDone) == streamDone ***REMOVED***
			// If it was already done, return.
			return
		***REMOVED***
		// The stream was unprocessed by the server.
		atomic.StoreUint32(&s.unprocessed, 1)
		s.write(recvMsg***REMOVED***err: err***REMOVED***)
		close(s.done)
		// If headerChan isn't closed, then close it.
		if atomic.CompareAndSwapUint32(&s.headerChanClosed, 0, 1) ***REMOVED***
			close(s.headerChan)
		***REMOVED***
	***REMOVED***
	hdr := &headerFrame***REMOVED***
		hf:        headerFields,
		endStream: false,
		initStream: func(id uint32) error ***REMOVED***
			t.mu.Lock()
			if state := t.state; state != reachable ***REMOVED***
				t.mu.Unlock()
				// Do a quick cleanup.
				err := error(errStreamDrain)
				if state == closing ***REMOVED***
					err = ErrConnClosing
				***REMOVED***
				cleanup(err)
				return err
			***REMOVED***
			t.activeStreams[id] = s
			if channelz.IsOn() ***REMOVED***
				atomic.AddInt64(&t.czData.streamsStarted, 1)
				atomic.StoreInt64(&t.czData.lastStreamCreatedTime, time.Now().UnixNano())
			***REMOVED***
			// If the keepalive goroutine has gone dormant, wake it up.
			if t.kpDormant ***REMOVED***
				t.kpDormancyCond.Signal()
			***REMOVED***
			t.mu.Unlock()
			return nil
		***REMOVED***,
		onOrphaned: cleanup,
		wq:         s.wq,
	***REMOVED***
	firstTry := true
	var ch chan struct***REMOVED******REMOVED***
	checkForStreamQuota := func(it interface***REMOVED******REMOVED***) bool ***REMOVED***
		if t.streamQuota <= 0 ***REMOVED*** // Can go negative if server decreases it.
			if firstTry ***REMOVED***
				t.waitingStreams++
			***REMOVED***
			ch = t.streamsQuotaAvailable
			return false
		***REMOVED***
		if !firstTry ***REMOVED***
			t.waitingStreams--
		***REMOVED***
		t.streamQuota--
		h := it.(*headerFrame)
		h.streamID = t.nextID
		t.nextID += 2
		s.id = h.streamID
		s.fc = &inFlow***REMOVED***limit: uint32(t.initialWindowSize)***REMOVED***
		if t.streamQuota > 0 && t.waitingStreams > 0 ***REMOVED***
			select ***REMOVED***
			case t.streamsQuotaAvailable <- struct***REMOVED******REMOVED******REMOVED******REMOVED***:
			default:
			***REMOVED***
		***REMOVED***
		return true
	***REMOVED***
	var hdrListSizeErr error
	checkForHeaderListSize := func(it interface***REMOVED******REMOVED***) bool ***REMOVED***
		if t.maxSendHeaderListSize == nil ***REMOVED***
			return true
		***REMOVED***
		hdrFrame := it.(*headerFrame)
		var sz int64
		for _, f := range hdrFrame.hf ***REMOVED***
			if sz += int64(f.Size()); sz > int64(*t.maxSendHeaderListSize) ***REMOVED***
				hdrListSizeErr = status.Errorf(codes.Internal, "header list size to send violates the maximum size (%d bytes) set by server", *t.maxSendHeaderListSize)
				return false
			***REMOVED***
		***REMOVED***
		return true
	***REMOVED***
	for ***REMOVED***
		success, err := t.controlBuf.executeAndPut(func(it interface***REMOVED******REMOVED***) bool ***REMOVED***
			if !checkForStreamQuota(it) ***REMOVED***
				return false
			***REMOVED***
			if !checkForHeaderListSize(it) ***REMOVED***
				return false
			***REMOVED***
			return true
		***REMOVED***, hdr)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if success ***REMOVED***
			break
		***REMOVED***
		if hdrListSizeErr != nil ***REMOVED***
			return nil, hdrListSizeErr
		***REMOVED***
		firstTry = false
		select ***REMOVED***
		case <-ch:
		case <-s.ctx.Done():
			return nil, ContextErr(s.ctx.Err())
		case <-t.goAway:
			return nil, errStreamDrain
		case <-t.ctx.Done():
			return nil, ErrConnClosing
		***REMOVED***
	***REMOVED***
	if t.statsHandler != nil ***REMOVED***
		header, ok := metadata.FromOutgoingContext(ctx)
		if ok ***REMOVED***
			header.Set("user-agent", t.userAgent)
		***REMOVED*** else ***REMOVED***
			header = metadata.Pairs("user-agent", t.userAgent)
		***REMOVED***
		// Note: The header fields are compressed with hpack after this call returns.
		// No WireLength field is set here.
		outHeader := &stats.OutHeader***REMOVED***
			Client:      true,
			FullMethod:  callHdr.Method,
			RemoteAddr:  t.remoteAddr,
			LocalAddr:   t.localAddr,
			Compression: callHdr.SendCompress,
			Header:      header,
		***REMOVED***
		t.statsHandler.HandleRPC(s.ctx, outHeader)
	***REMOVED***
	return s, nil
***REMOVED***

// CloseStream clears the footprint of a stream when the stream is not needed any more.
// This must not be executed in reader's goroutine.
func (t *http2Client) CloseStream(s *Stream, err error) ***REMOVED***
	var (
		rst     bool
		rstCode http2.ErrCode
	)
	if err != nil ***REMOVED***
		rst = true
		rstCode = http2.ErrCodeCancel
	***REMOVED***
	t.closeStream(s, err, rst, rstCode, status.Convert(err), nil, false)
***REMOVED***

func (t *http2Client) closeStream(s *Stream, err error, rst bool, rstCode http2.ErrCode, st *status.Status, mdata map[string][]string, eosReceived bool) ***REMOVED***
	// Set stream status to done.
	if s.swapState(streamDone) == streamDone ***REMOVED***
		// If it was already done, return.  If multiple closeStream calls
		// happen simultaneously, wait for the first to finish.
		<-s.done
		return
	***REMOVED***
	// status and trailers can be updated here without any synchronization because the stream goroutine will
	// only read it after it sees an io.EOF error from read or write and we'll write those errors
	// only after updating this.
	s.status = st
	if len(mdata) > 0 ***REMOVED***
		s.trailer = mdata
	***REMOVED***
	if err != nil ***REMOVED***
		// This will unblock reads eventually.
		s.write(recvMsg***REMOVED***err: err***REMOVED***)
	***REMOVED***
	// If headerChan isn't closed, then close it.
	if atomic.CompareAndSwapUint32(&s.headerChanClosed, 0, 1) ***REMOVED***
		s.noHeaders = true
		close(s.headerChan)
	***REMOVED***
	cleanup := &cleanupStream***REMOVED***
		streamID: s.id,
		onWrite: func() ***REMOVED***
			t.mu.Lock()
			if t.activeStreams != nil ***REMOVED***
				delete(t.activeStreams, s.id)
			***REMOVED***
			t.mu.Unlock()
			if channelz.IsOn() ***REMOVED***
				if eosReceived ***REMOVED***
					atomic.AddInt64(&t.czData.streamsSucceeded, 1)
				***REMOVED*** else ***REMOVED***
					atomic.AddInt64(&t.czData.streamsFailed, 1)
				***REMOVED***
			***REMOVED***
		***REMOVED***,
		rst:     rst,
		rstCode: rstCode,
	***REMOVED***
	addBackStreamQuota := func(interface***REMOVED******REMOVED***) bool ***REMOVED***
		t.streamQuota++
		if t.streamQuota > 0 && t.waitingStreams > 0 ***REMOVED***
			select ***REMOVED***
			case t.streamsQuotaAvailable <- struct***REMOVED******REMOVED******REMOVED******REMOVED***:
			default:
			***REMOVED***
		***REMOVED***
		return true
	***REMOVED***
	t.controlBuf.executeAndPut(addBackStreamQuota, cleanup)
	// This will unblock write.
	close(s.done)
***REMOVED***

// Close kicks off the shutdown process of the transport. This should be called
// only once on a transport. Once it is called, the transport should not be
// accessed any more.
//
// This method blocks until the addrConn that initiated this transport is
// re-connected. This happens because t.onClose() begins reconnect logic at the
// addrConn level and blocks until the addrConn is successfully connected.
func (t *http2Client) Close() error ***REMOVED***
	t.mu.Lock()
	// Make sure we only Close once.
	if t.state == closing ***REMOVED***
		t.mu.Unlock()
		return nil
	***REMOVED***
	// Call t.onClose before setting the state to closing to prevent the client
	// from attempting to create new streams ASAP.
	t.onClose()
	t.state = closing
	streams := t.activeStreams
	t.activeStreams = nil
	if t.kpDormant ***REMOVED***
		// If the keepalive goroutine is blocked on this condition variable, we
		// should unblock it so that the goroutine eventually exits.
		t.kpDormancyCond.Signal()
	***REMOVED***
	t.mu.Unlock()
	t.controlBuf.finish()
	t.cancel()
	err := t.conn.Close()
	if channelz.IsOn() ***REMOVED***
		channelz.RemoveEntry(t.channelzID)
	***REMOVED***
	// Notify all active streams.
	for _, s := range streams ***REMOVED***
		t.closeStream(s, ErrConnClosing, false, http2.ErrCodeNo, status.New(codes.Unavailable, ErrConnClosing.Desc), nil, false)
	***REMOVED***
	if t.statsHandler != nil ***REMOVED***
		connEnd := &stats.ConnEnd***REMOVED***
			Client: true,
		***REMOVED***
		t.statsHandler.HandleConn(t.ctx, connEnd)
	***REMOVED***
	return err
***REMOVED***

// GracefulClose sets the state to draining, which prevents new streams from
// being created and causes the transport to be closed when the last active
// stream is closed.  If there are no active streams, the transport is closed
// immediately.  This does nothing if the transport is already draining or
// closing.
func (t *http2Client) GracefulClose() ***REMOVED***
	t.mu.Lock()
	// Make sure we move to draining only from active.
	if t.state == draining || t.state == closing ***REMOVED***
		t.mu.Unlock()
		return
	***REMOVED***
	t.state = draining
	active := len(t.activeStreams)
	t.mu.Unlock()
	if active == 0 ***REMOVED***
		t.Close()
		return
	***REMOVED***
	t.controlBuf.put(&incomingGoAway***REMOVED******REMOVED***)
***REMOVED***

// Write formats the data into HTTP2 data frame(s) and sends it out. The caller
// should proceed only if Write returns nil.
func (t *http2Client) Write(s *Stream, hdr []byte, data []byte, opts *Options) error ***REMOVED***
	if opts.Last ***REMOVED***
		// If it's the last message, update stream state.
		if !s.compareAndSwapState(streamActive, streamWriteDone) ***REMOVED***
			return errStreamDone
		***REMOVED***
	***REMOVED*** else if s.getState() != streamActive ***REMOVED***
		return errStreamDone
	***REMOVED***
	df := &dataFrame***REMOVED***
		streamID:  s.id,
		endStream: opts.Last,
		h:         hdr,
		d:         data,
	***REMOVED***
	if hdr != nil || data != nil ***REMOVED*** // If it's not an empty data frame, check quota.
		if err := s.wq.get(int32(len(hdr) + len(data))); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return t.controlBuf.put(df)
***REMOVED***

func (t *http2Client) getStream(f http2.Frame) *Stream ***REMOVED***
	t.mu.Lock()
	s := t.activeStreams[f.Header().StreamID]
	t.mu.Unlock()
	return s
***REMOVED***

// adjustWindow sends out extra window update over the initial window size
// of stream if the application is requesting data larger in size than
// the window.
func (t *http2Client) adjustWindow(s *Stream, n uint32) ***REMOVED***
	if w := s.fc.maybeAdjust(n); w > 0 ***REMOVED***
		t.controlBuf.put(&outgoingWindowUpdate***REMOVED***streamID: s.id, increment: w***REMOVED***)
	***REMOVED***
***REMOVED***

// updateWindow adjusts the inbound quota for the stream.
// Window updates will be sent out when the cumulative quota
// exceeds the corresponding threshold.
func (t *http2Client) updateWindow(s *Stream, n uint32) ***REMOVED***
	if w := s.fc.onRead(n); w > 0 ***REMOVED***
		t.controlBuf.put(&outgoingWindowUpdate***REMOVED***streamID: s.id, increment: w***REMOVED***)
	***REMOVED***
***REMOVED***

// updateFlowControl updates the incoming flow control windows
// for the transport and the stream based on the current bdp
// estimation.
func (t *http2Client) updateFlowControl(n uint32) ***REMOVED***
	t.mu.Lock()
	for _, s := range t.activeStreams ***REMOVED***
		s.fc.newLimit(n)
	***REMOVED***
	t.mu.Unlock()
	updateIWS := func(interface***REMOVED******REMOVED***) bool ***REMOVED***
		t.initialWindowSize = int32(n)
		return true
	***REMOVED***
	t.controlBuf.executeAndPut(updateIWS, &outgoingWindowUpdate***REMOVED***streamID: 0, increment: t.fc.newLimit(n)***REMOVED***)
	t.controlBuf.put(&outgoingSettings***REMOVED***
		ss: []http2.Setting***REMOVED***
			***REMOVED***
				ID:  http2.SettingInitialWindowSize,
				Val: n,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***)
***REMOVED***

func (t *http2Client) handleData(f *http2.DataFrame) ***REMOVED***
	size := f.Header().Length
	var sendBDPPing bool
	if t.bdpEst != nil ***REMOVED***
		sendBDPPing = t.bdpEst.add(size)
	***REMOVED***
	// Decouple connection's flow control from application's read.
	// An update on connection's flow control should not depend on
	// whether user application has read the data or not. Such a
	// restriction is already imposed on the stream's flow control,
	// and therefore the sender will be blocked anyways.
	// Decoupling the connection flow control will prevent other
	// active(fast) streams from starving in presence of slow or
	// inactive streams.
	//
	if w := t.fc.onData(size); w > 0 ***REMOVED***
		t.controlBuf.put(&outgoingWindowUpdate***REMOVED***
			streamID:  0,
			increment: w,
		***REMOVED***)
	***REMOVED***
	if sendBDPPing ***REMOVED***
		// Avoid excessive ping detection (e.g. in an L7 proxy)
		// by sending a window update prior to the BDP ping.

		if w := t.fc.reset(); w > 0 ***REMOVED***
			t.controlBuf.put(&outgoingWindowUpdate***REMOVED***
				streamID:  0,
				increment: w,
			***REMOVED***)
		***REMOVED***

		t.controlBuf.put(bdpPing)
	***REMOVED***
	// Select the right stream to dispatch.
	s := t.getStream(f)
	if s == nil ***REMOVED***
		return
	***REMOVED***
	if size > 0 ***REMOVED***
		if err := s.fc.onData(size); err != nil ***REMOVED***
			t.closeStream(s, io.EOF, true, http2.ErrCodeFlowControl, status.New(codes.Internal, err.Error()), nil, false)
			return
		***REMOVED***
		if f.Header().Flags.Has(http2.FlagDataPadded) ***REMOVED***
			if w := s.fc.onRead(size - uint32(len(f.Data()))); w > 0 ***REMOVED***
				t.controlBuf.put(&outgoingWindowUpdate***REMOVED***s.id, w***REMOVED***)
			***REMOVED***
		***REMOVED***
		// TODO(bradfitz, zhaoq): A copy is required here because there is no
		// guarantee f.Data() is consumed before the arrival of next frame.
		// Can this copy be eliminated?
		if len(f.Data()) > 0 ***REMOVED***
			buffer := t.bufferPool.get()
			buffer.Reset()
			buffer.Write(f.Data())
			s.write(recvMsg***REMOVED***buffer: buffer***REMOVED***)
		***REMOVED***
	***REMOVED***
	// The server has closed the stream without sending trailers.  Record that
	// the read direction is closed, and set the status appropriately.
	if f.FrameHeader.Flags.Has(http2.FlagDataEndStream) ***REMOVED***
		t.closeStream(s, io.EOF, false, http2.ErrCodeNo, status.New(codes.Internal, "server closed the stream without sending trailers"), nil, true)
	***REMOVED***
***REMOVED***

func (t *http2Client) handleRSTStream(f *http2.RSTStreamFrame) ***REMOVED***
	s := t.getStream(f)
	if s == nil ***REMOVED***
		return
	***REMOVED***
	if f.ErrCode == http2.ErrCodeRefusedStream ***REMOVED***
		// The stream was unprocessed by the server.
		atomic.StoreUint32(&s.unprocessed, 1)
	***REMOVED***
	statusCode, ok := http2ErrConvTab[f.ErrCode]
	if !ok ***REMOVED***
		if logger.V(logLevel) ***REMOVED***
			logger.Warningf("transport: http2Client.handleRSTStream found no mapped gRPC status for the received http2 error %v", f.ErrCode)
		***REMOVED***
		statusCode = codes.Unknown
	***REMOVED***
	if statusCode == codes.Canceled ***REMOVED***
		if d, ok := s.ctx.Deadline(); ok && !d.After(time.Now()) ***REMOVED***
			// Our deadline was already exceeded, and that was likely the cause
			// of this cancelation.  Alter the status code accordingly.
			statusCode = codes.DeadlineExceeded
		***REMOVED***
	***REMOVED***
	t.closeStream(s, io.EOF, false, http2.ErrCodeNo, status.Newf(statusCode, "stream terminated by RST_STREAM with error code: %v", f.ErrCode), nil, false)
***REMOVED***

func (t *http2Client) handleSettings(f *http2.SettingsFrame, isFirst bool) ***REMOVED***
	if f.IsAck() ***REMOVED***
		return
	***REMOVED***
	var maxStreams *uint32
	var ss []http2.Setting
	var updateFuncs []func()
	f.ForeachSetting(func(s http2.Setting) error ***REMOVED***
		switch s.ID ***REMOVED***
		case http2.SettingMaxConcurrentStreams:
			maxStreams = new(uint32)
			*maxStreams = s.Val
		case http2.SettingMaxHeaderListSize:
			updateFuncs = append(updateFuncs, func() ***REMOVED***
				t.maxSendHeaderListSize = new(uint32)
				*t.maxSendHeaderListSize = s.Val
			***REMOVED***)
		default:
			ss = append(ss, s)
		***REMOVED***
		return nil
	***REMOVED***)
	if isFirst && maxStreams == nil ***REMOVED***
		maxStreams = new(uint32)
		*maxStreams = math.MaxUint32
	***REMOVED***
	sf := &incomingSettings***REMOVED***
		ss: ss,
	***REMOVED***
	if maxStreams != nil ***REMOVED***
		updateStreamQuota := func() ***REMOVED***
			delta := int64(*maxStreams) - int64(t.maxConcurrentStreams)
			t.maxConcurrentStreams = *maxStreams
			t.streamQuota += delta
			if delta > 0 && t.waitingStreams > 0 ***REMOVED***
				close(t.streamsQuotaAvailable) // wake all of them up.
				t.streamsQuotaAvailable = make(chan struct***REMOVED******REMOVED***, 1)
			***REMOVED***
		***REMOVED***
		updateFuncs = append(updateFuncs, updateStreamQuota)
	***REMOVED***
	t.controlBuf.executeAndPut(func(interface***REMOVED******REMOVED***) bool ***REMOVED***
		for _, f := range updateFuncs ***REMOVED***
			f()
		***REMOVED***
		return true
	***REMOVED***, sf)
***REMOVED***

func (t *http2Client) handlePing(f *http2.PingFrame) ***REMOVED***
	if f.IsAck() ***REMOVED***
		// Maybe it's a BDP ping.
		if t.bdpEst != nil ***REMOVED***
			t.bdpEst.calculate(f.Data)
		***REMOVED***
		return
	***REMOVED***
	pingAck := &ping***REMOVED***ack: true***REMOVED***
	copy(pingAck.data[:], f.Data[:])
	t.controlBuf.put(pingAck)
***REMOVED***

func (t *http2Client) handleGoAway(f *http2.GoAwayFrame) ***REMOVED***
	t.mu.Lock()
	if t.state == closing ***REMOVED***
		t.mu.Unlock()
		return
	***REMOVED***
	if f.ErrCode == http2.ErrCodeEnhanceYourCalm ***REMOVED***
		if logger.V(logLevel) ***REMOVED***
			logger.Infof("Client received GoAway with http2.ErrCodeEnhanceYourCalm.")
		***REMOVED***
	***REMOVED***
	id := f.LastStreamID
	if id > 0 && id%2 != 1 ***REMOVED***
		t.mu.Unlock()
		t.Close()
		return
	***REMOVED***
	// A client can receive multiple GoAways from the server (see
	// https://github.com/grpc/grpc-go/issues/1387).  The idea is that the first
	// GoAway will be sent with an ID of MaxInt32 and the second GoAway will be
	// sent after an RTT delay with the ID of the last stream the server will
	// process.
	//
	// Therefore, when we get the first GoAway we don't necessarily close any
	// streams. While in case of second GoAway we close all streams created after
	// the GoAwayId. This way streams that were in-flight while the GoAway from
	// server was being sent don't get killed.
	select ***REMOVED***
	case <-t.goAway: // t.goAway has been closed (i.e.,multiple GoAways).
		// If there are multiple GoAways the first one should always have an ID greater than the following ones.
		if id > t.prevGoAwayID ***REMOVED***
			t.mu.Unlock()
			t.Close()
			return
		***REMOVED***
	default:
		t.setGoAwayReason(f)
		close(t.goAway)
		t.controlBuf.put(&incomingGoAway***REMOVED******REMOVED***)
		// Notify the clientconn about the GOAWAY before we set the state to
		// draining, to allow the client to stop attempting to create streams
		// before disallowing new streams on this connection.
		t.onGoAway(t.goAwayReason)
		t.state = draining
	***REMOVED***
	// All streams with IDs greater than the GoAwayId
	// and smaller than the previous GoAway ID should be killed.
	upperLimit := t.prevGoAwayID
	if upperLimit == 0 ***REMOVED*** // This is the first GoAway Frame.
		upperLimit = math.MaxUint32 // Kill all streams after the GoAway ID.
	***REMOVED***
	for streamID, stream := range t.activeStreams ***REMOVED***
		if streamID > id && streamID <= upperLimit ***REMOVED***
			// The stream was unprocessed by the server.
			atomic.StoreUint32(&stream.unprocessed, 1)
			t.closeStream(stream, errStreamDrain, false, http2.ErrCodeNo, statusGoAway, nil, false)
		***REMOVED***
	***REMOVED***
	t.prevGoAwayID = id
	active := len(t.activeStreams)
	t.mu.Unlock()
	if active == 0 ***REMOVED***
		t.Close()
	***REMOVED***
***REMOVED***

// setGoAwayReason sets the value of t.goAwayReason based
// on the GoAway frame received.
// It expects a lock on transport's mutext to be held by
// the caller.
func (t *http2Client) setGoAwayReason(f *http2.GoAwayFrame) ***REMOVED***
	t.goAwayReason = GoAwayNoReason
	switch f.ErrCode ***REMOVED***
	case http2.ErrCodeEnhanceYourCalm:
		if string(f.DebugData()) == "too_many_pings" ***REMOVED***
			t.goAwayReason = GoAwayTooManyPings
		***REMOVED***
	***REMOVED***
***REMOVED***

func (t *http2Client) GetGoAwayReason() GoAwayReason ***REMOVED***
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.goAwayReason
***REMOVED***

func (t *http2Client) handleWindowUpdate(f *http2.WindowUpdateFrame) ***REMOVED***
	t.controlBuf.put(&incomingWindowUpdate***REMOVED***
		streamID:  f.Header().StreamID,
		increment: f.Increment,
	***REMOVED***)
***REMOVED***

// operateHeaders takes action on the decoded headers.
func (t *http2Client) operateHeaders(frame *http2.MetaHeadersFrame) ***REMOVED***
	s := t.getStream(frame)
	if s == nil ***REMOVED***
		return
	***REMOVED***
	endStream := frame.StreamEnded()
	atomic.StoreUint32(&s.bytesReceived, 1)
	initialHeader := atomic.LoadUint32(&s.headerChanClosed) == 0

	if !initialHeader && !endStream ***REMOVED***
		// As specified by gRPC over HTTP2, a HEADERS frame (and associated CONTINUATION frames) can only appear at the start or end of a stream. Therefore, second HEADERS frame must have EOS bit set.
		st := status.New(codes.Internal, "a HEADERS frame cannot appear in the middle of a stream")
		t.closeStream(s, st.Err(), true, http2.ErrCodeProtocol, st, nil, false)
		return
	***REMOVED***

	state := &decodeState***REMOVED******REMOVED***
	// Initialize isGRPC value to be !initialHeader, since if a gRPC Response-Headers has already been received, then it means that the peer is speaking gRPC and we are in gRPC mode.
	state.data.isGRPC = !initialHeader
	if h2code, err := state.decodeHeader(frame); err != nil ***REMOVED***
		t.closeStream(s, err, true, h2code, status.Convert(err), nil, endStream)
		return
	***REMOVED***

	isHeader := false
	defer func() ***REMOVED***
		if t.statsHandler != nil ***REMOVED***
			if isHeader ***REMOVED***
				inHeader := &stats.InHeader***REMOVED***
					Client:      true,
					WireLength:  int(frame.Header().Length),
					Header:      s.header.Copy(),
					Compression: s.recvCompress,
				***REMOVED***
				t.statsHandler.HandleRPC(s.ctx, inHeader)
			***REMOVED*** else ***REMOVED***
				inTrailer := &stats.InTrailer***REMOVED***
					Client:     true,
					WireLength: int(frame.Header().Length),
					Trailer:    s.trailer.Copy(),
				***REMOVED***
				t.statsHandler.HandleRPC(s.ctx, inTrailer)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	// If headerChan hasn't been closed yet
	if atomic.CompareAndSwapUint32(&s.headerChanClosed, 0, 1) ***REMOVED***
		s.headerValid = true
		if !endStream ***REMOVED***
			// HEADERS frame block carries a Response-Headers.
			isHeader = true
			// These values can be set without any synchronization because
			// stream goroutine will read it only after seeing a closed
			// headerChan which we'll close after setting this.
			s.recvCompress = state.data.encoding
			if len(state.data.mdata) > 0 ***REMOVED***
				s.header = state.data.mdata
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// HEADERS frame block carries a Trailers-Only.
			s.noHeaders = true
		***REMOVED***
		close(s.headerChan)
	***REMOVED***

	if !endStream ***REMOVED***
		return
	***REMOVED***

	// if client received END_STREAM from server while stream was still active, send RST_STREAM
	rst := s.getState() == streamActive
	t.closeStream(s, io.EOF, rst, http2.ErrCodeNo, state.status(), state.data.mdata, true)
***REMOVED***

// reader runs as a separate goroutine in charge of reading data from network
// connection.
//
// TODO(zhaoq): currently one reader per transport. Investigate whether this is
// optimal.
// TODO(zhaoq): Check the validity of the incoming frame sequence.
func (t *http2Client) reader() ***REMOVED***
	defer close(t.readerDone)
	// Check the validity of server preface.
	frame, err := t.framer.fr.ReadFrame()
	if err != nil ***REMOVED***
		t.Close() // this kicks off resetTransport, so must be last before return
		return
	***REMOVED***
	t.conn.SetReadDeadline(time.Time***REMOVED******REMOVED***) // reset deadline once we get the settings frame (we didn't time out, yay!)
	if t.keepaliveEnabled ***REMOVED***
		atomic.StoreInt64(&t.lastRead, time.Now().UnixNano())
	***REMOVED***
	sf, ok := frame.(*http2.SettingsFrame)
	if !ok ***REMOVED***
		t.Close() // this kicks off resetTransport, so must be last before return
		return
	***REMOVED***
	t.onPrefaceReceipt()
	t.handleSettings(sf, true)

	// loop to keep reading incoming messages on this transport.
	for ***REMOVED***
		t.controlBuf.throttle()
		frame, err := t.framer.fr.ReadFrame()
		if t.keepaliveEnabled ***REMOVED***
			atomic.StoreInt64(&t.lastRead, time.Now().UnixNano())
		***REMOVED***
		if err != nil ***REMOVED***
			// Abort an active stream if the http2.Framer returns a
			// http2.StreamError. This can happen only if the server's response
			// is malformed http2.
			if se, ok := err.(http2.StreamError); ok ***REMOVED***
				t.mu.Lock()
				s := t.activeStreams[se.StreamID]
				t.mu.Unlock()
				if s != nil ***REMOVED***
					// use error detail to provide better err message
					code := http2ErrConvTab[se.Code]
					errorDetail := t.framer.fr.ErrorDetail()
					var msg string
					if errorDetail != nil ***REMOVED***
						msg = errorDetail.Error()
					***REMOVED*** else ***REMOVED***
						msg = "received invalid frame"
					***REMOVED***
					t.closeStream(s, status.Error(code, msg), true, http2.ErrCodeProtocol, status.New(code, msg), nil, false)
				***REMOVED***
				continue
			***REMOVED*** else ***REMOVED***
				// Transport error.
				t.Close()
				return
			***REMOVED***
		***REMOVED***
		switch frame := frame.(type) ***REMOVED***
		case *http2.MetaHeadersFrame:
			t.operateHeaders(frame)
		case *http2.DataFrame:
			t.handleData(frame)
		case *http2.RSTStreamFrame:
			t.handleRSTStream(frame)
		case *http2.SettingsFrame:
			t.handleSettings(frame, false)
		case *http2.PingFrame:
			t.handlePing(frame)
		case *http2.GoAwayFrame:
			t.handleGoAway(frame)
		case *http2.WindowUpdateFrame:
			t.handleWindowUpdate(frame)
		default:
			if logger.V(logLevel) ***REMOVED***
				logger.Errorf("transport: http2Client.reader got unhandled frame type %v.", frame)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func minTime(a, b time.Duration) time.Duration ***REMOVED***
	if a < b ***REMOVED***
		return a
	***REMOVED***
	return b
***REMOVED***

// keepalive running in a separate goroutune makes sure the connection is alive by sending pings.
func (t *http2Client) keepalive() ***REMOVED***
	p := &ping***REMOVED***data: [8]byte***REMOVED******REMOVED******REMOVED***
	// True iff a ping has been sent, and no data has been received since then.
	outstandingPing := false
	// Amount of time remaining before which we should receive an ACK for the
	// last sent ping.
	timeoutLeft := time.Duration(0)
	// Records the last value of t.lastRead before we go block on the timer.
	// This is required to check for read activity since then.
	prevNano := time.Now().UnixNano()
	timer := time.NewTimer(t.kp.Time)
	for ***REMOVED***
		select ***REMOVED***
		case <-timer.C:
			lastRead := atomic.LoadInt64(&t.lastRead)
			if lastRead > prevNano ***REMOVED***
				// There has been read activity since the last time we were here.
				outstandingPing = false
				// Next timer should fire at kp.Time seconds from lastRead time.
				timer.Reset(time.Duration(lastRead) + t.kp.Time - time.Duration(time.Now().UnixNano()))
				prevNano = lastRead
				continue
			***REMOVED***
			if outstandingPing && timeoutLeft <= 0 ***REMOVED***
				t.Close()
				return
			***REMOVED***
			t.mu.Lock()
			if t.state == closing ***REMOVED***
				// If the transport is closing, we should exit from the
				// keepalive goroutine here. If not, we could have a race
				// between the call to Signal() from Close() and the call to
				// Wait() here, whereby the keepalive goroutine ends up
				// blocking on the condition variable which will never be
				// signalled again.
				t.mu.Unlock()
				return
			***REMOVED***
			if len(t.activeStreams) < 1 && !t.kp.PermitWithoutStream ***REMOVED***
				// If a ping was sent out previously (because there were active
				// streams at that point) which wasn't acked and its timeout
				// hadn't fired, but we got here and are about to go dormant,
				// we should make sure that we unconditionally send a ping once
				// we awaken.
				outstandingPing = false
				t.kpDormant = true
				t.kpDormancyCond.Wait()
			***REMOVED***
			t.kpDormant = false
			t.mu.Unlock()

			// We get here either because we were dormant and a new stream was
			// created which unblocked the Wait() call, or because the
			// keepalive timer expired. In both cases, we need to send a ping.
			if !outstandingPing ***REMOVED***
				if channelz.IsOn() ***REMOVED***
					atomic.AddInt64(&t.czData.kpCount, 1)
				***REMOVED***
				t.controlBuf.put(p)
				timeoutLeft = t.kp.Timeout
				outstandingPing = true
			***REMOVED***
			// The amount of time to sleep here is the minimum of kp.Time and
			// timeoutLeft. This will ensure that we wait only for kp.Time
			// before sending out the next ping (for cases where the ping is
			// acked).
			sleepDuration := minTime(t.kp.Time, timeoutLeft)
			timeoutLeft -= sleepDuration
			timer.Reset(sleepDuration)
		case <-t.ctx.Done():
			if !timer.Stop() ***REMOVED***
				<-timer.C
			***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (t *http2Client) Error() <-chan struct***REMOVED******REMOVED*** ***REMOVED***
	return t.ctx.Done()
***REMOVED***

func (t *http2Client) GoAway() <-chan struct***REMOVED******REMOVED*** ***REMOVED***
	return t.goAway
***REMOVED***

func (t *http2Client) ChannelzMetric() *channelz.SocketInternalMetric ***REMOVED***
	s := channelz.SocketInternalMetric***REMOVED***
		StreamsStarted:                  atomic.LoadInt64(&t.czData.streamsStarted),
		StreamsSucceeded:                atomic.LoadInt64(&t.czData.streamsSucceeded),
		StreamsFailed:                   atomic.LoadInt64(&t.czData.streamsFailed),
		MessagesSent:                    atomic.LoadInt64(&t.czData.msgSent),
		MessagesReceived:                atomic.LoadInt64(&t.czData.msgRecv),
		KeepAlivesSent:                  atomic.LoadInt64(&t.czData.kpCount),
		LastLocalStreamCreatedTimestamp: time.Unix(0, atomic.LoadInt64(&t.czData.lastStreamCreatedTime)),
		LastMessageSentTimestamp:        time.Unix(0, atomic.LoadInt64(&t.czData.lastMsgSentTime)),
		LastMessageReceivedTimestamp:    time.Unix(0, atomic.LoadInt64(&t.czData.lastMsgRecvTime)),
		LocalFlowControlWindow:          int64(t.fc.getSize()),
		SocketOptions:                   channelz.GetSocketOption(t.conn),
		LocalAddr:                       t.localAddr,
		RemoteAddr:                      t.remoteAddr,
		// RemoteName :
	***REMOVED***
	if au, ok := t.authInfo.(credentials.ChannelzSecurityInfo); ok ***REMOVED***
		s.Security = au.GetSecurityValue()
	***REMOVED***
	s.RemoteFlowControlWindow = t.getOutFlowWindow()
	return &s
***REMOVED***

func (t *http2Client) RemoteAddr() net.Addr ***REMOVED*** return t.remoteAddr ***REMOVED***

func (t *http2Client) IncrMsgSent() ***REMOVED***
	atomic.AddInt64(&t.czData.msgSent, 1)
	atomic.StoreInt64(&t.czData.lastMsgSentTime, time.Now().UnixNano())
***REMOVED***

func (t *http2Client) IncrMsgRecv() ***REMOVED***
	atomic.AddInt64(&t.czData.msgRecv, 1)
	atomic.StoreInt64(&t.czData.lastMsgRecvTime, time.Now().UnixNano())
***REMOVED***

func (t *http2Client) getOutFlowWindow() int64 ***REMOVED***
	resp := make(chan uint32, 1)
	timer := time.NewTimer(time.Second)
	defer timer.Stop()
	t.controlBuf.put(&outFlowControlSizeRequest***REMOVED***resp***REMOVED***)
	select ***REMOVED***
	case sz := <-resp:
		return int64(sz)
	case <-t.ctxDone:
		return -1
	case <-timer.C:
		return -2
	***REMOVED***
***REMOVED***
