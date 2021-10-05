// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http2

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/net/http/httpguts"
	"golang.org/x/net/http2/hpack"
)

// writeFramer is implemented by any type that is used to write frames.
type writeFramer interface ***REMOVED***
	writeFrame(writeContext) error

	// staysWithinBuffer reports whether this writer promises that
	// it will only write less than or equal to size bytes, and it
	// won't Flush the write context.
	staysWithinBuffer(size int) bool
***REMOVED***

// writeContext is the interface needed by the various frame writer
// types below. All the writeFrame methods below are scheduled via the
// frame writing scheduler (see writeScheduler in writesched.go).
//
// This interface is implemented by *serverConn.
//
// TODO: decide whether to a) use this in the client code (which didn't
// end up using this yet, because it has a simpler design, not
// currently implementing priorities), or b) delete this and
// make the server code a bit more concrete.
type writeContext interface ***REMOVED***
	Framer() *Framer
	Flush() error
	CloseConn() error
	// HeaderEncoder returns an HPACK encoder that writes to the
	// returned buffer.
	HeaderEncoder() (*hpack.Encoder, *bytes.Buffer)
***REMOVED***

// writeEndsStream reports whether w writes a frame that will transition
// the stream to a half-closed local state. This returns false for RST_STREAM,
// which closes the entire stream (not just the local half).
func writeEndsStream(w writeFramer) bool ***REMOVED***
	switch v := w.(type) ***REMOVED***
	case *writeData:
		return v.endStream
	case *writeResHeaders:
		return v.endStream
	case nil:
		// This can only happen if the caller reuses w after it's
		// been intentionally nil'ed out to prevent use. Keep this
		// here to catch future refactoring breaking it.
		panic("writeEndsStream called on nil writeFramer")
	***REMOVED***
	return false
***REMOVED***

type flushFrameWriter struct***REMOVED******REMOVED***

func (flushFrameWriter) writeFrame(ctx writeContext) error ***REMOVED***
	return ctx.Flush()
***REMOVED***

func (flushFrameWriter) staysWithinBuffer(max int) bool ***REMOVED*** return false ***REMOVED***

type writeSettings []Setting

func (s writeSettings) staysWithinBuffer(max int) bool ***REMOVED***
	const settingSize = 6 // uint16 + uint32
	return frameHeaderLen+settingSize*len(s) <= max

***REMOVED***

func (s writeSettings) writeFrame(ctx writeContext) error ***REMOVED***
	return ctx.Framer().WriteSettings([]Setting(s)...)
***REMOVED***

type writeGoAway struct ***REMOVED***
	maxStreamID uint32
	code        ErrCode
***REMOVED***

func (p *writeGoAway) writeFrame(ctx writeContext) error ***REMOVED***
	err := ctx.Framer().WriteGoAway(p.maxStreamID, p.code, nil)
	ctx.Flush() // ignore error: we're hanging up on them anyway
	return err
***REMOVED***

func (*writeGoAway) staysWithinBuffer(max int) bool ***REMOVED*** return false ***REMOVED*** // flushes

type writeData struct ***REMOVED***
	streamID  uint32
	p         []byte
	endStream bool
***REMOVED***

func (w *writeData) String() string ***REMOVED***
	return fmt.Sprintf("writeData(stream=%d, p=%d, endStream=%v)", w.streamID, len(w.p), w.endStream)
***REMOVED***

func (w *writeData) writeFrame(ctx writeContext) error ***REMOVED***
	return ctx.Framer().WriteData(w.streamID, w.endStream, w.p)
***REMOVED***

func (w *writeData) staysWithinBuffer(max int) bool ***REMOVED***
	return frameHeaderLen+len(w.p) <= max
***REMOVED***

// handlerPanicRST is the message sent from handler goroutines when
// the handler panics.
type handlerPanicRST struct ***REMOVED***
	StreamID uint32
***REMOVED***

func (hp handlerPanicRST) writeFrame(ctx writeContext) error ***REMOVED***
	return ctx.Framer().WriteRSTStream(hp.StreamID, ErrCodeInternal)
***REMOVED***

func (hp handlerPanicRST) staysWithinBuffer(max int) bool ***REMOVED*** return frameHeaderLen+4 <= max ***REMOVED***

func (se StreamError) writeFrame(ctx writeContext) error ***REMOVED***
	return ctx.Framer().WriteRSTStream(se.StreamID, se.Code)
***REMOVED***

func (se StreamError) staysWithinBuffer(max int) bool ***REMOVED*** return frameHeaderLen+4 <= max ***REMOVED***

type writePingAck struct***REMOVED*** pf *PingFrame ***REMOVED***

func (w writePingAck) writeFrame(ctx writeContext) error ***REMOVED***
	return ctx.Framer().WritePing(true, w.pf.Data)
