// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// HandshakeError describes an error with the handshake from the peer.
type HandshakeError struct ***REMOVED***
	message string
***REMOVED***

func (e HandshakeError) Error() string ***REMOVED*** return e.message ***REMOVED***

// Upgrader specifies parameters for upgrading an HTTP connection to a
// WebSocket connection.
//
// It is safe to call Upgrader's methods concurrently.
type Upgrader struct ***REMOVED***
	// HandshakeTimeout specifies the duration for the handshake to complete.
	HandshakeTimeout time.Duration

	// ReadBufferSize and WriteBufferSize specify I/O buffer sizes in bytes. If a buffer
	// size is zero, then buffers allocated by the HTTP server are used. The
	// I/O buffer sizes do not limit the size of the messages that can be sent
	// or received.
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

	// Subprotocols specifies the server's supported protocols in order of
	// preference. If this field is not nil, then the Upgrade method negotiates a
	// subprotocol by selecting the first match in this list with a protocol
	// requested by the client. If there's no match, then no protocol is
	// negotiated (the Sec-Websocket-Protocol header is not included in the
	// handshake response).
	Subprotocols []string

	// Error specifies the function for generating HTTP error responses. If Error
	// is nil, then http.Error is used to generate the HTTP response.
	Error func(w http.ResponseWriter, r *http.Request, status int, reason error)

	// CheckOrigin returns true if the request Origin header is acceptable. If
	// CheckOrigin is nil, then a safe default is used: return false if the
	// Origin request header is present and the origin host is not equal to
	// request Host header.
	//
	// A CheckOrigin function should carefully validate the request origin to
	// prevent cross-site request forgery.
	CheckOrigin func(r *http.Request) bool

	// EnableCompression specify if the server should attempt to negotiate per
	// message compression (RFC 7692). Setting this value to true does not
	// guarantee that compression will be supported. Currently only "no context
	// takeover" modes are supported.
	EnableCompression bool
***REMOVED***

func (u *Upgrader) returnError(w http.ResponseWriter, r *http.Request, status int, reason string) (*Conn, error) ***REMOVED***
	err := HandshakeError***REMOVED***reason***REMOVED***
	if u.Error != nil ***REMOVED***
		u.Error(w, r, status, err)
	***REMOVED*** else ***REMOVED***
		w.Header().Set("Sec-Websocket-Version", "13")
		http.Error(w, http.StatusText(status), status)
	***REMOVED***
	return nil, err
***REMOVED***

// checkSameOrigin returns true if the origin is not set or is equal to the request host.
func checkSameOrigin(r *http.Request) bool ***REMOVED***
	origin := r.Header["Origin"]
	if len(origin) == 0 ***REMOVED***
		return true
	***REMOVED***
	u, err := url.Parse(origin[0])
	if err != nil ***REMOVED***
		return false
	***REMOVED***
	return equalASCIIFold(u.Host, r.Host)
***REMOVED***

func (u *Upgrader) selectSubprotocol(r *http.Request, responseHeader http.Header) string ***REMOVED***
	if u.Subprotocols != nil ***REMOVED***
		clientProtocols := Subprotocols(r)
		for _, serverProtocol := range u.Subprotocols ***REMOVED***
			for _, clientProtocol := range clientProtocols ***REMOVED***
				if clientProtocol == serverProtocol ***REMOVED***
					return clientProtocol
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED*** else if responseHeader != nil ***REMOVED***
		return responseHeader.Get("Sec-Websocket-Protocol")
	***REMOVED***
	return ""
***REMOVED***

