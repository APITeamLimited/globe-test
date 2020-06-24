// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strings"
	"time"
)

// ErrBadHandshake is returned when the server response to opening handshake is
// invalid.
var ErrBadHandshake = errors.New("websocket: bad handshake")

var errInvalidCompression = errors.New("websocket: invalid compression negotiation")

// NewClient creates a new client connection using the given net connection.
// The URL u specifies the host and request URI. Use requestHeader to specify
// the origin (Origin), subprotocols (Sec-WebSocket-Protocol) and cookies
// (Cookie). Use the response.Header to get the selected subprotocol
// (Sec-WebSocket-Protocol) and cookies (Set-Cookie).
//
// If the WebSocket handshake fails, ErrBadHandshake is returned along with a
// non-nil *http.Response so that callers can handle redirects, authentication,
// etc.
//
// Deprecated: Use Dialer instead.
func NewClient(netConn net.Conn, u *url.URL, requestHeader http.Header, readBufSize, writeBufSize int) (c *Conn, response *http.Response, err error) ***REMOVED***
	d := Dialer***REMOVED***
		ReadBufferSize:  readBufSize,
		WriteBufferSize: writeBufSize,
		NetDial: func(net, addr string) (net.Conn, error) ***REMOVED***
			return netConn, nil
		***REMOVED***,
	***REMOVED***
	return d.Dial(u.String(), requestHeader)
***REMOVED***

// A Dialer contains options for connecting to WebSocket server.
type Dialer struct ***REMOVED***
	// NetDial specifies the dial function for creating TCP connections. If
	// NetDial is nil, net.Dial is used.
	NetDial func(network, addr string) (net.Conn, error)

	// NetDialContext specifies the dial function for creating TCP connections. If
	// NetDialContext is nil, net.DialContext is used.
	NetDialContext func(ctx context.Context, network, addr string) (net.Conn, error)

	// Proxy specifies a function to return a proxy for a given
	// Request. If the function returns a non-nil error, the
	// request is aborted with the provided error.
	// If Proxy is nil or returns a nil *URL, no proxy is used.
	Proxy func(*http.Request) (*url.URL, error)

	// TLSClientConfig specifies the TLS configuration to use with tls.Client.
	// If nil, the default configuration is used.
	TLSClientConfig *tls.Config

	// HandshakeTimeout specifies the duration for the handshake to complete.
	HandshakeTimeout time.Duration

	// ReadBufferSize and WriteBufferSize specify I/O buffer sizes in bytes. If a buffer
	// size is zero, then a useful default size is used. The I/O buffer sizes
	// do not limit the size of the messages that can be sent or received.
	ReadBufferSize, WriteBufferSize int

	// WriteBufferPool is a pool of buffers for write operations. If the value
	// is not set, then write buffers are allocated to the connection for the
	// lifetime of the connection.
	//
	// A pool is most useful when the application has a modest volume of writes
	// across a large number of connections.
	//
	// Applications should use a single pool for each unique value of
	// WriteBufferSize.
	WriteBufferPool BufferPool

	// Subprotocols specifies the client's requested subprotocols.
	Subprotocols []string

	// EnableCompression specifies if the client should attempt to negotiate
	// per message compression (RFC 7692). Setting this value to true does not
	// guarantee that compression will be supported. Currently only "no context
	// takeover" modes are supported.
	EnableCompression bool

	// Jar specifies the cookie jar.
	// If Jar is nil, cookies are not sent in requests and ignored
	// in responses.
	Jar http.CookieJar
***REMOVED***

// Dial creates a new client connection by calling DialContext with a background context.
func (d *Dialer) Dial(urlStr string, requestHeader http.Header) (*Conn, *http.Response, error) ***REMOVED***
	return d.DialContext(context.Background(), urlStr, requestHeader)
***REMOVED***

var errMalformedURL = errors.New("malformed ws or wss URL")

func hostPortNoPort(u *url.URL) (hostPort, hostNoPort string) ***REMOVED***
	hostPort = u.Host
	hostNoPort = u.Host
	if i := strings.LastIndex(u.Host, ":"); i > strings.LastIndex(u.Host, "]") ***REMOVED***
		hostNoPort = hostNoPort[:i]
	***REMOVED*** else ***REMOVED***
		switch u.Scheme ***REMOVED***
		case "wss":
			hostPort += ":443"
		case "https":
			hostPort += ":443"
		default:
			hostPort += ":80"
		***REMOVED***
	***REMOVED***
	return hostPort, hostNoPort
***REMOVED***

// DefaultDialer is a dialer with all fields set to the default values.
var DefaultDialer = &Dialer***REMOVED***
	Proxy:            http.ProxyFromEnvironment,
	HandshakeTimeout: 45 * time.Second,
***REMOVED***

