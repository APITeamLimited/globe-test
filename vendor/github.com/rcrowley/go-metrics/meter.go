package metrics

import (
	"math"
	"sync"
	"sync/atomic"
	"time"
)

// Meters count events to produce exponentially-weighted moving average rates
// at one-, five-, and fifteen-minutes and a mean rate.
type Meter interface ***REMOVED***
	Count() int64
	Mark(int64)
	Rate1() float64
	Rate5() float64
	Rate15() float64
	RateMean() float64
	Snapshot() Meter
	Stop()
***REMOVED***

// GetOrRegisterMeter returns an existing Meter or constructs and registers a
// new StandardMeter.
// Be sure to unregister the meter from the registry once it is of no use to
// allow for garbage collection.
func GetOrRegisterMeter(name string, r Registry) Meter ***REMOVED***
	if nil == r ***REMOVED***
		r = DefaultRegistry
	***REMOVED***
	return r.GetOrRegister(name, NewMeter).(Meter)
***REMOVED***

// NewMeter constructs a new StandardMeter and launches a goroutine.
// Be sure to call Stop() once the meter is of no use to allow for garbage collection.
func NewMeter() Meter ***REMOVED***
	if UseNilMetrics ***REMOVED***
		return NilMeter***REMOVED******REMOVED***
	***REMOVED***
	m := newStandardMeter()
	arbiter.Lock()
	defer arbiter.Unlock()
	arbiter.meters[m] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	if !arbiter.started ***REMOVED***
		arbiter.started = true
		go arbiter.tick()
	***REMOVED***
	return m
***REMOVED***

// NewMeter constructs and registers a new StandardMeter and launches a
// goroutine.
// Be sure to unregister the meter from the registry once it is of no use to
// allow for garbage collection.
func NewRegisteredMeter(name string, r Registry) Meter ***REMOVED***
	c := NewMeter()
	if nil == r ***REMOVED***
		r = DefaultRegistry
	***REMOVED***
	r.Register(name, c)
	return c
***REMOVED***

// MeterSnapshot is a read-only copy of another Meter.
type MeterSnapshot struct ***REMOVED***
	count                          int64
	rate1, rate5, rate15, rateMean uint64
***REMOVED***

// Count returns the count of events at the time the snapshot was taken.
func (m *MeterSnapshot) Count() int64 ***REMOVED*** return m.count ***REMOVED***

// Mark panics.
func (*MeterSnapshot) Mark(n int64) ***REMOVED***
	panic("Mark called on a MeterSnapshot")
***REMOVED***

// Rate1 returns the one-minute moving average rate of events per second at the
// time the snapshot was taken.
func (m *MeterSnapshot) Rate1() float64 ***REMOVED*** return math.Float64frombits(m.rate1) ***REMOVED***

// Rate5 returns the five-minute moving average rate of events per second at
// the time the snapshot was taken.
func (m *MeterSnapshot) Rate5() float64 ***REMOVED*** return math.Float64frombits(m.rate5) ***REMOVED***

// Rate15 returns the fifteen-minute moving average rate of events per second
// at the time the snapshot was taken.
func (m *MeterSnapshot) Rate15() float64 ***REMOVED*** return math.Float64frombits(m.rate15) ***REMOVED***

// RateMean returns the meter's mean rate of events per second at the time the
// snapshot was taken.
func (m *MeterSnapshot) RateMean() float64 ***REMOVED*** return math.Float64frombits(m.rateMean) ***REMOVED***

// Snapshot returns the snapshot.
func (m *MeterSnapshot) Snapshot() Meter ***REMOVED*** return m ***REMOVED***

// Stop is a no-op.
func (m *MeterSnapshot) Stop() ***REMOVED******REMOVED***

// NilMeter is a no-op Meter.
type NilMeter struct***REMOVED******REMOVED***

// Count is a no-op.
func (NilMeter) Count() int64 ***REMOVED*** return 0 ***REMOVED***

// Mark is a no-op.
func (NilMeter) Mark(n int64) ***REMOVED******REMOVED***

// Rate1 is a no-op.
func (NilMeter) Rate1() float64 ***REMOVED*** return 0.0 ***REMOVED***

// Rate5 is a no-op.
func (NilMeter) Rate5() float64 ***REMOVED*** return 0.0 ***REMOVED***

// Rate15is a no-op.
func (NilMeter) Rate15() float64 ***REMOVED*** return 0.0 ***REMOVED***

// RateMean is a no-op.
func (NilMeter) RateMean() float64 ***REMOVED*** return 0.0 ***REMOVED***

// Snapshot is a no-op.
func (NilMeter) Snapshot() Meter ***REMOVED*** return NilMeter***REMOVED******REMOVED*** ***REMOVED***

// Stop is a no-op.
func (NilMeter) Stop() ***REMOVED******REMOVED***

// StandardMeter is the standard implementation of a Meter.
type StandardMeter struct ***REMOVED***
	// Only used on stop.
	lock        sync.Mutex
	snapshot    *MeterSnapshot
	a1, a5, a15 EWMA
	startTime   time.Time
	stopped     uint32
