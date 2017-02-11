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
	"errors"
	"sort"
)

type Sink interface ***REMOVED***
	Add(s Sample)
	Format() map[string]float64
***REMOVED***

type CounterSink struct ***REMOVED***
	Value float64
***REMOVED***

func (c *CounterSink) Add(s Sample) ***REMOVED***
	c.Value += s.Value
***REMOVED***

func (c *CounterSink) Format() map[string]float64 ***REMOVED***
	return map[string]float64***REMOVED***"count": c.Value***REMOVED***
***REMOVED***

type GaugeSink struct ***REMOVED***
	Value float64
***REMOVED***

func (g *GaugeSink) Add(s Sample) ***REMOVED***
	g.Value = s.Value
***REMOVED***

func (g *GaugeSink) Format() map[string]float64 ***REMOVED***
	return map[string]float64***REMOVED***"value": g.Value***REMOVED***
***REMOVED***

type TrendSink struct ***REMOVED***
	Values []float64

	jumbled  bool
	count    uint64
	min, max float64
	sum, avg float64
	med      float64
***REMOVED***

func (t *TrendSink) Add(s Sample) ***REMOVED***
	t.Values = append(t.Values, s.Value)
	t.jumbled = true
	t.count += 1
	t.sum += s.Value
	t.avg = t.sum / float64(t.count)

	if s.Value > t.max ***REMOVED***
		t.max = s.Value
	***REMOVED***
	if s.Value < t.min || t.min == 0 ***REMOVED***
		t.min = s.Value
	***REMOVED***
***REMOVED***

func (t *TrendSink) P(pct float64) float64 ***REMOVED***
	switch t.count ***REMOVED***
	case 0:
		return 0
	case 1:
		return t.Values[0]
	case 2:
		if pct < 0.5 ***REMOVED***
			return t.Values[0]
		***REMOVED*** else ***REMOVED***
			return t.Values[1]
		***REMOVED***
	default:
		return t.Values[int(float64(t.count)*pct)]
	***REMOVED***
***REMOVED***

func (t *TrendSink) Format() map[string]float64 ***REMOVED***
	if t.jumbled ***REMOVED***
		sort.Float64s(t.Values)
		t.jumbled = false

		t.med = t.Values[t.count/2]
		if (t.count & 0x01) == 0 ***REMOVED***
			t.med = (t.med + t.Values[(t.count/2)-1]) / 2
		***REMOVED***
	***REMOVED***

	return map[string]float64***REMOVED***
		"min": t.min,
		"max": t.max,
		"avg": t.avg,
		"med": t.med,
		"p90": t.P(0.90),
		"p95": t.P(0.95),
	***REMOVED***
***REMOVED***

type RateSink struct ***REMOVED***
	Trues int64
	Total int64
***REMOVED***

func (r *RateSink) Add(s Sample) ***REMOVED***
	r.Total += 1
	if s.Value != 0 ***REMOVED***
		r.Trues += 1
	***REMOVED***
***REMOVED***

func (r RateSink) Format() map[string]float64 ***REMOVED***
	return map[string]float64***REMOVED***"rate": float64(r.Trues) / float64(r.Total)***REMOVED***
***REMOVED***

type DummySink map[string]float64

func (d DummySink) Add(s Sample) ***REMOVED***
	panic(errors.New("you can't add samples to a dummy sink"))
***REMOVED***

func (d DummySink) Format() map[string]float64 ***REMOVED***
	return map[string]float64(d)
***REMOVED***
