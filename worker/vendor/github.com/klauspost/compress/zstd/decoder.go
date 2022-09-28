// Copyright 2019+ Klaus Post. All rights reserved.
// License information can be found in the LICENSE file.
// Based on work by Yann Collet, released under BSD License.

package zstd

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"sync"

	"github.com/klauspost/compress/zstd/internal/xxhash"
)

// Decoder provides decoding of zstandard streams.
// The decoder has been designed to operate without allocations after a warmup.
// This means that you should store the decoder for best performance.
// To re-use a stream decoder, use the Reset(r io.Reader) error to switch to another stream.
// A decoder can safely be re-used even if the previous stream failed.
// To release the resources, you must call the Close() function on a decoder.
type Decoder struct ***REMOVED***
	o decoderOptions

	// Unreferenced decoders, ready for use.
	decoders chan *blockDec

	// Current read position used for Reader functionality.
	current decoderState

	// sync stream decoding
	syncStream struct ***REMOVED***
		decodedFrame uint64
		br           readerWrapper
		enabled      bool
		inFrame      bool
	***REMOVED***

	frame *frameDec

	// Custom dictionaries.
	// Always uses copies.
	dicts map[uint32]dict

	// streamWg is the waitgroup for all streams
	streamWg sync.WaitGroup
***REMOVED***

// decoderState is used for maintaining state when the decoder
// is used for streaming.
type decoderState struct ***REMOVED***
	// current block being written to stream.
	decodeOutput

	// output in order to be written to stream.
	output chan decodeOutput

	// cancel remaining output.
	cancel context.CancelFunc

	// crc of current frame
	crc *xxhash.Digest

	flushed bool
***REMOVED***

var (
	// Check the interfaces we want to support.
	_ = io.WriterTo(&Decoder***REMOVED******REMOVED***)
	_ = io.Reader(&Decoder***REMOVED******REMOVED***)
)

// NewReader creates a new decoder.
// A nil Reader can be provided in which case Reset can be used to start a decode.
//
// A Decoder can be used in two modes:
//
// 1) As a stream, or
// 2) For stateless decoding using DecodeAll.
//
// Only a single stream can be decoded concurrently, but the same decoder
// can run multiple concurrent stateless decodes. It is even possible to
// use stateless decodes while a stream is being decoded.
//
// The Reset function can be used to initiate a new stream, which is will considerably
// reduce the allocations normally caused by NewReader.
func NewReader(r io.Reader, opts ...DOption) (*Decoder, error) ***REMOVED***
	initPredefined()
	var d Decoder
	d.o.setDefault()
	for _, o := range opts ***REMOVED***
		err := o(&d.o)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	d.current.crc = xxhash.New()
	d.current.flushed = true

	if r == nil ***REMOVED***
		d.current.err = ErrDecoderNilInput
	***REMOVED***

	// Transfer option dicts.
	d.dicts = make(map[uint32]dict, len(d.o.dicts))
	for _, dc := range d.o.dicts ***REMOVED***
		d.dicts[dc.id] = dc
	***REMOVED***
	d.o.dicts = nil

	// Create decoders
	d.decoders = make(chan *blockDec, d.o.concurrent)
	for i := 0; i < d.o.concurrent; i++ ***REMOVED***
		dec := newBlockDec(d.o.lowMem)
		dec.localFrame = newFrameDec(d.o)
		d.decoders <- dec
	***REMOVED***

	if r == nil ***REMOVED***
		return &d, nil
	***REMOVED***
	return &d, d.Reset(r)
***REMOVED***

