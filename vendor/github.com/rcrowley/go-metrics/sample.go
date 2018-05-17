package metrics

import (
	"math"
	"math/rand"
	"sort"
	"sync"
	"time"
)

const rescaleThreshold = time.Hour

// Samples maintain a statistically-significant selection of values from
// a stream.
type Sample interface ***REMOVED***
	Clear()
	Count() int64
	Max() int64
	Mean() float64
	Min() int64
	Percentile(float64) float64
	Percentiles([]float64) []float64
	Size() int
	Snapshot() Sample
	StdDev() float64
	Sum() int64
	Update(int64)
	Values() []int64
	Variance() float64
***REMOVED***

// ExpDecaySample is an exponentially-decaying sample using a forward-decaying
// priority reservoir.  See Cormode et al's "Forward Decay: A Practical Time
// Decay Model for Streaming Systems".
//
// <http://dimacs.rutgers.edu/~graham/pubs/papers/fwddecay.pdf>
type ExpDecaySample struct ***REMOVED***
	alpha         float64
	count         int64
	mutex         sync.Mutex
	reservoirSize int
	t0, t1        time.Time
	values        *expDecaySampleHeap
***REMOVED***

// NewExpDecaySample constructs a new exponentially-decaying sample with the
// given reservoir size and alpha.
func NewExpDecaySample(reservoirSize int, alpha float64) Sample ***REMOVED***
	if UseNilMetrics ***REMOVED***
		return NilSample***REMOVED******REMOVED***
	***REMOVED***
	s := &ExpDecaySample***REMOVED***
		alpha:         alpha,
		reservoirSize: reservoirSize,
		t0:            time.Now(),
		values:        newExpDecaySampleHeap(reservoirSize),
	***REMOVED***
	s.t1 = s.t0.Add(rescaleThreshold)
	return s
***REMOVED***

// Clear clears all samples.
func (s *ExpDecaySample) Clear() ***REMOVED***
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.count = 0
	s.t0 = time.Now()
	s.t1 = s.t0.Add(rescaleThreshold)
	s.values.Clear()
***REMOVED***

// Count returns the number of samples recorded, which may exceed the
// reservoir size.
func (s *ExpDecaySample) Count() int64 ***REMOVED***
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.count
***REMOVED***

// Max returns the maximum value in the sample, which may not be the maximum
// value ever to be part of the sample.
func (s *ExpDecaySample) Max() int64 ***REMOVED***
	return SampleMax(s.Values())
***REMOVED***

// Mean returns the mean of the values in the sample.
func (s *ExpDecaySample) Mean() float64 ***REMOVED***
	return SampleMean(s.Values())
***REMOVED***

// Min returns the minimum value in the sample, which may not be the minimum
// value ever to be part of the sample.
func (s *ExpDecaySample) Min() int64 ***REMOVED***
	return SampleMin(s.Values())
***REMOVED***

// Percentile returns an arbitrary percentile of values in the sample.
func (s *ExpDecaySample) Percentile(p float64) float64 ***REMOVED***
	return SamplePercentile(s.Values(), p)
***REMOVED***

// Percentiles returns a slice of arbitrary percentiles of values in the
// sample.
func (s *ExpDecaySample) Percentiles(ps []float64) []float64 ***REMOVED***
	return SamplePercentiles(s.Values(), ps)
***REMOVED***

// Size returns the size of the sample, which is at most the reservoir size.
func (s *ExpDecaySample) Size() int ***REMOVED***
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.values.Size()
***REMOVED***

// Snapshot returns a read-only copy of the sample.
func (s *ExpDecaySample) Snapshot() Sample ***REMOVED***
	s.mutex.Lock()
	defer s.mutex.Unlock()
	vals := s.values.Values()
	values := make([]int64, len(vals))
	for i, v := range vals ***REMOVED***
		values[i] = v.v
	***REMOVED***
	return &SampleSnapshot***REMOVED***
		count:  s.count,
		values: values,
	***REMOVED***
***REMOVED***

// StdDev returns the standard deviation of the values in the sample.
func (s *ExpDecaySample) StdDev() float64 ***REMOVED***
	return SampleStdDev(s.Values())
