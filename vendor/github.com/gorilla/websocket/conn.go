// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"
	"unicode/utf8"
)

const (
	// Frame header byte 0 bits from Section 5.2 of RFC 6455
	finalBit = 1 << 7
	rsv1Bit  = 1 << 6
	rsv2Bit  = 1 << 5
	rsv3Bit  = 1 << 4

	// Frame header byte 1 bits from Section 5.2 of RFC 6455
	maskBit = 1 << 7

	maxFrameHeaderSize         = 2 + 8 + 4 // Fixed header + length + mask
	maxControlFramePayloadSize = 125

	writeWait = time.Second

	defaultReadBufferSize  = 4096
	defaultWriteBufferSize = 4096

	continuationFrame = 0
	noFrame           = -1
)

// Close codes defined in RFC 6455, section 11.7.
const (
	CloseNormalClosure           = 1000
	CloseGoingAway               = 1001
	CloseProtocolError           = 1002
	CloseUnsupportedData         = 1003
	CloseNoStatusReceived        = 1005
	CloseAbnormalClosure         = 1006
	CloseInvalidFramePayloadData = 1007
	ClosePolicyViolation         = 1008
	CloseMessageTooBig           = 1009
	CloseMandatoryExtension      = 1010
	CloseInternalServerErr       = 1011
	CloseServiceRestart          = 1012
	CloseTryAgainLater           = 1013
	CloseTLSHandshake            = 1015
)

// The message types are defined in RFC 6455, section 11.8.
const (
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	TextMessage = 1

	// BinaryMessage denotes a binary data message.
	BinaryMessage = 2

	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	CloseMessage = 8

	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PingMessage = 9

	// PongMessage denotes a pong control message. The optional message payload
	// is UTF-8 encoded text.
	PongMessage = 10
)

// ErrCloseSent is returned when the application writes a message to the
// connection after sending a close message.
var ErrCloseSent = errors.New("websocket: close sent")

// ErrReadLimit is returned when reading a message that is larger than the
// read limit set for the connection.
var ErrReadLimit = errors.New("websocket: read limit exceeded")

// netError satisfies the net Error interface.
type netError struct ***REMOVED***
	msg       string
	temporary bool
	timeout   bool
***REMOVED***

func (e *netError) Error() string   ***REMOVED*** return e.msg ***REMOVED***
func (e *netError) Temporary() bool ***REMOVED*** return e.temporary ***REMOVED***
func (e *netError) Timeout() bool   ***REMOVED*** return e.timeout ***REMOVED***

// CloseError represents a close message.
type CloseError struct ***REMOVED***
	// Code is defined in RFC 6455, section 11.7.
	Code int

	// Text is the optional text payload.
	Text string
***REMOVED***

func (e *CloseError) Error() string ***REMOVED***
	s := []byte("websocket: close ")
	s = strconv.AppendInt(s, int64(e.Code), 10)
	switch e.Code ***REMOVED***
	case CloseNormalClosure:
		s = append(s, " (normal)"...)
	case CloseGoingAway:
		s = append(s, " (going away)"...)
	case CloseProtocolError:
		s = append(s, " (protocol error)"...)
	case CloseUnsupportedData:
		s = append(s, " (unsupported data)"...)
	case CloseNoStatusReceived:
		s = append(s, " (no status)"...)
	case CloseAbnormalClosure:
		s = append(s, " (abnormal closure)"...)
	case CloseInvalidFramePayloadData:
		s = append(s, " (invalid payload data)"...)
	case ClosePolicyViolation:
		s = append(s, " (policy violation)"...)
	case CloseMessageTooBig:
		s = append(s, " (message too big)"...)
	case CloseMandatoryExtension:
		s = append(s, " (mandatory extension missing)"...)
	case CloseInternalServerErr:
		s = append(s, " (internal server error)"...)
	case CloseTLSHandshake:
		s = append(s, " (TLS handshake error)"...)
	***REMOVED***
	if e.Text != "" ***REMOVED***
		s = append(s, ": "...)
		s = append(s, e.Text...)
	***REMOVED***
	return string(s)
***REMOVED***

// IsCloseError returns boolean indicating whether the error is a *CloseError
// with one of the specified codes.
func IsCloseError(err error, codes ...int) bool ***REMOVED***
	if e, ok := err.(*CloseError); ok ***REMOVED***
		for _, code := range codes ***REMOVED***
			if e.Code == code ***REMOVED***
				return true
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// IsUnexpectedCloseError returns boolean indicating whether the error is a
// *CloseError with a code not in the list of expected codes.
func IsUnexpectedCloseError(err error, expectedCodes ...int) bool ***REMOVED***
	if e, ok := err.(*CloseError); ok ***REMOVED***
		for _, code := range expectedCodes ***REMOVED***
			if e.Code == code ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

var (
	errWriteTimeout        = &netError***REMOVED***msg: "websocket: write timeout", timeout: true, temporary: true***REMOVED***
	errUnexpectedEOF       = &CloseError***REMOVED***Code: CloseAbnormalClosure, Text: io.ErrUnexpectedEOF.Error()***REMOVED***
	errBadWriteOpCode      = errors.New("websocket: bad write message type")
	errWriteClosed         = errors.New("websocket: write closed")
	errInvalidControlFrame = errors.New("websocket: invalid control frame")
)