// Read bytes from the decompressed stream into p.
// Returns the number of bytes written and any error that occurred.
// When the stream is done, io.EOF will be returned.
func (d *Decoder) Read(p []byte) (int, error) ***REMOVED***
	var n int
	for ***REMOVED***
		if len(d.current.b) > 0 ***REMOVED***
			filled := copy(p, d.current.b)
			p = p[filled:]
			d.current.b = d.current.b[filled:]
			n += filled
		***REMOVED***
		if len(p) == 0 ***REMOVED***
			break
		***REMOVED***
		if len(d.current.b) == 0 ***REMOVED***
			// We have an error and no more data
			if d.current.err != nil ***REMOVED***
				break
			***REMOVED***
			if !d.nextBlock(n == 0) ***REMOVED***
				return n, d.current.err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if len(d.current.b) > 0 ***REMOVED***
		if debugDecoder ***REMOVED***
			println("returning", n, "still bytes left:", len(d.current.b))
		***REMOVED***
		// Only return error at end of block
		return n, nil
	***REMOVED***
	if d.current.err != nil ***REMOVED***
		d.drainOutput()
	***REMOVED***
	if debugDecoder ***REMOVED***
		println("returning", n, d.current.err, len(d.decoders))
	***REMOVED***
	return n, d.current.err
***REMOVED***

// Reset will reset the decoder the supplied stream after the current has finished processing.
// Note that this functionality cannot be used after Close has been called.
// Reset can be called with a nil reader to release references to the previous reader.
// After being called with a nil reader, no other operations than Reset or DecodeAll or Close
// should be used.
func (d *Decoder) Reset(r io.Reader) error ***REMOVED***
	if d.current.err == ErrDecoderClosed ***REMOVED***
		return d.current.err
	***REMOVED***

	d.drainOutput()

	d.syncStream.br.r = nil
	if r == nil ***REMOVED***
		d.current.err = ErrDecoderNilInput
		if len(d.current.b) > 0 ***REMOVED***
			d.current.b = d.current.b[:0]
		***REMOVED***
		d.current.flushed = true
		return nil
	***REMOVED***

	// If bytes buffer and < 5MB, do sync decoding anyway.
	if bb, ok := r.(byter); ok && bb.Len() < 5<<20 ***REMOVED***
		bb2 := bb
		if debugDecoder ***REMOVED***
			println("*bytes.Buffer detected, doing sync decode, len:", bb.Len())
		***REMOVED***
		b := bb2.Bytes()
		var dst []byte
		if cap(d.current.b) > 0 ***REMOVED***
			dst = d.current.b
		***REMOVED***

		dst, err := d.DecodeAll(b, dst[:0])
		if err == nil ***REMOVED***
			err = io.EOF
		***REMOVED***
		d.current.b = dst
		d.current.err = err
		d.current.flushed = true
		if debugDecoder ***REMOVED***
			println("sync decode to", len(dst), "bytes, err:", err)
		***REMOVED***
		return nil
	***REMOVED***
	// Remove current block.
	d.stashDecoder()
	d.current.decodeOutput = decodeOutput***REMOVED******REMOVED***
	d.current.err = nil
	d.current.flushed = false
	d.current.d = nil

	// Ensure no-one else is still running...
	d.streamWg.Wait()
	if d.frame == nil ***REMOVED***
		d.frame = newFrameDec(d.o)
	***REMOVED***

	if d.o.concurrent == 1 ***REMOVED***
		return d.startSyncDecoder(r)
	***REMOVED***

	d.current.output = make(chan decodeOutput, d.o.concurrent)
	ctx, cancel := context.WithCancel(context.Background())
	d.current.cancel = cancel
	d.streamWg.Add(1)
	go d.startStreamDecoder(ctx, r, d.current.output)

	return nil
***REMOVED***

// drainOutput will drain the output until errEndOfStream is sent.
func (d *Decoder) drainOutput() ***REMOVED***
	if d.current.cancel != nil ***REMOVED***
		if debugDecoder ***REMOVED***
			println("cancelling current")
		***REMOVED***
		d.current.cancel()
		d.current.cancel = nil
	***REMOVED***
	if d.current.d != nil ***REMOVED***
		if debugDecoder ***REMOVED***
			printf("re-adding current decoder %p, decoders: %d", d.current.d, len(d.decoders))
		***REMOVED***
		d.decoders <- d.current.d
		d.current.d = nil
		d.current.b = nil
	***REMOVED***
	if d.current.output == nil || d.current.flushed ***REMOVED***
		println("current already flushed")
		return
	***REMOVED***
	for v := range d.current.output ***REMOVED***
		if v.d != nil ***REMOVED***
			if debugDecoder ***REMOVED***
				printf("re-adding decoder %p", v.d)
			***REMOVED***
			d.decoders <- v.d
		***REMOVED***
	***REMOVED***
	d.current.output = nil
	d.current.flushed = true