// Upgrade upgrades the HTTP server connection to the WebSocket protocol.
//
// The responseHeader is included in the response to the client's upgrade
// request. Use the responseHeader to specify cookies (Set-Cookie). To specify
// subprotocols supported by the server, set Upgrader.Subprotocols directly.
//
// If the upgrade fails, then Upgrade replies to the client with an HTTP error
// response.
func (u *Upgrader) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*Conn, error) ***REMOVED***
	const badHandshake = "websocket: the client is not using the websocket protocol: "

	if !tokenListContainsValue(r.Header, "Connection", "upgrade") ***REMOVED***
		return u.returnError(w, r, http.StatusBadRequest, badHandshake+"'upgrade' token not found in 'Connection' header")
	***REMOVED***

	if !tokenListContainsValue(r.Header, "Upgrade", "websocket") ***REMOVED***
		return u.returnError(w, r, http.StatusBadRequest, badHandshake+"'websocket' token not found in 'Upgrade' header")
	***REMOVED***

	if r.Method != http.MethodGet ***REMOVED***
		return u.returnError(w, r, http.StatusMethodNotAllowed, badHandshake+"request method is not GET")
	***REMOVED***

	if !tokenListContainsValue(r.Header, "Sec-Websocket-Version", "13") ***REMOVED***
		return u.returnError(w, r, http.StatusBadRequest, "websocket: unsupported version: 13 not found in 'Sec-Websocket-Version' header")
	***REMOVED***

	if _, ok := responseHeader["Sec-Websocket-Extensions"]; ok ***REMOVED***
		return u.returnError(w, r, http.StatusInternalServerError, "websocket: application specific 'Sec-WebSocket-Extensions' headers are unsupported")
	***REMOVED***

	checkOrigin := u.CheckOrigin
	if checkOrigin == nil ***REMOVED***
		checkOrigin = checkSameOrigin
	***REMOVED***
	if !checkOrigin(r) ***REMOVED***
		return u.returnError(w, r, http.StatusForbidden, "websocket: request origin not allowed by Upgrader.CheckOrigin")
	***REMOVED***

	challengeKey := r.Header.Get("Sec-Websocket-Key")
	if challengeKey == "" ***REMOVED***
		return u.returnError(w, r, http.StatusBadRequest, "websocket: not a websocket handshake: 'Sec-WebSocket-Key' header is missing or blank")
	***REMOVED***

	subprotocol := u.selectSubprotocol(r, responseHeader)

	// Negotiate PMCE
	var compress bool
	if u.EnableCompression ***REMOVED***
		for _, ext := range parseExtensions(r.Header) ***REMOVED***
			if ext[""] != "permessage-deflate" ***REMOVED***
				continue
			***REMOVED***
			compress = true
			break
		***REMOVED***
	***REMOVED***

	h, ok := w.(http.Hijacker)
	if !ok ***REMOVED***
		return u.returnError(w, r, http.StatusInternalServerError, "websocket: response does not implement http.Hijacker")
	***REMOVED***
	var brw *bufio.ReadWriter
	netConn, brw, err := h.Hijack()
	if err != nil ***REMOVED***
		return u.returnError(w, r, http.StatusInternalServerError, err.Error())
	***REMOVED***

	if brw.Reader.Buffered() > 0 ***REMOVED***
		netConn.Close()
		return nil, errors.New("websocket: client sent data before handshake is complete")
	***REMOVED***

	var br *bufio.Reader
	if u.ReadBufferSize == 0 && bufioReaderSize(netConn, brw.Reader) > 256 ***REMOVED***
		// Reuse hijacked buffered reader as connection reader.
		br = brw.Reader
	***REMOVED***

	buf := bufioWriterBuffer(netConn, brw.Writer)

	var writeBuf []byte
	if u.WriteBufferPool == nil && u.WriteBufferSize == 0 && len(buf) >= maxFrameHeaderSize+256 ***REMOVED***
		// Reuse hijacked write buffer as connection buffer.
		writeBuf = buf
	***REMOVED***

	c := newConn(netConn, true, u.ReadBufferSize, u.WriteBufferSize, u.WriteBufferPool, br, writeBuf)
	c.subprotocol = subprotocol

	if compress ***REMOVED***
		c.newCompressionWriter = compressNoContextTakeover
		c.newDecompressionReader = decompressNoContextTakeover
	***REMOVED***

	// Use larger of hijacked buffer and connection write buffer for header.
	p := buf
	if len(c.writeBuf) > len(p) ***REMOVED***
		p = c.writeBuf
	***REMOVED***
	p = p[:0]

	p = append(p, "HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: "...)
	p = append(p, computeAcceptKey(challengeKey)...)
	p = append(p, "\r\n"...)
	if c.subprotocol != "" ***REMOVED***
		p = append(p, "Sec-WebSocket-Protocol: "...)
		p = append(p, c.subprotocol...)
		p = append(p, "\r\n"...)
	***REMOVED***
	if compress ***REMOVED***
		p = append(p, "Sec-WebSocket-Extensions: permessage-deflate; server_no_context_takeover; client_no_context_takeover\r\n"...)
	***REMOVED***
	for k, vs := range responseHeader ***REMOVED***
		if k == "Sec-Websocket-Protocol" ***REMOVED***
			continue
		***REMOVED***
		for _, v := range vs ***REMOVED***
			p = append(p, k...)
			p = append(p, ": "...)
			for i := 0; i < len(v); i++ ***REMOVED***
				b := v[i]
				if b <= 31 ***REMOVED***
					// prevent response splitting.
					b = ' '
				***REMOVED***
				p = append(p, b)
			***REMOVED***
			p = append(p, "\r\n"...)
		***REMOVED***
	***REMOVED***
	p = append(p, "\r\n"...)

	// Clear deadlines set by HTTP server.
	netConn.SetDeadline(time.Time***REMOVED******REMOVED***)

	if u.HandshakeTimeout > 0 ***REMOVED***
		netConn.SetWriteDeadline(time.Now().Add(u.HandshakeTimeout))
	***REMOVED***
	if _, err = netConn.Write(p); err != nil ***REMOVED***
		netConn.Close()
		return nil, err
	***REMOVED***
	if u.HandshakeTimeout > 0 ***REMOVED***
		netConn.SetWriteDeadline(time.Time***REMOVED******REMOVED***)
	***REMOVED***

	return c, nil
