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

package lib

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

//TODO: update test
/*
func TestSumStages(t *testing.T) ***REMOVED***
	testdata := map[string]struct ***REMOVED***
		Time   types.NullDuration
		Stages []Stage
	***REMOVED******REMOVED***
		"Blank":    ***REMOVED***types.NullDuration***REMOVED******REMOVED***, []Stage***REMOVED******REMOVED******REMOVED***,
		"Infinite": ***REMOVED***types.NullDuration***REMOVED******REMOVED***, []Stage***REMOVED******REMOVED******REMOVED******REMOVED******REMOVED***,
		"Limit": ***REMOVED***
			types.NullDurationFrom(10 * time.Second),
			[]Stage***REMOVED***
				***REMOVED***Duration: types.NullDurationFrom(5 * time.Second)***REMOVED***,
				***REMOVED***Duration: types.NullDurationFrom(5 * time.Second)***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"InfiniteTail": ***REMOVED***
			types.NullDuration***REMOVED***Duration: types.Duration(10 * time.Second), Valid: false***REMOVED***,
			[]Stage***REMOVED***
				***REMOVED***Duration: types.NullDurationFrom(5 * time.Second)***REMOVED***,
				***REMOVED***Duration: types.NullDurationFrom(5 * time.Second)***REMOVED***,
				***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	for name, data := range testdata ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			assert.Equal(t, data.Time, SumStages(data.Stages))
		***REMOVED***)
	***REMOVED***
***REMOVED***
*/

func TestSplitKV(t *testing.T) ***REMOVED***
	testdata := map[string]struct ***REMOVED***
		k string
		v string
	***REMOVED******REMOVED***
		"key=value":      ***REMOVED***"key", "value"***REMOVED***,
		"key=value=blah": ***REMOVED***"key", "value=blah"***REMOVED***,
		"key=":           ***REMOVED***"key", ""***REMOVED***,
		"key":            ***REMOVED***"key", ""***REMOVED***,
	***REMOVED***

	for s, data := range testdata ***REMOVED***
		t.Run(s, func(t *testing.T) ***REMOVED***
			k, v := SplitKV(s)
			assert.Equal(t, data.k, k)
			assert.Equal(t, data.v, v)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestLerp(t *testing.T) ***REMOVED***
	// data[x][y][t] = v
	data := map[int64]map[int64]map[float64]int64***REMOVED***
		0: ***REMOVED***
			0:   ***REMOVED***0.0: 0, 0.10: 0, 0.5: 0, 1.0: 0***REMOVED***,
			100: ***REMOVED***0.0: 0, 0.10: 10, 0.5: 50, 1.0: 100***REMOVED***,
			500: ***REMOVED***0.0: 0, 0.10: 50, 0.5: 250, 1.0: 500***REMOVED***,
		***REMOVED***,
		100: ***REMOVED***
			200: ***REMOVED***0.0: 100, 0.1: 110, 0.5: 150, 1.0: 200***REMOVED***,
			0:   ***REMOVED***0.0: 100, 0.1: 90, 0.5: 50, 1.0: 0***REMOVED***,
		***REMOVED***,
	***REMOVED***

	for x, data := range data ***REMOVED***
		t.Run("x="+strconv.FormatInt(x, 10), func(t *testing.T) ***REMOVED***
			for y, data := range data ***REMOVED***
				t.Run("y="+strconv.FormatInt(y, 10), func(t *testing.T) ***REMOVED***
					for t_, x1 := range data ***REMOVED***
						t.Run("t="+strconv.FormatFloat(t_, 'f', 2, 64), func(t *testing.T) ***REMOVED***
							assert.Equal(t, x1, Lerp(x, y, t_))
						***REMOVED***)
					***REMOVED***
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestClampf(t *testing.T) ***REMOVED***
	testdata := map[float64]map[struct ***REMOVED***
		Min, Max float64
	***REMOVED***]float64***REMOVED***
		-1.0: ***REMOVED***
			***REMOVED***0.0, 1.0***REMOVED***: 0.0,
			***REMOVED***0.5, 1.0***REMOVED***: 0.5,
			***REMOVED***1.0, 1.0***REMOVED***: 1.0,
			***REMOVED***0.0, 0.5***REMOVED***: 0.0,
		***REMOVED***,
		0.0: ***REMOVED***
			***REMOVED***0.0, 1.0***REMOVED***: 0.0,
			***REMOVED***0.5, 1.0***REMOVED***: 0.5,
			***REMOVED***1.0, 1.0***REMOVED***: 1.0,
			***REMOVED***0.0, 0.5***REMOVED***: 0.0,
		***REMOVED***,
		0.5: ***REMOVED***
			***REMOVED***0.0, 1.0***REMOVED***: 0.5,
			***REMOVED***0.5, 1.0***REMOVED***: 0.5,
			***REMOVED***1.0, 1.0***REMOVED***: 1.0,
			***REMOVED***0.0, 0.5***REMOVED***: 0.5,
		***REMOVED***,
		1.0: ***REMOVED***
			***REMOVED***0.0, 1.0***REMOVED***: 1.0,
			***REMOVED***0.5, 1.0***REMOVED***: 1.0,
			***REMOVED***1.0, 1.0***REMOVED***: 1.0,
			***REMOVED***0.0, 0.5***REMOVED***: 0.5,
		***REMOVED***,
		2.0: ***REMOVED***
			***REMOVED***0.0, 1.0***REMOVED***: 1.0,
			***REMOVED***0.5, 1.0***REMOVED***: 1.0,
			***REMOVED***1.0, 1.0***REMOVED***: 1.0,
			***REMOVED***0.0, 0.5***REMOVED***: 0.5,
		***REMOVED***,
	***REMOVED***

	for val, ranges := range testdata ***REMOVED***
		t.Run(fmt.Sprintf("val=%.1f", val), func(t *testing.T) ***REMOVED***
			for r, result := range ranges ***REMOVED***
				t.Run(fmt.Sprintf("min=%.1f,max=%.1f", r.Min, r.Max), func(t *testing.T) ***REMOVED***
					assert.Equal(t, result, Clampf(val, r.Min, r.Max))
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestMin(t *testing.T) ***REMOVED***
	assert.Equal(t, int64(10), Min(10, 100))
	assert.Equal(t, int64(10), Min(100, 10))
***REMOVED***

func TestMax(t *testing.T) ***REMOVED***
	assert.Equal(t, int64(100), Max(10, 100))
	assert.Equal(t, int64(100), Max(100, 10))
***REMOVED***