***REMOVED***

// WriteTo writes data to w until there's no more data to write or when an error occurs.
// The return value n is the number of bytes written.
// Any error encountered during the write is also returned.
func (d *Decoder) WriteTo(w io.Writer) (int64, error) ***REMOVED***
	var n int64
	for ***REMOVED***
		if len(d.current.b) > 0 ***REMOVED***
			n2, err2 := w.Write(d.current.b)
			n += int64(n2)
			if err2 != nil && (d.current.err == nil || d.current.err == io.EOF) ***REMOVED***
				d.current.err = err2
			***REMOVED*** else if n2 != len(d.current.b) ***REMOVED***
				d.current.err = io.ErrShortWrite
			***REMOVED***
		***REMOVED***
		if d.current.err != nil ***REMOVED***
			break
		***REMOVED***
		d.nextBlock(true)
	***REMOVED***
	err := d.current.err
	if err != nil ***REMOVED***
		d.drainOutput()
	***REMOVED***
	if err == io.EOF ***REMOVED***
		err = nil
	***REMOVED***
	return n, err
***REMOVED***

// DecodeAll allows stateless decoding of a blob of bytes.
// Output will be appended to dst, so if the destination size is known
// you can pre-allocate the destination slice to avoid allocations.
// DecodeAll can be used concurrently.
// The Decoder concurrency limits will be respected.
func (d *Decoder) DecodeAll(input, dst []byte) ([]byte, error) ***REMOVED***
	if d.decoders == nil ***REMOVED***
		return dst, ErrDecoderClosed
	***REMOVED***

	// Grab a block decoder and frame decoder.
	block := <-d.decoders
	frame := block.localFrame
	defer func() ***REMOVED***
		if debugDecoder ***REMOVED***
			printf("re-adding decoder: %p", block)
		***REMOVED***
		frame.rawInput = nil
		frame.bBuf = nil
		if frame.history.decoders.br != nil ***REMOVED***
			frame.history.decoders.br.in = nil
		***REMOVED***
		d.decoders <- block
	***REMOVED***()
	frame.bBuf = input

	for ***REMOVED***
		frame.history.reset()
		err := frame.reset(&frame.bBuf)
		if err != nil ***REMOVED***
			if err == io.EOF ***REMOVED***
				if debugDecoder ***REMOVED***
					println("frame reset return EOF")
				***REMOVED***
				return dst, nil
			***REMOVED***
			return dst, err
		***REMOVED***
		if frame.DictionaryID != nil ***REMOVED***
			dict, ok := d.dicts[*frame.DictionaryID]
			if !ok ***REMOVED***
				return nil, ErrUnknownDictionary
			***REMOVED***
			if debugDecoder ***REMOVED***
				println("setting dict", frame.DictionaryID)
			***REMOVED***
			frame.history.setDict(&dict)
		***REMOVED***
		if frame.WindowSize > d.o.maxWindowSize ***REMOVED***
			return dst, ErrWindowSizeExceeded
		***REMOVED***
		if frame.FrameContentSize != fcsUnknown ***REMOVED***
			if frame.FrameContentSize > d.o.maxDecodedSize-uint64(len(dst)) ***REMOVED***
				return dst, ErrDecoderSizeExceeded
			***REMOVED***
			if cap(dst)-len(dst) < int(frame.FrameContentSize) ***REMOVED***
				dst2 := make([]byte, len(dst), len(dst)+int(frame.FrameContentSize)+compressedBlockOverAlloc)
				copy(dst2, dst)
				dst = dst2
			***REMOVED***
		***REMOVED***

		if cap(dst) == 0 ***REMOVED***
			// Allocate len(input) * 2 by default if nothing is provided
			// and we didn't get frame content size.
			size := len(input) * 2
			// Cap to 1 MB.
			if size > 1<<20 ***REMOVED***
				size = 1 << 20
			***REMOVED***
			if uint64(size) > d.o.maxDecodedSize ***REMOVED***
				size = int(d.o.maxDecodedSize)
			***REMOVED***
			dst = make([]byte, 0, size)
		***REMOVED***

		dst, err = frame.runDecoder(dst, block)
		if err != nil ***REMOVED***
			return dst, err
		***REMOVED***
		if len(frame.bBuf) == 0 ***REMOVED***
			if debugDecoder ***REMOVED***
				println("frame dbuf empty")
			***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return dst, nil
