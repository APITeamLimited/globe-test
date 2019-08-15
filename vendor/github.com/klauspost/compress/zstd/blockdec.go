// Copyright 2019+ Klaus Post. All rights reserved.
// License information can be found in the LICENSE file.
// Based on work by Yann Collet, released under BSD License.

package zstd

import (
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/klauspost/compress/huff0"
)

type blockType uint8

//go:generate stringer -type=blockType,literalsBlockType,seqCompMode,tableIndex

const (
	blockTypeRaw blockType = iota
	blockTypeRLE
	blockTypeCompressed
	blockTypeReserved
)

type literalsBlockType uint8

const (
	literalsBlockRaw literalsBlockType = iota
	literalsBlockRLE
	literalsBlockCompressed
	literalsBlockTreeless
)

const (
	// maxCompressedBlockSize is the biggest allowed compressed block size (128KB)
	maxCompressedBlockSize = 128 << 10

	// Maximum possible block size (all Raw+Uncompressed).
	maxBlockSize = (1 << 21) - 1

	// https://github.com/facebook/zstd/blob/dev/doc/zstd_compression_format.md#literals_section_header
	maxCompressedLiteralSize = 1 << 18
	maxRLELiteralSize        = 1 << 20
	maxMatchLen              = 131074
	maxSequences             = 0x7f00 + 0xffff

	// We support slightly less than the reference decoder to be able to
	// use ints on 32 bit archs.
	maxOffsetBits = 30
)

var (
	huffDecoderPool = sync.Pool***REMOVED***New: func() interface***REMOVED******REMOVED*** ***REMOVED***
		return &huff0.Scratch***REMOVED******REMOVED***
	***REMOVED******REMOVED***

	fseDecoderPool = sync.Pool***REMOVED***New: func() interface***REMOVED******REMOVED*** ***REMOVED***
		return &fseDecoder***REMOVED******REMOVED***
	***REMOVED******REMOVED***
)

type blockDec struct ***REMOVED***
	// Raw source data of the block.
	data []byte

	// Destination of the decoded data.
	dst []byte

	// Buffer for literals data.
	literalBuf []byte

	// Window size of the block.
	WindowSize uint64
	Type       blockType
	RLESize    uint32

	// Is this the last block of a frame?
	Last bool

	// Use less memory
	lowMem      bool
	history     chan *history
	input       chan struct***REMOVED******REMOVED***
	result      chan decodeOutput
	sequenceBuf []seq
	tmp         [4]byte
	err         error
***REMOVED***

func (b *blockDec) String() string ***REMOVED***
	if b == nil ***REMOVED***
		return "<nil>"
	***REMOVED***
	return fmt.Sprintf("Steam Size: %d, Type: %v, Last: %t, Window: %d", len(b.data), b.Type, b.Last, b.WindowSize)
***REMOVED***

func newBlockDec(lowMem bool) *blockDec ***REMOVED***
	b := blockDec***REMOVED***
		lowMem:  lowMem,
		result:  make(chan decodeOutput, 1),
		input:   make(chan struct***REMOVED******REMOVED***, 1),
		history: make(chan *history, 1),
	***REMOVED***
	go b.startDecoder()
	return &b
***REMOVED***

