// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gzip

import (
	"errors"
	"fmt"
	"hash/crc32"
	"io"

	"github.com/klauspost/compress/flate"
)

// These constants are copied from the flate package, so that code that imports
// "compress/gzip" does not also have to import "compress/flate".
const (
	NoCompression       = flate.NoCompression
	BestSpeed           = flate.BestSpeed
	BestCompression     = flate.BestCompression
	DefaultCompression  = flate.DefaultCompression
	ConstantCompression = flate.ConstantCompression
	HuffmanOnly         = flate.HuffmanOnly

	// StatelessCompression will do compression but without maintaining any state
	// between Write calls.
	// There will be no memory kept between Write calls,
	// but compression and speed will be suboptimal.
	// Because of this, the size of actual Write calls will affect output size.
	StatelessCompression = -3
)

// A Writer is an io.WriteCloser.
// Writes to a Writer are compressed and written to w.
type Writer struct ***REMOVED***
	Header      // written at first call to Write, Flush, or Close
	w           io.Writer
	level       int
	err         error
	compressor  *flate.Writer
	digest      uint32 // CRC-32, IEEE polynomial (section 8)
	size        uint32 // Uncompressed size (section 2.3.1)
	wroteHeader bool
	closed      bool
	buf         [10]byte
***REMOVED***

// NewWriter returns a new Writer.
// Writes to the returned writer are compressed and written to w.
//
// It is the caller's responsibility to call Close on the WriteCloser when done.
// Writes may be buffered and not flushed until Close.
//
// Callers that wish to set the fields in Writer.Header must do so before
// the first call to Write, Flush, or Close.
func NewWriter(w io.Writer) *Writer ***REMOVED***
	z, _ := NewWriterLevel(w, DefaultCompression)
	return z
***REMOVED***

// NewWriterLevel is like NewWriter but specifies the compression level instead
// of assuming DefaultCompression.
//
// The compression level can be DefaultCompression, NoCompression, or any
// integer value between BestSpeed and BestCompression inclusive. The error
// returned will be nil if the level is valid.
func NewWriterLevel(w io.Writer, level int) (*Writer, error) ***REMOVED***
	if level < StatelessCompression || level > BestCompression ***REMOVED***
		return nil, fmt.Errorf("gzip: invalid compression level: %d", level)
	***REMOVED***
	z := new(Writer)
	z.init(w, level)
	return z, nil
***REMOVED***

func (z *Writer) init(w io.Writer, level int) ***REMOVED***
	compressor := z.compressor
	if level != StatelessCompression ***REMOVED***
		if compressor != nil ***REMOVED***
			compressor.Reset(w)
		***REMOVED***
	***REMOVED***

	*z = Writer***REMOVED***
		Header: Header***REMOVED***
			OS: 255, // unknown
		***REMOVED***,
		w:          w,
		level:      level,
		compressor: compressor,
	***REMOVED***
***REMOVED***

// Reset discards the Writer z's state and makes it equivalent to the
// result of its original state from NewWriter or NewWriterLevel, but
// writing to w instead. This permits reusing a Writer rather than
// allocating a new one.
func (z *Writer) Reset(w io.Writer) ***REMOVED***
	z.init(w, z.level)
***REMOVED***

// writeBytes writes a length-prefixed byte slice to z.w.
func (z *Writer) writeBytes(b []byte) error ***REMOVED***
	if len(b) > 0xffff ***REMOVED***
		return errors.New("gzip.Write: Extra data is too large")
	***REMOVED***
	le.PutUint16(z.buf[:2], uint16(len(b)))
	_, err := z.w.Write(z.buf[:2])
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	_, err = z.w.Write(b)
	return err
***REMOVED***

// writeString writes a UTF-8 string s in GZIP's format to z.w.
// GZIP (RFC 1952) specifies that strings are NUL-terminated ISO 8859-1 (Latin-1).
func (z *Writer) writeString(s string) (err error) ***REMOVED***
	// GZIP stores Latin-1 strings; error if non-Latin-1; convert if non-ASCII.
	needconv := false
	for _, v := range s ***REMOVED***
		if v == 0 || v > 0xff ***REMOVED***
			return errors.New("gzip.Write: non-Latin-1 header string")
		***REMOVED***
		if v > 0x7f ***REMOVED***
			needconv = true
		***REMOVED***
	***REMOVED***
	if needconv ***REMOVED***
		b := make([]byte, 0, len(s))
		for _, v := range s ***REMOVED***
			b = append(b, byte(v))
		***REMOVED***
		_, err = z.w.Write(b)
	***REMOVED*** else ***REMOVED***
		_, err = io.WriteString(z.w, s)
	***REMOVED***
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	// GZIP strings are NUL-terminated.
	z.buf[0] = 0
	_, err = z.w.Write(z.buf[:1])
	return err