***REMOVED***

// nextBlock returns the next block.
// If an error occurs d.err will be set.
// Optionally the function can block for new output.
// If non-blocking mode is used the returned boolean will be false
// if no data was available without blocking.
func (d *Decoder) nextBlock(blocking bool) (ok bool) ***REMOVED***
	if d.current.err != nil ***REMOVED***
		// Keep error state.
		return false
	***REMOVED***
	d.current.b = d.current.b[:0]

	// SYNC:
	if d.syncStream.enabled ***REMOVED***
		if !blocking ***REMOVED***
			return false
		***REMOVED***
		ok = d.nextBlockSync()
		if !ok ***REMOVED***
			d.stashDecoder()
		***REMOVED***
		return ok
	***REMOVED***

	//ASYNC:
	d.stashDecoder()
	if blocking ***REMOVED***
		d.current.decodeOutput, ok = <-d.current.output
	***REMOVED*** else ***REMOVED***
		select ***REMOVED***
		case d.current.decodeOutput, ok = <-d.current.output:
		default:
			return false
		***REMOVED***
	***REMOVED***
	if !ok ***REMOVED***
		// This should not happen, so signal error state...
		d.current.err = io.ErrUnexpectedEOF
		return false
	***REMOVED***
	next := d.current.decodeOutput
	if next.d != nil && next.d.async.newHist != nil ***REMOVED***
		d.current.crc.Reset()
	***REMOVED***
	if debugDecoder ***REMOVED***
		var tmp [4]byte
		binary.LittleEndian.PutUint32(tmp[:], uint32(xxhash.Sum64(next.b)))
		println("got", len(d.current.b), "bytes, error:", d.current.err, "data crc:", tmp)
	***REMOVED***

	if !d.o.ignoreChecksum && len(next.b) > 0 ***REMOVED***
		n, err := d.current.crc.Write(next.b)
		if err == nil ***REMOVED***
			if n != len(next.b) ***REMOVED***
				d.current.err = io.ErrShortWrite
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if next.err == nil && next.d != nil && len(next.d.checkCRC) != 0 ***REMOVED***
		got := d.current.crc.Sum64()
		var tmp [4]byte
		binary.LittleEndian.PutUint32(tmp[:], uint32(got))
		if !d.o.ignoreChecksum && !bytes.Equal(tmp[:], next.d.checkCRC) ***REMOVED***
			if debugDecoder ***REMOVED***
				println("CRC Check Failed:", tmp[:], " (got) !=", next.d.checkCRC, "(on stream)")
			***REMOVED***
			d.current.err = ErrCRCMismatch
		***REMOVED*** else ***REMOVED***
			if debugDecoder ***REMOVED***
				println("CRC ok", tmp[:])
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return true
***REMOVED***

