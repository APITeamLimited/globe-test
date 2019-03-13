/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2019 Load Impact
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

package scheduler

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckPercentagesSum(t *testing.T) ***REMOVED***
	t.Parallel()
	assert.NoError(t, checkPercentagesSum([]float64***REMOVED***100***REMOVED***))
	assert.NoError(t, checkPercentagesSum([]float64***REMOVED***50, 50***REMOVED***))
	assert.NoError(t, checkPercentagesSum([]float64***REMOVED***100.0 / 3, 100.0 / 3, 100.0 / 3***REMOVED***))
	assert.NoError(t, checkPercentagesSum([]float64***REMOVED***33.33, 33.33, 33.34***REMOVED***))

	assert.Error(t, checkPercentagesSum([]float64***REMOVED******REMOVED***))
	assert.Error(t, checkPercentagesSum([]float64***REMOVED***100 / 3, 100 / 3, 100 / 3***REMOVED***))
	assert.Error(t, checkPercentagesSum([]float64***REMOVED***33.33, 33.33, 33.33***REMOVED***))
	assert.Error(t, checkPercentagesSum([]float64***REMOVED***40, 40, 40***REMOVED***))
***REMOVED***

func TestStrictJSONUnmarshal(t *testing.T) ***REMOVED***
	t.Parallel()
	type someElement struct ***REMOVED***
		Data  int               `json:"data"`
		Props map[string]string `json:"props"`
	***REMOVED***

	testCases := []struct ***REMOVED***
		data           string
		expectedError  bool
		destination    interface***REMOVED******REMOVED***
		expectedResult interface***REMOVED******REMOVED***
	***REMOVED******REMOVED***
		***REMOVED***``, true, &someElement***REMOVED******REMOVED***, nil***REMOVED***,
		***REMOVED***`123`, true, &someElement***REMOVED******REMOVED***, nil***REMOVED***,
		***REMOVED***`"blah"`, true, &someElement***REMOVED******REMOVED***, nil***REMOVED***,
		***REMOVED***`null`, false, &someElement***REMOVED******REMOVED***, &someElement***REMOVED******REMOVED******REMOVED***,
		***REMOVED***
			`***REMOVED***"data": 123, "props": ***REMOVED***"test": "mest"***REMOVED******REMOVED***`, false, &someElement***REMOVED******REMOVED***,
			&someElement***REMOVED***123, map[string]string***REMOVED***"test": "mest"***REMOVED******REMOVED***,
		***REMOVED***,
		***REMOVED***`***REMOVED***"data": 123, "props": ***REMOVED***"test": "mest"***REMOVED******REMOVED***asdg`, true, &someElement***REMOVED******REMOVED***, nil***REMOVED***,
	***REMOVED***
	for i, tc := range testCases ***REMOVED***
		tc := tc
		t.Run(fmt.Sprintf("TestCase#%d", i), func(t *testing.T) ***REMOVED***
			err := strictJSONUnmarshal([]byte(tc.data), &tc.destination)
			if tc.expectedError ***REMOVED***
				require.Error(t, err)
				return
			***REMOVED***
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResult, tc.destination)
		***REMOVED***)
	***REMOVED***
***REMOVED***