***REMOVED***

func (w writePingAck) staysWithinBuffer(max int) bool ***REMOVED*** return frameHeaderLen+len(w.pf.Data) <= max ***REMOVED***

type writeSettingsAck struct***REMOVED******REMOVED***

func (writeSettingsAck) writeFrame(ctx writeContext) error ***REMOVED***
	return ctx.Framer().WriteSettingsAck()
***REMOVED***

func (writeSettingsAck) staysWithinBuffer(max int) bool ***REMOVED*** return frameHeaderLen <= max ***REMOVED***

// splitHeaderBlock splits headerBlock into fragments so that each fragment fits
// in a single frame, then calls fn for each fragment. firstFrag/lastFrag are true
// for the first/last fragment, respectively.
func splitHeaderBlock(ctx writeContext, headerBlock []byte, fn func(ctx writeContext, frag []byte, firstFrag, lastFrag bool) error) error ***REMOVED***
	// For now we're lazy and just pick the minimum MAX_FRAME_SIZE
	// that all peers must support (16KB). Later we could care
	// more and send larger frames if the peer advertised it, but
	// there's little point. Most headers are small anyway (so we
	// generally won't have CONTINUATION frames), and extra frames
	// only waste 9 bytes anyway.
	const maxFrameSize = 16384

	first := true
	for len(headerBlock) > 0 ***REMOVED***
		frag := headerBlock
		if len(frag) > maxFrameSize ***REMOVED***
			frag = frag[:maxFrameSize]
		***REMOVED***
		headerBlock = headerBlock[len(frag):]
		if err := fn(ctx, frag, first, len(headerBlock) == 0); err != nil ***REMOVED***
			return err
		***REMOVED***
		first = false
	***REMOVED***
	return nil
***REMOVED***

// writeResHeaders is a request to write a HEADERS and 0+ CONTINUATION frames
// for HTTP response headers or trailers from a server handler.
type writeResHeaders struct ***REMOVED***
	streamID    uint32
	httpResCode int         // 0 means no ":status" line
	h           http.Header // may be nil
	trailers    []string    // if non-nil, which keys of h to write. nil means all.
	endStream   bool

	date          string
	contentType   string
	contentLength string
***REMOVED***

func encKV(enc *hpack.Encoder, k, v string) ***REMOVED***
	if VerboseLogs ***REMOVED***
		log.Printf("http2: server encoding header %q = %q", k, v)
	***REMOVED***
	enc.WriteField(hpack.HeaderField***REMOVED***Name: k, Value: v***REMOVED***)
***REMOVED***

func (w *writeResHeaders) staysWithinBuffer(max int) bool ***REMOVED***
	// TODO: this is a common one. It'd be nice to return true
	// here and get into the fast path if we could be clever and
	// calculate the size fast enough, or at least a conservative
	// upper bound that usually fires. (Maybe if w.h and
	// w.trailers are nil, so we don't need to enumerate it.)
	// Otherwise I'm afraid that just calculating the length to
	// answer this question would be slower than the ~2µs benefit.
	return false
***REMOVED***

func (w *writeResHeaders) writeFrame(ctx writeContext) error ***REMOVED***
	enc, buf := ctx.HeaderEncoder()
	buf.Reset()

	if w.httpResCode != 0 ***REMOVED***
		encKV(enc, ":status", httpCodeString(w.httpResCode))
	***REMOVED***

	encodeHeaders(enc, w.h, w.trailers)

	if w.contentType != "" ***REMOVED***
		encKV(enc, "content-type", w.contentType)
	***REMOVED***
	if w.contentLength != "" ***REMOVED***
		encKV(enc, "content-length", w.contentLength)
	***REMOVED***
	if w.date != "" ***REMOVED***
		encKV(enc, "date", w.date)
	***REMOVED***

	headerBlock := buf.Bytes()
	if len(headerBlock) == 0 && w.trailers == nil ***REMOVED***
		panic("unexpected empty hpack")
	***REMOVED***

	return splitHeaderBlock(ctx, headerBlock, w.writeHeaderBlock)
***REMOVED***

func (w *writeResHeaders) writeHeaderBlock(ctx writeContext, frag []byte, firstFrag, lastFrag bool) error ***REMOVED***
	if firstFrag ***REMOVED***
		return ctx.Framer().WriteHeaders(HeadersFrameParam***REMOVED***
			StreamID:      w.streamID,
			BlockFragment: frag,
			EndStream:     w.endStream,
			EndHeaders:    lastFrag,
		***REMOVED***)
	***REMOVED*** else ***REMOVED***
		return ctx.Framer().WriteContinuation(w.streamID, lastFrag, frag)
	***REMOVED***
***REMOVED***

