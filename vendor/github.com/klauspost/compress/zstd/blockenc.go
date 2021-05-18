// Copyright 2019+ Klaus Post. All rights reserved.
// License information can be found in the LICENSE file.
// Based on work by Yann Collet, released under BSD License.

package zstd

import (
	"errors"
	"fmt"
	"math"
	"math/bits"

	"github.com/klauspost/compress/huff0"
)

type blockEnc struct ***REMOVED***
	size       int
	literals   []byte
	sequences  []seq
	coders     seqCoders
	litEnc     *huff0.Scratch
	dictLitEnc *huff0.Scratch
	wr         bitWriter

	extraLits         int
	output            []byte
	recentOffsets     [3]uint32
	prevRecentOffsets [3]uint32

	last   bool
	lowMem bool
***REMOVED***

// init should be used once the block has been created.
// If called more than once, the effect is the same as calling reset.
func (b *blockEnc) init() ***REMOVED***
	if b.lowMem ***REMOVED***
		// 1K literals
		if cap(b.literals) < 1<<10 ***REMOVED***
			b.literals = make([]byte, 0, 1<<10)
		***REMOVED***
		const defSeqs = 20
		if cap(b.sequences) < defSeqs ***REMOVED***
			b.sequences = make([]seq, 0, defSeqs)
		***REMOVED***
		// 1K
		if cap(b.output) < 1<<10 ***REMOVED***
			b.output = make([]byte, 0, 1<<10)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if cap(b.literals) < maxCompressedBlockSize ***REMOVED***
			b.literals = make([]byte, 0, maxCompressedBlockSize)
		***REMOVED***
		const defSeqs = 200
		if cap(b.sequences) < defSeqs ***REMOVED***
			b.sequences = make([]seq, 0, defSeqs)
		***REMOVED***
		if cap(b.output) < maxCompressedBlockSize ***REMOVED***
			b.output = make([]byte, 0, maxCompressedBlockSize)
		***REMOVED***
	***REMOVED***

	if b.coders.mlEnc == nil ***REMOVED***
		b.coders.mlEnc = &fseEncoder***REMOVED******REMOVED***
		b.coders.mlPrev = &fseEncoder***REMOVED******REMOVED***
		b.coders.ofEnc = &fseEncoder***REMOVED******REMOVED***
		b.coders.ofPrev = &fseEncoder***REMOVED******REMOVED***
		b.coders.llEnc = &fseEncoder***REMOVED******REMOVED***
		b.coders.llPrev = &fseEncoder***REMOVED******REMOVED***
	***REMOVED***
	b.litEnc = &huff0.Scratch***REMOVED***WantLogLess: 4***REMOVED***
	b.reset(nil)
***REMOVED***

// initNewEncode can be used to reset offsets and encoders to the initial state.
func (b *blockEnc) initNewEncode() ***REMOVED***
	b.recentOffsets = [3]uint32***REMOVED***1, 4, 8***REMOVED***
	b.litEnc.Reuse = huff0.ReusePolicyNone
	b.coders.setPrev(nil, nil, nil)
***REMOVED***

// reset will reset the block for a new encode, but in the same stream,
// meaning that state will be carried over, but the block content is reset.
// If a previous block is provided, the recent offsets are carried over.
func (b *blockEnc) reset(prev *blockEnc) ***REMOVED***
	b.extraLits = 0
	b.literals = b.literals[:0]
	b.size = 0
	b.sequences = b.sequences[:0]
	b.output = b.output[:0]
	b.last = false
	if prev != nil ***REMOVED***
		b.recentOffsets = prev.prevRecentOffsets
	***REMOVED***
	b.dictLitEnc = nil
***REMOVED***

// reset will reset the block for a new encode, but in the same stream,
// meaning that state will be carried over, but the block content is reset.
// If a previous block is provided, the recent offsets are carried over.
func (b *blockEnc) swapEncoders(prev *blockEnc) ***REMOVED***
	b.coders.swap(&prev.coders)
	b.litEnc, prev.litEnc = prev.litEnc, b.litEnc
***REMOVED***

// blockHeader contains the information for a block header.
type blockHeader uint32

