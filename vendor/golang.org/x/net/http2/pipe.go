// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http2

import (
	"errors"
	"io"
	"sync"
)

// pipe is a goroutine-safe io.Reader/io.Writer pair. It's like
// io.Pipe except there are no PipeReader/PipeWriter halves, and the
// underlying buffer is an interface. (io.Pipe is always unbuffered)
type pipe struct ***REMOVED***
	mu       sync.Mutex
	c        sync.Cond     // c.L lazily initialized to &p.mu
	b        pipeBuffer    // nil when done reading
	unread   int           // bytes unread when done
	err      error         // read error once empty. non-nil means closed.
	breakErr error         // immediate read error (caller doesn't see rest of b)
	donec    chan struct***REMOVED******REMOVED*** // closed on error
	readFn   func()        // optional code to run in Read before error
***REMOVED***

type pipeBuffer interface ***REMOVED***
	Len() int
	io.Writer
	io.Reader
***REMOVED***

func (p *pipe) Len() int ***REMOVED***
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.b == nil ***REMOVED***
		return p.unread
	***REMOVED***
	return p.b.Len()
***REMOVED***

// Read waits until data is available and copies bytes
// from the buffer into p.
func (p *pipe) Read(d []byte) (n int, err error) ***REMOVED***
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.c.L == nil ***REMOVED***
		p.c.L = &p.mu
	***REMOVED***
	for ***REMOVED***
		if p.breakErr != nil ***REMOVED***
			return 0, p.breakErr
		***REMOVED***
		if p.b != nil && p.b.Len() > 0 ***REMOVED***
			return p.b.Read(d)
		***REMOVED***
		if p.err != nil ***REMOVED***
			if p.readFn != nil ***REMOVED***
				p.readFn()     // e.g. copy trailers
				p.readFn = nil // not sticky like p.err
			***REMOVED***
			p.b = nil
			return 0, p.err
		***REMOVED***
		p.c.Wait()
	***REMOVED***
***REMOVED***

var errClosedPipeWrite = errors.New("write on closed buffer")

// Write copies bytes from p into the buffer and wakes a reader.
// It is an error to write more data than the buffer can hold.
func (p *pipe) Write(d []byte) (n int, err error) ***REMOVED***
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.c.L == nil ***REMOVED***
		p.c.L = &p.mu
	***REMOVED***
	defer p.c.Signal()
	if p.err != nil ***REMOVED***
		return 0, errClosedPipeWrite
	***REMOVED***
	if p.breakErr != nil ***REMOVED***
		p.unread += len(d)
		return len(d), nil // discard when there is no reader
	***REMOVED***
	return p.b.Write(d)
***REMOVED***

// CloseWithError causes the next Read (waking up a current blocked
// Read if needed) to return the provided err after all data has been
// read.
//
// The error must be non-nil.
func (p *pipe) CloseWithError(err error) ***REMOVED*** p.closeWithError(&p.err, err, nil) ***REMOVED***

// BreakWithError causes the next Read (waking up a current blocked
// Read if needed) to return the provided err immediately, without
// waiting for unread data.
func (p *pipe) BreakWithError(err error) ***REMOVED*** p.closeWithError(&p.breakErr, err, nil) ***REMOVED***

// closeWithErrorAndCode is like CloseWithError but also sets some code to run
// in the caller's goroutine before returning the error.
func (p *pipe) closeWithErrorAndCode(err error, fn func()) ***REMOVED*** p.closeWithError(&p.err, err, fn) ***REMOVED***

func (p *pipe) closeWithError(dst *error, err error, fn func()) ***REMOVED***
	if err == nil ***REMOVED***
		panic("err must be non-nil")
	***REMOVED***
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.c.L == nil ***REMOVED***
		p.c.L = &p.mu
	***REMOVED***
	defer p.c.Signal()
	if *dst != nil ***REMOVED***
		// Already been done.
		return
	***REMOVED***
	p.readFn = fn
	if dst == &p.breakErr ***REMOVED***
		if p.b != nil ***REMOVED***
			p.unread += p.b.Len()
		***REMOVED***
		p.b = nil
	***REMOVED***
	*dst = err
	p.closeDoneLocked()
***REMOVED***

// requires p.mu be held.
func (p *pipe) closeDoneLocked() ***REMOVED***
	if p.donec == nil ***REMOVED***
		return
	***REMOVED***
	// Close if unclosed. This isn't racy since we always
	// hold p.mu while closing.
	select ***REMOVED***
	case <-p.donec:
	default:
		close(p.donec)
	***REMOVED***
***REMOVED***

// Err returns the error (if any) first set by BreakWithError or CloseWithError.
func (p *pipe) Err() error ***REMOVED***
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.breakErr != nil ***REMOVED***
		return p.breakErr
	***REMOVED***
	return p.err
***REMOVED***

// Done returns a channel which is closed if and when this pipe is closed
// with CloseWithError.
func (p *pipe) Done() <-chan struct***REMOVED******REMOVED*** ***REMOVED***
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.donec == nil ***REMOVED***
		p.donec = make(chan struct***REMOVED******REMOVED***)
		if p.err != nil || p.breakErr != nil ***REMOVED***
			// Already hit an error.
			p.closeDoneLocked()
		***REMOVED***
	***REMOVED***
	return p.donec
***REMOVED***