***REMOVED***

// Write writes a compressed form of p to the underlying io.Writer. The
// compressed bytes are not necessarily flushed until the Writer is closed.
func (z *Writer) Write(p []byte) (int, error) ***REMOVED***
	if z.err != nil ***REMOVED***
		return 0, z.err
	***REMOVED***
	var n int
	// Write the GZIP header lazily.
	if !z.wroteHeader ***REMOVED***
		z.wroteHeader = true
		z.buf[0] = gzipID1
		z.buf[1] = gzipID2
		z.buf[2] = gzipDeflate
		z.buf[3] = 0
		if z.Extra != nil ***REMOVED***
			z.buf[3] |= 0x04
		***REMOVED***
		if z.Name != "" ***REMOVED***
			z.buf[3] |= 0x08
		***REMOVED***
		if z.Comment != "" ***REMOVED***
			z.buf[3] |= 0x10
		***REMOVED***
		le.PutUint32(z.buf[4:8], uint32(z.ModTime.Unix()))
		if z.level == BestCompression ***REMOVED***
			z.buf[8] = 2
		***REMOVED*** else if z.level == BestSpeed ***REMOVED***
			z.buf[8] = 4
		***REMOVED*** else ***REMOVED***
			z.buf[8] = 0
		***REMOVED***
		z.buf[9] = z.OS
		n, z.err = z.w.Write(z.buf[:10])
		if z.err != nil ***REMOVED***
			return n, z.err
		***REMOVED***
		if z.Extra != nil ***REMOVED***
			z.err = z.writeBytes(z.Extra)
			if z.err != nil ***REMOVED***
				return n, z.err
			***REMOVED***
		***REMOVED***
		if z.Name != "" ***REMOVED***
			z.err = z.writeString(z.Name)
			if z.err != nil ***REMOVED***
				return n, z.err
			***REMOVED***
		***REMOVED***
		if z.Comment != "" ***REMOVED***
			z.err = z.writeString(z.Comment)
			if z.err != nil ***REMOVED***
				return n, z.err
			***REMOVED***
		***REMOVED***

		if z.compressor == nil && z.level != StatelessCompression ***REMOVED***
			z.compressor, _ = flate.NewWriter(z.w, z.level)
		***REMOVED***
	***REMOVED***
	z.size += uint32(len(p))
	z.digest = crc32.Update(z.digest, crc32.IEEETable, p)
	if z.level == StatelessCompression ***REMOVED***
		return len(p), flate.StatelessDeflate(z.w, p, false, nil)
	***REMOVED***
	n, z.err = z.compressor.Write(p)
	return n, z.err
***REMOVED***

// Flush flushes any pending compressed data to the underlying writer.
//
// It is useful mainly in compressed network protocols, to ensure that
// a remote reader has enough data to reconstruct a packet. Flush does
// not return until the data has been written. If the underlying
// writer returns an error, Flush returns that error.
//
// In the terminology of the zlib library, Flush is equivalent to Z_SYNC_FLUSH.
func (z *Writer) Flush() error ***REMOVED***
	if z.err != nil ***REMOVED***
		return z.err
	***REMOVED***
	if z.closed || z.level == StatelessCompression ***REMOVED***
		return nil
	***REMOVED***
	if !z.wroteHeader ***REMOVED***
		z.Write(nil)
		if z.err != nil ***REMOVED***
			return z.err
		***REMOVED***
	***REMOVED***
	z.err = z.compressor.Flush()
	return z.err
***REMOVED***

// Close closes the Writer, flushing any unwritten data to the underlying
// io.Writer, but does not close the underlying io.Writer.
func (z *Writer) Close() error ***REMOVED***
	if z.err != nil ***REMOVED***
		return z.err
	***REMOVED***
	if z.closed ***REMOVED***
		return nil
	***REMOVED***
	z.closed = true
	if !z.wroteHeader ***REMOVED***
		z.Write(nil)
		if z.err != nil ***REMOVED***
			return z.err
		***REMOVED***
	***REMOVED***
	if z.level == StatelessCompression ***REMOVED***
		z.err = flate.StatelessDeflate(z.w, nil, true, nil)
	***REMOVED*** else ***REMOVED***
		z.err = z.compressor.Close()
	***REMOVED***
	if z.err != nil ***REMOVED***
		return z.err
	***REMOVED***
	le.PutUint32(z.buf[:4], z.digest)
	le.PutUint32(z.buf[4:8], z.size)
	_, z.err = z.w.Write(z.buf[:8])
	return z.err
***REMOVED***
