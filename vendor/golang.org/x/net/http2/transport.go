// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Transport code.

package http2

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rand"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	mathrand "math/rand"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/textproto"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/net/http/httpguts"
	"golang.org/x/net/http2/hpack"
	"golang.org/x/net/idna"
)

const (
	// transportDefaultConnFlow is how many connection-level flow control
	// tokens we give the server at start-up, past the default 64k.
	transportDefaultConnFlow = 1 << 30

	// transportDefaultStreamFlow is how many stream-level flow
	// control tokens we announce to the peer, and how many bytes
	// we buffer per stream.
	transportDefaultStreamFlow = 4 << 20

	// transportDefaultStreamMinRefresh is the minimum number of bytes we'll send
	// a stream-level WINDOW_UPDATE for at a time.
	transportDefaultStreamMinRefresh = 4 << 10

	defaultUserAgent = "Go-http-client/2.0"

	// initialMaxConcurrentStreams is a connections maxConcurrentStreams until
	// it's received servers initial SETTINGS frame, which corresponds with the
	// spec's minimum recommended value.
	initialMaxConcurrentStreams = 100

	// defaultMaxConcurrentStreams is a connections default maxConcurrentStreams
	// if the server doesn't include one in its initial SETTINGS frame.
	defaultMaxConcurrentStreams = 1000
)

// Transport is an HTTP/2 Transport.
//
// A Transport internally caches connections to servers. It is safe
// for concurrent use by multiple goroutines.
type Transport struct ***REMOVED***
	// DialTLS specifies an optional dial function for creating
	// TLS connections for requests.
	//
	// If DialTLS is nil, tls.Dial is used.
	//
	// If the returned net.Conn has a ConnectionState method like tls.Conn,
	// it will be used to set http.Response.TLS.
	DialTLS func(network, addr string, cfg *tls.Config) (net.Conn, error)

	// TLSClientConfig specifies the TLS configuration to use with
	// tls.Client. If nil, the default configuration is used.
	TLSClientConfig *tls.Config

	// ConnPool optionally specifies an alternate connection pool to use.
	// If nil, the default is used.
	ConnPool ClientConnPool

	// DisableCompression, if true, prevents the Transport from
	// requesting compression with an "Accept-Encoding: gzip"
	// request header when the Request contains no existing
	// Accept-Encoding value. If the Transport requests gzip on
	// its own and gets a gzipped response, it's transparently
	// decoded in the Response.Body. However, if the user
	// explicitly requested gzip it is not automatically
	// uncompressed.
	DisableCompression bool

	// AllowHTTP, if true, permits HTTP/2 requests using the insecure,
	// plain-text "http" scheme. Note that this does not enable h2c support.
	AllowHTTP bool

	// MaxHeaderListSize is the http2 SETTINGS_MAX_HEADER_LIST_SIZE to
	// send in the initial settings frame. It is how many bytes
	// of response headers are allowed. Unlike the http2 spec, zero here
	// means to use a default limit (currently 10MB). If you actually
	// want to advertise an unlimited value to the peer, Transport
	// interprets the highest possible value here (0xffffffff or 1<<32-1)
	// to mean no limit.
	MaxHeaderListSize uint32

	// StrictMaxConcurrentStreams controls whether the server's
	// SETTINGS_MAX_CONCURRENT_STREAMS should be respected
	// globally. If false, new TCP connections are created to the
	// server as needed to keep each under the per-connection
	// SETTINGS_MAX_CONCURRENT_STREAMS limit. If true, the
	// server's SETTINGS_MAX_CONCURRENT_STREAMS is interpreted as
	// a global limit and callers of RoundTrip block when needed,
	// waiting for their turn.
	StrictMaxConcurrentStreams bool

	// ReadIdleTimeout is the timeout after which a health check using ping
	// frame will be carried out if no frame is received on the connection.
	// Note that a ping response will is considered a received frame, so if
	// there is no other traffic on the connection, the health check will
	// be performed every ReadIdleTimeout interval.
	// If zero, no health check is performed.
	ReadIdleTimeout time.Duration

	// PingTimeout is the timeout after which the connection will be closed
	// if a response to Ping is not received.
	// Defaults to 15s.
	PingTimeout time.Duration

	// WriteByteTimeout is the timeout after which the connection will be
	// closed no data can be written to it. The timeout begins when data is
	// available to write, and is extended whenever any bytes are written.
	WriteByteTimeout time.Duration

	// CountError, if non-nil, is called on HTTP/2 transport errors.
	// It's intended to increment a metric for monitoring, such
	// as an expvar or Prometheus metric.
	// The errType consists of only ASCII word characters.
	CountError func(errType string)

	// t1, if non-nil, is the standard library Transport using
	// this transport. Its settings are used (but not its
	// RoundTrip method, etc).
	t1 *http.Transport

	connPoolOnce  sync.Once
	connPoolOrDef ClientConnPool // non-nil version of ConnPool
***REMOVED***

func (t *Transport) maxHeaderListSize() uint32 ***REMOVED***
	if t.MaxHeaderListSize == 0 ***REMOVED***
		return 10 << 20
	***REMOVED***
	if t.MaxHeaderListSize == 0xffffffff ***REMOVED***
		return 0
	***REMOVED***
	return t.MaxHeaderListSize
***REMOVED***

func (t *Transport) disableCompression() bool ***REMOVED***
	return t.DisableCompression || (t.t1 != nil && t.t1.DisableCompression)
***REMOVED***

func (t *Transport) pingTimeout() time.Duration ***REMOVED***
	if t.PingTimeout == 0 ***REMOVED***
		return 15 * time.Second
	***REMOVED***
	return t.PingTimeout

***REMOVED***

// ConfigureTransport configures a net/http HTTP/1 Transport to use HTTP/2.
// It returns an error if t1 has already been HTTP/2-enabled.
//
// Use ConfigureTransports instead to configure the HTTP/2 Transport.
func ConfigureTransport(t1 *http.Transport) error ***REMOVED***
	_, err := ConfigureTransports(t1)
	return err
***REMOVED***

// ConfigureTransports configures a net/http HTTP/1 Transport to use HTTP/2.
// It returns a new HTTP/2 Transport for further configuration.
// It returns an error if t1 has already been HTTP/2-enabled.
func ConfigureTransports(t1 *http.Transport) (*Transport, error) ***REMOVED***
	return configureTransports(t1)
***REMOVED***

func configureTransports(t1 *http.Transport) (*Transport, error) ***REMOVED***
	connPool := new(clientConnPool)
	t2 := &Transport***REMOVED***
		ConnPool: noDialClientConnPool***REMOVED***connPool***REMOVED***,
		t1:       t1,
	***REMOVED***
	connPool.t = t2
	if err := registerHTTPSProtocol(t1, noDialH2RoundTripper***REMOVED***t2***REMOVED***); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if t1.TLSClientConfig == nil ***REMOVED***
		t1.TLSClientConfig = new(tls.Config)
	***REMOVED***
	if !strSliceContains(t1.TLSClientConfig.NextProtos, "h2") ***REMOVED***
		t1.TLSClientConfig.NextProtos = append([]string***REMOVED***"h2"***REMOVED***, t1.TLSClientConfig.NextProtos...)
	***REMOVED***
	if !strSliceContains(t1.TLSClientConfig.NextProtos, "http/1.1") ***REMOVED***
		t1.TLSClientConfig.NextProtos = append(t1.TLSClientConfig.NextProtos, "http/1.1")
	***REMOVED***
	upgradeFn := func(authority string, c *tls.Conn) http.RoundTripper ***REMOVED***
		addr := authorityAddr("https", authority)
		if used, err := connPool.addConnIfNeeded(addr, t2, c); err != nil ***REMOVED***
			go c.Close()
			return erringRoundTripper***REMOVED***err***REMOVED***
		***REMOVED*** else if !used ***REMOVED***
			// Turns out we don't need this c.
			// For example, two goroutines made requests to the same host
			// at the same time, both kicking off TCP dials. (since protocol
			// was unknown)
			go c.Close()
		***REMOVED***
		return t2
	***REMOVED***
	if m := t1.TLSNextProto; len(m) == 0 ***REMOVED***
		t1.TLSNextProto = map[string]func(string, *tls.Conn) http.RoundTripper***REMOVED***
			"h2": upgradeFn,
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		m["h2"] = upgradeFn
	***REMOVED***
	return t2, nil
***REMOVED***

func (t *Transport) connPool() ClientConnPool ***REMOVED***
	t.connPoolOnce.Do(t.initConnPool)
	return t.connPoolOrDef
***REMOVED***

func (t *Transport) initConnPool() ***REMOVED***
	if t.ConnPool != nil ***REMOVED***
		t.connPoolOrDef = t.ConnPool
	***REMOVED*** else ***REMOVED***
		t.connPoolOrDef = &clientConnPool***REMOVED***t: t***REMOVED***
	***REMOVED***
***REMOVED***

// ClientConn is the state of a single HTTP/2 client connection to an
// HTTP/2 server.
type ClientConn struct ***REMOVED***
	t             *Transport
	tconn         net.Conn             // usually *tls.Conn, except specialized impls
	tlsState      *tls.ConnectionState // nil only for specialized impls
	reused        uint32               // whether conn is being reused; atomic
	singleUse     bool                 // whether being used for a single http.Request
	getConnCalled bool                 // used by clientConnPool

	// readLoop goroutine fields:
	readerDone chan struct***REMOVED******REMOVED*** // closed on error
	readerErr  error         // set before readerDone is closed

	idleTimeout time.Duration // or 0 for never
	idleTimer   *time.Timer

	mu              sync.Mutex // guards following
	cond            *sync.Cond // hold mu; broadcast on flow/closed changes
	flow            flow       // our conn-level flow control quota (cs.flow is per stream)
	inflow          flow       // peer's conn-level flow control
	doNotReuse      bool       // whether conn is marked to not be reused for any future requests
	closing         bool
	closed          bool
	seenSettings    bool                     // true if we've seen a settings frame, false otherwise
	wantSettingsAck bool                     // we sent a SETTINGS frame and haven't heard back
	goAway          *GoAwayFrame             // if non-nil, the GoAwayFrame we received
	goAwayDebug     string                   // goAway frame's debug data, retained as a string
	streams         map[uint32]*clientStream // client-initiated
	streamsReserved int                      // incr by ReserveNewRequest; decr on RoundTrip
	nextStreamID    uint32
	pendingRequests int                       // requests blocked and waiting to be sent because len(streams) == maxConcurrentStreams
	pings           map[[8]byte]chan struct***REMOVED******REMOVED*** // in flight ping data to notification channel
	br              *bufio.Reader
	lastActive      time.Time
	lastIdle        time.Time // time last idle
	// Settings from peer: (also guarded by wmu)
	maxFrameSize          uint32
	maxConcurrentStreams  uint32
	peerMaxHeaderListSize uint64
	initialWindowSize     uint32

	// reqHeaderMu is a 1-element semaphore channel controlling access to sending new requests.
	// Write to reqHeaderMu to lock it, read from it to unlock.
	// Lock reqmu BEFORE mu or wmu.
	reqHeaderMu chan struct***REMOVED******REMOVED***

	// wmu is held while writing.
	// Acquire BEFORE mu when holding both, to avoid blocking mu on network writes.
	// Only acquire both at the same time when changing peer settings.
	wmu  sync.Mutex
	bw   *bufio.Writer
	fr   *Framer
	werr error        // first write error that has occurred
	hbuf bytes.Buffer // HPACK encoder writes into this
	henc *hpack.Encoder
***REMOVED***

// clientStream is the state for a single HTTP/2 stream. One of these
// is created for each Transport.RoundTrip call.
type clientStream struct ***REMOVED***
	cc *ClientConn

	// Fields of Request that we may access even after the response body is closed.
	ctx       context.Context
	reqCancel <-chan struct***REMOVED******REMOVED***

	trace         *httptrace.ClientTrace // or nil
	ID            uint32
	bufPipe       pipe // buffered pipe with the flow-controlled response payload
	requestedGzip bool
	isHead        bool

	abortOnce sync.Once
	abort     chan struct***REMOVED******REMOVED*** // closed to signal stream should end immediately
	abortErr  error         // set if abort is closed

	peerClosed chan struct***REMOVED******REMOVED*** // closed when the peer sends an END_STREAM flag
	donec      chan struct***REMOVED******REMOVED*** // closed after the stream is in the closed state
	on100      chan struct***REMOVED******REMOVED*** // buffered; written to if a 100 is received

	respHeaderRecv chan struct***REMOVED******REMOVED***  // closed when headers are received
	res            *http.Response // set if respHeaderRecv is closed

	flow        flow  // guarded by cc.mu
	inflow      flow  // guarded by cc.mu
	bytesRemain int64 // -1 means unknown; owned by transportResponseBody.Read
	readErr     error // sticky read error; owned by transportResponseBody.Read

	reqBody              io.ReadCloser
	reqBodyContentLength int64 // -1 means unknown
	reqBodyClosed        bool  // body has been closed; guarded by cc.mu

	// owned by writeRequest:
	sentEndStream bool // sent an END_STREAM flag to the peer
	sentHeaders   bool

	// owned by clientConnReadLoop:
	firstByte    bool  // got the first response byte
	pastHeaders  bool  // got first MetaHeadersFrame (actual headers)
	pastTrailers bool  // got optional second MetaHeadersFrame (trailers)
	num1xx       uint8 // number of 1xx responses seen
	readClosed   bool  // peer sent an END_STREAM flag
	readAborted  bool  // read loop reset the stream

	trailer    http.Header  // accumulated trailers
	resTrailer *http.Header // client's Response.Trailer
***REMOVED***

var got1xxFuncForTests func(int, textproto.MIMEHeader) error

// get1xxTraceFunc returns the value of request's httptrace.ClientTrace.Got1xxResponse func,
// if any. It returns nil if not set or if the Go version is too old.
func (cs *clientStream) get1xxTraceFunc() func(int, textproto.MIMEHeader) error ***REMOVED***
	if fn := got1xxFuncForTests; fn != nil ***REMOVED***
		return fn
	***REMOVED***
	return traceGot1xxResponseFunc(cs.trace)
***REMOVED***

