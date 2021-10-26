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
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"gopkg.in/guregu/null.v3"
)

// NullDecoder converts data with expected type f to a guregu/null value
// of equivalent type t. It returns an error if a type mismatch occurs.
func NullDecoder(f reflect.Type, t reflect.Type, data interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	typeFrom := f.String()
	typeTo := t.String()

	expectedType := ""
	switch typeTo ***REMOVED***
	case "null.String":
		if typeFrom == reflect.String.String() ***REMOVED***
			return null.StringFrom(data.(string)), nil
		***REMOVED***
		expectedType = reflect.String.String()
	case "null.Bool":
		if typeFrom == reflect.Bool.String() ***REMOVED***
			return null.BoolFrom(data.(bool)), nil
		***REMOVED***
		expectedType = reflect.Bool.String()
	case "null.Int":
		if typeFrom == reflect.Int.String() ***REMOVED***
			return null.IntFrom(int64(data.(int))), nil
		***REMOVED***
		if typeFrom == reflect.Int32.String() ***REMOVED***
			return null.IntFrom(int64(data.(int32))), nil
		***REMOVED***
		if typeFrom == reflect.Int64.String() ***REMOVED***
			return null.IntFrom(data.(int64)), nil
		***REMOVED***
		expectedType = reflect.Int.String()
	case "null.Float":
		if typeFrom == reflect.Float32.String() ***REMOVED***
			return null.FloatFrom(float64(data.(float32))), nil
		***REMOVED***
		if typeFrom == reflect.Float64.String() ***REMOVED***
			return null.FloatFrom(data.(float64)), nil
		***REMOVED***
		expectedType = reflect.Float32.String() + " or " + reflect.Float64.String()
	case "types.NullDuration":
		if typeFrom == reflect.String.String() ***REMOVED***
			var d NullDuration
			err := d.UnmarshalText([]byte(data.(string)))
			return d, err
		***REMOVED***
		expectedType = reflect.String.String()
	***REMOVED***

	if expectedType != "" ***REMOVED***
		return data, fmt.Errorf("expected '%s', got '%s'", expectedType, typeFrom)
	***REMOVED***

	return data, nil
***REMOVED***

// TODO: something better that won't require so much boilerplate and casts for NullDuration values...

// Duration is an alias for time.Duration that de/serialises to JSON as human-readable strings.
type Duration time.Duration

func (d Duration) String() string ***REMOVED***
	return time.Duration(d).String()
***REMOVED***

// ParseExtendedDuration is a helper function that allows for string duration
// values containing days.
func ParseExtendedDuration(data string) (result time.Duration, err error) ***REMOVED***
	// Assume millisecond values if data is provided with no units
	if t, errp := strconv.ParseFloat(data, 64); errp == nil ***REMOVED***
		return time.Duration(t * float64(time.Millisecond)), nil
	***REMOVED***

	dPos := strings.IndexByte(data, 'd')
	if dPos < 0 ***REMOVED***
		return time.ParseDuration(data)
	***REMOVED***

	var hours time.Duration
	if dPos+1 < len(data) ***REMOVED*** // case "12d"
		hours, err = time.ParseDuration(data[dPos+1:])
		if err != nil ***REMOVED***
			return
		***REMOVED***
		if hours < 0 ***REMOVED***
			return 0, fmt.Errorf("invalid time format '%s'", data[dPos+1:])
		***REMOVED***
	***REMOVED***

	days, err := strconv.ParseInt(data[:dPos], 10, 64)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	if days < 0 ***REMOVED***
		hours = -hours
	***REMOVED***
	return time.Duration(days)*24*time.Hour + hours, nil
***REMOVED***

// UnmarshalText converts text data to Duration
func (d *Duration) UnmarshalText(data []byte) error ***REMOVED***
	v, err := ParseExtendedDuration(string(data))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	*d = Duration(v)
	return nil
***REMOVED***

// UnmarshalJSON converts JSON data to Duration
func (d *Duration) UnmarshalJSON(data []byte) error ***REMOVED***
	if len(data) > 0 && data[0] == '"' ***REMOVED***
		var str string
		if err := json.Unmarshal(data, &str); err != nil ***REMOVED***
			return err
		***REMOVED***

		v, err := ParseExtendedDuration(str)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		*d = Duration(v)
	***REMOVED*** else if t, errp := strconv.ParseFloat(string(data), 64); errp == nil ***REMOVED***
		*d = Duration(t * float64(time.Millisecond))
	***REMOVED*** else ***REMOVED***
		return fmt.Errorf("'%s' is not a valid duration value", string(data))
	***REMOVED***

	return nil
***REMOVED***

// MarshalJSON returns the JSON representation of d
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

// NullDurationFrom returns a new valid NullDuration from a time.Duration.
func NullDurationFrom(d time.Duration) NullDuration ***REMOVED***
	return NullDuration***REMOVED***Duration(d), true***REMOVED***
***REMOVED***

// UnmarshalText converts text data to a valid NullDuration
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

// UnmarshalJSON converts JSON data to a valid NullDuration
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

// MarshalJSON returns the JSON representation of d
func (d NullDuration) MarshalJSON() ([]byte, error) ***REMOVED***
	if !d.Valid ***REMOVED***
		return []byte(`null`), nil
	***REMOVED***
	return d.Duration.MarshalJSON()
***REMOVED***

// ValueOrZero returns the underlying Duration value of d if valid or
// its zero equivalent otherwise. It matches the existing guregu/null API.
func (d NullDuration) ValueOrZero() Duration ***REMOVED***
	if !d.Valid ***REMOVED***
		return Duration(0)
	***REMOVED***

	return d.Duration
***REMOVED***

// TimeDuration returns a NullDuration's value as a stdlib Duration.
func (d NullDuration) TimeDuration() time.Duration ***REMOVED***
	return time.Duration(d.Duration)
***REMOVED***

func getInt64(v interface***REMOVED******REMOVED***) (int64, error) ***REMOVED***
	switch n := v.(type) ***REMOVED***
	case int:
		return int64(n), nil
	case int8:
		return int64(n), nil
	case int16:
		return int64(n), nil
	case int32:
		return int64(n), nil
	case int64:
		return n, nil
	case uint:
		return int64(n), nil
	case uint8:
		return int64(n), nil
	case uint16:
		return int64(n), nil
	case uint32:
		return int64(n), nil
	case uint64:
		if n > math.MaxInt64 ***REMOVED***
			return 0, fmt.Errorf("%d is too big", n)
		***REMOVED***
		return int64(n), nil
	default:
		return 0, fmt.Errorf("unable to use type %T as a duration value", v)
	***REMOVED***
***REMOVED***

// GetDurationValue is a helper function that can convert a lot of different
// types to time.Duration.
//
// TODO: move to a separate package and check for integer overflows?
func GetDurationValue(v interface***REMOVED******REMOVED***) (time.Duration, error) ***REMOVED***
	switch d := v.(type) ***REMOVED***
	case time.Duration:
		return d, nil
	case string:
		return ParseExtendedDuration(d)
	case float32:
		return time.Duration(float64(d) * float64(time.Millisecond)), nil
	case float64:
		return time.Duration(d * float64(time.Millisecond)), nil
	default:
		n, err := getInt64(v)
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		return time.Duration(n) * time.Millisecond, nil
	***REMOVED***
***REMOVED***
