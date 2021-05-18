// Copyright 2019+ Klaus Post. All rights reserved.
// License information can be found in the LICENSE file.
// Based on work by Yann Collet, released under BSD License.

package zstd

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"

	"github.com/golang/snappy"
	"github.com/klauspost/compress/huff0"
)

const (
	snappyTagLiteral = 0x00
	snappyTagCopy1   = 0x01
	snappyTagCopy2   = 0x02
	snappyTagCopy4   = 0x03
)

const (
	snappyChecksumSize = 4
	snappyMagicBody    = "sNaPpY"

	// snappyMaxBlockSize is the maximum size of the input to encodeBlock. It is not
	// part of the wire format per se, but some parts of the encoder assume
	// that an offset fits into a uint16.
	//
	// Also, for the framing format (Writer type instead of Encode function),
	// https://github.com/google/snappy/blob/master/framing_format.txt says
	// that "the uncompressed data in a chunk must be no longer than 65536
	// bytes".
	snappyMaxBlockSize = 65536

	// snappyMaxEncodedLenOfMaxBlockSize equals MaxEncodedLen(snappyMaxBlockSize), but is
	// hard coded to be a const instead of a variable, so that obufLen can also
	// be a const. Their equivalence is confirmed by
	// TestMaxEncodedLenOfMaxBlockSize.
	snappyMaxEncodedLenOfMaxBlockSize = 76490
)

const (
	chunkTypeCompressedData   = 0x00
	chunkTypeUncompressedData = 0x01
	chunkTypePadding          = 0xfe
	chunkTypeStreamIdentifier = 0xff
)

var (
	// ErrSnappyCorrupt reports that the input is invalid.
	ErrSnappyCorrupt = errors.New("snappy: corrupt input")
	// ErrSnappyTooLarge reports that the uncompressed length is too large.
	ErrSnappyTooLarge = errors.New("snappy: decoded block is too large")
	// ErrSnappyUnsupported reports that the input isn't supported.
	ErrSnappyUnsupported = errors.New("snappy: unsupported input")

	errUnsupportedLiteralLength = errors.New("snappy: unsupported literal length")
)

// SnappyConverter can read SnappyConverter-compressed streams and convert them to zstd.
// Conversion is done by converting the stream directly from Snappy without intermediate
// full decoding.
// Therefore the compression ratio is much less than what can be done by a full decompression
// and compression, and a faulty Snappy stream may lead to a faulty Zstandard stream without
// any errors being generated.
// No CRC value is being generated and not all CRC values of the Snappy stream are checked.
// However, it provides really fast recompression of Snappy streams.
// The converter can be reused to avoid allocations, even after errors.
type SnappyConverter struct ***REMOVED***
	r     io.Reader
	err   error
	buf   []byte
	block *blockEnc
***REMOVED***

