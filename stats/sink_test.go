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

	"github.com/stretchr/testify/assert"
)

func TestTrendSink(t *testing.T) ***REMOVED***
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
			assert.Equal(t, 0.0, sink.P(0.0))
			assert.Equal(t, 55.0, sink.P(0.5))
			assert.Equal(t, 95.49999999999999, sink.P(0.95))
			assert.Equal(t, 99.1, sink.P(0.99))
			assert.Equal(t, 100.0, sink.P(1.0))
		***REMOVED***)
	***REMOVED***)
	t.Run("format", func(t *testing.T) ***REMOVED***
		sink := TrendSink***REMOVED******REMOVED***
		for _, s := range unsortedSamples10 ***REMOVED***
			sink.Add(Sample***REMOVED***Metric: &Metric***REMOVED******REMOVED***, Value: s***REMOVED***)
		***REMOVED***
		assert.Equal(t, map[string]float64***REMOVED***
			"min":   0.0,
			"max":   100.0,
			"avg":   54.0,
			"med":   55.0,
			"p(90)": 91.0,
			"p(95)": 95.49999999999999,
		***REMOVED***, sink.Format(0))
	***REMOVED***)
***REMOVED***

func TestDummySinkAddPanics(t *testing.T) ***REMOVED***
	assert.Panics(t, func() ***REMOVED***
		DummySink***REMOVED******REMOVED***.Add(Sample***REMOVED******REMOVED***)
	***REMOVED***)
***REMOVED***

func TestDummySinkFormatReturnsItself(t *testing.T) ***REMOVED***
	assert.Equal(t, map[string]float64***REMOVED***"a": 1***REMOVED***, DummySink***REMOVED***"a": 1***REMOVED***.Format(0))
***REMOVED***