func (d *Decoder) nextBlockSync() (ok bool) ***REMOVED***
	if d.current.d == nil ***REMOVED***
		d.current.d = <-d.decoders
	***REMOVED***
	for len(d.current.b) == 0 ***REMOVED***
		if !d.syncStream.inFrame ***REMOVED***
			d.frame.history.reset()
			d.current.err = d.frame.reset(&d.syncStream.br)
			if d.current.err != nil ***REMOVED***
				return false
			***REMOVED***
			if d.frame.DictionaryID != nil ***REMOVED***
				dict, ok := d.dicts[*d.frame.DictionaryID]
				if !ok ***REMOVED***
					d.current.err = ErrUnknownDictionary
					return false
				***REMOVED*** else ***REMOVED***
					d.frame.history.setDict(&dict)
				***REMOVED***
			***REMOVED***
			if d.frame.WindowSize > d.o.maxDecodedSize || d.frame.WindowSize > d.o.maxWindowSize ***REMOVED***
				d.current.err = ErrDecoderSizeExceeded
				return false
			***REMOVED***

			d.syncStream.decodedFrame = 0
			d.syncStream.inFrame = true
		***REMOVED***
		d.current.err = d.frame.next(d.current.d)
		if d.current.err != nil ***REMOVED***
			return false
		***REMOVED***
		d.frame.history.ensureBlock()
		if debugDecoder ***REMOVED***
			println("History trimmed:", len(d.frame.history.b), "decoded already:", d.syncStream.decodedFrame)
		***REMOVED***
		histBefore := len(d.frame.history.b)
		d.current.err = d.current.d.decodeBuf(&d.frame.history)

		if d.current.err != nil ***REMOVED***
			println("error after:", d.current.err)
			return false
		***REMOVED***
		d.current.b = d.frame.history.b[histBefore:]
		if debugDecoder ***REMOVED***
			println("history after:", len(d.frame.history.b))
		***REMOVED***

		// Check frame size (before CRC)
		d.syncStream.decodedFrame += uint64(len(d.current.b))
		if d.syncStream.decodedFrame > d.frame.FrameContentSize ***REMOVED***
			if debugDecoder ***REMOVED***
				printf("DecodedFrame (%d) > FrameContentSize (%d)\n", d.syncStream.decodedFrame, d.frame.FrameContentSize)
			***REMOVED***
			d.current.err = ErrFrameSizeExceeded
			return false
		***REMOVED***

		// Check FCS
		if d.current.d.Last && d.frame.FrameContentSize != fcsUnknown && d.syncStream.decodedFrame != d.frame.FrameContentSize ***REMOVED***
			if debugDecoder ***REMOVED***
				printf("DecodedFrame (%d) != FrameContentSize (%d)\n", d.syncStream.decodedFrame, d.frame.FrameContentSize)
			***REMOVED***
			d.current.err = ErrFrameSizeMismatch
			return false
		***REMOVED***

		// Update/Check CRC
		if d.frame.HasCheckSum ***REMOVED***
			if !d.o.ignoreChecksum ***REMOVED***
				d.frame.crc.Write(d.current.b)
			***REMOVED***
			if d.current.d.Last ***REMOVED***
				if !d.o.ignoreChecksum ***REMOVED***
					d.current.err = d.frame.checkCRC()
				***REMOVED*** else ***REMOVED***
					d.current.err = d.frame.consumeCRC()
				***REMOVED***
				if d.current.err != nil ***REMOVED***
					println("CRC error:", d.current.err)
					return false
				***REMOVED***
			***REMOVED***
		***REMOVED***
		d.syncStream.inFrame = !d.current.d.Last
	***REMOVED***
	return true
***REMOVED***

func (d *Decoder) stashDecoder() ***REMOVED***
	if d.current.d != nil ***REMOVED***
		if debugDecoder ***REMOVED***
			printf("re-adding current decoder %p", d.current.d)
		***REMOVED***
		d.decoders <- d.current.d
		d.current.d = nil
	***REMOVED***
***REMOVED***

