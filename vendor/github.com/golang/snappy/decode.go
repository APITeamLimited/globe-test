// Copyright 2011 The Snappy-Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package snappy

import (
	"encoding/binary"
	"errors"
	"io"
)

var (
	// ErrCorrupt reports that the input is invalid.
	ErrCorrupt = errors.New("snappy: corrupt input")
	// ErrTooLarge reports that the uncompressed length is too large.
	ErrTooLarge = errors.New("snappy: decoded block is too large")
	// ErrUnsupported reports that the input isn't supported.
	ErrUnsupported = errors.New("snappy: unsupported input")

	errUnsupportedLiteralLength = errors.New("snappy: unsupported literal length")
)

// DecodedLen returns the length of the decoded block.
func DecodedLen(src []byte) (int, error) ***REMOVED***
	v, _, err := decodedLen(src)
	return v, err
***REMOVED***

// decodedLen returns the length of the decoded block and the number of bytes
// that the length header occupied.
func decodedLen(src []byte) (blockLen, headerLen int, err error) ***REMOVED***
	v, n := binary.Uvarint(src)
	if n <= 0 || v > 0xffffffff ***REMOVED***
		return 0, 0, ErrCorrupt
	***REMOVED***

	const wordSize = 32 << (^uint(0) >> 32 & 1)
	if wordSize == 32 && v > 0x7fffffff ***REMOVED***
		return 0, 0, ErrTooLarge
	***REMOVED***
	return int(v), n, nil
***REMOVED***

const (
	decodeErrCodeCorrupt                  = 1
	decodeErrCodeUnsupportedLiteralLength = 2
)

// Decode returns the decoded form of src. The returned slice may be a sub-
// slice of dst if dst was large enough to hold the entire decoded block.
// Otherwise, a newly allocated slice will be returned.
//
// The dst and src must not overlap. It is valid to pass a nil dst.
func Decode(dst, src []byte) ([]byte, error) ***REMOVED***
	dLen, s, err := decodedLen(src)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if dLen <= len(dst) ***REMOVED***
		dst = dst[:dLen]
	***REMOVED*** else ***REMOVED***
		dst = make([]byte, dLen)
	***REMOVED***
	switch decode(dst, src[s:]) ***REMOVED***
	case 0:
		return dst, nil
	case decodeErrCodeUnsupportedLiteralLength:
		return nil, errUnsupportedLiteralLength
	***REMOVED***
	return nil, ErrCorrupt
***REMOVED***

// NewReader returns a new Reader that decompresses from r, using the framing
// format described at
// https://github.com/google/snappy/blob/master/framing_format.txt
func NewReader(r io.Reader) *Reader ***REMOVED***
	return &Reader***REMOVED***
		r:       r,
		decoded: make([]byte, maxBlockSize),
		buf:     make([]byte, maxEncodedLenOfMaxBlockSize+checksumSize),
	***REMOVED***
***REMOVED***

// Reader is an io.Reader that can read Snappy-compressed bytes.
type Reader struct ***REMOVED***
	r       io.Reader
	err     error
	decoded []byte
	buf     []byte
	// decoded[i:j] contains decoded bytes that have not yet been passed on.
	i, j       int
	readHeader bool
***REMOVED***

// Reset discards any buffered data, resets all state, and switches the Snappy
// reader to read from r. This permits reusing a Reader rather than allocating
// a new one.
func (r *Reader) Reset(reader io.Reader) ***REMOVED***
	r.r = reader
	r.err = nil
	r.i = 0
	r.j = 0
	r.readHeader = false
***REMOVED***

func (r *Reader) readFull(p []byte, allowEOF bool) (ok bool) ***REMOVED***
	if _, r.err = io.ReadFull(r.r, p); r.err != nil ***REMOVED***
		if r.err == io.ErrUnexpectedEOF || (r.err == io.EOF && !allowEOF) ***REMOVED***
			r.err = ErrCorrupt
		***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

