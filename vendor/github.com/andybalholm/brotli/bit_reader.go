package brotli

import "encoding/binary"

/* Copyright 2013 Google Inc. All Rights Reserved.

   Distributed under MIT license.
   See file LICENSE for detail or copy at https://opensource.org/licenses/MIT
*/

/* Bit reading helpers */

const shortFillBitWindowRead = (8 >> 1)

var kBitMask = [33]uint32***REMOVED***
	0x00000000,
	0x00000001,
	0x00000003,
	0x00000007,
	0x0000000F,
	0x0000001F,
	0x0000003F,
	0x0000007F,
	0x000000FF,
	0x000001FF,
	0x000003FF,
	0x000007FF,
	0x00000FFF,
	0x00001FFF,
	0x00003FFF,
	0x00007FFF,
	0x0000FFFF,
	0x0001FFFF,
	0x0003FFFF,
	0x0007FFFF,
	0x000FFFFF,
	0x001FFFFF,
	0x003FFFFF,
	0x007FFFFF,
	0x00FFFFFF,
	0x01FFFFFF,
	0x03FFFFFF,
	0x07FFFFFF,
	0x0FFFFFFF,
	0x1FFFFFFF,
	0x3FFFFFFF,
	0x7FFFFFFF,
	0xFFFFFFFF,
***REMOVED***

func bitMask(n uint32) uint32 ***REMOVED***
	return kBitMask[n]
***REMOVED***

type bitReader struct ***REMOVED***
	val_      uint64
	bit_pos_  uint32
	input     []byte
	input_len uint
	byte_pos  uint
***REMOVED***

type bitReaderState struct ***REMOVED***
	val_      uint64
	bit_pos_  uint32
	input     []byte
	input_len uint
	byte_pos  uint
***REMOVED***

/* Initializes the BrotliBitReader fields. */

/* Ensures that accumulator is not empty.
   May consume up to sizeof(brotli_reg_t) - 1 bytes of input.
   Returns false if data is required but there is no input available.
   For BROTLI_ALIGNED_READ this function also prepares bit reader for aligned
   reading. */
func bitReaderSaveState(from *bitReader, to *bitReaderState) ***REMOVED***
	to.val_ = from.val_
	to.bit_pos_ = from.bit_pos_
	to.input = from.input
	to.input_len = from.input_len
	to.byte_pos = from.byte_pos
***REMOVED***

func bitReaderRestoreState(to *bitReader, from *bitReaderState) ***REMOVED***
	to.val_ = from.val_
	to.bit_pos_ = from.bit_pos_
	to.input = from.input
	to.input_len = from.input_len
	to.byte_pos = from.byte_pos
***REMOVED***

func getAvailableBits(br *bitReader) uint32 ***REMOVED***
	return 64 - br.bit_pos_
***REMOVED***

/* Returns amount of unread bytes the bit reader still has buffered from the
   BrotliInput, including whole bytes in br->val_. */
func getRemainingBytes(br *bitReader) uint ***REMOVED***
	return uint(uint32(br.input_len-br.byte_pos) + (getAvailableBits(br) >> 3))
***REMOVED***

/* Checks if there is at least |num| bytes left in the input ring-buffer
   (excluding the bits remaining in br->val_). */
func checkInputAmount(br *bitReader, num uint) bool ***REMOVED***
	return br.input_len-br.byte_pos >= num
***REMOVED***

/* Guarantees that there are at least |n_bits| + 1 bits in accumulator.
   Precondition: accumulator contains at least 1 bit.
   |n_bits| should be in the range [1..24] for regular build. For portable
   non-64-bit little-endian build only 16 bits are safe to request. */
func fillBitWindow(br *bitReader, n_bits uint32) ***REMOVED***
	if br.bit_pos_ >= 32 ***REMOVED***
		br.val_ >>= 32
		br.bit_pos_ ^= 32 /* here same as -= 32 because of the if condition */
		br.val_ |= (uint64(binary.LittleEndian.Uint32(br.input[br.byte_pos:]))) << 32
		br.byte_pos += 4
	***REMOVED***
***REMOVED***

/* Mostly like BrotliFillBitWindow, but guarantees only 16 bits and reads no
   more than BROTLI_SHORT_FILL_BIT_WINDOW_READ bytes of input. */
func fillBitWindow16(br *bitReader) ***REMOVED***
	fillBitWindow(br, 17)
***REMOVED***

/* Tries to pull one byte of input to accumulator.
   Returns false if there is no input available. */
func pullByte(br *bitReader) bool ***REMOVED***
	if br.byte_pos == br.input_len ***REMOVED***
		return false
	***REMOVED***

	br.val_ >>= 8
	br.val_ |= (uint64(br.input[br.byte_pos])) << 56
	br.bit_pos_ -= 8
	br.byte_pos++
	return true
