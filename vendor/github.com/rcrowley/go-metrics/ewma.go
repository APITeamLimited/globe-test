package metrics

import (
	"math"
	"sync"
	"sync/atomic"
)

// EWMAs continuously calculate an exponentially-weighted moving average
// based on an outside source of clock ticks.
type EWMA interface ***REMOVED***
	Rate() float64
	Snapshot() EWMA
	Tick()
	Update(int64)
***REMOVED***

// NewEWMA constructs a new EWMA with the given alpha.
func NewEWMA(alpha float64) EWMA ***REMOVED***
	if UseNilMetrics ***REMOVED***
		return NilEWMA***REMOVED******REMOVED***
	***REMOVED***
	return &StandardEWMA***REMOVED***alpha: alpha***REMOVED***
***REMOVED***

// NewEWMA1 constructs a new EWMA for a one-minute moving average.
func NewEWMA1() EWMA ***REMOVED***
	return NewEWMA(1 - math.Exp(-5.0/60.0/1))
***REMOVED***

// NewEWMA5 constructs a new EWMA for a five-minute moving average.
func NewEWMA5() EWMA ***REMOVED***
	return NewEWMA(1 - math.Exp(-5.0/60.0/5))
***REMOVED***

// NewEWMA15 constructs a new EWMA for a fifteen-minute moving average.
func NewEWMA15() EWMA ***REMOVED***
	return NewEWMA(1 - math.Exp(-5.0/60.0/15))
***REMOVED***

// EWMASnapshot is a read-only copy of another EWMA.
type EWMASnapshot float64

// Rate returns the rate of events per second at the time the snapshot was
// taken.
func (a EWMASnapshot) Rate() float64 ***REMOVED*** return float64(a) ***REMOVED***

// Snapshot returns the snapshot.
func (a EWMASnapshot) Snapshot() EWMA ***REMOVED*** return a ***REMOVED***

// Tick panics.
func (EWMASnapshot) Tick() ***REMOVED***
	panic("Tick called on an EWMASnapshot")
***REMOVED***

// Update panics.
func (EWMASnapshot) Update(int64) ***REMOVED***
	panic("Update called on an EWMASnapshot")
***REMOVED***

// NilEWMA is a no-op EWMA.
type NilEWMA struct***REMOVED******REMOVED***

// Rate is a no-op.
func (NilEWMA) Rate() float64 ***REMOVED*** return 0.0 ***REMOVED***

// Snapshot is a no-op.
func (NilEWMA) Snapshot() EWMA ***REMOVED*** return NilEWMA***REMOVED******REMOVED*** ***REMOVED***

// Tick is a no-op.
func (NilEWMA) Tick() ***REMOVED******REMOVED***

// Update is a no-op.
func (NilEWMA) Update(n int64) ***REMOVED******REMOVED***

// StandardEWMA is the standard implementation of an EWMA and tracks the number
// of uncounted events and processes them on each tick.  It uses the
// sync/atomic package to manage uncounted events.
type StandardEWMA struct ***REMOVED***
	uncounted int64 // /!\ this should be the first member to ensure 64-bit alignment
	alpha     float64
	rate      uint64
	init      uint32
	mutex     sync.Mutex
***REMOVED***

// Rate returns the moving average rate of events per second.
func (a *StandardEWMA) Rate() float64 ***REMOVED***
	currentRate := math.Float64frombits(atomic.LoadUint64(&a.rate)) * float64(1e9)
	return currentRate
***REMOVED***

// Snapshot returns a read-only copy of the EWMA.
func (a *StandardEWMA) Snapshot() EWMA ***REMOVED***
	return EWMASnapshot(a.Rate())
***REMOVED***

// Tick ticks the clock to update the moving average.  It assumes it is called
// every five seconds.
func (a *StandardEWMA) Tick() ***REMOVED***
	// Optimization to avoid mutex locking in the hot-path.
	if atomic.LoadUint32(&a.init) == 1 ***REMOVED***
		a.updateRate(a.fetchInstantRate())
	***REMOVED*** else ***REMOVED***
		// Slow-path: this is only needed on the first Tick() and preserves transactional updating
		// of init and rate in the else block. The first conditional is needed below because
		// a different thread could have set a.init = 1 between the time of the first atomic load and when
		// the lock was acquired.
		a.mutex.Lock()
		if atomic.LoadUint32(&a.init) == 1 ***REMOVED***
			// The fetchInstantRate() uses atomic loading, which is unecessary in this critical section
			// but again, this section is only invoked on the first successful Tick() operation.
			a.updateRate(a.fetchInstantRate())
		***REMOVED*** else ***REMOVED***
			atomic.StoreUint32(&a.init, 1)
			atomic.StoreUint64(&a.rate, math.Float64bits(a.fetchInstantRate()))
		***REMOVED***
		a.mutex.Unlock()
	***REMOVED***
***REMOVED***

func (a *StandardEWMA) fetchInstantRate() float64 ***REMOVED***
	count := atomic.LoadInt64(&a.uncounted)
	atomic.AddInt64(&a.uncounted, -count)
	instantRate := float64(count) / float64(5e9)
	return instantRate
***REMOVED***

func (a *StandardEWMA) updateRate(instantRate float64) ***REMOVED***
	currentRate := math.Float64frombits(atomic.LoadUint64(&a.rate))
	currentRate += a.alpha * (instantRate - currentRate)
	atomic.StoreUint64(&a.rate, math.Float64bits(currentRate))
***REMOVED***

// Update adds n uncounted events.
func (a *StandardEWMA) Update(n int64) ***REMOVED***
	atomic.AddInt64(&a.uncounted, n)
***REMOVED***
