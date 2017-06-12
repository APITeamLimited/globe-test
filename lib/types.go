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
	"bytes"
	"encoding/json"
	"time"
)

type Duration time.Duration

func (d Duration) String() string ***REMOVED***
	return time.Duration(d).String()
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

type NullDuration struct ***REMOVED***
	Duration
	Valid bool
***REMOVED***

func NullDurationFrom(d time.Duration) NullDuration ***REMOVED***
	return NullDuration***REMOVED***Duration(d), true***REMOVED***
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
