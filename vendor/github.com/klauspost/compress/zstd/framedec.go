// Copyright 2019+ Klaus Post. All rights reserved.
// License information can be found in the LICENSE file.
// Based on work by Yann Collet, released under BSD License.

package zstd

import (
	"bytes"
	"encoding/hex"
	"errors"
	"hash"
	"io"
	"sync"

	"github.com/klauspost/compress/zstd/internal/xxhash"
)

type frameDec struct ***REMOVED***
	o      decoderOptions
	crc    hash.Hash64
	offset int64

	WindowSize uint64

	// maxWindowSize is the maximum windows size to support.
	// should never be bigger than max-int.
	maxWindowSize uint64

	// In order queue of blocks being decoded.
	decoding chan *blockDec

	// Frame history passed between blocks
	history history

	rawInput byteBuffer

	// Byte buffer that can be reused for small input blocks.
	bBuf byteBuf

	FrameContentSize uint64
	frameDone        sync.WaitGroup

	DictionaryID  *uint32
	HasCheckSum   bool
	SingleSegment bool

	// asyncRunning indicates whether the async routine processes input on 'decoding'.
	asyncRunningMu sync.Mutex
	asyncRunning   bool
***REMOVED***

const (
	// The minimum Window_Size is 1 KB.
	MinWindowSize = 1 << 10
	MaxWindowSize = 1 << 29
)

var (
	frameMagic          = []byte***REMOVED***0x28, 0xb5, 0x2f, 0xfd***REMOVED***
	skippableFrameMagic = []byte***REMOVED***0x2a, 0x4d, 0x18***REMOVED***
)

func newFrameDec(o decoderOptions) *frameDec ***REMOVED***
	d := frameDec***REMOVED***
		o:             o,
		maxWindowSize: MaxWindowSize,
	***REMOVED***
	if d.maxWindowSize > o.maxDecodedSize ***REMOVED***
		d.maxWindowSize = o.maxDecodedSize
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
	var b []byte
	for ***REMOVED***
		b = br.readSmall(4)
		if b == nil ***REMOVED***
			return io.EOF
		***REMOVED***
		if !bytes.Equal(b[1:4], skippableFrameMagic) || b[0]&0xf0 != 0x50 ***REMOVED***
			if debug ***REMOVED***
				println("Not skippable", hex.EncodeToString(b), hex.EncodeToString(skippableFrameMagic))
			***REMOVED***
			// Break if not skippable frame.
			break
		***REMOVED***
		// Read size to skip
		b = br.readSmall(4)
		if b == nil ***REMOVED***
			println("Reading Frame Size EOF")
			return io.ErrUnexpectedEOF
		***REMOVED***
		n := uint32(b[0]) | (uint32(b[1]) << 8) | (uint32(b[2]) << 16) | (uint32(b[3]) << 24)
		println("Skipping frame with", n, "bytes.")
		err := br.skipN(int(n))
		if err != nil ***REMOVED***
			if debug ***REMOVED***
				println("Reading discarded frame", err)
			***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if !bytes.Equal(b, frameMagic) ***REMOVED***
		println("Got magic numbers: ", b, "want:", frameMagic)
		return ErrMagicMismatch
	***REMOVED***

	// Read Frame_Header_Descriptor
	fhd, err := br.readByte()
	if err != nil ***REMOVED***
		println("Reading Frame_Header_Descriptor", err)
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
			println("Reading Window_Descriptor", err)
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
		b = br.readSmall(int(size))
		if b == nil ***REMOVED***
			if debug ***REMOVED***
				println("Reading Dictionary_ID", io.ErrUnexpectedEOF)
			***REMOVED***
			return io.ErrUnexpectedEOF
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
		if debug ***REMOVED***
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
	d.FrameContentSize = 0
	if fcsSize > 0 ***REMOVED***
		b := br.readSmall(fcsSize)
		if b == nil ***REMOVED***
			println("Reading Frame content", io.ErrUnexpectedEOF)
			return io.ErrUnexpectedEOF
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
		if debug ***REMOVED***
			println("field size bits:", v, "fcsSize:", fcsSize, "FrameContentSize:", d.FrameContentSize, hex.EncodeToString(b[:fcsSize]), "singleseg:", d.SingleSegment, "window:", d.WindowSize)
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

	if d.WindowSize > d.maxWindowSize ***REMOVED***
		printf("window size %d > max %d\n", d.WindowSize, d.maxWindowSize)
		return ErrWindowSizeExceeded
	***REMOVED***
	// The minimum Window_Size is 1 KB.
	if d.WindowSize < MinWindowSize ***REMOVED***
		println("got window size: ", d.WindowSize)
		return ErrWindowSizeTooSmall
	***REMOVED***
	d.history.windowSize = int(d.WindowSize)
	if d.o.lowMem && d.history.windowSize < maxBlockSize ***REMOVED***
		d.history.maxSize = d.history.windowSize * 2
	***REMOVED*** else ***REMOVED***
		d.history.maxSize = d.history.windowSize + maxBlockSize
	***REMOVED***
	// history contains input - maybe we do something
	d.rawInput = br
	return nil
