package metrics

import "sync/atomic"

// Gauges hold an int64 value that can be set arbitrarily.
type Gauge interface ***REMOVED***
	Snapshot() Gauge
	Update(int64)
	Value() int64
***REMOVED***

// GetOrRegisterGauge returns an existing Gauge or constructs and registers a
// new StandardGauge.
func GetOrRegisterGauge(name string, r Registry) Gauge ***REMOVED***
	if nil == r ***REMOVED***
		r = DefaultRegistry
	***REMOVED***
	return r.GetOrRegister(name, NewGauge).(Gauge)
***REMOVED***

// NewGauge constructs a new StandardGauge.
func NewGauge() Gauge ***REMOVED***
	if UseNilMetrics ***REMOVED***
		return NilGauge***REMOVED******REMOVED***
	***REMOVED***
	return &StandardGauge***REMOVED***0***REMOVED***
***REMOVED***

// NewRegisteredGauge constructs and registers a new StandardGauge.
func NewRegisteredGauge(name string, r Registry) Gauge ***REMOVED***
	c := NewGauge()
	if nil == r ***REMOVED***
		r = DefaultRegistry
	***REMOVED***
	r.Register(name, c)
	return c
***REMOVED***

// NewFunctionalGauge constructs a new FunctionalGauge.
func NewFunctionalGauge(f func() int64) Gauge ***REMOVED***
	if UseNilMetrics ***REMOVED***
		return NilGauge***REMOVED******REMOVED***
	***REMOVED***
	return &FunctionalGauge***REMOVED***value: f***REMOVED***
***REMOVED***

// NewRegisteredFunctionalGauge constructs and registers a new StandardGauge.
func NewRegisteredFunctionalGauge(name string, r Registry, f func() int64) Gauge ***REMOVED***
	c := NewFunctionalGauge(f)
	if nil == r ***REMOVED***
		r = DefaultRegistry
	***REMOVED***
	r.Register(name, c)
	return c
***REMOVED***

// GaugeSnapshot is a read-only copy of another Gauge.
type GaugeSnapshot int64

// Snapshot returns the snapshot.
func (g GaugeSnapshot) Snapshot() Gauge ***REMOVED*** return g ***REMOVED***

// Update panics.
func (GaugeSnapshot) Update(int64) ***REMOVED***
	panic("Update called on a GaugeSnapshot")
***REMOVED***

// Value returns the value at the time the snapshot was taken.
func (g GaugeSnapshot) Value() int64 ***REMOVED*** return int64(g) ***REMOVED***

// NilGauge is a no-op Gauge.
type NilGauge struct***REMOVED******REMOVED***

// Snapshot is a no-op.
func (NilGauge) Snapshot() Gauge ***REMOVED*** return NilGauge***REMOVED******REMOVED*** ***REMOVED***

// Update is a no-op.
func (NilGauge) Update(v int64) ***REMOVED******REMOVED***

// Value is a no-op.
func (NilGauge) Value() int64 ***REMOVED*** return 0 ***REMOVED***

// StandardGauge is the standard implementation of a Gauge and uses the
// sync/atomic package to manage a single int64 value.
type StandardGauge struct ***REMOVED***
	value int64
***REMOVED***

// Snapshot returns a read-only copy of the gauge.
func (g *StandardGauge) Snapshot() Gauge ***REMOVED***
	return GaugeSnapshot(g.Value())
***REMOVED***

// Update updates the gauge's value.
func (g *StandardGauge) Update(v int64) ***REMOVED***
	atomic.StoreInt64(&g.value, v)
***REMOVED***

// Value returns the gauge's current value.
func (g *StandardGauge) Value() int64 ***REMOVED***
	return atomic.LoadInt64(&g.value)
***REMOVED***

// FunctionalGauge returns value from given function
type FunctionalGauge struct ***REMOVED***
	value func() int64
***REMOVED***

// Value returns the gauge's current value.
func (g FunctionalGauge) Value() int64 ***REMOVED***
	return g.value()
***REMOVED***

// Snapshot returns the snapshot.
func (g FunctionalGauge) Snapshot() Gauge ***REMOVED*** return GaugeSnapshot(g.Value()) ***REMOVED***

// Update panics.
func (FunctionalGauge) Update(int64) ***REMOVED***
	panic("Update called on a FunctionalGauge")
***REMOVED***