// setLast sets the 'last' indicator on a block.
func (h *blockHeader) setLast(b bool) ***REMOVED***
	if b ***REMOVED***
		*h = *h | 1
	***REMOVED*** else ***REMOVED***
		const mask = (1 << 24) - 2
		*h = *h & mask
	***REMOVED***
***REMOVED***

// setSize will store the compressed size of a block.
func (h *blockHeader) setSize(v uint32) ***REMOVED***
	const mask = 7
	*h = (*h)&mask | blockHeader(v<<3)
***REMOVED***

// setType sets the block type.
func (h *blockHeader) setType(t blockType) ***REMOVED***
	const mask = 1 | (((1 << 24) - 1) ^ 7)
	*h = (*h & mask) | blockHeader(t<<1)
***REMOVED***

// appendTo will append the block header to a slice.
func (h blockHeader) appendTo(b []byte) []byte ***REMOVED***
	return append(b, uint8(h), uint8(h>>8), uint8(h>>16))
***REMOVED***

// String returns a string representation of the block.
func (h blockHeader) String() string ***REMOVED***
	return fmt.Sprintf("Type: %d, Size: %d, Last:%t", (h>>1)&3, h>>3, h&1 == 1)
***REMOVED***

// literalsHeader contains literals header information.
type literalsHeader uint64

// setType can be used to set the type of literal block.
func (h *literalsHeader) setType(t literalsBlockType) ***REMOVED***
	const mask = math.MaxUint64 - 3
	*h = (*h & mask) | literalsHeader(t)
***REMOVED***

// setSize can be used to set a single size, for uncompressed and RLE content.
func (h *literalsHeader) setSize(regenLen int) ***REMOVED***
	inBits := bits.Len32(uint32(regenLen))
	// Only retain 2 bits
	const mask = 3
	lh := uint64(*h & mask)
	switch ***REMOVED***
	case inBits < 5:
		lh |= (uint64(regenLen) << 3) | (1 << 60)
		if debug ***REMOVED***
			got := int(lh>>3) & 0xff
			if got != regenLen ***REMOVED***
				panic(fmt.Sprint("litRegenSize = ", regenLen, "(want) != ", got, "(got)"))
			***REMOVED***
		***REMOVED***
	case inBits < 12:
		lh |= (1 << 2) | (uint64(regenLen) << 4) | (2 << 60)
	case inBits < 20:
		lh |= (3 << 2) | (uint64(regenLen) << 4) | (3 << 60)
	default:
		panic(fmt.Errorf("internal error: block too big (%d)", regenLen))
	***REMOVED***
	*h = literalsHeader(lh)
***REMOVED***

// setSizes will set the size of a compressed literals section and the input length.
func (h *literalsHeader) setSizes(compLen, inLen int, single bool) ***REMOVED***
	compBits, inBits := bits.Len32(uint32(compLen)), bits.Len32(uint32(inLen))
	// Only retain 2 bits
	const mask = 3
	lh := uint64(*h & mask)
	switch ***REMOVED***
	case compBits <= 10 && inBits <= 10:
		if !single ***REMOVED***
			lh |= 1 << 2
		***REMOVED***
		lh |= (uint64(inLen) << 4) | (uint64(compLen) << (10 + 4)) | (3 << 60)
		if debug ***REMOVED***
			const mmask = (1 << 24) - 1
			n := (lh >> 4) & mmask
			if int(n&1023) != inLen ***REMOVED***
				panic(fmt.Sprint("regensize:", int(n&1023), "!=", inLen, inBits))
			***REMOVED***
			if int(n>>10) != compLen ***REMOVED***
				panic(fmt.Sprint("compsize:", int(n>>10), "!=", compLen, compBits))
			***REMOVED***
		***REMOVED***
	case compBits <= 14 && inBits <= 14:
		lh |= (2 << 2) | (uint64(inLen) << 4) | (uint64(compLen) << (14 + 4)) | (4 << 60)
		if single ***REMOVED***
			panic("single stream used with more than 10 bits length.")
		***REMOVED***
	case compBits <= 18 && inBits <= 18:
		lh |= (3 << 2) | (uint64(inLen) << 4) | (uint64(compLen) << (18 + 4)) | (5 << 60)
		if single ***REMOVED***
			panic("single stream used with more than 10 bits length.")
		***REMOVED***
	default:
		panic("internal error: block too big")
	***REMOVED***
	*h = literalsHeader(lh)
