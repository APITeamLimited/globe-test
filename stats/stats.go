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
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"gopkg.in/guregu/null.v3"
)

const (
	counterString = `"counter"`
	gaugeString   = `"gauge"`
	trendString   = `"trend"`
	rateString    = `"rate"`

	defaultString = `"default"`
	timeString    = `"time"`
	dataString    = `"data"`
)

// Possible values for MetricType.
const (
	Counter = MetricType(iota) // A counter that sums its data points
	Gauge                      // A gauge that displays the latest value
	Trend                      // A trend, min/max/avg/med are interesting
	Rate                       // A rate, displays % of values that aren't 0
)

// Possible values for ValueType.
const (
	Default = ValueType(iota) // Values are presented as-is
	Time                      // Values are timestamps (nanoseconds)
	Data                      // Values are data amounts (bytes)
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
	case Data:
		return []byte(dataString), nil
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
	case dataString:
		*t = Data
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
	case Data:
		return dataString
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
	Name       string       `json:"name"`
	Type       MetricType   `json:"type"`
	Contains   ValueType    `json:"contains"`
	Tainted    null.Bool    `json:"tainted"`
	Thresholds Thresholds   `json:"thresholds"`
	Submetrics []*Submetric `json:"submetrics"`
	Sub        Submetric    `json:"sub,omitempty"`
	Sink       Sink         `json:"-"`
***REMOVED***

func New(name string, typ MetricType, t ...ValueType) *Metric ***REMOVED***
	vt := Default
	if len(t) > 0 ***REMOVED***
		vt = t[0]
	***REMOVED***
	var sink Sink
	switch typ ***REMOVED***
	case Counter:
		sink = &CounterSink***REMOVED******REMOVED***
	case Gauge:
		sink = &GaugeSink***REMOVED******REMOVED***
	case Trend:
		sink = &TrendSink***REMOVED******REMOVED***
	case Rate:
		sink = &RateSink***REMOVED******REMOVED***
	default:
		return nil
	***REMOVED***
	return &Metric***REMOVED***Name: name, Type: typ, Contains: vt, Sink: sink***REMOVED***
***REMOVED***

func (m *Metric) HumanizeValue(v float64) string ***REMOVED***
	switch m.Type ***REMOVED***
	case Rate:
		// Truncate instead of round when decreasing precision to 2 decimal places
		return strconv.FormatFloat(float64(int(v*100*100))/100, 'f', 2, 64) + "%"
	default:
		switch m.Contains ***REMOVED***
		case Time:
			d := ToD(v)
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
		case Data:
			return humanize.Bytes(uint64(v))
		default:
			return humanize.Ftoa(v)
		***REMOVED***
	***REMOVED***
***REMOVED***

// A Submetric represents a filtered dataset based on a parent metric.
type Submetric struct ***REMOVED***
	Name   string            `json:"name"`
	Parent string            `json:"parent"`
	Suffix string            `json:"suffix"`
	Tags   map[string]string `json:"tags"`
	Metric *Metric           `json:"-"`
***REMOVED***

// Creates a submetric from a name.
func NewSubmetric(name string) (parentName string, sm *Submetric) ***REMOVED***
	parts := strings.SplitN(strings.TrimSuffix(name, "***REMOVED***"), "***REMOVED***", 2)
	if len(parts) == 1 ***REMOVED***
		return parts[0], &Submetric***REMOVED***Name: name***REMOVED***
	***REMOVED***

	kvs := strings.Split(parts[1], ",")
	tags := make(map[string]string, len(kvs))
	for _, kv := range kvs ***REMOVED***
		if kv == "" ***REMOVED***
			continue
		***REMOVED***
		parts := strings.SplitN(kv, ":", 2)

		key := strings.TrimSpace(strings.Trim(parts[0], `"'`))
		if len(parts) != 2 ***REMOVED***
			tags[key] = ""
			continue
		***REMOVED***

		value := strings.TrimSpace(strings.Trim(parts[1], `"'`))
		tags[key] = value
	***REMOVED***
	return parts[0], &Submetric***REMOVED***Name: name, Parent: parts[0], Suffix: parts[1], Tags: tags***REMOVED***
***REMOVED***

func (m *Metric) Summary(t time.Duration) *Summary ***REMOVED***
	return &Summary***REMOVED***
		Metric:  m,
		Summary: m.Sink.Format(t),
	***REMOVED***
***REMOVED***

type Summary struct ***REMOVED***
	Metric  *Metric            `json:"metric"`
	Summary map[string]float64 `json:"summary"`
***REMOVED***
