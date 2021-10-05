// Copyright 2019+ Klaus Post. All rights reserved.
// License information can be found in the LICENSE file.
// Based on work by Yann Collet, released under BSD License.

package zstd

import (
	"crypto/rand"
	"fmt"
	"io"
	rdebug "runtime/debug"
	"sync"

	"github.com/klauspost/compress/zstd/internal/xxhash"
)

// Encoder provides encoding to Zstandard.
// An Encoder can be used for either compressing a stream via the
// io.WriteCloser interface supported by the Encoder or as multiple independent
// tasks via the EncodeAll function.
// Smaller encodes are encouraged to use the EncodeAll function.
// Use NewWriter to create a new instance.
type Encoder struct ***REMOVED***
	o        encoderOptions
	encoders chan encoder
	state    encoderState
	init     sync.Once
***REMOVED***

type encoder interface ***REMOVED***
	Encode(blk *blockEnc, src []byte)
	EncodeNoHist(blk *blockEnc, src []byte)
	Block() *blockEnc
	CRC() *xxhash.Digest
	AppendCRC([]byte) []byte
	WindowSize(size int64) int32
	UseBlock(*blockEnc)
	Reset(d *dict, singleBlock bool)
***REMOVED***

type encoderState struct ***REMOVED***
	w                io.Writer
	filling          []byte
	current          []byte
	previous         []byte
	encoder          encoder
	writing          *blockEnc
	err              error
	writeErr         error
	nWritten         int64
	nInput           int64
	frameContentSize int64
	headerWritten    bool
	eofWritten       bool
	fullFrameWritten bool

	// This waitgroup indicates an encode is running.
	wg sync.WaitGroup
	// This waitgroup indicates we have a block encoding/writing.
	wWg sync.WaitGroup
***REMOVED***

// NewWriter will create a new Zstandard encoder.
// If the encoder will be used for encoding blocks a nil writer can be used.
func NewWriter(w io.Writer, opts ...EOption) (*Encoder, error) ***REMOVED***
	initPredefined()
	var e Encoder
	e.o.setDefault()
	for _, o := range opts ***REMOVED***
		err := o(&e.o)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	if w != nil ***REMOVED***
		e.Reset(w)
	***REMOVED***
	return &e, nil
***REMOVED***

func (e *Encoder) initialize() ***REMOVED***
	if e.o.concurrent == 0 ***REMOVED***
		e.o.setDefault()
	***REMOVED***
	e.encoders = make(chan encoder, e.o.concurrent)
	for i := 0; i < e.o.concurrent; i++ ***REMOVED***
		enc := e.o.encoder()
		e.encoders <- enc
	***REMOVED***
***REMOVED***

// Reset will re-initialize the writer and new writes will encode to the supplied writer
// as a new, independent stream.
func (e *Encoder) Reset(w io.Writer) ***REMOVED***
	s := &e.state
	s.wg.Wait()
	s.wWg.Wait()
	if cap(s.filling) == 0 ***REMOVED***
		s.filling = make([]byte, 0, e.o.blockSize)
	***REMOVED***
	if cap(s.current) == 0 ***REMOVED***
		s.current = make([]byte, 0, e.o.blockSize)
	***REMOVED***
	if cap(s.previous) == 0 ***REMOVED***
		s.previous = make([]byte, 0, e.o.blockSize)
	***REMOVED***
	if s.encoder == nil ***REMOVED***
		s.encoder = e.o.encoder()
	***REMOVED***
	if s.writing == nil ***REMOVED***
		s.writing = &blockEnc***REMOVED***lowMem: e.o.lowMem***REMOVED***
		s.writing.init()
	***REMOVED***
	s.writing.initNewEncode()
	s.filling = s.filling[:0]
	s.current = s.current[:0]
	s.previous = s.previous[:0]
	s.encoder.Reset(e.o.dict, false)
	s.headerWritten = false
	s.eofWritten = false
	s.fullFrameWritten = false
	s.w = w
	s.err = nil
	s.nWritten = 0
	s.nInput = 0
	s.writeErr = nil
	s.frameContentSize = 0
