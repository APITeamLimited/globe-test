// Copyright 2019+ Klaus Post. All rights reserved.
// License information can be found in the LICENSE file.
// Based on work by Yann Collet, released under BSD License.

package zstd

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/klauspost/compress/huff0"
	"github.com/klauspost/compress/zstd/internal/xxhash"
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

	compressedBlockOverAlloc    = 16
	maxCompressedBlockSizeAlloc = 128<<10 + compressedBlockOverAlloc

	// Maximum possible block size (all Raw+Uncompressed).
	maxBlockSize = (1 << 21) - 1

	maxMatchLen  = 131074
	maxSequences = 0x7f00 + 0xffff

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
	data        []byte
	dataStorage []byte

	// Destination of the decoded data.
	dst []byte

	// Buffer for literals data.
	literalBuf []byte

	// Window size of the block.
	WindowSize uint64

	err error

	// Check against this crc
	checkCRC []byte

	// Frame to use for singlethreaded decoding.
	// Should not be used by the decoder itself since parent may be another frame.
	localFrame *frameDec

	sequence []seqVals

	async struct ***REMOVED***
		newHist  *history
		literals []byte
		seqData  []byte
		seqSize  int // Size of uncompressed sequences
		fcs      uint64
	***REMOVED***

	// Block is RLE, this is the size.
	RLESize uint32

	Type blockType

	// Is this the last block of a frame?
	Last bool

	// Use less memory
	lowMem bool
***REMOVED***

func (b *blockDec) String() string ***REMOVED***
	if b == nil ***REMOVED***
		return "<nil>"
	***REMOVED***
	return fmt.Sprintf("Steam Size: %d, Type: %v, Last: %t, Window: %d", len(b.data), b.Type, b.Last, b.WindowSize)
***REMOVED***

func newBlockDec(lowMem bool) *blockDec ***REMOVED***
	b := blockDec***REMOVED***
		lowMem: lowMem,
	***REMOVED***
	return &b
***REMOVED***

// reset will reset the block.
// Input must be a start of a block and will be at the end of the block when returned.
func (b *blockDec) reset(br byteBuffer, windowSize uint64) error ***REMOVED***
	b.WindowSize = windowSize
	tmp, err := br.readSmall(3)
	if err != nil ***REMOVED***
		println("Reading block header:", err)
		return err
	***REMOVED***
	bh := uint32(tmp[0]) | (uint32(tmp[1]) << 8) | (uint32(tmp[2]) << 16)
	b.Last = bh&1 != 0
	b.Type = blockType((bh >> 1) & 3)
	// find size.
	cSize := int(bh >> 3)
	maxSize := maxCompressedBlockSizeAlloc
	switch b.Type ***REMOVED***
	case blockTypeReserved:
		return ErrReservedBlockType
	case blockTypeRLE:
		if cSize > maxCompressedBlockSize || cSize > int(b.WindowSize) ***REMOVED***
			if debugDecoder ***REMOVED***
				printf("rle block too big: csize:%d block: %+v\n", uint64(cSize), b)
			***REMOVED***
			return ErrWindowSizeExceeded
		***REMOVED***
		b.RLESize = uint32(cSize)
		if b.lowMem ***REMOVED***
			maxSize = cSize
		***REMOVED***
		cSize = 1
	case blockTypeCompressed:
		if debugDecoder ***REMOVED***
			println("Data size on stream:", cSize)
		***REMOVED***
		b.RLESize = 0
		maxSize = maxCompressedBlockSizeAlloc
		if windowSize < maxCompressedBlockSize && b.lowMem ***REMOVED***
			maxSize = int(windowSize) + compressedBlockOverAlloc
		***REMOVED***
		if cSize > maxCompressedBlockSize || uint64(cSize) > b.WindowSize ***REMOVED***
			if debugDecoder ***REMOVED***
				printf("compressed block too big: csize:%d block: %+v\n", uint64(cSize), b)
			***REMOVED***
			return ErrCompressedSizeTooBig
		***REMOVED***
		// Empty compressed blocks must at least be 2 bytes
		// for Literals_Block_Type and one for Sequences_Section_Header.
		if cSize < 2 ***REMOVED***
			return ErrBlockTooSmall
		***REMOVED***
	case blockTypeRaw:
		if cSize > maxCompressedBlockSize || cSize > int(b.WindowSize) ***REMOVED***
			if debugDecoder ***REMOVED***
				printf("rle block too big: csize:%d block: %+v\n", uint64(cSize), b)
			***REMOVED***
			return ErrWindowSizeExceeded
		***REMOVED***

		b.RLESize = 0
		// We do not need a destination for raw blocks.
		maxSize = -1
	default:
		panic("Invalid block type")
	***REMOVED***

	// Read block data.
	if cap(b.dataStorage) < cSize ***REMOVED***
		if b.lowMem || cSize > maxCompressedBlockSize ***REMOVED***
			b.dataStorage = make([]byte, 0, cSize+compressedBlockOverAlloc)
		***REMOVED*** else ***REMOVED***
			b.dataStorage = make([]byte, 0, maxCompressedBlockSizeAlloc)
		***REMOVED***
	***REMOVED***
	if cap(b.dst) <= maxSize ***REMOVED***
		b.dst = make([]byte, 0, maxSize+1)
	***REMOVED***
	b.data, err = br.readBig(cSize, b.dataStorage)
	if err != nil ***REMOVED***
		if debugDecoder ***REMOVED***
			println("Reading block:", err, "(", cSize, ")", len(b.data))
			printf("%T", br)
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
***REMOVED***

