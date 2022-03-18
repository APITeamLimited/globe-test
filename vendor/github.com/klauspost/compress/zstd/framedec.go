// Copyright 2019+ Klaus Post. All rights reserved.
// License information can be found in the LICENSE file.
// Based on work by Yann Collet, released under BSD License.

package zstd

import (
	"bytes"
	"encoding/hex"
	"errors"
	"io"

	"github.com/klauspost/compress/zstd/internal/xxhash"
)

type frameDec struct ***REMOVED***
	o   decoderOptions
	crc *xxhash.Digest

	WindowSize uint64

	// Frame history passed between blocks
	history history

	rawInput byteBuffer

	// Byte buffer that can be reused for small input blocks.
	bBuf byteBuf

	FrameContentSize uint64

	DictionaryID  *uint32
	HasCheckSum   bool
	SingleSegment bool
***REMOVED***

const (
	// MinWindowSize is the minimum Window Size, which is 1 KB.
	MinWindowSize = 1 << 10

	// MaxWindowSize is the maximum encoder window size
	// and the default decoder maximum window size.
	MaxWindowSize = 1 << 29
)

var (
	frameMagic          = []byte***REMOVED***0x28, 0xb5, 0x2f, 0xfd***REMOVED***
	skippableFrameMagic = []byte***REMOVED***0x2a, 0x4d, 0x18***REMOVED***
)

func newFrameDec(o decoderOptions) *frameDec ***REMOVED***
	if o.maxWindowSize > o.maxDecodedSize ***REMOVED***
		o.maxWindowSize = o.maxDecodedSize
	***REMOVED***
	d := frameDec***REMOVED***
		o: o,
	***REMOVED***
	return &d
***REMOVED***

