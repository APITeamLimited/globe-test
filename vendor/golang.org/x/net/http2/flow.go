// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Flow control

package http2

// flow is the flow control window's size.
type flow struct ***REMOVED***
	// n is the number of DATA bytes we're allowed to send.
	// A flow is kept both on a conn and a per-stream.
	n int32

	// conn points to the shared connection-level flow that is
	// shared by all streams on that conn. It is nil for the flow
	// that's on the conn directly.
	conn *flow
***REMOVED***

func (f *flow) setConnFlow(cf *flow) ***REMOVED*** f.conn = cf ***REMOVED***

func (f *flow) available() int32 ***REMOVED***
	n := f.n
	if f.conn != nil && f.conn.n < n ***REMOVED***
		n = f.conn.n
	***REMOVED***
	return n
***REMOVED***

func (f *flow) take(n int32) ***REMOVED***
	if n > f.available() ***REMOVED***
		panic("internal error: took too much")
	***REMOVED***
	f.n -= n
	if f.conn != nil ***REMOVED***
		f.conn.n -= n
	***REMOVED***
***REMOVED***

// add adds n bytes (positive or negative) to the flow control window.
// It returns false if the sum would exceed 2^31-1.
func (f *flow) add(n int32) bool ***REMOVED***
	sum := f.n + n
	if (sum > n) == (f.n > 0) ***REMOVED***
		f.n = sum
		return true
	***REMOVED***
	return false
***REMOVED***
