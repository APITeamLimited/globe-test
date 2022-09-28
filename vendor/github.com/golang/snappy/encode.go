// Copyright 2011 The Snappy-Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package snappy

import (
	"encoding/binary"
	"errors"
	"io"
)

// Encode returns the encoded form of src. The returned slice may be a sub-
// slice of dst if dst was large enough to hold the entire encoded block.
// Otherwise, a newly allocated slice will be returned.
//
// The dst and src must not overlap. It is valid to pass a nil dst.
func Encode(dst, src []byte) []byte ***REMOVED***
	if n := MaxEncodedLen(len(src)); n < 0 ***REMOVED***
		panic(ErrTooLarge)
	***REMOVED*** else if len(dst) < n ***REMOVED***
		dst = make([]byte, n)
	***REMOVED***

	// The block starts with the varint-encoded length of the decompressed bytes.
	d := binary.PutUvarint(dst, uint64(len(src)))

	for len(src) > 0 ***REMOVED***
		p := src
		src = nil
		if len(p) > maxBlockSize ***REMOVED***
			p, src = p[:maxBlockSize], p[maxBlockSize:]
		***REMOVED***
		if len(p) < minNonLiteralBlockSize ***REMOVED***
			d += emitLiteral(dst[d:], p)
		***REMOVED*** else ***REMOVED***
			d += encodeBlock(dst[d:], p)
		***REMOVED***
	***REMOVED***
	return dst[:d]
***REMOVED***

// inputMargin is the minimum number of extra input bytes to keep, inside
// encodeBlock's inner loop. On some architectures, this margin lets us
// implement a fast path for emitLiteral, where the copy of short (<= 16 byte)
// literals can be implemented as a single load to and store from a 16-byte
// register. That literal's actual length can be as short as 1 byte, so this
// can copy up to 15 bytes too much, but that's OK as subsequent iterations of
// the encoding loop will fix up the copy overrun, and this inputMargin ensures
// that we don't overrun the dst and src buffers.
const inputMargin = 16 - 1

// minNonLiteralBlockSize is the minimum size of the input to encodeBlock that
// could be encoded with a copy tag. This is the minimum with respect to the
// algorithm used by encodeBlock, not a minimum enforced by the file format.
//
// The encoded output must start with at least a 1 byte literal, as there are
// no previous bytes to copy. A minimal (1 byte) copy after that, generated
// from an emitCopy call in encodeBlock's main loop, would require at least
// another inputMargin bytes, for the reason above: we want any emitLiteral
// calls inside encodeBlock's main loop to use the fast path if possible, which
// requires being able to overrun by inputMargin bytes. Thus,
// minNonLiteralBlockSize equals 1 + 1 + inputMargin.
//
// The C++ code doesn't use this exact threshold, but it could, as discussed at
// https://groups.google.com/d/topic/snappy-compression/oGbhsdIJSJ8/discussion
// The difference between Go (2+inputMargin) and C++ (inputMargin) is purely an
// optimization. It should not affect the encoded form. This is tested by
// TestSameEncodingAsCppShortCopies.
const minNonLiteralBlockSize = 1 + 1 + inputMargin

// MaxEncodedLen returns the maximum length of a snappy block, given its
// uncompressed length.
//
// It will return a negative value if srcLen is too large to encode.
func MaxEncodedLen(srcLen int) int ***REMOVED***
	n := uint64(srcLen)
	if n > 0xffffffff ***REMOVED***
		return -1
	***REMOVED***
	// Compressed data can be defined as:
	//    compressed := item* literal*
	//    item       := literal* copy
	//
	// The trailing literal sequence has a space blowup of at most 62/60
	// since a literal of length 60 needs one tag byte + one extra byte
	// for length information.
	//
	// Item blowup is trickier to measure. Suppose the "copy" op copies
	// 4 bytes of data. Because of a special check in the encoding code,
	// we produce a 4-byte copy only if the offset is < 65536. Therefore
	// the copy op takes 3 bytes to encode, and this type of item leads
	// to at most the 62/60 blowup for representing literals.
	//
	// Suppose the "copy" op copies 5 bytes of data. If the offset is big
	// enough, it will take 5 bytes to encode the copy op. Therefore the
	// worst case here is a one-byte literal followed by a five-byte copy.
	// That is, 6 bytes of input turn into 7 bytes of "compressed" data.
	//
	// This last factor dominates the blowup, so the final estimate is:
	n = 32 + n + n/6
	if n > 0xffffffff ***REMOVED***
		return -1
	***REMOVED***
	return int(n)
***REMOVED***

var errClosed = errors.New("snappy: Writer is closed")

// NewWriter returns a new Writer that compresses to w.
//
// The Writer returned does not buffer writes. There is no need to Flush or
// Close such a Writer.
//
// Deprecated: the Writer returned is not suitable for many small writes, only
// for few large writes. Use NewBufferedWriter instead, which is efficient
// regardless of the frequency and shape of the writes, and remember to Close
// that Writer when done.
func NewWriter(w io.Writer) *Writer ***REMOVED***
	return &Writer***REMOVED***
		w:    w,
		obuf: make([]byte, obufLen),
	***REMOVED***
***REMOVED***

// NewBufferedWriter returns a new Writer that compresses to w, using the
// framing format described at
// https://github.com/google/snappy/blob/master/framing_format.txt
//
// The Writer returned buffers writes. Users must call Close to guarantee all
// data has been forwarded to the underlying io.Writer. They may also call
// Flush zero or more times before calling Close.
func NewBufferedWriter(w io.Writer) *Writer ***REMOVED***
	return &Writer***REMOVED***
		w:    w,
		ibuf: make([]byte, 0, maxBlockSize),
		obuf: make([]byte, obufLen),
	***REMOVED***
