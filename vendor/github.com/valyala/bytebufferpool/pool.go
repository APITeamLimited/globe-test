package bytebufferpool

import (
	"sort"
	"sync"
	"sync/atomic"
)

const (
	minBitSize = 6 // 2**6=64 is a CPU cache line size
	steps      = 20

	minSize = 1 << minBitSize
	maxSize = 1 << (minBitSize + steps - 1)

	calibrateCallsThreshold = 42000
	maxPercentile           = 0.95
)

// Pool represents byte buffer pool.
//
// Distinct pools may be used for distinct types of byte buffers.
// Properly determined byte buffer types with their own pools may help reducing
// memory waste.
type Pool struct ***REMOVED***
	calls       [steps]uint64
	calibrating uint64

	defaultSize uint64
	maxSize     uint64

	pool sync.Pool
***REMOVED***

var defaultPool Pool

// Get returns an empty byte buffer from the pool.
//
// Got byte buffer may be returned to the pool via Put call.
// This reduces the number of memory allocations required for byte buffer
// management.
func Get() *ByteBuffer ***REMOVED*** return defaultPool.Get() ***REMOVED***

// Get returns new byte buffer with zero length.
//
// The byte buffer may be returned to the pool via Put after the use
// in order to minimize GC overhead.
func (p *Pool) Get() *ByteBuffer ***REMOVED***
	v := p.pool.Get()
	if v != nil ***REMOVED***
		return v.(*ByteBuffer)
	***REMOVED***
	return &ByteBuffer***REMOVED***
		B: make([]byte, 0, atomic.LoadUint64(&p.defaultSize)),
	***REMOVED***
***REMOVED***

// Put returns byte buffer to the pool.
//
// ByteBuffer.B mustn't be touched after returning it to the pool.
// Otherwise data races will occur.
func Put(b *ByteBuffer) ***REMOVED*** defaultPool.Put(b) ***REMOVED***

// Put releases byte buffer obtained via Get to the pool.
//
// The buffer mustn't be accessed after returning to the pool.
func (p *Pool) Put(b *ByteBuffer) ***REMOVED***
	idx := index(len(b.B))

	if atomic.AddUint64(&p.calls[idx], 1) > calibrateCallsThreshold ***REMOVED***
		p.calibrate()
	***REMOVED***

	maxSize := int(atomic.LoadUint64(&p.maxSize))
	if maxSize == 0 || cap(b.B) <= maxSize ***REMOVED***
		b.Reset()
		p.pool.Put(b)
	***REMOVED***
***REMOVED***

func (p *Pool) calibrate() ***REMOVED***
	if !atomic.CompareAndSwapUint64(&p.calibrating, 0, 1) ***REMOVED***
		return
	***REMOVED***

	a := make(callSizes, 0, steps)
	var callsSum uint64
	for i := uint64(0); i < steps; i++ ***REMOVED***
		calls := atomic.SwapUint64(&p.calls[i], 0)
		callsSum += calls
		a = append(a, callSize***REMOVED***
			calls: calls,
			size:  minSize << i,
		***REMOVED***)
	***REMOVED***
	sort.Sort(a)

	defaultSize := a[0].size
	maxSize := defaultSize

	maxSum := uint64(float64(callsSum) * maxPercentile)
	callsSum = 0
	for i := 0; i < steps; i++ ***REMOVED***
		if callsSum > maxSum ***REMOVED***
			break
		***REMOVED***
		callsSum += a[i].calls
		size := a[i].size
		if size > maxSize ***REMOVED***
			maxSize = size
		***REMOVED***
	***REMOVED***

	atomic.StoreUint64(&p.defaultSize, defaultSize)
	atomic.StoreUint64(&p.maxSize, maxSize)

	atomic.StoreUint64(&p.calibrating, 0)
***REMOVED***

type callSize struct ***REMOVED***
	calls uint64
	size  uint64
***REMOVED***

type callSizes []callSize

func (ci callSizes) Len() int ***REMOVED***
	return len(ci)
***REMOVED***

func (ci callSizes) Less(i, j int) bool ***REMOVED***
	return ci[i].calls > ci[j].calls
***REMOVED***

func (ci callSizes) Swap(i, j int) ***REMOVED***
	ci[i], ci[j] = ci[j], ci[i]
***REMOVED***

func index(n int) int ***REMOVED***
	n--
	n >>= minBitSize
	idx := 0
	for n > 0 ***REMOVED***
		n >>= 1
		idx++
	***REMOVED***
	if idx >= steps ***REMOVED***
		idx = steps - 1
	***REMOVED***
	return idx
***REMOVED***