***REMOVED***

// ResetContentSize will reset and set a content size for the next stream.
// If the bytes written does not match the size given an error will be returned
// when calling Close().
// This is removed when Reset is called.
// Sizes <= 0 results in no content size set.
func (e *Encoder) ResetContentSize(w io.Writer, size int64) ***REMOVED***
	e.Reset(w)
	if size >= 0 ***REMOVED***
		e.state.frameContentSize = size
	***REMOVED***
***REMOVED***

// Write data to the encoder.
// Input data will be buffered and as the buffer fills up
// content will be compressed and written to the output.
// When done writing, use Close to flush the remaining output
// and write CRC if requested.
func (e *Encoder) Write(p []byte) (n int, err error) ***REMOVED***
	s := &e.state
	for len(p) > 0 ***REMOVED***
		if len(p)+len(s.filling) < e.o.blockSize ***REMOVED***
			if e.o.crc ***REMOVED***
				_, _ = s.encoder.CRC().Write(p)
			***REMOVED***
			s.filling = append(s.filling, p...)
			return n + len(p), nil
		***REMOVED***
		add := p
		if len(p)+len(s.filling) > e.o.blockSize ***REMOVED***
			add = add[:e.o.blockSize-len(s.filling)]
		***REMOVED***
		if e.o.crc ***REMOVED***
			_, _ = s.encoder.CRC().Write(add)
		***REMOVED***
		s.filling = append(s.filling, add...)
		p = p[len(add):]
		n += len(add)
		if len(s.filling) < e.o.blockSize ***REMOVED***
			return n, nil
		***REMOVED***
		err := e.nextBlock(false)
		if err != nil ***REMOVED***
			return n, err
		***REMOVED***
		if debugAsserts && len(s.filling) > 0 ***REMOVED***
			panic(len(s.filling))
		***REMOVED***
	***REMOVED***
	return n, nil
***REMOVED***

