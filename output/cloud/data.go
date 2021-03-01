/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2018 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package cloud

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/lib/netext/httpext"
	"github.com/loadimpact/k6/stats"
)

// DataType constants
const (
	DataTypeSingle             = "Point"
	DataTypeMap                = "Points"
	DataTypeAggregatedHTTPReqs = "AggregatedPoints"
)

//go:generate easyjson -pkg -no_std_marshalers -gen_build_flags -mod=mod .

func toMicroSecond(t time.Time) int64 ***REMOVED***
	return t.UnixNano() / 1000
***REMOVED***

// Sample is the generic struct that contains all types of data that we send to the cloud.
//easyjson:json
type Sample struct ***REMOVED***
	Type   string      `json:"type"`
	Metric string      `json:"metric"`
	Data   interface***REMOVED******REMOVED*** `json:"data"`
***REMOVED***

// UnmarshalJSON decodes the Data into the corresponding struct
func (ct *Sample) UnmarshalJSON(p []byte) error ***REMOVED***
	var tmpSample struct ***REMOVED***
		Type   string          `json:"type"`
		Metric string          `json:"metric"`
		Data   json.RawMessage `json:"data"`
	***REMOVED***
	if err := json.Unmarshal(p, &tmpSample); err != nil ***REMOVED***
		return err
	***REMOVED***
	s := Sample***REMOVED***
		Type:   tmpSample.Type,
		Metric: tmpSample.Metric,
	***REMOVED***

	switch tmpSample.Type ***REMOVED***
	case DataTypeSingle:
		s.Data = new(SampleDataSingle)
	case DataTypeMap:
		s.Data = new(SampleDataMap)
	case DataTypeAggregatedHTTPReqs:
		s.Data = new(SampleDataAggregatedHTTPReqs)
	default:
		return fmt.Errorf("unknown sample type '%s'", tmpSample.Type)
	***REMOVED***

	if err := json.Unmarshal(tmpSample.Data, &s.Data); err != nil ***REMOVED***
		return err
	***REMOVED***

	*ct = s
	return nil
***REMOVED***

// SampleDataSingle is used in all simple un-aggregated single-value samples.
//easyjson:json
type SampleDataSingle struct ***REMOVED***
	Time  int64             `json:"time,string"`
	Type  stats.MetricType  `json:"type"`
	Tags  *stats.SampleTags `json:"tags,omitempty"`
	Value float64           `json:"value"`
***REMOVED***

// SampleDataMap is used by samples that contain multiple values, currently
// that's only iteration metrics (`iter_li_all`) and unaggregated HTTP
// requests (`http_req_li_all`).
//easyjson:json
type SampleDataMap struct ***REMOVED***
	Time   int64              `json:"time,string"`
	Type   stats.MetricType   `json:"type"`
	Tags   *stats.SampleTags  `json:"tags,omitempty"`
	Values map[string]float64 `json:"values,omitempty"`
***REMOVED***

// NewSampleFromTrail just creates a ready-to-send Sample instance
// directly from a httpext.Trail.
func NewSampleFromTrail(trail *httpext.Trail) *Sample ***REMOVED***
	length := 8
	if trail.Failed.Valid ***REMOVED***
		length++
	***REMOVED***

	values := make(map[string]float64, length)
	values[metrics.HTTPReqs.Name] = 1
	values[metrics.HTTPReqDuration.Name] = stats.D(trail.Duration)
	values[metrics.HTTPReqBlocked.Name] = stats.D(trail.Blocked)
	values[metrics.HTTPReqConnecting.Name] = stats.D(trail.Connecting)
	values[metrics.HTTPReqTLSHandshaking.Name] = stats.D(trail.TLSHandshaking)
	values[metrics.HTTPReqSending.Name] = stats.D(trail.Sending)
	values[metrics.HTTPReqWaiting.Name] = stats.D(trail.Waiting)
	values[metrics.HTTPReqReceiving.Name] = stats.D(trail.Receiving)
	if trail.Failed.Valid ***REMOVED*** // this is done so the adding of 1 map element doesn't reexpand the map as this is a hotpath
		values[metrics.HTTPReqFailed.Name] = stats.B(trail.Failed.Bool)
	***REMOVED***
	return &Sample***REMOVED***
		Type:   DataTypeMap,
		Metric: "http_req_li_all",
		Data: &SampleDataMap***REMOVED***
			Time:   toMicroSecond(trail.GetTime()),
			Tags:   trail.GetTags(),
			Values: values,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// SampleDataAggregatedHTTPReqs is used in aggregated samples for HTTP requests.