***REMOVED***

// appendTo will append the literals header to a byte slice.
func (h literalsHeader) appendTo(b []byte) []byte ***REMOVED***
	size := uint8(h >> 60)
	switch size ***REMOVED***
	case 1:
		b = append(b, uint8(h))
	case 2:
		b = append(b, uint8(h), uint8(h>>8))
	case 3:
		b = append(b, uint8(h), uint8(h>>8), uint8(h>>16))
	case 4:
		b = append(b, uint8(h), uint8(h>>8), uint8(h>>16), uint8(h>>24))
	case 5:
		b = append(b, uint8(h), uint8(h>>8), uint8(h>>16), uint8(h>>24), uint8(h>>32))
	default:
		panic(fmt.Errorf("internal error: literalsHeader has invalid size (%d)", size))
	***REMOVED***
	return b
***REMOVED***

// size returns the output size with currently set values.
func (h literalsHeader) size() int ***REMOVED***
	return int(h >> 60)
***REMOVED***

func (h literalsHeader) String() string ***REMOVED***
	return fmt.Sprintf("Type: %d, SizeFormat: %d, Size: 0x%d, Bytes:%d", literalsBlockType(h&3), (h>>2)&3, h&((1<<60)-1)>>4, h>>60)
***REMOVED***

// pushOffsets will push the recent offsets to the backup store.
func (b *blockEnc) pushOffsets() ***REMOVED***
	b.prevRecentOffsets = b.recentOffsets
***REMOVED***

// pushOffsets will push the recent offsets to the backup store.
func (b *blockEnc) popOffsets() ***REMOVED***
	b.recentOffsets = b.prevRecentOffsets
***REMOVED***

// matchOffset will adjust recent offsets and return the adjusted one,
// if it matches a previous offset.
func (b *blockEnc) matchOffset(offset, lits uint32) uint32 ***REMOVED***
	// Check if offset is one of the recent offsets.
	// Adjusts the output offset accordingly.
	// Gives a tiny bit of compression, typically around 1%.
	if true ***REMOVED***
		if lits > 0 ***REMOVED***
			switch offset ***REMOVED***
			case b.recentOffsets[0]:
				offset = 1
			case b.recentOffsets[1]:
				b.recentOffsets[1] = b.recentOffsets[0]
				b.recentOffsets[0] = offset
				offset = 2
			case b.recentOffsets[2]:
				b.recentOffsets[2] = b.recentOffsets[1]
				b.recentOffsets[1] = b.recentOffsets[0]
				b.recentOffsets[0] = offset
				offset = 3
			default:
				b.recentOffsets[2] = b.recentOffsets[1]
				b.recentOffsets[1] = b.recentOffsets[0]
				b.recentOffsets[0] = offset
				offset += 3
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			switch offset ***REMOVED***
			case b.recentOffsets[1]:
				b.recentOffsets[1] = b.recentOffsets[0]
				b.recentOffsets[0] = offset
				offset = 1
			case b.recentOffsets[2]:
				b.recentOffsets[2] = b.recentOffsets[1]
				b.recentOffsets[1] = b.recentOffsets[0]
				b.recentOffsets[0] = offset
				offset = 2
			case b.recentOffsets[0] - 1:
				b.recentOffsets[2] = b.recentOffsets[1]
				b.recentOffsets[1] = b.recentOffsets[0]
				b.recentOffsets[0] = offset
				offset = 3
			default:
				b.recentOffsets[2] = b.recentOffsets[1]
				b.recentOffsets[1] = b.recentOffsets[0]
				b.recentOffsets[0] = offset
				offset += 3
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		offset += 3
	***REMOVED***
	return offset
***REMOVED***

// encodeRaw can be used to set the output to a raw representation of supplied bytes.
func (b *blockEnc) encodeRaw(a []byte) ***REMOVED***
	var bh blockHeader
	bh.setLast(b.last)
	bh.setSize(uint32(len(a)))
	bh.setType(blockTypeRaw)
	b.output = bh.appendTo(b.output[:0])
	b.output = append(b.output, a...)
	if debug ***REMOVED***
		println("Adding RAW block, length", len(a), "last:", b.last)
	***REMOVED***
