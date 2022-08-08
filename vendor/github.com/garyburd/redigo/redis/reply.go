// Copyright 2012 Gary Burd
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package redis

import (
	"errors"
	"fmt"
	"strconv"
)

// ErrNil indicates that a reply value is nil.
var ErrNil = errors.New("redigo: nil returned")

// Int is a helper that converts a command reply to an integer. If err is not
// equal to nil, then Int returns 0, err. Otherwise, Int converts the
// reply to an int as follows:
//
//  Reply type    Result
//  integer       int(reply), nil
//  bulk string   parsed reply, nil
//  nil           0, ErrNil
//  other         0, error
func Int(reply interface***REMOVED******REMOVED***, err error) (int, error) ***REMOVED***
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	switch reply := reply.(type) ***REMOVED***
	case int64:
		x := int(reply)
		if int64(x) != reply ***REMOVED***
			return 0, strconv.ErrRange
		***REMOVED***
		return x, nil
	case []byte:
		n, err := strconv.ParseInt(string(reply), 10, 0)
		return int(n), err
	case nil:
		return 0, ErrNil
	case Error:
		return 0, reply
	***REMOVED***
	return 0, fmt.Errorf("redigo: unexpected type for Int, got type %T", reply)
***REMOVED***

// Int64 is a helper that converts a command reply to 64 bit integer. If err is
// not equal to nil, then Int returns 0, err. Otherwise, Int64 converts the
// reply to an int64 as follows:
//
//  Reply type    Result
//  integer       reply, nil
//  bulk string   parsed reply, nil
//  nil           0, ErrNil
//  other         0, error
func Int64(reply interface***REMOVED******REMOVED***, err error) (int64, error) ***REMOVED***
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	switch reply := reply.(type) ***REMOVED***
	case int64:
		return reply, nil
	case []byte:
		n, err := strconv.ParseInt(string(reply), 10, 64)
		return n, err
	case nil:
		return 0, ErrNil
	case Error:
		return 0, reply
	***REMOVED***
	return 0, fmt.Errorf("redigo: unexpected type for Int64, got type %T", reply)
***REMOVED***

var errNegativeInt = errors.New("redigo: unexpected value for Uint64")

// Uint64 is a helper that converts a command reply to 64 bit integer. If err is
// not equal to nil, then Int returns 0, err. Otherwise, Int64 converts the
// reply to an int64 as follows:
//
//  Reply type    Result
//  integer       reply, nil
//  bulk string   parsed reply, nil
//  nil           0, ErrNil
//  other         0, error
func Uint64(reply interface***REMOVED******REMOVED***, err error) (uint64, error) ***REMOVED***
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	switch reply := reply.(type) ***REMOVED***
	case int64:
		if reply < 0 ***REMOVED***
			return 0, errNegativeInt
		***REMOVED***
		return uint64(reply), nil
	case []byte:
		n, err := strconv.ParseUint(string(reply), 10, 64)
		return n, err
	case nil:
		return 0, ErrNil
	case Error:
		return 0, reply
	***REMOVED***
	return 0, fmt.Errorf("redigo: unexpected type for Uint64, got type %T", reply)
***REMOVED***

// Float64 is a helper that converts a command reply to 64 bit float. If err is
// not equal to nil, then Float64 returns 0, err. Otherwise, Float64 converts
// the reply to an int as follows:
//
//  Reply type    Result
//  bulk string   parsed reply, nil
//  nil           0, ErrNil
//  other         0, error
func Float64(reply interface***REMOVED******REMOVED***, err error) (float64, error) ***REMOVED***
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	switch reply := reply.(type) ***REMOVED***
	case []byte:
		n, err := strconv.ParseFloat(string(reply), 64)
		return n, err
	case nil:
		return 0, ErrNil
	case Error:
		return 0, reply
	***REMOVED***
	return 0, fmt.Errorf("redigo: unexpected type for Float64, got type %T", reply)
***REMOVED***

