// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http2

import "fmt"

// WriteScheduler is the interface implemented by HTTP/2 write schedulers.
// Methods are never called concurrently.
type WriteScheduler interface ***REMOVED***
	// OpenStream opens a new stream in the write scheduler.
	// It is illegal to call this with streamID=0 or with a streamID that is
	// already open -- the call may panic.
	OpenStream(streamID uint32, options OpenStreamOptions)

	// CloseStream closes a stream in the write scheduler. Any frames queued on
	// this stream should be discarded. It is illegal to call this on a stream
	// that is not open -- the call may panic.
	CloseStream(streamID uint32)

	// AdjustStream adjusts the priority of the given stream. This may be called
	// on a stream that has not yet been opened or has been closed. Note that
	// RFC 7540 allows PRIORITY frames to be sent on streams in any state. See:
	// https://tools.ietf.org/html/rfc7540#section-5.1
	AdjustStream(streamID uint32, priority PriorityParam)

	// Push queues a frame in the scheduler. In most cases, this will not be
	// called with wr.StreamID()!=0 unless that stream is currently open. The one
	// exception is RST_STREAM frames, which may be sent on idle or closed streams.
	Push(wr FrameWriteRequest)

	// Pop dequeues the next frame to write. Returns false if no frames can
	// be written. Frames with a given wr.StreamID() are Pop'd in the same
	// order they are Push'd, except RST_STREAM frames. No frames should be
	// discarded except by CloseStream.
	Pop() (wr FrameWriteRequest, ok bool)
***REMOVED***

// OpenStreamOptions specifies extra options for WriteScheduler.OpenStream.
type OpenStreamOptions struct ***REMOVED***
	// PusherID is zero if the stream was initiated by the client. Otherwise,
	// PusherID names the stream that pushed the newly opened stream.
	PusherID uint32
***REMOVED***

// FrameWriteRequest is a request to write a frame.
type FrameWriteRequest struct ***REMOVED***
	// write is the interface value that does the writing, once the
	// WriteScheduler has selected this frame to write. The write
	// functions are all defined in write.go.
	write writeFramer

	// stream is the stream on which this frame will be written.
	// nil for non-stream frames like PING and SETTINGS.
	// nil for RST_STREAM streams, which use the StreamError.StreamID field instead.
	stream *stream

	// done, if non-nil, must be a buffered channel with space for
	// 1 message and is sent the return value from write (or an
	// earlier error) when the frame has been written.
	done chan error
***REMOVED***

// StreamID returns the id of the stream this frame will be written to.
// 0 is used for non-stream frames such as PING and SETTINGS.
func (wr FrameWriteRequest) StreamID() uint32 ***REMOVED***
	if wr.stream == nil ***REMOVED***
		if se, ok := wr.write.(StreamError); ok ***REMOVED***
			// (*serverConn).resetStream doesn't set
			// stream because it doesn't necessarily have
			// one. So special case this type of write
			// message.
			return se.StreamID
		***REMOVED***
		return 0
	***REMOVED***
	return wr.stream.id
***REMOVED***

// isControl reports whether wr is a control frame for MaxQueuedControlFrames
// purposes. That includes non-stream frames and RST_STREAM frames.
func (wr FrameWriteRequest) isControl() bool ***REMOVED***
	return wr.stream == nil
***REMOVED***

// DataSize returns the number of flow control bytes that must be consumed
// to write this entire frame. This is 0 for non-DATA frames.
func (wr FrameWriteRequest) DataSize() int ***REMOVED***
	if wd, ok := wr.write.(*writeData); ok ***REMOVED***
		return len(wd.p)
	***REMOVED***
	return 0
***REMOVED***

