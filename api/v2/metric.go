package v2

import (
	"bytes"
	"encoding/json"
	"github.com/loadimpact/k6/stats"
)

type NullMetricType struct ***REMOVED***
	Type  stats.MetricType
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
	Type  stats.ValueType
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
	Name     string         `json:"-"`
	Type     NullMetricType `json:"type"`
	Contains NullValueType  `json:"contains"`
***REMOVED***

func (m Metric) GetID() string ***REMOVED***
	return m.Name
***REMOVED***

func (m *Metric) SetID(id string) error ***REMOVED***
	m.Name = id
	return nil
***REMOVED***