// Read satisfies the io.Reader interface.
func (r *Reader) Read(p []byte) (int, error) ***REMOVED***
	if r.err != nil ***REMOVED***
		return 0, r.err
	***REMOVED***
	for ***REMOVED***
		if r.i < r.j ***REMOVED***
			n := copy(p, r.decoded[r.i:r.j])
			r.i += n
			return n, nil
		***REMOVED***
		if !r.readFull(r.buf[:4], true) ***REMOVED***
			return 0, r.err
		***REMOVED***
		chunkType := r.buf[0]
		if !r.readHeader ***REMOVED***
			if chunkType != chunkTypeStreamIdentifier ***REMOVED***
				r.err = ErrCorrupt
				return 0, r.err
			***REMOVED***
			r.readHeader = true
		***REMOVED***
		chunkLen := int(r.buf[1]) | int(r.buf[2])<<8 | int(r.buf[3])<<16
		if chunkLen > len(r.buf) ***REMOVED***
			r.err = ErrUnsupported
			return 0, r.err
		***REMOVED***

		// The chunk types are specified at
		// https://github.com/google/snappy/blob/master/framing_format.txt
		switch chunkType ***REMOVED***
		case chunkTypeCompressedData:
			// Section 4.2. Compressed data (chunk type 0x00).
			if chunkLen < checksumSize ***REMOVED***
				r.err = ErrCorrupt
				return 0, r.err
			***REMOVED***
			buf := r.buf[:chunkLen]
			if !r.readFull(buf, false) ***REMOVED***
				return 0, r.err
			***REMOVED***
			checksum := uint32(buf[0]) | uint32(buf[1])<<8 | uint32(buf[2])<<16 | uint32(buf[3])<<24
			buf = buf[checksumSize:]

			n, err := DecodedLen(buf)
			if err != nil ***REMOVED***
				r.err = err
				return 0, r.err
			***REMOVED***
			if n > len(r.decoded) ***REMOVED***
				r.err = ErrCorrupt
				return 0, r.err
			***REMOVED***
			if _, err := Decode(r.decoded, buf); err != nil ***REMOVED***
				r.err = err
				return 0, r.err
			***REMOVED***
			if crc(r.decoded[:n]) != checksum ***REMOVED***
				r.err = ErrCorrupt
				return 0, r.err
			***REMOVED***
			r.i, r.j = 0, n
			continue

		case chunkTypeUncompressedData:
			// Section 4.3. Uncompressed data (chunk type 0x01).
			if chunkLen < checksumSize ***REMOVED***
				r.err = ErrCorrupt
				return 0, r.err
			***REMOVED***
			buf := r.buf[:checksumSize]
			if !r.readFull(buf, false) ***REMOVED***
				return 0, r.err
			***REMOVED***
			checksum := uint32(buf[0]) | uint32(buf[1])<<8 | uint32(buf[2])<<16 | uint32(buf[3])<<24
			// Read directly into r.decoded instead of via r.buf.
			n := chunkLen - checksumSize
			if n > len(r.decoded) ***REMOVED***
				r.err = ErrCorrupt
				return 0, r.err
			***REMOVED***
			if !r.readFull(r.decoded[:n], false) ***REMOVED***
				return 0, r.err
			***REMOVED***
			if crc(r.decoded[:n]) != checksum ***REMOVED***
				r.err = ErrCorrupt
				return 0, r.err
			***REMOVED***
			r.i, r.j = 0, n
			continue

		case chunkTypeStreamIdentifier:
			// Section 4.1. Stream identifier (chunk type 0xff).
			if chunkLen != len(magicBody) ***REMOVED***
				r.err = ErrCorrupt
				return 0, r.err
			***REMOVED***
			if !r.readFull(r.buf[:len(magicBody)], false) ***REMOVED***
				return 0, r.err
			***REMOVED***
			for i := 0; i < len(magicBody); i++ ***REMOVED***
				if r.buf[i] != magicBody[i] ***REMOVED***
					r.err = ErrCorrupt
					return 0, r.err
				***REMOVED***
			***REMOVED***
			continue
		***REMOVED***

		if chunkType <= 0x7f ***REMOVED***
			// Section 4.5. Reserved unskippable chunks (chunk types 0x02-0x7f).
			r.err = ErrUnsupported
			return 0, r.err
		***REMOVED***
		// Section 4.4 Padding (chunk type 0xfe).
		// Section 4.6. Reserved skippable chunks (chunk types 0x80-0xfd).
		if !r.readFull(r.buf[:chunkLen], false) ***REMOVED***
			return 0, r.err
		***REMOVED***
	***REMOVED***
***REMOVED***
