// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flate

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

const (
	// bits 0-16  	xoffset = offset - MIN_OFFSET_SIZE, or literal - 16 bits
	// bits 16-22	offsetcode - 5 bits
	// bits 22-30   xlength = length - MIN_MATCH_LENGTH - 8 bits
	// bits 30-32   type   0 = literal  1=EOF  2=Match   3=Unused - 2 bits
	lengthShift         = 22
	offsetMask          = 1<<lengthShift - 1
	typeMask            = 3 << 30
	literalType         = 0 << 30
	matchType           = 1 << 30
	matchOffsetOnlyMask = 0xffff
)

// The length code for length X (MIN_MATCH_LENGTH <= X <= MAX_MATCH_LENGTH)
// is lengthCodes[length - MIN_MATCH_LENGTH]
var lengthCodes = [256]uint8***REMOVED***
	0, 1, 2, 3, 4, 5, 6, 7, 8, 8,
	9, 9, 10, 10, 11, 11, 12, 12, 12, 12,
	13, 13, 13, 13, 14, 14, 14, 14, 15, 15,
	15, 15, 16, 16, 16, 16, 16, 16, 16, 16,
	17, 17, 17, 17, 17, 17, 17, 17, 18, 18,
	18, 18, 18, 18, 18, 18, 19, 19, 19, 19,
	19, 19, 19, 19, 20, 20, 20, 20, 20, 20,
	20, 20, 20, 20, 20, 20, 20, 20, 20, 20,
	21, 21, 21, 21, 21, 21, 21, 21, 21, 21,
	21, 21, 21, 21, 21, 21, 22, 22, 22, 22,
	22, 22, 22, 22, 22, 22, 22, 22, 22, 22,
	22, 22, 23, 23, 23, 23, 23, 23, 23, 23,
	23, 23, 23, 23, 23, 23, 23, 23, 24, 24,
	24, 24, 24, 24, 24, 24, 24, 24, 24, 24,
	24, 24, 24, 24, 24, 24, 24, 24, 24, 24,
	24, 24, 24, 24, 24, 24, 24, 24, 24, 24,
	25, 25, 25, 25, 25, 25, 25, 25, 25, 25,
	25, 25, 25, 25, 25, 25, 25, 25, 25, 25,
	25, 25, 25, 25, 25, 25, 25, 25, 25, 25,
	25, 25, 26, 26, 26, 26, 26, 26, 26, 26,
	26, 26, 26, 26, 26, 26, 26, 26, 26, 26,
	26, 26, 26, 26, 26, 26, 26, 26, 26, 26,
	26, 26, 26, 26, 27, 27, 27, 27, 27, 27,
	27, 27, 27, 27, 27, 27, 27, 27, 27, 27,
	27, 27, 27, 27, 27, 27, 27, 27, 27, 27,
	27, 27, 27, 27, 27, 28,
***REMOVED***

// lengthCodes1 is length codes, but starting at 1.
var lengthCodes1 = [256]uint8***REMOVED***
	1, 2, 3, 4, 5, 6, 7, 8, 9, 9,
	10, 10, 11, 11, 12, 12, 13, 13, 13, 13,
	14, 14, 14, 14, 15, 15, 15, 15, 16, 16,
	16, 16, 17, 17, 17, 17, 17, 17, 17, 17,
	18, 18, 18, 18, 18, 18, 18, 18, 19, 19,
	19, 19, 19, 19, 19, 19, 20, 20, 20, 20,
	20, 20, 20, 20, 21, 21, 21, 21, 21, 21,
	21, 21, 21, 21, 21, 21, 21, 21, 21, 21,
	22, 22, 22, 22, 22, 22, 22, 22, 22, 22,
	22, 22, 22, 22, 22, 22, 23, 23, 23, 23,
	23, 23, 23, 23, 23, 23, 23, 23, 23, 23,
	23, 23, 24, 24, 24, 24, 24, 24, 24, 24,
	24, 24, 24, 24, 24, 24, 24, 24, 25, 25,
	25, 25, 25, 25, 25, 25, 25, 25, 25, 25,
	25, 25, 25, 25, 25, 25, 25, 25, 25, 25,
	25, 25, 25, 25, 25, 25, 25, 25, 25, 25,
	26, 26, 26, 26, 26, 26, 26, 26, 26, 26,
	26, 26, 26, 26, 26, 26, 26, 26, 26, 26,
	26, 26, 26, 26, 26, 26, 26, 26, 26, 26,
	26, 26, 27, 27, 27, 27, 27, 27, 27, 27,
	27, 27, 27, 27, 27, 27, 27, 27, 27, 27,
	27, 27, 27, 27, 27, 27, 27, 27, 27, 27,
	27, 27, 27, 27, 28, 28, 28, 28, 28, 28,
	28, 28, 28, 28, 28, 28, 28, 28, 28, 28,
	28, 28, 28, 28, 28, 28, 28, 28, 28, 28,
	28, 28, 28, 28, 28, 29,