// reset will reset the block.
// Input must be a start of a block and will be at the end of the block when returned.
func (b *blockDec) reset(br byteBuffer, windowSize uint64) error ***REMOVED***
	b.WindowSize = windowSize
	tmp := br.readSmall(3)
	if tmp == nil ***REMOVED***
		if debug ***REMOVED***
			println("Reading block header:", io.ErrUnexpectedEOF)
		***REMOVED***
		return io.ErrUnexpectedEOF
	***REMOVED***
	bh := uint32(tmp[0]) | (uint32(tmp[1]) << 8) | (uint32(tmp[2]) << 16)
	b.Last = bh&1 != 0
	b.Type = blockType((bh >> 1) & 3)
	// find size.
	cSize := int(bh >> 3)
	switch b.Type ***REMOVED***
	case blockTypeReserved:
		return ErrReservedBlockType
	case blockTypeRLE:
		b.RLESize = uint32(cSize)
		cSize = 1
	case blockTypeCompressed:
		if debug ***REMOVED***
			println("Data size on stream:", cSize)
		***REMOVED***
		b.RLESize = 0
		if cSize > maxCompressedBlockSize || uint64(cSize) > b.WindowSize ***REMOVED***
			if debug ***REMOVED***
				printf("compressed block too big: csize:%d block: %+v\n", uint64(cSize), b)
			***REMOVED***
			return ErrCompressedSizeTooBig
		***REMOVED***
	default:
		b.RLESize = 0
	***REMOVED***

	// Read block data.
	if cap(b.data) < cSize ***REMOVED***
		if b.lowMem ***REMOVED***
			b.data = make([]byte, 0, cSize)
		***REMOVED*** else ***REMOVED***
			b.data = make([]byte, 0, maxBlockSize)
		***REMOVED***
	***REMOVED***
	if cap(b.dst) <= maxBlockSize ***REMOVED***
		b.dst = make([]byte, 0, maxBlockSize+1)
	***REMOVED***
	var err error
	b.data, err = br.readBig(cSize, b.data[:0])
	if err != nil ***REMOVED***
		if debug ***REMOVED***
			println("Reading block:", err)
		***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// sendEOF will make the decoder send EOF on this frame.
func (b *blockDec) sendErr(err error) ***REMOVED***
	b.Last = true
	b.Type = blockTypeReserved
	b.err = err
	b.input <- struct***REMOVED******REMOVED******REMOVED******REMOVED***
***REMOVED***

// Close will release resources.
// Closed blockDec cannot be reset.
func (b *blockDec) Close() ***REMOVED***
	close(b.input)
	close(b.history)
	close(b.result)
***REMOVED***

// decodeAsync will prepare decoding the block when it receives input.
// This will separate output and history.
func (b *blockDec) startDecoder() ***REMOVED***
	for range b.input ***REMOVED***
		//println("blockDec: Got block input")
		switch b.Type ***REMOVED***
		case blockTypeRLE:
			if cap(b.dst) < int(b.RLESize) ***REMOVED***
				if b.lowMem ***REMOVED***
					b.dst = make([]byte, b.RLESize)
				***REMOVED*** else ***REMOVED***
					b.dst = make([]byte, maxBlockSize)
				***REMOVED***
			***REMOVED***
			o := decodeOutput***REMOVED***
				d:   b,
				b:   b.dst[:b.RLESize],
				err: nil,
			***REMOVED***
			v := b.data[0]
			for i := range o.b ***REMOVED***
				o.b[i] = v
			***REMOVED***
			hist := <-b.history
			hist.append(o.b)
			b.result <- o
		case blockTypeRaw:
			o := decodeOutput***REMOVED***
				d:   b,
				b:   b.data,
				err: nil,
			***REMOVED***
			hist := <-b.history
			hist.append(o.b)
			b.result <- o
		case blockTypeCompressed:
			b.dst = b.dst[:0]
			err := b.decodeCompressed(nil)
			o := decodeOutput***REMOVED***
				d:   b,
				b:   b.dst,
				err: err,
			***REMOVED***
			if debug ***REMOVED***
				println("Decompressed to", len(b.dst), "bytes, error:", err)
			***REMOVED***
			b.result <- o
		case blockTypeReserved:
			// Used for returning errors.
			<-b.history
			b.result <- decodeOutput***REMOVED***
				d:   b,
				b:   nil,
				err: b.err,
			***REMOVED***
		default:
			panic("Invalid block type")
		***REMOVED***
		if debug ***REMOVED***
			println("blockDec: Finished block")
		***REMOVED***
	***REMOVED***
***REMOVED***

// decodeAsync will prepare decoding the block when it receives the history.
// If history is provided, it will not fetch it from the channel.
func (b *blockDec) decodeBuf(hist *history) error ***REMOVED***
	switch b.Type ***REMOVED***
	case blockTypeRLE:
		if cap(b.dst) < int(b.RLESize) ***REMOVED***
			if b.lowMem ***REMOVED***
				b.dst = make([]byte, b.RLESize)
			***REMOVED*** else ***REMOVED***
				b.dst = make([]byte, maxBlockSize)
			***REMOVED***
		***REMOVED***
		b.dst = b.dst[:b.RLESize]
		v := b.data[0]
		for i := range b.dst ***REMOVED***
			b.dst[i] = v
		***REMOVED***
		hist.appendKeep(b.dst)
		return nil
	case blockTypeRaw:
		hist.appendKeep(b.data)
		return nil
	case blockTypeCompressed:
		saved := b.dst
		b.dst = hist.b
		hist.b = nil
		err := b.decodeCompressed(hist)
		if debug ***REMOVED***
			println("Decompressed to total", len(b.dst), "bytes, error:", err)
		***REMOVED***
		hist.b = b.dst
		b.dst = saved
		return err
	case blockTypeReserved:
		// Used for returning errors.
		return b.err
	default:
		panic("Invalid block type")
	***REMOVED***