// Convert the Snappy stream supplied in 'in' and write the zStandard stream to 'w'.
// If any error is detected on the Snappy stream it is returned.
// The number of bytes written is returned.
func (r *SnappyConverter) Convert(in io.Reader, w io.Writer) (int64, error) ***REMOVED***
	initPredefined()
	r.err = nil
	r.r = in
	if r.block == nil ***REMOVED***
		r.block = &blockEnc***REMOVED******REMOVED***
		r.block.init()
	***REMOVED***
	r.block.initNewEncode()
	if len(r.buf) != snappyMaxEncodedLenOfMaxBlockSize+snappyChecksumSize ***REMOVED***
		r.buf = make([]byte, snappyMaxEncodedLenOfMaxBlockSize+snappyChecksumSize)
	***REMOVED***
	r.block.litEnc.Reuse = huff0.ReusePolicyNone
	var written int64
	var readHeader bool
	***REMOVED***
		var header []byte
		var n int
		header, r.err = frameHeader***REMOVED***WindowSize: snappyMaxBlockSize***REMOVED***.appendTo(r.buf[:0])

		n, r.err = w.Write(header)
		if r.err != nil ***REMOVED***
			return written, r.err
		***REMOVED***
		written += int64(n)
	***REMOVED***

	for ***REMOVED***
		if !r.readFull(r.buf[:4], true) ***REMOVED***
			// Add empty last block
			r.block.reset(nil)
			r.block.last = true
			err := r.block.encodeLits(r.block.literals, false)
			if err != nil ***REMOVED***
				return written, err
			***REMOVED***
			n, err := w.Write(r.block.output)
			if err != nil ***REMOVED***
				return written, err
			***REMOVED***
			written += int64(n)

			return written, r.err
		***REMOVED***
		chunkType := r.buf[0]
		if !readHeader ***REMOVED***
			if chunkType != chunkTypeStreamIdentifier ***REMOVED***
				println("chunkType != chunkTypeStreamIdentifier", chunkType)
				r.err = ErrSnappyCorrupt
				return written, r.err
			***REMOVED***
			readHeader = true
		***REMOVED***
		chunkLen := int(r.buf[1]) | int(r.buf[2])<<8 | int(r.buf[3])<<16
		if chunkLen > len(r.buf) ***REMOVED***
			println("chunkLen > len(r.buf)", chunkType)
			r.err = ErrSnappyUnsupported
			return written, r.err
		***REMOVED***

		// The chunk types are specified at
		// https://github.com/google/snappy/blob/master/framing_format.txt
		switch chunkType ***REMOVED***
		case chunkTypeCompressedData:
			// Section 4.2. Compressed data (chunk type 0x00).
			if chunkLen < snappyChecksumSize ***REMOVED***
				println("chunkLen < snappyChecksumSize", chunkLen, snappyChecksumSize)
				r.err = ErrSnappyCorrupt
				return written, r.err
			***REMOVED***
			buf := r.buf[:chunkLen]
			if !r.readFull(buf, false) ***REMOVED***
				return written, r.err
			***REMOVED***
			//checksum := uint32(buf[0]) | uint32(buf[1])<<8 | uint32(buf[2])<<16 | uint32(buf[3])<<24
			buf = buf[snappyChecksumSize:]

			n, hdr, err := snappyDecodedLen(buf)
			if err != nil ***REMOVED***
				r.err = err
				return written, r.err
			***REMOVED***
			buf = buf[hdr:]
			if n > snappyMaxBlockSize ***REMOVED***
				println("n > snappyMaxBlockSize", n, snappyMaxBlockSize)
				r.err = ErrSnappyCorrupt
				return written, r.err
			***REMOVED***
			r.block.reset(nil)
			r.block.pushOffsets()
			if err := decodeSnappy(r.block, buf); err != nil ***REMOVED***
				r.err = err
				return written, r.err
			***REMOVED***
			if r.block.size+r.block.extraLits != n ***REMOVED***
				printf("invalid size, want %d, got %d\n", n, r.block.size+r.block.extraLits)
				r.err = ErrSnappyCorrupt
				return written, r.err
			***REMOVED***
			err = r.block.encode(nil, false, false)
			switch err ***REMOVED***
			case errIncompressible:
				r.block.popOffsets()
				r.block.reset(nil)
				r.block.literals, err = snappy.Decode(r.block.literals[:n], r.buf[snappyChecksumSize:chunkLen])
				if err != nil ***REMOVED***
					return written, err
				***REMOVED***
				err = r.block.encodeLits(r.block.literals, false)
				if err != nil ***REMOVED***
					return written, err
				***REMOVED***
			case nil:
			default:
				return written, err
			***REMOVED***

			n, r.err = w.Write(r.block.output)
			if r.err != nil ***REMOVED***
				return written, err
			***REMOVED***
			written += int64(n)
			continue
		case chunkTypeUncompressedData:
			if debug ***REMOVED***
				println("Uncompressed, chunklen", chunkLen)
			***REMOVED***
			// Section 4.3. Uncompressed data (chunk type 0x01).
			if chunkLen < snappyChecksumSize ***REMOVED***
				println("chunkLen < snappyChecksumSize", chunkLen, snappyChecksumSize)
				r.err = ErrSnappyCorrupt
				return written, r.err
			***REMOVED***
			r.block.reset(nil)
			buf := r.buf[:snappyChecksumSize]
			if !r.readFull(buf, false) ***REMOVED***
				return written, r.err
			***REMOVED***
			checksum := uint32(buf[0]) | uint32(buf[1])<<8 | uint32(buf[2])<<16 | uint32(buf[3])<<24
			// Read directly into r.decoded instead of via r.buf.
			n := chunkLen - snappyChecksumSize
			if n > snappyMaxBlockSize ***REMOVED***
				println("n > snappyMaxBlockSize", n, snappyMaxBlockSize)
				r.err = ErrSnappyCorrupt
				return written, r.err
			***REMOVED***
			r.block.literals = r.block.literals[:n]
			if !r.readFull(r.block.literals, false) ***REMOVED***
				return written, r.err
			***REMOVED***
			if snappyCRC(r.block.literals) != checksum ***REMOVED***
				println("literals crc mismatch")
				r.err = ErrSnappyCorrupt
				return written, r.err
			***REMOVED***
			err := r.block.encodeLits(r.block.literals, false)
			if err != nil ***REMOVED***
				return written, err
			***REMOVED***
			n, r.err = w.Write(r.block.output)
			if r.err != nil ***REMOVED***
				return written, err
			***REMOVED***
			written += int64(n)
			continue

		case chunkTypeStreamIdentifier:
			if debug ***REMOVED***
				println("stream id", chunkLen, len(snappyMagicBody))
			***REMOVED***
			// Section 4.1. Stream identifier (chunk type 0xff).
			if chunkLen != len(snappyMagicBody) ***REMOVED***
				println("chunkLen != len(snappyMagicBody)", chunkLen, len(snappyMagicBody))
				r.err = ErrSnappyCorrupt
				return written, r.err
			***REMOVED***
			if !r.readFull(r.buf[:len(snappyMagicBody)], false) ***REMOVED***
				return written, r.err
			***REMOVED***
			for i := 0; i < len(snappyMagicBody); i++ ***REMOVED***
				if r.buf[i] != snappyMagicBody[i] ***REMOVED***
					println("r.buf[i] != snappyMagicBody[i]", r.buf[i], snappyMagicBody[i], i)
					r.err = ErrSnappyCorrupt
					return written, r.err
				***REMOVED***
			***REMOVED***
			continue
		***REMOVED***

		if chunkType <= 0x7f ***REMOVED***
			// Section 4.5. Reserved unskippable chunks (chunk types 0x02-0x7f).
			println("chunkType <= 0x7f")
			r.err = ErrSnappyUnsupported
			return written, r.err
		***REMOVED***
		// Section 4.4 Padding (chunk type 0xfe).
		// Section 4.6. Reserved skippable chunks (chunk types 0x80-0xfd).
		if !r.readFull(r.buf[:chunkLen], false) ***REMOVED***
			return written, r.err
		***REMOVED***
	***REMOVED***