func newMaskKey() [4]byte ***REMOVED***
	n := rand.Uint32()
	return [4]byte***REMOVED***byte(n), byte(n >> 8), byte(n >> 16), byte(n >> 24)***REMOVED***
***REMOVED***

func hideTempErr(err error) error ***REMOVED***
	if e, ok := err.(net.Error); ok && e.Temporary() ***REMOVED***
		err = &netError***REMOVED***msg: e.Error(), timeout: e.Timeout()***REMOVED***
	***REMOVED***
	return err
***REMOVED***

func isControl(frameType int) bool ***REMOVED***
	return frameType == CloseMessage || frameType == PingMessage || frameType == PongMessage
***REMOVED***

func isData(frameType int) bool ***REMOVED***
	return frameType == TextMessage || frameType == BinaryMessage
***REMOVED***

var validReceivedCloseCodes = map[int]bool***REMOVED***
	// see http://www.iana.org/assignments/websocket/websocket.xhtml#close-code-number

	CloseNormalClosure:           true,
	CloseGoingAway:               true,
	CloseProtocolError:           true,
	CloseUnsupportedData:         true,
	CloseNoStatusReceived:        false,
	CloseAbnormalClosure:         false,
	CloseInvalidFramePayloadData: true,
	ClosePolicyViolation:         true,
	CloseMessageTooBig:           true,
	CloseMandatoryExtension:      true,
	CloseInternalServerErr:       true,
	CloseServiceRestart:          true,
	CloseTryAgainLater:           true,
	CloseTLSHandshake:            false,
***REMOVED***

func isValidReceivedCloseCode(code int) bool ***REMOVED***
	return validReceivedCloseCodes[code] || (code >= 3000 && code <= 4999)
***REMOVED***

// BufferPool represents a pool of buffers. The *sync.Pool type satisfies this
// interface.  The type of the value stored in a pool is not specified.
type BufferPool interface ***REMOVED***
	// Get gets a value from the pool or returns nil if the pool is empty.
	Get() interface***REMOVED******REMOVED***
	// Put adds a value to the pool.
	Put(interface***REMOVED******REMOVED***)
***REMOVED***

// writePoolData is the type added to the write buffer pool. This wrapper is
// used to prevent applications from peeking at and depending on the values
// added to the pool.
type writePoolData struct***REMOVED*** buf []byte ***REMOVED***

// The Conn type represents a WebSocket connection.
type Conn struct ***REMOVED***
	conn        net.Conn
	isServer    bool
	subprotocol string

	// Write fields
	mu            chan struct***REMOVED******REMOVED*** // used as mutex to protect write to conn
	writeBuf      []byte        // frame is constructed in this buffer.
	writePool     BufferPool
	writeBufSize  int
	writeDeadline time.Time
	writer        io.WriteCloser // the current writer returned to the application
	isWriting     bool           // for best-effort concurrent write detection

	writeErrMu sync.Mutex
	writeErr   error

	enableWriteCompression bool
	compressionLevel       int
	newCompressionWriter   func(io.WriteCloser, int) io.WriteCloser

	// Read fields
	reader  io.ReadCloser // the current reader returned to the application
	readErr error
	br      *bufio.Reader
	// bytes remaining in current frame.
	// set setReadRemaining to safely update this value and prevent overflow
	readRemaining int64
	readFinal     bool  // true the current message has more frames.
	readLength    int64 // Message size.
	readLimit     int64 // Maximum message size.
	readMaskPos   int
	readMaskKey   [4]byte
	handlePong    func(string) error
	handlePing    func(string) error
	handleClose   func(int, string) error
	readErrCount  int
	messageReader *messageReader // the current low-level reader

	readDecompress         bool // whether last read frame had RSV1 set
	newDecompressionReader func(io.Reader) io.ReadCloser
***REMOVED***

func newConn(conn net.Conn, isServer bool, readBufferSize, writeBufferSize int, writeBufferPool BufferPool, br *bufio.Reader, writeBuf []byte) *Conn ***REMOVED***

	if br == nil ***REMOVED***
		if readBufferSize == 0 ***REMOVED***
			readBufferSize = defaultReadBufferSize
		***REMOVED*** else if readBufferSize < maxControlFramePayloadSize ***REMOVED***
			// must be large enough for control frame
			readBufferSize = maxControlFramePayloadSize
		***REMOVED***
		br = bufio.NewReaderSize(conn, readBufferSize)
	***REMOVED***

	if writeBufferSize <= 0 ***REMOVED***
		writeBufferSize = defaultWriteBufferSize
	***REMOVED***
	writeBufferSize += maxFrameHeaderSize

	if writeBuf == nil && writeBufferPool == nil ***REMOVED***
		writeBuf = make([]byte, writeBufferSize)
	***REMOVED***

	mu := make(chan struct***REMOVED******REMOVED***, 1)
	mu <- struct***REMOVED******REMOVED******REMOVED******REMOVED***
	c := &Conn***REMOVED***
		isServer:               isServer,
		br:                     br,
		conn:                   conn,
		mu:                     mu,
		readFinal:              true,
		writeBuf:               writeBuf,
		writePool:              writeBufferPool,
		writeBufSize:           writeBufferSize,
		enableWriteCompression: true,
		compressionLevel:       defaultCompressionLevel,
	***REMOVED***
	c.SetCloseHandler(nil)
	c.SetPingHandler(nil)
	c.SetPongHandler(nil)
	return c
