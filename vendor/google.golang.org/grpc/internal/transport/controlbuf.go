/*
 *
 * Copyright 2014 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package transport

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
	"google.golang.org/grpc/internal/grpcutil"
	"google.golang.org/grpc/status"
)

var updateHeaderTblSize = func(e *hpack.Encoder, v uint32) ***REMOVED***
	e.SetMaxDynamicTableSizeLimit(v)
***REMOVED***

type itemNode struct ***REMOVED***
	it   interface***REMOVED******REMOVED***
	next *itemNode
***REMOVED***

type itemList struct ***REMOVED***
	head *itemNode
	tail *itemNode
***REMOVED***

func (il *itemList) enqueue(i interface***REMOVED******REMOVED***) ***REMOVED***
	n := &itemNode***REMOVED***it: i***REMOVED***
	if il.tail == nil ***REMOVED***
		il.head, il.tail = n, n
		return
	***REMOVED***
	il.tail.next = n
	il.tail = n
***REMOVED***

// peek returns the first item in the list without removing it from the
// list.
func (il *itemList) peek() interface***REMOVED******REMOVED*** ***REMOVED***
	return il.head.it
***REMOVED***

func (il *itemList) dequeue() interface***REMOVED******REMOVED*** ***REMOVED***
	if il.head == nil ***REMOVED***
		return nil
	***REMOVED***
	i := il.head.it
	il.head = il.head.next
	if il.head == nil ***REMOVED***
		il.tail = nil
	***REMOVED***
	return i
***REMOVED***

func (il *itemList) dequeueAll() *itemNode ***REMOVED***
	h := il.head
	il.head, il.tail = nil, nil
	return h
***REMOVED***

func (il *itemList) isEmpty() bool ***REMOVED***
	return il.head == nil
***REMOVED***

// The following defines various control items which could flow through
// the control buffer of transport. They represent different aspects of
// control tasks, e.g., flow control, settings, streaming resetting, etc.

// maxQueuedTransportResponseFrames is the most queued "transport response"
// frames we will buffer before preventing new reads from occurring on the
// transport.  These are control frames sent in response to client requests,
// such as RST_STREAM due to bad headers or settings acks.
const maxQueuedTransportResponseFrames = 50

type cbItem interface ***REMOVED***
	isTransportResponseFrame() bool
***REMOVED***

// registerStream is used to register an incoming stream with loopy writer.
type registerStream struct ***REMOVED***
	streamID uint32
	wq       *writeQuota
***REMOVED***

func (*registerStream) isTransportResponseFrame() bool ***REMOVED*** return false ***REMOVED***

// headerFrame is also used to register stream on the client-side.
type headerFrame struct ***REMOVED***
	streamID   uint32
	hf         []hpack.HeaderField
	endStream  bool               // Valid on server side.
	initStream func(uint32) error // Used only on the client side.
	onWrite    func()
	wq         *writeQuota    // write quota for the stream created.
	cleanup    *cleanupStream // Valid on the server side.
	onOrphaned func(error)    // Valid on client-side
***REMOVED***

func (h *headerFrame) isTransportResponseFrame() bool ***REMOVED***
	return h.cleanup != nil && h.cleanup.rst // Results in a RST_STREAM
***REMOVED***

type cleanupStream struct ***REMOVED***
	streamID uint32
	rst      bool
	rstCode  http2.ErrCode
	onWrite  func()
***REMOVED***

func (c *cleanupStream) isTransportResponseFrame() bool ***REMOVED*** return c.rst ***REMOVED*** // Results in a RST_STREAM

type earlyAbortStream struct ***REMOVED***
	httpStatus     uint32
	streamID       uint32
	contentSubtype string
	status         *status.Status
***REMOVED***

func (*earlyAbortStream) isTransportResponseFrame() bool ***REMOVED*** return false ***REMOVED***

type dataFrame struct ***REMOVED***
	streamID  uint32
	endStream bool
	h         []byte
	d         []byte
	// onEachWrite is called every time
	// a part of d is written out.
	onEachWrite func()
***REMOVED***

func (*dataFrame) isTransportResponseFrame() bool ***REMOVED*** return false ***REMOVED***

type incomingWindowUpdate struct ***REMOVED***
	streamID  uint32
	increment uint32
***REMOVED***

func (*incomingWindowUpdate) isTransportResponseFrame() bool ***REMOVED*** return false ***REMOVED***

type outgoingWindowUpdate struct ***REMOVED***
	streamID  uint32
	increment uint32
***REMOVED***

func (*outgoingWindowUpdate) isTransportResponseFrame() bool ***REMOVED***
	return false // window updates are throttled by thresholds
***REMOVED***

type incomingSettings struct ***REMOVED***
	ss []http2.Setting
***REMOVED***

func (*incomingSettings) isTransportResponseFrame() bool ***REMOVED*** return true ***REMOVED*** // Results in a settings ACK

type outgoingSettings struct ***REMOVED***
	ss []http2.Setting
***REMOVED***

func (*outgoingSettings) isTransportResponseFrame() bool ***REMOVED*** return false ***REMOVED***

type incomingGoAway struct ***REMOVED***
***REMOVED***

func (*incomingGoAway) isTransportResponseFrame() bool ***REMOVED*** return false ***REMOVED***

type goAway struct ***REMOVED***
	code      http2.ErrCode
	debugData []byte
	headsUp   bool
	closeConn bool
***REMOVED***

func (*goAway) isTransportResponseFrame() bool ***REMOVED*** return false ***REMOVED***

type ping struct ***REMOVED***
	ack  bool
	data [8]byte
***REMOVED***

func (*ping) isTransportResponseFrame() bool ***REMOVED*** return true ***REMOVED***

type outFlowControlSizeRequest struct ***REMOVED***
	resp chan uint32
***REMOVED***

func (*outFlowControlSizeRequest) isTransportResponseFrame() bool ***REMOVED*** return false ***REMOVED***

type outStreamState int

const (
	active outStreamState = iota
	empty
	waitingOnStreamQuota
)

type outStream struct ***REMOVED***
	id               uint32
	state            outStreamState
	itl              *itemList
	bytesOutStanding int
	wq               *writeQuota

	next *outStream
	prev *outStream
***REMOVED***

func (s *outStream) deleteSelf() ***REMOVED***
	if s.prev != nil ***REMOVED***
		s.prev.next = s.next
	***REMOVED***
	if s.next != nil ***REMOVED***
		s.next.prev = s.prev
	***REMOVED***
	s.next, s.prev = nil, nil
***REMOVED***

type outStreamList struct ***REMOVED***
	// Following are sentinel objects that mark the
	// beginning and end of the list. They do not
	// contain any item lists. All valid objects are
	// inserted in between them.
	// This is needed so that an outStream object can
	// deleteSelf() in O(1) time without knowing which
	// list it belongs to.
	head *outStream
	tail *outStream
***REMOVED***

func newOutStreamList() *outStreamList ***REMOVED***
	head, tail := new(outStream), new(outStream)
	head.next = tail
	tail.prev = head
	return &outStreamList***REMOVED***
		head: head,
		tail: tail,
	***REMOVED***
***REMOVED***

func (l *outStreamList) enqueue(s *outStream) ***REMOVED***
	e := l.tail.prev
	e.next = s
	s.prev = e
	s.next = l.tail
	l.tail.prev = s
***REMOVED***

// remove from the beginning of the list.
func (l *outStreamList) dequeue() *outStream ***REMOVED***
	b := l.head.next
	if b == l.tail ***REMOVED***
		return nil
	***REMOVED***
	b.deleteSelf()
	return b
***REMOVED***

// controlBuffer is a way to pass information to loopy.
// Information is passed as specific struct types called control frames.
// A control frame not only represents data, messages or headers to be sent out
// but can also be used to instruct loopy to update its internal state.
// It shouldn't be confused with an HTTP2 frame, although some of the control frames
// like dataFrame and headerFrame do go out on wire as HTTP2 frames.
type controlBuffer struct ***REMOVED***
	ch              chan struct***REMOVED******REMOVED***
	done            <-chan struct***REMOVED******REMOVED***
	mu              sync.Mutex
	consumerWaiting bool
	list            *itemList
	err             error

	// transportResponseFrames counts the number of queued items that represent
	// the response of an action initiated by the peer.  trfChan is created
	// when transportResponseFrames >= maxQueuedTransportResponseFrames and is
	// closed and nilled when transportResponseFrames drops below the
	// threshold.  Both fields are protected by mu.
	transportResponseFrames int
	trfChan                 atomic.Value // chan struct***REMOVED******REMOVED***
***REMOVED***

func newControlBuffer(done <-chan struct***REMOVED******REMOVED***) *controlBuffer ***REMOVED***
	return &controlBuffer***REMOVED***
		ch:   make(chan struct***REMOVED******REMOVED***, 1),
		list: &itemList***REMOVED******REMOVED***,
		done: done,
	***REMOVED***
***REMOVED***

// throttle blocks if there are too many incomingSettings/cleanupStreams in the
// controlbuf.
func (c *controlBuffer) throttle() ***REMOVED***
	ch, _ := c.trfChan.Load().(chan struct***REMOVED******REMOVED***)
	if ch != nil ***REMOVED***
		select ***REMOVED***
		case <-ch:
		case <-c.done:
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *controlBuffer) put(it cbItem) error ***REMOVED***
	_, err := c.executeAndPut(nil, it)
	return err
***REMOVED***

func (c *controlBuffer) executeAndPut(f func(it interface***REMOVED******REMOVED***) bool, it cbItem) (bool, error) ***REMOVED***
	var wakeUp bool
	c.mu.Lock()
	if c.err != nil ***REMOVED***
		c.mu.Unlock()
		return false, c.err
	***REMOVED***
	if f != nil ***REMOVED***
		if !f(it) ***REMOVED*** // f wasn't successful
			c.mu.Unlock()
			return false, nil
		***REMOVED***
	***REMOVED***
	if c.consumerWaiting ***REMOVED***
		wakeUp = true
		c.consumerWaiting = false
	***REMOVED***
	c.list.enqueue(it)
	if it.isTransportResponseFrame() ***REMOVED***
		c.transportResponseFrames++
		if c.transportResponseFrames == maxQueuedTransportResponseFrames ***REMOVED***
			// We are adding the frame that puts us over the threshold; create
			// a throttling channel.
			c.trfChan.Store(make(chan struct***REMOVED******REMOVED***))
		***REMOVED***
	***REMOVED***
	c.mu.Unlock()
	if wakeUp ***REMOVED***
		select ***REMOVED***
		case c.ch <- struct***REMOVED******REMOVED******REMOVED******REMOVED***:
		default:
		***REMOVED***
	***REMOVED***
	return true, nil
***REMOVED***

// Note argument f should never be nil.
func (c *controlBuffer) execute(f func(it interface***REMOVED******REMOVED***) bool, it interface***REMOVED******REMOVED***) (bool, error) ***REMOVED***
	c.mu.Lock()
	if c.err != nil ***REMOVED***
		c.mu.Unlock()
		return false, c.err
	***REMOVED***
	if !f(it) ***REMOVED*** // f wasn't successful
		c.mu.Unlock()
		return false, nil
	***REMOVED***
	c.mu.Unlock()
	return true, nil
***REMOVED***

func (c *controlBuffer) get(block bool) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	for ***REMOVED***
		c.mu.Lock()
		if c.err != nil ***REMOVED***
			c.mu.Unlock()
			return nil, c.err
		***REMOVED***
		if !c.list.isEmpty() ***REMOVED***
			h := c.list.dequeue().(cbItem)
			if h.isTransportResponseFrame() ***REMOVED***
				if c.transportResponseFrames == maxQueuedTransportResponseFrames ***REMOVED***
					// We are removing the frame that put us over the
					// threshold; close and clear the throttling channel.
					ch := c.trfChan.Load().(chan struct***REMOVED******REMOVED***)
					close(ch)
					c.trfChan.Store((chan struct***REMOVED******REMOVED***)(nil))
				***REMOVED***
				c.transportResponseFrames--
			***REMOVED***
			c.mu.Unlock()
			return h, nil
		***REMOVED***
		if !block ***REMOVED***
			c.mu.Unlock()
			return nil, nil
		***REMOVED***
		c.consumerWaiting = true
		c.mu.Unlock()
		select ***REMOVED***
		case <-c.ch:
		case <-c.done:
			return nil, ErrConnClosing
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *controlBuffer) finish() ***REMOVED***
	c.mu.Lock()
	if c.err != nil ***REMOVED***
		c.mu.Unlock()
		return
	***REMOVED***
	c.err = ErrConnClosing
	// There may be headers for streams in the control buffer.
	// These streams need to be cleaned out since the transport
	// is still not aware of these yet.
	for head := c.list.dequeueAll(); head != nil; head = head.next ***REMOVED***
		hdr, ok := head.it.(*headerFrame)
		if !ok ***REMOVED***
			continue
		***REMOVED***
		if hdr.onOrphaned != nil ***REMOVED*** // It will be nil on the server-side.
			hdr.onOrphaned(ErrConnClosing)
		***REMOVED***
	***REMOVED***
	// In case throttle() is currently in flight, it needs to be unblocked.
	// Otherwise, the transport may not close, since the transport is closed by
	// the reader encountering the connection error.
	ch, _ := c.trfChan.Load().(chan struct***REMOVED******REMOVED***)
	if ch != nil ***REMOVED***
		close(ch)
	***REMOVED***
	c.trfChan.Store((chan struct***REMOVED******REMOVED***)(nil))
	c.mu.Unlock()
***REMOVED***

type side int

const (
	clientSide side = iota
	serverSide
)

// Loopy receives frames from the control buffer.
// Each frame is handled individually; most of the work done by loopy goes
// into handling data frames. Loopy maintains a queue of active streams, and each
// stream maintains a queue of data frames; as loopy receives data frames
// it gets added to the queue of the relevant stream.
// Loopy goes over this list of active streams by processing one node every iteration,
// thereby closely resemebling to a round-robin scheduling over all streams. While
// processing a stream, loopy writes out data bytes from this stream capped by the min
// of http2MaxFrameLen, connection-level flow control and stream-level flow control.
type loopyWriter struct ***REMOVED***
	side      side
	cbuf      *controlBuffer
	sendQuota uint32
	oiws      uint32 // outbound initial window size.
	// estdStreams is map of all established streams that are not cleaned-up yet.
	// On client-side, this is all streams whose headers were sent out.
	// On server-side, this is all streams whose headers were received.
	estdStreams map[uint32]*outStream // Established streams.
	// activeStreams is a linked-list of all streams that have data to send and some
	// stream-level flow control quota.
	// Each of these streams internally have a list of data items(and perhaps trailers
	// on the server-side) to be sent out.
	activeStreams *outStreamList
	framer        *framer
	hBuf          *bytes.Buffer  // The buffer for HPACK encoding.
	hEnc          *hpack.Encoder // HPACK encoder.
	bdpEst        *bdpEstimator
	draining      bool

	// Side-specific handlers
	ssGoAwayHandler func(*goAway) (bool, error)
***REMOVED***

func newLoopyWriter(s side, fr *framer, cbuf *controlBuffer, bdpEst *bdpEstimator) *loopyWriter ***REMOVED***
	var buf bytes.Buffer
	l := &loopyWriter***REMOVED***
		side:          s,
		cbuf:          cbuf,
		sendQuota:     defaultWindowSize,
		oiws:          defaultWindowSize,
		estdStreams:   make(map[uint32]*outStream),
		activeStreams: newOutStreamList(),
		framer:        fr,
		hBuf:          &buf,
		hEnc:          hpack.NewEncoder(&buf),
		bdpEst:        bdpEst,
	***REMOVED***
	return l
***REMOVED***

const minBatchSize = 1000

// run should be run in a separate goroutine.
// It reads control frames from controlBuf and processes them by:
// 1. Updating loopy's internal state, or/and
// 2. Writing out HTTP2 frames on the wire.
//
// Loopy keeps all active streams with data to send in a linked-list.
// All streams in the activeStreams linked-list must have both:
// 1. Data to send, and
// 2. Stream level flow control quota available.
//
// In each iteration of run loop, other than processing the incoming control
// frame, loopy calls processData, which processes one node from the activeStreams linked-list.
// This results in writing of HTTP2 frames into an underlying write buffer.
// When there's no more control frames to read from controlBuf, loopy flushes the write buffer.
// As an optimization, to increase the batch size for each flush, loopy yields the processor, once
// if the batch size is too low to give stream goroutines a chance to fill it up.
func (l *loopyWriter) run() (err error) ***REMOVED***
	defer func() ***REMOVED***
		if err == ErrConnClosing ***REMOVED***
			// Don't log ErrConnClosing as error since it happens
			// 1. When the connection is closed by some other known issue.
			// 2. User closed the connection.
			// 3. A graceful close of connection.
			if logger.V(logLevel) ***REMOVED***
				logger.Infof("transport: loopyWriter.run returning. %v", err)
			***REMOVED***
			err = nil
		***REMOVED***
	***REMOVED***()
	for ***REMOVED***
		it, err := l.cbuf.get(true)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if err = l.handle(it); err != nil ***REMOVED***
			return err
		***REMOVED***
		if _, err = l.processData(); err != nil ***REMOVED***
			return err
		***REMOVED***
		gosched := true
	hasdata:
		for ***REMOVED***
			it, err := l.cbuf.get(false)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if it != nil ***REMOVED***
				if err = l.handle(it); err != nil ***REMOVED***
					return err
				***REMOVED***
				if _, err = l.processData(); err != nil ***REMOVED***
					return err
				***REMOVED***
				continue hasdata
			***REMOVED***
			isEmpty, err := l.processData()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if !isEmpty ***REMOVED***
				continue hasdata
			***REMOVED***
			if gosched ***REMOVED***
				gosched = false
				if l.framer.writer.offset < minBatchSize ***REMOVED***
					runtime.Gosched()
					continue hasdata
				***REMOVED***
			***REMOVED***
			l.framer.writer.Flush()
			break hasdata

		***REMOVED***
	***REMOVED***
***REMOVED***

func (l *loopyWriter) outgoingWindowUpdateHandler(w *outgoingWindowUpdate) error ***REMOVED***
	return l.framer.fr.WriteWindowUpdate(w.streamID, w.increment)
***REMOVED***

func (l *loopyWriter) incomingWindowUpdateHandler(w *incomingWindowUpdate) error ***REMOVED***
	// Otherwise update the quota.
	if w.streamID == 0 ***REMOVED***
		l.sendQuota += w.increment
		return nil
	***REMOVED***
	// Find the stream and update it.
	if str, ok := l.estdStreams[w.streamID]; ok ***REMOVED***
		str.bytesOutStanding -= int(w.increment)
		if strQuota := int(l.oiws) - str.bytesOutStanding; strQuota > 0 && str.state == waitingOnStreamQuota ***REMOVED***
			str.state = active
			l.activeStreams.enqueue(str)
			return nil
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (l *loopyWriter) outgoingSettingsHandler(s *outgoingSettings) error ***REMOVED***
	return l.framer.fr.WriteSettings(s.ss...)
***REMOVED***

func (l *loopyWriter) incomingSettingsHandler(s *incomingSettings) error ***REMOVED***
	if err := l.applySettings(s.ss); err != nil ***REMOVED***
		return err
	***REMOVED***
	return l.framer.fr.WriteSettingsAck()
***REMOVED***

func (l *loopyWriter) registerStreamHandler(h *registerStream) error ***REMOVED***
	str := &outStream***REMOVED***
		id:    h.streamID,
		state: empty,
		itl:   &itemList***REMOVED******REMOVED***,
		wq:    h.wq,
	***REMOVED***
	l.estdStreams[h.streamID] = str
	return nil
***REMOVED***

func (l *loopyWriter) headerHandler(h *headerFrame) error ***REMOVED***
	if l.side == serverSide ***REMOVED***
		str, ok := l.estdStreams[h.streamID]
		if !ok ***REMOVED***
			if logger.V(logLevel) ***REMOVED***
				logger.Warningf("transport: loopy doesn't recognize the stream: %d", h.streamID)
			***REMOVED***
			return nil
		***REMOVED***
		// Case 1.A: Server is responding back with headers.
		if !h.endStream ***REMOVED***
			return l.writeHeader(h.streamID, h.endStream, h.hf, h.onWrite)
		***REMOVED***
		// else:  Case 1.B: Server wants to close stream.

		if str.state != empty ***REMOVED*** // either active or waiting on stream quota.
			// add it str's list of items.
			str.itl.enqueue(h)
			return nil
		***REMOVED***
		if err := l.writeHeader(h.streamID, h.endStream, h.hf, h.onWrite); err != nil ***REMOVED***
			return err
		***REMOVED***
		return l.cleanupStreamHandler(h.cleanup)
	***REMOVED***
	// Case 2: Client wants to originate stream.
	str := &outStream***REMOVED***
		id:    h.streamID,
		state: empty,
		itl:   &itemList***REMOVED******REMOVED***,
		wq:    h.wq,
	***REMOVED***
	str.itl.enqueue(h)
	return l.originateStream(str)
***REMOVED***

func (l *loopyWriter) originateStream(str *outStream) error ***REMOVED***
	hdr := str.itl.dequeue().(*headerFrame)
	if err := hdr.initStream(str.id); err != nil ***REMOVED***
		if err == ErrConnClosing ***REMOVED***
			return err
		***REMOVED***
		// Other errors(errStreamDrain) need not close transport.
		return nil
	***REMOVED***
	if err := l.writeHeader(str.id, hdr.endStream, hdr.hf, hdr.onWrite); err != nil ***REMOVED***
		return err
	***REMOVED***
	l.estdStreams[str.id] = str
	return nil
***REMOVED***

func (l *loopyWriter) writeHeader(streamID uint32, endStream bool, hf []hpack.HeaderField, onWrite func()) error ***REMOVED***
	if onWrite != nil ***REMOVED***
		onWrite()
	***REMOVED***
	l.hBuf.Reset()
	for _, f := range hf ***REMOVED***
		if err := l.hEnc.WriteField(f); err != nil ***REMOVED***
			if logger.V(logLevel) ***REMOVED***
				logger.Warningf("transport: loopyWriter.writeHeader encountered error while encoding headers: %v", err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	var (
		err               error
		endHeaders, first bool
	)
	first = true
	for !endHeaders ***REMOVED***
		size := l.hBuf.Len()
		if size > http2MaxFrameLen ***REMOVED***
			size = http2MaxFrameLen
		***REMOVED*** else ***REMOVED***
			endHeaders = true
		***REMOVED***
		if first ***REMOVED***
			first = false
			err = l.framer.fr.WriteHeaders(http2.HeadersFrameParam***REMOVED***
				StreamID:      streamID,
				BlockFragment: l.hBuf.Next(size),
				EndStream:     endStream,
				EndHeaders:    endHeaders,
			***REMOVED***)
		***REMOVED*** else ***REMOVED***
			err = l.framer.fr.WriteContinuation(
				streamID,
				endHeaders,
				l.hBuf.Next(size),
			)
		***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (l *loopyWriter) preprocessData(df *dataFrame) error ***REMOVED***
	str, ok := l.estdStreams[df.streamID]
	if !ok ***REMOVED***
		return nil
	***REMOVED***
	// If we got data for a stream it means that
	// stream was originated and the headers were sent out.
	str.itl.enqueue(df)
	if str.state == empty ***REMOVED***
		str.state = active
		l.activeStreams.enqueue(str)
	***REMOVED***
	return nil
***REMOVED***

func (l *loopyWriter) pingHandler(p *ping) error ***REMOVED***
	if !p.ack ***REMOVED***
		l.bdpEst.timesnap(p.data)
	***REMOVED***
	return l.framer.fr.WritePing(p.ack, p.data)

***REMOVED***

func (l *loopyWriter) outFlowControlSizeRequestHandler(o *outFlowControlSizeRequest) error ***REMOVED***
	o.resp <- l.sendQuota
	return nil
***REMOVED***

func (l *loopyWriter) cleanupStreamHandler(c *cleanupStream) error ***REMOVED***
	c.onWrite()
	if str, ok := l.estdStreams[c.streamID]; ok ***REMOVED***
		// On the server side it could be a trailers-only response or
		// a RST_STREAM before stream initialization thus the stream might
		// not be established yet.
		delete(l.estdStreams, c.streamID)
		str.deleteSelf()
	***REMOVED***
	if c.rst ***REMOVED*** // If RST_STREAM needs to be sent.
		if err := l.framer.fr.WriteRSTStream(c.streamID, c.rstCode); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if l.side == clientSide && l.draining && len(l.estdStreams) == 0 ***REMOVED***
		return ErrConnClosing
	***REMOVED***
	return nil
***REMOVED***

func (l *loopyWriter) earlyAbortStreamHandler(eas *earlyAbortStream) error ***REMOVED***
	if l.side == clientSide ***REMOVED***
		return errors.New("earlyAbortStream not handled on client")
	***REMOVED***
	// In case the caller forgets to set the http status, default to 200.
	if eas.httpStatus == 0 ***REMOVED***
		eas.httpStatus = 200
	***REMOVED***
	headerFields := []hpack.HeaderField***REMOVED***
		***REMOVED***Name: ":status", Value: strconv.Itoa(int(eas.httpStatus))***REMOVED***,
		***REMOVED***Name: "content-type", Value: grpcutil.ContentType(eas.contentSubtype)***REMOVED***,
		***REMOVED***Name: "grpc-status", Value: strconv.Itoa(int(eas.status.Code()))***REMOVED***,
		***REMOVED***Name: "grpc-message", Value: encodeGrpcMessage(eas.status.Message())***REMOVED***,
	***REMOVED***

	if err := l.writeHeader(eas.streamID, true, headerFields, nil); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func (l *loopyWriter) incomingGoAwayHandler(*incomingGoAway) error ***REMOVED***
	if l.side == clientSide ***REMOVED***
		l.draining = true
		if len(l.estdStreams) == 0 ***REMOVED***
			return ErrConnClosing
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (l *loopyWriter) goAwayHandler(g *goAway) error ***REMOVED***
	// Handling of outgoing GoAway is very specific to side.
	if l.ssGoAwayHandler != nil ***REMOVED***
		draining, err := l.ssGoAwayHandler(g)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		l.draining = draining
	***REMOVED***
	return nil
***REMOVED***

func (l *loopyWriter) handle(i interface***REMOVED******REMOVED***) error ***REMOVED***
	switch i := i.(type) ***REMOVED***
	case *incomingWindowUpdate:
		return l.incomingWindowUpdateHandler(i)
	case *outgoingWindowUpdate:
		return l.outgoingWindowUpdateHandler(i)
	case *incomingSettings:
		return l.incomingSettingsHandler(i)
	case *outgoingSettings:
		return l.outgoingSettingsHandler(i)
	case *headerFrame:
		return l.headerHandler(i)
	case *registerStream:
		return l.registerStreamHandler(i)
	case *cleanupStream:
		return l.cleanupStreamHandler(i)
	case *earlyAbortStream:
		return l.earlyAbortStreamHandler(i)
	case *incomingGoAway:
		return l.incomingGoAwayHandler(i)
	case *dataFrame:
		return l.preprocessData(i)
	case *ping:
		return l.pingHandler(i)
	case *goAway:
		return l.goAwayHandler(i)
	case *outFlowControlSizeRequest:
		return l.outFlowControlSizeRequestHandler(i)
	default:
		return fmt.Errorf("transport: unknown control message type %T", i)
	***REMOVED***
***REMOVED***

func (l *loopyWriter) applySettings(ss []http2.Setting) error ***REMOVED***
	for _, s := range ss ***REMOVED***
		switch s.ID ***REMOVED***
		case http2.SettingInitialWindowSize:
			o := l.oiws
			l.oiws = s.Val
			if o < l.oiws ***REMOVED***
				// If the new limit is greater make all depleted streams active.
				for _, stream := range l.estdStreams ***REMOVED***
					if stream.state == waitingOnStreamQuota ***REMOVED***
						stream.state = active
						l.activeStreams.enqueue(stream)
					***REMOVED***
				***REMOVED***
			***REMOVED***
		case http2.SettingHeaderTableSize:
			updateHeaderTblSize(l.hEnc, s.Val)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// processData removes the first stream from active streams, writes out at most 16KB
// of its data and then puts it at the end of activeStreams if there's still more data
// to be sent and stream has some stream-level flow control.
func (l *loopyWriter) processData() (bool, error) ***REMOVED***
	if l.sendQuota == 0 ***REMOVED***
		return true, nil
	***REMOVED***
	str := l.activeStreams.dequeue() // Remove the first stream.
	if str == nil ***REMOVED***
		return true, nil
	***REMOVED***
	dataItem := str.itl.peek().(*dataFrame) // Peek at the first data item this stream.
	// A data item is represented by a dataFrame, since it later translates into
	// multiple HTTP2 data frames.
	// Every dataFrame has two buffers; h that keeps grpc-message header and d that is acutal data.
	// As an optimization to keep wire traffic low, data from d is copied to h to make as big as the
	// maximum possilbe HTTP2 frame size.

	if len(dataItem.h) == 0 && len(dataItem.d) == 0 ***REMOVED*** // Empty data frame
		// Client sends out empty data frame with endStream = true
		if err := l.framer.fr.WriteData(dataItem.streamID, dataItem.endStream, nil); err != nil ***REMOVED***
			return false, err
		***REMOVED***
		str.itl.dequeue() // remove the empty data item from stream
		if str.itl.isEmpty() ***REMOVED***
			str.state = empty
		***REMOVED*** else if trailer, ok := str.itl.peek().(*headerFrame); ok ***REMOVED*** // the next item is trailers.
			if err := l.writeHeader(trailer.streamID, trailer.endStream, trailer.hf, trailer.onWrite); err != nil ***REMOVED***
				return false, err
			***REMOVED***
			if err := l.cleanupStreamHandler(trailer.cleanup); err != nil ***REMOVED***
				return false, nil
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			l.activeStreams.enqueue(str)
		***REMOVED***
		return false, nil
	***REMOVED***
	var (
		buf []byte
	)
	// Figure out the maximum size we can send
	maxSize := http2MaxFrameLen
	if strQuota := int(l.oiws) - str.bytesOutStanding; strQuota <= 0 ***REMOVED*** // stream-level flow control.
		str.state = waitingOnStreamQuota
		return false, nil
	***REMOVED*** else if maxSize > strQuota ***REMOVED***
		maxSize = strQuota
	***REMOVED***
	if maxSize > int(l.sendQuota) ***REMOVED*** // connection-level flow control.
		maxSize = int(l.sendQuota)
	***REMOVED***
	// Compute how much of the header and data we can send within quota and max frame length
	hSize := min(maxSize, len(dataItem.h))
	dSize := min(maxSize-hSize, len(dataItem.d))
	if hSize != 0 ***REMOVED***
		if dSize == 0 ***REMOVED***
			buf = dataItem.h
		***REMOVED*** else ***REMOVED***
			// We can add some data to grpc message header to distribute bytes more equally across frames.
			// Copy on the stack to avoid generating garbage
			var localBuf [http2MaxFrameLen]byte
			copy(localBuf[:hSize], dataItem.h)
			copy(localBuf[hSize:], dataItem.d[:dSize])
			buf = localBuf[:hSize+dSize]
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		buf = dataItem.d
	***REMOVED***

	size := hSize + dSize

	// Now that outgoing flow controls are checked we can replenish str's write quota
	str.wq.replenish(size)
	var endStream bool
	// If this is the last data message on this stream and all of it can be written in this iteration.
	if dataItem.endStream && len(dataItem.h)+len(dataItem.d) <= size ***REMOVED***
		endStream = true
	***REMOVED***
	if dataItem.onEachWrite != nil ***REMOVED***
		dataItem.onEachWrite()
	***REMOVED***
	if err := l.framer.fr.WriteData(dataItem.streamID, endStream, buf[:size]); err != nil ***REMOVED***
		return false, err
	***REMOVED***
	str.bytesOutStanding += size
	l.sendQuota -= uint32(size)
	dataItem.h = dataItem.h[hSize:]
	dataItem.d = dataItem.d[dSize:]

	if len(dataItem.h) == 0 && len(dataItem.d) == 0 ***REMOVED*** // All the data from that message was written out.
		str.itl.dequeue()
	***REMOVED***
	if str.itl.isEmpty() ***REMOVED***
		str.state = empty
	***REMOVED*** else if trailer, ok := str.itl.peek().(*headerFrame); ok ***REMOVED*** // The next item is trailers.
		if err := l.writeHeader(trailer.streamID, trailer.endStream, trailer.hf, trailer.onWrite); err != nil ***REMOVED***
			return false, err
		***REMOVED***
		if err := l.cleanupStreamHandler(trailer.cleanup); err != nil ***REMOVED***
			return false, err
		***REMOVED***
	***REMOVED*** else if int(l.oiws)-str.bytesOutStanding <= 0 ***REMOVED*** // Ran out of stream quota.
		str.state = waitingOnStreamQuota
	***REMOVED*** else ***REMOVED*** // Otherwise add it back to the list of active streams.
		l.activeStreams.enqueue(str)
	***REMOVED***
	return false, nil
***REMOVED***

func min(a, b int) int ***REMOVED***
	if a < b ***REMOVED***
		return a
	***REMOVED***
	return b
***REMOVED***
