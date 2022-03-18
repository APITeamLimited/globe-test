// Copyright 2018 Klaus Post. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Based on work Copyright (c) 2013, Yann Collet, released under BSD License.

package huff0

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// bitReader reads a bitstream in reverse.
// The last set bit indicates the start of the stream and is used
// for aligning the input.
type bitReaderBytes struct ***REMOVED***
	in       []byte
	off      uint // next byte to read is at in[off - 1]
	value    uint64
	bitsRead uint8
***REMOVED***

// init initializes and resets the bit reader.
func (b *bitReaderBytes) init(in []byte) error ***REMOVED***
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
	if len(in) >= 8 ***REMOVED***
		b.fillFastStart()
	***REMOVED*** else ***REMOVED***
		b.fill()
		b.fill()
	***REMOVED***
	b.advance(8 - uint8(highBit32(uint32(v))))
	return nil
***REMOVED***

// peekBitsFast requires that at least one bit is requested every time.
// There are no checks if the buffer is filled.
func (b *bitReaderBytes) peekByteFast() uint8 ***REMOVED***
	got := uint8(b.value >> 56)
	return got
***REMOVED***

func (b *bitReaderBytes) advance(n uint8) ***REMOVED***
	b.bitsRead += n
	b.value <<= n & 63
***REMOVED***

// fillFast() will make sure at least 32 bits are available.
// There must be at least 4 bytes available.
func (b *bitReaderBytes) fillFast() ***REMOVED***
	if b.bitsRead < 32 ***REMOVED***
		return
	***REMOVED***

	// 2 bounds checks.
	v := b.in[b.off-4 : b.off]
	v = v[:4]
	low := (uint32(v[0])) | (uint32(v[1]) << 8) | (uint32(v[2]) << 16) | (uint32(v[3]) << 24)
	b.value |= uint64(low) << (b.bitsRead - 32)
	b.bitsRead -= 32
	b.off -= 4
***REMOVED***

// fillFastStart() assumes the bitReaderBytes is empty and there is at least 8 bytes to read.
func (b *bitReaderBytes) fillFastStart() ***REMOVED***
	// Do single re-slice to avoid bounds checks.
	b.value = binary.LittleEndian.Uint64(b.in[b.off-8:])
	b.bitsRead = 0
	b.off -= 8
***REMOVED***

// fill() will make sure at least 32 bits are available.
func (b *bitReaderBytes) fill() ***REMOVED***
	if b.bitsRead < 32 ***REMOVED***
		return
	***REMOVED***
	if b.off > 4 ***REMOVED***
		v := b.in[b.off-4:]
		v = v[:4]
		low := (uint32(v[0])) | (uint32(v[1]) << 8) | (uint32(v[2]) << 16) | (uint32(v[3]) << 24)
		b.value |= uint64(low) << (b.bitsRead - 32)
		b.bitsRead -= 32
		b.off -= 4
		return
	***REMOVED***
	for b.off > 0 ***REMOVED***
		b.value |= uint64(b.in[b.off-1]) << (b.bitsRead - 8)
		b.bitsRead -= 8
		b.off--
	***REMOVED***
***REMOVED***

// finished returns true if all bits have been read from the bit stream.
func (b *bitReaderBytes) finished() bool ***REMOVED***
	return b.off == 0 && b.bitsRead >= 64
***REMOVED***

func (b *bitReaderBytes) remaining() uint ***REMOVED***
	return b.off*8 + uint(64-b.bitsRead)
***REMOVED***

// close the bitstream and returns an error if out-of-buffer reads occurred.
func (b *bitReaderBytes) close() error ***REMOVED***
	// Release reference.
	b.in = nil
	if b.remaining() > 0 ***REMOVED***
		return fmt.Errorf("corrupt input: %d bits remain on stream", b.remaining())
	***REMOVED***
	if b.bitsRead > 64 ***REMOVED***
		return io.ErrUnexpectedEOF
	***REMOVED***
	return nil
***REMOVED***