***REMOVED***

// setReadRemaining tracks the number of bytes remaining on the connection. If n
// overflows, an ErrReadLimit is returned.
func (c *Conn) setReadRemaining(n int64) error ***REMOVED***
	if n < 0 ***REMOVED***
		return ErrReadLimit
	***REMOVED***

	c.readRemaining = n
	return nil
***REMOVED***

// Subprotocol returns the negotiated protocol for the connection.
func (c *Conn) Subprotocol() string ***REMOVED***
	return c.subprotocol
***REMOVED***

// Close closes the underlying network connection without sending or waiting
// for a close message.
func (c *Conn) Close() error ***REMOVED***
	return c.conn.Close()
***REMOVED***

// LocalAddr returns the local network address.
func (c *Conn) LocalAddr() net.Addr ***REMOVED***
	return c.conn.LocalAddr()
***REMOVED***

// RemoteAddr returns the remote network address.
func (c *Conn) RemoteAddr() net.Addr ***REMOVED***
	return c.conn.RemoteAddr()
***REMOVED***

// Write methods

func (c *Conn) writeFatal(err error) error ***REMOVED***
	err = hideTempErr(err)
	c.writeErrMu.Lock()
	if c.writeErr == nil ***REMOVED***
		c.writeErr = err
	***REMOVED***
	c.writeErrMu.Unlock()
	return err
***REMOVED***

func (c *Conn) read(n int) ([]byte, error) ***REMOVED***
	p, err := c.br.Peek(n)
	if err == io.EOF ***REMOVED***
		err = errUnexpectedEOF
	***REMOVED***
	c.br.Discard(len(p))
	return p, err
***REMOVED***

func (c *Conn) write(frameType int, deadline time.Time, buf0, buf1 []byte) error ***REMOVED***
	<-c.mu
	defer func() ***REMOVED*** c.mu <- struct***REMOVED******REMOVED******REMOVED******REMOVED*** ***REMOVED***()

	c.writeErrMu.Lock()
	err := c.writeErr
	c.writeErrMu.Unlock()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	c.conn.SetWriteDeadline(deadline)
	if len(buf1) == 0 ***REMOVED***
		_, err = c.conn.Write(buf0)
	***REMOVED*** else ***REMOVED***
		err = c.writeBufs(buf0, buf1)
	***REMOVED***
	if err != nil ***REMOVED***
		return c.writeFatal(err)
	***REMOVED***
	if frameType == CloseMessage ***REMOVED***
		c.writeFatal(ErrCloseSent)
	***REMOVED***
	return nil
***REMOVED***

// WriteControl writes a control message with the given deadline. The allowed
// message types are CloseMessage, PingMessage and PongMessage.
func (c *Conn) WriteControl(messageType int, data []byte, deadline time.Time) error ***REMOVED***
	if !isControl(messageType) ***REMOVED***
		return errBadWriteOpCode
	***REMOVED***
	if len(data) > maxControlFramePayloadSize ***REMOVED***
		return errInvalidControlFrame
	***REMOVED***

	b0 := byte(messageType) | finalBit
	b1 := byte(len(data))
	if !c.isServer ***REMOVED***
		b1 |= maskBit
	***REMOVED***

	buf := make([]byte, 0, maxFrameHeaderSize+maxControlFramePayloadSize)
	buf = append(buf, b0, b1)

	if c.isServer ***REMOVED***
		buf = append(buf, data...)
	***REMOVED*** else ***REMOVED***
		key := newMaskKey()
		buf = append(buf, key[:]...)
		buf = append(buf, data...)
		maskBytes(key, 0, buf[6:])
	***REMOVED***

	d := 1000 * time.Hour
	if !deadline.IsZero() ***REMOVED***
		d = deadline.Sub(time.Now())
		if d < 0 ***REMOVED***
			return errWriteTimeout
		***REMOVED***
	***REMOVED***

	timer := time.NewTimer(d)
	select ***REMOVED***
	case <-c.mu:
		timer.Stop()
	case <-timer.C:
		return errWriteTimeout
	***REMOVED***
	defer func() ***REMOVED*** c.mu <- struct***REMOVED******REMOVED******REMOVED******REMOVED*** ***REMOVED***()

	c.writeErrMu.Lock()
	err := c.writeErr
	c.writeErrMu.Unlock()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	c.conn.SetWriteDeadline(deadline)
	_, err = c.conn.Write(buf)
	if err != nil ***REMOVED***
		return c.writeFatal(err)
	***REMOVED***
	if messageType == CloseMessage ***REMOVED***
		c.writeFatal(ErrCloseSent)
	***REMOVED***
	return err
