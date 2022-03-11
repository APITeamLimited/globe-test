// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flate

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	// The largest offset code.
	offsetCodeCount = 30

	// The special code used to mark the end of a block.
	endBlockMarker = 256

	// The first length code.
	lengthCodesStart = 257

	// The number of codegen codes.
	codegenCodeCount = 19
	badCode          = 255

	// bufferFlushSize indicates the buffer size
	// after which bytes are flushed to the writer.
	// Should preferably be a multiple of 6, since
	// we accumulate 6 bytes between writes to the buffer.
	bufferFlushSize = 246

	// bufferSize is the actual output byte buffer size.
	// It must have additional headroom for a flush
	// which can contain up to 8 bytes.
	bufferSize = bufferFlushSize + 8
)

// The number of extra bits needed by length code X - LENGTH_CODES_START.
var lengthExtraBits = [32]int8***REMOVED***
	/* 257 */ 0, 0, 0,
	/* 260 */ 0, 0, 0, 0, 0, 1, 1, 1, 1, 2,
	/* 270 */ 2, 2, 2, 3, 3, 3, 3, 4, 4, 4,
	/* 280 */ 4, 5, 5, 5, 5, 0,
***REMOVED***

// The length indicated by length code X - LENGTH_CODES_START.
var lengthBase = [32]uint8***REMOVED***
	0, 1, 2, 3, 4, 5, 6, 7, 8, 10,
	12, 14, 16, 20, 24, 28, 32, 40, 48, 56,
	64, 80, 96, 112, 128, 160, 192, 224, 255,
***REMOVED***

// offset code word extra bits.
var offsetExtraBits = [64]int8***REMOVED***
	0, 0, 0, 0, 1, 1, 2, 2, 3, 3,
	4, 4, 5, 5, 6, 6, 7, 7, 8, 8,
	9, 9, 10, 10, 11, 11, 12, 12, 13, 13,
	/* extended window */
	14, 14, 15, 15, 16, 16, 17, 17, 18, 18, 19, 19, 20, 20,
***REMOVED***

var offsetCombined = [32]uint32***REMOVED******REMOVED***

func init() ***REMOVED***
	var offsetBase = [64]uint32***REMOVED***
		/* normal deflate */
		0x000000, 0x000001, 0x000002, 0x000003, 0x000004,
		0x000006, 0x000008, 0x00000c, 0x000010, 0x000018,
		0x000020, 0x000030, 0x000040, 0x000060, 0x000080,
		0x0000c0, 0x000100, 0x000180, 0x000200, 0x000300,
		0x000400, 0x000600, 0x000800, 0x000c00, 0x001000,
		0x001800, 0x002000, 0x003000, 0x004000, 0x006000,

		/* extended window */
		0x008000, 0x00c000, 0x010000, 0x018000, 0x020000,
		0x030000, 0x040000, 0x060000, 0x080000, 0x0c0000,
		0x100000, 0x180000, 0x200000, 0x300000,
	***REMOVED***

	for i := range offsetCombined[:] ***REMOVED***
		// Don't use extended window values...
		if offsetBase[i] > 0x006000 ***REMOVED***
			continue
		***REMOVED***
		offsetCombined[i] = uint32(offsetExtraBits[i])<<16 | (offsetBase[i])
	***REMOVED***
***REMOVED***

// The odd order in which the codegen code sizes are written.
var codegenOrder = []uint32***REMOVED***16, 17, 18, 0, 8, 7, 9, 6, 10, 5, 11, 4, 12, 3, 13, 2, 14, 1, 15***REMOVED***

type huffmanBitWriter struct ***REMOVED***
	// writer is the underlying writer.
	// Do not use it directly; use the write method, which ensures
	// that Write errors are sticky.
	writer io.Writer

	// Data waiting to be written is bytes[0:nbytes]
	// and then the low nbits of bits.
	bits            uint64
	nbits           uint16
	nbytes          uint8
	lastHuffMan     bool
	literalEncoding *huffmanEncoder
	tmpLitEncoding  *huffmanEncoder
	offsetEncoding  *huffmanEncoder
	codegenEncoding *huffmanEncoder
	err             error
	lastHeader      int
	// Set between 0 (reused block can be up to 2x the size)
	logNewTablePenalty uint
	bytes              [256 + 8]byte
	literalFreq        [lengthCodesStart + 32]uint16
	offsetFreq         [32]uint16
	codegenFreq        [codegenCodeCount]uint16

	// codegen must have an extra space for the final symbol.
	codegen [literalCount + offsetCodeCount + 1]uint8
***REMOVED***

