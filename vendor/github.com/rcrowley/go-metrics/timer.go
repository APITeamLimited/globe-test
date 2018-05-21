package metrics

import (
	"sync"
	"time"
)

// Timers capture the duration and rate of events.
type Timer interface ***REMOVED***
	Count() int64
	Max() int64
	Mean() float64
	Min() int64
	Percentile(float64) float64
	Percentiles([]float64) []float64
	Rate1() float64
	Rate5() float64
	Rate15() float64
	RateMean() float64
	Snapshot() Timer
	StdDev() float64
	Stop()
	Sum() int64
	Time(func())
	Update(time.Duration)
	UpdateSince(time.Time)
	Variance() float64
***REMOVED***

// GetOrRegisterTimer returns an existing Timer or constructs and registers a
// new StandardTimer.
// Be sure to unregister the meter from the registry once it is of no use to
// allow for garbage collection.
func GetOrRegisterTimer(name string, r Registry) Timer ***REMOVED***
	if nil == r ***REMOVED***
		r = DefaultRegistry
	***REMOVED***
	return r.GetOrRegister(name, NewTimer).(Timer)
***REMOVED***

// NewCustomTimer constructs a new StandardTimer from a Histogram and a Meter.
// Be sure to call Stop() once the timer is of no use to allow for garbage collection.
func NewCustomTimer(h Histogram, m Meter) Timer ***REMOVED***
	if UseNilMetrics ***REMOVED***
		return NilTimer***REMOVED******REMOVED***
	***REMOVED***
	return &StandardTimer***REMOVED***
		histogram: h,
		meter:     m,
	***REMOVED***
***REMOVED***

// NewRegisteredTimer constructs and registers a new StandardTimer.
// Be sure to unregister the meter from the registry once it is of no use to
// allow for garbage collection.
func NewRegisteredTimer(name string, r Registry) Timer ***REMOVED***
	c := NewTimer()
	if nil == r ***REMOVED***
		r = DefaultRegistry
	***REMOVED***
	r.Register(name, c)
	return c
***REMOVED***

// NewTimer constructs a new StandardTimer using an exponentially-decaying
// sample with the same reservoir size and alpha as UNIX load averages.
// Be sure to call Stop() once the timer is of no use to allow for garbage collection.
func NewTimer() Timer ***REMOVED***
	if UseNilMetrics ***REMOVED***
		return NilTimer***REMOVED******REMOVED***
	***REMOVED***
	return &StandardTimer***REMOVED***
		histogram: NewHistogram(NewExpDecaySample(1028, 0.015)),
		meter:     NewMeter(),
	***REMOVED***
***REMOVED***

// NilTimer is a no-op Timer.
type NilTimer struct ***REMOVED***
	h Histogram
	m Meter
***REMOVED***

// Count is a no-op.
func (NilTimer) Count() int64 ***REMOVED*** return 0 ***REMOVED***

// Max is a no-op.
func (NilTimer) Max() int64 ***REMOVED*** return 0 ***REMOVED***

// Mean is a no-op.
func (NilTimer) Mean() float64 ***REMOVED*** return 0.0 ***REMOVED***

// Min is a no-op.
func (NilTimer) Min() int64 ***REMOVED*** return 0 ***REMOVED***

// Percentile is a no-op.
func (NilTimer) Percentile(p float64) float64 ***REMOVED*** return 0.0 ***REMOVED***

// Percentiles is a no-op.
func (NilTimer) Percentiles(ps []float64) []float64 ***REMOVED***
	return make([]float64, len(ps))
***REMOVED***

// Rate1 is a no-op.
func (NilTimer) Rate1() float64 ***REMOVED*** return 0.0 ***REMOVED***

// Rate5 is a no-op.
func (NilTimer) Rate5() float64 ***REMOVED*** return 0.0 ***REMOVED***

// Rate15 is a no-op.
func (NilTimer) Rate15() float64 ***REMOVED*** return 0.0 ***REMOVED***

// RateMean is a no-op.
func (NilTimer) RateMean() float64 ***REMOVED*** return 0.0 ***REMOVED***