***REMOVED***

// Sum returns the sum of the values in the sample.
func (s *ExpDecaySample) Sum() int64 ***REMOVED***
	return SampleSum(s.Values())
***REMOVED***

// Update samples a new value.
func (s *ExpDecaySample) Update(v int64) ***REMOVED***
	s.update(time.Now(), v)
***REMOVED***

// Values returns a copy of the values in the sample.
func (s *ExpDecaySample) Values() []int64 ***REMOVED***
	s.mutex.Lock()
	defer s.mutex.Unlock()
	vals := s.values.Values()
	values := make([]int64, len(vals))
	for i, v := range vals ***REMOVED***
		values[i] = v.v
	***REMOVED***
	return values
***REMOVED***

// Variance returns the variance of the values in the sample.
func (s *ExpDecaySample) Variance() float64 ***REMOVED***
	return SampleVariance(s.Values())
***REMOVED***

// update samples a new value at a particular timestamp.  This is a method all
// its own to facilitate testing.
func (s *ExpDecaySample) update(t time.Time, v int64) ***REMOVED***
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.count++
	if s.values.Size() == s.reservoirSize ***REMOVED***
		s.values.Pop()
	***REMOVED***
	s.values.Push(expDecaySample***REMOVED***
		k: math.Exp(t.Sub(s.t0).Seconds()*s.alpha) / rand.Float64(),
		v: v,
	***REMOVED***)
	if t.After(s.t1) ***REMOVED***
		values := s.values.Values()
		t0 := s.t0
		s.values.Clear()
		s.t0 = t
		s.t1 = s.t0.Add(rescaleThreshold)
		for _, v := range values ***REMOVED***
			v.k = v.k * math.Exp(-s.alpha*s.t0.Sub(t0).Seconds())
			s.values.Push(v)
		***REMOVED***
	***REMOVED***
***REMOVED***

// NilSample is a no-op Sample.
type NilSample struct***REMOVED******REMOVED***

// Clear is a no-op.
func (NilSample) Clear() ***REMOVED******REMOVED***

// Count is a no-op.
func (NilSample) Count() int64 ***REMOVED*** return 0 ***REMOVED***

// Max is a no-op.
func (NilSample) Max() int64 ***REMOVED*** return 0 ***REMOVED***

// Mean is a no-op.
func (NilSample) Mean() float64 ***REMOVED*** return 0.0 ***REMOVED***

// Min is a no-op.
func (NilSample) Min() int64 ***REMOVED*** return 0 ***REMOVED***

// Percentile is a no-op.
func (NilSample) Percentile(p float64) float64 ***REMOVED*** return 0.0 ***REMOVED***

// Percentiles is a no-op.
func (NilSample) Percentiles(ps []float64) []float64 ***REMOVED***
	return make([]float64, len(ps))
***REMOVED***

// Size is a no-op.
func (NilSample) Size() int ***REMOVED*** return 0 ***REMOVED***

// Sample is a no-op.
func (NilSample) Snapshot() Sample ***REMOVED*** return NilSample***REMOVED******REMOVED*** ***REMOVED***

// StdDev is a no-op.
func (NilSample) StdDev() float64 ***REMOVED*** return 0.0 ***REMOVED***

// Sum is a no-op.
func (NilSample) Sum() int64 ***REMOVED*** return 0 ***REMOVED***

// Update is a no-op.
func (NilSample) Update(v int64) ***REMOVED******REMOVED***

// Values is a no-op.
func (NilSample) Values() []int64 ***REMOVED*** return []int64***REMOVED******REMOVED*** ***REMOVED***

// Variance is a no-op.
func (NilSample) Variance() float64 ***REMOVED*** return 0.0 ***REMOVED***

// SampleMax returns the maximum value of the slice of int64.
func SampleMax(values []int64) int64 ***REMOVED***
	if 0 == len(values) ***REMOVED***
		return 0
	***REMOVED***
	var max int64 = math.MinInt64
	for _, v := range values ***REMOVED***
		if max < v ***REMOVED***
			max = v
		***REMOVED***
	***REMOVED***
	return max
***REMOVED***