***REMOVED***

// beginMessage prepares a connection and message writer for a new message.
func (c *Conn) beginMessage(mw *messageWriter, messageType int) error ***REMOVED***
	// Close previous writer if not already closed by the application. It's
	// probably better to return an error in this situation, but we cannot
	// change this without breaking existing applications.
	if c.writer != nil ***REMOVED***
		c.writer.Close()
		c.writer = nil
	***REMOVED***

	if !isControl(messageType) && !isData(messageType) ***REMOVED***
		return errBadWriteOpCode
	***REMOVED***

	c.writeErrMu.Lock()
	err := c.writeErr
	c.writeErrMu.Unlock()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	mw.c = c
	mw.frameType = messageType
	mw.pos = maxFrameHeaderSize

	if c.writeBuf == nil ***REMOVED***
		wpd, ok := c.writePool.Get().(writePoolData)
		if ok ***REMOVED***
			c.writeBuf = wpd.buf
		***REMOVED*** else ***REMOVED***
			c.writeBuf = make([]byte, c.writeBufSize)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// NextWriter returns a writer for the next message to send. The writer's Close
// method flushes the complete message to the network.
//
// There can be at most one open writer on a connection. NextWriter closes the
// previous writer if the application has not already done so.
//
// All message types (TextMessage, BinaryMessage, CloseMessage, PingMessage and
// PongMessage) are supported.
func (c *Conn) NextWriter(messageType int) (io.WriteCloser, error) ***REMOVED***
	var mw messageWriter
	if err := c.beginMessage(&mw, messageType); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	c.writer = &mw
	if c.newCompressionWriter != nil && c.enableWriteCompression && isData(messageType) ***REMOVED***
		w := c.newCompressionWriter(c.writer, c.compressionLevel)
		mw.compress = true
		c.writer = w
	***REMOVED***
	return c.writer, nil
***REMOVED***

type messageWriter struct ***REMOVED***
	c         *Conn
	compress  bool // whether next call to flushFrame should set RSV1
	pos       int  // end of data in writeBuf.
	frameType int  // type of the current frame.
	err       error
***REMOVED***

func (w *messageWriter) endMessage(err error) error ***REMOVED***
	if w.err != nil ***REMOVED***
		return err
	***REMOVED***
	c := w.c
	w.err = err
	c.writer = nil
	if c.writePool != nil ***REMOVED***
		c.writePool.Put(writePoolData***REMOVED***buf: c.writeBuf***REMOVED***)
		c.writeBuf = nil
	***REMOVED***
	return err
***REMOVED***

// flushFrame writes buffered data and extra as a frame to the network. The
// final argument indicates that this is the last frame in the message.
func (w *messageWriter) flushFrame(final bool, extra []byte) error ***REMOVED***
	c := w.c
	length := w.pos - maxFrameHeaderSize + len(extra)

	// Check for invalid control frames.
	if isControl(w.frameType) &&
		(!final || length > maxControlFramePayloadSize) ***REMOVED***
		return w.endMessage(errInvalidControlFrame)
	***REMOVED***

	b0 := byte(w.frameType)
	if final ***REMOVED***
		b0 |= finalBit
	***REMOVED***
	if w.compress ***REMOVED***
		b0 |= rsv1Bit
	***REMOVED***
	w.compress = false

	b1 := byte(0)
	if !c.isServer ***REMOVED***
		b1 |= maskBit
	***REMOVED***

	// Assume that the frame starts at beginning of c.writeBuf.
	framePos := 0
	if c.isServer ***REMOVED***
		// Adjust up if mask not included in the header.
		framePos = 4
	***REMOVED***

	switch ***REMOVED***
	case length >= 65536:
		c.writeBuf[framePos] = b0
		c.writeBuf[framePos+1] = b1 | 127
		binary.BigEndian.PutUint64(c.writeBuf[framePos+2:], uint64(length))
	case length > 125:
		framePos += 6
		c.writeBuf[framePos] = b0
		c.writeBuf[framePos+1] = b1 | 126
		binary.BigEndian.PutUint16(c.writeBuf[framePos+2:], uint16(length))
	default:
		framePos += 8
		c.writeBuf[framePos] = b0
		c.writeBuf[framePos+1] = b1 | byte(length)
	***REMOVED***

	if !c.isServer ***REMOVED***
		key := newMaskKey()
		copy(c.writeBuf[maxFrameHeaderSize-4:], key[:])
		maskBytes(key, 0, c.writeBuf[maxFrameHeaderSize:w.pos])
		if len(extra) > 0 ***REMOVED***
			return w.endMessage(c.writeFatal(errors.New("websocket: internal error, extra used in client mode")))
		***REMOVED***
	***REMOVED***

	// Write the buffers to the connection with best-effort detection of
	// concurrent writes. See the concurrency section in the package
	// documentation for more info.

	if c.isWriting ***REMOVED***
		panic("concurrent write to websocket connection")
	***REMOVED***
	c.isWriting = true

	err := c.write(w.frameType, c.writeDeadline, c.writeBuf[framePos:w.pos], extra)

	if !c.isWriting ***REMOVED***
		panic("concurrent write to websocket connection")
	***REMOVED***
	c.isWriting = false

	if err != nil ***REMOVED***
		return w.endMessage(err)
	***REMOVED***

	if final ***REMOVED***
		w.endMessage(errWriteClosed)
		return nil
	***REMOVED***

	// Setup for next frame.
	w.pos = maxFrameHeaderSize
	w.frameType = continuationFrame
	return nil
