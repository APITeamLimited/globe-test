// Package xxhash implements the 64-bit variant of xxHash (XXH64) as described
// at http://cyan4973.github.io/xxHash/.
// THIS IS VENDORED: Go to github.com/cespare/xxhash for original package.

package xxhash

import (
	"encoding/binary"
	"errors"
	"math/bits"
)

const (
	prime1 uint64 = 11400714785074694791
	prime2 uint64 = 14029467366897019727
	prime3 uint64 = 1609587929392839161
	prime4 uint64 = 9650029242287828579
	prime5 uint64 = 2870177450012600261
)

// NOTE(caleb): I'm using both consts and vars of the primes. Using consts where
// possible in the Go code is worth a small (but measurable) performance boost
// by avoiding some MOVQs. Vars are needed for the asm and also are useful for
// convenience in the Go code in a few places where we need to intentionally
// avoid constant arithmetic (e.g., v1 := prime1 + prime2 fails because the
// result overflows a uint64).
var (
	prime1v = prime1
	prime2v = prime2
	prime3v = prime3
	prime4v = prime4
	prime5v = prime5
)

// Digest implements hash.Hash64.
type Digest struct ***REMOVED***
	v1    uint64
	v2    uint64
	v3    uint64
	v4    uint64
	total uint64
	mem   [32]byte
	n     int // how much of mem is used
***REMOVED***

// New creates a new Digest that computes the 64-bit xxHash algorithm.
func New() *Digest ***REMOVED***
	var d Digest
	d.Reset()
	return &d
***REMOVED***

// Reset clears the Digest's state so that it can be reused.
func (d *Digest) Reset() ***REMOVED***
	d.v1 = prime1v + prime2
	d.v2 = prime2
	d.v3 = 0
	d.v4 = -prime1v
	d.total = 0
	d.n = 0
***REMOVED***

// Size always returns 8 bytes.
func (d *Digest) Size() int ***REMOVED*** return 8 ***REMOVED***

// BlockSize always returns 32 bytes.
func (d *Digest) BlockSize() int ***REMOVED*** return 32 ***REMOVED***

// Write adds more data to d. It always returns len(b), nil.
func (d *Digest) Write(b []byte) (n int, err error) ***REMOVED***
	n = len(b)
	d.total += uint64(n)

	if d.n+n < 32 ***REMOVED***
		// This new data doesn't even fill the current block.
		copy(d.mem[d.n:], b)
		d.n += n
		return
	***REMOVED***

	if d.n > 0 ***REMOVED***
		// Finish off the partial block.
		copy(d.mem[d.n:], b)
		d.v1 = round(d.v1, u64(d.mem[0:8]))
		d.v2 = round(d.v2, u64(d.mem[8:16]))
		d.v3 = round(d.v3, u64(d.mem[16:24]))
		d.v4 = round(d.v4, u64(d.mem[24:32]))
		b = b[32-d.n:]
		d.n = 0
	***REMOVED***

	if len(b) >= 32 ***REMOVED***
		// One or more full blocks left.
		nw := writeBlocks(d, b)
		b = b[nw:]
	***REMOVED***

	// Store any remaining partial block.
	copy(d.mem[:], b)
	d.n = len(b)

	return
***REMOVED***

// Sum appends the current hash to b and returns the resulting slice.
func (d *Digest) Sum(b []byte) []byte ***REMOVED***
	s := d.Sum64()
	return append(
		b,
		byte(s>>56),
		byte(s>>48),
		byte(s>>40),
		byte(s>>32),
		byte(s>>24),
		byte(s>>16),
		byte(s>>8),
		byte(s),
	)
***REMOVED***