***REMOVED***

// next will start decoding the next block from stream.
func (d *frameDec) next(block *blockDec) error ***REMOVED***
	if debug ***REMOVED***
		printf("decoding new block %p:%p", block, block.data)
	***REMOVED***
	err := block.reset(d.rawInput, d.WindowSize)
	if err != nil ***REMOVED***
		println("block error:", err)
		// Signal the frame decoder we have a problem.
		d.sendErr(block, err)
		return err
	***REMOVED***
	block.input <- struct***REMOVED******REMOVED******REMOVED******REMOVED***
	if debug ***REMOVED***
		println("next block:", block)
	***REMOVED***
	d.asyncRunningMu.Lock()
	defer d.asyncRunningMu.Unlock()
	if !d.asyncRunning ***REMOVED***
		return nil
	***REMOVED***
	if block.Last ***REMOVED***
		// We indicate the frame is done by sending io.EOF
		d.decoding <- block
		return io.EOF
	***REMOVED***
	d.decoding <- block
	return nil
***REMOVED***

// sendEOF will queue an error block on the frame.
// This will cause the frame decoder to return when it encounters the block.
// Returns true if the decoder was added.
func (d *frameDec) sendErr(block *blockDec, err error) bool ***REMOVED***
	d.asyncRunningMu.Lock()
	defer d.asyncRunningMu.Unlock()
	if !d.asyncRunning ***REMOVED***
		return false
	***REMOVED***

	println("sending error", err.Error())
	block.sendErr(err)
	d.decoding <- block
	return true
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
	want := d.rawInput.readSmall(4)
	if want == nil ***REMOVED***
		println("CRC missing?")
		return io.ErrUnexpectedEOF
	***REMOVED***

	if !bytes.Equal(tmp[:], want) ***REMOVED***
		if debug ***REMOVED***
			println("CRC Check Failed:", tmp[:], "!=", want)
		***REMOVED***
		return ErrCRCMismatch
	***REMOVED***
	if debug ***REMOVED***
		println("CRC ok", tmp[:])
	***REMOVED***
	return nil
***REMOVED***

func (d *frameDec) initAsync() ***REMOVED***
	if !d.o.lowMem && !d.SingleSegment ***REMOVED***
		// set max extra size history to 10MB.
		d.history.maxSize = d.history.windowSize + maxBlockSize*5
	***REMOVED***
	// re-alloc if more than one extra block size.
	if d.o.lowMem && cap(d.history.b) > d.history.maxSize+maxBlockSize ***REMOVED***
		d.history.b = make([]byte, 0, d.history.maxSize)
	***REMOVED***
	if cap(d.history.b) < d.history.maxSize ***REMOVED***
		d.history.b = make([]byte, 0, d.history.maxSize)
	***REMOVED***
	if cap(d.decoding) < d.o.concurrent ***REMOVED***
		d.decoding = make(chan *blockDec, d.o.concurrent)
	***REMOVED***
	if debug ***REMOVED***
		h := d.history
		printf("history init. len: %d, cap: %d", len(h.b), cap(h.b))
	***REMOVED***
	d.asyncRunningMu.Lock()
	d.asyncRunning = true
	d.asyncRunningMu.Unlock()
***REMOVED***