// Consume consumes min(n, available) bytes from this frame, where available
// is the number of flow control bytes available on the stream. Consume returns
// 0, 1, or 2 frames, where the integer return value gives the number of frames
// returned.
//
// If flow control prevents consuming any bytes, this returns (_, _, 0). If
// the entire frame was consumed, this returns (wr, _, 1). Otherwise, this
// returns (consumed, rest, 2), where 'consumed' contains the consumed bytes and
// 'rest' contains the remaining bytes. The consumed bytes are deducted from the
// underlying stream's flow control budget.
func (wr FrameWriteRequest) Consume(n int32) (FrameWriteRequest, FrameWriteRequest, int) ***REMOVED***
	var empty FrameWriteRequest

	// Non-DATA frames are always consumed whole.
	wd, ok := wr.write.(*writeData)
	if !ok || len(wd.p) == 0 ***REMOVED***
		return wr, empty, 1
	***REMOVED***

	// Might need to split after applying limits.
	allowed := wr.stream.flow.available()
	if n < allowed ***REMOVED***
		allowed = n
	***REMOVED***
	if wr.stream.sc.maxFrameSize < allowed ***REMOVED***
		allowed = wr.stream.sc.maxFrameSize
	***REMOVED***
	if allowed <= 0 ***REMOVED***
		return empty, empty, 0
	***REMOVED***
	if len(wd.p) > int(allowed) ***REMOVED***
		wr.stream.flow.take(allowed)
		consumed := FrameWriteRequest***REMOVED***
			stream: wr.stream,
			write: &writeData***REMOVED***
				streamID: wd.streamID,
				p:        wd.p[:allowed],
				// Even if the original had endStream set, there
				// are bytes remaining because len(wd.p) > allowed,
				// so we know endStream is false.
				endStream: false,
			***REMOVED***,
			// Our caller is blocking on the final DATA frame, not
			// this intermediate frame, so no need to wait.
			done: nil,
		***REMOVED***
		rest := FrameWriteRequest***REMOVED***
			stream: wr.stream,
			write: &writeData***REMOVED***
				streamID:  wd.streamID,
				p:         wd.p[allowed:],
				endStream: wd.endStream,
			***REMOVED***,
			done: wr.done,
		***REMOVED***
		return consumed, rest, 2
	***REMOVED***

	// The frame is consumed whole.
	// NB: This cast cannot overflow because allowed is <= math.MaxInt32.
	wr.stream.flow.take(int32(len(wd.p)))
	return wr, empty, 1
***REMOVED***

// String is for debugging only.
func (wr FrameWriteRequest) String() string ***REMOVED***
	var des string
	if s, ok := wr.write.(fmt.Stringer); ok ***REMOVED***
		des = s.String()
	***REMOVED*** else ***REMOVED***
		des = fmt.Sprintf("%T", wr.write)
	***REMOVED***
	return fmt.Sprintf("[FrameWriteRequest stream=%d, ch=%v, writer=%v]", wr.StreamID(), wr.done != nil, des)
***REMOVED***

// replyToWriter sends err to wr.done and panics if the send must block
// This does nothing if wr.done is nil.
func (wr *FrameWriteRequest) replyToWriter(err error) ***REMOVED***
	if wr.done == nil ***REMOVED***
		return
	***REMOVED***
	select ***REMOVED***
	case wr.done <- err:
	default:
		panic(fmt.Sprintf("unbuffered done channel passed in for type %T", wr.write))
	***REMOVED***
	wr.write = nil // prevent use (assume it's tainted after wr.done send)
***REMOVED***

// writeQueue is used by implementations of WriteScheduler.
type writeQueue struct ***REMOVED***
	s []FrameWriteRequest
***REMOVED***

func (q *writeQueue) empty() bool ***REMOVED*** return len(q.s) == 0 ***REMOVED***

func (q *writeQueue) push(wr FrameWriteRequest) ***REMOVED***
	q.s = append(q.s, wr)
***REMOVED***

func (q *writeQueue) shift() FrameWriteRequest ***REMOVED***
	if len(q.s) == 0 ***REMOVED***
		panic("invalid use of queue")
	***REMOVED***
	wr := q.s[0]
	// TODO: less copy-happy queue.
	copy(q.s, q.s[1:])
	q.s[len(q.s)-1] = FrameWriteRequest***REMOVED******REMOVED***
	q.s = q.s[:len(q.s)-1]
	return wr
***REMOVED***

// consume consumes up to n bytes from q.s[0]. If the frame is
// entirely consumed, it is removed from the queue. If the frame
// is partially consumed, the frame is kept with the consumed
// bytes removed. Returns true iff any bytes were consumed.
func (q *writeQueue) consume(n int32) (FrameWriteRequest, bool) ***REMOVED***
	if len(q.s) == 0 ***REMOVED***
		return FrameWriteRequest***REMOVED******REMOVED***, false
	***REMOVED***
	consumed, rest, numresult := q.s[0].Consume(n)
	switch numresult ***REMOVED***
	case 0:
		return FrameWriteRequest***REMOVED******REMOVED***, false
	case 1:
		q.shift()
	case 2:
		q.s[0] = rest
	***REMOVED***
	return consumed, true
***REMOVED***

type writeQueuePool []*writeQueue

// put inserts an unused writeQueue into the pool.
func (p *writeQueuePool) put(q *writeQueue) ***REMOVED***
	for i := range q.s ***REMOVED***
		q.s[i] = FrameWriteRequest***REMOVED******REMOVED***
	***REMOVED***
	q.s = q.s[:0]
	*p = append(*p, q)
***REMOVED***

// get returns an empty writeQueue.
func (p *writeQueuePool) get() *writeQueue ***REMOVED***
	ln := len(*p)
	if ln == 0 ***REMOVED***
		return new(writeQueue)
	***REMOVED***
	x := ln - 1
	q := (*p)[x]
	(*p)[x] = nil
	*p = (*p)[:x]
	return q
***REMOVED***
