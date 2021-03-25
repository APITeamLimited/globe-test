// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package norm

import "io"

type normWriter struct ***REMOVED***
	rb  reorderBuffer
	w   io.Writer
	buf []byte
***REMOVED***

// Write implements the standard write interface.  If the last characters are
// not at a normalization boundary, the bytes will be buffered for the next
// write. The remaining bytes will be written on close.
func (w *normWriter) Write(data []byte) (n int, err error) ***REMOVED***
	// Process data in pieces to keep w.buf size bounded.
	const chunk = 4000

	for len(data) > 0 ***REMOVED***
		// Normalize into w.buf.
		m := len(data)
		if m > chunk ***REMOVED***
			m = chunk
		***REMOVED***
		w.rb.src = inputBytes(data[:m])
		w.rb.nsrc = m
		w.buf = doAppend(&w.rb, w.buf, 0)
		data = data[m:]
		n += m

		// Write out complete prefix, save remainder.
		// Note that lastBoundary looks back at most 31 runes.
		i := lastBoundary(&w.rb.f, w.buf)
		if i == -1 ***REMOVED***
			i = 0
		***REMOVED***
		if i > 0 ***REMOVED***
			if _, err = w.w.Write(w.buf[:i]); err != nil ***REMOVED***
				break
			***REMOVED***
			bn := copy(w.buf, w.buf[i:])
			w.buf = w.buf[:bn]
		***REMOVED***
	***REMOVED***
	return n, err
***REMOVED***

// Close forces data that remains in the buffer to be written.
func (w *normWriter) Close() error ***REMOVED***
	if len(w.buf) > 0 ***REMOVED***
		_, err := w.w.Write(w.buf)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Writer returns a new writer that implements Write(b)
// by writing f(b) to w. The returned writer may use an
// internal buffer to maintain state across Write calls.
// Calling its Close method writes any buffered data to w.
func (f Form) Writer(w io.Writer) io.WriteCloser ***REMOVED***
	wr := &normWriter***REMOVED***rb: reorderBuffer***REMOVED******REMOVED***, w: w***REMOVED***
	wr.rb.init(f, nil)
	return wr
***REMOVED***

type normReader struct ***REMOVED***
	rb           reorderBuffer
	r            io.Reader
	inbuf        []byte
	outbuf       []byte
	bufStart     int
	lastBoundary int
	err          error
***REMOVED***

// Read implements the standard read interface.
func (r *normReader) Read(p []byte) (int, error) ***REMOVED***
	for ***REMOVED***
		if r.lastBoundary-r.bufStart > 0 ***REMOVED***
			n := copy(p, r.outbuf[r.bufStart:r.lastBoundary])
			r.bufStart += n
			if r.lastBoundary-r.bufStart > 0 ***REMOVED***
				return n, nil
			***REMOVED***
			return n, r.err
		***REMOVED***
		if r.err != nil ***REMOVED***
			return 0, r.err
		***REMOVED***
		outn := copy(r.outbuf, r.outbuf[r.lastBoundary:])
		r.outbuf = r.outbuf[0:outn]
		r.bufStart = 0

		n, err := r.r.Read(r.inbuf)
		r.rb.src = inputBytes(r.inbuf[0:n])
		r.rb.nsrc, r.err = n, err
		if n > 0 ***REMOVED***
			r.outbuf = doAppend(&r.rb, r.outbuf, 0)
		***REMOVED***
		if err == io.EOF ***REMOVED***
			r.lastBoundary = len(r.outbuf)
		***REMOVED*** else ***REMOVED***
			r.lastBoundary = lastBoundary(&r.rb.f, r.outbuf)
			if r.lastBoundary == -1 ***REMOVED***
				r.lastBoundary = 0
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// Reader returns a new reader that implements Read
// by reading data from r and returning f(data).
func (f Form) Reader(r io.Reader) io.Reader ***REMOVED***
	const chunk = 4000
	buf := make([]byte, chunk)
	rr := &normReader***REMOVED***rb: reorderBuffer***REMOVED******REMOVED***, r: r, inbuf: buf***REMOVED***
	rr.rb.init(f, buf)
	return rr
***REMOVED***
