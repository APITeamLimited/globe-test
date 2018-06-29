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

package types

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	null "gopkg.in/guregu/null.v3"
)

func TestNullDecoder(t *testing.T) ***REMOVED***
	type foo struct ***REMOVED***
		Strs      []string
		Str       null.String
		Boolean   null.Bool
		Integer   null.Int
		Integer32 null.Int
		Integer64 null.Int
		Float32   null.Float
		Float64   null.Float
		Dur       NullDuration
	***REMOVED***
	f := foo***REMOVED******REMOVED***
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig***REMOVED***
		DecodeHook: NullDecoder,
		Result:     &f,
	***REMOVED***)
	require.NoError(t, err)

	conf := map[string]interface***REMOVED******REMOVED******REMOVED***
		"strs":      []string***REMOVED***"fake"***REMOVED***,
		"str":       "bar",
		"boolean":   true,
		"integer":   42,
		"integer32": int32(42),
		"integer64": int64(42),
		"float32":   float32(3.14),
		"float64":   float64(3.14),
		"dur":       "1m",
	***REMOVED***

	err = dec.Decode(conf)
	require.NoError(t, err)

	require.Equal(t, foo***REMOVED***
		Strs:      []string***REMOVED***"fake"***REMOVED***,
		Str:       null.StringFrom("bar"),
		Boolean:   null.BoolFrom(true),
		Integer:   null.IntFrom(42),
		Integer32: null.IntFrom(42),
		Integer64: null.IntFrom(42),
		Float32:   null.FloatFrom(3.140000104904175),
		Float64:   null.FloatFrom(3.14),
		Dur:       NewNullDuration(1*time.Minute, true),
	***REMOVED***, f)

	input := map[string][]interface***REMOVED******REMOVED******REMOVED***
		"Str":       ***REMOVED***true, "string", "bool"***REMOVED***,
		"Boolean":   ***REMOVED***"invalid", "bool", "string"***REMOVED***,
		"Integer":   ***REMOVED***"invalid", "int", "string"***REMOVED***,
		"Integer32": ***REMOVED***true, "int", "bool"***REMOVED***,
		"Integer64": ***REMOVED***"invalid", "int", "string"***REMOVED***,
		"Float32":   ***REMOVED***true, "float32 or float64", "bool"***REMOVED***,
		"Float64":   ***REMOVED***"invalid", "float32 or float64", "string"***REMOVED***,
		"Dur":       ***REMOVED***10, "string", "int"***REMOVED***,
	***REMOVED***

	for k, v := range input ***REMOVED***
		t.Run("Error Message/"+k, func(t *testing.T) ***REMOVED***
			err = dec.Decode(map[string]interface***REMOVED******REMOVED******REMOVED***
				k: v[0],
			***REMOVED***)
			assert.EqualError(t, err, fmt.Sprintf("1 error(s) decoding:\n\n* error decoding '%s': expected '%s', got '%s'", k, v[1], v[2]))
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestDuration(t *testing.T) ***REMOVED***
	t.Run("String", func(t *testing.T) ***REMOVED***
		assert.Equal(t, "1m15s", Duration(75*time.Second).String())
	***REMOVED***)
	t.Run("JSON", func(t *testing.T) ***REMOVED***
		t.Run("Unmarshal", func(t *testing.T) ***REMOVED***
			t.Run("Number", func(t *testing.T) ***REMOVED***
				var d Duration
				assert.NoError(t, json.Unmarshal([]byte(`75000000000`), &d))
				assert.Equal(t, Duration(75*time.Second), d)
			***REMOVED***)
			t.Run("Seconds", func(t *testing.T) ***REMOVED***
				var d Duration
				assert.NoError(t, json.Unmarshal([]byte(`"75s"`), &d))
				assert.Equal(t, Duration(75*time.Second), d)
			***REMOVED***)
			t.Run("String", func(t *testing.T) ***REMOVED***
				var d Duration
				assert.NoError(t, json.Unmarshal([]byte(`"1m15s"`), &d))
				assert.Equal(t, Duration(75*time.Second), d)
			***REMOVED***)
		***REMOVED***)
		t.Run("Marshal", func(t *testing.T) ***REMOVED***
			d := Duration(75 * time.Second)
			data, err := json.Marshal(d)
			assert.NoError(t, err)
			assert.Equal(t, `"1m15s"`, string(data))
		***REMOVED***)
	***REMOVED***)
	t.Run("Text", func(t *testing.T) ***REMOVED***
		var d Duration
		assert.NoError(t, d.UnmarshalText([]byte(`10s`)))
		assert.Equal(t, Duration(10*time.Second), d)
	***REMOVED***)
***REMOVED***

func TestNullDuration(t *testing.T) ***REMOVED***
	t.Run("String", func(t *testing.T) ***REMOVED***
		assert.Equal(t, "1m15s", Duration(75*time.Second).String())
	***REMOVED***)
	t.Run("JSON", func(t *testing.T) ***REMOVED***
		t.Run("Unmarshal", func(t *testing.T) ***REMOVED***
			t.Run("Number", func(t *testing.T) ***REMOVED***
				var d NullDuration
				assert.NoError(t, json.Unmarshal([]byte(`75000000000`), &d))
				assert.Equal(t, NullDuration***REMOVED***Duration(75 * time.Second), true***REMOVED***, d)
			***REMOVED***)
			t.Run("Seconds", func(t *testing.T) ***REMOVED***
				var d NullDuration
				assert.NoError(t, json.Unmarshal([]byte(`"75s"`), &d))
				assert.Equal(t, NullDuration***REMOVED***Duration(75 * time.Second), true***REMOVED***, d)
			***REMOVED***)
			t.Run("String", func(t *testing.T) ***REMOVED***
				var d NullDuration
				assert.NoError(t, json.Unmarshal([]byte(`"1m15s"`), &d))
				assert.Equal(t, NullDuration***REMOVED***Duration(75 * time.Second), true***REMOVED***, d)
			***REMOVED***)
			t.Run("Null", func(t *testing.T) ***REMOVED***
				var d NullDuration
				assert.NoError(t, json.Unmarshal([]byte(`null`), &d))
				assert.Equal(t, NullDuration***REMOVED***Duration(0), false***REMOVED***, d)
			***REMOVED***)
		***REMOVED***)
		t.Run("Marshal", func(t *testing.T) ***REMOVED***
			t.Run("Valid", func(t *testing.T) ***REMOVED***
				d := NullDuration***REMOVED***Duration(75 * time.Second), true***REMOVED***
				data, err := json.Marshal(d)
				assert.NoError(t, err)
				assert.Equal(t, `"1m15s"`, string(data))
			***REMOVED***)
			t.Run("null", func(t *testing.T) ***REMOVED***
				var d NullDuration
				data, err := json.Marshal(d)
				assert.NoError(t, err)
				assert.Equal(t, `null`, string(data))
			***REMOVED***)
		***REMOVED***)
	***REMOVED***)
	t.Run("Text", func(t *testing.T) ***REMOVED***
		var d NullDuration
		assert.NoError(t, d.UnmarshalText([]byte(`10s`)))
		assert.Equal(t, NullDurationFrom(10*time.Second), d)

		t.Run("Empty", func(t *testing.T) ***REMOVED***
			var d NullDuration
			assert.NoError(t, d.UnmarshalText([]byte(``)))
			assert.Equal(t, NullDuration***REMOVED******REMOVED***, d)
		***REMOVED***)
	***REMOVED***)
***REMOVED***

func TestNullDurationFrom(t *testing.T) ***REMOVED***
	assert.Equal(t, NullDuration***REMOVED***Duration(10 * time.Second), true***REMOVED***, NullDurationFrom(10*time.Second))
***REMOVED***