// nilDialer is dialer to use when receiver is nil.
var nilDialer = *DefaultDialer

// DialContext creates a new client connection. Use requestHeader to specify the
// origin (Origin), subprotocols (Sec-WebSocket-Protocol) and cookies (Cookie).
// Use the response.Header to get the selected subprotocol
// (Sec-WebSocket-Protocol) and cookies (Set-Cookie).
//
// The context will be used in the request and in the Dialer.
//
// If the WebSocket handshake fails, ErrBadHandshake is returned along with a
// non-nil *http.Response so that callers can handle redirects, authentication,
// etcetera. The response body may not contain the entire response and does not
// need to be closed by the application.
func (d *Dialer) DialContext(ctx context.Context, urlStr string, requestHeader http.Header) (*Conn, *http.Response, error) ***REMOVED***
	if d == nil ***REMOVED***
		d = &nilDialer
	***REMOVED***

	challengeKey, err := generateChallengeKey()
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	u, err := url.Parse(urlStr)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	switch u.Scheme ***REMOVED***
	case "ws":
		u.Scheme = "http"
	case "wss":
		u.Scheme = "https"
	default:
		return nil, nil, errMalformedURL
	***REMOVED***

	if u.User != nil ***REMOVED***
		// User name and password are not allowed in websocket URIs.
		return nil, nil, errMalformedURL
	***REMOVED***

	req := &http.Request***REMOVED***
		Method:     "GET",
		URL:        u,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Host:       u.Host,
	***REMOVED***
	req = req.WithContext(ctx)

	// Set the cookies present in the cookie jar of the dialer
	if d.Jar != nil ***REMOVED***
		for _, cookie := range d.Jar.Cookies(u) ***REMOVED***
			req.AddCookie(cookie)
		***REMOVED***
	***REMOVED***

	// Set the request headers using the capitalization for names and values in
	// RFC examples. Although the capitalization shouldn't matter, there are
	// servers that depend on it. The Header.Set method is not used because the
	// method canonicalizes the header names.
	req.Header["Upgrade"] = []string***REMOVED***"websocket"***REMOVED***
	req.Header["Connection"] = []string***REMOVED***"Upgrade"***REMOVED***
	req.Header["Sec-WebSocket-Key"] = []string***REMOVED***challengeKey***REMOVED***
	req.Header["Sec-WebSocket-Version"] = []string***REMOVED***"13"***REMOVED***
	if len(d.Subprotocols) > 0 ***REMOVED***
		req.Header["Sec-WebSocket-Protocol"] = []string***REMOVED***strings.Join(d.Subprotocols, ", ")***REMOVED***
	***REMOVED***
	for k, vs := range requestHeader ***REMOVED***
		switch ***REMOVED***
		case k == "Host":
			if len(vs) > 0 ***REMOVED***
				req.Host = vs[0]
			***REMOVED***
		case k == "Upgrade" ||
			k == "Connection" ||
			k == "Sec-Websocket-Key" ||
			k == "Sec-Websocket-Version" ||
			k == "Sec-Websocket-Extensions" ||
			(k == "Sec-Websocket-Protocol" && len(d.Subprotocols) > 0):
			return nil, nil, errors.New("websocket: duplicate header not allowed: " + k)
		case k == "Sec-Websocket-Protocol":
			req.Header["Sec-WebSocket-Protocol"] = vs
		default:
			req.Header[k] = vs
		***REMOVED***
	***REMOVED***

	if d.EnableCompression ***REMOVED***
		req.Header["Sec-WebSocket-Extensions"] = []string***REMOVED***"permessage-deflate; server_no_context_takeover; client_no_context_takeover"***REMOVED***
	***REMOVED***

	if d.HandshakeTimeout != 0 ***REMOVED***
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, d.HandshakeTimeout)
		defer cancel()
	***REMOVED***

	// Get network dial function.
	var netDial func(network, add string) (net.Conn, error)

	if d.NetDialContext != nil ***REMOVED***
		netDial = func(network, addr string) (net.Conn, error) ***REMOVED***
			return d.NetDialContext(ctx, network, addr)
		***REMOVED***
	***REMOVED*** else if d.NetDial != nil ***REMOVED***
		netDial = d.NetDial
	***REMOVED*** else ***REMOVED***
		netDialer := &net.Dialer***REMOVED******REMOVED***
		netDial = func(network, addr string) (net.Conn, error) ***REMOVED***
			return netDialer.DialContext(ctx, network, addr)
		***REMOVED***
	***REMOVED***

	// If needed, wrap the dial function to set the connection deadline.
	if deadline, ok := ctx.Deadline(); ok ***REMOVED***
		forwardDial := netDial
		netDial = func(network, addr string) (net.Conn, error) ***REMOVED***
			c, err := forwardDial(network, addr)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			err = c.SetDeadline(deadline)
			if err != nil ***REMOVED***
				c.Close()
				return nil, err
			***REMOVED***
			return c, nil
		***REMOVED***
	***REMOVED***

	// If needed, wrap the dial function to connect through a proxy.
	if d.Proxy != nil ***REMOVED***
		proxyURL, err := d.Proxy(req)
		if err != nil ***REMOVED***
			return nil, nil, err
		***REMOVED***
		if proxyURL != nil ***REMOVED***
			dialer, err := proxy_FromURL(proxyURL, netDialerFunc(netDial))
			if err != nil ***REMOVED***
				return nil, nil, err
			***REMOVED***
			netDial = dialer.Dial
		***REMOVED***
	***REMOVED***

	hostPort, hostNoPort := hostPortNoPort(u)
	trace := httptrace.ContextClientTrace(ctx)
	if trace != nil && trace.GetConn != nil ***REMOVED***
		trace.GetConn(hostPort)
	***REMOVED***

	netConn, err := netDial("tcp", hostPort)
	if trace != nil && trace.GotConn != nil ***REMOVED***
		trace.GotConn(httptrace.GotConnInfo***REMOVED***
			Conn: netConn,
		***REMOVED***)
	***REMOVED***
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	defer func() ***REMOVED***
		if netConn != nil ***REMOVED***
			netConn.Close()
		***REMOVED***
	***REMOVED***()

	if u.Scheme == "https" ***REMOVED***
		cfg := cloneTLSConfig(d.TLSClientConfig)
		if cfg.ServerName == "" ***REMOVED***
			cfg.ServerName = hostNoPort
		***REMOVED***
		tlsConn := tls.Client(netConn, cfg)
		netConn = tlsConn

		var err error
		if trace != nil ***REMOVED***
			err = doHandshakeWithTrace(trace, tlsConn, cfg)
		***REMOVED*** else ***REMOVED***
			err = doHandshake(tlsConn, cfg)
		***REMOVED***

		if err != nil ***REMOVED***
			return nil, nil, err
		***REMOVED***
	***REMOVED***

	conn := newConn(netConn, false, d.ReadBufferSize, d.WriteBufferSize, d.WriteBufferPool, nil, nil)

	if err := req.Write(netConn); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	if trace != nil && trace.GotFirstResponseByte != nil ***REMOVED***
		if peek, err := conn.br.Peek(1); err == nil && len(peek) == 1 ***REMOVED***
			trace.GotFirstResponseByte()
		***REMOVED***
	***REMOVED***

	resp, err := http.ReadResponse(conn.br, req)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	if d.Jar != nil ***REMOVED***
		if rc := resp.Cookies(); len(rc) > 0 ***REMOVED***
			d.Jar.SetCookies(u, rc)
		***REMOVED***
	***REMOVED***

	if resp.StatusCode != 101 ||
		!strings.EqualFold(resp.Header.Get("Upgrade"), "websocket") ||
		!strings.EqualFold(resp.Header.Get("Connection"), "upgrade") ||
		resp.Header.Get("Sec-Websocket-Accept") != computeAcceptKey(challengeKey) ***REMOVED***
		// Before closing the network connection on return from this
		// function, slurp up some of the response to aid application
		// debugging.
		buf := make([]byte, 1024)
		n, _ := io.ReadFull(resp.Body, buf)
		resp.Body = ioutil.NopCloser(bytes.NewReader(buf[:n]))
		return nil, resp, ErrBadHandshake
	***REMOVED***

	for _, ext := range parseExtensions(resp.Header) ***REMOVED***
		if ext[""] != "permessage-deflate" ***REMOVED***
			continue
		***REMOVED***
		_, snct := ext["server_no_context_takeover"]
		_, cnct := ext["client_no_context_takeover"]
		if !snct || !cnct ***REMOVED***
			return nil, resp, errInvalidCompression
		***REMOVED***
		conn.newCompressionWriter = compressNoContextTakeover
		conn.newDecompressionReader = decompressNoContextTakeover
		break
	***REMOVED***

	resp.Body = ioutil.NopCloser(bytes.NewReader([]byte***REMOVED******REMOVED***))
	conn.subprotocol = resp.Header.Get("Sec-Websocket-Protocol")

	netConn.SetDeadline(time.Time***REMOVED******REMOVED***)
	netConn = nil // to avoid close in defer.
	return conn, resp, nil
***REMOVED***

func doHandshake(tlsConn *tls.Conn, cfg *tls.Config) error ***REMOVED***
	if err := tlsConn.Handshake(); err != nil ***REMOVED***
		return err
	***REMOVED***
	if !cfg.InsecureSkipVerify ***REMOVED***
		if err := tlsConn.VerifyHostname(cfg.ServerName); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
