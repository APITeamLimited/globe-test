// Package xxHash32 implements the very fast xxHash hashing algorithm (32 bits version).
// (https://github.com/Cyan4973/xxHash/)
package xxHash32

import "hash"

const (
	prime32_1 = 2654435761
	prime32_2 = 2246822519
	prime32_3 = 3266489917
	prime32_4 = 668265263
	prime32_5 = 374761393
)

type xxHash struct ***REMOVED***
	seed     uint32
	v1       uint32
	v2       uint32
	v3       uint32
	v4       uint32
	totalLen uint64
	buf      [16]byte
	bufused  int
***REMOVED***

// New returns a new Hash32 instance.
func New(seed uint32) hash.Hash32 ***REMOVED***
	xxh := &xxHash***REMOVED***seed: seed***REMOVED***
	xxh.Reset()
	return xxh
***REMOVED***

// Sum appends the current hash to b and returns the resulting slice.
// It does not change the underlying hash state.
func (xxh xxHash) Sum(b []byte) []byte ***REMOVED***
	h32 := xxh.Sum32()
	return append(b, byte(h32), byte(h32>>8), byte(h32>>16), byte(h32>>24))
***REMOVED***

// Reset resets the Hash to its initial state.
func (xxh *xxHash) Reset() ***REMOVED***
	xxh.v1 = xxh.seed + prime32_1 + prime32_2
	xxh.v2 = xxh.seed + prime32_2
	xxh.v3 = xxh.seed
	xxh.v4 = xxh.seed - prime32_1
	xxh.totalLen = 0
	xxh.bufused = 0
***REMOVED***

// Size returns the number of bytes returned by Sum().
func (xxh *xxHash) Size() int ***REMOVED***
	return 4
***REMOVED***

// BlockSize gives the minimum number of bytes accepted by Write().
func (xxh *xxHash) BlockSize() int ***REMOVED***
	return 1
***REMOVED***

// Write adds input bytes to the Hash.
// It never returns an error.
func (xxh *xxHash) Write(input []byte) (int, error) ***REMOVED***
	n := len(input)
	m := xxh.bufused

	xxh.totalLen += uint64(n)

	r := len(xxh.buf) - m
	if n < r ***REMOVED***
		copy(xxh.buf[m:], input)
		xxh.bufused += len(input)
		return n, nil
	***REMOVED***

	p := 0
	if m > 0 ***REMOVED***
		// some data left from previous update
		copy(xxh.buf[xxh.bufused:], input[:r])
		xxh.bufused += len(input) - r

		// fast rotl(13)
		p32 := xxh.v1 + (uint32(xxh.buf[p+3])<<24|uint32(xxh.buf[p+2])<<16|uint32(xxh.buf[p+1])<<8|uint32(xxh.buf[p]))*prime32_2
		xxh.v1 = (p32<<13 | p32>>19) * prime32_1
		p += 4
		p32 = xxh.v2 + (uint32(xxh.buf[p+3])<<24|uint32(xxh.buf[p+2])<<16|uint32(xxh.buf[p+1])<<8|uint32(xxh.buf[p]))*prime32_2
		xxh.v2 = (p32<<13 | p32>>19) * prime32_1
		p += 4
		p32 = xxh.v3 + (uint32(xxh.buf[p+3])<<24|uint32(xxh.buf[p+2])<<16|uint32(xxh.buf[p+1])<<8|uint32(xxh.buf[p]))*prime32_2
		xxh.v3 = (p32<<13 | p32>>19) * prime32_1
		p += 4
		p32 = xxh.v4 + (uint32(xxh.buf[p+3])<<24|uint32(xxh.buf[p+2])<<16|uint32(xxh.buf[p+1])<<8|uint32(xxh.buf[p]))*prime32_2
		xxh.v4 = (p32<<13 | p32>>19) * prime32_1

		p = r
		xxh.bufused = 0
	***REMOVED***

	for n := n - 16; p <= n; ***REMOVED***
		p32 := xxh.v1 + (uint32(input[p+3])<<24|uint32(input[p+2])<<16|uint32(input[p+1])<<8|uint32(input[p]))*prime32_2
		xxh.v1 = (p32<<13 | p32>>19) * prime32_1
		p += 4
		p32 = xxh.v2 + (uint32(input[p+3])<<24|uint32(input[p+2])<<16|uint32(input[p+1])<<8|uint32(input[p]))*prime32_2
		xxh.v2 = (p32<<13 | p32>>19) * prime32_1
		p += 4
		p32 = xxh.v3 + (uint32(input[p+3])<<24|uint32(input[p+2])<<16|uint32(input[p+1])<<8|uint32(input[p]))*prime32_2
		xxh.v3 = (p32<<13 | p32>>19) * prime32_1
		p += 4
		p32 = xxh.v4 + (uint32(input[p+3])<<24|uint32(input[p+2])<<16|uint32(input[p+1])<<8|uint32(input[p]))*prime32_2
		xxh.v4 = (p32<<13 | p32>>19) * prime32_1
		p += 4
	***REMOVED***

	copy(xxh.buf[xxh.bufused:], input[p:])
	xxh.bufused += len(input) - p

	return n, nil
