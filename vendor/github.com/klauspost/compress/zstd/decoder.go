// Copyright 2019+ Klaus Post. All rights reserved.
// License information can be found in the LICENSE file.
// Based on work by Yann Collet, released under BSD License.

package zstd

import (
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

	// Streams ready to be decoded.
	stream chan decodeStream

	// Current read position used for Reader functionality.
	current decoderState

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
	d.current.output = make(chan decodeOutput, d.o.concurrent)
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
	if d.stream == nil ***REMOVED***
		return 0, ErrDecoderNilInput
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
			if !d.nextBlock(n == 0) ***REMOVED***
				return n, nil
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if len(d.current.b) > 0 ***REMOVED***
		if debug ***REMOVED***
			println("returning", n, "still bytes left:", len(d.current.b))
		***REMOVED***
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
// Reset can be called with a nil reader to release references to the previous reader.
// After being called with a nil reader, no other operations than Reset or DecodeAll or Close
// should be used.
func (d *Decoder) Reset(r io.Reader) error ***REMOVED***
	if d.current.err == ErrDecoderClosed ***REMOVED***
		return d.current.err
	***REMOVED***

	d.drainOutput()

	if r == nil ***REMOVED***
		d.current.err = ErrDecoderNilInput
		d.current.flushed = true
		return nil
	***REMOVED***

	if d.stream == nil ***REMOVED***
		d.stream = make(chan decodeStream, 1)
		d.streamWg.Add(1)
		go d.startStreamDecoder(d.stream)
	***REMOVED***

	// If bytes buffer and < 1MB, do sync decoding anyway.
	if bb, ok := r.(byter); ok && bb.Len() < 1<<20 ***REMOVED***
		var bb2 byter
		bb2 = bb
		if debug ***REMOVED***
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
		if debug ***REMOVED***
			println("sync decode to", len(dst), "bytes, err:", err)
		***REMOVED***
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
		if debug ***REMOVED***
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
	for ***REMOVED***
		select ***REMOVED***
		case v := <-d.current.output:
			if v.d != nil ***REMOVED***
				if debug ***REMOVED***
					printf("re-adding decoder %p", v.d)
				***REMOVED***
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
		return 0, ErrDecoderNilInput
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
	if d.current.err == ErrDecoderClosed ***REMOVED***
		return dst, ErrDecoderClosed
	***REMOVED***

	// Grab a block decoder and frame decoder.
	block := <-d.decoders
	frame := block.localFrame
	defer func() ***REMOVED***
		if debug ***REMOVED***
			printf("re-adding decoder: %p", block)
		***REMOVED***
		frame.rawInput = nil
		frame.bBuf = nil
		d.decoders <- block
	***REMOVED***()
	frame.bBuf = input

	for ***REMOVED***
		frame.history.reset()
		err := frame.reset(&frame.bBuf)
		if err == io.EOF ***REMOVED***
			if debug ***REMOVED***
				println("frame reset return EOF")
			***REMOVED***
			return dst, nil
		***REMOVED***
		if frame.DictionaryID != nil ***REMOVED***
			dict, ok := d.dicts[*frame.DictionaryID]
			if !ok ***REMOVED***
				return nil, ErrUnknownDictionary
			***REMOVED***
			frame.history.setDict(&dict)
		***REMOVED***
		if err != nil ***REMOVED***
			return dst, err
		***REMOVED***
		if frame.FrameContentSize > d.o.maxDecodedSize-uint64(len(dst)) ***REMOVED***
			return dst, ErrDecoderSizeExceeded
		***REMOVED***
		if frame.FrameContentSize > 0 && frame.FrameContentSize < 1<<30 ***REMOVED***
			// Never preallocate moe than 1 GB up front.
			if cap(dst)-len(dst) < int(frame.FrameContentSize) ***REMOVED***
				dst2 := make([]byte, len(dst), len(dst)+int(frame.FrameContentSize))
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
			if debug ***REMOVED***
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
	if d.current.d != nil ***REMOVED***
		if debug ***REMOVED***
			printf("re-adding current decoder %p", d.current.d)
		***REMOVED***
		d.decoders <- d.current.d
		d.current.d = nil
	***REMOVED***
	if d.current.err != nil ***REMOVED***
		// Keep error state.
		return blocking
	***REMOVED***

	if blocking ***REMOVED***
		d.current.decodeOutput = <-d.current.output
	***REMOVED*** else ***REMOVED***
		select ***REMOVED***
		case d.current.decodeOutput = <-d.current.output:
		default:
			return false
		***REMOVED***
	***REMOVED***
	if debug ***REMOVED***
		println("got", len(d.current.b), "bytes, error:", d.current.err)
	***REMOVED***
	return true
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
		if debug ***REMOVED***
			println("got new stream")
		***REMOVED***
		br := readerWrapper***REMOVED***r: stream.r***REMOVED***
	decodeStream:
		for ***REMOVED***
			frame.history.reset()
			err := frame.reset(&br)
			if debug && err != nil ***REMOVED***
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
