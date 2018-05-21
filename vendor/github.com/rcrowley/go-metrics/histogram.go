package metrics

// Histograms calculate distribution statistics from a series of int64 values.
type Histogram interface ***REMOVED***
	Clear()
	Count() int64
	Max() int64
	Mean() float64
	Min() int64
	Percentile(float64) float64
	Percentiles([]float64) []float64
	Sample() Sample
	Snapshot() Histogram
	StdDev() float64
	Sum() int64
	Update(int64)
	Variance() float64
***REMOVED***

// GetOrRegisterHistogram returns an existing Histogram or constructs and
// registers a new StandardHistogram.
func GetOrRegisterHistogram(name string, r Registry, s Sample) Histogram ***REMOVED***
	if nil == r ***REMOVED***
		r = DefaultRegistry
	***REMOVED***
	return r.GetOrRegister(name, func() Histogram ***REMOVED*** return NewHistogram(s) ***REMOVED***).(Histogram)
***REMOVED***

// NewHistogram constructs a new StandardHistogram from a Sample.
func NewHistogram(s Sample) Histogram ***REMOVED***
	if UseNilMetrics ***REMOVED***
		return NilHistogram***REMOVED******REMOVED***
	***REMOVED***
	return &StandardHistogram***REMOVED***sample: s***REMOVED***
***REMOVED***

// NewRegisteredHistogram constructs and registers a new StandardHistogram from
// a Sample.
func NewRegisteredHistogram(name string, r Registry, s Sample) Histogram ***REMOVED***
	c := NewHistogram(s)
	if nil == r ***REMOVED***
		r = DefaultRegistry
	***REMOVED***
	r.Register(name, c)
	return c
***REMOVED***

// HistogramSnapshot is a read-only copy of another Histogram.
type HistogramSnapshot struct ***REMOVED***
	sample *SampleSnapshot
***REMOVED***

// Clear panics.
func (*HistogramSnapshot) Clear() ***REMOVED***
	panic("Clear called on a HistogramSnapshot")
***REMOVED***

// Count returns the number of samples recorded at the time the snapshot was
// taken.
func (h *HistogramSnapshot) Count() int64 ***REMOVED*** return h.sample.Count() ***REMOVED***

// Max returns the maximum value in the sample at the time the snapshot was
// taken.
func (h *HistogramSnapshot) Max() int64 ***REMOVED*** return h.sample.Max() ***REMOVED***

// Mean returns the mean of the values in the sample at the time the snapshot
// was taken.
func (h *HistogramSnapshot) Mean() float64 ***REMOVED*** return h.sample.Mean() ***REMOVED***

// Min returns the minimum value in the sample at the time the snapshot was
// taken.
func (h *HistogramSnapshot) Min() int64 ***REMOVED*** return h.sample.Min() ***REMOVED***

// Percentile returns an arbitrary percentile of values in the sample at the
// time the snapshot was taken.
func (h *HistogramSnapshot) Percentile(p float64) float64 ***REMOVED***
	return h.sample.Percentile(p)
***REMOVED***

// Percentiles returns a slice of arbitrary percentiles of values in the sample
// at the time the snapshot was taken.
func (h *HistogramSnapshot) Percentiles(ps []float64) []float64 ***REMOVED***
	return h.sample.Percentiles(ps)
***REMOVED***

// Sample returns the Sample underlying the histogram.
func (h *HistogramSnapshot) Sample() Sample ***REMOVED*** return h.sample ***REMOVED***

// Snapshot returns the snapshot.
func (h *HistogramSnapshot) Snapshot() Histogram ***REMOVED*** return h ***REMOVED***

// StdDev returns the standard deviation of the values in the sample at the
// time the snapshot was taken.
func (h *HistogramSnapshot) StdDev() float64 ***REMOVED*** return h.sample.StdDev() ***REMOVED***

// Sum returns the sum in the sample at the time the snapshot was taken.
func (h *HistogramSnapshot) Sum() int64 ***REMOVED*** return h.sample.Sum() ***REMOVED***

// Update panics.
func (*HistogramSnapshot) Update(int64) ***REMOVED***
	panic("Update called on a HistogramSnapshot")
***REMOVED***

// Variance returns the variance of inputs at the time the snapshot was taken.
func (h *HistogramSnapshot) Variance() float64 ***REMOVED*** return h.sample.Variance() ***REMOVED***

