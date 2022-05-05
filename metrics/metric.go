package metrics

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gopkg.in/guregu/null.v3"
)

// A Metric defines the shape of a set of data.
type Metric struct ***REMOVED***
	Name     string     `json:"name"`
	Type     MetricType `json:"type"`
	Contains ValueType  `json:"contains"`

	// TODO: decouple the metrics from the sinks and thresholds... have them
	// linked, but not in the same struct?
	Tainted    null.Bool    `json:"tainted"`
	Thresholds Thresholds   `json:"thresholds"`
	Submetrics []*Submetric `json:"submetrics"`
	Sub        *Submetric   `json:"-"`
	Sink       Sink         `json:"-"`
	Observed   bool         `json:"-"`
***REMOVED***

// Sample samples the metric at the given time, with the provided tags and value
func (m *Metric) Sample(t time.Time, tags *SampleTags, value float64) Sample ***REMOVED***
	return Sample***REMOVED***
		Time:   t,
		Tags:   tags,
		Value:  value,
		Metric: m,
	***REMOVED***
***REMOVED***

// newMetric instantiates a new Metric
func newMetric(name string, mt MetricType, vt ...ValueType) *Metric ***REMOVED***
	valueType := Default
	if len(vt) > 0 ***REMOVED***
		valueType = vt[0]
	***REMOVED***

	var sink Sink
	switch mt ***REMOVED***
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

	return &Metric***REMOVED***
		Name:     name,
		Type:     mt,
		Contains: valueType,
		Sink:     sink,
	***REMOVED***
***REMOVED***

// A Submetric represents a filtered dataset based on a parent metric.
type Submetric struct ***REMOVED***
	Name   string      `json:"name"`
	Suffix string      `json:"suffix"` // TODO: rename?
	Tags   *SampleTags `json:"tags"`

	Metric *Metric `json:"-"`
	Parent *Metric `json:"-"`
***REMOVED***

// AddSubmetric creates a new submetric from the key:value threshold definition
// and adds it to the metric's submetrics list.
func (m *Metric) AddSubmetric(keyValues string) (*Submetric, error) ***REMOVED***
	keyValues = strings.TrimSpace(keyValues)
	if len(keyValues) == 0 ***REMOVED***
		return nil, fmt.Errorf("submetric criteria for metric '%s' cannot be empty", m.Name)
	***REMOVED***
	kvs := strings.Split(keyValues, ",")
	rawTags := make(map[string]string, len(kvs))
	for _, kv := range kvs ***REMOVED***
		if kv == "" ***REMOVED***
			continue
		***REMOVED***
		parts := strings.SplitN(kv, ":", 2)

		key := strings.Trim(strings.TrimSpace(parts[0]), `"'`)
		if len(parts) != 2 ***REMOVED***
			rawTags[key] = ""
			continue
		***REMOVED***

		value := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
		rawTags[key] = value
	***REMOVED***

	tags := IntoSampleTags(&rawTags)

	for _, sm := range m.Submetrics ***REMOVED***
		if sm.Tags.IsEqual(tags) ***REMOVED***
			return nil, fmt.Errorf(
				"sub-metric with params '%s' already exists for metric %s: %s",
				keyValues, m.Name, sm.Name,
			)
		***REMOVED***
	***REMOVED***

	subMetric := &Submetric***REMOVED***
		Name:   m.Name + "***REMOVED***" + keyValues + "***REMOVED***",
		Suffix: keyValues,
		Tags:   tags,
		Parent: m,
	***REMOVED***
	subMetricMetric := newMetric(subMetric.Name, m.Type, m.Contains)
	subMetricMetric.Sub = subMetric // sigh
	subMetric.Metric = subMetricMetric

	m.Submetrics = append(m.Submetrics, subMetric)

	return subMetric, nil
***REMOVED***

// ErrMetricNameParsing indicates parsing a metric name failed
var ErrMetricNameParsing = errors.New("parsing metric name failed")

// ParseMetricName parses a metric name expression of the form metric_name***REMOVED***tag_key:tag_value,...***REMOVED***
// Its first return value is the parsed metric name, second are parsed tags as as slice
// of "key:value" strings. On failure, it returns an error containing the `ErrMetricNameParsing` in its chain.
func ParseMetricName(name string) (string, []string, error) ***REMOVED***
	openingTokenPos := strings.IndexByte(name, '***REMOVED***')
	closingTokenPos := strings.LastIndexByte(name, '***REMOVED***')
	containsOpeningToken := openingTokenPos != -1
	containsClosingToken := closingTokenPos != -1

	// Neither the opening '***REMOVED***' token nor the closing '***REMOVED***' token
	// are present, thus the metric name only consists of a literal.
	if !containsOpeningToken && !containsClosingToken ***REMOVED***
		return name, nil, nil
	***REMOVED***

	// If the name contains an opening or closing token, but not
	// its counterpart, the expression is malformed.
	if (containsOpeningToken && !containsClosingToken) ||
		(!containsOpeningToken && containsClosingToken) ***REMOVED***
		return "", nil, fmt.Errorf(
			"%w, metric %q has unmatched opening/close curly brace",
			ErrMetricNameParsing, name,
		)
	***REMOVED***

	// If the closing brace token appears before the opening one,
	// the expression is malformed
	if closingTokenPos < openingTokenPos ***REMOVED***
		return "", nil, fmt.Errorf("%w, metric %q closing curly brace appears before opening one", ErrMetricNameParsing, name)
	***REMOVED***

	// If the last character is not a closing brace token,
	// the expression is malformed.
	if closingTokenPos != (len(name) - 1) ***REMOVED***
		err := fmt.Errorf(
			"%w, metric %q lacks a closing curly brace in its last position",
			ErrMetricNameParsing,
			name,
		)
		return "", nil, err
	***REMOVED***

	// We already know the position of the opening and closing curly brace
	// tokens. Thus, we extract the string in between them, and split its
	// content to obtain the tags key values.
	tags := strings.Split(name[openingTokenPos+1:closingTokenPos], ",")

	// For each tag definition, ensure it is correctly formed
	for i, t := range tags ***REMOVED***
		keyValue := strings.SplitN(t, ":", 2)

		if len(keyValue) != 2 || keyValue[1] == "" ***REMOVED***
			return "", nil, fmt.Errorf("%w, metric %q tag expression is malformed", ErrMetricNameParsing, t)
		***REMOVED***

		tags[i] = strings.TrimSpace(t)
	***REMOVED***

	return name[0:openingTokenPos], tags, nil
***REMOVED***