***REMOVED***

// decodeSnappy writes the decoding of src to dst. It assumes that the varint-encoded
// length of the decompressed bytes has already been read.
func decodeSnappy(blk *blockEnc, src []byte) error ***REMOVED***
	//decodeRef(make([]byte, snappyMaxBlockSize), src)
	var s, length int
	lits := blk.extraLits
	var offset uint32
	for s < len(src) ***REMOVED***
		switch src[s] & 0x03 ***REMOVED***
		case snappyTagLiteral:
			x := uint32(src[s] >> 2)
			switch ***REMOVED***
			case x < 60:
				s++
			case x == 60:
				s += 2
				if uint(s) > uint(len(src)) ***REMOVED*** // The uint conversions catch overflow from the previous line.
					println("uint(s) > uint(len(src)", s, src)
					return ErrSnappyCorrupt
				***REMOVED***
				x = uint32(src[s-1])
			case x == 61:
				s += 3
				if uint(s) > uint(len(src)) ***REMOVED*** // The uint conversions catch overflow from the previous line.
					println("uint(s) > uint(len(src)", s, src)
					return ErrSnappyCorrupt
				***REMOVED***
				x = uint32(src[s-2]) | uint32(src[s-1])<<8
			case x == 62:
				s += 4
				if uint(s) > uint(len(src)) ***REMOVED*** // The uint conversions catch overflow from the previous line.
					println("uint(s) > uint(len(src)", s, src)
					return ErrSnappyCorrupt
				***REMOVED***
				x = uint32(src[s-3]) | uint32(src[s-2])<<8 | uint32(src[s-1])<<16
			case x == 63:
				s += 5
				if uint(s) > uint(len(src)) ***REMOVED*** // The uint conversions catch overflow from the previous line.
					println("uint(s) > uint(len(src)", s, src)
					return ErrSnappyCorrupt
				***REMOVED***
				x = uint32(src[s-4]) | uint32(src[s-3])<<8 | uint32(src[s-2])<<16 | uint32(src[s-1])<<24
			***REMOVED***
			if x > snappyMaxBlockSize ***REMOVED***
				println("x > snappyMaxBlockSize", x, snappyMaxBlockSize)
				return ErrSnappyCorrupt
			***REMOVED***
			length = int(x) + 1
			if length <= 0 ***REMOVED***
				println("length <= 0 ", length)

				return errUnsupportedLiteralLength
			***REMOVED***
			//if length > snappyMaxBlockSize-d || uint32(length) > len(src)-s ***REMOVED***
			//	return ErrSnappyCorrupt
			//***REMOVED***

			blk.literals = append(blk.literals, src[s:s+length]...)
			//println(length, "litLen")
			lits += length
			s += length
			continue

		case snappyTagCopy1:
			s += 2
			if uint(s) > uint(len(src)) ***REMOVED*** // The uint conversions catch overflow from the previous line.
				println("uint(s) > uint(len(src)", s, len(src))
				return ErrSnappyCorrupt
			***REMOVED***
			length = 4 + int(src[s-2])>>2&0x7
			offset = uint32(src[s-2])&0xe0<<3 | uint32(src[s-1])

		case snappyTagCopy2:
			s += 3
			if uint(s) > uint(len(src)) ***REMOVED*** // The uint conversions catch overflow from the previous line.
				println("uint(s) > uint(len(src)", s, len(src))
				return ErrSnappyCorrupt
			***REMOVED***
			length = 1 + int(src[s-3])>>2
			offset = uint32(src[s-2]) | uint32(src[s-1])<<8

		case snappyTagCopy4:
			s += 5
			if uint(s) > uint(len(src)) ***REMOVED*** // The uint conversions catch overflow from the previous line.
				println("uint(s) > uint(len(src)", s, len(src))
				return ErrSnappyCorrupt
			***REMOVED***
			length = 1 + int(src[s-5])>>2
			offset = uint32(src[s-4]) | uint32(src[s-3])<<8 | uint32(src[s-2])<<16 | uint32(src[s-1])<<24
		***REMOVED***

		if offset <= 0 || blk.size+lits < int(offset) /*|| length > len(blk)-d */ ***REMOVED***
			println("offset <= 0 || blk.size+lits < int(offset)", offset, blk.size+lits, int(offset), blk.size, lits)

			return ErrSnappyCorrupt
		***REMOVED***

		// Check if offset is one of the recent offsets.
		// Adjusts the output offset accordingly.
		// Gives a tiny bit of compression, typically around 1%.
		if false ***REMOVED***
			offset = blk.matchOffset(offset, uint32(lits))
		***REMOVED*** else ***REMOVED***
			offset += 3
		***REMOVED***

		blk.sequences = append(blk.sequences, seq***REMOVED***
			litLen:   uint32(lits),
			offset:   offset,
			matchLen: uint32(length) - zstdMinMatch,
		***REMOVED***)
		blk.size += length + lits
		lits = 0
	***REMOVED***
	blk.extraLits = lits
	return nil