//easyjson:json
type SampleDataAggregatedHTTPReqs struct ***REMOVED***
	Time   int64             `json:"time,string"`
	Type   string            `json:"type"`
	Count  uint64            `json:"count"`
	Tags   *stats.SampleTags `json:"tags,omitempty"`
	Values struct ***REMOVED***
		Duration       AggregatedMetric `json:"http_req_duration"`
		Blocked        AggregatedMetric `json:"http_req_blocked"`
		Connecting     AggregatedMetric `json:"http_req_connecting"`
		TLSHandshaking AggregatedMetric `json:"http_req_tls_handshaking"`
		Sending        AggregatedMetric `json:"http_req_sending"`
		Waiting        AggregatedMetric `json:"http_req_waiting"`
		Receiving      AggregatedMetric `json:"http_req_receiving"`
		Failed         AggregatedRate   `json:"http_req_failed,omitempty"`
	***REMOVED*** `json:"values"`
***REMOVED***

// Add updates all agregated values with the supplied trail data
func (sdagg *SampleDataAggregatedHTTPReqs) Add(trail *httpext.Trail) ***REMOVED***
	sdagg.Count++
	sdagg.Values.Duration.Add(trail.Duration)
	sdagg.Values.Blocked.Add(trail.Blocked)
	sdagg.Values.Connecting.Add(trail.Connecting)
	sdagg.Values.TLSHandshaking.Add(trail.TLSHandshaking)
	sdagg.Values.Sending.Add(trail.Sending)
	sdagg.Values.Waiting.Add(trail.Waiting)
	sdagg.Values.Receiving.Add(trail.Receiving)
	if trail.Failed.Valid ***REMOVED***
		sdagg.Values.Failed.Add(trail.Failed.Bool)
	***REMOVED***
***REMOVED***

// CalcAverages calculates and sets all `Avg` properties in the `Values` struct
func (sdagg *SampleDataAggregatedHTTPReqs) CalcAverages() ***REMOVED***
	count := float64(sdagg.Count)
	sdagg.Values.Duration.Calc(count)
	sdagg.Values.Blocked.Calc(count)
	sdagg.Values.Connecting.Calc(count)
	sdagg.Values.TLSHandshaking.Calc(count)
	sdagg.Values.Sending.Calc(count)
	sdagg.Values.Waiting.Calc(count)
	sdagg.Values.Receiving.Calc(count)
***REMOVED***

// AggregatedRate is an aggregation of a Rate metric
type AggregatedRate struct ***REMOVED***
	Count   float64 `json:"count"`
	NzCount float64 `json:"nz_count"`
***REMOVED***

// Add a boolean to the aggregated rate
func (ar *AggregatedRate) Add(b bool) ***REMOVED***
	ar.Count++
	if b ***REMOVED***
		ar.NzCount++
	***REMOVED***
***REMOVED***

// IsDefined implements easyjson.Optional
func (ar AggregatedRate) IsDefined() bool ***REMOVED***
	return ar.Count != 0
***REMOVED***

// AggregatedMetric is used to store aggregated information for a
// particular metric in an SampleDataAggregatedMap.
type AggregatedMetric struct ***REMOVED***
	// Used by Add() to keep working state
	minD time.Duration
	maxD time.Duration
	sumD time.Duration
	// Updated by Calc() and used in the JSON output
	Min float64 `json:"min"`
	Max float64 `json:"max"`
	Avg float64 `json:"avg"`
***REMOVED***

// Add the new duration to the internal sum and update Min and Max if necessary
func (am *AggregatedMetric) Add(t time.Duration) ***REMOVED***
	if am.sumD == 0 || am.minD > t ***REMOVED***
		am.minD = t
	***REMOVED***
	if am.maxD < t ***REMOVED***
		am.maxD = t
	***REMOVED***
	am.sumD += t
***REMOVED***

// Calc populates the float fields for min and max and calculates the average value
func (am *AggregatedMetric) Calc(count float64) ***REMOVED***
	am.Min = stats.D(am.minD)
	am.Max = stats.D(am.maxD)
	am.Avg = stats.D(am.sumD) / count
***REMOVED***

type aggregationBucket map[*stats.SampleTags][]*httpext.Trail

type durations []time.Duration

func (d durations) Len() int           ***REMOVED*** return len(d) ***REMOVED***
func (d durations) Swap(i, j int)      ***REMOVED*** d[i], d[j] = d[j], d[i] ***REMOVED***
func (d durations) Less(i, j int) bool ***REMOVED*** return d[i] < d[j] ***REMOVED***

