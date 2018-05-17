package metrics

import "sync/atomic"

// Counters hold an int64 value that can be incremented and decremented.
type Counter interface ***REMOVED***
	Clear()
	Count() int64
	Dec(int64)
	Inc(int64)
	Snapshot() Counter
***REMOVED***

// GetOrRegisterCounter returns an existing Counter or constructs and registers
// a new StandardCounter.
func GetOrRegisterCounter(name string, r Registry) Counter ***REMOVED***
	if nil == r ***REMOVED***
		r = DefaultRegistry
	***REMOVED***
	return r.GetOrRegister(name, NewCounter).(Counter)
***REMOVED***

// NewCounter constructs a new StandardCounter.
func NewCounter() Counter ***REMOVED***
	if UseNilMetrics ***REMOVED***
		return NilCounter***REMOVED******REMOVED***
	***REMOVED***
	return &StandardCounter***REMOVED***0***REMOVED***
***REMOVED***

// NewRegisteredCounter constructs and registers a new StandardCounter.
func NewRegisteredCounter(name string, r Registry) Counter ***REMOVED***
	c := NewCounter()
	if nil == r ***REMOVED***
		r = DefaultRegistry
	***REMOVED***
	r.Register(name, c)
	return c
***REMOVED***

// CounterSnapshot is a read-only copy of another Counter.
type CounterSnapshot int64

// Clear panics.
func (CounterSnapshot) Clear() ***REMOVED***
	panic("Clear called on a CounterSnapshot")
***REMOVED***

// Count returns the count at the time the snapshot was taken.
func (c CounterSnapshot) Count() int64 ***REMOVED*** return int64(c) ***REMOVED***

// Dec panics.
func (CounterSnapshot) Dec(int64) ***REMOVED***
	panic("Dec called on a CounterSnapshot")
***REMOVED***

// Inc panics.
func (CounterSnapshot) Inc(int64) ***REMOVED***
	panic("Inc called on a CounterSnapshot")
***REMOVED***

// Snapshot returns the snapshot.
func (c CounterSnapshot) Snapshot() Counter ***REMOVED*** return c ***REMOVED***

// NilCounter is a no-op Counter.
type NilCounter struct***REMOVED******REMOVED***

// Clear is a no-op.
func (NilCounter) Clear() ***REMOVED******REMOVED***

// Count is a no-op.
func (NilCounter) Count() int64 ***REMOVED*** return 0 ***REMOVED***

// Dec is a no-op.
func (NilCounter) Dec(i int64) ***REMOVED******REMOVED***

// Inc is a no-op.
func (NilCounter) Inc(i int64) ***REMOVED******REMOVED***

// Snapshot is a no-op.
func (NilCounter) Snapshot() Counter ***REMOVED*** return NilCounter***REMOVED******REMOVED*** ***REMOVED***

// StandardCounter is the standard implementation of a Counter and uses the
// sync/atomic package to manage a single int64 value.
type StandardCounter struct ***REMOVED***
	count int64
***REMOVED***

// Clear sets the counter to zero.
func (c *StandardCounter) Clear() ***REMOVED***
	atomic.StoreInt64(&c.count, 0)
***REMOVED***

// Count returns the current count.
func (c *StandardCounter) Count() int64 ***REMOVED***
	return atomic.LoadInt64(&c.count)
***REMOVED***

// Dec decrements the counter by the given amount.
func (c *StandardCounter) Dec(i int64) ***REMOVED***
	atomic.AddInt64(&c.count, -i)
***REMOVED***

// Inc increments the counter by the given amount.
func (c *StandardCounter) Inc(i int64) ***REMOVED***
	atomic.AddInt64(&c.count, i)
***REMOVED***

// Snapshot returns a read-only copy of the counter.
func (c *StandardCounter) Snapshot() Counter ***REMOVED***
	return CounterSnapshot(c.Count())
***REMOVED***
