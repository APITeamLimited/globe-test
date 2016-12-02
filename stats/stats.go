package stats

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	counterString = `"counter"`
	gaugeString   = `"gauge"`
	trendString   = `"trend"`
	rateString    = `"rate"`

	defaultString = `"default"`
	timeString    = `"time"`
)

// Possible values for MetricType.
const (
	NoType  = MetricType(iota) // No type; metrics like this are ignored
	Counter                    // A counter that sums its data points
	Gauge                      // A gauge that displays the latest value
	Trend                      // A trend, min/max/avg/med are interesting
	Rate                       // A rate, displays % of values that aren't 0
)

// Possible values for ValueType.
const (
	Default = ValueType(iota) // Values are presented as-is
	Time                      // Values are timestamps (nanoseconds)
)

// The serialized metric type is invalid.
var ErrInvalidMetricType = errors.New("Invalid metric type")

// The serialized value type is invalid.
var ErrInvalidValueType = errors.New("Invalid value type")

// A MetricType specifies the type of a metric.
type MetricType int

// MarshalJSON serializes a MetricType as a human readable string.
func (t MetricType) MarshalJSON() ([]byte, error) ***REMOVED***
	switch t ***REMOVED***
	case Counter:
		return []byte(counterString), nil
	case Gauge:
		return []byte(gaugeString), nil
	case Trend:
		return []byte(trendString), nil
	case Rate:
		return []byte(rateString), nil
	default:
		return nil, ErrInvalidMetricType
	***REMOVED***
***REMOVED***

// UnmarshalJSON deserializes a MetricType from a string representation.
func (t *MetricType) UnmarshalJSON(data []byte) error ***REMOVED***
	switch string(data) ***REMOVED***
	case counterString:
		*t = Counter
	case gaugeString:
		*t = Gauge
	case trendString:
		*t = Trend
	case rateString:
		*t = Rate
	default:
		return ErrInvalidMetricType
	***REMOVED***

	return nil
***REMOVED***

func (t MetricType) String() string ***REMOVED***
	switch t ***REMOVED***
	case Counter:
		return counterString
	case Gauge:
		return gaugeString
	case Trend:
		return trendString
	case Rate:
		return rateString
	default:
		return "[INVALID]"
	***REMOVED***
***REMOVED***

// The type of values a metric contains.
type ValueType int

// MarshalJSON serializes a ValueType as a human readable string.
func (t ValueType) MarshalJSON() ([]byte, error) ***REMOVED***
	switch t ***REMOVED***
	case Default:
		return []byte(defaultString), nil
	case Time:
		return []byte(timeString), nil
	default:
		return nil, ErrInvalidValueType
	***REMOVED***
***REMOVED***

// UnmarshalJSON deserializes a ValueType from a string representation.
func (t *ValueType) UnmarshalJSON(data []byte) error ***REMOVED***
	switch string(data) ***REMOVED***
	case defaultString:
		*t = Default
	case timeString:
		*t = Time
	default:
		return ErrInvalidValueType
	***REMOVED***

	return nil
***REMOVED***

func (t ValueType) String() string ***REMOVED***
	switch t ***REMOVED***
	case Default:
		return defaultString
	case Time:
		return timeString
	default:
		return "[INVALID]"
	***REMOVED***
***REMOVED***

// A Sample is a single measurement.
type Sample struct ***REMOVED***
	Metric *Metric
	Time   time.Time
	Tags   map[string]string
	Value  float64
***REMOVED***

// A Metric defines the shape of a set of data.
type Metric struct ***REMOVED***
	Name     string     `json:"-"`
	Type     MetricType `json:"type"`
	Contains ValueType  `json:"contains"`

	// Filled in by the API when requested, the server side cannot count on its presence.
	Sample map[string]float64 `json:"sample"`

	// Set to true if the metric has failed a threshold.
	Tainted bool
***REMOVED***

func New(name string, typ MetricType, t ...ValueType) *Metric ***REMOVED***
	vt := Default
	if len(t) > 0 ***REMOVED***
		vt = t[0]
	***REMOVED***
	return &Metric***REMOVED***Name: name, Type: typ, Contains: vt***REMOVED***
***REMOVED***

func (m Metric) Humanize() string ***REMOVED***
	sample := m.Sample
	switch len(sample) ***REMOVED***
	case 0:
		return ""
	case 1:
		for _, v := range sample ***REMOVED***
			return m.HumanizeValue(v)
		***REMOVED***
		return ""
	default:
		parts := make([]string, 0, len(m.Sample))
		for key, val := range m.Sample ***REMOVED***
			parts = append(parts, fmt.Sprintf("%s=%s", key, m.HumanizeValue(val)))
		***REMOVED***
		sort.Strings(parts)
		return strings.Join(parts, ", ")
	***REMOVED***
***REMOVED***

func (m Metric) HumanizeValue(v float64) string ***REMOVED***
	switch m.Type ***REMOVED***
	case Rate:
		return strconv.FormatFloat(100*v, 'f', 2, 64) + "%"
	default:
		switch m.Contains ***REMOVED***
		case Time:
			d := time.Duration(v)
			switch ***REMOVED***
			case d > time.Minute:
				d -= d % (1 * time.Second)
			case d > time.Second:
				d -= d % (10 * time.Millisecond)
			case d > time.Millisecond:
				d -= d % (10 * time.Microsecond)
			case d > time.Microsecond:
				d -= d % (10 * time.Nanosecond)
			***REMOVED***
			return d.String()
		default:
			return strconv.FormatFloat(v, 'f', -1, 64)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (m Metric) GetID() string ***REMOVED***
	return m.Name
***REMOVED***

func (m *Metric) SetID(id string) error ***REMOVED***
	m.Name = id
	return nil
***REMOVED***
