package stats

import (
	"sort"
	"time"
)

// The type of a metric.
type MetricType int

// Possible values for MetricType.
const (
	NoType  = MetricType(iota) // No type; metrics like this are ignored
	Counter                    // A counter that sums its data points
	Gauge                      // A gauge that displays the latest value
	Trend                      // A trend, min/max/avg/med are interesting
)

// The type of values a metric contains.
type ValueType int

// Possible values for ValueType.
const (
	Default = ValueType(iota) // Values are presented as-is
	Time                      // Values are timestamps (nanoseconds)
)

// A Sample is a single measurement.
type Sample struct ***REMOVED***
	Time  time.Time
	Tags  map[string]string
	Value float64
***REMOVED***

// An MSample is a Sample tagged with a Metric, to make returning samples easier.
type FatSample struct ***REMOVED***
	Sample
	Metric *Metric
***REMOVED***

// A Metric defines the shape of a set of data.
type Metric struct ***REMOVED***
	Name     string
	Type     MetricType
	Contains ValueType
***REMOVED***

func New(name string, typ MetricType, t ...ValueType) *Metric ***REMOVED***
	vt := Default
	if len(t) > 0 ***REMOVED***
		vt = t[0]
	***REMOVED***
	return &Metric***REMOVED***Name: name, Type: typ, Contains: vt***REMOVED***
***REMOVED***

func (m *Metric) Format(samples []Sample) map[string]float64 ***REMOVED***
	switch m.Type ***REMOVED***
	case Counter:
		var total float64
		for _, s := range samples ***REMOVED***
			total += s.Value
		***REMOVED***
		return map[string]float64***REMOVED***"count": total***REMOVED***
	case Gauge:
		l := len(samples)
		if l == 0 ***REMOVED***
			return map[string]float64***REMOVED***"value": 0***REMOVED***
		***REMOVED***
		return map[string]float64***REMOVED***"value": samples[l-1].Value***REMOVED***
	case Trend:
		values := make([]float64, len(samples))
		for i, s := range samples ***REMOVED***
			values[i] = s.Value
		***REMOVED***
		sort.Float64s(values)

		var min, max, avg, med, sum float64
		for i, v := range values ***REMOVED***
			if v < min || i == 0 ***REMOVED***
				min = v
			***REMOVED***
			if v > max ***REMOVED***
				max = v
			***REMOVED***
			sum += v
		***REMOVED***

		l := len(values)
		switch l ***REMOVED***
		case 0:
		case 1:
			avg = values[0]
			med = values[0]
		default:
			avg = sum / float64(l)
			med = values[l/2]

			// Median for an even number of values is the average of the middle two
			if (l & 0x01) == 0 ***REMOVED***
				med = (med + values[(l/2)-1]) / 2
			***REMOVED***
		***REMOVED***

		return map[string]float64***REMOVED***
			"min": min,
			"max": max,
			"med": med,
			"avg": avg,
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***