// startDecoder will start decoding blocks and write them to the writer.
// The decoder will stop as soon as an error occurs or at end of frame.
// When the frame has finished decoding the *bufio.Reader
// containing the remaining input will be sent on frameDec.frameDone.
func (d *frameDec) startDecoder(output chan decodeOutput) ***REMOVED***
	written := int64(0)

	defer func() ***REMOVED***
		d.asyncRunningMu.Lock()
		d.asyncRunning = false
		d.asyncRunningMu.Unlock()

		// Drain the currently decoding.
		d.history.error = true
	flushdone:
		for ***REMOVED***
			select ***REMOVED***
			case b := <-d.decoding:
				b.history <- &d.history
				output <- <-b.result
			default:
				break flushdone
			***REMOVED***
		***REMOVED***
		println("frame decoder done, signalling done")
		d.frameDone.Done()
	***REMOVED***()
	// Get decoder for first block.
	block := <-d.decoding
	block.history <- &d.history
	for ***REMOVED***
		var next *blockDec
		// Get result
		r := <-block.result
		if r.err != nil ***REMOVED***
			println("Result contained error", r.err)
			output <- r
			return
		***REMOVED***
		if debug ***REMOVED***
			println("got result, from ", d.offset, "to", d.offset+int64(len(r.b)))
			d.offset += int64(len(r.b))
		***REMOVED***
		if !block.Last ***REMOVED***
			// Send history to next block
			select ***REMOVED***
			case next = <-d.decoding:
				if debug ***REMOVED***
					println("Sending ", len(d.history.b), "bytes as history")
				***REMOVED***
				next.history <- &d.history
			default:
				// Wait until we have sent the block, so
				// other decoders can potentially get the decoder.
				next = nil
			***REMOVED***
		***REMOVED***

		// Add checksum, async to decoding.
		if d.HasCheckSum ***REMOVED***
			n, err := d.crc.Write(r.b)
			if err != nil ***REMOVED***
				r.err = err
				if n != len(r.b) ***REMOVED***
					r.err = io.ErrShortWrite
				***REMOVED***
				output <- r
				return
			***REMOVED***
		***REMOVED***
		written += int64(len(r.b))
		if d.SingleSegment && uint64(written) > d.FrameContentSize ***REMOVED***
			println("runDecoder: single segment and", uint64(written), ">", d.FrameContentSize)
			r.err = ErrFrameSizeExceeded
			output <- r
			return
		***REMOVED***
		if block.Last ***REMOVED***
			r.err = d.checkCRC()
			output <- r
			return
		***REMOVED***
		output <- r
		if next == nil ***REMOVED***
			// There was no decoder available, we wait for one now that we have sent to the writer.
			if debug ***REMOVED***
				println("Sending ", len(d.history.b), " bytes as history")
			***REMOVED***
			next = <-d.decoding
			next.history <- &d.history
		***REMOVED***
		block = next
	***REMOVED***
***REMOVED***

// runDecoder will create a sync decoder that will decode a block of data.
func (d *frameDec) runDecoder(dst []byte, dec *blockDec) ([]byte, error) ***REMOVED***
	saved := d.history.b

	// We use the history for output to avoid copying it.
	d.history.b = dst
	// Store input length, so we only check new data.
	crcStart := len(dst)
	var err error
	for ***REMOVED***
		err = dec.reset(d.rawInput, d.WindowSize)
		if err != nil ***REMOVED***
			break
		***REMOVED***
		if debug ***REMOVED***
			println("next block:", dec)
		***REMOVED***
		err = dec.decodeBuf(&d.history)
		if err != nil || dec.Last ***REMOVED***
			break
		***REMOVED***
		if uint64(len(d.history.b)) > d.o.maxDecodedSize ***REMOVED***
			err = ErrDecoderSizeExceeded
			break
		***REMOVED***
		if d.SingleSegment && uint64(len(d.history.b)) > d.o.maxDecodedSize ***REMOVED***
			println("runDecoder: single segment and", uint64(len(d.history.b)), ">", d.o.maxDecodedSize)
			err = ErrFrameSizeExceeded
			break
		***REMOVED***
	***REMOVED***
	dst = d.history.b
	if err == nil ***REMOVED***
		if d.HasCheckSum ***REMOVED***
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
