package metrics

import (
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
