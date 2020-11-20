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
func ConfigureTransport(t1 *http.Transport) error ***REMOVED***
	_, err := configureTransport(t1)
	return err
***REMOVED***

func configureTransport(t1 *http.Transport) (*Transport, error) ***REMOVED***
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
	t         *Transport
	tconn     net.Conn             // usually *tls.Conn, except specialized impls
	tlsState  *tls.ConnectionState // nil only for specialized impls
	reused    uint32               // whether conn is being reused; atomic
	singleUse bool                 // whether being used for a single http.Request

	// readLoop goroutine fields:
	readerDone chan struct***REMOVED******REMOVED*** // closed on error
	readerErr  error         // set before readerDone is closed

	idleTimeout time.Duration // or 0 for never
	idleTimer   *time.Timer

	mu              sync.Mutex // guards following
	cond            *sync.Cond // hold mu; broadcast on flow/closed changes
	flow            flow       // our conn-level flow control quota (cs.flow is per stream)
	inflow          flow       // peer's conn-level flow control
	closing         bool
	closed          bool
	wantSettingsAck bool                     // we sent a SETTINGS frame and haven't heard back
	goAway          *GoAwayFrame             // if non-nil, the GoAwayFrame we received
	goAwayDebug     string                   // goAway frame's debug data, retained as a string
	streams         map[uint32]*clientStream // client-initiated
	nextStreamID    uint32
	pendingRequests int                       // requests blocked and waiting to be sent because len(streams) == maxConcurrentStreams
	pings           map[[8]byte]chan struct***REMOVED******REMOVED*** // in flight ping data to notification channel
	bw              *bufio.Writer
	br              *bufio.Reader
	fr              *Framer
	lastActive      time.Time
	lastIdle        time.Time // time last idle
	// Settings from peer: (also guarded by mu)
	maxFrameSize          uint32
	maxConcurrentStreams  uint32
	peerMaxHeaderListSize uint64
	initialWindowSize     uint32

	hbuf    bytes.Buffer // HPACK encoder writes into this
	henc    *hpack.Encoder
	freeBuf [][]byte

	wmu  sync.Mutex // held while writing; acquire AFTER mu if holding both
	werr error      // first write error that has occurred
***REMOVED***

// clientStream is the state for a single HTTP/2 stream. One of these
// is created for each Transport.RoundTrip call.
type clientStream struct ***REMOVED***
	cc            *ClientConn
	req           *http.Request
	trace         *httptrace.ClientTrace // or nil
	ID            uint32
	resc          chan resAndError
	bufPipe       pipe // buffered pipe with the flow-controlled response payload
	startedWrite  bool // started request body write; guarded by cc.mu
	requestedGzip bool
	on100         func() // optional code to run if get a 100 continue response

	flow        flow  // guarded by cc.mu
	inflow      flow  // guarded by cc.mu
	bytesRemain int64 // -1 means unknown; owned by transportResponseBody.Read
	readErr     error // sticky read error; owned by transportResponseBody.Read
	stopReqBody error // if non-nil, stop writing req body; guarded by cc.mu
	didReset    bool  // whether we sent a RST_STREAM to the server; guarded by cc.mu

	peerReset chan struct***REMOVED******REMOVED*** // closed on peer reset
	resetErr  error         // populated before peerReset is closed

	done chan struct***REMOVED******REMOVED*** // closed when stream remove from cc.streams map; close calls guarded by cc.mu

	// owned by clientConnReadLoop:
	firstByte    bool  // got the first response byte
	pastHeaders  bool  // got first MetaHeadersFrame (actual headers)
	pastTrailers bool  // got optional second MetaHeadersFrame (trailers)
	num1xx       uint8 // number of 1xx responses seen

	trailer    http.Header  // accumulated trailers
	resTrailer *http.Header // client's Response.Trailer
***REMOVED***

// awaitRequestCancel waits for the user to cancel a request or for the done
// channel to be signaled. A non-nil error is returned only if the request was
// canceled.
func awaitRequestCancel(req *http.Request, done <-chan struct***REMOVED******REMOVED***) error ***REMOVED***
	ctx := req.Context()
	if req.Cancel == nil && ctx.Done() == nil ***REMOVED***
		return nil
	***REMOVED***
	select ***REMOVED***
	case <-req.Cancel:
		return errRequestCanceled
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	***REMOVED***
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

// awaitRequestCancel waits for the user to cancel a request, its context to
// expire, or for the request to be done (any way it might be removed from the
// cc.streams map: peer reset, successful completion, TCP connection breakage,
// etc). If the request is canceled, then cs will be canceled and closed.
func (cs *clientStream) awaitRequestCancel(req *http.Request) ***REMOVED***
	if err := awaitRequestCancel(req, cs.done); err != nil ***REMOVED***
		cs.cancelStream()
		cs.bufPipe.CloseWithError(err)
	***REMOVED***
***REMOVED***

func (cs *clientStream) cancelStream() ***REMOVED***
	cc := cs.cc
	cc.mu.Lock()
	didReset := cs.didReset
	cs.didReset = true
	cc.mu.Unlock()

	if !didReset ***REMOVED***
		cc.writeStreamReset(cs.ID, ErrCodeCancel, nil)
		cc.forgetStreamID(cs.ID)
	***REMOVED***
***REMOVED***

// checkResetOrDone reports any error sent in a RST_STREAM frame by the
// server, or errStreamClosed if the stream is complete.
func (cs *clientStream) checkResetOrDone() error ***REMOVED***
	select ***REMOVED***
	case <-cs.peerReset:
		return cs.resetErr
	case <-cs.done:
		return errStreamClosed
	default:
		return nil
	***REMOVED***