***REMOVED***

func (r *SnappyConverter) readFull(p []byte, allowEOF bool) (ok bool) ***REMOVED***
	if _, r.err = io.ReadFull(r.r, p); r.err != nil ***REMOVED***
		if r.err == io.ErrUnexpectedEOF || (r.err == io.EOF && !allowEOF) ***REMOVED***
			r.err = ErrSnappyCorrupt
		***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

var crcTable = crc32.MakeTable(crc32.Castagnoli)

// crc implements the checksum specified in section 3 of
// https://github.com/google/snappy/blob/master/framing_format.txt
func snappyCRC(b []byte) uint32 ***REMOVED***
	c := crc32.Update(0, crcTable, b)
	return c>>15 | c<<17 + 0xa282ead8
***REMOVED***

// snappyDecodedLen returns the length of the decoded block and the number of bytes
// that the length header occupied.
func snappyDecodedLen(src []byte) (blockLen, headerLen int, err error) ***REMOVED***
	v, n := binary.Uvarint(src)
	if n <= 0 || v > 0xffffffff ***REMOVED***
		return 0, 0, ErrSnappyCorrupt
	***REMOVED***

	const wordSize = 32 << (^uint(0) >> 32 & 1)
	if wordSize == 32 && v > 0x7fffffff ***REMOVED***
		return 0, 0, ErrSnappyTooLarge
	***REMOVED***
	return int(v), n, nil
***REMOVED***
