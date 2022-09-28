// Copyright 2016 The Snappy-Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !amd64 appengine !gc noasm

package snappy

// decode writes the decoding of src to dst. It assumes that the varint-encoded
// length of the decompressed bytes has already been read, and that len(dst)
// equals that length.
//
// It returns 0 on success or a decodeErrCodeXxx error code on failure.
func decode(dst, src []byte) int ***REMOVED***
	var d, s, offset, length int
	for s < len(src) ***REMOVED***
		switch src[s] & 0x03 ***REMOVED***
		case tagLiteral:
			x := uint32(src[s] >> 2)
			switch ***REMOVED***
			case x < 60:
				s++
			case x == 60:
				s += 2
				if uint(s) > uint(len(src)) ***REMOVED*** // The uint conversions catch overflow from the previous line.
					return decodeErrCodeCorrupt
				***REMOVED***
				x = uint32(src[s-1])
			case x == 61:
				s += 3
				if uint(s) > uint(len(src)) ***REMOVED*** // The uint conversions catch overflow from the previous line.
					return decodeErrCodeCorrupt
				***REMOVED***
				x = uint32(src[s-2]) | uint32(src[s-1])<<8
			case x == 62:
				s += 4
				if uint(s) > uint(len(src)) ***REMOVED*** // The uint conversions catch overflow from the previous line.
					return decodeErrCodeCorrupt
				***REMOVED***
				x = uint32(src[s-3]) | uint32(src[s-2])<<8 | uint32(src[s-1])<<16
			case x == 63:
				s += 5
				if uint(s) > uint(len(src)) ***REMOVED*** // The uint conversions catch overflow from the previous line.
					return decodeErrCodeCorrupt
				***REMOVED***
				x = uint32(src[s-4]) | uint32(src[s-3])<<8 | uint32(src[s-2])<<16 | uint32(src[s-1])<<24
			***REMOVED***
			length = int(x) + 1
			if length <= 0 ***REMOVED***
				return decodeErrCodeUnsupportedLiteralLength
			***REMOVED***
			if length > len(dst)-d || length > len(src)-s ***REMOVED***
				return decodeErrCodeCorrupt
			***REMOVED***
			copy(dst[d:], src[s:s+length])
			d += length
			s += length
			continue

		case tagCopy1:
			s += 2
			if uint(s) > uint(len(src)) ***REMOVED*** // The uint conversions catch overflow from the previous line.
				return decodeErrCodeCorrupt
			***REMOVED***
			length = 4 + int(src[s-2])>>2&0x7
			offset = int(uint32(src[s-2])&0xe0<<3 | uint32(src[s-1]))

		case tagCopy2:
			s += 3
			if uint(s) > uint(len(src)) ***REMOVED*** // The uint conversions catch overflow from the previous line.
				return decodeErrCodeCorrupt
			***REMOVED***
			length = 1 + int(src[s-3])>>2
			offset = int(uint32(src[s-2]) | uint32(src[s-1])<<8)

		case tagCopy4:
			s += 5
			if uint(s) > uint(len(src)) ***REMOVED*** // The uint conversions catch overflow from the previous line.
				return decodeErrCodeCorrupt
			***REMOVED***
			length = 1 + int(src[s-5])>>2
			offset = int(uint32(src[s-4]) | uint32(src[s-3])<<8 | uint32(src[s-2])<<16 | uint32(src[s-1])<<24)
		***REMOVED***

		if offset <= 0 || d < offset || length > len(dst)-d ***REMOVED***
			return decodeErrCodeCorrupt
		***REMOVED***
		// Copy from an earlier sub-slice of dst to a later sub-slice. Unlike
		// the built-in copy function, this byte-by-byte copy always runs
		// forwards, even if the slices overlap. Conceptually, this is:
		//
		// d += forwardCopy(dst[d:d+length], dst[d-offset:])
		for end := d + length; d != end; d++ ***REMOVED***
			dst[d] = dst[d-offset]
		***REMOVED***
	***REMOVED***
	if d != len(dst) ***REMOVED***
		return decodeErrCodeCorrupt
	***REMOVED***
	return 0
***REMOVED***