// Close will release all resources.
// It is NOT possible to reuse the decoder after this.
func (d *Decoder) Close() ***REMOVED***
	if d.current.err == ErrDecoderClosed ***REMOVED***
		return
	***REMOVED***
	d.drainOutput()
	if d.current.cancel != nil ***REMOVED***
		d.current.cancel()
		d.streamWg.Wait()
		d.current.cancel = nil
	***REMOVED***
	if d.decoders != nil ***REMOVED***
		close(d.decoders)
		for dec := range d.decoders ***REMOVED***
			dec.Close()
		***REMOVED***
		d.decoders = nil
	***REMOVED***
	if d.current.d != nil ***REMOVED***
		d.current.d.Close()
		d.current.d = nil
	***REMOVED***
	d.current.err = ErrDecoderClosed
***REMOVED***

// IOReadCloser returns the decoder as an io.ReadCloser for convenience.
// Any changes to the decoder will be reflected, so the returned ReadCloser
// can be reused along with the decoder.
// io.WriterTo is also supported by the returned ReadCloser.
func (d *Decoder) IOReadCloser() io.ReadCloser ***REMOVED***
	return closeWrapper***REMOVED***d: d***REMOVED***
***REMOVED***

// closeWrapper wraps a function call as a closer.
type closeWrapper struct ***REMOVED***
	d *Decoder
***REMOVED***

// WriteTo forwards WriteTo calls to the decoder.
func (c closeWrapper) WriteTo(w io.Writer) (n int64, err error) ***REMOVED***
	return c.d.WriteTo(w)
***REMOVED***

// Read forwards read calls to the decoder.
func (c closeWrapper) Read(p []byte) (n int, err error) ***REMOVED***
	return c.d.Read(p)
***REMOVED***

// Close closes the decoder.
func (c closeWrapper) Close() error ***REMOVED***
	c.d.Close()
	return nil
***REMOVED***

type decodeOutput struct ***REMOVED***
	d   *blockDec
	b   []byte
	err error
***REMOVED***

func (d *Decoder) startSyncDecoder(r io.Reader) error ***REMOVED***
	d.frame.history.reset()
	d.syncStream.br = readerWrapper***REMOVED***r: r***REMOVED***
	d.syncStream.inFrame = false
	d.syncStream.enabled = true
	d.syncStream.decodedFrame = 0
	return nil
***REMOVED***

