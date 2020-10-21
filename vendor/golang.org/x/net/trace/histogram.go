// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trace

// This file implements histogramming for RPC statistics collection.

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"math"
	"sync"

	"golang.org/x/net/internal/timeseries"
)

const (
	bucketCount = 38
)

// histogram keeps counts of values in buckets that are spaced
// out in powers of 2: 0-1, 2-3, 4-7...
// histogram implements timeseries.Observable
type histogram struct ***REMOVED***
	sum          int64   // running total of measurements
	sumOfSquares float64 // square of running total
	buckets      []int64 // bucketed values for histogram
	value        int     // holds a single value as an optimization
	valueCount   int64   // number of values recorded for single value
***REMOVED***

// AddMeasurement records a value measurement observation to the histogram.
func (h *histogram) addMeasurement(value int64) ***REMOVED***
	// TODO: assert invariant
	h.sum += value
	h.sumOfSquares += float64(value) * float64(value)

	bucketIndex := getBucket(value)

	if h.valueCount == 0 || (h.valueCount > 0 && h.value == bucketIndex) ***REMOVED***
		h.value = bucketIndex
		h.valueCount++
	***REMOVED*** else ***REMOVED***
		h.allocateBuckets()
		h.buckets[bucketIndex]++
	***REMOVED***
***REMOVED***

func (h *histogram) allocateBuckets() ***REMOVED***
	if h.buckets == nil ***REMOVED***
		h.buckets = make([]int64, bucketCount)
		h.buckets[h.value] = h.valueCount
		h.value = 0
		h.valueCount = -1
	***REMOVED***
***REMOVED***

func log2(i int64) int ***REMOVED***
	n := 0
	for ; i >= 0x100; i >>= 8 ***REMOVED***
		n += 8
	***REMOVED***
	for ; i > 0; i >>= 1 ***REMOVED***
		n += 1
	***REMOVED***
	return n
***REMOVED***

func getBucket(i int64) (index int) ***REMOVED***
	index = log2(i) - 1
	if index < 0 ***REMOVED***
		index = 0
	***REMOVED***
	if index >= bucketCount ***REMOVED***
		index = bucketCount - 1
	***REMOVED***
	return
***REMOVED***

// Total returns the number of recorded observations.
func (h *histogram) total() (total int64) ***REMOVED***
	if h.valueCount >= 0 ***REMOVED***
		total = h.valueCount
	***REMOVED***
	for _, val := range h.buckets ***REMOVED***
		total += int64(val)
	***REMOVED***
	return
***REMOVED***

// Average returns the average value of recorded observations.
func (h *histogram) average() float64 ***REMOVED***
	t := h.total()
	if t == 0 ***REMOVED***
		return 0
	***REMOVED***
	return float64(h.sum) / float64(t)
***REMOVED***

// Variance returns the variance of recorded observations.
func (h *histogram) variance() float64 ***REMOVED***
	t := float64(h.total())
	if t == 0 ***REMOVED***
		return 0
	***REMOVED***
	s := float64(h.sum) / t
	return h.sumOfSquares/t - s*s
***REMOVED***

// StandardDeviation returns the standard deviation of recorded observations.
func (h *histogram) standardDeviation() float64 ***REMOVED***
	return math.Sqrt(h.variance())
***REMOVED***

// PercentileBoundary estimates the value that the given fraction of recorded
// observations are less than.
func (h *histogram) percentileBoundary(percentile float64) int64 ***REMOVED***
	total := h.total()

	// Corner cases (make sure result is strictly less than Total())
	if total == 0 ***REMOVED***
		return 0
	***REMOVED*** else if total == 1 ***REMOVED***
		return int64(h.average())
	***REMOVED***

	percentOfTotal := round(float64(total) * percentile)
	var runningTotal int64

	for i := range h.buckets ***REMOVED***
		value := h.buckets[i]
		runningTotal += value
		if runningTotal == percentOfTotal ***REMOVED***
			// We hit an exact bucket boundary. If the next bucket has data, it is a
			// good estimate of the value. If the bucket is empty, we interpolate the
			// midpoint between the next bucket's boundary and the next non-zero
			// bucket. If the remaining buckets are all empty, then we use the
			// boundary for the next bucket as the estimate.
			j := uint8(i + 1)
			min := bucketBoundary(j)
			if runningTotal < total ***REMOVED***
				for h.buckets[j] == 0 ***REMOVED***
					j++
				***REMOVED***
			***REMOVED***
			max := bucketBoundary(j)
			return min + round(float64(max-min)/2)
		***REMOVED*** else if runningTotal > percentOfTotal ***REMOVED***
			// The value is in this bucket. Interpolate the value.
			delta := runningTotal - percentOfTotal
			percentBucket := float64(value-delta) / float64(value)
			bucketMin := bucketBoundary(uint8(i))
			nextBucketMin := bucketBoundary(uint8(i + 1))
			bucketSize := nextBucketMin - bucketMin
			return bucketMin + round(percentBucket*float64(bucketSize))
		***REMOVED***
	***REMOVED***
	return bucketBoundary(bucketCount - 1)
