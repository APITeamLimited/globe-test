package metrics

import (
	"math"
	"sync/atomic"
)

// GaugeFloat64s hold a float64 value that can be set arbitrarily.
type GaugeFloat64 interface ***REMOVED***
	Snapshot() GaugeFloat64
	Update(float64)
	Value() float64
***REMOVED***

// GetOrRegisterGaugeFloat64 returns an existing GaugeFloat64 or constructs and registers a
// new StandardGaugeFloat64.
func GetOrRegisterGaugeFloat64(name string, r Registry) GaugeFloat64 ***REMOVED***
	if nil == r ***REMOVED***
		r = DefaultRegistry
	***REMOVED***
	return r.GetOrRegister(name, NewGaugeFloat64()).(GaugeFloat64)
***REMOVED***

// NewGaugeFloat64 constructs a new StandardGaugeFloat64.
func NewGaugeFloat64() GaugeFloat64 ***REMOVED***
	if UseNilMetrics ***REMOVED***
		return NilGaugeFloat64***REMOVED******REMOVED***
	***REMOVED***
	return &StandardGaugeFloat64***REMOVED***
		value: 0.0,
	***REMOVED***
***REMOVED***

// NewRegisteredGaugeFloat64 constructs and registers a new StandardGaugeFloat64.
func NewRegisteredGaugeFloat64(name string, r Registry) GaugeFloat64 ***REMOVED***
	c := NewGaugeFloat64()
	if nil == r ***REMOVED***
		r = DefaultRegistry
	***REMOVED***
	r.Register(name, c)
	return c
***REMOVED***

// NewFunctionalGauge constructs a new FunctionalGauge.
func NewFunctionalGaugeFloat64(f func() float64) GaugeFloat64 ***REMOVED***
	if UseNilMetrics ***REMOVED***
		return NilGaugeFloat64***REMOVED******REMOVED***
	***REMOVED***
	return &FunctionalGaugeFloat64***REMOVED***value: f***REMOVED***
***REMOVED***

// NewRegisteredFunctionalGauge constructs and registers a new StandardGauge.
func NewRegisteredFunctionalGaugeFloat64(name string, r Registry, f func() float64) GaugeFloat64 ***REMOVED***
	c := NewFunctionalGaugeFloat64(f)
	if nil == r ***REMOVED***
		r = DefaultRegistry
	***REMOVED***
	r.Register(name, c)
	return c
***REMOVED***

// GaugeFloat64Snapshot is a read-only copy of another GaugeFloat64.
type GaugeFloat64Snapshot float64

// Snapshot returns the snapshot.
func (g GaugeFloat64Snapshot) Snapshot() GaugeFloat64 ***REMOVED*** return g ***REMOVED***

// Update panics.
func (GaugeFloat64Snapshot) Update(float64) ***REMOVED***
	panic("Update called on a GaugeFloat64Snapshot")
***REMOVED***

// Value returns the value at the time the snapshot was taken.
func (g GaugeFloat64Snapshot) Value() float64 ***REMOVED*** return float64(g) ***REMOVED***

// NilGauge is a no-op Gauge.
type NilGaugeFloat64 struct***REMOVED******REMOVED***

// Snapshot is a no-op.
func (NilGaugeFloat64) Snapshot() GaugeFloat64 ***REMOVED*** return NilGaugeFloat64***REMOVED******REMOVED*** ***REMOVED***

// Update is a no-op.
func (NilGaugeFloat64) Update(v float64) ***REMOVED******REMOVED***

// Value is a no-op.
func (NilGaugeFloat64) Value() float64 ***REMOVED*** return 0.0 ***REMOVED***

// StandardGaugeFloat64 is the standard implementation of a GaugeFloat64 and uses
// sync.Mutex to manage a single float64 value.
type StandardGaugeFloat64 struct ***REMOVED***
	value uint64
***REMOVED***

// Snapshot returns a read-only copy of the gauge.
func (g *StandardGaugeFloat64) Snapshot() GaugeFloat64 ***REMOVED***
	return GaugeFloat64Snapshot(g.Value())
***REMOVED***

// Update updates the gauge's value.
func (g *StandardGaugeFloat64) Update(v float64) ***REMOVED***
	atomic.StoreUint64(&g.value, math.Float64bits(v))
***REMOVED***

// Value returns the gauge's current value.
func (g *StandardGaugeFloat64) Value() float64 ***REMOVED***
	return math.Float64frombits(atomic.LoadUint64(&g.value))
***REMOVED***

// FunctionalGaugeFloat64 returns value from given function
type FunctionalGaugeFloat64 struct ***REMOVED***
	value func() float64
***REMOVED***

// Value returns the gauge's current value.
func (g FunctionalGaugeFloat64) Value() float64 ***REMOVED***
	return g.value()
***REMOVED***

// Snapshot returns the snapshot.
func (g FunctionalGaugeFloat64) Snapshot() GaugeFloat64 ***REMOVED*** return GaugeFloat64Snapshot(g.Value()) ***REMOVED***

// Update panics.
func (FunctionalGaugeFloat64) Update(float64) ***REMOVED***
	panic("Update called on a FunctionalGaugeFloat64")
***REMOVED***
