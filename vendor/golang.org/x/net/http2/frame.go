// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http2

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"

	"golang.org/x/net/http/httpguts"
	"golang.org/x/net/http2/hpack"
)

const frameHeaderLen = 9

var padZeros = make([]byte, 255) // zeros for padding

// A FrameType is a registered frame type as defined in
// http://http2.github.io/http2-spec/#rfc.section.11.2
type FrameType uint8

const (
	FrameData         FrameType = 0x0
	FrameHeaders      FrameType = 0x1
	FramePriority     FrameType = 0x2
	FrameRSTStream    FrameType = 0x3
	FrameSettings     FrameType = 0x4
	FramePushPromise  FrameType = 0x5
	FramePing         FrameType = 0x6
	FrameGoAway       FrameType = 0x7
	FrameWindowUpdate FrameType = 0x8
	FrameContinuation FrameType = 0x9
)

var frameName = map[FrameType]string***REMOVED***
	FrameData:         "DATA",
	FrameHeaders:      "HEADERS",
	FramePriority:     "PRIORITY",
	FrameRSTStream:    "RST_STREAM",
	FrameSettings:     "SETTINGS",
	FramePushPromise:  "PUSH_PROMISE",
	FramePing:         "PING",
	FrameGoAway:       "GOAWAY",
	FrameWindowUpdate: "WINDOW_UPDATE",
	FrameContinuation: "CONTINUATION",
***REMOVED***

func (t FrameType) String() string ***REMOVED***
	if s, ok := frameName[t]; ok ***REMOVED***
		return s
	***REMOVED***
	return fmt.Sprintf("UNKNOWN_FRAME_TYPE_%d", uint8(t))
***REMOVED***

// Flags is a bitmask of HTTP/2 flags.
// The meaning of flags varies depending on the frame type.
type Flags uint8

// Has reports whether f contains all (0 or more) flags in v.
func (f Flags) Has(v Flags) bool ***REMOVED***
	return (f & v) == v
***REMOVED***

// Frame-specific FrameHeader flag bits.
const (
	// Data Frame
	FlagDataEndStream Flags = 0x1
	FlagDataPadded    Flags = 0x8

	// Headers Frame
	FlagHeadersEndStream  Flags = 0x1
	FlagHeadersEndHeaders Flags = 0x4
	FlagHeadersPadded     Flags = 0x8
	FlagHeadersPriority   Flags = 0x20

	// Settings Frame
	FlagSettingsAck Flags = 0x1

	// Ping Frame
	FlagPingAck Flags = 0x1

	// Continuation Frame
	FlagContinuationEndHeaders Flags = 0x4

	FlagPushPromiseEndHeaders Flags = 0x4
	FlagPushPromisePadded     Flags = 0x8
)

var flagName = map[FrameType]map[Flags]string***REMOVED***
	FrameData: ***REMOVED***
		FlagDataEndStream: "END_STREAM",
		FlagDataPadded:    "PADDED",
	***REMOVED***,
	FrameHeaders: ***REMOVED***
		FlagHeadersEndStream:  "END_STREAM",
		FlagHeadersEndHeaders: "END_HEADERS",
		FlagHeadersPadded:     "PADDED",
		FlagHeadersPriority:   "PRIORITY",
	***REMOVED***,
	FrameSettings: ***REMOVED***
		FlagSettingsAck: "ACK",
	***REMOVED***,
	FramePing: ***REMOVED***
		FlagPingAck: "ACK",
	***REMOVED***,
	FrameContinuation: ***REMOVED***
		FlagContinuationEndHeaders: "END_HEADERS",
	***REMOVED***,
	FramePushPromise: ***REMOVED***
		FlagPushPromiseEndHeaders: "END_HEADERS",
		FlagPushPromisePadded:     "PADDED",
	***REMOVED***,
***REMOVED***

// a frameParser parses a frame given its FrameHeader and payload
// bytes. The length of payload will always equal fh.Length (which
// might be 0).
type frameParser func(fc *frameCache, fh FrameHeader, countError func(string), payload []byte) (Frame, error)

var frameParsers = map[FrameType]frameParser***REMOVED***
	FrameData:         parseDataFrame,
	FrameHeaders:      parseHeadersFrame,
	FramePriority:     parsePriorityFrame,
	FrameRSTStream:    parseRSTStreamFrame,
	FrameSettings:     parseSettingsFrame,
	FramePushPromise:  parsePushPromise,
	FramePing:         parsePingFrame,
	FrameGoAway:       parseGoAwayFrame,
	FrameWindowUpdate: parseWindowUpdateFrame,
	FrameContinuation: parseContinuationFrame,
***REMOVED***

func typeFrameParser(t FrameType) frameParser ***REMOVED***
	if f := frameParsers[t]; f != nil ***REMOVED***
		return f
	***REMOVED***
	return parseUnknownFrame
***REMOVED***

// A FrameHeader is the 9 byte header of all HTTP/2 frames.
//
// See http://http2.github.io/http2-spec/#FrameHeader
type FrameHeader struct ***REMOVED***
	valid bool // caller can access []byte fields in the Frame

	// Type is the 1 byte frame type. There are ten standard frame
	// types, but extension frame types may be written by WriteRawFrame
	// and will be returned by ReadFrame (as UnknownFrame).
	Type FrameType

	// Flags are the 1 byte of 8 potential bit flags per frame.
	// They are specific to the frame type.
	Flags Flags

	// Length is the length of the frame, not including the 9 byte header.
	// The maximum size is one byte less than 16MB (uint24), but only
	// frames up to 16KB are allowed without peer agreement.
	Length uint32

	// StreamID is which stream this frame is for. Certain frames
	// are not stream-specific, in which case this field is 0.
	StreamID uint32
***REMOVED***

// Header returns h. It exists so FrameHeaders can be embedded in other
// specific frame types and implement the Frame interface.
func (h FrameHeader) Header() FrameHeader ***REMOVED*** return h ***REMOVED***

func (h FrameHeader) String() string ***REMOVED***
	var buf bytes.Buffer
	buf.WriteString("[FrameHeader ")
	h.writeDebug(&buf)
	buf.WriteByte(']')
	return buf.String()
***REMOVED***

func (h FrameHeader) writeDebug(buf *bytes.Buffer) ***REMOVED***
	buf.WriteString(h.Type.String())
	if h.Flags != 0 ***REMOVED***
		buf.WriteString(" flags=")
		set := 0
		for i := uint8(0); i < 8; i++ ***REMOVED***
			if h.Flags&(1<<i) == 0 ***REMOVED***
				continue
			***REMOVED***
			set++
			if set > 1 ***REMOVED***
				buf.WriteByte('|')
			***REMOVED***
			name := flagName[h.Type][Flags(1<<i)]
			if name != "" ***REMOVED***
				buf.WriteString(name)
			***REMOVED*** else ***REMOVED***
				fmt.Fprintf(buf, "0x%x", 1<<i)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if h.StreamID != 0 ***REMOVED***
		fmt.Fprintf(buf, " stream=%d", h.StreamID)
	***REMOVED***
	fmt.Fprintf(buf, " len=%d", h.Length)
***REMOVED***

func (h *FrameHeader) checkValid() ***REMOVED***
	if !h.valid ***REMOVED***
		panic("Frame accessor called on non-owned Frame")
	***REMOVED***
***REMOVED***

func (h *FrameHeader) invalidate() ***REMOVED*** h.valid = false ***REMOVED***

// frame header bytes.
// Used only by ReadFrameHeader.
var fhBytes = sync.Pool***REMOVED***
	New: func() interface***REMOVED******REMOVED*** ***REMOVED***
		buf := make([]byte, frameHeaderLen)
		return &buf
	***REMOVED***,
***REMOVED***

// ReadFrameHeader reads 9 bytes from r and returns a FrameHeader.
// Most users should use Framer.ReadFrame instead.
func ReadFrameHeader(r io.Reader) (FrameHeader, error) ***REMOVED***
	bufp := fhBytes.Get().(*[]byte)
	defer fhBytes.Put(bufp)
	return readFrameHeader(*bufp, r)