// SampleMean returns the mean value of the slice of int64.
func SampleMean(values []int64) float64 ***REMOVED***
	if 0 == len(values) ***REMOVED***
		return 0.0
	***REMOVED***
	return float64(SampleSum(values)) / float64(len(values))
***REMOVED***

// SampleMin returns the minimum value of the slice of int64.
func SampleMin(values []int64) int64 ***REMOVED***
	if 0 == len(values) ***REMOVED***
		return 0
	***REMOVED***
	var min int64 = math.MaxInt64
	for _, v := range values ***REMOVED***
		if min > v ***REMOVED***
			min = v
		***REMOVED***
	***REMOVED***
	return min
***REMOVED***

// SamplePercentiles returns an arbitrary percentile of the slice of int64.
func SamplePercentile(values int64Slice, p float64) float64 ***REMOVED***
	return SamplePercentiles(values, []float64***REMOVED***p***REMOVED***)[0]
***REMOVED***

// SamplePercentiles returns a slice of arbitrary percentiles of the slice of
// int64.
func SamplePercentiles(values int64Slice, ps []float64) []float64 ***REMOVED***
	scores := make([]float64, len(ps))
	size := len(values)
	if size > 0 ***REMOVED***
		sort.Sort(values)
		for i, p := range ps ***REMOVED***
			pos := p * float64(size+1)
			if pos < 1.0 ***REMOVED***
				scores[i] = float64(values[0])
			***REMOVED*** else if pos >= float64(size) ***REMOVED***
				scores[i] = float64(values[size-1])
			***REMOVED*** else ***REMOVED***
				lower := float64(values[int(pos)-1])
				upper := float64(values[int(pos)])
				scores[i] = lower + (pos-math.Floor(pos))*(upper-lower)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return scores
***REMOVED***

// SampleSnapshot is a read-only copy of another Sample.
type SampleSnapshot struct ***REMOVED***
	count  int64
	values []int64
***REMOVED***

func NewSampleSnapshot(count int64, values []int64) *SampleSnapshot ***REMOVED***
	return &SampleSnapshot***REMOVED***
		count:  count,
		values: values,
	***REMOVED***
***REMOVED***

// Clear panics.
func (*SampleSnapshot) Clear() ***REMOVED***
	panic("Clear called on a SampleSnapshot")
***REMOVED***

// Count returns the count of inputs at the time the snapshot was taken.
func (s *SampleSnapshot) Count() int64 ***REMOVED*** return s.count ***REMOVED***

// Max returns the maximal value at the time the snapshot was taken.
func (s *SampleSnapshot) Max() int64 ***REMOVED*** return SampleMax(s.values) ***REMOVED***

// Mean returns the mean value at the time the snapshot was taken.
func (s *SampleSnapshot) Mean() float64 ***REMOVED*** return SampleMean(s.values) ***REMOVED***

// Min returns the minimal value at the time the snapshot was taken.
func (s *SampleSnapshot) Min() int64 ***REMOVED*** return SampleMin(s.values) ***REMOVED***

// Percentile returns an arbitrary percentile of values at the time the
// snapshot was taken.
func (s *SampleSnapshot) Percentile(p float64) float64 ***REMOVED***
	return SamplePercentile(s.values, p)
***REMOVED***

// Percentiles returns a slice of arbitrary percentiles of values at the time
// the snapshot was taken.
func (s *SampleSnapshot) Percentiles(ps []float64) []float64 ***REMOVED***
	return SamplePercentiles(s.values, ps)
***REMOVED***

// Size returns the size of the sample at the time the snapshot was taken.
func (s *SampleSnapshot) Size() int ***REMOVED*** return len(s.values) ***REMOVED***

// Snapshot returns the snapshot.
func (s *SampleSnapshot) Snapshot() Sample ***REMOVED*** return s ***REMOVED***

// StdDev returns the standard deviation of values at the time the snapshot was
// taken.
func (s *SampleSnapshot) StdDev() float64 ***REMOVED*** return SampleStdDev(s.values) ***REMOVED***