// String is a helper that converts a command reply to a string. If err is not
// equal to nil, then String returns "", err. Otherwise String converts the
// reply to a string as follows:
//
//  Reply type      Result
//  bulk string     string(reply), nil
//  simple string   reply, nil
//  nil             "",  ErrNil
//  other           "",  error
func String(reply interface***REMOVED******REMOVED***, err error) (string, error) ***REMOVED***
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	switch reply := reply.(type) ***REMOVED***
	case []byte:
		return string(reply), nil
	case string:
		return reply, nil
	case nil:
		return "", ErrNil
	case Error:
		return "", reply
	***REMOVED***
	return "", fmt.Errorf("redigo: unexpected type for String, got type %T", reply)
***REMOVED***

// Bytes is a helper that converts a command reply to a slice of bytes. If err
// is not equal to nil, then Bytes returns nil, err. Otherwise Bytes converts
// the reply to a slice of bytes as follows:
//
//  Reply type      Result
//  bulk string     reply, nil
//  simple string   []byte(reply), nil
//  nil             nil, ErrNil
//  other           nil, error
func Bytes(reply interface***REMOVED******REMOVED***, err error) ([]byte, error) ***REMOVED***
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	switch reply := reply.(type) ***REMOVED***
	case []byte:
		return reply, nil
	case string:
		return []byte(reply), nil
	case nil:
		return nil, ErrNil
	case Error:
		return nil, reply
	***REMOVED***
	return nil, fmt.Errorf("redigo: unexpected type for Bytes, got type %T", reply)
***REMOVED***

// Bool is a helper that converts a command reply to a boolean. If err is not
// equal to nil, then Bool returns false, err. Otherwise Bool converts the
// reply to boolean as follows:
//
//  Reply type      Result
//  integer         value != 0, nil
//  bulk string     strconv.ParseBool(reply)
//  nil             false, ErrNil
//  other           false, error
func Bool(reply interface***REMOVED******REMOVED***, err error) (bool, error) ***REMOVED***
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	switch reply := reply.(type) ***REMOVED***
	case int64:
		return reply != 0, nil
	case []byte:
		return strconv.ParseBool(string(reply))
	case nil:
		return false, ErrNil
	case Error:
		return false, reply
	***REMOVED***
	return false, fmt.Errorf("redigo: unexpected type for Bool, got type %T", reply)
***REMOVED***

// MultiBulk is a helper that converts an array command reply to a []interface***REMOVED******REMOVED***.
//
// Deprecated: Use Values instead.
func MultiBulk(reply interface***REMOVED******REMOVED***, err error) ([]interface***REMOVED******REMOVED***, error) ***REMOVED*** return Values(reply, err) ***REMOVED***

// Values is a helper that converts an array command reply to a []interface***REMOVED******REMOVED***.
// If err is not equal to nil, then Values returns nil, err. Otherwise, Values
// converts the reply as follows:
//
//  Reply type      Result
//  array           reply, nil
//  nil             nil, ErrNil
//  other           nil, error
func Values(reply interface***REMOVED******REMOVED***, err error) ([]interface***REMOVED******REMOVED***, error) ***REMOVED***
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	switch reply := reply.(type) ***REMOVED***
	case []interface***REMOVED******REMOVED***:
		return reply, nil
	case nil:
		return nil, ErrNil
	case Error:
		return nil, reply
	***REMOVED***
	return nil, fmt.Errorf("redigo: unexpected type for Values, got type %T", reply)
***REMOVED***

func sliceHelper(reply interface***REMOVED******REMOVED***, err error, name string, makeSlice func(int), assign func(int, interface***REMOVED******REMOVED***) error) error ***REMOVED***
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	switch reply := reply.(type) ***REMOVED***
	case []interface***REMOVED******REMOVED***:
		makeSlice(len(reply))
		for i := range reply ***REMOVED***
			if reply[i] == nil ***REMOVED***
				continue
			***REMOVED***
			if err := assign(i, reply[i]); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		return nil
	case nil:
		return ErrNil
	case Error:
		return reply
	***REMOVED***
	return fmt.Errorf("redigo: unexpected type for %s, got type %T", name, reply)
***REMOVED***