***REMOVED***

var offsetCodes = [256]uint32***REMOVED***
	0, 1, 2, 3, 4, 4, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7,
	8, 8, 8, 8, 8, 8, 8, 8, 9, 9, 9, 9, 9, 9, 9, 9,
	10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10,
	11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11,
	12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12,
	12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12,
	13, 13, 13, 13, 13, 13, 13, 13, 13, 13, 13, 13, 13, 13, 13, 13,
	13, 13, 13, 13, 13, 13, 13, 13, 13, 13, 13, 13, 13, 13, 13, 13,
	14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14,
	14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14,
	14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14,
	14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
***REMOVED***

// offsetCodes14 are offsetCodes, but with 14 added.
var offsetCodes14 = [256]uint32***REMOVED***
	14, 15, 16, 17, 18, 18, 19, 19, 20, 20, 20, 20, 21, 21, 21, 21,
	22, 22, 22, 22, 22, 22, 22, 22, 23, 23, 23, 23, 23, 23, 23, 23,
	24, 24, 24, 24, 24, 24, 24, 24, 24, 24, 24, 24, 24, 24, 24, 24,
	25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25,
	26, 26, 26, 26, 26, 26, 26, 26, 26, 26, 26, 26, 26, 26, 26, 26,
	26, 26, 26, 26, 26, 26, 26, 26, 26, 26, 26, 26, 26, 26, 26, 26,
	27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27,
	27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27,
	28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28,
	28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28,
	28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28,
	28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28,
	29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29,
	29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29,
	29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29,
	29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29,
***REMOVED***

type token uint32

type tokens struct ***REMOVED***
	extraHist [32]uint16  // codes 256->maxnumlit
	offHist   [32]uint16  // offset codes
	litHist   [256]uint16 // codes 0->255
	nFilled   int
	n         uint16 // Must be able to contain maxStoreBlockSize
	tokens    [maxStoreBlockSize + 1]token
***REMOVED***

func (t *tokens) Reset() ***REMOVED***
	if t.n == 0 ***REMOVED***
		return
	***REMOVED***
	t.n = 0
	t.nFilled = 0
	for i := range t.litHist[:] ***REMOVED***
		t.litHist[i] = 0
	***REMOVED***
	for i := range t.extraHist[:] ***REMOVED***
		t.extraHist[i] = 0
	***REMOVED***
	for i := range t.offHist[:] ***REMOVED***
		t.offHist[i] = 0
	***REMOVED***
***REMOVED***

func (t *tokens) Fill() ***REMOVED***
	if t.n == 0 ***REMOVED***
		return
	***REMOVED***
	for i, v := range t.litHist[:] ***REMOVED***
		if v == 0 ***REMOVED***
			t.litHist[i] = 1
			t.nFilled++
		***REMOVED***
	***REMOVED***
	for i, v := range t.extraHist[:literalCount-256] ***REMOVED***
		if v == 0 ***REMOVED***
			t.nFilled++
			t.extraHist[i] = 1
		***REMOVED***
	***REMOVED***
	for i, v := range t.offHist[:offsetCodeCount] ***REMOVED***
		if v == 0 ***REMOVED***
			t.offHist[i] = 1
		***REMOVED***
	***REMOVED***
