// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
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

	// ReadBufferSize and WriteBufferSize specify I/O buffer sizes. If a buffer
	// size is zero, then a useful default size is used. The I/O buffer sizes
	// do not limit the size of the messages that can be sent or received.
	ReadBufferSize, WriteBufferSize int

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

// DefaultDialer is a dialer with all fields set to the default zero values.
var DefaultDialer = &Dialer***REMOVED***
	Proxy: http.ProxyFromEnvironment,
***REMOVED***

// Dial creates a new client connection. Use requestHeader to specify the
// origin (Origin), subprotocols (Sec-WebSocket-Protocol) and cookies (Cookie).
// Use the response.Header to get the selected subprotocol
// (Sec-WebSocket-Protocol) and cookies (Set-Cookie).
//
// If the WebSocket handshake fails, ErrBadHandshake is returned along with a
// non-nil *http.Response so that callers can handle redirects, authentication,
// etcetera. The response body may not contain the entire response and does not
// need to be closed by the application.
func (d *Dialer) Dial(urlStr string, requestHeader http.Header) (*Conn, *http.Response, error) ***REMOVED***

	if d == nil ***REMOVED***
		d = &Dialer***REMOVED***
			Proxy: http.ProxyFromEnvironment,
		***REMOVED***
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
		default:
			req.Header[k] = vs
		***REMOVED***
	***REMOVED***

	if d.EnableCompression ***REMOVED***
		req.Header.Set("Sec-Websocket-Extensions", "permessage-deflate; server_no_context_takeover; client_no_context_takeover")
	***REMOVED***

	hostPort, hostNoPort := hostPortNoPort(u)

	var proxyURL *url.URL
	// Check wether the proxy method has been configured
	if d.Proxy != nil ***REMOVED***
		proxyURL, err = d.Proxy(req)
	***REMOVED***
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	var targetHostPort string
	if proxyURL != nil ***REMOVED***
		targetHostPort, _ = hostPortNoPort(proxyURL)
	***REMOVED*** else ***REMOVED***
		targetHostPort = hostPort
	***REMOVED***

	var deadline time.Time
	if d.HandshakeTimeout != 0 ***REMOVED***
		deadline = time.Now().Add(d.HandshakeTimeout)
	***REMOVED***

	netDial := d.NetDial
	if netDial == nil ***REMOVED***
		netDialer := &net.Dialer***REMOVED***Deadline: deadline***REMOVED***
		netDial = netDialer.Dial
	***REMOVED***

	netConn, err := netDial("tcp", targetHostPort)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	defer func() ***REMOVED***
		if netConn != nil ***REMOVED***
			netConn.Close()
		***REMOVED***
	***REMOVED***()

	if err := netConn.SetDeadline(deadline); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	if proxyURL != nil ***REMOVED***
		connectHeader := make(http.Header)
		if user := proxyURL.User; user != nil ***REMOVED***
			proxyUser := user.Username()
			if proxyPassword, passwordSet := user.Password(); passwordSet ***REMOVED***
				credential := base64.StdEncoding.EncodeToString([]byte(proxyUser + ":" + proxyPassword))
				connectHeader.Set("Proxy-Authorization", "Basic "+credential)
			***REMOVED***
		***REMOVED***
		connectReq := &http.Request***REMOVED***
			Method: "CONNECT",
			URL:    &url.URL***REMOVED***Opaque: hostPort***REMOVED***,
			Host:   hostPort,
			Header: connectHeader,
		***REMOVED***

		connectReq.Write(netConn)

		// Read response.
		// Okay to use and discard buffered reader here, because
		// TLS server will not speak until spoken to.
		br := bufio.NewReader(netConn)
		resp, err := http.ReadResponse(br, connectReq)
		if err != nil ***REMOVED***
			return nil, nil, err
		***REMOVED***
		if resp.StatusCode != 200 ***REMOVED***
			f := strings.SplitN(resp.Status, " ", 2)
			return nil, nil, errors.New(f[1])
		***REMOVED***
	***REMOVED***

	if u.Scheme == "https" ***REMOVED***
		cfg := cloneTLSConfig(d.TLSClientConfig)
		if cfg.ServerName == "" ***REMOVED***
			cfg.ServerName = hostNoPort
		***REMOVED***
		tlsConn := tls.Client(netConn, cfg)
		netConn = tlsConn
		if err := tlsConn.Handshake(); err != nil ***REMOVED***
			return nil, nil, err
		***REMOVED***
		if !cfg.InsecureSkipVerify ***REMOVED***
			if err := tlsConn.VerifyHostname(cfg.ServerName); err != nil ***REMOVED***
				return nil, nil, err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	conn := newConn(netConn, false, d.ReadBufferSize, d.WriteBufferSize)

	if err := req.Write(netConn); err != nil ***REMOVED***
		return nil, nil, err
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