***REMOVED***

// encodeRaw can be used to set the output to a raw representation of supplied bytes.
func (b *blockEnc) encodeRawTo(dst, src []byte) []byte ***REMOVED***
	var bh blockHeader
	bh.setLast(b.last)
	bh.setSize(uint32(len(src)))
	bh.setType(blockTypeRaw)
	dst = bh.appendTo(dst)
	dst = append(dst, src...)
	if debug ***REMOVED***
		println("Adding RAW block, length", len(src), "last:", b.last)
	***REMOVED***
	return dst
***REMOVED***

// encodeLits can be used if the block is only litLen.
func (b *blockEnc) encodeLits(lits []byte, raw bool) error ***REMOVED***
	var bh blockHeader
	bh.setLast(b.last)
	bh.setSize(uint32(len(lits)))

	// Don't compress extremely small blocks
	if len(lits) < 8 || (len(lits) < 32 && b.dictLitEnc == nil) || raw ***REMOVED***
		if debug ***REMOVED***
			println("Adding RAW block, length", len(lits), "last:", b.last)
		***REMOVED***
		bh.setType(blockTypeRaw)
		b.output = bh.appendTo(b.output)
		b.output = append(b.output, lits...)
		return nil
	***REMOVED***

	var (
		out            []byte
		reUsed, single bool
		err            error
	)
	if b.dictLitEnc != nil ***REMOVED***
		b.litEnc.TransferCTable(b.dictLitEnc)
		b.litEnc.Reuse = huff0.ReusePolicyAllow
		b.dictLitEnc = nil
	***REMOVED***
	if len(lits) >= 1024 ***REMOVED***
		// Use 4 Streams.
		out, reUsed, err = huff0.Compress4X(lits, b.litEnc)
	***REMOVED*** else if len(lits) > 32 ***REMOVED***
		// Use 1 stream
		single = true
		out, reUsed, err = huff0.Compress1X(lits, b.litEnc)
	***REMOVED*** else ***REMOVED***
		err = huff0.ErrIncompressible
	***REMOVED***

	switch err ***REMOVED***
	case huff0.ErrIncompressible:
		if debug ***REMOVED***
			println("Adding RAW block, length", len(lits), "last:", b.last)
		***REMOVED***
		bh.setType(blockTypeRaw)
		b.output = bh.appendTo(b.output)
		b.output = append(b.output, lits...)
		return nil
	case huff0.ErrUseRLE:
		if debug ***REMOVED***
			println("Adding RLE block, length", len(lits))
		***REMOVED***
		bh.setType(blockTypeRLE)
		b.output = bh.appendTo(b.output)
		b.output = append(b.output, lits[0])
		return nil
	case nil:
	default:
		return err
	***REMOVED***
	// Compressed...
	// Now, allow reuse
	b.litEnc.Reuse = huff0.ReusePolicyAllow
	bh.setType(blockTypeCompressed)
	var lh literalsHeader
	if reUsed ***REMOVED***
		if debug ***REMOVED***
			println("Reused tree, compressed to", len(out))
		***REMOVED***
		lh.setType(literalsBlockTreeless)
	***REMOVED*** else ***REMOVED***
		if debug ***REMOVED***
			println("New tree, compressed to", len(out), "tree size:", len(b.litEnc.OutTable))
		***REMOVED***
		lh.setType(literalsBlockCompressed)
	***REMOVED***
	// Set sizes
	lh.setSizes(len(out), len(lits), single)
	bh.setSize(uint32(len(out) + lh.size() + 1))

	// Write block headers.
	b.output = bh.appendTo(b.output)
	b.output = lh.appendTo(b.output)
	// Add compressed data.
	b.output = append(b.output, out...)
	// No sequences.
	b.output = append(b.output, 0)
	return nil
***REMOVED***

