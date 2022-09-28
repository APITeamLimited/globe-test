// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http2

import (
	"errors"
	"fmt"
	"sync"
)

// Buffer chunks are allocated from a pool to reduce pressure on GC.
// The maximum wasted space per dataBuffer is 2x the largest size class,
// which happens when the dataBuffer has multiple chunks and there is
// one unread byte in both the first and last chunks. We use a few size
// classes to minimize overheads for servers that typically receive very
// small request bodies.
//
// TODO: Benchmark to determine if the pools are necessary. The GC may have
// improved enough that we can instead allocate chunks like this:
// make([]byte, max(16<<10, expectedBytesRemaining))
var (
	dataChunkSizeClasses = []int***REMOVED***
		1 << 10,
		2 << 10,
		4 << 10,
		8 << 10,
		16 << 10,
	***REMOVED***
	dataChunkPools = [...]sync.Pool***REMOVED***
		***REMOVED***New: func() interface***REMOVED******REMOVED*** ***REMOVED*** return make([]byte, 1<<10) ***REMOVED******REMOVED***,
		***REMOVED***New: func() interface***REMOVED******REMOVED*** ***REMOVED*** return make([]byte, 2<<10) ***REMOVED******REMOVED***,
		***REMOVED***New: func() interface***REMOVED******REMOVED*** ***REMOVED*** return make([]byte, 4<<10) ***REMOVED******REMOVED***,
		***REMOVED***New: func() interface***REMOVED******REMOVED*** ***REMOVED*** return make([]byte, 8<<10) ***REMOVED******REMOVED***,
		***REMOVED***New: func() interface***REMOVED******REMOVED*** ***REMOVED*** return make([]byte, 16<<10) ***REMOVED******REMOVED***,
	***REMOVED***
)

func getDataBufferChunk(size int64) []byte ***REMOVED***
	i := 0
	for ; i < len(dataChunkSizeClasses)-1; i++ ***REMOVED***
		if size <= int64(dataChunkSizeClasses[i]) ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return dataChunkPools[i].Get().([]byte)
***REMOVED***

func putDataBufferChunk(p []byte) ***REMOVED***
	for i, n := range dataChunkSizeClasses ***REMOVED***
		if len(p) == n ***REMOVED***
			dataChunkPools[i].Put(p)
			return
		***REMOVED***
	***REMOVED***
	panic(fmt.Sprintf("unexpected buffer len=%v", len(p)))
***REMOVED***

// dataBuffer is an io.ReadWriter backed by a list of data chunks.
// Each dataBuffer is used to read DATA frames on a single stream.
// The buffer is divided into chunks so the server can limit the
// total memory used by a single connection without limiting the
// request body size on any single stream.
type dataBuffer struct ***REMOVED***
	chunks   [][]byte
	r        int   // next byte to read is chunks[0][r]
	w        int   // next byte to write is chunks[len(chunks)-1][w]
	size     int   // total buffered bytes
	expected int64 // we expect at least this many bytes in future Write calls (ignored if <= 0)
***REMOVED***

var errReadEmpty = errors.New("read from empty dataBuffer")

// Read copies bytes from the buffer into p.
// It is an error to read when no data is available.
func (b *dataBuffer) Read(p []byte) (int, error) ***REMOVED***
	if b.size == 0 ***REMOVED***
		return 0, errReadEmpty
	***REMOVED***
	var ntotal int
	for len(p) > 0 && b.size > 0 ***REMOVED***
		readFrom := b.bytesFromFirstChunk()
		n := copy(p, readFrom)
		p = p[n:]
		ntotal += n
		b.r += n
		b.size -= n
		// If the first chunk has been consumed, advance to the next chunk.
		if b.r == len(b.chunks[0]) ***REMOVED***
			putDataBufferChunk(b.chunks[0])
			end := len(b.chunks) - 1
			copy(b.chunks[:end], b.chunks[1:])
			b.chunks[end] = nil
			b.chunks = b.chunks[:end]
			b.r = 0
		***REMOVED***
	***REMOVED***
	return ntotal, nil
***REMOVED***

func (b *dataBuffer) bytesFromFirstChunk() []byte ***REMOVED***
	if len(b.chunks) == 1 ***REMOVED***
		return b.chunks[0][b.r:b.w]
	***REMOVED***
	return b.chunks[0][b.r:]
***REMOVED***

// Len returns the number of bytes of the unread portion of the buffer.
func (b *dataBuffer) Len() int ***REMOVED***
	return b.size
***REMOVED***

// Write appends p to the buffer.
func (b *dataBuffer) Write(p []byte) (int, error) ***REMOVED***
	ntotal := len(p)
	for len(p) > 0 ***REMOVED***
		// If the last chunk is empty, allocate a new chunk. Try to allocate
		// enough to fully copy p plus any additional bytes we expect to
		// receive. However, this may allocate less than len(p).
		want := int64(len(p))
		if b.expected > want ***REMOVED***
			want = b.expected
		***REMOVED***
		chunk := b.lastChunkOrAlloc(want)
		n := copy(chunk[b.w:], p)
		p = p[n:]
		b.w += n
		b.size += n
		b.expected -= int64(n)
	***REMOVED***
	return ntotal, nil
***REMOVED***

func (b *dataBuffer) lastChunkOrAlloc(want int64) []byte ***REMOVED***
	if len(b.chunks) != 0 ***REMOVED***
		last := b.chunks[len(b.chunks)-1]
		if b.w < len(last) ***REMOVED***
			return last
		***REMOVED***
	***REMOVED***
	chunk := getDataBufferChunk(want)
	b.chunks = append(b.chunks, chunk)
	b.w = 0
	return chunk
***REMOVED***