***REMOVED***

func indexTokens(in []token) tokens ***REMOVED***
	var t tokens
	t.indexTokens(in)
	return t
***REMOVED***

func (t *tokens) indexTokens(in []token) ***REMOVED***
	t.Reset()
	for _, tok := range in ***REMOVED***
		if tok < matchType ***REMOVED***
			t.AddLiteral(tok.literal())
			continue
		***REMOVED***
		t.AddMatch(uint32(tok.length()), tok.offset()&matchOffsetOnlyMask)
	***REMOVED***
***REMOVED***

// emitLiteral writes a literal chunk and returns the number of bytes written.
func emitLiteral(dst *tokens, lit []byte) ***REMOVED***
	for _, v := range lit ***REMOVED***
		dst.tokens[dst.n] = token(v)
		dst.litHist[v]++
		dst.n++
	***REMOVED***
***REMOVED***

func (t *tokens) AddLiteral(lit byte) ***REMOVED***
	t.tokens[t.n] = token(lit)
	t.litHist[lit]++
	t.n++
***REMOVED***

// from https://stackoverflow.com/a/28730362
func mFastLog2(val float32) float32 ***REMOVED***
	ux := int32(math.Float32bits(val))
	log2 := (float32)(((ux >> 23) & 255) - 128)
	ux &= -0x7f800001
	ux += 127 << 23
	uval := math.Float32frombits(uint32(ux))
	log2 += ((-0.34484843)*uval+2.02466578)*uval - 0.67487759
	return log2
***REMOVED***

// EstimatedBits will return an minimum size estimated by an *optimal*
// compression of the block.
// The size of the block
func (t *tokens) EstimatedBits() int ***REMOVED***
	shannon := float32(0)
	bits := int(0)
	nMatches := 0
	total := int(t.n) + t.nFilled
	if total > 0 ***REMOVED***
		invTotal := 1.0 / float32(total)
		for _, v := range t.litHist[:] ***REMOVED***
			if v > 0 ***REMOVED***
				n := float32(v)
				shannon += atLeastOne(-mFastLog2(n*invTotal)) * n
			***REMOVED***
		***REMOVED***
		// Just add 15 for EOB
		shannon += 15
		for i, v := range t.extraHist[1 : literalCount-256] ***REMOVED***
			if v > 0 ***REMOVED***
				n := float32(v)
				shannon += atLeastOne(-mFastLog2(n*invTotal)) * n
				bits += int(lengthExtraBits[i&31]) * int(v)
				nMatches += int(v)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if nMatches > 0 ***REMOVED***
		invTotal := 1.0 / float32(nMatches)
		for i, v := range t.offHist[:offsetCodeCount] ***REMOVED***
			if v > 0 ***REMOVED***
				n := float32(v)
				shannon += atLeastOne(-mFastLog2(n*invTotal)) * n
				bits += int(offsetExtraBits[i&31]) * int(v)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return int(shannon) + bits
***REMOVED***

// AddMatch adds a match to the tokens.
// This function is very sensitive to inlining and right on the border.
func (t *tokens) AddMatch(xlength uint32, xoffset uint32) ***REMOVED***
	if debugDeflate ***REMOVED***
		if xlength >= maxMatchLength+baseMatchLength ***REMOVED***
			panic(fmt.Errorf("invalid length: %v", xlength))
		***REMOVED***
		if xoffset >= maxMatchOffset+baseMatchOffset ***REMOVED***
			panic(fmt.Errorf("invalid offset: %v", xoffset))
		***REMOVED***
	***REMOVED***
	oCode := offsetCode(xoffset)
	xoffset |= oCode << 16

	t.extraHist[lengthCodes1[uint8(xlength)]]++
	t.offHist[oCode&31]++
	t.tokens[t.n] = token(matchType | xlength<<lengthShift | xoffset)
	t.n++