// Sum returns the sum of values at the time the snapshot was taken.
func (s *SampleSnapshot) Sum() int64 ***REMOVED*** return SampleSum(s.values) ***REMOVED***

// Update panics.
func (*SampleSnapshot) Update(int64) ***REMOVED***
	panic("Update called on a SampleSnapshot")
***REMOVED***

// Values returns a copy of the values in the sample.
func (s *SampleSnapshot) Values() []int64 ***REMOVED***
	values := make([]int64, len(s.values))
	copy(values, s.values)
	return values
***REMOVED***

// Variance returns the variance of values at the time the snapshot was taken.
func (s *SampleSnapshot) Variance() float64 ***REMOVED*** return SampleVariance(s.values) ***REMOVED***

// SampleStdDev returns the standard deviation of the slice of int64.
func SampleStdDev(values []int64) float64 ***REMOVED***
	return math.Sqrt(SampleVariance(values))
***REMOVED***

// SampleSum returns the sum of the slice of int64.
func SampleSum(values []int64) int64 ***REMOVED***
	var sum int64
	for _, v := range values ***REMOVED***
		sum += v
	***REMOVED***
	return sum
***REMOVED***

// SampleVariance returns the variance of the slice of int64.
func SampleVariance(values []int64) float64 ***REMOVED***
	if 0 == len(values) ***REMOVED***
		return 0.0
	***REMOVED***
	m := SampleMean(values)
	var sum float64
	for _, v := range values ***REMOVED***
		d := float64(v) - m
		sum += d * d
	***REMOVED***
	return sum / float64(len(values))
***REMOVED***

// A uniform sample using Vitter's Algorithm R.
//
// <http://www.cs.umd.edu/~samir/498/vitter.pdf>
type UniformSample struct ***REMOVED***
	count         int64
	mutex         sync.Mutex
	reservoirSize int
	values        []int64
***REMOVED***

// NewUniformSample constructs a new uniform sample with the given reservoir
// size.
func NewUniformSample(reservoirSize int) Sample ***REMOVED***
	if UseNilMetrics ***REMOVED***
		return NilSample***REMOVED******REMOVED***
	***REMOVED***
	return &UniformSample***REMOVED***
		reservoirSize: reservoirSize,
		values:        make([]int64, 0, reservoirSize),
	***REMOVED***
***REMOVED***

// Clear clears all samples.
func (s *UniformSample) Clear() ***REMOVED***
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.count = 0
	s.values = make([]int64, 0, s.reservoirSize)
***REMOVED***

// Count returns the number of samples recorded, which may exceed the
// reservoir size.
func (s *UniformSample) Count() int64 ***REMOVED***
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.count
***REMOVED***

// Max returns the maximum value in the sample, which may not be the maximum
// value ever to be part of the sample.
func (s *UniformSample) Max() int64 ***REMOVED***
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return SampleMax(s.values)
***REMOVED***

// Mean returns the mean of the values in the sample.
func (s *UniformSample) Mean() float64 ***REMOVED***
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return SampleMean(s.values)
***REMOVED***

// Min returns the minimum value in the sample, which may not be the minimum
// value ever to be part of the sample.
func (s *UniformSample) Min() int64 ***REMOVED***
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return SampleMin(s.values)
***REMOVED***

// Percentile returns an arbitrary percentile of values in the sample.
func (s *UniformSample) Percentile(p float64) float64 ***REMOVED***
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return SamplePercentile(s.values, p)
***REMOVED***

// Percentiles returns a slice of arbitrary percentiles of values in the
// sample.
func (s *UniformSample) Percentiles(ps []float64) []float64 ***REMOVED***
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return SamplePercentiles(s.values, ps)
***REMOVED***

// Size returns the size of the sample, which is at most the reservoir size.
func (s *UniformSample) Size() int ***REMOVED***
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return len(s.values)
***REMOVED***

// Snapshot returns a read-only copy of the sample.
func (s *UniformSample) Snapshot() Sample ***REMOVED***
	s.mutex.Lock()
	defer s.mutex.Unlock()
	values := make([]int64, len(s.values))
	copy(values, s.values)
	return &SampleSnapshot***REMOVED***
		count:  s.count,
		values: values,
	***REMOVED***