***REMOVED***

// Median returns the estimated median of the observed values.
func (h *histogram) median() int64 ***REMOVED***
	return h.percentileBoundary(0.5)
***REMOVED***

// Add adds other to h.
func (h *histogram) Add(other timeseries.Observable) ***REMOVED***
	o := other.(*histogram)
	if o.valueCount == 0 ***REMOVED***
		// Other histogram is empty
	***REMOVED*** else if h.valueCount >= 0 && o.valueCount > 0 && h.value == o.value ***REMOVED***
		// Both have a single bucketed value, aggregate them
		h.valueCount += o.valueCount
	***REMOVED*** else ***REMOVED***
		// Two different values necessitate buckets in this histogram
		h.allocateBuckets()
		if o.valueCount >= 0 ***REMOVED***
			h.buckets[o.value] += o.valueCount
		***REMOVED*** else ***REMOVED***
			for i := range h.buckets ***REMOVED***
				h.buckets[i] += o.buckets[i]
			***REMOVED***
		***REMOVED***
	***REMOVED***
	h.sumOfSquares += o.sumOfSquares
	h.sum += o.sum
***REMOVED***

// Clear resets the histogram to an empty state, removing all observed values.
func (h *histogram) Clear() ***REMOVED***
	h.buckets = nil
	h.value = 0
	h.valueCount = 0
	h.sum = 0
	h.sumOfSquares = 0
***REMOVED***

// CopyFrom copies from other, which must be a *histogram, into h.
func (h *histogram) CopyFrom(other timeseries.Observable) ***REMOVED***
	o := other.(*histogram)
	if o.valueCount == -1 ***REMOVED***
		h.allocateBuckets()
		copy(h.buckets, o.buckets)
	***REMOVED***
	h.sum = o.sum
	h.sumOfSquares = o.sumOfSquares
	h.value = o.value
	h.valueCount = o.valueCount
***REMOVED***

// Multiply scales the histogram by the specified ratio.
func (h *histogram) Multiply(ratio float64) ***REMOVED***
	if h.valueCount == -1 ***REMOVED***
		for i := range h.buckets ***REMOVED***
			h.buckets[i] = int64(float64(h.buckets[i]) * ratio)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		h.valueCount = int64(float64(h.valueCount) * ratio)
	***REMOVED***
	h.sum = int64(float64(h.sum) * ratio)
	h.sumOfSquares = h.sumOfSquares * ratio
***REMOVED***

// New creates a new histogram.
func (h *histogram) New() timeseries.Observable ***REMOVED***
	r := new(histogram)
	r.Clear()
	return r
***REMOVED***

func (h *histogram) String() string ***REMOVED***
	return fmt.Sprintf("%d, %f, %d, %d, %v",
		h.sum, h.sumOfSquares, h.value, h.valueCount, h.buckets)
***REMOVED***

// round returns the closest int64 to the argument
func round(in float64) int64 ***REMOVED***
	return int64(math.Floor(in + 0.5))
***REMOVED***

// bucketBoundary returns the first value in the bucket.
func bucketBoundary(bucket uint8) int64 ***REMOVED***
	if bucket == 0 ***REMOVED***
		return 0
	***REMOVED***
	return 1 << bucket
***REMOVED***

// bucketData holds data about a specific bucket for use in distTmpl.
type bucketData struct ***REMOVED***
	Lower, Upper       int64
	N                  int64
	Pct, CumulativePct float64
	GraphWidth         int
***REMOVED***

// data holds data about a Distribution for use in distTmpl.
type data struct ***REMOVED***
	Buckets                 []*bucketData
	Count, Median           int64
	Mean, StandardDeviation float64
***REMOVED***

// maxHTMLBarWidth is the maximum width of the HTML bar for visualizing buckets.
const maxHTMLBarWidth = 350.0