// Huffman reuse.
//
// The huffmanBitWriter supports reusing huffman tables and thereby combining block sections.
//
// This is controlled by several variables:
//
// If lastHeader is non-zero the Huffman table can be reused.
// This also indicates that a Huffman table has been generated that can output all
// possible symbols.
// It also indicates that an EOB has not yet been emitted, so if a new tabel is generated
// an EOB with the previous table must be written.
//
// If lastHuffMan is set, a table for outputting literals has been generated and offsets are invalid.
//
// An incoming block estimates the output size of a new table using a 'fresh' by calculating the
// optimal size and adding a penalty in 'logNewTablePenalty'.
// A Huffman table is not optimal, which is why we add a penalty, and generating a new table
// is slower both for compression and decompression.

func newHuffmanBitWriter(w io.Writer) *huffmanBitWriter ***REMOVED***
	return &huffmanBitWriter***REMOVED***
		writer:          w,
		literalEncoding: newHuffmanEncoder(literalCount),
		tmpLitEncoding:  newHuffmanEncoder(literalCount),
		codegenEncoding: newHuffmanEncoder(codegenCodeCount),
		offsetEncoding:  newHuffmanEncoder(offsetCodeCount),
	***REMOVED***
***REMOVED***

func (w *huffmanBitWriter) reset(writer io.Writer) ***REMOVED***
	w.writer = writer
	w.bits, w.nbits, w.nbytes, w.err = 0, 0, 0, nil
	w.lastHeader = 0
	w.lastHuffMan = false
***REMOVED***

func (w *huffmanBitWriter) canReuse(t *tokens) (offsets, lits bool) ***REMOVED***
	offsets, lits = true, true
	a := t.offHist[:offsetCodeCount]
	b := w.offsetFreq[:len(a)]
	for i := range a ***REMOVED***
		if b[i] == 0 && a[i] != 0 ***REMOVED***
			offsets = false
			break
		***REMOVED***
	***REMOVED***

	a = t.extraHist[:literalCount-256]
	b = w.literalFreq[256:literalCount]
	b = b[:len(a)]
	for i := range a ***REMOVED***
		if b[i] == 0 && a[i] != 0 ***REMOVED***
			lits = false
			break
		***REMOVED***
	***REMOVED***
	if lits ***REMOVED***
		a = t.litHist[:]
		b = w.literalFreq[:len(a)]
		for i := range a ***REMOVED***
			if b[i] == 0 && a[i] != 0 ***REMOVED***
				lits = false
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (w *huffmanBitWriter) flush() ***REMOVED***
	if w.err != nil ***REMOVED***
		w.nbits = 0
		return
	***REMOVED***
	if w.lastHeader > 0 ***REMOVED***
		// We owe an EOB
		w.writeCode(w.literalEncoding.codes[endBlockMarker])
		w.lastHeader = 0
	***REMOVED***
	n := w.nbytes
	for w.nbits != 0 ***REMOVED***
		w.bytes[n] = byte(w.bits)
		w.bits >>= 8
		if w.nbits > 8 ***REMOVED*** // Avoid underflow
			w.nbits -= 8
		***REMOVED*** else ***REMOVED***
			w.nbits = 0
		***REMOVED***
		n++
	***REMOVED***
	w.bits = 0
	w.write(w.bytes[:n])
	w.nbytes = 0
***REMOVED***

func (w *huffmanBitWriter) write(b []byte) ***REMOVED***
	if w.err != nil ***REMOVED***
		return
	***REMOVED***
	_, w.err = w.writer.Write(b)
***REMOVED***

func (w *huffmanBitWriter) writeBits(b int32, nb uint16) ***REMOVED***
	w.bits |= uint64(b) << w.nbits
	w.nbits += nb
	if w.nbits >= 48 ***REMOVED***
		w.writeOutBits()
	***REMOVED***
***REMOVED***

func (w *huffmanBitWriter) writeBytes(bytes []byte) ***REMOVED***
	if w.err != nil ***REMOVED***
		return
	***REMOVED***
	n := w.nbytes
	if w.nbits&7 != 0 ***REMOVED***
		w.err = InternalError("writeBytes with unfinished bits")
		return
	***REMOVED***
	for w.nbits != 0 ***REMOVED***
		w.bytes[n] = byte(w.bits)
		w.bits >>= 8
		w.nbits -= 8
		n++
	***REMOVED***
	if n != 0 ***REMOVED***
		w.write(w.bytes[:n])
	***REMOVED***
	w.nbytes = 0
	w.write(bytes)
***REMOVED***

