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
	Name string `json:"-"`

	Type     NullMetricType `json:"type"`
	Contains NullValueType  `json:"contains"`
***REMOVED***

func NewMetric(m stats.Metric) Metric ***REMOVED***
	return Metric***REMOVED***
		Name:     m.Name,
		Type:     NullMetricType***REMOVED***m.Type, true***REMOVED***,
		Contains: NullValueType***REMOVED***m.Contains, true***REMOVED***,
	***REMOVED***
***REMOVED***

func (m Metric) GetID() string ***REMOVED***
	return m.Name
***REMOVED***

func (m *Metric) SetID(id string) error ***REMOVED***
	m.Name = id
	return nil
***REMOVED***