***REMOVED***

func readFrameHeader(buf []byte, r io.Reader) (FrameHeader, error) ***REMOVED***
	_, err := io.ReadFull(r, buf[:frameHeaderLen])
	if err != nil ***REMOVED***
		return FrameHeader***REMOVED******REMOVED***, err
	***REMOVED***
	return FrameHeader***REMOVED***
		Length:   (uint32(buf[0])<<16 | uint32(buf[1])<<8 | uint32(buf[2])),
		Type:     FrameType(buf[3]),
		Flags:    Flags(buf[4]),
		StreamID: binary.BigEndian.Uint32(buf[5:]) & (1<<31 - 1),
		valid:    true,
	***REMOVED***, nil
***REMOVED***

// A Frame is the base interface implemented by all frame types.
// Callers will generally type-assert the specific frame type:
// *HeadersFrame, *SettingsFrame, *WindowUpdateFrame, etc.
//
// Frames are only valid until the next call to Framer.ReadFrame.
type Frame interface ***REMOVED***
	Header() FrameHeader

	// invalidate is called by Framer.ReadFrame to make this
	// frame's buffers as being invalid, since the subsequent
	// frame will reuse them.
	invalidate()
***REMOVED***

// A Framer reads and writes Frames.
type Framer struct ***REMOVED***
	r         io.Reader
	lastFrame Frame
	errDetail error

	// countError is a non-nil func that's called on a frame parse
	// error with some unique error path token. It's initialized
	// from Transport.CountError or Server.CountError.
	countError func(errToken string)

	// lastHeaderStream is non-zero if the last frame was an
	// unfinished HEADERS/CONTINUATION.
	lastHeaderStream uint32

	maxReadSize uint32
	headerBuf   [frameHeaderLen]byte

	// TODO: let getReadBuf be configurable, and use a less memory-pinning
	// allocator in server.go to minimize memory pinned for many idle conns.
	// Will probably also need to make frame invalidation have a hook too.
	getReadBuf func(size uint32) []byte
	readBuf    []byte // cache for default getReadBuf

	maxWriteSize uint32 // zero means unlimited; TODO: implement

	w    io.Writer
	wbuf []byte

	// AllowIllegalWrites permits the Framer's Write methods to
	// write frames that do not conform to the HTTP/2 spec. This
	// permits using the Framer to test other HTTP/2
	// implementations' conformance to the spec.
	// If false, the Write methods will prefer to return an error
	// rather than comply.
	AllowIllegalWrites bool

	// AllowIllegalReads permits the Framer's ReadFrame method
	// to return non-compliant frames or frame orders.
	// This is for testing and permits using the Framer to test
	// other HTTP/2 implementations' conformance to the spec.
	// It is not compatible with ReadMetaHeaders.
	AllowIllegalReads bool

	// ReadMetaHeaders if non-nil causes ReadFrame to merge
	// HEADERS and CONTINUATION frames together and return
	// MetaHeadersFrame instead.
	ReadMetaHeaders *hpack.Decoder

	// MaxHeaderListSize is the http2 MAX_HEADER_LIST_SIZE.
	// It's used only if ReadMetaHeaders is set; 0 means a sane default
	// (currently 16MB)
	// If the limit is hit, MetaHeadersFrame.Truncated is set true.
	MaxHeaderListSize uint32

	// TODO: track which type of frame & with which flags was sent
	// last. Then return an error (unless AllowIllegalWrites) if
	// we're in the middle of a header block and a
	// non-Continuation or Continuation on a different stream is
	// attempted to be written.

	logReads, logWrites bool

	debugFramer       *Framer // only use for logging written writes
	debugFramerBuf    *bytes.Buffer
	debugReadLoggerf  func(string, ...interface***REMOVED******REMOVED***)
	debugWriteLoggerf func(string, ...interface***REMOVED******REMOVED***)

	frameCache *frameCache // nil if frames aren't reused (default)
***REMOVED***

func (fr *Framer) maxHeaderListSize() uint32 ***REMOVED***
	if fr.MaxHeaderListSize == 0 ***REMOVED***
		return 16 << 20 // sane default, per docs
	***REMOVED***
	return fr.MaxHeaderListSize
***REMOVED***

func (f *Framer) startWrite(ftype FrameType, flags Flags, streamID uint32) ***REMOVED***
	// Write the FrameHeader.
	f.wbuf = append(f.wbuf[:0],
		0, // 3 bytes of length, filled in in endWrite
		0,
		0,
		byte(ftype),
		byte(flags),
		byte(streamID>>24),
		byte(streamID>>16),
		byte(streamID>>8),
		byte(streamID))
***REMOVED***

func (f *Framer) endWrite() error ***REMOVED***
	// Now that we know the final size, fill in the FrameHeader in
	// the space previously reserved for it. Abuse append.
	length := len(f.wbuf) - frameHeaderLen
	if length >= (1 << 24) ***REMOVED***
		return ErrFrameTooLarge
	***REMOVED***
	_ = append(f.wbuf[:0],
		byte(length>>16),
		byte(length>>8),
		byte(length))
	if f.logWrites ***REMOVED***
		f.logWrite()
	***REMOVED***

	n, err := f.w.Write(f.wbuf)
	if err == nil && n != len(f.wbuf) ***REMOVED***
		err = io.ErrShortWrite
	***REMOVED***
	return err
***REMOVED***

func (f *Framer) logWrite() ***REMOVED***
	if f.debugFramer == nil ***REMOVED***
		f.debugFramerBuf = new(bytes.Buffer)
		f.debugFramer = NewFramer(nil, f.debugFramerBuf)
		f.debugFramer.logReads = false // we log it ourselves, saying "wrote" below
		// Let us read anything, even if we accidentally wrote it
		// in the wrong order:
		f.debugFramer.AllowIllegalReads = true
	***REMOVED***
	f.debugFramerBuf.Write(f.wbuf)
	fr, err := f.debugFramer.ReadFrame()
	if err != nil ***REMOVED***
		f.debugWriteLoggerf("http2: Framer %p: failed to decode just-written frame", f)
		return
	***REMOVED***
	f.debugWriteLoggerf("http2: Framer %p: wrote %v", f, summarizeFrame(fr))
***REMOVED***

