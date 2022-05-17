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

package metrics

import (
	"errors"
	"math"
	"sort"
	"time"
)

var (
	_ Sink = &CounterSink***REMOVED******REMOVED***
	_ Sink = &GaugeSink***REMOVED******REMOVED***
	_ Sink = &TrendSink***REMOVED******REMOVED***
	_ Sink = &RateSink***REMOVED******REMOVED***
	_ Sink = &DummySink***REMOVED******REMOVED***
)

type Sink interface ***REMOVED***
	Add(s Sample)                              // Add a sample to the sink.
	Calc()                                     // Make final calculations.
	Format(t time.Duration) map[string]float64 // Data for thresholds.
	IsEmpty() bool                             // Check if the Sink is empty.
***REMOVED***

type CounterSink struct ***REMOVED***
	Value float64
	First time.Time
***REMOVED***

func (c *CounterSink) Add(s Sample) ***REMOVED***
	c.Value += s.Value
	if c.First.IsZero() ***REMOVED***
		c.First = s.Time
	***REMOVED***
***REMOVED***

// IsEmpty indicates whether the CounterSink is empty.
func (c *CounterSink) IsEmpty() bool ***REMOVED*** return c.First.IsZero() ***REMOVED***

func (c *CounterSink) Calc() ***REMOVED******REMOVED***

func (c *CounterSink) Format(t time.Duration) map[string]float64 ***REMOVED***
	return map[string]float64***REMOVED***
		"count": c.Value,
		"rate":  c.Value / (float64(t) / float64(time.Second)),
	***REMOVED***
***REMOVED***

type GaugeSink struct ***REMOVED***
	Value    float64
	Max, Min float64
	minSet   bool
***REMOVED***

// IsEmpty indicates whether the GaugeSink is empty.
func (g *GaugeSink) IsEmpty() bool ***REMOVED*** return !g.minSet ***REMOVED***

func (g *GaugeSink) Add(s Sample) ***REMOVED***
	g.Value = s.Value
	if s.Value > g.Max ***REMOVED***
		g.Max = s.Value
	***REMOVED***
	if s.Value < g.Min || !g.minSet ***REMOVED***
		g.Min = s.Value
		g.minSet = true
	***REMOVED***
***REMOVED***

func (g *GaugeSink) Calc() ***REMOVED******REMOVED***

func (g *GaugeSink) Format(t time.Duration) map[string]float64 ***REMOVED***
	return map[string]float64***REMOVED***"value": g.Value***REMOVED***
***REMOVED***

type TrendSink struct ***REMOVED***
	Values  []float64
	jumbled bool

	Count    uint64
	Min, Max float64
	Sum, Avg float64
	Med      float64
***REMOVED***

// IsEmpty indicates whether the TrendSink is empty.
func (t *TrendSink) IsEmpty() bool ***REMOVED*** return t.Count == 0 ***REMOVED***

func (t *TrendSink) Add(s Sample) ***REMOVED***
	t.Values = append(t.Values, s.Value)
	t.jumbled = true
	t.Count += 1
	t.Sum += s.Value
	t.Avg = t.Sum / float64(t.Count)

	if s.Value > t.Max ***REMOVED***
		t.Max = s.Value
	***REMOVED***
	if s.Value < t.Min || t.Count == 1 ***REMOVED***
		t.Min = s.Value
	***REMOVED***
***REMOVED***

// P calculates the given percentile from sink values.
func (t *TrendSink) P(pct float64) float64 ***REMOVED***
	switch t.Count ***REMOVED***
	case 0:
		return 0
	case 1:
		return t.Values[0]
	default:
		// If percentile falls on a value in Values slice, we return that value.
		// If percentile does not fall on a value in Values slice, we calculate (linear interpolation)
		// the value that would fall at percentile, given the values above and below that percentile.
		t.Calc()
		i := pct * (float64(t.Count) - 1.0)
		j := t.Values[int(math.Floor(i))]
		k := t.Values[int(math.Ceil(i))]
		f := i - math.Floor(i)
		return j + (k-j)*f
	***REMOVED***
***REMOVED***

func (t *TrendSink) Calc() ***REMOVED***
	if !t.jumbled ***REMOVED***
		return
	***REMOVED***

	sort.Float64s(t.Values)
	t.jumbled = false

	// The median of an even number of values is the average of the middle two.
	if (t.Count & 0x01) == 0 ***REMOVED***
		t.Med = (t.Values[(t.Count/2)-1] + t.Values[(t.Count/2)]) / 2
	***REMOVED*** else ***REMOVED***
		t.Med = t.Values[t.Count/2]
	***REMOVED***
***REMOVED***

func (t *TrendSink) Format(tt time.Duration) map[string]float64 ***REMOVED***
	t.Calc()
	// TODO: respect the summaryTrendStats for REST API
	return map[string]float64***REMOVED***
		"min":   t.Min,
		"max":   t.Max,
		"avg":   t.Avg,
		"med":   t.Med,
		"p(90)": t.P(0.90),
		"p(95)": t.P(0.95),
	***REMOVED***
***REMOVED***

type RateSink struct ***REMOVED***
	Trues int64
	Total int64
***REMOVED***

// IsEmpty indicates whether the RateSink is empty.
func (r *RateSink) IsEmpty() bool ***REMOVED*** return r.Total == 0 ***REMOVED***

func (r *RateSink) Add(s Sample) ***REMOVED***
	r.Total += 1
	if s.Value != 0 ***REMOVED***
		r.Trues += 1
	***REMOVED***
***REMOVED***

func (r RateSink) Calc() ***REMOVED******REMOVED***

func (r RateSink) Format(t time.Duration) map[string]float64 ***REMOVED***
	var rate float64
	if r.Total > 0 ***REMOVED***
		rate = float64(r.Trues) / float64(r.Total)
	***REMOVED***

	return map[string]float64***REMOVED***"rate": rate***REMOVED***
***REMOVED***

type DummySink map[string]float64

// IsEmpty indicates whether the DummySink is empty.
func (d DummySink) IsEmpty() bool ***REMOVED*** return len(d) == 0 ***REMOVED***

func (d DummySink) Add(s Sample) ***REMOVED***
	panic(errors.New("you can't add samples to a dummy sink"))
***REMOVED***

func (d DummySink) Calc() ***REMOVED******REMOVED***

func (d DummySink) Format(t time.Duration) map[string]float64 ***REMOVED***
	return map[string]float64(d)
***REMOVED***