***REMOVED***

func (w *messageWriter) ncopy(max int) (int, error) ***REMOVED***
	n := len(w.c.writeBuf) - w.pos
	if n <= 0 ***REMOVED***
		if err := w.flushFrame(false, nil); err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		n = len(w.c.writeBuf) - w.pos
	***REMOVED***
	if n > max ***REMOVED***
		n = max
	***REMOVED***
	return n, nil
***REMOVED***

func (w *messageWriter) Write(p []byte) (int, error) ***REMOVED***
	if w.err != nil ***REMOVED***
		return 0, w.err
	***REMOVED***

	if len(p) > 2*len(w.c.writeBuf) && w.c.isServer ***REMOVED***
		// Don't buffer large messages.
		err := w.flushFrame(false, p)
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		return len(p), nil
	***REMOVED***

	nn := len(p)
	for len(p) > 0 ***REMOVED***
		n, err := w.ncopy(len(p))
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		copy(w.c.writeBuf[w.pos:], p[:n])
		w.pos += n
		p = p[n:]
	***REMOVED***
	return nn, nil
***REMOVED***

func (w *messageWriter) WriteString(p string) (int, error) ***REMOVED***
	if w.err != nil ***REMOVED***
		return 0, w.err
	***REMOVED***

	nn := len(p)
	for len(p) > 0 ***REMOVED***
		n, err := w.ncopy(len(p))
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		copy(w.c.writeBuf[w.pos:], p[:n])
		w.pos += n
		p = p[n:]
	***REMOVED***
	return nn, nil
***REMOVED***

func (w *messageWriter) ReadFrom(r io.Reader) (nn int64, err error) ***REMOVED***
	if w.err != nil ***REMOVED***
		return 0, w.err
	***REMOVED***
	for ***REMOVED***
		if w.pos == len(w.c.writeBuf) ***REMOVED***
			err = w.flushFrame(false, nil)
			if err != nil ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
		var n int
		n, err = r.Read(w.c.writeBuf[w.pos:])
		w.pos += n
		nn += int64(n)
		if err != nil ***REMOVED***
			if err == io.EOF ***REMOVED***
				err = nil
			***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return nn, err
***REMOVED***

func (w *messageWriter) Close() error ***REMOVED***
	if w.err != nil ***REMOVED***
		return w.err
	***REMOVED***
	return w.flushFrame(true, nil)
***REMOVED***

