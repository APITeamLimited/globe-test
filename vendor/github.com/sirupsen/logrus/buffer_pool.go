package logrus

import (
	"bytes"
	"sync"
)

var (
	bufferPool BufferPool
)

type BufferPool interface ***REMOVED***
	Put(*bytes.Buffer)
	Get() *bytes.Buffer
***REMOVED***

type defaultPool struct ***REMOVED***
	pool *sync.Pool
***REMOVED***

func (p *defaultPool) Put(buf *bytes.Buffer) ***REMOVED***
	p.pool.Put(buf)
***REMOVED***

func (p *defaultPool) Get() *bytes.Buffer ***REMOVED***
	return p.pool.Get().(*bytes.Buffer)
***REMOVED***

func getBuffer() *bytes.Buffer ***REMOVED***
	return bufferPool.Get()
***REMOVED***

func putBuffer(buf *bytes.Buffer) ***REMOVED***
	buf.Reset()
	bufferPool.Put(buf)
***REMOVED***

// SetBufferPool allows to replace the default logrus buffer pool
// to better meets the specific needs of an application.
func SetBufferPool(bp BufferPool) ***REMOVED***
	bufferPool = bp
***REMOVED***

func init() ***REMOVED***
	SetBufferPool(&defaultPool***REMOVED***
		pool: &sync.Pool***REMOVED***
			New: func() interface***REMOVED******REMOVED*** ***REMOVED***
				return new(bytes.Buffer)
			***REMOVED***,
		***REMOVED***,
	***REMOVED***)
***REMOVED***