// bitReaderShifted reads a bitstream in reverse.
// The last set bit indicates the start of the stream and is used
// for aligning the input.
type bitReaderShifted struct ***REMOVED***
	in       []byte
	off      uint // next byte to read is at in[off - 1]
	value    uint64
	bitsRead uint8
***REMOVED***

// init initializes and resets the bit reader.
func (b *bitReaderShifted) init(in []byte) error ***REMOVED***
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
	if len(in) >= 8 ***REMOVED***
		b.fillFastStart()
	***REMOVED*** else ***REMOVED***
		b.fill()
		b.fill()
	***REMOVED***
	b.advance(8 - uint8(highBit32(uint32(v))))
	return nil
***REMOVED***

// peekBitsFast requires that at least one bit is requested every time.
// There are no checks if the buffer is filled.
func (b *bitReaderShifted) peekBitsFast(n uint8) uint16 ***REMOVED***
	return uint16(b.value >> ((64 - n) & 63))
***REMOVED***

// peekTopBits(n) is equvialent to peekBitFast(64 - n)
func (b *bitReaderShifted) peekTopBits(n uint8) uint16 ***REMOVED***
	return uint16(b.value >> n)
***REMOVED***

func (b *bitReaderShifted) advance(n uint8) ***REMOVED***
	b.bitsRead += n
	b.value <<= n & 63
***REMOVED***

// fillFast() will make sure at least 32 bits are available.
// There must be at least 4 bytes available.
func (b *bitReaderShifted) fillFast() ***REMOVED***
	if b.bitsRead < 32 ***REMOVED***
		return
	***REMOVED***

	// 2 bounds checks.
	v := b.in[b.off-4 : b.off]
	v = v[:4]
	low := (uint32(v[0])) | (uint32(v[1]) << 8) | (uint32(v[2]) << 16) | (uint32(v[3]) << 24)
	b.value |= uint64(low) << ((b.bitsRead - 32) & 63)
	b.bitsRead -= 32
	b.off -= 4
***REMOVED***

// fillFastStart() assumes the bitReaderShifted is empty and there is at least 8 bytes to read.
func (b *bitReaderShifted) fillFastStart() ***REMOVED***
	// Do single re-slice to avoid bounds checks.
	b.value = binary.LittleEndian.Uint64(b.in[b.off-8:])
	b.bitsRead = 0
	b.off -= 8
***REMOVED***

// fill() will make sure at least 32 bits are available.
func (b *bitReaderShifted) fill() ***REMOVED***
	if b.bitsRead < 32 ***REMOVED***
		return
	***REMOVED***
	if b.off > 4 ***REMOVED***
		v := b.in[b.off-4:]
		v = v[:4]
		low := (uint32(v[0])) | (uint32(v[1]) << 8) | (uint32(v[2]) << 16) | (uint32(v[3]) << 24)
		b.value |= uint64(low) << ((b.bitsRead - 32) & 63)
		b.bitsRead -= 32
		b.off -= 4
		return
	***REMOVED***
	for b.off > 0 ***REMOVED***
		b.value |= uint64(b.in[b.off-1]) << ((b.bitsRead - 8) & 63)
		b.bitsRead -= 8
		b.off--
	***REMOVED***
***REMOVED***

// finished returns true if all bits have been read from the bit stream.
func (b *bitReaderShifted) finished() bool ***REMOVED***
	return b.off == 0 && b.bitsRead >= 64
***REMOVED***

func (b *bitReaderShifted) remaining() uint ***REMOVED***
	return b.off*8 + uint(64-b.bitsRead)
***REMOVED***

// close the bitstream and returns an error if out-of-buffer reads occurred.
func (b *bitReaderShifted) close() error ***REMOVED***
	// Release reference.
	b.in = nil
	if b.remaining() > 0 ***REMOVED***
		return fmt.Errorf("corrupt input: %d bits remain on stream", b.remaining())
	***REMOVED***
	if b.bitsRead > 64 ***REMOVED***
		return io.ErrUnexpectedEOF
	***REMOVED***
	return nil
***REMOVED***