// WritePreparedMessage writes prepared message into connection.
func (c *Conn) WritePreparedMessage(pm *PreparedMessage) error ***REMOVED***
	frameType, frameData, err := pm.frame(prepareKey***REMOVED***
		isServer:         c.isServer,
		compress:         c.newCompressionWriter != nil && c.enableWriteCompression && isData(pm.messageType),
		compressionLevel: c.compressionLevel,
	***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if c.isWriting ***REMOVED***
		panic("concurrent write to websocket connection")
	***REMOVED***
	c.isWriting = true
	err = c.write(frameType, c.writeDeadline, frameData, nil)
	if !c.isWriting ***REMOVED***
		panic("concurrent write to websocket connection")
	***REMOVED***
	c.isWriting = false
	return err
***REMOVED***

// WriteMessage is a helper method for getting a writer using NextWriter,
// writing the message and closing the writer.
func (c *Conn) WriteMessage(messageType int, data []byte) error ***REMOVED***

	if c.isServer && (c.newCompressionWriter == nil || !c.enableWriteCompression) ***REMOVED***
		// Fast path with no allocations and single frame.

		var mw messageWriter
		if err := c.beginMessage(&mw, messageType); err != nil ***REMOVED***
			return err
		***REMOVED***
		n := copy(c.writeBuf[mw.pos:], data)
		mw.pos += n
		data = data[n:]
		return mw.flushFrame(true, data)
	***REMOVED***

	w, err := c.NextWriter(messageType)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if _, err = w.Write(data); err != nil ***REMOVED***
		return err
	***REMOVED***
	return w.Close()
***REMOVED***

// SetWriteDeadline sets the write deadline on the underlying network
// connection. After a write has timed out, the websocket state is corrupt and
// all future writes will return an error. A zero value for t means writes will
// not time out.
func (c *Conn) SetWriteDeadline(t time.Time) error ***REMOVED***
	c.writeDeadline = t
	return nil
***REMOVED***

// Read methods

func (c *Conn) advanceFrame() (int, error) ***REMOVED***
	// 1. Skip remainder of previous frame.

	if c.readRemaining > 0 ***REMOVED***
		if _, err := io.CopyN(ioutil.Discard, c.br, c.readRemaining); err != nil ***REMOVED***
			return noFrame, err
		***REMOVED***
	***REMOVED***

	// 2. Read and parse first two bytes of frame header.

	p, err := c.read(2)
	if err != nil ***REMOVED***
		return noFrame, err
	***REMOVED***

	final := p[0]&finalBit != 0
	frameType := int(p[0] & 0xf)
	mask := p[1]&maskBit != 0
	c.setReadRemaining(int64(p[1] & 0x7f))

	c.readDecompress = false
	if c.newDecompressionReader != nil && (p[0]&rsv1Bit) != 0 ***REMOVED***
		c.readDecompress = true
		p[0] &^= rsv1Bit
	***REMOVED***

	if rsv := p[0] & (rsv1Bit | rsv2Bit | rsv3Bit); rsv != 0 ***REMOVED***
		return noFrame, c.handleProtocolError("unexpected reserved bits 0x" + strconv.FormatInt(int64(rsv), 16))
	***REMOVED***

	switch frameType ***REMOVED***
	case CloseMessage, PingMessage, PongMessage:
		if c.readRemaining > maxControlFramePayloadSize ***REMOVED***
			return noFrame, c.handleProtocolError("control frame length > 125")
		***REMOVED***
		if !final ***REMOVED***
			return noFrame, c.handleProtocolError("control frame not final")
		***REMOVED***
	case TextMessage, BinaryMessage:
		if !c.readFinal ***REMOVED***
			return noFrame, c.handleProtocolError("message start before final message frame")
		***REMOVED***
		c.readFinal = final
	case continuationFrame:
		if c.readFinal ***REMOVED***
			return noFrame, c.handleProtocolError("continuation after final message frame")
		***REMOVED***
		c.readFinal = final
	default:
		return noFrame, c.handleProtocolError("unknown opcode " + strconv.Itoa(frameType))
	***REMOVED***

	// 3. Read and parse frame length as per
	// https://tools.ietf.org/html/rfc6455#section-5.2
	//
	// The length of the "Payload data", in bytes: if 0-125, that is the payload
	// length.
	// - If 126, the following 2 bytes interpreted as a 16-bit unsigned
	// integer are the payload length.
	// - If 127, the following 8 bytes interpreted as
	// a 64-bit unsigned integer (the most significant bit MUST be 0) are the
	// payload length. Multibyte length quantities are expressed in network byte
	// order.

	switch c.readRemaining ***REMOVED***
	case 126:
		p, err := c.read(2)
		if err != nil ***REMOVED***
			return noFrame, err
		***REMOVED***

		if err := c.setReadRemaining(int64(binary.BigEndian.Uint16(p))); err != nil ***REMOVED***
			return noFrame, err
		***REMOVED***
	case 127:
		p, err := c.read(8)
		if err != nil ***REMOVED***
			return noFrame, err
		***REMOVED***

		if err := c.setReadRemaining(int64(binary.BigEndian.Uint64(p))); err != nil ***REMOVED***
			return noFrame, err
		***REMOVED***
	***REMOVED***

	// 4. Handle frame masking.

	if mask != c.isServer ***REMOVED***
		return noFrame, c.handleProtocolError("incorrect mask flag")
	***REMOVED***

	if mask ***REMOVED***
		c.readMaskPos = 0
		p, err := c.read(len(c.readMaskKey))
		if err != nil ***REMOVED***
			return noFrame, err
		***REMOVED***
		copy(c.readMaskKey[:], p)
	***REMOVED***

	// 5. For text and binary messages, enforce read limit and return.

	if frameType == continuationFrame || frameType == TextMessage || frameType == BinaryMessage ***REMOVED***

		c.readLength += c.readRemaining
		// Don't allow readLength to overflow in the presence of a large readRemaining
		// counter.
		if c.readLength < 0 ***REMOVED***
			return noFrame, ErrReadLimit
		***REMOVED***

		if c.readLimit > 0 && c.readLength > c.readLimit ***REMOVED***
			c.WriteControl(CloseMessage, FormatCloseMessage(CloseMessageTooBig, ""), time.Now().Add(writeWait))
			return noFrame, ErrReadLimit
		***REMOVED***

		return frameType, nil
	***REMOVED***

	// 6. Read control frame payload.

	var payload []byte
	if c.readRemaining > 0 ***REMOVED***
		payload, err = c.read(int(c.readRemaining))
		c.setReadRemaining(0)
		if err != nil ***REMOVED***
			return noFrame, err
		***REMOVED***
		if c.isServer ***REMOVED***
			maskBytes(c.readMaskKey, 0, payload)
		***REMOVED***
	***REMOVED***

	// 7. Process control frame payload.

	switch frameType ***REMOVED***
	case PongMessage:
		if err := c.handlePong(string(payload)); err != nil ***REMOVED***
			return noFrame, err
		***REMOVED***
	case PingMessage:
		if err := c.handlePing(string(payload)); err != nil ***REMOVED***
			return noFrame, err
		***REMOVED***
	case CloseMessage:
		closeCode := CloseNoStatusReceived
		closeText := ""
		if len(payload) >= 2 ***REMOVED***
			closeCode = int(binary.BigEndian.Uint16(payload))
			if !isValidReceivedCloseCode(closeCode) ***REMOVED***
				return noFrame, c.handleProtocolError("invalid close code")
			***REMOVED***
			closeText = string(payload[2:])
			if !utf8.ValidString(closeText) ***REMOVED***
				return noFrame, c.handleProtocolError("invalid utf8 payload in close frame")
			***REMOVED***
		***REMOVED***
		if err := c.handleClose(closeCode, closeText); err != nil ***REMOVED***
			return noFrame, err
		***REMOVED***
		return noFrame, &CloseError***REMOVED***Code: closeCode, Text: closeText***REMOVED***
	***REMOVED***

	return frameType, nil
***REMOVED***

func (c *Conn) handleProtocolError(message string) error ***REMOVED***
	c.WriteControl(CloseMessage, FormatCloseMessage(CloseProtocolError, message), time.Now().Add(writeWait))
	return errors.New("websocket: " + message)
***REMOVED***

// NextReader returns the next data message received from the peer. The
// returned messageType is either TextMessage or BinaryMessage.
//
// There can be at most one open reader on a connection. NextReader discards
// the previous message if the application has not already consumed it.
//
// Applications must break out of the application's read loop when this method
// returns a non-nil error value. Errors returned from this method are
// permanent. Once this method returns a non-nil error, all subsequent calls to
// this method return the same error.
func (c *Conn) NextReader() (messageType int, r io.Reader, err error) ***REMOVED***
	// Close previous reader, only relevant for decompression.
	if c.reader != nil ***REMOVED***
		c.reader.Close()
		c.reader = nil
	***REMOVED***

	c.messageReader = nil
	c.readLength = 0

	for c.readErr == nil ***REMOVED***
		frameType, err := c.advanceFrame()
		if err != nil ***REMOVED***
			c.readErr = hideTempErr(err)
			break
		***REMOVED***

		if frameType == TextMessage || frameType == BinaryMessage ***REMOVED***
			c.messageReader = &messageReader***REMOVED***c***REMOVED***
			c.reader = c.messageReader
			if c.readDecompress ***REMOVED***
				c.reader = c.newDecompressionReader(c.reader)
			***REMOVED***
			return frameType, c.reader, nil
		***REMOVED***
	***REMOVED***

	// Applications that do handle the error returned from this method spin in
	// tight loop on connection failure. To help application developers detect
	// this error, panic on repeated reads to the failed connection.
	c.readErrCount++
	if c.readErrCount >= 1000 ***REMOVED***
		panic("repeated read on failed websocket connection")
	***REMOVED***

	return noFrame, nil, c.readErr
***REMOVED***

type messageReader struct***REMOVED*** c *Conn ***REMOVED***

func (r *messageReader) Read(b []byte) (int, error) ***REMOVED***
	c := r.c
	if c.messageReader != r ***REMOVED***
		return 0, io.EOF
	***REMOVED***

	for c.readErr == nil ***REMOVED***

		if c.readRemaining > 0 ***REMOVED***
			if int64(len(b)) > c.readRemaining ***REMOVED***
				b = b[:c.readRemaining]
			***REMOVED***
			n, err := c.br.Read(b)
			c.readErr = hideTempErr(err)
			if c.isServer ***REMOVED***
				c.readMaskPos = maskBytes(c.readMaskKey, c.readMaskPos, b[:n])
			***REMOVED***
			rem := c.readRemaining
			rem -= int64(n)
			c.setReadRemaining(rem)
			if c.readRemaining > 0 && c.readErr == io.EOF ***REMOVED***
				c.readErr = errUnexpectedEOF
			***REMOVED***
			return n, c.readErr
		***REMOVED***

		if c.readFinal ***REMOVED***
			c.messageReader = nil
			return 0, io.EOF
		***REMOVED***

		frameType, err := c.advanceFrame()
		switch ***REMOVED***
		case err != nil:
			c.readErr = hideTempErr(err)
		case frameType == TextMessage || frameType == BinaryMessage:
			c.readErr = errors.New("websocket: internal error, unexpected text or binary in Reader")
		***REMOVED***
	***REMOVED***

	err := c.readErr
	if err == io.EOF && c.messageReader == r ***REMOVED***
		err = errUnexpectedEOF
	***REMOVED***
	return 0, err
***REMOVED***

func (r *messageReader) Close() error ***REMOVED***
	return nil
***REMOVED***

// ReadMessage is a helper method for getting a reader using NextReader and
// reading from that reader to a buffer.
func (c *Conn) ReadMessage() (messageType int, p []byte, err error) ***REMOVED***
	var r io.Reader
	messageType, r, err = c.NextReader()
	if err != nil ***REMOVED***
		return messageType, nil, err
	***REMOVED***
	p, err = ioutil.ReadAll(r)
	return messageType, p, err
***REMOVED***

// SetReadDeadline sets the read deadline on the underlying network connection.
// After a read has timed out, the websocket connection state is corrupt and
// all future reads will return an error. A zero value for t means reads will
// not time out.
func (c *Conn) SetReadDeadline(t time.Time) error ***REMOVED***
	return c.conn.SetReadDeadline(t)
***REMOVED***

// SetReadLimit sets the maximum size in bytes for a message read from the peer. If a
// message exceeds the limit, the connection sends a close message to the peer
// and returns ErrReadLimit to the application.
func (c *Conn) SetReadLimit(limit int64) ***REMOVED***
	c.readLimit = limit
***REMOVED***

// CloseHandler returns the current close handler
func (c *Conn) CloseHandler() func(code int, text string) error ***REMOVED***
	return c.handleClose
***REMOVED***

// SetCloseHandler sets the handler for close messages received from the peer.
// The code argument to h is the received close code or CloseNoStatusReceived
// if the close message is empty. The default close handler sends a close
// message back to the peer.
//
// The handler function is called from the NextReader, ReadMessage and message
// reader Read methods. The application must read the connection to process
// close messages as described in the section on Control Messages above.
//
// The connection read methods return a CloseError when a close message is
// received. Most applications should handle close messages as part of their
// normal error handling. Applications should only set a close handler when the
// application must perform some action before sending a close message back to
// the peer.
func (c *Conn) SetCloseHandler(h func(code int, text string) error) ***REMOVED***
	if h == nil ***REMOVED***
		h = func(code int, text string) error ***REMOVED***
			message := FormatCloseMessage(code, "")
			c.WriteControl(CloseMessage, message, time.Now().Add(writeWait))
			return nil
		***REMOVED***
	***REMOVED***
	c.handleClose = h
***REMOVED***

// PingHandler returns the current ping handler
func (c *Conn) PingHandler() func(appData string) error ***REMOVED***
	return c.handlePing
***REMOVED***

// SetPingHandler sets the handler for ping messages received from the peer.
// The appData argument to h is the PING message application data. The default
// ping handler sends a pong to the peer.
//
// The handler function is called from the NextReader, ReadMessage and message
// reader Read methods. The application must read the connection to process
// ping messages as described in the section on Control Messages above.
func (c *Conn) SetPingHandler(h func(appData string) error) ***REMOVED***
	if h == nil ***REMOVED***
		h = func(message string) error ***REMOVED***
			err := c.WriteControl(PongMessage, []byte(message), time.Now().Add(writeWait))
			if err == ErrCloseSent ***REMOVED***
				return nil
			***REMOVED*** else if e, ok := err.(net.Error); ok && e.Temporary() ***REMOVED***
				return nil
			***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	c.handlePing = h
***REMOVED***

// PongHandler returns the current pong handler
func (c *Conn) PongHandler() func(appData string) error ***REMOVED***
	return c.handlePong
***REMOVED***

// SetPongHandler sets the handler for pong messages received from the peer.
// The appData argument to h is the PONG message application data. The default
// pong handler does nothing.
//
// The handler function is called from the NextReader, ReadMessage and message
// reader Read methods. The application must read the connection to process
// pong messages as described in the section on Control Messages above.
func (c *Conn) SetPongHandler(h func(appData string) error) ***REMOVED***
	if h == nil ***REMOVED***
		h = func(string) error ***REMOVED*** return nil ***REMOVED***
	***REMOVED***
	c.handlePong = h
***REMOVED***

// UnderlyingConn returns the internal net.Conn. This can be used to further
// modifications to connection specific flags.
func (c *Conn) UnderlyingConn() net.Conn ***REMOVED***
	return c.conn
***REMOVED***

// EnableWriteCompression enables and disables write compression of
// subsequent text and binary messages. This function is a noop if
// compression was not negotiated with the peer.
func (c *Conn) EnableWriteCompression(enable bool) ***REMOVED***
	c.enableWriteCompression = enable
***REMOVED***

// SetCompressionLevel sets the flate compression level for subsequent text and
// binary messages. This function is a noop if compression was not negotiated
// with the peer. See the compress/flate package for a description of
// compression levels.
func (c *Conn) SetCompressionLevel(level int) error ***REMOVED***
	if !isValidCompressionLevel(level) ***REMOVED***
		return errors.New("websocket: invalid compression level")
	***REMOVED***
	c.compressionLevel = level
	return nil
***REMOVED***

// FormatCloseMessage formats closeCode and text as a WebSocket close message.
// An empty message is returned for code CloseNoStatusReceived.
func FormatCloseMessage(closeCode int, text string) []byte ***REMOVED***
	if closeCode == CloseNoStatusReceived ***REMOVED***
		// Return empty message because it's illegal to send
		// CloseNoStatusReceived. Return non-nil value in case application
		// checks for nil.
		return []byte***REMOVED******REMOVED***
	***REMOVED***
	buf := make([]byte, 2+len(text))
	binary.BigEndian.PutUint16(buf, uint16(closeCode))
	copy(buf[2:], text)
	return buf
***REMOVED***