// Close will release resources.
// Closed blockDec cannot be reset.
func (b *blockDec) Close() ***REMOVED***
***REMOVED***

// decodeBuf
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
		// Append directly to history
		if hist.ignoreBuffer == 0 ***REMOVED***
			b.dst = hist.b
			hist.b = nil
		***REMOVED*** else ***REMOVED***
			b.dst = b.dst[:0]
		***REMOVED***
		err := b.decodeCompressed(hist)
		if debugDecoder ***REMOVED***
			println("Decompressed to total", len(b.dst), "bytes, hash:", xxhash.Sum64(b.dst), "error:", err)
		***REMOVED***
		if hist.ignoreBuffer == 0 ***REMOVED***
			hist.b = b.dst
			b.dst = saved
		***REMOVED*** else ***REMOVED***
			hist.appendKeep(b.dst)
		***REMOVED***
		return err
	case blockTypeReserved:
		// Used for returning errors.
		return b.err
	default:
		panic("Invalid block type")
	***REMOVED***
***REMOVED***

func (b *blockDec) decodeLiterals(in []byte, hist *history) (remain []byte, err error) ***REMOVED***
	// There must be at least one byte for Literals_Block_Type and one for Sequences_Section_Header
	if len(in) < 2 ***REMOVED***
		return in, ErrBlockTooSmall
	***REMOVED***

	litType := literalsBlockType(in[0] & 3)
	var litRegenSize int
	var litCompSize int
	sizeFormat := (in[0] >> 2) & 3
	var fourStreams bool
	var literals []byte
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
				return in, ErrBlockTooSmall
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
				return in, ErrBlockTooSmall
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
				return in, ErrBlockTooSmall
			***REMOVED***
			n := uint64(in[0]>>4) + (uint64(in[1]) << 4) + (uint64(in[2]) << 12) + (uint64(in[3]) << 20)
			litRegenSize = int(n & 16383)
			litCompSize = int(n >> 14)
			in = in[4:]
		case 3:
			fourStreams = true
			if len(in) < 5 ***REMOVED***
				println("too small: litType:", litType, " sizeFormat", sizeFormat, len(in))
				return in, ErrBlockTooSmall
			***REMOVED***
			n := uint64(in[0]>>4) + (uint64(in[1]) << 4) + (uint64(in[2]) << 12) + (uint64(in[3]) << 20) + (uint64(in[4]) << 28)
			litRegenSize = int(n & 262143)
			litCompSize = int(n >> 18)
			in = in[5:]
		***REMOVED***
	***REMOVED***
	if debugDecoder ***REMOVED***
		println("literals type:", litType, "litRegenSize:", litRegenSize, "litCompSize:", litCompSize, "sizeFormat:", sizeFormat, "4X:", fourStreams)
	***REMOVED***
	if litRegenSize > int(b.WindowSize) || litRegenSize > maxCompressedBlockSize ***REMOVED***
		return in, ErrWindowSizeExceeded
	***REMOVED***

	switch litType ***REMOVED***
	case literalsBlockRaw:
		if len(in) < litRegenSize ***REMOVED***
			println("too small: litType:", litType, " sizeFormat", sizeFormat, "remain:", len(in), "want:", litRegenSize)
			return in, ErrBlockTooSmall
		***REMOVED***
		literals = in[:litRegenSize]
		in = in[litRegenSize:]
		//printf("Found %d uncompressed literals\n", litRegenSize)
	case literalsBlockRLE:
		if len(in) < 1 ***REMOVED***
			println("too small: litType:", litType, " sizeFormat", sizeFormat, "remain:", len(in), "want:", 1)
			return in, ErrBlockTooSmall
		***REMOVED***
		if cap(b.literalBuf) < litRegenSize ***REMOVED***
			if b.lowMem ***REMOVED***
				b.literalBuf = make([]byte, litRegenSize, litRegenSize+compressedBlockOverAlloc)
			***REMOVED*** else ***REMOVED***
				b.literalBuf = make([]byte, litRegenSize, maxCompressedBlockSize+compressedBlockOverAlloc)
			***REMOVED***
		***REMOVED***
		literals = b.literalBuf[:litRegenSize]
		v := in[0]
		for i := range literals ***REMOVED***
			literals[i] = v
		***REMOVED***
		in = in[1:]
		if debugDecoder ***REMOVED***
			printf("Found %d RLE compressed literals\n", litRegenSize)
		***REMOVED***
	case literalsBlockTreeless:
		if len(in) < litCompSize ***REMOVED***
			println("too small: litType:", litType, " sizeFormat", sizeFormat, "remain:", len(in), "want:", litCompSize)
			return in, ErrBlockTooSmall
		***REMOVED***
		// Store compressed literals, so we defer decoding until we get history.
		literals = in[:litCompSize]
		in = in[litCompSize:]
		if debugDecoder ***REMOVED***
			printf("Found %d compressed literals\n", litCompSize)
		***REMOVED***
		huff := hist.huffTree
		if huff == nil ***REMOVED***
			return in, errors.New("literal block was treeless, but no history was defined")
		***REMOVED***
		// Ensure we have space to store it.
		if cap(b.literalBuf) < litRegenSize ***REMOVED***
			if b.lowMem ***REMOVED***
				b.literalBuf = make([]byte, 0, litRegenSize+compressedBlockOverAlloc)
			***REMOVED*** else ***REMOVED***
				b.literalBuf = make([]byte, 0, maxCompressedBlockSize+compressedBlockOverAlloc)
			***REMOVED***
		***REMOVED***
		var err error
		// Use our out buffer.
		huff.MaxDecodedSize = litRegenSize
		if fourStreams ***REMOVED***
			literals, err = huff.Decoder().Decompress4X(b.literalBuf[:0:litRegenSize], literals)
		***REMOVED*** else ***REMOVED***
			literals, err = huff.Decoder().Decompress1X(b.literalBuf[:0:litRegenSize], literals)
		***REMOVED***
		// Make sure we don't leak our literals buffer
		if err != nil ***REMOVED***
			println("decompressing literals:", err)
			return in, err
		***REMOVED***
		if len(literals) != litRegenSize ***REMOVED***
			return in, fmt.Errorf("literal output size mismatch want %d, got %d", litRegenSize, len(literals))
		***REMOVED***

	case literalsBlockCompressed:
		if len(in) < litCompSize ***REMOVED***
			println("too small: litType:", litType, " sizeFormat", sizeFormat, "remain:", len(in), "want:", litCompSize)
			return in, ErrBlockTooSmall
		***REMOVED***
		literals = in[:litCompSize]
		in = in[litCompSize:]
		// Ensure we have space to store it.
		if cap(b.literalBuf) < litRegenSize ***REMOVED***
			if b.lowMem ***REMOVED***
				b.literalBuf = make([]byte, 0, litRegenSize+compressedBlockOverAlloc)
			***REMOVED*** else ***REMOVED***
				b.literalBuf = make([]byte, 0, maxCompressedBlockSize+compressedBlockOverAlloc)
			***REMOVED***
		***REMOVED***
		huff := hist.huffTree
		if huff == nil || (hist.dict != nil && huff == hist.dict.litEnc) ***REMOVED***
			huff = huffDecoderPool.Get().(*huff0.Scratch)
			if huff == nil ***REMOVED***
				huff = &huff0.Scratch***REMOVED******REMOVED***
			***REMOVED***
		***REMOVED***
		var err error
		huff, literals, err = huff0.ReadTable(literals, huff)
		if err != nil ***REMOVED***
			println("reading huffman table:", err)
			return in, err
		***REMOVED***
		hist.huffTree = huff
		huff.MaxDecodedSize = litRegenSize
		// Use our out buffer.
		if fourStreams ***REMOVED***
			literals, err = huff.Decoder().Decompress4X(b.literalBuf[:0:litRegenSize], literals)
		***REMOVED*** else ***REMOVED***
			literals, err = huff.Decoder().Decompress1X(b.literalBuf[:0:litRegenSize], literals)
		***REMOVED***
		if err != nil ***REMOVED***
			println("decoding compressed literals:", err)
			return in, err
		***REMOVED***
		// Make sure we don't leak our literals buffer
		if len(literals) != litRegenSize ***REMOVED***
			return in, fmt.Errorf("literal output size mismatch want %d, got %d", litRegenSize, len(literals))
		***REMOVED***
		// Re-cap to get extra size.
		literals = b.literalBuf[:len(literals)]
		if debugDecoder ***REMOVED***
			printf("Decompressed %d literals into %d bytes\n", litCompSize, litRegenSize)
		***REMOVED***
	***REMOVED***
	hist.decoders.literals = literals
	return in, nil