// fuzzFseEncoder can be used to fuzz the FSE encoder.
func fuzzFseEncoder(data []byte) int ***REMOVED***
	if len(data) > maxSequences || len(data) < 2 ***REMOVED***
		return 0
	***REMOVED***
	enc := fseEncoder***REMOVED******REMOVED***
	hist := enc.Histogram()[:256]
	maxSym := uint8(0)
	for i, v := range data ***REMOVED***
		v = v & 63
		data[i] = v
		hist[v]++
		if v > maxSym ***REMOVED***
			maxSym = v
		***REMOVED***
	***REMOVED***
	if maxSym == 0 ***REMOVED***
		// All 0
		return 0
	***REMOVED***
	maxCount := func(a []uint32) int ***REMOVED***
		var max uint32
		for _, v := range a ***REMOVED***
			if v > max ***REMOVED***
				max = v
			***REMOVED***
		***REMOVED***
		return int(max)
	***REMOVED***
	cnt := maxCount(hist[:maxSym])
	if cnt == len(data) ***REMOVED***
		// RLE
		return 0
	***REMOVED***
	enc.HistogramFinished(maxSym, cnt)
	err := enc.normalizeCount(len(data))
	if err != nil ***REMOVED***
		return 0
	***REMOVED***
	_, err = enc.writeCount(nil)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return 1
***REMOVED***