***REMOVED***

// AddMatchLong adds a match to the tokens, potentially longer than max match length.
// Length should NOT have the base subtracted, only offset should.
func (t *tokens) AddMatchLong(xlength int32, xoffset uint32) ***REMOVED***
	if debugDeflate ***REMOVED***
		if xoffset >= maxMatchOffset+baseMatchOffset ***REMOVED***
			panic(fmt.Errorf("invalid offset: %v", xoffset))
		***REMOVED***
	***REMOVED***
	oc := offsetCode(xoffset)
	xoffset |= oc << 16
	for xlength > 0 ***REMOVED***
		xl := xlength
		if xl > 258 ***REMOVED***
			// We need to have at least baseMatchLength left over for next loop.
			if xl > 258+baseMatchLength ***REMOVED***
				xl = 258
			***REMOVED*** else ***REMOVED***
				xl = 258 - baseMatchLength
			***REMOVED***
		***REMOVED***
		xlength -= xl
		xl -= baseMatchLength
		t.extraHist[lengthCodes1[uint8(xl)]]++
		t.offHist[oc&31]++
		t.tokens[t.n] = token(matchType | uint32(xl)<<lengthShift | xoffset)
		t.n++
	***REMOVED***
***REMOVED***

func (t *tokens) AddEOB() ***REMOVED***
	t.tokens[t.n] = token(endBlockMarker)
	t.extraHist[0]++
	t.n++
***REMOVED***

func (t *tokens) Slice() []token ***REMOVED***
	return t.tokens[:t.n]
***REMOVED***

// VarInt returns the tokens as varint encoded bytes.
func (t *tokens) VarInt() []byte ***REMOVED***
	var b = make([]byte, binary.MaxVarintLen32*int(t.n))
	var off int
	for _, v := range t.tokens[:t.n] ***REMOVED***
		off += binary.PutUvarint(b[off:], uint64(v))
	***REMOVED***
	return b[:off]
***REMOVED***

// FromVarInt restores t to the varint encoded tokens provided.
// Any data in t is removed.
func (t *tokens) FromVarInt(b []byte) error ***REMOVED***
	var buf = bytes.NewReader(b)
	var toks []token
	for ***REMOVED***
		r, err := binary.ReadUvarint(buf)
		if err == io.EOF ***REMOVED***
			break
		***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		toks = append(toks, token(r))
	***REMOVED***
	t.indexTokens(toks)
	return nil
***REMOVED***

// Returns the type of a token
func (t token) typ() uint32 ***REMOVED*** return uint32(t) & typeMask ***REMOVED***

// Returns the literal of a literal token
func (t token) literal() uint8 ***REMOVED*** return uint8(t) ***REMOVED***

// Returns the extra offset of a match token
func (t token) offset() uint32 ***REMOVED*** return uint32(t) & offsetMask ***REMOVED***

func (t token) length() uint8 ***REMOVED*** return uint8(t >> lengthShift) ***REMOVED***

// Convert length to code.
func lengthCode(len uint8) uint8 ***REMOVED*** return lengthCodes[len] ***REMOVED***

// Returns the offset code corresponding to a specific offset
func offsetCode(off uint32) uint32 ***REMOVED***
	if false ***REMOVED***
		if off < uint32(len(offsetCodes)) ***REMOVED***
			return offsetCodes[off&255]
		***REMOVED*** else if off>>7 < uint32(len(offsetCodes)) ***REMOVED***
			return offsetCodes[(off>>7)&255] + 14
		***REMOVED*** else ***REMOVED***
			return offsetCodes[(off>>14)&255] + 28
		***REMOVED***
	***REMOVED***
	if off < uint32(len(offsetCodes)) ***REMOVED***
		return offsetCodes[uint8(off)]
	***REMOVED***
	return offsetCodes14[uint8(off>>7)]
***REMOVED***
