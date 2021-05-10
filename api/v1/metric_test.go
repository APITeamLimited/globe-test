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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/stats"
)

func TestNullMetricTypeJSON(t *testing.T) ***REMOVED***
	values := map[NullMetricType]string***REMOVED***
		***REMOVED******REMOVED***:                    `null`,
		***REMOVED***stats.Counter, true***REMOVED***: `"counter"`,
		***REMOVED***stats.Gauge, true***REMOVED***:   `"gauge"`,
		***REMOVED***stats.Trend, true***REMOVED***:   `"trend"`,
		***REMOVED***stats.Rate, true***REMOVED***:    `"rate"`,
	***REMOVED***
	t.Run("Marshal", func(t *testing.T) ***REMOVED***
		for mt, val := range values ***REMOVED***
			t.Run(val, func(t *testing.T) ***REMOVED***
				data, err := json.Marshal(mt)
				assert.NoError(t, err)
				assert.Equal(t, val, string(data))
			***REMOVED***)
		***REMOVED***
	***REMOVED***)
	t.Run("Unmarshal", func(t *testing.T) ***REMOVED***
		for mt, val := range values ***REMOVED***
			t.Run(val, func(t *testing.T) ***REMOVED***
				var value NullMetricType
				assert.NoError(t, json.Unmarshal([]byte(val), &value))
				assert.Equal(t, mt, value)
			***REMOVED***)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestNullValueTypeJSON(t *testing.T) ***REMOVED***
	values := map[NullValueType]string***REMOVED***
		***REMOVED******REMOVED***:                    `null`,
		***REMOVED***stats.Default, true***REMOVED***: `"default"`,
		***REMOVED***stats.Time, true***REMOVED***:    `"time"`,
	***REMOVED***
	t.Run("Marshal", func(t *testing.T) ***REMOVED***
		for mt, val := range values ***REMOVED***
			t.Run(val, func(t *testing.T) ***REMOVED***
				data, err := json.Marshal(mt)
				assert.NoError(t, err)
				assert.Equal(t, val, string(data))
			***REMOVED***)
		***REMOVED***
	***REMOVED***)
	t.Run("Unmarshal", func(t *testing.T) ***REMOVED***
		for mt, val := range values ***REMOVED***
			t.Run(val, func(t *testing.T) ***REMOVED***
				var value NullValueType
				assert.NoError(t, json.Unmarshal([]byte(val), &value))
				assert.Equal(t, mt, value)
			***REMOVED***)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestNewMetric(t *testing.T) ***REMOVED***
	old := stats.New("name", stats.Trend, stats.Time)
	old.Tainted = null.BoolFrom(true)
	m := NewMetric(old, 0)
	assert.Equal(t, "name", m.Name)
	assert.True(t, m.Type.Valid)
	assert.Equal(t, stats.Trend, m.Type.Type)
	assert.True(t, m.Contains.Valid)
	assert.True(t, m.Tainted.Bool)
	assert.True(t, m.Tainted.Valid)
	assert.Equal(t, stats.Time, m.Contains.Type)
	assert.NotEmpty(t, m.Sample)
***REMOVED***