// Used when there are fewer samples in the bucket (so we can interpolate)
// and for benchmark comparisons and verification of the quickselect
// algorithm (it should return exactly the same results if interpolation isn't used).
func (d durations) SortGetNormalBounds(radius, iqrLowerCoef, iqrUpperCoef float64, interpolate bool) (min, max time.Duration) ***REMOVED***
	if len(d) == 0 ***REMOVED***
		return
	***REMOVED***
	sort.Sort(d)
	last := float64(len(d) - 1)

	getValue := func(percentile float64) time.Duration ***REMOVED***
		i := percentile * last
		// If interpolation is not enabled, we round the resulting slice position
		// and return the value there.
		if !interpolate ***REMOVED***
			return d[int(math.Round(i))]
		***REMOVED***

		// Otherwise, we calculate (with linear interpolation) the value that
		// should fall at that percentile, given the values above and below it.
		floor := d[int(math.Floor(i))]
		ceil := d[int(math.Ceil(i))]
		posDiff := i - math.Floor(i)
		return floor + time.Duration(float64(ceil-floor)*posDiff)
	***REMOVED***

	// See https://en.wikipedia.org/wiki/Quartile#Outliers for details
	radius = math.Min(0.5, radius) // guard against a radius greater than 50%, see AggregationOutlierIqrRadius
	q1 := getValue(0.5 - radius)   // get Q1, the (interpolated) value at a `radius` distance before the median
	q3 := getValue(0.5 + radius)   // get Q3, the (interpolated) value at a `radius` distance after the median
	iqr := float64(q3 - q1)        // calculate the interquartile range (IQR)

	min = q1 - time.Duration(iqrLowerCoef*iqr) // lower fence, anything below this is an outlier
	max = q3 + time.Duration(iqrUpperCoef*iqr) // upper fence, anything above this is an outlier
	return
***REMOVED***

// Reworked and translated in Go from:
// https://github.com/haifengl/smile/blob/master/math/src/main/java/smile/sort/QuickSelect.java
// Originally Copyright (c) 2010 Haifeng Li
// Licensed under the Apache License, Version 2.0
//
// This could potentially be implemented as a standalone function
// that only depends on the sort.Interface methods, but that would
// probably introduce some performance overhead because of the
// dynamic dispatch.
func (d durations) quickSelect(k int) time.Duration ***REMOVED***
	n := len(d)
	l := 0
	ir := n - 1

	var i, j, mid int
	for ***REMOVED***
		if ir <= l+1 ***REMOVED***
			if ir == l+1 && d[ir] < d[l] ***REMOVED***
				d.Swap(l, ir)
			***REMOVED***
			return d[k]
		***REMOVED***
		mid = (l + ir) >> 1
		d.Swap(mid, l+1)
		if d[l] > d[ir] ***REMOVED***
			d.Swap(l, ir)
		***REMOVED***
		if d[l+1] > d[ir] ***REMOVED***
			d.Swap(l+1, ir)
		***REMOVED***
		if d[l] > d[l+1] ***REMOVED***
			d.Swap(l, l+1)
		***REMOVED***
		i = l + 1
		j = ir
		for ***REMOVED***
			for i++; d[i] < d[l+1]; i++ ***REMOVED***
			***REMOVED***
			for j--; d[j] > d[l+1]; j-- ***REMOVED***
			***REMOVED***
			if j < i ***REMOVED***
				break
			***REMOVED***
			d.Swap(i, j)
		***REMOVED***
		d.Swap(l+1, j)
		if j >= k ***REMOVED***
			ir = j - 1
		***REMOVED***
		if j <= k ***REMOVED***
			l = i
		***REMOVED***
	***REMOVED***
***REMOVED***

// Uses Quickselect to avoid sorting the whole slice
// https://en.wikipedia.org/wiki/Quickselect
func (d durations) SelectGetNormalBounds(radius, iqrLowerCoef, iqrUpperCoef float64) (min, max time.Duration) ***REMOVED***
	if len(d) == 0 ***REMOVED***
		return
	***REMOVED***
	radius = math.Min(0.5, radius)
	last := float64(len(d) - 1)

	q1 := d.quickSelect(int(math.Round(last * (0.5 - radius))))
	q3 := d.quickSelect(int(math.Round(last * (0.5 + radius))))

	iqr := float64(q3 - q1)
	min = q1 - time.Duration(iqrLowerCoef*iqr)
	max = q3 + time.Duration(iqrUpperCoef*iqr)
	return
***REMOVED***

//easyjson:json
type samples []*Sample
