// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package md4 implements the MD4 hash algorithm as defined in RFC 1320.
//
// Deprecated: MD4 is cryptographically broken and should should only be used
// where compatibility with legacy systems, not security, is the goal. Instead,
// use a secure hash like SHA-256 (from crypto/sha256).
package md4 // import "golang.org/x/crypto/md4"

import (
	"crypto"
	"hash"
)

func init() ***REMOVED***
	crypto.RegisterHash(crypto.MD4, New)
***REMOVED***

// The size of an MD4 checksum in bytes.
const Size = 16

// The blocksize of MD4 in bytes.
const BlockSize = 64

const (
	_Chunk = 64
	_Init0 = 0x67452301
	_Init1 = 0xEFCDAB89
	_Init2 = 0x98BADCFE
	_Init3 = 0x10325476
)

// digest represents the partial evaluation of a checksum.
type digest struct ***REMOVED***
	s   [4]uint32
	x   [_Chunk]byte
	nx  int
	len uint64
***REMOVED***

func (d *digest) Reset() ***REMOVED***
	d.s[0] = _Init0
	d.s[1] = _Init1
	d.s[2] = _Init2
	d.s[3] = _Init3
	d.nx = 0
	d.len = 0
***REMOVED***

// New returns a new hash.Hash computing the MD4 checksum.
func New() hash.Hash ***REMOVED***
	d := new(digest)
	d.Reset()
	return d
***REMOVED***

func (d *digest) Size() int ***REMOVED*** return Size ***REMOVED***

func (d *digest) BlockSize() int ***REMOVED*** return BlockSize ***REMOVED***

func (d *digest) Write(p []byte) (nn int, err error) ***REMOVED***
	nn = len(p)
	d.len += uint64(nn)
	if d.nx > 0 ***REMOVED***
		n := len(p)
		if n > _Chunk-d.nx ***REMOVED***
			n = _Chunk - d.nx
		***REMOVED***
		for i := 0; i < n; i++ ***REMOVED***
			d.x[d.nx+i] = p[i]
		***REMOVED***
		d.nx += n
		if d.nx == _Chunk ***REMOVED***
			_Block(d, d.x[0:])
			d.nx = 0
		***REMOVED***
		p = p[n:]
	***REMOVED***
	n := _Block(d, p)
	p = p[n:]
	if len(p) > 0 ***REMOVED***
		d.nx = copy(d.x[:], p)
	***REMOVED***
	return
***REMOVED***

func (d0 *digest) Sum(in []byte) []byte ***REMOVED***
	// Make a copy of d0, so that caller can keep writing and summing.
	d := new(digest)
	*d = *d0

	// Padding.  Add a 1 bit and 0 bits until 56 bytes mod 64.
	len := d.len
	var tmp [64]byte
	tmp[0] = 0x80
	if len%64 < 56 ***REMOVED***
		d.Write(tmp[0 : 56-len%64])
	***REMOVED*** else ***REMOVED***
		d.Write(tmp[0 : 64+56-len%64])
	***REMOVED***

	// Length in bits.
	len <<= 3
	for i := uint(0); i < 8; i++ ***REMOVED***
		tmp[i] = byte(len >> (8 * i))
	***REMOVED***
	d.Write(tmp[0:8])

	if d.nx != 0 ***REMOVED***
		panic("d.nx != 0")
	***REMOVED***

	for _, s := range d.s ***REMOVED***
		in = append(in, byte(s>>0))
		in = append(in, byte(s>>8))
		in = append(in, byte(s>>16))
		in = append(in, byte(s>>24))
	***REMOVED***
	return in
***REMOVED***
