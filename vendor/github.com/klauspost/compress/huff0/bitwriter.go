// Copyright 2018 Klaus Post. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Based on work Copyright (c) 2013, Yann Collet, released under BSD License.

package huff0

// bitWriter will write bits.
// First bit will be LSB of the first byte of output.
type bitWriter struct ***REMOVED***
	bitContainer uint64
	nBits        uint8
	out          []byte
***REMOVED***

// bitMask16 is bitmasks. Has extra to avoid bounds check.
var bitMask16 = [32]uint16***REMOVED***
	0, 1, 3, 7, 0xF, 0x1F,
	0x3F, 0x7F, 0xFF, 0x1FF, 0x3FF, 0x7FF,
	0xFFF, 0x1FFF, 0x3FFF, 0x7FFF, 0xFFFF, 0xFFFF,
	0xFFFF, 0xFFFF, 0xFFFF, 0xFFFF, 0xFFFF, 0xFFFF,
	0xFFFF, 0xFFFF***REMOVED*** /* up to 16 bits */

// addBits16Clean will add up to 16 bits. value may not contain more set bits than indicated.
// It will not check if there is space for them, so the caller must ensure that it has flushed recently.
func (b *bitWriter) addBits16Clean(value uint16, bits uint8) ***REMOVED***
	b.bitContainer |= uint64(value) << (b.nBits & 63)
	b.nBits += bits
***REMOVED***

// encSymbol will add up to 16 bits. value may not contain more set bits than indicated.
// It will not check if there is space for them, so the caller must ensure that it has flushed recently.
func (b *bitWriter) encSymbol(ct cTable, symbol byte) ***REMOVED***
	enc := ct[symbol]
	b.bitContainer |= uint64(enc.val) << (b.nBits & 63)
	if false ***REMOVED***
		if enc.nBits == 0 ***REMOVED***
			panic("nbits 0")
		***REMOVED***
	***REMOVED***
	b.nBits += enc.nBits
***REMOVED***

// encTwoSymbols will add up to 32 bits. value may not contain more set bits than indicated.
// It will not check if there is space for them, so the caller must ensure that it has flushed recently.
func (b *bitWriter) encTwoSymbols(ct cTable, av, bv byte) ***REMOVED***
	encA := ct[av]
	encB := ct[bv]
	sh := b.nBits & 63
	combined := uint64(encA.val) | (uint64(encB.val) << (encA.nBits & 63))
	b.bitContainer |= combined << sh
	if false ***REMOVED***
		if encA.nBits == 0 ***REMOVED***
			panic("nbitsA 0")
		***REMOVED***
		if encB.nBits == 0 ***REMOVED***
			panic("nbitsB 0")
		***REMOVED***
	***REMOVED***
	b.nBits += encA.nBits + encB.nBits
***REMOVED***

// flush32 will flush out, so there are at least 32 bits available for writing.
func (b *bitWriter) flush32() ***REMOVED***
	if b.nBits < 32 ***REMOVED***
		return
	***REMOVED***
	b.out = append(b.out,
		byte(b.bitContainer),
		byte(b.bitContainer>>8),
		byte(b.bitContainer>>16),
		byte(b.bitContainer>>24))
	b.nBits -= 32
	b.bitContainer >>= 32
***REMOVED***

// flushAlign will flush remaining full bytes and align to next byte boundary.
func (b *bitWriter) flushAlign() ***REMOVED***
	nbBytes := (b.nBits + 7) >> 3
	for i := uint8(0); i < nbBytes; i++ ***REMOVED***
		b.out = append(b.out, byte(b.bitContainer>>(i*8)))
	***REMOVED***
	b.nBits = 0
	b.bitContainer = 0
***REMOVED***

// close will write the alignment bit and write the final byte(s)
// to the output.
func (b *bitWriter) close() error ***REMOVED***
	// End mark
	b.addBits16Clean(1, 1)
	// flush until next byte.
	b.flushAlign()
	return nil
***REMOVED***