// RFC 1951 3.2.7 specifies a special run-length encoding for specifying
// the literal and offset lengths arrays (which are concatenated into a single
// array).  This method generates that run-length encoding.
//
// The result is written into the codegen array, and the frequencies
// of each code is written into the codegenFreq array.
// Codes 0-15 are single byte codes. Codes 16-18 are followed by additional
// information. Code badCode is an end marker
//
//  numLiterals      The number of literals in literalEncoding
//  numOffsets       The number of offsets in offsetEncoding
//  litenc, offenc   The literal and offset encoder to use
func (w *huffmanBitWriter) generateCodegen(numLiterals int, numOffsets int, litEnc, offEnc *huffmanEncoder) ***REMOVED***
	for i := range w.codegenFreq ***REMOVED***
		w.codegenFreq[i] = 0
	***REMOVED***
	// Note that we are using codegen both as a temporary variable for holding
	// a copy of the frequencies, and as the place where we put the result.
	// This is fine because the output is always shorter than the input used
	// so far.
	codegen := w.codegen[:] // cache
	// Copy the concatenated code sizes to codegen. Put a marker at the end.
	cgnl := codegen[:numLiterals]
	for i := range cgnl ***REMOVED***
		cgnl[i] = uint8(litEnc.codes[i].len)
	***REMOVED***

	cgnl = codegen[numLiterals : numLiterals+numOffsets]
	for i := range cgnl ***REMOVED***
		cgnl[i] = uint8(offEnc.codes[i].len)
	***REMOVED***
	codegen[numLiterals+numOffsets] = badCode

	size := codegen[0]
	count := 1
	outIndex := 0
	for inIndex := 1; size != badCode; inIndex++ ***REMOVED***
		// INVARIANT: We have seen "count" copies of size that have not yet
		// had output generated for them.
		nextSize := codegen[inIndex]
		if nextSize == size ***REMOVED***
			count++
			continue
		***REMOVED***
		// We need to generate codegen indicating "count" of size.
		if size != 0 ***REMOVED***
			codegen[outIndex] = size
			outIndex++
			w.codegenFreq[size]++
			count--
			for count >= 3 ***REMOVED***
				n := 6
				if n > count ***REMOVED***
					n = count
				***REMOVED***
				codegen[outIndex] = 16
				outIndex++
				codegen[outIndex] = uint8(n - 3)
				outIndex++
				w.codegenFreq[16]++
				count -= n
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			for count >= 11 ***REMOVED***
				n := 138
				if n > count ***REMOVED***
					n = count
				***REMOVED***
				codegen[outIndex] = 18
				outIndex++
				codegen[outIndex] = uint8(n - 11)
				outIndex++
				w.codegenFreq[18]++
				count -= n
			***REMOVED***
			if count >= 3 ***REMOVED***
				// count >= 3 && count <= 10
				codegen[outIndex] = 17
				outIndex++
				codegen[outIndex] = uint8(count - 3)
				outIndex++
				w.codegenFreq[17]++
				count = 0
			***REMOVED***
		***REMOVED***
		count--
		for ; count >= 0; count-- ***REMOVED***
			codegen[outIndex] = size
			outIndex++
			w.codegenFreq[size]++
		***REMOVED***
		// Set up invariant for next time through the loop.
		size = nextSize
		count = 1
	***REMOVED***
	// Marker indicating the end of the codegen.
	codegen[outIndex] = badCode
***REMOVED***

func (w *huffmanBitWriter) codegens() int ***REMOVED***
	numCodegens := len(w.codegenFreq)
	for numCodegens > 4 && w.codegenFreq[codegenOrder[numCodegens-1]] == 0 ***REMOVED***
		numCodegens--
	***REMOVED***
	return numCodegens
***REMOVED***

func (w *huffmanBitWriter) headerSize() (size, numCodegens int) ***REMOVED***
	numCodegens = len(w.codegenFreq)
	for numCodegens > 4 && w.codegenFreq[codegenOrder[numCodegens-1]] == 0 ***REMOVED***
		numCodegens--
	***REMOVED***
	return 3 + 5 + 5 + 4 + (3 * numCodegens) +
		w.codegenEncoding.bitLength(w.codegenFreq[:]) +
		int(w.codegenFreq[16])*2 +
		int(w.codegenFreq[17])*3 +
		int(w.codegenFreq[18])*7, numCodegens
***REMOVED***

// dynamicSize returns the size of dynamically encoded data in bits.
func (w *huffmanBitWriter) dynamicReuseSize(litEnc, offEnc *huffmanEncoder) (size int) ***REMOVED***
	size = litEnc.bitLength(w.literalFreq[:]) +
		offEnc.bitLength(w.offsetFreq[:])
	return size
***REMOVED***

// dynamicSize returns the size of dynamically encoded data in bits.
func (w *huffmanBitWriter) dynamicSize(litEnc, offEnc *huffmanEncoder, extraBits int) (size, numCodegens int) ***REMOVED***
	header, numCodegens := w.headerSize()
	size = header +
		litEnc.bitLength(w.literalFreq[:]) +
		offEnc.bitLength(w.offsetFreq[:]) +
		extraBits
	return size, numCodegens
***REMOVED***

// extraBitSize will return the number of bits that will be written
// as "extra" bits on matches.
func (w *huffmanBitWriter) extraBitSize() int ***REMOVED***
	total := 0
	for i, n := range w.literalFreq[257:literalCount] ***REMOVED***
		total += int(n) * int(lengthExtraBits[i&31])
	***REMOVED***
	for i, n := range w.offsetFreq[:offsetCodeCount] ***REMOVED***
		total += int(n) * int(offsetExtraBits[i&31])
	***REMOVED***
	return total