func (f *Framer) writeByte(v byte)     ***REMOVED*** f.wbuf = append(f.wbuf, v) ***REMOVED***
func (f *Framer) writeBytes(v []byte)  ***REMOVED*** f.wbuf = append(f.wbuf, v...) ***REMOVED***
func (f *Framer) writeUint16(v uint16) ***REMOVED*** f.wbuf = append(f.wbuf, byte(v>>8), byte(v)) ***REMOVED***
func (f *Framer) writeUint32(v uint32) ***REMOVED***
	f.wbuf = append(f.wbuf, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
***REMOVED***

const (
	minMaxFrameSize = 1 << 14
	maxFrameSize    = 1<<24 - 1
)

// SetReuseFrames allows the Framer to reuse Frames.
// If called on a Framer, Frames returned by calls to ReadFrame are only
// valid until the next call to ReadFrame.
func (fr *Framer) SetReuseFrames() ***REMOVED***
	if fr.frameCache != nil ***REMOVED***
		return
	***REMOVED***
	fr.frameCache = &frameCache***REMOVED******REMOVED***
***REMOVED***

type frameCache struct ***REMOVED***
	dataFrame DataFrame
***REMOVED***

func (fc *frameCache) getDataFrame() *DataFrame ***REMOVED***
	if fc == nil ***REMOVED***
		return &DataFrame***REMOVED******REMOVED***
	***REMOVED***
	return &fc.dataFrame
***REMOVED***

// NewFramer returns a Framer that writes frames to w and reads them from r.
func NewFramer(w io.Writer, r io.Reader) *Framer ***REMOVED***
	fr := &Framer***REMOVED***
		w:                 w,
		r:                 r,
		countError:        func(string) ***REMOVED******REMOVED***,
		logReads:          logFrameReads,
		logWrites:         logFrameWrites,
		debugReadLoggerf:  log.Printf,
		debugWriteLoggerf: log.Printf,
	***REMOVED***
	fr.getReadBuf = func(size uint32) []byte ***REMOVED***
		if cap(fr.readBuf) >= int(size) ***REMOVED***
			return fr.readBuf[:size]
		***REMOVED***
		fr.readBuf = make([]byte, size)
		return fr.readBuf
	***REMOVED***
	fr.SetMaxReadFrameSize(maxFrameSize)
	return fr
***REMOVED***

// SetMaxReadFrameSize sets the maximum size of a frame
// that will be read by a subsequent call to ReadFrame.
// It is the caller's responsibility to advertise this
// limit with a SETTINGS frame.
func (fr *Framer) SetMaxReadFrameSize(v uint32) ***REMOVED***
	if v > maxFrameSize ***REMOVED***
		v = maxFrameSize
	***REMOVED***
	fr.maxReadSize = v
***REMOVED***

// ErrorDetail returns a more detailed error of the last error
// returned by Framer.ReadFrame. For instance, if ReadFrame
// returns a StreamError with code PROTOCOL_ERROR, ErrorDetail
// will say exactly what was invalid. ErrorDetail is not guaranteed
// to return a non-nil value and like the rest of the http2 package,
// its return value is not protected by an API compatibility promise.
// ErrorDetail is reset after the next call to ReadFrame.
func (fr *Framer) ErrorDetail() error ***REMOVED***
	return fr.errDetail
***REMOVED***

// ErrFrameTooLarge is returned from Framer.ReadFrame when the peer
// sends a frame that is larger than declared with SetMaxReadFrameSize.
var ErrFrameTooLarge = errors.New("http2: frame too large")

// terminalReadFrameError reports whether err is an unrecoverable
// error from ReadFrame and no other frames should be read.
func terminalReadFrameError(err error) bool ***REMOVED***
	if _, ok := err.(StreamError); ok ***REMOVED***
		return false
	***REMOVED***
	return err != nil
***REMOVED***

// ReadFrame reads a single frame. The returned Frame is only valid
// until the next call to ReadFrame.
//
// If the frame is larger than previously set with SetMaxReadFrameSize, the
// returned error is ErrFrameTooLarge. Other errors may be of type
// ConnectionError, StreamError, or anything else from the underlying
// reader.
func (fr *Framer) ReadFrame() (Frame, error) ***REMOVED***
	fr.errDetail = nil
	if fr.lastFrame != nil ***REMOVED***
		fr.lastFrame.invalidate()
	***REMOVED***
	fh, err := readFrameHeader(fr.headerBuf[:], fr.r)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if fh.Length > fr.maxReadSize ***REMOVED***
		return nil, ErrFrameTooLarge
	***REMOVED***
	payload := fr.getReadBuf(fh.Length)
	if _, err := io.ReadFull(fr.r, payload); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	f, err := typeFrameParser(fh.Type)(fr.frameCache, fh, fr.countError, payload)
	if err != nil ***REMOVED***
		if ce, ok := err.(connError); ok ***REMOVED***
			return nil, fr.connError(ce.Code, ce.Reason)
		***REMOVED***
		return nil, err
	***REMOVED***
	if err := fr.checkFrameOrder(f); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if fr.logReads ***REMOVED***
		fr.debugReadLoggerf("http2: Framer %p: read %v", fr, summarizeFrame(f))
	***REMOVED***
	if fh.Type == FrameHeaders && fr.ReadMetaHeaders != nil ***REMOVED***
		return fr.readMetaFrame(f.(*HeadersFrame))
	***REMOVED***
	return f, nil
***REMOVED***

// connError returns ConnectionError(code) but first
// stashes away a public reason to the caller can optionally relay it
// to the peer before hanging up on them. This might help others debug
// their implementations.
func (fr *Framer) connError(code ErrCode, reason string) error ***REMOVED***
	fr.errDetail = errors.New(reason)
	return ConnectionError(code)
***REMOVED***

// checkFrameOrder reports an error if f is an invalid frame to return
// next from ReadFrame. Mostly it checks whether HEADERS and
// CONTINUATION frames are contiguous.
func (fr *Framer) checkFrameOrder(f Frame) error ***REMOVED***
	last := fr.lastFrame
	fr.lastFrame = f
	if fr.AllowIllegalReads ***REMOVED***
		return nil
	***REMOVED***

	fh := f.Header()
	if fr.lastHeaderStream != 0 ***REMOVED***
		if fh.Type != FrameContinuation ***REMOVED***
			return fr.connError(ErrCodeProtocol,
				fmt.Sprintf("got %s for stream %d; expected CONTINUATION following %s for stream %d",
					fh.Type, fh.StreamID,
					last.Header().Type, fr.lastHeaderStream))
		***REMOVED***
		if fh.StreamID != fr.lastHeaderStream ***REMOVED***
			return fr.connError(ErrCodeProtocol,
				fmt.Sprintf("got CONTINUATION for stream %d; expected stream %d",
					fh.StreamID, fr.lastHeaderStream))
		***REMOVED***
	***REMOVED*** else if fh.Type == FrameContinuation ***REMOVED***
		return fr.connError(ErrCodeProtocol, fmt.Sprintf("unexpected CONTINUATION for stream %d", fh.StreamID))
	***REMOVED***

	switch fh.Type ***REMOVED***
	case FrameHeaders, FrameContinuation:
		if fh.Flags.Has(FlagHeadersEndHeaders) ***REMOVED***
			fr.lastHeaderStream = 0
		***REMOVED*** else ***REMOVED***
			fr.lastHeaderStream = fh.StreamID
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// A DataFrame conveys arbitrary, variable-length sequences of octets
// associated with a stream.
// See http://http2.github.io/http2-spec/#rfc.section.6.1
type DataFrame struct ***REMOVED***
	FrameHeader
	data []byte
***REMOVED***

func (f *DataFrame) StreamEnded() bool ***REMOVED***
	return f.FrameHeader.Flags.Has(FlagDataEndStream)
***REMOVED***

// Data returns the frame's data octets, not including any padding
// size byte or padding suffix bytes.
// The caller must not retain the returned memory past the next
// call to ReadFrame.
func (f *DataFrame) Data() []byte ***REMOVED***
	f.checkValid()
	return f.data
***REMOVED***

func parseDataFrame(fc *frameCache, fh FrameHeader, countError func(string), payload []byte) (Frame, error) ***REMOVED***
	if fh.StreamID == 0 ***REMOVED***
		// DATA frames MUST be associated with a stream. If a
		// DATA frame is received whose stream identifier
		// field is 0x0, the recipient MUST respond with a
		// connection error (Section 5.4.1) of type
		// PROTOCOL_ERROR.
		countError("frame_data_stream_0")
		return nil, connError***REMOVED***ErrCodeProtocol, "DATA frame with stream ID 0"***REMOVED***
	***REMOVED***
	f := fc.getDataFrame()
	f.FrameHeader = fh

	var padSize byte
	if fh.Flags.Has(FlagDataPadded) ***REMOVED***
		var err error
		payload, padSize, err = readByte(payload)
		if err != nil ***REMOVED***
			countError("frame_data_pad_byte_short")
			return nil, err
		***REMOVED***
	***REMOVED***
	if int(padSize) > len(payload) ***REMOVED***
		// If the length of the padding is greater than the
		// length of the frame payload, the recipient MUST
		// treat this as a connection error.
		// Filed: https://github.com/http2/http2-spec/issues/610
		countError("frame_data_pad_too_big")
		return nil, connError***REMOVED***ErrCodeProtocol, "pad size larger than data payload"***REMOVED***
	***REMOVED***
	f.data = payload[:len(payload)-int(padSize)]
	return f, nil
***REMOVED***

var (
	errStreamID    = errors.New("invalid stream ID")
	errDepStreamID = errors.New("invalid dependent stream ID")
	errPadLength   = errors.New("pad length too large")
	errPadBytes    = errors.New("padding bytes must all be zeros unless AllowIllegalWrites is enabled")
)

func validStreamIDOrZero(streamID uint32) bool ***REMOVED***
	return streamID&(1<<31) == 0
***REMOVED***

func validStreamID(streamID uint32) bool ***REMOVED***
	return streamID != 0 && streamID&(1<<31) == 0
***REMOVED***

// WriteData writes a DATA frame.
//
// It will perform exactly one Write to the underlying Writer.
// It is the caller's responsibility not to violate the maximum frame size
// and to not call other Write methods concurrently.
func (f *Framer) WriteData(streamID uint32, endStream bool, data []byte) error ***REMOVED***
	return f.WriteDataPadded(streamID, endStream, data, nil)
***REMOVED***

// WriteDataPadded writes a DATA frame with optional padding.
//
// If pad is nil, the padding bit is not sent.
// The length of pad must not exceed 255 bytes.
// The bytes of pad must all be zero, unless f.AllowIllegalWrites is set.
//
// It will perform exactly one Write to the underlying Writer.
// It is the caller's responsibility not to violate the maximum frame size
// and to not call other Write methods concurrently.
func (f *Framer) WriteDataPadded(streamID uint32, endStream bool, data, pad []byte) error ***REMOVED***
	if !validStreamID(streamID) && !f.AllowIllegalWrites ***REMOVED***
		return errStreamID
	***REMOVED***
	if len(pad) > 0 ***REMOVED***
		if len(pad) > 255 ***REMOVED***
			return errPadLength
		***REMOVED***
		if !f.AllowIllegalWrites ***REMOVED***
			for _, b := range pad ***REMOVED***
				if b != 0 ***REMOVED***
					// "Padding octets MUST be set to zero when sending."
					return errPadBytes
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	var flags Flags
	if endStream ***REMOVED***
		flags |= FlagDataEndStream
	***REMOVED***
	if pad != nil ***REMOVED***
		flags |= FlagDataPadded
	***REMOVED***
	f.startWrite(FrameData, flags, streamID)
	if pad != nil ***REMOVED***
		f.wbuf = append(f.wbuf, byte(len(pad)))
	***REMOVED***
	f.wbuf = append(f.wbuf, data...)
	f.wbuf = append(f.wbuf, pad...)
	return f.endWrite()
***REMOVED***

// A SettingsFrame conveys configuration parameters that affect how
// endpoints communicate, such as preferences and constraints on peer
// behavior.
//
// See http://http2.github.io/http2-spec/#SETTINGS
type SettingsFrame struct ***REMOVED***
	FrameHeader
	p []byte
***REMOVED***

func parseSettingsFrame(_ *frameCache, fh FrameHeader, countError func(string), p []byte) (Frame, error) ***REMOVED***
	if fh.Flags.Has(FlagSettingsAck) && fh.Length > 0 ***REMOVED***
		// When this (ACK 0x1) bit is set, the payload of the
		// SETTINGS frame MUST be empty. Receipt of a
		// SETTINGS frame with the ACK flag set and a length
		// field value other than 0 MUST be treated as a
		// connection error (Section 5.4.1) of type
		// FRAME_SIZE_ERROR.
		countError("frame_settings_ack_with_length")
		return nil, ConnectionError(ErrCodeFrameSize)
	***REMOVED***
	if fh.StreamID != 0 ***REMOVED***
		// SETTINGS frames always apply to a connection,
		// never a single stream. The stream identifier for a
		// SETTINGS frame MUST be zero (0x0).  If an endpoint
		// receives a SETTINGS frame whose stream identifier
		// field is anything other than 0x0, the endpoint MUST
		// respond with a connection error (Section 5.4.1) of
		// type PROTOCOL_ERROR.
		countError("frame_settings_has_stream")
		return nil, ConnectionError(ErrCodeProtocol)
	***REMOVED***
	if len(p)%6 != 0 ***REMOVED***
		countError("frame_settings_mod_6")
		// Expecting even number of 6 byte settings.
		return nil, ConnectionError(ErrCodeFrameSize)
	***REMOVED***
	f := &SettingsFrame***REMOVED***FrameHeader: fh, p: p***REMOVED***
	if v, ok := f.Value(SettingInitialWindowSize); ok && v > (1<<31)-1 ***REMOVED***
		countError("frame_settings_window_size_too_big")
		// Values above the maximum flow control window size of 2^31 - 1 MUST
		// be treated as a connection error (Section 5.4.1) of type
		// FLOW_CONTROL_ERROR.
		return nil, ConnectionError(ErrCodeFlowControl)
	***REMOVED***
	return f, nil
***REMOVED***

func (f *SettingsFrame) IsAck() bool ***REMOVED***
	return f.FrameHeader.Flags.Has(FlagSettingsAck)
***REMOVED***

func (f *SettingsFrame) Value(id SettingID) (v uint32, ok bool) ***REMOVED***
	f.checkValid()
	for i := 0; i < f.NumSettings(); i++ ***REMOVED***
		if s := f.Setting(i); s.ID == id ***REMOVED***
			return s.Val, true
		***REMOVED***
	***REMOVED***
	return 0, false
***REMOVED***

// Setting returns the setting from the frame at the given 0-based index.
// The index must be >= 0 and less than f.NumSettings().
func (f *SettingsFrame) Setting(i int) Setting ***REMOVED***
	buf := f.p
	return Setting***REMOVED***
		ID:  SettingID(binary.BigEndian.Uint16(buf[i*6 : i*6+2])),
		Val: binary.BigEndian.Uint32(buf[i*6+2 : i*6+6]),
	***REMOVED***
***REMOVED***

func (f *SettingsFrame) NumSettings() int ***REMOVED*** return len(f.p) / 6 ***REMOVED***

// HasDuplicates reports whether f contains any duplicate setting IDs.
func (f *SettingsFrame) HasDuplicates() bool ***REMOVED***
	num := f.NumSettings()
	if num == 0 ***REMOVED***
		return false
	***REMOVED***
	// If it's small enough (the common case), just do the n^2
	// thing and avoid a map allocation.
	if num < 10 ***REMOVED***
		for i := 0; i < num; i++ ***REMOVED***
			idi := f.Setting(i).ID
			for j := i + 1; j < num; j++ ***REMOVED***
				idj := f.Setting(j).ID
				if idi == idj ***REMOVED***
					return true
				***REMOVED***
			***REMOVED***
		***REMOVED***
		return false
	***REMOVED***
	seen := map[SettingID]bool***REMOVED******REMOVED***
	for i := 0; i < num; i++ ***REMOVED***
		id := f.Setting(i).ID
		if seen[id] ***REMOVED***
			return true
		***REMOVED***
		seen[id] = true
	***REMOVED***
	return false
***REMOVED***

// ForeachSetting runs fn for each setting.
// It stops and returns the first error.
func (f *SettingsFrame) ForeachSetting(fn func(Setting) error) error ***REMOVED***
	f.checkValid()
	for i := 0; i < f.NumSettings(); i++ ***REMOVED***
		if err := fn(f.Setting(i)); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// WriteSettings writes a SETTINGS frame with zero or more settings
// specified and the ACK bit not set.
//
// It will perform exactly one Write to the underlying Writer.
// It is the caller's responsibility to not call other Write methods concurrently.
func (f *Framer) WriteSettings(settings ...Setting) error ***REMOVED***
	f.startWrite(FrameSettings, 0, 0)
	for _, s := range settings ***REMOVED***
		f.writeUint16(uint16(s.ID))
		f.writeUint32(s.Val)
	***REMOVED***
	return f.endWrite()
***REMOVED***

// WriteSettingsAck writes an empty SETTINGS frame with the ACK bit set.
//
// It will perform exactly one Write to the underlying Writer.
// It is the caller's responsibility to not call other Write methods concurrently.
func (f *Framer) WriteSettingsAck() error ***REMOVED***
	f.startWrite(FrameSettings, FlagSettingsAck, 0)
	return f.endWrite()
***REMOVED***

// A PingFrame is a mechanism for measuring a minimal round trip time
// from the sender, as well as determining whether an idle connection
// is still functional.
// See http://http2.github.io/http2-spec/#rfc.section.6.7
type PingFrame struct ***REMOVED***
	FrameHeader
	Data [8]byte
***REMOVED***

func (f *PingFrame) IsAck() bool ***REMOVED*** return f.Flags.Has(FlagPingAck) ***REMOVED***

func parsePingFrame(_ *frameCache, fh FrameHeader, countError func(string), payload []byte) (Frame, error) ***REMOVED***
	if len(payload) != 8 ***REMOVED***
		countError("frame_ping_length")
		return nil, ConnectionError(ErrCodeFrameSize)
	***REMOVED***
	if fh.StreamID != 0 ***REMOVED***
		countError("frame_ping_has_stream")
		return nil, ConnectionError(ErrCodeProtocol)
	***REMOVED***
	f := &PingFrame***REMOVED***FrameHeader: fh***REMOVED***
	copy(f.Data[:], payload)
	return f, nil
***REMOVED***

func (f *Framer) WritePing(ack bool, data [8]byte) error ***REMOVED***
	var flags Flags
	if ack ***REMOVED***
		flags = FlagPingAck
	***REMOVED***
	f.startWrite(FramePing, flags, 0)
	f.writeBytes(data[:])
	return f.endWrite()
***REMOVED***

// A GoAwayFrame informs the remote peer to stop creating streams on this connection.
// See http://http2.github.io/http2-spec/#rfc.section.6.8
type GoAwayFrame struct ***REMOVED***
	FrameHeader
	LastStreamID uint32
	ErrCode      ErrCode
	debugData    []byte
***REMOVED***

// DebugData returns any debug data in the GOAWAY frame. Its contents
// are not defined.
// The caller must not retain the returned memory past the next
// call to ReadFrame.
func (f *GoAwayFrame) DebugData() []byte ***REMOVED***
	f.checkValid()
	return f.debugData
***REMOVED***

func parseGoAwayFrame(_ *frameCache, fh FrameHeader, countError func(string), p []byte) (Frame, error) ***REMOVED***
	if fh.StreamID != 0 ***REMOVED***
		countError("frame_goaway_has_stream")
		return nil, ConnectionError(ErrCodeProtocol)
	***REMOVED***
	if len(p) < 8 ***REMOVED***
		countError("frame_goaway_short")
		return nil, ConnectionError(ErrCodeFrameSize)
	***REMOVED***
	return &GoAwayFrame***REMOVED***
		FrameHeader:  fh,
		LastStreamID: binary.BigEndian.Uint32(p[:4]) & (1<<31 - 1),
		ErrCode:      ErrCode(binary.BigEndian.Uint32(p[4:8])),
		debugData:    p[8:],
	***REMOVED***, nil
***REMOVED***

func (f *Framer) WriteGoAway(maxStreamID uint32, code ErrCode, debugData []byte) error ***REMOVED***
	f.startWrite(FrameGoAway, 0, 0)
	f.writeUint32(maxStreamID & (1<<31 - 1))
	f.writeUint32(uint32(code))
	f.writeBytes(debugData)
	return f.endWrite()
***REMOVED***

// An UnknownFrame is the frame type returned when the frame type is unknown
// or no specific frame type parser exists.
type UnknownFrame struct ***REMOVED***
	FrameHeader
	p []byte
***REMOVED***

// Payload returns the frame's payload (after the header).  It is not
// valid to call this method after a subsequent call to
// Framer.ReadFrame, nor is it valid to retain the returned slice.
// The memory is owned by the Framer and is invalidated when the next
// frame is read.
func (f *UnknownFrame) Payload() []byte ***REMOVED***
	f.checkValid()
	return f.p
***REMOVED***

func parseUnknownFrame(_ *frameCache, fh FrameHeader, countError func(string), p []byte) (Frame, error) ***REMOVED***
	return &UnknownFrame***REMOVED***fh, p***REMOVED***, nil
***REMOVED***

// A WindowUpdateFrame is used to implement flow control.
// See http://http2.github.io/http2-spec/#rfc.section.6.9
type WindowUpdateFrame struct ***REMOVED***
	FrameHeader
	Increment uint32 // never read with high bit set
***REMOVED***

func parseWindowUpdateFrame(_ *frameCache, fh FrameHeader, countError func(string), p []byte) (Frame, error) ***REMOVED***
	if len(p) != 4 ***REMOVED***
		countError("frame_windowupdate_bad_len")
		return nil, ConnectionError(ErrCodeFrameSize)
	***REMOVED***
	inc := binary.BigEndian.Uint32(p[:4]) & 0x7fffffff // mask off high reserved bit
	if inc == 0 ***REMOVED***
		// A receiver MUST treat the receipt of a
		// WINDOW_UPDATE frame with an flow control window
		// increment of 0 as a stream error (Section 5.4.2) of
		// type PROTOCOL_ERROR; errors on the connection flow
		// control window MUST be treated as a connection
		// error (Section 5.4.1).
		if fh.StreamID == 0 ***REMOVED***
			countError("frame_windowupdate_zero_inc_conn")
			return nil, ConnectionError(ErrCodeProtocol)
		***REMOVED***
		countError("frame_windowupdate_zero_inc_stream")
		return nil, streamError(fh.StreamID, ErrCodeProtocol)
	***REMOVED***
	return &WindowUpdateFrame***REMOVED***
		FrameHeader: fh,
		Increment:   inc,
	***REMOVED***, nil
***REMOVED***

// WriteWindowUpdate writes a WINDOW_UPDATE frame.
// The increment value must be between 1 and 2,147,483,647, inclusive.
// If the Stream ID is zero, the window update applies to the
// connection as a whole.
func (f *Framer) WriteWindowUpdate(streamID, incr uint32) error ***REMOVED***
	// "The legal range for the increment to the flow control window is 1 to 2^31-1 (2,147,483,647) octets."
	if (incr < 1 || incr > 2147483647) && !f.AllowIllegalWrites ***REMOVED***
		return errors.New("illegal window increment value")
	***REMOVED***
	f.startWrite(FrameWindowUpdate, 0, streamID)
	f.writeUint32(incr)
	return f.endWrite()
***REMOVED***

// A HeadersFrame is used to open a stream and additionally carries a
// header block fragment.
type HeadersFrame struct ***REMOVED***
	FrameHeader

	// Priority is set if FlagHeadersPriority is set in the FrameHeader.
	Priority PriorityParam

	headerFragBuf []byte // not owned
***REMOVED***

func (f *HeadersFrame) HeaderBlockFragment() []byte ***REMOVED***
	f.checkValid()
	return f.headerFragBuf
***REMOVED***

func (f *HeadersFrame) HeadersEnded() bool ***REMOVED***
	return f.FrameHeader.Flags.Has(FlagHeadersEndHeaders)
***REMOVED***

func (f *HeadersFrame) StreamEnded() bool ***REMOVED***
	return f.FrameHeader.Flags.Has(FlagHeadersEndStream)
***REMOVED***

func (f *HeadersFrame) HasPriority() bool ***REMOVED***
	return f.FrameHeader.Flags.Has(FlagHeadersPriority)
***REMOVED***

func parseHeadersFrame(_ *frameCache, fh FrameHeader, countError func(string), p []byte) (_ Frame, err error) ***REMOVED***
	hf := &HeadersFrame***REMOVED***
		FrameHeader: fh,
	***REMOVED***
	if fh.StreamID == 0 ***REMOVED***
		// HEADERS frames MUST be associated with a stream. If a HEADERS frame
		// is received whose stream identifier field is 0x0, the recipient MUST
		// respond with a connection error (Section 5.4.1) of type
		// PROTOCOL_ERROR.
		countError("frame_headers_zero_stream")
		return nil, connError***REMOVED***ErrCodeProtocol, "HEADERS frame with stream ID 0"***REMOVED***
	***REMOVED***
	var padLength uint8
	if fh.Flags.Has(FlagHeadersPadded) ***REMOVED***
		if p, padLength, err = readByte(p); err != nil ***REMOVED***
			countError("frame_headers_pad_short")
			return
		***REMOVED***
	***REMOVED***
	if fh.Flags.Has(FlagHeadersPriority) ***REMOVED***
		var v uint32
		p, v, err = readUint32(p)
		if err != nil ***REMOVED***
			countError("frame_headers_prio_short")
			return nil, err
		***REMOVED***
		hf.Priority.StreamDep = v & 0x7fffffff
		hf.Priority.Exclusive = (v != hf.Priority.StreamDep) // high bit was set
		p, hf.Priority.Weight, err = readByte(p)
		if err != nil ***REMOVED***
			countError("frame_headers_prio_weight_short")
			return nil, err
		***REMOVED***
	***REMOVED***
	if len(p)-int(padLength) < 0 ***REMOVED***
		countError("frame_headers_pad_too_big")
		return nil, streamError(fh.StreamID, ErrCodeProtocol)
	***REMOVED***
	hf.headerFragBuf = p[:len(p)-int(padLength)]
	return hf, nil
***REMOVED***

// HeadersFrameParam are the parameters for writing a HEADERS frame.
type HeadersFrameParam struct ***REMOVED***
	// StreamID is the required Stream ID to initiate.
	StreamID uint32
	// BlockFragment is part (or all) of a Header Block.
	BlockFragment []byte

	// EndStream indicates that the header block is the last that
	// the endpoint will send for the identified stream. Setting
	// this flag causes the stream to enter one of "half closed"
	// states.
	EndStream bool

	// EndHeaders indicates that this frame contains an entire
	// header block and is not followed by any
	// CONTINUATION frames.
	EndHeaders bool

	// PadLength is the optional number of bytes of zeros to add
	// to this frame.
	PadLength uint8

	// Priority, if non-zero, includes stream priority information
	// in the HEADER frame.
	Priority PriorityParam
***REMOVED***

// WriteHeaders writes a single HEADERS frame.
//
// This is a low-level header writing method. Encoding headers and
// splitting them into any necessary CONTINUATION frames is handled
// elsewhere.
//
// It will perform exactly one Write to the underlying Writer.
// It is the caller's responsibility to not call other Write methods concurrently.
func (f *Framer) WriteHeaders(p HeadersFrameParam) error ***REMOVED***
	if !validStreamID(p.StreamID) && !f.AllowIllegalWrites ***REMOVED***
		return errStreamID
	***REMOVED***
	var flags Flags
	if p.PadLength != 0 ***REMOVED***
		flags |= FlagHeadersPadded
	***REMOVED***
	if p.EndStream ***REMOVED***
		flags |= FlagHeadersEndStream
	***REMOVED***
	if p.EndHeaders ***REMOVED***
		flags |= FlagHeadersEndHeaders
	***REMOVED***
	if !p.Priority.IsZero() ***REMOVED***
		flags |= FlagHeadersPriority
	***REMOVED***
	f.startWrite(FrameHeaders, flags, p.StreamID)
	if p.PadLength != 0 ***REMOVED***
		f.writeByte(p.PadLength)
	***REMOVED***
	if !p.Priority.IsZero() ***REMOVED***
		v := p.Priority.StreamDep
		if !validStreamIDOrZero(v) && !f.AllowIllegalWrites ***REMOVED***
			return errDepStreamID
		***REMOVED***
		if p.Priority.Exclusive ***REMOVED***
			v |= 1 << 31
		***REMOVED***
		f.writeUint32(v)
		f.writeByte(p.Priority.Weight)
	***REMOVED***
	f.wbuf = append(f.wbuf, p.BlockFragment...)
	f.wbuf = append(f.wbuf, padZeros[:p.PadLength]...)
	return f.endWrite()
***REMOVED***

// A PriorityFrame specifies the sender-advised priority of a stream.
// See http://http2.github.io/http2-spec/#rfc.section.6.3
type PriorityFrame struct ***REMOVED***
	FrameHeader
	PriorityParam
***REMOVED***

// PriorityParam are the stream prioritzation parameters.
type PriorityParam struct ***REMOVED***
	// StreamDep is a 31-bit stream identifier for the
	// stream that this stream depends on. Zero means no
	// dependency.
	StreamDep uint32

	// Exclusive is whether the dependency is exclusive.
	Exclusive bool

	// Weight is the stream's zero-indexed weight. It should be
	// set together with StreamDep, or neither should be set. Per
	// the spec, "Add one to the value to obtain a weight between
	// 1 and 256."
	Weight uint8
***REMOVED***

func (p PriorityParam) IsZero() bool ***REMOVED***
	return p == PriorityParam***REMOVED******REMOVED***
***REMOVED***

func parsePriorityFrame(_ *frameCache, fh FrameHeader, countError func(string), payload []byte) (Frame, error) ***REMOVED***
	if fh.StreamID == 0 ***REMOVED***
		countError("frame_priority_zero_stream")
		return nil, connError***REMOVED***ErrCodeProtocol, "PRIORITY frame with stream ID 0"***REMOVED***
	***REMOVED***
	if len(payload) != 5 ***REMOVED***
		countError("frame_priority_bad_length")
		return nil, connError***REMOVED***ErrCodeFrameSize, fmt.Sprintf("PRIORITY frame payload size was %d; want 5", len(payload))***REMOVED***
	***REMOVED***
	v := binary.BigEndian.Uint32(payload[:4])
	streamID := v & 0x7fffffff // mask off high bit
	return &PriorityFrame***REMOVED***
		FrameHeader: fh,
		PriorityParam: PriorityParam***REMOVED***
			Weight:    payload[4],
			StreamDep: streamID,
			Exclusive: streamID != v, // was high bit set?
		***REMOVED***,
	***REMOVED***, nil
***REMOVED***

// WritePriority writes a PRIORITY frame.
//
// It will perform exactly one Write to the underlying Writer.
// It is the caller's responsibility to not call other Write methods concurrently.
func (f *Framer) WritePriority(streamID uint32, p PriorityParam) error ***REMOVED***
	if !validStreamID(streamID) && !f.AllowIllegalWrites ***REMOVED***
		return errStreamID
	***REMOVED***
	if !validStreamIDOrZero(p.StreamDep) ***REMOVED***
		return errDepStreamID
	***REMOVED***
	f.startWrite(FramePriority, 0, streamID)
	v := p.StreamDep
	if p.Exclusive ***REMOVED***
		v |= 1 << 31
	***REMOVED***
	f.writeUint32(v)
	f.writeByte(p.Weight)
	return f.endWrite()
***REMOVED***

// A RSTStreamFrame allows for abnormal termination of a stream.
// See http://http2.github.io/http2-spec/#rfc.section.6.4
type RSTStreamFrame struct ***REMOVED***
	FrameHeader
	ErrCode ErrCode
***REMOVED***

func parseRSTStreamFrame(_ *frameCache, fh FrameHeader, countError func(string), p []byte) (Frame, error) ***REMOVED***
	if len(p) != 4 ***REMOVED***
		countError("frame_rststream_bad_len")
		return nil, ConnectionError(ErrCodeFrameSize)
	***REMOVED***
	if fh.StreamID == 0 ***REMOVED***
		countError("frame_rststream_zero_stream")
		return nil, ConnectionError(ErrCodeProtocol)
	***REMOVED***
	return &RSTStreamFrame***REMOVED***fh, ErrCode(binary.BigEndian.Uint32(p[:4]))***REMOVED***, nil
***REMOVED***

// WriteRSTStream writes a RST_STREAM frame.
//
// It will perform exactly one Write to the underlying Writer.
// It is the caller's responsibility to not call other Write methods concurrently.
func (f *Framer) WriteRSTStream(streamID uint32, code ErrCode) error ***REMOVED***
	if !validStreamID(streamID) && !f.AllowIllegalWrites ***REMOVED***
		return errStreamID
	***REMOVED***
	f.startWrite(FrameRSTStream, 0, streamID)
	f.writeUint32(uint32(code))
	return f.endWrite()
***REMOVED***

// A ContinuationFrame is used to continue a sequence of header block fragments.
// See http://http2.github.io/http2-spec/#rfc.section.6.10
type ContinuationFrame struct ***REMOVED***
	FrameHeader
	headerFragBuf []byte
***REMOVED***

func parseContinuationFrame(_ *frameCache, fh FrameHeader, countError func(string), p []byte) (Frame, error) ***REMOVED***
	if fh.StreamID == 0 ***REMOVED***
		countError("frame_continuation_zero_stream")
		return nil, connError***REMOVED***ErrCodeProtocol, "CONTINUATION frame with stream ID 0"***REMOVED***
	***REMOVED***
	return &ContinuationFrame***REMOVED***fh, p***REMOVED***, nil
***REMOVED***

func (f *ContinuationFrame) HeaderBlockFragment() []byte ***REMOVED***
	f.checkValid()
	return f.headerFragBuf
***REMOVED***

func (f *ContinuationFrame) HeadersEnded() bool ***REMOVED***
	return f.FrameHeader.Flags.Has(FlagContinuationEndHeaders)
***REMOVED***

// WriteContinuation writes a CONTINUATION frame.
//
// It will perform exactly one Write to the underlying Writer.
// It is the caller's responsibility to not call other Write methods concurrently.
func (f *Framer) WriteContinuation(streamID uint32, endHeaders bool, headerBlockFragment []byte) error ***REMOVED***
	if !validStreamID(streamID) && !f.AllowIllegalWrites ***REMOVED***
		return errStreamID
	***REMOVED***
	var flags Flags
	if endHeaders ***REMOVED***
		flags |= FlagContinuationEndHeaders
	***REMOVED***
	f.startWrite(FrameContinuation, flags, streamID)
	f.wbuf = append(f.wbuf, headerBlockFragment...)
	return f.endWrite()
***REMOVED***

// A PushPromiseFrame is used to initiate a server stream.
// See http://http2.github.io/http2-spec/#rfc.section.6.6
type PushPromiseFrame struct ***REMOVED***
	FrameHeader
	PromiseID     uint32
	headerFragBuf []byte // not owned
***REMOVED***

func (f *PushPromiseFrame) HeaderBlockFragment() []byte ***REMOVED***
	f.checkValid()
	return f.headerFragBuf
***REMOVED***

func (f *PushPromiseFrame) HeadersEnded() bool ***REMOVED***
	return f.FrameHeader.Flags.Has(FlagPushPromiseEndHeaders)
***REMOVED***

func parsePushPromise(_ *frameCache, fh FrameHeader, countError func(string), p []byte) (_ Frame, err error) ***REMOVED***
	pp := &PushPromiseFrame***REMOVED***
		FrameHeader: fh,
	***REMOVED***
	if pp.StreamID == 0 ***REMOVED***
		// PUSH_PROMISE frames MUST be associated with an existing,
		// peer-initiated stream. The stream identifier of a
		// PUSH_PROMISE frame indicates the stream it is associated
		// with. If the stream identifier field specifies the value
		// 0x0, a recipient MUST respond with a connection error
		// (Section 5.4.1) of type PROTOCOL_ERROR.
		countError("frame_pushpromise_zero_stream")
		return nil, ConnectionError(ErrCodeProtocol)
	***REMOVED***
	// The PUSH_PROMISE frame includes optional padding.
	// Padding fields and flags are identical to those defined for DATA frames
	var padLength uint8
	if fh.Flags.Has(FlagPushPromisePadded) ***REMOVED***
		if p, padLength, err = readByte(p); err != nil ***REMOVED***
			countError("frame_pushpromise_pad_short")
			return
		***REMOVED***
	***REMOVED***

	p, pp.PromiseID, err = readUint32(p)
	if err != nil ***REMOVED***
		countError("frame_pushpromise_promiseid_short")
		return
	***REMOVED***
	pp.PromiseID = pp.PromiseID & (1<<31 - 1)

	if int(padLength) > len(p) ***REMOVED***
		// like the DATA frame, error out if padding is longer than the body.
		countError("frame_pushpromise_pad_too_big")
		return nil, ConnectionError(ErrCodeProtocol)
	***REMOVED***
	pp.headerFragBuf = p[:len(p)-int(padLength)]
	return pp, nil
***REMOVED***

// PushPromiseParam are the parameters for writing a PUSH_PROMISE frame.
type PushPromiseParam struct ***REMOVED***
	// StreamID is the required Stream ID to initiate.
	StreamID uint32

	// PromiseID is the required Stream ID which this
	// Push Promises
	PromiseID uint32

	// BlockFragment is part (or all) of a Header Block.
	BlockFragment []byte

	// EndHeaders indicates that this frame contains an entire
	// header block and is not followed by any
	// CONTINUATION frames.
	EndHeaders bool

	// PadLength is the optional number of bytes of zeros to add
	// to this frame.
	PadLength uint8
***REMOVED***

// WritePushPromise writes a single PushPromise Frame.
//
// As with Header Frames, This is the low level call for writing
// individual frames. Continuation frames are handled elsewhere.
//
// It will perform exactly one Write to the underlying Writer.
// It is the caller's responsibility to not call other Write methods concurrently.
func (f *Framer) WritePushPromise(p PushPromiseParam) error ***REMOVED***
	if !validStreamID(p.StreamID) && !f.AllowIllegalWrites ***REMOVED***
		return errStreamID
	***REMOVED***
	var flags Flags
	if p.PadLength != 0 ***REMOVED***
		flags |= FlagPushPromisePadded
	***REMOVED***
	if p.EndHeaders ***REMOVED***
		flags |= FlagPushPromiseEndHeaders
	***REMOVED***
	f.startWrite(FramePushPromise, flags, p.StreamID)
	if p.PadLength != 0 ***REMOVED***
		f.writeByte(p.PadLength)
	***REMOVED***
	if !validStreamID(p.PromiseID) && !f.AllowIllegalWrites ***REMOVED***
		return errStreamID
	***REMOVED***
	f.writeUint32(p.PromiseID)
	f.wbuf = append(f.wbuf, p.BlockFragment...)
	f.wbuf = append(f.wbuf, padZeros[:p.PadLength]...)
	return f.endWrite()
***REMOVED***

// WriteRawFrame writes a raw frame. This can be used to write
// extension frames unknown to this package.
func (f *Framer) WriteRawFrame(t FrameType, flags Flags, streamID uint32, payload []byte) error ***REMOVED***
	f.startWrite(t, flags, streamID)
	f.writeBytes(payload)
	return f.endWrite()
***REMOVED***

func readByte(p []byte) (remain []byte, b byte, err error) ***REMOVED***
	if len(p) == 0 ***REMOVED***
		return nil, 0, io.ErrUnexpectedEOF
	***REMOVED***
	return p[1:], p[0], nil
***REMOVED***

func readUint32(p []byte) (remain []byte, v uint32, err error) ***REMOVED***
	if len(p) < 4 ***REMOVED***
		return nil, 0, io.ErrUnexpectedEOF
	***REMOVED***
	return p[4:], binary.BigEndian.Uint32(p[:4]), nil
***REMOVED***

type streamEnder interface ***REMOVED***
	StreamEnded() bool
***REMOVED***

type headersEnder interface ***REMOVED***
	HeadersEnded() bool
***REMOVED***

type headersOrContinuation interface ***REMOVED***
	headersEnder
	HeaderBlockFragment() []byte
***REMOVED***

// A MetaHeadersFrame is the representation of one HEADERS frame and
// zero or more contiguous CONTINUATION frames and the decoding of
// their HPACK-encoded contents.
//
// This type of frame does not appear on the wire and is only returned
// by the Framer when Framer.ReadMetaHeaders is set.
type MetaHeadersFrame struct ***REMOVED***
	*HeadersFrame

	// Fields are the fields contained in the HEADERS and
	// CONTINUATION frames. The underlying slice is owned by the
	// Framer and must not be retained after the next call to
	// ReadFrame.
	//
	// Fields are guaranteed to be in the correct http2 order and
	// not have unknown pseudo header fields or invalid header
	// field names or values. Required pseudo header fields may be
	// missing, however. Use the MetaHeadersFrame.Pseudo accessor
	// method access pseudo headers.
	Fields []hpack.HeaderField

	// Truncated is whether the max header list size limit was hit
	// and Fields is incomplete. The hpack decoder state is still
	// valid, however.
	Truncated bool
***REMOVED***

// PseudoValue returns the given pseudo header field's value.
// The provided pseudo field should not contain the leading colon.
func (mh *MetaHeadersFrame) PseudoValue(pseudo string) string ***REMOVED***
	for _, hf := range mh.Fields ***REMOVED***
		if !hf.IsPseudo() ***REMOVED***
			return ""
		***REMOVED***
		if hf.Name[1:] == pseudo ***REMOVED***
			return hf.Value
		***REMOVED***
	***REMOVED***
	return ""
***REMOVED***

// RegularFields returns the regular (non-pseudo) header fields of mh.
// The caller does not own the returned slice.
func (mh *MetaHeadersFrame) RegularFields() []hpack.HeaderField ***REMOVED***
	for i, hf := range mh.Fields ***REMOVED***
		if !hf.IsPseudo() ***REMOVED***
			return mh.Fields[i:]
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// PseudoFields returns the pseudo header fields of mh.
// The caller does not own the returned slice.
func (mh *MetaHeadersFrame) PseudoFields() []hpack.HeaderField ***REMOVED***
	for i, hf := range mh.Fields ***REMOVED***
		if !hf.IsPseudo() ***REMOVED***
			return mh.Fields[:i]
		***REMOVED***
	***REMOVED***
	return mh.Fields
***REMOVED***

func (mh *MetaHeadersFrame) checkPseudos() error ***REMOVED***
	var isRequest, isResponse bool
	pf := mh.PseudoFields()
	for i, hf := range pf ***REMOVED***
		switch hf.Name ***REMOVED***
		case ":method", ":path", ":scheme", ":authority":
			isRequest = true
		case ":status":
			isResponse = true
		default:
			return pseudoHeaderError(hf.Name)
		***REMOVED***
		// Check for duplicates.
		// This would be a bad algorithm, but N is 4.
		// And this doesn't allocate.
		for _, hf2 := range pf[:i] ***REMOVED***
			if hf.Name == hf2.Name ***REMOVED***
				return duplicatePseudoHeaderError(hf.Name)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if isRequest && isResponse ***REMOVED***
		return errMixPseudoHeaderTypes
	***REMOVED***
	return nil
***REMOVED***

func (fr *Framer) maxHeaderStringLen() int ***REMOVED***
	v := fr.maxHeaderListSize()
	if uint32(int(v)) == v ***REMOVED***
		return int(v)
	***REMOVED***
	// They had a crazy big number for MaxHeaderBytes anyway,
	// so give them unlimited header lengths:
	return 0
***REMOVED***

// readMetaFrame returns 0 or more CONTINUATION frames from fr and
// merge them into the provided hf and returns a MetaHeadersFrame
// with the decoded hpack values.
func (fr *Framer) readMetaFrame(hf *HeadersFrame) (*MetaHeadersFrame, error) ***REMOVED***
	if fr.AllowIllegalReads ***REMOVED***
		return nil, errors.New("illegal use of AllowIllegalReads with ReadMetaHeaders")
	***REMOVED***
	mh := &MetaHeadersFrame***REMOVED***
		HeadersFrame: hf,
	***REMOVED***
	var remainSize = fr.maxHeaderListSize()
	var sawRegular bool

	var invalid error // pseudo header field errors
	hdec := fr.ReadMetaHeaders
	hdec.SetEmitEnabled(true)
	hdec.SetMaxStringLength(fr.maxHeaderStringLen())
	hdec.SetEmitFunc(func(hf hpack.HeaderField) ***REMOVED***
		if VerboseLogs && fr.logReads ***REMOVED***
			fr.debugReadLoggerf("http2: decoded hpack field %+v", hf)
		***REMOVED***
		if !httpguts.ValidHeaderFieldValue(hf.Value) ***REMOVED***
			// Don't include the value in the error, because it may be sensitive.
			invalid = headerFieldValueError(hf.Name)
		***REMOVED***
		isPseudo := strings.HasPrefix(hf.Name, ":")
		if isPseudo ***REMOVED***
			if sawRegular ***REMOVED***
				invalid = errPseudoAfterRegular
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			sawRegular = true
			if !validWireHeaderFieldName(hf.Name) ***REMOVED***
				invalid = headerFieldNameError(hf.Name)
			***REMOVED***
		***REMOVED***

		if invalid != nil ***REMOVED***
			hdec.SetEmitEnabled(false)
			return
		***REMOVED***

		size := hf.Size()
		if size > remainSize ***REMOVED***
			hdec.SetEmitEnabled(false)
			mh.Truncated = true
			return
		***REMOVED***
		remainSize -= size

		mh.Fields = append(mh.Fields, hf)
	***REMOVED***)
	// Lose reference to MetaHeadersFrame:
	defer hdec.SetEmitFunc(func(hf hpack.HeaderField) ***REMOVED******REMOVED***)

	var hc headersOrContinuation = hf
	for ***REMOVED***
		frag := hc.HeaderBlockFragment()
		if _, err := hdec.Write(frag); err != nil ***REMOVED***
			return nil, ConnectionError(ErrCodeCompression)
		***REMOVED***

		if hc.HeadersEnded() ***REMOVED***
			break
		***REMOVED***
		if f, err := fr.ReadFrame(); err != nil ***REMOVED***
			return nil, err
		***REMOVED*** else ***REMOVED***
			hc = f.(*ContinuationFrame) // guaranteed by checkFrameOrder
		***REMOVED***
	***REMOVED***

	mh.HeadersFrame.headerFragBuf = nil
	mh.HeadersFrame.invalidate()

	if err := hdec.Close(); err != nil ***REMOVED***
		return nil, ConnectionError(ErrCodeCompression)
	***REMOVED***
	if invalid != nil ***REMOVED***
		fr.errDetail = invalid
		if VerboseLogs ***REMOVED***
			log.Printf("http2: invalid header: %v", invalid)
		***REMOVED***
		return nil, StreamError***REMOVED***mh.StreamID, ErrCodeProtocol, invalid***REMOVED***
	***REMOVED***
	if err := mh.checkPseudos(); err != nil ***REMOVED***
		fr.errDetail = err
		if VerboseLogs ***REMOVED***
			log.Printf("http2: invalid pseudo headers: %v", err)
		***REMOVED***
		return nil, StreamError***REMOVED***mh.StreamID, ErrCodeProtocol, err***REMOVED***
	***REMOVED***
	return mh, nil
***REMOVED***

func summarizeFrame(f Frame) string ***REMOVED***
	var buf bytes.Buffer
	f.Header().writeDebug(&buf)
	switch f := f.(type) ***REMOVED***
	case *SettingsFrame:
		n := 0
		f.ForeachSetting(func(s Setting) error ***REMOVED***
			n++
			if n == 1 ***REMOVED***
				buf.WriteString(", settings:")
			***REMOVED***
			fmt.Fprintf(&buf, " %v=%v,", s.ID, s.Val)
			return nil
		***REMOVED***)
		if n > 0 ***REMOVED***
			buf.Truncate(buf.Len() - 1) // remove trailing comma
		***REMOVED***
	case *DataFrame:
		data := f.Data()
		const max = 256
		if len(data) > max ***REMOVED***
			data = data[:max]
		***REMOVED***
		fmt.Fprintf(&buf, " data=%q", data)
		if len(f.Data()) > max ***REMOVED***
			fmt.Fprintf(&buf, " (%d bytes omitted)", len(f.Data())-max)
		***REMOVED***
	case *WindowUpdateFrame:
		if f.StreamID == 0 ***REMOVED***
			buf.WriteString(" (conn)")
		***REMOVED***
		fmt.Fprintf(&buf, " incr=%v", f.Increment)
	case *PingFrame:
		fmt.Fprintf(&buf, " ping=%q", f.Data[:])
	case *GoAwayFrame:
		fmt.Fprintf(&buf, " LastStreamID=%v ErrCode=%v Debug=%q",
			f.LastStreamID, f.ErrCode, f.debugData)
	case *RSTStreamFrame:
		fmt.Fprintf(&buf, " ErrCode=%v", f.ErrCode)
	***REMOVED***
	return buf.String()
***REMOVED***