func (cs *clientStream) abortStream(err error) ***REMOVED***
	cs.cc.mu.Lock()
	defer cs.cc.mu.Unlock()
	cs.abortStreamLocked(err)
***REMOVED***

func (cs *clientStream) abortStreamLocked(err error) ***REMOVED***
	cs.abortOnce.Do(func() ***REMOVED***
		cs.abortErr = err
		close(cs.abort)
	***REMOVED***)
	if cs.reqBody != nil && !cs.reqBodyClosed ***REMOVED***
		cs.reqBody.Close()
		cs.reqBodyClosed = true
	***REMOVED***
	// TODO(dneil): Clean up tests where cs.cc.cond is nil.
	if cs.cc.cond != nil ***REMOVED***
		// Wake up writeRequestBody if it is waiting on flow control.
		cs.cc.cond.Broadcast()
	***REMOVED***
***REMOVED***

func (cs *clientStream) abortRequestBodyWrite() ***REMOVED***
	cc := cs.cc
	cc.mu.Lock()
	defer cc.mu.Unlock()
	if cs.reqBody != nil && !cs.reqBodyClosed ***REMOVED***
		cs.reqBody.Close()
		cs.reqBodyClosed = true
		cc.cond.Broadcast()
	***REMOVED***
***REMOVED***

type stickyErrWriter struct ***REMOVED***
	conn    net.Conn
	timeout time.Duration
	err     *error
***REMOVED***

func (sew stickyErrWriter) Write(p []byte) (n int, err error) ***REMOVED***
	if *sew.err != nil ***REMOVED***
		return 0, *sew.err
	***REMOVED***
	for ***REMOVED***
		if sew.timeout != 0 ***REMOVED***
			sew.conn.SetWriteDeadline(time.Now().Add(sew.timeout))
		***REMOVED***
		nn, err := sew.conn.Write(p[n:])
		n += nn
		if n < len(p) && nn > 0 && errors.Is(err, os.ErrDeadlineExceeded) ***REMOVED***
			// Keep extending the deadline so long as we're making progress.
			continue
		***REMOVED***
		if sew.timeout != 0 ***REMOVED***
			sew.conn.SetWriteDeadline(time.Time***REMOVED******REMOVED***)
		***REMOVED***
		*sew.err = err
		return n, err
	***REMOVED***
***REMOVED***

// noCachedConnError is the concrete type of ErrNoCachedConn, which
// needs to be detected by net/http regardless of whether it's its
// bundled version (in h2_bundle.go with a rewritten type name) or
// from a user's x/net/http2. As such, as it has a unique method name
// (IsHTTP2NoCachedConnError) that net/http sniffs for via func
// isNoCachedConnError.
type noCachedConnError struct***REMOVED******REMOVED***

func (noCachedConnError) IsHTTP2NoCachedConnError() ***REMOVED******REMOVED***
func (noCachedConnError) Error() string             ***REMOVED*** return "http2: no cached connection was available" ***REMOVED***

// isNoCachedConnError reports whether err is of type noCachedConnError
// or its equivalent renamed type in net/http2's h2_bundle.go. Both types
// may coexist in the same running program.
func isNoCachedConnError(err error) bool ***REMOVED***
	_, ok := err.(interface***REMOVED*** IsHTTP2NoCachedConnError() ***REMOVED***)
	return ok
***REMOVED***

var ErrNoCachedConn error = noCachedConnError***REMOVED******REMOVED***

// RoundTripOpt are options for the Transport.RoundTripOpt method.
type RoundTripOpt struct ***REMOVED***
	// OnlyCachedConn controls whether RoundTripOpt may
	// create a new TCP connection. If set true and
	// no cached connection is available, RoundTripOpt
	// will return ErrNoCachedConn.
	OnlyCachedConn bool
***REMOVED***

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) ***REMOVED***
	return t.RoundTripOpt(req, RoundTripOpt***REMOVED******REMOVED***)
***REMOVED***

// authorityAddr returns a given authority (a host/IP, or host:port / ip:port)
// and returns a host:port. The port 443 is added if needed.
func authorityAddr(scheme string, authority string) (addr string) ***REMOVED***
	host, port, err := net.SplitHostPort(authority)
	if err != nil ***REMOVED*** // authority didn't have a port
		port = "443"
		if scheme == "http" ***REMOVED***
			port = "80"
		***REMOVED***
		host = authority
	***REMOVED***
	if a, err := idna.ToASCII(host); err == nil ***REMOVED***
		host = a
	***REMOVED***
	// IPv6 address literal, without a port:
	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") ***REMOVED***
		return host + ":" + port
	***REMOVED***
	return net.JoinHostPort(host, port)
***REMOVED***

// RoundTripOpt is like RoundTrip, but takes options.
func (t *Transport) RoundTripOpt(req *http.Request, opt RoundTripOpt) (*http.Response, error) ***REMOVED***
	if !(req.URL.Scheme == "https" || (req.URL.Scheme == "http" && t.AllowHTTP)) ***REMOVED***
		return nil, errors.New("http2: unsupported scheme")
	***REMOVED***

	addr := authorityAddr(req.URL.Scheme, req.URL.Host)
	for retry := 0; ; retry++ ***REMOVED***
		cc, err := t.connPool().GetClientConn(req, addr)
		if err != nil ***REMOVED***
			t.vlogf("http2: Transport failed to get client conn for %s: %v", addr, err)
			return nil, err
		***REMOVED***
		reused := !atomic.CompareAndSwapUint32(&cc.reused, 0, 1)
		traceGotConn(req, cc, reused)
		res, err := cc.RoundTrip(req)
		if err != nil && retry <= 6 ***REMOVED***
			if req, err = shouldRetryRequest(req, err); err == nil ***REMOVED***
				// After the first retry, do exponential backoff with 10% jitter.
				if retry == 0 ***REMOVED***
					continue
				***REMOVED***
				backoff := float64(uint(1) << (uint(retry) - 1))
				backoff += backoff * (0.1 * mathrand.Float64())
				select ***REMOVED***
				case <-time.After(time.Second * time.Duration(backoff)):
					continue
				case <-req.Context().Done():
					err = req.Context().Err()
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if err != nil ***REMOVED***
			t.vlogf("RoundTrip failure: %v", err)
			return nil, err
		***REMOVED***
		return res, nil
	***REMOVED***
***REMOVED***

// CloseIdleConnections closes any connections which were previously
// connected from previous requests but are now sitting idle.
// It does not interrupt any connections currently in use.
func (t *Transport) CloseIdleConnections() ***REMOVED***
	if cp, ok := t.connPool().(clientConnPoolIdleCloser); ok ***REMOVED***
		cp.closeIdleConnections()
	***REMOVED***
***REMOVED***

var (
	errClientConnClosed    = errors.New("http2: client conn is closed")
	errClientConnUnusable  = errors.New("http2: client conn not usable")
	errClientConnGotGoAway = errors.New("http2: Transport received Server's graceful shutdown GOAWAY")
)

// shouldRetryRequest is called by RoundTrip when a request fails to get
// response headers. It is always called with a non-nil error.
// It returns either a request to retry (either the same request, or a
// modified clone), or an error if the request can't be replayed.
func shouldRetryRequest(req *http.Request, err error) (*http.Request, error) ***REMOVED***
	if !canRetryError(err) ***REMOVED***
		return nil, err
	***REMOVED***
	// If the Body is nil (or http.NoBody), it's safe to reuse
	// this request and its Body.
	if req.Body == nil || req.Body == http.NoBody ***REMOVED***
		return req, nil
	***REMOVED***

	// If the request body can be reset back to its original
	// state via the optional req.GetBody, do that.
	if req.GetBody != nil ***REMOVED***
		body, err := req.GetBody()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		newReq := *req
		newReq.Body = body
		return &newReq, nil
	***REMOVED***

	// The Request.Body can't reset back to the beginning, but we
	// don't seem to have started to read from it yet, so reuse
	// the request directly.
	if err == errClientConnUnusable ***REMOVED***
		return req, nil
	***REMOVED***

	return nil, fmt.Errorf("http2: Transport: cannot retry err [%v] after Request.Body was written; define Request.GetBody to avoid this error", err)
***REMOVED***

func canRetryError(err error) bool ***REMOVED***
	if err == errClientConnUnusable || err == errClientConnGotGoAway ***REMOVED***
		return true
	***REMOVED***
	if se, ok := err.(StreamError); ok ***REMOVED***
		if se.Code == ErrCodeProtocol && se.Cause == errFromPeer ***REMOVED***
			// See golang/go#47635, golang/go#42777
			return true
		***REMOVED***
		return se.Code == ErrCodeRefusedStream
	***REMOVED***
	return false
***REMOVED***

