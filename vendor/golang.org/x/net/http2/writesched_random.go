// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http2

import "math"

// NewRandomWriteScheduler constructs a WriteScheduler that ignores HTTP/2
// priorities. Control frames like SETTINGS and PING are written before DATA
// frames, but if no control frames are queued and multiple streams have queued
// HEADERS or DATA frames, Pop selects a ready stream arbitrarily.
func NewRandomWriteScheduler() WriteScheduler ***REMOVED***
	return &randomWriteScheduler***REMOVED***sq: make(map[uint32]*writeQueue)***REMOVED***
***REMOVED***

type randomWriteScheduler struct ***REMOVED***
	// zero are frames not associated with a specific stream.
	zero writeQueue

	// sq contains the stream-specific queues, keyed by stream ID.
	// When a stream is idle, closed, or emptied, it's deleted
	// from the map.
	sq map[uint32]*writeQueue

	// pool of empty queues for reuse.
	queuePool writeQueuePool
***REMOVED***

func (ws *randomWriteScheduler) OpenStream(streamID uint32, options OpenStreamOptions) ***REMOVED***
	// no-op: idle streams are not tracked
***REMOVED***

func (ws *randomWriteScheduler) CloseStream(streamID uint32) ***REMOVED***
	q, ok := ws.sq[streamID]
	if !ok ***REMOVED***
		return
	***REMOVED***
	delete(ws.sq, streamID)
	ws.queuePool.put(q)
***REMOVED***

func (ws *randomWriteScheduler) AdjustStream(streamID uint32, priority PriorityParam) ***REMOVED***
	// no-op: priorities are ignored
***REMOVED***

func (ws *randomWriteScheduler) Push(wr FrameWriteRequest) ***REMOVED***
	id := wr.StreamID()
	if id == 0 ***REMOVED***
		ws.zero.push(wr)
		return
	***REMOVED***
	q, ok := ws.sq[id]
	if !ok ***REMOVED***
		q = ws.queuePool.get()
		ws.sq[id] = q
	***REMOVED***
	q.push(wr)
***REMOVED***

func (ws *randomWriteScheduler) Pop() (FrameWriteRequest, bool) ***REMOVED***
	// Control frames first.
	if !ws.zero.empty() ***REMOVED***
		return ws.zero.shift(), true
	***REMOVED***
	// Iterate over all non-idle streams until finding one that can be consumed.
	for streamID, q := range ws.sq ***REMOVED***
		if wr, ok := q.consume(math.MaxInt32); ok ***REMOVED***
			if q.empty() ***REMOVED***
				delete(ws.sq, streamID)
				ws.queuePool.put(q)
			***REMOVED***
			return wr, true
		***REMOVED***
	***REMOVED***
	return FrameWriteRequest***REMOVED******REMOVED***, false
***REMOVED***
