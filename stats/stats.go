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
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
var ErrInvalidMetricType = errors.New("invalid metric type")

// The serialized value type is invalid.
var ErrInvalidValueType = errors.New("invalid value type")

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

// SampleTags is an immutable string[string] map for tags. Once a tag
// set is created, direct modification is prohibited. It has
// copy-on-write semantics and uses pointers for faster comparison
// between maps, since the same tag set is often used for multiple samples.
// All methods should not panic, even if they are called on a nil pointer.
type SampleTags struct ***REMOVED***
	tags map[string]string
	json []byte
***REMOVED***

// Get returns an empty string and false if the the requested key is not
// present or its value and true if it is.
func (st *SampleTags) Get(key string) (string, bool) ***REMOVED***
	if st == nil ***REMOVED***
		return "", false
	***REMOVED***
	val, ok := st.tags[key]
	return val, ok
***REMOVED***

// IsEmpty checks for a nil pointer or zero tags.
// It's necessary because of this envconfig issue: https://github.com/kelseyhightower/envconfig/issues/113
func (st *SampleTags) IsEmpty() bool ***REMOVED***
	return st == nil || len(st.tags) == 0
***REMOVED***

// IsEqual tries to compare two tag sets with maximum efficiency.
func (st *SampleTags) IsEqual(other *SampleTags) bool ***REMOVED***
	if st == other ***REMOVED***
		return true
	***REMOVED***
	if st == nil || other == nil || len(st.tags) != len(other.tags) ***REMOVED***
		return false
	***REMOVED***
	for k, v := range st.tags ***REMOVED***
		if otherv, ok := other.tags[k]; !ok || v != otherv ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func (st *SampleTags) Contains(other *SampleTags) bool ***REMOVED***
	if st == other || other == nil ***REMOVED***
		return true
	***REMOVED***
	if st == nil || len(st.tags) < len(other.tags) ***REMOVED***
		return false
	***REMOVED***

	for k, v := range other.tags ***REMOVED***
		if st.tags[k] != v ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	return true
***REMOVED***

// MarshalJSON serializes SampleTags to a JSON string and caches
// the result. It is not thread safe in the sense that the Go race
// detector will complain if it's used concurrently, but no data
// should be corrupted.
func (st *SampleTags) MarshalJSON() ([]byte, error) ***REMOVED***
	if st.IsEmpty() ***REMOVED***
		return []byte("null"), nil
	***REMOVED***
	if st.json != nil ***REMOVED***
		return st.json, nil
	***REMOVED***
	res, err := json.Marshal(st.tags)
	if err != nil ***REMOVED***
		return res, err
	***REMOVED***
	st.json = res
	return res, nil
***REMOVED***

// UnmarshalJSON deserializes SampleTags from a JSON string.
func (st *SampleTags) UnmarshalJSON(data []byte) error ***REMOVED***
	if st == nil ***REMOVED***
		*st = SampleTags***REMOVED******REMOVED***
	***REMOVED***
	return json.Unmarshal(data, &st.tags)
***REMOVED***

// CloneTags copies the underlying set of a sample tags and
// returns it. If the receiver is nil, it returns an empty non-nil map.
func (st *SampleTags) CloneTags() map[string]string ***REMOVED***
	if st == nil ***REMOVED***
		return map[string]string***REMOVED******REMOVED***
	***REMOVED***
	res := make(map[string]string, len(st.tags))
	for k, v := range st.tags ***REMOVED***
		res[k] = v
	***REMOVED***
	return res
***REMOVED***

// NewSampleTags *copies* the supplied tag set and returns a new SampleTags
// instance with the key-value pairs from it.
func NewSampleTags(data map[string]string) *SampleTags ***REMOVED***
	if len(data) == 0 ***REMOVED***
		return nil
	***REMOVED***

	tags := map[string]string***REMOVED******REMOVED***
	for k, v := range data ***REMOVED***
		tags[k] = v
	***REMOVED***
	return &SampleTags***REMOVED***tags: tags***REMOVED***
***REMOVED***

// IntoSampleTags "consumes" the passed map and creates a new SampleTags
// struct with the data. The map is set to nil as a hint that it shouldn't
// be changed after it has been transformed into an "immutable" tag set.
// Oh, how I miss Rust and move semantics... :)
func IntoSampleTags(data *map[string]string) *SampleTags ***REMOVED***
	if len(*data) == 0 ***REMOVED***
		return nil
	***REMOVED***

	res := SampleTags***REMOVED***tags: *data***REMOVED***
	*data = nil
	return &res
***REMOVED***

// A Sample is a single measurement.
type Sample struct ***REMOVED***
	Metric *Metric
	Time   time.Time
	Tags   *SampleTags
	Value  float64
***REMOVED***

