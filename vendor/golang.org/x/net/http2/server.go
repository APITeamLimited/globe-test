// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// TODO: turn off the serve goroutine when idle, so
// an idle conn only has the readFrames goroutine active. (which could
// also be optimized probably to pin less memory in crypto/tls). This
// would involve tracking when the serve goroutine is active (atomic
// int32 read/CAS probably?) and starting it up when frames arrive,
// and shutting it down when all handlers exit. the occasional PING
// packets could use time.AfterFunc to call sc.wakeStartServeLoop()
// (which is a no-op if already running) and then queue the PING write
// as normal. The serve loop would then exit in most cases (if no
// Handlers running) and not be woken up again until the PING packet
// returns.

// TODO (maybe): add a mechanism for Handlers to going into
// half-closed-local mode (rw.(io.Closer) test?) but not exit their
// handler, and continue to be able to read from the
// Request.Body. This would be a somewhat semantic change from HTTP/1
// (or at least what we expose in net/http), so I'd probably want to
// add it there too. For now, this package says that returning from
// the Handler ServeHTTP function means you're both done reading and
// done writing, without a way to stop just one or the other.

package http2

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/http/httpguts"
	"golang.org/x/net/http2/hpack"
)

const (
	prefaceTimeout         = 10 * time.Second
	firstSettingsTimeout   = 2 * time.Second // should be in-flight with preface anyway
	handlerChunkWriteSize  = 4 << 10
	defaultMaxStreams      = 250 // TODO: make this 100 as the GFE seems to?
	maxQueuedControlFrames = 10000
)

var (
	errClientDisconnected = errors.New("client disconnected")
	errClosedBody         = errors.New("body closed by handler")
	errHandlerComplete    = errors.New("http2: request body closed due to handler exiting")
	errStreamClosed       = errors.New("http2: stream closed")
)

var responseWriterStatePool = sync.Pool***REMOVED***
	New: func() interface***REMOVED******REMOVED*** ***REMOVED***
		rws := &responseWriterState***REMOVED******REMOVED***
		rws.bw = bufio.NewWriterSize(chunkWriter***REMOVED***rws***REMOVED***, handlerChunkWriteSize)
		return rws
	***REMOVED***,
***REMOVED***

// Test hooks.
var (
	testHookOnConn        func()
	testHookGetServerConn func(*serverConn)
	testHookOnPanicMu     *sync.Mutex // nil except in tests
	testHookOnPanic       func(sc *serverConn, panicVal interface***REMOVED******REMOVED***) (rePanic bool)
)

// Server is an HTTP/2 server.
type Server struct ***REMOVED***
	// MaxHandlers limits the number of http.Handler ServeHTTP goroutines
	// which may run at a time over all connections.
	// Negative or zero no limit.
	// TODO: implement
	MaxHandlers int

	// MaxConcurrentStreams optionally specifies the number of
	// concurrent streams that each client may have open at a
	// time. This is unrelated to the number of http.Handler goroutines
	// which may be active globally, which is MaxHandlers.
	// If zero, MaxConcurrentStreams defaults to at least 100, per
	// the HTTP/2 spec's recommendations.
	MaxConcurrentStreams uint32

	// MaxReadFrameSize optionally specifies the largest frame
	// this server is willing to read. A valid value is between
	// 16k and 16M, inclusive. If zero or otherwise invalid, a
	// default value is used.
	MaxReadFrameSize uint32

	// PermitProhibitedCipherSuites, if true, permits the use of
	// cipher suites prohibited by the HTTP/2 spec.
	PermitProhibitedCipherSuites bool

	// IdleTimeout specifies how long until idle clients should be
	// closed with a GOAWAY frame. PING frames are not considered
	// activity for the purposes of IdleTimeout.
	IdleTimeout time.Duration

	// MaxUploadBufferPerConnection is the size of the initial flow
	// control window for each connections. The HTTP/2 spec does not
	// allow this to be smaller than 65535 or larger than 2^32-1.
	// If the value is outside this range, a default value will be
	// used instead.
	MaxUploadBufferPerConnection int32

	// MaxUploadBufferPerStream is the size of the initial flow control
	// window for each stream. The HTTP/2 spec does not allow this to
	// be larger than 2^32-1. If the value is zero or larger than the
	// maximum, a default value will be used instead.
	MaxUploadBufferPerStream int32

	// NewWriteScheduler constructs a write scheduler for a connection.
	// If nil, a default scheduler is chosen.
	NewWriteScheduler func() WriteScheduler

	// Internal state. This is a pointer (rather than embedded directly)
	// so that we don't embed a Mutex in this struct, which will make the
	// struct non-copyable, which might break some callers.
	state *serverInternalState
***REMOVED***

func (s *Server) initialConnRecvWindowSize() int32 ***REMOVED***
	if s.MaxUploadBufferPerConnection > initialWindowSize ***REMOVED***
		return s.MaxUploadBufferPerConnection
	***REMOVED***
	return 1 << 20
***REMOVED***

func (s *Server) initialStreamRecvWindowSize() int32 ***REMOVED***
	if s.MaxUploadBufferPerStream > 0 ***REMOVED***
		return s.MaxUploadBufferPerStream
	***REMOVED***
	return 1 << 20
***REMOVED***

func (s *Server) maxReadFrameSize() uint32 ***REMOVED***
	if v := s.MaxReadFrameSize; v >= minMaxFrameSize && v <= maxFrameSize ***REMOVED***
		return v
	***REMOVED***
	return defaultMaxReadFrameSize
***REMOVED***

func (s *Server) maxConcurrentStreams() uint32 ***REMOVED***
	if v := s.MaxConcurrentStreams; v > 0 ***REMOVED***
		return v
	***REMOVED***
	return defaultMaxStreams
***REMOVED***

// maxQueuedControlFrames is the maximum number of control frames like
// SETTINGS, PING and RST_STREAM that will be queued for writing before
// the connection is closed to prevent memory exhaustion attacks.
func (s *Server) maxQueuedControlFrames() int ***REMOVED***
	// TODO: if anybody asks, add a Server field, and remember to define the
	// behavior of negative values.
	return maxQueuedControlFrames
***REMOVED***

type serverInternalState struct ***REMOVED***
	mu          sync.Mutex
	activeConns map[*serverConn]struct***REMOVED******REMOVED***
***REMOVED***

func (s *serverInternalState) registerConn(sc *serverConn) ***REMOVED***
	if s == nil ***REMOVED***
		return // if the Server was used without calling ConfigureServer
	***REMOVED***
	s.mu.Lock()
	s.activeConns[sc] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	s.mu.Unlock()
***REMOVED***

func (s *serverInternalState) unregisterConn(sc *serverConn) ***REMOVED***
	if s == nil ***REMOVED***
		return // if the Server was used without calling ConfigureServer
	***REMOVED***
	s.mu.Lock()
	delete(s.activeConns, sc)
	s.mu.Unlock()
***REMOVED***

func (s *serverInternalState) startGracefulShutdown() ***REMOVED***
	if s == nil ***REMOVED***
		return // if the Server was used without calling ConfigureServer
	***REMOVED***
	s.mu.Lock()
	for sc := range s.activeConns ***REMOVED***
		sc.startGracefulShutdown()
	***REMOVED***
	s.mu.Unlock()
***REMOVED***

// ConfigureServer adds HTTP/2 support to a net/http Server.
//
// The configuration conf may be nil.
//
// ConfigureServer must be called before s begins serving.
func ConfigureServer(s *http.Server, conf *Server) error ***REMOVED***
	if s == nil ***REMOVED***
		panic("nil *http.Server")
	***REMOVED***
	if conf == nil ***REMOVED***
		conf = new(Server)
	***REMOVED***
	conf.state = &serverInternalState***REMOVED***activeConns: make(map[*serverConn]struct***REMOVED******REMOVED***)***REMOVED***
	if h1, h2 := s, conf; h2.IdleTimeout == 0 ***REMOVED***
		if h1.IdleTimeout != 0 ***REMOVED***
			h2.IdleTimeout = h1.IdleTimeout
		***REMOVED*** else ***REMOVED***
			h2.IdleTimeout = h1.ReadTimeout
		***REMOVED***
	***REMOVED***
	s.RegisterOnShutdown(conf.state.startGracefulShutdown)

	if s.TLSConfig == nil ***REMOVED***
		s.TLSConfig = new(tls.Config)
	***REMOVED*** else if s.TLSConfig.CipherSuites != nil ***REMOVED***
		// If they already provided a CipherSuite list, return
		// an error if it has a bad order or is missing
		// ECDHE_RSA_WITH_AES_128_GCM_SHA256 or ECDHE_ECDSA_WITH_AES_128_GCM_SHA256.
		haveRequired := false
		sawBad := false
		for i, cs := range s.TLSConfig.CipherSuites ***REMOVED***
			switch cs ***REMOVED***
			case tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				// Alternative MTI cipher to not discourage ECDSA-only servers.
				// See http://golang.org/cl/30721 for further information.
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256:
				haveRequired = true
			***REMOVED***
			if isBadCipher(cs) ***REMOVED***
				sawBad = true
			***REMOVED*** else if sawBad ***REMOVED***
				return fmt.Errorf("http2: TLSConfig.CipherSuites index %d contains an HTTP/2-approved cipher suite (%#04x), but it comes after unapproved cipher suites. With this configuration, clients that don't support previous, approved cipher suites may be given an unapproved one and reject the connection.", i, cs)
			***REMOVED***
		***REMOVED***
		if !haveRequired ***REMOVED***
			return fmt.Errorf("http2: TLSConfig.CipherSuites is missing an HTTP/2-required AES_128_GCM_SHA256 cipher (need at least one of TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 or TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256).")
		***REMOVED***
	***REMOVED***

	// Note: not setting MinVersion to tls.VersionTLS12,
	// as we don't want to interfere with HTTP/1.1 traffic
	// on the user's server. We enforce TLS 1.2 later once
	// we accept a connection. Ideally this should be done
	// during next-proto selection, but using TLS <1.2 with
	// HTTP/2 is still the client's bug.

	s.TLSConfig.PreferServerCipherSuites = true

	haveNPN := false
	for _, p := range s.TLSConfig.NextProtos ***REMOVED***
		if p == NextProtoTLS ***REMOVED***
			haveNPN = true
			break
		***REMOVED***
	***REMOVED***
	if !haveNPN ***REMOVED***
		s.TLSConfig.NextProtos = append(s.TLSConfig.NextProtos, NextProtoTLS)
	***REMOVED***

	if s.TLSNextProto == nil ***REMOVED***
		s.TLSNextProto = map[string]func(*http.Server, *tls.Conn, http.Handler)***REMOVED******REMOVED***
	***REMOVED***
	protoHandler := func(hs *http.Server, c *tls.Conn, h http.Handler) ***REMOVED***
		if testHookOnConn != nil ***REMOVED***
			testHookOnConn()
		***REMOVED***
		// The TLSNextProto interface predates contexts, so
		// the net/http package passes down its per-connection
		// base context via an exported but unadvertised
		// method on the Handler. This is for internal
		// net/http<=>http2 use only.
		var ctx context.Context
		type baseContexter interface ***REMOVED***
			BaseContext() context.Context
		***REMOVED***
		if bc, ok := h.(baseContexter); ok ***REMOVED***
			ctx = bc.BaseContext()
		***REMOVED***
		conf.ServeConn(c, &ServeConnOpts***REMOVED***
			Context:    ctx,
			Handler:    h,
			BaseConfig: hs,
		***REMOVED***)
	***REMOVED***
	s.TLSNextProto[NextProtoTLS] = protoHandler
	return nil
***REMOVED***

// ServeConnOpts are options for the Server.ServeConn method.
type ServeConnOpts struct ***REMOVED***
	// Context is the base context to use.
	// If nil, context.Background is used.
	Context context.Context

	// BaseConfig optionally sets the base configuration
	// for values. If nil, defaults are used.
	BaseConfig *http.Server

	// Handler specifies which handler to use for processing
	// requests. If nil, BaseConfig.Handler is used. If BaseConfig
	// or BaseConfig.Handler is nil, http.DefaultServeMux is used.
	Handler http.Handler
***REMOVED***

func (o *ServeConnOpts) context() context.Context ***REMOVED***
	if o != nil && o.Context != nil ***REMOVED***
		return o.Context
	***REMOVED***
	return context.Background()
***REMOVED***

func (o *ServeConnOpts) baseConfig() *http.Server ***REMOVED***
	if o != nil && o.BaseConfig != nil ***REMOVED***
		return o.BaseConfig
	***REMOVED***
	return new(http.Server)
***REMOVED***

func (o *ServeConnOpts) handler() http.Handler ***REMOVED***
	if o != nil ***REMOVED***
		if o.Handler != nil ***REMOVED***
			return o.Handler
		***REMOVED***
		if o.BaseConfig != nil && o.BaseConfig.Handler != nil ***REMOVED***
			return o.BaseConfig.Handler
		***REMOVED***
	***REMOVED***
	return http.DefaultServeMux
***REMOVED***