***REMOVED***

// decodeCompressed will start decompressing a block.
func (b *blockDec) decodeCompressed(hist *history) error ***REMOVED***
	in := b.data
	in, err := b.decodeLiterals(in, hist)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	err = b.prepareSequences(in, hist)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if hist.decoders.nSeqs == 0 ***REMOVED***
		b.dst = append(b.dst, hist.decoders.literals...)
		return nil
	***REMOVED***
	before := len(hist.decoders.out)
	err = hist.decoders.decodeSync(hist.b[hist.ignoreBuffer:])
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if hist.decoders.maxSyncLen > 0 ***REMOVED***
		hist.decoders.maxSyncLen += uint64(before)
		hist.decoders.maxSyncLen -= uint64(len(hist.decoders.out))
	***REMOVED***
	b.dst = hist.decoders.out
	hist.recentOffsets = hist.decoders.prevOffset
	return nil
***REMOVED***

func (b *blockDec) prepareSequences(in []byte, hist *history) (err error) ***REMOVED***
	if debugDecoder ***REMOVED***
		printf("prepareSequences: %d byte(s) input\n", len(in))
	***REMOVED***
	// Decode Sequences
	// https://github.com/facebook/zstd/blob/dev/doc/zstd_compression_format.md#sequences-section
	if len(in) < 1 ***REMOVED***
		return ErrBlockTooSmall
	***REMOVED***
	var nSeqs int
	seqHeader := in[0]
	switch ***REMOVED***
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
	if nSeqs == 0 && len(in) != 0 ***REMOVED***
		// When no sequences, there should not be any more data...
		if debugDecoder ***REMOVED***
			printf("prepareSequences: 0 sequences, but %d byte(s) left on stream\n", len(in))
		***REMOVED***
		return ErrUnexpectedBlockSize
	***REMOVED***

	var seqs = &hist.decoders
	seqs.nSeqs = nSeqs
	if nSeqs > 0 ***REMOVED***
		if len(in) < 1 ***REMOVED***
			return ErrBlockTooSmall
		***REMOVED***
		br := byteReader***REMOVED***b: in, off: 0***REMOVED***
		compMode := br.Uint8()
		br.advance(1)
		if debugDecoder ***REMOVED***
			printf("Compression modes: 0b%b", compMode)
		***REMOVED***
		for i := uint(0); i < 3; i++ ***REMOVED***
			mode := seqCompMode((compMode >> (6 - i*2)) & 3)
			if debugDecoder ***REMOVED***
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
				if seq.fse != nil && !seq.fse.preDefined ***REMOVED***
					fseDecoderPool.Put(seq.fse)
				***REMOVED***
				seq.fse = &fsePredef[i]
			case compModeRLE:
				if br.remain() < 1 ***REMOVED***
					return ErrBlockTooSmall
				***REMOVED***
				v := br.Uint8()
				br.advance(1)
				if seq.fse == nil || seq.fse.preDefined ***REMOVED***
					seq.fse = fseDecoderPool.Get().(*fseDecoder)
				***REMOVED***
				symb, err := decSymbolValue(v, symbolTableX[i])
				if err != nil ***REMOVED***
					printf("RLE Transform table (%v) error: %v", tableIndex(i), err)
					return err
				***REMOVED***
				seq.fse.setRLE(symb)
				if debugDecoder ***REMOVED***
					printf("RLE set to %+v, code: %v", symb, v)
				***REMOVED***
			case compModeFSE:
				println("Reading table for", tableIndex(i))
				if seq.fse == nil || seq.fse.preDefined ***REMOVED***
					seq.fse = fseDecoderPool.Get().(*fseDecoder)
				***REMOVED***
				err := seq.fse.readNCount(&br, uint16(maxTableSymbol[i]))
				if err != nil ***REMOVED***
					println("Read table error:", err)
					return err
				***REMOVED***
				err = seq.fse.transform(symbolTableX[i])
				if err != nil ***REMOVED***
					println("Transform table error:", err)
					return err
				***REMOVED***
				if debugDecoder ***REMOVED***
					println("Read table ok", "symbolLen:", seq.fse.symbolLen)
				***REMOVED***
			case compModeRepeat:
				seq.repeat = true
			***REMOVED***
			if br.overread() ***REMOVED***
				return io.ErrUnexpectedEOF
			***REMOVED***
		***REMOVED***
		in = br.unread()
	***REMOVED***
	if debugDecoder ***REMOVED***
		println("Literals:", len(seqs.literals), "hash:", xxhash.Sum64(seqs.literals), "and", seqs.nSeqs, "sequences.")
	***REMOVED***

	if nSeqs == 0 ***REMOVED***
		if len(b.sequence) > 0 ***REMOVED***
			b.sequence = b.sequence[:0]
		***REMOVED***
		return nil
	***REMOVED***
	br := seqs.br
	if br == nil ***REMOVED***
		br = &bitReader***REMOVED******REMOVED***
	***REMOVED***
	if err := br.init(in); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := seqs.initialize(br, hist, b.dst); err != nil ***REMOVED***
		println("initializing sequences:", err)
		return err
	***REMOVED***
	// Extract blocks...
	if false && hist.dict == nil ***REMOVED***
		fatalErr := func(err error) ***REMOVED***
			if err != nil ***REMOVED***
				panic(err)
			***REMOVED***
		***REMOVED***
		fn := fmt.Sprintf("n-%d-lits-%d-prev-%d-%d-%d-win-%d.blk", hist.decoders.nSeqs, len(hist.decoders.literals), hist.recentOffsets[0], hist.recentOffsets[1], hist.recentOffsets[2], hist.windowSize)
		var buf bytes.Buffer
		fatalErr(binary.Write(&buf, binary.LittleEndian, hist.decoders.litLengths.fse))
		fatalErr(binary.Write(&buf, binary.LittleEndian, hist.decoders.matchLengths.fse))
		fatalErr(binary.Write(&buf, binary.LittleEndian, hist.decoders.offsets.fse))
		buf.Write(in)
		ioutil.WriteFile(filepath.Join("testdata", "seqs", fn), buf.Bytes(), os.ModePerm)
	***REMOVED***

	return nil