// SampleContainer is a simple abstraction that allows sample
// producers to attach extra information to samples they return
type SampleContainer interface ***REMOVED***
	GetSamples() []Sample
***REMOVED***

// Samples is just the simplest SampleContainer implementation
// that will be used when there's no need for extra information
type Samples []Sample

// GetSamples just implements the SampleContainer interface
func (s Samples) GetSamples() []Sample ***REMOVED***
	return s
***REMOVED***

// ConnectedSampleContainer is an extension of the SampleContainer
// interface that should be implemented when emitted samples
// are connected and share the same time and tags.
type ConnectedSampleContainer interface ***REMOVED***
	SampleContainer
	GetTags() *SampleTags
	GetTime() time.Time
***REMOVED***

// ConnectedSamples is the simplest ConnectedSampleContainer
// implementation that will be used when there's no need for
// extra information
type ConnectedSamples struct ***REMOVED***
	Samples []Sample
	Tags    *SampleTags
	Time    time.Time
***REMOVED***

// GetSamples implements the SampleContainer and ConnectedSampleContainer
// interfaces and returns the stored slice with samples.
func (cs ConnectedSamples) GetSamples() []Sample ***REMOVED***
	return cs.Samples
***REMOVED***

// GetTags implements ConnectedSampleContainer interface and returns stored tags.
func (cs ConnectedSamples) GetTags() *SampleTags ***REMOVED***
	return cs.Tags
***REMOVED***

// GetTime implements ConnectedSampleContainer interface and returns stored time.
func (cs ConnectedSamples) GetTime() time.Time ***REMOVED***
	return cs.Time
***REMOVED***

// GetSamples implement the ConnectedSampleContainer interface
// for a single Sample, since it's obviously connected with itself :)
func (s Sample) GetSamples() []Sample ***REMOVED***
	return []Sample***REMOVED***s***REMOVED***
***REMOVED***

// GetTags implements ConnectedSampleContainer interface
// and returns the sample's tags.
func (s Sample) GetTags() *SampleTags ***REMOVED***
	return s.Tags
***REMOVED***

// GetTime just implements ConnectedSampleContainer interface
// and returns the sample's time.
func (s Sample) GetTime() time.Time ***REMOVED***
	return s.Time
***REMOVED***

// Ensure that interfaces are implemented correctly
var _ SampleContainer = Sample***REMOVED******REMOVED***
var _ SampleContainer = Samples***REMOVED******REMOVED***

var _ ConnectedSampleContainer = Sample***REMOVED******REMOVED***
var _ ConnectedSampleContainer = ConnectedSamples***REMOVED******REMOVED***

// GetBufferedSamples will read all present (i.e. buffered or currently being pushed)
// values in the input channel and return them as a slice.
func GetBufferedSamples(input <-chan SampleContainer) (result []SampleContainer) ***REMOVED***
	for ***REMOVED***
		select ***REMOVED***
		case val, ok := <-input:
			if !ok ***REMOVED***
				return
			***REMOVED***
			result = append(result, val)
		default:
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// PushIfNotDone first checks if the supplied context is done and doesn't push
// the sample container if it is.
func PushIfNotDone(ctx context.Context, output chan<- SampleContainer, sample SampleContainer) bool ***REMOVED***
	if ctx.Err() != nil ***REMOVED***
		return false
	***REMOVED***
	output <- sample
	return true
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

var unitMap = map[string][]interface***REMOVED******REMOVED******REMOVED***
	"s":  ***REMOVED***"s", time.Second***REMOVED***,
	"ms": ***REMOVED***"ms", time.Millisecond***REMOVED***,
	"us": ***REMOVED***"µs", time.Microsecond***REMOVED***,
***REMOVED***

func (m *Metric) HumanizeValue(v float64, timeUnit string) string ***REMOVED***
	switch m.Type ***REMOVED***
	case Rate:
		// Truncate instead of round when decreasing precision to 2 decimal places
		return strconv.FormatFloat(float64(int(v*100*100))/100, 'f', 2, 64) + "%"
	default:
		switch m.Contains ***REMOVED***
		case Time:
			d := ToD(v)

			if v, ok := unitMap[timeUnit]; ok ***REMOVED***
				return fmt.Sprintf("%.2f%s", float64(d.Nanoseconds())/float64(v[1].(time.Duration)), v[0])
			***REMOVED***

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
	Name   string      `json:"name"`
	Parent string      `json:"parent"`
	Suffix string      `json:"suffix"`
	Tags   *SampleTags `json:"tags"`
	Metric *Metric     `json:"-"`
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
	return parts[0], &Submetric***REMOVED***Name: name, Parent: parts[0], Suffix: parts[1], Tags: IntoSampleTags(&tags)***REMOVED***
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