// Float64s is a helper that converts an array command reply to a []float64. If
// err is not equal to nil, then Float64s returns nil, err. Nil array items are
// converted to 0 in the output slice. Floats64 returns an error if an array
// item is not a bulk string or nil.
func Float64s(reply interface***REMOVED******REMOVED***, err error) ([]float64, error) ***REMOVED***
	var result []float64
	err = sliceHelper(reply, err, "Float64s", func(n int) ***REMOVED*** result = make([]float64, n) ***REMOVED***, func(i int, v interface***REMOVED******REMOVED***) error ***REMOVED***
		p, ok := v.([]byte)
		if !ok ***REMOVED***
			return fmt.Errorf("redigo: unexpected element type for Floats64, got type %T", v)
		***REMOVED***
		f, err := strconv.ParseFloat(string(p), 64)
		result[i] = f
		return err
	***REMOVED***)
	return result, err
***REMOVED***

// Strings is a helper that converts an array command reply to a []string. If
// err is not equal to nil, then Strings returns nil, err. Nil array items are
// converted to "" in the output slice. Strings returns an error if an array
// item is not a bulk string or nil.
func Strings(reply interface***REMOVED******REMOVED***, err error) ([]string, error) ***REMOVED***
	var result []string
	err = sliceHelper(reply, err, "Strings", func(n int) ***REMOVED*** result = make([]string, n) ***REMOVED***, func(i int, v interface***REMOVED******REMOVED***) error ***REMOVED***
		switch v := v.(type) ***REMOVED***
		case string:
			result[i] = v
			return nil
		case []byte:
			result[i] = string(v)
			return nil
		default:
			return fmt.Errorf("redigo: unexpected element type for Strings, got type %T", v)
		***REMOVED***
	***REMOVED***)
	return result, err
***REMOVED***

// ByteSlices is a helper that converts an array command reply to a [][]byte.
// If err is not equal to nil, then ByteSlices returns nil, err. Nil array
// items are stay nil. ByteSlices returns an error if an array item is not a
// bulk string or nil.
func ByteSlices(reply interface***REMOVED******REMOVED***, err error) ([][]byte, error) ***REMOVED***
	var result [][]byte
	err = sliceHelper(reply, err, "ByteSlices", func(n int) ***REMOVED*** result = make([][]byte, n) ***REMOVED***, func(i int, v interface***REMOVED******REMOVED***) error ***REMOVED***
		p, ok := v.([]byte)
		if !ok ***REMOVED***
			return fmt.Errorf("redigo: unexpected element type for ByteSlices, got type %T", v)
		***REMOVED***
		result[i] = p
		return nil
	***REMOVED***)
	return result, err
***REMOVED***

// Int64s is a helper that converts an array command reply to a []int64.
// If err is not equal to nil, then Int64s returns nil, err. Nil array
// items are stay nil. Int64s returns an error if an array item is not a
// bulk string or nil.
func Int64s(reply interface***REMOVED******REMOVED***, err error) ([]int64, error) ***REMOVED***
	var result []int64
	err = sliceHelper(reply, err, "Int64s", func(n int) ***REMOVED*** result = make([]int64, n) ***REMOVED***, func(i int, v interface***REMOVED******REMOVED***) error ***REMOVED***
		switch v := v.(type) ***REMOVED***
		case int64:
			result[i] = v
			return nil
		case []byte:
			n, err := strconv.ParseInt(string(v), 10, 64)
			result[i] = n
			return err
		default:
			return fmt.Errorf("redigo: unexpected element type for Int64s, got type %T", v)
		***REMOVED***
	***REMOVED***)
	return result, err
***REMOVED***

// Ints is a helper that converts an array command reply to a []in.
// If err is not equal to nil, then Ints returns nil, err. Nil array
// items are stay nil. Ints returns an error if an array item is not a
// bulk string or nil.
func Ints(reply interface***REMOVED******REMOVED***, err error) ([]int, error) ***REMOVED***
	var result []int
	err = sliceHelper(reply, err, "Ints", func(n int) ***REMOVED*** result = make([]int, n) ***REMOVED***, func(i int, v interface***REMOVED******REMOVED***) error ***REMOVED***
		switch v := v.(type) ***REMOVED***
		case int64:
			n := int(v)
			if int64(n) != v ***REMOVED***
				return strconv.ErrRange
			***REMOVED***
			result[i] = n
			return nil
		case []byte:
			n, err := strconv.Atoi(string(v))
			result[i] = n
			return err
		default:
			return fmt.Errorf("redigo: unexpected element type for Ints, got type %T", v)
		***REMOVED***
	***REMOVED***)
	return result, err