func (t *Transport) dialClientConn(ctx context.Context, addr string, singleUse bool) (*ClientConn, error) ***REMOVED***
	host, _, err := net.SplitHostPort(addr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	tconn, err := t.dialTLS(ctx)("tcp", addr, t.newTLSConfig(host))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return t.newClientConn(tconn, singleUse)
***REMOVED***

func (t *Transport) newTLSConfig(host string) *tls.Config ***REMOVED***
	cfg := new(tls.Config)
	if t.TLSClientConfig != nil ***REMOVED***
		*cfg = *t.TLSClientConfig.Clone()
	***REMOVED***
	if !strSliceContains(cfg.NextProtos, NextProtoTLS) ***REMOVED***
		cfg.NextProtos = append([]string***REMOVED***NextProtoTLS***REMOVED***, cfg.NextProtos...)
	***REMOVED***
	if cfg.ServerName == "" ***REMOVED***
		cfg.ServerName = host
	***REMOVED***
	return cfg
***REMOVED***

func (t *Transport) dialTLS(ctx context.Context) func(string, string, *tls.Config) (net.Conn, error) ***REMOVED***
	if t.DialTLS != nil ***REMOVED***
		return t.DialTLS
	***REMOVED***
	return func(network, addr string, cfg *tls.Config) (net.Conn, error) ***REMOVED***
		tlsCn, err := t.dialTLSWithContext(ctx, network, addr, cfg)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		state := tlsCn.ConnectionState()
		if p := state.NegotiatedProtocol; p != NextProtoTLS ***REMOVED***
			return nil, fmt.Errorf("http2: unexpected ALPN protocol %q; want %q", p, NextProtoTLS)
		***REMOVED***
		if !state.NegotiatedProtocolIsMutual ***REMOVED***
			return nil, errors.New("http2: could not negotiate protocol mutually")
		***REMOVED***
		return tlsCn, nil
	***REMOVED***
***REMOVED***

// disableKeepAlives reports whether connections should be closed as
// soon as possible after handling the first request.
func (t *Transport) disableKeepAlives() bool ***REMOVED***
	return t.t1 != nil && t.t1.DisableKeepAlives
***REMOVED***

func (t *Transport) expectContinueTimeout() time.Duration ***REMOVED***
	if t.t1 == nil ***REMOVED***
		return 0
	***REMOVED***
	return t.t1.ExpectContinueTimeout
***REMOVED***

func (t *Transport) NewClientConn(c net.Conn) (*ClientConn, error) ***REMOVED***
	return t.newClientConn(c, t.disableKeepAlives())
***REMOVED***

func (t *Transport) newClientConn(c net.Conn, singleUse bool) (*ClientConn, error) ***REMOVED***
	cc := &ClientConn***REMOVED***
		t:                     t,
		tconn:                 c,
		readerDone:            make(chan struct***REMOVED******REMOVED***),
		nextStreamID:          1,
		maxFrameSize:          16 << 10,                    // spec default
		initialWindowSize:     65535,                       // spec default
		maxConcurrentStreams:  initialMaxConcurrentStreams, // "infinite", per spec. Use a smaller value until we have received server settings.
		peerMaxHeaderListSize: 0xffffffffffffffff,          // "infinite", per spec. Use 2^64-1 instead.
		streams:               make(map[uint32]*clientStream),
		singleUse:             singleUse,
		wantSettingsAck:       true,
		pings:                 make(map[[8]byte]chan struct***REMOVED******REMOVED***),
		reqHeaderMu:           make(chan struct***REMOVED******REMOVED***, 1),
	***REMOVED***
	if d := t.idleConnTimeout(); d != 0 ***REMOVED***
		cc.idleTimeout = d
		cc.idleTimer = time.AfterFunc(d, cc.onIdleTimeout)
	***REMOVED***
	if VerboseLogs ***REMOVED***
		t.vlogf("http2: Transport creating client conn %p to %v", cc, c.RemoteAddr())
	***REMOVED***

	cc.cond = sync.NewCond(&cc.mu)
	cc.flow.add(int32(initialWindowSize))

	// TODO: adjust this writer size to account for frame size +
	// MTU + crypto/tls record padding.
	cc.bw = bufio.NewWriter(stickyErrWriter***REMOVED***
		conn:    c,
		timeout: t.WriteByteTimeout,
		err:     &cc.werr,
	***REMOVED***)
	cc.br = bufio.NewReader(c)
	cc.fr = NewFramer(cc.bw, cc.br)
	if t.CountError != nil ***REMOVED***
		cc.fr.countError = t.CountError
	***REMOVED***
	cc.fr.ReadMetaHeaders = hpack.NewDecoder(initialHeaderTableSize, nil)
	cc.fr.MaxHeaderListSize = t.maxHeaderListSize()

	// TODO: SetMaxDynamicTableSize, SetMaxDynamicTableSizeLimit on
	// henc in response to SETTINGS frames?
	cc.henc = hpack.NewEncoder(&cc.hbuf)

	if t.AllowHTTP ***REMOVED***
		cc.nextStreamID = 3
	***REMOVED***

	if cs, ok := c.(connectionStater); ok ***REMOVED***
		state := cs.ConnectionState()
		cc.tlsState = &state
	***REMOVED***

	initialSettings := []Setting***REMOVED***
		***REMOVED***ID: SettingEnablePush, Val: 0***REMOVED***,
		***REMOVED***ID: SettingInitialWindowSize, Val: transportDefaultStreamFlow***REMOVED***,
	***REMOVED***
	if max := t.maxHeaderListSize(); max != 0 ***REMOVED***
		initialSettings = append(initialSettings, Setting***REMOVED***ID: SettingMaxHeaderListSize, Val: max***REMOVED***)
	***REMOVED***

	cc.bw.Write(clientPreface)
	cc.fr.WriteSettings(initialSettings...)
	cc.fr.WriteWindowUpdate(0, transportDefaultConnFlow)
	cc.inflow.add(transportDefaultConnFlow + initialWindowSize)
	cc.bw.Flush()
	if cc.werr != nil ***REMOVED***
		cc.Close()
		return nil, cc.werr
	***REMOVED***

	go cc.readLoop()
	return cc, nil
***REMOVED***

func (cc *ClientConn) healthCheck() ***REMOVED***
	pingTimeout := cc.t.pingTimeout()
	// We don't need to periodically ping in the health check, because the readLoop of ClientConn will
	// trigger the healthCheck again if there is no frame received.
	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()
	err := cc.Ping(ctx)
	if err != nil ***REMOVED***
		cc.closeForLostPing()
		return
	***REMOVED***
***REMOVED***

// SetDoNotReuse marks cc as not reusable for future HTTP requests.
func (cc *ClientConn) SetDoNotReuse() ***REMOVED***
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.doNotReuse = true
***REMOVED***

func (cc *ClientConn) setGoAway(f *GoAwayFrame) ***REMOVED***
	cc.mu.Lock()
	defer cc.mu.Unlock()

	old := cc.goAway
	cc.goAway = f

	// Merge the previous and current GoAway error frames.
	if cc.goAwayDebug == "" ***REMOVED***
		cc.goAwayDebug = string(f.DebugData())
	***REMOVED***
	if old != nil && old.ErrCode != ErrCodeNo ***REMOVED***
		cc.goAway.ErrCode = old.ErrCode
	***REMOVED***
	last := f.LastStreamID
	for streamID, cs := range cc.streams ***REMOVED***
		if streamID > last ***REMOVED***
			cs.abortStreamLocked(errClientConnGotGoAway)
		***REMOVED***
	***REMOVED***
***REMOVED***

// CanTakeNewRequest reports whether the connection can take a new request,
// meaning it has not been closed or received or sent a GOAWAY.
//
// If the caller is going to immediately make a new request on this
// connection, use ReserveNewRequest instead.
func (cc *ClientConn) CanTakeNewRequest() bool ***REMOVED***
	cc.mu.Lock()
	defer cc.mu.Unlock()
	return cc.canTakeNewRequestLocked()
***REMOVED***

// ReserveNewRequest is like CanTakeNewRequest but also reserves a
// concurrent stream in cc. The reservation is decremented on the
// next call to RoundTrip.
func (cc *ClientConn) ReserveNewRequest() bool ***REMOVED***
	cc.mu.Lock()
	defer cc.mu.Unlock()
	if st := cc.idleStateLocked(); !st.canTakeNewRequest ***REMOVED***
		return false
	***REMOVED***
	cc.streamsReserved++
	return true
***REMOVED***

// ClientConnState describes the state of a ClientConn.
type ClientConnState struct ***REMOVED***
	// Closed is whether the connection is closed.
	Closed bool

	// Closing is whether the connection is in the process of
	// closing. It may be closing due to shutdown, being a
	// single-use connection, being marked as DoNotReuse, or
	// having received a GOAWAY frame.
	Closing bool

	// StreamsActive is how many streams are active.
	StreamsActive int

	// StreamsReserved is how many streams have been reserved via
	// ClientConn.ReserveNewRequest.
	StreamsReserved int

	// StreamsPending is how many requests have been sent in excess
	// of the peer's advertised MaxConcurrentStreams setting and
	// are waiting for other streams to complete.
	StreamsPending int

	// MaxConcurrentStreams is how many concurrent streams the
	// peer advertised as acceptable. Zero means no SETTINGS
	// frame has been received yet.
	MaxConcurrentStreams uint32

	// LastIdle, if non-zero, is when the connection last
	// transitioned to idle state.
	LastIdle time.Time
***REMOVED***

// State returns a snapshot of cc's state.
func (cc *ClientConn) State() ClientConnState ***REMOVED***
	cc.wmu.Lock()
	maxConcurrent := cc.maxConcurrentStreams
	if !cc.seenSettings ***REMOVED***
		maxConcurrent = 0
	***REMOVED***
	cc.wmu.Unlock()

	cc.mu.Lock()
	defer cc.mu.Unlock()
	return ClientConnState***REMOVED***
		Closed:               cc.closed,
		Closing:              cc.closing || cc.singleUse || cc.doNotReuse || cc.goAway != nil,
		StreamsActive:        len(cc.streams),
		StreamsReserved:      cc.streamsReserved,
		StreamsPending:       cc.pendingRequests,
		LastIdle:             cc.lastIdle,
		MaxConcurrentStreams: maxConcurrent,
	***REMOVED***
***REMOVED***

// clientConnIdleState describes the suitability of a client
// connection to initiate a new RoundTrip request.
type clientConnIdleState struct ***REMOVED***
	canTakeNewRequest bool
***REMOVED***

func (cc *ClientConn) idleState() clientConnIdleState ***REMOVED***
	cc.mu.Lock()
	defer cc.mu.Unlock()
	return cc.idleStateLocked()
***REMOVED***

func (cc *ClientConn) idleStateLocked() (st clientConnIdleState) ***REMOVED***
	if cc.singleUse && cc.nextStreamID > 1 ***REMOVED***
		return
	***REMOVED***
	var maxConcurrentOkay bool
	if cc.t.StrictMaxConcurrentStreams ***REMOVED***
		// We'll tell the caller we can take a new request to
		// prevent the caller from dialing a new TCP
		// connection, but then we'll block later before
		// writing it.
		maxConcurrentOkay = true
	***REMOVED*** else ***REMOVED***
		maxConcurrentOkay = int64(len(cc.streams)+cc.streamsReserved+1) <= int64(cc.maxConcurrentStreams)
	***REMOVED***

	st.canTakeNewRequest = cc.goAway == nil && !cc.closed && !cc.closing && maxConcurrentOkay &&
		!cc.doNotReuse &&
		int64(cc.nextStreamID)+2*int64(cc.pendingRequests) < math.MaxInt32 &&
		!cc.tooIdleLocked()
	return
***REMOVED***

func (cc *ClientConn) canTakeNewRequestLocked() bool ***REMOVED***
	st := cc.idleStateLocked()
	return st.canTakeNewRequest
***REMOVED***

// tooIdleLocked reports whether this connection has been been sitting idle
// for too much wall time.
func (cc *ClientConn) tooIdleLocked() bool ***REMOVED***
	// The Round(0) strips the monontonic clock reading so the
	// times are compared based on their wall time. We don't want
	// to reuse a connection that's been sitting idle during
	// VM/laptop suspend if monotonic time was also frozen.
	return cc.idleTimeout != 0 && !cc.lastIdle.IsZero() && time.Since(cc.lastIdle.Round(0)) > cc.idleTimeout
***REMOVED***

// onIdleTimeout is called from a time.AfterFunc goroutine. It will
// only be called when we're idle, but because we're coming from a new
// goroutine, there could be a new request coming in at the same time,
// so this simply calls the synchronized closeIfIdle to shut down this
// connection. The timer could just call closeIfIdle, but this is more
// clear.
func (cc *ClientConn) onIdleTimeout() ***REMOVED***
	cc.closeIfIdle()
***REMOVED***

func (cc *ClientConn) closeConn() error ***REMOVED***
	t := time.AfterFunc(250*time.Millisecond, cc.forceCloseConn)
	defer t.Stop()
	return cc.tconn.Close()
***REMOVED***

// A tls.Conn.Close can hang for a long time if the peer is unresponsive.
// Try to shut it down more aggressively.
func (cc *ClientConn) forceCloseConn() ***REMOVED***
	tc, ok := cc.tconn.(*tls.Conn)
	if !ok ***REMOVED***
		return
	***REMOVED***
	if nc := tlsUnderlyingConn(tc); nc != nil ***REMOVED***
		nc.Close()
	***REMOVED***
***REMOVED***

func (cc *ClientConn) closeIfIdle() ***REMOVED***
	cc.mu.Lock()
	if len(cc.streams) > 0 || cc.streamsReserved > 0 ***REMOVED***
		cc.mu.Unlock()
		return
	***REMOVED***
	cc.closed = true
	nextID := cc.nextStreamID
	// TODO: do clients send GOAWAY too? maybe? Just Close:
	cc.mu.Unlock()

	if VerboseLogs ***REMOVED***
		cc.vlogf("http2: Transport closing idle conn %p (forSingleUse=%v, maxStream=%v)", cc, cc.singleUse, nextID-2)
	***REMOVED***
	cc.closeConn()
***REMOVED***

func (cc *ClientConn) isDoNotReuseAndIdle() bool ***REMOVED***
	cc.mu.Lock()
	defer cc.mu.Unlock()
	return cc.doNotReuse && len(cc.streams) == 0
***REMOVED***

var shutdownEnterWaitStateHook = func() ***REMOVED******REMOVED***

// Shutdown gracefully closes the client connection, waiting for running streams to complete.
func (cc *ClientConn) Shutdown(ctx context.Context) error ***REMOVED***
	if err := cc.sendGoAway(); err != nil ***REMOVED***
		return err
	***REMOVED***
	// Wait for all in-flight streams to complete or connection to close
	done := make(chan struct***REMOVED******REMOVED***)
	cancelled := false // guarded by cc.mu
	go func() ***REMOVED***
		cc.mu.Lock()
		defer cc.mu.Unlock()
		for ***REMOVED***
			if len(cc.streams) == 0 || cc.closed ***REMOVED***
				cc.closed = true
				close(done)
				break
			***REMOVED***
			if cancelled ***REMOVED***
				break
			***REMOVED***
			cc.cond.Wait()
		***REMOVED***
	***REMOVED***()
	shutdownEnterWaitStateHook()
	select ***REMOVED***
	case <-done:
		return cc.closeConn()
	case <-ctx.Done():
		cc.mu.Lock()
		// Free the goroutine above
		cancelled = true
		cc.cond.Broadcast()
		cc.mu.Unlock()
		return ctx.Err()
	***REMOVED***
***REMOVED***

func (cc *ClientConn) sendGoAway() error ***REMOVED***
	cc.mu.Lock()
	closing := cc.closing
	cc.closing = true
	maxStreamID := cc.nextStreamID
	cc.mu.Unlock()
	if closing ***REMOVED***
		// GOAWAY sent already
		return nil
	***REMOVED***

	cc.wmu.Lock()
	defer cc.wmu.Unlock()
	// Send a graceful shutdown frame to server
	if err := cc.fr.WriteGoAway(maxStreamID, ErrCodeNo, nil); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := cc.bw.Flush(); err != nil ***REMOVED***
		return err
	***REMOVED***
	// Prevent new requests
	return nil
***REMOVED***

// closes the client connection immediately. In-flight requests are interrupted.
// err is sent to streams.
func (cc *ClientConn) closeForError(err error) error ***REMOVED***
	cc.mu.Lock()
	cc.closed = true
	for _, cs := range cc.streams ***REMOVED***
		cs.abortStreamLocked(err)
	***REMOVED***
	cc.cond.Broadcast()
	cc.mu.Unlock()
	return cc.closeConn()
***REMOVED***

// Close closes the client connection immediately.
//
// In-flight requests are interrupted. For a graceful shutdown, use Shutdown instead.
func (cc *ClientConn) Close() error ***REMOVED***
	err := errors.New("http2: client connection force closed via ClientConn.Close")
	return cc.closeForError(err)
***REMOVED***

// closes the client connection immediately. In-flight requests are interrupted.
func (cc *ClientConn) closeForLostPing() error ***REMOVED***
	err := errors.New("http2: client connection lost")
	if f := cc.t.CountError; f != nil ***REMOVED***
		f("conn_close_lost_ping")
	***REMOVED***
	return cc.closeForError(err)
***REMOVED***

// errRequestCanceled is a copy of net/http's errRequestCanceled because it's not
// exported. At least they'll be DeepEqual for h1-vs-h2 comparisons tests.
var errRequestCanceled = errors.New("net/http: request canceled")

func commaSeparatedTrailers(req *http.Request) (string, error) ***REMOVED***
	keys := make([]string, 0, len(req.Trailer))
	for k := range req.Trailer ***REMOVED***
		k = http.CanonicalHeaderKey(k)
		switch k ***REMOVED***
		case "Transfer-Encoding", "Trailer", "Content-Length":
			return "", fmt.Errorf("invalid Trailer key %q", k)
		***REMOVED***
		keys = append(keys, k)
	***REMOVED***
	if len(keys) > 0 ***REMOVED***
		sort.Strings(keys)
		return strings.Join(keys, ","), nil
	***REMOVED***
	return "", nil
***REMOVED***

func (cc *ClientConn) responseHeaderTimeout() time.Duration ***REMOVED***
	if cc.t.t1 != nil ***REMOVED***
		return cc.t.t1.ResponseHeaderTimeout
	***REMOVED***
	// No way to do this (yet?) with just an http2.Transport. Probably
	// no need. Request.Cancel this is the new way. We only need to support
	// this for compatibility with the old http.Transport fields when
	// we're doing transparent http2.
	return 0
***REMOVED***

// checkConnHeaders checks whether req has any invalid connection-level headers.
// per RFC 7540 section 8.1.2.2: Connection-Specific Header Fields.
// Certain headers are special-cased as okay but not transmitted later.
func checkConnHeaders(req *http.Request) error ***REMOVED***
	if v := req.Header.Get("Upgrade"); v != "" ***REMOVED***
		return fmt.Errorf("http2: invalid Upgrade request header: %q", req.Header["Upgrade"])
	***REMOVED***
	if vv := req.Header["Transfer-Encoding"]; len(vv) > 0 && (len(vv) > 1 || vv[0] != "" && vv[0] != "chunked") ***REMOVED***
		return fmt.Errorf("http2: invalid Transfer-Encoding request header: %q", vv)
	***REMOVED***
	if vv := req.Header["Connection"]; len(vv) > 0 && (len(vv) > 1 || vv[0] != "" && !asciiEqualFold(vv[0], "close") && !asciiEqualFold(vv[0], "keep-alive")) ***REMOVED***
		return fmt.Errorf("http2: invalid Connection request header: %q", vv)
	***REMOVED***
	return nil
***REMOVED***

// actualContentLength returns a sanitized version of
// req.ContentLength, where 0 actually means zero (not unknown) and -1
// means unknown.
func actualContentLength(req *http.Request) int64 ***REMOVED***
	if req.Body == nil || req.Body == http.NoBody ***REMOVED***
		return 0
	***REMOVED***
	if req.ContentLength != 0 ***REMOVED***
		return req.ContentLength
	***REMOVED***
	return -1
***REMOVED***

func (cc *ClientConn) decrStreamReservations() ***REMOVED***
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.decrStreamReservationsLocked()
***REMOVED***

func (cc *ClientConn) decrStreamReservationsLocked() ***REMOVED***
	if cc.streamsReserved > 0 ***REMOVED***
		cc.streamsReserved--
	***REMOVED***
***REMOVED***

func (cc *ClientConn) RoundTrip(req *http.Request) (*http.Response, error) ***REMOVED***
	ctx := req.Context()
	cs := &clientStream***REMOVED***
		cc:                   cc,
		ctx:                  ctx,
		reqCancel:            req.Cancel,
		isHead:               req.Method == "HEAD",
		reqBody:              req.Body,
		reqBodyContentLength: actualContentLength(req),
		trace:                httptrace.ContextClientTrace(ctx),
		peerClosed:           make(chan struct***REMOVED******REMOVED***),
		abort:                make(chan struct***REMOVED******REMOVED***),
		respHeaderRecv:       make(chan struct***REMOVED******REMOVED***),
		donec:                make(chan struct***REMOVED******REMOVED***),
	***REMOVED***
	go cs.doRequest(req)

	waitDone := func() error ***REMOVED***
		select ***REMOVED***
		case <-cs.donec:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		case <-cs.reqCancel:
			return errRequestCanceled
		***REMOVED***
	***REMOVED***

	handleResponseHeaders := func() (*http.Response, error) ***REMOVED***
		res := cs.res
		if res.StatusCode > 299 ***REMOVED***
			// On error or status code 3xx, 4xx, 5xx, etc abort any
			// ongoing write, assuming that the server doesn't care
			// about our request body. If the server replied with 1xx or
			// 2xx, however, then assume the server DOES potentially
			// want our body (e.g. full-duplex streaming:
			// golang.org/issue/13444). If it turns out the server
			// doesn't, they'll RST_STREAM us soon enough. This is a
			// heuristic to avoid adding knobs to Transport. Hopefully
			// we can keep it.
			cs.abortRequestBodyWrite()
		***REMOVED***
		res.Request = req
		res.TLS = cc.tlsState
		if res.Body == noBody && actualContentLength(req) == 0 ***REMOVED***
			// If there isn't a request or response body still being
			// written, then wait for the stream to be closed before
			// RoundTrip returns.
			if err := waitDone(); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED***
		return res, nil
	***REMOVED***

	for ***REMOVED***
		select ***REMOVED***
		case <-cs.respHeaderRecv:
			return handleResponseHeaders()
		case <-cs.abort:
			select ***REMOVED***
			case <-cs.respHeaderRecv:
				// If both cs.respHeaderRecv and cs.abort are signaling,
				// pick respHeaderRecv. The server probably wrote the
				// response and immediately reset the stream.
				// golang.org/issue/49645
				return handleResponseHeaders()
			default:
				waitDone()
				return nil, cs.abortErr
			***REMOVED***
		case <-ctx.Done():
			err := ctx.Err()
			cs.abortStream(err)
			return nil, err
		case <-cs.reqCancel:
			cs.abortStream(errRequestCanceled)
			return nil, errRequestCanceled
		***REMOVED***
	***REMOVED***
***REMOVED***

// doRequest runs for the duration of the request lifetime.
//
// It sends the request and performs post-request cleanup (closing Request.Body, etc.).
func (cs *clientStream) doRequest(req *http.Request) ***REMOVED***
	err := cs.writeRequest(req)
	cs.cleanupWriteRequest(err)
***REMOVED***

// writeRequest sends a request.
//
// It returns nil after the request is written, the response read,
// and the request stream is half-closed by the peer.
//
// It returns non-nil if the request ends otherwise.
// If the returned error is StreamError, the error Code may be used in resetting the stream.
func (cs *clientStream) writeRequest(req *http.Request) (err error) ***REMOVED***
	cc := cs.cc
	ctx := cs.ctx

	if err := checkConnHeaders(req); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Acquire the new-request lock by writing to reqHeaderMu.
	// This lock guards the critical section covering allocating a new stream ID
	// (requires mu) and creating the stream (requires wmu).
	if cc.reqHeaderMu == nil ***REMOVED***
		panic("RoundTrip on uninitialized ClientConn") // for tests
	***REMOVED***
	select ***REMOVED***
	case cc.reqHeaderMu <- struct***REMOVED******REMOVED******REMOVED******REMOVED***:
	case <-cs.reqCancel:
		return errRequestCanceled
	case <-ctx.Done():
		return ctx.Err()
	***REMOVED***

	cc.mu.Lock()
	if cc.idleTimer != nil ***REMOVED***
		cc.idleTimer.Stop()
	***REMOVED***
	cc.decrStreamReservationsLocked()
	if err := cc.awaitOpenSlotForStreamLocked(cs); err != nil ***REMOVED***
		cc.mu.Unlock()
		<-cc.reqHeaderMu
		return err
	***REMOVED***
	cc.addStreamLocked(cs) // assigns stream ID
	if isConnectionCloseRequest(req) ***REMOVED***
		cc.doNotReuse = true
	***REMOVED***
	cc.mu.Unlock()

	// TODO(bradfitz): this is a copy of the logic in net/http. Unify somewhere?
	if !cc.t.disableCompression() &&
		req.Header.Get("Accept-Encoding") == "" &&
		req.Header.Get("Range") == "" &&
		!cs.isHead ***REMOVED***
		// Request gzip only, not deflate. Deflate is ambiguous and
		// not as universally supported anyway.
		// See: https://zlib.net/zlib_faq.html#faq39
		//
		// Note that we don't request this for HEAD requests,
		// due to a bug in nginx:
		//   http://trac.nginx.org/nginx/ticket/358
		//   https://golang.org/issue/5522
		//
		// We don't request gzip if the request is for a range, since
		// auto-decoding a portion of a gzipped document will just fail
		// anyway. See https://golang.org/issue/8923
		cs.requestedGzip = true
	***REMOVED***

	continueTimeout := cc.t.expectContinueTimeout()
	if continueTimeout != 0 ***REMOVED***
		if !httpguts.HeaderValuesContainsToken(req.Header["Expect"], "100-continue") ***REMOVED***
			continueTimeout = 0
		***REMOVED*** else ***REMOVED***
			cs.on100 = make(chan struct***REMOVED******REMOVED***, 1)
		***REMOVED***
	***REMOVED***

	// Past this point (where we send request headers), it is possible for
	// RoundTrip to return successfully. Since the RoundTrip contract permits
	// the caller to "mutate or reuse" the Request after closing the Response's Body,
	// we must take care when referencing the Request from here on.
	err = cs.encodeAndWriteHeaders(req)
	<-cc.reqHeaderMu
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	hasBody := cs.reqBodyContentLength != 0
	if !hasBody ***REMOVED***
		cs.sentEndStream = true
	***REMOVED*** else ***REMOVED***
		if continueTimeout != 0 ***REMOVED***
			traceWait100Continue(cs.trace)
			timer := time.NewTimer(continueTimeout)
			select ***REMOVED***
			case <-timer.C:
				err = nil
			case <-cs.on100:
				err = nil
			case <-cs.abort:
				err = cs.abortErr
			case <-ctx.Done():
				err = ctx.Err()
			case <-cs.reqCancel:
				err = errRequestCanceled
			***REMOVED***
			timer.Stop()
			if err != nil ***REMOVED***
				traceWroteRequest(cs.trace, err)
				return err
			***REMOVED***
		***REMOVED***

		if err = cs.writeRequestBody(req); err != nil ***REMOVED***
			if err != errStopReqBodyWrite ***REMOVED***
				traceWroteRequest(cs.trace, err)
				return err
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			cs.sentEndStream = true
		***REMOVED***
	***REMOVED***

	traceWroteRequest(cs.trace, err)

	var respHeaderTimer <-chan time.Time
	var respHeaderRecv chan struct***REMOVED******REMOVED***
	if d := cc.responseHeaderTimeout(); d != 0 ***REMOVED***
		timer := time.NewTimer(d)
		defer timer.Stop()
		respHeaderTimer = timer.C
		respHeaderRecv = cs.respHeaderRecv
	***REMOVED***
	// Wait until the peer half-closes its end of the stream,
	// or until the request is aborted (via context, error, or otherwise),
	// whichever comes first.
	for ***REMOVED***
		select ***REMOVED***
		case <-cs.peerClosed:
			return nil
		case <-respHeaderTimer:
			return errTimeout
		case <-respHeaderRecv:
			respHeaderRecv = nil
			respHeaderTimer = nil // keep waiting for END_STREAM
		case <-cs.abort:
			return cs.abortErr
		case <-ctx.Done():
			return ctx.Err()
		case <-cs.reqCancel:
			return errRequestCanceled
		***REMOVED***
	***REMOVED***
***REMOVED***

func (cs *clientStream) encodeAndWriteHeaders(req *http.Request) error ***REMOVED***
	cc := cs.cc
	ctx := cs.ctx

	cc.wmu.Lock()
	defer cc.wmu.Unlock()

	// If the request was canceled while waiting for cc.mu, just quit.
	select ***REMOVED***
	case <-cs.abort:
		return cs.abortErr
	case <-ctx.Done():
		return ctx.Err()
	case <-cs.reqCancel:
		return errRequestCanceled
	default:
	***REMOVED***

	// Encode headers.
	//
	// we send: HEADERS***REMOVED***1***REMOVED***, CONTINUATION***REMOVED***0,***REMOVED*** + DATA***REMOVED***0,***REMOVED*** (DATA is
	// sent by writeRequestBody below, along with any Trailers,
	// again in form HEADERS***REMOVED***1***REMOVED***, CONTINUATION***REMOVED***0,***REMOVED***)
	trailers, err := commaSeparatedTrailers(req)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	hasTrailers := trailers != ""
	contentLen := actualContentLength(req)
	hasBody := contentLen != 0
	hdrs, err := cc.encodeHeaders(req, cs.requestedGzip, trailers, contentLen)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Write the request.
	endStream := !hasBody && !hasTrailers
	cs.sentHeaders = true
	err = cc.writeHeaders(cs.ID, endStream, int(cc.maxFrameSize), hdrs)
	traceWroteHeaders(cs.trace)
	return err
***REMOVED***

// cleanupWriteRequest performs post-request tasks.
//
// If err (the result of writeRequest) is non-nil and the stream is not closed,
// cleanupWriteRequest will send a reset to the peer.
func (cs *clientStream) cleanupWriteRequest(err error) ***REMOVED***
	cc := cs.cc

	if cs.ID == 0 ***REMOVED***
		// We were canceled before creating the stream, so return our reservation.
		cc.decrStreamReservations()
	***REMOVED***

	// TODO: write h12Compare test showing whether
	// Request.Body is closed by the Transport,
	// and in multiple cases: server replies <=299 and >299
	// while still writing request body
	cc.mu.Lock()
	bodyClosed := cs.reqBodyClosed
	cs.reqBodyClosed = true
	cc.mu.Unlock()
	if !bodyClosed && cs.reqBody != nil ***REMOVED***
		cs.reqBody.Close()
	***REMOVED***

	if err != nil && cs.sentEndStream ***REMOVED***
		// If the connection is closed immediately after the response is read,
		// we may be aborted before finishing up here. If the stream was closed
		// cleanly on both sides, there is no error.
		select ***REMOVED***
		case <-cs.peerClosed:
			err = nil
		default:
		***REMOVED***
	***REMOVED***
	if err != nil ***REMOVED***
		cs.abortStream(err) // possibly redundant, but harmless
		if cs.sentHeaders ***REMOVED***
			if se, ok := err.(StreamError); ok ***REMOVED***
				if se.Cause != errFromPeer ***REMOVED***
					cc.writeStreamReset(cs.ID, se.Code, err)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				cc.writeStreamReset(cs.ID, ErrCodeCancel, err)
			***REMOVED***
		***REMOVED***
		cs.bufPipe.CloseWithError(err) // no-op if already closed
	***REMOVED*** else ***REMOVED***
		if cs.sentHeaders && !cs.sentEndStream ***REMOVED***
			cc.writeStreamReset(cs.ID, ErrCodeNo, nil)
		***REMOVED***
		cs.bufPipe.CloseWithError(errRequestCanceled)
	***REMOVED***
	if cs.ID != 0 ***REMOVED***
		cc.forgetStreamID(cs.ID)
	***REMOVED***

	cc.wmu.Lock()
	werr := cc.werr
	cc.wmu.Unlock()
	if werr != nil ***REMOVED***
		cc.Close()
	***REMOVED***

	close(cs.donec)
***REMOVED***

// awaitOpenSlotForStream waits until len(streams) < maxConcurrentStreams.
// Must hold cc.mu.
func (cc *ClientConn) awaitOpenSlotForStreamLocked(cs *clientStream) error ***REMOVED***
	for ***REMOVED***
		cc.lastActive = time.Now()
		if cc.closed || !cc.canTakeNewRequestLocked() ***REMOVED***
			return errClientConnUnusable
		***REMOVED***
		cc.lastIdle = time.Time***REMOVED******REMOVED***
		if int64(len(cc.streams)) < int64(cc.maxConcurrentStreams) ***REMOVED***
			return nil
		***REMOVED***
		cc.pendingRequests++
		cc.cond.Wait()
		cc.pendingRequests--
		select ***REMOVED***
		case <-cs.abort:
			return cs.abortErr
		default:
		***REMOVED***
	***REMOVED***
***REMOVED***

// requires cc.wmu be held
func (cc *ClientConn) writeHeaders(streamID uint32, endStream bool, maxFrameSize int, hdrs []byte) error ***REMOVED***
	first := true // first frame written (HEADERS is first, then CONTINUATION)
	for len(hdrs) > 0 && cc.werr == nil ***REMOVED***
		chunk := hdrs
		if len(chunk) > maxFrameSize ***REMOVED***
			chunk = chunk[:maxFrameSize]
		***REMOVED***
		hdrs = hdrs[len(chunk):]
		endHeaders := len(hdrs) == 0
		if first ***REMOVED***
			cc.fr.WriteHeaders(HeadersFrameParam***REMOVED***
				StreamID:      streamID,
				BlockFragment: chunk,
				EndStream:     endStream,
				EndHeaders:    endHeaders,
			***REMOVED***)
			first = false
		***REMOVED*** else ***REMOVED***
			cc.fr.WriteContinuation(streamID, endHeaders, chunk)
		***REMOVED***
	***REMOVED***
	cc.bw.Flush()
	return cc.werr
***REMOVED***

// internal error values; they don't escape to callers
var (
	// abort request body write; don't send cancel
	errStopReqBodyWrite = errors.New("http2: aborting request body write")

	// abort request body write, but send stream reset of cancel.
	errStopReqBodyWriteAndCancel = errors.New("http2: canceling request")

	errReqBodyTooLong = errors.New("http2: request body larger than specified content length")
)

// frameScratchBufferLen returns the length of a buffer to use for
// outgoing request bodies to read/write to/from.
//
// It returns max(1, min(peer's advertised max frame size,
// Request.ContentLength+1, 512KB)).
func (cs *clientStream) frameScratchBufferLen(maxFrameSize int) int ***REMOVED***
	const max = 512 << 10
	n := int64(maxFrameSize)
	if n > max ***REMOVED***
		n = max
	***REMOVED***
	if cl := cs.reqBodyContentLength; cl != -1 && cl+1 < n ***REMOVED***
		// Add an extra byte past the declared content-length to
		// give the caller's Request.Body io.Reader a chance to
		// give us more bytes than they declared, so we can catch it
		// early.
		n = cl + 1
	***REMOVED***
	if n < 1 ***REMOVED***
		return 1
	***REMOVED***
	return int(n) // doesn't truncate; max is 512K
***REMOVED***

var bufPool sync.Pool // of *[]byte

func (cs *clientStream) writeRequestBody(req *http.Request) (err error) ***REMOVED***
	cc := cs.cc
	body := cs.reqBody
	sentEnd := false // whether we sent the final DATA frame w/ END_STREAM

	hasTrailers := req.Trailer != nil
	remainLen := cs.reqBodyContentLength
	hasContentLen := remainLen != -1

	cc.mu.Lock()
	maxFrameSize := int(cc.maxFrameSize)
	cc.mu.Unlock()

	// Scratch buffer for reading into & writing from.
	scratchLen := cs.frameScratchBufferLen(maxFrameSize)
	var buf []byte
	if bp, ok := bufPool.Get().(*[]byte); ok && len(*bp) >= scratchLen ***REMOVED***
		defer bufPool.Put(bp)
		buf = *bp
	***REMOVED*** else ***REMOVED***
		buf = make([]byte, scratchLen)
		defer bufPool.Put(&buf)
	***REMOVED***

	var sawEOF bool
	for !sawEOF ***REMOVED***
		n, err := body.Read(buf[:len(buf)])
		if hasContentLen ***REMOVED***
			remainLen -= int64(n)
			if remainLen == 0 && err == nil ***REMOVED***
				// The request body's Content-Length was predeclared and
				// we just finished reading it all, but the underlying io.Reader
				// returned the final chunk with a nil error (which is one of
				// the two valid things a Reader can do at EOF). Because we'd prefer
				// to send the END_STREAM bit early, double-check that we're actually
				// at EOF. Subsequent reads should return (0, EOF) at this point.
				// If either value is different, we return an error in one of two ways below.
				var scratch [1]byte
				var n1 int
				n1, err = body.Read(scratch[:])
				remainLen -= int64(n1)
			***REMOVED***
			if remainLen < 0 ***REMOVED***
				err = errReqBodyTooLong
				return err
			***REMOVED***
		***REMOVED***
		if err != nil ***REMOVED***
			cc.mu.Lock()
			bodyClosed := cs.reqBodyClosed
			cc.mu.Unlock()
			switch ***REMOVED***
			case bodyClosed:
				return errStopReqBodyWrite
			case err == io.EOF:
				sawEOF = true
				err = nil
			default:
				return err
			***REMOVED***
		***REMOVED***

		remain := buf[:n]
		for len(remain) > 0 && err == nil ***REMOVED***
			var allowed int32
			allowed, err = cs.awaitFlowControl(len(remain))
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			cc.wmu.Lock()
			data := remain[:allowed]
			remain = remain[allowed:]
			sentEnd = sawEOF && len(remain) == 0 && !hasTrailers
			err = cc.fr.WriteData(cs.ID, sentEnd, data)
			if err == nil ***REMOVED***
				// TODO(bradfitz): this flush is for latency, not bandwidth.
				// Most requests won't need this. Make this opt-in or
				// opt-out?  Use some heuristic on the body type? Nagel-like
				// timers?  Based on 'n'? Only last chunk of this for loop,
				// unless flow control tokens are low? For now, always.
				// If we change this, see comment below.
				err = cc.bw.Flush()
			***REMOVED***
			cc.wmu.Unlock()
		***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if sentEnd ***REMOVED***
		// Already sent END_STREAM (which implies we have no
		// trailers) and flushed, because currently all
		// WriteData frames above get a flush. So we're done.
		return nil
	***REMOVED***

	// Since the RoundTrip contract permits the caller to "mutate or reuse"
	// a request after the Response's Body is closed, verify that this hasn't
	// happened before accessing the trailers.
	cc.mu.Lock()
	trailer := req.Trailer
	err = cs.abortErr
	cc.mu.Unlock()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	cc.wmu.Lock()
	defer cc.wmu.Unlock()
	var trls []byte
	if len(trailer) > 0 ***REMOVED***
		trls, err = cc.encodeTrailers(trailer)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// Two ways to send END_STREAM: either with trailers, or
	// with an empty DATA frame.
	if len(trls) > 0 ***REMOVED***
		err = cc.writeHeaders(cs.ID, true, maxFrameSize, trls)
	***REMOVED*** else ***REMOVED***
		err = cc.fr.WriteData(cs.ID, true, nil)
	***REMOVED***
	if ferr := cc.bw.Flush(); ferr != nil && err == nil ***REMOVED***
		err = ferr
	***REMOVED***
	return err
***REMOVED***

// awaitFlowControl waits for [1, min(maxBytes, cc.cs.maxFrameSize)] flow
// control tokens from the server.
// It returns either the non-zero number of tokens taken or an error
// if the stream is dead.
func (cs *clientStream) awaitFlowControl(maxBytes int) (taken int32, err error) ***REMOVED***
	cc := cs.cc
	ctx := cs.ctx
	cc.mu.Lock()
	defer cc.mu.Unlock()
	for ***REMOVED***
		if cc.closed ***REMOVED***
			return 0, errClientConnClosed
		***REMOVED***
		if cs.reqBodyClosed ***REMOVED***
			return 0, errStopReqBodyWrite
		***REMOVED***
		select ***REMOVED***
		case <-cs.abort:
			return 0, cs.abortErr
		case <-ctx.Done():
			return 0, ctx.Err()
		case <-cs.reqCancel:
			return 0, errRequestCanceled
		default:
		***REMOVED***
		if a := cs.flow.available(); a > 0 ***REMOVED***
			take := a
			if int(take) > maxBytes ***REMOVED***

				take = int32(maxBytes) // can't truncate int; take is int32
			***REMOVED***
			if take > int32(cc.maxFrameSize) ***REMOVED***
				take = int32(cc.maxFrameSize)
			***REMOVED***
			cs.flow.take(take)
			return take, nil
		***REMOVED***
		cc.cond.Wait()
	***REMOVED***
***REMOVED***

var errNilRequestURL = errors.New("http2: Request.URI is nil")

// requires cc.wmu be held.
func (cc *ClientConn) encodeHeaders(req *http.Request, addGzipHeader bool, trailers string, contentLength int64) ([]byte, error) ***REMOVED***
	cc.hbuf.Reset()
	if req.URL == nil ***REMOVED***
		return nil, errNilRequestURL
	***REMOVED***

	host := req.Host
	if host == "" ***REMOVED***
		host = req.URL.Host
	***REMOVED***
	host, err := httpguts.PunycodeHostPort(host)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var path string
	if req.Method != "CONNECT" ***REMOVED***
		path = req.URL.RequestURI()
		if !validPseudoPath(path) ***REMOVED***
			orig := path
			path = strings.TrimPrefix(path, req.URL.Scheme+"://"+host)
			if !validPseudoPath(path) ***REMOVED***
				if req.URL.Opaque != "" ***REMOVED***
					return nil, fmt.Errorf("invalid request :path %q from URL.Opaque = %q", orig, req.URL.Opaque)
				***REMOVED*** else ***REMOVED***
					return nil, fmt.Errorf("invalid request :path %q", orig)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Check for any invalid headers and return an error before we
	// potentially pollute our hpack state. (We want to be able to
	// continue to reuse the hpack encoder for future requests)
	for k, vv := range req.Header ***REMOVED***
		if !httpguts.ValidHeaderFieldName(k) ***REMOVED***
			return nil, fmt.Errorf("invalid HTTP header name %q", k)
		***REMOVED***
		for _, v := range vv ***REMOVED***
			if !httpguts.ValidHeaderFieldValue(v) ***REMOVED***
				return nil, fmt.Errorf("invalid HTTP header value %q for header %q", v, k)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	enumerateHeaders := func(f func(name, value string)) ***REMOVED***
		// 8.1.2.3 Request Pseudo-Header Fields
		// The :path pseudo-header field includes the path and query parts of the
		// target URI (the path-absolute production and optionally a '?' character
		// followed by the query production (see Sections 3.3 and 3.4 of
		// [RFC3986]).
		f(":authority", host)
		m := req.Method
		if m == "" ***REMOVED***
			m = http.MethodGet
		***REMOVED***
		f(":method", m)
		if req.Method != "CONNECT" ***REMOVED***
			f(":path", path)
			f(":scheme", req.URL.Scheme)
		***REMOVED***
		if trailers != "" ***REMOVED***
			f("trailer", trailers)
		***REMOVED***

		var didUA bool
		for k, vv := range req.Header ***REMOVED***
			if asciiEqualFold(k, "host") || asciiEqualFold(k, "content-length") ***REMOVED***
				// Host is :authority, already sent.
				// Content-Length is automatic, set below.
				continue
			***REMOVED*** else if asciiEqualFold(k, "connection") ||
				asciiEqualFold(k, "proxy-connection") ||
				asciiEqualFold(k, "transfer-encoding") ||
				asciiEqualFold(k, "upgrade") ||
				asciiEqualFold(k, "keep-alive") ***REMOVED***
				// Per 8.1.2.2 Connection-Specific Header
				// Fields, don't send connection-specific
				// fields. We have already checked if any
				// are error-worthy so just ignore the rest.
				continue
			***REMOVED*** else if asciiEqualFold(k, "user-agent") ***REMOVED***
				// Match Go's http1 behavior: at most one
				// User-Agent. If set to nil or empty string,
				// then omit it. Otherwise if not mentioned,
				// include the default (below).
				didUA = true
				if len(vv) < 1 ***REMOVED***
					continue
				***REMOVED***
				vv = vv[:1]
				if vv[0] == "" ***REMOVED***
					continue
				***REMOVED***
			***REMOVED*** else if asciiEqualFold(k, "cookie") ***REMOVED***
				// Per 8.1.2.5 To allow for better compression efficiency, the
				// Cookie header field MAY be split into separate header fields,
				// each with one or more cookie-pairs.
				for _, v := range vv ***REMOVED***
					for ***REMOVED***
						p := strings.IndexByte(v, ';')
						if p < 0 ***REMOVED***
							break
						***REMOVED***
						f("cookie", v[:p])
						p++
						// strip space after semicolon if any.
						for p+1 <= len(v) && v[p] == ' ' ***REMOVED***
							p++
						***REMOVED***
						v = v[p:]
					***REMOVED***
					if len(v) > 0 ***REMOVED***
						f("cookie", v)
					***REMOVED***
				***REMOVED***
				continue
			***REMOVED***

			for _, v := range vv ***REMOVED***
				f(k, v)
			***REMOVED***
		***REMOVED***
		if shouldSendReqContentLength(req.Method, contentLength) ***REMOVED***
			f("content-length", strconv.FormatInt(contentLength, 10))
		***REMOVED***
		if addGzipHeader ***REMOVED***
			f("accept-encoding", "gzip")
		***REMOVED***
		if !didUA ***REMOVED***
			f("user-agent", defaultUserAgent)
		***REMOVED***
	***REMOVED***

	// Do a first pass over the headers counting bytes to ensure
	// we don't exceed cc.peerMaxHeaderListSize. This is done as a
	// separate pass before encoding the headers to prevent
	// modifying the hpack state.
	hlSize := uint64(0)
	enumerateHeaders(func(name, value string) ***REMOVED***
		hf := hpack.HeaderField***REMOVED***Name: name, Value: value***REMOVED***
		hlSize += uint64(hf.Size())
	***REMOVED***)

	if hlSize > cc.peerMaxHeaderListSize ***REMOVED***
		return nil, errRequestHeaderListSize
	***REMOVED***

	trace := httptrace.ContextClientTrace(req.Context())
	traceHeaders := traceHasWroteHeaderField(trace)

	// Header list size is ok. Write the headers.
	enumerateHeaders(func(name, value string) ***REMOVED***
		name, ascii := asciiToLower(name)
		if !ascii ***REMOVED***
			// Skip writing invalid headers. Per RFC 7540, Section 8.1.2, header
			// field names have to be ASCII characters (just as in HTTP/1.x).
			return
		***REMOVED***
		cc.writeHeader(name, value)
		if traceHeaders ***REMOVED***
			traceWroteHeaderField(trace, name, value)
		***REMOVED***
	***REMOVED***)

	return cc.hbuf.Bytes(), nil
***REMOVED***

// shouldSendReqContentLength reports whether the http2.Transport should send
// a "content-length" request header. This logic is basically a copy of the net/http
// transferWriter.shouldSendContentLength.
// The contentLength is the corrected contentLength (so 0 means actually 0, not unknown).
// -1 means unknown.
func shouldSendReqContentLength(method string, contentLength int64) bool ***REMOVED***
	if contentLength > 0 ***REMOVED***
		return true
	***REMOVED***
	if contentLength < 0 ***REMOVED***
		return false
	***REMOVED***
	// For zero bodies, whether we send a content-length depends on the method.
	// It also kinda doesn't matter for http2 either way, with END_STREAM.
	switch method ***REMOVED***
	case "POST", "PUT", "PATCH":
		return true
	default:
		return false
	***REMOVED***
***REMOVED***

// requires cc.wmu be held.
func (cc *ClientConn) encodeTrailers(trailer http.Header) ([]byte, error) ***REMOVED***
	cc.hbuf.Reset()

	hlSize := uint64(0)
	for k, vv := range trailer ***REMOVED***
		for _, v := range vv ***REMOVED***
			hf := hpack.HeaderField***REMOVED***Name: k, Value: v***REMOVED***
			hlSize += uint64(hf.Size())
		***REMOVED***
	***REMOVED***
	if hlSize > cc.peerMaxHeaderListSize ***REMOVED***
		return nil, errRequestHeaderListSize
	***REMOVED***

	for k, vv := range trailer ***REMOVED***
		lowKey, ascii := asciiToLower(k)
		if !ascii ***REMOVED***
			// Skip writing invalid headers. Per RFC 7540, Section 8.1.2, header
			// field names have to be ASCII characters (just as in HTTP/1.x).
			continue
		***REMOVED***
		// Transfer-Encoding, etc.. have already been filtered at the
		// start of RoundTrip
		for _, v := range vv ***REMOVED***
			cc.writeHeader(lowKey, v)
		***REMOVED***
	***REMOVED***
	return cc.hbuf.Bytes(), nil
***REMOVED***

func (cc *ClientConn) writeHeader(name, value string) ***REMOVED***
	if VerboseLogs ***REMOVED***
		log.Printf("http2: Transport encoding header %q = %q", name, value)
	***REMOVED***
	cc.henc.WriteField(hpack.HeaderField***REMOVED***Name: name, Value: value***REMOVED***)
***REMOVED***

type resAndError struct ***REMOVED***
	_   incomparable
	res *http.Response
	err error
***REMOVED***

// requires cc.mu be held.
func (cc *ClientConn) addStreamLocked(cs *clientStream) ***REMOVED***
	cs.flow.add(int32(cc.initialWindowSize))
	cs.flow.setConnFlow(&cc.flow)
	cs.inflow.add(transportDefaultStreamFlow)
	cs.inflow.setConnFlow(&cc.inflow)
	cs.ID = cc.nextStreamID
	cc.nextStreamID += 2
	cc.streams[cs.ID] = cs
	if cs.ID == 0 ***REMOVED***
		panic("assigned stream ID 0")
	***REMOVED***
***REMOVED***

func (cc *ClientConn) forgetStreamID(id uint32) ***REMOVED***
	cc.mu.Lock()
	slen := len(cc.streams)
	delete(cc.streams, id)
	if len(cc.streams) != slen-1 ***REMOVED***
		panic("forgetting unknown stream id")
	***REMOVED***
	cc.lastActive = time.Now()
	if len(cc.streams) == 0 && cc.idleTimer != nil ***REMOVED***
		cc.idleTimer.Reset(cc.idleTimeout)
		cc.lastIdle = time.Now()
	***REMOVED***
	// Wake up writeRequestBody via clientStream.awaitFlowControl and
	// wake up RoundTrip if there is a pending request.
	cc.cond.Broadcast()

	closeOnIdle := cc.singleUse || cc.doNotReuse || cc.t.disableKeepAlives()
	if closeOnIdle && cc.streamsReserved == 0 && len(cc.streams) == 0 ***REMOVED***
		if VerboseLogs ***REMOVED***
			cc.vlogf("http2: Transport closing idle conn %p (forSingleUse=%v, maxStream=%v)", cc, cc.singleUse, cc.nextStreamID-2)
		***REMOVED***
		cc.closed = true
		defer cc.closeConn()
	***REMOVED***

	cc.mu.Unlock()
***REMOVED***

// clientConnReadLoop is the state owned by the clientConn's frame-reading readLoop.
type clientConnReadLoop struct ***REMOVED***
	_  incomparable
	cc *ClientConn
***REMOVED***

// readLoop runs in its own goroutine and reads and dispatches frames.
func (cc *ClientConn) readLoop() ***REMOVED***
	rl := &clientConnReadLoop***REMOVED***cc: cc***REMOVED***
	defer rl.cleanup()
	cc.readerErr = rl.run()
	if ce, ok := cc.readerErr.(ConnectionError); ok ***REMOVED***
		cc.wmu.Lock()
		cc.fr.WriteGoAway(0, ErrCode(ce), nil)
		cc.wmu.Unlock()
	***REMOVED***
***REMOVED***

// GoAwayError is returned by the Transport when the server closes the
// TCP connection after sending a GOAWAY frame.
type GoAwayError struct ***REMOVED***
	LastStreamID uint32
	ErrCode      ErrCode
	DebugData    string
***REMOVED***

func (e GoAwayError) Error() string ***REMOVED***
	return fmt.Sprintf("http2: server sent GOAWAY and closed the connection; LastStreamID=%v, ErrCode=%v, debug=%q",
		e.LastStreamID, e.ErrCode, e.DebugData)
***REMOVED***

func isEOFOrNetReadError(err error) bool ***REMOVED***
	if err == io.EOF ***REMOVED***
		return true
	***REMOVED***
	ne, ok := err.(*net.OpError)
	return ok && ne.Op == "read"
***REMOVED***

func (rl *clientConnReadLoop) cleanup() ***REMOVED***
	cc := rl.cc
	cc.t.connPool().MarkDead(cc)
	defer cc.closeConn()
	defer close(cc.readerDone)

	if cc.idleTimer != nil ***REMOVED***
		cc.idleTimer.Stop()
	***REMOVED***

	// Close any response bodies if the server closes prematurely.
	// TODO: also do this if we've written the headers but not
	// gotten a response yet.
	err := cc.readerErr
	cc.mu.Lock()
	if cc.goAway != nil && isEOFOrNetReadError(err) ***REMOVED***
		err = GoAwayError***REMOVED***
			LastStreamID: cc.goAway.LastStreamID,
			ErrCode:      cc.goAway.ErrCode,
			DebugData:    cc.goAwayDebug,
		***REMOVED***
	***REMOVED*** else if err == io.EOF ***REMOVED***
		err = io.ErrUnexpectedEOF
	***REMOVED***
	cc.closed = true
	for _, cs := range cc.streams ***REMOVED***
		select ***REMOVED***
		case <-cs.peerClosed:
			// The server closed the stream before closing the conn,
			// so no need to interrupt it.
		default:
			cs.abortStreamLocked(err)
		***REMOVED***
	***REMOVED***
	cc.cond.Broadcast()
	cc.mu.Unlock()
***REMOVED***

// countReadFrameError calls Transport.CountError with a string
// representing err.
func (cc *ClientConn) countReadFrameError(err error) ***REMOVED***
	f := cc.t.CountError
	if f == nil || err == nil ***REMOVED***
		return
	***REMOVED***
	if ce, ok := err.(ConnectionError); ok ***REMOVED***
		errCode := ErrCode(ce)
		f(fmt.Sprintf("read_frame_conn_error_%s", errCode.stringToken()))
		return
	***REMOVED***
	if errors.Is(err, io.EOF) ***REMOVED***
		f("read_frame_eof")
		return
	***REMOVED***
	if errors.Is(err, io.ErrUnexpectedEOF) ***REMOVED***
		f("read_frame_unexpected_eof")
		return
	***REMOVED***
	if errors.Is(err, ErrFrameTooLarge) ***REMOVED***
		f("read_frame_too_large")
		return
	***REMOVED***
	f("read_frame_other")
***REMOVED***

func (rl *clientConnReadLoop) run() error ***REMOVED***
	cc := rl.cc
	gotSettings := false
	readIdleTimeout := cc.t.ReadIdleTimeout
	var t *time.Timer
	if readIdleTimeout != 0 ***REMOVED***
		t = time.AfterFunc(readIdleTimeout, cc.healthCheck)
		defer t.Stop()
	***REMOVED***
	for ***REMOVED***
		f, err := cc.fr.ReadFrame()
		if t != nil ***REMOVED***
			t.Reset(readIdleTimeout)
		***REMOVED***
		if err != nil ***REMOVED***
			cc.vlogf("http2: Transport readFrame error on conn %p: (%T) %v", cc, err, err)
		***REMOVED***
		if se, ok := err.(StreamError); ok ***REMOVED***
			if cs := rl.streamByID(se.StreamID); cs != nil ***REMOVED***
				if se.Cause == nil ***REMOVED***
					se.Cause = cc.fr.errDetail
				***REMOVED***
				rl.endStreamError(cs, se)
			***REMOVED***
			continue
		***REMOVED*** else if err != nil ***REMOVED***
			cc.countReadFrameError(err)
			return err
		***REMOVED***
		if VerboseLogs ***REMOVED***
			cc.vlogf("http2: Transport received %s", summarizeFrame(f))
		***REMOVED***
		if !gotSettings ***REMOVED***
			if _, ok := f.(*SettingsFrame); !ok ***REMOVED***
				cc.logf("protocol error: received %T before a SETTINGS frame", f)
				return ConnectionError(ErrCodeProtocol)
			***REMOVED***
			gotSettings = true
		***REMOVED***

		switch f := f.(type) ***REMOVED***
		case *MetaHeadersFrame:
			err = rl.processHeaders(f)
		case *DataFrame:
			err = rl.processData(f)
		case *GoAwayFrame:
			err = rl.processGoAway(f)
		case *RSTStreamFrame:
			err = rl.processResetStream(f)
		case *SettingsFrame:
			err = rl.processSettings(f)
		case *PushPromiseFrame:
			err = rl.processPushPromise(f)
		case *WindowUpdateFrame:
			err = rl.processWindowUpdate(f)
		case *PingFrame:
			err = rl.processPing(f)
		default:
			cc.logf("Transport: unhandled response frame type %T", f)
		***REMOVED***
		if err != nil ***REMOVED***
			if VerboseLogs ***REMOVED***
				cc.vlogf("http2: Transport conn %p received error from processing frame %v: %v", cc, summarizeFrame(f), err)
			***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
***REMOVED***

func (rl *clientConnReadLoop) processHeaders(f *MetaHeadersFrame) error ***REMOVED***
	cs := rl.streamByID(f.StreamID)
	if cs == nil ***REMOVED***
		// We'd get here if we canceled a request while the
		// server had its response still in flight. So if this
		// was just something we canceled, ignore it.
		return nil
	***REMOVED***
	if cs.readClosed ***REMOVED***
		rl.endStreamError(cs, StreamError***REMOVED***
			StreamID: f.StreamID,
			Code:     ErrCodeProtocol,
			Cause:    errors.New("protocol error: headers after END_STREAM"),
		***REMOVED***)
		return nil
	***REMOVED***
	if !cs.firstByte ***REMOVED***
		if cs.trace != nil ***REMOVED***
			// TODO(bradfitz): move first response byte earlier,
			// when we first read the 9 byte header, not waiting
			// until all the HEADERS+CONTINUATION frames have been
			// merged. This works for now.
			traceFirstResponseByte(cs.trace)
		***REMOVED***
		cs.firstByte = true
	***REMOVED***
	if !cs.pastHeaders ***REMOVED***
		cs.pastHeaders = true
	***REMOVED*** else ***REMOVED***
		return rl.processTrailers(cs, f)
	***REMOVED***

	res, err := rl.handleResponse(cs, f)
	if err != nil ***REMOVED***
		if _, ok := err.(ConnectionError); ok ***REMOVED***
			return err
		***REMOVED***
		// Any other error type is a stream error.
		rl.endStreamError(cs, StreamError***REMOVED***
			StreamID: f.StreamID,
			Code:     ErrCodeProtocol,
			Cause:    err,
		***REMOVED***)
		return nil // return nil from process* funcs to keep conn alive
	***REMOVED***
	if res == nil ***REMOVED***
		// (nil, nil) special case. See handleResponse docs.
		return nil
	***REMOVED***
	cs.resTrailer = &res.Trailer
	cs.res = res
	close(cs.respHeaderRecv)
	if f.StreamEnded() ***REMOVED***
		rl.endStream(cs)
	***REMOVED***
	return nil
***REMOVED***

// may return error types nil, or ConnectionError. Any other error value
// is a StreamError of type ErrCodeProtocol. The returned error in that case
// is the detail.
//
// As a special case, handleResponse may return (nil, nil) to skip the
// frame (currently only used for 1xx responses).
func (rl *clientConnReadLoop) handleResponse(cs *clientStream, f *MetaHeadersFrame) (*http.Response, error) ***REMOVED***
	if f.Truncated ***REMOVED***
		return nil, errResponseHeaderListSize
	***REMOVED***

	status := f.PseudoValue("status")
	if status == "" ***REMOVED***
		return nil, errors.New("malformed response from server: missing status pseudo header")
	***REMOVED***
	statusCode, err := strconv.Atoi(status)
	if err != nil ***REMOVED***
		return nil, errors.New("malformed response from server: malformed non-numeric status pseudo header")
	***REMOVED***

	regularFields := f.RegularFields()
	strs := make([]string, len(regularFields))
	header := make(http.Header, len(regularFields))
	res := &http.Response***REMOVED***
		Proto:      "HTTP/2.0",
		ProtoMajor: 2,
		Header:     header,
		StatusCode: statusCode,
		Status:     status + " " + http.StatusText(statusCode),
	***REMOVED***
	for _, hf := range regularFields ***REMOVED***
		key := http.CanonicalHeaderKey(hf.Name)
		if key == "Trailer" ***REMOVED***
			t := res.Trailer
			if t == nil ***REMOVED***
				t = make(http.Header)
				res.Trailer = t
			***REMOVED***
			foreachHeaderElement(hf.Value, func(v string) ***REMOVED***
				t[http.CanonicalHeaderKey(v)] = nil
			***REMOVED***)
		***REMOVED*** else ***REMOVED***
			vv := header[key]
			if vv == nil && len(strs) > 0 ***REMOVED***
				// More than likely this will be a single-element key.
				// Most headers aren't multi-valued.
				// Set the capacity on strs[0] to 1, so any future append
				// won't extend the slice into the other strings.
				vv, strs = strs[:1:1], strs[1:]
				vv[0] = hf.Value
				header[key] = vv
			***REMOVED*** else ***REMOVED***
				header[key] = append(vv, hf.Value)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if statusCode >= 100 && statusCode <= 199 ***REMOVED***
		if f.StreamEnded() ***REMOVED***
			return nil, errors.New("1xx informational response with END_STREAM flag")
		***REMOVED***
		cs.num1xx++
		const max1xxResponses = 5 // arbitrary bound on number of informational responses, same as net/http
		if cs.num1xx > max1xxResponses ***REMOVED***
			return nil, errors.New("http2: too many 1xx informational responses")
		***REMOVED***
		if fn := cs.get1xxTraceFunc(); fn != nil ***REMOVED***
			if err := fn(statusCode, textproto.MIMEHeader(header)); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED***
		if statusCode == 100 ***REMOVED***
			traceGot100Continue(cs.trace)
			select ***REMOVED***
			case cs.on100 <- struct***REMOVED******REMOVED******REMOVED******REMOVED***:
			default:
			***REMOVED***
		***REMOVED***
		cs.pastHeaders = false // do it all again
		return nil, nil
	***REMOVED***

	res.ContentLength = -1
	if clens := res.Header["Content-Length"]; len(clens) == 1 ***REMOVED***
		if cl, err := strconv.ParseUint(clens[0], 10, 63); err == nil ***REMOVED***
			res.ContentLength = int64(cl)
		***REMOVED*** else ***REMOVED***
			// TODO: care? unlike http/1, it won't mess up our framing, so it's
			// more safe smuggling-wise to ignore.
		***REMOVED***
	***REMOVED*** else if len(clens) > 1 ***REMOVED***
		// TODO: care? unlike http/1, it won't mess up our framing, so it's
		// more safe smuggling-wise to ignore.
	***REMOVED*** else if f.StreamEnded() && !cs.isHead ***REMOVED***
		res.ContentLength = 0
	***REMOVED***

	if cs.isHead ***REMOVED***
		res.Body = noBody
		return res, nil
	***REMOVED***

	if f.StreamEnded() ***REMOVED***
		if res.ContentLength > 0 ***REMOVED***
			res.Body = missingBody***REMOVED******REMOVED***
		***REMOVED*** else ***REMOVED***
			res.Body = noBody
		***REMOVED***
		return res, nil
	***REMOVED***

	cs.bufPipe.setBuffer(&dataBuffer***REMOVED***expected: res.ContentLength***REMOVED***)
	cs.bytesRemain = res.ContentLength
	res.Body = transportResponseBody***REMOVED***cs***REMOVED***

	if cs.requestedGzip && asciiEqualFold(res.Header.Get("Content-Encoding"), "gzip") ***REMOVED***
		res.Header.Del("Content-Encoding")
		res.Header.Del("Content-Length")
		res.ContentLength = -1
		res.Body = &gzipReader***REMOVED***body: res.Body***REMOVED***
		res.Uncompressed = true
	***REMOVED***
	return res, nil
***REMOVED***

func (rl *clientConnReadLoop) processTrailers(cs *clientStream, f *MetaHeadersFrame) error ***REMOVED***
	if cs.pastTrailers ***REMOVED***
		// Too many HEADERS frames for this stream.
		return ConnectionError(ErrCodeProtocol)
	***REMOVED***
	cs.pastTrailers = true
	if !f.StreamEnded() ***REMOVED***
		// We expect that any headers for trailers also
		// has END_STREAM.
		return ConnectionError(ErrCodeProtocol)
	***REMOVED***
	if len(f.PseudoFields()) > 0 ***REMOVED***
		// No pseudo header fields are defined for trailers.
		// TODO: ConnectionError might be overly harsh? Check.
		return ConnectionError(ErrCodeProtocol)
	***REMOVED***

	trailer := make(http.Header)
	for _, hf := range f.RegularFields() ***REMOVED***
		key := http.CanonicalHeaderKey(hf.Name)
		trailer[key] = append(trailer[key], hf.Value)
	***REMOVED***
	cs.trailer = trailer

	rl.endStream(cs)
	return nil
***REMOVED***

// transportResponseBody is the concrete type of Transport.RoundTrip's
// Response.Body. It is an io.ReadCloser.
type transportResponseBody struct ***REMOVED***
	cs *clientStream
***REMOVED***

func (b transportResponseBody) Read(p []byte) (n int, err error) ***REMOVED***
	cs := b.cs
	cc := cs.cc

	if cs.readErr != nil ***REMOVED***
		return 0, cs.readErr
	***REMOVED***
	n, err = b.cs.bufPipe.Read(p)
	if cs.bytesRemain != -1 ***REMOVED***
		if int64(n) > cs.bytesRemain ***REMOVED***
			n = int(cs.bytesRemain)
			if err == nil ***REMOVED***
				err = errors.New("net/http: server replied with more than declared Content-Length; truncated")
				cs.abortStream(err)
			***REMOVED***
			cs.readErr = err
			return int(cs.bytesRemain), err
		***REMOVED***
		cs.bytesRemain -= int64(n)
		if err == io.EOF && cs.bytesRemain > 0 ***REMOVED***
			err = io.ErrUnexpectedEOF
			cs.readErr = err
			return n, err
		***REMOVED***
	***REMOVED***
	if n == 0 ***REMOVED***
		// No flow control tokens to send back.
		return
	***REMOVED***

	cc.mu.Lock()
	var connAdd, streamAdd int32
	// Check the conn-level first, before the stream-level.
	if v := cc.inflow.available(); v < transportDefaultConnFlow/2 ***REMOVED***
		connAdd = transportDefaultConnFlow - v
		cc.inflow.add(connAdd)
	***REMOVED***
	if err == nil ***REMOVED*** // No need to refresh if the stream is over or failed.
		// Consider any buffered body data (read from the conn but not
		// consumed by the client) when computing flow control for this
		// stream.
		v := int(cs.inflow.available()) + cs.bufPipe.Len()
		if v < transportDefaultStreamFlow-transportDefaultStreamMinRefresh ***REMOVED***
			streamAdd = int32(transportDefaultStreamFlow - v)
			cs.inflow.add(streamAdd)
		***REMOVED***
	***REMOVED***
	cc.mu.Unlock()

	if connAdd != 0 || streamAdd != 0 ***REMOVED***
		cc.wmu.Lock()
		defer cc.wmu.Unlock()
		if connAdd != 0 ***REMOVED***
			cc.fr.WriteWindowUpdate(0, mustUint31(connAdd))
		***REMOVED***
		if streamAdd != 0 ***REMOVED***
			cc.fr.WriteWindowUpdate(cs.ID, mustUint31(streamAdd))
		***REMOVED***
		cc.bw.Flush()
	***REMOVED***
	return
***REMOVED***

var errClosedResponseBody = errors.New("http2: response body closed")

func (b transportResponseBody) Close() error ***REMOVED***
	cs := b.cs
	cc := cs.cc

	unread := cs.bufPipe.Len()
	if unread > 0 ***REMOVED***
		cc.mu.Lock()
		// Return connection-level flow control.
		if unread > 0 ***REMOVED***
			cc.inflow.add(int32(unread))
		***REMOVED***
		cc.mu.Unlock()

		// TODO(dneil): Acquiring this mutex can block indefinitely.
		// Move flow control return to a goroutine?
		cc.wmu.Lock()
		// Return connection-level flow control.
		if unread > 0 ***REMOVED***
			cc.fr.WriteWindowUpdate(0, uint32(unread))
		***REMOVED***
		cc.bw.Flush()
		cc.wmu.Unlock()
	***REMOVED***

	cs.bufPipe.BreakWithError(errClosedResponseBody)
	cs.abortStream(errClosedResponseBody)

	select ***REMOVED***
	case <-cs.donec:
	case <-cs.ctx.Done():
		// See golang/go#49366: The net/http package can cancel the
		// request context after the response body is fully read.
		// Don't treat this as an error.
		return nil
	case <-cs.reqCancel:
		return errRequestCanceled
	***REMOVED***
	return nil
***REMOVED***

func (rl *clientConnReadLoop) processData(f *DataFrame) error ***REMOVED***
	cc := rl.cc
	cs := rl.streamByID(f.StreamID)
	data := f.Data()
	if cs == nil ***REMOVED***
		cc.mu.Lock()
		neverSent := cc.nextStreamID
		cc.mu.Unlock()
		if f.StreamID >= neverSent ***REMOVED***
			// We never asked for this.
			cc.logf("http2: Transport received unsolicited DATA frame; closing connection")
			return ConnectionError(ErrCodeProtocol)
		***REMOVED***
		// We probably did ask for this, but canceled. Just ignore it.
		// TODO: be stricter here? only silently ignore things which
		// we canceled, but not things which were closed normally
		// by the peer? Tough without accumulating too much state.

		// But at least return their flow control:
		if f.Length > 0 ***REMOVED***
			cc.mu.Lock()
			cc.inflow.add(int32(f.Length))
			cc.mu.Unlock()

			cc.wmu.Lock()
			cc.fr.WriteWindowUpdate(0, uint32(f.Length))
			cc.bw.Flush()
			cc.wmu.Unlock()
		***REMOVED***
		return nil
	***REMOVED***
	if cs.readClosed ***REMOVED***
		cc.logf("protocol error: received DATA after END_STREAM")
		rl.endStreamError(cs, StreamError***REMOVED***
			StreamID: f.StreamID,
			Code:     ErrCodeProtocol,
		***REMOVED***)
		return nil
	***REMOVED***
	if !cs.firstByte ***REMOVED***
		cc.logf("protocol error: received DATA before a HEADERS frame")
		rl.endStreamError(cs, StreamError***REMOVED***
			StreamID: f.StreamID,
			Code:     ErrCodeProtocol,
		***REMOVED***)
		return nil
	***REMOVED***
	if f.Length > 0 ***REMOVED***
		if cs.isHead && len(data) > 0 ***REMOVED***
			cc.logf("protocol error: received DATA on a HEAD request")
			rl.endStreamError(cs, StreamError***REMOVED***
				StreamID: f.StreamID,
				Code:     ErrCodeProtocol,
			***REMOVED***)
			return nil
		***REMOVED***
		// Check connection-level flow control.
		cc.mu.Lock()
		if cs.inflow.available() >= int32(f.Length) ***REMOVED***
			cs.inflow.take(int32(f.Length))
		***REMOVED*** else ***REMOVED***
			cc.mu.Unlock()
			return ConnectionError(ErrCodeFlowControl)
		***REMOVED***
		// Return any padded flow control now, since we won't
		// refund it later on body reads.
		var refund int
		if pad := int(f.Length) - len(data); pad > 0 ***REMOVED***
			refund += pad
		***REMOVED***

		didReset := false
		var err error
		if len(data) > 0 ***REMOVED***
			if _, err = cs.bufPipe.Write(data); err != nil ***REMOVED***
				// Return len(data) now if the stream is already closed,
				// since data will never be read.
				didReset = true
				refund += len(data)
			***REMOVED***
		***REMOVED***

		if refund > 0 ***REMOVED***
			cc.inflow.add(int32(refund))
			if !didReset ***REMOVED***
				cs.inflow.add(int32(refund))
			***REMOVED***
		***REMOVED***
		cc.mu.Unlock()

		if refund > 0 ***REMOVED***
			cc.wmu.Lock()
			cc.fr.WriteWindowUpdate(0, uint32(refund))
			if !didReset ***REMOVED***
				cc.fr.WriteWindowUpdate(cs.ID, uint32(refund))
			***REMOVED***
			cc.bw.Flush()
			cc.wmu.Unlock()
		***REMOVED***

		if err != nil ***REMOVED***
			rl.endStreamError(cs, err)
			return nil
		***REMOVED***
	***REMOVED***

	if f.StreamEnded() ***REMOVED***
		rl.endStream(cs)
	***REMOVED***
	return nil
***REMOVED***

func (rl *clientConnReadLoop) endStream(cs *clientStream) ***REMOVED***
	// TODO: check that any declared content-length matches, like
	// server.go's (*stream).endStream method.
	if !cs.readClosed ***REMOVED***
		cs.readClosed = true
		// Close cs.bufPipe and cs.peerClosed with cc.mu held to avoid a
		// race condition: The caller can read io.EOF from Response.Body
		// and close the body before we close cs.peerClosed, causing
		// cleanupWriteRequest to send a RST_STREAM.
		rl.cc.mu.Lock()
		defer rl.cc.mu.Unlock()
		cs.bufPipe.closeWithErrorAndCode(io.EOF, cs.copyTrailers)
		close(cs.peerClosed)
	***REMOVED***
***REMOVED***

func (rl *clientConnReadLoop) endStreamError(cs *clientStream, err error) ***REMOVED***
	cs.readAborted = true
	cs.abortStream(err)
***REMOVED***

func (rl *clientConnReadLoop) streamByID(id uint32) *clientStream ***REMOVED***
	rl.cc.mu.Lock()
	defer rl.cc.mu.Unlock()
	cs := rl.cc.streams[id]
	if cs != nil && !cs.readAborted ***REMOVED***
		return cs
	***REMOVED***
	return nil
***REMOVED***

func (cs *clientStream) copyTrailers() ***REMOVED***
	for k, vv := range cs.trailer ***REMOVED***
		t := cs.resTrailer
		if *t == nil ***REMOVED***
			*t = make(http.Header)
		***REMOVED***
		(*t)[k] = vv
	***REMOVED***
***REMOVED***

func (rl *clientConnReadLoop) processGoAway(f *GoAwayFrame) error ***REMOVED***
	cc := rl.cc
	cc.t.connPool().MarkDead(cc)
	if f.ErrCode != 0 ***REMOVED***
		// TODO: deal with GOAWAY more. particularly the error code
		cc.vlogf("transport got GOAWAY with error code = %v", f.ErrCode)
		if fn := cc.t.CountError; fn != nil ***REMOVED***
			fn("recv_goaway_" + f.ErrCode.stringToken())
		***REMOVED***

	***REMOVED***
	cc.setGoAway(f)
	return nil
***REMOVED***

func (rl *clientConnReadLoop) processSettings(f *SettingsFrame) error ***REMOVED***
	cc := rl.cc
	// Locking both mu and wmu here allows frame encoding to read settings with only wmu held.
	// Acquiring wmu when f.IsAck() is unnecessary, but convenient and mostly harmless.
	cc.wmu.Lock()
	defer cc.wmu.Unlock()

	if err := rl.processSettingsNoWrite(f); err != nil ***REMOVED***
		return err
	***REMOVED***
	if !f.IsAck() ***REMOVED***
		cc.fr.WriteSettingsAck()
		cc.bw.Flush()
	***REMOVED***
	return nil
***REMOVED***

func (rl *clientConnReadLoop) processSettingsNoWrite(f *SettingsFrame) error ***REMOVED***
	cc := rl.cc
	cc.mu.Lock()
	defer cc.mu.Unlock()

	if f.IsAck() ***REMOVED***
		if cc.wantSettingsAck ***REMOVED***
			cc.wantSettingsAck = false
			return nil
		***REMOVED***
		return ConnectionError(ErrCodeProtocol)
	***REMOVED***

	var seenMaxConcurrentStreams bool
	err := f.ForeachSetting(func(s Setting) error ***REMOVED***
		switch s.ID ***REMOVED***
		case SettingMaxFrameSize:
			cc.maxFrameSize = s.Val
		case SettingMaxConcurrentStreams:
			cc.maxConcurrentStreams = s.Val
			seenMaxConcurrentStreams = true
		case SettingMaxHeaderListSize:
			cc.peerMaxHeaderListSize = uint64(s.Val)
		case SettingInitialWindowSize:
			// Values above the maximum flow-control
			// window size of 2^31-1 MUST be treated as a
			// connection error (Section 5.4.1) of type
			// FLOW_CONTROL_ERROR.
			if s.Val > math.MaxInt32 ***REMOVED***
				return ConnectionError(ErrCodeFlowControl)
			***REMOVED***

			// Adjust flow control of currently-open
			// frames by the difference of the old initial
			// window size and this one.
			delta := int32(s.Val) - int32(cc.initialWindowSize)
			for _, cs := range cc.streams ***REMOVED***
				cs.flow.add(delta)
			***REMOVED***
			cc.cond.Broadcast()

			cc.initialWindowSize = s.Val
		default:
			// TODO(bradfitz): handle more settings? SETTINGS_HEADER_TABLE_SIZE probably.
			cc.vlogf("Unhandled Setting: %v", s)
		***REMOVED***
		return nil
	***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if !cc.seenSettings ***REMOVED***
		if !seenMaxConcurrentStreams ***REMOVED***
			// This was the servers initial SETTINGS frame and it
			// didn't contain a MAX_CONCURRENT_STREAMS field so
			// increase the number of concurrent streams this
			// connection can establish to our default.
			cc.maxConcurrentStreams = defaultMaxConcurrentStreams
		***REMOVED***
		cc.seenSettings = true
	***REMOVED***

	return nil
***REMOVED***

func (rl *clientConnReadLoop) processWindowUpdate(f *WindowUpdateFrame) error ***REMOVED***
	cc := rl.cc
	cs := rl.streamByID(f.StreamID)
	if f.StreamID != 0 && cs == nil ***REMOVED***
		return nil
	***REMOVED***

	cc.mu.Lock()
	defer cc.mu.Unlock()

	fl := &cc.flow
	if cs != nil ***REMOVED***
		fl = &cs.flow
	***REMOVED***
	if !fl.add(int32(f.Increment)) ***REMOVED***
		return ConnectionError(ErrCodeFlowControl)
	***REMOVED***
	cc.cond.Broadcast()
	return nil
***REMOVED***

func (rl *clientConnReadLoop) processResetStream(f *RSTStreamFrame) error ***REMOVED***
	cs := rl.streamByID(f.StreamID)
	if cs == nil ***REMOVED***
		// TODO: return error if server tries to RST_STREAM an idle stream
		return nil
	***REMOVED***
	serr := streamError(cs.ID, f.ErrCode)
	serr.Cause = errFromPeer
	if f.ErrCode == ErrCodeProtocol ***REMOVED***
		rl.cc.SetDoNotReuse()
	***REMOVED***
	if fn := cs.cc.t.CountError; fn != nil ***REMOVED***
		fn("recv_rststream_" + f.ErrCode.stringToken())
	***REMOVED***
	cs.abortStream(serr)

	cs.bufPipe.CloseWithError(serr)
	return nil
***REMOVED***

// Ping sends a PING frame to the server and waits for the ack.
func (cc *ClientConn) Ping(ctx context.Context) error ***REMOVED***
	c := make(chan struct***REMOVED******REMOVED***)
	// Generate a random payload
	var p [8]byte
	for ***REMOVED***
		if _, err := rand.Read(p[:]); err != nil ***REMOVED***
			return err
		***REMOVED***
		cc.mu.Lock()
		// check for dup before insert
		if _, found := cc.pings[p]; !found ***REMOVED***
			cc.pings[p] = c
			cc.mu.Unlock()
			break
		***REMOVED***
		cc.mu.Unlock()
	***REMOVED***
	errc := make(chan error, 1)
	go func() ***REMOVED***
		cc.wmu.Lock()
		defer cc.wmu.Unlock()
		if err := cc.fr.WritePing(false, p); err != nil ***REMOVED***
			errc <- err
			return
		***REMOVED***
		if err := cc.bw.Flush(); err != nil ***REMOVED***
			errc <- err
			return
		***REMOVED***
	***REMOVED***()
	select ***REMOVED***
	case <-c:
		return nil
	case err := <-errc:
		return err
	case <-ctx.Done():
		return ctx.Err()
	case <-cc.readerDone:
		// connection closed
		return cc.readerErr
	***REMOVED***
***REMOVED***

func (rl *clientConnReadLoop) processPing(f *PingFrame) error ***REMOVED***
	if f.IsAck() ***REMOVED***
		cc := rl.cc
		cc.mu.Lock()
		defer cc.mu.Unlock()
		// If ack, notify listener if any
		if c, ok := cc.pings[f.Data]; ok ***REMOVED***
			close(c)
			delete(cc.pings, f.Data)
		***REMOVED***
		return nil
	***REMOVED***
	cc := rl.cc
	cc.wmu.Lock()
	defer cc.wmu.Unlock()
	if err := cc.fr.WritePing(true, f.Data); err != nil ***REMOVED***
		return err
	***REMOVED***
	return cc.bw.Flush()
***REMOVED***

func (rl *clientConnReadLoop) processPushPromise(f *PushPromiseFrame) error ***REMOVED***
	// We told the peer we don't want them.
	// Spec says:
	// "PUSH_PROMISE MUST NOT be sent if the SETTINGS_ENABLE_PUSH
	// setting of the peer endpoint is set to 0. An endpoint that
	// has set this setting and has received acknowledgement MUST
	// treat the receipt of a PUSH_PROMISE frame as a connection
	// error (Section 5.4.1) of type PROTOCOL_ERROR."
	return ConnectionError(ErrCodeProtocol)
***REMOVED***

func (cc *ClientConn) writeStreamReset(streamID uint32, code ErrCode, err error) ***REMOVED***
	// TODO: map err to more interesting error codes, once the
	// HTTP community comes up with some. But currently for
	// RST_STREAM there's no equivalent to GOAWAY frame's debug
	// data, and the error codes are all pretty vague ("cancel").
	cc.wmu.Lock()
	cc.fr.WriteRSTStream(streamID, code)
	cc.bw.Flush()
	cc.wmu.Unlock()
***REMOVED***

var (
	errResponseHeaderListSize = errors.New("http2: response header list larger than advertised limit")
	errRequestHeaderListSize  = errors.New("http2: request header list larger than peer's advertised limit")
)

func (cc *ClientConn) logf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	cc.t.logf(format, args...)
***REMOVED***

func (cc *ClientConn) vlogf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	cc.t.vlogf(format, args...)
***REMOVED***

func (t *Transport) vlogf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if VerboseLogs ***REMOVED***
		t.logf(format, args...)
	***REMOVED***
***REMOVED***

func (t *Transport) logf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	log.Printf(format, args...)
***REMOVED***

var noBody io.ReadCloser = ioutil.NopCloser(bytes.NewReader(nil))

type missingBody struct***REMOVED******REMOVED***

func (missingBody) Close() error             ***REMOVED*** return nil ***REMOVED***
func (missingBody) Read([]byte) (int, error) ***REMOVED*** return 0, io.ErrUnexpectedEOF ***REMOVED***

func strSliceContains(ss []string, s string) bool ***REMOVED***
	for _, v := range ss ***REMOVED***
		if v == s ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

type erringRoundTripper struct***REMOVED*** err error ***REMOVED***

func (rt erringRoundTripper) RoundTripErr() error                             ***REMOVED*** return rt.err ***REMOVED***
func (rt erringRoundTripper) RoundTrip(*http.Request) (*http.Response, error) ***REMOVED*** return nil, rt.err ***REMOVED***

// gzipReader wraps a response body so it can lazily
// call gzip.NewReader on the first call to Read
type gzipReader struct ***REMOVED***
	_    incomparable
	body io.ReadCloser // underlying Response.Body
	zr   *gzip.Reader  // lazily-initialized gzip reader
	zerr error         // sticky error
***REMOVED***

func (gz *gzipReader) Read(p []byte) (n int, err error) ***REMOVED***
	if gz.zerr != nil ***REMOVED***
		return 0, gz.zerr
	***REMOVED***
	if gz.zr == nil ***REMOVED***
		gz.zr, err = gzip.NewReader(gz.body)
		if err != nil ***REMOVED***
			gz.zerr = err
			return 0, err
		***REMOVED***
	***REMOVED***
	return gz.zr.Read(p)
***REMOVED***

func (gz *gzipReader) Close() error ***REMOVED***
	return gz.body.Close()
***REMOVED***

type errorReader struct***REMOVED*** err error ***REMOVED***

func (r errorReader) Read(p []byte) (int, error) ***REMOVED*** return 0, r.err ***REMOVED***

// isConnectionCloseRequest reports whether req should use its own
// connection for a single request and then close the connection.
func isConnectionCloseRequest(req *http.Request) bool ***REMOVED***
	return req.Close || httpguts.HeaderValuesContainsToken(req.Header["Connection"], "close")
***REMOVED***

// registerHTTPSProtocol calls Transport.RegisterProtocol but
// converting panics into errors.
func registerHTTPSProtocol(t *http.Transport, rt noDialH2RoundTripper) (err error) ***REMOVED***
	defer func() ***REMOVED***
		if e := recover(); e != nil ***REMOVED***
			err = fmt.Errorf("%v", e)
		***REMOVED***
	***REMOVED***()
	t.RegisterProtocol("https", rt)
	return nil
***REMOVED***

// noDialH2RoundTripper is a RoundTripper which only tries to complete the request
// if there's already has a cached connection to the host.
// (The field is exported so it can be accessed via reflect from net/http; tested
// by TestNoDialH2RoundTripperType)
type noDialH2RoundTripper struct***REMOVED*** *Transport ***REMOVED***

func (rt noDialH2RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) ***REMOVED***
	res, err := rt.Transport.RoundTrip(req)
	if isNoCachedConnError(err) ***REMOVED***
		return nil, http.ErrSkipAltProtocol
	***REMOVED***
	return res, err
***REMOVED***

func (t *Transport) idleConnTimeout() time.Duration ***REMOVED***
	if t.t1 != nil ***REMOVED***
		return t.t1.IdleConnTimeout
	***REMOVED***
	return 0
***REMOVED***

func traceGetConn(req *http.Request, hostPort string) ***REMOVED***
	trace := httptrace.ContextClientTrace(req.Context())
	if trace == nil || trace.GetConn == nil ***REMOVED***
		return
	***REMOVED***
	trace.GetConn(hostPort)
***REMOVED***

func traceGotConn(req *http.Request, cc *ClientConn, reused bool) ***REMOVED***
	trace := httptrace.ContextClientTrace(req.Context())
	if trace == nil || trace.GotConn == nil ***REMOVED***
		return
	***REMOVED***
	ci := httptrace.GotConnInfo***REMOVED***Conn: cc.tconn***REMOVED***
	ci.Reused = reused
	cc.mu.Lock()
	ci.WasIdle = len(cc.streams) == 0 && reused
	if ci.WasIdle && !cc.lastActive.IsZero() ***REMOVED***
		ci.IdleTime = time.Now().Sub(cc.lastActive)
	***REMOVED***
	cc.mu.Unlock()

	trace.GotConn(ci)
***REMOVED***

func traceWroteHeaders(trace *httptrace.ClientTrace) ***REMOVED***
	if trace != nil && trace.WroteHeaders != nil ***REMOVED***
		trace.WroteHeaders()
	***REMOVED***
***REMOVED***

func traceGot100Continue(trace *httptrace.ClientTrace) ***REMOVED***
	if trace != nil && trace.Got100Continue != nil ***REMOVED***
		trace.Got100Continue()
	***REMOVED***
***REMOVED***

func traceWait100Continue(trace *httptrace.ClientTrace) ***REMOVED***
	if trace != nil && trace.Wait100Continue != nil ***REMOVED***
		trace.Wait100Continue()
	***REMOVED***
***REMOVED***

func traceWroteRequest(trace *httptrace.ClientTrace, err error) ***REMOVED***
	if trace != nil && trace.WroteRequest != nil ***REMOVED***
		trace.WroteRequest(httptrace.WroteRequestInfo***REMOVED***Err: err***REMOVED***)
	***REMOVED***
***REMOVED***

func traceFirstResponseByte(trace *httptrace.ClientTrace) ***REMOVED***
	if trace != nil && trace.GotFirstResponseByte != nil ***REMOVED***
		trace.GotFirstResponseByte()
	***REMOVED***
***REMOVED***