// newData returns data representing h for use in distTmpl.
func (h *histogram) newData() *data ***REMOVED***
	// Force the allocation of buckets to simplify the rendering implementation
	h.allocateBuckets()
	// We scale the bars on the right so that the largest bar is
	// maxHTMLBarWidth pixels in width.
	maxBucket := int64(0)
	for _, n := range h.buckets ***REMOVED***
		if n > maxBucket ***REMOVED***
			maxBucket = n
		***REMOVED***
	***REMOVED***
	total := h.total()
	barsizeMult := maxHTMLBarWidth / float64(maxBucket)
	var pctMult float64
	if total == 0 ***REMOVED***
		pctMult = 1.0
	***REMOVED*** else ***REMOVED***
		pctMult = 100.0 / float64(total)
	***REMOVED***

	buckets := make([]*bucketData, len(h.buckets))
	runningTotal := int64(0)
	for i, n := range h.buckets ***REMOVED***
		if n == 0 ***REMOVED***
			continue
		***REMOVED***
		runningTotal += n
		var upperBound int64
		if i < bucketCount-1 ***REMOVED***
			upperBound = bucketBoundary(uint8(i + 1))
		***REMOVED*** else ***REMOVED***
			upperBound = math.MaxInt64
		***REMOVED***
		buckets[i] = &bucketData***REMOVED***
			Lower:         bucketBoundary(uint8(i)),
			Upper:         upperBound,
			N:             n,
			Pct:           float64(n) * pctMult,
			CumulativePct: float64(runningTotal) * pctMult,
			GraphWidth:    int(float64(n) * barsizeMult),
		***REMOVED***
	***REMOVED***
	return &data***REMOVED***
		Buckets:           buckets,
		Count:             total,
		Median:            h.median(),
		Mean:              h.average(),
		StandardDeviation: h.standardDeviation(),
	***REMOVED***
***REMOVED***

func (h *histogram) html() template.HTML ***REMOVED***
	buf := new(bytes.Buffer)
	if err := distTmpl().Execute(buf, h.newData()); err != nil ***REMOVED***
		buf.Reset()
		log.Printf("net/trace: couldn't execute template: %v", err)
	***REMOVED***
	return template.HTML(buf.String())
***REMOVED***

var distTmplCache *template.Template
var distTmplOnce sync.Once

func distTmpl() *template.Template ***REMOVED***
	distTmplOnce.Do(func() ***REMOVED***
		// Input: data
		distTmplCache = template.Must(template.New("distTmpl").Parse(`
<table>
<tr>
    <td style="padding:0.25em">Count: ***REMOVED******REMOVED***.Count***REMOVED******REMOVED***</td>
    <td style="padding:0.25em">Mean: ***REMOVED******REMOVED***printf "%.0f" .Mean***REMOVED******REMOVED***</td>
    <td style="padding:0.25em">StdDev: ***REMOVED******REMOVED***printf "%.0f" .StandardDeviation***REMOVED******REMOVED***</td>
    <td style="padding:0.25em">Median: ***REMOVED******REMOVED***.Median***REMOVED******REMOVED***</td>
</tr>
</table>
<hr>
<table>
***REMOVED******REMOVED***range $b := .Buckets***REMOVED******REMOVED***
***REMOVED******REMOVED***if $b***REMOVED******REMOVED***
  <tr>
    <td style="padding:0 0 0 0.25em">[</td>
    <td style="text-align:right;padding:0 0.25em">***REMOVED******REMOVED***.Lower***REMOVED******REMOVED***,</td>
    <td style="text-align:right;padding:0 0.25em">***REMOVED******REMOVED***.Upper***REMOVED******REMOVED***)</td>
    <td style="text-align:right;padding:0 0.25em">***REMOVED******REMOVED***.N***REMOVED******REMOVED***</td>
    <td style="text-align:right;padding:0 0.25em">***REMOVED******REMOVED***printf "%#.3f" .Pct***REMOVED******REMOVED***%</td>
    <td style="text-align:right;padding:0 0.25em">***REMOVED******REMOVED***printf "%#.3f" .CumulativePct***REMOVED******REMOVED***%</td>
    <td><div style="background-color: blue; height: 1em; width: ***REMOVED******REMOVED***.GraphWidth***REMOVED******REMOVED***;"></div></td>
  </tr>
***REMOVED******REMOVED***end***REMOVED******REMOVED***
***REMOVED******REMOVED***end***REMOVED******REMOVED***
</table>
`))
	***REMOVED***)
	return distTmplCache
***REMOVED***