***REMOVED***

func (cs *clientStream) getStartedWrite() bool ***REMOVED***
	cc := cs.cc
	cc.mu.Lock()
	defer cc.mu.Unlock()
	return cs.startedWrite
***REMOVED***

func (cs *clientStream) abortRequestBodyWrite(err error) ***REMOVED***
	if err == nil ***REMOVED***
		panic("nil error")
	***REMOVED***
	cc := cs.cc
	cc.mu.Lock()
	cs.stopReqBody = err
	cc.cond.Broadcast()
	cc.mu.Unlock()
***REMOVED***

type stickyErrWriter struct ***REMOVED***
	w   io.Writer
	err *error
***REMOVED***

func (sew stickyErrWriter) Write(p []byte) (n int, err error) ***REMOVED***
	if *sew.err != nil ***REMOVED***
		return 0, *sew.err
	***REMOVED***
	n, err = sew.w.Write(p)
	*sew.err = err
	return
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
		res, gotErrAfterReqBodyWrite, err := cc.roundTrip(req)
		if err != nil && retry <= 6 ***REMOVED***
			if req, err = shouldRetryRequest(req, err, gotErrAfterReqBodyWrite); err == nil ***REMOVED***
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
					return nil, req.Context().Err()
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
func shouldRetryRequest(req *http.Request, err error, afterBodyWrite bool) (*http.Request, error) ***REMOVED***
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
		// TODO: consider a req.Body.Close here? or audit that all caller paths do?
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
	// the request directly. The "afterBodyWrite" means the
	// bodyWrite process has started, which becomes true before
	// the first Read.
	if !afterBodyWrite ***REMOVED***
		return req, nil
	***REMOVED***

	return nil, fmt.Errorf("http2: Transport: cannot retry err [%v] after Request.Body was written; define Request.GetBody to avoid this error", err)
***REMOVED***

func canRetryError(err error) bool ***REMOVED***
	if err == errClientConnUnusable || err == errClientConnGotGoAway ***REMOVED***
		return true
	***REMOVED***
	if se, ok := err.(StreamError); ok ***REMOVED***
		return se.Code == ErrCodeRefusedStream
	***REMOVED***
	return false
***REMOVED***

