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
	"strings"
)

// Returns the total sum of time taken by the given set of stages.
func SumStages(stages []Stage) (d NullDuration) ***REMOVED***
	for _, stage := range stages ***REMOVED***
		d.Valid = stage.Duration.Valid
		if stage.Duration.Valid ***REMOVED***
			d.Duration += stage.Duration.Duration
		***REMOVED***
	***REMOVED***
	return d
***REMOVED***

// Splits a string in the form "key=value".
func SplitKV(s string) (key, value string) ***REMOVED***
	parts := strings.SplitN(s, "=", 2)
	if len(parts) == 1 ***REMOVED***
		return parts[0], ""
	***REMOVED***
	return parts[0], parts[1]
***REMOVED***

// Lerp is a linear interpolation between two values x and y, returning the value at the point t,
// where t is a fraction in the range [0.0 - 1.0].
func Lerp(x, y int64, t float64) int64 ***REMOVED***
	return x + int64(t*float64(y-x))
***REMOVED***

// Clampf returns the given value, "clamped" to the range [min, max].
func Clampf(val, min, max float64) float64 ***REMOVED***
	switch ***REMOVED***
	case val < min:
		return min
	case val > max:
		return max
	default:
		return val
	***REMOVED***
***REMOVED***

// Returns the maximum value of a and b.
func Max(a, b int64) int64 ***REMOVED***
	if a > b ***REMOVED***
		return a
	***REMOVED***
	return b
***REMOVED***

// Returns the minimum value of a and b.
func Min(a, b int64) int64 ***REMOVED***
	if a < b ***REMOVED***
		return a
	***REMOVED***
	return b
***REMOVED***