***REMOVED***

func newStandardMeter() *StandardMeter ***REMOVED***
	return &StandardMeter***REMOVED***
		snapshot:  &MeterSnapshot***REMOVED******REMOVED***,
		a1:        NewEWMA1(),
		a5:        NewEWMA5(),
		a15:       NewEWMA15(),
		startTime: time.Now(),
	***REMOVED***
***REMOVED***

// Stop stops the meter, Mark() will be a no-op if you use it after being stopped.
func (m *StandardMeter) Stop() ***REMOVED***
	m.lock.Lock()
	stopped := m.stopped
	m.stopped = 1
	m.lock.Unlock()
	if stopped != 1 ***REMOVED***
		arbiter.Lock()
		delete(arbiter.meters, m)
		arbiter.Unlock()
	***REMOVED***
***REMOVED***

// Count returns the number of events recorded.
func (m *StandardMeter) Count() int64 ***REMOVED***
	return atomic.LoadInt64(&m.snapshot.count)
***REMOVED***

// Mark records the occurance of n events.
func (m *StandardMeter) Mark(n int64) ***REMOVED***
	if atomic.LoadUint32(&m.stopped) == 1 ***REMOVED***
		return
	***REMOVED***

	atomic.AddInt64(&m.snapshot.count, n)

	m.a1.Update(n)
	m.a5.Update(n)
	m.a15.Update(n)
	m.updateSnapshot()
***REMOVED***

// Rate1 returns the one-minute moving average rate of events per second.
func (m *StandardMeter) Rate1() float64 ***REMOVED***
	return math.Float64frombits(atomic.LoadUint64(&m.snapshot.rate1))
***REMOVED***

// Rate5 returns the five-minute moving average rate of events per second.
func (m *StandardMeter) Rate5() float64 ***REMOVED***
	return math.Float64frombits(atomic.LoadUint64(&m.snapshot.rate5))
***REMOVED***

// Rate15 returns the fifteen-minute moving average rate of events per second.
func (m *StandardMeter) Rate15() float64 ***REMOVED***
	return math.Float64frombits(atomic.LoadUint64(&m.snapshot.rate15))
***REMOVED***

// RateMean returns the meter's mean rate of events per second.
func (m *StandardMeter) RateMean() float64 ***REMOVED***
	return math.Float64frombits(atomic.LoadUint64(&m.snapshot.rateMean))
***REMOVED***

// Snapshot returns a read-only copy of the meter.
func (m *StandardMeter) Snapshot() Meter ***REMOVED***
	copiedSnapshot := MeterSnapshot***REMOVED***
		count:    atomic.LoadInt64(&m.snapshot.count),
		rate1:    atomic.LoadUint64(&m.snapshot.rate1),
		rate5:    atomic.LoadUint64(&m.snapshot.rate5),
		rate15:   atomic.LoadUint64(&m.snapshot.rate15),
		rateMean: atomic.LoadUint64(&m.snapshot.rateMean),
	***REMOVED***
	return &copiedSnapshot
***REMOVED***

func (m *StandardMeter) updateSnapshot() ***REMOVED***
	rate1 := math.Float64bits(m.a1.Rate())
	rate5 := math.Float64bits(m.a5.Rate())
	rate15 := math.Float64bits(m.a15.Rate())
	rateMean := math.Float64bits(float64(m.Count()) / time.Since(m.startTime).Seconds())

	atomic.StoreUint64(&m.snapshot.rate1, rate1)
	atomic.StoreUint64(&m.snapshot.rate5, rate5)
	atomic.StoreUint64(&m.snapshot.rate15, rate15)
	atomic.StoreUint64(&m.snapshot.rateMean, rateMean)
***REMOVED***

func (m *StandardMeter) tick() ***REMOVED***
	m.a1.Tick()
	m.a5.Tick()
	m.a15.Tick()
	m.updateSnapshot()
***REMOVED***

// meterArbiter ticks meters every 5s from a single goroutine.
// meters are references in a set for future stopping.
type meterArbiter struct ***REMOVED***
	sync.RWMutex
	started bool
	meters  map[*StandardMeter]struct***REMOVED******REMOVED***
	ticker  *time.Ticker
***REMOVED***

var arbiter = meterArbiter***REMOVED***ticker: time.NewTicker(5e9), meters: make(map[*StandardMeter]struct***REMOVED******REMOVED***)***REMOVED***

// Ticks meters on the scheduled interval
func (ma *meterArbiter) tick() ***REMOVED***
	for ***REMOVED***
		select ***REMOVED***
		case <-ma.ticker.C:
			ma.tickMeters()
		***REMOVED***
	***REMOVED***
***REMOVED***

func (ma *meterArbiter) tickMeters() ***REMOVED***
	ma.RLock()
	defer ma.RUnlock()
	for meter := range ma.meters ***REMOVED***
		meter.tick()
	***REMOVED***
***REMOVED***