***REMOVED***

// fixedSize returns the size of dynamically encoded data in bits.
func (w *huffmanBitWriter) fixedSize(extraBits int) int ***REMOVED***
	return 3 +
		fixedLiteralEncoding.bitLength(w.literalFreq[:]) +
		fixedOffsetEncoding.bitLength(w.offsetFreq[:]) +
		extraBits
***REMOVED***

// storedSize calculates the stored size, including header.
// The function returns the size in bits and whether the block
// fits inside a single block.
func (w *huffmanBitWriter) storedSize(in []byte) (int, bool) ***REMOVED***
	if in == nil ***REMOVED***
		return 0, false
	***REMOVED***
	if len(in) <= maxStoreBlockSize ***REMOVED***
		return (len(in) + 5) * 8, true
	***REMOVED***
	return 0, false
***REMOVED***

func (w *huffmanBitWriter) writeCode(c hcode) ***REMOVED***
	// The function does not get inlined if we "& 63" the shift.
	w.bits |= uint64(c.code) << w.nbits
	w.nbits += c.len
	if w.nbits >= 48 ***REMOVED***
		w.writeOutBits()
	***REMOVED***
***REMOVED***

// writeOutBits will write bits to the buffer.
func (w *huffmanBitWriter) writeOutBits() ***REMOVED***
	bits := w.bits
	w.bits >>= 48
	w.nbits -= 48
	n := w.nbytes

	// We over-write, but faster...
	binary.LittleEndian.PutUint64(w.bytes[n:], bits)
	n += 6

	if n >= bufferFlushSize ***REMOVED***
		if w.err != nil ***REMOVED***
			n = 0
			return
		***REMOVED***
		w.write(w.bytes[:n])
		n = 0
	***REMOVED***

	w.nbytes = n
***REMOVED***

// Write the header of a dynamic Huffman block to the output stream.
//
//  numLiterals  The number of literals specified in codegen
//  numOffsets   The number of offsets specified in codegen
//  numCodegens  The number of codegens used in codegen
func (w *huffmanBitWriter) writeDynamicHeader(numLiterals int, numOffsets int, numCodegens int, isEof bool) ***REMOVED***
	if w.err != nil ***REMOVED***
		return
	***REMOVED***
	var firstBits int32 = 4
	if isEof ***REMOVED***
		firstBits = 5
	***REMOVED***
	w.writeBits(firstBits, 3)
	w.writeBits(int32(numLiterals-257), 5)
	w.writeBits(int32(numOffsets-1), 5)
	w.writeBits(int32(numCodegens-4), 4)

	for i := 0; i < numCodegens; i++ ***REMOVED***
		value := uint(w.codegenEncoding.codes[codegenOrder[i]].len)
		w.writeBits(int32(value), 3)
	***REMOVED***

	i := 0
	for ***REMOVED***
		var codeWord = uint32(w.codegen[i])
		i++
		if codeWord == badCode ***REMOVED***
			break
		***REMOVED***
		w.writeCode(w.codegenEncoding.codes[codeWord])

		switch codeWord ***REMOVED***
		case 16:
			w.writeBits(int32(w.codegen[i]), 2)
			i++
		case 17:
			w.writeBits(int32(w.codegen[i]), 3)
			i++
		case 18:
			w.writeBits(int32(w.codegen[i]), 7)
			i++
		***REMOVED***
	***REMOVED***
***REMOVED***

// writeStoredHeader will write a stored header.
// If the stored block is only used for EOF,
// it is replaced with a fixed huffman block.
func (w *huffmanBitWriter) writeStoredHeader(length int, isEof bool) ***REMOVED***
	if w.err != nil ***REMOVED***
		return
	***REMOVED***
	if w.lastHeader > 0 ***REMOVED***
		// We owe an EOB
		w.writeCode(w.literalEncoding.codes[endBlockMarker])
		w.lastHeader = 0
	***REMOVED***

	// To write EOF, use a fixed encoding block. 10 bits instead of 5 bytes.
	if length == 0 && isEof ***REMOVED***
		w.writeFixedHeader(isEof)
		// EOB: 7 bits, value: 0
		w.writeBits(0, 7)
		w.flush()
		return
	***REMOVED***

	var flag int32
	if isEof ***REMOVED***
		flag = 1
	***REMOVED***
	w.writeBits(flag, 3)
	w.flush()
	w.writeBits(int32(length), 16)
	w.writeBits(int32(^uint16(length)), 16)
***REMOVED***