// ServeConn serves HTTP/2 requests on the provided connection and
// blocks until the connection is no longer readable.
//
// ServeConn starts speaking HTTP/2 assuming that c has not had any
// reads or writes. It writes its initial settings frame and expects
// to be able to read the preface and settings frame from the
// client. If c has a ConnectionState method like a *tls.Conn, the
// ConnectionState is used to verify the TLS ciphersuite and to set
// the Request.TLS field in Handlers.
//
// ServeConn does not support h2c by itself. Any h2c support must be
// implemented in terms of providing a suitably-behaving net.Conn.
//
// The opts parameter is optional. If nil, default values are used.
func (s *Server) ServeConn(c net.Conn, opts *ServeConnOpts) ***REMOVED***
	baseCtx, cancel := serverConnBaseContext(c, opts)
	defer cancel()

	sc := &serverConn***REMOVED***
		srv:                         s,
		hs:                          opts.baseConfig(),
		conn:                        c,
		baseCtx:                     baseCtx,
		remoteAddrStr:               c.RemoteAddr().String(),
		bw:                          newBufferedWriter(c),
		handler:                     opts.handler(),
		streams:                     make(map[uint32]*stream),
		readFrameCh:                 make(chan readFrameResult),
		wantWriteFrameCh:            make(chan FrameWriteRequest, 8),
		serveMsgCh:                  make(chan interface***REMOVED******REMOVED***, 8),
		wroteFrameCh:                make(chan frameWriteResult, 1), // buffered; one send in writeFrameAsync
		bodyReadCh:                  make(chan bodyReadMsg),         // buffering doesn't matter either way
		doneServing:                 make(chan struct***REMOVED******REMOVED***),
		clientMaxStreams:            math.MaxUint32, // Section 6.5.2: "Initially, there is no limit to this value"
		advMaxStreams:               s.maxConcurrentStreams(),
		initialStreamSendWindowSize: initialWindowSize,
		maxFrameSize:                initialMaxFrameSize,
		headerTableSize:             initialHeaderTableSize,
		serveG:                      newGoroutineLock(),
		pushEnabled:                 true,
	***REMOVED***

	s.state.registerConn(sc)
	defer s.state.unregisterConn(sc)

	// The net/http package sets the write deadline from the
	// http.Server.WriteTimeout during the TLS handshake, but then
	// passes the connection off to us with the deadline already set.
	// Write deadlines are set per stream in serverConn.newStream.
	// Disarm the net.Conn write deadline here.
	if sc.hs.WriteTimeout != 0 ***REMOVED***
		sc.conn.SetWriteDeadline(time.Time***REMOVED******REMOVED***)
	***REMOVED***

	if s.NewWriteScheduler != nil ***REMOVED***
		sc.writeSched = s.NewWriteScheduler()
	***REMOVED*** else ***REMOVED***
		sc.writeSched = NewRandomWriteScheduler()
	***REMOVED***

	// These start at the RFC-specified defaults. If there is a higher
	// configured value for inflow, that will be updated when we send a
	// WINDOW_UPDATE shortly after sending SETTINGS.
	sc.flow.add(initialWindowSize)
	sc.inflow.add(initialWindowSize)
	sc.hpackEncoder = hpack.NewEncoder(&sc.headerWriteBuf)

	fr := NewFramer(sc.bw, c)
	fr.ReadMetaHeaders = hpack.NewDecoder(initialHeaderTableSize, nil)
	fr.MaxHeaderListSize = sc.maxHeaderListSize()
	fr.SetMaxReadFrameSize(s.maxReadFrameSize())
	sc.framer = fr

	if tc, ok := c.(connectionStater); ok ***REMOVED***
		sc.tlsState = new(tls.ConnectionState)
		*sc.tlsState = tc.ConnectionState()
		// 9.2 Use of TLS Features
		// An implementation of HTTP/2 over TLS MUST use TLS
		// 1.2 or higher with the restrictions on feature set
		// and cipher suite described in this section. Due to
		// implementation limitations, it might not be
		// possible to fail TLS negotiation. An endpoint MUST
		// immediately terminate an HTTP/2 connection that
		// does not meet the TLS requirements described in
		// this section with a connection error (Section
		// 5.4.1) of type INADEQUATE_SECURITY.
		if sc.tlsState.Version < tls.VersionTLS12 ***REMOVED***
			sc.rejectConn(ErrCodeInadequateSecurity, "TLS version too low")
			return
		***REMOVED***

		if sc.tlsState.ServerName == "" ***REMOVED***
			// Client must use SNI, but we don't enforce that anymore,
			// since it was causing problems when connecting to bare IP
			// addresses during development.
			//
			// TODO: optionally enforce? Or enforce at the time we receive
			// a new request, and verify the ServerName matches the :authority?
			// But that precludes proxy situations, perhaps.
			//
			// So for now, do nothing here again.
		***REMOVED***

		if !s.PermitProhibitedCipherSuites && isBadCipher(sc.tlsState.CipherSuite) ***REMOVED***
			// "Endpoints MAY choose to generate a connection error
			// (Section 5.4.1) of type INADEQUATE_SECURITY if one of
			// the prohibited cipher suites are negotiated."
			//
			// We choose that. In my opinion, the spec is weak
			// here. It also says both parties must support at least
			// TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 so there's no
			// excuses here. If we really must, we could allow an
			// "AllowInsecureWeakCiphers" option on the server later.
			// Let's see how it plays out first.
			sc.rejectConn(ErrCodeInadequateSecurity, fmt.Sprintf("Prohibited TLS 1.2 Cipher Suite: %x", sc.tlsState.CipherSuite))
			return
		***REMOVED***
	***REMOVED***

	if hook := testHookGetServerConn; hook != nil ***REMOVED***
		hook(sc)
	***REMOVED***
	sc.serve()
***REMOVED***

func serverConnBaseContext(c net.Conn, opts *ServeConnOpts) (ctx context.Context, cancel func()) ***REMOVED***
	ctx, cancel = context.WithCancel(opts.context())
	ctx = context.WithValue(ctx, http.LocalAddrContextKey, c.LocalAddr())
	if hs := opts.baseConfig(); hs != nil ***REMOVED***
		ctx = context.WithValue(ctx, http.ServerContextKey, hs)
	***REMOVED***
	return
***REMOVED***

func (sc *serverConn) rejectConn(err ErrCode, debug string) ***REMOVED***
	sc.vlogf("http2: server rejecting conn: %v, %s", err, debug)
	// ignoring errors. hanging up anyway.
	sc.framer.WriteGoAway(0, err, []byte(debug))
	sc.bw.Flush()
	sc.conn.Close()
***REMOVED***

type serverConn struct ***REMOVED***
	// Immutable:
	srv              *Server
	hs               *http.Server
	conn             net.Conn
	bw               *bufferedWriter // writing to conn
	handler          http.Handler
	baseCtx          context.Context
	framer           *Framer
	doneServing      chan struct***REMOVED******REMOVED***          // closed when serverConn.serve ends
	readFrameCh      chan readFrameResult   // written by serverConn.readFrames
	wantWriteFrameCh chan FrameWriteRequest // from handlers -> serve
	wroteFrameCh     chan frameWriteResult  // from writeFrameAsync -> serve, tickles more frame writes
	bodyReadCh       chan bodyReadMsg       // from handlers -> serve
	serveMsgCh       chan interface***REMOVED******REMOVED***       // misc messages & code to send to / run on the serve loop
	flow             flow                   // conn-wide (not stream-specific) outbound flow control
	inflow           flow                   // conn-wide inbound flow control
	tlsState         *tls.ConnectionState   // shared by all handlers, like net/http
	remoteAddrStr    string
	writeSched       WriteScheduler

	// Everything following is owned by the serve loop; use serveG.check():
	serveG                      goroutineLock // used to verify funcs are on serve()
	pushEnabled                 bool
	sawFirstSettings            bool // got the initial SETTINGS frame after the preface
	needToSendSettingsAck       bool
	unackedSettings             int    // how many SETTINGS have we sent without ACKs?
	queuedControlFrames         int    // control frames in the writeSched queue
	clientMaxStreams            uint32 // SETTINGS_MAX_CONCURRENT_STREAMS from client (our PUSH_PROMISE limit)
	advMaxStreams               uint32 // our SETTINGS_MAX_CONCURRENT_STREAMS advertised the client
	curClientStreams            uint32 // number of open streams initiated by the client
	curPushedStreams            uint32 // number of open streams initiated by server push
	maxClientStreamID           uint32 // max ever seen from client (odd), or 0 if there have been no client requests
	maxPushPromiseID            uint32 // ID of the last push promise (even), or 0 if there have been no pushes
	streams                     map[uint32]*stream
	initialStreamSendWindowSize int32
	maxFrameSize                int32
	headerTableSize             uint32
	peerMaxHeaderListSize       uint32            // zero means unknown (default)
	canonHeader                 map[string]string // http2-lower-case -> Go-Canonical-Case
	writingFrame                bool              // started writing a frame (on serve goroutine or separate)
	writingFrameAsync           bool              // started a frame on its own goroutine but haven't heard back on wroteFrameCh
	needsFrameFlush             bool              // last frame write wasn't a flush
	inGoAway                    bool              // we've started to or sent GOAWAY
	inFrameScheduleLoop         bool              // whether we're in the scheduleFrameWrite loop
	needToSendGoAway            bool              // we need to schedule a GOAWAY frame write
	goAwayCode                  ErrCode
	shutdownTimer               *time.Timer // nil until used
	idleTimer                   *time.Timer // nil if unused

	// Owned by the writeFrameAsync goroutine:
	headerWriteBuf bytes.Buffer
	hpackEncoder   *hpack.Encoder

	// Used by startGracefulShutdown.
	shutdownOnce sync.Once
***REMOVED***

func (sc *serverConn) maxHeaderListSize() uint32 ***REMOVED***
	n := sc.hs.MaxHeaderBytes
	if n <= 0 ***REMOVED***
		n = http.DefaultMaxHeaderBytes
	***REMOVED***
	// http2's count is in a slightly different unit and includes 32 bytes per pair.
	// So, take the net/http.Server value and pad it up a bit, assuming 10 headers.
	const perFieldOverhead = 32 // per http2 spec
	const typicalHeaders = 10   // conservative
	return uint32(n + typicalHeaders*perFieldOverhead)
***REMOVED***

func (sc *serverConn) curOpenStreams() uint32 ***REMOVED***
	sc.serveG.check()
	return sc.curClientStreams + sc.curPushedStreams
***REMOVED***

// stream represents a stream. This is the minimal metadata needed by
// the serve goroutine. Most of the actual stream state is owned by
// the http.Handler's goroutine in the responseWriter. Because the
// responseWriter's responseWriterState is recycled at the end of a
// handler, this struct intentionally has no pointer to the
// *responseWriter***REMOVED***,State***REMOVED*** itself, as the Handler ending nils out the
// responseWriter's state field.
type stream struct ***REMOVED***
	// immutable:
	sc        *serverConn
	id        uint32
	body      *pipe       // non-nil if expecting DATA frames
	cw        closeWaiter // closed wait stream transitions to closed state
	ctx       context.Context
	cancelCtx func()

	// owned by serverConn's serve loop:
	bodyBytes        int64 // body bytes seen so far
	declBodyBytes    int64 // or -1 if undeclared
	flow             flow  // limits writing from Handler to client
	inflow           flow  // what the client is allowed to POST/etc to us
	state            streamState
	resetQueued      bool        // RST_STREAM queued for write; set by sc.resetStream
	gotTrailerHeader bool        // HEADER frame for trailers was seen
	wroteHeaders     bool        // whether we wrote headers (not status 100)
	writeDeadline    *time.Timer // nil if unused

	trailer    http.Header // accumulated trailers
	reqTrailer http.Header // handler's Request.Trailer
***REMOVED***

func (sc *serverConn) Framer() *Framer  ***REMOVED*** return sc.framer ***REMOVED***
func (sc *serverConn) CloseConn() error ***REMOVED*** return sc.conn.Close() ***REMOVED***
func (sc *serverConn) Flush() error     ***REMOVED*** return sc.bw.Flush() ***REMOVED***
func (sc *serverConn) HeaderEncoder() (*hpack.Encoder, *bytes.Buffer) ***REMOVED***
	return sc.hpackEncoder, &sc.headerWriteBuf
***REMOVED***

func (sc *serverConn) state(streamID uint32) (streamState, *stream) ***REMOVED***
	sc.serveG.check()
	// http://tools.ietf.org/html/rfc7540#section-5.1
	if st, ok := sc.streams[streamID]; ok ***REMOVED***
		return st.state, st
	***REMOVED***
	// "The first use of a new stream identifier implicitly closes all
	// streams in the "idle" state that might have been initiated by
	// that peer with a lower-valued stream identifier. For example, if
	// a client sends a HEADERS frame on stream 7 without ever sending a
	// frame on stream 5, then stream 5 transitions to the "closed"
	// state when the first frame for stream 7 is sent or received."
	if streamID%2 == 1 ***REMOVED***
		if streamID <= sc.maxClientStreamID ***REMOVED***
			return stateClosed, nil
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if streamID <= sc.maxPushPromiseID ***REMOVED***
			return stateClosed, nil
		***REMOVED***
	***REMOVED***
	return stateIdle, nil
***REMOVED***

// setConnState calls the net/http ConnState hook for this connection, if configured.
// Note that the net/http package does StateNew and StateClosed for us.
// There is currently no plan for StateHijacked or hijacking HTTP/2 connections.
func (sc *serverConn) setConnState(state http.ConnState) ***REMOVED***
	if sc.hs.ConnState != nil ***REMOVED***
		sc.hs.ConnState(sc.conn, state)
	***REMOVED***
***REMOVED***

func (sc *serverConn) vlogf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if VerboseLogs ***REMOVED***
		sc.logf(format, args...)
	***REMOVED***
***REMOVED***

func (sc *serverConn) logf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if lg := sc.hs.ErrorLog; lg != nil ***REMOVED***
		lg.Printf(format, args...)
	***REMOVED*** else ***REMOVED***
		log.Printf(format, args...)
	***REMOVED***
***REMOVED***

// errno returns v's underlying uintptr, else 0.
//
// TODO: remove this helper function once http2 can use build
// tags. See comment in isClosedConnError.
func errno(v error) uintptr ***REMOVED***
	if rv := reflect.ValueOf(v); rv.Kind() == reflect.Uintptr ***REMOVED***
		return uintptr(rv.Uint())
	***REMOVED***
	return 0
***REMOVED***