***REMOVED***

/* Returns currently available bits.
   The number of valid bits could be calculated by BrotliGetAvailableBits. */
func getBitsUnmasked(br *bitReader) uint64 ***REMOVED***
	return br.val_ >> br.bit_pos_
***REMOVED***

/* Like BrotliGetBits, but does not mask the result.
   The result contains at least 16 valid bits. */
func get16BitsUnmasked(br *bitReader) uint32 ***REMOVED***
	fillBitWindow(br, 16)
	return uint32(getBitsUnmasked(br))
***REMOVED***

/* Returns the specified number of bits from |br| without advancing bit
   position. */
func getBits(br *bitReader, n_bits uint32) uint32 ***REMOVED***
	fillBitWindow(br, n_bits)
	return uint32(getBitsUnmasked(br)) & bitMask(n_bits)
***REMOVED***

/* Tries to peek the specified amount of bits. Returns false, if there
   is not enough input. */
func safeGetBits(br *bitReader, n_bits uint32, val *uint32) bool ***REMOVED***
	for getAvailableBits(br) < n_bits ***REMOVED***
		if !pullByte(br) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	*val = uint32(getBitsUnmasked(br)) & bitMask(n_bits)
	return true
***REMOVED***

/* Advances the bit pos by |n_bits|. */
func dropBits(br *bitReader, n_bits uint32) ***REMOVED***
	br.bit_pos_ += n_bits
***REMOVED***

func bitReaderUnload(br *bitReader) ***REMOVED***
	var unused_bytes uint32 = getAvailableBits(br) >> 3
	var unused_bits uint32 = unused_bytes << 3
	br.byte_pos -= uint(unused_bytes)
	if unused_bits == 64 ***REMOVED***
		br.val_ = 0
	***REMOVED*** else ***REMOVED***
		br.val_ <<= unused_bits
	***REMOVED***

	br.bit_pos_ += unused_bits
***REMOVED***

/* Reads the specified number of bits from |br| and advances the bit pos.
   Precondition: accumulator MUST contain at least |n_bits|. */
func takeBits(br *bitReader, n_bits uint32, val *uint32) ***REMOVED***
	*val = uint32(getBitsUnmasked(br)) & bitMask(n_bits)
	dropBits(br, n_bits)
***REMOVED***

/* Reads the specified number of bits from |br| and advances the bit pos.
   Assumes that there is enough input to perform BrotliFillBitWindow. */
func readBits(br *bitReader, n_bits uint32) uint32 ***REMOVED***
	var val uint32
	fillBitWindow(br, n_bits)
	takeBits(br, n_bits, &val)
	return val
***REMOVED***

/* Tries to read the specified amount of bits. Returns false, if there
   is not enough input. |n_bits| MUST be positive. */
func safeReadBits(br *bitReader, n_bits uint32, val *uint32) bool ***REMOVED***
	for getAvailableBits(br) < n_bits ***REMOVED***
		if !pullByte(br) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	takeBits(br, n_bits, val)
	return true
***REMOVED***

/* Advances the bit reader position to the next byte boundary and verifies
   that any skipped bits are set to zero. */
func bitReaderJumpToByteBoundary(br *bitReader) bool ***REMOVED***
	var pad_bits_count uint32 = getAvailableBits(br) & 0x7
	var pad_bits uint32 = 0
	if pad_bits_count != 0 ***REMOVED***
		takeBits(br, pad_bits_count, &pad_bits)
	***REMOVED***

	return pad_bits == 0
***REMOVED***

/* Copies remaining input bytes stored in the bit reader to the output. Value
   |num| may not be larger than BrotliGetRemainingBytes. The bit reader must be
   warmed up again after this. */
func copyBytes(dest []byte, br *bitReader, num uint) ***REMOVED***
	for getAvailableBits(br) >= 8 && num > 0 ***REMOVED***
		dest[0] = byte(getBitsUnmasked(br))
		dropBits(br, 8)
		dest = dest[1:]
		num--
	***REMOVED***

	copy(dest, br.input[br.byte_pos:][:num])
	br.byte_pos += num
***REMOVED***

func initBitReader(br *bitReader) ***REMOVED***
	br.val_ = 0
	br.bit_pos_ = 64
***REMOVED***

func warmupBitReader(br *bitReader) bool ***REMOVED***
	/* Fixing alignment after unaligned BrotliFillWindow would result accumulator
	   overflow. If unalignment is caused by BrotliSafeReadBits, then there is
	   enough space in accumulator to fix alignment. */
	if getAvailableBits(br) == 0 ***REMOVED***
		if !pullByte(br) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	return true
***REMOVED***