// Create Decoder:
// ASYNC:
// Spawn 3 go routines.
// 0: Read frames and decode block literals.
// 1: Decode sequences.
// 2: Execute sequences, send to output.
func (d *Decoder) startStreamDecoder(ctx context.Context, r io.Reader, output chan decodeOutput) ***REMOVED***
	defer d.streamWg.Done()
	br := readerWrapper***REMOVED***r: r***REMOVED***

	var seqDecode = make(chan *blockDec, d.o.concurrent)
	var seqExecute = make(chan *blockDec, d.o.concurrent)

	// Async 1: Decode sequences...
	go func() ***REMOVED***
		var hist history
		var hasErr bool

		for block := range seqDecode ***REMOVED***
			if hasErr ***REMOVED***
				if block != nil ***REMOVED***
					seqExecute <- block
				***REMOVED***
				continue
			***REMOVED***
			if block.async.newHist != nil ***REMOVED***
				if debugDecoder ***REMOVED***
					println("Async 1: new history, recent:", block.async.newHist.recentOffsets)
				***REMOVED***
				hist.decoders = block.async.newHist.decoders
				hist.recentOffsets = block.async.newHist.recentOffsets
				hist.windowSize = block.async.newHist.windowSize
				if block.async.newHist.dict != nil ***REMOVED***
					hist.setDict(block.async.newHist.dict)
				***REMOVED***
			***REMOVED***
			if block.err != nil || block.Type != blockTypeCompressed ***REMOVED***
				hasErr = block.err != nil
				seqExecute <- block
				continue
			***REMOVED***

			hist.decoders.literals = block.async.literals
			block.err = block.prepareSequences(block.async.seqData, &hist)
			if debugDecoder && block.err != nil ***REMOVED***
				println("prepareSequences returned:", block.err)
			***REMOVED***
			hasErr = block.err != nil
			if block.err == nil ***REMOVED***
				block.err = block.decodeSequences(&hist)
				if debugDecoder && block.err != nil ***REMOVED***
					println("decodeSequences returned:", block.err)
				***REMOVED***
				hasErr = block.err != nil
				//				block.async.sequence = hist.decoders.seq[:hist.decoders.nSeqs]
				block.async.seqSize = hist.decoders.seqSize
			***REMOVED***
			seqExecute <- block
		***REMOVED***
		close(seqExecute)
	***REMOVED***()

	var wg sync.WaitGroup
	wg.Add(1)

	// Async 3: Execute sequences...
	frameHistCache := d.frame.history.b
	go func() ***REMOVED***
		var hist history
		var decodedFrame uint64
		var fcs uint64
		var hasErr bool
		for block := range seqExecute ***REMOVED***
			out := decodeOutput***REMOVED***err: block.err, d: block***REMOVED***
			if block.err != nil || hasErr ***REMOVED***
				hasErr = true
				output <- out
				continue
			***REMOVED***
			if block.async.newHist != nil ***REMOVED***
				if debugDecoder ***REMOVED***
					println("Async 2: new history")
				***REMOVED***
				hist.windowSize = block.async.newHist.windowSize
				hist.allocFrameBuffer = block.async.newHist.allocFrameBuffer
				if block.async.newHist.dict != nil ***REMOVED***
					hist.setDict(block.async.newHist.dict)
				***REMOVED***

				if cap(hist.b) < hist.allocFrameBuffer ***REMOVED***
					if cap(frameHistCache) >= hist.allocFrameBuffer ***REMOVED***
						hist.b = frameHistCache
					***REMOVED*** else ***REMOVED***
						hist.b = make([]byte, 0, hist.allocFrameBuffer)
						println("Alloc history sized", hist.allocFrameBuffer)
					***REMOVED***
				***REMOVED***
				hist.b = hist.b[:0]
				fcs = block.async.fcs
				decodedFrame = 0
			***REMOVED***
			do := decodeOutput***REMOVED***err: block.err, d: block***REMOVED***
			switch block.Type ***REMOVED***
			case blockTypeRLE:
				if debugDecoder ***REMOVED***
					println("add rle block length:", block.RLESize)
				***REMOVED***

				if cap(block.dst) < int(block.RLESize) ***REMOVED***
					if block.lowMem ***REMOVED***
						block.dst = make([]byte, block.RLESize)
					***REMOVED*** else ***REMOVED***
						block.dst = make([]byte, maxBlockSize)
					***REMOVED***
				***REMOVED***
				block.dst = block.dst[:block.RLESize]
				v := block.data[0]
				for i := range block.dst ***REMOVED***
					block.dst[i] = v
				***REMOVED***
				hist.append(block.dst)
				do.b = block.dst
			case blockTypeRaw:
				if debugDecoder ***REMOVED***
					println("add raw block length:", len(block.data))
				***REMOVED***
				hist.append(block.data)
				do.b = block.data
			case blockTypeCompressed:
				if debugDecoder ***REMOVED***
					println("execute with history length:", len(hist.b), "window:", hist.windowSize)
				***REMOVED***
				hist.decoders.seqSize = block.async.seqSize
				hist.decoders.literals = block.async.literals
				do.err = block.executeSequences(&hist)
				hasErr = do.err != nil
				if debugDecoder && hasErr ***REMOVED***
					println("executeSequences returned:", do.err)
				***REMOVED***
				do.b = block.dst
			***REMOVED***
			if !hasErr ***REMOVED***
				decodedFrame += uint64(len(do.b))
				if decodedFrame > fcs ***REMOVED***
					println("fcs exceeded", block.Last, fcs, decodedFrame)
					do.err = ErrFrameSizeExceeded
					hasErr = true
				***REMOVED*** else if block.Last && fcs != fcsUnknown && decodedFrame != fcs ***REMOVED***
					do.err = ErrFrameSizeMismatch
					hasErr = true
				***REMOVED*** else ***REMOVED***
					if debugDecoder ***REMOVED***
						println("fcs ok", block.Last, fcs, decodedFrame)
					***REMOVED***
				***REMOVED***
			***REMOVED***
			output <- do
		***REMOVED***
		close(output)
		frameHistCache = hist.b
		wg.Done()
		if debugDecoder ***REMOVED***
			println("decoder goroutines finished")
		***REMOVED***
	***REMOVED***()