// nextBlock will synchronize and start compressing input in e.state.filling.
// If an error has occurred during encoding it will be returned.
func (e *Encoder) nextBlock(final bool) error ***REMOVED***
	s := &e.state
	// Wait for current block.
	s.wg.Wait()
	if s.err != nil ***REMOVED***
		return s.err
	***REMOVED***
	if len(s.filling) > e.o.blockSize ***REMOVED***
		return fmt.Errorf("block > maxStoreBlockSize")
	***REMOVED***
	if !s.headerWritten ***REMOVED***
		// If we have a single block encode, do a sync compression.
		if final && len(s.filling) == 0 && !e.o.fullZero ***REMOVED***
			s.headerWritten = true
			s.fullFrameWritten = true
			s.eofWritten = true
			return nil
		***REMOVED***
		if final && len(s.filling) > 0 ***REMOVED***
			s.current = e.EncodeAll(s.filling, s.current[:0])
			var n2 int
			n2, s.err = s.w.Write(s.current)
			if s.err != nil ***REMOVED***
				return s.err
			***REMOVED***
			s.nWritten += int64(n2)
			s.nInput += int64(len(s.filling))
			s.current = s.current[:0]
			s.filling = s.filling[:0]
			s.headerWritten = true
			s.fullFrameWritten = true
			s.eofWritten = true
			return nil
		***REMOVED***

		var tmp [maxHeaderSize]byte
		fh := frameHeader***REMOVED***
			ContentSize:   uint64(s.frameContentSize),
			WindowSize:    uint32(s.encoder.WindowSize(s.frameContentSize)),
			SingleSegment: false,
			Checksum:      e.o.crc,
			DictID:        e.o.dict.ID(),
		***REMOVED***

		dst, err := fh.appendTo(tmp[:0])
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		s.headerWritten = true
		s.wWg.Wait()
		var n2 int
		n2, s.err = s.w.Write(dst)
		if s.err != nil ***REMOVED***
			return s.err
		***REMOVED***
		s.nWritten += int64(n2)
	***REMOVED***
	if s.eofWritten ***REMOVED***
		// Ensure we only write it once.
		final = false
	***REMOVED***

	if len(s.filling) == 0 ***REMOVED***
		// Final block, but no data.
		if final ***REMOVED***
			enc := s.encoder
			blk := enc.Block()
			blk.reset(nil)
			blk.last = true
			blk.encodeRaw(nil)
			s.wWg.Wait()
			_, s.err = s.w.Write(blk.output)
			s.nWritten += int64(len(blk.output))
			s.eofWritten = true
		***REMOVED***
		return s.err
	***REMOVED***

	// Move blocks forward.
	s.filling, s.current, s.previous = s.previous[:0], s.filling, s.current
	s.nInput += int64(len(s.current))
	s.wg.Add(1)
	go func(src []byte) ***REMOVED***
		if debugEncoder ***REMOVED***
			println("Adding block,", len(src), "bytes, final:", final)
		***REMOVED***
		defer func() ***REMOVED***
			if r := recover(); r != nil ***REMOVED***
				s.err = fmt.Errorf("panic while encoding: %v", r)
				rdebug.PrintStack()
			***REMOVED***
			s.wg.Done()
		***REMOVED***()
		enc := s.encoder
		blk := enc.Block()
		enc.Encode(blk, src)
		blk.last = final
		if final ***REMOVED***
			s.eofWritten = true
		***REMOVED***
		// Wait for pending writes.
		s.wWg.Wait()
		if s.writeErr != nil ***REMOVED***
			s.err = s.writeErr
			return
		***REMOVED***
		// Transfer encoders from previous write block.
		blk.swapEncoders(s.writing)
		// Transfer recent offsets to next.
		enc.UseBlock(s.writing)
		s.writing = blk
		s.wWg.Add(1)
		go func() ***REMOVED***
			defer func() ***REMOVED***
				if r := recover(); r != nil ***REMOVED***
					s.writeErr = fmt.Errorf("panic while encoding/writing: %v", r)
					rdebug.PrintStack()
				***REMOVED***
				s.wWg.Done()
			***REMOVED***()
			err := errIncompressible
			// If we got the exact same number of literals as input,
			// assume the literals cannot be compressed.
			if len(src) != len(blk.literals) || len(src) != e.o.blockSize ***REMOVED***
				err = blk.encode(src, e.o.noEntropy, !e.o.allLitEntropy)
			***REMOVED***
			switch err ***REMOVED***
			case errIncompressible:
				if debugEncoder ***REMOVED***
					println("Storing incompressible block as raw")
				***REMOVED***
				blk.encodeRaw(src)
				// In fast mode, we do not transfer offsets, so we don't have to deal with changing the.
			case nil:
			default:
				s.writeErr = err
				return
			***REMOVED***
			_, s.writeErr = s.w.Write(blk.output)
			s.nWritten += int64(len(blk.output))
		***REMOVED***()
	***REMOVED***(s.current)
	return nil
***REMOVED***

// ReadFrom reads data from r until EOF or error.
// The return value n is the number of bytes read.
// Any error except io.EOF encountered during the read is also returned.
//
// The Copy function uses ReaderFrom if available.
func (e *Encoder) ReadFrom(r io.Reader) (n int64, err error) ***REMOVED***
	if debugEncoder ***REMOVED***
		println("Using ReadFrom")
	***REMOVED***

	// Flush any current writes.
	if len(e.state.filling) > 0 ***REMOVED***
		if err := e.nextBlock(false); err != nil ***REMOVED***
			return 0, err
		***REMOVED***
	***REMOVED***
	e.state.filling = e.state.filling[:e.o.blockSize]
	src := e.state.filling
	for ***REMOVED***
		n2, err := r.Read(src)
		if e.o.crc ***REMOVED***
			_, _ = e.state.encoder.CRC().Write(src[:n2])
		***REMOVED***
		// src is now the unfilled part...
		src = src[n2:]
		n += int64(n2)
		switch err ***REMOVED***
		case io.EOF:
			e.state.filling = e.state.filling[:len(e.state.filling)-len(src)]
			if debugEncoder ***REMOVED***
				println("ReadFrom: got EOF final block:", len(e.state.filling))
			***REMOVED***
			return n, nil
		case nil:
		default:
			if debugEncoder ***REMOVED***
				println("ReadFrom: got error:", err)
			***REMOVED***
			e.state.err = err
			return n, err
		***REMOVED***
		if len(src) > 0 ***REMOVED***
			if debugEncoder ***REMOVED***
				println("ReadFrom: got space left in source:", len(src))
			***REMOVED***
			continue
		***REMOVED***
		err = e.nextBlock(false)
		if err != nil ***REMOVED***
			return n, err
		***REMOVED***
		e.state.filling = e.state.filling[:e.o.blockSize]
		src = e.state.filling
	***REMOVED***