// NilHistogram is a no-op Histogram.
type NilHistogram struct***REMOVED******REMOVED***

// Clear is a no-op.
func (NilHistogram) Clear() ***REMOVED******REMOVED***

// Count is a no-op.
func (NilHistogram) Count() int64 ***REMOVED*** return 0 ***REMOVED***

// Max is a no-op.
func (NilHistogram) Max() int64 ***REMOVED*** return 0 ***REMOVED***

// Mean is a no-op.
func (NilHistogram) Mean() float64 ***REMOVED*** return 0.0 ***REMOVED***

// Min is a no-op.
func (NilHistogram) Min() int64 ***REMOVED*** return 0 ***REMOVED***

// Percentile is a no-op.
func (NilHistogram) Percentile(p float64) float64 ***REMOVED*** return 0.0 ***REMOVED***

// Percentiles is a no-op.
func (NilHistogram) Percentiles(ps []float64) []float64 ***REMOVED***
	return make([]float64, len(ps))
***REMOVED***

// Sample is a no-op.
func (NilHistogram) Sample() Sample ***REMOVED*** return NilSample***REMOVED******REMOVED*** ***REMOVED***

// Snapshot is a no-op.
func (NilHistogram) Snapshot() Histogram ***REMOVED*** return NilHistogram***REMOVED******REMOVED*** ***REMOVED***

// StdDev is a no-op.
func (NilHistogram) StdDev() float64 ***REMOVED*** return 0.0 ***REMOVED***

// Sum is a no-op.
func (NilHistogram) Sum() int64 ***REMOVED*** return 0 ***REMOVED***

// Update is a no-op.
func (NilHistogram) Update(v int64) ***REMOVED******REMOVED***

// Variance is a no-op.
func (NilHistogram) Variance() float64 ***REMOVED*** return 0.0 ***REMOVED***

// StandardHistogram is the standard implementation of a Histogram and uses a
// Sample to bound its memory use.
type StandardHistogram struct ***REMOVED***
	sample Sample
***REMOVED***

// Clear clears the histogram and its sample.
func (h *StandardHistogram) Clear() ***REMOVED*** h.sample.Clear() ***REMOVED***

// Count returns the number of samples recorded since the histogram was last
// cleared.
func (h *StandardHistogram) Count() int64 ***REMOVED*** return h.sample.Count() ***REMOVED***

// Max returns the maximum value in the sample.
func (h *StandardHistogram) Max() int64 ***REMOVED*** return h.sample.Max() ***REMOVED***

// Mean returns the mean of the values in the sample.
func (h *StandardHistogram) Mean() float64 ***REMOVED*** return h.sample.Mean() ***REMOVED***

// Min returns the minimum value in the sample.
func (h *StandardHistogram) Min() int64 ***REMOVED*** return h.sample.Min() ***REMOVED***

// Percentile returns an arbitrary percentile of the values in the sample.
func (h *StandardHistogram) Percentile(p float64) float64 ***REMOVED***
	return h.sample.Percentile(p)
***REMOVED***

// Percentiles returns a slice of arbitrary percentiles of the values in the
// sample.
func (h *StandardHistogram) Percentiles(ps []float64) []float64 ***REMOVED***
	return h.sample.Percentiles(ps)
***REMOVED***

// Sample returns the Sample underlying the histogram.
func (h *StandardHistogram) Sample() Sample ***REMOVED*** return h.sample ***REMOVED***

// Snapshot returns a read-only copy of the histogram.
func (h *StandardHistogram) Snapshot() Histogram ***REMOVED***
	return &HistogramSnapshot***REMOVED***sample: h.sample.Snapshot().(*SampleSnapshot)***REMOVED***
***REMOVED***

// StdDev returns the standard deviation of the values in the sample.
func (h *StandardHistogram) StdDev() float64 ***REMOVED*** return h.sample.StdDev() ***REMOVED***

// Sum returns the sum in the sample.
func (h *StandardHistogram) Sum() int64 ***REMOVED*** return h.sample.Sum() ***REMOVED***

// Update samples a new value.
func (h *StandardHistogram) Update(v int64) ***REMOVED*** h.sample.Update(v) ***REMOVED***

// Variance returns the variance of the values in the sample.
func (h *StandardHistogram) Variance() float64 ***REMOVED*** return h.sample.Variance() ***REMOVED***