decodeStream:
	for ***REMOVED***
		var hist history
		var hasErr bool

		decodeBlock := func(block *blockDec) ***REMOVED***
			if hasErr ***REMOVED***
				if block != nil ***REMOVED***
					seqDecode <- block
				***REMOVED***
				return
			***REMOVED***
			if block.err != nil || block.Type != blockTypeCompressed ***REMOVED***
				hasErr = block.err != nil
				seqDecode <- block
				return
			***REMOVED***

			remain, err := block.decodeLiterals(block.data, &hist)
			block.err = err
			hasErr = block.err != nil
			if err == nil ***REMOVED***
				block.async.literals = hist.decoders.literals
				block.async.seqData = remain
			***REMOVED*** else if debugDecoder ***REMOVED***
				println("decodeLiterals error:", err)
			***REMOVED***
			seqDecode <- block
		***REMOVED***
		frame := d.frame
		if debugDecoder ***REMOVED***
			println("New frame...")
		***REMOVED***
		var historySent bool
		frame.history.reset()
		err := frame.reset(&br)
		if debugDecoder && err != nil ***REMOVED***
			println("Frame decoder returned", err)
		***REMOVED***
		if err == nil && frame.DictionaryID != nil ***REMOVED***
			dict, ok := d.dicts[*frame.DictionaryID]
			if !ok ***REMOVED***
				err = ErrUnknownDictionary
			***REMOVED*** else ***REMOVED***
				frame.history.setDict(&dict)
			***REMOVED***
		***REMOVED***
		if err == nil && d.frame.WindowSize > d.o.maxWindowSize ***REMOVED***
			err = ErrDecoderSizeExceeded
		***REMOVED***
		if err != nil ***REMOVED***
			select ***REMOVED***
			case <-ctx.Done():
			case dec := <-d.decoders:
				dec.sendErr(err)
				decodeBlock(dec)
			***REMOVED***
			break decodeStream
		***REMOVED***

		// Go through all blocks of the frame.
		for ***REMOVED***
			var dec *blockDec
			select ***REMOVED***
			case <-ctx.Done():
				break decodeStream
			case dec = <-d.decoders:
				// Once we have a decoder, we MUST return it.
			***REMOVED***
			err := frame.next(dec)
			if !historySent ***REMOVED***
				h := frame.history
				if debugDecoder ***REMOVED***
					println("Alloc History:", h.allocFrameBuffer)
				***REMOVED***
				hist.reset()
				if h.dict != nil ***REMOVED***
					hist.setDict(h.dict)
				***REMOVED***
				dec.async.newHist = &h
				dec.async.fcs = frame.FrameContentSize
				historySent = true
			***REMOVED*** else ***REMOVED***
				dec.async.newHist = nil
			***REMOVED***
			if debugDecoder && err != nil ***REMOVED***
				println("next block returned error:", err)
			***REMOVED***
			dec.err = err
			dec.checkCRC = nil
			if dec.Last && frame.HasCheckSum && err == nil ***REMOVED***
				crc, err := frame.rawInput.readSmall(4)
				if err != nil ***REMOVED***
					println("CRC missing?", err)
					dec.err = err
				***REMOVED***
				var tmp [4]byte
				copy(tmp[:], crc)
				dec.checkCRC = tmp[:]
				if debugDecoder ***REMOVED***
					println("found crc to check:", dec.checkCRC)
				***REMOVED***
			***REMOVED***
			err = dec.err
			last := dec.Last
			decodeBlock(dec)
			if err != nil ***REMOVED***
				break decodeStream
			***REMOVED***
			if last ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
	close(seqDecode)
	wg.Wait()
	d.frame.history.b = frameHistCache
***REMOVED***