***REMOVED***

// decodeCompressed will start decompressing a block.
// If no history is supplied the decoder will decodeAsync as much as possible
// before fetching from blockDec.history
func (b *blockDec) decodeCompressed(hist *history) error ***REMOVED***
	in := b.data
	delayedHistory := hist == nil

	if delayedHistory ***REMOVED***
		// We must always grab history.
		defer func() ***REMOVED***
			if hist == nil ***REMOVED***
				<-b.history
			***REMOVED***
		***REMOVED***()
	***REMOVED***
	// There must be at least one byte for Literals_Block_Type and one for Sequences_Section_Header
	if len(in) < 2 ***REMOVED***
		return ErrBlockTooSmall
	***REMOVED***
	litType := literalsBlockType(in[0] & 3)
	var litRegenSize int
	var litCompSize int
	sizeFormat := (in[0] >> 2) & 3
	var fourStreams bool
	switch litType ***REMOVED***
	case literalsBlockRaw, literalsBlockRLE:
		switch sizeFormat ***REMOVED***
		case 0, 2:
			// Regenerated_Size uses 5 bits (0-31). Literals_Section_Header uses 1 byte.
			litRegenSize = int(in[0] >> 3)
			in = in[1:]
		case 1:
			// Regenerated_Size uses 12 bits (0-4095). Literals_Section_Header uses 2 bytes.
			litRegenSize = int(in[0]>>4) + (int(in[1]) << 4)
			in = in[2:]
		case 3:
			//  Regenerated_Size uses 20 bits (0-1048575). Literals_Section_Header uses 3 bytes.
			if len(in) < 3 ***REMOVED***
				println("too small: litType:", litType, " sizeFormat", sizeFormat, len(in))
				return ErrBlockTooSmall
			***REMOVED***
			litRegenSize = int(in[0]>>4) + (int(in[1]) << 4) + (int(in[2]) << 12)
			in = in[3:]
		***REMOVED***
	case literalsBlockCompressed, literalsBlockTreeless:
		switch sizeFormat ***REMOVED***
		case 0, 1:
			// Both Regenerated_Size and Compressed_Size use 10 bits (0-1023).
			if len(in) < 3 ***REMOVED***
				println("too small: litType:", litType, " sizeFormat", sizeFormat, len(in))
				return ErrBlockTooSmall
			***REMOVED***
			n := uint64(in[0]>>4) + (uint64(in[1]) << 4) + (uint64(in[2]) << 12)
			litRegenSize = int(n & 1023)
			litCompSize = int(n >> 10)
			fourStreams = sizeFormat == 1
			in = in[3:]
		case 2:
			fourStreams = true
			if len(in) < 4 ***REMOVED***
				println("too small: litType:", litType, " sizeFormat", sizeFormat, len(in))
				return ErrBlockTooSmall
			***REMOVED***
			n := uint64(in[0]>>4) + (uint64(in[1]) << 4) + (uint64(in[2]) << 12) + (uint64(in[3]) << 20)
			litRegenSize = int(n & 16383)
			litCompSize = int(n >> 14)
			in = in[4:]
		case 3:
			fourStreams = true
			if len(in) < 5 ***REMOVED***
				println("too small: litType:", litType, " sizeFormat", sizeFormat, len(in))
				return ErrBlockTooSmall
			***REMOVED***
			n := uint64(in[0]>>4) + (uint64(in[1]) << 4) + (uint64(in[2]) << 12) + (uint64(in[3]) << 20) + (uint64(in[4]) << 28)
			litRegenSize = int(n & 262143)
			litCompSize = int(n >> 18)
			in = in[5:]
		***REMOVED***
	***REMOVED***
	if debug ***REMOVED***
		println("literals type:", litType, "litRegenSize:", litRegenSize, "litCompSize", litCompSize)
	***REMOVED***
	var literals []byte
	var huff *huff0.Scratch
	switch litType ***REMOVED***
	case literalsBlockRaw:
		if len(in) < litRegenSize ***REMOVED***
			println("too small: litType:", litType, " sizeFormat", sizeFormat, "remain:", len(in), "want:", litRegenSize)
			return ErrBlockTooSmall
		***REMOVED***
		literals = in[:litRegenSize]
		in = in[litRegenSize:]
		//printf("Found %d uncompressed literals\n", litRegenSize)
	case literalsBlockRLE:
		if len(in) < 1 ***REMOVED***
			println("too small: litType:", litType, " sizeFormat", sizeFormat, "remain:", len(in), "want:", 1)
			return ErrBlockTooSmall
		***REMOVED***
		if cap(b.literalBuf) < litRegenSize ***REMOVED***
			if b.lowMem ***REMOVED***
				b.literalBuf = make([]byte, litRegenSize)
			***REMOVED*** else ***REMOVED***
				if litRegenSize > maxCompressedLiteralSize ***REMOVED***
					// Exceptional
					b.literalBuf = make([]byte, litRegenSize)
				***REMOVED*** else ***REMOVED***
					b.literalBuf = make([]byte, litRegenSize, maxCompressedLiteralSize)

				***REMOVED***
			***REMOVED***
		***REMOVED***
		literals = b.literalBuf[:litRegenSize]
		v := in[0]
		for i := range literals ***REMOVED***
			literals[i] = v
		***REMOVED***
		in = in[1:]
		if debug ***REMOVED***
			printf("Found %d RLE compressed literals\n", litRegenSize)
		***REMOVED***
	case literalsBlockTreeless:
		if len(in) < litCompSize ***REMOVED***
			println("too small: litType:", litType, " sizeFormat", sizeFormat, "remain:", len(in), "want:", litCompSize)
			return ErrBlockTooSmall
		***REMOVED***
		// Store compressed literals, so we defer decoding until we get history.
		literals = in[:litCompSize]
		in = in[litCompSize:]
		if debug ***REMOVED***
			printf("Found %d compressed literals\n", litCompSize)
		***REMOVED***
	case literalsBlockCompressed:
		if len(in) < litCompSize ***REMOVED***
			println("too small: litType:", litType, " sizeFormat", sizeFormat, "remain:", len(in), "want:", litCompSize)
			return ErrBlockTooSmall
		***REMOVED***
		literals = in[:litCompSize]
		in = in[litCompSize:]

		huff = huffDecoderPool.Get().(*huff0.Scratch)
		var err error
		// Ensure we have space to store it.
		if cap(b.literalBuf) < litRegenSize ***REMOVED***
			if b.lowMem ***REMOVED***
				b.literalBuf = make([]byte, 0, litRegenSize)
			***REMOVED*** else ***REMOVED***
				b.literalBuf = make([]byte, 0, maxCompressedLiteralSize)
			***REMOVED***
		***REMOVED***
		if huff == nil ***REMOVED***
			huff = &huff0.Scratch***REMOVED******REMOVED***
		***REMOVED***
		huff.Out = b.literalBuf[:0]
		huff, literals, err = huff0.ReadTable(literals, huff)
		if err != nil ***REMOVED***
			println("reading huffman table:", err)
			return err
		***REMOVED***
		// Use our out buffer.
		huff.Out = b.literalBuf[:0]
		if fourStreams ***REMOVED***
			literals, err = huff.Decompress4X(literals, litRegenSize)
		***REMOVED*** else ***REMOVED***
			literals, err = huff.Decompress1X(literals)
		***REMOVED***
		if err != nil ***REMOVED***
			println("decoding compressed literals:", err)
			return err
		***REMOVED***
		// Make sure we don't leak our literals buffer
		huff.Out = nil
		if len(literals) != litRegenSize ***REMOVED***
			return fmt.Errorf("literal output size mismatch want %d, got %d", litRegenSize, len(literals))
		***REMOVED***
		if debug ***REMOVED***
			printf("Decompressed %d literals into %d bytes\n", litCompSize, litRegenSize)
		***REMOVED***
	***REMOVED***

	// Decode Sequences
	// https://github.com/facebook/zstd/blob/dev/doc/zstd_compression_format.md#sequences-section
	if len(in) < 1 ***REMOVED***
		return ErrBlockTooSmall
	***REMOVED***
	seqHeader := in[0]
	nSeqs := 0
	switch ***REMOVED***
	case seqHeader == 0:
		in = in[1:]
	case seqHeader < 128:
		nSeqs = int(seqHeader)
		in = in[1:]
	case seqHeader < 255:
		if len(in) < 2 ***REMOVED***
			return ErrBlockTooSmall
		***REMOVED***
		nSeqs = int(seqHeader-128)<<8 | int(in[1])
		in = in[2:]
	case seqHeader == 255:
		if len(in) < 3 ***REMOVED***
			return ErrBlockTooSmall
		***REMOVED***
		nSeqs = 0x7f00 + int(in[1]) + (int(in[2]) << 8)
		in = in[3:]
	***REMOVED***
	// Allocate sequences
	if cap(b.sequenceBuf) < nSeqs ***REMOVED***
		if b.lowMem ***REMOVED***
			b.sequenceBuf = make([]seq, nSeqs)
		***REMOVED*** else ***REMOVED***
			// Allocate max
			b.sequenceBuf = make([]seq, nSeqs, maxSequences)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// Reuse buffer
		b.sequenceBuf = b.sequenceBuf[:nSeqs]
	***REMOVED***
	var seqs = &sequenceDecs***REMOVED******REMOVED***
	if nSeqs > 0 ***REMOVED***
		if len(in) < 1 ***REMOVED***
			return ErrBlockTooSmall
		***REMOVED***
		br := byteReader***REMOVED***b: in, off: 0***REMOVED***
		compMode := br.Uint8()
		br.advance(1)
		if debug ***REMOVED***
			printf("Compression modes: 0b%b", compMode)
		***REMOVED***
		for i := uint(0); i < 3; i++ ***REMOVED***
			mode := seqCompMode((compMode >> (6 - i*2)) & 3)
			if debug ***REMOVED***
				println("Table", tableIndex(i), "is", mode)
			***REMOVED***
			var seq *sequenceDec
			switch tableIndex(i) ***REMOVED***
			case tableLiteralLengths:
				seq = &seqs.litLengths
			case tableOffsets:
				seq = &seqs.offsets
			case tableMatchLengths:
				seq = &seqs.matchLengths
			default:
				panic("unknown table")
			***REMOVED***
			switch mode ***REMOVED***
			case compModePredefined:
				seq.fse = &fsePredef[i]
			case compModeRLE:
				if br.remain() < 1 ***REMOVED***
					return ErrBlockTooSmall
				***REMOVED***
				v := br.Uint8()
				br.advance(1)
				dec := fseDecoderPool.Get().(*fseDecoder)
				symb, err := decSymbolValue(v, symbolTableX[i])
				if err != nil ***REMOVED***
					printf("RLE Transform table (%v) error: %v", tableIndex(i), err)
					return err
				***REMOVED***
				dec.setRLE(symb)
				seq.fse = dec
				if debug ***REMOVED***
					printf("RLE set to %+v, code: %v", symb, v)
				***REMOVED***
			case compModeFSE:
				println("Reading table for", tableIndex(i))
				dec := fseDecoderPool.Get().(*fseDecoder)
				err := dec.readNCount(&br, uint16(maxTableSymbol[i]))
				if err != nil ***REMOVED***
					println("Read table error:", err)
					return err
				***REMOVED***
				err = dec.transform(symbolTableX[i])
				if err != nil ***REMOVED***
					println("Transform table error:", err)
					return err
				***REMOVED***
				if debug ***REMOVED***
					println("Read table ok", "symbolLen:", dec.symbolLen)
				***REMOVED***
				seq.fse = dec
			case compModeRepeat:
				seq.repeat = true
			***REMOVED***
			if br.overread() ***REMOVED***
				return io.ErrUnexpectedEOF
			***REMOVED***
		***REMOVED***
		in = br.unread()
	***REMOVED***

	// Wait for history.
	// All time spent after this is critical since it is strictly sequential.
	if hist == nil ***REMOVED***
		hist = <-b.history
		if hist.error ***REMOVED***
			return ErrDecoderClosed
		***REMOVED***
	***REMOVED***

	// Decode treeless literal block.
	if litType == literalsBlockTreeless ***REMOVED***
		// TODO: We could send the history early WITHOUT the stream history.
		//   This would allow decoding treeless literials before the byte history is available.
		//   Silencia stats: Treeless 4393, with: 32775, total: 37168, 11% treeless.
		//   So not much obvious gain here.

		if hist.huffTree == nil ***REMOVED***
			return errors.New("literal block was treeless, but no history was defined")
		***REMOVED***
		// Ensure we have space to store it.
		if cap(b.literalBuf) < litRegenSize ***REMOVED***
			if b.lowMem ***REMOVED***
				b.literalBuf = make([]byte, 0, litRegenSize)
			***REMOVED*** else ***REMOVED***
				b.literalBuf = make([]byte, 0, maxCompressedLiteralSize)
			***REMOVED***
		***REMOVED***
		var err error
		// Use our out buffer.
		huff = hist.huffTree
		huff.Out = b.literalBuf[:0]
		if fourStreams ***REMOVED***
			literals, err = huff.Decompress4X(literals, litRegenSize)
		***REMOVED*** else ***REMOVED***
			literals, err = huff.Decompress1X(literals)
		***REMOVED***
		// Make sure we don't leak our literals buffer
		huff.Out = nil
		if err != nil ***REMOVED***
			println("decompressing literals:", err)
			return err
		***REMOVED***
		if len(literals) != litRegenSize ***REMOVED***
			return fmt.Errorf("literal output size mismatch want %d, got %d", litRegenSize, len(literals))
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if hist.huffTree != nil && huff != nil ***REMOVED***
			huffDecoderPool.Put(hist.huffTree)
			hist.huffTree = nil
		***REMOVED***
	***REMOVED***
	if huff != nil ***REMOVED***
		huff.Out = nil
		hist.huffTree = huff
	***REMOVED***
	if debug ***REMOVED***
		println("Final literals:", len(literals), "and", nSeqs, "sequences.")
	***REMOVED***

	if nSeqs == 0 ***REMOVED***
		// Decompressed content is defined entirely as Literals Section content.
		b.dst = append(b.dst, literals...)
		if delayedHistory ***REMOVED***
			hist.append(literals)
		***REMOVED***
		return nil
	***REMOVED***

	seqs, err := seqs.mergeHistory(&hist.decoders)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if debug ***REMOVED***
		println("History merged ok")
	***REMOVED***
	br := &bitReader***REMOVED******REMOVED***
	if err := br.init(in); err != nil ***REMOVED***
		return err
	***REMOVED***

	// TODO: Investigate if sending history without decoders are faster.
	//   This would allow the sequences to be decoded async and only have to construct stream history.
	//   If only recent offsets were not transferred, this would be an obvious win.
	// 	 Also, if first 3 sequences don't reference recent offsets, all sequences can be decoded.

	if err := seqs.initialize(br, hist, literals, b.dst); err != nil ***REMOVED***
		println("initializing sequences:", err)
		return err
	***REMOVED***

	err = seqs.decode(nSeqs, br, hist.b)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if !br.finished() ***REMOVED***
		return fmt.Errorf("%d extra bits on block, should be 0", br.remain())
	***REMOVED***

	err = br.close()
	if err != nil ***REMOVED***
		printf("Closing sequences: %v, %+v\n", err, *br)
	***REMOVED***
	if len(b.data) > maxCompressedBlockSize ***REMOVED***
		return fmt.Errorf("compressed block size too large (%d)", len(b.data))
	***REMOVED***
	// Set output and release references.
	b.dst = seqs.out
	seqs.out, seqs.literals, seqs.hist = nil, nil, nil

	if !delayedHistory ***REMOVED***
		// If we don't have delayed history, no need to update.
		hist.recentOffsets = seqs.prevOffset
		return nil
	***REMOVED***
	if b.Last ***REMOVED***
		// if last block we don't care about history.
		println("Last block, no history returned")
		hist.b = hist.b[:0]
		return nil
	***REMOVED***
	hist.append(b.dst)
	hist.recentOffsets = seqs.prevOffset
	if debug ***REMOVED***
		println("Finished block with literals:", len(literals), "and", nSeqs, "sequences.")
	***REMOVED***

	return nil
***REMOVED***