func (t *Transport) dialClientConn(addr string, singleUse bool) (*ClientConn, error) ***REMOVED***
	host, _, err := net.SplitHostPort(addr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	tconn, err := t.dialTLS()("tcp", addr, t.newTLSConfig(host))
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

func (t *Transport) dialTLS() func(string, string, *tls.Config) (net.Conn, error) ***REMOVED***
	if t.DialTLS != nil ***REMOVED***
		return t.DialTLS
	***REMOVED***
	return t.dialTLSDefault
***REMOVED***

func (t *Transport) dialTLSDefault(network, addr string, cfg *tls.Config) (net.Conn, error) ***REMOVED***
	cn, err := tls.Dial(network, addr, cfg)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := cn.Handshake(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if !cfg.InsecureSkipVerify ***REMOVED***
		if err := cn.VerifyHostname(cfg.ServerName); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	state := cn.ConnectionState()
	if p := state.NegotiatedProtocol; p != NextProtoTLS ***REMOVED***
		return nil, fmt.Errorf("http2: unexpected ALPN protocol %q; want %q", p, NextProtoTLS)
	***REMOVED***
	if !state.NegotiatedProtocolIsMutual ***REMOVED***
		return nil, errors.New("http2: could not negotiate protocol mutually")
	***REMOVED***
	return cn, nil
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
		maxFrameSize:          16 << 10,           // spec default
		initialWindowSize:     65535,              // spec default
		maxConcurrentStreams:  1000,               // "infinite", per spec. 1000 seems good enough.
		peerMaxHeaderListSize: 0xffffffffffffffff, // "infinite", per spec. Use 2^64-1 instead.
		streams:               make(map[uint32]*clientStream),
		singleUse:             singleUse,
		wantSettingsAck:       true,
		pings:                 make(map[[8]byte]chan struct***REMOVED******REMOVED***),
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
	cc.bw = bufio.NewWriter(stickyErrWriter***REMOVED***c, &cc.werr***REMOVED***)
	cc.br = bufio.NewReader(c)
	cc.fr = NewFramer(cc.bw, cc.br)
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
		cc.t.connPool().MarkDead(cc)
		return
	***REMOVED***
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
			select ***REMOVED***
			case cs.resc <- resAndError***REMOVED***err: errClientConnGotGoAway***REMOVED***:
			default:
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// CanTakeNewRequest reports whether the connection can take a new request,
// meaning it has not been closed or received or sent a GOAWAY.
func (cc *ClientConn) CanTakeNewRequest() bool ***REMOVED***
	cc.mu.Lock()
	defer cc.mu.Unlock()
	return cc.canTakeNewRequestLocked()
***REMOVED***

// clientConnIdleState describes the suitability of a client
// connection to initiate a new RoundTrip request.
type clientConnIdleState struct ***REMOVED***
	canTakeNewRequest bool
	freshConn         bool // whether it's unused by any previous request
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
		maxConcurrentOkay = int64(len(cc.streams)+1) < int64(cc.maxConcurrentStreams)
	***REMOVED***

	st.canTakeNewRequest = cc.goAway == nil && !cc.closed && !cc.closing && maxConcurrentOkay &&
		int64(cc.nextStreamID)+2*int64(cc.pendingRequests) < math.MaxInt32 &&
		!cc.tooIdleLocked()
	st.freshConn = cc.nextStreamID == 1 && st.canTakeNewRequest
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

func (cc *ClientConn) closeIfIdle() ***REMOVED***
	cc.mu.Lock()
	if len(cc.streams) > 0 ***REMOVED***
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
	cc.tconn.Close()
***REMOVED***

var shutdownEnterWaitStateHook = func() ***REMOVED******REMOVED***

// Shutdown gracefully close the client connection, waiting for running streams to complete.
func (cc *ClientConn) Shutdown(ctx context.Context) error ***REMOVED***
	if err := cc.sendGoAway(); err != nil ***REMOVED***
		return err
	***REMOVED***
	// Wait for all in-flight streams to complete or connection to close
	done := make(chan error, 1)
	cancelled := false // guarded by cc.mu
	go func() ***REMOVED***
		cc.mu.Lock()
		defer cc.mu.Unlock()
		for ***REMOVED***
			if len(cc.streams) == 0 || cc.closed ***REMOVED***
				cc.closed = true
				done <- cc.tconn.Close()
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
	case err := <-done:
		return err
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
	defer cc.mu.Unlock()
	cc.wmu.Lock()
	defer cc.wmu.Unlock()
	if cc.closing ***REMOVED***
		// GOAWAY sent already
		return nil
	***REMOVED***
	// Send a graceful shutdown frame to server
	maxStreamID := cc.nextStreamID
	if err := cc.fr.WriteGoAway(maxStreamID, ErrCodeNo, nil); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := cc.bw.Flush(); err != nil ***REMOVED***
		return err
	***REMOVED***
	// Prevent new requests
	cc.closing = true
	return nil
***REMOVED***

// closes the client connection immediately. In-flight requests are interrupted.
// err is sent to streams.
func (cc *ClientConn) closeForError(err error) error ***REMOVED***
	cc.mu.Lock()
	defer cc.cond.Broadcast()
	defer cc.mu.Unlock()
	for id, cs := range cc.streams ***REMOVED***
		select ***REMOVED***
		case cs.resc <- resAndError***REMOVED***err: err***REMOVED***:
		default:
		***REMOVED***
		cs.bufPipe.CloseWithError(err)
		delete(cc.streams, id)
	***REMOVED***
	cc.closed = true
	return cc.tconn.Close()
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
	return cc.closeForError(err)
***REMOVED***

const maxAllocFrameSize = 512 << 10

// frameBuffer returns a scratch buffer suitable for writing DATA frames.
// They're capped at the min of the peer's max frame size or 512KB
// (kinda arbitrarily), but definitely capped so we don't allocate 4GB
// bufers.
func (cc *ClientConn) frameScratchBuffer() []byte ***REMOVED***
	cc.mu.Lock()
	size := cc.maxFrameSize
	if size > maxAllocFrameSize ***REMOVED***
		size = maxAllocFrameSize
	***REMOVED***
	for i, buf := range cc.freeBuf ***REMOVED***
		if len(buf) >= int(size) ***REMOVED***
			cc.freeBuf[i] = nil
			cc.mu.Unlock()
			return buf[:size]
		***REMOVED***
	***REMOVED***
	cc.mu.Unlock()
	return make([]byte, size)
***REMOVED***

func (cc *ClientConn) putFrameScratchBuffer(buf []byte) ***REMOVED***
	cc.mu.Lock()
	defer cc.mu.Unlock()
	const maxBufs = 4 // arbitrary; 4 concurrent requests per conn? investigate.
	if len(cc.freeBuf) < maxBufs ***REMOVED***
		cc.freeBuf = append(cc.freeBuf, buf)
		return
	***REMOVED***
	for i, old := range cc.freeBuf ***REMOVED***
		if old == nil ***REMOVED***
			cc.freeBuf[i] = buf
			return
		***REMOVED***
	***REMOVED***
	// forget about it.
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
	if vv := req.Header["Connection"]; len(vv) > 0 && (len(vv) > 1 || vv[0] != "" && !strings.EqualFold(vv[0], "close") && !strings.EqualFold(vv[0], "keep-alive")) ***REMOVED***
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

func (cc *ClientConn) RoundTrip(req *http.Request) (*http.Response, error) ***REMOVED***
	resp, _, err := cc.roundTrip(req)
	return resp, err
***REMOVED***

func (cc *ClientConn) roundTrip(req *http.Request) (res *http.Response, gotErrAfterReqBodyWrite bool, err error) ***REMOVED***
	if err := checkConnHeaders(req); err != nil ***REMOVED***
		return nil, false, err
	***REMOVED***
	if cc.idleTimer != nil ***REMOVED***
		cc.idleTimer.Stop()
	***REMOVED***

	trailers, err := commaSeparatedTrailers(req)
	if err != nil ***REMOVED***
		return nil, false, err
	***REMOVED***
	hasTrailers := trailers != ""

	cc.mu.Lock()
	if err := cc.awaitOpenSlotForRequest(req); err != nil ***REMOVED***
		cc.mu.Unlock()
		return nil, false, err
	***REMOVED***

	body := req.Body
	contentLen := actualContentLength(req)
	hasBody := contentLen != 0

	// TODO(bradfitz): this is a copy of the logic in net/http. Unify somewhere?
	var requestedGzip bool
	if !cc.t.disableCompression() &&
		req.Header.Get("Accept-Encoding") == "" &&
		req.Header.Get("Range") == "" &&
		req.Method != "HEAD" ***REMOVED***
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
		requestedGzip = true
	***REMOVED***

	// we send: HEADERS***REMOVED***1***REMOVED***, CONTINUATION***REMOVED***0,***REMOVED*** + DATA***REMOVED***0,***REMOVED*** (DATA is
	// sent by writeRequestBody below, along with any Trailers,
	// again in form HEADERS***REMOVED***1***REMOVED***, CONTINUATION***REMOVED***0,***REMOVED***)
	hdrs, err := cc.encodeHeaders(req, requestedGzip, trailers, contentLen)
	if err != nil ***REMOVED***
		cc.mu.Unlock()
		return nil, false, err
	***REMOVED***

	cs := cc.newStream()
	cs.req = req
	cs.trace = httptrace.ContextClientTrace(req.Context())
	cs.requestedGzip = requestedGzip
	bodyWriter := cc.t.getBodyWriterState(cs, body)
	cs.on100 = bodyWriter.on100

	cc.wmu.Lock()
	endStream := !hasBody && !hasTrailers
	werr := cc.writeHeaders(cs.ID, endStream, int(cc.maxFrameSize), hdrs)
	cc.wmu.Unlock()
	traceWroteHeaders(cs.trace)
	cc.mu.Unlock()

	if werr != nil ***REMOVED***
		if hasBody ***REMOVED***
			req.Body.Close() // per RoundTripper contract
			bodyWriter.cancel()
		***REMOVED***
		cc.forgetStreamID(cs.ID)
		// Don't bother sending a RST_STREAM (our write already failed;
		// no need to keep writing)
		traceWroteRequest(cs.trace, werr)
		return nil, false, werr
	***REMOVED***

	var respHeaderTimer <-chan time.Time
	if hasBody ***REMOVED***
		bodyWriter.scheduleBodyWrite()
	***REMOVED*** else ***REMOVED***
		traceWroteRequest(cs.trace, nil)
		if d := cc.responseHeaderTimeout(); d != 0 ***REMOVED***
			timer := time.NewTimer(d)
			defer timer.Stop()
			respHeaderTimer = timer.C
		***REMOVED***
	***REMOVED***

	readLoopResCh := cs.resc
	bodyWritten := false
	ctx := req.Context()

	handleReadLoopResponse := func(re resAndError) (*http.Response, bool, error) ***REMOVED***
		res := re.res
		if re.err != nil || res.StatusCode > 299 ***REMOVED***
			// On error or status code 3xx, 4xx, 5xx, etc abort any
			// ongoing write, assuming that the server doesn't care
			// about our request body. If the server replied with 1xx or
			// 2xx, however, then assume the server DOES potentially
			// want our body (e.g. full-duplex streaming:
			// golang.org/issue/13444). If it turns out the server
			// doesn't, they'll RST_STREAM us soon enough. This is a
			// heuristic to avoid adding knobs to Transport. Hopefully
			// we can keep it.
			bodyWriter.cancel()
			cs.abortRequestBodyWrite(errStopReqBodyWrite)
		***REMOVED***
		if re.err != nil ***REMOVED***
			cc.forgetStreamID(cs.ID)
			return nil, cs.getStartedWrite(), re.err
		***REMOVED***
		res.Request = req
		res.TLS = cc.tlsState
		return res, false, nil
	***REMOVED***

	for ***REMOVED***
		select ***REMOVED***
		case re := <-readLoopResCh:
			return handleReadLoopResponse(re)
		case <-respHeaderTimer:
			if !hasBody || bodyWritten ***REMOVED***
				cc.writeStreamReset(cs.ID, ErrCodeCancel, nil)
			***REMOVED*** else ***REMOVED***
				bodyWriter.cancel()
				cs.abortRequestBodyWrite(errStopReqBodyWriteAndCancel)
			***REMOVED***
			cc.forgetStreamID(cs.ID)
			return nil, cs.getStartedWrite(), errTimeout
		case <-ctx.Done():
			if !hasBody || bodyWritten ***REMOVED***
				cc.writeStreamReset(cs.ID, ErrCodeCancel, nil)
			***REMOVED*** else ***REMOVED***
				bodyWriter.cancel()
				cs.abortRequestBodyWrite(errStopReqBodyWriteAndCancel)
			***REMOVED***
			cc.forgetStreamID(cs.ID)
			return nil, cs.getStartedWrite(), ctx.Err()
		case <-req.Cancel:
			if !hasBody || bodyWritten ***REMOVED***
				cc.writeStreamReset(cs.ID, ErrCodeCancel, nil)
			***REMOVED*** else ***REMOVED***
				bodyWriter.cancel()
				cs.abortRequestBodyWrite(errStopReqBodyWriteAndCancel)
			***REMOVED***
			cc.forgetStreamID(cs.ID)
			return nil, cs.getStartedWrite(), errRequestCanceled
		case <-cs.peerReset:
			// processResetStream already removed the
			// stream from the streams map; no need for
			// forgetStreamID.
			return nil, cs.getStartedWrite(), cs.resetErr
		case err := <-bodyWriter.resc:
			// Prefer the read loop's response, if available. Issue 16102.
			select ***REMOVED***
			case re := <-readLoopResCh:
				return handleReadLoopResponse(re)
			default:
			***REMOVED***
			if err != nil ***REMOVED***
				cc.forgetStreamID(cs.ID)
				return nil, cs.getStartedWrite(), err
			***REMOVED***
			bodyWritten = true
			if d := cc.responseHeaderTimeout(); d != 0 ***REMOVED***
				timer := time.NewTimer(d)
				defer timer.Stop()
				respHeaderTimer = timer.C
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// awaitOpenSlotForRequest waits until len(streams) < maxConcurrentStreams.
// Must hold cc.mu.
func (cc *ClientConn) awaitOpenSlotForRequest(req *http.Request) error ***REMOVED***
	var waitingForConn chan struct***REMOVED******REMOVED***
	var waitingForConnErr error // guarded by cc.mu
	for ***REMOVED***
		cc.lastActive = time.Now()
		if cc.closed || !cc.canTakeNewRequestLocked() ***REMOVED***
			if waitingForConn != nil ***REMOVED***
				close(waitingForConn)
			***REMOVED***
			return errClientConnUnusable
		***REMOVED***
		cc.lastIdle = time.Time***REMOVED******REMOVED***
		if int64(len(cc.streams))+1 <= int64(cc.maxConcurrentStreams) ***REMOVED***
			if waitingForConn != nil ***REMOVED***
				close(waitingForConn)
			***REMOVED***
			return nil
		***REMOVED***
		// Unfortunately, we cannot wait on a condition variable and channel at
		// the same time, so instead, we spin up a goroutine to check if the
		// request is canceled while we wait for a slot to open in the connection.
		if waitingForConn == nil ***REMOVED***
			waitingForConn = make(chan struct***REMOVED******REMOVED***)
			go func() ***REMOVED***
				if err := awaitRequestCancel(req, waitingForConn); err != nil ***REMOVED***
					cc.mu.Lock()
					waitingForConnErr = err
					cc.cond.Broadcast()
					cc.mu.Unlock()
				***REMOVED***
			***REMOVED***()
		***REMOVED***
		cc.pendingRequests++
		cc.cond.Wait()
		cc.pendingRequests--
		if waitingForConnErr != nil ***REMOVED***
			return waitingForConnErr
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
	// TODO(bradfitz): this Flush could potentially block (as
	// could the WriteHeaders call(s) above), which means they
	// wouldn't respond to Request.Cancel being readable. That's
	// rare, but this should probably be in a goroutine.
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

func (cs *clientStream) writeRequestBody(body io.Reader, bodyCloser io.Closer) (err error) ***REMOVED***
	cc := cs.cc
	sentEnd := false // whether we sent the final DATA frame w/ END_STREAM
	buf := cc.frameScratchBuffer()
	defer cc.putFrameScratchBuffer(buf)

	defer func() ***REMOVED***
		traceWroteRequest(cs.trace, err)
		// TODO: write h12Compare test showing whether
		// Request.Body is closed by the Transport,
		// and in multiple cases: server replies <=299 and >299
		// while still writing request body
		cerr := bodyCloser.Close()
		if err == nil ***REMOVED***
			err = cerr
		***REMOVED***
	***REMOVED***()

	req := cs.req
	hasTrailers := req.Trailer != nil
	remainLen := actualContentLength(req)
	hasContentLen := remainLen != -1

	var sawEOF bool
	for !sawEOF ***REMOVED***
		n, err := body.Read(buf[:len(buf)-1])
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
				var n1 int
				n1, err = body.Read(buf[n:])
				remainLen -= int64(n1)
			***REMOVED***
			if remainLen < 0 ***REMOVED***
				err = errReqBodyTooLong
				cc.writeStreamReset(cs.ID, ErrCodeCancel, err)
				return err
			***REMOVED***
		***REMOVED***
		if err == io.EOF ***REMOVED***
			sawEOF = true
			err = nil
		***REMOVED*** else if err != nil ***REMOVED***
			cc.writeStreamReset(cs.ID, ErrCodeCancel, err)
			return err
		***REMOVED***

		remain := buf[:n]
		for len(remain) > 0 && err == nil ***REMOVED***
			var allowed int32
			allowed, err = cs.awaitFlowControl(len(remain))
			switch ***REMOVED***
			case err == errStopReqBodyWrite:
				return err
			case err == errStopReqBodyWriteAndCancel:
				cc.writeStreamReset(cs.ID, ErrCodeCancel, nil)
				return err
			case err != nil:
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

	var trls []byte
	if hasTrailers ***REMOVED***
		cc.mu.Lock()
		trls, err = cc.encodeTrailers(req)
		cc.mu.Unlock()
		if err != nil ***REMOVED***
			cc.writeStreamReset(cs.ID, ErrCodeInternal, err)
			cc.forgetStreamID(cs.ID)
			return err
		***REMOVED***
	***REMOVED***

	cc.mu.Lock()
	maxFrameSize := int(cc.maxFrameSize)
	cc.mu.Unlock()

	cc.wmu.Lock()
	defer cc.wmu.Unlock()

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
	cc.mu.Lock()
	defer cc.mu.Unlock()
	for ***REMOVED***
		if cc.closed ***REMOVED***
			return 0, errClientConnClosed
		***REMOVED***
		if cs.stopReqBody != nil ***REMOVED***
			return 0, cs.stopReqBody
		***REMOVED***
		if err := cs.checkResetOrDone(); err != nil ***REMOVED***
			return 0, err
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

// requires cc.mu be held.
func (cc *ClientConn) encodeHeaders(req *http.Request, addGzipHeader bool, trailers string, contentLength int64) ([]byte, error) ***REMOVED***
	cc.hbuf.Reset()

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
			if strings.EqualFold(k, "host") || strings.EqualFold(k, "content-length") ***REMOVED***
				// Host is :authority, already sent.
				// Content-Length is automatic, set below.
				continue
			***REMOVED*** else if strings.EqualFold(k, "connection") || strings.EqualFold(k, "proxy-connection") ||
				strings.EqualFold(k, "transfer-encoding") || strings.EqualFold(k, "upgrade") ||
				strings.EqualFold(k, "keep-alive") ***REMOVED***
				// Per 8.1.2.2 Connection-Specific Header
				// Fields, don't send connection-specific
				// fields. We have already checked if any
				// are error-worthy so just ignore the rest.
				continue
			***REMOVED*** else if strings.EqualFold(k, "user-agent") ***REMOVED***
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
			***REMOVED*** else if strings.EqualFold(k, "cookie") ***REMOVED***
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
		name = strings.ToLower(name)
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

// requires cc.mu be held.
func (cc *ClientConn) encodeTrailers(req *http.Request) ([]byte, error) ***REMOVED***
	cc.hbuf.Reset()

	hlSize := uint64(0)
	for k, vv := range req.Trailer ***REMOVED***
		for _, v := range vv ***REMOVED***
			hf := hpack.HeaderField***REMOVED***Name: k, Value: v***REMOVED***
			hlSize += uint64(hf.Size())
		***REMOVED***
	***REMOVED***
	if hlSize > cc.peerMaxHeaderListSize ***REMOVED***
		return nil, errRequestHeaderListSize
	***REMOVED***

	for k, vv := range req.Trailer ***REMOVED***
		// Transfer-Encoding, etc.. have already been filtered at the
		// start of RoundTrip
		lowKey := strings.ToLower(k)
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
func (cc *ClientConn) newStream() *clientStream ***REMOVED***
	cs := &clientStream***REMOVED***
		cc:        cc,
		ID:        cc.nextStreamID,
		resc:      make(chan resAndError, 1),
		peerReset: make(chan struct***REMOVED******REMOVED***),
		done:      make(chan struct***REMOVED******REMOVED***),
	***REMOVED***
	cs.flow.add(int32(cc.initialWindowSize))
	cs.flow.setConnFlow(&cc.flow)
	cs.inflow.add(transportDefaultStreamFlow)
	cs.inflow.setConnFlow(&cc.inflow)
	cc.nextStreamID += 2
	cc.streams[cs.ID] = cs
	return cs
***REMOVED***

func (cc *ClientConn) forgetStreamID(id uint32) ***REMOVED***
	cc.streamByID(id, true)
***REMOVED***

func (cc *ClientConn) streamByID(id uint32, andRemove bool) *clientStream ***REMOVED***
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cs := cc.streams[id]
	if andRemove && cs != nil && !cc.closed ***REMOVED***
		cc.lastActive = time.Now()
		delete(cc.streams, id)
		if len(cc.streams) == 0 && cc.idleTimer != nil ***REMOVED***
			cc.idleTimer.Reset(cc.idleTimeout)
			cc.lastIdle = time.Now()
		***REMOVED***
		close(cs.done)
		// Wake up checkResetOrDone via clientStream.awaitFlowControl and
		// wake up RoundTrip if there is a pending request.
		cc.cond.Broadcast()
	***REMOVED***
	return cs
***REMOVED***

// clientConnReadLoop is the state owned by the clientConn's frame-reading readLoop.
type clientConnReadLoop struct ***REMOVED***
	_             incomparable
	cc            *ClientConn
	closeWhenIdle bool
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
	defer cc.tconn.Close()
	defer cc.t.connPool().MarkDead(cc)
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
	for _, cs := range cc.streams ***REMOVED***
		cs.bufPipe.CloseWithError(err) // no-op if already closed
		select ***REMOVED***
		case cs.resc <- resAndError***REMOVED***err: err***REMOVED***:
		default:
		***REMOVED***
		close(cs.done)
	***REMOVED***
	cc.closed = true
	cc.cond.Broadcast()
	cc.mu.Unlock()
***REMOVED***

func (rl *clientConnReadLoop) run() error ***REMOVED***
	cc := rl.cc
	rl.closeWhenIdle = cc.t.disableKeepAlives() || cc.singleUse
	gotReply := false // ever saw a HEADERS reply
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
			if cs := cc.streamByID(se.StreamID, false); cs != nil ***REMOVED***
				cs.cc.writeStreamReset(cs.ID, se.Code, err)
				cs.cc.forgetStreamID(cs.ID)
				if se.Cause == nil ***REMOVED***
					se.Cause = cc.fr.errDetail
				***REMOVED***
				rl.endStreamError(cs, se)
			***REMOVED***
			continue
		***REMOVED*** else if err != nil ***REMOVED***
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
		maybeIdle := false // whether frame might transition us to idle

		switch f := f.(type) ***REMOVED***
		case *MetaHeadersFrame:
			err = rl.processHeaders(f)
			maybeIdle = true
			gotReply = true
		case *DataFrame:
			err = rl.processData(f)
			maybeIdle = true
		case *GoAwayFrame:
			err = rl.processGoAway(f)
			maybeIdle = true
		case *RSTStreamFrame:
			err = rl.processResetStream(f)
			maybeIdle = true
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
		if rl.closeWhenIdle && gotReply && maybeIdle ***REMOVED***
			cc.closeIfIdle()
		***REMOVED***
	***REMOVED***
***REMOVED***

func (rl *clientConnReadLoop) processHeaders(f *MetaHeadersFrame) error ***REMOVED***
	cc := rl.cc
	cs := cc.streamByID(f.StreamID, false)
	if cs == nil ***REMOVED***
		// We'd get here if we canceled a request while the
		// server had its response still in flight. So if this
		// was just something we canceled, ignore it.
		return nil
	***REMOVED***
	if f.StreamEnded() ***REMOVED***
		// Issue 20521: If the stream has ended, streamByID() causes
		// clientStream.done to be closed, which causes the request's bodyWriter
		// to be closed with an errStreamClosed, which may be received by
		// clientConn.RoundTrip before the result of processing these headers.
		// Deferring stream closure allows the header processing to occur first.
		// clientConn.RoundTrip may still receive the bodyWriter error first, but
		// the fix for issue 16102 prioritises any response.
		//
		// Issue 22413: If there is no request body, we should close the
		// stream before writing to cs.resc so that the stream is closed
		// immediately once RoundTrip returns.
		if cs.req.Body != nil ***REMOVED***
			defer cc.forgetStreamID(f.StreamID)
		***REMOVED*** else ***REMOVED***
			cc.forgetStreamID(f.StreamID)
		***REMOVED***
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
		cs.cc.writeStreamReset(f.StreamID, ErrCodeProtocol, err)
		cc.forgetStreamID(cs.ID)
		cs.resc <- resAndError***REMOVED***err: err***REMOVED***
		return nil // return nil from process* funcs to keep conn alive
	***REMOVED***
	if res == nil ***REMOVED***
		// (nil, nil) special case. See handleResponse docs.
		return nil
	***REMOVED***
	cs.resTrailer = &res.Trailer
	cs.resc <- resAndError***REMOVED***res: res***REMOVED***
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
			if cs.on100 != nil ***REMOVED***
				cs.on100() // forces any write delay timer to fire
			***REMOVED***
		***REMOVED***
		cs.pastHeaders = false // do it all again
		return nil, nil
	***REMOVED***

	streamEnded := f.StreamEnded()
	isHead := cs.req.Method == "HEAD"
	if !streamEnded || isHead ***REMOVED***
		res.ContentLength = -1
		if clens := res.Header["Content-Length"]; len(clens) == 1 ***REMOVED***
			if clen64, err := strconv.ParseInt(clens[0], 10, 64); err == nil ***REMOVED***
				res.ContentLength = clen64
			***REMOVED*** else ***REMOVED***
				// TODO: care? unlike http/1, it won't mess up our framing, so it's
				// more safe smuggling-wise to ignore.
			***REMOVED***
		***REMOVED*** else if len(clens) > 1 ***REMOVED***
			// TODO: care? unlike http/1, it won't mess up our framing, so it's
			// more safe smuggling-wise to ignore.
		***REMOVED***
	***REMOVED***

	if streamEnded || isHead ***REMOVED***
		res.Body = noBody
		return res, nil
	***REMOVED***

	cs.bufPipe = pipe***REMOVED***b: &dataBuffer***REMOVED***expected: res.ContentLength***REMOVED******REMOVED***
	cs.bytesRemain = res.ContentLength
	res.Body = transportResponseBody***REMOVED***cs***REMOVED***
	go cs.awaitRequestCancel(cs.req)

	if cs.requestedGzip && res.Header.Get("Content-Encoding") == "gzip" ***REMOVED***
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
// Response.Body. It is an io.ReadCloser. On Read, it reads from cs.body.
// On Close it sends RST_STREAM if EOF wasn't already seen.
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
				cc.writeStreamReset(cs.ID, ErrCodeProtocol, err)
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
	defer cc.mu.Unlock()

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

	serverSentStreamEnd := cs.bufPipe.Err() == io.EOF
	unread := cs.bufPipe.Len()

	if unread > 0 || !serverSentStreamEnd ***REMOVED***
		cc.mu.Lock()
		cc.wmu.Lock()
		if !serverSentStreamEnd ***REMOVED***
			cc.fr.WriteRSTStream(cs.ID, ErrCodeCancel)
			cs.didReset = true
		***REMOVED***
		// Return connection-level flow control.
		if unread > 0 ***REMOVED***
			cc.inflow.add(int32(unread))
			cc.fr.WriteWindowUpdate(0, uint32(unread))
		***REMOVED***
		cc.bw.Flush()
		cc.wmu.Unlock()
		cc.mu.Unlock()
	***REMOVED***

	cs.bufPipe.BreakWithError(errClosedResponseBody)
	cc.forgetStreamID(cs.ID)
	return nil
***REMOVED***

func (rl *clientConnReadLoop) processData(f *DataFrame) error ***REMOVED***
	cc := rl.cc
	cs := cc.streamByID(f.StreamID, f.StreamEnded())
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
	if !cs.firstByte ***REMOVED***
		cc.logf("protocol error: received DATA before a HEADERS frame")
		rl.endStreamError(cs, StreamError***REMOVED***
			StreamID: f.StreamID,
			Code:     ErrCodeProtocol,
		***REMOVED***)
		return nil
	***REMOVED***
	if f.Length > 0 ***REMOVED***
		if cs.req.Method == "HEAD" && len(data) > 0 ***REMOVED***
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
		// Return len(data) now if the stream is already closed,
		// since data will never be read.
		didReset := cs.didReset
		if didReset ***REMOVED***
			refund += len(data)
		***REMOVED***
		if refund > 0 ***REMOVED***
			cc.inflow.add(int32(refund))
			cc.wmu.Lock()
			cc.fr.WriteWindowUpdate(0, uint32(refund))
			if !didReset ***REMOVED***
				cs.inflow.add(int32(refund))
				cc.fr.WriteWindowUpdate(cs.ID, uint32(refund))
			***REMOVED***
			cc.bw.Flush()
			cc.wmu.Unlock()
		***REMOVED***
		cc.mu.Unlock()

		if len(data) > 0 && !didReset ***REMOVED***
			if _, err := cs.bufPipe.Write(data); err != nil ***REMOVED***
				rl.endStreamError(cs, err)
				return err
			***REMOVED***
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
	rl.endStreamError(cs, nil)
***REMOVED***

func (rl *clientConnReadLoop) endStreamError(cs *clientStream, err error) ***REMOVED***
	var code func()
	if err == nil ***REMOVED***
		err = io.EOF
		code = cs.copyTrailers
	***REMOVED***
	if isConnectionCloseRequest(cs.req) ***REMOVED***
		rl.closeWhenIdle = true
	***REMOVED***
	cs.bufPipe.closeWithErrorAndCode(err, code)

	select ***REMOVED***
	case cs.resc <- resAndError***REMOVED***err: err***REMOVED***:
	default:
	***REMOVED***
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
	***REMOVED***
	cc.setGoAway(f)
	return nil
***REMOVED***

func (rl *clientConnReadLoop) processSettings(f *SettingsFrame) error ***REMOVED***
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

	err := f.ForeachSetting(func(s Setting) error ***REMOVED***
		switch s.ID ***REMOVED***
		case SettingMaxFrameSize:
			cc.maxFrameSize = s.Val
		case SettingMaxConcurrentStreams:
			cc.maxConcurrentStreams = s.Val
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

	cc.wmu.Lock()
	defer cc.wmu.Unlock()

	cc.fr.WriteSettingsAck()
	cc.bw.Flush()
	return cc.werr
***REMOVED***

func (rl *clientConnReadLoop) processWindowUpdate(f *WindowUpdateFrame) error ***REMOVED***
	cc := rl.cc
	cs := cc.streamByID(f.StreamID, false)
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
	cs := rl.cc.streamByID(f.StreamID, true)
	if cs == nil ***REMOVED***
		// TODO: return error if server tries to RST_STEAM an idle stream
		return nil
	***REMOVED***
	select ***REMOVED***
	case <-cs.peerReset:
		// Already reset.
		// This is the only goroutine
		// which closes this, so there
		// isn't a race.
	default:
		err := streamError(cs.ID, f.ErrCode)
		cs.resetErr = err
		close(cs.peerReset)
		cs.bufPipe.CloseWithError(err)
		cs.cc.cond.Broadcast() // wake up checkResetOrDone via clientStream.awaitFlowControl
	***REMOVED***
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
	cc.wmu.Lock()
	if err := cc.fr.WritePing(false, p); err != nil ***REMOVED***
		cc.wmu.Unlock()
		return err
	***REMOVED***
	if err := cc.bw.Flush(); err != nil ***REMOVED***
		cc.wmu.Unlock()
		return err
	***REMOVED***
	cc.wmu.Unlock()
	select ***REMOVED***
	case <-c:
		return nil
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

func strSliceContains(ss []string, s string) bool ***REMOVED***
	for _, v := range ss ***REMOVED***
		if v == s ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

type erringRoundTripper struct***REMOVED*** err error ***REMOVED***

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

// bodyWriterState encapsulates various state around the Transport's writing
// of the request body, particularly regarding doing delayed writes of the body
// when the request contains "Expect: 100-continue".
type bodyWriterState struct ***REMOVED***
	cs     *clientStream
	timer  *time.Timer   // if non-nil, we're doing a delayed write
	fnonce *sync.Once    // to call fn with
	fn     func()        // the code to run in the goroutine, writing the body
	resc   chan error    // result of fn's execution
	delay  time.Duration // how long we should delay a delayed write for
***REMOVED***

func (t *Transport) getBodyWriterState(cs *clientStream, body io.Reader) (s bodyWriterState) ***REMOVED***
	s.cs = cs
	if body == nil ***REMOVED***
		return
	***REMOVED***
	resc := make(chan error, 1)
	s.resc = resc
	s.fn = func() ***REMOVED***
		cs.cc.mu.Lock()
		cs.startedWrite = true
		cs.cc.mu.Unlock()
		resc <- cs.writeRequestBody(body, cs.req.Body)
	***REMOVED***
	s.delay = t.expectContinueTimeout()
	if s.delay == 0 ||
		!httpguts.HeaderValuesContainsToken(
			cs.req.Header["Expect"],
			"100-continue") ***REMOVED***
		return
	***REMOVED***
	s.fnonce = new(sync.Once)

	// Arm the timer with a very large duration, which we'll
	// intentionally lower later. It has to be large now because
	// we need a handle to it before writing the headers, but the
	// s.delay value is defined to not start until after the
	// request headers were written.
	const hugeDuration = 365 * 24 * time.Hour
	s.timer = time.AfterFunc(hugeDuration, func() ***REMOVED***
		s.fnonce.Do(s.fn)
	***REMOVED***)
	return
***REMOVED***

func (s bodyWriterState) cancel() ***REMOVED***
	if s.timer != nil ***REMOVED***
		s.timer.Stop()
	***REMOVED***
***REMOVED***

func (s bodyWriterState) on100() ***REMOVED***
	if s.timer == nil ***REMOVED***
		// If we didn't do a delayed write, ignore the server's
		// bogus 100 continue response.
		return
	***REMOVED***
	s.timer.Stop()
	go func() ***REMOVED*** s.fnonce.Do(s.fn) ***REMOVED***()
***REMOVED***

// scheduleBodyWrite starts writing the body, either immediately (in
// the common case) or after the delay timeout. It should not be
// called until after the headers have been written.
func (s bodyWriterState) scheduleBodyWrite() ***REMOVED***
	if s.timer == nil ***REMOVED***
		// We're not doing a delayed write (see
		// getBodyWriterState), so just start the writing
		// goroutine immediately.
		go s.fn()
		return
	***REMOVED***
	traceWait100Continue(s.cs.trace)
	if s.timer.Stop() ***REMOVED***
		s.timer.Reset(s.delay)
	***REMOVED***
***REMOVED***

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
