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
	"bytes"
	"encoding/json"
	"time"
)

// Duration is an alias for time.Duration that de/serialises to JSON as human-readable strings.
type Duration time.Duration

func (d Duration) String() string ***REMOVED***
	return time.Duration(d).String()
***REMOVED***

func (d *Duration) UnmarshalText(data []byte) error ***REMOVED***
	v, err := time.ParseDuration(string(data))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	*d = Duration(v)
	return nil
***REMOVED***

func (d *Duration) UnmarshalJSON(data []byte) error ***REMOVED***
	if len(data) > 0 && data[0] == '"' ***REMOVED***
		var str string
		if err := json.Unmarshal(data, &str); err != nil ***REMOVED***
			return err
		***REMOVED***

		v, err := time.ParseDuration(str)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		*d = Duration(v)
	***REMOVED*** else ***REMOVED***
		var v time.Duration
		if err := json.Unmarshal(data, &v); err != nil ***REMOVED***
			return err
		***REMOVED***
		*d = Duration(v)
	***REMOVED***

	return nil
***REMOVED***

func (d Duration) MarshalJSON() ([]byte, error) ***REMOVED***
	return json.Marshal(d.String())
***REMOVED***

// NullDuration is a nullable Duration, in the same vein as the nullable types provided by
// package gopkg.in/guregu/null.v3.
type NullDuration struct ***REMOVED***
	Duration
	Valid bool
***REMOVED***

// NewNullDuration is a simple helper constructor function
func NewNullDuration(d time.Duration, valid bool) NullDuration ***REMOVED***
	return NullDuration***REMOVED***Duration(d), valid***REMOVED***
***REMOVED***

// Creates a valid NullDuration from a time.Duration.
func NullDurationFrom(d time.Duration) NullDuration ***REMOVED***
	return NullDuration***REMOVED***Duration(d), true***REMOVED***
***REMOVED***

func (d *NullDuration) UnmarshalText(data []byte) error ***REMOVED***
	if len(data) == 0 ***REMOVED***
		*d = NullDuration***REMOVED******REMOVED***
		return nil
	***REMOVED***
	if err := d.Duration.UnmarshalText(data); err != nil ***REMOVED***
		return err
	***REMOVED***
	d.Valid = true
	return nil
***REMOVED***

func (d *NullDuration) UnmarshalJSON(data []byte) error ***REMOVED***
	if bytes.Equal(data, []byte(`null`)) ***REMOVED***
		d.Valid = false
		return nil
	***REMOVED***
	if err := json.Unmarshal(data, &d.Duration); err != nil ***REMOVED***
		return err
	***REMOVED***
	d.Valid = true
	return nil
***REMOVED***

func (d NullDuration) MarshalJSON() ([]byte, error) ***REMOVED***
	if !d.Valid ***REMOVED***
		return []byte(`null`), nil
	***REMOVED***
	return d.Duration.MarshalJSON()
***REMOVED***