// Snapshot is a no-op.
func (NilTimer) Snapshot() Timer ***REMOVED*** return NilTimer***REMOVED******REMOVED*** ***REMOVED***

// StdDev is a no-op.
func (NilTimer) StdDev() float64 ***REMOVED*** return 0.0 ***REMOVED***

// Stop is a no-op.
func (NilTimer) Stop() ***REMOVED******REMOVED***

// Sum is a no-op.
func (NilTimer) Sum() int64 ***REMOVED*** return 0 ***REMOVED***

// Time is a no-op.
func (NilTimer) Time(func()) ***REMOVED******REMOVED***

// Update is a no-op.
func (NilTimer) Update(time.Duration) ***REMOVED******REMOVED***

// UpdateSince is a no-op.
func (NilTimer) UpdateSince(time.Time) ***REMOVED******REMOVED***

// Variance is a no-op.
func (NilTimer) Variance() float64 ***REMOVED*** return 0.0 ***REMOVED***

// StandardTimer is the standard implementation of a Timer and uses a Histogram
// and Meter.
type StandardTimer struct ***REMOVED***
	histogram Histogram
	meter     Meter
	mutex     sync.Mutex
***REMOVED***

// Count returns the number of events recorded.
func (t *StandardTimer) Count() int64 ***REMOVED***
	return t.histogram.Count()
***REMOVED***

// Max returns the maximum value in the sample.
func (t *StandardTimer) Max() int64 ***REMOVED***
	return t.histogram.Max()
***REMOVED***

// Mean returns the mean of the values in the sample.
func (t *StandardTimer) Mean() float64 ***REMOVED***
	return t.histogram.Mean()
***REMOVED***

// Min returns the minimum value in the sample.
func (t *StandardTimer) Min() int64 ***REMOVED***
	return t.histogram.Min()
***REMOVED***

// Percentile returns an arbitrary percentile of the values in the sample.
func (t *StandardTimer) Percentile(p float64) float64 ***REMOVED***
	return t.histogram.Percentile(p)
***REMOVED***

// Percentiles returns a slice of arbitrary percentiles of the values in the
// sample.
func (t *StandardTimer) Percentiles(ps []float64) []float64 ***REMOVED***
	return t.histogram.Percentiles(ps)
***REMOVED***

// Rate1 returns the one-minute moving average rate of events per second.
func (t *StandardTimer) Rate1() float64 ***REMOVED***
	return t.meter.Rate1()
***REMOVED***

// Rate5 returns the five-minute moving average rate of events per second.
func (t *StandardTimer) Rate5() float64 ***REMOVED***
	return t.meter.Rate5()
***REMOVED***

// Rate15 returns the fifteen-minute moving average rate of events per second.
func (t *StandardTimer) Rate15() float64 ***REMOVED***
	return t.meter.Rate15()
***REMOVED***

// RateMean returns the meter's mean rate of events per second.
func (t *StandardTimer) RateMean() float64 ***REMOVED***
	return t.meter.RateMean()
***REMOVED***

// Snapshot returns a read-only copy of the timer.
func (t *StandardTimer) Snapshot() Timer ***REMOVED***
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return &TimerSnapshot***REMOVED***
		histogram: t.histogram.Snapshot().(*HistogramSnapshot),
		meter:     t.meter.Snapshot().(*MeterSnapshot),
	***REMOVED***
***REMOVED***

// StdDev returns the standard deviation of the values in the sample.
func (t *StandardTimer) StdDev() float64 ***REMOVED***
	return t.histogram.StdDev()
***REMOVED***

// Stop stops the meter.
func (t *StandardTimer) Stop() ***REMOVED***
	t.meter.Stop()
***REMOVED***

// Sum returns the sum in the sample.
func (t *StandardTimer) Sum() int64 ***REMOVED***
	return t.histogram.Sum()
***REMOVED***

// Record the duration of the execution of the given function.
func (t *StandardTimer) Time(f func()) ***REMOVED***
	ts := time.Now()
	f()
	t.Update(time.Since(ts))
***REMOVED***

// Record the duration of an event.
func (t *StandardTimer) Update(d time.Duration) ***REMOVED***
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.histogram.Update(int64(d))
	t.meter.Mark(1)
