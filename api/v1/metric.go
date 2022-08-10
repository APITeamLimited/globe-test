package v1

import (
	"bytes"
	"encoding/json"
	"time"

	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/metrics"
)

type NullMetricType struct ***REMOVED***
	Type  metrics.MetricType
	Valid bool
***REMOVED***

func (t NullMetricType) MarshalJSON() ([]byte, error) ***REMOVED***
	if !t.Valid ***REMOVED***
		return []byte("null"), nil
	***REMOVED***
	return t.Type.MarshalJSON()
***REMOVED***

func (t *NullMetricType) UnmarshalJSON(data []byte) error ***REMOVED***
	if bytes.Equal(data, []byte("null")) ***REMOVED***
		t.Valid = false
		return nil
	***REMOVED***
	t.Valid = true
	return json.Unmarshal(data, &t.Type)
***REMOVED***

type NullValueType struct ***REMOVED***
	Type  metrics.ValueType
	Valid bool
***REMOVED***

func (t NullValueType) MarshalJSON() ([]byte, error) ***REMOVED***
	if !t.Valid ***REMOVED***
		return []byte("null"), nil
	***REMOVED***
	return t.Type.MarshalJSON()
***REMOVED***

func (t *NullValueType) UnmarshalJSON(data []byte) error ***REMOVED***
	if bytes.Equal(data, []byte("null")) ***REMOVED***
		t.Valid = false
		return nil
	***REMOVED***
	t.Valid = true
	return json.Unmarshal(data, &t.Type)
***REMOVED***

type Metric struct ***REMOVED***
	Name string `json:"-" yaml:"name"`

	Type     NullMetricType `json:"type" yaml:"type"`
	Contains NullValueType  `json:"contains" yaml:"contains"`
	Tainted  null.Bool      `json:"tainted" yaml:"tainted"`

	Sample map[string]float64 `json:"sample" yaml:"sample"`
***REMOVED***

// NewMetric constructs a new Metric
func NewMetric(m *metrics.Metric, t time.Duration) Metric ***REMOVED***
	return Metric***REMOVED***
		Name:     m.Name,
		Type:     NullMetricType***REMOVED***m.Type, true***REMOVED***,
		Contains: NullValueType***REMOVED***m.Contains, true***REMOVED***,
		Tainted:  m.Tainted,
		Sample:   m.Sink.Format(t),
	***REMOVED***
***REMOVED***