func (w *huffmanBitWriter) writeFixedHeader(isEof bool) ***REMOVED***
	if w.err != nil ***REMOVED***
		return
	***REMOVED***
	if w.lastHeader > 0 ***REMOVED***
		// We owe an EOB
		w.writeCode(w.literalEncoding.codes[endBlockMarker])
		w.lastHeader = 0
	***REMOVED***

	// Indicate that we are a fixed Huffman block
	var value int32 = 2
	if isEof ***REMOVED***
		value = 3
	***REMOVED***
	w.writeBits(value, 3)
***REMOVED***

// writeBlock will write a block of tokens with the smallest encoding.
// The original input can be supplied, and if the huffman encoded data
// is larger than the original bytes, the data will be written as a
// stored block.
// If the input is nil, the tokens will always be Huffman encoded.
func (w *huffmanBitWriter) writeBlock(tokens *tokens, eof bool, input []byte) ***REMOVED***
	if w.err != nil ***REMOVED***
		return
	***REMOVED***

	tokens.AddEOB()
	if w.lastHeader > 0 ***REMOVED***
		// We owe an EOB
		w.writeCode(w.literalEncoding.codes[endBlockMarker])
		w.lastHeader = 0
	***REMOVED***
	numLiterals, numOffsets := w.indexTokens(tokens, false)
	w.generate(tokens)
	var extraBits int
	storedSize, storable := w.storedSize(input)
	if storable ***REMOVED***
		extraBits = w.extraBitSize()
	***REMOVED***

	// Figure out smallest code.
	// Fixed Huffman baseline.
	var literalEncoding = fixedLiteralEncoding
	var offsetEncoding = fixedOffsetEncoding
	var size = w.fixedSize(extraBits)

	// Dynamic Huffman?
	var numCodegens int

	// Generate codegen and codegenFrequencies, which indicates how to encode
	// the literalEncoding and the offsetEncoding.
	w.generateCodegen(numLiterals, numOffsets, w.literalEncoding, w.offsetEncoding)
	w.codegenEncoding.generate(w.codegenFreq[:], 7)
	dynamicSize, numCodegens := w.dynamicSize(w.literalEncoding, w.offsetEncoding, extraBits)

	if dynamicSize < size ***REMOVED***
		size = dynamicSize
		literalEncoding = w.literalEncoding
		offsetEncoding = w.offsetEncoding
	***REMOVED***

	// Stored bytes?
	if storable && storedSize < size ***REMOVED***
		w.writeStoredHeader(len(input), eof)
		w.writeBytes(input)
		return
	***REMOVED***

	// Huffman.
	if literalEncoding == fixedLiteralEncoding ***REMOVED***
		w.writeFixedHeader(eof)
	***REMOVED*** else ***REMOVED***
		w.writeDynamicHeader(numLiterals, numOffsets, numCodegens, eof)
	***REMOVED***

	// Write the tokens.
	w.writeTokens(tokens.Slice(), literalEncoding.codes, offsetEncoding.codes)
***REMOVED***

// writeBlockDynamic encodes a block using a dynamic Huffman table.
// This should be used if the symbols used have a disproportionate
// histogram distribution.
// If input is supplied and the compression savings are below 1/16th of the
// input size the block is stored.
func (w *huffmanBitWriter) writeBlockDynamic(tokens *tokens, eof bool, input []byte, sync bool) ***REMOVED***
	if w.err != nil ***REMOVED***
		return
	***REMOVED***

	sync = sync || eof
	if sync ***REMOVED***
		tokens.AddEOB()
	***REMOVED***

	// We cannot reuse pure huffman table, and must mark as EOF.
	if (w.lastHuffMan || eof) && w.lastHeader > 0 ***REMOVED***
		// We will not try to reuse.
		w.writeCode(w.literalEncoding.codes[endBlockMarker])
		w.lastHeader = 0
		w.lastHuffMan = false
	***REMOVED***
	if !sync ***REMOVED***
		tokens.Fill()
	***REMOVED***
	numLiterals, numOffsets := w.indexTokens(tokens, !sync)

	var size int
	// Check if we should reuse.
	if w.lastHeader > 0 ***REMOVED***
		// Estimate size for using a new table.
		// Use the previous header size as the best estimate.
		newSize := w.lastHeader + tokens.EstimatedBits()
		newSize += newSize >> w.logNewTablePenalty

		// The estimated size is calculated as an optimal table.
		// We add a penalty to make it more realistic and re-use a bit more.
		reuseSize := w.dynamicReuseSize(w.literalEncoding, w.offsetEncoding) + w.extraBitSize()

		// Check if a new table is better.
		if newSize < reuseSize ***REMOVED***
			// Write the EOB we owe.
			w.writeCode(w.literalEncoding.codes[endBlockMarker])
			size = newSize
			w.lastHeader = 0
		***REMOVED*** else ***REMOVED***
			size = reuseSize
		***REMOVED***
		// Check if we get a reasonable size decrease.
		if ssize, storable := w.storedSize(input); storable && ssize < (size+size>>4) ***REMOVED***
			w.writeStoredHeader(len(input), eof)
			w.writeBytes(input)
			w.lastHeader = 0
			return
		***REMOVED***
	***REMOVED***

	// We want a new block/table
	if w.lastHeader == 0 ***REMOVED***
		w.generate(tokens)
		// Generate codegen and codegenFrequencies, which indicates how to encode
		// the literalEncoding and the offsetEncoding.
		w.generateCodegen(numLiterals, numOffsets, w.literalEncoding, w.offsetEncoding)
		w.codegenEncoding.generate(w.codegenFreq[:], 7)
		var numCodegens int
		size, numCodegens = w.dynamicSize(w.literalEncoding, w.offsetEncoding, w.extraBitSize())
		// Store bytes, if we don't get a reasonable improvement.
		if ssize, storable := w.storedSize(input); storable && ssize < (size+size>>4) ***REMOVED***
			w.writeStoredHeader(len(input), eof)
			w.writeBytes(input)
			w.lastHeader = 0
			return
		***REMOVED***

		// Write Huffman table.
		w.writeDynamicHeader(numLiterals, numOffsets, numCodegens, eof)
		w.lastHeader, _ = w.headerSize()
		w.lastHuffMan = false
	***REMOVED***

	if sync ***REMOVED***
		w.lastHeader = 0
	***REMOVED***
	// Write the tokens.
	w.writeTokens(tokens.Slice(), w.literalEncoding.codes, w.offsetEncoding.codes)
