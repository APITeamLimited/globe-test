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

package lib

import (
	"strings"
	"testing"
	"time"

	"go.k6.io/k6/lib/consts"
)

func TestTimeoutError(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		stage, expectedStrContain string
		d                         time.Duration
	***REMOVED******REMOVED***
		***REMOVED***consts.SetupFn, "1 seconds", time.Second***REMOVED***,
		***REMOVED***consts.TeardownFn, "2 seconds", time.Second * 2***REMOVED***,
		***REMOVED***"", "0 seconds", time.Duration(0)***REMOVED***,
	***REMOVED***

	for _, tc := range tests ***REMOVED***
		te := NewTimeoutError(tc.stage, tc.d)
		if !strings.Contains(te.String(), tc.expectedStrContain) ***REMOVED***
			t.Errorf("Expected error contains %s, but got: %s", tc.expectedStrContain, te.String())
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestTimeoutErrorHint(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		stage string
		empty bool
	***REMOVED******REMOVED***
		***REMOVED***consts.SetupFn, false***REMOVED***,
		***REMOVED***consts.TeardownFn, false***REMOVED***,
		***REMOVED***"not handle", true***REMOVED***,
	***REMOVED***

	for _, tc := range tests ***REMOVED***
		te := NewTimeoutError(tc.stage, time.Second)
		if tc.empty && te.Hint() != "" ***REMOVED***
			t.Errorf("Expected empty hint, got: %s", te.Hint())
		***REMOVED***
		if !tc.empty && te.Hint() == "" ***REMOVED***
			t.Errorf("Expected non-empty hint, got empty")
		***REMOVED***
	***REMOVED***
***REMOVED***
