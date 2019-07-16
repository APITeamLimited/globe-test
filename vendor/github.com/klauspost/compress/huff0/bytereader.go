// Copyright 2018 Klaus Post. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Based on work Copyright (c) 2013, Yann Collet, released under BSD License.

package huff0

// byteReader provides a byte reader that reads
// little endian values from a byte stream.
// The input stream is manually advanced.
// The reader performs no bounds checks.
type byteReader struct ***REMOVED***
	b   []byte
	off int
***REMOVED***

// init will initialize the reader and set the input.
func (b *byteReader) init(in []byte) ***REMOVED***
	b.b = in
	b.off = 0
***REMOVED***

// advance the stream b n bytes.
func (b *byteReader) advance(n uint) ***REMOVED***
	b.off += int(n)
***REMOVED***

// Int32 returns a little endian int32 starting at current offset.
func (b byteReader) Int32() int32 ***REMOVED***
	v3 := int32(b.b[b.off+3])
	v2 := int32(b.b[b.off+2])
	v1 := int32(b.b[b.off+1])
	v0 := int32(b.b[b.off])
	return (v3 << 24) | (v2 << 16) | (v1 << 8) | v0
***REMOVED***

// Uint32 returns a little endian uint32 starting at current offset.
func (b byteReader) Uint32() uint32 ***REMOVED***
	v3 := uint32(b.b[b.off+3])
	v2 := uint32(b.b[b.off+2])
	v1 := uint32(b.b[b.off+1])
	v0 := uint32(b.b[b.off])
	return (v3 << 24) | (v2 << 16) | (v1 << 8) | v0
***REMOVED***

// unread returns the unread portion of the input.
func (b byteReader) unread() []byte ***REMOVED***
	return b.b[b.off:]
***REMOVED***

// remain will return the number of bytes remaining.
func (b byteReader) remain() int ***REMOVED***
	return len(b.b) - b.off
***REMOVED***
