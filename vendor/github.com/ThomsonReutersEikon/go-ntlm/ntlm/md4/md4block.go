// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// MD4 block step.
// In its own file so that a faster assembly or C version
// can be substituted easily.

package md4

var shift1 = []uint***REMOVED***3, 7, 11, 19***REMOVED***
var shift2 = []uint***REMOVED***3, 5, 9, 13***REMOVED***
var shift3 = []uint***REMOVED***3, 9, 11, 15***REMOVED***

var xIndex2 = []uint***REMOVED***0, 4, 8, 12, 1, 5, 9, 13, 2, 6, 10, 14, 3, 7, 11, 15***REMOVED***
var xIndex3 = []uint***REMOVED***0, 8, 4, 12, 2, 10, 6, 14, 1, 9, 5, 13, 3, 11, 7, 15***REMOVED***

func _Block(dig *digest, p []byte) int ***REMOVED***
	a := dig.s[0]
	b := dig.s[1]
	c := dig.s[2]
	d := dig.s[3]
	n := 0
	var X [16]uint32
	for len(p) >= _Chunk ***REMOVED***
		aa, bb, cc, dd := a, b, c, d

		j := 0
		for i := 0; i < 16; i++ ***REMOVED***
			X[i] = uint32(p[j]) | uint32(p[j+1])<<8 | uint32(p[j+2])<<16 | uint32(p[j+3])<<24
			j += 4
		***REMOVED***

		// If this needs to be made faster in the future,
		// the usual trick is to unroll each of these
		// loops by a factor of 4; that lets you replace
		// the shift[] lookups with constants and,
		// with suitable variable renaming in each
		// unrolled body, delete the a, b, c, d = d, a, b, c
		// (or you can let the optimizer do the renaming).
		//
		// The index variables are uint so that % by a power
		// of two can be optimized easily by a compiler.

		// Round 1.
		for i := uint(0); i < 16; i++ ***REMOVED***
			x := i
			s := shift1[i%4]
			f := ((c ^ d) & b) ^ d
			a += f + X[x]
			a = a<<s | a>>(32-s)
			a, b, c, d = d, a, b, c
		***REMOVED***

		// Round 2.
		for i := uint(0); i < 16; i++ ***REMOVED***
			x := xIndex2[i]
			s := shift2[i%4]
			g := (b & c) | (b & d) | (c & d)
			a += g + X[x] + 0x5a827999
			a = a<<s | a>>(32-s)
			a, b, c, d = d, a, b, c
		***REMOVED***

		// Round 3.
		for i := uint(0); i < 16; i++ ***REMOVED***
			x := xIndex3[i]
			s := shift3[i%4]
			h := b ^ c ^ d
			a += h + X[x] + 0x6ed9eba1
			a = a<<s | a>>(32-s)
			a, b, c, d = d, a, b, c
		***REMOVED***

		a += aa
		b += bb
		c += cc
		d += dd

		p = p[_Chunk:]
		n += _Chunk
	***REMOVED***

	dig.s[0] = a
	dig.s[1] = b
	dig.s[2] = c
	dig.s[3] = d
	return n
***REMOVED***
