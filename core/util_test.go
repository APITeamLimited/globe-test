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

package core

import (
	"testing"

	"time"

	"github.com/loadimpact/k6/lib"
	"github.com/stretchr/testify/assert"
)

func TestSumStages(t *testing.T) ***REMOVED***
	testdata := map[string]struct ***REMOVED***
		Time   lib.NullDuration
		Stages []lib.Stage
	***REMOVED******REMOVED***
		"Blank":    ***REMOVED***lib.NullDuration***REMOVED******REMOVED***, []lib.Stage***REMOVED******REMOVED******REMOVED***,
		"Infinite": ***REMOVED***lib.NullDuration***REMOVED******REMOVED***, []lib.Stage***REMOVED******REMOVED******REMOVED******REMOVED******REMOVED***,
		"Limit": ***REMOVED***
			lib.NullDurationFrom(10 * time.Second),
			[]lib.Stage***REMOVED***
				***REMOVED***Duration: lib.NullDurationFrom(5 * time.Second)***REMOVED***,
				***REMOVED***Duration: lib.NullDurationFrom(5 * time.Second)***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"InfiniteTail": ***REMOVED***
			lib.NullDuration***REMOVED***Duration: lib.Duration(10 * time.Second), Valid: false***REMOVED***,
			[]lib.Stage***REMOVED***
				***REMOVED***Duration: lib.NullDurationFrom(5 * time.Second)***REMOVED***,
				***REMOVED***Duration: lib.NullDurationFrom(5 * time.Second)***REMOVED***,
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
