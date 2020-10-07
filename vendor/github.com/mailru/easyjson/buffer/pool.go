// Package buffer implements a buffer for serialization, consisting of a chain of []byte-s to
// reduce copying and to allow reuse of individual chunks.
package buffer

import (
	"io"
	"net"
	"sync"
)

// PoolConfig contains configuration for the allocation and reuse strategy.
type PoolConfig struct ***REMOVED***
	StartSize  int // Minimum chunk size that is allocated.
	PooledSize int // Minimum chunk size that is reused, reusing chunks too small will result in overhead.
	MaxSize    int // Maximum chunk size that will be allocated.
***REMOVED***

var config = PoolConfig***REMOVED***
	StartSize:  128,
	PooledSize: 512,
	MaxSize:    32768,
***REMOVED***

// Reuse pool: chunk size -> pool.
var buffers = map[int]*sync.Pool***REMOVED******REMOVED***

func initBuffers() ***REMOVED***
	for l := config.PooledSize; l <= config.MaxSize; l *= 2 ***REMOVED***
		buffers[l] = new(sync.Pool)
	***REMOVED***
***REMOVED***

func init() ***REMOVED***
	initBuffers()
***REMOVED***

// Init sets up a non-default pooling and allocation strategy. Should be run before serialization is done.
func Init(cfg PoolConfig) ***REMOVED***
	config = cfg
	initBuffers()
***REMOVED***

// putBuf puts a chunk to reuse pool if it can be reused.
func putBuf(buf []byte) ***REMOVED***
	size := cap(buf)
	if size < config.PooledSize ***REMOVED***
		return
	***REMOVED***
	if c := buffers[size]; c != nil ***REMOVED***
		c.Put(buf[:0])
	***REMOVED***
***REMOVED***

// getBuf gets a chunk from reuse pool or creates a new one if reuse failed.
func getBuf(size int) []byte ***REMOVED***
	if size >= config.PooledSize ***REMOVED***
		if c := buffers[size]; c != nil ***REMOVED***
			v := c.Get()
			if v != nil ***REMOVED***
				return v.([]byte)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return make([]byte, 0, size)
***REMOVED***

// Buffer is a buffer optimized for serialization without extra copying.
type Buffer struct ***REMOVED***

	// Buf is the current chunk that can be used for serialization.
	Buf []byte

	toPool []byte
	bufs   [][]byte
***REMOVED***

// EnsureSpace makes sure that the current chunk contains at least s free bytes,
// possibly creating a new chunk.
func (b *Buffer) EnsureSpace(s int) ***REMOVED***
	if cap(b.Buf)-len(b.Buf) < s ***REMOVED***
		b.ensureSpaceSlow(s)
	***REMOVED***
***REMOVED***

func (b *Buffer) ensureSpaceSlow(s int) ***REMOVED***
	l := len(b.Buf)
	if l > 0 ***REMOVED***
		if cap(b.toPool) != cap(b.Buf) ***REMOVED***
			// Chunk was reallocated, toPool can be pooled.
			putBuf(b.toPool)
		***REMOVED***
		if cap(b.bufs) == 0 ***REMOVED***
			b.bufs = make([][]byte, 0, 8)
		***REMOVED***
		b.bufs = append(b.bufs, b.Buf)
		l = cap(b.toPool) * 2
	***REMOVED*** else ***REMOVED***
		l = config.StartSize
	***REMOVED***

	if l > config.MaxSize ***REMOVED***
		l = config.MaxSize
	***REMOVED***
	b.Buf = getBuf(l)
	b.toPool = b.Buf
***REMOVED***

// AppendByte appends a single byte to buffer.
func (b *Buffer) AppendByte(data byte) ***REMOVED***
	b.EnsureSpace(1)
	b.Buf = append(b.Buf, data)
***REMOVED***

// AppendBytes appends a byte slice to buffer.
func (b *Buffer) AppendBytes(data []byte) ***REMOVED***
	if len(data) <= cap(b.Buf)-len(b.Buf) ***REMOVED***
		b.Buf = append(b.Buf, data...) // fast path
	***REMOVED*** else ***REMOVED***
		b.appendBytesSlow(data)
	***REMOVED***
***REMOVED***

func (b *Buffer) appendBytesSlow(data []byte) ***REMOVED***
	for len(data) > 0 ***REMOVED***
		b.EnsureSpace(1)

		sz := cap(b.Buf) - len(b.Buf)
		if sz > len(data) ***REMOVED***
			sz = len(data)
		***REMOVED***

		b.Buf = append(b.Buf, data[:sz]...)
		data = data[sz:]
	***REMOVED***
