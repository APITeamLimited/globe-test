package stats

import (
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
