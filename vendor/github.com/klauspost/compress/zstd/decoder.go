// Copyright 2019+ Klaus Post. All rights reserved.
// License information can be found in the LICENSE file.
// Based on work by Yann Collet, released under BSD License.

package zstd

import (
	"bytes"
	"errors"
	"io"
	"sync"
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

	// Unreferenced decoders, ready for use.
	frames chan *frameDec

	// Streams ready to be decoded.
	stream chan decodeStream

	// Current read position used for Reader functionality.
	current decoderState

	// Custom dictionaries
	dicts map[uint32]struct***REMOVED******REMOVED***

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
	cancel chan struct***REMOVED******REMOVED***

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
// 2) For stateless decoding using DecodeAll or DecodeBuffer.
//
// Only a single stream can be decoded concurrently, but the same decoder
// can run multiple concurrent stateless decodes. It is even possible to
// use stateless decodes while a stream is being decoded.
//
// The Reset function can be used to initiate a new stream, which is will considerably
// reduce the allocations normally caused by NewReader.
func NewReader(r io.Reader, opts ...DOption) (*Decoder, error) ***REMOVED***
	var d Decoder
	d.o.setDefault()
	for _, o := range opts ***REMOVED***
		err := o(&d.o)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	d.current.output = make(chan decodeOutput, d.o.concurrent)
	d.current.flushed = true

	// Create decoders
	d.decoders = make(chan *blockDec, d.o.concurrent)
	d.frames = make(chan *frameDec, d.o.concurrent)
	for i := 0; i < d.o.concurrent; i++ ***REMOVED***
		d.frames <- newFrameDec(d.o)
		d.decoders <- newBlockDec(d.o.lowMem)
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
	if d.stream == nil ***REMOVED***
		return 0, errors.New("no input has been initialized")
	***REMOVED***
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
			d.nextBlock()
		***REMOVED***
	***REMOVED***
	if len(d.current.b) > 0 ***REMOVED***
		// Only return error at end of block
		return n, nil
	***REMOVED***
	if d.current.err != nil ***REMOVED***
		d.drainOutput()
	***REMOVED***
	if debug ***REMOVED***
		println("returning", n, d.current.err, len(d.decoders))
	***REMOVED***
	return n, d.current.err
***REMOVED***

// Reset will reset the decoder the supplied stream after the current has finished processing.
// Note that this functionality cannot be used after Close has been called.
func (d *Decoder) Reset(r io.Reader) error ***REMOVED***
	if d.current.err == ErrDecoderClosed ***REMOVED***
		return d.current.err
	***REMOVED***
	if r == nil ***REMOVED***
		return errors.New("nil Reader sent as input")
	***REMOVED***

	if d.stream == nil ***REMOVED***
		d.stream = make(chan decodeStream, 1)
		d.streamWg.Add(1)
		go d.startStreamDecoder(d.stream)
	***REMOVED***

	d.drainOutput()

	// If bytes buffer and < 1MB, do sync decoding anyway.
	if bb, ok := r.(*bytes.Buffer); ok && bb.Len() < 1<<20 ***REMOVED***
		b := bb.Bytes()
		dst, err := d.DecodeAll(b, nil)
		if err == nil ***REMOVED***
			err = io.EOF
		***REMOVED***
		d.current.b = dst
		d.current.err = err
		d.current.flushed = true
		return nil
	***REMOVED***

	// Remove current block.
	d.current.decodeOutput = decodeOutput***REMOVED******REMOVED***
	d.current.err = nil
	d.current.cancel = make(chan struct***REMOVED******REMOVED***)
	d.current.flushed = false
	d.current.d = nil

	d.stream <- decodeStream***REMOVED***
		r:      r,
		output: d.current.output,
		cancel: d.current.cancel,
	***REMOVED***
	return nil
***REMOVED***

// drainOutput will drain the output until errEndOfStream is sent.
func (d *Decoder) drainOutput() ***REMOVED***
	if d.current.cancel != nil ***REMOVED***
		println("cancelling current")
		close(d.current.cancel)
		d.current.cancel = nil
	***REMOVED***
	if d.current.d != nil ***REMOVED***
		println("re-adding current decoder", d.current.d, len(d.decoders))
		d.decoders <- d.current.d
		d.current.d = nil
		d.current.b = nil
	***REMOVED***
	if d.current.output == nil || d.current.flushed ***REMOVED***
		println("current already flushed")
		return
	***REMOVED***
	for ***REMOVED***
		select ***REMOVED***
		case v := <-d.current.output:
			if v.d != nil ***REMOVED***
				println("got decoder", v.d)
				d.decoders <- v.d
			***REMOVED***
			if v.err == errEndOfStream ***REMOVED***
				println("current flushed")
				d.current.flushed = true
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// WriteTo writes data to w until there's no more data to write or when an error occurs.
// The return value n is the number of bytes written.
// Any error encountered during the write is also returned.
func (d *Decoder) WriteTo(w io.Writer) (int64, error) ***REMOVED***
	if d.stream == nil ***REMOVED***
		return 0, errors.New("no input has been initialized")
	***REMOVED***
	var n int64
	for ***REMOVED***
		if len(d.current.b) > 0 ***REMOVED***
			n2, err2 := w.Write(d.current.b)
			n += int64(n2)
			if err2 != nil && d.current.err == nil ***REMOVED***
				d.current.err = err2
				break
			***REMOVED***
		***REMOVED***
		if d.current.err != nil ***REMOVED***
			break
		***REMOVED***
		d.nextBlock()
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
	if d.current.err == ErrDecoderClosed ***REMOVED***
		return dst, ErrDecoderClosed
	***REMOVED***
	//println(len(d.frames), len(d.decoders), d.current)
	block, frame := <-d.decoders, <-d.frames
	defer func() ***REMOVED***
		d.decoders <- block
		frame.rawInput = nil
		d.frames <- frame
	***REMOVED***()
	if cap(dst) == 0 ***REMOVED***
		// Allocate 1MB by default.
		dst = make([]byte, 0, 1<<20)
	***REMOVED***
	br := byteBuf(input)
	for ***REMOVED***
		err := frame.reset(&br)
		if err == io.EOF ***REMOVED***
			return dst, nil
		***REMOVED***
		if err != nil ***REMOVED***
			return dst, err
		***REMOVED***
		if frame.FrameContentSize > d.o.maxDecodedSize-uint64(len(dst)) ***REMOVED***
			return dst, ErrDecoderSizeExceeded
		***REMOVED***
		if frame.FrameContentSize > 0 && frame.FrameContentSize < 1<<30 ***REMOVED***
			// Never preallocate moe than 1 GB up front.
			if uint64(cap(dst)) < frame.FrameContentSize ***REMOVED***
				dst2 := make([]byte, len(dst), len(dst)+int(frame.FrameContentSize))
				copy(dst2, dst)
				dst = dst2
			***REMOVED***
		***REMOVED***
		dst, err = frame.runDecoder(dst, block)
		if err != nil ***REMOVED***
			return dst, err
		***REMOVED***
		if len(br) == 0 ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return dst, nil
***REMOVED***

// nextBlock returns the next block.
// If an error occurs d.err will be set.
func (d *Decoder) nextBlock() ***REMOVED***
	if d.current.d != nil ***REMOVED***
		d.decoders <- d.current.d
		d.current.d = nil
	***REMOVED***
	if d.current.err != nil ***REMOVED***
		// Keep error state.
		return
	***REMOVED***
	d.current.decodeOutput = <-d.current.output
	if debug ***REMOVED***
		println("got", len(d.current.b), "bytes, error:", d.current.err)
	***REMOVED***
***REMOVED***

// Close will release all resources.
// It is NOT possible to reuse the decoder after this.
func (d *Decoder) Close() ***REMOVED***
	if d.current.err == ErrDecoderClosed ***REMOVED***
		return
	***REMOVED***
	d.drainOutput()
	if d.stream != nil ***REMOVED***
		close(d.stream)
		d.streamWg.Wait()
		d.stream = nil
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

type decodeOutput struct ***REMOVED***
	d   *blockDec
	b   []byte
	err error
***REMOVED***

type decodeStream struct ***REMOVED***
	r io.Reader

	// Blocks ready to be written to output.
	output chan decodeOutput

	// cancel reading from the input
	cancel chan struct***REMOVED******REMOVED***
***REMOVED***

// errEndOfStream indicates that everything from the stream was read.
var errEndOfStream = errors.New("end-of-stream")

// Create Decoder:
// Spawn n block decoders. These accept tasks to decode a block.
// Create goroutine that handles stream processing, this will send history to decoders as they are available.
// Decoders update the history as they decode.
// When a block is returned:
// 		a) history is sent to the next decoder,
// 		b) content written to CRC.
// 		c) return data to WRITER.
// 		d) wait for next block to return data.
// Once WRITTEN, the decoders reused by the writer frame decoder for re-use.
func (d *Decoder) startStreamDecoder(inStream chan decodeStream) ***REMOVED***
	defer d.streamWg.Done()
	frame := newFrameDec(d.o)
	for stream := range inStream ***REMOVED***
		br := readerWrapper***REMOVED***r: stream.r***REMOVED***
	decodeStream:
		for ***REMOVED***
			err := frame.reset(&br)
			if debug && err != nil ***REMOVED***
				println("Frame decoder returned", err)
			***REMOVED***
			if err != nil ***REMOVED***
				stream.output <- decodeOutput***REMOVED***
					err: err,
				***REMOVED***
				break
			***REMOVED***
			if debug ***REMOVED***
				println("starting frame decoder")
			***REMOVED***

			// This goroutine will forward history between frames.
			frame.frameDone.Add(1)
			frame.initAsync()

			go frame.startDecoder(stream.output)
		decodeFrame:
			// Go through all blocks of the frame.
			for ***REMOVED***
				dec := <-d.decoders
				select ***REMOVED***
				case <-stream.cancel:
					if !frame.sendErr(dec, io.EOF) ***REMOVED***
						// To not let the decoder dangle, send it back.
						stream.output <- decodeOutput***REMOVED***d: dec***REMOVED***
					***REMOVED***
					break decodeStream
				default:
				***REMOVED***
				err := frame.next(dec)
				switch err ***REMOVED***
				case io.EOF:
					// End of current frame, no error
					println("EOF on next block")
					break decodeFrame
				case nil:
					continue
				default:
					println("block decoder returned", err)
					break decodeStream
				***REMOVED***
			***REMOVED***
			// All blocks have started decoding, check if there are more frames.
			println("waiting for done")
			frame.frameDone.Wait()
			println("done waiting...")
		***REMOVED***
		frame.frameDone.Wait()
		println("Sending EOS")
		stream.output <- decodeOutput***REMOVED***err: errEndOfStream***REMOVED***
	***REMOVED***
***REMOVED***