// writePushPromise is a request to write a PUSH_PROMISE and 0+ CONTINUATION frames.
type writePushPromise struct ***REMOVED***
	streamID uint32   // pusher stream
	method   string   // for :method
	url      *url.URL // for :scheme, :authority, :path
	h        http.Header

	// Creates an ID for a pushed stream. This runs on serveG just before
	// the frame is written. The returned ID is copied to promisedID.
	allocatePromisedID func() (uint32, error)
	promisedID         uint32
***REMOVED***

func (w *writePushPromise) staysWithinBuffer(max int) bool ***REMOVED***
	// TODO: see writeResHeaders.staysWithinBuffer
	return false
***REMOVED***

func (w *writePushPromise) writeFrame(ctx writeContext) error ***REMOVED***
	enc, buf := ctx.HeaderEncoder()
	buf.Reset()

	encKV(enc, ":method", w.method)
	encKV(enc, ":scheme", w.url.Scheme)
	encKV(enc, ":authority", w.url.Host)
	encKV(enc, ":path", w.url.RequestURI())
	encodeHeaders(enc, w.h, nil)

	headerBlock := buf.Bytes()
	if len(headerBlock) == 0 ***REMOVED***
		panic("unexpected empty hpack")
	***REMOVED***

	return splitHeaderBlock(ctx, headerBlock, w.writeHeaderBlock)
***REMOVED***

func (w *writePushPromise) writeHeaderBlock(ctx writeContext, frag []byte, firstFrag, lastFrag bool) error ***REMOVED***
	if firstFrag ***REMOVED***
		return ctx.Framer().WritePushPromise(PushPromiseParam***REMOVED***
			StreamID:      w.streamID,
			PromiseID:     w.promisedID,
			BlockFragment: frag,
			EndHeaders:    lastFrag,
		***REMOVED***)
	***REMOVED*** else ***REMOVED***
		return ctx.Framer().WriteContinuation(w.streamID, lastFrag, frag)
	***REMOVED***
***REMOVED***

type write100ContinueHeadersFrame struct ***REMOVED***
	streamID uint32
***REMOVED***

func (w write100ContinueHeadersFrame) writeFrame(ctx writeContext) error ***REMOVED***
	enc, buf := ctx.HeaderEncoder()
	buf.Reset()
	encKV(enc, ":status", "100")
	return ctx.Framer().WriteHeaders(HeadersFrameParam***REMOVED***
		StreamID:      w.streamID,
		BlockFragment: buf.Bytes(),
		EndStream:     false,
		EndHeaders:    true,
	***REMOVED***)
***REMOVED***

func (w write100ContinueHeadersFrame) staysWithinBuffer(max int) bool ***REMOVED***
	// Sloppy but conservative:
	return 9+2*(len(":status")+len("100")) <= max
***REMOVED***

type writeWindowUpdate struct ***REMOVED***
	streamID uint32 // or 0 for conn-level
	n        uint32
***REMOVED***

func (wu writeWindowUpdate) staysWithinBuffer(max int) bool ***REMOVED*** return frameHeaderLen+4 <= max ***REMOVED***

func (wu writeWindowUpdate) writeFrame(ctx writeContext) error ***REMOVED***
	return ctx.Framer().WriteWindowUpdate(wu.streamID, wu.n)
***REMOVED***

// encodeHeaders encodes an http.Header. If keys is not nil, then (k, h[k])
// is encoded only if k is in keys.
func encodeHeaders(enc *hpack.Encoder, h http.Header, keys []string) ***REMOVED***
	if keys == nil ***REMOVED***
		sorter := sorterPool.Get().(*sorter)
		// Using defer here, since the returned keys from the
		// sorter.Keys method is only valid until the sorter
		// is returned:
		defer sorterPool.Put(sorter)
		keys = sorter.Keys(h)
	***REMOVED***
	for _, k := range keys ***REMOVED***
		vv := h[k]
		k, ascii := lowerHeader(k)
		if !ascii ***REMOVED***
			// Skip writing invalid headers. Per RFC 7540, Section 8.1.2, header
			// field names have to be ASCII characters (just as in HTTP/1.x).
			continue
		***REMOVED***
		if !validWireHeaderFieldName(k) ***REMOVED***
			// Skip it as backup paranoia. Per
			// golang.org/issue/14048, these should
			// already be rejected at a higher level.
			continue
		***REMOVED***
		isTE := k == "transfer-encoding"
		for _, v := range vv ***REMOVED***
			if !httpguts.ValidHeaderFieldValue(v) ***REMOVED***
				// TODO: return an error? golang.org/issue/14048
				// For now just omit it.
				continue
			***REMOVED***
			// TODO: more of "8.1.2.2 Connection-Specific Header Fields"
			if isTE && v != "trailers" ***REMOVED***
				continue
			***REMOVED***
			encKV(enc, k, v)
		***REMOVED***
	***REMOVED***
***REMOVED***