***REMOVED***

// Sum32 returns the 32 bits Hash value.
func (xxh *xxHash) Sum32() uint32 ***REMOVED***
	h32 := uint32(xxh.totalLen)
	if xxh.totalLen >= 16 ***REMOVED***
		h32 += ((xxh.v1 << 1) | (xxh.v1 >> 31)) +
			((xxh.v2 << 7) | (xxh.v2 >> 25)) +
			((xxh.v3 << 12) | (xxh.v3 >> 20)) +
			((xxh.v4 << 18) | (xxh.v4 >> 14))
	***REMOVED*** else ***REMOVED***
		h32 += xxh.seed + prime32_5
	***REMOVED***

	p := 0
	n := xxh.bufused
	for n := n - 4; p <= n; p += 4 ***REMOVED***
		h32 += (uint32(xxh.buf[p+3])<<24 | uint32(xxh.buf[p+2])<<16 | uint32(xxh.buf[p+1])<<8 | uint32(xxh.buf[p])) * prime32_3
		h32 = ((h32 << 17) | (h32 >> 15)) * prime32_4
	***REMOVED***
	for ; p < n; p++ ***REMOVED***
		h32 += uint32(xxh.buf[p]) * prime32_5
		h32 = ((h32 << 11) | (h32 >> 21)) * prime32_1
	***REMOVED***

	h32 ^= h32 >> 15
	h32 *= prime32_2
	h32 ^= h32 >> 13
	h32 *= prime32_3
	h32 ^= h32 >> 16

	return h32
***REMOVED***

// Checksum returns the 32bits Hash value.
func Checksum(input []byte, seed uint32) uint32 ***REMOVED***
	n := len(input)
	h32 := uint32(n)

	if n < 16 ***REMOVED***
		h32 += seed + prime32_5
	***REMOVED*** else ***REMOVED***
		v1 := seed + prime32_1 + prime32_2
		v2 := seed + prime32_2
		v3 := seed
		v4 := seed - prime32_1
		p := 0
		for p <= n-16 ***REMOVED***
			v1 += (uint32(input[p+3])<<24 | uint32(input[p+2])<<16 | uint32(input[p+1])<<8 | uint32(input[p])) * prime32_2
			v1 = (v1<<13 | v1>>19) * prime32_1
			p += 4
			v2 += (uint32(input[p+3])<<24 | uint32(input[p+2])<<16 | uint32(input[p+1])<<8 | uint32(input[p])) * prime32_2
			v2 = (v2<<13 | v2>>19) * prime32_1
			p += 4
			v3 += (uint32(input[p+3])<<24 | uint32(input[p+2])<<16 | uint32(input[p+1])<<8 | uint32(input[p])) * prime32_2
			v3 = (v3<<13 | v3>>19) * prime32_1
			p += 4
			v4 += (uint32(input[p+3])<<24 | uint32(input[p+2])<<16 | uint32(input[p+1])<<8 | uint32(input[p])) * prime32_2
			v4 = (v4<<13 | v4>>19) * prime32_1
			p += 4
		***REMOVED***
		input = input[p:]
		n -= p
		h32 += ((v1 << 1) | (v1 >> 31)) +
			((v2 << 7) | (v2 >> 25)) +
			((v3 << 12) | (v3 >> 20)) +
			((v4 << 18) | (v4 >> 14))
	***REMOVED***

	p := 0
	for p <= n-4 ***REMOVED***
		h32 += (uint32(input[p+3])<<24 | uint32(input[p+2])<<16 | uint32(input[p+1])<<8 | uint32(input[p])) * prime32_3
		h32 = ((h32 << 17) | (h32 >> 15)) * prime32_4
		p += 4
	***REMOVED***
	for p < n ***REMOVED***
		h32 += uint32(input[p]) * prime32_5
		h32 = ((h32 << 11) | (h32 >> 21)) * prime32_1
		p++
	***REMOVED***

	h32 ^= h32 >> 15
	h32 *= prime32_2
	h32 ^= h32 >> 13
	h32 *= prime32_3
	h32 ^= h32 >> 16

	return h32
***REMOVED***