***REMOVED***

// StringMap is a helper that converts an array of strings (alternating key, value)
// into a map[string]string. The HGETALL and CONFIG GET commands return replies in this format.
// Requires an even number of values in result.
func StringMap(result interface***REMOVED******REMOVED***, err error) (map[string]string, error) ***REMOVED***
	values, err := Values(result, err)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if len(values)%2 != 0 ***REMOVED***
		return nil, errors.New("redigo: StringMap expects even number of values result")
	***REMOVED***
	m := make(map[string]string, len(values)/2)
	for i := 0; i < len(values); i += 2 ***REMOVED***
		key, okKey := values[i].([]byte)
		value, okValue := values[i+1].([]byte)
		if !okKey || !okValue ***REMOVED***
			return nil, errors.New("redigo: StringMap key not a bulk string value")
		***REMOVED***
		m[string(key)] = string(value)
	***REMOVED***
	return m, nil
***REMOVED***

// IntMap is a helper that converts an array of strings (alternating key, value)
// into a map[string]int. The HGETALL commands return replies in this format.
// Requires an even number of values in result.
func IntMap(result interface***REMOVED******REMOVED***, err error) (map[string]int, error) ***REMOVED***
	values, err := Values(result, err)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if len(values)%2 != 0 ***REMOVED***
		return nil, errors.New("redigo: IntMap expects even number of values result")
	***REMOVED***
	m := make(map[string]int, len(values)/2)
	for i := 0; i < len(values); i += 2 ***REMOVED***
		key, ok := values[i].([]byte)
		if !ok ***REMOVED***
			return nil, errors.New("redigo: IntMap key not a bulk string value")
		***REMOVED***
		value, err := Int(values[i+1], nil)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		m[string(key)] = value
	***REMOVED***
	return m, nil
***REMOVED***

// Int64Map is a helper that converts an array of strings (alternating key, value)
// into a map[string]int64. The HGETALL commands return replies in this format.
// Requires an even number of values in result.
func Int64Map(result interface***REMOVED******REMOVED***, err error) (map[string]int64, error) ***REMOVED***
	values, err := Values(result, err)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if len(values)%2 != 0 ***REMOVED***
		return nil, errors.New("redigo: Int64Map expects even number of values result")
	***REMOVED***
	m := make(map[string]int64, len(values)/2)
	for i := 0; i < len(values); i += 2 ***REMOVED***
		key, ok := values[i].([]byte)
		if !ok ***REMOVED***
			return nil, errors.New("redigo: Int64Map key not a bulk string value")
		***REMOVED***
		value, err := Int64(values[i+1], nil)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		m[string(key)] = value
	***REMOVED***
	return m, nil
***REMOVED***

// Positions is a helper that converts an array of positions (lat, long)
// into a [][2]float64. The GEOPOS command returns replies in this format.
func Positions(result interface***REMOVED******REMOVED***, err error) ([]*[2]float64, error) ***REMOVED***
	values, err := Values(result, err)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	positions := make([]*[2]float64, len(values))
	for i := range values ***REMOVED***
		if values[i] == nil ***REMOVED***
			continue
		***REMOVED***
		p, ok := values[i].([]interface***REMOVED******REMOVED***)
		if !ok ***REMOVED***
			return nil, fmt.Errorf("redigo: unexpected element type for interface slice, got type %T", values[i])
		***REMOVED***
		if len(p) != 2 ***REMOVED***
			return nil, fmt.Errorf("redigo: unexpected number of values for a member position, got %d", len(p))
		***REMOVED***
		lat, err := Float64(p[0], nil)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		long, err := Float64(p[1], nil)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		positions[i] = &[2]float64***REMOVED***lat, long***REMOVED***
	***REMOVED***
	return positions, nil
***REMOVED***