// encode will encode the block and append the output in b.output.
// Previous offset codes must be pushed if more blocks are expected.
func (b *blockEnc) encode(org []byte, raw, rawAllLits bool) error ***REMOVED***
	if len(b.sequences) == 0 ***REMOVED***
		return b.encodeLits(b.literals, rawAllLits)
	***REMOVED***
	// We want some difference to at least account for the headers.
	saved := b.size - len(b.literals) - (b.size >> 5)
	if saved < 16 ***REMOVED***
		if org == nil ***REMOVED***
			return errIncompressible
		***REMOVED***
		b.popOffsets()
		return b.encodeLits(org, rawAllLits)
	***REMOVED***

	var bh blockHeader
	var lh literalsHeader
	bh.setLast(b.last)
	bh.setType(blockTypeCompressed)
	// Store offset of the block header. Needed when we know the size.
	bhOffset := len(b.output)
	b.output = bh.appendTo(b.output)

	var (
		out            []byte
		reUsed, single bool
		err            error
	)
	if b.dictLitEnc != nil ***REMOVED***
		b.litEnc.TransferCTable(b.dictLitEnc)
		b.litEnc.Reuse = huff0.ReusePolicyAllow
		b.dictLitEnc = nil
	***REMOVED***
	if len(b.literals) >= 1024 && !raw ***REMOVED***
		// Use 4 Streams.
		out, reUsed, err = huff0.Compress4X(b.literals, b.litEnc)
	***REMOVED*** else if len(b.literals) > 32 && !raw ***REMOVED***
		// Use 1 stream
		single = true
		out, reUsed, err = huff0.Compress1X(b.literals, b.litEnc)
	***REMOVED*** else ***REMOVED***
		err = huff0.ErrIncompressible
	***REMOVED***

	switch err ***REMOVED***
	case huff0.ErrIncompressible:
		lh.setType(literalsBlockRaw)
		lh.setSize(len(b.literals))
		b.output = lh.appendTo(b.output)
		b.output = append(b.output, b.literals...)
		if debug ***REMOVED***
			println("Adding literals RAW, length", len(b.literals))
		***REMOVED***
	case huff0.ErrUseRLE:
		lh.setType(literalsBlockRLE)
		lh.setSize(len(b.literals))
		b.output = lh.appendTo(b.output)
		b.output = append(b.output, b.literals[0])
		if debug ***REMOVED***
			println("Adding literals RLE")
		***REMOVED***
	case nil:
		// Compressed litLen...
		if reUsed ***REMOVED***
			if debug ***REMOVED***
				println("reused tree")
			***REMOVED***
			lh.setType(literalsBlockTreeless)
		***REMOVED*** else ***REMOVED***
			if debug ***REMOVED***
				println("new tree, size:", len(b.litEnc.OutTable))
			***REMOVED***
			lh.setType(literalsBlockCompressed)
			if debug ***REMOVED***
				_, _, err := huff0.ReadTable(out, nil)
				if err != nil ***REMOVED***
					panic(err)
				***REMOVED***
			***REMOVED***
		***REMOVED***
		lh.setSizes(len(out), len(b.literals), single)
		if debug ***REMOVED***
			printf("Compressed %d literals to %d bytes", len(b.literals), len(out))
			println("Adding literal header:", lh)
		***REMOVED***
		b.output = lh.appendTo(b.output)
		b.output = append(b.output, out...)
		b.litEnc.Reuse = huff0.ReusePolicyAllow
		if debug ***REMOVED***
			println("Adding literals compressed")
		***REMOVED***
	default:
		if debug ***REMOVED***
			println("Adding literals ERROR:", err)
		***REMOVED***
		return err
	***REMOVED***
	// Sequence compression

	// Write the number of sequences
	switch ***REMOVED***
	case len(b.sequences) < 128:
		b.output = append(b.output, uint8(len(b.sequences)))
	case len(b.sequences) < 0x7f00: // TODO: this could be wrong
		n := len(b.sequences)
		b.output = append(b.output, 128+uint8(n>>8), uint8(n))
	default:
		n := len(b.sequences) - 0x7f00
		b.output = append(b.output, 255, uint8(n), uint8(n>>8))
	***REMOVED***
	if debug ***REMOVED***
		println("Encoding", len(b.sequences), "sequences")
	***REMOVED***
	b.genCodes()
	llEnc := b.coders.llEnc
	ofEnc := b.coders.ofEnc
	mlEnc := b.coders.mlEnc
	err = llEnc.normalizeCount(len(b.sequences))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	err = ofEnc.normalizeCount(len(b.sequences))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	err = mlEnc.normalizeCount(len(b.sequences))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Choose the best compression mode for each type.
	// Will evaluate the new vs predefined and previous.
	chooseComp := func(cur, prev, preDef *fseEncoder) (*fseEncoder, seqCompMode) ***REMOVED***
		// See if predefined/previous is better
		hist := cur.count[:cur.symbolLen]
		nSize := cur.approxSize(hist) + cur.maxHeaderSize()
		predefSize := preDef.approxSize(hist)
		prevSize := prev.approxSize(hist)

		// Add a small penalty for new encoders.
		// Don't bother with extremely small (<2 byte gains).
		nSize = nSize + (nSize+2*8*16)>>4
		switch ***REMOVED***
		case predefSize <= prevSize && predefSize <= nSize || forcePreDef:
			if debug ***REMOVED***
				println("Using predefined", predefSize>>3, "<=", nSize>>3)
			***REMOVED***
			return preDef, compModePredefined
		case prevSize <= nSize:
			if debug ***REMOVED***
				println("Using previous", prevSize>>3, "<=", nSize>>3)
			***REMOVED***
			return prev, compModeRepeat
		default:
			if debug ***REMOVED***
				println("Using new, predef", predefSize>>3, ". previous:", prevSize>>3, ">", nSize>>3, "header max:", cur.maxHeaderSize()>>3, "bytes")
				println("tl:", cur.actualTableLog, "symbolLen:", cur.symbolLen, "norm:", cur.norm[:cur.symbolLen], "hist", cur.count[:cur.symbolLen])
			***REMOVED***
			return cur, compModeFSE
		***REMOVED***
	***REMOVED***

	// Write compression mode
	var mode uint8
	if llEnc.useRLE ***REMOVED***
		mode |= uint8(compModeRLE) << 6
		llEnc.setRLE(b.sequences[0].llCode)
		if debug ***REMOVED***
			println("llEnc.useRLE")
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		var m seqCompMode
		llEnc, m = chooseComp(llEnc, b.coders.llPrev, &fsePredefEnc[tableLiteralLengths])
		mode |= uint8(m) << 6
	***REMOVED***
	if ofEnc.useRLE ***REMOVED***
		mode |= uint8(compModeRLE) << 4
		ofEnc.setRLE(b.sequences[0].ofCode)
		if debug ***REMOVED***
			println("ofEnc.useRLE")
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		var m seqCompMode
		ofEnc, m = chooseComp(ofEnc, b.coders.ofPrev, &fsePredefEnc[tableOffsets])
		mode |= uint8(m) << 4
	***REMOVED***

	if mlEnc.useRLE ***REMOVED***
		mode |= uint8(compModeRLE) << 2
		mlEnc.setRLE(b.sequences[0].mlCode)
		if debug ***REMOVED***
			println("mlEnc.useRLE, code: ", b.sequences[0].mlCode, "value", b.sequences[0].matchLen)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		var m seqCompMode
		mlEnc, m = chooseComp(mlEnc, b.coders.mlPrev, &fsePredefEnc[tableMatchLengths])
		mode |= uint8(m) << 2
	***REMOVED***
	b.output = append(b.output, mode)
	if debug ***REMOVED***
		printf("Compression modes: 0b%b", mode)
	***REMOVED***
	b.output, err = llEnc.writeCount(b.output)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	start := len(b.output)
	b.output, err = ofEnc.writeCount(b.output)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if false ***REMOVED***
		println("block:", b.output[start:], "tablelog", ofEnc.actualTableLog, "maxcount:", ofEnc.maxCount)
		fmt.Printf("selected TableLog: %d, Symbol length: %d\n", ofEnc.actualTableLog, ofEnc.symbolLen)
		for i, v := range ofEnc.norm[:ofEnc.symbolLen] ***REMOVED***
			fmt.Printf("%3d: %5d -> %4d \n", i, ofEnc.count[i], v)
		***REMOVED***
	***REMOVED***
	b.output, err = mlEnc.writeCount(b.output)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Maybe in block?
	wr := &b.wr
	wr.reset(b.output)

	var ll, of, ml cState

	// Current sequence
	seq := len(b.sequences) - 1
	s := b.sequences[seq]
	llEnc.setBits(llBitsTable[:])
	mlEnc.setBits(mlBitsTable[:])
	ofEnc.setBits(nil)

	llTT, ofTT, mlTT := llEnc.ct.symbolTT[:256], ofEnc.ct.symbolTT[:256], mlEnc.ct.symbolTT[:256]

	// We have 3 bounds checks here (and in the loop).
	// Since we are iterating backwards it is kinda hard to avoid.
	llB, ofB, mlB := llTT[s.llCode], ofTT[s.ofCode], mlTT[s.mlCode]
	ll.init(wr, &llEnc.ct, llB)
	of.init(wr, &ofEnc.ct, ofB)
	wr.flush32()
	ml.init(wr, &mlEnc.ct, mlB)

	// Each of these lookups also generates a bounds check.
	wr.addBits32NC(s.litLen, llB.outBits)
	wr.addBits32NC(s.matchLen, mlB.outBits)
	wr.flush32()
	wr.addBits32NC(s.offset, ofB.outBits)
	if debugSequences ***REMOVED***
		println("Encoded seq", seq, s, "codes:", s.llCode, s.mlCode, s.ofCode, "states:", ll.state, ml.state, of.state, "bits:", llB, mlB, ofB)
	***REMOVED***
	seq--
	if llEnc.maxBits+mlEnc.maxBits+ofEnc.maxBits <= 32 ***REMOVED***
		// No need to flush (common)
		for seq >= 0 ***REMOVED***
			s = b.sequences[seq]
			wr.flush32()
			llB, ofB, mlB := llTT[s.llCode], ofTT[s.ofCode], mlTT[s.mlCode]
			// tabelog max is 8 for all.
			of.encode(ofB)
			ml.encode(mlB)
			ll.encode(llB)
			wr.flush32()

			// We checked that all can stay within 32 bits
			wr.addBits32NC(s.litLen, llB.outBits)
			wr.addBits32NC(s.matchLen, mlB.outBits)
			wr.addBits32NC(s.offset, ofB.outBits)

			if debugSequences ***REMOVED***
				println("Encoded seq", seq, s)
			***REMOVED***

			seq--
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		for seq >= 0 ***REMOVED***
			s = b.sequences[seq]
			wr.flush32()
			llB, ofB, mlB := llTT[s.llCode], ofTT[s.ofCode], mlTT[s.mlCode]
			// tabelog max is below 8 for each.
			of.encode(ofB)
			ml.encode(mlB)
			ll.encode(llB)
			wr.flush32()

			// ml+ll = max 32 bits total
			wr.addBits32NC(s.litLen, llB.outBits)
			wr.addBits32NC(s.matchLen, mlB.outBits)
			wr.flush32()
			wr.addBits32NC(s.offset, ofB.outBits)

			if debugSequences ***REMOVED***
				println("Encoded seq", seq, s)
			***REMOVED***

			seq--
		***REMOVED***
	***REMOVED***
	ml.flush(mlEnc.actualTableLog)
	of.flush(ofEnc.actualTableLog)
	ll.flush(llEnc.actualTableLog)
	err = wr.close()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	b.output = wr.out

	if len(b.output)-3-bhOffset >= b.size ***REMOVED***
		// Maybe even add a bigger margin.
		b.litEnc.Reuse = huff0.ReusePolicyNone
		return errIncompressible
	***REMOVED***

	// Size is output minus block header.
	bh.setSize(uint32(len(b.output)-bhOffset) - 3)
	if debug ***REMOVED***
		println("Rewriting block header", bh)
	***REMOVED***
	_ = bh.appendTo(b.output[bhOffset:bhOffset])
	b.coders.setPrev(llEnc, mlEnc, ofEnc)
	return nil
