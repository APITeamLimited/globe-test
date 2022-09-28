package bpool

import (
	"bytes"
)

// BufferPool implements a pool of bytes.Buffers in the form of a bounded
// channel.
type BufferPool struct ***REMOVED***
	c chan *bytes.Buffer
***REMOVED***

// NewBufferPool creates a new BufferPool bounded to the given size.
func NewBufferPool(size int) (bp *BufferPool) ***REMOVED***
	return &BufferPool***REMOVED***
		c: make(chan *bytes.Buffer, size),
	***REMOVED***
***REMOVED***

// Get gets a Buffer from the BufferPool, or creates a new one if none are
// available in the pool.
func (bp *BufferPool) Get() (b *bytes.Buffer) ***REMOVED***
	select ***REMOVED***
	case b = <-bp.c:
	// reuse existing buffer
	default:
		// create new buffer
		b = bytes.NewBuffer([]byte***REMOVED******REMOVED***)
	***REMOVED***
	return
***REMOVED***

// Put returns the given Buffer to the BufferPool.
func (bp *BufferPool) Put(b *bytes.Buffer) ***REMOVED***
	b.Reset()
	select ***REMOVED***
	case bp.c <- b:
	default: // Discard the buffer if the pool is full.
	***REMOVED***
***REMOVED***

// NumPooled returns the number of items currently pooled.
func (bp *BufferPool) NumPooled() int ***REMOVED***
	return len(bp.c)
***REMOVED***
