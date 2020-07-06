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
	"fmt"
	"time"

	"github.com/loadimpact/k6/lib/consts"
)

// TimeoutError is used when somethings timeouts
type TimeoutError struct ***REMOVED***
	place string
	d     time.Duration
***REMOVED***

// NewTimeoutError returns a new TimeoutError reporting that timeout has happened
// at the given place and given duration.
func NewTimeoutError(place string, d time.Duration) TimeoutError ***REMOVED***
	return TimeoutError***REMOVED***place: place, d: d***REMOVED***
***REMOVED***

// String returns timeout error in human readable format.
func (t TimeoutError) String() string ***REMOVED***
	return fmt.Sprintf("%s() execution timed out after %.f seconds", t.place, t.d.Seconds())
***REMOVED***

// Error implements error interface.
func (t TimeoutError) Error() string ***REMOVED***
	return t.String()
***REMOVED***

// Place returns the place where timeout occurred.
func (t TimeoutError) Place() string ***REMOVED***
	return t.place
***REMOVED***

// Hint returns a hint message for logging with given stage.
func (t TimeoutError) Hint() string ***REMOVED***
	hint := ""

	switch t.place ***REMOVED***
	case consts.SetupFn:
		hint = "You can increase the time limit via the setupTimeout option"
	case consts.TeardownFn:
		hint = "You can increase the time limit via the teardownTimeout option"
	***REMOVED***
	return hint
***REMOVED***