***REMOVED***

// AppendString appends a string to buffer.
func (b *Buffer) AppendString(data string) ***REMOVED***
	if len(data) <= cap(b.Buf)-len(b.Buf) ***REMOVED***
		b.Buf = append(b.Buf, data...) // fast path
	***REMOVED*** else ***REMOVED***
		b.appendStringSlow(data)
	***REMOVED***
***REMOVED***

func (b *Buffer) appendStringSlow(data string) ***REMOVED***
	for len(data) > 0 ***REMOVED***
		b.EnsureSpace(1)

		sz := cap(b.Buf) - len(b.Buf)
		if sz > len(data) ***REMOVED***
			sz = len(data)
		***REMOVED***

		b.Buf = append(b.Buf, data[:sz]...)
		data = data[sz:]
	***REMOVED***
***REMOVED***

// Size computes the size of a buffer by adding sizes of every chunk.
func (b *Buffer) Size() int ***REMOVED***
	size := len(b.Buf)
	for _, buf := range b.bufs ***REMOVED***
		size += len(buf)
	***REMOVED***
	return size
***REMOVED***

// DumpTo outputs the contents of a buffer to a writer and resets the buffer.
func (b *Buffer) DumpTo(w io.Writer) (written int, err error) ***REMOVED***
	bufs := net.Buffers(b.bufs)
	if len(b.Buf) > 0 ***REMOVED***
		bufs = append(bufs, b.Buf)
	***REMOVED***
	n, err := bufs.WriteTo(w)

	for _, buf := range b.bufs ***REMOVED***
		putBuf(buf)
	***REMOVED***
	putBuf(b.toPool)

	b.bufs = nil
	b.Buf = nil
	b.toPool = nil

	return int(n), err
***REMOVED***

// BuildBytes creates a single byte slice with all the contents of the buffer. Data is
// copied if it does not fit in a single chunk. You can optionally provide one byte
// slice as argument that it will try to reuse.
func (b *Buffer) BuildBytes(reuse ...[]byte) []byte ***REMOVED***
	if len(b.bufs) == 0 ***REMOVED***
		ret := b.Buf
		b.toPool = nil
		b.Buf = nil
		return ret
	***REMOVED***

	var ret []byte
	size := b.Size()

	// If we got a buffer as argument and it is big enough, reuse it.
	if len(reuse) == 1 && cap(reuse[0]) >= size ***REMOVED***
		ret = reuse[0][:0]
	***REMOVED*** else ***REMOVED***
		ret = make([]byte, 0, size)
	***REMOVED***
	for _, buf := range b.bufs ***REMOVED***
		ret = append(ret, buf...)
		putBuf(buf)
	***REMOVED***

	ret = append(ret, b.Buf...)
	putBuf(b.toPool)

	b.bufs = nil
	b.toPool = nil
	b.Buf = nil

	return ret
***REMOVED***

type readCloser struct ***REMOVED***
	offset int
	bufs   [][]byte
***REMOVED***

func (r *readCloser) Read(p []byte) (n int, err error) ***REMOVED***
	for _, buf := range r.bufs ***REMOVED***
		// Copy as much as we can.
		x := copy(p[n:], buf[r.offset:])
		n += x // Increment how much we filled.

		// Did we empty the whole buffer?
		if r.offset+x == len(buf) ***REMOVED***
			// On to the next buffer.
			r.offset = 0
			r.bufs = r.bufs[1:]

			// We can release this buffer.
			putBuf(buf)
		***REMOVED*** else ***REMOVED***
			r.offset += x
		***REMOVED***

		if n == len(p) ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	// No buffers left or nothing read?
	if len(r.bufs) == 0 ***REMOVED***
		err = io.EOF
	***REMOVED***
	return
***REMOVED***

func (r *readCloser) Close() error ***REMOVED***
	// Release all remaining buffers.
	for _, buf := range r.bufs ***REMOVED***
		putBuf(buf)
	***REMOVED***
	// In case Close gets called multiple times.
	r.bufs = nil

	return nil
***REMOVED***

// ReadCloser creates an io.ReadCloser with all the contents of the buffer.
func (b *Buffer) ReadCloser() io.ReadCloser ***REMOVED***
	ret := &readCloser***REMOVED***0, append(b.bufs, b.Buf)***REMOVED***

	b.bufs = nil
	b.toPool = nil
	b.Buf = nil

	return ret
***REMOVED***