// reset will read the frame header and prepare for block decoding.
// If nothing can be read from the input, io.EOF will be returned.
// Any other error indicated that the stream contained data, but
// there was a problem.
func (d *frameDec) reset(br byteBuffer) error ***REMOVED***
	d.HasCheckSum = false
	d.WindowSize = 0
	var signature [4]byte
	for ***REMOVED***
		var err error
		// Check if we can read more...
		b, err := br.readSmall(1)
		switch err ***REMOVED***
		case io.EOF, io.ErrUnexpectedEOF:
			return io.EOF
		default:
			return err
		case nil:
			signature[0] = b[0]
		***REMOVED***
		// Read the rest, don't allow io.ErrUnexpectedEOF
		b, err = br.readSmall(3)
		switch err ***REMOVED***
		case io.EOF:
			return io.EOF
		default:
			return err
		case nil:
			copy(signature[1:], b)
		***REMOVED***

		if !bytes.Equal(signature[1:4], skippableFrameMagic) || signature[0]&0xf0 != 0x50 ***REMOVED***
			if debugDecoder ***REMOVED***
				println("Not skippable", hex.EncodeToString(signature[:]), hex.EncodeToString(skippableFrameMagic))
			***REMOVED***
			// Break if not skippable frame.
			break
		***REMOVED***
		// Read size to skip
		b, err = br.readSmall(4)
		if err != nil ***REMOVED***
			if debugDecoder ***REMOVED***
				println("Reading Frame Size", err)
			***REMOVED***
			return err
		***REMOVED***
		n := uint32(b[0]) | (uint32(b[1]) << 8) | (uint32(b[2]) << 16) | (uint32(b[3]) << 24)
		println("Skipping frame with", n, "bytes.")
		err = br.skipN(int(n))
		if err != nil ***REMOVED***
			if debugDecoder ***REMOVED***
				println("Reading discarded frame", err)
			***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if !bytes.Equal(signature[:], frameMagic) ***REMOVED***
		if debugDecoder ***REMOVED***
			println("Got magic numbers: ", signature, "want:", frameMagic)
		***REMOVED***
		return ErrMagicMismatch
	***REMOVED***

	// Read Frame_Header_Descriptor
	fhd, err := br.readByte()
	if err != nil ***REMOVED***
		if debugDecoder ***REMOVED***
			println("Reading Frame_Header_Descriptor", err)
		***REMOVED***
		return err
	***REMOVED***
	d.SingleSegment = fhd&(1<<5) != 0

	if fhd&(1<<3) != 0 ***REMOVED***
		return errors.New("reserved bit set on frame header")
	***REMOVED***

	// Read Window_Descriptor
	// https://github.com/facebook/zstd/blob/dev/doc/zstd_compression_format.md#window_descriptor
	d.WindowSize = 0
	if !d.SingleSegment ***REMOVED***
		wd, err := br.readByte()
		if err != nil ***REMOVED***
			if debugDecoder ***REMOVED***
				println("Reading Window_Descriptor", err)
			***REMOVED***
			return err
		***REMOVED***
		printf("raw: %x, mantissa: %d, exponent: %d\n", wd, wd&7, wd>>3)
		windowLog := 10 + (wd >> 3)
		windowBase := uint64(1) << windowLog
		windowAdd := (windowBase / 8) * uint64(wd&0x7)
		d.WindowSize = windowBase + windowAdd
	***REMOVED***

	// Read Dictionary_ID
	// https://github.com/facebook/zstd/blob/dev/doc/zstd_compression_format.md#dictionary_id
	d.DictionaryID = nil
	if size := fhd & 3; size != 0 ***REMOVED***
		if size == 3 ***REMOVED***
			size = 4
		***REMOVED***

		b, err := br.readSmall(int(size))
		if err != nil ***REMOVED***
			println("Reading Dictionary_ID", err)
			return err
		***REMOVED***
		var id uint32
		switch size ***REMOVED***
		case 1:
			id = uint32(b[0])
		case 2:
			id = uint32(b[0]) | (uint32(b[1]) << 8)
		case 4:
			id = uint32(b[0]) | (uint32(b[1]) << 8) | (uint32(b[2]) << 16) | (uint32(b[3]) << 24)
		***REMOVED***
		if debugDecoder ***REMOVED***
			println("Dict size", size, "ID:", id)
		***REMOVED***
		if id > 0 ***REMOVED***
			// ID 0 means "sorry, no dictionary anyway".
			// https://github.com/facebook/zstd/blob/dev/doc/zstd_compression_format.md#dictionary-format
			d.DictionaryID = &id
		***REMOVED***
	***REMOVED***

	// Read Frame_Content_Size
	// https://github.com/facebook/zstd/blob/dev/doc/zstd_compression_format.md#frame_content_size
	var fcsSize int
	v := fhd >> 6
	switch v ***REMOVED***
	case 0:
		if d.SingleSegment ***REMOVED***
			fcsSize = 1
		***REMOVED***
	default:
		fcsSize = 1 << v
	***REMOVED***
	d.FrameContentSize = fcsUnknown
	if fcsSize > 0 ***REMOVED***
		b, err := br.readSmall(fcsSize)
		if err != nil ***REMOVED***
			println("Reading Frame content", err)
			return err
		***REMOVED***
		switch fcsSize ***REMOVED***
		case 1:
			d.FrameContentSize = uint64(b[0])
		case 2:
			// When FCS_Field_Size is 2, the offset of 256 is added.
			d.FrameContentSize = uint64(b[0]) | (uint64(b[1]) << 8) + 256
		case 4:
			d.FrameContentSize = uint64(b[0]) | (uint64(b[1]) << 8) | (uint64(b[2]) << 16) | (uint64(b[3]) << 24)
		case 8:
			d1 := uint32(b[0]) | (uint32(b[1]) << 8) | (uint32(b[2]) << 16) | (uint32(b[3]) << 24)
			d2 := uint32(b[4]) | (uint32(b[5]) << 8) | (uint32(b[6]) << 16) | (uint32(b[7]) << 24)
			d.FrameContentSize = uint64(d1) | (uint64(d2) << 32)
		***REMOVED***
		if debugDecoder ***REMOVED***
			println("Read FCS:", d.FrameContentSize)
		***REMOVED***
	***REMOVED***

	// Move this to shared.
	d.HasCheckSum = fhd&(1<<2) != 0
	if d.HasCheckSum ***REMOVED***
		if d.crc == nil ***REMOVED***
			d.crc = xxhash.New()
		***REMOVED***
		d.crc.Reset()
	***REMOVED***

	if d.WindowSize == 0 && d.SingleSegment ***REMOVED***
		// We may not need window in this case.
		d.WindowSize = d.FrameContentSize
		if d.WindowSize < MinWindowSize ***REMOVED***
			d.WindowSize = MinWindowSize
		***REMOVED***
	***REMOVED***

	if d.WindowSize > uint64(d.o.maxWindowSize) ***REMOVED***
		if debugDecoder ***REMOVED***
			printf("window size %d > max %d\n", d.WindowSize, d.o.maxWindowSize)
		***REMOVED***
		return ErrWindowSizeExceeded
	***REMOVED***
	// The minimum Window_Size is 1 KB.
	if d.WindowSize < MinWindowSize ***REMOVED***
		if debugDecoder ***REMOVED***
			println("got window size: ", d.WindowSize)
		***REMOVED***
		return ErrWindowSizeTooSmall
	***REMOVED***
	d.history.windowSize = int(d.WindowSize)
	if d.o.lowMem && d.history.windowSize < maxBlockSize ***REMOVED***
		d.history.allocFrameBuffer = d.history.windowSize * 2
		// TODO: Maybe use FrameContent size
	***REMOVED*** else ***REMOVED***
		d.history.allocFrameBuffer = d.history.windowSize + maxBlockSize
	***REMOVED***

	if debugDecoder ***REMOVED***
		println("Frame: Dict:", d.DictionaryID, "FrameContentSize:", d.FrameContentSize, "singleseg:", d.SingleSegment, "window:", d.WindowSize, "crc:", d.HasCheckSum)
	***REMOVED***

	// history contains input - maybe we do something
	d.rawInput = br
	return nil