// isClosedConnError reports whether err is an error from use of a closed
// network connection.
func isClosedConnError(err error) bool ***REMOVED***
	if err == nil ***REMOVED***
		return false
	***REMOVED***

	// TODO: remove this string search and be more like the Windows
	// case below. That might involve modifying the standard library
	// to return better error types.
	str := err.Error()
	if strings.Contains(str, "use of closed network connection") ***REMOVED***
		return true
	***REMOVED***

	// TODO(bradfitz): x/tools/cmd/bundle doesn't really support
	// build tags, so I can't make an http2_windows.go file with
	// Windows-specific stuff. Fix that and move this, once we
	// have a way to bundle this into std's net/http somehow.
	if runtime.GOOS == "windows" ***REMOVED***
		if oe, ok := err.(*net.OpError); ok && oe.Op == "read" ***REMOVED***
			if se, ok := oe.Err.(*os.SyscallError); ok && se.Syscall == "wsarecv" ***REMOVED***
				const WSAECONNABORTED = 10053
				const WSAECONNRESET = 10054
				if n := errno(se.Err); n == WSAECONNRESET || n == WSAECONNABORTED ***REMOVED***
					return true
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func (sc *serverConn) condlogf(err error, format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if err == nil ***REMOVED***
		return
	***REMOVED***
	if err == io.EOF || err == io.ErrUnexpectedEOF || isClosedConnError(err) || err == errPrefaceTimeout ***REMOVED***
		// Boring, expected errors.
		sc.vlogf(format, args...)
	***REMOVED*** else ***REMOVED***
		sc.logf(format, args...)
	***REMOVED***
***REMOVED***

func (sc *serverConn) canonicalHeader(v string) string ***REMOVED***
	sc.serveG.check()
	buildCommonHeaderMapsOnce()
	cv, ok := commonCanonHeader[v]
	if ok ***REMOVED***
		return cv
	***REMOVED***
	cv, ok = sc.canonHeader[v]
	if ok ***REMOVED***
		return cv
	***REMOVED***
	if sc.canonHeader == nil ***REMOVED***
		sc.canonHeader = make(map[string]string)
	***REMOVED***
	cv = http.CanonicalHeaderKey(v)
	sc.canonHeader[v] = cv
	return cv
***REMOVED***

type readFrameResult struct ***REMOVED***
	f   Frame // valid until readMore is called
	err error

	// readMore should be called once the consumer no longer needs or
	// retains f. After readMore, f is invalid and more frames can be
	// read.
	readMore func()
***REMOVED***

// readFrames is the loop that reads incoming frames.
// It takes care to only read one frame at a time, blocking until the
// consumer is done with the frame.
// It's run on its own goroutine.
func (sc *serverConn) readFrames() ***REMOVED***
	gate := make(gate)
	gateDone := gate.Done
	for ***REMOVED***
		f, err := sc.framer.ReadFrame()
		select ***REMOVED***
		case sc.readFrameCh <- readFrameResult***REMOVED***f, err, gateDone***REMOVED***:
		case <-sc.doneServing:
			return
		***REMOVED***
		select ***REMOVED***
		case <-gate:
		case <-sc.doneServing:
			return
		***REMOVED***
		if terminalReadFrameError(err) ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// frameWriteResult is the message passed from writeFrameAsync to the serve goroutine.
type frameWriteResult struct ***REMOVED***
	_   incomparable
	wr  FrameWriteRequest // what was written (or attempted)
	err error             // result of the writeFrame call
***REMOVED***

// writeFrameAsync runs in its own goroutine and writes a single frame
// and then reports when it's done.
// At most one goroutine can be running writeFrameAsync at a time per
// serverConn.
func (sc *serverConn) writeFrameAsync(wr FrameWriteRequest) ***REMOVED***
	err := wr.write.writeFrame(sc)
	sc.wroteFrameCh <- frameWriteResult***REMOVED***wr: wr, err: err***REMOVED***
***REMOVED***

func (sc *serverConn) closeAllStreamsOnConnClose() ***REMOVED***
	sc.serveG.check()
	for _, st := range sc.streams ***REMOVED***
		sc.closeStream(st, errClientDisconnected)
	***REMOVED***
***REMOVED***

func (sc *serverConn) stopShutdownTimer() ***REMOVED***
	sc.serveG.check()
	if t := sc.shutdownTimer; t != nil ***REMOVED***
		t.Stop()
	***REMOVED***
***REMOVED***

func (sc *serverConn) notePanic() ***REMOVED***
	// Note: this is for serverConn.serve panicking, not http.Handler code.
	if testHookOnPanicMu != nil ***REMOVED***
		testHookOnPanicMu.Lock()
		defer testHookOnPanicMu.Unlock()
	***REMOVED***
	if testHookOnPanic != nil ***REMOVED***
		if e := recover(); e != nil ***REMOVED***
			if testHookOnPanic(sc, e) ***REMOVED***
				panic(e)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (sc *serverConn) serve() ***REMOVED***
	sc.serveG.check()
	defer sc.notePanic()
	defer sc.conn.Close()
	defer sc.closeAllStreamsOnConnClose()
	defer sc.stopShutdownTimer()
	defer close(sc.doneServing) // unblocks handlers trying to send

	if VerboseLogs ***REMOVED***
		sc.vlogf("http2: server connection from %v on %p", sc.conn.RemoteAddr(), sc.hs)
	***REMOVED***

	sc.writeFrame(FrameWriteRequest***REMOVED***
		write: writeSettings***REMOVED***
			***REMOVED***SettingMaxFrameSize, sc.srv.maxReadFrameSize()***REMOVED***,
			***REMOVED***SettingMaxConcurrentStreams, sc.advMaxStreams***REMOVED***,
			***REMOVED***SettingMaxHeaderListSize, sc.maxHeaderListSize()***REMOVED***,
			***REMOVED***SettingInitialWindowSize, uint32(sc.srv.initialStreamRecvWindowSize())***REMOVED***,
		***REMOVED***,
	***REMOVED***)
	sc.unackedSettings++

	// Each connection starts with intialWindowSize inflow tokens.
	// If a higher value is configured, we add more tokens.
	if diff := sc.srv.initialConnRecvWindowSize() - initialWindowSize; diff > 0 ***REMOVED***
		sc.sendWindowUpdate(nil, int(diff))
	***REMOVED***

	if err := sc.readPreface(); err != nil ***REMOVED***
		sc.condlogf(err, "http2: server: error reading preface from client %v: %v", sc.conn.RemoteAddr(), err)
		return
	***REMOVED***
	// Now that we've got the preface, get us out of the
	// "StateNew" state. We can't go directly to idle, though.
	// Active means we read some data and anticipate a request. We'll
	// do another Active when we get a HEADERS frame.
	sc.setConnState(http.StateActive)
	sc.setConnState(http.StateIdle)

	if sc.srv.IdleTimeout != 0 ***REMOVED***
		sc.idleTimer = time.AfterFunc(sc.srv.IdleTimeout, sc.onIdleTimer)
		defer sc.idleTimer.Stop()
	***REMOVED***

	go sc.readFrames() // closed by defer sc.conn.Close above

	settingsTimer := time.AfterFunc(firstSettingsTimeout, sc.onSettingsTimer)
	defer settingsTimer.Stop()

	loopNum := 0
	for ***REMOVED***
		loopNum++
		select ***REMOVED***
		case wr := <-sc.wantWriteFrameCh:
			if se, ok := wr.write.(StreamError); ok ***REMOVED***
				sc.resetStream(se)
				break
			***REMOVED***
			sc.writeFrame(wr)
		case res := <-sc.wroteFrameCh:
			sc.wroteFrame(res)
		case res := <-sc.readFrameCh:
			if !sc.processFrameFromReader(res) ***REMOVED***
				return
			***REMOVED***
			res.readMore()
			if settingsTimer != nil ***REMOVED***
				settingsTimer.Stop()
				settingsTimer = nil
			***REMOVED***
		case m := <-sc.bodyReadCh:
			sc.noteBodyRead(m.st, m.n)
		case msg := <-sc.serveMsgCh:
			switch v := msg.(type) ***REMOVED***
			case func(int):
				v(loopNum) // for testing
			case *serverMessage:
				switch v ***REMOVED***
				case settingsTimerMsg:
					sc.logf("timeout waiting for SETTINGS frames from %v", sc.conn.RemoteAddr())
					return
				case idleTimerMsg:
					sc.vlogf("connection is idle")
					sc.goAway(ErrCodeNo)
				case shutdownTimerMsg:
					sc.vlogf("GOAWAY close timer fired; closing conn from %v", sc.conn.RemoteAddr())
					return
				case gracefulShutdownMsg:
					sc.startGracefulShutdownInternal()
				default:
					panic("unknown timer")
				***REMOVED***
			case *startPushRequest:
				sc.startPush(v)
			default:
				panic(fmt.Sprintf("unexpected type %T", v))
			***REMOVED***
		***REMOVED***

		// If the peer is causing us to generate a lot of control frames,
		// but not reading them from us, assume they are trying to make us
		// run out of memory.
		if sc.queuedControlFrames > sc.srv.maxQueuedControlFrames() ***REMOVED***
			sc.vlogf("http2: too many control frames in send queue, closing connection")
			return
		***REMOVED***

		// Start the shutdown timer after sending a GOAWAY. When sending GOAWAY
		// with no error code (graceful shutdown), don't start the timer until
		// all open streams have been completed.
		sentGoAway := sc.inGoAway && !sc.needToSendGoAway && !sc.writingFrame
		gracefulShutdownComplete := sc.goAwayCode == ErrCodeNo && sc.curOpenStreams() == 0
		if sentGoAway && sc.shutdownTimer == nil && (sc.goAwayCode != ErrCodeNo || gracefulShutdownComplete) ***REMOVED***
			sc.shutDownIn(goAwayTimeout)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (sc *serverConn) awaitGracefulShutdown(sharedCh <-chan struct***REMOVED******REMOVED***, privateCh chan struct***REMOVED******REMOVED***) ***REMOVED***
	select ***REMOVED***
	case <-sc.doneServing:
	case <-sharedCh:
		close(privateCh)
	***REMOVED***
***REMOVED***

type serverMessage int

// Message values sent to serveMsgCh.
var (
	settingsTimerMsg    = new(serverMessage)
	idleTimerMsg        = new(serverMessage)
	shutdownTimerMsg    = new(serverMessage)
	gracefulShutdownMsg = new(serverMessage)
)

func (sc *serverConn) onSettingsTimer() ***REMOVED*** sc.sendServeMsg(settingsTimerMsg) ***REMOVED***
func (sc *serverConn) onIdleTimer()     ***REMOVED*** sc.sendServeMsg(idleTimerMsg) ***REMOVED***
func (sc *serverConn) onShutdownTimer() ***REMOVED*** sc.sendServeMsg(shutdownTimerMsg) ***REMOVED***

func (sc *serverConn) sendServeMsg(msg interface***REMOVED******REMOVED***) ***REMOVED***
	sc.serveG.checkNotOn() // NOT
	select ***REMOVED***
	case sc.serveMsgCh <- msg:
	case <-sc.doneServing:
	***REMOVED***
***REMOVED***

var errPrefaceTimeout = errors.New("timeout waiting for client preface")

// readPreface reads the ClientPreface greeting from the peer or
// returns errPrefaceTimeout on timeout, or an error if the greeting
// is invalid.
func (sc *serverConn) readPreface() error ***REMOVED***
	errc := make(chan error, 1)
	go func() ***REMOVED***
		// Read the client preface
		buf := make([]byte, len(ClientPreface))
		if _, err := io.ReadFull(sc.conn, buf); err != nil ***REMOVED***
			errc <- err
		***REMOVED*** else if !bytes.Equal(buf, clientPreface) ***REMOVED***
			errc <- fmt.Errorf("bogus greeting %q", buf)
		***REMOVED*** else ***REMOVED***
			errc <- nil
		***REMOVED***
	***REMOVED***()
	timer := time.NewTimer(prefaceTimeout) // TODO: configurable on *Server?
	defer timer.Stop()
	select ***REMOVED***
	case <-timer.C:
		return errPrefaceTimeout
	case err := <-errc:
		if err == nil ***REMOVED***
			if VerboseLogs ***REMOVED***
				sc.vlogf("http2: server: client %v said hello", sc.conn.RemoteAddr())
			***REMOVED***
		***REMOVED***
		return err
	***REMOVED***
***REMOVED***

var errChanPool = sync.Pool***REMOVED***
	New: func() interface***REMOVED******REMOVED*** ***REMOVED*** return make(chan error, 1) ***REMOVED***,
***REMOVED***

var writeDataPool = sync.Pool***REMOVED***
	New: func() interface***REMOVED******REMOVED*** ***REMOVED*** return new(writeData) ***REMOVED***,
***REMOVED***

// writeDataFromHandler writes DATA response frames from a handler on
// the given stream.
func (sc *serverConn) writeDataFromHandler(stream *stream, data []byte, endStream bool) error ***REMOVED***
	ch := errChanPool.Get().(chan error)
	writeArg := writeDataPool.Get().(*writeData)
	*writeArg = writeData***REMOVED***stream.id, data, endStream***REMOVED***
	err := sc.writeFrameFromHandler(FrameWriteRequest***REMOVED***
		write:  writeArg,
		stream: stream,
		done:   ch,
	***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	var frameWriteDone bool // the frame write is done (successfully or not)
	select ***REMOVED***
	case err = <-ch:
		frameWriteDone = true
	case <-sc.doneServing:
		return errClientDisconnected
	case <-stream.cw:
		// If both ch and stream.cw were ready (as might
		// happen on the final Write after an http.Handler
		// ends), prefer the write result. Otherwise this
		// might just be us successfully closing the stream.
		// The writeFrameAsync and serve goroutines guarantee
		// that the ch send will happen before the stream.cw
		// close.
		select ***REMOVED***
		case err = <-ch:
			frameWriteDone = true
		default:
			return errStreamClosed
		***REMOVED***
	***REMOVED***
	errChanPool.Put(ch)
	if frameWriteDone ***REMOVED***
		writeDataPool.Put(writeArg)
	***REMOVED***
	return err
***REMOVED***

// writeFrameFromHandler sends wr to sc.wantWriteFrameCh, but aborts
// if the connection has gone away.
//
// This must not be run from the serve goroutine itself, else it might
// deadlock writing to sc.wantWriteFrameCh (which is only mildly
// buffered and is read by serve itself). If you're on the serve
// goroutine, call writeFrame instead.
func (sc *serverConn) writeFrameFromHandler(wr FrameWriteRequest) error ***REMOVED***
	sc.serveG.checkNotOn() // NOT
	select ***REMOVED***
	case sc.wantWriteFrameCh <- wr:
		return nil
	case <-sc.doneServing:
		// Serve loop is gone.
		// Client has closed their connection to the server.
		return errClientDisconnected
	***REMOVED***
***REMOVED***

// writeFrame schedules a frame to write and sends it if there's nothing
// already being written.
//
// There is no pushback here (the serve goroutine never blocks). It's
// the http.Handlers that block, waiting for their previous frames to
// make it onto the wire
//
// If you're not on the serve goroutine, use writeFrameFromHandler instead.
func (sc *serverConn) writeFrame(wr FrameWriteRequest) ***REMOVED***
	sc.serveG.check()

	// If true, wr will not be written and wr.done will not be signaled.
	var ignoreWrite bool

	// We are not allowed to write frames on closed streams. RFC 7540 Section
	// 5.1.1 says: "An endpoint MUST NOT send frames other than PRIORITY on
	// a closed stream." Our server never sends PRIORITY, so that exception
	// does not apply.
	//
	// The serverConn might close an open stream while the stream's handler
	// is still running. For example, the server might close a stream when it
	// receives bad data from the client. If this happens, the handler might
	// attempt to write a frame after the stream has been closed (since the
	// handler hasn't yet been notified of the close). In this case, we simply
	// ignore the frame. The handler will notice that the stream is closed when
	// it waits for the frame to be written.
	//
	// As an exception to this rule, we allow sending RST_STREAM after close.
	// This allows us to immediately reject new streams without tracking any
	// state for those streams (except for the queued RST_STREAM frame). This
	// may result in duplicate RST_STREAMs in some cases, but the client should
	// ignore those.
	if wr.StreamID() != 0 ***REMOVED***
		_, isReset := wr.write.(StreamError)
		if state, _ := sc.state(wr.StreamID()); state == stateClosed && !isReset ***REMOVED***
			ignoreWrite = true
		***REMOVED***
	***REMOVED***

	// Don't send a 100-continue response if we've already sent headers.
	// See golang.org/issue/14030.
	switch wr.write.(type) ***REMOVED***
	case *writeResHeaders:
		wr.stream.wroteHeaders = true
	case write100ContinueHeadersFrame:
		if wr.stream.wroteHeaders ***REMOVED***
			// We do not need to notify wr.done because this frame is
			// never written with wr.done != nil.
			if wr.done != nil ***REMOVED***
				panic("wr.done != nil for write100ContinueHeadersFrame")
			***REMOVED***
			ignoreWrite = true
		***REMOVED***
	***REMOVED***

	if !ignoreWrite ***REMOVED***
		if wr.isControl() ***REMOVED***
			sc.queuedControlFrames++
			// For extra safety, detect wraparounds, which should not happen,
			// and pull the plug.
			if sc.queuedControlFrames < 0 ***REMOVED***
				sc.conn.Close()
			***REMOVED***
		***REMOVED***
		sc.writeSched.Push(wr)
	***REMOVED***
	sc.scheduleFrameWrite()
***REMOVED***

// startFrameWrite starts a goroutine to write wr (in a separate
// goroutine since that might block on the network), and updates the
// serve goroutine's state about the world, updated from info in wr.
func (sc *serverConn) startFrameWrite(wr FrameWriteRequest) ***REMOVED***
	sc.serveG.check()
	if sc.writingFrame ***REMOVED***
		panic("internal error: can only be writing one frame at a time")
	***REMOVED***

	st := wr.stream
	if st != nil ***REMOVED***
		switch st.state ***REMOVED***
		case stateHalfClosedLocal:
			switch wr.write.(type) ***REMOVED***
			case StreamError, handlerPanicRST, writeWindowUpdate:
				// RFC 7540 Section 5.1 allows sending RST_STREAM, PRIORITY, and WINDOW_UPDATE
				// in this state. (We never send PRIORITY from the server, so that is not checked.)
			default:
				panic(fmt.Sprintf("internal error: attempt to send frame on a half-closed-local stream: %v", wr))
			***REMOVED***
		case stateClosed:
			panic(fmt.Sprintf("internal error: attempt to send frame on a closed stream: %v", wr))
		***REMOVED***
	***REMOVED***
	if wpp, ok := wr.write.(*writePushPromise); ok ***REMOVED***
		var err error
		wpp.promisedID, err = wpp.allocatePromisedID()
		if err != nil ***REMOVED***
			sc.writingFrameAsync = false
			wr.replyToWriter(err)
			return
		***REMOVED***
	***REMOVED***

	sc.writingFrame = true
	sc.needsFrameFlush = true
	if wr.write.staysWithinBuffer(sc.bw.Available()) ***REMOVED***
		sc.writingFrameAsync = false
		err := wr.write.writeFrame(sc)
		sc.wroteFrame(frameWriteResult***REMOVED***wr: wr, err: err***REMOVED***)
	***REMOVED*** else ***REMOVED***
		sc.writingFrameAsync = true
		go sc.writeFrameAsync(wr)
	***REMOVED***
***REMOVED***

// errHandlerPanicked is the error given to any callers blocked in a read from
// Request.Body when the main goroutine panics. Since most handlers read in the
// main ServeHTTP goroutine, this will show up rarely.
var errHandlerPanicked = errors.New("http2: handler panicked")

// wroteFrame is called on the serve goroutine with the result of
// whatever happened on writeFrameAsync.
func (sc *serverConn) wroteFrame(res frameWriteResult) ***REMOVED***
	sc.serveG.check()
	if !sc.writingFrame ***REMOVED***
		panic("internal error: expected to be already writing a frame")
	***REMOVED***
	sc.writingFrame = false
	sc.writingFrameAsync = false

	wr := res.wr

	if writeEndsStream(wr.write) ***REMOVED***
		st := wr.stream
		if st == nil ***REMOVED***
			panic("internal error: expecting non-nil stream")
		***REMOVED***
		switch st.state ***REMOVED***
		case stateOpen:
			// Here we would go to stateHalfClosedLocal in
			// theory, but since our handler is done and
			// the net/http package provides no mechanism
			// for closing a ResponseWriter while still
			// reading data (see possible TODO at top of
			// this file), we go into closed state here
			// anyway, after telling the peer we're
			// hanging up on them. We'll transition to
			// stateClosed after the RST_STREAM frame is
			// written.
			st.state = stateHalfClosedLocal
			// Section 8.1: a server MAY request that the client abort
			// transmission of a request without error by sending a
			// RST_STREAM with an error code of NO_ERROR after sending
			// a complete response.
			sc.resetStream(streamError(st.id, ErrCodeNo))
		case stateHalfClosedRemote:
			sc.closeStream(st, errHandlerComplete)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		switch v := wr.write.(type) ***REMOVED***
		case StreamError:
			// st may be unknown if the RST_STREAM was generated to reject bad input.
			if st, ok := sc.streams[v.StreamID]; ok ***REMOVED***
				sc.closeStream(st, v)
			***REMOVED***
		case handlerPanicRST:
			sc.closeStream(wr.stream, errHandlerPanicked)
		***REMOVED***
	***REMOVED***

	// Reply (if requested) to unblock the ServeHTTP goroutine.
	wr.replyToWriter(res.err)

	sc.scheduleFrameWrite()
***REMOVED***

// scheduleFrameWrite tickles the frame writing scheduler.
//
// If a frame is already being written, nothing happens. This will be called again
// when the frame is done being written.
//
// If a frame isn't being written and we need to send one, the best frame
// to send is selected by writeSched.
//
// If a frame isn't being written and there's nothing else to send, we
// flush the write buffer.
func (sc *serverConn) scheduleFrameWrite() ***REMOVED***
	sc.serveG.check()
	if sc.writingFrame || sc.inFrameScheduleLoop ***REMOVED***
		return
	***REMOVED***
	sc.inFrameScheduleLoop = true
	for !sc.writingFrameAsync ***REMOVED***
		if sc.needToSendGoAway ***REMOVED***
			sc.needToSendGoAway = false
			sc.startFrameWrite(FrameWriteRequest***REMOVED***
				write: &writeGoAway***REMOVED***
					maxStreamID: sc.maxClientStreamID,
					code:        sc.goAwayCode,
				***REMOVED***,
			***REMOVED***)
			continue
		***REMOVED***
		if sc.needToSendSettingsAck ***REMOVED***
			sc.needToSendSettingsAck = false
			sc.startFrameWrite(FrameWriteRequest***REMOVED***write: writeSettingsAck***REMOVED******REMOVED******REMOVED***)
			continue
		***REMOVED***
		if !sc.inGoAway || sc.goAwayCode == ErrCodeNo ***REMOVED***
			if wr, ok := sc.writeSched.Pop(); ok ***REMOVED***
				if wr.isControl() ***REMOVED***
					sc.queuedControlFrames--
				***REMOVED***
				sc.startFrameWrite(wr)
				continue
			***REMOVED***
		***REMOVED***
		if sc.needsFrameFlush ***REMOVED***
			sc.startFrameWrite(FrameWriteRequest***REMOVED***write: flushFrameWriter***REMOVED******REMOVED******REMOVED***)
			sc.needsFrameFlush = false // after startFrameWrite, since it sets this true
			continue
		***REMOVED***
		break
	***REMOVED***
	sc.inFrameScheduleLoop = false
***REMOVED***

// startGracefulShutdown gracefully shuts down a connection. This
// sends GOAWAY with ErrCodeNo to tell the client we're gracefully
// shutting down. The connection isn't closed until all current
// streams are done.
//
// startGracefulShutdown returns immediately; it does not wait until
// the connection has shut down.
func (sc *serverConn) startGracefulShutdown() ***REMOVED***
	sc.serveG.checkNotOn() // NOT
	sc.shutdownOnce.Do(func() ***REMOVED*** sc.sendServeMsg(gracefulShutdownMsg) ***REMOVED***)
***REMOVED***

// After sending GOAWAY, the connection will close after goAwayTimeout.
// If we close the connection immediately after sending GOAWAY, there may
// be unsent data in our kernel receive buffer, which will cause the kernel
// to send a TCP RST on close() instead of a FIN. This RST will abort the
// connection immediately, whether or not the client had received the GOAWAY.
//
// Ideally we should delay for at least 1 RTT + epsilon so the client has
// a chance to read the GOAWAY and stop sending messages. Measuring RTT
// is hard, so we approximate with 1 second. See golang.org/issue/18701.
//
// This is a var so it can be shorter in tests, where all requests uses the
// loopback interface making the expected RTT very small.
//
// TODO: configurable?
var goAwayTimeout = 1 * time.Second

func (sc *serverConn) startGracefulShutdownInternal() ***REMOVED***
	sc.goAway(ErrCodeNo)
***REMOVED***

func (sc *serverConn) goAway(code ErrCode) ***REMOVED***
	sc.serveG.check()
	if sc.inGoAway ***REMOVED***
		return
	***REMOVED***
	sc.inGoAway = true
	sc.needToSendGoAway = true
	sc.goAwayCode = code
	sc.scheduleFrameWrite()
***REMOVED***

func (sc *serverConn) shutDownIn(d time.Duration) ***REMOVED***
	sc.serveG.check()
	sc.shutdownTimer = time.AfterFunc(d, sc.onShutdownTimer)
***REMOVED***

func (sc *serverConn) resetStream(se StreamError) ***REMOVED***
	sc.serveG.check()
	sc.writeFrame(FrameWriteRequest***REMOVED***write: se***REMOVED***)
	if st, ok := sc.streams[se.StreamID]; ok ***REMOVED***
		st.resetQueued = true
	***REMOVED***
***REMOVED***

// processFrameFromReader processes the serve loop's read from readFrameCh from the
// frame-reading goroutine.
// processFrameFromReader returns whether the connection should be kept open.
func (sc *serverConn) processFrameFromReader(res readFrameResult) bool ***REMOVED***
	sc.serveG.check()
	err := res.err
	if err != nil ***REMOVED***
		if err == ErrFrameTooLarge ***REMOVED***
			sc.goAway(ErrCodeFrameSize)
			return true // goAway will close the loop
		***REMOVED***
		clientGone := err == io.EOF || err == io.ErrUnexpectedEOF || isClosedConnError(err)
		if clientGone ***REMOVED***
			// TODO: could we also get into this state if
			// the peer does a half close
			// (e.g. CloseWrite) because they're done
			// sending frames but they're still wanting
			// our open replies?  Investigate.
			// TODO: add CloseWrite to crypto/tls.Conn first
			// so we have a way to test this? I suppose
			// just for testing we could have a non-TLS mode.
			return false
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		f := res.f
		if VerboseLogs ***REMOVED***
			sc.vlogf("http2: server read frame %v", summarizeFrame(f))
		***REMOVED***
		err = sc.processFrame(f)
		if err == nil ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***

	switch ev := err.(type) ***REMOVED***
	case StreamError:
		sc.resetStream(ev)
		return true
	case goAwayFlowError:
		sc.goAway(ErrCodeFlowControl)
		return true
	case ConnectionError:
		sc.logf("http2: server connection error from %v: %v", sc.conn.RemoteAddr(), ev)
		sc.goAway(ErrCode(ev))
		return true // goAway will handle shutdown
	default:
		if res.err != nil ***REMOVED***
			sc.vlogf("http2: server closing client connection; error reading frame from client %s: %v", sc.conn.RemoteAddr(), err)
		***REMOVED*** else ***REMOVED***
			sc.logf("http2: server closing client connection: %v", err)
		***REMOVED***
		return false
	***REMOVED***
***REMOVED***

func (sc *serverConn) processFrame(f Frame) error ***REMOVED***
	sc.serveG.check()

	// First frame received must be SETTINGS.
	if !sc.sawFirstSettings ***REMOVED***
		if _, ok := f.(*SettingsFrame); !ok ***REMOVED***
			return ConnectionError(ErrCodeProtocol)
		***REMOVED***
		sc.sawFirstSettings = true
	***REMOVED***

	switch f := f.(type) ***REMOVED***
	case *SettingsFrame:
		return sc.processSettings(f)
	case *MetaHeadersFrame:
		return sc.processHeaders(f)
	case *WindowUpdateFrame:
		return sc.processWindowUpdate(f)
	case *PingFrame:
		return sc.processPing(f)
	case *DataFrame:
		return sc.processData(f)
	case *RSTStreamFrame:
		return sc.processResetStream(f)
	case *PriorityFrame:
		return sc.processPriority(f)
	case *GoAwayFrame:
		return sc.processGoAway(f)
	case *PushPromiseFrame:
		// A client cannot push. Thus, servers MUST treat the receipt of a PUSH_PROMISE
		// frame as a connection error (Section 5.4.1) of type PROTOCOL_ERROR.
		return ConnectionError(ErrCodeProtocol)
	default:
		sc.vlogf("http2: server ignoring frame: %v", f.Header())
		return nil
	***REMOVED***
***REMOVED***

func (sc *serverConn) processPing(f *PingFrame) error ***REMOVED***
	sc.serveG.check()
	if f.IsAck() ***REMOVED***
		// 6.7 PING: " An endpoint MUST NOT respond to PING frames
		// containing this flag."
		return nil
	***REMOVED***
	if f.StreamID != 0 ***REMOVED***
		// "PING frames are not associated with any individual
		// stream. If a PING frame is received with a stream
		// identifier field value other than 0x0, the recipient MUST
		// respond with a connection error (Section 5.4.1) of type
		// PROTOCOL_ERROR."
		return ConnectionError(ErrCodeProtocol)
	***REMOVED***
	if sc.inGoAway && sc.goAwayCode != ErrCodeNo ***REMOVED***
		return nil
	***REMOVED***
	sc.writeFrame(FrameWriteRequest***REMOVED***write: writePingAck***REMOVED***f***REMOVED******REMOVED***)
	return nil
***REMOVED***

func (sc *serverConn) processWindowUpdate(f *WindowUpdateFrame) error ***REMOVED***
	sc.serveG.check()
	switch ***REMOVED***
	case f.StreamID != 0: // stream-level flow control
		state, st := sc.state(f.StreamID)
		if state == stateIdle ***REMOVED***
			// Section 5.1: "Receiving any frame other than HEADERS
			// or PRIORITY on a stream in this state MUST be
			// treated as a connection error (Section 5.4.1) of
			// type PROTOCOL_ERROR."
			return ConnectionError(ErrCodeProtocol)
		***REMOVED***
		if st == nil ***REMOVED***
			// "WINDOW_UPDATE can be sent by a peer that has sent a
			// frame bearing the END_STREAM flag. This means that a
			// receiver could receive a WINDOW_UPDATE frame on a "half
			// closed (remote)" or "closed" stream. A receiver MUST
			// NOT treat this as an error, see Section 5.1."
			return nil
		***REMOVED***
		if !st.flow.add(int32(f.Increment)) ***REMOVED***
			return streamError(f.StreamID, ErrCodeFlowControl)
		***REMOVED***
	default: // connection-level flow control
		if !sc.flow.add(int32(f.Increment)) ***REMOVED***
			return goAwayFlowError***REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***
	sc.scheduleFrameWrite()
	return nil
***REMOVED***

func (sc *serverConn) processResetStream(f *RSTStreamFrame) error ***REMOVED***
	sc.serveG.check()

	state, st := sc.state(f.StreamID)
	if state == stateIdle ***REMOVED***
		// 6.4 "RST_STREAM frames MUST NOT be sent for a
		// stream in the "idle" state. If a RST_STREAM frame
		// identifying an idle stream is received, the
		// recipient MUST treat this as a connection error
		// (Section 5.4.1) of type PROTOCOL_ERROR.
		return ConnectionError(ErrCodeProtocol)
	***REMOVED***
	if st != nil ***REMOVED***
		st.cancelCtx()
		sc.closeStream(st, streamError(f.StreamID, f.ErrCode))
	***REMOVED***
	return nil
***REMOVED***

func (sc *serverConn) closeStream(st *stream, err error) ***REMOVED***
	sc.serveG.check()
	if st.state == stateIdle || st.state == stateClosed ***REMOVED***
		panic(fmt.Sprintf("invariant; can't close stream in state %v", st.state))
	***REMOVED***
	st.state = stateClosed
	if st.writeDeadline != nil ***REMOVED***
		st.writeDeadline.Stop()
	***REMOVED***
	if st.isPushed() ***REMOVED***
		sc.curPushedStreams--
	***REMOVED*** else ***REMOVED***
		sc.curClientStreams--
	***REMOVED***
	delete(sc.streams, st.id)
	if len(sc.streams) == 0 ***REMOVED***
		sc.setConnState(http.StateIdle)
		if sc.srv.IdleTimeout != 0 ***REMOVED***
			sc.idleTimer.Reset(sc.srv.IdleTimeout)
		***REMOVED***
		if h1ServerKeepAlivesDisabled(sc.hs) ***REMOVED***
			sc.startGracefulShutdownInternal()
		***REMOVED***
	***REMOVED***
	if p := st.body; p != nil ***REMOVED***
		// Return any buffered unread bytes worth of conn-level flow control.
		// See golang.org/issue/16481
		sc.sendWindowUpdate(nil, p.Len())

		p.CloseWithError(err)
	***REMOVED***
	st.cw.Close() // signals Handler's CloseNotifier, unblocks writes, etc
	sc.writeSched.CloseStream(st.id)
***REMOVED***

func (sc *serverConn) processSettings(f *SettingsFrame) error ***REMOVED***
	sc.serveG.check()
	if f.IsAck() ***REMOVED***
		sc.unackedSettings--
		if sc.unackedSettings < 0 ***REMOVED***
			// Why is the peer ACKing settings we never sent?
			// The spec doesn't mention this case, but
			// hang up on them anyway.
			return ConnectionError(ErrCodeProtocol)
		***REMOVED***
		return nil
	***REMOVED***
	if f.NumSettings() > 100 || f.HasDuplicates() ***REMOVED***
		// This isn't actually in the spec, but hang up on
		// suspiciously large settings frames or those with
		// duplicate entries.
		return ConnectionError(ErrCodeProtocol)
	***REMOVED***
	if err := f.ForeachSetting(sc.processSetting); err != nil ***REMOVED***
		return err
	***REMOVED***
	// TODO: judging by RFC 7540, Section 6.5.3 each SETTINGS frame should be
	// acknowledged individually, even if multiple are received before the ACK.
	sc.needToSendSettingsAck = true
	sc.scheduleFrameWrite()
	return nil
***REMOVED***

func (sc *serverConn) processSetting(s Setting) error ***REMOVED***
	sc.serveG.check()
	if err := s.Valid(); err != nil ***REMOVED***
		return err
	***REMOVED***
	if VerboseLogs ***REMOVED***
		sc.vlogf("http2: server processing setting %v", s)
	***REMOVED***
	switch s.ID ***REMOVED***
	case SettingHeaderTableSize:
		sc.headerTableSize = s.Val
		sc.hpackEncoder.SetMaxDynamicTableSize(s.Val)
	case SettingEnablePush:
		sc.pushEnabled = s.Val != 0
	case SettingMaxConcurrentStreams:
		sc.clientMaxStreams = s.Val
	case SettingInitialWindowSize:
		return sc.processSettingInitialWindowSize(s.Val)
	case SettingMaxFrameSize:
		sc.maxFrameSize = int32(s.Val) // the maximum valid s.Val is < 2^31
	case SettingMaxHeaderListSize:
		sc.peerMaxHeaderListSize = s.Val
	default:
		// Unknown setting: "An endpoint that receives a SETTINGS
		// frame with any unknown or unsupported identifier MUST
		// ignore that setting."
		if VerboseLogs ***REMOVED***
			sc.vlogf("http2: server ignoring unknown setting %v", s)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (sc *serverConn) processSettingInitialWindowSize(val uint32) error ***REMOVED***
	sc.serveG.check()
	// Note: val already validated to be within range by
	// processSetting's Valid call.

	// "A SETTINGS frame can alter the initial flow control window
	// size for all current streams. When the value of
	// SETTINGS_INITIAL_WINDOW_SIZE changes, a receiver MUST
	// adjust the size of all stream flow control windows that it
	// maintains by the difference between the new value and the
	// old value."
	old := sc.initialStreamSendWindowSize
	sc.initialStreamSendWindowSize = int32(val)
	growth := int32(val) - old // may be negative
	for _, st := range sc.streams ***REMOVED***
		if !st.flow.add(growth) ***REMOVED***
			// 6.9.2 Initial Flow Control Window Size
			// "An endpoint MUST treat a change to
			// SETTINGS_INITIAL_WINDOW_SIZE that causes any flow
			// control window to exceed the maximum size as a
			// connection error (Section 5.4.1) of type
			// FLOW_CONTROL_ERROR."
			return ConnectionError(ErrCodeFlowControl)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (sc *serverConn) processData(f *DataFrame) error ***REMOVED***
	sc.serveG.check()
	if sc.inGoAway && sc.goAwayCode != ErrCodeNo ***REMOVED***
		return nil
	***REMOVED***
	data := f.Data()

	// "If a DATA frame is received whose stream is not in "open"
	// or "half closed (local)" state, the recipient MUST respond
	// with a stream error (Section 5.4.2) of type STREAM_CLOSED."
	id := f.Header().StreamID
	state, st := sc.state(id)
	if id == 0 || state == stateIdle ***REMOVED***
		// Section 5.1: "Receiving any frame other than HEADERS
		// or PRIORITY on a stream in this state MUST be
		// treated as a connection error (Section 5.4.1) of
		// type PROTOCOL_ERROR."
		return ConnectionError(ErrCodeProtocol)
	***REMOVED***
	if st == nil || state != stateOpen || st.gotTrailerHeader || st.resetQueued ***REMOVED***
		// This includes sending a RST_STREAM if the stream is
		// in stateHalfClosedLocal (which currently means that
		// the http.Handler returned, so it's done reading &
		// done writing). Try to stop the client from sending
		// more DATA.

		// But still enforce their connection-level flow control,
		// and return any flow control bytes since we're not going
		// to consume them.
		if sc.inflow.available() < int32(f.Length) ***REMOVED***
			return streamError(id, ErrCodeFlowControl)
		***REMOVED***
		// Deduct the flow control from inflow, since we're
		// going to immediately add it back in
		// sendWindowUpdate, which also schedules sending the
		// frames.
		sc.inflow.take(int32(f.Length))
		sc.sendWindowUpdate(nil, int(f.Length)) // conn-level

		if st != nil && st.resetQueued ***REMOVED***
			// Already have a stream error in flight. Don't send another.
			return nil
		***REMOVED***
		return streamError(id, ErrCodeStreamClosed)
	***REMOVED***
	if st.body == nil ***REMOVED***
		panic("internal error: should have a body in this state")
	***REMOVED***

	// Sender sending more than they'd declared?
	if st.declBodyBytes != -1 && st.bodyBytes+int64(len(data)) > st.declBodyBytes ***REMOVED***
		st.body.CloseWithError(fmt.Errorf("sender tried to send more than declared Content-Length of %d bytes", st.declBodyBytes))
		// RFC 7540, sec 8.1.2.6: A request or response is also malformed if the
		// value of a content-length header field does not equal the sum of the
		// DATA frame payload lengths that form the body.
		return streamError(id, ErrCodeProtocol)
	***REMOVED***
	if f.Length > 0 ***REMOVED***
		// Check whether the client has flow control quota.
		if st.inflow.available() < int32(f.Length) ***REMOVED***
			return streamError(id, ErrCodeFlowControl)
		***REMOVED***
		st.inflow.take(int32(f.Length))

		if len(data) > 0 ***REMOVED***
			wrote, err := st.body.Write(data)
			if err != nil ***REMOVED***
				sc.sendWindowUpdate(nil, int(f.Length)-wrote)
				return streamError(id, ErrCodeStreamClosed)
			***REMOVED***
			if wrote != len(data) ***REMOVED***
				panic("internal error: bad Writer")
			***REMOVED***
			st.bodyBytes += int64(len(data))
		***REMOVED***

		// Return any padded flow control now, since we won't
		// refund it later on body reads.
		if pad := int32(f.Length) - int32(len(data)); pad > 0 ***REMOVED***
			sc.sendWindowUpdate32(nil, pad)
			sc.sendWindowUpdate32(st, pad)
		***REMOVED***
	***REMOVED***
	if f.StreamEnded() ***REMOVED***
		st.endStream()
	***REMOVED***
	return nil
***REMOVED***

func (sc *serverConn) processGoAway(f *GoAwayFrame) error ***REMOVED***
	sc.serveG.check()
	if f.ErrCode != ErrCodeNo ***REMOVED***
		sc.logf("http2: received GOAWAY %+v, starting graceful shutdown", f)
	***REMOVED*** else ***REMOVED***
		sc.vlogf("http2: received GOAWAY %+v, starting graceful shutdown", f)
	***REMOVED***
	sc.startGracefulShutdownInternal()
	// http://tools.ietf.org/html/rfc7540#section-6.8
	// We should not create any new streams, which means we should disable push.
	sc.pushEnabled = false
	return nil
***REMOVED***

// isPushed reports whether the stream is server-initiated.
func (st *stream) isPushed() bool ***REMOVED***
	return st.id%2 == 0
***REMOVED***

// endStream closes a Request.Body's pipe. It is called when a DATA
// frame says a request body is over (or after trailers).
func (st *stream) endStream() ***REMOVED***
	sc := st.sc
	sc.serveG.check()

	if st.declBodyBytes != -1 && st.declBodyBytes != st.bodyBytes ***REMOVED***
		st.body.CloseWithError(fmt.Errorf("request declared a Content-Length of %d but only wrote %d bytes",
			st.declBodyBytes, st.bodyBytes))
	***REMOVED*** else ***REMOVED***
		st.body.closeWithErrorAndCode(io.EOF, st.copyTrailersToHandlerRequest)
		st.body.CloseWithError(io.EOF)
	***REMOVED***
	st.state = stateHalfClosedRemote
***REMOVED***

// copyTrailersToHandlerRequest is run in the Handler's goroutine in
// its Request.Body.Read just before it gets io.EOF.
func (st *stream) copyTrailersToHandlerRequest() ***REMOVED***
	for k, vv := range st.trailer ***REMOVED***
		if _, ok := st.reqTrailer[k]; ok ***REMOVED***
			// Only copy it over it was pre-declared.
			st.reqTrailer[k] = vv
		***REMOVED***
	***REMOVED***
***REMOVED***

// onWriteTimeout is run on its own goroutine (from time.AfterFunc)
// when the stream's WriteTimeout has fired.
func (st *stream) onWriteTimeout() ***REMOVED***
	st.sc.writeFrameFromHandler(FrameWriteRequest***REMOVED***write: streamError(st.id, ErrCodeInternal)***REMOVED***)
***REMOVED***

func (sc *serverConn) processHeaders(f *MetaHeadersFrame) error ***REMOVED***
	sc.serveG.check()
	id := f.StreamID
	if sc.inGoAway ***REMOVED***
		// Ignore.
		return nil
	***REMOVED***
	// http://tools.ietf.org/html/rfc7540#section-5.1.1
	// Streams initiated by a client MUST use odd-numbered stream
	// identifiers. [...] An endpoint that receives an unexpected
	// stream identifier MUST respond with a connection error
	// (Section 5.4.1) of type PROTOCOL_ERROR.
	if id%2 != 1 ***REMOVED***
		return ConnectionError(ErrCodeProtocol)
	***REMOVED***
	// A HEADERS frame can be used to create a new stream or
	// send a trailer for an open one. If we already have a stream
	// open, let it process its own HEADERS frame (trailers at this
	// point, if it's valid).
	if st := sc.streams[f.StreamID]; st != nil ***REMOVED***
		if st.resetQueued ***REMOVED***
			// We're sending RST_STREAM to close the stream, so don't bother
			// processing this frame.
			return nil
		***REMOVED***
		// RFC 7540, sec 5.1: If an endpoint receives additional frames, other than
		// WINDOW_UPDATE, PRIORITY, or RST_STREAM, for a stream that is in
		// this state, it MUST respond with a stream error (Section 5.4.2) of
		// type STREAM_CLOSED.
		if st.state == stateHalfClosedRemote ***REMOVED***
			return streamError(id, ErrCodeStreamClosed)
		***REMOVED***
		return st.processTrailerHeaders(f)
	***REMOVED***

	// [...] The identifier of a newly established stream MUST be
	// numerically greater than all streams that the initiating
	// endpoint has opened or reserved. [...]  An endpoint that
	// receives an unexpected stream identifier MUST respond with
	// a connection error (Section 5.4.1) of type PROTOCOL_ERROR.
	if id <= sc.maxClientStreamID ***REMOVED***
		return ConnectionError(ErrCodeProtocol)
	***REMOVED***
	sc.maxClientStreamID = id

	if sc.idleTimer != nil ***REMOVED***
		sc.idleTimer.Stop()
	***REMOVED***

	// http://tools.ietf.org/html/rfc7540#section-5.1.2
	// [...] Endpoints MUST NOT exceed the limit set by their peer. An
	// endpoint that receives a HEADERS frame that causes their
	// advertised concurrent stream limit to be exceeded MUST treat
	// this as a stream error (Section 5.4.2) of type PROTOCOL_ERROR
	// or REFUSED_STREAM.
	if sc.curClientStreams+1 > sc.advMaxStreams ***REMOVED***
		if sc.unackedSettings == 0 ***REMOVED***
			// They should know better.
			return streamError(id, ErrCodeProtocol)
		***REMOVED***
		// Assume it's a network race, where they just haven't
		// received our last SETTINGS update. But actually
		// this can't happen yet, because we don't yet provide
		// a way for users to adjust server parameters at
		// runtime.
		return streamError(id, ErrCodeRefusedStream)
	***REMOVED***

	initialState := stateOpen
	if f.StreamEnded() ***REMOVED***
		initialState = stateHalfClosedRemote
	***REMOVED***
	st := sc.newStream(id, 0, initialState)

	if f.HasPriority() ***REMOVED***
		if err := checkPriority(f.StreamID, f.Priority); err != nil ***REMOVED***
			return err
		***REMOVED***
		sc.writeSched.AdjustStream(st.id, f.Priority)
	***REMOVED***

	rw, req, err := sc.newWriterAndRequest(st, f)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	st.reqTrailer = req.Trailer
	if st.reqTrailer != nil ***REMOVED***
		st.trailer = make(http.Header)
	***REMOVED***
	st.body = req.Body.(*requestBody).pipe // may be nil
	st.declBodyBytes = req.ContentLength

	handler := sc.handler.ServeHTTP
	if f.Truncated ***REMOVED***
		// Their header list was too long. Send a 431 error.
		handler = handleHeaderListTooLong
	***REMOVED*** else if err := checkValidHTTP2RequestHeaders(req.Header); err != nil ***REMOVED***
		handler = new400Handler(err)
	***REMOVED***

	// The net/http package sets the read deadline from the
	// http.Server.ReadTimeout during the TLS handshake, but then
	// passes the connection off to us with the deadline already
	// set. Disarm it here after the request headers are read,
	// similar to how the http1 server works. Here it's
	// technically more like the http1 Server's ReadHeaderTimeout
	// (in Go 1.8), though. That's a more sane option anyway.
	if sc.hs.ReadTimeout != 0 ***REMOVED***
		sc.conn.SetReadDeadline(time.Time***REMOVED******REMOVED***)
	***REMOVED***

	go sc.runHandler(rw, req, handler)
	return nil
***REMOVED***

func (st *stream) processTrailerHeaders(f *MetaHeadersFrame) error ***REMOVED***
	sc := st.sc
	sc.serveG.check()
	if st.gotTrailerHeader ***REMOVED***
		return ConnectionError(ErrCodeProtocol)
	***REMOVED***
	st.gotTrailerHeader = true
	if !f.StreamEnded() ***REMOVED***
		return streamError(st.id, ErrCodeProtocol)
	***REMOVED***

	if len(f.PseudoFields()) > 0 ***REMOVED***
		return streamError(st.id, ErrCodeProtocol)
	***REMOVED***
	if st.trailer != nil ***REMOVED***
		for _, hf := range f.RegularFields() ***REMOVED***
			key := sc.canonicalHeader(hf.Name)
			if !httpguts.ValidTrailerHeader(key) ***REMOVED***
				// TODO: send more details to the peer somehow. But http2 has
				// no way to send debug data at a stream level. Discuss with
				// HTTP folk.
				return streamError(st.id, ErrCodeProtocol)
			***REMOVED***
			st.trailer[key] = append(st.trailer[key], hf.Value)
		***REMOVED***
	***REMOVED***
	st.endStream()
	return nil
***REMOVED***

func checkPriority(streamID uint32, p PriorityParam) error ***REMOVED***
	if streamID == p.StreamDep ***REMOVED***
		// Section 5.3.1: "A stream cannot depend on itself. An endpoint MUST treat
		// this as a stream error (Section 5.4.2) of type PROTOCOL_ERROR."
		// Section 5.3.3 says that a stream can depend on one of its dependencies,
		// so it's only self-dependencies that are forbidden.
		return streamError(streamID, ErrCodeProtocol)
	***REMOVED***
	return nil
***REMOVED***

func (sc *serverConn) processPriority(f *PriorityFrame) error ***REMOVED***
	if sc.inGoAway ***REMOVED***
		return nil
	***REMOVED***
	if err := checkPriority(f.StreamID, f.PriorityParam); err != nil ***REMOVED***
		return err
	***REMOVED***
	sc.writeSched.AdjustStream(f.StreamID, f.PriorityParam)
	return nil
***REMOVED***

func (sc *serverConn) newStream(id, pusherID uint32, state streamState) *stream ***REMOVED***
	sc.serveG.check()
	if id == 0 ***REMOVED***
		panic("internal error: cannot create stream with id 0")
	***REMOVED***

	ctx, cancelCtx := context.WithCancel(sc.baseCtx)
	st := &stream***REMOVED***
		sc:        sc,
		id:        id,
		state:     state,
		ctx:       ctx,
		cancelCtx: cancelCtx,
	***REMOVED***
	st.cw.Init()
	st.flow.conn = &sc.flow // link to conn-level counter
	st.flow.add(sc.initialStreamSendWindowSize)
	st.inflow.conn = &sc.inflow // link to conn-level counter
	st.inflow.add(sc.srv.initialStreamRecvWindowSize())
	if sc.hs.WriteTimeout != 0 ***REMOVED***
		st.writeDeadline = time.AfterFunc(sc.hs.WriteTimeout, st.onWriteTimeout)
	***REMOVED***

	sc.streams[id] = st
	sc.writeSched.OpenStream(st.id, OpenStreamOptions***REMOVED***PusherID: pusherID***REMOVED***)
	if st.isPushed() ***REMOVED***
		sc.curPushedStreams++
	***REMOVED*** else ***REMOVED***
		sc.curClientStreams++
	***REMOVED***
	if sc.curOpenStreams() == 1 ***REMOVED***
		sc.setConnState(http.StateActive)
	***REMOVED***

	return st
***REMOVED***

func (sc *serverConn) newWriterAndRequest(st *stream, f *MetaHeadersFrame) (*responseWriter, *http.Request, error) ***REMOVED***
	sc.serveG.check()

	rp := requestParam***REMOVED***
		method:    f.PseudoValue("method"),
		scheme:    f.PseudoValue("scheme"),
		authority: f.PseudoValue("authority"),
		path:      f.PseudoValue("path"),
	***REMOVED***

	isConnect := rp.method == "CONNECT"
	if isConnect ***REMOVED***
		if rp.path != "" || rp.scheme != "" || rp.authority == "" ***REMOVED***
			return nil, nil, streamError(f.StreamID, ErrCodeProtocol)
		***REMOVED***
	***REMOVED*** else if rp.method == "" || rp.path == "" || (rp.scheme != "https" && rp.scheme != "http") ***REMOVED***
		// See 8.1.2.6 Malformed Requests and Responses:
		//
		// Malformed requests or responses that are detected
		// MUST be treated as a stream error (Section 5.4.2)
		// of type PROTOCOL_ERROR."
		//
		// 8.1.2.3 Request Pseudo-Header Fields
		// "All HTTP/2 requests MUST include exactly one valid
		// value for the :method, :scheme, and :path
		// pseudo-header fields"
		return nil, nil, streamError(f.StreamID, ErrCodeProtocol)
	***REMOVED***

	bodyOpen := !f.StreamEnded()
	if rp.method == "HEAD" && bodyOpen ***REMOVED***
		// HEAD requests can't have bodies
		return nil, nil, streamError(f.StreamID, ErrCodeProtocol)
	***REMOVED***

	rp.header = make(http.Header)
	for _, hf := range f.RegularFields() ***REMOVED***
		rp.header.Add(sc.canonicalHeader(hf.Name), hf.Value)
	***REMOVED***
	if rp.authority == "" ***REMOVED***
		rp.authority = rp.header.Get("Host")
	***REMOVED***

	rw, req, err := sc.newWriterAndRequestNoBody(st, rp)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	if bodyOpen ***REMOVED***
		if vv, ok := rp.header["Content-Length"]; ok ***REMOVED***
			if cl, err := strconv.ParseUint(vv[0], 10, 63); err == nil ***REMOVED***
				req.ContentLength = int64(cl)
			***REMOVED*** else ***REMOVED***
				req.ContentLength = 0
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			req.ContentLength = -1
		***REMOVED***
		req.Body.(*requestBody).pipe = &pipe***REMOVED***
			b: &dataBuffer***REMOVED***expected: req.ContentLength***REMOVED***,
		***REMOVED***
	***REMOVED***
	return rw, req, nil
***REMOVED***

type requestParam struct ***REMOVED***
	method                  string
	scheme, authority, path string
	header                  http.Header
***REMOVED***

func (sc *serverConn) newWriterAndRequestNoBody(st *stream, rp requestParam) (*responseWriter, *http.Request, error) ***REMOVED***
	sc.serveG.check()

	var tlsState *tls.ConnectionState // nil if not scheme https
	if rp.scheme == "https" ***REMOVED***
		tlsState = sc.tlsState
	***REMOVED***

	needsContinue := rp.header.Get("Expect") == "100-continue"
	if needsContinue ***REMOVED***
		rp.header.Del("Expect")
	***REMOVED***
	// Merge Cookie headers into one "; "-delimited value.
	if cookies := rp.header["Cookie"]; len(cookies) > 1 ***REMOVED***
		rp.header.Set("Cookie", strings.Join(cookies, "; "))
	***REMOVED***

	// Setup Trailers
	var trailer http.Header
	for _, v := range rp.header["Trailer"] ***REMOVED***
		for _, key := range strings.Split(v, ",") ***REMOVED***
			key = http.CanonicalHeaderKey(textproto.TrimString(key))
			switch key ***REMOVED***
			case "Transfer-Encoding", "Trailer", "Content-Length":
				// Bogus. (copy of http1 rules)
				// Ignore.
			default:
				if trailer == nil ***REMOVED***
					trailer = make(http.Header)
				***REMOVED***
				trailer[key] = nil
			***REMOVED***
		***REMOVED***
	***REMOVED***
	delete(rp.header, "Trailer")

	var url_ *url.URL
	var requestURI string
	if rp.method == "CONNECT" ***REMOVED***
		url_ = &url.URL***REMOVED***Host: rp.authority***REMOVED***
		requestURI = rp.authority // mimic HTTP/1 server behavior
	***REMOVED*** else ***REMOVED***
		var err error
		url_, err = url.ParseRequestURI(rp.path)
		if err != nil ***REMOVED***
			return nil, nil, streamError(st.id, ErrCodeProtocol)
		***REMOVED***
		requestURI = rp.path
	***REMOVED***

	body := &requestBody***REMOVED***
		conn:          sc,
		stream:        st,
		needsContinue: needsContinue,
	***REMOVED***
	req := &http.Request***REMOVED***
		Method:     rp.method,
		URL:        url_,
		RemoteAddr: sc.remoteAddrStr,
		Header:     rp.header,
		RequestURI: requestURI,
		Proto:      "HTTP/2.0",
		ProtoMajor: 2,
		ProtoMinor: 0,
		TLS:        tlsState,
		Host:       rp.authority,
		Body:       body,
		Trailer:    trailer,
	***REMOVED***
	req = req.WithContext(st.ctx)

	rws := responseWriterStatePool.Get().(*responseWriterState)
	bwSave := rws.bw
	*rws = responseWriterState***REMOVED******REMOVED*** // zero all the fields
	rws.conn = sc
	rws.bw = bwSave
	rws.bw.Reset(chunkWriter***REMOVED***rws***REMOVED***)
	rws.stream = st
	rws.req = req
	rws.body = body

	rw := &responseWriter***REMOVED***rws: rws***REMOVED***
	return rw, req, nil
***REMOVED***

// Run on its own goroutine.
func (sc *serverConn) runHandler(rw *responseWriter, req *http.Request, handler func(http.ResponseWriter, *http.Request)) ***REMOVED***
	didPanic := true
	defer func() ***REMOVED***
		rw.rws.stream.cancelCtx()
		if didPanic ***REMOVED***
			e := recover()
			sc.writeFrameFromHandler(FrameWriteRequest***REMOVED***
				write:  handlerPanicRST***REMOVED***rw.rws.stream.id***REMOVED***,
				stream: rw.rws.stream,
			***REMOVED***)
			// Same as net/http:
			if e != nil && e != http.ErrAbortHandler ***REMOVED***
				const size = 64 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				sc.logf("http2: panic serving %v: %v\n%s", sc.conn.RemoteAddr(), e, buf)
			***REMOVED***
			return
		***REMOVED***
		rw.handlerDone()
	***REMOVED***()
	handler(rw, req)
	didPanic = false
***REMOVED***

func handleHeaderListTooLong(w http.ResponseWriter, r *http.Request) ***REMOVED***
	// 10.5.1 Limits on Header Block Size:
	// .. "A server that receives a larger header block than it is
	// willing to handle can send an HTTP 431 (Request Header Fields Too
	// Large) status code"
	const statusRequestHeaderFieldsTooLarge = 431 // only in Go 1.6+
	w.WriteHeader(statusRequestHeaderFieldsTooLarge)
	io.WriteString(w, "<h1>HTTP Error 431</h1><p>Request Header Field(s) Too Large</p>")
***REMOVED***

// called from handler goroutines.
// h may be nil.
func (sc *serverConn) writeHeaders(st *stream, headerData *writeResHeaders) error ***REMOVED***
	sc.serveG.checkNotOn() // NOT on
	var errc chan error
	if headerData.h != nil ***REMOVED***
		// If there's a header map (which we don't own), so we have to block on
		// waiting for this frame to be written, so an http.Flush mid-handler
		// writes out the correct value of keys, before a handler later potentially
		// mutates it.
		errc = errChanPool.Get().(chan error)
	***REMOVED***
	if err := sc.writeFrameFromHandler(FrameWriteRequest***REMOVED***
		write:  headerData,
		stream: st,
		done:   errc,
	***REMOVED***); err != nil ***REMOVED***
		return err
	***REMOVED***
	if errc != nil ***REMOVED***
		select ***REMOVED***
		case err := <-errc:
			errChanPool.Put(errc)
			return err
		case <-sc.doneServing:
			return errClientDisconnected
		case <-st.cw:
			return errStreamClosed
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// called from handler goroutines.
func (sc *serverConn) write100ContinueHeaders(st *stream) ***REMOVED***
	sc.writeFrameFromHandler(FrameWriteRequest***REMOVED***
		write:  write100ContinueHeadersFrame***REMOVED***st.id***REMOVED***,
		stream: st,
	***REMOVED***)
***REMOVED***

// A bodyReadMsg tells the server loop that the http.Handler read n
// bytes of the DATA from the client on the given stream.
type bodyReadMsg struct ***REMOVED***
	st *stream
	n  int
***REMOVED***

// called from handler goroutines.
// Notes that the handler for the given stream ID read n bytes of its body
// and schedules flow control tokens to be sent.
func (sc *serverConn) noteBodyReadFromHandler(st *stream, n int, err error) ***REMOVED***
	sc.serveG.checkNotOn() // NOT on
	if n > 0 ***REMOVED***
		select ***REMOVED***
		case sc.bodyReadCh <- bodyReadMsg***REMOVED***st, n***REMOVED***:
		case <-sc.doneServing:
		***REMOVED***
	***REMOVED***
***REMOVED***

func (sc *serverConn) noteBodyRead(st *stream, n int) ***REMOVED***
	sc.serveG.check()
	sc.sendWindowUpdate(nil, n) // conn-level
	if st.state != stateHalfClosedRemote && st.state != stateClosed ***REMOVED***
		// Don't send this WINDOW_UPDATE if the stream is closed
		// remotely.
		sc.sendWindowUpdate(st, n)
	***REMOVED***
***REMOVED***

// st may be nil for conn-level
func (sc *serverConn) sendWindowUpdate(st *stream, n int) ***REMOVED***
	sc.serveG.check()
	// "The legal range for the increment to the flow control
	// window is 1 to 2^31-1 (2,147,483,647) octets."
	// A Go Read call on 64-bit machines could in theory read
	// a larger Read than this. Very unlikely, but we handle it here
	// rather than elsewhere for now.
	const maxUint31 = 1<<31 - 1
	for n >= maxUint31 ***REMOVED***
		sc.sendWindowUpdate32(st, maxUint31)
		n -= maxUint31
	***REMOVED***
	sc.sendWindowUpdate32(st, int32(n))
***REMOVED***

// st may be nil for conn-level
func (sc *serverConn) sendWindowUpdate32(st *stream, n int32) ***REMOVED***
	sc.serveG.check()
	if n == 0 ***REMOVED***
		return
	***REMOVED***
	if n < 0 ***REMOVED***
		panic("negative update")
	***REMOVED***
	var streamID uint32
	if st != nil ***REMOVED***
		streamID = st.id
	***REMOVED***
	sc.writeFrame(FrameWriteRequest***REMOVED***
		write:  writeWindowUpdate***REMOVED***streamID: streamID, n: uint32(n)***REMOVED***,
		stream: st,
	***REMOVED***)
	var ok bool
	if st == nil ***REMOVED***
		ok = sc.inflow.add(n)
	***REMOVED*** else ***REMOVED***
		ok = st.inflow.add(n)
	***REMOVED***
	if !ok ***REMOVED***
		panic("internal error; sent too many window updates without decrements?")
	***REMOVED***
***REMOVED***

// requestBody is the Handler's Request.Body type.
// Read and Close may be called concurrently.
type requestBody struct ***REMOVED***
	_             incomparable
	stream        *stream
	conn          *serverConn
	closed        bool  // for use by Close only
	sawEOF        bool  // for use by Read only
	pipe          *pipe // non-nil if we have a HTTP entity message body
	needsContinue bool  // need to send a 100-continue
***REMOVED***

func (b *requestBody) Close() error ***REMOVED***
	if b.pipe != nil && !b.closed ***REMOVED***
		b.pipe.BreakWithError(errClosedBody)
	***REMOVED***
	b.closed = true
	return nil
***REMOVED***

func (b *requestBody) Read(p []byte) (n int, err error) ***REMOVED***
	if b.needsContinue ***REMOVED***
		b.needsContinue = false
		b.conn.write100ContinueHeaders(b.stream)
	***REMOVED***
	if b.pipe == nil || b.sawEOF ***REMOVED***
		return 0, io.EOF
	***REMOVED***
	n, err = b.pipe.Read(p)
	if err == io.EOF ***REMOVED***
		b.sawEOF = true
	***REMOVED***
	if b.conn == nil && inTests ***REMOVED***
		return
	***REMOVED***
	b.conn.noteBodyReadFromHandler(b.stream, n, err)
	return
***REMOVED***

// responseWriter is the http.ResponseWriter implementation. It's
// intentionally small (1 pointer wide) to minimize garbage. The
// responseWriterState pointer inside is zeroed at the end of a
// request (in handlerDone) and calls on the responseWriter thereafter
// simply crash (caller's mistake), but the much larger responseWriterState
// and buffers are reused between multiple requests.
type responseWriter struct ***REMOVED***
	rws *responseWriterState
***REMOVED***

// Optional http.ResponseWriter interfaces implemented.
var (
	_ http.CloseNotifier = (*responseWriter)(nil)
	_ http.Flusher       = (*responseWriter)(nil)
	_ stringWriter       = (*responseWriter)(nil)
)

type responseWriterState struct ***REMOVED***
	// immutable within a request:
	stream *stream
	req    *http.Request
	body   *requestBody // to close at end of request, if DATA frames didn't
	conn   *serverConn

	// TODO: adjust buffer writing sizes based on server config, frame size updates from peer, etc
	bw *bufio.Writer // writing to a chunkWriter***REMOVED***this *responseWriterState***REMOVED***

	// mutated by http.Handler goroutine:
	handlerHeader http.Header // nil until called
	snapHeader    http.Header // snapshot of handlerHeader at WriteHeader time
	trailers      []string    // set in writeChunk
	status        int         // status code passed to WriteHeader
	wroteHeader   bool        // WriteHeader called (explicitly or implicitly). Not necessarily sent to user yet.
	sentHeader    bool        // have we sent the header frame?
	handlerDone   bool        // handler has finished
	dirty         bool        // a Write failed; don't reuse this responseWriterState

	sentContentLen int64 // non-zero if handler set a Content-Length header
	wroteBytes     int64

	closeNotifierMu sync.Mutex // guards closeNotifierCh
	closeNotifierCh chan bool  // nil until first used
***REMOVED***

type chunkWriter struct***REMOVED*** rws *responseWriterState ***REMOVED***

func (cw chunkWriter) Write(p []byte) (n int, err error) ***REMOVED*** return cw.rws.writeChunk(p) ***REMOVED***

func (rws *responseWriterState) hasTrailers() bool ***REMOVED*** return len(rws.trailers) > 0 ***REMOVED***

func (rws *responseWriterState) hasNonemptyTrailers() bool ***REMOVED***
	for _, trailer := range rws.trailers ***REMOVED***
		if _, ok := rws.handlerHeader[trailer]; ok ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// declareTrailer is called for each Trailer header when the
// response header is written. It notes that a header will need to be
// written in the trailers at the end of the response.
func (rws *responseWriterState) declareTrailer(k string) ***REMOVED***
	k = http.CanonicalHeaderKey(k)
	if !httpguts.ValidTrailerHeader(k) ***REMOVED***
		// Forbidden by RFC 7230, section 4.1.2.
		rws.conn.logf("ignoring invalid trailer %q", k)
		return
	***REMOVED***
	if !strSliceContains(rws.trailers, k) ***REMOVED***
		rws.trailers = append(rws.trailers, k)
	***REMOVED***
***REMOVED***

// writeChunk writes chunks from the bufio.Writer. But because
// bufio.Writer may bypass its chunking, sometimes p may be
// arbitrarily large.
//
// writeChunk is also responsible (on the first chunk) for sending the
// HEADER response.
func (rws *responseWriterState) writeChunk(p []byte) (n int, err error) ***REMOVED***
	if !rws.wroteHeader ***REMOVED***
		rws.writeHeader(200)
	***REMOVED***

	isHeadResp := rws.req.Method == "HEAD"
	if !rws.sentHeader ***REMOVED***
		rws.sentHeader = true
		var ctype, clen string
		if clen = rws.snapHeader.Get("Content-Length"); clen != "" ***REMOVED***
			rws.snapHeader.Del("Content-Length")
			if cl, err := strconv.ParseUint(clen, 10, 63); err == nil ***REMOVED***
				rws.sentContentLen = int64(cl)
			***REMOVED*** else ***REMOVED***
				clen = ""
			***REMOVED***
		***REMOVED***
		if clen == "" && rws.handlerDone && bodyAllowedForStatus(rws.status) && (len(p) > 0 || !isHeadResp) ***REMOVED***
			clen = strconv.Itoa(len(p))
		***REMOVED***
		_, hasContentType := rws.snapHeader["Content-Type"]
		// If the Content-Encoding is non-blank, we shouldn't
		// sniff the body. See Issue golang.org/issue/31753.
		ce := rws.snapHeader.Get("Content-Encoding")
		hasCE := len(ce) > 0
		if !hasCE && !hasContentType && bodyAllowedForStatus(rws.status) && len(p) > 0 ***REMOVED***
			ctype = http.DetectContentType(p)
		***REMOVED***
		var date string
		if _, ok := rws.snapHeader["Date"]; !ok ***REMOVED***
			// TODO(bradfitz): be faster here, like net/http? measure.
			date = time.Now().UTC().Format(http.TimeFormat)
		***REMOVED***

		for _, v := range rws.snapHeader["Trailer"] ***REMOVED***
			foreachHeaderElement(v, rws.declareTrailer)
		***REMOVED***

		// "Connection" headers aren't allowed in HTTP/2 (RFC 7540, 8.1.2.2),
		// but respect "Connection" == "close" to mean sending a GOAWAY and tearing
		// down the TCP connection when idle, like we do for HTTP/1.
		// TODO: remove more Connection-specific header fields here, in addition
		// to "Connection".
		if _, ok := rws.snapHeader["Connection"]; ok ***REMOVED***
			v := rws.snapHeader.Get("Connection")
			delete(rws.snapHeader, "Connection")
			if v == "close" ***REMOVED***
				rws.conn.startGracefulShutdown()
			***REMOVED***
		***REMOVED***

		endStream := (rws.handlerDone && !rws.hasTrailers() && len(p) == 0) || isHeadResp
		err = rws.conn.writeHeaders(rws.stream, &writeResHeaders***REMOVED***
			streamID:      rws.stream.id,
			httpResCode:   rws.status,
			h:             rws.snapHeader,
			endStream:     endStream,
			contentType:   ctype,
			contentLength: clen,
			date:          date,
		***REMOVED***)
		if err != nil ***REMOVED***
			rws.dirty = true
			return 0, err
		***REMOVED***
		if endStream ***REMOVED***
			return 0, nil
		***REMOVED***
	***REMOVED***
	if isHeadResp ***REMOVED***
		return len(p), nil
	***REMOVED***
	if len(p) == 0 && !rws.handlerDone ***REMOVED***
		return 0, nil
	***REMOVED***

	if rws.handlerDone ***REMOVED***
		rws.promoteUndeclaredTrailers()
	***REMOVED***

	// only send trailers if they have actually been defined by the
	// server handler.
	hasNonemptyTrailers := rws.hasNonemptyTrailers()
	endStream := rws.handlerDone && !hasNonemptyTrailers
	if len(p) > 0 || endStream ***REMOVED***
		// only send a 0 byte DATA frame if we're ending the stream.
		if err := rws.conn.writeDataFromHandler(rws.stream, p, endStream); err != nil ***REMOVED***
			rws.dirty = true
			return 0, err
		***REMOVED***
	***REMOVED***

	if rws.handlerDone && hasNonemptyTrailers ***REMOVED***
		err = rws.conn.writeHeaders(rws.stream, &writeResHeaders***REMOVED***
			streamID:  rws.stream.id,
			h:         rws.handlerHeader,
			trailers:  rws.trailers,
			endStream: true,
		***REMOVED***)
		if err != nil ***REMOVED***
			rws.dirty = true
		***REMOVED***
		return len(p), err
	***REMOVED***
	return len(p), nil
***REMOVED***

// TrailerPrefix is a magic prefix for ResponseWriter.Header map keys
// that, if present, signals that the map entry is actually for
// the response trailers, and not the response headers. The prefix
// is stripped after the ServeHTTP call finishes and the values are
// sent in the trailers.
//
// This mechanism is intended only for trailers that are not known
// prior to the headers being written. If the set of trailers is fixed
// or known before the header is written, the normal Go trailers mechanism
// is preferred:
//    https://golang.org/pkg/net/http/#ResponseWriter
//    https://golang.org/pkg/net/http/#example_ResponseWriter_trailers
const TrailerPrefix = "Trailer:"

// promoteUndeclaredTrailers permits http.Handlers to set trailers
// after the header has already been flushed. Because the Go
// ResponseWriter interface has no way to set Trailers (only the
// Header), and because we didn't want to expand the ResponseWriter
// interface, and because nobody used trailers, and because RFC 7230
// says you SHOULD (but not must) predeclare any trailers in the
// header, the official ResponseWriter rules said trailers in Go must
// be predeclared, and then we reuse the same ResponseWriter.Header()
// map to mean both Headers and Trailers. When it's time to write the
// Trailers, we pick out the fields of Headers that were declared as
// trailers. That worked for a while, until we found the first major
// user of Trailers in the wild: gRPC (using them only over http2),
// and gRPC libraries permit setting trailers mid-stream without
// predeclaring them. So: change of plans. We still permit the old
// way, but we also permit this hack: if a Header() key begins with
// "Trailer:", the suffix of that key is a Trailer. Because ':' is an
// invalid token byte anyway, there is no ambiguity. (And it's already
// filtered out) It's mildly hacky, but not terrible.
//
// This method runs after the Handler is done and promotes any Header
// fields to be trailers.
func (rws *responseWriterState) promoteUndeclaredTrailers() ***REMOVED***
	for k, vv := range rws.handlerHeader ***REMOVED***
		if !strings.HasPrefix(k, TrailerPrefix) ***REMOVED***
			continue
		***REMOVED***
		trailerKey := strings.TrimPrefix(k, TrailerPrefix)
		rws.declareTrailer(trailerKey)
		rws.handlerHeader[http.CanonicalHeaderKey(trailerKey)] = vv
	***REMOVED***

	if len(rws.trailers) > 1 ***REMOVED***
		sorter := sorterPool.Get().(*sorter)
		sorter.SortStrings(rws.trailers)
		sorterPool.Put(sorter)
	***REMOVED***
***REMOVED***

func (w *responseWriter) Flush() ***REMOVED***
	rws := w.rws
	if rws == nil ***REMOVED***
		panic("Header called after Handler finished")
	***REMOVED***
	if rws.bw.Buffered() > 0 ***REMOVED***
		if err := rws.bw.Flush(); err != nil ***REMOVED***
			// Ignore the error. The frame writer already knows.
			return
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// The bufio.Writer won't call chunkWriter.Write
		// (writeChunk with zero bytes, so we have to do it
		// ourselves to force the HTTP response header and/or
		// final DATA frame (with END_STREAM) to be sent.
		rws.writeChunk(nil)
	***REMOVED***
***REMOVED***

func (w *responseWriter) CloseNotify() <-chan bool ***REMOVED***
	rws := w.rws
	if rws == nil ***REMOVED***
		panic("CloseNotify called after Handler finished")
	***REMOVED***
	rws.closeNotifierMu.Lock()
	ch := rws.closeNotifierCh
	if ch == nil ***REMOVED***
		ch = make(chan bool, 1)
		rws.closeNotifierCh = ch
		cw := rws.stream.cw
		go func() ***REMOVED***
			cw.Wait() // wait for close
			ch <- true
		***REMOVED***()
	***REMOVED***
	rws.closeNotifierMu.Unlock()
	return ch
***REMOVED***

func (w *responseWriter) Header() http.Header ***REMOVED***
	rws := w.rws
	if rws == nil ***REMOVED***
		panic("Header called after Handler finished")
	***REMOVED***
	if rws.handlerHeader == nil ***REMOVED***
		rws.handlerHeader = make(http.Header)
	***REMOVED***
	return rws.handlerHeader
***REMOVED***

// checkWriteHeaderCode is a copy of net/http's checkWriteHeaderCode.
func checkWriteHeaderCode(code int) ***REMOVED***
	// Issue 22880: require valid WriteHeader status codes.
	// For now we only enforce that it's three digits.
	// In the future we might block things over 599 (600 and above aren't defined
	// at http://httpwg.org/specs/rfc7231.html#status.codes)
	// and we might block under 200 (once we have more mature 1xx support).
	// But for now any three digits.
	//
	// We used to send "HTTP/1.1 000 0" on the wire in responses but there's
	// no equivalent bogus thing we can realistically send in HTTP/2,
	// so we'll consistently panic instead and help people find their bugs
	// early. (We can't return an error from WriteHeader even if we wanted to.)
	if code < 100 || code > 999 ***REMOVED***
		panic(fmt.Sprintf("invalid WriteHeader code %v", code))
	***REMOVED***
***REMOVED***

func (w *responseWriter) WriteHeader(code int) ***REMOVED***
	rws := w.rws
	if rws == nil ***REMOVED***
		panic("WriteHeader called after Handler finished")
	***REMOVED***
	rws.writeHeader(code)
***REMOVED***

func (rws *responseWriterState) writeHeader(code int) ***REMOVED***
	if !rws.wroteHeader ***REMOVED***
		checkWriteHeaderCode(code)
		rws.wroteHeader = true
		rws.status = code
		if len(rws.handlerHeader) > 0 ***REMOVED***
			rws.snapHeader = cloneHeader(rws.handlerHeader)
		***REMOVED***
	***REMOVED***
***REMOVED***

func cloneHeader(h http.Header) http.Header ***REMOVED***
	h2 := make(http.Header, len(h))
	for k, vv := range h ***REMOVED***
		vv2 := make([]string, len(vv))
		copy(vv2, vv)
		h2[k] = vv2
	***REMOVED***
	return h2
***REMOVED***

// The Life Of A Write is like this:
//
// * Handler calls w.Write or w.WriteString ->
// * -> rws.bw (*bufio.Writer) ->
// * (Handler might call Flush)
// * -> chunkWriter***REMOVED***rws***REMOVED***
// * -> responseWriterState.writeChunk(p []byte)
// * -> responseWriterState.writeChunk (most of the magic; see comment there)
func (w *responseWriter) Write(p []byte) (n int, err error) ***REMOVED***
	return w.write(len(p), p, "")
***REMOVED***

func (w *responseWriter) WriteString(s string) (n int, err error) ***REMOVED***
	return w.write(len(s), nil, s)
***REMOVED***

// either dataB or dataS is non-zero.
func (w *responseWriter) write(lenData int, dataB []byte, dataS string) (n int, err error) ***REMOVED***
	rws := w.rws
	if rws == nil ***REMOVED***
		panic("Write called after Handler finished")
	***REMOVED***
	if !rws.wroteHeader ***REMOVED***
		w.WriteHeader(200)
	***REMOVED***
	if !bodyAllowedForStatus(rws.status) ***REMOVED***
		return 0, http.ErrBodyNotAllowed
	***REMOVED***
	rws.wroteBytes += int64(len(dataB)) + int64(len(dataS)) // only one can be set
	if rws.sentContentLen != 0 && rws.wroteBytes > rws.sentContentLen ***REMOVED***
		// TODO: send a RST_STREAM
		return 0, errors.New("http2: handler wrote more than declared Content-Length")
	***REMOVED***

	if dataB != nil ***REMOVED***
		return rws.bw.Write(dataB)
	***REMOVED*** else ***REMOVED***
		return rws.bw.WriteString(dataS)
	***REMOVED***
***REMOVED***

func (w *responseWriter) handlerDone() ***REMOVED***
	rws := w.rws
	dirty := rws.dirty
	rws.handlerDone = true
	w.Flush()
	w.rws = nil
	if !dirty ***REMOVED***
		// Only recycle the pool if all prior Write calls to
		// the serverConn goroutine completed successfully. If
		// they returned earlier due to resets from the peer
		// there might still be write goroutines outstanding
		// from the serverConn referencing the rws memory. See
		// issue 20704.
		responseWriterStatePool.Put(rws)
	***REMOVED***
***REMOVED***

// Push errors.
var (
	ErrRecursivePush    = errors.New("http2: recursive push not allowed")
	ErrPushLimitReached = errors.New("http2: push would exceed peer's SETTINGS_MAX_CONCURRENT_STREAMS")
)

var _ http.Pusher = (*responseWriter)(nil)

func (w *responseWriter) Push(target string, opts *http.PushOptions) error ***REMOVED***
	st := w.rws.stream
	sc := st.sc
	sc.serveG.checkNotOn()

	// No recursive pushes: "PUSH_PROMISE frames MUST only be sent on a peer-initiated stream."
	// http://tools.ietf.org/html/rfc7540#section-6.6
	if st.isPushed() ***REMOVED***
		return ErrRecursivePush
	***REMOVED***

	if opts == nil ***REMOVED***
		opts = new(http.PushOptions)
	***REMOVED***

	// Default options.
	if opts.Method == "" ***REMOVED***
		opts.Method = "GET"
	***REMOVED***
	if opts.Header == nil ***REMOVED***
		opts.Header = http.Header***REMOVED******REMOVED***
	***REMOVED***
	wantScheme := "http"
	if w.rws.req.TLS != nil ***REMOVED***
		wantScheme = "https"
	***REMOVED***

	// Validate the request.
	u, err := url.Parse(target)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if u.Scheme == "" ***REMOVED***
		if !strings.HasPrefix(target, "/") ***REMOVED***
			return fmt.Errorf("target must be an absolute URL or an absolute path: %q", target)
		***REMOVED***
		u.Scheme = wantScheme
		u.Host = w.rws.req.Host
	***REMOVED*** else ***REMOVED***
		if u.Scheme != wantScheme ***REMOVED***
			return fmt.Errorf("cannot push URL with scheme %q from request with scheme %q", u.Scheme, wantScheme)
		***REMOVED***
		if u.Host == "" ***REMOVED***
			return errors.New("URL must have a host")
		***REMOVED***
	***REMOVED***
	for k := range opts.Header ***REMOVED***
		if strings.HasPrefix(k, ":") ***REMOVED***
			return fmt.Errorf("promised request headers cannot include pseudo header %q", k)
		***REMOVED***
		// These headers are meaningful only if the request has a body,
		// but PUSH_PROMISE requests cannot have a body.
		// http://tools.ietf.org/html/rfc7540#section-8.2
		// Also disallow Host, since the promised URL must be absolute.
		switch strings.ToLower(k) ***REMOVED***
		case "content-length", "content-encoding", "trailer", "te", "expect", "host":
			return fmt.Errorf("promised request headers cannot include %q", k)
		***REMOVED***
	***REMOVED***
	if err := checkValidHTTP2RequestHeaders(opts.Header); err != nil ***REMOVED***
		return err
	***REMOVED***

	// The RFC effectively limits promised requests to GET and HEAD:
	// "Promised requests MUST be cacheable [GET, HEAD, or POST], and MUST be safe [GET or HEAD]"
	// http://tools.ietf.org/html/rfc7540#section-8.2
	if opts.Method != "GET" && opts.Method != "HEAD" ***REMOVED***
		return fmt.Errorf("method %q must be GET or HEAD", opts.Method)
	***REMOVED***

	msg := &startPushRequest***REMOVED***
		parent: st,
		method: opts.Method,
		url:    u,
		header: cloneHeader(opts.Header),
		done:   errChanPool.Get().(chan error),
	***REMOVED***

	select ***REMOVED***
	case <-sc.doneServing:
		return errClientDisconnected
	case <-st.cw:
		return errStreamClosed
	case sc.serveMsgCh <- msg:
	***REMOVED***

	select ***REMOVED***
	case <-sc.doneServing:
		return errClientDisconnected
	case <-st.cw:
		return errStreamClosed
	case err := <-msg.done:
		errChanPool.Put(msg.done)
		return err
	***REMOVED***
***REMOVED***

type startPushRequest struct ***REMOVED***
	parent *stream
	method string
	url    *url.URL
	header http.Header
	done   chan error
***REMOVED***

func (sc *serverConn) startPush(msg *startPushRequest) ***REMOVED***
	sc.serveG.check()

	// http://tools.ietf.org/html/rfc7540#section-6.6.
	// PUSH_PROMISE frames MUST only be sent on a peer-initiated stream that
	// is in either the "open" or "half-closed (remote)" state.
	if msg.parent.state != stateOpen && msg.parent.state != stateHalfClosedRemote ***REMOVED***
		// responseWriter.Push checks that the stream is peer-initiated.
		msg.done <- errStreamClosed
		return
	***REMOVED***

	// http://tools.ietf.org/html/rfc7540#section-6.6.
	if !sc.pushEnabled ***REMOVED***
		msg.done <- http.ErrNotSupported
		return
	***REMOVED***

	// PUSH_PROMISE frames must be sent in increasing order by stream ID, so
	// we allocate an ID for the promised stream lazily, when the PUSH_PROMISE
	// is written. Once the ID is allocated, we start the request handler.
	allocatePromisedID := func() (uint32, error) ***REMOVED***
		sc.serveG.check()

		// Check this again, just in case. Technically, we might have received
		// an updated SETTINGS by the time we got around to writing this frame.
		if !sc.pushEnabled ***REMOVED***
			return 0, http.ErrNotSupported
		***REMOVED***
		// http://tools.ietf.org/html/rfc7540#section-6.5.2.
		if sc.curPushedStreams+1 > sc.clientMaxStreams ***REMOVED***
			return 0, ErrPushLimitReached
		***REMOVED***

		// http://tools.ietf.org/html/rfc7540#section-5.1.1.
		// Streams initiated by the server MUST use even-numbered identifiers.
		// A server that is unable to establish a new stream identifier can send a GOAWAY
		// frame so that the client is forced to open a new connection for new streams.
		if sc.maxPushPromiseID+2 >= 1<<31 ***REMOVED***
			sc.startGracefulShutdownInternal()
			return 0, ErrPushLimitReached
		***REMOVED***
		sc.maxPushPromiseID += 2
		promisedID := sc.maxPushPromiseID

		// http://tools.ietf.org/html/rfc7540#section-8.2.
		// Strictly speaking, the new stream should start in "reserved (local)", then
		// transition to "half closed (remote)" after sending the initial HEADERS, but
		// we start in "half closed (remote)" for simplicity.
		// See further comments at the definition of stateHalfClosedRemote.
		promised := sc.newStream(promisedID, msg.parent.id, stateHalfClosedRemote)
		rw, req, err := sc.newWriterAndRequestNoBody(promised, requestParam***REMOVED***
			method:    msg.method,
			scheme:    msg.url.Scheme,
			authority: msg.url.Host,
			path:      msg.url.RequestURI(),
			header:    cloneHeader(msg.header), // clone since handler runs concurrently with writing the PUSH_PROMISE
		***REMOVED***)
		if err != nil ***REMOVED***
			// Should not happen, since we've already validated msg.url.
			panic(fmt.Sprintf("newWriterAndRequestNoBody(%+v): %v", msg.url, err))
		***REMOVED***

		go sc.runHandler(rw, req, sc.handler.ServeHTTP)
		return promisedID, nil
	***REMOVED***

	sc.writeFrame(FrameWriteRequest***REMOVED***
		write: &writePushPromise***REMOVED***
			streamID:           msg.parent.id,
			method:             msg.method,
			url:                msg.url,
			h:                  msg.header,
			allocatePromisedID: allocatePromisedID,
		***REMOVED***,
		stream: msg.parent,
		done:   msg.done,
	***REMOVED***)
***REMOVED***

// foreachHeaderElement splits v according to the "#rule" construction
// in RFC 7230 section 7 and calls fn for each non-empty element.
func foreachHeaderElement(v string, fn func(string)) ***REMOVED***
	v = textproto.TrimString(v)
	if v == "" ***REMOVED***
		return
	***REMOVED***
	if !strings.Contains(v, ",") ***REMOVED***
		fn(v)
		return
	***REMOVED***
	for _, f := range strings.Split(v, ",") ***REMOVED***
		if f = textproto.TrimString(f); f != "" ***REMOVED***
			fn(f)
		***REMOVED***
	***REMOVED***
***REMOVED***

// From http://httpwg.org/specs/rfc7540.html#rfc.section.8.1.2.2
var connHeaders = []string***REMOVED***
	"Connection",
	"Keep-Alive",
	"Proxy-Connection",
	"Transfer-Encoding",
	"Upgrade",
***REMOVED***

// checkValidHTTP2RequestHeaders checks whether h is a valid HTTP/2 request,
// per RFC 7540 Section 8.1.2.2.
// The returned error is reported to users.
func checkValidHTTP2RequestHeaders(h http.Header) error ***REMOVED***
	for _, k := range connHeaders ***REMOVED***
		if _, ok := h[k]; ok ***REMOVED***
			return fmt.Errorf("request header %q is not valid in HTTP/2", k)
		***REMOVED***
	***REMOVED***
	te := h["Te"]
	if len(te) > 0 && (len(te) > 1 || (te[0] != "trailers" && te[0] != "")) ***REMOVED***
		return errors.New(`request header "TE" may only be "trailers" in HTTP/2`)
	***REMOVED***
	return nil
***REMOVED***

func new400Handler(err error) http.HandlerFunc ***REMOVED***
	return func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		http.Error(w, err.Error(), http.StatusBadRequest)
	***REMOVED***
***REMOVED***

// h1ServerKeepAlivesDisabled reports whether hs has its keep-alives
// disabled. See comments on h1ServerShutdownChan above for why
// the code is written this way.
func h1ServerKeepAlivesDisabled(hs *http.Server) bool ***REMOVED***
	var x interface***REMOVED******REMOVED*** = hs
	type I interface ***REMOVED***
		doKeepAlives() bool
	***REMOVED***
	if hs, ok := x.(I); ok ***REMOVED***
		return !hs.doKeepAlives()
	***REMOVED***
	return false
***REMOVED***
