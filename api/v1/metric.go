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

package v1

import (
	"bytes"
	"encoding/json"
	"time"

	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/stats"
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
	Name string `json:"-" yaml:"name"`

	Type     NullMetricType `json:"type" yaml:"type"`
	Contains NullValueType  `json:"contains" yaml:"contains"`
	Tainted  null.Bool      `json:"tainted" yaml:"tainted"`

	Sample map[string]float64 `json:"sample" yaml:"sample"`
***REMOVED***

func NewMetric(m *stats.Metric, t time.Duration) Metric ***REMOVED***
	return Metric***REMOVED***
		Name:     m.Name,
		Type:     NullMetricType***REMOVED***m.Type, true***REMOVED***,
		Contains: NullValueType***REMOVED***m.Contains, true***REMOVED***,
		Tainted:  m.Tainted,
		Sample:   m.Sink.Format(t),
	***REMOVED***
***REMOVED***

// GetID gets a metric ID (name)
// Deprecated: use instead m.Name directly
// This method will be removed with the one of the PRs of (https://github.com/grafana/k6/issues/911)
func (m Metric) GetID() string ***REMOVED***
	return m.Name
***REMOVED***

// SetID sets a metric ID (name)
// Deprecated: use instead m.Name directly
// This method will be removed with the one of the PRs of (https://github.com/grafana/k6/issues/911)
func (m *Metric) SetID(id string) error ***REMOVED***
	m.Name = id
	return nil
***REMOVED***