***REMOVED***

// next will start decoding the next block from stream.
func (d *frameDec) next(block *blockDec) error ***REMOVED***
	if debugDecoder ***REMOVED***
		println("decoding new block")
	***REMOVED***
	err := block.reset(d.rawInput, d.WindowSize)
	if err != nil ***REMOVED***
		println("block error:", err)
		// Signal the frame decoder we have a problem.
		block.sendErr(err)
		return err
	***REMOVED***
	return nil
***REMOVED***

// checkCRC will check the checksum if the frame has one.
// Will return ErrCRCMismatch if crc check failed, otherwise nil.
func (d *frameDec) checkCRC() error ***REMOVED***
	if !d.HasCheckSum ***REMOVED***
		return nil
	***REMOVED***
	var tmp [4]byte
	got := d.crc.Sum64()
	// Flip to match file order.
	tmp[0] = byte(got >> 0)
	tmp[1] = byte(got >> 8)
	tmp[2] = byte(got >> 16)
	tmp[3] = byte(got >> 24)

	// We can overwrite upper tmp now
	want, err := d.rawInput.readSmall(4)
	if err != nil ***REMOVED***
		println("CRC missing?", err)
		return err
	***REMOVED***

	if !bytes.Equal(tmp[:], want) && !ignoreCRC ***REMOVED***
		if debugDecoder ***REMOVED***
			println("CRC Check Failed:", tmp[:], "!=", want)
		***REMOVED***
		return ErrCRCMismatch
	***REMOVED***
	if debugDecoder ***REMOVED***
		println("CRC ok", tmp[:])
	***REMOVED***
	return nil
***REMOVED***

// runDecoder will create a sync decoder that will decode a block of data.
func (d *frameDec) runDecoder(dst []byte, dec *blockDec) ([]byte, error) ***REMOVED***
	saved := d.history.b

	// We use the history for output to avoid copying it.
	d.history.b = dst
	d.history.ignoreBuffer = len(dst)
	// Store input length, so we only check new data.
	crcStart := len(dst)
	var err error
	for ***REMOVED***
		err = dec.reset(d.rawInput, d.WindowSize)
		if err != nil ***REMOVED***
			break
		***REMOVED***
		if debugDecoder ***REMOVED***
			println("next block:", dec)
		***REMOVED***
		err = dec.decodeBuf(&d.history)
		if err != nil ***REMOVED***
			break
		***REMOVED***
		if uint64(len(d.history.b)) > d.o.maxDecodedSize ***REMOVED***
			err = ErrDecoderSizeExceeded
			break
		***REMOVED***
		if uint64(len(d.history.b)-crcStart) > d.FrameContentSize ***REMOVED***
			println("runDecoder: FrameContentSize exceeded", uint64(len(d.history.b)-crcStart), ">", d.FrameContentSize)
			err = ErrFrameSizeExceeded
			break
		***REMOVED***
		if dec.Last ***REMOVED***
			break
		***REMOVED***
		if debugDecoder ***REMOVED***
			println("runDecoder: FrameContentSize", uint64(len(d.history.b)-crcStart), "<=", d.FrameContentSize)
		***REMOVED***
	***REMOVED***
	dst = d.history.b
	if err == nil ***REMOVED***
		if d.FrameContentSize != fcsUnknown && uint64(len(d.history.b)-crcStart) != d.FrameContentSize ***REMOVED***
			err = ErrFrameSizeMismatch
		***REMOVED*** else if d.HasCheckSum ***REMOVED***
			var n int
			n, err = d.crc.Write(dst[crcStart:])
			if err == nil ***REMOVED***
				if n != len(dst)-crcStart ***REMOVED***
					err = io.ErrShortWrite
				***REMOVED*** else ***REMOVED***
					err = d.checkCRC()
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	d.history.b = saved
	return dst, err
***REMOVED***
