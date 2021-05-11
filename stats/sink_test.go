/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
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

package stats

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCounterSink(t *testing.T) ***REMOVED***
	samples10 := []float64***REMOVED***1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 100.0***REMOVED***
	now := time.Now()

	t.Run("add", func(t *testing.T) ***REMOVED***
		t.Run("one value", func(t *testing.T) ***REMOVED***
			sink := CounterSink***REMOVED******REMOVED***
			sink.Add(Sample***REMOVED***Metric: &Metric***REMOVED******REMOVED***, Value: 1.0, Time: now***REMOVED***)
			assert.Equal(t, 1.0, sink.Value)
			assert.Equal(t, now, sink.First)
		***REMOVED***)
		t.Run("values", func(t *testing.T) ***REMOVED***
			sink := CounterSink***REMOVED******REMOVED***
			for _, s := range samples10 ***REMOVED***
				sink.Add(Sample***REMOVED***Metric: &Metric***REMOVED******REMOVED***, Value: s, Time: now***REMOVED***)
			***REMOVED***
			assert.Equal(t, 145.0, sink.Value)
			assert.Equal(t, now, sink.First)
		***REMOVED***)
	***REMOVED***)
	t.Run("calc", func(t *testing.T) ***REMOVED***
		sink := CounterSink***REMOVED******REMOVED***
		sink.Calc()
		assert.Equal(t, 0.0, sink.Value)
		assert.Equal(t, time.Time***REMOVED******REMOVED***, sink.First)
	***REMOVED***)
	t.Run("format", func(t *testing.T) ***REMOVED***
		sink := CounterSink***REMOVED******REMOVED***
		for _, s := range samples10 ***REMOVED***
			sink.Add(Sample***REMOVED***Metric: &Metric***REMOVED******REMOVED***, Value: s, Time: now***REMOVED***)
		***REMOVED***
		assert.Equal(t, map[string]float64***REMOVED***"count": 145, "rate": 145.0***REMOVED***, sink.Format(1*time.Second))
	***REMOVED***)
***REMOVED***

func TestGaugeSink(t *testing.T) ***REMOVED***
	samples6 := []float64***REMOVED***1.0, 2.0, 3.0, 4.0, 10.0, 5.0***REMOVED***

	t.Run("add", func(t *testing.T) ***REMOVED***
		t.Run("one value", func(t *testing.T) ***REMOVED***
			sink := GaugeSink***REMOVED******REMOVED***
			sink.Add(Sample***REMOVED***Metric: &Metric***REMOVED******REMOVED***, Value: 1.0***REMOVED***)
			assert.Equal(t, 1.0, sink.Value)
			assert.Equal(t, 1.0, sink.Min)
			assert.Equal(t, true, sink.minSet)
			assert.Equal(t, 1.0, sink.Max)
		***REMOVED***)
		t.Run("values", func(t *testing.T) ***REMOVED***
			sink := GaugeSink***REMOVED******REMOVED***
			for _, s := range samples6 ***REMOVED***
				sink.Add(Sample***REMOVED***Metric: &Metric***REMOVED******REMOVED***, Value: s***REMOVED***)
			***REMOVED***
			assert.Equal(t, 5.0, sink.Value)
			assert.Equal(t, 1.0, sink.Min)
			assert.Equal(t, true, sink.minSet)
			assert.Equal(t, 10.0, sink.Max)
		***REMOVED***)
	***REMOVED***)
	t.Run("calc", func(t *testing.T) ***REMOVED***
		sink := GaugeSink***REMOVED******REMOVED***
		sink.Calc()
		assert.Equal(t, 0.0, sink.Value)
		assert.Equal(t, 0.0, sink.Min)
		assert.Equal(t, false, sink.minSet)
		assert.Equal(t, 0.0, sink.Max)
	***REMOVED***)
	t.Run("format", func(t *testing.T) ***REMOVED***
		sink := GaugeSink***REMOVED******REMOVED***
		for _, s := range samples6 ***REMOVED***
			sink.Add(Sample***REMOVED***Metric: &Metric***REMOVED******REMOVED***, Value: s***REMOVED***)
		***REMOVED***
		assert.Equal(t, map[string]float64***REMOVED***"value": 5.0***REMOVED***, sink.Format(0))
	***REMOVED***)
***REMOVED***

