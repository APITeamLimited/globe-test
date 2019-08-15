// Copyright 2018 Klaus Post. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Based on work Copyright (c) 2013, Yann Collet, released under BSD License.

package huff0

import (
	"errors"
	"io"
)

// bitReader reads a bitstream in reverse.
// The last set bit indicates the start of the stream and is used
// for aligning the input.
type bitReader struct ***REMOVED***
	in       []byte
	off      uint // next byte to read is at in[off - 1]
	value    uint64
	bitsRead uint8
***REMOVED***

// init initializes and resets the bit reader.
func (b *bitReader) init(in []byte) error ***REMOVED***
	if len(in) < 1 ***REMOVED***
		return errors.New("corrupt stream: too short")
	***REMOVED***
	b.in = in
	b.off = uint(len(in))
	// The highest bit of the last byte indicates where to start
	v := in[len(in)-1]
	if v == 0 ***REMOVED***
		return errors.New("corrupt stream, did not find end of stream")
	***REMOVED***
	b.bitsRead = 64
	b.value = 0
	b.fill()
	b.fill()
	b.bitsRead += 8 - uint8(highBit32(uint32(v)))
	return nil
***REMOVED***

// getBits will return n bits. n can be 0.
func (b *bitReader) getBits(n uint8) uint16 ***REMOVED***
	if n == 0 || b.bitsRead >= 64 ***REMOVED***
		return 0
	***REMOVED***
	return b.getBitsFast(n)
***REMOVED***

// getBitsFast requires that at least one bit is requested every time.
// There are no checks if the buffer is filled.
func (b *bitReader) getBitsFast(n uint8) uint16 ***REMOVED***
	const regMask = 64 - 1
	v := uint16((b.value << (b.bitsRead & regMask)) >> ((regMask + 1 - n) & regMask))
	b.bitsRead += n
	return v
***REMOVED***

// peekBitsFast requires that at least one bit is requested every time.
// There are no checks if the buffer is filled.
func (b *bitReader) peekBitsFast(n uint8) uint16 ***REMOVED***
	const regMask = 64 - 1
	v := uint16((b.value << (b.bitsRead & regMask)) >> ((regMask + 1 - n) & regMask))
	return v
***REMOVED***

// fillFast() will make sure at least 32 bits are available.
// There must be at least 4 bytes available.
func (b *bitReader) fillFast() ***REMOVED***
	if b.bitsRead < 32 ***REMOVED***
		return
	***REMOVED***
	// Do single re-slice to avoid bounds checks.
	v := b.in[b.off-4 : b.off]
	low := (uint32(v[0])) | (uint32(v[1]) << 8) | (uint32(v[2]) << 16) | (uint32(v[3]) << 24)
	b.value = (b.value << 32) | uint64(low)
	b.bitsRead -= 32
	b.off -= 4
***REMOVED***

// fill() will make sure at least 32 bits are available.
func (b *bitReader) fill() ***REMOVED***
	if b.bitsRead < 32 ***REMOVED***
		return
	***REMOVED***
	if b.off > 4 ***REMOVED***
		v := b.in[b.off-4 : b.off]
		low := (uint32(v[0])) | (uint32(v[1]) << 8) | (uint32(v[2]) << 16) | (uint32(v[3]) << 24)
		b.value = (b.value << 32) | uint64(low)
		b.bitsRead -= 32
		b.off -= 4
		return
	***REMOVED***
	for b.off > 0 ***REMOVED***
		b.value = (b.value << 8) | uint64(b.in[b.off-1])
		b.bitsRead -= 8
		b.off--
	***REMOVED***
***REMOVED***

// finished returns true if all bits have been read from the bit stream.
func (b *bitReader) finished() bool ***REMOVED***
	return b.off == 0 && b.bitsRead >= 64
***REMOVED***

// close the bitstream and returns an error if out-of-buffer reads occurred.
func (b *bitReader) close() error ***REMOVED***
	// Release reference.
	b.in = nil
	if b.bitsRead > 64 ***REMOVED***
		return io.ErrUnexpectedEOF
	***REMOVED***
	return nil
***REMOVED***