***REMOVED***

// Upgrade upgrades the HTTP server connection to the WebSocket protocol.
//
// Deprecated: Use websocket.Upgrader instead.
//
// Upgrade does not perform origin checking. The application is responsible for
// checking the Origin header before calling Upgrade. An example implementation
// of the same origin policy check is:
//
//	if req.Header.Get("Origin") != "http://"+req.Host ***REMOVED***
//		http.Error(w, "Origin not allowed", http.StatusForbidden)
//		return
//	***REMOVED***
//
// If the endpoint supports subprotocols, then the application is responsible
// for negotiating the protocol used on the connection. Use the Subprotocols()
// function to get the subprotocols requested by the client. Use the
// Sec-Websocket-Protocol response header to specify the subprotocol selected
// by the application.
//
// The responseHeader is included in the response to the client's upgrade
// request. Use the responseHeader to specify cookies (Set-Cookie) and the
// negotiated subprotocol (Sec-Websocket-Protocol).
//
// The connection buffers IO to the underlying network connection. The
// readBufSize and writeBufSize parameters specify the size of the buffers to
// use. Messages can be larger than the buffers.
//
// If the request is not a valid WebSocket handshake, then Upgrade returns an
// error of type HandshakeError. Applications should handle this error by
// replying to the client with an HTTP error response.
func Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header, readBufSize, writeBufSize int) (*Conn, error) ***REMOVED***
	u := Upgrader***REMOVED***ReadBufferSize: readBufSize, WriteBufferSize: writeBufSize***REMOVED***
	u.Error = func(w http.ResponseWriter, r *http.Request, status int, reason error) ***REMOVED***
		// don't return errors to maintain backwards compatibility
	***REMOVED***
	u.CheckOrigin = func(r *http.Request) bool ***REMOVED***
		// allow all connections by default
		return true
	***REMOVED***
	return u.Upgrade(w, r, responseHeader)
***REMOVED***

// Subprotocols returns the subprotocols requested by the client in the
// Sec-Websocket-Protocol header.
func Subprotocols(r *http.Request) []string ***REMOVED***
	h := strings.TrimSpace(r.Header.Get("Sec-Websocket-Protocol"))
	if h == "" ***REMOVED***
		return nil
	***REMOVED***
	protocols := strings.Split(h, ",")
	for i := range protocols ***REMOVED***
		protocols[i] = strings.TrimSpace(protocols[i])
	***REMOVED***
	return protocols
***REMOVED***

// IsWebSocketUpgrade returns true if the client requested upgrade to the
// WebSocket protocol.
func IsWebSocketUpgrade(r *http.Request) bool ***REMOVED***
	return tokenListContainsValue(r.Header, "Connection", "upgrade") &&
		tokenListContainsValue(r.Header, "Upgrade", "websocket")
***REMOVED***

// bufioReaderSize size returns the size of a bufio.Reader.
func bufioReaderSize(originalReader io.Reader, br *bufio.Reader) int ***REMOVED***
	// This code assumes that peek on a reset reader returns
	// bufio.Reader.buf[:0].
	// TODO: Use bufio.Reader.Size() after Go 1.10
	br.Reset(originalReader)
	if p, err := br.Peek(0); err == nil ***REMOVED***
		return cap(p)
	***REMOVED***
	return 0
***REMOVED***

// writeHook is an io.Writer that records the last slice passed to it vio
// io.Writer.Write.
type writeHook struct ***REMOVED***
	p []byte
***REMOVED***

func (wh *writeHook) Write(p []byte) (int, error) ***REMOVED***
	wh.p = p
	return len(p), nil
***REMOVED***

// bufioWriterBuffer grabs the buffer from a bufio.Writer.
func bufioWriterBuffer(originalWriter io.Writer, bw *bufio.Writer) []byte ***REMOVED***
	// This code assumes that bufio.Writer.buf[:1] is passed to the
	// bufio.Writer's underlying writer.
	var wh writeHook
	bw.Reset(&wh)
	bw.WriteByte(0)
	bw.Flush()

	bw.Reset(originalWriter)

	return wh.p[:cap(wh.p)]
***REMOVED***