***REMOVED***

// indexTokens indexes a slice of tokens, and updates
// literalFreq and offsetFreq, and generates literalEncoding
// and offsetEncoding.
// The number of literal and offset tokens is returned.
func (w *huffmanBitWriter) indexTokens(t *tokens, filled bool) (numLiterals, numOffsets int) ***REMOVED***
	copy(w.literalFreq[:], t.litHist[:])
	copy(w.literalFreq[256:], t.extraHist[:])
	copy(w.offsetFreq[:], t.offHist[:offsetCodeCount])

	if t.n == 0 ***REMOVED***
		return
	***REMOVED***
	if filled ***REMOVED***
		return maxNumLit, maxNumDist
	***REMOVED***
	// get the number of literals
	numLiterals = len(w.literalFreq)
	for w.literalFreq[numLiterals-1] == 0 ***REMOVED***
		numLiterals--
	***REMOVED***
	// get the number of offsets
	numOffsets = len(w.offsetFreq)
	for numOffsets > 0 && w.offsetFreq[numOffsets-1] == 0 ***REMOVED***
		numOffsets--
	***REMOVED***
	if numOffsets == 0 ***REMOVED***
		// We haven't found a single match. If we want to go with the dynamic encoding,
		// we should count at least one offset to be sure that the offset huffman tree could be encoded.
		w.offsetFreq[0] = 1
		numOffsets = 1
	***REMOVED***
	return
***REMOVED***

func (w *huffmanBitWriter) generate(t *tokens) ***REMOVED***
	w.literalEncoding.generate(w.literalFreq[:literalCount], 15)
	w.offsetEncoding.generate(w.offsetFreq[:offsetCodeCount], 15)
***REMOVED***