***REMOVED***

// Record the duration of an event that started at a time and ends now.
func (t *StandardTimer) UpdateSince(ts time.Time) ***REMOVED***
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.histogram.Update(int64(time.Since(ts)))
	t.meter.Mark(1)
***REMOVED***

// Variance returns the variance of the values in the sample.
func (t *StandardTimer) Variance() float64 ***REMOVED***
	return t.histogram.Variance()
***REMOVED***

// TimerSnapshot is a read-only copy of another Timer.
type TimerSnapshot struct ***REMOVED***
	histogram *HistogramSnapshot
	meter     *MeterSnapshot
***REMOVED***

// Count returns the number of events recorded at the time the snapshot was
// taken.
func (t *TimerSnapshot) Count() int64 ***REMOVED*** return t.histogram.Count() ***REMOVED***

// Max returns the maximum value at the time the snapshot was taken.
func (t *TimerSnapshot) Max() int64 ***REMOVED*** return t.histogram.Max() ***REMOVED***

// Mean returns the mean value at the time the snapshot was taken.
func (t *TimerSnapshot) Mean() float64 ***REMOVED*** return t.histogram.Mean() ***REMOVED***

// Min returns the minimum value at the time the snapshot was taken.
func (t *TimerSnapshot) Min() int64 ***REMOVED*** return t.histogram.Min() ***REMOVED***

// Percentile returns an arbitrary percentile of sampled values at the time the
// snapshot was taken.
func (t *TimerSnapshot) Percentile(p float64) float64 ***REMOVED***
	return t.histogram.Percentile(p)
***REMOVED***

// Percentiles returns a slice of arbitrary percentiles of sampled values at
// the time the snapshot was taken.
func (t *TimerSnapshot) Percentiles(ps []float64) []float64 ***REMOVED***
	return t.histogram.Percentiles(ps)
***REMOVED***

// Rate1 returns the one-minute moving average rate of events per second at the
// time the snapshot was taken.
func (t *TimerSnapshot) Rate1() float64 ***REMOVED*** return t.meter.Rate1() ***REMOVED***

// Rate5 returns the five-minute moving average rate of events per second at
// the time the snapshot was taken.
func (t *TimerSnapshot) Rate5() float64 ***REMOVED*** return t.meter.Rate5() ***REMOVED***

// Rate15 returns the fifteen-minute moving average rate of events per second
// at the time the snapshot was taken.
func (t *TimerSnapshot) Rate15() float64 ***REMOVED*** return t.meter.Rate15() ***REMOVED***

// RateMean returns the meter's mean rate of events per second at the time the
// snapshot was taken.
func (t *TimerSnapshot) RateMean() float64 ***REMOVED*** return t.meter.RateMean() ***REMOVED***

// Snapshot returns the snapshot.
func (t *TimerSnapshot) Snapshot() Timer ***REMOVED*** return t ***REMOVED***

// StdDev returns the standard deviation of the values at the time the snapshot
// was taken.
func (t *TimerSnapshot) StdDev() float64 ***REMOVED*** return t.histogram.StdDev() ***REMOVED***

// Stop is a no-op.
func (t *TimerSnapshot) Stop() ***REMOVED******REMOVED***

// Sum returns the sum at the time the snapshot was taken.
func (t *TimerSnapshot) Sum() int64 ***REMOVED*** return t.histogram.Sum() ***REMOVED***

// Time panics.
func (*TimerSnapshot) Time(func()) ***REMOVED***
	panic("Time called on a TimerSnapshot")
***REMOVED***

// Update panics.
func (*TimerSnapshot) Update(time.Duration) ***REMOVED***
	panic("Update called on a TimerSnapshot")
***REMOVED***

// UpdateSince panics.
func (*TimerSnapshot) UpdateSince(time.Time) ***REMOVED***
	panic("UpdateSince called on a TimerSnapshot")
***REMOVED***

// Variance returns the variance of the values at the time the snapshot was
// taken.
func (t *TimerSnapshot) Variance() float64 ***REMOVED*** return t.histogram.Variance() ***REMOVED***