***REMOVED***

// StdDev returns the standard deviation of the values in the sample.
func (s *UniformSample) StdDev() float64 ***REMOVED***
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return SampleStdDev(s.values)
***REMOVED***

// Sum returns the sum of the values in the sample.
func (s *UniformSample) Sum() int64 ***REMOVED***
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return SampleSum(s.values)
***REMOVED***

// Update samples a new value.
func (s *UniformSample) Update(v int64) ***REMOVED***
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.count++
	if len(s.values) < s.reservoirSize ***REMOVED***
		s.values = append(s.values, v)
	***REMOVED*** else ***REMOVED***
		r := rand.Int63n(s.count)
		if r < int64(len(s.values)) ***REMOVED***
			s.values[int(r)] = v
		***REMOVED***
	***REMOVED***
***REMOVED***

// Values returns a copy of the values in the sample.
func (s *UniformSample) Values() []int64 ***REMOVED***
	s.mutex.Lock()
	defer s.mutex.Unlock()
	values := make([]int64, len(s.values))
	copy(values, s.values)
	return values
***REMOVED***

// Variance returns the variance of the values in the sample.
func (s *UniformSample) Variance() float64 ***REMOVED***
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return SampleVariance(s.values)
***REMOVED***

// expDecaySample represents an individual sample in a heap.
type expDecaySample struct ***REMOVED***
	k float64
	v int64
***REMOVED***

func newExpDecaySampleHeap(reservoirSize int) *expDecaySampleHeap ***REMOVED***
	return &expDecaySampleHeap***REMOVED***make([]expDecaySample, 0, reservoirSize)***REMOVED***
***REMOVED***

// expDecaySampleHeap is a min-heap of expDecaySamples.
// The internal implementation is copied from the standard library's container/heap
type expDecaySampleHeap struct ***REMOVED***
	s []expDecaySample
***REMOVED***

func (h *expDecaySampleHeap) Clear() ***REMOVED***
	h.s = h.s[:0]
***REMOVED***

func (h *expDecaySampleHeap) Push(s expDecaySample) ***REMOVED***
	n := len(h.s)
	h.s = h.s[0 : n+1]
	h.s[n] = s
	h.up(n)
***REMOVED***

func (h *expDecaySampleHeap) Pop() expDecaySample ***REMOVED***
	n := len(h.s) - 1
	h.s[0], h.s[n] = h.s[n], h.s[0]
	h.down(0, n)

	n = len(h.s)
	s := h.s[n-1]
	h.s = h.s[0 : n-1]
	return s
***REMOVED***

func (h *expDecaySampleHeap) Size() int ***REMOVED***
	return len(h.s)
***REMOVED***

func (h *expDecaySampleHeap) Values() []expDecaySample ***REMOVED***
	return h.s
***REMOVED***

func (h *expDecaySampleHeap) up(j int) ***REMOVED***
	for ***REMOVED***
		i := (j - 1) / 2 // parent
		if i == j || !(h.s[j].k < h.s[i].k) ***REMOVED***
			break
		***REMOVED***
		h.s[i], h.s[j] = h.s[j], h.s[i]
		j = i
	***REMOVED***
***REMOVED***

func (h *expDecaySampleHeap) down(i, n int) ***REMOVED***
	for ***REMOVED***
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 ***REMOVED*** // j1 < 0 after int overflow
			break
		***REMOVED***
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && !(h.s[j1].k < h.s[j2].k) ***REMOVED***
			j = j2 // = 2*i + 2  // right child
		***REMOVED***
		if !(h.s[j].k < h.s[i].k) ***REMOVED***
			break
		***REMOVED***
		h.s[i], h.s[j] = h.s[j], h.s[i]
		i = j
	***REMOVED***
***REMOVED***

type int64Slice []int64

func (p int64Slice) Len() int           ***REMOVED*** return len(p) ***REMOVED***
func (p int64Slice) Less(i, j int) bool ***REMOVED*** return p[i] < p[j] ***REMOVED***
func (p int64Slice) Swap(i, j int)      ***REMOVED*** p[i], p[j] = p[j], p[i] ***REMOVED***