// writeTokens writes a slice of tokens to the output.
// codes for literal and offset encoding must be supplied.
func (w *huffmanBitWriter) writeTokens(tokens []token, leCodes, oeCodes []hcode) ***REMOVED***
	if w.err != nil ***REMOVED***
		return
	***REMOVED***
	if len(tokens) == 0 ***REMOVED***
		return
	***REMOVED***

	// Only last token should be endBlockMarker.
	var deferEOB bool
	if tokens[len(tokens)-1] == endBlockMarker ***REMOVED***
		tokens = tokens[:len(tokens)-1]
		deferEOB = true
	***REMOVED***

	// Create slices up to the next power of two to avoid bounds checks.
	lits := leCodes[:256]
	offs := oeCodes[:32]
	lengths := leCodes[lengthCodesStart:]
	lengths = lengths[:32]

	// Go 1.16 LOVES having these on stack.
	bits, nbits, nbytes := w.bits, w.nbits, w.nbytes

	for _, t := range tokens ***REMOVED***
		if t < matchType ***REMOVED***
			//w.writeCode(lits[t.literal()])
			c := lits[t.literal()]
			bits |= uint64(c.code) << nbits
			nbits += c.len
			if nbits >= 48 ***REMOVED***
				binary.LittleEndian.PutUint64(w.bytes[nbytes:], bits)
				//*(*uint64)(unsafe.Pointer(&w.bytes[nbytes])) = bits
				bits >>= 48
				nbits -= 48
				nbytes += 6
				if nbytes >= bufferFlushSize ***REMOVED***
					if w.err != nil ***REMOVED***
						nbytes = 0
						return
					***REMOVED***
					_, w.err = w.writer.Write(w.bytes[:nbytes])
					nbytes = 0
				***REMOVED***
			***REMOVED***
			continue
		***REMOVED***

		// Write the length
		length := t.length()
		lengthCode := lengthCode(length)
		if false ***REMOVED***
			w.writeCode(lengths[lengthCode&31])
		***REMOVED*** else ***REMOVED***
			// inlined
			c := lengths[lengthCode&31]
			bits |= uint64(c.code) << nbits
			nbits += c.len
			if nbits >= 48 ***REMOVED***
				binary.LittleEndian.PutUint64(w.bytes[nbytes:], bits)
				//*(*uint64)(unsafe.Pointer(&w.bytes[nbytes])) = bits
				bits >>= 48
				nbits -= 48
				nbytes += 6
				if nbytes >= bufferFlushSize ***REMOVED***
					if w.err != nil ***REMOVED***
						nbytes = 0
						return
					***REMOVED***
					_, w.err = w.writer.Write(w.bytes[:nbytes])
					nbytes = 0
				***REMOVED***
			***REMOVED***
		***REMOVED***

		extraLengthBits := uint16(lengthExtraBits[lengthCode&31])
		if extraLengthBits > 0 ***REMOVED***
			//w.writeBits(extraLength, extraLengthBits)
			extraLength := int32(length - lengthBase[lengthCode&31])
			bits |= uint64(extraLength) << nbits
			nbits += extraLengthBits
			if nbits >= 48 ***REMOVED***
				binary.LittleEndian.PutUint64(w.bytes[nbytes:], bits)
				//*(*uint64)(unsafe.Pointer(&w.bytes[nbytes])) = bits
				bits >>= 48
				nbits -= 48
				nbytes += 6
				if nbytes >= bufferFlushSize ***REMOVED***
					if w.err != nil ***REMOVED***
						nbytes = 0
						return
					***REMOVED***
					_, w.err = w.writer.Write(w.bytes[:nbytes])
					nbytes = 0
				***REMOVED***
			***REMOVED***
		***REMOVED***
		// Write the offset
		offset := t.offset()
		offsetCode := offset >> 16
		offset &= matchOffsetOnlyMask
		if false ***REMOVED***
			w.writeCode(offs[offsetCode&31])
		***REMOVED*** else ***REMOVED***
			// inlined
			c := offs[offsetCode]
			bits |= uint64(c.code) << nbits
			nbits += c.len
			if nbits >= 48 ***REMOVED***
				binary.LittleEndian.PutUint64(w.bytes[nbytes:], bits)
				//*(*uint64)(unsafe.Pointer(&w.bytes[nbytes])) = bits
				bits >>= 48
				nbits -= 48
				nbytes += 6
				if nbytes >= bufferFlushSize ***REMOVED***
					if w.err != nil ***REMOVED***
						nbytes = 0
						return
					***REMOVED***
					_, w.err = w.writer.Write(w.bytes[:nbytes])
					nbytes = 0
				***REMOVED***
			***REMOVED***
		***REMOVED***
		offsetComb := offsetCombined[offsetCode]
		if offsetComb > 1<<16 ***REMOVED***
			//w.writeBits(extraOffset, extraOffsetBits)
			bits |= uint64(offset&matchOffsetOnlyMask-(offsetComb&0xffff)) << nbits
			nbits += uint16(offsetComb >> 16)
			if nbits >= 48 ***REMOVED***
				binary.LittleEndian.PutUint64(w.bytes[nbytes:], bits)
				//*(*uint64)(unsafe.Pointer(&w.bytes[nbytes])) = bits
				bits >>= 48
				nbits -= 48
				nbytes += 6
				if nbytes >= bufferFlushSize ***REMOVED***
					if w.err != nil ***REMOVED***
						nbytes = 0
						return
					***REMOVED***
					_, w.err = w.writer.Write(w.bytes[:nbytes])
					nbytes = 0
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	// Restore...
	w.bits, w.nbits, w.nbytes = bits, nbits, nbytes

	if deferEOB ***REMOVED***
		w.writeCode(leCodes[endBlockMarker])
	***REMOVED***
***REMOVED***

// huffOffset is a static offset encoder used for huffman only encoding.
// It can be reused since we will not be encoding offset values.
var huffOffset *huffmanEncoder

func init() ***REMOVED***
	w := newHuffmanBitWriter(nil)
	w.offsetFreq[0] = 1
	huffOffset = newHuffmanEncoder(offsetCodeCount)
	huffOffset.generate(w.offsetFreq[:offsetCodeCount], 15)
***REMOVED***

