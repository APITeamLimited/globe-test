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
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseExtendedDuration(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		durStr string
		expErr bool
		expDur time.Duration
	***REMOVED******REMOVED***
		***REMOVED***"", true, 0***REMOVED***,
		***REMOVED***"d", true, 0***REMOVED***,
		***REMOVED***"d2h", true, 0***REMOVED***,
		***REMOVED***"d2h", true, 0***REMOVED***,
		***REMOVED***"2.1d", true, 0***REMOVED***,
		***REMOVED***"2d-2h", true, 0***REMOVED***,
		***REMOVED***"-2d-2h", true, 0***REMOVED***,
		***REMOVED***"2+d", true, 0***REMOVED***,
		***REMOVED***"2da", true, 0***REMOVED***,
		***REMOVED***"2-d", true, 0***REMOVED***,
		***REMOVED***"1.12s", false, 1120 * time.Millisecond***REMOVED***,
		***REMOVED***"0d1.12s", false, 1120 * time.Millisecond***REMOVED***,
		***REMOVED***"10d1.12s", false, 240*time.Hour + 1120*time.Millisecond***REMOVED***,
		***REMOVED***"1s", false, 1 * time.Second***REMOVED***,
		***REMOVED***"1d", false, 24 * time.Hour***REMOVED***,
		***REMOVED***"20d", false, 480 * time.Hour***REMOVED***,
		***REMOVED***"1d23h", false, 47 * time.Hour***REMOVED***,
		***REMOVED***"1d24h15m", false, 48*time.Hour + 15*time.Minute***REMOVED***,
		***REMOVED***"1d25h80m", false, 50*time.Hour + 20*time.Minute***REMOVED***,
		***REMOVED***"0d25h120m80s", false, 27*time.Hour + 80*time.Second***REMOVED***,
		***REMOVED***"-1d2h", false, -26 * time.Hour***REMOVED***,
		***REMOVED***"-1d24h", false, -48 * time.Hour***REMOVED***,
		***REMOVED***"2d1ns", false, 48*time.Hour + 1***REMOVED***,
		***REMOVED***"-2562047h47m16.854775807s", false, time.Duration(math.MinInt64 + 1)***REMOVED***,
		***REMOVED***"-106751d23h47m16.854775807s", false, time.Duration(math.MinInt64 + 1)***REMOVED***,
		***REMOVED***"2562047h47m16.854775807s", false, time.Duration(math.MaxInt64)***REMOVED***,
		***REMOVED***"106751d23h47m16.854775807s", false, time.Duration(math.MaxInt64)***REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		tc := tc
		t.Run(fmt.Sprintf("tc_%s_exp", tc.durStr), func(t *testing.T) ***REMOVED***
			result, err := ParseExtendedDuration(tc.durStr)
			if tc.expErr ***REMOVED***
				assert.Error(t, err)
			***REMOVED*** else ***REMOVED***
				assert.NoError(t, err)
				assert.Equal(t, tc.expDur, result)
			***REMOVED***
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
				assert.NoError(t, json.Unmarshal([]byte(`75000`), &d))
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
			t.Run("Extended", func(t *testing.T) ***REMOVED***
				var d Duration
				assert.NoError(t, json.Unmarshal([]byte(`"1d2h1m15s"`), &d))
				assert.Equal(t, Duration(26*time.Hour+75*time.Second), d)
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
				assert.NoError(t, json.Unmarshal([]byte(`75000`), &d))
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

func TestGetDurationValue(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		val      interface***REMOVED******REMOVED***
		expError bool
		exp      time.Duration
	***REMOVED******REMOVED***
		***REMOVED***false, true, 0***REMOVED***,
		***REMOVED***time.Now(), true, 0***REMOVED***,
		***REMOVED***"invalid", true, 0***REMOVED***,
		***REMOVED***uint64(math.MaxInt64) + 1, true, 0***REMOVED***,

		***REMOVED***int8(100), false, 100 * time.Millisecond***REMOVED***,
		***REMOVED***uint8(100), false, 100 * time.Millisecond***REMOVED***,
		***REMOVED***uint(1000), false, time.Second***REMOVED***,
		***REMOVED***int(1000), false, time.Second***REMOVED***,
		***REMOVED***uint16(1000), false, time.Second***REMOVED***,
		***REMOVED***int16(1000), false, time.Second***REMOVED***,
		***REMOVED***uint32(1000), false, time.Second***REMOVED***,
		***REMOVED***int32(1000), false, time.Second***REMOVED***,
		***REMOVED***uint(1000), false, time.Second***REMOVED***,
		***REMOVED***int(1000), false, time.Second***REMOVED***,
		***REMOVED***uint64(1000), false, time.Second***REMOVED***,
		***REMOVED***int64(1000), false, time.Second***REMOVED***,
		***REMOVED***1000, false, time.Second***REMOVED***,
		***REMOVED***1000, false, time.Second***REMOVED***,
		***REMOVED***1000.0, false, time.Second***REMOVED***,
		***REMOVED***float32(1000.0), false, time.Second***REMOVED***,
		***REMOVED***float64(1000.001), false, time.Second + time.Microsecond***REMOVED***,
		***REMOVED***"1000", false, time.Second***REMOVED***,
		***REMOVED***"1000.001", false, time.Second + time.Microsecond***REMOVED***,
		***REMOVED***"1s", false, time.Second***REMOVED***,
		***REMOVED***"1.5s", false, 1500 * time.Millisecond***REMOVED***,
		***REMOVED***time.Second, false, time.Second***REMOVED***,
		***REMOVED***"1d3h1s", false, 27*time.Hour + time.Second***REMOVED***,
		// TODO: test for int overflows when that's implemented
	***REMOVED***

	for i, tc := range testCases ***REMOVED***
		i, tc := i, tc
		t.Run(fmt.Sprintf("testcase_%02d", i), func(t *testing.T) ***REMOVED***
			res, err := GetDurationValue(tc.val)
			if tc.expError ***REMOVED***
				assert.Error(t, err)
			***REMOVED*** else ***REMOVED***
				assert.NoError(t, err)
				assert.Equal(t, tc.exp, res)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