// Sum64 returns the current hash.
func (d *Digest) Sum64() uint64 ***REMOVED***
	var h uint64

	if d.total >= 32 ***REMOVED***
		v1, v2, v3, v4 := d.v1, d.v2, d.v3, d.v4
		h = rol1(v1) + rol7(v2) + rol12(v3) + rol18(v4)
		h = mergeRound(h, v1)
		h = mergeRound(h, v2)
		h = mergeRound(h, v3)
		h = mergeRound(h, v4)
	***REMOVED*** else ***REMOVED***
		h = d.v3 + prime5
	***REMOVED***

	h += d.total

	i, end := 0, d.n
	for ; i+8 <= end; i += 8 ***REMOVED***
		k1 := round(0, u64(d.mem[i:i+8]))
		h ^= k1
		h = rol27(h)*prime1 + prime4
	***REMOVED***
	if i+4 <= end ***REMOVED***
		h ^= uint64(u32(d.mem[i:i+4])) * prime1
		h = rol23(h)*prime2 + prime3
		i += 4
	***REMOVED***
	for i < end ***REMOVED***
		h ^= uint64(d.mem[i]) * prime5
		h = rol11(h) * prime1
		i++
	***REMOVED***

	h ^= h >> 33
	h *= prime2
	h ^= h >> 29
	h *= prime3
	h ^= h >> 32

	return h
***REMOVED***

const (
	magic         = "xxh\x06"
	marshaledSize = len(magic) + 8*5 + 32
)

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (d *Digest) MarshalBinary() ([]byte, error) ***REMOVED***
	b := make([]byte, 0, marshaledSize)
	b = append(b, magic...)
	b = appendUint64(b, d.v1)
	b = appendUint64(b, d.v2)
	b = appendUint64(b, d.v3)
	b = appendUint64(b, d.v4)
	b = appendUint64(b, d.total)
	b = append(b, d.mem[:d.n]...)
	b = b[:len(b)+len(d.mem)-d.n]
	return b, nil
***REMOVED***

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (d *Digest) UnmarshalBinary(b []byte) error ***REMOVED***
	if len(b) < len(magic) || string(b[:len(magic)]) != magic ***REMOVED***
		return errors.New("xxhash: invalid hash state identifier")
	***REMOVED***
	if len(b) != marshaledSize ***REMOVED***
		return errors.New("xxhash: invalid hash state size")
	***REMOVED***
	b = b[len(magic):]
	b, d.v1 = consumeUint64(b)
	b, d.v2 = consumeUint64(b)
	b, d.v3 = consumeUint64(b)
	b, d.v4 = consumeUint64(b)
	b, d.total = consumeUint64(b)
	copy(d.mem[:], b)
	b = b[len(d.mem):]
	d.n = int(d.total % uint64(len(d.mem)))
	return nil
***REMOVED***

func appendUint64(b []byte, x uint64) []byte ***REMOVED***
	var a [8]byte
	binary.LittleEndian.PutUint64(a[:], x)
	return append(b, a[:]...)
***REMOVED***

func consumeUint64(b []byte) ([]byte, uint64) ***REMOVED***
	x := u64(b)
	return b[8:], x
***REMOVED***

func u64(b []byte) uint64 ***REMOVED*** return binary.LittleEndian.Uint64(b) ***REMOVED***
func u32(b []byte) uint32 ***REMOVED*** return binary.LittleEndian.Uint32(b) ***REMOVED***

func round(acc, input uint64) uint64 ***REMOVED***
	acc += input * prime2
	acc = rol31(acc)
	acc *= prime1
	return acc
***REMOVED***

func mergeRound(acc, val uint64) uint64 ***REMOVED***
	val = round(0, val)
	acc ^= val
	acc = acc*prime1 + prime4
	return acc
***REMOVED***

func rol1(x uint64) uint64  ***REMOVED*** return bits.RotateLeft64(x, 1) ***REMOVED***
func rol7(x uint64) uint64  ***REMOVED*** return bits.RotateLeft64(x, 7) ***REMOVED***
func rol11(x uint64) uint64 ***REMOVED*** return bits.RotateLeft64(x, 11) ***REMOVED***
func rol12(x uint64) uint64 ***REMOVED*** return bits.RotateLeft64(x, 12) ***REMOVED***
func rol18(x uint64) uint64 ***REMOVED*** return bits.RotateLeft64(x, 18) ***REMOVED***
func rol23(x uint64) uint64 ***REMOVED*** return bits.RotateLeft64(x, 23) ***REMOVED***
func rol27(x uint64) uint64 ***REMOVED*** return bits.RotateLeft64(x, 27) ***REMOVED***
func rol31(x uint64) uint64 ***REMOVED*** return bits.RotateLeft64(x, 31) ***REMOVED***