***REMOVED***

// Flush will send the currently written data to output
// and block until everything has been written.
// This should only be used on rare occasions where pushing the currently queued data is critical.
func (e *Encoder) Flush() error ***REMOVED***
	s := &e.state
	if len(s.filling) > 0 ***REMOVED***
		err := e.nextBlock(false)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	s.wg.Wait()
	s.wWg.Wait()
	if s.err != nil ***REMOVED***
		return s.err
	***REMOVED***
	return s.writeErr
***REMOVED***

// Close will flush the final output and close the stream.
// The function will block until everything has been written.
// The Encoder can still be re-used after calling this.
func (e *Encoder) Close() error ***REMOVED***
	s := &e.state
	if s.encoder == nil ***REMOVED***
		return nil
	***REMOVED***
	err := e.nextBlock(true)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if s.frameContentSize > 0 ***REMOVED***
		if s.nInput != s.frameContentSize ***REMOVED***
			return fmt.Errorf("frame content size %d given, but %d bytes was written", s.frameContentSize, s.nInput)
		***REMOVED***
	***REMOVED***
	if e.state.fullFrameWritten ***REMOVED***
		return s.err
	***REMOVED***
	s.wg.Wait()
	s.wWg.Wait()

	if s.err != nil ***REMOVED***
		return s.err
	***REMOVED***
	if s.writeErr != nil ***REMOVED***
		return s.writeErr
	***REMOVED***

	// Write CRC
	if e.o.crc && s.err == nil ***REMOVED***
		// heap alloc.
		var tmp [4]byte
		_, s.err = s.w.Write(s.encoder.AppendCRC(tmp[:0]))
		s.nWritten += 4
	***REMOVED***

	// Add padding with content from crypto/rand.Reader
	if s.err == nil && e.o.pad > 0 ***REMOVED***
		add := calcSkippableFrame(s.nWritten, int64(e.o.pad))
		frame, err := skippableFrame(s.filling[:0], add, rand.Reader)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		_, s.err = s.w.Write(frame)
	***REMOVED***
	return s.err
***REMOVED***