***REMOVED***

func (b *blockDec) decodeSequences(hist *history) error ***REMOVED***
	if cap(b.sequence) < hist.decoders.nSeqs ***REMOVED***
		if b.lowMem ***REMOVED***
			b.sequence = make([]seqVals, 0, hist.decoders.nSeqs)
		***REMOVED*** else ***REMOVED***
			b.sequence = make([]seqVals, 0, 0x7F00+0xffff)
		***REMOVED***
	***REMOVED***
	b.sequence = b.sequence[:hist.decoders.nSeqs]
	if hist.decoders.nSeqs == 0 ***REMOVED***
		hist.decoders.seqSize = len(hist.decoders.literals)
		return nil
	***REMOVED***
	hist.decoders.windowSize = hist.windowSize
	hist.decoders.prevOffset = hist.recentOffsets

	err := hist.decoders.decode(b.sequence)
	hist.recentOffsets = hist.decoders.prevOffset
	return err
***REMOVED***

func (b *blockDec) executeSequences(hist *history) error ***REMOVED***
	hbytes := hist.b
	if len(hbytes) > hist.windowSize ***REMOVED***
		hbytes = hbytes[len(hbytes)-hist.windowSize:]
		// We do not need history anymore.
		if hist.dict != nil ***REMOVED***
			hist.dict.content = nil
		***REMOVED***
	***REMOVED***
	hist.decoders.windowSize = hist.windowSize
	hist.decoders.out = b.dst[:0]
	err := hist.decoders.execute(b.sequence, hbytes)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return b.updateHistory(hist)
***REMOVED***

func (b *blockDec) updateHistory(hist *history) error ***REMOVED***
	if len(b.data) > maxCompressedBlockSize ***REMOVED***
		return fmt.Errorf("compressed block size too large (%d)", len(b.data))
	***REMOVED***
	// Set output and release references.
	b.dst = hist.decoders.out
	hist.recentOffsets = hist.decoders.prevOffset

	if b.Last ***REMOVED***
		// if last block we don't care about history.
		println("Last block, no history returned")
		hist.b = hist.b[:0]
		return nil
	***REMOVED*** else ***REMOVED***
		hist.append(b.dst)
		if debugDecoder ***REMOVED***
			println("Finished block with ", len(b.sequence), "sequences. Added", len(b.dst), "to history, now length", len(hist.b))
		***REMOVED***
	***REMOVED***
	hist.decoders.out, hist.decoders.literals = nil, nil

	return nil
***REMOVED***
