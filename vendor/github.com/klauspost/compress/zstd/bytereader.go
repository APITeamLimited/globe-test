// Copyright 2019+ Klaus Post. All rights reserved.
// License information can be found in the LICENSE file.
// Based on work by Yann Collet, released under BSD License.

package zstd

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

// overread returns whether we have advanced too far.
func (b *byteReader) overread() bool ***REMOVED***
	return b.off > len(b.b)
***REMOVED***

// Int32 returns a little endian int32 starting at current offset.
func (b byteReader) Int32() int32 ***REMOVED***
	b2 := b.b[b.off:]
	b2 = b2[:4]
	v3 := int32(b2[3])
	v2 := int32(b2[2])
	v1 := int32(b2[1])
	v0 := int32(b2[0])
	return v0 | (v1 << 8) | (v2 << 16) | (v3 << 24)
***REMOVED***

// Uint8 returns the next byte
func (b *byteReader) Uint8() uint8 ***REMOVED***
	v := b.b[b.off]
	return v
***REMOVED***

// Uint32 returns a little endian uint32 starting at current offset.
func (b byteReader) Uint32() uint32 ***REMOVED***
	if r := b.remain(); r < 4 ***REMOVED***
		// Very rare
		v := uint32(0)
		for i := 1; i <= r; i++ ***REMOVED***
			v = (v << 8) | uint32(b.b[len(b.b)-i])
		***REMOVED***
		return v
	***REMOVED***
	b2 := b.b[b.off:]
	b2 = b2[:4]
	v3 := uint32(b2[3])
	v2 := uint32(b2[2])
	v1 := uint32(b2[1])
	v0 := uint32(b2[0])
	return v0 | (v1 << 8) | (v2 << 16) | (v3 << 24)
***REMOVED***

// Uint32NC returns a little endian uint32 starting at current offset.
// The caller must be sure if there are at least 4 bytes left.
func (b byteReader) Uint32NC() uint32 ***REMOVED***
	b2 := b.b[b.off:]
	b2 = b2[:4]
	v3 := uint32(b2[3])
	v2 := uint32(b2[2])
	v1 := uint32(b2[1])
	v0 := uint32(b2[0])
	return v0 | (v1 << 8) | (v2 << 16) | (v3 << 24)
***REMOVED***

// unread returns the unread portion of the input.
func (b byteReader) unread() []byte ***REMOVED***
	return b.b[b.off:]
***REMOVED***

// remain will return the number of bytes remaining.
func (b byteReader) remain() int ***REMOVED***
	return len(b.b) - b.off
***REMOVED***