// writeBlockHuff encodes a block of bytes as either
// Huffman encoded literals or uncompressed bytes if the
// results only gains very little from compression.
func (w *huffmanBitWriter) writeBlockHuff(eof bool, input []byte, sync bool) ***REMOVED***
	if w.err != nil ***REMOVED***
		return
	***REMOVED***

	// Clear histogram
	for i := range w.literalFreq[:] ***REMOVED***
		w.literalFreq[i] = 0
	***REMOVED***
	if !w.lastHuffMan ***REMOVED***
		for i := range w.offsetFreq[:] ***REMOVED***
			w.offsetFreq[i] = 0
		***REMOVED***
	***REMOVED***

	// Fill is rarely better...
	const fill = false
	const numLiterals = endBlockMarker + 1
	const numOffsets = 1

	// Add everything as literals
	// We have to estimate the header size.
	// Assume header is around 70 bytes:
	// https://stackoverflow.com/a/25454430
	const guessHeaderSizeBits = 70 * 8
	histogram(input, w.literalFreq[:numLiterals], fill)
	w.literalFreq[endBlockMarker] = 1
	w.tmpLitEncoding.generate(w.literalFreq[:numLiterals], 15)
	if fill ***REMOVED***
		// Clear fill...
		for i := range w.literalFreq[:numLiterals] ***REMOVED***
			w.literalFreq[i] = 0
		***REMOVED***
		histogram(input, w.literalFreq[:numLiterals], false)
	***REMOVED***
	estBits := w.tmpLitEncoding.canReuseBits(w.literalFreq[:numLiterals])
	estBits += w.lastHeader
	if w.lastHeader == 0 ***REMOVED***
		estBits += guessHeaderSizeBits
	***REMOVED***
	estBits += estBits >> w.logNewTablePenalty

	// Store bytes, if we don't get a reasonable improvement.
	ssize, storable := w.storedSize(input)
	if storable && ssize <= estBits ***REMOVED***
		w.writeStoredHeader(len(input), eof)
		w.writeBytes(input)
		return
	***REMOVED***

	if w.lastHeader > 0 ***REMOVED***
		reuseSize := w.literalEncoding.canReuseBits(w.literalFreq[:256])

		if estBits < reuseSize ***REMOVED***
			if debugDeflate ***REMOVED***
				//fmt.Println("not reusing, reuse:", reuseSize/8, "> new:", estBits/8, "- header est:", w.lastHeader/8)
			***REMOVED***
			// We owe an EOB
			w.writeCode(w.literalEncoding.codes[endBlockMarker])
			w.lastHeader = 0
		***REMOVED*** else if debugDeflate ***REMOVED***
			fmt.Println("reusing, reuse:", reuseSize/8, "> new:", estBits/8, "- header est:", w.lastHeader/8)
		***REMOVED***
	***REMOVED***

	count := 0
	if w.lastHeader == 0 ***REMOVED***
		// Use the temp encoding, so swap.
		w.literalEncoding, w.tmpLitEncoding = w.tmpLitEncoding, w.literalEncoding
		// Generate codegen and codegenFrequencies, which indicates how to encode
		// the literalEncoding and the offsetEncoding.
		w.generateCodegen(numLiterals, numOffsets, w.literalEncoding, huffOffset)
		w.codegenEncoding.generate(w.codegenFreq[:], 7)
		numCodegens := w.codegens()

		// Huffman.
		w.writeDynamicHeader(numLiterals, numOffsets, numCodegens, eof)
		w.lastHuffMan = true
		w.lastHeader, _ = w.headerSize()
		if debugDeflate ***REMOVED***
			count += w.lastHeader
			fmt.Println("header:", count/8)
		***REMOVED***
	***REMOVED***

	encoding := w.literalEncoding.codes[:256]
	// Go 1.16 LOVES having these on stack. At least 1.5x the speed.
	bits, nbits, nbytes := w.bits, w.nbits, w.nbytes
	for _, t := range input ***REMOVED***
		// Bitwriting inlined, ~30% speedup
		c := encoding[t]
		bits |= uint64(c.code) << nbits
		nbits += c.len
		if debugDeflate ***REMOVED***
			count += int(c.len)
		***REMOVED***
		if nbits >= 48 ***REMOVED***
			binary.LittleEndian.PutUint64(w.bytes[nbytes:], bits)
			//*(*uint64)(unsafe.Pointer(&w.bytes[nbytes])) = bits
			bits >>= 48
			nbits -= 48
			nbytes += 6
			if nbytes >= bufferFlushSize ***REMOVED***
				if w.err != nil ***REMOVED***
					nbytes = 0
					return
				***REMOVED***
				_, w.err = w.writer.Write(w.bytes[:nbytes])
				nbytes = 0
			***REMOVED***
		***REMOVED***
	***REMOVED***
	// Restore...
	w.bits, w.nbits, w.nbytes = bits, nbits, nbytes

	if debugDeflate ***REMOVED***
		fmt.Println("wrote", count/8, "bytes")
	***REMOVED***
	if eof || sync ***REMOVED***
		w.writeCode(w.literalEncoding.codes[endBlockMarker])
		w.lastHeader = 0
		w.lastHuffMan = false
	***REMOVED***
***REMOVED***
