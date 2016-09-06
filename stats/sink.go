package stats

import (
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

func (t *TrendSink) Format() map[string]float64 ***REMOVED***
	if t.jumbled ***REMOVED***
		sort.Float64s(t.Values)
		t.jumbled = false

		t.med = t.Values[t.count/2]
		if (t.count & 0x01) == 0 ***REMOVED***
			t.med = (t.med + t.Values[(t.count/2)-1]) / 2
		***REMOVED***
	***REMOVED***

	return map[string]float64***REMOVED***"min": t.min, "max": t.max, "avg": t.avg, "med": t.med***REMOVED***
***REMOVED***
