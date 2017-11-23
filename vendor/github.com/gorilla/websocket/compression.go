// Copyright 2017 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"compress/flate"
	"errors"
	"io"
	"strings"
	"sync"
)

const (
	minCompressionLevel     = -2 // flate.HuffmanOnly not defined in Go < 1.6
	maxCompressionLevel     = flate.BestCompression
	defaultCompressionLevel = 1
)

var (
	flateWriterPools [maxCompressionLevel - minCompressionLevel + 1]sync.Pool
	flateReaderPool  = sync.Pool***REMOVED***New: func() interface***REMOVED******REMOVED*** ***REMOVED***
		return flate.NewReader(nil)
	***REMOVED******REMOVED***
)

func decompressNoContextTakeover(r io.Reader) io.ReadCloser ***REMOVED***
	const tail =
	// Add four bytes as specified in RFC
	"\x00\x00\xff\xff" +
		// Add final block to squelch unexpected EOF error from flate reader.
		"\x01\x00\x00\xff\xff"

	fr, _ := flateReaderPool.Get().(io.ReadCloser)
	fr.(flate.Resetter).Reset(io.MultiReader(r, strings.NewReader(tail)), nil)
	return &flateReadWrapper***REMOVED***fr***REMOVED***
***REMOVED***

func isValidCompressionLevel(level int) bool ***REMOVED***
	return minCompressionLevel <= level && level <= maxCompressionLevel
***REMOVED***

func compressNoContextTakeover(w io.WriteCloser, level int) io.WriteCloser ***REMOVED***
	p := &flateWriterPools[level-minCompressionLevel]
	tw := &truncWriter***REMOVED***w: w***REMOVED***
	fw, _ := p.Get().(*flate.Writer)
	if fw == nil ***REMOVED***
		fw, _ = flate.NewWriter(tw, level)
	***REMOVED*** else ***REMOVED***
		fw.Reset(tw)
	***REMOVED***
	return &flateWriteWrapper***REMOVED***fw: fw, tw: tw, p: p***REMOVED***
***REMOVED***

// truncWriter is an io.Writer that writes all but the last four bytes of the
// stream to another io.Writer.
type truncWriter struct ***REMOVED***
	w io.WriteCloser
	n int
	p [4]byte
***REMOVED***

func (w *truncWriter) Write(p []byte) (int, error) ***REMOVED***
	n := 0

	// fill buffer first for simplicity.
	if w.n < len(w.p) ***REMOVED***
		n = copy(w.p[w.n:], p)
		p = p[n:]
		w.n += n
		if len(p) == 0 ***REMOVED***
			return n, nil
		***REMOVED***
	***REMOVED***

	m := len(p)
	if m > len(w.p) ***REMOVED***
		m = len(w.p)
	***REMOVED***

	if nn, err := w.w.Write(w.p[:m]); err != nil ***REMOVED***
		return n + nn, err
	***REMOVED***

	copy(w.p[:], w.p[m:])
	copy(w.p[len(w.p)-m:], p[len(p)-m:])
	nn, err := w.w.Write(p[:len(p)-m])
	return n + nn, err
***REMOVED***

type flateWriteWrapper struct ***REMOVED***
	fw *flate.Writer
	tw *truncWriter
	p  *sync.Pool
***REMOVED***

func (w *flateWriteWrapper) Write(p []byte) (int, error) ***REMOVED***
	if w.fw == nil ***REMOVED***
		return 0, errWriteClosed
	***REMOVED***
	return w.fw.Write(p)
***REMOVED***

func (w *flateWriteWrapper) Close() error ***REMOVED***
	if w.fw == nil ***REMOVED***
		return errWriteClosed
	***REMOVED***
	err1 := w.fw.Flush()
	w.p.Put(w.fw)
	w.fw = nil
	if w.tw.p != [4]byte***REMOVED***0, 0, 0xff, 0xff***REMOVED*** ***REMOVED***
		return errors.New("websocket: internal error, unexpected bytes at end of flate stream")
	***REMOVED***
	err2 := w.tw.w.Close()
	if err1 != nil ***REMOVED***
		return err1
	***REMOVED***
	return err2
***REMOVED***

type flateReadWrapper struct ***REMOVED***
	fr io.ReadCloser
***REMOVED***

func (r *flateReadWrapper) Read(p []byte) (int, error) ***REMOVED***
	if r.fr == nil ***REMOVED***
		return 0, io.ErrClosedPipe
	***REMOVED***
	n, err := r.fr.Read(p)
	if err == io.EOF ***REMOVED***
		// Preemptively place the reader back in the pool. This helps with
		// scenarios where the application does not call NextReader() soon after
		// this final read.
		r.Close()
	***REMOVED***
	return n, err
***REMOVED***

func (r *flateReadWrapper) Close() error ***REMOVED***
	if r.fr == nil ***REMOVED***
		return io.ErrClosedPipe
	***REMOVED***
	err := r.fr.Close()
	flateReaderPool.Put(r.fr)
	r.fr = nil
	return err
***REMOVED***