***REMOVED***

// Writer is an io.Writer that can write Snappy-compressed bytes.
type Writer struct ***REMOVED***
	w   io.Writer
	err error

	// ibuf is a buffer for the incoming (uncompressed) bytes.
	//
	// Its use is optional. For backwards compatibility, Writers created by the
	// NewWriter function have ibuf == nil, do not buffer incoming bytes, and
	// therefore do not need to be Flush'ed or Close'd.
	ibuf []byte

	// obuf is a buffer for the outgoing (compressed) bytes.
	obuf []byte

	// wroteStreamHeader is whether we have written the stream header.
	wroteStreamHeader bool
***REMOVED***

// Reset discards the writer's state and switches the Snappy writer to write to
// w. This permits reusing a Writer rather than allocating a new one.
func (w *Writer) Reset(writer io.Writer) ***REMOVED***
	w.w = writer
	w.err = nil
	if w.ibuf != nil ***REMOVED***
		w.ibuf = w.ibuf[:0]
	***REMOVED***
	w.wroteStreamHeader = false
***REMOVED***

// Write satisfies the io.Writer interface.
func (w *Writer) Write(p []byte) (nRet int, errRet error) ***REMOVED***
	if w.ibuf == nil ***REMOVED***
		// Do not buffer incoming bytes. This does not perform or compress well
		// if the caller of Writer.Write writes many small slices. This
		// behavior is therefore deprecated, but still supported for backwards
		// compatibility with code that doesn't explicitly Flush or Close.
		return w.write(p)
	***REMOVED***

	// The remainder of this method is based on bufio.Writer.Write from the
	// standard library.

	for len(p) > (cap(w.ibuf)-len(w.ibuf)) && w.err == nil ***REMOVED***
		var n int
		if len(w.ibuf) == 0 ***REMOVED***
			// Large write, empty buffer.
			// Write directly from p to avoid copy.
			n, _ = w.write(p)
		***REMOVED*** else ***REMOVED***
			n = copy(w.ibuf[len(w.ibuf):cap(w.ibuf)], p)
			w.ibuf = w.ibuf[:len(w.ibuf)+n]
			w.Flush()
		***REMOVED***
		nRet += n
		p = p[n:]
	***REMOVED***
	if w.err != nil ***REMOVED***
		return nRet, w.err
	***REMOVED***
	n := copy(w.ibuf[len(w.ibuf):cap(w.ibuf)], p)
	w.ibuf = w.ibuf[:len(w.ibuf)+n]
	nRet += n
	return nRet, nil
***REMOVED***

func (w *Writer) write(p []byte) (nRet int, errRet error) ***REMOVED***
	if w.err != nil ***REMOVED***
		return 0, w.err
	***REMOVED***
	for len(p) > 0 ***REMOVED***
		obufStart := len(magicChunk)
		if !w.wroteStreamHeader ***REMOVED***
			w.wroteStreamHeader = true
			copy(w.obuf, magicChunk)
			obufStart = 0
		***REMOVED***

		var uncompressed []byte
		if len(p) > maxBlockSize ***REMOVED***
			uncompressed, p = p[:maxBlockSize], p[maxBlockSize:]
		***REMOVED*** else ***REMOVED***
			uncompressed, p = p, nil
		***REMOVED***
		checksum := crc(uncompressed)

		// Compress the buffer, discarding the result if the improvement
		// isn't at least 12.5%.
		compressed := Encode(w.obuf[obufHeaderLen:], uncompressed)
		chunkType := uint8(chunkTypeCompressedData)
		chunkLen := 4 + len(compressed)
		obufEnd := obufHeaderLen + len(compressed)
		if len(compressed) >= len(uncompressed)-len(uncompressed)/8 ***REMOVED***
			chunkType = chunkTypeUncompressedData
			chunkLen = 4 + len(uncompressed)
			obufEnd = obufHeaderLen
		***REMOVED***

		// Fill in the per-chunk header that comes before the body.
		w.obuf[len(magicChunk)+0] = chunkType
		w.obuf[len(magicChunk)+1] = uint8(chunkLen >> 0)
		w.obuf[len(magicChunk)+2] = uint8(chunkLen >> 8)
		w.obuf[len(magicChunk)+3] = uint8(chunkLen >> 16)
		w.obuf[len(magicChunk)+4] = uint8(checksum >> 0)
		w.obuf[len(magicChunk)+5] = uint8(checksum >> 8)
		w.obuf[len(magicChunk)+6] = uint8(checksum >> 16)
		w.obuf[len(magicChunk)+7] = uint8(checksum >> 24)

		if _, err := w.w.Write(w.obuf[obufStart:obufEnd]); err != nil ***REMOVED***
			w.err = err
			return nRet, err
		***REMOVED***
		if chunkType == chunkTypeUncompressedData ***REMOVED***
			if _, err := w.w.Write(uncompressed); err != nil ***REMOVED***
				w.err = err
				return nRet, err
			***REMOVED***
		***REMOVED***
		nRet += len(uncompressed)
	***REMOVED***
	return nRet, nil
***REMOVED***

// Flush flushes the Writer to its underlying io.Writer.
func (w *Writer) Flush() error ***REMOVED***
	if w.err != nil ***REMOVED***
		return w.err
	***REMOVED***
	if len(w.ibuf) == 0 ***REMOVED***
		return nil
	***REMOVED***
	w.write(w.ibuf)
	w.ibuf = w.ibuf[:0]
	return w.err
***REMOVED***

// Close calls Flush and then closes the Writer.
func (w *Writer) Close() error ***REMOVED***
	w.Flush()
	ret := w.err
	if w.err == nil ***REMOVED***
		w.err = errClosed
	***REMOVED***
	return ret
***REMOVED***