// EncodeAll will encode all input in src and append it to dst.
// This function can be called concurrently, but each call will only run on a single goroutine.
// If empty input is given, nothing is returned, unless WithZeroFrames is specified.
// Encoded blocks can be concatenated and the result will be the combined input stream.
// Data compressed with EncodeAll can be decoded with the Decoder,
// using either a stream or DecodeAll.
func (e *Encoder) EncodeAll(src, dst []byte) []byte ***REMOVED***
	if len(src) == 0 ***REMOVED***
		if e.o.fullZero ***REMOVED***
			// Add frame header.
			fh := frameHeader***REMOVED***
				ContentSize:   0,
				WindowSize:    MinWindowSize,
				SingleSegment: true,
				// Adding a checksum would be a waste of space.
				Checksum: false,
				DictID:   0,
			***REMOVED***
			dst, _ = fh.appendTo(dst)

			// Write raw block as last one only.
			var blk blockHeader
			blk.setSize(0)
			blk.setType(blockTypeRaw)
			blk.setLast(true)
			dst = blk.appendTo(dst)
		***REMOVED***
		return dst
	***REMOVED***
	e.init.Do(e.initialize)
	enc := <-e.encoders
	defer func() ***REMOVED***
		// Release encoder reference to last block.
		// If a non-single block is needed the encoder will reset again.
		e.encoders <- enc
	***REMOVED***()
	// Use single segments when above minimum window and below 1MB.
	single := len(src) < 1<<20 && len(src) > MinWindowSize
	if e.o.single != nil ***REMOVED***
		single = *e.o.single
	***REMOVED***
	fh := frameHeader***REMOVED***
		ContentSize:   uint64(len(src)),
		WindowSize:    uint32(enc.WindowSize(int64(len(src)))),
		SingleSegment: single,
		Checksum:      e.o.crc,
		DictID:        e.o.dict.ID(),
	***REMOVED***

	// If less than 1MB, allocate a buffer up front.
	if len(dst) == 0 && cap(dst) == 0 && len(src) < 1<<20 && !e.o.lowMem ***REMOVED***
		dst = make([]byte, 0, len(src))
	***REMOVED***
	dst, err := fh.appendTo(dst)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	// If we can do everything in one block, prefer that.
	if len(src) <= maxCompressedBlockSize ***REMOVED***
		enc.Reset(e.o.dict, true)
		// Slightly faster with no history and everything in one block.
		if e.o.crc ***REMOVED***
			_, _ = enc.CRC().Write(src)
		***REMOVED***
		blk := enc.Block()
		blk.last = true
		if e.o.dict == nil ***REMOVED***
			enc.EncodeNoHist(blk, src)
		***REMOVED*** else ***REMOVED***
			enc.Encode(blk, src)
		***REMOVED***

		// If we got the exact same number of literals as input,
		// assume the literals cannot be compressed.
		err := errIncompressible
		oldout := blk.output
		if len(blk.literals) != len(src) || len(src) != e.o.blockSize ***REMOVED***
			// Output directly to dst
			blk.output = dst
			err = blk.encode(src, e.o.noEntropy, !e.o.allLitEntropy)
		***REMOVED***

		switch err ***REMOVED***
		case errIncompressible:
			if debugEncoder ***REMOVED***
				println("Storing incompressible block as raw")
			***REMOVED***
			dst = blk.encodeRawTo(dst, src)
		case nil:
			dst = blk.output
		default:
			panic(err)
		***REMOVED***
		blk.output = oldout
	***REMOVED*** else ***REMOVED***
		enc.Reset(e.o.dict, false)
		blk := enc.Block()
		for len(src) > 0 ***REMOVED***
			todo := src
			if len(todo) > e.o.blockSize ***REMOVED***
				todo = todo[:e.o.blockSize]
			***REMOVED***
			src = src[len(todo):]
			if e.o.crc ***REMOVED***
				_, _ = enc.CRC().Write(todo)
			***REMOVED***
			blk.pushOffsets()
			enc.Encode(blk, todo)
			if len(src) == 0 ***REMOVED***
				blk.last = true
			***REMOVED***
			err := errIncompressible
			// If we got the exact same number of literals as input,
			// assume the literals cannot be compressed.
			if len(blk.literals) != len(todo) || len(todo) != e.o.blockSize ***REMOVED***
				err = blk.encode(todo, e.o.noEntropy, !e.o.allLitEntropy)
			***REMOVED***

			switch err ***REMOVED***
			case errIncompressible:
				if debugEncoder ***REMOVED***
					println("Storing incompressible block as raw")
				***REMOVED***
				dst = blk.encodeRawTo(dst, todo)
				blk.popOffsets()
			case nil:
				dst = append(dst, blk.output...)
			default:
				panic(err)
			***REMOVED***
			blk.reset(nil)
		***REMOVED***
	***REMOVED***
	if e.o.crc ***REMOVED***
		dst = enc.AppendCRC(dst)
	***REMOVED***
	// Add padding with content from crypto/rand.Reader
	if e.o.pad > 0 ***REMOVED***
		add := calcSkippableFrame(int64(len(dst)), int64(e.o.pad))
		dst, err = skippableFrame(dst, add, rand.Reader)
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***
	***REMOVED***
	return dst
***REMOVED***