***REMOVED***

var errIncompressible = errors.New("incompressible")

func (b *blockEnc) genCodes() ***REMOVED***
	if len(b.sequences) == 0 ***REMOVED***
		// nothing to do
		return
	***REMOVED***

	if len(b.sequences) > math.MaxUint16 ***REMOVED***
		panic("can only encode up to 64K sequences")
	***REMOVED***
	// No bounds checks after here:
	llH := b.coders.llEnc.Histogram()[:256]
	ofH := b.coders.ofEnc.Histogram()[:256]
	mlH := b.coders.mlEnc.Histogram()[:256]
	for i := range llH ***REMOVED***
		llH[i] = 0
	***REMOVED***
	for i := range ofH ***REMOVED***
		ofH[i] = 0
	***REMOVED***
	for i := range mlH ***REMOVED***
		mlH[i] = 0
	***REMOVED***

	var llMax, ofMax, mlMax uint8
	for i, seq := range b.sequences ***REMOVED***
		v := llCode(seq.litLen)
		seq.llCode = v
		llH[v]++
		if v > llMax ***REMOVED***
			llMax = v
		***REMOVED***

		v = ofCode(seq.offset)
		seq.ofCode = v
		ofH[v]++
		if v > ofMax ***REMOVED***
			ofMax = v
		***REMOVED***

		v = mlCode(seq.matchLen)
		seq.mlCode = v
		mlH[v]++
		if v > mlMax ***REMOVED***
			mlMax = v
			if debugAsserts && mlMax > maxMatchLengthSymbol ***REMOVED***
				panic(fmt.Errorf("mlMax > maxMatchLengthSymbol (%d), matchlen: %d", mlMax, seq.matchLen))
			***REMOVED***
		***REMOVED***
		b.sequences[i] = seq
	***REMOVED***
	maxCount := func(a []uint32) int ***REMOVED***
		var max uint32
		for _, v := range a ***REMOVED***
			if v > max ***REMOVED***
				max = v
			***REMOVED***
		***REMOVED***
		return int(max)
	***REMOVED***
	if debugAsserts && mlMax > maxMatchLengthSymbol ***REMOVED***
		panic(fmt.Errorf("mlMax > maxMatchLengthSymbol (%d)", mlMax))
	***REMOVED***
	if debugAsserts && ofMax > maxOffsetBits ***REMOVED***
		panic(fmt.Errorf("ofMax > maxOffsetBits (%d)", ofMax))
	***REMOVED***
	if debugAsserts && llMax > maxLiteralLengthSymbol ***REMOVED***
		panic(fmt.Errorf("llMax > maxLiteralLengthSymbol (%d)", llMax))
	***REMOVED***

	b.coders.mlEnc.HistogramFinished(mlMax, maxCount(mlH[:mlMax+1]))
	b.coders.ofEnc.HistogramFinished(ofMax, maxCount(ofH[:ofMax+1]))
	b.coders.llEnc.HistogramFinished(llMax, maxCount(llH[:llMax+1]))
***REMOVED***