func TestTrendSink(t *testing.T) ***REMOVED***
	unsortedSamples5 := []float64***REMOVED***0.0, 5.0, 10.0, 3.0, 1.0***REMOVED***
	unsortedSamples10 := []float64***REMOVED***0.0, 100.0, 30.0, 80.0, 70.0, 60.0, 50.0, 40.0, 90.0, 20.0***REMOVED***

	t.Run("add", func(t *testing.T) ***REMOVED***
		t.Run("one value", func(t *testing.T) ***REMOVED***
			sink := TrendSink***REMOVED******REMOVED***
			sink.Add(Sample***REMOVED***Metric: &Metric***REMOVED******REMOVED***, Value: 7.0***REMOVED***)
			assert.Equal(t, uint64(1), sink.Count)
			assert.Equal(t, true, sink.jumbled)
			assert.Equal(t, 7.0, sink.Min)
			assert.Equal(t, 7.0, sink.Max)
			assert.Equal(t, 7.0, sink.Avg)
			assert.Equal(t, 0.0, sink.Med) // calculated in Calc()
		***REMOVED***)
		t.Run("values", func(t *testing.T) ***REMOVED***
			sink := TrendSink***REMOVED******REMOVED***
			for _, s := range unsortedSamples10 ***REMOVED***
				sink.Add(Sample***REMOVED***Metric: &Metric***REMOVED******REMOVED***, Value: s***REMOVED***)
			***REMOVED***
			assert.Equal(t, uint64(len(unsortedSamples10)), sink.Count)
			assert.Equal(t, true, sink.jumbled)
			assert.Equal(t, 0.0, sink.Min)
			assert.Equal(t, 100.0, sink.Max)
			assert.Equal(t, 54.0, sink.Avg)
			assert.Equal(t, 0.0, sink.Med) // calculated in Calc()
		***REMOVED***)
	***REMOVED***)
	t.Run("calc", func(t *testing.T) ***REMOVED***
		t.Run("no values", func(t *testing.T) ***REMOVED***
			sink := TrendSink***REMOVED******REMOVED***
			sink.Calc()
			assert.Equal(t, uint64(0), sink.Count)
			assert.Equal(t, false, sink.jumbled)
			assert.Equal(t, 0.0, sink.Med)
		***REMOVED***)
		t.Run("odd number of samples median", func(t *testing.T) ***REMOVED***
			sink := TrendSink***REMOVED******REMOVED***
			for _, s := range unsortedSamples5 ***REMOVED***
				sink.Add(Sample***REMOVED***Metric: &Metric***REMOVED******REMOVED***, Value: s***REMOVED***)
			***REMOVED***
			sink.Calc()
			assert.Equal(t, uint64(len(unsortedSamples5)), sink.Count)
			assert.Equal(t, false, sink.jumbled)
			assert.Equal(t, 3.0, sink.Med)
		***REMOVED***)
		t.Run("sorted", func(t *testing.T) ***REMOVED***
			sink := TrendSink***REMOVED******REMOVED***
			for _, s := range unsortedSamples10 ***REMOVED***
				sink.Add(Sample***REMOVED***Metric: &Metric***REMOVED******REMOVED***, Value: s***REMOVED***)
			***REMOVED***
			sink.Calc()
			assert.Equal(t, uint64(len(unsortedSamples10)), sink.Count)
			assert.Equal(t, false, sink.jumbled)
			assert.Equal(t, 55.0, sink.Med)
			assert.Equal(t, 0.0, sink.Min)
			assert.Equal(t, 100.0, sink.Max)
			assert.Equal(t, 54.0, sink.Avg)
		***REMOVED***)
	***REMOVED***)

	tolerance := 0.000001
	t.Run("percentile", func(t *testing.T) ***REMOVED***
		t.Run("no values", func(t *testing.T) ***REMOVED***
			sink := TrendSink***REMOVED******REMOVED***
			for i := 1; i <= 100; i++ ***REMOVED***
				assert.Equal(t, 0.0, sink.P(float64(i)/100.0))
			***REMOVED***
		***REMOVED***)
		t.Run("one value", func(t *testing.T) ***REMOVED***
			sink := TrendSink***REMOVED******REMOVED***
			sink.Add(Sample***REMOVED***Metric: &Metric***REMOVED******REMOVED***, Value: 10.0***REMOVED***)
			for i := 1; i <= 100; i++ ***REMOVED***
				assert.Equal(t, 10.0, sink.P(float64(i)/100.0))
			***REMOVED***
		***REMOVED***)
		t.Run("two values", func(t *testing.T) ***REMOVED***
			sink := TrendSink***REMOVED******REMOVED***
			sink.Add(Sample***REMOVED***Metric: &Metric***REMOVED******REMOVED***, Value: 5.0***REMOVED***)
			sink.Add(Sample***REMOVED***Metric: &Metric***REMOVED******REMOVED***, Value: 10.0***REMOVED***)
			assert.Equal(t, 5.0, sink.P(0.0))
			assert.Equal(t, 7.5, sink.P(0.5))
			assert.Equal(t, 5+(10-5)*0.95, sink.P(0.95))
			assert.Equal(t, 5+(10-5)*0.99, sink.P(0.99))
			assert.Equal(t, 10.0, sink.P(1.0))
		***REMOVED***)
		t.Run("more than 2", func(t *testing.T) ***REMOVED***
			sink := TrendSink***REMOVED******REMOVED***
			for _, s := range unsortedSamples10 ***REMOVED***
				sink.Add(Sample***REMOVED***Metric: &Metric***REMOVED******REMOVED***, Value: s***REMOVED***)
			***REMOVED***
			assert.InDelta(t, 0.0, sink.P(0.0), tolerance)
			assert.InDelta(t, 55.0, sink.P(0.5), tolerance)
			assert.InDelta(t, 95.5, sink.P(0.95), tolerance)
			assert.InDelta(t, 99.1, sink.P(0.99), tolerance)
			assert.InDelta(t, 100.0, sink.P(1.0), tolerance)
		***REMOVED***)
	***REMOVED***)
	t.Run("format", func(t *testing.T) ***REMOVED***
		sink := TrendSink***REMOVED******REMOVED***
		for _, s := range unsortedSamples10 ***REMOVED***
			sink.Add(Sample***REMOVED***Metric: &Metric***REMOVED******REMOVED***, Value: s***REMOVED***)
		***REMOVED***
		expected := map[string]float64***REMOVED***
			"min":   0.0,
			"max":   100.0,
			"avg":   54.0,
			"med":   55.0,
			"p(90)": 91.0,
			"p(95)": 95.5,
		***REMOVED***
		result := sink.Format(0)
		require.Equal(t, len(expected), len(result))
		for k, expV := range expected ***REMOVED***
			assert.Contains(t, result, k)
			assert.InDelta(t, expV, result[k], tolerance)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestRateSink(t *testing.T) ***REMOVED***
	samples6 := []float64***REMOVED***1.0, 0.0, 1.0, 0.0, 0.0, 1.0***REMOVED***

	t.Run("add", func(t *testing.T) ***REMOVED***
		t.Run("one true", func(t *testing.T) ***REMOVED***
			sink := RateSink***REMOVED******REMOVED***
			sink.Add(Sample***REMOVED***Metric: &Metric***REMOVED******REMOVED***, Value: 1.0***REMOVED***)
			assert.Equal(t, int64(1), sink.Total)
			assert.Equal(t, int64(1), sink.Trues)
		***REMOVED***)
		t.Run("one false", func(t *testing.T) ***REMOVED***
			sink := RateSink***REMOVED******REMOVED***
			sink.Add(Sample***REMOVED***Metric: &Metric***REMOVED******REMOVED***, Value: 0.0***REMOVED***)
			assert.Equal(t, int64(1), sink.Total)
			assert.Equal(t, int64(0), sink.Trues)
		***REMOVED***)
		t.Run("values", func(t *testing.T) ***REMOVED***
			sink := RateSink***REMOVED******REMOVED***
			for _, s := range samples6 ***REMOVED***
				sink.Add(Sample***REMOVED***Metric: &Metric***REMOVED******REMOVED***, Value: s***REMOVED***)
			***REMOVED***
			assert.Equal(t, int64(6), sink.Total)
			assert.Equal(t, int64(3), sink.Trues)
		***REMOVED***)
	***REMOVED***)
	t.Run("calc", func(t *testing.T) ***REMOVED***
		sink := RateSink***REMOVED******REMOVED***
		sink.Calc()
		assert.Equal(t, int64(0), sink.Total)
		assert.Equal(t, int64(0), sink.Trues)
	***REMOVED***)
	t.Run("format", func(t *testing.T) ***REMOVED***
		sink := RateSink***REMOVED******REMOVED***
		for _, s := range samples6 ***REMOVED***
			sink.Add(Sample***REMOVED***Metric: &Metric***REMOVED******REMOVED***, Value: s***REMOVED***)
		***REMOVED***
		assert.Equal(t, map[string]float64***REMOVED***"rate": 0.5***REMOVED***, sink.Format(0))
	***REMOVED***)
***REMOVED***

func TestDummySinkAddPanics(t *testing.T) ***REMOVED***
	assert.Panics(t, func() ***REMOVED***
		DummySink***REMOVED******REMOVED***.Add(Sample***REMOVED******REMOVED***)
	***REMOVED***)
***REMOVED***

func TestDummySinkCalcDoesNothing(t *testing.T) ***REMOVED***
	sink := DummySink***REMOVED***"a": 1***REMOVED***
	sink.Calc()
	assert.Equal(t, 1.0, sink["a"])
***REMOVED***

func TestDummySinkFormatReturnsItself(t *testing.T) ***REMOVED***
	assert.Equal(t, map[string]float64***REMOVED***"a": 1***REMOVED***, DummySink***REMOVED***"a": 1***REMOVED***.Format(0))
***REMOVED***
